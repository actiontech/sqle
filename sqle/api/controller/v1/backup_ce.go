//go:build !enterprise
// +build !enterprise

package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/labstack/echo/v4"
)

func getBackupSqlList(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.NewNotSupportBackupErr())
}

func createRollbackWorkflow(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.NewNotSupportBackupErr())
}

func updateSqlBackupStrategy(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.NewNotSupportBackupErr())
}

func updateTaskBackupStrategy(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.NewNotSupportBackupErr())
}
