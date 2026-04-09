//go:build !enterprise
// +build !enterprise

package dashboard

import (
	"context"

	"github.com/actiontech/sqle/sqle/model"
)

func init() {
	RegisterDetector(func(ctx context.Context, userID string) (bool, error) {
		storage := model.GetStorage()
		if storage == nil {
			return false, nil
		}

		// CE 版仅关注待审核/待上线的工单
		// 对应 V1 的状态过滤逻辑
		statusList := []string{
			model.WorkflowStatusWaitForAudit,
			model.WorkflowStatusWaitForExecution,
			model.WorkflowStatusReject,
		}

		data := map[string]interface{}{
			"filter_status_list":                  statusList,
			"filter_current_step_assignee_user_id": userID,
		}

		count, err := storage.GetGlobalWorkflowTotalNum(data)
		if err != nil {
			return false, err
		}

		return count > 0, nil
	})
}
