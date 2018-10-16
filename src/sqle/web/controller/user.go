package controller

import "sqle/storage"

type UserController struct {
	BaseController
}

type UserReq struct {
	Name     string `form:"user"`
	Password string `form:"password"`
}

func (c *BaseController) AddUser() {
	req := &UserReq{}
	c.validForm(req)

	user := &storage.User{
		Name: req.Name,
	}
	exist, err := c.storage.Exist(user)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if exist {
		c.CustomAbort(500, "user exist")
	}

	user.Password = req.Password
	err = c.storage.Create(user)
	if nil != err {
		c.CustomAbort(500, err.Error())
	}
	c.Ctx.WriteString("ok")
	return
}

func (c *BaseController) UserList() {
	users := []*storage.User{}
	users, err := c.storage.GetUsers()
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.serveJson(users)
}
