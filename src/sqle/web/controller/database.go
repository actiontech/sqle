package controller

import (
	"sqle/executor"
	"sqle/storage"
)

type DbController struct {
	BaseController
}

type DatabaseReq struct {
	DbType   int    `form:"db_type"`
	User     string `form:"user"`
	Host     string `form:"host"`
	Port     string `form:"port"`
	Password string `form:"password"`
	Alias    string `form:"alias"`
}

func (c *BaseController) AddDatabase() {
	req := &DatabaseReq{}
	c.validForm(req)
	database := &storage.Db{
		DbType:   req.DbType,
		User:     req.User,
		Host:     req.Host,
		Port:     req.Port,
		Password: req.Password,
		Alias:    req.Alias,
	}
	exist, err := c.storage.Exist(database)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if exist {
		c.CustomAbort(500, "user exist")
	}

	err = c.storage.Save(database)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.Ctx.WriteString("ok")
	return
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
	c.Ctx.WriteString("ok")
	return
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
