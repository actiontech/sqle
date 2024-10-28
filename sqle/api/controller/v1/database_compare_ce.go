//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var ErrCommunityEditionNotSupportDatabaseCompare = errors.New(errors.EnterpriseEditionFeatures, e.New("database compare is enterprise version feature"))

func getDatabaseCompare(c echo.Context) error {
	return ErrCommunityEditionNotSupportDatabaseCompare
}

func getCompareStatement(c echo.Context) error {
	return ErrCommunityEditionNotSupportDatabaseCompare
}

func genDatabaseDiffModifySQLs(c echo.Context) error {
	return ErrCommunityEditionNotSupportDatabaseCompare
}
