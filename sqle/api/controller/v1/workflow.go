package v1

import (
	"context"
	e "errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/labstack/echo/v4"
)

var ErrWorkflowNoAccess = errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist or you can't access it"))

var ErrForbidMyBatisXMLTask = func(taskId uint) error {
	return errors.New(errors.DataConflict,
		fmt.Errorf("the task for audit mybatis xml file is not allow to create workflow. taskId=%v", taskId))
}

var ErrCanNotTerminateExecute = func(workflowStatus, taskStatus string) error {
	return errors.NewDataInvalidErr(
		"workflow status is %s and task status is %s, termination can not be performed",
		workflowStatus, taskStatus)
}

var ErrWorkflowExecuteTimeIncorrect = errors.New(errors.TaskActionInvalid, fmt.Errorf("please go online during instance operation and maintenance time"))

type GetWorkflowTemplateResV1 struct {
	controller.BaseRes
	Data *WorkflowTemplateDetailResV1 `json:"data"`
}

type WorkflowTemplateDetailResV1 struct {
	Name                          string                       `json:"workflow_template_name"`
	Desc                          string                       `json:"desc,omitempty"`
	AllowSubmitWhenLessAuditLevel string                       `json:"allow_submit_when_less_audit_level" enums:"normal,notice,warn,error"`
	Steps                         []*WorkFlowStepTemplateResV1 `json:"workflow_step_template_list"`
	UpdateTime                    time.Time                    `json:"update_time"`
}

type WorkFlowStepTemplateResV1 struct {
	Number               int      `json:"number"`
	Typ                  string   `json:"type"`
	Desc                 string   `json:"desc,omitempty"`
	ApprovedByAuthorized bool     `json:"approved_by_authorized"`
	ExecuteByAuthorized  bool     `json:"execute_by_authorized"`
	Users                []string `json:"assignee_user_id_list"`
}

// @Summary 获取审批流程模板详情
// @Description get workflow template detail
// @Tags workflow
// @Id getWorkflowTemplateV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetWorkflowTemplateResV1
// @router /v1/projects/{project_name}/workflow_template [get]
func GetWorkflowTemplate(c echo.Context) error {
	return getWorkflowTemplate(c)
}

func convertWorkflowTemplateToRes(template *model.WorkflowTemplate) *WorkflowTemplateDetailResV1 {
	res := &WorkflowTemplateDetailResV1{
		Name:                          template.Name,
		Desc:                          template.Desc,
		AllowSubmitWhenLessAuditLevel: template.AllowSubmitWhenLessAuditLevel,
		UpdateTime:                    template.UpdatedAt,
	}
	stepsRes := make([]*WorkFlowStepTemplateResV1, 0, len(template.Steps))
	for _, step := range template.Steps {
		stepRes := &WorkFlowStepTemplateResV1{
			Number:               int(step.Number),
			ApprovedByAuthorized: step.ApprovedByAuthorized.Bool,
			ExecuteByAuthorized:  step.ExecuteByAuthorized.Bool,
			Typ:                  step.Typ,
			Desc:                 step.Desc,
		}
		stepRes.Users = make([]string, 0)
		if step.Users != "" {
			stepRes.Users = strings.Split(step.Users, ",")
		}
		stepsRes = append(stepsRes, stepRes)
	}
	res.Steps = stepsRes

	// instanceNames, err := s.GetInstanceNamesByWorkflowTemplateId(template.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// res.Instances = instanceNames
	return res
}

type WorkFlowStepTemplateReqV1 struct {
	Type                 string   `json:"type" form:"type" valid:"oneof=sql_review sql_execute" enums:"sql_review,sql_execute"`
	Desc                 string   `json:"desc" form:"desc"`
	ApprovedByAuthorized bool     `json:"approved_by_authorized"`
	ExecuteByAuthorized  bool     `json:"execute_by_authorized"`
	Users                []string `json:"assignee_user_id_list" form:"assignee_user_id_list"`
}

type UpdateWorkflowTemplateReqV1 struct {
	Desc                          *string                      `json:"desc" form:"desc"`
	AllowSubmitWhenLessAuditLevel *string                      `json:"allow_submit_when_less_audit_level" enums:"normal,notice,warn,error"`
	Steps                         []*WorkFlowStepTemplateReqV1 `json:"workflow_step_template_list" form:"workflow_step_template_list"`
}

// @Summary 更新Sql审批流程模板
// @Description update the workflow template
// @Tags workflow
// @Id updateWorkflowTemplateV1
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param project_name path string true "project name"
// @Param instance body v1.UpdateWorkflowTemplateReqV1 true "create workflow template"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflow_template [patch]
func UpdateWorkflowTemplate(c echo.Context) error {
	return updateWorkflowTemplate(c)
}

