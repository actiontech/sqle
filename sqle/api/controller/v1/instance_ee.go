//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

func listTableBySchema(c echo.Context) error {
	// 功能已经废弃，但是接口还在，后续可能在其他功能上次引入。
	return controller.JSONBaseErrorReq(c, errors.New(errors.FeatureNotImplemented, fmt.Errorf("ListTableBySchema not implement")))
}

func getTableMetadata(c echo.Context) error {
	// 功能已经废弃，但是接口还在，后续可能在其他功能上会次引入。
	return controller.JSONBaseErrorReq(c, errors.New(errors.FeatureNotImplemented, fmt.Errorf("GetTableMetadata not implement")))
}
