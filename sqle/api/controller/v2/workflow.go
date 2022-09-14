package v2

import (
	"time"

	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type CreateWorkflowReqV2 struct {
	Subject string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc    string `json:"desc" form:"desc"`
	TaskIds []uint `json:"task_ids" form:"task_ids" valid:"required"`
}

// CreateWorkflowV2
// @Summary 创建工单
// @Description create workflow
// @Accept json
// @Produce json
// @Tags workflow
// @Id createWorkflowV2
// @Security ApiKeyAuth
// @Param instance body v2.CreateWorkflowReqV2 true "create workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v2/workflows [post]
func CreateWorkflowV2(c echo.Context) error {
	return nil
}

type GetWorkflowsResV2 struct {
	controller.BaseRes
	Data      []*WorkflowDetailResV2 `json:"data"`
	TotalNums uint64                 `json:"total_nums"`
}

type WorkflowDetailResV2 struct {
	Id                      uint       `json:"workflow_id"`
	Subject                 string     `json:"subject"`
	Desc                    string     `json:"desc"`
	CreateUser              string     `json:"create_user_name"`
	CreateTime              *time.Time `json:"create_time"`
	CurrentStepType         string     `json:"current_step_type,omitempty" enums:"sql_review,sql_execute"`
	CurrentStepAssigneeUser []string   `json:"current_step_assignee_user_name_list,omitempty"`
	Status                  string     `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,finished"`
}

// GetWorkflowsV2
// @Summary 获取工单列表
// @Description get workflow list
// @Tags workflow
// @Id getWorkflowsV2
// @Security ApiKeyAuth
// @Param filter_subject query string false "filter subject"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param filter_create_user_name query string false "filter create user name"
// @Param filter_status query string false "filter workflow status" Enums(wait_for_audit, wait_for_execution, rejected, canceled, exec_failed, finished)
// @Param filter_current_step_assignee_user_name query string false "filter current step assignee user name"
// @Param filter_task_instance_name query string false "filter instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v2.GetWorkflowsResV2
// @router /v2/workflows [get]
func GetWorkflowsV2(c echo.Context) error {
	return nil
}

type GetWorkflowResV2 struct {
	controller.BaseRes
	Data *WorkflowResV2 `json:"data"`
}

type WorkflowTaskItem struct {
	Id uint `json:"task_ids"`
}

type WorkflowRecordResV2 struct {
	TaskIds           []*WorkflowTaskItem     `json:"task_ids"`
	CurrentStepNumber uint                    `json:"current_step_number,omitempty"`
	Status            string                  `json:"status" enums:"wait_for_audit, wait_for_execution, rejected, canceled, exec_failed, finished"`
	Steps             []*v1.WorkflowStepResV1 `json:"workflow_step_list,omitempty"`
}

type WorkflowResV2 struct {
	Id            uint                   `json:"workflow_id"`
	Subject       string                 `json:"subject"`
	Desc          string                 `json:"desc,omitempty"`
	Mode          string                 `json:"mode" enums:"same_sqls,different_sqls"`
	CreateUser    string                 `json:"create_user_name"`
	CreateTime    *time.Time             `json:"create_time"`
	Record        *WorkflowRecordResV2   `json:"record"`
	RecordHistory []*WorkflowRecordResV2 `json:"record_history_list,omitempty"`
}

// GetWorkflowV2
// @Summary 获取工单详情
// @Description get workflow detail
// @Tags workflow
// @Id getWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path integer true "workflow id"
// @Success 200 {object} v2.GetWorkflowResV2
// @router /v2/workflows/{workflow_id}/ [get]
func GetWorkflowV2(c echo.Context) error {
	return nil
}

type UpdateWorkflowReqV2 struct {
	TaskIds []uint `json:"task_ids" form:"task_ids" valid:"required"`
}

// UpdateWorkflowV2
// @Summary 更新工单（驳回后才可更新）
// @Description update workflow when it is rejected to creator.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param instance body v2.UpdateWorkflowReqV2 true "update workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v2/workflows/{workflow_id}/ [patch]
func UpdateWorkflowV2(c echo.Context) error {
	return nil
}

// UpdateWorkflowScheduleV2
// @Summary 设置工单定时上线时间（设置为空则代表取消定时时间，需要SQL审核流程都通过后才可以设置）
// @Description update workflow schedule.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowScheduleV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param task_id path string true "task id"
// @Param instance body v1.UpdateWorkflowScheduleV1 true "update workflow schedule request"
// @Success 200 {object} controller.BaseRes
// @router /v2/workflows/{workflow_id}/task/{task_id}/schedule [put]
func UpdateWorkflowScheduleV2(c echo.Context) error {
	return nil
}