type WorkflowStepResV1 struct {
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

// @Deprecated
// @Summary 审批通过
// @Description approve workflow
// @Tags workflow
// @Id approveWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_name path string true "workflow name"
// @Param workflow_step_id path string true "workflow step id"
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/{workflow_name}/steps/{workflow_step_id}/approve [post]
func ApproveWorkflow(c echo.Context) error {
	return nil
}

type RejectWorkflowReqV1 struct {
	Reason string `json:"reason" form:"reason"`
}

// @Deprecated
// @Summary 审批驳回
// @Description reject workflow
// @Tags workflow
// @Id rejectWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_name path string true "workflow name"
// @Param project_name path string true "project name"
// @Param workflow_step_id path string true "workflow step id"
// @param workflow_approve body v1.RejectWorkflowReqV1 true "workflow approve request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/{workflow_name}/steps/{workflow_step_id}/reject [post]
func RejectWorkflow(c echo.Context) error {
	return nil
}

// @Deprecated
// @Summary 审批关闭（中止）
// @Description cancel workflow
// @Tags workflow
// @Id cancelWorkflowV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param workflow_name path string true "workflow name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/{workflow_name}/cancel [post]
func CancelWorkflow(c echo.Context) error {
	return nil
}

type BatchCancelWorkflowsReqV1 struct {
	WorkflowNames []string `json:"workflow_names" form:"workflow_names"`
}

// BatchCancelWorkflows batch cancel workflows.
// @Deprecated
// @Summary 批量取消工单
// @Description batch cancel workflows
// @Tags workflow
// @Id batchCancelWorkflowsV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param BatchCancelWorkflowsReqV1 body v1.BatchCancelWorkflowsReqV1 true "batch cancel workflows request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/cancel [post]
func BatchCancelWorkflows(c echo.Context) error {
	return nil
}

type BatchCompleteWorkflowsReqV1 struct {
	WorkflowNames []string `json:"workflow_names" form:"workflow_names"`
}

// BatchCompleteWorkflows complete workflows.
// @Deprecated
// @Summary 批量完成工单
// @Description this api will directly change the work order status to finished without real online operation
// @Tags workflow
// @Id batchCompleteWorkflowsV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param data body v1.BatchCompleteWorkflowsReqV1 true "batch complete workflows request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/complete [post]
func BatchCompleteWorkflows(c echo.Context) error {
	return nil

}

func FormatStringToInt(s string) (ret int, err error) {
	if s == "" {
		return 0, nil
	} else {
		ret, err = strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
	}
	return ret, nil
}

func FormatStringToUint64(s string) (ret uint64, err error) {
	if s == "" {
		return 0, nil
	} else {
		ret, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
	}
	return ret, nil
}

// ExecuteOneTaskOnWorkflowV1
// @Deprecated
// @Summary 工单提交单个数据源上线
// @Description execute one task on workflow
// @Tags workflow
// @Id executeOneTaskOnWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_name path string true "workflow name"
// @Param project_name path string true "project name"
// @Param task_id path string true "task id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/{workflow_name}/tasks/{task_id}/execute [post]
func ExecuteOneTaskOnWorkflowV1(c echo.Context) error {
	return nil
}

func IsTaskCanExecute(s *model.Storage, taskId string) (bool, error) {
	task, err := getTaskById(context.Background(), taskId)
	if err != nil {
		return false, fmt.Errorf("get task by id failed. taskId=%v err=%v", taskId, err)
	}

	if task.Instance == nil {
		return false, fmt.Errorf("task instance is nil. taskId=%v", taskId)
	}

	inst := task.Instance
	if len(inst.MaintenancePeriod) > 0 && !inst.MaintenancePeriod.IsWithinScope(time.Now()) {
		return false, nil
	}

	instanceRecord, err := s.GetWorkInstanceRecordByTaskId(taskId)
	if err != nil {
		return false, fmt.Errorf("get work instance record by task id failed. taskId=%v err=%v", taskId, err)
	}

	if instanceRecord.ScheduledAt != nil || instanceRecord.IsSQLExecuted {
		return false, nil
	}

	return true, nil
}

func GetNeedExecTaskIds(workflow *model.Workflow, user *model.User) (taskIds map[uint] /*task id*/ string /*user id*/, err error) {
	instances := make([]*model.Instance, 0, len(workflow.Record.InstanceRecords))
	for _, item := range workflow.Record.InstanceRecords {
		instances = append(instances, item.Instance)
	}

	// 有不在运维时间内的instances报错
	var cannotExecuteInstanceNames []string
	for _, inst := range instances {
		if len(inst.MaintenancePeriod) != 0 && !inst.MaintenancePeriod.IsWithinScope(time.Now()) {
			cannotExecuteInstanceNames = append(cannotExecuteInstanceNames, inst.Name)
		}
	}
	if len(cannotExecuteInstanceNames) > 0 {
		return nil, errors.New(errors.TaskActionInvalid,
			fmt.Errorf("please go online during instance operation and maintenance time. these instances are not in maintenance time[%v]", strings.Join(cannotExecuteInstanceNames, ",")))
	}

	// 定时的instances和已上线的跳过
	needExecTaskIds := make(map[uint]string)
	for _, instRecord := range workflow.Record.InstanceRecords {
		if instRecord.ScheduledAt != nil || instRecord.IsSQLExecuted {
			continue
		}
		needExecTaskIds[instRecord.TaskId] = user.GetIDStr()
	}
	return needExecTaskIds, nil
}

func PrepareForWorkflowExecution(c echo.Context, projectUid string, workflow *model.Workflow, user *model.User) error {
	err := CheckCurrentUserCanOperateWorkflow(c, projectUid, workflow, []dmsV1.OpPermissionType{})
	if err != nil {
		return err
	}

	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return errors.New(errors.DataInvalid, fmt.Errorf("workflow current step not found"))
	}

	if workflow.Record.Status != model.WorkflowStatusWaitForExecution {
		return errors.New(errors.DataInvalid,
			fmt.Errorf("workflow need to be approved first"))
	}

	err = server.CheckUserCanOperateStep(user, workflow, int(currentStep.ID))
	if err != nil {
		return errors.New(errors.DataInvalid, err)
	}
	return nil
}

func PrepareForTaskExecution(c echo.Context, projectID string, workflow *model.Workflow, user *model.User, TaskId int) error {
	if workflow.Record.Status != model.WorkflowStatusWaitForExecution {
		return errors.New(errors.DataInvalid, e.New("workflow need to be approved first"))
	}

	err := CheckCurrentUserCanOperateTasks(c, projectID, workflow, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeExecuteWorkflow}, []uint{uint(TaskId)})
	if err != nil {
		return err
	}

	for _, record := range workflow.Record.InstanceRecords {
		if record.TaskId != uint(TaskId) {
			continue
		}

		for _, u := range strings.Split(record.ExecutionAssignees, ",") {
			if u == user.GetIDStr() {
				return nil
			}
		}
	}

	return e.New("you are not allow to execute the task")
}

type GetWorkflowTasksResV1 struct {
	controller.BaseRes
	Data []*GetWorkflowTasksItemV1 `json:"data"`
}

type GetWorkflowTasksItemV1 struct {
	TaskId                   uint                    `json:"task_id"`
	InstanceName             string                  `json:"instance_name"`
	Status                   string                  `json:"status" enums:"wait_for_audit,wait_for_execution,exec_scheduled,exec_failed,exec_succeeded,executing,manually_executed"`
	ExecStartTime            *time.Time              `json:"exec_start_time,omitempty"`
	ExecEndTime              *time.Time              `json:"exec_end_time,omitempty"`
	ScheduleTime             *time.Time              `json:"schedule_time,omitempty"`
	CurrentStepAssigneeUser  []string                `json:"current_step_assignee_user_name_list,omitempty"`
	TaskPassRate             float64                 `json:"task_pass_rate"`
	TaskScore                int32                   `json:"task_score"`
	InstanceMaintenanceTimes []*MaintenanceTimeResV1 `json:"instance_maintenance_times"`
	ExecutionUserName        string                  `json:"execution_user_name"`
}

// GetSummaryOfWorkflowTasksV1
// @Deprecated
// @Summary 获取工单数据源任务概览
// @Description get summary of workflow instance tasks
// @Tags workflow
// @Id getSummaryOfInstanceTasksV1
// @Security ApiKeyAuth
// @Param workflow_name path string true "workflow name"
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetWorkflowTasksResV1
// @router /v1/projects/{project_name}/workflows/{workflow_name}/tasks [get]
func GetSummaryOfWorkflowTasksV1(c echo.Context) error {
	return nil
}

