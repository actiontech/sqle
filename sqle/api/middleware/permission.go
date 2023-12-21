package middleware

import (
	"context"
	"net/http"

	dmsJWT "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/labstack/echo/v4"
)

// AdminUserAllowed is a `echo` middleware, only allow admin user to access next.
func AdminUserAllowed() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, err := dmsJWT.GetUserUidStrFromContextWithOldJwt(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			up, err := dms.NewUserPermission(uid, "700300" /*TODO 支持不传空间 */)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			if up.IsAdmin() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}

func ProjectAdminUserAllowed() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, err := dmsJWT.GetUserUidStrFromContextWithOldJwt(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			up, err := dms.NewUserPermission(uid, projectUid)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			if up.IsAdmin() || up.IsProjectAdmin() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}

func ProjectMemberAllowed() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, err := dmsJWT.GetUserUidStrFromContextWithOldJwt(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			up, err := dms.NewUserPermission(uid, projectUid)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			if up.IsAdmin() || up.IsProjectAdmin() || up.IsProjectMember() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}
