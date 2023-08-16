package v2

// import (
// 	"net/http"

// 	"github.com/actiontech/sqle/sqle/api/controller"
// 	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
// 	"github.com/labstack/echo/v4"
// )

// type UserLoginReqV2 struct {
// 	UserName string `json:"username" form:"username" example:"test" valid:"required"`
// 	Password string `json:"password" form:"password" example:"123456" valid:"required"`
// }

// // @Summary 用户登录
// // @Description user login
// // @Tags user
// // @Id loginV2
// // @Param user body v1.UserLoginReqV1 true "user login request"
// // @Success 200 {object} controller.BaseRes
// // @router /v2/login [post]
// func LoginV2(c echo.Context) error {
// 	req := new(UserLoginReqV2)
// 	if err := controller.BindAndValidateReq(c, req); err != nil {
// 		return err
// 	}

// 	_, err := v1.Login(c, req.UserName, req.Password)
// 	if err != nil {
// 		return controller.JSONBaseErrorReq(c, err)
// 	}

// 	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
// }
