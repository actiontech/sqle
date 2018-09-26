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
	beego.Router("/test", &BaseController{}, "GET:Test")
	beego.Router("/user", &BaseController{}, "GET:AddUser")
	beego.Router("/users", &BaseController{}, "GET:UserList")
	beego.Router("/database", &BaseController{}, "GET:AddDatabase")
	beego.Router("/databases", &BaseController{}, "GET:DatabaseList")
	beego.Router("/task", &BaseController{}, "GET:AddTask")
	beego.Router("/tasks", &BaseController{}, "GET:TaskList")
}
