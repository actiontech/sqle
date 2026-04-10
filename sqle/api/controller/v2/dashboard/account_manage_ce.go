//go:build !dms
// +build !dms

package dashboard

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

func getGlobalAccountListV2(c echo.Context) error {
	return controller.JSONOnlySupportForEnterpriseVersionErr(c)
}

func getGlobalAccountStatisticsV2(c echo.Context) error {
	return controller.JSONOnlySupportForEnterpriseVersionErr(c)
}
