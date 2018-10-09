package web

import (
	"fmt"
	"github.com/astaxie/beego"
)

func StartBeego(port int, beegoExitChan chan struct{}) {
	addRouter()
	beego.Run(fmt.Sprintf(":%v", port))
	close(beegoExitChan)
}

func addRouter() {
	beego.Router("/user", &BaseController{}, "POST:AddUser")
	beego.Router("/users", &BaseController{}, "GET:UserList")
	beego.Router("/database", &BaseController{}, "POST:AddDatabase")
	beego.Router("/database/:dbId/schemas", &BaseController{}, "GET:GetDatabaseSchemas")
	beego.Router("/database/ping", &BaseController{}, "GET:PingDatabase")
	beego.Router("/databases", &BaseController{}, "GET:DatabaseList")
	beego.Router("/task", &BaseController{}, "POST:AddTask")
	beego.Router("/tasks", &BaseController{}, "GET:TaskList")
	beego.Router("/task/:taskId/inspect", &BaseController{}, "POST:Inspect")
	beego.Router("/task/:taskId/commit", &BaseController{}, "POST:Commit")
}
