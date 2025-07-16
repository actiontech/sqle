//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var ErrCommunityEditionNotSupportSqlInsight = errors.New(errors.EnterpriseEditionFeatures, e.New("sql insight is enterprise version feature"))

func getSqlPerformanceInsights(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlInsight
}

func getSqlPerformanceInsightsRelatedSQL(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlInsight
}

func getSqlPerformanceInsightsRelatedTransaction(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlInsight
}
