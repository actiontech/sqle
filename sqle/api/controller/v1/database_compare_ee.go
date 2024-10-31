//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

func getDatabaseComparison(c echo.Context) error {

	return c.JSON(http.StatusOK, &DatabaseComparisonResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    nil,
	})
}

func getComparisonStatement(c echo.Context) error {

	return c.JSON(http.StatusOK, &DatabaseComparisonStatementsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    nil,
	})
}

func genDatabaseDiffModifySQLs(c echo.Context) error {

	return c.JSON(http.StatusOK, &GenModifySQLResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    nil,
	})
}
