package controller

import (
	"sqle/executor"
	"sqle/storage"
	"fmt"
)

func (c *BaseController) AddUser() {
	user := &storage.User{}
	user.Name = c.GetString("name")
	exist, err := c.storage.Exist(user)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if exist {
		c.CustomAbort(500, "user exist")
	}

	user.Password = c.GetString("password")
	err = c.storage.Create(user)
	if nil != err {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) UserList() {
	users := []*storage.User{}
	users, err := c.storage.GetUsers()
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
	err := c.storage.Save(database)
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

func (c *BaseController) GetDatabaseSchemas() {
	dbId := c.Ctx.Input.Param(":dbId")
	db, err := c.storage.GetDatabaseById(dbId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	schemas, err := executor.ShowDatabase(db)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(schemas)
}

func (c *BaseController) DatabaseList() {
	databases, err := c.storage.GetDatabases()
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

	user, err := c.storage.GetUserById(userId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	approver, err := c.storage.GetUserById(approverId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	database, err := c.storage.GetDatabaseById(dbId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}

	task.User = *user
	task.Approver = *approver
	task.Db = *database
	task.ReqSql = c.GetString("sql")
	err = c.storage.Save(&task)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) TaskList() {
	fmt.Println("tasks")
	tasks, err := c.storage.GetTasks()
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(tasks)
}

func (c *BaseController) Inspect() {
	taskId := c.Ctx.Input.Param(":taskId")

	err := c.storage.InspectTask(taskId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) Commit() {
	taskId := c.Ctx.Input.Param(":taskId")
	err := c.storage.CommitTask(taskId)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}
