//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionNotSupportCompanyNotice = errors.New(errors.EnterpriseEditionFeatures, e.New("company notice is enterprise version feature"))

func getCompanyNotice(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCompanyNotice)
}

func updateCompanyNotice(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCompanyNotice)
}
