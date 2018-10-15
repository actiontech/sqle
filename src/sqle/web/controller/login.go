package controller

import (
	"sqle/storage"
)

type LoginController struct {
	InitController
}

func (c *LoginController) Login() {
	userName := c.GetString("user")
	password := c.GetString("password")

	var user = &storage.User{Name: userName}
	ok, err := c.storage.Exist(user)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	if !ok {
		c.CustomAbort(500, "user not exist")
	}
	user, err = c.storage.GetUserByName(userName)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	// TODO: Password needs to be encrypted
	if user.Password != password {
		c.CustomAbort(500, "password is invalid")
	}
	c.SetSession("user", user)
	c.CustomAbort(200, "")
}

func (c *LoginController) UnLogin() {
	c.DelSession("user")
	c.Ctx.Redirect(302, "/login")
}
