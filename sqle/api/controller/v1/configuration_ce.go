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
	errCommunityEditionNotSupportFeishuAudit         = errors.New(errors.EnterpriseEditionFeatures, e.New("feishu audit is enterprise version feature"))
	errCommunityEditionNotSupportDingDingAudit       = errors.New(errors.EnterpriseEditionFeatures, e.New("dingding audit is enterprise version feature"))
	errCommunityEditionNotSupportWechatAudit         = errors.New(errors.EnterpriseEditionFeatures, e.New("wechat audit is enterprise version feature"))
	errCommunityEditionDoesNotSupportScheduledNotify = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support workflow scheduled notify"))
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

func getDingTalkConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportDingDingAudit)
}

func updateDingTalkConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportDingDingAudit)
}

func testDingTalkConfigV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportDingDingAudit)
}

func getWechatAuditConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func updateWechatAuditConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func testWechatAuditConfigV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func getScheduledTaskDefaultOptionV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportScheduledNotify)
}
