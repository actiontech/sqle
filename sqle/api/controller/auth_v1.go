package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const JWTSecret = "secret"

type UserLoginReq struct {
	UserName string `json:"username" form:"username" example:"test" valid:"required"`
	Password string `json:"password" form:"password" example:"123456" valid:"required"`
}

type UserLoginRes struct {
	Token string `json:"token" example:"this is a jwt token string"`
}

// @Summary 用户登录
// @Description user login
// @Param username formData string true "user name"
// @param password formData string true "user password"
// @Success 200 {object} UserLoginRes
// @router /v1/login [post]
func Login(c echo.Context) error {
	req := new(UserLoginReq)
	if err := BindAndValidateReq(c, req); err != nil {
		return err
	}
	if req.UserName == "admin" && req.Password == "admin" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "admin"
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		t, err := token.SignedString([]byte(JWTSecret))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return c.JSON(http.StatusOK, &UserLoginRes{
			Token: t,
		})
	}
	return echo.ErrUnauthorized
}
