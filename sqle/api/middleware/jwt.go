package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTTokenAdapter is a `echo` middleware,　by rewriting the header, the jwt token support header
// "Authorization: {token}" and "Authorization: Bearer {token}".
func JWTTokenAdapter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			if auth != "" && !strings.HasPrefix(auth, middleware.DefaultJWTConfig.AuthScheme) {
				c.Request().Header.Set(echo.HeaderAuthorization,
					fmt.Sprintf("%s %s", middleware.DefaultJWTConfig.AuthScheme, auth))
			}
			// sqle-token为空时，可能是cookie过期被清理了，希望返回的错误是http.StatusUnauthorized
			// 但sqle-token为空时jwt返回的错误是http.StatusBadRequest
			_, err := c.Cookie("sqle-token")
			if err == http.ErrNoCookie && auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "can not find sqle-token")
			}

			return next(c)
		}
	}
}

func JWTWithConfig(key interface{}) echo.MiddlewareFunc {
	c := middleware.DefaultJWTConfig
	c.SigningKey = key
	c.TokenLookup = "cookie:sqle-token,header:Authorization" // tell the middleware where to get token: from cookie and header
	return middleware.JWTWithConfig(c)
}

var errAuditPlanMisMatch = errors.New("audit plan name don't match the token or audit plan not found")

// ScannerVerifier is a `echo` middleware. Every audit plan should be
// scanner-scoped which means that scanner-A should not push SQL to scanner-B.
func ScannerVerifier() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWT parser expect no 'Bearer' ahead of token, so
			// we cut the leading auth schema.
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			parts := strings.Split(auth, " ")
			token := parts[0]
			if len(parts) == 2 {
				token = parts[1]
			}

			apnInToken, userName, err := utils.ParseAuditPlanToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			projectName := c.Param("project_name")
			apnInParam := c.Param("audit_plan_name")
			// verify user
			user, isExist, err := model.GetStorage().GetUserByName(userName)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if !isExist {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("user is not exist"))
			}
			if user.IsDisabled() {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("current user is disabled"))
			}
			// verify audit plan
			// 由于对生成的JWT Token的负载使用MD5算法进行预处理，因此在验证的时候也需要对param中的apn使用MD5处理
			// 为了兼容老版本的JWT Token需要增加不经MD5处理的apnInParam和apnInToken的判断
			if apnInToken != apnInParam && apnInToken != utils.Md5(apnInParam) {
				return echo.NewHTTPError(http.StatusInternalServerError, errAuditPlanMisMatch.Error())
			}

			apn, apnExist, err := model.GetStorage().GetAuditPlanFromProjectByName(projectName, apnInParam)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			// verify token in audit plan
			if !apnExist || apn.Token != token {
				return echo.NewHTTPError(http.StatusInternalServerError, errAuditPlanMisMatch.Error())
			}

			return next(c)
		}
	}
}
