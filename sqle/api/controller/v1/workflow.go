package v1

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
)

var ErrWorkflowNoAccess = errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist or you can't access it"))

var ErrForbidMyBatisXMLTask = func(taskId uint) error {
	return errors.New(errors.DataConflict,
		fmt.Errorf("the task for audit mybatis xml file is not allow to create workflow. taskId=%v", taskId))
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
	Instances                     []string                     `json:"instance_name_list,omitempty"`
}

type WorkFlowStepTemplateResV1 struct {
	Number               int      `json:"number"`
	Typ                  string   `json:"type"`
	Desc                 string   `json:"desc,omitempty"`
	ApprovedByAuthorized bool     `json:"approved_by_authorized"`
	ExecuteByAuthorized  bool     `json:"execute_by_authorized"`
	Users                []string `json:"assignee_user_name_list"`
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
	s := model.GetStorage()

	projectName := c.Param("project_name")

	err := CheckIsProjectMember(controller.GetUserName(c), projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}

	template, exist, err := s.GetWorkflowTemplateById(project.WorkflowTemplateId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("workflow template is not exist")))
	}

	res, err := getWorkflowTemplateDetailByTemplate(template)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetWorkflowTemplateResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    res,
	})
}

func getWorkflowTemplateDetailByTemplate(template *model.WorkflowTemplate) (*WorkflowTemplateDetailResV1, error) {
	s := model.GetStorage()
	steps, err := s.GetWorkflowStepsDetailByTemplateId(template.ID)
	if err != nil {
		return nil, err
	}
	template.Steps = steps
	res := &WorkflowTemplateDetailResV1{
		Name:                          template.Name,
		Desc:                          template.Desc,
		AllowSubmitWhenLessAuditLevel: template.AllowSubmitWhenLessAuditLevel,
	}
	stepsRes := make([]*WorkFlowStepTemplateResV1, 0, len(steps))
	for _, step := range steps {
		stepRes := &WorkFlowStepTemplateResV1{
			Number:               int(step.Number),
			ApprovedByAuthorized: step.ApprovedByAuthorized.Bool,
			ExecuteByAuthorized:  step.ExecuteByAuthorized.Bool,
			Typ:                  step.Typ,
			Desc:                 step.Desc,
		}
		users := []string{}
		if step.Users != nil {
			for _, user := range step.Users {
				users = append(users, user.Name)
			}
		}
		stepRes.Users = users
		stepsRes = append(stepsRes, stepRes)
	}
	res.Steps = stepsRes

	instanceNames, err := s.GetInstanceNamesByWorkflowTemplateId(template.ID)
	if err != nil {
		return nil, err
	}
	res.Instances = instanceNames
	return res, nil
}

type WorkFlowStepTemplateReqV1 struct {
	Type                 string   `json:"type" form:"type" valid:"oneof=sql_review sql_execute" enums:"sql_review,sql_execute"`
	Desc                 string   `json:"desc" form:"desc"`
	ApprovedByAuthorized bool     `json:"approved_by_authorized"`
	ExecuteByAuthorized  bool     `json:"execute_by_authorized"`
	Users                []string `json:"assignee_user_name_list" form:"assignee_user_name_list"`
}

func validWorkflowTemplateReq(steps []*WorkFlowStepTemplateReqV1) error {
	if len(steps) == 0 {
		return fmt.Errorf("workflow steps cannot be empty")
	}
	if len(steps) > 5 {
		return fmt.Errorf("workflow steps length must be less than 6")
	}

	for i, step := range steps {
		isLastStep := i == len(steps)-1
		if isLastStep && step.Type != model.WorkflowStepTypeSQLExecute {
			return fmt.Errorf("the last workflow step type must be sql_execute")
		}
		if !isLastStep && step.Type == model.WorkflowStepTypeSQLExecute {
			return fmt.Errorf("workflow step type sql_execute just be used in last step")
		}
		if len(step.Users) == 0 && !step.ApprovedByAuthorized && !step.ExecuteByAuthorized {
			return fmt.Errorf("the assignee is empty for step %s", step.Desc)
		}
		if len(step.Users) > 3 {
			return fmt.Errorf("the assignee for step cannot be more than 3")
		}
	}
	return nil
}

