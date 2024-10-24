//go:build !enterprise
// +build !enterprise

package v1

import (
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/labstack/echo/v4"
)

func init() {
	RegisterRedDotModules(GlobalDashBoardModule{})
}

type GlobalDashBoardModule struct{}

func (m GlobalDashBoardModule) Name() string {
	return "global_dashboard"
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
	// 将用户权限信息，转化为全局待处理清单统一的用户可见性
	userVisibility := getGlobalDashBoardVisibilityOfUser(isAdmin, permissions)
	has, err := HasOthersWorkflowToHandle(ctx.Request().Context(), user, userVisibility)
	if err != nil {
		return false, err
	}
	if has {
		return true, nil
	}

	has, err = HasMyWorkflowToHandle(ctx.Request().Context(), user, userVisibility)
	if err != nil {
		return false, err
	}
	return has, nil
}
