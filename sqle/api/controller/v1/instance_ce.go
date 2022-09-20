//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"net/http"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionDoesNotSupportListTables = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support list tables"))
var errCommunityEditionDoesNotSupportGetTablesMetadata = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support get table metadata"))

func getInstanceTips(c echo.Context) error {
	req := new(InstanceTipReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var instances []*model.Instance
	switch req.FunctionalModule {
	case create_audit_plan:
		instances, err = s.GetInstanceTipsByUserAndOperation(user, req.FilterDBType, model.OP_AUDIT_PLAN_SAVE)
	case sql_query:
		instances, err = s.GetInstanceTipsByUser(user, req.FilterDBType)
	default: // create_workflow case
		instances, err = s.GetInstancesTipsByUserAndTypeAndTempId(user, req.FilterDBType, req.FilterWorkflowTemplateId)
	}
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceTipsResV1 := make([]InstanceTipResV1, 0, len(instances))
	for _, inst := range instances {
		instanceTipRes := InstanceTipResV1{
			Name:               inst.Name,
			Type:               inst.DbType,
			WorkflowTemplateId: uint32(inst.WorkflowTemplateId),
		}
		instanceTipsResV1 = append(instanceTipsResV1, instanceTipRes)
	}

	return c.JSON(http.StatusOK, &GetInstanceTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    instanceTipsResV1,
	})
}

func listTableBySchema(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportListTables)
}

func getTableMetadata(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportGetTablesMetadata)
}
