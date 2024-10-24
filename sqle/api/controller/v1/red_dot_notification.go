package v1

import (
	"context"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type RedDotModule interface {
	Name() string
	HasRedDot(ctx echo.Context) (bool, error)
}

type RedDot struct {
	ModuleName string
	HasRedDot  bool
}

var redDotList []RedDotModule

func RegisterRedDotModules(redDotModule ...RedDotModule) {
	redDotList = append(redDotList, redDotModule...)
}

func GetSystemModuleRedDotsList(ctx echo.Context) ([]*RedDot, error) {
	redDots := make([]*RedDot, len(redDotList))
	for i, rd := range redDotList {
		hasRedDot, err := rd.HasRedDot(ctx)
		if err != nil {
			return nil, err
		}
		redDots[i] = &RedDot{
			ModuleName: rd.Name(),
			HasRedDot:  hasRedDot,
		}
	}
	return redDots, nil
}

var statusOfGlobalWorkflowToBeFiltered []string = []string{
	model.WorkflowStatusWaitForExecution,
	model.WorkflowStatusExecFailed,
	model.WorkflowStatusWaitForAudit,
	model.WorkflowStatusReject,
}

// 查询待关注工单，是否有未处理的工单
func HasOthersWorkflowToHandle(ctx context.Context, user *model.User, userVisibility GlobalDashBoardVisibility) (has bool, err error) {
	filter, err := constructGlobalWorkflowBasicFilter(ctx, user, userVisibility, &globalWorkflowBasicFilter{
		FilterStatusList: statusOfGlobalWorkflowToBeFiltered,
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
	return false, nil
}

// 查询创建的工单，是否有未处理的工单
func HasMyWorkflowToHandle(ctx context.Context, user *model.User, userVisibility GlobalDashBoardVisibility) (has bool, err error) {
	filter, err := constructGlobalWorkflowBasicFilter(ctx, user, userVisibility, &globalWorkflowBasicFilter{
		FilterStatusList:   statusOfGlobalWorkflowToBeFiltered,
		FilterCreateUserId: user.GetIDStr(),
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
	return false, nil
}
