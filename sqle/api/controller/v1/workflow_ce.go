//go:build !enterprise
// +build !enterprise

package v1

import (
	"context"
	e "errors"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
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

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	td := model.DefaultWorkflowTemplate(projectUid)

	return c.JSON(http.StatusOK, &GetWorkflowTemplateResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowTemplateToRes(td),
	})
}

func updateWorkflowTemplate(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportWorkflowTemplate)
}
