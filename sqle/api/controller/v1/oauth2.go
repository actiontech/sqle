package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

// @Summary oauth2通过此链接跳转到第三方登录网址
// @Description oauth2 link
// @Id Oauth2Link
// @Tags oauth2
// @router /v1/oauth2/link [get]
func Oauth2Link(c echo.Context) error {
	return oauth2Link(c)
}

// Oauth2Callback is a hidden interface for third-party platform callbacks for oauth2 verification
func Oauth2Callback(c echo.Context) error {
	return oauth2Callback(c)
}

type BindOauth2UserReqV1 struct {
	UserName     string `json:"user_name" from:"user_name" valid:"required"`
	Pwd          string `json:"pwd" from:"pwd" valid:"required"`
	Oauth2UserID string `json:"oauth2_user_id" from:"oauth2_user_id" valid:"required"`
}

type BindOauth2UserResV1 struct {
	controller.BaseRes
	Data BindOauth2UserResDataV1 `json:"data"`
}

type BindOauth2UserResDataV1 struct {
	Token string `json:"token"`
}

// @Summary 绑定 Oauth2 和 sqle用户
// @Description bind Oauth2 user to sqle
// @Id bindOauth2User
// @Tags oauth2
// @Param conf body v1.BindOauth2UserReqV1 true "bind oauth2 user req"
// @Success 200 {object} v1.BindOauth2UserResDataV1
// @router /v1/oauth2/user/bind [post]
func BindOauth2User(c echo.Context) error {
	return bindOauth2User(c)
}
