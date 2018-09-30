package web

import (
	"github.com/astaxie/beego"
	"sqle"
	"sqle/executor"
	"sqle/storage"
)

type BaseController struct {
	beego.Controller
	user  *storage.User
	sqled *sqle.Sqled
}

func (c *BaseController) Prepare() {
	c.sqled = sqle.GetSqled()
}

func (c *BaseController) serveJson(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}

func (c *BaseController) AddUser() {
	user := &storage.User{}
	user.Name = c.GetString("name")
	exist, err := c.sqled.Storage.Exist(user)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if exist {
		c.CustomAbort(500, "user exist")
	}

	user.Password = c.GetString("password")
	err = c.sqled.Storage.Create(user)
	if nil != err {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) UserList() {
	users := []*storage.User{}
	users, err := c.sqled.Storage.GetUsers()
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(users)
}

func (c *BaseController) AddDatabase() {
	database := &storage.Db{}
	database.User = c.GetString("user")
	database.DbType, _ = c.GetInt("db_type", 0)
	database.Host = c.GetString("host")
	database.Port = c.GetString("port")
	database.Password = c.GetString("password")
	database.Alias = c.GetString("alias")
	err := c.sqled.Storage.Save(database)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) PingDatabase() {
	database := &storage.Db{}
	database.User = c.GetString("user")
	database.DbType, _ = c.GetInt("db_type", 0)
	database.Host = c.GetString("host")
	database.Port = c.GetString("port")
	database.Password = c.GetString("password")
	err := executor.Ping(database)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) DatabaseList() {
	databases, err := c.sqled.Storage.GetDatabases()
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(databases)
}

func (c *BaseController) AddTask() {
	task := storage.Task{}
	userId := c.GetString("user_id")
	dbId := c.GetString("db_id")
	approverId := c.GetString("approver_id")

	user, err := c.sqled.Storage.GetUserById(userId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	approver, err := c.sqled.Storage.GetUserById(approverId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	database, err := c.sqled.Storage.GetDatabaseById(dbId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}

	task.User = *user
	task.Approver = *approver
	task.Db = *database
	task.ReqSql = c.GetString("sql")
	err = c.sqled.Storage.Save(&task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) TaskList() {
	tasks, err := c.sqled.Storage.GetTasks()
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(tasks)
}

func (c *BaseController) Inspect() {
	taskId := c.Ctx.Input.Param(":taskId")
	task, err := c.sqled.Storage.GetTaskById(taskId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	err = c.sqled.Inspect(task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) Commit() {
	taskId := c.Ctx.Input.Param(":taskId")
	task, err := c.sqled.Storage.GetTaskById(taskId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	err = c.sqled.Commit(task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}
