package v2

import (
	"time"

	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type WorkflowStepResV2 struct {
	Id            uint       `json:"workflow_step_id,omitempty"`
	Number        uint       `json:"number"`
	Type          string     `json:"type" enums:"create_workflow,update_workflow,sql_review,sql_execute"`
	Desc          string     `json:"desc,omitempty"`
	Users         []string   `json:"assignee_user_name_list,omitempty"`
	OperationUser string     `json:"operation_user_name,omitempty"`
	OperationTime *time.Time `json:"operation_time,omitempty"`
	State         string     `json:"state,omitempty" enums:"initialized,approved,rejected"`
	Reason        string     `json:"reason,omitempty"`
}

// @Summary 审批通过
// @Description approve workflow
// @Tags workflow
// @Id approveWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param workflow_step_id path string true "workflow step id"
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/steps/{workflow_step_id}/approve [post]
func ApproveWorkflowV2(c echo.Context) error {
	return nil
}

type RejectWorkflowReqV2 struct {
	Reason string `json:"reason" form:"reason"`
}

// @Summary 审批驳回
// @Description reject workflow
// @Tags workflow
// @Id rejectWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Param workflow_step_id path string true "workflow step id"
// @param workflow_approve body v2.RejectWorkflowReqV2 true "workflow approve request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/steps/{workflow_step_id}/reject [post]
func RejectWorkflowV2(c echo.Context) error {
	return nil
}

// @Summary 审批关闭（中止）
// @Description cancel workflow
// @Tags workflow
// @Id cancelWorkflowV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param workflow_id path string true "workflow id"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/cancel [post]
func CancelWorkflowV2(c echo.Context) error {
	return nil
}

type BatchCancelWorkflowsReqV2 struct {
	WorkflowIDList []string `json:"workflow_id_list" form:"workflow_id_list"`
}

// BatchCancelWorkflowsV2 batch cancel workflows.
// @Summary 批量取消工单
// @Description batch cancel workflows
// @Tags workflow
// @Id batchCancelWorkflowsV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param BatchCancelWorkflowsReqV2 body v2.BatchCancelWorkflowsReqV2 true "batch cancel workflows request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/cancel [post]
func BatchCancelWorkflowsV2(c echo.Context) error {
	return nil
}

type BatchCompleteWorkflowsReqV2 struct {
	WorkflowIDList []string `json:"workflow_id_list" form:"workflow_id_list"`
}

// BatchCompleteWorkflowsV2 complete workflows.
// @Summary 批量完成工单
// @Description this api will directly change the work order status to finished without real online operation
// @Tags workflow
// @Id batchCompleteWorkflowsV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param data body v2.BatchCompleteWorkflowsReqV2 true "batch complete workflows request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/complete [post]
func BatchCompleteWorkflowsV2(c echo.Context) error {
	return nil
}

// ExecuteOneTaskOnWorkflowV2
// @Summary 工单提交单个数据源上线
// @Description execute one task on workflow
// @Tags workflow
// @Id executeOneTaskOnWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Param task_id path string true "task id"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/execute [post]
func ExecuteOneTaskOnWorkflowV2(c echo.Context) error {
	return nil
}

type GetWorkflowTasksResV2 struct {
	controller.BaseRes
	Data []*GetWorkflowTasksItemV2 `json:"data"`
}

type GetWorkflowTasksItemV2 struct {
	TaskId                   uint                       `json:"task_id"`
	InstanceName             string                     `json:"instance_name"`
	Status                   string                     `json:"status" enums:"wait_for_audit,wait_for_execution,exec_scheduled,exec_failed,exec_succeeded,executing,manually_executed"`
	ExecStartTime            *time.Time                 `json:"exec_start_time,omitempty"`
	ExecEndTime              *time.Time                 `json:"exec_end_time,omitempty"`
	ScheduleTime             *time.Time                 `json:"schedule_time,omitempty"`
	CurrentStepAssigneeUser  []string                   `json:"current_step_assignee_user_name_list,omitempty"`
	TaskPassRate             float64                    `json:"task_pass_rate"`
	TaskScore                int32                      `json:"task_score"`
	InstanceMaintenanceTimes []*v1.MaintenanceTimeResV1 `json:"instance_maintenance_times"`
	ExecutionUserName        string                     `json:"execution_user_name"`
}

// GetSummaryOfWorkflowTasksV2
// @Summary 获取工单数据源任务概览
// @Description get summary of workflow instance tasks
// @Tags workflow
// @Id getSummaryOfInstanceTasksV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Success 200 {object} v2.GetWorkflowTasksResV2
// @router /v2/projects/{project_name}/workflows/{workflow_id}/tasks [get]
func GetSummaryOfWorkflowTasksV2(c echo.Context) error {
	return nil
}

type CreateWorkflowReqV2 struct {
	Subject    string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	WorkflowId string `json:"workflow_id" form:"workflow_id" valid:"required"`
	Desc       string `json:"desc" form:"desc"`
	TaskIds    []uint `json:"task_ids" form:"task_ids" valid:"required"`
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
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows [post]
func CreateWorkflowV2(c echo.Context) error {
	return nil
}

type GetWorkflowsReqV2 struct {
	FilterSubject                     string `json:"filter_subject" query:"filter_subject"`
	FilterCreateTimeFrom              string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo                string `json:"filter_create_time_to" query:"filter_create_time_to"`
	FilterCreateUserName              string `json:"filter_create_user_name" query:"filter_create_user_name"`
	FilterStatus                      string `json:"filter_status" query:"filter_status" valid:"omitempty,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterCurrentStepAssigneeUserName string `json:"filter_current_step_assignee_user_name" query:"filter_current_step_assignee_user_name"`
	FilterTaskInstanceName            string `json:"filter_task_instance_name" query:"filter_task_instance_name"`
	FilterTaskExecuteStartTimeFrom    string `json:"filter_task_execute_start_time_from" query:"filter_task_execute_start_time_from"`
	FilterTaskExecuteStartTimeTo      string `json:"filter_task_execute_start_time_to" query:"filter_task_execute_start_time_to"`
	PageIndex                         uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize                          uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetWorkflowsResV2 struct {
	controller.BaseRes
	Data      []*WorkflowDetailResV2 `json:"data"`
	TotalNums uint64                 `json:"total_nums"`
}

type WorkflowDetailResV2 struct {
	ProjectName             string     `json:"project_name"`
	Name                    string     `json:"workflow_name"`
	WorkflowId              string     `json:"workflow_id" `
	Desc                    string     `json:"desc"`
	CreateUser              string     `json:"create_user_name"`
	CreateTime              *time.Time `json:"create_time"`
	CurrentStepType         string     `json:"current_step_type,omitempty" enums:"sql_review,sql_execute"`
	CurrentStepAssigneeUser []string   `json:"current_step_assignee_user_name_list,omitempty"`
	Status                  string     `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,executing,finished"`
}

// GetGlobalWorkflowsV2
// @Summary 获取全局工单列表
// @Description get global workflow list
// @Tags workflow
// @Id getGlobalWorkflowsV2
// @Security ApiKeyAuth
// @Param filter_subject query string false "filter subject"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param filter_task_execute_start_time_from query string false "filter_task_execute_start_time_from"
// @Param filter_task_execute_start_time_to query string false "filter_task_execute_start_time_to"
// @Param filter_create_user_name query string false "filter create user name"
// @Param filter_status query string false "filter workflow status" Enums(wait_for_audit,wait_for_execution,rejected,executing,canceled,exec_failed,finished)
// @Param filter_current_step_assignee_user_name query string false "filter current step assignee user name"
// @Param filter_task_instance_name query string false "filter instance name"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v2.GetWorkflowsResV2
// @router /v2/workflows [get]
func GetGlobalWorkflowsV2(c echo.Context) error {
	return nil
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
// @Param filter_task_execute_start_time_from query string false "filter_task_execute_start_time_from"
// @Param filter_task_execute_start_time_to query string false "filter_task_execute_start_time_to"
// @Param filter_create_user_name query string false "filter create user name"
// @Param filter_status query string false "filter workflow status" Enums(wait_for_audit,wait_for_execution,rejected,executing,canceled,exec_failed,finished)
// @Param filter_current_step_assignee_user_name query string false "filter current step assignee user name"
// @Param filter_task_instance_name query string false "filter instance name"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Param project_name path string true "project name"
// @Success 200 {object} v2.GetWorkflowsResV2
// @router /v2/projects/{project_name}/workflows [get]
func GetWorkflowsV2(c echo.Context) error {
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
// @Param project_name path string true "project name"
// @Param instance body v2.UpdateWorkflowReqV2 true "update workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/ [patch]
func UpdateWorkflowV2(c echo.Context) error {
	return nil
}

type UpdateWorkflowScheduleReqV2 struct {
	ScheduleTime *time.Time `json:"schedule_time"`
}

// UpdateWorkflowScheduleV2
// @Summary 设置工单数据源定时上线时间（设置为空则代表取消定时时间，需要SQL审核流程都通过后才可以设置）
// @Description update workflow schedule.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowScheduleV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param task_id path string true "task id"
// @Param project_name path string true "project name"
// @Param instance body v2.UpdateWorkflowScheduleReqV2 true "update workflow schedule request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/schedule [put]
func UpdateWorkflowScheduleV2(c echo.Context) error {
	return nil
}

// ExecuteTasksOnWorkflowV2
// @Summary 多数据源批量上线
// @Description execute tasks on workflow
// @Tags workflow
// @Id executeTasksOnWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/tasks/execute [post]
func ExecuteTasksOnWorkflowV2(c echo.Context) error {
	return nil
}

type GetWorkflowResV2 struct {
	controller.BaseRes
	Data *WorkflowResV2 `json:"data"`
}

type WorkflowTaskItem struct {
	Id uint `json:"task_id"`
}

type WorkflowRecordResV2 struct {
	Tasks             []*WorkflowTaskItem  `json:"tasks"`
	CurrentStepNumber uint                 `json:"current_step_number,omitempty"`
	Status            string               `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,executing,finished"`
	Steps             []*WorkflowStepResV2 `json:"workflow_step_list,omitempty"`
}

type WorkflowResV2 struct {
	Name          string                 `json:"workflow_name"`
	WorkflowID    string                 `json:"workflow_id"`
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
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Success 200 {object} GetWorkflowResV2
// @router /v2/projects/{project_name}/workflows/{workflow_id}/ [get]
func GetWorkflowV2(c echo.Context) error {
	return nil
}