const (
	taskDisplayStatusWaitForAudit       = "wait_for_audit"
	taskDisplayStatusWaitForExecution   = "wait_for_execution"
	taskDisplayStatusExecFailed         = "exec_failed"
	taskDisplayStatusExecSucceeded      = "exec_succeeded"
	taskStatusManuallyExecuted          = "manually_executed"
	taskDisplayStatusExecuting          = "executing"
	taskDisplayStatusScheduled          = "exec_scheduled"
	taskDisplayStatusTerminating        = "terminating"
	taskDisplayStatusTerminateSucceeded = "terminate_succeeded"
	taskDisplayStatusTerminateFailed    = "terminate_failed"
)

func GetTaskStatusRes(workflowStatus string, taskStatus string, scheduleAt *time.Time) (status string) {
	if workflowStatus == model.WorkflowStatusWaitForAudit {
		return taskDisplayStatusWaitForAudit
	}

	if scheduleAt != nil && taskStatus == model.TaskStatusAudited {
		return taskDisplayStatusScheduled
	}

	switch taskStatus {
	case model.TaskStatusAudited:
		return taskDisplayStatusWaitForExecution
	case model.TaskStatusExecuteSucceeded:
		return taskDisplayStatusExecSucceeded
	case model.TaskStatusExecuteFailed:
		return taskDisplayStatusExecFailed
	case model.TaskStatusExecuting:
		return taskDisplayStatusExecuting
	case model.TaskStatusManuallyExecuted:
		return taskStatusManuallyExecuted
	case model.TaskStatusTerminating:
		return taskDisplayStatusTerminating
	case model.TaskStatusTerminateSucc:
		return taskDisplayStatusTerminateSucceeded
	case model.TaskStatusTerminateFail:
		return taskDisplayStatusTerminateFailed
	}
	return ""
}

type CreateWorkflowReqV1 struct {
	Subject string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc    string `json:"desc" form:"desc"`
	TaskIds []uint `json:"task_ids" form:"task_ids" valid:"required"`
}

// CreateWorkflowV1
// @Deprecated
// @Summary 创建工单
// @Description create workflow
// @Accept json
// @Produce json
// @Tags workflow
// @Id createWorkflowV1
// @Security ApiKeyAuth
// @Param instance body v1.CreateWorkflowReqV1 true "create workflow request"
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows [post]
func CreateWorkflowV1(c echo.Context) error {
	return nil
}

func CheckWorkflowCanCommit(template *model.WorkflowTemplate, tasks []*model.Task) error {
	allowLevel := driverV2.RuleLevelError
	if template.AllowSubmitWhenLessAuditLevel != "" {
		allowLevel = driverV2.RuleLevel(template.AllowSubmitWhenLessAuditLevel)
	}
	for _, task := range tasks {
		if driverV2.RuleLevel(task.AuditLevel).More(allowLevel) {
			return errors.New(errors.DataInvalid,
				fmt.Errorf("there is an audit result with an error level higher than the allowable submission level(%v), please modify it before submitting. taskId=%v", allowLevel, task.ID))
		}
	}
	return nil
}

type GetWorkflowsReqV1 struct {
	FilterSubject                   string                `json:"filter_subject" query:"filter_subject"`
	FilterWorkflowID                string                `json:"filter_workflow_id" query:"filter_workflow_id"`
	FilterCreateTimeFrom            string                `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo              string                `json:"filter_create_time_to" query:"filter_create_time_to"`
	FilterCreateUserId              string                `json:"filter_create_user_id" query:"filter_create_user_id"`
	FilterStatus                    string                `json:"filter_status" query:"filter_status" valid:"omitempty,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterStatusList                []string              `json:"filter_status_list" query:"filter_status_list" validate:"dive,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterCurrentStepAssigneeUserId string                `json:"filter_current_step_assignee_user_id" query:"filter_current_step_assignee_user_id"`
	FilterTaskInstanceId            string                `json:"filter_task_instance_id" query:"filter_task_instance_id"`
	FilterTaskExecuteStartTimeFrom  string                `json:"filter_task_execute_start_time_from" query:"filter_task_execute_start_time_from"`
	FilterTaskExecuteStartTimeTo    string                `json:"filter_task_execute_start_time_to" query:"filter_task_execute_start_time_to"`
	FilterSqlVersionID              *uint                 `json:"filter_sql_version_id" query:"filter_sql_version_id"`
	FilterProjectUid                string                `json:"filter_project_uid" query:"filter_project_uid"`
	FilterInstanceId                string                `json:"filter_instance_id" query:"filter_instance_id"`
	FilterProjectPriority           dmsV1.ProjectPriority `json:"filter_project_priority" query:"filter_project_priority"  valid:"omitempty,oneof=high medium low"`
	PageIndex                       uint32                `json:"page_index" query:"page_index" valid:"required"`
	PageSize                        uint32                `json:"page_size" query:"page_size" valid:"required"`
	FuzzyKeyword                    string                `json:"fuzzy_keyword" query:"fuzzy_keyword"`
}

type GetWorkflowsResV1 struct {
	controller.BaseRes
	Data      []*WorkflowDetailResV1 `json:"data"`
	TotalNums uint64                 `json:"total_nums"`
}

type WorkflowDetailResV1 struct {
	ProjectName             string                `json:"project_name"`
	ProjectUid              string                `json:"project_uid,omitempty"`
	ProjectPriority         dmsV1.ProjectPriority `json:"project_priority"`
	Name                    string                `json:"workflow_name"`
	WorkflowId              string                `json:"workflow_id" `
	Desc                    string                `json:"desc"`
	SqlVersionName          []string              `json:"sql_version_name,omitempty"`
	CreateUser              string                `json:"create_user_name"`
	CreateTime              *time.Time            `json:"create_time"`
	CurrentStepType         string                `json:"current_step_type,omitempty" enums:"sql_review,sql_execute"`
	CurrentStepAssigneeUser []string              `json:"current_step_assignee_user_name_list,omitempty"`
	Status                  string                `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,executing,finished"`
	InstanceInfo            []InstanceInfo        `json:"instance_info,omitempty"`
}

type InstanceInfo struct {
	InstanceId   string `json:"instance_id,omitempty"`
	InstanceName string `json:"instance_name,omitempty"`
}

