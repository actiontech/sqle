//go:build !enterprise
// +build !enterprise

package v1

import (
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func init() {
	RegisterRedDotModules(GlobalDashBoardModule{})
}

type GlobalDashBoardModule struct{}

func (m GlobalDashBoardModule) Name() string {
	return "global_dashboard"
}

var statusOfGlobalWorkflowRequireAttention []string = []string{
	model.WorkflowStatusWaitForExecution,
	model.WorkflowStatusExecFailed,
	model.WorkflowStatusWaitForAudit,
	model.WorkflowStatusReject,
}

func (m GlobalDashBoardModule) HasRedDot(ctx echo.Context) (bool, error) {
	user, err := controller.GetCurrentUser(ctx, dms.GetUser)
	if err != nil {
		return false, err
	}
	permissions, isAdmin, err := dmsobject.GetUserOpPermission(ctx.Request().Context(), "", user.GetIDStr(), dms.GetDMSServerAddress())
	if err != nil {
		return false, err
	}
	// 2. 将用户权限信息，转化为全局待处理清单统一的用户可见性
	userVisibility := getGlobalDashBoardVisibilityOfUser(isAdmin, permissions)
	// 查询待关注工单，是否有未处理的工单
	filter, err := constructGlobalWorkflowBasicFilter(ctx.Request().Context(), user, userVisibility, permissions, &globalWorkflowBasicFilter{
		FilterStatusList: statusOfGlobalWorkflowRequireAttention,
	})
	if err != nil {
		return false, err
	}
	s := model.GetStorage()
	total, err := s.GetGlobalWorkflowTotalNum(filter)
	if err != nil {
		return false, err
	}
	if total > 0 {
		return true, nil
	}
	// 查询创建的工单，是否有未处理的工单
	userVisibility = GlobalDashBoardVisibilityGlobal
	filter, err = constructGlobalWorkflowBasicFilter(ctx.Request().Context(), user, userVisibility, permissions, &globalWorkflowBasicFilter{
		FilterStatusList:   statusOfGlobalWorkflowRequireAttention,
		FilterCreateUserId: user.GetIDStr(),
	})
	if err != nil {
		return false, err
	}
	total, err = s.GetGlobalWorkflowTotalNum(filter)
	if err != nil {
		return false, err
	}
	if total > 0 {
		return true, nil
	}
	return false, nil
}
