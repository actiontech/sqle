//go:build enterprise
// +build enterprise

package v1

import (
	"context"

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

func (m GlobalDashBoardModule) HasRedDot(ctx echo.Context) (bool, error) {
	user, err := controller.GetCurrentUser(ctx, dms.GetUser)
	if err != nil {
		return false, err
	}
	permissions, isAdmin, err := dmsobject.GetUserOpPermission(ctx.Request().Context(), "", user.GetIDStr(), dms.GetDMSServerAddress())
	if err != nil {
		return false, err
	}
	// 将用户权限信息，转化为全局待处理清单统一的用户可视范围
	userVisibility := getGlobalDashBoardVisibilityOfUser(isAdmin, permissions)

	has, err := HasSqlManageToHandle(ctx.Request().Context(), user, userVisibility)
	if err != nil {
		return false, err
	}
	if has {
		return true, nil
	}
	has, err = HasOthersWorkflowToHandle(ctx.Request().Context(), user, userVisibility)
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

// 查询待关注工单，是否有未处理的工单
func HasSqlManageToHandle(ctx context.Context, user *model.User, userVisibility GlobalDashBoardVisibility) (has bool, err error) {
	filter, err := constructGlobalSqlManageBasicFilter(ctx, user, userVisibility, &globalSqlManageBasicFilter{})
	if err != nil {
		return false, err
	}
	s := model.GetStorage()
	total, err := s.GetGlobalSqlManageTotalNum(filter)
	if err != nil {
		return false, err
	}
	if total > 0 {
		return true, nil
	}
	return false, nil
}
