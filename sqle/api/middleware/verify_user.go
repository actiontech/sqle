package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
)

func VerifyUserIsDisabled() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userName := controller.GetUserName(c)
			user, isExist, err := model.GetStorage().GetUserDetailByName(userName)
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			if !isExist {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("user is not exist")))
			}
			if user.IsDisabled {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataDisabled, fmt.Errorf("current user is disabled.")))
			}
			return next(c)
		}
	}
}
