package v2

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/labstack/echo/v4"
)

type UserLoginReqV2 struct {
	UserName          string `json:"username" form:"username" example:"test"`
	Password          string `json:"password" form:"password" example:"123456"`
	EncryptedUsername string `json:"encrypted_username" form:"encrypted_username"`
	EncryptedPassword string `json:"encrypted_password" form:"encrypted_password"`
	KeyID             string `json:"key_id" form:"key_id"`
}

// @Summary 用户登录
// @Description user login
// @Tags user
// @Id loginV2
// @Param user body v2.UserLoginReqV2 true "user login request"
// @Success 200 {object} controller.BaseRes
// @router /v2/login [post]
func LoginV2(c echo.Context) error {
	req := new(UserLoginReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	userName, err := v1.ResolveLoginUsername(req.UserName, req.EncryptedUsername, req.KeyID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	password, err := v1.ResolveLoginPassword(req.Password, req.EncryptedPassword, req.KeyID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	_, err = v1.Login(c, userName, password)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}
