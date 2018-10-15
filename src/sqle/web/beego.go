package web

import (
	"fmt"
	"github.com/astaxie/beego"
	"math/rand"
	"sqle/web/controller"
)

func init() {
	beego.Router("/login", &controller.LoginController{}, "POST:Login")
	beego.Router("/user", &controller.BaseController{}, "POST:AddUser")
	beego.Router("/users", &controller.BaseController{}, "GET:UserList")
	beego.Router("/database", &controller.BaseController{}, "POST:AddDatabase")
	beego.Router("/database/:dbId/schemas", &controller.BaseController{}, "GET:GetDatabaseSchemas")
	beego.Router("/database/ping", &controller.BaseController{}, "GET:PingDatabase")
	beego.Router("/databases", &controller.BaseController{}, "GET:DatabaseList")
	beego.Router("/task", &controller.BaseController{}, "POST:AddTask")
	beego.Router("/tasks", &controller.BaseController{}, "GET:TaskList")
	beego.Router("/task/:taskId/inspect", &controller.BaseController{}, "POST:Inspect")
	beego.Router("/task/:taskId/commit", &controller.BaseController{}, "POST:Commit")
}

func StartBeego(port int, beegoExitChan chan struct{}) {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = fmt.Sprintf("sqle-%v-%x", port, rand.Int())
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 0
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 1800

	beego.Run(fmt.Sprintf(":%v", port))
	close(beegoExitChan)
}
