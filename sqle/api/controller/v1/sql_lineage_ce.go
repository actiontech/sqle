//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var ErrCommunityEditionNotSupportSQLLineage = errors.New(
	errors.EnterpriseEditionFeatures,
	e.New("sql lineage analysis is enterprise version feature"),
)

func sqlLineageAnalyze(c echo.Context) error {
	return ErrCommunityEditionNotSupportSQLLineage
}