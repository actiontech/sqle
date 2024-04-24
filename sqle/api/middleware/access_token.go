package middleware

import (
	"context"
	"fmt"
	"net/http"

	jwtPkg "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const AccessTokenLogin = "access_token_login"

func CheckLatestAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, exist, err := GetTokenFromContext(c)
			if err != nil {
				return err
			}
			if !exist {
				return next(c)
			}
			uid, exist, err := GetUidFromAccessToken(token)
			if err != nil {
				return err
			}
			if !exist {
				return next(c)
			}

			userInfo, err := dmsobject.GetUser(context.TODO(), uid, controller.GetDMSServerAddress())
			if err != nil {
				return err
			}
			if userInfo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "access token: cannot get user info")
			}

			if userInfo.AccessTokenInfo.AccessToken != token.Raw {
				return echo.NewHTTPError(http.StatusUnauthorized, "access token is not latest")
			}

			return next(c)
		}
	}
}

func GetTokenFromContext(c echo.Context) (token *jwt.Token, exist bool, err error) {
	user := c.Get("user")
	if user == nil {
		return nil, false, nil
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return nil, true, echo.NewHTTPError(http.StatusBadRequest, "failed to convert user from jwt token")
	}

	return token, true, nil
}

func GetUidFromAccessToken(token *jwt.Token) (uid string, exist bool, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", true, echo.NewHTTPError(http.StatusBadRequest, "failed to convert token claims to jwt")
	}

	// 如果不存在JWTLoginType字段，代表是账号密码登录获取的token或者是扫描任务的凭证，不进行校验
	loginType, ok := claims[jwtPkg.JWTLoginType]
	if !ok {
		return "", false, nil
	}
	if loginType != AccessTokenLogin {
		return "", true, echo.NewHTTPError(http.StatusUnauthorized, "access token login type is error")
	}
	uid = fmt.Sprintf("%v", claims[jwtPkg.JWTUserId])
	return uid, true, nil
}
