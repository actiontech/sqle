package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

type UserLoginReqV1 struct {
	UserName string `json:"username" form:"username" example:"test" valid:"required"`
	Password string `json:"password" form:"password" example:"123456" valid:"required"`
}

type GetUserLoginResV1 struct {
	controller.BaseRes
	Data UserLoginResV1 `json:"data"`
}

type UserLoginResV1 struct {
	Token string `json:"token" example:"this is a jwt token string"`
}

// @Summary 用户登录
// @Description user login
// @Tags user
// @Id loginV1
// @Param user body v1.UserLoginReqV1 true "user login request"
// @Success 200 {object} v1.GetUserLoginResV1
// @router /v1/login [post]
func Login(c echo.Context) error {
	req := new(UserLoginReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	user, exist, err := s.GetUserByName(req.UserName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist || !(req.UserName == user.Name && req.Password == user.Password) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.LoginAuthFail,
			fmt.Errorf("password is wrong or user does not exist")))
	}

	j := utils.NewJWT([]byte(utils.JWTSecret))
	t, err := j.CreateToken(req.UserName, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, &GetUserLoginResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: UserLoginResV1{
			Token: t,
		},
	})
}
