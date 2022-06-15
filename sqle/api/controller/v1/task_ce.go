//go:build !enterprise
// +build !enterprise

package v1

import (
	"fmt"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errNotSupportGetTaskAnalysisData = errors.New(errors.EnterpriseEditionFeatures, fmt.Errorf("get task analysis data is enterprise version functions"))

func getTaskAnalysisData(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errNotSupportGetTaskAnalysisData)
}
