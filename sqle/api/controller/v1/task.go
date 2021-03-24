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

type CreateTaskReqV1 struct {
	InstanceName   string `json:"instance_name" form:"instance_name" example:"inst_1" valid:"required"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema" example:"db1" valid:"-"`
	Sql            string `json:"sql" example:"alter table tb1 drop columns c1" valid:"-"`
}

type GetTaskResV1 struct {
	controller.BaseRes
	Data *TaskResV1 `json:"data"`
}

type TaskResV1 struct {
	Id             uint    `json:"task_id"`
	InstanceName   string  `json:"instance_name"`
	InstanceSchema string  `json:"instance_schema" example:"db1"`
	PassRate       float64 `json:"pass_rate"`
	Status         string  `json:"status" enums:"initialized, audited, executing, exec_success, exec_failed"`
}

func createTask(c echo.Context) (*model.Task, controller.BaseRes) {
	req := new(CreateTaskReqV1)
	if err := c.Bind(req); err != nil {
		return nil, controller.NewBaseReq(err)
	}
	if err := c.Validate(req); err != nil {
		return nil, controller.NewBaseReq(err)
	}
	return createTaskByRequestParam(req)
}

func createTaskByRequestParam(req *CreateTaskReqV1) (*model.Task, controller.BaseRes) {
	s := model.GetStorage()
	instance, exist, err := s.GetInstanceByName(req.InstanceName)
	if err != nil {
		return nil, controller.NewBaseReq(err)
	}
	if !exist {
		return nil, controller.NewBaseReq(errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist")))
	}
	if err := executor.Ping(log.NewEntry(), instance); err != nil {
		return nil, controller.NewBaseReq(err)
	}

	task := &model.Task{
		Schema:      req.InstanceSchema,
		InstanceId:  instance.ID,
		Instance:    instance,
		ExecuteSQLs: []*model.ExecuteSQL{},
	}

	createAt := time.Now()
	task.CreatedAt = createAt

	nodes, err := inspector.NewInspector(log.NewEntry(), inspector.NewContext(nil), task, nil, nil).
		ParseSql(req.Sql)
	if err != nil {
		return nil, controller.NewBaseReq(err)
	}
	for n, node := range nodes {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(n + 1),
				Content: node.Text(),
			},
		})
	}
	task.Instance = nil // if task instance is not nil, gorm will update instance when save task.
	err = s.Save(task)
	if err != nil {
		return nil, controller.NewBaseReq(err)
	}
	task.Instance = instance
	return task, controller.NewBaseReq(nil)
}

func convertTaskToRes(task *model.Task) *TaskResV1 {
	return &TaskResV1{
		Id:             task.ID,
		InstanceName:   task.Instance.Name,
		InstanceSchema: task.Schema,
		PassRate:       task.PassRate,
		Status:         task.Status,
	}
}

// @Summary 创建Sql审核任务
// @Description create a task
// @Accept json
// @Produce json
// @Tags task
// @Id createTaskV1
// @Security ApiKeyAuth
// @Param instance body v1.CreateTaskReqV1 true "add task request"
// @Success 200 {object} v1.GetTaskResV1
// @router /v1/tasks [post]
func CreateTask(c echo.Context) error {
	task, res := createTask(c)
	if res.Code != 0 {
		return c.JSON(http.StatusOK, res)
	}
	return c.JSON(http.StatusOK, &GetTaskResV1{
		BaseRes: res,
		Data:    convertTaskToRes(task),
	})
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
// @Success 200 {object} v1.GetTaskResV1
// @router /v1/task/audit [post]
func CreateAndAuditTask(c echo.Context) error {
	//check params
	req := new(CreateTaskReqV1)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}

	if "" == req.Sql {
		_, sqls, err := controller.ReadFileToByte(c, "input_sql_file")
		if err != nil {
			return c.JSON(http.StatusOK, controller.NewBaseReq(err))
		}

		req.Sql = string(sqls)
	}

	task, res := createTaskByRequestParam(req)
	if res.Code != 0 {
		return c.JSON(http.StatusOK, res)
	}

	task, err := server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%d", task.ID), model.TASK_ACTION_AUDIT)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetTaskResV1{
		BaseRes: res,
		Data:    convertTaskToRes(task),
	})
}

// @Summary 获取Sql审核任务信息
// @Description get task
// @Tags task
// @Id getTaskV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 {object} v1.GetTaskResV1
// @router /v1/tasks/{task_id}/ [get]
func GetTask(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, controller.NewBaseReq(
			errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))))
	}
	return c.JSON(http.StatusOK, &GetTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
}

