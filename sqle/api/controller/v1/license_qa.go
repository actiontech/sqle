//go:build !release
// +build !release

package v1

import (
	e "errors"
	"net/http"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

var ErrNoLicenseRequired = errors.New(errors.EnterpriseEditionFeatures, e.New("sqle-qa version has unlimited resources does not need to set license"))

func getLicense(c echo.Context) error {
	return c.JSON(http.StatusOK, GetLicenseResV1{
		BaseRes: controller.NewBaseReq(nil),
		License: []LicenseItem{
			{
				Description: "实例数",
				Name:        "instance_num",
				Limit:       "无限制",
			},
			{
				Description: "用户数",
				Name:        "user",
				Limit:       "无限制",
			},
			{
				Description: "授权运行时长(天)",
				Name:        "work duration day",
				Limit:       "无限制",
			},
		},
	})
}

func getSQLELicenseInfo(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, ErrNoLicenseRequired)
}

func setLicense(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, ErrNoLicenseRequired)
}

func checkLicense(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, ErrNoLicenseRequired)
}
