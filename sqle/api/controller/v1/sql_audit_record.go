package v1

import (
	"archive/zip"
	"context"
	"encoding/json"
	e "errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	xmlParser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

type CreateSQLAuditRecordReqV1 struct {
	DbType         string `json:"db_type" form:"db_type" example:"MySQL"`
	InstanceName   string `json:"instance_name" form:"instance_name" example:"inst_1"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema" example:"db1"`
	Sqls           string `json:"sqls" form:"sqls" example:"alter table tb1 drop columns c1; select * from tb"`
}

type CreateSQLAuditRecordResV1 struct {
	controller.BaseRes
	Data *SQLAuditRecordResData `json:"data"`
}

type SQLAuditRecordResData struct {
	Id   string          `json:"sql_audit_record_id"`
	Task *AuditTaskResV1 `json:"task"`
}

// 10M
var maxZipFileSize int64 = 1024 * 1024 * 10

// CreateSQLAuditRecord
// @Summary SQL审核
// @Id CreateSQLAuditRecordV1
// @Description SQL audit
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Description 4. file[input_zip_file]: it is ZIP file that sql will be parsed from xml or sql file inside it.
// @Accept mpfd
// @Produce json
// @Tags sql_audit_record
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name formData string false "instance name"
// @Param instance_schema formData string false "schema of instance"
// @Param db_type formData string false "db type of instance"
// @Param sqls formData string false "sqls for audit"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Success 200 {object} v1.CreateSQLAuditRecordResV1
// @router /v1/projects/{project_name}/sql_audit_records [post]
func CreateSQLAuditRecord(c echo.Context) error {
	req := new(CreateSQLAuditRecordReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if req.DbType == "" && req.InstanceName == "" {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, e.New("db_type and instance_name can't both be empty")))
	}
	projectName := c.Param("project_name")

	s := model.GetStorage()
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}
	if project.IsArchived() {
		return controller.JSONBaseErrorReq(c, ErrProjectArchived)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err := CheckIsProjectMember(user.Name, project.Name); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var sqls string
	var source string
	if req.Sqls != "" {
		sqls, source = req.Sqls, model.TaskSQLSourceFromFormData
	} else {
		sqls, source, err = getSQLFromFile(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	var task *model.Task
	if req.InstanceName != "" {
		task, err = buildOnlineTaskForAudit(c, s, user.ID, req.InstanceName, req.InstanceSchema, projectName, source, sqls)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	} else {
		task, err = buildOfflineTaskForAudit(user.ID, req.DbType, source, sqls)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	recordId, err := utils.GenUid()
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("generate audit record id failed: %v", err))
	}
	record := model.SQLAuditRecord{
		ProjectId:     project.ID,
		CreatorId:     user.ID,
		AuditRecordId: recordId,
		TaskId:        task.ID,
		Task:          task,
	}
	if err := s.Save(&record); err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("save sql audit record failed: %v", err))
	}

	task, err = server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%d", task.ID), server.ActionTypeAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &CreateSQLAuditRecordResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &SQLAuditRecordResData{
			Id: record.AuditRecordId,
			Task: &AuditTaskResV1{
				Id:             task.ID,
				InstanceName:   task.InstanceName(),
				InstanceDbType: task.DBType,
				InstanceSchema: req.InstanceSchema,
				AuditLevel:     task.AuditLevel,
				Score:          task.Score,
				PassRate:       task.PassRate,
				Status:         task.Status,
				SQLSource:      task.SQLSource,
				ExecStartTime:  task.ExecStartAt,
				ExecEndTime:    task.ExecEndAt,
			},
		},
	})
}

func buildOnlineTaskForAudit(c echo.Context, s *model.Storage, userId uint, instanceName, instanceSchema, projectName, sourceType, sqls string) (*model.Task, error) {
	instance, exist, err := s.GetInstanceByNameAndProjectName(instanceName, projectName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrInstanceNoAccess
	}
	can, err := checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return nil, err
	}
	if !can {
		return nil, ErrInstanceNoAccess
	}

	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	if err := plugin.Ping(context.TODO()); err != nil {
		return nil, err
	}

	task := &model.Task{
		Schema:       instanceSchema,
		InstanceId:   instance.ID,
		Instance:     instance,
		CreateUserId: userId,
		ExecuteSQLs:  []*model.ExecuteSQL{},
		SQLSource:    sourceType,
		DBType:       instance.DbType,
	}
	createAt := time.Now()
	task.CreatedAt = createAt

	nodes, err := plugin.Parse(context.TODO(), sqls)
	if err != nil {
		return nil, err
	}
	for n, node := range nodes {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(n + 1),
				Content: node.Text,
			},
		})
	}
	return task, nil
}

