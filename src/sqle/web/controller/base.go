package controller

import (
	"fmt"
	"github.com/astaxie/beego"
	"sqle"
	"sqle/storage"
)

type InitController struct {
	beego.Controller
	storage *storage.Storage
}

func (c *InitController) Prepare() {
	c.storage = sqle.GetSqled().Storage
}

func (c *InitController) serveJson(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}

func (c *InitController) redirectAndAbort(status int, localUrl string) {
	c.Ctx.Redirect(status, localUrl)
	panic(beego.ErrAbort)
}

// BaseController is a Controller for user whom has logged in.
type BaseController struct {
	InitController
	currentUser *storage.User
}

func (c *BaseController) Prepare() {
	c.InitController.Prepare()

	// load user
	user := c.Ctx.Input.Session("user")
	if user == nil {
		fmt.Println(1)
		c.redirectAndAbort(302, "/login")
	}
	u, ok := user.(*storage.User)
	if !ok {
		fmt.Println(2)
		c.redirectAndAbort(302, "/login")
	}

	exist, err := c.storage.Exist(u)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if !exist {
		fmt.Println(3)
		c.redirectAndAbort(302, "/login")
	}

	if user, err := c.storage.GetUserByName(u.Name); err != nil {
		c.CustomAbort(500, err.Error())
	} else {
		c.currentUser = user
	}
}
