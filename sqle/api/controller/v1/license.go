//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/labstack/echo/v4"
)

var NoLicenseRequiredError = errors.New(errors.ErrAccessDeniedError, e.New("sqle-ce no license required"))

type GetLicenseResV1 struct {
	LicenseContent string `json:"license_content" example:"This license is for: &{ExpireDate:2022-02-10 Version:99.99.99 AgentCount:2};;iVWLgIfzYtIFlMEIMTxX2~S8lgXsNqT4Ccug23GybWsiP0i1SW8GaorcbRvLGdD4X1v4VbFU77zqg1_1TisP;;U7gAUCECm86~kodfMDQSUdEd3QHR5MXMKp2KFFcjb8_NliBt"`
}

// GetLicense get sqle license
// @Summary 获取 sqle license
// @Description get sqle license
// @Id getSQLELicenseV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetLicenseResV1
// @router /v1/configurations/license [get]
func GetLicense(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, NoLicenseRequiredError)
}

// GetSQLEServerInfo get information about the machine where SQLE is located
// @Summary 获取 sqle 所在机器的信息
// @Description get information about the machine where SQLE is located
// @Id GetSQLEServerInfoV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 file 1 "server info"
// @router /v1/configurations/sqle_server_info [get]
func GetSQLEServerInfo(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, NoLicenseRequiredError)
}

// SetLicense set sqle license
// @Summary 导入 sqle license
// @Description set sqle license
// @Id setSQLELicenseV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param license_file formData file true "SQLE license file"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/license [post]
func SetLicense(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, NoLicenseRequiredError)
}
