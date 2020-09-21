package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"actiontech.cloud/universe/sqle/v4/sqle/executor"

	"actiontech.cloud/universe/sqle/v4/sqle/api/server"
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"actiontech.cloud/universe/sqle/v4/sqle/inspector"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/sqle/v4/sqle/model"

	"github.com/labstack/echo/v4"
)

type CreateTaskReq struct {
	Name     string `json:"name" example:"test" valid:"required"`
	Desc     string `json:"desc" example:"this is a test task" valid:"-"`
	InstName string `json:"inst_name" form:"inst_name" example:"inst_1" valid:"required"`
	Schema   string `json:"schema" example:"db1" valid:"-"`
	Sql      string `json:"sql" example:"alter table tb1 drop columns c1" valid:"-"`
}

type GetTaskRes struct {
	BaseRes
	Data model.TaskDetail `json:"data"`
}

// @Summary 创建Sql审核任务
// @Description create a task
// @Accept json
// @Produce json
// @Param instance body controller.CreateTaskReq true "add task"
// @Success 200 {object} controller.GetTaskRes
// @router /tasks [post]
func CreateTask(c echo.Context) error {
	task, res := createTask(c)
	if res.Code != 0 {
		return c.JSON(200, res)
	}
	return c.JSON(200, &GetTaskRes{
		BaseRes: res,
		Data:    task.Detail(),
	})
}

func createTask(c echo.Context) (*model.Task, BaseRes) {
	req := new(CreateTaskReq)
	if err := c.Bind(req); err != nil {
		return nil, NewBaseReq(err)
	}
	if err := c.Validate(req); err != nil {
		return nil, NewBaseReq(err)
	}

	params := []*string{&req.Name, &req.Desc, &req.InstName, &req.Schema, &req.Sql}
	if err := unescapeParamString(params); nil != err {
		return nil, NewBaseReq(err)
	}

	return createTaskByRequestParam(req)
}

func createTaskByRequestParam(req *CreateTaskReq) (*model.Task, BaseRes) {
	s := model.GetStorage()
	inst, exist, err := s.GetInstByName(req.InstName)
	if err != nil {
		return nil, NewBaseReq(err)
	}
	if !exist {
		return nil, INSTANCE_NOT_EXIST_ERROR
	}
	if err := executor.Ping(log.NewEntry(), inst); err != nil {
		return nil, NewBaseReq(err)
	}

	task := &model.Task{
		Name:         req.Name,
		Desc:         req.Desc,
		Schema:       req.Schema,
		InstanceId:   inst.ID,
		Instance:     inst,
		CommitSqls:   []*model.CommitSql{},
		InstanceName: req.InstName,
	}

	nodes, err := inspector.NewInspector(log.NewEntry(), inspector.NewContext(nil), task, nil, nil).
		ParseSql(req.Sql)
	if err != nil {
		return nil, NewBaseReq(err)
	}
	for n, node := range nodes {
		task.CommitSqls = append(task.CommitSqls, &model.CommitSql{
			Sql: model.Sql{
				Number:  uint(n + 1),
				Content: node.Text(),
			},
		})
	}
	task.Instance = nil
	err = s.Save(task)
	if err != nil {
		return nil, NewBaseReq(err)
	}
	return task, NewBaseReq(nil)
}

// @Summary 上传 SQL 文件
// @Description upload SQL file
// @Accept mpfd
// @Param task_id path string true "Task ID"
// @Param sql_file formData file true "SQL file"
// @Success 200 {object} controller.BaseRes
// @router /tasks/{task_id}/upload_sql_file [post]
func UploadSqlFile(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, TASK_NOT_EXIST)
	}
	if task.HasDoingCommit() {
		return c.JSON(http.StatusOK, errors.New(errors.TASK_ACTION_INVALID,
			fmt.Errorf("task has commit, not allow update sql")))
	}
	_, sql, err := readFileToByte(c, "sql_file")
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}
	nodes, err := inspector.NewInspector(log.NewEntry(), inspector.NewContext(nil), task, nil, nil).
		ParseSql(string(sql))
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	commitSqls := make([]*model.CommitSql, 0, len(nodes))
	for n, node := range nodes {
		commitSqls = append(commitSqls, &model.CommitSql{
			Sql: model.Sql{
				Number:  uint(n + 1),
				Content: node.Text(),
			},
		})
	}
	s.UpdateCommitSql(task, commitSqls)
	return c.JSON(http.StatusOK, NewBaseReq(nil))
}

// @Summary 获取Sql审核任务信息
// @Description get task
// @Param task_id path string true "Task ID"
// @Success 200 {object} controller.GetTaskRes
// @router /tasks/{task_id}/ [get]
func GetTask(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, TASK_NOT_EXIST)
	}
	return c.JSON(http.StatusOK, &GetTaskRes{
		BaseRes: NewBaseReq(nil),
		Data:    task.Detail(),
	})
}

