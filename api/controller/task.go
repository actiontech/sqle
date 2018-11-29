package controller

import (
	"github.com/labstack/echo"
	"net/http"
	"sqle/api/server"
	"sqle/inspector"
	"sqle/model"
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
	s := model.GetStorage()
	req := new(CreateTaskReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	inst, exist, err := s.GetInstByName(req.InstName)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, INSTANCE_NOT_EXIST_ERROR)
	}

	task := &model.Task{
		Name:       req.Name,
		Desc:       req.Desc,
		Schema:     req.Schema,
		InstanceId: inst.ID,
		Instance:   inst,
		CommitSqls: []*model.CommitSql{},
	}
	sqlArray, err := inspector.NewInspector(task).SplitSql(req.Sql)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	for n, sql := range sqlArray {
		task.CommitSqls = append(task.CommitSqls, &model.CommitSql{
			Sql: model.Sql{
				Number:  uint(n + 1),
				Content: sql,
			},
		})
	}
	task.Instance = nil
	err = s.Save(task)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	return c.JSON(200, &GetTaskRes{
		BaseRes: NewBaseReq(nil),
		Data:    task.Detail(),
	})
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
	_, sql, err := readFileToByte(c, "sql_file")
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}
	sqlArray, err := inspector.NewInspector(task).SplitSql(string(sql))
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	commitSqls := make([]*model.CommitSql, 0, len(sqlArray))
	for n, sql := range sqlArray {
		commitSqls = append(commitSqls, &model.CommitSql{
			Sql: model.Sql{
				Number:  uint(n + 1),
				Content: sql,
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
	return c.JSON(http.StatusOK, &GetTaskRes{
		BaseRes: NewBaseReq(nil),
		Data:    task.Detail(),
	})
}

type GetAllTaskRes struct {
	BaseRes
	Data []model.Task `json:"data"`
}

// @Summary Sql审核列表
// @Description get all tasks
// @Success 200 {object} controller.GetAllTaskRes
// @router /tasks [get]
func GetTasks(c echo.Context) error {
	s := model.GetStorage()
	if s == nil {
		c.String(500, "nil")
	}
	tasks, err := s.GetTasks()
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

// @Summary Sql提交上线
// @Description commit sql
// @Param task_id path string true "Task ID"
// @Success 200 {object} controller.BaseRes
// @router /tasks/{task_id}/commit [post]
func CommitTask(c echo.Context) error {
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
	err = server.GetSqled().AddTask(taskId, model.TASK_ACTION_COMMIT)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, NewBaseReq(nil))
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
