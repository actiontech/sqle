//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var (
	errCommunityEditionDoesNotSupportFeatureExportWorkflowList = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support feature export workflow list"))
	errCommunityEditionDoesNotSupportWorkflowTemplate          = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support workflow template"))
)

func exportWorkflowV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportFeatureExportWorkflowList)
}

func getWorkflowTemplate(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportWorkflowTemplate)
}

func updateWorkflowTemplate(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportWorkflowTemplate)
}
