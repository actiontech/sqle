//go:build !enterprise
// +build !enterprise

package dashboard

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

func getGlobalWorkflowStatisticsV2(c echo.Context) error {
	return controller.JSONOnlySupportForEnterpriseVersionErr(c)
}

func getGlobalWorkflowListV2(c echo.Context) error {
	return controller.JSONOnlySupportForEnterpriseVersionErr(c)
}
