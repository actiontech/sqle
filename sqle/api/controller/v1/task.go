package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	mybatis_parser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

var ErrTooManyDataSource = errors.New(errors.DataConflict, fmt.Errorf("the number of data sources must be less than %v", MaximumDataSourceNum))

type CreateAuditTaskReqV1 struct {
	InstanceName   string `json:"instance_name" form:"instance_name" example:"inst_1" valid:"required"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema" example:"db1"`
	Sql            string `json:"sql" form:"sql" example:"alter table tb1 drop columns c1"`
}

type GetAuditTaskResV1 struct {
	controller.BaseRes
	Data *AuditTaskResV1 `json:"data"`
}

type AuditTaskResV1 struct {
	Id             uint       `json:"task_id"`
	InstanceName   string     `json:"instance_name"`
	InstanceDbType string     `json:"instance_db_type"`
	InstanceSchema string     `json:"instance_schema" example:"db1"`
	AuditLevel     string     `json:"audit_level" enums:"normal,notice,warn,error,"`
	Score          int32      `json:"score"`
	PassRate       float64    `json:"pass_rate"`
	Status         string     `json:"status" enums:"initialized,audited,executing,exec_success,exec_failed,manually_executed"`
	SQLSource      string     `json:"sql_source" enums:"form_data,sql_file,mybatis_xml_file,audit_plan"`
	ExecStartTime  *time.Time `json:"exec_start_time,omitempty"`
	ExecEndTime    *time.Time `json:"exec_end_time,omitempty"`
}

func convertTaskToRes(task *model.Task) *AuditTaskResV1 {
	return &AuditTaskResV1{
		Id:             task.ID,
		InstanceName:   task.InstanceName(),
		InstanceDbType: task.DBType,
		InstanceSchema: task.Schema,
		AuditLevel:     task.AuditLevel,
		Score:          task.Score,
		PassRate:       task.PassRate,
		Status:         task.Status,
		SQLSource:      task.SQLSource,
		ExecStartTime:  task.ExecStartAt,
		ExecEndTime:    task.ExecEndAt,
	}
}

const (
	InputSQLFileName        = "input_sql_file"
	InputMyBatisXMLFileName = "input_mybatis_xml_file"
	InputZipFileName        = "input_zip_file"
	GitHttpURL              = "git_http_url"
	GitUserName             = "git_user_name"
	GitPassword             = "git_user_password"
)

func getSQLFromFile(c echo.Context) (getSQLFromFileResp, error) {
	// Read it from sql file.
	fileName, sqlsFromSQLFile, exist, err := controller.ReadFile(c, InputSQLFileName)
	if err != nil {
		return getSQLFromFileResp{}, err
	}
	if exist {
		return getSQLFromFileResp{
			SourceType: model.TaskSQLSourceFromSQLFile,
			SQLsFromSQLFiles: []SQLsFromSQLFile{{
				FilePath: fileName,
				SQLs:     sqlsFromSQLFile}},
		}, nil
	}

	// If sql_file is not exist, read it from mybatis xml file.
	fileName, data, exist, err := controller.ReadFile(c, InputMyBatisXMLFileName)
	if err != nil {
		return getSQLFromFileResp{}, err
	}
	if exist {
		sqls, err := mybatis_parser.ParseXMLs([]mybatis_parser.XmlFile{{Content: data}}, true)
		if err != nil {
			return getSQLFromFileResp{}, errors.New(errors.ParseMyBatisXMLFileError, err)
		}
		sqlsFromXMLs := make([]SQLFromXML, len(sqls))
		for i := range sqls {
			sqlsFromXMLs[i] = SQLFromXML{
				FilePath:  fileName,
				StartLine: sqls[i].StartLine,
				SQL:       sqls[i].SQL,
			}
		}
		return getSQLFromFileResp{
			SourceType:   model.TaskSQLSourceFromMyBatisXMLFile,
			SQLsFromXMLs: sqlsFromXMLs,
		}, nil
	}

	// If mybatis xml file is not exist, read it from zip file.
	sqlsFromSQLFiles, sqlsFromXML, exist, err := getSqlsFromZip(c)
	if err != nil {
		return getSQLFromFileResp{}, err
	}
	if exist {
		return getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromZipFile,
			SQLsFromSQLFiles: sqlsFromSQLFiles,
			SQLsFromXMLs:     sqlsFromXML,
		}, nil
	}

	// If zip file is not exist, read it from git repository
	sqlsFromSQLFiles, sqlsFromJavaFiles, sqlsFromXMLs, exist, err := getSqlsFromGit(c)
	if err != nil {
		return getSQLFromFileResp{}, err
	}
	if exist {
		return getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromGitRepository,
			SQLsFromSQLFiles: append(sqlsFromSQLFiles, sqlsFromJavaFiles...),
			SQLsFromXMLs:     sqlsFromXMLs,
		}, nil
	}
	return getSQLFromFileResp{}, errors.New(errors.DataInvalid, fmt.Errorf("input sql is empty"))
}

