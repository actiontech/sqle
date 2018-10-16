package controller

import (
	"github.com/astaxie/beego/validation"
	"sqle/storage"
)

type LoginController struct {
	InitController
}

type LoginReq struct {
	Name     string `form:"user"`
	Password string `form:"password"`
}

func (r *LoginReq) Valid(valid *validation.Validation) {
	valid.Required(r.Name, "name").Message("不能为空")
	valid.Required(r.Password, "password").Message("不能为空")
}

func (c *LoginController) Login() {
	r := &LoginReq{}
	c.validForm(r)

	var user = &storage.User{Name: r.Name}
	ok, err := c.storage.Exist(user)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if !ok {
		c.CustomAbort(500, "user not exist")
	}
	user, err = c.storage.GetUserByName(r.Name)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	// TODO: Password needs to be encrypted
	if user.Password != r.Password {
		c.CustomAbort(500, "password is invalid")
	}
	c.SetSession("user", user)
	c.CustomAbort(200, "")
}

func (c *LoginController) UnLogin() {
	c.DelSession("user")
	c.Ctx.Redirect(302, "/login")
}
