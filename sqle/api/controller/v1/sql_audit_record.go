package v1

import (
	"archive/zip"
	"context"
	"encoding/json"
	e "errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	javaParser "github.com/actiontech/java-sql-extractor/parser"
	xmlParser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"
	goGit "github.com/go-git/go-git/v5"
	goGitTransport "github.com/go-git/go-git/v5/plumbing/transport/http"

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
// @Description 5. formData[git_http_url]:the url which scheme is http(s) and end with .git.
// @Description 6. formData[git_user_name]:The name of the user who owns the repository read access.
// @Description 7. formData[git_user_password]:The password corresponding to git_user_name.
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
// @Param git_http_url formData string false "git repository url"
// @Param git_user_name formData string false "the name of user to clone the repository"
// @Param git_user_password formData string false "the password corresponding to git_user_name"
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
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	sqls := getSQLFromFileResp{}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if req.Sqls != "" {
		sqls = getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromFormData,
			SQLsFromFormData: req.Sqls,
		}
	} else {
		sqls, err = getSQLFromFile(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	var task *model.Task
	if req.InstanceName != "" {
		task, err = buildOnlineTaskForAudit(c, s, uint64(user.ID), req.InstanceName, req.InstanceSchema, projectUid, sqls)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	} else {
		task, err = buildOfflineTaskForAudit(uint64(user.ID), req.DbType, sqls)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	// if task instance is not nil, gorm will update instance when save task.
	task.Instance = nil

	recordId, err := utils.GenUid()
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("generate audit record id failed: %v", err))
	}
	record := model.SQLAuditRecord{
		ProjectId:     projectUid,
		CreatorId:     user.GetIDStr(),
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

type getSQLFromFileResp struct {
	SourceType       string
	SQLsFromFormData string
	SQLsFromSQLFiles []SQLsFromSQLFile
	SQLsFromXMLs     []SQLFromXML
}

type SQLsFromSQLFile struct {
	FilePath string
	SQLs     string
}

type SQLFromXML struct {
	FilePath  string
	StartLine uint64
	SQL       string
}

func addSQLsFromFileToTasks(sqls getSQLFromFileResp, task *model.Task, plugin driver.Plugin) error {
	var num uint = 1

	fileTask := func(sqlsText, filePath string, defaultStartLine uint64) error {
		nodes, err := plugin.Parse(context.TODO(), sqlsText)
		if err != nil {
			return fmt.Errorf("parse sqls failed: %v", err)
		}
		for _, node := range nodes {
			startLine := defaultStartLine
			if startLine == 0 {
				startLine = node.StartLine
			}
			task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
				BaseSQL: model.BaseSQL{
					Number:     num,
					Content:    node.Text,
					SourceFile: filePath,
					StartLine:  startLine,
					SQLType:    node.Type,
				},
			})
			num++
		}
		return nil
	}

	// 处理从页面输入的SQL
	err := fileTask(sqls.SQLsFromFormData, "", 0)
	if err != nil {
		return fmt.Errorf("parse sqls failed: %v", err)
	}

	for _, sqlsFromOneFile := range sqls.SQLsFromSQLFiles {
		// SQL文件和Java文件里的SQL都在这里处理
		err := fileTask(sqlsFromOneFile.SQLs, sqlsFromOneFile.FilePath, 0)
		if err != nil {
			return fmt.Errorf("parse sqls failed: %v", err)
		}
	}

	for _, sqlsFromOneFile := range sqls.SQLsFromXMLs {
		// 这里遍历的一个元素里只有一条SQL
		err := fileTask(sqlsFromOneFile.SQL, sqlsFromOneFile.FilePath, sqlsFromOneFile.StartLine)
		if err != nil {
			return fmt.Errorf("parse sqls failed: %v", err)
		}
	}
	return nil
}