// GetGlobalWorkflowsV1
// @Summary 获取全局工单列表
// @Description get global workflow list
// @Tags workflow
// @Id getGlobalWorkflowsV1
// @Security ApiKeyAuth
// @Param filter_create_user_id query string false "filter create user id"
// @Param filter_status_list query []string false "filter by workflow status,, support using many status" Enums(wait_for_audit,wait_for_execution,rejected,executing,canceled,exec_failed,finished)
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id in project"
// @Param filter_project_priority query string false "filter by project priority" Enums(high,medium,low)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetWorkflowsResV1
// @router /v1/workflows [get]
func GetGlobalWorkflowsV1(c echo.Context) error {
	req := new(GetWorkflowsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 1. 获取用户权限信息
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	permissions, isAdmin, err := dmsobject.GetUserOpPermission(c.Request().Context(), "", user.GetIDStr(), dms.GetDMSServerAddress())
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 2. 将用户权限信息，转化为全局待处理清单统一的用户可视范围
	userVisibility := getGlobalDashBoardVisibilityOfUser(isAdmin, permissions)

	// 3. 将用户可视范围、接口请求以及用户的权限范围，构造为全局工单的基础的过滤器，满足全局工单统一的过滤逻辑
	filter, err := constructGlobalWorkflowBasicFilter(c.Request().Context(), user, userVisibility,
		&globalWorkflowBasicFilter{
			FilterCreateUserId:    req.FilterCreateUserId,
			FilterStatusList:      req.FilterStatusList,
			FilterProjectUid:      req.FilterProjectUid,
			FilterInstanceId:      req.FilterInstanceId,
			FilterProjectPriority: req.FilterProjectPriority,
		})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 4. 过滤器增加分页
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	filter["limit"] = limit
	filter["offset"] = offset
	// 5. 根据筛选项筛选工单信息
	s := model.GetStorage()
	workflows, count, err := s.GetWorkflowsByReq(filter)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if len(workflows) == 0 {
		return c.JSON(http.StatusOK, GetWorkflowsResV1{
			BaseRes:   controller.NewBaseReq(nil),
			Data:      []*WorkflowDetailResV1{},
			TotalNums: count,
		})
	}
	// 6. 从dms获取工单对应的项目信息
	var projectMap = make(ProjectMap)
	if req.FilterProjectPriority != "" {
		_, projectMap, err = loadProjectsByPriority(c.Request().Context(), req.FilterProjectPriority)
	} else {
		projectMap, err = loadProjectsByWorkflows(c.Request().Context(), workflows)
	}
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 7. 从dms获取工单对应的数据源信息
	instanceMap, err := loadInstanceByWorkflows(c.Request().Context(), workflows)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, GetWorkflowsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      toGlobalWorkflowRes(workflows, projectMap, instanceMap),
		TotalNums: count,
	})
}

type ProjectMap map[string] /* project uid */ *dmsV1.ListProject

func (m ProjectMap) ProjectName(projectUid string) string {
	if m == nil {
		return ""
	}
	if project, exist := m[projectUid]; exist && project != nil {
		return project.Name
	}
	return ""
}

func (m ProjectMap) ProjectPriority(projectUid string) dmsV1.ProjectPriority {
	if m == nil {
		return dmsV1.ProjectPriorityUnknown
	}
	if project, exist := m[projectUid]; exist && project != nil {
		return project.ProjectPriority
	}
	return dmsV1.ProjectPriorityUnknown
}

type InstanceMap map[string] /* instance id */ *dmsV1.ListDBService

func (m InstanceMap) InstanceName(instanceId string) string {
	if m == nil {
		return ""
	}
	if instance, exist := m[instanceId]; exist && instance != nil {
		return instance.Name
	}
	return ""
}

func toGlobalWorkflowRes(workflows []*model.WorkflowListDetail, projectMap ProjectMap, instanceMap InstanceMap) (workflowsResV1 []*WorkflowDetailResV1) {
	workflowsResV1 = make([]*WorkflowDetailResV1, 0, len(workflows))
	for _, workflow := range workflows {
		instanceInfos := make([]InstanceInfo, 0, len(workflow.InstanceIds))
		for _, id := range workflow.InstanceIds {
			instanceInfos = append(instanceInfos, InstanceInfo{
				InstanceId:   id,
				InstanceName: instanceMap.InstanceName(id),
			})
		}
		workflowRes := &WorkflowDetailResV1{
			ProjectName:             projectMap.ProjectName(workflow.ProjectId),
			ProjectUid:              workflow.ProjectId,
			ProjectPriority:         projectMap.ProjectPriority(workflow.ProjectId),
			InstanceInfo:            instanceInfos,
			Name:                    workflow.Subject,
			WorkflowId:              workflow.WorkflowId,
			Desc:                    workflow.Desc,
			CreateUser:              utils.AddDelTag(workflow.CreateUserDeletedAt, workflow.CreateUser.String),
			CreateTime:              workflow.CreateTime,
			CurrentStepType:         workflow.CurrentStepType.String,
			CurrentStepAssigneeUser: strings.Split(workflow.CurrentStepAssigneeUserIds.String, ","),
			Status:                  workflow.Status,
		}
		workflowsResV1 = append(workflowsResV1, workflowRes)
	}
	return workflowsResV1
}

type GetGlobalWorkflowStatisticsReqV1 struct {
	FilterCreateUserId    string                `json:"filter_create_user_id" query:"filter_create_user_id"`
	FilterStatusList      []string              `json:"filter_status_list" query:"filter_status_list" validate:"dive,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterProjectUid      string                `json:"filter_project_uid" query:"filter_project_uid"`
	FilterInstanceId      string                `json:"filter_instance_id" query:"filter_instance_id"`
	FilterProjectPriority dmsV1.ProjectPriority `json:"filter_project_priority" query:"filter_project_priority"  valid:"omitempty,oneof=high medium low"`
}

type GlobalWorkflowStatisticsResV1 struct {
	controller.BaseRes
	TotalNums uint64 `json:"total_nums"`
}

// GetGlobalWorkflowStatistics
// @Summary 获取全局工单统计数据
// @Description get global workflow statistics
// @Tags workflow
// @Id GetGlobalWorkflowStatistics
// @Security ApiKeyAuth
// @Param filter_create_user_id query string false "filter create user id"
// @Param filter_status_list query []string false "filter by workflow status,, support using many status" Enums(wait_for_audit,wait_for_execution,rejected,executing,canceled,exec_failed,finished)
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id in project"
// @Param filter_project_priority query string false "filter by project priority" Enums(high,medium,low)
// @Success 200 {object} v1.GlobalWorkflowStatisticsResV1
// @router /v1/workflows/statistics [get]
func GetGlobalWorkflowStatistics(c echo.Context) error {
	req := new(GetGlobalWorkflowStatisticsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 1. 获取用户权限信息
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	permissions, isAdmin, err := dmsobject.GetUserOpPermission(c.Request().Context(), "", user.GetIDStr(), dms.GetDMSServerAddress())
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 2. 将用户权限信息，转化为全局待处理清单统一的用户可视范围
	userVisibility := getGlobalDashBoardVisibilityOfUser(isAdmin, permissions)

	// 3. 将用户可视范围、接口请求以及用户的权限范围，构造为全局工单的基础的过滤器，满足全局工单统一的过滤逻辑
	filter, err := constructGlobalWorkflowBasicFilter(c.Request().Context(), user, userVisibility,
		&globalWorkflowBasicFilter{
			FilterCreateUserId:    req.FilterCreateUserId,
			FilterStatusList:      req.FilterStatusList,
			FilterProjectUid:      req.FilterProjectUid,
			FilterInstanceId:      req.FilterInstanceId,
			FilterProjectPriority: req.FilterProjectPriority,
		})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 4. 根据筛选项获取工单数量
	s := model.GetStorage()
	count, err := s.GetGlobalWorkflowTotalNum(filter)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, GlobalWorkflowStatisticsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		TotalNums: count,
	})
}

type globalWorkflowBasicFilter struct {
	FilterCreateUserId    string                `json:"filter_create_user_id" query:"filter_create_user_id"`
	FilterStatusList      []string              `json:"filter_status_list" query:"filter_status_list" validate:"dive,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterProjectUid      string                `json:"filter_project_uid" query:"filter_project_uid"`
	FilterInstanceId      string                `json:"filter_instance_id" query:"filter_instance_id"`
	FilterProjectPriority dmsV1.ProjectPriority `json:"filter_project_priority" query:"filter_project_priority"  valid:"omitempty,oneof=high medium low"`
}