// @Summary 删除审核任务
// @Description delete task
// @Param task_id path string true "Task ID"
// @Success 200 {object} controller.BaseRes
// @router /tasks/{task_id}/ [delete]
func DeleteTask(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, TASK_NOT_EXIST)
	}

	// must check task not running

	err = s.Delete(task)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, NewBaseReq(nil))
}

// @Summary 批量删除审核任务
// @Description delete tasks by ids
// @Accept x-www-form-urlencoded
// @Param task_ids formData string true "remove tasks by ids(interlaced by ',')"
// @Success 200 {object} controller.BaseRes
// @router /tasks/remove_by_task_ids [post]
func DeleteTasks(c echo.Context) error {
	s := model.GetStorage()
	taskIds, err := url.QueryUnescape(c.FormValue("task_ids"))
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	deleteTaskIds := strings.Split(strings.TrimRight(taskIds, ","), ",")

	err = s.HardDeleteTasksByIds(deleteTaskIds)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	err = s.HardDeleteRollbackSqlByTaskIds(deleteTaskIds)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	err = s.HardDeleteSqlCommittingResultByTaskIds(deleteTaskIds)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, NewBaseReq(nil))
}

type GetAllTaskRes struct {
	BaseRes
	Data []model.Task `json:"data"`
}

// @Summary Sql审核列表
// @Description get all tasks
// @Param task_ids query string false "get task by ids(interlaced by ',')"
// @Success 200 {object} controller.GetAllTaskRes
// @router /tasks [get]
func GetTasks(c echo.Context) error {
	s := model.GetStorage()
	if s == nil {
		c.String(500, "nil")
	}
	var err error
	var tasks []model.Task

	taskIds, err := url.QueryUnescape(c.QueryParam("task_ids"))
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if taskIds == "" {
		tasks, err = s.GetTasks()
	} else {
		tasks, err = s.GetTasksByIds(strings.Split(strings.TrimRight(taskIds, ","), ","))
	}
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetAllTaskRes{
		BaseRes: NewBaseReq(nil),
		Data:    tasks,
	})
}

// @Summary Sql提交审核
// @Description inspect sql
// @Param task_id path string true "Task ID"
// @Success 200 {object} controller.GetTaskRes
// @router /tasks/{task_id}/inspection [post]
func InspectTask(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, TASK_NOT_EXIST)
	}
	if task.Instance == nil {
		return c.JSON(http.StatusOK, INSTANCE_NOT_EXIST_ERROR)
	}
	task, err = server.GetSqled().AddTaskWaitResult(taskId, model.TASK_ACTION_INSPECT)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetTaskRes{
		BaseRes: NewBaseReq(nil),
		Data:    task.Detail(),
	})
}

type CommitTaskRes struct {
	BaseRes
	Data CommitTaskResult `json:"data"`
}

type CommitTaskResult struct {
	TaskExecStatus string `json:"task_exec_status"`
}

// @Summary Sql提交上线
// @Description commit sql
// @Accept x-www-form-urlencoded
// @Param task_id path string true "Task ID"
// @Param is_sync formData boolean false "the request is sync or async."
// @Success 200 {object} controller.CommitTaskRes
// @router /tasks/{task_id}/commit [post]
func CommitTask(c echo.Context) error {
	s := model.GetStorage()
	isSync := c.FormValue("is_sync")
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, TASK_NOT_EXIST)
	}
	if task.Instance == nil {
		return c.JSON(http.StatusOK, INSTANCE_NOT_EXIST_ERROR)
	}
	sqledServer := server.GetSqled()
	taskRes := &model.Task{}
	if isSync == "true" {
		taskRes, err = sqledServer.AddTaskWaitResult(taskId, model.TASK_ACTION_COMMIT)
	} else {
		err = sqledServer.AddTask(taskId, model.TASK_ACTION_COMMIT)
	}
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	return c.JSON(http.StatusOK, &CommitTaskRes{
		BaseRes: NewBaseReq(nil),
		Data:    CommitTaskResult{TaskExecStatus: taskRes.ExecStatus},
	})
}

// @Summary Sql提交回滚
// @Description rollback sql
// @Param task_id path string true "Task ID"
// @Success 200 {object} controller.BaseRes
// @router /tasks/{task_id}/rollback [post]
func RollbackTask(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, TASK_NOT_EXIST)
	}
	if task.Instance == nil {
		return c.JSON(http.StatusOK, INSTANCE_NOT_EXIST_ERROR)
	}
	err = server.GetSqled().AddTask(taskId, model.TASK_ACTION_ROLLBACK)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, NewBaseReq(nil))
}