func buildOnlineTaskForAudit(c echo.Context, s *model.Storage, userId uint64, instanceName, instanceSchema, projectUid string, sqls getSQLFromFileResp) (*model.Task, error) {
	instance, exist, err := dms.GetInstanceInProjectByName(c.Request().Context(), projectUid, instanceName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrInstanceNoAccess
	}
	can, err := CheckCurrentUserCanAccessInstances(c.Request().Context(), projectUid, controller.GetUserID(c), []*model.Instance{instance})
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

	if err = plugin.Ping(context.TODO()); err != nil {
		return nil, err
	}

	task := &model.Task{
		Schema:       instanceSchema,
		InstanceId:   instance.ID,
		Instance:     instance,
		CreateUserId: userId,
		ExecuteSQLs:  []*model.ExecuteSQL{},
		SQLSource:    sqls.SourceType,
		DBType:       instance.DbType,
	}
	createAt := time.Now()
	task.CreatedAt = createAt
	err = addSQLsFromFileToTasks(sqls, task, plugin)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func buildOfflineTaskForAudit(userId uint64, dbType string, sqls getSQLFromFileResp) (*model.Task, error) {
	task := &model.Task{
		CreateUserId: userId,
		ExecuteSQLs:  []*model.ExecuteSQL{},
		SQLSource:    sqls.SourceType,
		DBType:       dbType,
	}
	var err error
	plugin, err := common.NewDriverManagerWithoutCfg(log.NewEntry(), dbType)
	if err != nil {
		return nil, fmt.Errorf("open plugin failed: %v", err)
	}
	defer plugin.Close(context.TODO())

	err = addSQLsFromFileToTasks(sqls, task, plugin)
	if err != nil {
		return nil, err
	}

	createAt := time.Now()
	task.CreatedAt = createAt

	return task, nil
}

func getSqlsFromZip(c echo.Context) (sqlsFromSQLFile []SQLsFromSQLFile, sqlsFromXML []SQLFromXML, exist bool, err error) {
	file, err := c.FormFile(InputZipFileName)
	if err == http.ErrMissingFile {
		return nil, nil, false, nil
	}
	if err != nil {
		return nil, nil, false, err
	}
	f, err := file.Open()
	if err != nil {
		return nil, nil, false, err
	}
	defer f.Close()

	currentPos, err := f.Seek(0, io.SeekEnd) // get size of zip file
	if err != nil {
		return nil, nil, false, err
	}
	size := currentPos + 1
	if size > maxZipFileSize {
		return nil, nil, false, fmt.Errorf("file can't be bigger than %vM", maxZipFileSize/1024/1024)
	}
	r, err := zip.NewReader(f, size)
	if err != nil {
		return nil, nil, false, err
	}

	var xmlContents []xmlParser.XmlFile
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
			return nil, nil, false, fmt.Errorf("open src file failed:  %v", err)
		}
		content, err := io.ReadAll(r)
		if err != nil {
			return nil, nil, false, fmt.Errorf("read src file failed:  %v", err)
		}

		if strings.HasSuffix(srcFile.Name, ".xml") {
			xmlContents = append(xmlContents, xmlParser.XmlFile{
				FilePath: srcFile.Name,
				Content:  string(content),
			})
		} else if strings.HasSuffix(srcFile.Name, ".sql") {
			sqlsFromSQLFile = append(sqlsFromSQLFile, SQLsFromSQLFile{
				FilePath: srcFile.Name,
				SQLs:     string(content),
			})
		}
	}

	// parse xml content
	// xml文件需要把所有文件内容同时解析，否则会无法解析跨namespace引用的SQL
	{
		sqlsFromXmls, err := parseXMLsWithFilePath(xmlContents)
		if err != nil {
			return nil, nil, false, err
		}
		sqlsFromXML = append(sqlsFromXML, sqlsFromXmls...)
	}

	return sqlsFromSQLFile, sqlsFromXML, true, nil
}

func parseXMLsWithFilePath(xmlContents []xmlParser.XmlFile) ([]SQLFromXML, error) {
	allStmtsFromXml, err := xmlParser.ParseXMLs(xmlContents, false)
	if err != nil {
		return nil, fmt.Errorf("parse sqls from xml failed: %v", err)
	}

	sqls := []SQLFromXML{}
	for _, stmtInfo := range allStmtsFromXml {
		sqls = append(sqls, SQLFromXML{
			FilePath:  stmtInfo.FilePath,
			StartLine: stmtInfo.StartLine,
			SQL:       stmtInfo.SQL,
		})
	}
	return sqls, nil
}