// @Summary 创建Sql扫描任务并提交审核
// @Description create and audit a task, you can upload sql content in three ways, any one can be used, but only one is effective.
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Accept mpfd
// @Produce json
// @Tags task
// @Id createAndAuditTaskV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name formData string true "instance name"
// @Param instance_schema formData string false "schema of instance"
// @Param sql formData string false "sqls for audit"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Success 200 {object} v1.GetAuditTaskResV1
// @router /v1/projects/{project_name}/tasks/audits [post]
func CreateAndAuditTask(c echo.Context) error {
	req := new(CreateAuditTaskReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	var sqls getSQLFromFileResp
	var err error

	if req.Sql != "" {
		sqls = getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromFormData,
			SQLsFromFormData: req.Sql,
		}
	} else {
		sqls, err = getSQLFromFile(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	task, err := buildOnlineTaskForAudit(c, s, uint64(user.ID), req.InstanceName, req.InstanceSchema, projectUid, sqls)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// if task instance is not nil, gorm will update instance when save task.
	tmpInst := *task.Instance
	task.Instance = nil

	taskGroup := model.TaskGroup{Tasks: []*model.Task{task}}
	err = s.Save(&taskGroup)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	task.Instance = &tmpInst
	task, err = server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%d", task.ID), server.ActionTypeAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetAuditTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
}

// @Summary 获取Sql扫描任务信息
// @Description get task
// @Tags task
// @Id getAuditTaskV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 {object} v1.GetAuditTaskResV1
// @router /v1/tasks/audits/{task_id}/ [get]
func GetTask(c echo.Context) error {
	taskId := c.Param("task_id")
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
}
func GetTaskById(ctx context.Context, taskId string) (*model.Task, error) {
	return getTaskById(ctx, taskId)
}

func getTaskById(ctx context.Context, taskId string) (*model.Task, error) {
	s := model.GetStorage()
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.NewTaskNoExistOrNoAccessErr()
	}

	instance, exist, err := dms.GetInstancesById(ctx, task.InstanceId)
	if err != nil {
		return nil, err
	}
	task.Instance = instance

	return task, nil
}

type GetAuditTaskSQLsReqV1 struct {
	FilterExecStatus  string `json:"filter_exec_status" query:"filter_exec_status"`
	FilterAuditStatus string `json:"filter_audit_status" query:"filter_audit_status"`
	FilterAuditLevel  string `json:"filter_audit_level" query:"filter_audit_level"`
	NoDuplicate       bool   `json:"no_duplicate" query:"no_duplicate"`
	PageIndex         uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize          uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditTaskSQLsResV1 struct {
	controller.BaseRes
	Data      []*AuditTaskSQLResV1 `json:"data"`
	TotalNums uint64               `json:"total_nums"`
}

type AuditTaskSQLResV1 struct {
	Number      uint   `json:"number"`
	ExecSQL     string `json:"exec_sql"`
	AuditResult string `json:"audit_result"`
	AuditLevel  string `json:"audit_level"`
	AuditStatus string `json:"audit_status"`
	ExecResult  string `json:"exec_result"`
	ExecStatus  string `json:"exec_status"`
	RollbackSQL string `json:"rollback_sql,omitempty"`
	Description string `json:"description"`
}

// @Summary 获取指定扫描任务的SQLs信息
// @Description get information of all SQLs belong to the specified audit task
// @Tags task
// @Id getAuditTaskSQLsV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed,manually_executed)
// @Param filter_audit_status query string false "filter: audit status of task sql" Enums(initialized,doing,finished)
// @Param filter_audit_level query string false "filter: audit level of task sql" Enums(normal,notice,warn,error)
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v1.GetAuditTaskSQLsResV1
// @router /v1/tasks/audits/{task_id}/sqls [get]
func GetTaskSQLs(c echo.Context) error {
	req := new(GetAuditTaskSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTaskDMS(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"task_id":             taskId,
		"filter_exec_status":  req.FilterExecStatus,
		"filter_audit_status": req.FilterAuditStatus,
		"filter_audit_level":  req.FilterAuditLevel,
		"no_duplicate":        req.NoDuplicate,
		"limit":               req.PageSize,
		"offset":              offset,
	}

	taskSQLs, count, err := s.GetTaskSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskSQLsRes := make([]*AuditTaskSQLResV1, 0, len(taskSQLs))
	for _, taskSQL := range taskSQLs {
		taskSQLRes := &AuditTaskSQLResV1{
			Number:      taskSQL.Number,
			Description: taskSQL.Description,
			ExecSQL:     taskSQL.ExecSQL,
			AuditResult: taskSQL.GetAuditResults(),
			AuditLevel:  taskSQL.AuditLevel,
			AuditStatus: taskSQL.AuditStatus,
			ExecResult:  taskSQL.ExecResult,
			ExecStatus:  taskSQL.ExecStatus,
			RollbackSQL: taskSQL.RollbackSQL.String,
		}
		taskSQLsRes = append(taskSQLsRes, taskSQLRes)
	}

	return c.JSON(http.StatusOK, &GetAuditTaskSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      taskSQLsRes,
		TotalNums: count,
	})
}

type DownloadAuditTaskSQLsFileReqV1 struct {
	NoDuplicate bool `json:"no_duplicate" query:"no_duplicate"`
}

// @Summary 下载指定扫描任务的SQLs信息报告
// @Description download report file of all SQLs information belong to the specified audit task
// @Tags task
// @Id downloadAuditTaskSQLReportV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Success 200 file 1 "sql report csv file"
// @router /v1/tasks/audits/{task_id}/sql_report [get]
func DownloadTaskSQLReportFile(c echo.Context) error {
	req := new(DownloadAuditTaskSQLsFileReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := map[string]interface{}{
		"task_id":      taskId,
		"no_duplicate": req.NoDuplicate,
	}

	taskSQLsDetail, _, err := s.GetTaskSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	buff := &bytes.Buffer{}
	buff.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	cw := csv.NewWriter(buff)
	err = cw.Write([]string{"序号", "SQL", "SQL审核状态", "SQL审核结果", "SQL执行状态", "SQL执行结果", "SQL对应的回滚语句", "SQL描述"})
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.WriteDataToTheFileError, err))
	}
	for _, td := range taskSQLsDetail {
		taskSql := &model.ExecuteSQL{
			AuditResults: td.AuditResults,
			AuditStatus:  td.AuditStatus,
		}
		taskSql.ExecStatus = td.ExecStatus
		err := cw.Write([]string{
			strconv.FormatUint(uint64(td.Number), 10),
			td.ExecSQL,
			taskSql.GetAuditStatusDesc(),
			taskSql.GetAuditResultDesc(),
			taskSql.GetExecStatusDesc(),
			td.ExecResult,
			td.RollbackSQL.String,
			td.Description,
		})
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.WriteDataToTheFileError, err))
		}
	}
	cw.Flush()
	fileName := fmt.Sprintf("SQL审核报告_%v_%v.csv", task.InstanceName(), taskId)
	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": fileName}))
	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}

