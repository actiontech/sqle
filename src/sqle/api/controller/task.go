package controller

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"log"
	"sqle/storage"
)

type AddTaskReq struct {
	DbId         string `form:"db_id"`
	ApproverName string `form:"approver"`
	Sql          string `form:"sql"`
}

func (r *AddTaskReq) Valid(valid *validation.Validation) {
	valid.Required(r.DbId, "db_id").Message("不能为空")
	valid.Required(r.ApproverName, "approver").Message("不能为空")
}

func (c *BaseController) AddTask() {
	req := &AddTaskReq{}
	c.validForm(req)

	approver, err := c.storage.GetUserByName(req.ApproverName)
	if err != nil {
		fmt.Println(1)
		c.CustomAbort(500, err.Error())
	}
	database, _, err := c.storage.GetDatabaseById(req.DbId)
	if err != nil {
		fmt.Println(2)
		c.CustomAbort(500, err.Error())
	}

	task := storage.Task{
		User:     *c.currentUser,
		Db:       *database,
		Approver: *approver,
		ReqSql:   req.Sql,
	}
	err = c.storage.Save(&task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.Ctx.WriteString("ok")
	return
}

func (c *BaseController) TaskList() {
	tasks, err := c.storage.GetTasks()
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(tasks)
}

func (c *BaseController) Inspect() {
	defer func() {
		log.Println("done")
		// Println executes normally even if there is a panic
		if err := recover(); err != nil {
			log.Printf("run time panic: %v", err)
		}
	}()
	taskId := c.Ctx.Input.Param(":taskId")
	task, err := c.storage.GetTaskById(taskId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	err = c.storage.InspectTask(task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.Ctx.WriteString("ok")
	return
}

func (c *BaseController) Approve() {
	taskId := c.Ctx.Input.Param(":taskId")
	task, err := c.storage.GetTaskById(taskId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if task.Approver.ID != c.currentUser.ID {
		c.CustomAbort(500, "you are not approver")
	}
	if task.Approved {
		c.CustomAbort(500, "task has approved")
	}
	err = c.storage.ApproveTask(task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.Ctx.WriteString("ok")
	return
}

func (c *BaseController) Commit() {
	taskId := c.Ctx.Input.Param(":taskId")
	task, err := c.storage.GetTaskById(taskId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	err = c.storage.CommitTask(task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.Ctx.WriteString("ok")
	return
}

func (c *BaseController) Rollback() {
	taskId := c.Ctx.Input.Param(":taskId")
	task, err := c.storage.GetTaskById(taskId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	err = c.storage.RollbackTask(task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.Ctx.WriteString("ok")
	return
}
