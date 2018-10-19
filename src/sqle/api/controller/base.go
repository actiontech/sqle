package controller

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"sqle/storage"
	"strings"
)

// ValidFormer valid interface
type ValidFormer interface {
	Valid(valid *validation.Validation)
}

type InitController struct {
	beego.Controller
	storage *storage.Storage
}

func (c *InitController) Prepare() {
	//c.storage = sqle.GetSqled().Storage
}

func (c *InitController) serveJson(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}

func (c *InitController) redirectAndAbort(status int, localUrl string) {
	c.Ctx.Redirect(status, localUrl)
	panic(beego.ErrAbort)
}

func (c *InitController) validForm(obj interface{}) {
	if err := c.ParseForm(obj); err != nil {
		c.CustomAbort(500, err.Error())
	}
	form, ok := obj.(ValidFormer)
	if !ok {
		return
	}
	valid := &validation.Validation{}
	form.Valid(valid)
	if valid.HasErrors() {
		msgs := []string{}
		for _, err := range valid.Errors {
			msgs = append(msgs, fmt.Sprintf("%s:%s", err.Key, err.Message))
		}
		c.CustomAbort(500, strings.Join(msgs, ", "))
	}
}

// BaseController is a Controller for user whom has logged in.
type BaseController struct {
	InitController
	currentUser *storage.User
}

func (c *BaseController) Prepare() {
	c.InitController.Prepare()

	//// auth and load user
	//user := c.Ctx.Input.Session("user")
	//if user == nil {
	//	c.redirectAndAbort(302, "/login")
	//}
	//u, ok := user.(*storage.User)
	//if !ok {
	//	c.redirectAndAbort(302, "/login")
	//}
	//
	//exist, err := c.storage.Exist(u)
	//if err != nil {
	//	c.CustomAbort(500, err.Error())
	//}
	//if !exist {
	//	c.redirectAndAbort(302, "/login")
	//}
	//
	//if user, err := c.storage.GetUserByName(u.Name); err != nil {
	//	c.CustomAbort(500, err.Error())
	//} else {
	//	c.currentUser = user
	//}
}

type BaseReq struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewBaseReq(code int, message string) BaseReq {
	return BaseReq{
		Code:    code,
		Message: message,
	}
}
