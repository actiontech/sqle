//go:build !enterprise
// +build !enterprise

package v1

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/labstack/echo/v4"
)

var errNotSupportOauth2 = errors.New(errors.EnterpriseEditionFeatures, fmt.Errorf("oauth2 related functions are enterprise version functions"))

func oauth2Link(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errNotSupportOauth2)
}

func oauth2Callback(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errNotSupportOauth2)
}

func bindOauth2User(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errNotSupportOauth2)
}
