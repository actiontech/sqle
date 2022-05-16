//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionDoesNotSupportSqlQuery = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support SQL Query"))

func prepareSQLQuery(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportSqlQuery)
}

func getSQLResult(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportSqlQuery)
}

func getSQLQueryHistory(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportSqlQuery)
}