// 将用户可视范围、接口请求以及用户的权限范围，构造为全局工单的基础的过滤器，满足全局工单统一的过滤逻辑
func constructGlobalWorkflowBasicFilter(ctx context.Context, user *model.User, userVisibility GlobalDashBoardVisibility, req *globalWorkflowBasicFilter) (map[string]interface{}, error) {
	// 1. 基本筛选项
	data := map[string]interface{}{
		"filter_create_user_id": req.FilterCreateUserId, // 根据创建人ID筛选用户自己创建的工单
		"filter_status_list":    req.FilterStatusList,   // 根据SQL工单的状态筛选多个状态的工单
		"filter_project_id":     req.FilterProjectUid,   // 根据项目id筛选某些一个项目下的多个工单
		"filter_instance_id":    req.FilterInstanceId,   // 根据工单记录的数据源id，筛选包含该数据源的工单，多数据源情况下，一旦包含该数据源，则被选中
	}
	// 1.1 页面筛选项：如果根据项目优先级筛选，则先筛选出对应优先级下的项目
	var projectIdsByPriority []string
	var err error
	if req.FilterProjectPriority != "" {
		projectIdsByPriority, _, err = loadProjectsByPriority(ctx, req.FilterProjectPriority)
		if err != nil {
			return nil, err
		}
		data["filter_project_id_list"] = projectIdsByPriority
	}
	// 2. 发起的工单页面，根据当前用户筛选在所有项目下该用户创建的工单，因此将用户可视范围调整为全局
	if req.FilterCreateUserId != "" {
		userVisibility.VisibilityType = GlobalDashBoardVisibilityGlobal
	}
	// 3. 待处理工单页面，根据当前用户的可视范围筛选
	switch userVisibility.ViewType() {
	case GlobalDashBoardVisibilityProjects:
		// 3.1 当用户的可视范围为多项目，则根据项目id筛选
		if req.FilterProjectPriority != "" {
			// 若根据项目优先级筛选，则将可查看的项目和项目优先级筛选后的项目的集合取交集
			data["filter_project_id_list"] = utils.IntersectionStringSlice(projectIdsByPriority, userVisibility.ViewRange())
		} else {
			// 若不根据项目优先级筛选，则通过用户的有权限的项目进行筛选
			data["filter_project_id_list"] = userVisibility.ViewRange()
		}
	case GlobalDashBoardVisibilityAssignee:
		// 3.2 若用户可视范围为受让人，则查看分配给他的工单
		data["filter_current_step_assignee_user_id"] = user.GetIDStr()
	}
	return data, nil
}

type VisibilityType string

const GlobalDashBoardVisibilityGlobal VisibilityType = "global"     // 全局可见
const GlobalDashBoardVisibilityProjects VisibilityType = "projects" // 多项目可见
const GlobalDashBoardVisibilityAssignee VisibilityType = "assignee" // 仅可见授予自己的

type GlobalDashBoardVisibility struct {
	VisibilityType  VisibilityType
	VisibilityRange []string // 对于项目是项目id
}

func (v GlobalDashBoardVisibility) ViewType() VisibilityType {
	return v.VisibilityType
}

func (v GlobalDashBoardVisibility) ViewRange() []string {
	return v.VisibilityRange
}

// 将用户权限信息，转化为全局待处理清单统一的用户可视范围
func getGlobalDashBoardVisibilityOfUser(isAdmin bool, permissions []dmsV1.OpPermissionItem) GlobalDashBoardVisibility {
	// 角色：全局管理员，全局可查看者
	if isAdmin {
		return GlobalDashBoardVisibility{
			VisibilityType: GlobalDashBoardVisibilityGlobal,
		}
	}
	for _, permission := range permissions {
		if permission.OpPermissionType == dmsV1.OpPermissionTypeGlobalView || permission.OpPermissionType == dmsV1.OpPermissionTypeGlobalManagement {
			return GlobalDashBoardVisibility{
				VisibilityType: GlobalDashBoardVisibilityGlobal,
			}
		}
	}
	// 角色：多项目管理者
	var projectRange []string
	for _, permission := range permissions {
		if permission.OpPermissionType == dmsV1.OpPermissionTypeProjectAdmin {
			projectRange = append(projectRange, permission.RangeUids...)
		}
	}
	if len(projectRange) > 0 {
		return GlobalDashBoardVisibility{
			VisibilityType:  GlobalDashBoardVisibilityProjects,
			VisibilityRange: projectRange,
		}
	}
	// 角色：受让人，事件处理者
	return GlobalDashBoardVisibility{
		VisibilityType: GlobalDashBoardVisibilityAssignee,
	}
}

// 根据项目优先级从 dms 系统中获取相应的项目列表，并返回项目ID列表和项目映射
func loadProjectsByPriority(ctx context.Context, priority dmsV1.ProjectPriority) (projectIds []string, projectMap ProjectMap, err error) {
	projectMap = make(ProjectMap)
	// 如果根据项目优先级筛选SQL工单，则先获取项目优先级，根据优先级对应的项目ID进行筛选
	projects, _, err := dmsobject.ListProjects(ctx, controller.GetDMSServerAddress(), dmsV1.ListProjectReq{
		PageSize:                999,
		PageIndex:               1,
		FilterByProjectPriority: priority,
	})
	if err != nil {
		return nil, nil, err
	}
	for _, project := range projects {
		if _, exist := projectMap[project.ProjectUid]; !exist {
			projectMap[project.ProjectUid] = project
		}
		projectIds = append(projectIds, project.ProjectUid)
	}
	return projectIds, projectMap, nil
}

// 根据工单列表中的项目ID从 dms 系统中获取对应的项目信息，并返回项目映射
func loadProjectsByWorkflows(ctx context.Context, workflows []*model.WorkflowListDetail) (projectMap ProjectMap, err error) {
	projectMap = make(ProjectMap)
	if len(workflows) == 0 {
		return projectMap, nil
	}

	var projectIds []string
	for _, workflow := range workflows {
		if _, exist := projectMap[workflow.ProjectId]; !exist {
			projectIds = append(projectIds, workflow.ProjectId)
			projectMap[workflow.ProjectId] = nil
		}
	}
	return loadProjectsByProjectIds(ctx, projectIds)
}

