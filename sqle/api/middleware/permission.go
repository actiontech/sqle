package middleware

import (
	"context"
	"net/http"

	dmsJWT "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/labstack/echo/v4"
)

func OpGlobalAllowed() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, err := dmsJWT.GetUserUidStrFromContextWithOldJwt(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			up, err := dms.NewUserPermission(uid, "")
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			if up.CanOpGlobal() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}

func OpProjectAllowed() echo.MiddlewareFunc {
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

			if up.CanOpProject() {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}

func ViewGlobalAllowed() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, err := dmsJWT.GetUserUidStrFromContextWithOldJwt(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			up, err := dms.NewUserPermission(uid, "")
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			if up.CanViewGlobal() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}

func ViewProjectAllowed() echo.MiddlewareFunc {
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

			if up.CanViewProject() {
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

func ProjectMemberOpAllowed() echo.MiddlewareFunc {
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
			if up.CanOpProject() || up.IsProjectMember() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}

func ProjectMemberViewAllowed() echo.MiddlewareFunc {
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
			if up.CanViewProject() || up.IsProjectMember() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}