//// @Summary 删除审核任务
//// @Description delete task
//// @Param task_id path string true "Task ID"
//// @Success 200 {object} controller.BaseRes
//// @router /tasks/{task_id}/ [delete]
//func DeleteTask(c echo.Context) error {
//	s := model.GetStorage()
//	taskId := c.Param("task_id")
//	task, exist, err := s.GetTaskById(taskId)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	if !exist {
//		return c.JSON(http.StatusOK, TASK_NOT_EXIST)
//	}
//
//	// must check task not running
//
//	err = s.Delete(task)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
//}

//// @Summary 批量删除审核任务
//// @Description delete tasks by ids
//// @Accept x-www-form-urlencoded
//// @Param task_ids formData string true "remove tasks by ids(interlaced by ',')"
//// @Success 200 {object} controller.BaseRes
//// @router /tasks/remove_by_task_ids [post]
//func DeleteTasks(c echo.Context) error {
//	s := model.GetStorage()
//	taskIds, err := url.QueryUnescape(c.FormValue("task_ids"))
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	deleteTaskIds := strings.Split(strings.TrimRight(taskIds, ","), ",")
//
//	err = s.HardDeleteTasksByIds(deleteTaskIds)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	//err = s.HardDeleteRollbackSqlByTaskIds(deleteTaskIds)
//	//if err != nil {
//	//	return c.JSON(http.StatusOK, NewBaseReq(err))
//	//}
//	//err = s.HardDeleteSqlCommittingResultByTaskIds(deleteTaskIds)
//	//if err != nil {
//	//	return c.JSON(http.StatusOK, NewBaseReq(err))
//	//}
//	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
//}

//type GetAllTaskRes struct {
//	controller.BaseRes
//	Data GetAllTaskResult `json:"data"`
//}
//
//type GetAllTaskResult struct {
//	Tasks     []*model.TaskListDetail `json:"tasks"`
//	TotalNums uint64                  `json:"total_nums"`
//}
//
//// @Summary Sql审核列表
//// @Description get all tasks
//// @Param filter_task_id query uint32 false "filter instance name"
//// @Param filter_instance_name query string false "filter instance name"
//// @Param filter_task_type query string false "filter task type"
//// @Param filter_task_status query string false "filter instance status"
//// @Param filter_create_time_from query string false "filter create time from"
//// @Param filter_create_time_to query string false "filter create time to"
//// @Param page_index query uint32 false "page index"
//// @Param page_size query uint32 false "size of per page"
//// @Success 200 {object} controller.GetAllTaskRes
//// @router /tasks [get]
//func GetTasks(c echo.Context) error {
//	req := new(model.GetTasksReq)
//	if err := c.Bind(req); err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	if err := c.Validate(req); err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//
//	s := model.GetStorage()
//	tasks, totalNums, err := s.GetTasksByReq(req)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	return c.JSON(http.StatusOK, &GetAllTaskRes{
//		BaseRes: controller.NewBaseReq(nil),
//		Data: GetAllTaskResult{
//			Tasks:     tasks,
//			TotalNums: totalNums,
//		},
//	})
//}

// @Summary Sql提交审核
// @Description audit sql
// @Tags task
// @Id auditTaskV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 {object} v1.GetTaskResV1
// @router /tasks/{task_id}/audit [post]
func AuditTask(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, controller.NewBaseReq(
			errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))))
	}
	if task.Instance == nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(
			errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))))
	}
	task, err = server.GetSqled().AddTaskWaitResult(taskId, model.TASK_ACTION_AUDIT)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTaskToRes(task),
	})
}

