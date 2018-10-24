package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sqle"
	"sqle/model"
)

type CreateTaskReq struct {
	Name     string `form:"name" example:"test"`
	Desc     string `form:"desc" example:"this is a test task"`
	InstName string `json:"inst_name" form:"inst_name" example:"inst_1"`
	Schema   string `form:"schema" example:"db1"`
	Sql      string `form:"sql" example:"alter table tb1 drop columns c1"`
}

type CreateTaskRes struct {
	BaseRes
	Data model.Task `json:"data"`
}

// @Summary 创建Sql审核单
// @Description create a task
// @Accept json
// @Accept json
// @Param instance body controller.CreateTaskReq true "add task"
// @Success 200 {object} controller.CreateTaskRes
//// @router /tasks [post]
func CreateTask(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateTaskReq)
	if err := c.Bind(req); err != nil {
		return err
	}

	inst, exist, err := s.GetInstByName(req.InstName)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	if !exist {
		return c.JSON(200, NewBaseReq(-1, fmt.Sprintf("instance %s is not exist", req.InstName)))
	}

	sql := model.Sql{
		Sql: req.Sql,
	}

	task := &model.Task{
		Name:   req.Name,
		Desc:   req.Desc,
		Schema: req.Schema,
		InstId: inst.ID,
		Sql:    sql,
	}
	err = s.Save(task)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	return c.JSON(200, NewBaseReq(0, "ok"))
}

type GetAllTaskRes struct {
	BaseRes
	Data []model.Task `json:"data"`
}

// @Summary Sql审核列表
// @Description get all tasks
// @Success 200 {object} controller.GetAllTaskRes
//// @router /tasks [get]
func GetTasks(c echo.Context) error {
	s := model.GetStorage()
	if s == nil {
		c.String(500, "nil")
	}
	tasks, err := s.GetTasks()
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.JSON(http.StatusOK, &GetAllTaskRes{
		BaseRes: NewBaseReq(0, "ok"),
		Data:    tasks,
	})
}

// @Summary Sql提交审核
// @Description inspect sql
// @Param task_id path string true "Task ID"
// @Success 200 {object} controller.BaseRes
// @router /tasks/{task_id}/inspection [post]
func InspectTask(c echo.Context) error {
	s := model.GetStorage()
	taskId := c.Param("task_id")
	_, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(-1, err.Error()))
	}
	if !exist {
		return c.JSON(http.StatusOK, NewBaseReq(-1, "task not exist"))
	}
	actionChan, err := sqle.GetSqled().AddTask(taskId, model.TASK_ACTION_INSPECT)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(-1, err.Error()))
	}
	action := <-actionChan
	if action.Error != nil {
		return c.JSON(http.StatusOK, NewBaseReq(-1, action.Error.Error()))
	}
	return c.JSON(200, action.Task.Sql)
}

// @Summary Sql提交上线
// @Description commit sql
// @Success 200 {object} controller.BaseRes
// @router /tasks/{task_id}/commit [post]
func CommitTask(c echo.Context) error {
	return nil
}

// @Summary Sql提交回滚
// @Description rollback sql
// @Param manual_execute query boolean false "manual execute rollback sql"
// @Success 200 {object} controller.BaseRes
// @router /tasks/{task_id}/rollback [post]
func RollbackTask(c echo.Context) error {
	return nil
}

type GetTaskReq struct {
	BaseRes
	Data model.Task `json:"data"`
}

// @Summary 获取Sql审核单信息
// @Description get task
// @Success 200 {object} controller.GetTaskReq
// @router /tasks/{task_id} [get]
func GetTask(c echo.Context) error {
	return nil
}