func loadProjectsByProjectIds(ctx context.Context, projectIds []string) (projectMap ProjectMap, err error) {
	// get project priority from dms
	projects, _, err := dmsobject.ListProjects(ctx, controller.GetDMSServerAddress(), dmsV1.ListProjectReq{
		PageSize:            uint32(len(projectIds)),
		PageIndex:           1,
		FilterByProjectUids: projectIds,
	})
	if err != nil {
		return nil, err
	}
	projectMap = make(map[string] /* project uid */ *dmsV1.ListProject)
	for _, project := range projects {
		projectMap[project.ProjectUid] = project
	}
	return projectMap, nil
}

// 根据工单列表中的实例ID从 dms 系统中获取对应的数据源实例信息，并返回实例映射
func loadInstanceByWorkflows(ctx context.Context, workflows []*model.WorkflowListDetail) (instanceMap InstanceMap, err error) {
	instanceMap = make(InstanceMap)
	if len(workflows) == 0 {
		return instanceMap, nil
	}

	var instanceIdList []string
	for _, workflow := range workflows {
		for _, id := range workflow.InstanceIds {
			if _, exist := instanceMap[id]; !exist {
				instanceIdList = append(instanceIdList, id)
				instanceMap[id] = nil
			}
		}
	}
	return loadInstanceByInstanceIds(ctx, instanceIdList)
}

func loadInstanceByInstanceIds(ctx context.Context, instanceIds []string) (instanceMap InstanceMap, err error) {
	// get instances from dms
	instanceMap = make(InstanceMap)
	instances, _, err := dmsobject.ListDbServices(ctx, controller.GetDMSServerAddress(), dmsV1.ListDBServiceReq{
		PageSize:             uint32(len(instanceIds)),
		PageIndex:            1,
		FilterByDBServiceIds: instanceIds,
	})
	if err != nil {
		return nil, err
	}
	for _, instance := range instances {
		instanceMap[instance.DBServiceUid] = instance
	}
	return instanceMap, nil
}

// GetWorkflowsV1
// @Summary 获取工单列表
// @Description get workflow list
// @Tags workflow
// @Id getWorkflowsV1
// @Security ApiKeyAuth
// @Param filter_subject query string false "filter subject"
// @Param filter_workflow_id query string false "filter by workflow_id"
// @Param fuzzy_search_workflow_desc query string false "fuzzy search by workflow description"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param filter_task_execute_start_time_from query string false "filter_task_execute_start_time_from"
// @Param filter_task_execute_start_time_to query string false "filter_task_execute_start_time_to"
// @Param filter_create_user_id query string false "filter create user id"
// @Param filter_status query string false "filter workflow status" Enums(wait_for_audit,wait_for_execution,rejected,executing,canceled,exec_failed,finished)
// @Param filter_current_step_assignee_user_id query string false "filter current step assignee user id"
// @Param filter_task_instance_id query string false "filter instance id"
// @Param filter_sql_version_id query string false "filter sql version id"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Param project_name path string true "project name"
// @Param fuzzy_keyword query string false "fuzzy matching subject/workflow_id"
// @Success 200 {object} v1.GetWorkflowsResV1
// @router /v1/projects/{project_name}/workflows [get]
func GetWorkflowsV1(c echo.Context) error {
	req := new(GetWorkflowsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	up, err := dms.NewUserPermission(user.GetIDStr(), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	data := map[string]interface{}{
		"filter_workflow_id":                   req.FilterWorkflowID,
		"filter_sql_version_id":                req.FilterSqlVersionID,
		"filter_subject":                       req.FilterSubject,
		"filter_create_time_from":              req.FilterCreateTimeFrom,
		"filter_create_time_to":                req.FilterCreateTimeTo,
		"filter_create_user_id":                req.FilterCreateUserId,
		"filter_task_execute_start_time_from":  req.FilterTaskExecuteStartTimeFrom,
		"filter_task_execute_start_time_to":    req.FilterTaskExecuteStartTimeTo,
		"filter_status":                        req.FilterStatus,
		"filter_current_step_assignee_user_id": req.FilterCurrentStepAssigneeUserId,
		"filter_task_instance_id":              req.FilterTaskInstanceId,
		"filter_project_id":                    projectUid,
		"current_user_id":                      user.ID,
		"check_user_can_access":                !up.CanViewProject(),
		"limit":                                limit,
		"offset":                               offset,
	}
	if req.FuzzyKeyword != "" {
		data["fuzzy_keyword"] = fmt.Sprintf("%%%s%%", req.FuzzyKeyword)
	}

	if !up.CanViewProject() {
		data["viewable_instance_ids"] = strings.Join(up.GetInstancesByOP(dmsV1.OpPermissionTypeViewOthersWorkflow), ",")
	}

	workflows, count, err := s.GetWorkflowsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowsResV1 := make([]*WorkflowDetailResV1, 0, len(workflows))
	for _, workflow := range workflows {
		// TODO DMS提供根据ID批量查询用户接口，demo阶段使用GetUser实现
		CurrentStepAssigneeUserNames := make([]string, 0)
		for _, currentStepAssigneeUser := range strings.Split(workflow.CurrentStepAssigneeUserIds.String, ",") {
			if currentStepAssigneeUser == "" {
				continue
			}
			CurrentStepAssigneeUserNames = append(CurrentStepAssigneeUserNames, dms.GetUserNameWithDelTag(currentStepAssigneeUser))
		}
		workflowRes := &WorkflowDetailResV1{
			ProjectName:             workflow.ProjectId, // dms-todo: 暂时使用id代替name
			Name:                    workflow.Subject,
			WorkflowId:              workflow.WorkflowId,
			Desc:                    workflow.Desc,
			CreateUser:              dms.GetUserNameWithDelTag(workflow.CreateUser.String),
			CreateTime:              workflow.CreateTime,
			CurrentStepType:         workflow.CurrentStepType.String,
			CurrentStepAssigneeUser: CurrentStepAssigneeUserNames,
			Status:                  workflow.Status,
			SqlVersionName:          workflow.SqlVersionNames,
		}
		workflowsResV1 = append(workflowsResV1, workflowRes)
	}

	return c.JSON(http.StatusOK, GetWorkflowsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      workflowsResV1,
		TotalNums: count,
	})
}

type GetWorkflowStatisticOfInstancesResV1 struct {
	controller.BaseRes
	Data []*WorkflowStatisticOfInstance `json:"data"`
}

type WorkflowStatisticOfInstance struct {
	InstanceId      int64 `json:"instance_id"`
	UnfinishedCount int64 `json:"unfinished_count"`
}

// GetWorkflowStatisticOfInstances
// @Summary 获取实例上工单的统计信息
// @Description Get Workflows Statistic Of Instances
// @Tags workflow
// @Id GetWorkflowStatisticOfInstances
// @Security ApiKeyAuth
// @Param instance_id query string true "instance id"
// @Success 200 {object} v1.GetWorkflowStatisticOfInstancesResV1
// @router /v1/workflows/statistic_of_instances [get]
func GetWorkflowStatisticOfInstances(c echo.Context) error {
	instanceIds := c.QueryParams()["instance_id"]
	if len(instanceIds) == 0 {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("query param instance_id requied"))
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if user.Name != model.DefaultAdminUser && user.Name != model.DefaultSysUser {
		// dms-todo: 后续需要通过dms来判断权限，没这么做的原因是：目前dms接口获取用户权限时必需projectUid参数
		return controller.JSONBaseErrorReq(c, fmt.Errorf("permission denied"))
	}

	s := model.GetStorage()
	unfinishedStatuses := []string{model.WorkflowStatusWaitForAudit, model.WorkflowStatusWaitForExecution, model.WorkflowStatusExecuting}
	results, err := s.GetWorkflowStatusesCountOfInstances(unfinishedStatuses, instanceIds)
	if err != nil {
		return err
	}

	workflowsInfoV1 := make([]*WorkflowStatisticOfInstance, len(results))
	for k, v := range results {
		workflowsInfoV1[k] = &WorkflowStatisticOfInstance{
			InstanceId:      v.InstanceId,
			UnfinishedCount: v.Count,
		}
	}

	return c.JSON(http.StatusOK, GetWorkflowStatisticOfInstancesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    workflowsInfoV1,
	})
}