// @Summary 下载指定扫描任务的SQL文件
// @Description download SQL file for the audit task
// @Tags task
// @Id downloadAuditTaskSQLFileV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 file 1 "sql file"
// @router /v1/tasks/audits/{task_id}/sql_file [get]
func DownloadTaskSQLFile(c echo.Context) error {
	taskId := c.Param("task_id")
	s := model.GetStorage()
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	content, err := s.GetTaskExecuteSQLContent(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	fileName := fmt.Sprintf("exec_sql_%s_%s.sql", task.InstanceName(), taskId)
	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": fileName}))

	return c.Blob(http.StatusOK, echo.MIMETextPlain, content)
}

type GetAuditTaskSQLContentResV1 struct {
	controller.BaseRes
	Data *AuditTaskSQLContentResV1 `json:"data"`
}

type AuditTaskSQLContentResV1 struct {
	Sql string `json:"sql" example:"alter table tb1 drop columns c1"`
}

// @Summary 获取指定扫描任务的SQL内容
// @Description get SQL content for the audit task
// @Tags task
// @Id getAuditTaskSQLContentV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 {object} v1.GetAuditTaskSQLContentResV1
// @router /v1/tasks/audits/{task_id}/sql_content [get]
func GetAuditTaskSQLContent(c echo.Context) error {
	taskId := c.Param("task_id")
	s := model.GetStorage()

	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	content, err := s.GetTaskExecuteSQLContent(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetAuditTaskSQLContentResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &AuditTaskSQLContentResV1{
			Sql: string(content),
		},
	})
}

type UpdateAuditTaskSQLsReqV1 struct {
	Description string `json:"description"`
}

