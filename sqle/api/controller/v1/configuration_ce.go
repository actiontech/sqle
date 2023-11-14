//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

var (
	errCommunityEditionNotSupportFeishuAudit = errors.New(errors.EnterpriseEditionFeatures, e.New("feishu audit is enterprise version feature"))
)

func updateFeishuAuditConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportFeishuAudit)
}

func getFeishuAuditConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportFeishuAudit)
}

func testFeishuAuditConfigV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportFeishuAudit)
}
