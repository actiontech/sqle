package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type GetDashboardResV1 struct {
	controller.BaseRes
	Data *DashboardResV1 `json:"data"`
}

type DashboardResV1 struct {
	WorkflowStatistics *WorkflowStatisticsResV1 `json:"workflow_statistics"`
}

type WorkflowStatisticsResV1 struct {
	MyWorkflowNumber            uint64 `json:"my_on_process_workflow_number"`
	MyRejectedWorkflowNumber    uint64 `json:"my_rejected_workflow_number"`
	MyNeedReviewWorkflowNumber  uint64 `json:"my_need_review_workflow_number"`
	MyNeedExecuteWorkflowNumber uint64 `json:"my_need_execute_workflow_number"`
	NeedMeReviewNumber          uint64 `json:"need_me_to_review_workflow_number"`
	NeedMeExecuteNumber         uint64 `json:"need_me_to_execute_workflow_number"`
}

type GetDashboardReqV1 struct {
	FilterProjectName string `json:"filter_project_id" query:"filter_project_id" form:"filter_project_id"`
}

// @Summary 获取 dashboard 信息
// @Description get dashboard info
// @Id getDashboardV1
// @Tags dashboard
// @Security ApiKeyAuth
// @Param filter_project_name query string false "filter project name"
// @Produce json
// @Success 200 {object} v1.GetDashboardResV1
// @router /v1/dashboard [get]
func Dashboard(c echo.Context) error {
	req := new(GetDashboardReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), req.FilterProjectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()

	createdNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_id":     projectUid,
		"filter_create_user_id": user.GetIDStr(),
		"filter_status":         model.WorkflowStatusWaitForAudit,
		"check_user_can_access": false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	rejectedNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_id":     projectUid,
		"filter_create_user_id": user.GetIDStr(),
		"filter_status":         model.WorkflowStatusReject,
		"check_user_can_access": false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	myNeedReviewNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_id":        projectUid,
		"filter_status":            model.WorkflowStatusWaitForAudit,
		"filter_current_step_type": model.WorkflowStepTypeSQLReview,
		"filter_create_user_id":    user.GetIDStr(),
		"check_user_can_access":    false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	myNeedExecuteNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_id":        projectUid,
		"filter_status":            model.WorkflowStatusWaitForExecution,
		"filter_current_step_type": model.WorkflowStepTypeSQLExecute,
		"filter_create_user_id":    user.GetIDStr(),
		"check_user_can_access":    false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	reviewNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_id":                    projectUid,
		"filter_status":                        model.WorkflowStatusWaitForAudit,
		"filter_current_step_type":             model.WorkflowStepTypeSQLReview,
		"filter_current_step_assignee_user_id": user.GetIDStr(),
		"check_user_can_access":                false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	executeNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_id":                    projectUid,
		"filter_status":                        model.WorkflowStatusWaitForExecution,
		"filter_current_step_type":             model.WorkflowStepTypeSQLExecute,
		"filter_current_step_assignee_user_id": user.GetIDStr(),
		"check_user_can_access":                false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflowStatisticsRes := &WorkflowStatisticsResV1{
		MyWorkflowNumber:            createdNumber, // todo 这个返回字段没有再用到了，可以在V2移除
		MyRejectedWorkflowNumber:    rejectedNumber,
		MyNeedReviewWorkflowNumber:  myNeedReviewNumber,
		MyNeedExecuteWorkflowNumber: myNeedExecuteNumber,
		NeedMeReviewNumber:          reviewNumber,
		NeedMeExecuteNumber:         executeNumber,
	}
	return c.JSON(http.StatusOK, &GetDashboardResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &DashboardResV1{
			WorkflowStatistics: workflowStatisticsRes,
		},
	})
}