// @Summary 修改扫描任务中某条SQL的相关信息
// @Description modify the relevant information of a certain SQL in the audit task
// @Tags task
// @Id updateAuditTaskSQLsV1
// @Accept json
// @Param task_id path string true "task id"
// @Param number path string true "sql number"
// @Param audit_plan body v1.UpdateAuditTaskSQLsReqV1 true "modify the relevant information of a certain SQL in the audit task"
// @Security ApiKeyAuth
// @Success 200 {object} controller.BaseRes
// @router /v1/tasks/audits/{task_id}/sqls/{number} [patch]
func UpdateAuditTaskSQLs(c echo.Context) error {
	req := new(UpdateAuditTaskSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	taskId := c.Param("task_id")
	number := c.Param("number")

	s := model.GetStorage()
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	taskSql, exist, err := s.GetTaskSQLByNumber(taskId, number)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("sql number not found")))
	}
	// the user may leave the description blank to clear the description, so no processing is performed
	taskSql.Description = req.Description
	err = s.Save(taskSql)
	return controller.JSONBaseErrorReq(c, err)
}

func CheckCurrentUserCanViewTask(c echo.Context, task *model.Task) (err error) {
	return checkCurrentUserCanAccessTask(c, task, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeViewOthersWorkflow})
}

// TODO 使用DMS的权限校验
func CheckCurrentUserCanViewTaskDMS(c echo.Context, task *model.Task) error {
	_, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return err
	}
	return nil
}

type SQLExplain struct {
	SQL     string `json:"sql"`
	Message string `json:"message"`
	// explain result in table format
	ClassicResult ExplainClassicResult `json:"classic_result"`
}

type ExplainClassicResult struct {
	Rows []map[string] /* head name */ string `json:"rows"`
	Head []TableMetaItemHeadResV1             `json:"head"`
}

type GetTaskAnalysisDataResItemV1 struct {
	SQLExplain SQLExplain  `json:"sql_explain"`
	TableMetas []TableMeta `json:"table_metas"`
}

type GetTaskAnalysisDataResV1 struct {
	controller.BaseRes
	Data GetTaskAnalysisDataResItemV1 `json:"data"`
}

// GetTaskAnalysisData get SQL explain and related table metadata for analysis
// @Summary 获取task相关的SQL执行计划和表元数据
// @Description get SQL explain and related table metadata for analysis
// @Id getTaskAnalysisData
// @Tags task
// @Param task_id path string true "task id"
// @Param number path uint true "sql number"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskAnalysisDataResV1
// @router /v1/tasks/audits/{task_id}/sqls/{number}/analysis [get]
func GetTaskAnalysisData(c echo.Context) error {
	return getTaskAnalysisData(c)
}

type CreateAuditTasksGroupReqV1 struct {
	Instances []*InstanceForCreatingTask `json:"instances" valid:"dive,required"`
}

type InstanceForCreatingTask struct {
	InstanceName   string `json:"instance_name" form:"instance_name" valid:"required"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema"`
}

type CreateAuditTasksGroupResV1 struct {
	controller.BaseRes
	Data AuditTasksGroupResV1 `json:"data"`
}

type AuditTasksGroupResV1 struct {
	TaskGroupId uint `json:"task_group_id" form:"task_group_id" valid:"required"`
}

const MaximumDataSourceNum = 10

// CreateAuditTasksGroupV1
// @Summary 创建审核任务组
// @Description create tasks group.
// @Accept json
// @Produce json
// @Tags task
// @Id createAuditTasksV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param req body v1.CreateAuditTasksGroupReqV1 true "parameters for creating audit tasks group"
// @Success 200 {object} v1.CreateAuditTasksGroupResV1
// @router /v1/projects/{project_name}/task_groups [post]
func CreateAuditTasksGroupV1(c echo.Context) error {
	req := new(CreateAuditTasksGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 数据源个数最大为10
	if len(req.Instances) > MaximumDataSourceNum {
		return controller.JSONBaseErrorReq(c, ErrTooManyDataSource)
	}

	instNames := make([]string, len(req.Instances))
	for i, instance := range req.Instances {
		instNames[i] = instance.InstanceName
	}

	distinctInstNames := utils.RemoveDuplicate(instNames)

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	instances, err := dms.GetInstancesInProjectByNames(c.Request().Context(), projectUid, distinctInstNames)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	nameInstanceMap := make(map[string]*model.Instance, len(req.Instances))
	for _, inst := range instances {
		// https://github.com/actiontech/sqle/issues/1673
		inst, exist, err := dms.GetInstanceInProjectByName(c.Request().Context(), projectUid, inst.Name)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
		}

		can, err := CheckCurrentUserCanAccessInstances(c.Request().Context(), projectUid, user.GetIDStr(), instances)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !can {
			return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
		}

		if err := common.CheckInstanceIsConnectable(inst); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		nameInstanceMap[inst.Name] = inst
	}

	tasks := make([]*model.Task, len(req.Instances))
	for i, reqInstance := range req.Instances {
		tasks[i] = &model.Task{
			Schema:       reqInstance.InstanceSchema,
			InstanceId:   nameInstanceMap[reqInstance.InstanceName].ID,
			CreateUserId: uint64(user.ID),
			DBType:       nameInstanceMap[reqInstance.InstanceName].DbType,
		}
		tasks[i].CreatedAt = time.Now()
	}

	taskGroup := model.TaskGroup{Tasks: tasks}
	if err := s.Save(&taskGroup); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, CreateAuditTasksGroupResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditTasksGroupResV1{
			TaskGroupId: taskGroup.ID,
		},
	})
}