func getSqlsFromGit(c echo.Context) (sqlsFromSQLFiles, sqlsFromJavaFiles []SQLsFromSQLFile, sqlsFromXMLs []SQLFromXML, exist bool, err error) {
	// make a temp dir and clean up befor return
	dir, err := os.MkdirTemp("./", "git-repo-")
	if err != nil {
		return nil, nil, nil, false, err
	}
	defer os.RemoveAll(dir)
	// read http url from form and check if it's a git url
	url := c.FormValue(GitHttpURL)
	if !utils.IsGitHttpURL(url) {
		return nil, nil, nil, false, errors.New(errors.DataInvalid, fmt.Errorf("url is not a git url"))
	}
	cloneOpts := &goGit.CloneOptions{
		URL: url,
	}
	// public repository do not require an user name and password
	userName := c.FormValue(GitUserName)
	password := c.FormValue(GitPassword)
	if userName != "" {
		cloneOpts.Auth = &goGitTransport.BasicAuth{
			Username: userName,
			Password: password,
		}
	}
	// clone from git
	_, err = goGit.PlainCloneContext(c.Request().Context(), dir, false, cloneOpts)
	if err != nil {
		return nil, nil, nil, false, err
	}
	l := log.NewEntry().WithField("function", "getSqlsFromGit")
	var xmlContents []xmlParser.XmlFile
	// traverse the repository, parse and put SQL into sqlBuffer
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		gitPath := strings.TrimPrefix(path, strings.TrimPrefix(dir, "./"))
		if !info.IsDir() {
			var sqlBuffer strings.Builder
			switch {
			case strings.HasSuffix(path, ".xml"):
				content, err := os.ReadFile(path)
				if err != nil {
					l.Errorf("skip file [%v]. because read file failed: %v", path, err)
					return nil
				}
				xmlContents = append(xmlContents, xmlParser.XmlFile{
					FilePath: gitPath,
					Content:  string(content),
				})
			case strings.HasSuffix(path, ".sql"):
				content, err := os.ReadFile(path)
				if err != nil {
					l.Errorf("skip file [%v]. because read file failed: %v", path, err)
					return nil
				}
				sqlsFromSQLFiles = append(sqlsFromSQLFiles, SQLsFromSQLFile{
					FilePath: gitPath,
					SQLs:     string(content),
				})
			case strings.HasSuffix(path, ".java"):
				sqls, err := javaParser.GetSqlFromJavaFile(path)
				if err != nil {
					l.Errorf("skip file [%v]. because get sql from java file failed: %v", path, err)
					return nil
				}
				for _, sql := range sqls {
					if !strings.HasSuffix(sql, ";") {
						sql += ";"
					}
					_, err = sqlBuffer.WriteString(sql)
					if err != nil {
						return fmt.Errorf("gather sqls from java file failed: %v", err)
					}
				}
				sqlsFromJavaFiles = append(sqlsFromJavaFiles, SQLsFromSQLFile{
					FilePath: gitPath,
					SQLs:     sqlBuffer.String(),
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, nil, false, err
	}

	// parse xml content
	// xml文件需要把所有文件内容同时解析，否则会无法解析跨namespace引用的SQL
	sqlsFromXMLs, err = parseXMLsWithFilePath(xmlContents)
	if err != nil {
		return nil, nil, nil, false, err
	}

	return sqlsFromSQLFiles, sqlsFromJavaFiles, sqlsFromXMLs, true, nil
}

type UpdateSQLAuditRecordReqV1 struct {
	Tags []string `json:"tags" valid:"dive,tag_name"`
}

// UpdateSQLAuditRecordV1
// @Summary 更新SQL审核记录
// @Description update SQL audit record
// @Accept json
// @Id updateSQLAuditRecordV1
// @Tags sql_audit_record
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_audit_record_id path string true "sql audit record id"
// @Param param body v1.UpdateSQLAuditRecordReqV1 true "update SQL audit record"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_audit_records/{sql_audit_record_id}/ [patch]
func UpdateSQLAuditRecordV1(c echo.Context) error {
	req := new(UpdateSQLAuditRecordReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditRecordId := c.Param("sql_audit_record_id")

	s := model.GetStorage()

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	record, exist, err := s.GetSQLAuditRecordById(projectUid, auditRecordId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.ErrAccessDeniedError, fmt.Errorf("sql audit record id %v not exist", auditRecordId)))
	}
	instance, exist, err := dms.GetInstanceInProjectById(c.Request().Context(), projectUid, record.Task.InstanceId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		record.Task.Instance = instance
	}

	if req.Tags != nil {
		up, err := dms.NewUserPermission(user.GetIDStr(), projectUid)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden)
		}
		if err != nil {
			return errors.New(errors.ConnectStorageError, fmt.Errorf("check project manager failed: %v", err))
		}
		if !up.IsProjectAdmin() {
			if yes, err := s.IsSQLAuditRecordBelongToCurrentUser(user.GetIDStr(), projectUid, auditRecordId); err != nil {
				return controller.JSONBaseErrorReq(c, fmt.Errorf("check privilege failed: %v", err))
			} else if !yes {
				return controller.JSONBaseErrorReq(c, errors.New(errors.ErrAccessDeniedError, errors.NewAccessDeniedErr("you can't update SQL audit record that created by others")))
			}
		}

		data := model.SQLAuditRecordUpdateData{Tags: req.Tags}
		if err = s.UpdateSQLAuditRecordById(auditRecordId, data); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		go func() {
			err = syncSqlManage(record, req.Tags)
			if err != nil {
				log.NewEntry().WithField("sync_sql_audit_record", auditRecordId).Errorf("sync sql manager failed: %v", err)
			}
		}()
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetSQLAuditRecordsReqV1 struct {
	FuzzySearchTags         string `json:"fuzzy_search_tags" query:"fuzzy_search_tags"` // todo issue1811
	FilterSQLAuditStatus    string `json:"filter_sql_audit_status" query:"filter_sql_audit_status" enums:"auditing,successfully,"`
	FilterInstanceId        uint64 `json:"filter_instance_id" query:"filter_instance_id"`
	FilterCreateTimeFrom    string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo      string `json:"filter_create_time_to" query:"filter_create_time_to"`
	FilterSqlAuditRecordIDs string `json:"filter_sql_audit_record_ids" query:"filter_sql_audit_record_ids" example:"1711247438821462016,1711246967037759488"`
	PageIndex               uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize                uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type SQLAuditRecord struct {
	Creator          string                  `json:"creator"`
	SQLAuditRecordId string                  `json:"sql_audit_record_id"`
	SQLAuditStatus   string                  `json:"sql_audit_status"`
	Tags             []string                `json:"tags"`
	Instance         *SQLAuditRecordInstance `json:"instance,omitempty"`
	Task             *AuditTaskResV1         `json:"task,omitempty"`
	CreatedAt        *time.Time              `json:"created_at,omitempty"`
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
// @Param filter_instance_id query uint64 false "filter instance id"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param filter_sql_audit_record_ids query string false "filter sql audit record ids"
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
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), false)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	up, err := dms.NewUserPermission(user.GetIDStr(), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("check project manager failed: %v", err))
	}

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}

	data := map[string]interface{}{
		"filter_project_id":       projectUid,
		"filter_creator_id":       user.ID,
		"fuzzy_search_tags":       req.FuzzySearchTags,
		"filter_instance_id":      req.FilterInstanceId,
		"filter_create_time_from": req.FilterCreateTimeFrom,
		"filter_create_time_to":   req.FilterCreateTimeTo,
		"check_user_can_access":   !up.IsProjectAdmin(),
		"filter_audit_record_ids": req.FilterSqlAuditRecordIDs,
		"limit":                   req.PageSize,
		"offset":                  offset,
	}
	if req.FilterSQLAuditStatus == SQLAuditRecordStatusAuditing {
		data["filter_task_status_exclude"] = model.TaskStatusAudited
	} else if req.FilterSQLAuditStatus == SQLAuditRecordStatusSuccessfully {
		data["filter_task_status"] = model.TaskStatusAudited
	}

	records, total, err := s.GetSQLAuditRecordsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// batch get instance info
	instanceIds := make([]uint64, 0)
	for _, record := range records {
		instanceIds = append(instanceIds, record.InstanceId)
	}
	instances, err := dms.GetInstancesByIds(c.Request().Context(), instanceIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	idInstMap := make(map[uint64] /*instance id*/ *model.Instance /*instance */, 0)
	for _, instance := range instances {
		idInstMap[instance.ID] = instance
	}

	resData := make([]SQLAuditRecord, len(records))
	for i := range records {
		record := records[i]
		status := SQLAuditRecordStatusAuditing
		if record.TaskStatus == model.TaskStatusAudited {
			status = SQLAuditRecordStatusSuccessfully
		}
		var tags []string
		if record.Tags.Valid {
			if err := json.Unmarshal([]byte(record.Tags.String), &tags); err != nil {
				log.NewEntry().Errorf("parse tags failed,tags:%v , err: %v", record.Tags, err)
			}
		}

		resData[i] = SQLAuditRecord{
			Creator:          dms.GetUserNameWithDelTag(record.CreatorId),
			SQLAuditRecordId: record.AuditRecordId,
			SQLAuditStatus:   status,
			Tags:             tags,
			CreatedAt:        record.RecordCreatedAt,
			Task: &AuditTaskResV1{
				Id:             record.TaskId,
				InstanceDbType: record.DbType,
				InstanceSchema: record.InstanceSchema,
				AuditLevel:     record.AuditLevel.String,
				Score:          record.AuditScore.Int32,
				PassRate:       record.AuditPassRate.Float64,
				Status:         record.TaskStatus,
				SQLSource:      record.SQLSource,
			},
		}
		if inst, ok := idInstMap[record.InstanceId]; ok {
			resData[i].Instance = &SQLAuditRecordInstance{
				Host: inst.Host,
				Port: inst.Port,
			}
			resData[i].Task.InstanceName = idInstMap[record.InstanceId].Name
		}
	}
	return c.JSON(http.StatusOK, &GetSQLAuditRecordsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      resData,
		TotalNums: total,
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
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), false)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	auditRecordId := c.Param("sql_audit_record_id")

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	up, err := dms.NewUserPermission(user.GetIDStr(), projectUid)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden)
	}
	if err != nil {
		return errors.New(errors.ConnectStorageError, fmt.Errorf("check project manager failed: %v", err))
	}

	s := model.GetStorage()

	if !up.IsProjectAdmin() {
		if yes, err := s.IsSQLAuditRecordBelongToCurrentUser(user.GetIDStr(), projectUid, auditRecordId); err != nil {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("check privilege failed: %v", err))
		} else if !yes {
			return controller.JSONBaseErrorReq(c, errors.New(errors.ErrAccessDeniedError, errors.NewAccessDeniedErr("you can't see the SQL audit record because it isn't created by you")))
		}
	}

	record, exist, err := s.GetSQLAuditRecordById(projectUid, auditRecordId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, e.New("can not find record")))
	}
	// build ret info
	status := SQLAuditRecordStatusAuditing
	if record.Task.Status == model.TaskStatusAudited {
		status = SQLAuditRecordStatusSuccessfully
	}
	var tags []string
	if err := json.Unmarshal([]byte(record.Tags), &tags); err != nil {
		log.NewEntry().Errorf("parse tags failed,tags:%v , err: %v", record.Tags, err)
	}
	task := &AuditTaskResV1{
		Id:             record.Task.ID,
		InstanceDbType: record.Task.DBType,
		InstanceSchema: record.Task.Schema,
		AuditLevel:     record.Task.AuditLevel,
		Score:          record.Task.Score,
		PassRate:       record.Task.PassRate,
		Status:         record.Task.Status,
		SQLSource:      record.Task.SQLSource,
		ExecStartTime:  record.Task.ExecStartAt,
		ExecEndTime:    record.Task.ExecEndAt,
	}
	instance, exist, err := dms.GetInstancesById(c.Request().Context(), record.Task.InstanceId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		task.InstanceName = instance.Name
	}

	return c.JSON(http.StatusOK, &GetSQLAuditRecordResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: SQLAuditRecord{
			SQLAuditStatus:   status,
			Tags:             tags,
			SQLAuditRecordId: auditRecordId,
			Task:             task,
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
