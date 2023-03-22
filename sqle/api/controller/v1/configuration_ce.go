//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionNotSupportCostumeLogo = e.New("costume logo is enterprise version feature")

func uploadLogo(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCostumeLogo)
}

func getLogo(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCostumeLogo)
}

func updatePersonaliseConfig(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCostumeLogo)
}
