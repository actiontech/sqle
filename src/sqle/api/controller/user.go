package controller

//import "sqle/model"
//
//type UserController struct {
//	BaseController
//}
//
//type UserReq struct {
//	Name     string `form:"user"`
//	Password string `form:"password"`
//}
//
//func (c *BaseController) AddUser() {
//	req := &UserReq{}
//	c.validForm(req)
//
//	user := &model.User{
//		Name: req.Name,
//	}
//	exist, err := c.model.Exist(user)
//	if err != nil {
//		c.CustomAbort(500, err.Error())
//	}
//	if exist {
//		c.CustomAbort(500, "user exist")
//	}
//
//	user.Password = req.Password
//	err = c.model.Create(user)
//	if nil != err {
//		c.CustomAbort(500, err.Error())
//	}
//	c.Ctx.WriteString("ok")
//	return
//}
//
//func (c *BaseController) UserList() {
//	users := []*model.User{}
//	users, err := c.model.GetUsers()
//	if err != nil {
//		c.CustomAbort(500, err.Error())
//	}
//	c.serveJson(users)
//}
