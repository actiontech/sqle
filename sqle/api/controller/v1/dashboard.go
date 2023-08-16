package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
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
	FilterProjectName string `json:"filter_project_name" query:"filter_project_name" form:"filter_project_name"`
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

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()

	createdNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_name":     req.FilterProjectName,
		"filter_create_user_name": user.Name,
		"filter_status":           model.WorkflowStatusWaitForAudit,
		"check_user_can_access":   false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	rejectedNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_name":     req.FilterProjectName,
		"filter_create_user_name": user.Name,
		"filter_status":           model.WorkflowStatusReject,
		"check_user_can_access":   false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	myNeedReviewNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_name":      req.FilterProjectName,
		"filter_status":            model.WorkflowStatusWaitForAudit,
		"filter_current_step_type": model.WorkflowStepTypeSQLReview,
		"filter_create_user_name":  user.Name,
		"check_user_can_access":    false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	myNeedExecuteNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_name":      req.FilterProjectName,
		"filter_status":            model.WorkflowStatusWaitForExecution,
		"filter_current_step_type": model.WorkflowStepTypeSQLExecute,
		"filter_create_user_name":  user.Name,
		"check_user_can_access":    false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	reviewNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_name":                    req.FilterProjectName,
		"filter_status":                          model.WorkflowStatusWaitForAudit,
		"filter_current_step_type":               model.WorkflowStepTypeSQLReview,
		"filter_current_step_assignee_user_name": user.Name,
		"check_user_can_access":                  false,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	executeNumber, err := s.GetWorkflowCountByReq(map[string]interface{}{
		"filter_project_name":                    req.FilterProjectName,
		"filter_status":                          model.WorkflowStatusWaitForExecution,
		"filter_current_step_type":               model.WorkflowStepTypeSQLExecute,
		"filter_current_step_assignee_user_name": user.Name,
		"check_user_can_access":                  false,
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

type DashboardProjectTipV1 struct {
	Name                    string `json:"project_name"`
	UnfinishedWorkflowCount int    `json:"unfinished_workflow_count"` // 只统计与当前用户相关的未完成工单
}

type GetDashboardProjectTipsResV1 struct {
	controller.BaseRes
	Data []*DashboardProjectTipV1 `json:"data"`
}

// DashboardProjectTipsV1
// @Summary 获取dashboard项目提示列表
// @Description get dashboard project tips
// @Tags dashboard
// @Id getDashboardProjectTipsV1
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} v1.GetDashboardProjectTipsResV1
// @router /v1/dashboard/project_tips [get]
func DashboardProjectTipsV1(c echo.Context) error {
	// TODO 暂不使用，避免页面报错
	// user, err := controller.GetCurrentUserFromDMS(c)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// s := model.GetStorage()
	// allProjectsByCurrentUser, err := s.GetProjectTips(controller.GetUserName(c))
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }
	// createdByMeWorkflowCounts, err := s.GetWorkflowCountForDashboardProjectTipsByReq(map[string]interface{}{
	// 	"filter_create_user_name": user.Name,
	// 	"filter_status":           []string{model.WorkflowStatusReject, model.WorkflowStatusWaitForAudit, model.WorkflowStatusWaitForExecution},
	// 	"check_user_can_access":   false,
	// })
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// needMeHandleWorkflowCounts, err := s.GetWorkflowCountForDashboardProjectTipsByReq(map[string]interface{}{
	// 	"filter_status":                          []string{model.WorkflowStatusWaitForAudit, model.WorkflowStatusWaitForExecution},
	// 	"filter_current_step_assignee_user_name": user.Name,
	// 	"check_user_can_access":                  false,
	// })
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// projectToWorkflowCount := make(map[string]int)
	// summingUpWorkflowCount := func(projectName string, records []*model.ProjectWorkflowCount) {
	// 	for _, record := range records {
	// 		if record.ProjectName != projectName {
	// 			continue
	// 		}
	// 		if workflowCount, ok := projectToWorkflowCount[record.ProjectName]; ok {
	// 			projectToWorkflowCount[record.ProjectName] = workflowCount + record.WorkflowCount
	// 		} else {
	// 			projectToWorkflowCount[record.ProjectName] = record.WorkflowCount
	// 		}
	// 	}
	// }

	// for _, p := range allProjectsByCurrentUser {
	// 	projectToWorkflowCount[p.Name] = 0
	// 	summingUpWorkflowCount(p.Name, createdByMeWorkflowCounts)
	// 	summingUpWorkflowCount(p.Name, needMeHandleWorkflowCounts)
	// }

	// // sort by unfinished workflow count, project name
	// i := 0
	// projectTips := make(dashboardProjectTipSort, len(projectToWorkflowCount))
	// for pName, count := range projectToWorkflowCount {
	// 	projectTips[i] = &DashboardProjectTipV1{
	// 		Name:                    pName,
	// 		UnfinishedWorkflowCount: count,
	// 	}
	// 	i++
	// }
	// sort.Sort(projectTips)

	// data := make([]*DashboardProjectTipV1, len(projectTips))
	// for j, projectTip := range projectTips {
	// 	data[j] = projectTip
	// }

	return c.JSON(http.StatusOK, &GetDashboardProjectTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    nil,
	})
}

type dashboardProjectTipSort []*DashboardProjectTipV1

func (m dashboardProjectTipSort) Len() int {
	return len(m)
}

func (m dashboardProjectTipSort) Less(i, j int) bool {
	return m[i].UnfinishedWorkflowCount < m[j].UnfinishedWorkflowCount
}

func (m dashboardProjectTipSort) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
