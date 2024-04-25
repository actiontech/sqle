//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

func sqlOptimizate(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}

func getOptimizationRecord(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}

func getOptimizationRecords(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}

func getOptimizationSQL(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}

func getOptimizationSQLs(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}

func getDBPerformanceImproveOverview(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}

func getOptimizationRecordOverview(c echo.Context) error {
	return errors.New(errors.SQLOptimizationCommunityNotSupported, e.New("sql optimization community not supported"))
}
