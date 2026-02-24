//go:build !enterprise
// +build !enterprise

package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

// ========== AI 智能中心 CE 版本（未实现） ==========

// getAIHubBanner CE 版本 - 未实现
func getAIHubBanner(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.New(errors.EnterpriseEditionFeatures, nil))
}

// getAIHubStrategicValue CE 版本 - 未实现
func getAIHubStrategicValue(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.New(errors.EnterpriseEditionFeatures, nil))
}

// getAIHubManagementView CE 版本 - 未实现
func getAIHubManagementView(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.New(errors.EnterpriseEditionFeatures, nil))
}

// getAIHubExecutionData CE 版本 - 未实现
func getAIHubExecutionData(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.New(errors.EnterpriseEditionFeatures, nil))
}
