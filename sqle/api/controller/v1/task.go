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

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"

	mybatis_parser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/labstack/echo/v4"
)

var ErrTaskNoAccess = errors.New(errors.DataNotExist, fmt.Errorf("task is not exist or you can't access it"))

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
	InstanceSchema string     `json:"instance_schema" example:"db1"`
	AuditLevel     string     `json:"audit_level" enums:"normal,notice,warn,error"`
	Score          int32      `json:"score"`
	PassRate       float64    `json:"pass_rate"`
	Status         string     `json:"status" enums:"initialized,audited,executing,exec_success,exec_failed"`
	SQLSource      string     `json:"sql_source" enums:"form_data,sql_file,mybatis_xml_file,audit_plan"`
	ExecStartTime  *time.Time `json:"exec_start_time,omitempty"`
	ExecEndTime    *time.Time `json:"exec_end_time,omitempty"`
}

func convertTaskToRes(task *model.Task) *AuditTaskResV1 {
	return &AuditTaskResV1{
		Id:             task.ID,
		InstanceName:   task.InstanceName(),
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
)

func getSQLFromFile(c echo.Context) (string, string, error) {
	// Read it from sql file.
	sql, exist, err := controller.ReadFileContent(c, InputSQLFileName)
	if err != nil {
		return "", model.TaskSQLSourceFromSQLFile, err
	}
	if exist {
		return sql, model.TaskSQLSourceFromSQLFile, nil
	}

	// If sql_file is not exist, read it from mybatis xml file.
	data, exist, err := controller.ReadFileContent(c, InputMyBatisXMLFileName)
	if err != nil {
		return "", model.TaskSQLSourceFromMyBatisXMLFile, err
	}
	if exist {
		sql, err := mybatis_parser.ParseXML(data)
		if err != nil {
			return "", model.TaskSQLSourceFromMyBatisXMLFile, errors.New(errors.ParseMyBatisXMLFileError, err)
		}
		return sql, model.TaskSQLSourceFromMyBatisXMLFile, nil
	}
	return "", "", errors.New(errors.DataInvalid, fmt.Errorf("input sql is empty"))
}

// @Summary 创建Sql审核任务并提交审核
// @Description create and audit a task, you can upload sql content in three ways, any one can be used, but only one is effective.
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Accept mpfd
// @Produce json
// @Tags task
// @Id createAndAuditTaskV1
// @Security ApiKeyAuth
// @Param instance_name formData string true "instance name"
// @Param instance_schema formData string false "schema of instance"
// @Param sql formData string false "sqls for audit"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Success 200 {object} v1.GetAuditTaskResV1
// @router /v1/tasks/audits [post]
func CreateAndAuditTask(c echo.Context) error {
	req := new(CreateAuditTaskReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	var sql string
	var source string
	var err error

	if req.Sql != "" {
		sql, source = req.Sql, model.TaskSQLSourceFromFormData
	} else {
		sql, source, err = getSQLFromFile(c)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	s := model.GetStorage()
	instance, exist, err := s.GetInstanceByName(req.InstanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errInstanceNoAccess)
	}

	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	d, err := newDriverWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return err
	}
	defer d.Close(context.TODO())
	if err := d.Ping(context.TODO()); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	task := &model.Task{
		Schema:       req.InstanceSchema,
		InstanceId:   instance.ID,
		Instance:     instance,
		CreateUserId: user.ID,
		ExecuteSQLs:  []*model.ExecuteSQL{},
		SQLSource:    source,
		DBType:       instance.DbType,
	}
	createAt := time.Now()
	task.CreatedAt = createAt

	nodes, err := d.Parse(context.TODO(), sql)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	for n, node := range nodes {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(n + 1),
				Content: node.Text,
			},
		})
	}
	// if task instance is not nil, gorm will update instance when save task.
	task.Instance = nil
	err = s.Save(task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	task.Instance = instance
	task, err = server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%d", task.ID), server.ActionTypeAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetAuditTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
}

func checkCurrentUserCanAccessTask(c echo.Context, task *model.Task) error {
	if controller.GetUserName(c) == model.DefaultAdminUser {
		return nil
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	if user.ID == task.CreateUserId {
		return nil
	}
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByTaskId(task.ID)
	if err != nil {
		return err
	}
	if !exist {
		return ErrTaskNoAccess
	}
	access, err := s.UserCanAccessWorkflow(user, workflow)
	if err != nil {
		return err
	}
	if !access {
		return ErrTaskNoAccess
	}
	return nil
}

// @Summary 获取Sql审核任务信息
// @Description get task
// @Tags task
// @Id getAuditTaskV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 {object} v1.GetAuditTaskResV1
// @router /v1/tasks/audits/{task_id}/ [get]
func GetTask(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}
	err = checkCurrentUserCanAccessTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
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

// @Summary 获取指定审核任务的SQLs信息
// @Description get information of all SQLs belong to the specified audit task
// @Tags task
// @Id getAuditTaskSQLsV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed)
// @Param filter_audit_status query string false "filter: audit status of task sql" Enums(initialized,doing,finished)
// @Param filter_audit_level query string false "filter: audit level of task sql" Enums(normal,notice,warn,error)
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Param page_index query string false "page index"
// @Param page_size query string false "page size"
// @Success 200 {object} v1.GetAuditTaskSQLsResV1
// @router /v1/tasks/audits/{task_id}/sqls [get]
func GetTaskSQLs(c echo.Context) error {
	req := new(GetAuditTaskSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}
	err = checkCurrentUserCanAccessTask(c, task)
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
			AuditResult: taskSQL.AuditResult,
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

// @Summary 下载指定审核任务的SQLs信息报告
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
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}
	err = checkCurrentUserCanAccessTask(c, task)
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
			AuditResult: td.AuditResult,
			AuditStatus: td.AuditStatus,
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

// @Summary 下载指定审核任务的SQL文件
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
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}
	err = checkCurrentUserCanAccessTask(c, task)
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

// @Summary 获取指定审核任务的SQL内容
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
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}
	err = checkCurrentUserCanAccessTask(c, task)
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

// @Summary 修改审核任务中某条SQL的相关信息
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
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}
	err = checkCurrentUserCanAccessTask(c, task)
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
