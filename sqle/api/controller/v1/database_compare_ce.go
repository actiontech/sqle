//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var ErrCommunityEditionNotSupportDatabaseStructComparison = errors.New(errors.EnterpriseEditionFeatures, e.New("database struct comparison is enterprise version feature"))

func getDatabaseComparison(c echo.Context) error {
	return ErrCommunityEditionNotSupportDatabaseStructComparison
}

func getComparisonStatement(c echo.Context) error {
	return ErrCommunityEditionNotSupportDatabaseStructComparison
}

func genDatabaseDiffModifySQLs(c echo.Context) error {
	return ErrCommunityEditionNotSupportDatabaseStructComparison
}
