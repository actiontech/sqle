//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionDoesNotSupportSyncInstance = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support sync instance"))

func createSyncInstanceTask(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportSyncInstance)
}

func updateSyncInstanceTask(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportSyncInstance)
}

func deleteSyncInstanceTask(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportSyncInstance)
}

func triggerSyncInstance(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportSyncInstance)
}