type UpdateWorkflowReqV1 struct {
	TaskIds []uint `json:"task_ids" form:"task_ids" valid:"required"`
}

// UpdateWorkflowV1
// @Deprecated
// @Summary 更新工单（驳回后才可更新）
// @Description update workflow when it is rejected to creator.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_name path string true "workflow name"
// @Param project_name path string true "project name"
// @Param instance body v1.UpdateWorkflowReqV1 true "update workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/{workflow_name}/ [patch]
func UpdateWorkflowV1(c echo.Context) error {
	return nil
}

type UpdateWorkflowScheduleReqV1 struct {
	ScheduleTime *time.Time `json:"schedule_time"`
}

// UpdateWorkflowScheduleV1
// @Deprecated
// @Summary 设置工单数据源定时上线时间（设置为空则代表取消定时时间，需要SQL审核流程都通过后才可以设置）
// @Description update workflow schedule.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowScheduleV1
// @Security ApiKeyAuth
// @Param workflow_name path string true "workflow name"
// @Param task_id path string true "task id"
// @Param project_name path string true "project name"
// @Param instance body v1.UpdateWorkflowScheduleReqV1 true "update workflow schedule request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/{workflow_name}/tasks/{task_id}/schedule [put]
func UpdateWorkflowScheduleV1(c echo.Context) error {
	return nil
}

// ExecuteTasksOnWorkflowV1
// @Deprecated
// @Summary 多数据源批量上线
// @Description execute tasks on workflow
// @Tags workflow
// @Id executeTasksOnWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_name path string true "workflow name"
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/{workflow_name}/tasks/execute [post]
func ExecuteTasksOnWorkflowV1(c echo.Context) error {
	return nil
}

type GetWorkflowResV1 struct {
	controller.BaseRes
	Data *WorkflowResV1 `json:"data"`
}

type WorkflowTaskItem struct {
	Id uint `json:"task_id"`
}

type WorkflowRecordResV1 struct {
	Tasks             []*WorkflowTaskItem  `json:"tasks"`
	CurrentStepNumber uint                 `json:"current_step_number,omitempty"`
	Status            string               `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,executing,finished"`
	Steps             []*WorkflowStepResV1 `json:"workflow_step_list,omitempty"`
}

type WorkflowResV1 struct {
	Name          string                 `json:"workflow_name"`
	Desc          string                 `json:"desc,omitempty"`
	Mode          string                 `json:"mode" enums:"same_sqls,different_sqls"`
	CreateUser    string                 `json:"create_user_name"`
	CreateTime    *time.Time             `json:"create_time"`
	Record        *WorkflowRecordResV1   `json:"record"`
	RecordHistory []*WorkflowRecordResV1 `json:"record_history_list,omitempty"`
}

// GetWorkflowV1
// @Deprecated
// @Summary 获取工单详情
// @Description get workflow detail
// @Tags workflow
// @Id getWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_name path string true "workflow name"
// @Param project_name path string true "project name"
// @Success 200 {object} GetWorkflowResV1
// @router /v1/projects/{project_name}/workflows/{workflow_name}/ [get]
func GetWorkflowV1(c echo.Context) error {
	return nil
}

