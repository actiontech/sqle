package web

import (
	"github.com/astaxie/beego"
	"sqle"
	"sqle/storage"
)

type BaseController struct {
	beego.Controller
	sqled *sqle.Sqled
}

func (c *BaseController) Prepare() {
	c.sqled = sqle.GetSqled()
}

func (c *BaseController) serveJson(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}

func (c *BaseController) Test() {
	c.CustomAbort(200, "ok")
}

func (c *BaseController) AddUser() {
	user := &storage.User{}
	user.Name = c.GetString("name")
	exist, err := c.sqled.Exist(user)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if exist {
		c.CustomAbort(500, "user exist")
	}

	user.Password = c.GetString("password")
	err = c.sqled.Db.Create(user).Error
	if nil != err {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) UserList() {
	users := []*storage.User{}
	err := c.sqled.Db.Find(&users).Error
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
	err := c.sqled.Db.Save(database).Error
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200, "ok")
}

func (c *BaseController) DatabaseList() {
	databases := []*storage.Db{}
	err := c.sqled.Db.Find(&databases).Error
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(databases)
}

func (c *BaseController) AddTask() {
	task := storage.Task{}
	user := storage.User{}
	p := storage.User{}
	d := storage.Db{}
	userId, _ := c.GetInt("user_id")
	dbId, _ := c.GetInt("db_id")
	approverId,_:=c.GetInt("approver_id")

	err := c.sqled.Db.First(&user, userId).Error
	if err != nil {
		c.CustomAbort(500, err.Error())
	}

	err = c.sqled.Db.First(&d, dbId).Error
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	err = c.sqled.Db.First(&p, approverId).Error
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	task.User = user
	task.Approver = p
	task.Db = d
	task.ReqSql = "select * from sqle.dbs"
	err = c.sqled.Db.Save(&task).Error
	if err!=nil{
		c.CustomAbort(500, err.Error())
	}
	c.CustomAbort(200,"ok")
}

func (c *BaseController) TaskList() {
	tasks := []*storage.Task{}
	err := c.sqled.Db.Preload("User").Preload("Approver").Preload("Db").Find(&tasks).Error
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(tasks)
}