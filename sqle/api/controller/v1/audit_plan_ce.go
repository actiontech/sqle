//go:build !enterprise
// +build !enterprise

package v1

import (
	"fmt"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errNotSupportGetAuditPlanAnalysisData = errors.New(errors.EnterpriseEditionFeatures, fmt.Errorf("get audit plan analysis data is enterprise version functions"))

func getAuditPlanAnalysisData(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errNotSupportGetAuditPlanAnalysisData)
}
