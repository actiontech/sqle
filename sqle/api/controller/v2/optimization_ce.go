//go:build !enterprise
// +build !enterprise

package v2

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

func sqlOptimize(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}

func getOptimizationRecords(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}

func getOptimizationSQL(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}
