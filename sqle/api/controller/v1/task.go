package v1

import (
	"actiontech.cloud/universe/sqle/v4/sqle/api/controller"
	"bytes"
	"encoding/csv"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"actiontech.cloud/universe/sqle/v4/sqle/executor"

	"actiontech.cloud/universe/sqle/v4/sqle/api/server"
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"actiontech.cloud/universe/sqle/v4/sqle/inspector"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/sqle/v4/sqle/model"

	"github.com/labstack/echo/v4"
)

const (
	SqlAuditTaskExpiredTime = "720h"
)

var TaskNoAccessError = errors.New(errors.DataNotExist, fmt.Errorf("task is not exist or you can't access it"))

type CreateAuditTaskReqV1 struct {
	InstanceName   string `json:"instance_name" form:"instance_name" example:"inst_1" valid:"required"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema" example:"db1" valid:"-"`
	Sql            string `json:"sql" example:"alter table tb1 drop columns c1" valid:"-"`
}

type GetAuditTaskResV1 struct {
	controller.BaseRes
	Data *AuditTaskResV1 `json:"data"`
}

type AuditTaskResV1 struct {
	Id             uint    `json:"task_id"`
	InstanceName   string  `json:"instance_name"`
	InstanceSchema string  `json:"instance_schema" example:"db1"`
	PassRate       float64 `json:"pass_rate"`
	Status         string  `json:"status" enums:"initialized,audited,executing,exec_success,exec_failed"`
}

func convertTaskToRes(task *model.Task) *AuditTaskResV1 {
	return &AuditTaskResV1{
		Id:             task.ID,
		InstanceName:   task.Instance.Name,
		InstanceSchema: task.Schema,
		PassRate:       task.PassRate,
		Status:         task.Status,
	}
}

// @Summary 创建Sql审核任务并提交审核
// @Description create and audit a task. NOTE: it will create a task with sqls from "sql" if "sql" isn't empty
// @Accept mpfd
// @Produce json
// @Tags task
// @Id createAndAuditTaskV1
// @Security ApiKeyAuth
// @Param instance_name formData string true "instance name"
// @Param instance_schema formData string false "schema of instance"
// @Param sql formData string false "sqls for audit"
// @Param input_sql_file formData file false "input SQL file"
// @Success 200 {object} v1.GetAuditTaskResV1
// @router /v1/tasks/audits [post]
func CreateAndAuditTask(c echo.Context) error {
	req := new(CreateAuditTaskReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	if req.Sql == "" {
		_, sqls, err := controller.ReadFileToByte(c, "input_sql_file")
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		req.Sql = string(sqls)
	}

	s := model.GetStorage()
	instance, exist, err := s.GetInstanceByName(req.InstanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, instanceNoAccessError)
	}

	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := executor.Ping(log.NewEntry(), instance); err != nil {
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
	}

	createAt := time.Now()
	task.CreatedAt = createAt

	nodes, err := inspector.NewInspector(log.NewEntry(), inspector.NewContext(nil), task, nil, nil).
		ParseSql(req.Sql)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	for n, node := range nodes {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(n + 1),
				Content: node.Text(),
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
	task, err = server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%d", task.ID), model.TASK_ACTION_AUDIT)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetAuditTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
}

func checkCurrentUserCanAccessTask(c echo.Context, task *model.Task) error {
	if controller.GetUserName(c) == defaultAdminUser {
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
	workflow, exist, err := s.GetWorkflowByTaskId(fmt.Sprintf("%d", task.ID))
	if err != nil {
		return err
	}
	if !exist {
		return TaskNoAccessError
	}
	access, err := s.UserCanAccessWorkflow(user, workflow)
	if err != nil {
		return err
	}
	if access {
		return nil
	}
	return TaskNoAccessError
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
		return controller.JSONBaseErrorReq(c, TaskNoAccessError)
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
	NoDuplicate       bool   `json:"no_duplicate" query:"no_duplicate"`
	PageIndex         uint32 `json:"page_index" query:"page_index" valid:"required,int"`
	PageSize          uint32 `json:"page_size" query:"page_size" valid:"required,int"`
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
}

// @Summary 获取指定审核任务的SQLs信息
// @Description get information of all SQLs belong to the specified audit task
// @Tags task
// @Id getAuditTaskSQLsV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed)
// @Param filter_audit_status query string false "filter: audit status of task sql" Enums(initialized,doing,finished)
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
		return controller.JSONBaseErrorReq(c, TaskNoAccessError)
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
	NoDuplicate string `json:"no_duplicate" query:"no_duplicate"`
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
		return controller.JSONBaseErrorReq(c, TaskNoAccessError)
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
	cw.Write([]string{"序号", "SQL", "SQL审核状态", "SQL审核结果", "SQL执行状态", "SQL执行结果", "SQL对应的回滚语句"})
	for _, td := range taskSQLsDetail {
		taskSql := &model.ExecuteSQL{
			AuditResult: td.AuditResult,
			AuditStatus: td.AuditStatus,
		}
		taskSql.ExecStatus = td.ExecStatus
		cw.Write([]string{
			strconv.FormatUint(uint64(td.Number), 10),
			td.ExecSQL,
			taskSql.GetAuditStatusDesc(),
			taskSql.GetAuditResultDesc(),
			taskSql.GetExecStatusDesc(),
			td.ExecResult,
			td.RollbackSQL.String,
		})
	}
	cw.Flush()
	fileName := fmt.Sprintf("SQL审核报告_%v_%v.csv", task.Instance.Name, taskId)
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
		return controller.JSONBaseErrorReq(c, TaskNoAccessError)
	}
	err = checkCurrentUserCanAccessTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	sqls, err := s.GetTaskExecuteSQLContent(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	fileName := fmt.Sprintf("exec_sql_%s_%s.sql", task.Instance.Name, taskId)
	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": fileName}))

	buff := &bytes.Buffer{}
	for _, sql := range sqls {
		buff.WriteString(strings.TrimRight(sql, ";"))
		buff.WriteString(";\n")
	}
	return c.Blob(http.StatusOK, echo.MIMETextPlain, buff.Bytes())
}
