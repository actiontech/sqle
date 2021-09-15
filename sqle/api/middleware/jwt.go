package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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
			return next(c)
		}
	}
}

var errAuditPlanMisMatch = errors.New("audit plan name don't match the token")

// AuditPlanVerifyAdapter is a `echo` middleware. Every audit plan should be
// scanner-scoped which means that scanner-A should not push SQL to scanner-B.
func AuditPlanVerifyAdapter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			apnInToken, err := utils.ParseAuditPlanName(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			apInPath := c.Param("audit_plan_name")
			if apnInToken != apInPath {
				return echo.NewHTTPError(http.StatusInternalServerError, errAuditPlanMisMatch.Error())
			}
			return next(c)
		}
	}
}
