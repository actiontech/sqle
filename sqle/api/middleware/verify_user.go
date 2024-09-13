package middleware

import (
	dmsJWT "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	dmsObject "github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

func VerifyUserIsDisabled() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, err := dmsJWT.GetUserUidStrFromContextWithOldJwt(c)
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			user, err := dmsObject.GetUser(c.Request().Context(), uid, controller.GetDMSServerAddress())
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			if user.Stat != "正常" && user.Stat != "Normal" { // todo i18n user Stat
				return controller.JSONBaseErrorReq(c, errors.NewUserDisabledErr("current user status is %s", user.Stat))
			}
			return next(c)
		}
	}
}
