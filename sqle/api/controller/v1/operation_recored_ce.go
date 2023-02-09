//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionDoesNotSupportFeatureOperationRecord = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support feature operation record"))

func getOperationTypeNameList(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportFeatureOperationRecord)
}

func getOperationActionList(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportFeatureOperationRecord)
}