type UpdateWorkflowTemplateReqV1 struct {
	Desc                          *string                      `json:"desc" form:"desc"`
	AllowSubmitWhenLessAuditLevel *string                      `json:"allow_submit_when_less_audit_level" enums:"normal,notice,warn,error"`
	Steps                         []*WorkFlowStepTemplateReqV1 `json:"workflow_step_template_list" form:"workflow_step_template_list"`
	Instances                     []string                     `json:"instance_name_list" form:"instance_name_list"`
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
	req := new(UpdateWorkflowTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")

	s := model.GetStorage()
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}
	if project.IsArchived() {
		return controller.JSONBaseErrorReq(c, ErrProjectArchived)
	}
	userName := controller.GetUserName(c)

	err = CheckIsProjectManager(userName, project.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowTemplate, exist, err := s.GetWorkflowTemplateById(project.WorkflowTemplateId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("workflow template is not exist")))
	}

	var instances []*model.Instance
	if req.Instances != nil && len(req.Instances) > 0 {
		instances, err = s.GetAndCheckInstanceExist(req.Instances, project.Name)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Steps != nil {
		err = validWorkflowTemplateReq(req.Steps)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
		}
		userNames := []string{}
		for _, step := range req.Steps {
			userNames = append(userNames, step.Users...)
		}

		users, err := s.GetAndCheckUserExist(userNames)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		userMap := map[string]*model.User{}
		for _, user := range users {
			userMap[user.Name] = user
		}

		steps := make([]*model.WorkflowStepTemplate, 0, len(req.Steps))
		for i, step := range req.Steps {
			s := &model.WorkflowStepTemplate{
				Number: uint(i + 1),
				ApprovedByAuthorized: sql.NullBool{
					Bool:  step.ApprovedByAuthorized,
					Valid: true,
				},
				ExecuteByAuthorized: sql.NullBool{
					Bool:  step.ExecuteByAuthorized,
					Valid: true,
				},
				Typ:  step.Type,
				Desc: step.Desc,
			}
			stepUsers := make([]*model.User, 0, len(step.Users))
			for _, userName := range step.Users {
				stepUsers = append(stepUsers, userMap[userName])
			}
			s.Users = stepUsers
			steps = append(steps, s)
		}
		err = s.UpdateWorkflowTemplateSteps(workflowTemplate.ID, steps)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Desc != nil {
		workflowTemplate.Desc = *req.Desc
	}

	if req.AllowSubmitWhenLessAuditLevel != nil {
		workflowTemplate.AllowSubmitWhenLessAuditLevel = *req.AllowSubmitWhenLessAuditLevel
	}

	err = s.Save(workflowTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.Instances != nil {
		err = s.UpdateWorkflowTemplateInstances(workflowTemplate, instances...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
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

func CheckUserCanOperateStep(user *model.User, workflow *model.Workflow, stepId int) error {
	if workflow.Record.Status != model.WorkflowStatusWaitForAudit && workflow.Record.Status != model.WorkflowStatusWaitForExecution {
		return fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)
	}
	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return fmt.Errorf("workflow current step not found")
	}
	if uint(stepId) != workflow.CurrentStep().ID {
		return fmt.Errorf("workflow current step is not %d", stepId)
	}

	if !workflow.IsOperationUser(user) {
		return fmt.Errorf("you are not allow to operate the workflow")
	}
	return nil
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
	task, exist, err := s.GetTaskById(taskId)
	if err != nil {
		return false, fmt.Errorf("get task by id failed. taskId=%v err=%v", taskId, err)
	}
	if !exist {
		return false, fmt.Errorf("task not exist. taskId=%v", taskId)
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

func GetNeedExecTaskIds(s *model.Storage, workflow *model.Workflow, user *model.User) (taskIds map[uint] /*task id*/ uint /*user id*/, err error) {
	instances, err := s.GetInstancesByWorkflowID(workflow.ID)
	if err != nil {
		return nil, err
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
	needExecTaskIds := make(map[uint]uint)
	for _, instRecord := range workflow.Record.InstanceRecords {
		if instRecord.ScheduledAt != nil || instRecord.IsSQLExecuted {
			continue
		}
		needExecTaskIds[instRecord.TaskId] = user.ID
	}
	return needExecTaskIds, nil
}

func PrepareForWorkflowExecution(c echo.Context, projectName string, workflow *model.Workflow, user *model.User) error {
	err := CheckCurrentUserCanOperateWorkflow(c, &model.Project{Name: projectName}, workflow, []uint{})
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

	err = CheckUserCanOperateStep(user, workflow, int(currentStep.ID))
	if err != nil {
		return errors.New(errors.DataInvalid, err)
	}
	return nil
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
	FilterSubject                     string `json:"filter_subject" query:"filter_subject"`
	FilterWorkflowID                  string `json:"filter_workflow_id" query:"filter_workflow_id"`
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

type GetWorkflowsResV1 struct {
	controller.BaseRes
	Data      []*WorkflowDetailResV1 `json:"data"`
	TotalNums uint64                 `json:"total_nums"`
}

type WorkflowDetailResV1 struct {
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

// GetGlobalWorkflowsV1
// @Summary 获取全局工单列表
// @Description get global workflow list
// @Tags workflow
// @Id getGlobalWorkflowsV1
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
// @Success 200 {object} v1.GetWorkflowsResV1
// @router /v1/workflows [get]
func GetGlobalWorkflowsV1(c echo.Context) error {
	req := new(GetWorkflowsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}

	data := map[string]interface{}{
		"filter_subject":                         req.FilterSubject,
		"filter_create_time_from":                req.FilterCreateTimeFrom,
		"filter_create_time_to":                  req.FilterCreateTimeTo,
		"filter_create_user_name":                req.FilterCreateUserName,
		"filter_task_execute_start_time_from":    req.FilterTaskExecuteStartTimeFrom,
		"filter_task_execute_start_time_to":      req.FilterTaskExecuteStartTimeTo,
		"filter_status":                          req.FilterStatus,
		"filter_current_step_assignee_user_name": req.FilterCurrentStepAssigneeUserName,
		"filter_task_instance_name":              req.FilterTaskInstanceName,
		"current_user_id":                        user.ID,
		"check_user_can_access":                  user.Name != model.DefaultAdminUser,
		"limit":                                  req.PageSize,
		"offset":                                 offset,
	}

	s := model.GetStorage()
	workflows, count, err := s.GetWorkflowsByReq(data, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowsResV1 := make([]*WorkflowDetailResV1, 0, len(workflows))
	for _, workflow := range workflows {
		workflowRes := &WorkflowDetailResV1{
			ProjectName:             workflow.ProjectName,
			Name:                    workflow.Subject,
			WorkflowId:              workflow.WorkflowId,
			Desc:                    workflow.Desc,
			CreateUser:              utils.AddDelTag(workflow.CreateUserDeletedAt, workflow.CreateUser.String),
			CreateTime:              workflow.CreateTime,
			CurrentStepType:         workflow.CurrentStepType.String,
			CurrentStepAssigneeUser: workflow.CurrentStepAssigneeUser,
			Status:                  workflow.Status,
		}
		workflowsResV1 = append(workflowsResV1, workflowRes)
	}

	return c.JSON(http.StatusOK, GetWorkflowsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      workflowsResV1,
		TotalNums: count,
	})
}

// GetWorkflowsV1
// @Summary 获取工单列表
// @Description get workflow list
// @Tags workflow
// @Id getWorkflowsV1
// @Security ApiKeyAuth
// @Param filter_subject query string false "filter subject"
// @Param filter_workflow_id query string false "filter by workflow_id"
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
// @Success 200 {object} v1.GetWorkflowsResV1
// @router /v1/projects/{project_name}/workflows [get]
func GetWorkflowsV1(c echo.Context) error {
	req := new(GetWorkflowsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectName := c.Param("project_name")

	s := model.GetStorage()

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := CheckIsProjectMember(user.Name, project.Name); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}

	data := map[string]interface{}{
		"filter_workflow_id":                     req.FilterWorkflowID,
		"filter_subject":                         req.FilterSubject,
		"filter_create_time_from":                req.FilterCreateTimeFrom,
		"filter_create_time_to":                  req.FilterCreateTimeTo,
		"filter_create_user_name":                req.FilterCreateUserName,
		"filter_task_execute_start_time_from":    req.FilterTaskExecuteStartTimeFrom,
		"filter_task_execute_start_time_to":      req.FilterTaskExecuteStartTimeTo,
		"filter_status":                          req.FilterStatus,
		"filter_current_step_assignee_user_name": req.FilterCurrentStepAssigneeUserName,
		"filter_task_instance_name":              req.FilterTaskInstanceName,
		"filter_project_name":                    project.Name,
		"current_user_id":                        user.ID,
		"check_user_can_access":                  CheckIsProjectManager(user.Name, project.Name) != nil,
		"limit":                                  req.PageSize,
		"offset":                                 offset,
	}

	workflows, count, err := s.GetWorkflowsByReq(data, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowsResV1 := make([]*WorkflowDetailResV1, 0, len(workflows))
	for _, workflow := range workflows {
		workflowRes := &WorkflowDetailResV1{
			ProjectName:             workflow.ProjectName,
			Name:                    workflow.Subject,
			WorkflowId:              workflow.WorkflowId,
			Desc:                    workflow.Desc,
			CreateUser:              utils.AddDelTag(workflow.CreateUserDeletedAt, workflow.CreateUser.String),
			CreateTime:              workflow.CreateTime,
			CurrentStepType:         workflow.CurrentStepType.String,
			CurrentStepAssigneeUser: workflow.CurrentStepAssigneeUser,
			Status:                  workflow.Status,
		}
		workflowsResV1 = append(workflowsResV1, workflowRes)
	}

	return c.JSON(http.StatusOK, GetWorkflowsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      workflowsResV1,
		TotalNums: count,
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
	FilterSubject                     string `json:"filter_subject" query:"filter_subject"`
	FilterCreateTimeFrom              string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo                string `json:"filter_create_time_to" query:"filter_create_time_to"`
	FilterCreateUserName              string `json:"filter_create_user_name" query:"filter_create_user_name"`
	FilterStatus                      string `json:"filter_status" query:"filter_status" valid:"omitempty,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterCurrentStepAssigneeUserName string `json:"filter_current_step_assignee_user_name" query:"filter_current_step_assignee_user_name"`
	FilterTaskInstanceName            string `json:"filter_task_instance_name" query:"filter_task_instance_name"`
	FilterTaskExecuteStartTimeFrom    string `json:"filter_task_execute_start_time_from" query:"filter_task_execute_start_time_from"`
	FilterTaskExecuteStartTimeTo      string `json:"filter_task_execute_start_time_to" query:"filter_task_execute_start_time_to"`
}

// ExportWorkflowV1
// @Summary 导出工单
// @Description export workflow
// @Id exportWorkflowV1
// @Tags workflow
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
// @Param project_name path string true "project name"
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

	projectName := c.Param("project_name")
	workflowID := c.Param("workflow_id")
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()

	var workflow *model.Workflow
	{
		var exist bool
		workflow, exist, err = s.GetWorkflowDetailByWorkflowID(projectName, workflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
		}
	}

	// check workflow permission
	{
		err = checkBeforeTasksTermination(c, projectName, workflow, user)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	terminatingTaskIDs := getTerminatingTaskIDs(s, workflow, user.ID)

	err = s.UpdateTaskStatusByIDs(terminatingTaskIDs,
		map[string]string{"status": model.TaskStatusTerminating})

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
	projectName := c.Param("project_name")
	workflowID := c.Param("workflow_id")
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()

	var workflow *model.Workflow
	{
		var exist bool
		workflow, exist, err = s.GetWorkflowDetailByWorkflowID(projectName, workflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
		}
	}

	// check workflow permission
	{
		err := checkBeforeTasksTermination(c, projectName, workflow, user)
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
		map[string]string{"status": model.TaskStatusTerminating})

	return c.JSON(http.StatusOK, controller.NewBaseReq(err))
}

func checkBeforeTasksTermination(c echo.Context, projectName string,
	workflow *model.Workflow, user *model.User) error {

	if workflow.Record.Status != model.WorkflowStatusExecuting {
		return errors.NewDataInvalidErr(
			"workflow status is %s, termination can not be performed",
			workflow.Record.Status)
	}

	err := CheckCurrentUserCanOperateWorkflow(c,
		&model.Project{Name: projectName}, workflow, []uint{model.OP_WORKFLOW_EXECUTE})
	if err != nil {
		return err
	}

	return nil
}

func isTaskCanBeTerminate(s *model.Storage, taskID string) (bool, error) {
	task, exist, err := s.GetTaskById(taskID)
	if err != nil {
		return false, fmt.Errorf("get task by id failed. taskID=%v err=%v", taskID, err)
	}
	if !exist {
		return false, fmt.Errorf("task not exist. taskID=%v", taskID)
	}
	if task.Instance == nil {
		return false, fmt.Errorf("task instance is nil. taskID=%v", taskID)
	}

	if task.Status == model.TaskStatusExecuting {
		return true, nil
	}

	return false, nil
}

func getTerminatingTaskIDs(s *model.Storage, workflow *model.Workflow, userID uint) (
	taskIDs []uint) {

	taskIDs = make([]uint, 0)
	for i := range workflow.Record.InstanceRecords {
		instRecord := workflow.Record.InstanceRecords[i]
		if instRecord.Task.Status == model.TaskStatusExecuting {
			taskIDs = append(taskIDs, instRecord.TaskId)
		}
	}
	return taskIDs
}