//type CommitTaskRes struct {
//	controller.BaseRes
//	Data CommitTaskResult `json:"data"`
//}
//
//type CommitTaskResult struct {
//	TaskExecStatus string `json:"task_exec_status"`
//}
//
//// @Summary Sql提交上线
//// @Description commit sql
//// @Accept x-www-form-urlencoded
//// @Param task_id path string true "Task ID"
//// @Param is_sync formData boolean false "the request is sync or async."
//// @Success 200 {object} controller.CommitTaskRes
//// @router /tasks/{task_id}/commit [post]
//func CommitTask(c echo.Context) error {
//	s := model.GetStorage()
//	isSync := c.FormValue("is_sync")
//	taskId := c.Param("task_id")
//	task, exist, err := s.GetTaskDetailById(taskId)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	if !exist {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(
//			errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))))
//	}
//	if task.Instance == nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(
//			errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))))
//	}
//
//	// if instance is not connectable, exec sql must be failed;
//	// commit action unable to retry, so don't to exec it.
//	if err := executor.Ping(log.NewEntry(), task.Instance); err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//
//	sqledServer := server.GetSqled()
//	taskRes := &model.Task{}
//	if isSync == "true" {
//		taskRes, err = sqledServer.AddTaskWaitResult(taskId, model.TASK_ACTION_COMMIT)
//	} else {
//		err = sqledServer.AddTask(taskId, model.TASK_ACTION_COMMIT)
//	}
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//
//	return c.JSON(http.StatusOK, &CommitTaskRes{
//		BaseRes: controller.NewBaseReq(nil),
//		Data:    CommitTaskResult{TaskExecStatus: taskRes.ExecStatus},
//	})
//}

//// @Summary Sql提交回滚
//// @Description rollback sql
//// @Param task_id path string true "Task ID"
//// @Success 200 {object} controller.BaseRes
//// @router /tasks/{task_id}/rollback [post]
//func RollbackTask(c echo.Context) error {
//	s := model.GetStorage()
//	taskId := c.Param("task_id")
//	task, exist, err := s.GetTaskDetailById(taskId)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	if !exist {
//		return c.JSON(http.StatusOK, TASK_NOT_EXIST)
//	}
//	if task.Instance == nil {
//		return c.JSON(http.StatusOK, INSTANCE_NOT_EXIST_ERROR)
//	}
//	err = server.GetSqled().AddTask(taskId, model.TASK_ACTION_ROLLBACK)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
//}

//// @Summary 驳回工单
//// @Description reject task
//// @Param task_id path string true "Task ID"
//// @Success 200 {object} controller.BaseRes
//// @router /tasks/{task_id}/reject [post]
//func RejectTask(c echo.Context) error {
//	s := model.GetStorage()
//	taskId := c.Param("task_id")
//	task, exist, err := s.GetTaskDetailById(taskId)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	if !exist {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(
//			errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))))
//	}
//	if task.TaskStatus == model.TaskStatusExecuting ||
//		task.TaskStatus == model.TaskStatusSucceeded ||
//		task.TaskStatus == model.TaskStatusFailed {
//		return c.JSON(http.StatusOK, errors.New(errors.TASK_ACTION_INVALID,
//			fmt.Errorf("task has been executing")))
//	} else if task.TaskStatus == model.TaskStatusReject {
//		return c.JSON(http.StatusOK, errors.New(errors.TASK_ACTION_INVALID,
//			fmt.Errorf("task has been rejected")))
//	}
//	err = s.UpdateTask(task, map[string]interface{}{"task_status": model.TaskStatusReject})
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
//}

type GetTaskSQLsReqV1 struct {
	FilterExecStatus  string `json:"filter_exec_status" query:"filter_exec_status"`
	FilterAuditStatus string `json:"filter_audit_status" query:"filter_audit_status"`
	NoDuplicate       bool   `json:"no_duplicate" query:"no_duplicate"`
	PageIndex         uint32 `json:"page_index" query:"page_index" valid:"required,int"`
	PageSize          uint32 `json:"page_size" query:"page_size" valid:"required,int"`
}

type GetTaskSQLsResV1 struct {
	controller.BaseRes
	Data      []*TaskSQLResV1 `json:"data"`
	TotalNums uint64          `json:"total_nums"`
}

type TaskSQLResV1 struct {
	Number      uint   `json:"number"`
	ExecSQL     string `json:"exec_sql"`
	AuditResult string `json:"audit_result"`
	AuditLevel  string `json:"audit_level"`
	AuditStatus string `json:"audit_status"`
	ExecResult  string `json:"exec_result"`
	ExecStatus  string `json:"exec_status"`
	RollbackSQL string `json:"rollback_sql"`
}

