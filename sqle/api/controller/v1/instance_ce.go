//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionDoesNotSupportListTables = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support list tables"))
var errCommunityEditionDoesNotSupportGetTablesMetadata = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support get table metadata"))

func listTableBySchema(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportListTables)
}

func getTableMetadata(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportGetTablesMetadata)
}