func buildOfflineTaskForAudit(userId uint, dbType, sourceType, sqls string) (*model.Task, error) {
	task := &model.Task{
		CreateUserId: userId,
		ExecuteSQLs:  []*model.ExecuteSQL{},
		SQLSource:    sourceType,
		DBType:       dbType,
	}
	var err error
	var nodes []driverV2.Node
	plugin, err := common.NewDriverManagerWithoutCfg(log.NewEntry(), dbType)
	if err != nil {
		return nil, fmt.Errorf("open plugin failed: %v", err)
	}
	defer plugin.Close(context.TODO())

	nodes, err = plugin.Parse(context.TODO(), sqls)
	if err != nil {
		return nil, fmt.Errorf("parse sqls failed: %v", err)
	}

	createAt := time.Now()
	task.CreatedAt = createAt

	for n, node := range nodes {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(n + 1),
				Content: node.Text,
			},
		})
	}
	return task, nil
}

func getSqlsFromZip(c echo.Context) (sqls string, exist bool, err error) {
	file, err := c.FormFile(InputZipFileName)
	if err == http.ErrMissingFile {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	f, err := file.Open()
	if err != nil {
		return "", false, err
	}
	defer f.Close()

	currentPos, err := f.Seek(0, io.SeekEnd) // get size of zip file
	if err != nil {
		return "", false, err
	}
	size := currentPos + 1
	if size > maxZipFileSize {
		return "", false, fmt.Errorf("file can't be bigger than %vM", maxZipFileSize/1024/1024)
	}
	r, err := zip.NewReader(f, size)
	if err != nil {
		return "", false, err
	}
	var sqlBuffer strings.Builder
	xmlContents := make([]string, len(r.File))
	for i := range r.File {
		srcFile := r.File[i]
		if srcFile == nil {
			continue
		}
		if !strings.HasSuffix(srcFile.Name, ".xml") && !strings.HasSuffix(srcFile.Name, ".sql") {
			continue
		}

		r, err := srcFile.Open()
		if err != nil {
			return "", false, fmt.Errorf("open src file failed:  %v", err)
		}
		content, err := io.ReadAll(r)
		if err != nil {
			return "", false, fmt.Errorf("read src file failed:  %v", err)
		}

		if strings.HasSuffix(srcFile.Name, ".xml") {
			xmlContents[i] = string(content)
		} else if strings.HasSuffix(srcFile.Name, ".sql") {
			if _, err = sqlBuffer.Write(content); err != nil {
				return "", false, fmt.Errorf("gather sqls from sql file failed: %v", err)
			}
		}
	}

	// parse xml content
	ss, err := xmlParser.ParseXMLs(xmlContents, false)
	if err != nil {
		return "", false, fmt.Errorf("parse sqls from xml failed: %v", err)
	}
	for i := range ss {
		if !strings.HasSuffix(sqlBuffer.String(), ";") {
			if _, err = sqlBuffer.WriteString(";"); err != nil {
				return "", false, fmt.Errorf("gather sqls from xml file failed: %v", err)
			}
		}
		if _, err = sqlBuffer.WriteString(ss[i]); err != nil {
			return "", false, fmt.Errorf("gather sqls from xml file failed: %v", err)
		}
	}

	return sqlBuffer.String(), true, nil
}

type UpdateSQLAuditRecordReqV1 struct {
	Tags *[]string `json:"tags" valid:"dive,name"`
}

// UpdateSQLAuditRecordV1
// @Summary 更新SQL审核记录
// @Description update SQL audit record
// @Accept json
// @Id updateSQLAuditRecordV1
// @Tags sql_audit_record
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param param body v1.UpdateSQLAuditRecordReqV1 true "update SQL audit record"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_audit_records/{sql_audit_record_id}/ [patch]
func UpdateSQLAuditRecordV1(c echo.Context) error {
	req := new(UpdateSQLAuditRecordReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectName := c.Param("project_name")
	auditRecordId := c.Param("sql_audit_record_id")

	s := model.GetStorage()
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.Tags != nil {
		if yes, err := s.IsSQLAuditRecordBelongToCurrentUser(user.ID, project.ID, auditRecordId); err != nil {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("check privilege failed: %v", err))
		} else if !yes {
			return controller.JSONBaseErrorReq(c, errors.New(errors.ErrAccessDeniedError, errors.NewAccessDeniedErr("you can't update SQL audit record that created by others")))
		}

		data := model.SQLAuditRecordUpdateData{Tags: *req.Tags}
		if err = s.UpdateSQLAuditRecordById(auditRecordId, data); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetSQLAuditRecordsReqV1 struct {
	FuzzySearchTags      string `json:"fuzzy_search_tags" query:"fuzzy_search_tags"` // todo issue1811
	FilterSQLAuditStatus string `json:"filter_sql_audit_status" query:"filter_sql_audit_status" enums:"auditing,successfully,"`
	FilterInstanceName   string `json:"filter_instance_name" query:"filter_instance_name"`
	FilterCreateTimeFrom string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo   string `json:"filter_create_time_to" query:"filter_create_time_to"`
	PageIndex            uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize             uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type SQLAuditRecord struct {
	Creator          string                 `json:"creator"`
	SQLAuditRecordId string                 `json:"sql_audit_record_id"`
	SQLAuditStatus   string                 `json:"sql_audit_status"`
	Tags             []string               `json:"tags"`
	Instance         SQLAuditRecordInstance `json:"instance"`
	Task             AuditTaskResV1         `json:"task"`
	CreatedAt        *time.Time             `json:"created_at"`
}

type SQLAuditRecordInstance struct {
	Host string `json:"db_host" example:"10.10.10.10"`
	Port string `json:"db_port" example:"3306"`
}

type GetSQLAuditRecordsResV1 struct {
	controller.BaseRes
	Data      []SQLAuditRecord `json:"data"`
	TotalNums uint64           `json:"total_nums"`
}

const (
	SQLAuditRecordStatusAuditing     = "auditing"
	SQLAuditRecordStatusSuccessfully = "successfully"
)

// GetSQLAuditRecordsV1
// @Summary 获取SQL审核记录列表
// @Description get sql audit records
// @Tags sql_audit_record
// @Id getSQLAuditRecordsV1
// @Security ApiKeyAuth
// @Param fuzzy_search_tags query string false "fuzzy search tags"
// @Param filter_sql_audit_status query string false "filter sql audit status" Enums(auditing,successfully)
// @Param filter_instance_name query string false "filter instance name"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetSQLAuditRecordsResV1
// @router /v1/projects/{project_name}/sql_audit_records [get]
func GetSQLAuditRecordsV1(c echo.Context) error {
	req := new(GetSQLAuditRecordsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectName := c.Param("project_name")

	s := model.GetStorage()
	_, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	isManager, err := s.IsProjectManager(user.Name, projectName)
	if err != nil {
		return err
	}

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}

	data := map[string]interface{}{
		"filter_project_name":     projectName,
		"filter_creator_id":       user.ID,
		"fuzzy_search_tags":       req.FuzzySearchTags,
		"filter_instance_name":    req.FilterInstanceName,
		"filter_create_time_from": req.FilterCreateTimeFrom,
		"filter_create_time_to":   req.FilterCreateTimeTo,
		"check_user_can_access":   !isManager,
		"limit":                   req.PageSize,
		"offset":                  offset,
	}
	if req.FilterSQLAuditStatus == SQLAuditRecordStatusAuditing {
		data["filter_task_status_exclude"] = model.TaskStatusAudited
	} else if req.FilterSQLAuditStatus == SQLAuditRecordStatusSuccessfully {
		data["filter_task_status"] = model.TaskStatusAudited
	}

	records, err := s.GetSQLAuditRecordsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	resData := make([]SQLAuditRecord, len(records))
	for i := range records {
		record := records[i]
		status := SQLAuditRecordStatusAuditing
		if record.TaskStatus == model.TaskStatusAudited {
			status = SQLAuditRecordStatusSuccessfully
		}
		var tags []string
		if err := json.Unmarshal([]byte(record.Tags.String), &tags); err != nil {
			log.NewEntry().Errorf("parse tags failed,tags:%v , err: %v", record.Tags, err)
		}
		resData[i] = SQLAuditRecord{
			Creator:          record.CreatorName,
			SQLAuditRecordId: record.AuditRecordId,
			SQLAuditStatus:   status,
			Tags:             tags,
			CreatedAt:        record.RecordCreatedAt,
			Instance: SQLAuditRecordInstance{
				Host: record.InstanceHost.String,
				Port: record.InstancePort.String,
			},
			Task: AuditTaskResV1{
				Id:             record.TaskId,
				InstanceName:   record.InstanceName.String,
				InstanceDbType: record.DbType,
				InstanceSchema: record.InstanceSchema,
				AuditLevel:     record.AuditLevel.String,
				Score:          record.AuditScore.Int32,
				PassRate:       record.AuditPassRate.Float64,
				Status:         record.TaskStatus,
				SQLSource:      record.SQLSource,
			},
		}
	}
	return c.JSON(http.StatusOK, &GetSQLAuditRecordsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      resData,
		TotalNums: uint64(len(records)),
	})
}

type GetSQLAuditRecordResV1 struct {
	controller.BaseRes
	Data SQLAuditRecord `json:"data"`
}

// GetSQLAuditRecordV1
// @Summary 获取SQL审核记录信息
// @Description get sql audit record info
// @Tags sql_audit_record
// @Id getSQLAuditRecordV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_audit_record_id path string true "sql audit record id"
// @Success 200 {object} v1.GetSQLAuditRecordResV1
// @router /v1/projects/{project_name}/sql_audit_records/{sql_audit_record_id}/ [get]
func GetSQLAuditRecordV1(c echo.Context) error {
	projectName := c.Param("project_name")
	auditRecordId := c.Param("sql_audit_record_id")

	s := model.GetStorage()
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if yes, err := s.IsSQLAuditRecordBelongToCurrentUser(user.ID, project.ID, auditRecordId); err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("check privilege failed: %v", err))
	} else if !yes {
		return controller.JSONBaseErrorReq(c, errors.New(errors.ErrAccessDeniedError, errors.NewAccessDeniedErr("you can't see the SQL audit record because it isn't created by you")))
	}

	record, exist, err := s.GetSQLAuditRecordById(project.ID, auditRecordId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, e.New("can not find record")))
	}

	return c.JSON(http.StatusOK, &CreateSQLAuditRecordResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &SQLAuditRecordResData{
			Id: auditRecordId,
			Task: &AuditTaskResV1{
				Id:             record.Task.ID,
				InstanceName:   record.Task.InstanceName(),
				InstanceDbType: record.Task.DBType,
				InstanceSchema: record.Task.Schema,
				AuditLevel:     record.Task.AuditLevel,
				Score:          record.Task.Score,
				PassRate:       record.Task.PassRate,
				Status:         record.Task.Status,
				SQLSource:      record.Task.SQLSource,
				ExecStartTime:  record.Task.ExecStartAt,
				ExecEndTime:    record.Task.ExecEndAt,
			},
		},
	})
}

type GetSQLAuditRecordTagTipsResV1 struct {
	controller.BaseRes
	Tags []string `json:"data"`
}

// GetSQLAuditRecordTagTipsV1
// @Summary 获取SQL审核记录标签列表
// @Description get sql audit record tag tips
// @Tags sql_audit_record
// @Id GetSQLAuditRecordTagTipsV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetSQLAuditRecordTagTipsResV1
// @router /v1/projects/{project_name}/sql_audit_records/tag_tips [get]
func GetSQLAuditRecordTagTipsV1(c echo.Context) error {
	return c.JSON(http.StatusOK, &GetSQLAuditRecordTagTipsResV1{
		BaseRes: controller.BaseRes{},
		Tags:    []string{"全量", "增量"},
	})
}