// @Summary 获取指定task的SQLs信息
// @Description get information of all SQLs belong to the specified task
// @Tags task
// @Id getTaskSQLsV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed)
// @Param filter_audit_status query string false "filter: audit status of task sql" Enums(initialized,doing,finished)
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Param page_index query string false "page index"
// @Param page_size query string false "page size"
// @Success 200 {object} v1.GetTaskSQLsResV1
// @router /v1/tasks/{task_id}/sqls [get]
func GetTaskSQLs(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	_, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, controller.NewBaseReq(
			errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))))
	}

	req := new(GetTaskSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
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

	taskSQLsRes := make([]*TaskSQLResV1, 0, len(taskSQLs))
	for _, taskSQL := range taskSQLs {
		taskSQLRes := &TaskSQLResV1{
			Number:      taskSQL.Number,
			ExecSQL:     taskSQL.ExecSQL,
			AuditResult: taskSQL.AuditResult,
			AuditLevel:  taskSQL.AuditLevel,
			AuditStatus: taskSQL.AuditStatus,
			ExecResult:  taskSQL.ExecResult,
			ExecStatus:  taskSQL.ExecStatus,
			RollbackSQL: taskSQL.RollbackSQL,
		}
		taskSQLsRes = append(taskSQLsRes, taskSQLRes)
	}
	return c.JSON(http.StatusOK, &GetTaskSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      taskSQLsRes,
		TotalNums: count,
	})
}

type DownloadTaskSQLsFileReqV1 struct {
	NoDuplicate string `json:"no_duplicate" query:"no_duplicate"`
}

// @Summary 下载指定task的SQLs信息报告
// @Description download report file of all SQLs information belong to the specified task
// @Tags task
// @Id downloadTaskSQLReportV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Success 200 file 1 "sql report csv file"
// @router /v1/tasks/{task_id}/sql_report [get]
func DownloadTaskSQLReportFile(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, controller.NewBaseReq(
			errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))))
	}
	req := new(DownloadTaskSQLsFileReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
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
			td.RollbackSQL,
		})
	}
	cw.Flush()
	fileName := fmt.Sprintf("SQL审核报告_%v_%v.csv", task.Instance.Name, taskId)
	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": fileName}))
	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}

// @Summary 下载指定task的SQL文件
// @Description download SQL file for the task
// @Tags task
// @Id downloadTaskSQLFileV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Success 200 200 file 1 "sql file"
// @router /v1/tasks/{task_id}/sql_file [get]
func DownloadTaskSQLFile(c echo.Context) error {
	taskId := c.Param("task_id")
	s := model.GetStorage()
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, controller.NewBaseReq(
			errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))))
	}
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	sqls, err := s.GetTaskExecuteSQLContent(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
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

//type GetExecErrUploadedSqlsRes struct {
//	controller.BaseRes
//	ExecErrCommitSqls []model.ExecuteSQL `json:"exec_error_commit_sql_list"`
//}
//
//// @Summary 获取指定task 执行异常的SQLs信息
//// @Description get information of execute error SQLs belong to the specified task
//// @Param task_id path string true "task id"
//// @Success 200 {object} controller.GetExecErrUploadedSqlsRes
//// @router /tasks/{task_id}/execute_error_uploaded_sqls [get]
//func GetExecErrUploadedSqls(c echo.Context) error {
//	s := model.GetStorage()
//	taskId := c.Param("task_id")
//	execErrCommitSqls, err := s.GetExecErrorCommitSqlsByTaskId(taskId)
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	return c.JSON(http.StatusOK, &GetExecErrUploadedSqlsRes{
//		BaseRes:           controller.NewBaseReq(nil),
//		ExecErrCommitSqls: execErrCommitSqls,
//	})
//}
//
//type GetTaskListTipsRes struct {
//	controller.BaseRes
//	Data GetTaskListTipsData `json:"data"`
//}
//
//type GetTaskListTipsData struct {
//	InstancesName []string `json:"instance_name_list"`
//}
//
//// @Summary 获取sql审核工单的提示信息
//// @Description show task tips for task list
//// @Success 200 {object} controller.GetTaskListTipsRes
//// @router /task_tips [get]
//func GetTaskListTips(c echo.Context) error {
//	s := model.GetStorage()
//	names, err := s.GetTasksInstanceName()
//	if err != nil {
//		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
//	}
//	return c.JSON(http.StatusOK, &GetTaskListTipsRes{
//		BaseRes: controller.NewBaseReq(nil),
//		Data: GetTaskListTipsData{
//			InstancesName: names,
//		},
//	})
//}
//
//func FormatStringToInt(s string) (ret int, err error) {
//	if s == "" {
//		return 0, nil
//	} else {
//		ret, err = strconv.Atoi(s)
//		if err != nil {
//			return 0, err
//		}
//	}
//	return ret, nil
//}
//
//func FormatStringToBoolean(s string) (ret bool) {
//	switch s {
//	case "TRUE", "true", "True", "1", "ON", "on", "On":
//		return true
//	default:
//		return false
//	}
//}