type ExportWorkflowReqV1 struct {
	FilterSubject                   string `json:"filter_subject" query:"filter_subject"`
	FilterCreateTimeFrom            string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo              string `json:"filter_create_time_to" query:"filter_create_time_to"`
	FilterCreateUserID              string `json:"filter_create_user_id" query:"filter_create_user_id"`
	FilterStatus                    string `json:"filter_status" query:"filter_status" valid:"omitempty,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterCurrentStepAssigneeUserId string `json:"filter_current_step_assignee_user_id" query:"filter_current_step_assignee_user_id"`
	FilterTaskInstanceId            string `json:"filter_task_instance_id" query:"filter_task_instance_id"`
	FilterTaskExecuteStartTimeFrom  string `json:"filter_task_execute_start_time_from" query:"filter_task_execute_start_time_from"`
	FilterTaskExecuteStartTimeTo    string `json:"filter_task_execute_start_time_to" query:"filter_task_execute_start_time_to"`
	FuzzyKeyword                    string `json:"fuzzy_keyword" query:"fuzzy_keyword"`
}

// ExportWorkflowV1
// @Summary 导出工单
// @Description export workflow
// @Id exportWorkflowV1
// @Tags workflow
// @Security ApiKeyAuth
// @Param filter_subject query string false "filter subject"
// @Param fuzzy_search_workflow_desc query string false "fuzzy search by workflow description"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param filter_task_execute_start_time_from query string false "filter_task_execute_start_time_from"
// @Param filter_task_execute_start_time_to query string false "filter_task_execute_start_time_to"
// @Param filter_create_user_id query string false "filter create user id"
// @Param filter_status query string false "filter workflow status" Enums(wait_for_audit,wait_for_execution,rejected,executing,canceled,exec_failed,finished)
// @Param filter_current_step_assignee_user_id query string false "filter current step assignee user id"
// @Param filter_task_instance_id query string false "filter instance id"
// @Param project_name path string true "project name"
// @Param fuzzy_keyword query string false "fuzzy matching subject/workflow_id/desc"
// @Success 200 {file} file "export workflow"
// @Router /v1/projects/{project_name}/workflows/exports [get]
func ExportWorkflowV1(c echo.Context) error {
	return exportWorkflowV1(c)
}

// TerminateMultipleTaskByWorkflowV1
// @Summary 终止工单下多个上线任务
// @Description terminate multiple task by project and workflow
// @Tags workflow
// @Id terminateMultipleTaskByWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @Router /v1/projects/{project_name}/workflows/{workflow_id}/tasks/terminate [post]
func TerminateMultipleTaskByWorkflowV1(c echo.Context) error {

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowID := c.Param("workflow_id")
	// user, err := controller.GetCurrentUser(c,dms.GetUser)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }
	s := model.GetStorage()

	var workflow *model.Workflow
	{
		workflow, err = dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	terminatingTaskIDs := getTerminatingTaskIDs(workflow)

	// check workflow permission
	{
		err := checkBeforeTasksTermination(c, projectUid, workflow, terminatingTaskIDs)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	err = s.UpdateTaskStatusByIDs(terminatingTaskIDs,
		map[string]interface{}{"status": model.TaskStatusTerminating})

	return c.JSON(http.StatusOK, controller.NewBaseReq(err))
}

// TerminateSingleTaskByWorkflowV1
// @Summary 终止单个上线任务
// @Description execute one task on workflow
// @Tags workflow
// @Id terminateSingleTaskByWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Param task_id path string true "task id"
// @Success 200 {object} controller.BaseRes
// @Router /v1/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/terminate [post]
func TerminateSingleTaskByWorkflowV1(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowID := c.Param("workflow_id")
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// user, err := controller.GetCurrentUser(c,dms.GetUser)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }
	s := model.GetStorage()

	var workflow *model.Workflow
	{
		workflow, err = dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	// check workflow permission
	{
		err := checkBeforeTasksTermination(c, projectUid, workflow, []uint{uint(taskID)})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	// check task
	{
		ok, err := isTaskCanBeTerminate(s, taskIDStr)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !ok {
			return controller.JSONBaseErrorReq(c,
				fmt.Errorf("task can not be terminated. taskId=%v workflowId=%v", taskID, workflowID))
		}
	}

	err = s.UpdateTaskStatusByIDs([]uint{uint(taskID)},
		map[string]interface{}{"status": model.TaskStatusTerminating})

	return c.JSON(http.StatusOK, controller.NewBaseReq(err))
}

func checkBeforeTasksTermination(c echo.Context, projectId string, workflow *model.Workflow, needTerminatedTaskIdList []uint) error {
	needTerminatedTaskIdMap := make(map[uint]struct{}, len(needTerminatedTaskIdList))
	for _, taskID := range needTerminatedTaskIdList {
		needTerminatedTaskIdMap[taskID] = struct{}{}
	}

	for _, record := range workflow.Record.InstanceRecords {
		if _, ok := needTerminatedTaskIdMap[record.TaskId]; !ok {
			continue
		}

		isWorkflowWaitForExecution := workflow.Record.Status == model.WorkflowStatusWaitForExecution
		isWorkflowExecuting := workflow.Record.Status == model.WorkflowStatusExecuting
		isTaskExecuting := record.Task.Status == model.TaskStatusExecuting

		if !(isWorkflowWaitForExecution || isWorkflowExecuting) {
			return ErrCanNotTerminateExecute(workflow.Record.Status, record.Task.Status)
		}

		if !isTaskExecuting {
			return ErrCanNotTerminateExecute(workflow.Record.Status, record.Task.Status)
		}

		return nil
	}

	err := CheckCurrentUserCanOperateTasks(c,
		projectId, workflow, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeViewOthersWorkflow}, needTerminatedTaskIdList)
	if err != nil {
		return err
	}

	return nil
}

func isTaskCanBeTerminate(s *model.Storage, taskID string) (bool, error) {
	task, err := getTaskById(context.Background(), taskID)
	if err != nil {
		return false, fmt.Errorf("get task by id failed. taskID=%v err=%v", taskID, err)
	}

	if task.Instance == nil {
		return false, fmt.Errorf("task instance is nil. taskID=%v", taskID)
	}

	if task.Status == model.TaskStatusExecuting {
		return true, nil
	}

	return false, nil
}

func getTerminatingTaskIDs(workflow *model.Workflow) (taskIDs []uint) {

	taskIDs = make([]uint, 0)
	for i := range workflow.Record.InstanceRecords {
		instRecord := workflow.Record.InstanceRecords[i]
		if instRecord.Task.Status == model.TaskStatusExecuting {
			taskIDs = append(taskIDs, instRecord.TaskId)
		}
	}
	return taskIDs
}

type FileToSort struct {
	FileID   uint `json:"file_id" valid:"required"`
	NewIndex uint `json:"new_index" valid:"required"`
}

type UpdateSqlFileOrderV1Req struct {
	FilesToSort []FileToSort `json:"files_to_sort"`
}

// UpdateSqlFileOrderV1
// @Summary 修改文件上线顺序
// @Description update sql file order
// @Accept json
// @Produce json
// @Tags task
// @Id updateSqlFileOrderV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param workflow_id path string true "workflow id"
// @Param task_id path string true "task id"
// @Param instance body v1.UpdateSqlFileOrderV1Req true "instance body v1.UpdateSqlFileOrderV1Req true"
// @Success 200 {object} v1.GetSqlFileOrderMethodResV1
// @router /v1/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/order_file [post]
func UpdateSqlFileOrderByWorkflowV1(c echo.Context) error {
	return updateSqlFileOrderByWorkflow(c)
}

// GetWorkflowAttachment
// @Summary 获取工单的task附件
// @Description get workflow attachment
// @Tags workflow
// @Id getWorkflowAttachment
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param workflow_id path string true "workflow id"
// @Param task_id path string true "task id"
// @Success 200 {file} file "get workflow attachment"
// @Router /v1/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/attachment [get]
func GetWorkflowTaskAuditFile(c echo.Context) error {
	taskId := c.Param("task_id")
	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	attachment, exist, err := s.GetParentFileByTaskId(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("can not find any file in this task"))
	}

	if attachment.FileHost != config.GetOptions().SqleOptions.ReportHost {
		log.NewEntry().Infof("try to reverse to sqle due to file.FileHost %v this host %v", attachment.FileHost, config.GetOptions().SqleOptions.ReportHost)
		err = ReverseToSqle(c, c.Request().URL.Path, attachment.FileHost)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	} else {
		filePath := model.DefaultFilePath(attachment.UniqueName)
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		c.Response().Header().Set(echo.HeaderContentDisposition, mime.FormatMediaType("inline", map[string]string{"filename": attachment.FileName}))
		err = c.Blob(http.StatusOK, echo.MIMEOctetStream, fileData)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	return c.NoContent(http.StatusOK)
}
