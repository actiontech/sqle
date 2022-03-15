//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/labstack/echo/v4"
)

var errCommunityEditionDoesNotSupportWeChatConfiguration = errors.New(errors.ErrAccessDeniedError, e.New("community edition does not support WeChat configuration"))

func updateWeChatConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportWeChatConfiguration)
}

func getWeChatConfiguration(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportWeChatConfiguration)
}