// @Summary 创建Sql审核任务并提交审核
// @Description create a task and inspect. NOTE: it will create a task with sqls from "sql" if "sql" isn't empty
// @Accept mpfd
// @Produce json
// @Param name formData string true "task name"
// @Param desc formData string false "description of task"
// @Param inst_name formData string true "instance name"
// @Param schema formData string false "schema of instance"
// @Param sql formData string false "sqls for audit"
// @Param uploaded_sql_file formData file false "uploaded SQL file"
// @Success 200 {object} controller.GetTaskRes
// @router /task/create_inspect [post]
func CreateAndInspectTask(c echo.Context) error {
	//check params
	req := new(CreateTaskReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	//todo @luowei unescape params using function of echo
	params := []*string{&req.Name, &req.Desc, &req.InstName, &req.Schema, &req.Sql}
	if err := unescapeParamString(params); nil != err {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	if "" == req.Sql {
		_, sqls, err := readFileToByte(c, "uploaded_sql_file")
		if err != nil {
			return c.JSON(http.StatusOK, NewBaseReq(err))
		}

		req.Sql = string(sqls)
	}

	task, res := createTaskByRequestParam(req)
	if res.Code != 0 {
		return c.JSON(200, res)
	}

	task, err := server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%d", task.ID), model.TASK_ACTION_INSPECT)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetTaskRes{
		BaseRes: NewBaseReq(nil),
		Data:    task.Detail(),
	})
}

type GetUploadedSqlsRes struct {
	BaseRes
	CommitSqls   []model.CommitSql   `json:"commit_sql_list"`
	TotalNums    uint32              `json:"total_nums"`
	RollbackSqls []model.RollbackSql `json:"rollback_sql_list"`
}

// @Summary 获取指定task的SQLs信息
// @Description get information of all SQLs belong to the specified task
// @Param task_id path string true "task id"
// @Param page_index query string false "page index"
// @Param page_size query string false "page size"
// @Param filter_sql_execution_status query string false "filter: execution status of task uploaded sql" Enums(finished, initialized, doing, failed)
// @Param filter_sql_audit_status query string false "filter: audit status of task uploaded sql" Enums(doing, finished)
// @Success 200 {object} controller.GetUploadedSqlsRes
// @router /tasks/{task_id}/uploaded_sqls [get]
func GetUploadedSqls(c echo.Context) error {
	s := model.GetStorage()
	var pageIndex, pageSize int
	taskId := c.Param("task_id")
	index, err := url.QueryUnescape(c.QueryParam("page_index"))
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	pageIndex, err = FormatStringToInt(index)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	size, err := url.QueryUnescape(c.QueryParam("page_size"))
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	pageSize, err = FormatStringToInt(size)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	filterSqlExecutionStatus, err := url.QueryUnescape(c.QueryParam("filter_sql_execution_status"))
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	filterSqlAuditStatus, err := url.QueryUnescape(c.QueryParam("filter_sql_audit_status"))
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	commitSqls, totalNums, err := s.GetUploadedSqls(taskId, filterSqlExecutionStatus, filterSqlAuditStatus, pageIndex, pageSize)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	commitSqlNum := make([]string, 0)
	for _, commitSql := range commitSqls {
		commitSqlNum = append(commitSqlNum, strconv.FormatUint(uint64(commitSql.Number), 10))
	}
	rollbackSqls, err := s.GetRollbackSqlByTaskId(taskId, commitSqlNum)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetUploadedSqlsRes{
		BaseRes:      NewBaseReq(nil),
		CommitSqls:   commitSqls,
		RollbackSqls: rollbackSqls,
		TotalNums:    totalNums,
	})
}

type GetExecErrUploadedSqlsRes struct {
	BaseRes
	ExecErrCommitSqls []model.CommitSql `json:"exec_error_commit_sql_list"`
}

// @Summary 获取指定task 执行异常的SQLs信息
// @Description get information of execute error SQLs belong to the specified task
// @Param task_id path string true "task id"
// @Success 200 {object} controller.GetExecErrUploadedSqlsRes
// @router /tasks/{task_id}/execute_error_uploaded_sqls [get]
func GetExecErrUploadedSqls(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	execErrCommitSqls, err := s.GetExecErrorCommitSqlsByTaskId(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetExecErrUploadedSqlsRes{
		BaseRes:           NewBaseReq(nil),
		ExecErrCommitSqls: execErrCommitSqls,
	})
}

func FormatStringToInt(s string) (ret int, err error) {
	if s == "" {
		return 0, nil
	} else {
		ret, err = strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
	}
	return ret, nil
}
