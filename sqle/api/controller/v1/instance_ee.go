//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

const ( // InstanceTipReqV1.FunctionalModule Enums
	sql_query = "sql_query"
)

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
		instances, err = s.GetInstanceTipsByUserAndOperation(user, req.FilterDBType, model.OP_SQL_QUERY_QUERY)
	default:
		instances, err = s.GetInstanceTipsByUser(user, req.FilterDBType)
	}
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	instanceTipsResV1 := make([]InstanceTipResV1, 0, len(instances))

	for _, inst := range instances {
		instanceTipRes := InstanceTipResV1{
			Name: inst.Name,
			Type: inst.DbType,
		}
		instanceTipsResV1 = append(instanceTipsResV1, instanceTipRes)
	}
	return c.JSON(http.StatusOK, &GetInstanceTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    instanceTipsResV1,
	})
}
