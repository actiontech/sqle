//go:build enterprise
// +build enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/labstack/echo/v4"
)

var ErrNoLicenseRequired = errors.New(errors.ErrAccessDeniedError, e.New("sqle-ce no license required"))

func getLicense(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, ErrNoLicenseRequired)
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