type AuditTaskGroupReqV1 struct {
	TaskGroupId uint   `json:"task_group_id" form:"task_group_id" valid:"required"`
	Sql         string `json:"sql" form:"sql" example:"alter table tb1 drop columns c1"`
}

type AuditTaskGroupRes struct {
	TaskGroupId uint              `json:"task_group_id"`
	Tasks       []*AuditTaskResV1 `json:"tasks"`
}

type AuditTaskGroupResV1 struct {
	controller.BaseRes
	Data AuditTaskGroupRes `json:"data"`
}

// AuditTaskGroupV1
// @Summary 审核任务组
// @Description audit task group.
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Description 4. file[input_zip_file]: it is zip file, sql will be parsed from it.
// @Accept mpfd
// @Produce json
// @Tags task
// @Id auditTaskGroupIdV1
// @Security ApiKeyAuth
// @Param task_group_id formData uint true "group id of tasks"
// @Param sql formData string false "sqls for audit"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Success 200 {object} v1.AuditTaskGroupResV1
// @router /v1/task_groups/audit [post]
func AuditTaskGroupV1(c echo.Context) error {
	req := new(AuditTaskGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	var err error
	var sqls getSQLFromFileResp
	if req.Sql != "" {
		sqls = getSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromFormData,
			SQLsFromFormData: req.Sql,
		}
	} else {
		sqls, err = getSQLFromFile(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	s := model.GetStorage()
	taskGroup, err := s.GetTaskGroupByGroupId(req.TaskGroupId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	tasks := taskGroup.Tasks

	{
		instanceIds := make([]uint64, 0, len(tasks))
		for _, task := range tasks {
			instanceIds = append(instanceIds, task.InstanceId)
		}

		instances, err := dms.GetInstancesByIds(c.Request().Context(), instanceIds)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		// 因为这个接口数据源属于同一个项目,取第一个DB所属项目
		projectId := instances[0].ProjectId
		can, err := CheckCurrentUserCanAccessInstances(c.Request().Context(), projectId, controller.GetUserID(c), instances)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !can {
			return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
		}

		l := log.NewEntry()

		// 因为这个接口数据源必然相同，所以只取第一个实例的DbType即可
		dbType := instances[0].DbType
		plugin, err := common.NewDriverManagerWithoutCfg(l, dbType)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		defer plugin.Close(context.TODO())

		for _, task := range tasks {
			err := addSQLsFromFileToTasks(sqls, task, plugin)
			if err != nil {
				return controller.JSONBaseErrorReq(c, errors.New(errors.GenericError, fmt.Errorf("add sqls from file to task failed: %v", err)))
			}
		}
	}

	if err := s.Save(taskGroup); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	for i, task := range tasks {
		if task.Status != model.TaskStatusInit {
			continue
		}

		tasks[i], err = server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%d", task.ID), server.ActionTypeAudit)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	tasksRes := make([]*AuditTaskResV1, len(tasks))
	for i, task := range tasks {
		tasksRes[i] = &AuditTaskResV1{
			Id:             task.ID,
			InstanceName:   task.InstanceName(),
			InstanceDbType: task.DBType,
			InstanceSchema: task.Schema,
			AuditLevel:     task.AuditLevel,
			Score:          task.Score,
			PassRate:       task.PassRate,
			Status:         task.Status,
			SQLSource:      task.SQLSource,
			ExecStartTime:  task.ExecStartAt,
			ExecEndTime:    task.ExecEndAt,
		}
	}

	return c.JSON(http.StatusOK, AuditTaskGroupResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditTaskGroupRes{
			TaskGroupId: taskGroup.ID,
			Tasks:       tasksRes,
		},
	})
}
