package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	e "errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsV2 "github.com/actiontech/dms/pkg/dms-common/api/dms/v2"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"
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

type CheckedWorkflowInfo struct {
	WorkflowId    string
	User          *model.User
	Tasks         []*model.Task
	InstanceIds   []uint64
	StepTemplates []*model.WorkflowStepTemplate
	ProjectId     model.ProjectUID
	GetOpExecUser func([]*model.Task) (canAuditUsers [][]*model.User, canExecUsers [][]*model.User)
}

func CheckWorkflowCreationPrerequisites(c echo.Context, projectName string, taskIdsToBindWithWorkflow []uint) (*CheckedWorkflowInfo, error) {
	// check project
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), projectName, true)
	if err != nil {
		return nil, err
	}

	s := model.GetStorage()
	// check user
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return nil, err
	}

	// new workflow check id duplicated
	// dms-todo: 与 dms 生成uid保持一致
	workflowId, err := utils.GenUid()
	if err != nil {
		return nil, err
	}

	_, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, workflowId)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New(errors.DataExist, fmt.Errorf("workflow[%v] is exist", workflowId))
	}
	// check task exist
	taskIds := utils.RemoveDuplicateUint(taskIdsToBindWithWorkflow)
	if len(taskIds) > MaximumDataSourceNum {
		return nil, errors.New(errors.DataConflict, fmt.Errorf("the max task count of a workflow is %v", MaximumDataSourceNum))
	}
	tasks, foundAllTasks, err := s.GetTasksByIds(taskIds)
	if err != nil {
		return nil, err
	}
	if !foundAllTasks {
		return nil, errors.NewTaskNoExistOrNoAccessErr()
	}
	// check instances exist
	instanceIdsOfWorkflowTasks := make([]uint64, 0, len(tasks))
	for _, task := range tasks {
		instanceIdsOfWorkflowTasks = append(instanceIdsOfWorkflowTasks, task.InstanceId)
	}

	instancesOfWorkflowInProject, err := dms.GetInstancesInProjectByIds(c.Request().Context(), projectUid, instanceIdsOfWorkflowTasks)
	if err != nil {
		return nil, err
	}

	projectInstanceMap := map[uint64]*model.Instance{}
	for _, instance := range instancesOfWorkflowInProject {
		projectInstanceMap[instance.ID] = instance
	}
	// check template of workflow exist
	workflowTemplate, exist, err := s.GetWorkflowTemplateByProjectId(model.ProjectUID(projectUid))
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(errors.DataNotExist, fmt.Errorf("the task instance is not bound workflow template"))
	}
	// check tasks instance
	for _, task := range tasks {
		if instance, ok := projectInstanceMap[task.InstanceId]; ok {
			task.Instance = instance
		}

		if task.Instance == nil {
			return nil, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist. taskId=%v", task.ID))
		}

		if task.Instance.ProjectId != projectUid {
			return nil, errors.New(errors.DataNotExist, fmt.Errorf("instance is not in project. taskId=%v", task.ID))
		}

		count, err := s.GetTaskSQLCountByTaskID(task.ID)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, errors.New(errors.DataInvalid, fmt.Errorf("workflow's execute sql is null. taskId=%v", task.ID))
		}

		if task.CreateUserId != uint64(user.ID) {
			return nil, errors.New(errors.DataConflict,
				fmt.Errorf("the task is not created by yourself. taskId=%v", task.ID))
		}

		if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
			return nil, ErrForbidMyBatisXMLTask(task.ID)
		}
	}

	// check user role operations
	{

		canOperationInstance, err := CheckCurrentUserCanCreateWorkflow(c.Request().Context(), projectUid, user, tasks)
		if err != nil {
			return nil, err
		}
		if !canOperationInstance {
			return nil, fmt.Errorf("can't operation instance")
		}

	}
	// check if task been used
	count, err := s.GetWorkflowRecordCountByTaskIds(taskIds)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New(errors.DataConflict, fmt.Errorf("task has been used in other workflow"))
	}

	stepTemplates, err := s.GetWorkflowStepsByTemplateId(workflowTemplate.ID)
	if err != nil {
		return nil, err
	}

	memberWithPermissions, _, err := dmsobject.ListMembersInProject(c.Request().Context(), controller.GetDMSServerAddress(), dmsV1.ListMembersForInternalReq{
		ProjectUid: projectUid,
		PageSize:   999,
		PageIndex:  1,
	})
	if err != nil {
		return nil, err
	}
	return &CheckedWorkflowInfo{
		WorkflowId:    workflowId,
		User:          user,
		Tasks:         tasks,
		StepTemplates: stepTemplates,
		InstanceIds:   instanceIdsOfWorkflowTasks,
		ProjectId:     model.ProjectUID(projectUid),
		GetOpExecUser: func(tasks []*model.Task) (auditWorkflowUsers, canExecUser [][]*model.User) {
			auditWorkflowUsers = make([][]*model.User, len(tasks))
			executorWorkflowUsers := make([][]*model.User, len(tasks))
			for i, task := range tasks {
				auditWorkflowUsers[i], err = GetCanOpInstanceUsers(memberWithPermissions, task.Instance, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeAuditWorkflow})
				if err != nil {
					return
				}
				executorWorkflowUsers[i], err = GetCanOpInstanceUsers(memberWithPermissions, task.Instance, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeExecuteWorkflow})
				if err != nil {
					return
				}
			}
			return auditWorkflowUsers, executorWorkflowUsers
		},
	}, nil
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
	Status                  string                `json:"status" enums:"wait_for_audit,wait_for_approve,wait_for_execution,wait_for_export,rejected,canceled,cancel,exec_failed,failed,executing,exporting,finished,finish"`
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
// @Param filter_current_step_assignee_user_id query string false "filter current step assignee user id"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetWorkflowsResV1
// @router /v1/dashboard/workflows [get]
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
			FilterCreateUserId:              req.FilterCreateUserId,
			FilterStatusList:                req.FilterStatusList,
			FilterProjectUid:                req.FilterProjectUid,
			FilterInstanceId:                req.FilterInstanceId,
			FilterProjectPriority:           req.FilterProjectPriority,
			FilterCurrentStepAssigneeUserId: req.FilterCurrentStepAssigneeUserId,
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
	projectMap := make(ProjectMap)
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

type InstanceMap map[string] /* instance id */ *dmsV2.ListDBService

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

		CurrentStepAssigneeUserNames := make([]string, 0)
		for _, currentStepAssigneeUser := range strings.Split(workflow.CurrentStepAssigneeUserIds.String, ",") {
			if currentStepAssigneeUser == "" {
				continue
			}
			CurrentStepAssigneeUserNames = append(CurrentStepAssigneeUserNames, dms.GetUserNameWithDelTag(currentStepAssigneeUser))
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
			CurrentStepAssigneeUser: CurrentStepAssigneeUserNames,
			Status:                  workflow.Status,
		}
		workflowsResV1 = append(workflowsResV1, workflowRes)
	}
	return workflowsResV1
}

type GetGlobalWorkflowStatisticsReqV1 struct {
	FilterCreateUserId              string                `json:"filter_create_user_id" query:"filter_create_user_id"`
	FilterStatusList                []string              `json:"filter_status_list" query:"filter_status_list" validate:"dive,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterProjectUid                string                `json:"filter_project_uid" query:"filter_project_uid"`
	FilterInstanceId                string                `json:"filter_instance_id" query:"filter_instance_id"`
	FilterProjectPriority           dmsV1.ProjectPriority `json:"filter_project_priority" query:"filter_project_priority"  valid:"omitempty,oneof=high medium low"`
	FilterCurrentStepAssigneeUserId string                `json:"filter_current_step_assignee_user_id" query:"filter_current_step_assignee_user_id"`
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
// @Param filter_current_step_assignee_user_id query string false "filter current step assignee user id"
// @Success 200 {object} v1.GlobalWorkflowStatisticsResV1
// @router /v1/dashboard/workflows/statistics [get]
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
			FilterCreateUserId:              req.FilterCreateUserId,
			FilterStatusList:                req.FilterStatusList,
			FilterProjectUid:                req.FilterProjectUid,
			FilterInstanceId:                req.FilterInstanceId,
			FilterProjectPriority:           req.FilterProjectPriority,
			FilterCurrentStepAssigneeUserId: req.FilterCurrentStepAssigneeUserId,
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

// GetGlobalDataExportWorkflowsV1
// @Summary 获取全局导出工单列表
// @Description get global data export workflows list
// @Tags workflow
// @Id getGlobalDataExportWorkflowsV1
// @Security ApiKeyAuth
// @Param filter_create_user_id query string false "filter create user id"
// @Param filter_status_list query []string false "filter by workflow status,support using many status" Enums(wait_for_approve,wait_for_export,exporting,failed,rejected,cancel,finish)
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id in project"
// @Param filter_project_priority query string false "filter by project priority" Enums(high,medium,low)
// @Param filter_current_step_assignee_user_id query string false "filter current step assignee user id"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetWorkflowsResV1
// @router /v1/dashboard/data_export_workflows [get]
func GetGlobalDataExportWorkflowsV1(c echo.Context) error {
	return getGlobalDataExportWorkflowsV1(c)
}

// GetGlobalDataExportWorkflowStatisticsV1
// @Summary 获取全局导出工单统计数据
// @Description get global data export workflows statistics
// @Tags workflow
// @Id getGlobalDataExportWorkflowStatisticsV1
// @Security ApiKeyAuth
// @Param filter_create_user_id query string false "filter create user id"
// @Param filter_status_list query []string false "filter by workflow status,support using many status" Enums(wait_for_approve,wait_for_export,exporting,failed,rejected,cancel,finish)
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id in project"
// @Param filter_project_priority query string false "filter by project priority" Enums(high,medium,low)
// @Param filter_current_step_assignee_user_id query string false "filter current step assignee user id"
// @Success 200 {object} v1.GlobalWorkflowStatisticsResV1
// @router /v1/dashboard/data_export_workflows/statistics [get]
func GetGlobalDataExportWorkflowStatisticsV1(c echo.Context) error {
	return getGlobalDataExportWorkflowStatisticsV1(c)
}

type globalWorkflowBasicFilter struct {
	FilterCreateUserId              string                `json:"filter_create_user_id" query:"filter_create_user_id"`
	FilterStatusList                []string              `json:"filter_status_list" query:"filter_status_list" validate:"dive,oneof=wait_for_audit wait_for_execution rejected canceled executing exec_failed finished"`
	FilterProjectUid                string                `json:"filter_project_uid" query:"filter_project_uid"`
	FilterInstanceId                string                `json:"filter_instance_id" query:"filter_instance_id"`
	FilterProjectPriority           dmsV1.ProjectPriority `json:"filter_project_priority" query:"filter_project_priority"  valid:"omitempty,oneof=high medium low"`
	FilterCurrentStepAssigneeUserId string                `json:"filter_current_step_assignee_user_id" query:"filter_current_step_assignee_user_id"`
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
	// 1.0 如果指定了待操作人筛选，则使用指定的待操作人ID
	if req.FilterCurrentStepAssigneeUserId != "" {
		data["filter_current_step_assignee_user_id"] = req.FilterCurrentStepAssigneeUserId
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
		// 如果请求中已经指定了待操作人筛选，则使用请求中的值，否则使用当前用户ID
		if req.FilterCurrentStepAssigneeUserId == "" {
			data["filter_current_step_assignee_user_id"] = user.GetIDStr()
		}
	}
	return data, nil
}

type VisibilityType string

const (
	GlobalDashBoardVisibilityGlobal   VisibilityType = "global"   // 全局可见
	GlobalDashBoardVisibilityProjects VisibilityType = "projects" // 多项目可见
	GlobalDashBoardVisibilityAssignee VisibilityType = "assignee" // 仅可见授予自己的
)

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
	instances, _, err := dmsobject.ListDbServices(ctx, controller.GetDMSServerAddress(), dmsV2.ListDBServiceReq{
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
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
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

type ReExecuteTaskOnWorkflowReq struct {
	ExecSqlIds []uint `json:"exec_sql_ids" form:"exec_sql_ids" valid:"required"`
}

// ReExecuteTaskOnWorkflowV1
// @Summary 单数据源SQL重新上线
// @Description re-execute task on workflow
// @Tags workflow
// @Id reExecuteTaskOnWorkflowV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param workflow_id path string true "workflow id"
// @Param task_id path string true "task id"
// @Param instance body v1.ReExecuteTaskOnWorkflowReq true "re-execute task on workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/re_execute [post]
func ReExecuteTaskOnWorkflowV1(c echo.Context) error {
	req := new(ReExecuteTaskOnWorkflowReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflowId := c.Param("workflow_id")
	taskId := c.Param("task_id")
	reExecSqlIds := req.ExecSqlIds

	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	task, exist, err := s.GetTaskDetailById(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("task is not exist"))
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := PrepareForTaskReExecution(c, projectUid, workflow, user, task, reExecSqlIds); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = server.ReExecuteTaskSQLs(workflow, task, reExecSqlIds, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func PrepareForTaskReExecution(c echo.Context, projectID string, workflow *model.Workflow, user *model.User, task *model.Task, reExecSqlIds []uint) error {
	// 只有上线失败的工单可以重新上线sql
	if workflow.Record.Status != model.WorkflowStatusExecFailed {
		return errors.New(errors.DataInvalid, e.New("workflow status is not exec failed"))
	}

	if task.Status != model.TaskStatusExecuteFailed {
		return errors.New(errors.DataInvalid, e.New("task status is not execute failed"))
	}

	err := CheckCurrentUserCanOperateTasks(c, projectID, workflow, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeExecuteWorkflow}, []uint{task.ID})
	if err != nil {
		return err
	}

	for _, record := range workflow.Record.InstanceRecords {
		if record.TaskId != task.ID {
			continue
		}

		for _, u := range strings.Split(record.ExecutionAssignees, ",") {
			if u == user.GetIDStr() {
				goto CheckReExecSqlIds
			}
		}
	}

	return e.New("you are not allow to execute the task")

CheckReExecSqlIds:
	// 校验reExecSqlIds对应的SQL状态是否都为SQLExecuteStatusFailed
	if len(reExecSqlIds) == 0 {
		return errors.New(errors.DataInvalid, e.New("re-execute sql ids cannot be empty"))
	}

	// 创建一个map用于快速查找ExecuteSQLs中的SQL
	execSqlMap := make(map[uint]*model.ExecuteSQL)
	for _, execSql := range task.ExecuteSQLs {
		execSqlMap[execSql.ID] = execSql
	}

	// 检查每个reExecSqlId
	for _, sqlId := range reExecSqlIds {
		execSql, exists := execSqlMap[sqlId]
		if !exists {
			return errors.New(errors.DataInvalid, fmt.Errorf("execute sql id %d not found in task", sqlId))
		}

		if execSql.ExecStatus != model.SQLExecuteStatusFailed && execSql.ExecStatus != model.SQLExecuteStatusInitialized {
			return errors.New(errors.DataInvalid, fmt.Errorf("execute sql id %d status is %s, only failed or initialized sql can be re-executed", sqlId, execSql.ExecStatus))
		}
	}

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
	ExportFormat                    string `json:"export_format" query:"export_format" enums:"csv,excel" example:"excel"` // 导出格式：csv 或 excel，默认为 excel
}

// ExportWorkflowV1
// @Summary 导出工单
// @Description export workflow as CSV or Excel
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
// @Param export_format query string false "export format" Enums(csv,excel) "export format: csv or excel, default is excel"
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
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"), true)
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
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"), true)
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

type CreateRollbackWorkflowReq struct {
	Subject        string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc           string `json:"desc" form:"desc"`
	SqlVersionID   *uint  `json:"sql_version_id" form:"sql_version_id"`
	TaskIds        []uint `json:"task_ids" form:"task_ids" valid:"required"`
	RollbackSqlIds []uint `json:"rollback_sql_ids" form:"rollback_sql_ids" valid:"required"`
}

type CreateRollbackWorkflowRes struct {
	controller.BaseRes
	Data *CreateRollbackWorkflowResData `json:"data"`
}

type CreateRollbackWorkflowResData struct {
	WorkflowID string `json:"workflow_id"`
}

// CreateRollbackWorkflow
// @Summary 创建回滚工单
// @Description create rollback workflow
// @Accept json
// @Produce json
// @Tags workflow
// @Id CreateRollbackWorkflow
// @Security ApiKeyAuth
// @Param instance body v1.CreateRollbackWorkflowReq true "create rollback workflow request"
// @Param project_name path string true "project name"
// @Param workflow_id path string true "origin workflow id to rollback"
// @Success 200 {object} CreateRollbackWorkflowRes
// @router /v1/projects/{project_name}/workflows/{workflow_id}/create_rollback_workflow [post]
func CreateRollbackWorkflow(c echo.Context) error {
	return createRollbackWorkflow(c)
}

type AutoCreateAndExecuteWorkflowReqV1 struct {
	// 创建task group的参数
	Instances       []*InstanceForCreatingTask `json:"instances" valid:"dive,required"`
	ExecMode        string                     `json:"exec_mode" enums:"sql_file,sqls"`
	FileOrderMethod string                     `json:"file_order_method"`
	// 审核task group的参数
	Sql string `json:"sql" form:"sql" example:"alter table tb1 drop columns c1"`
	// 创建工单的参数
	Subject string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc    string `json:"desc" form:"desc"`
}

type AutoCreateAndExecuteWorkflowResV1 struct {
	controller.BaseRes
	Data *AutoCreateAndExecuteWorkflowResV1Data `json:"data"`
}

type AutoCreateAndExecuteWorkflowResV1Data struct {
	WorkflowID     string `json:"workflow_id"`
	WorkFlowStatus string `json:"workflow_status"`
}

// AutoCreateAndExecuteWorkflowV1
// @Summary 自动创建工单、审核SQL、审批和上线工单（仅sys用户）
// @Description auto create task group, audit SQL, create workflow, approve and execute workflow (sys user only)
// @Accept mpfd
// @Produce json
// @Tags workflow
// @Id autoCreateAndExecuteWorkflowV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instances formData string true "instances JSON array" example:"[{\"instance_name\":\"inst_1\",\"instance_schema\":\"db1\"}]"
// @Param exec_mode formData string false "exec mode" Enums(sql_file,sqls)
// @Param file_order_method formData string false "file order method"
// @Param sql formData string false "sqls for audit"
// @Param workflow_subject formData string true "workflow subject"
// @Param desc formData string false "workflow description"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Success 200 {object} v1.AutoCreateAndExecuteWorkflowResV1
// @router /v1/projects/{project_name}/workflows/auto_create_and_execute [post]
func AutoCreateAndExecuteWorkflowV1(c echo.Context) error {
	// 1. 权限校验
	user, err := validateSysUserPermission(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectName := c.Param("project_name")
	projectUid, err := dms.GetProjectUIDByName(c.Request().Context(), projectName, true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	// 2. 解析请求参数
	req, err := parseAutoCreateWorkflowRequest(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 3. 准备和验证实例
	nameInstanceMap, err := prepareAndValidateInstances(c.Request().Context(), projectUid, req.Instances)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 4. 创建任务组
	tasks, err := createTaskGroup(s, user, req, nameInstanceMap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 5. 审核任务组
	if err := auditTaskGroup(c, s, tasks, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 6. 验证非 DQL SQL 的备份配置
	if err := validateBackupForNonDQLSQLs(s, tasks); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 7. 检查工单创建前置条件
	taskIds := make([]uint, 0, len(tasks))
	for _, task := range tasks {
		taskIds = append(taskIds, task.ID)
	}

	w, err := CheckWorkflowCreationPrerequisites(c, projectName, taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 8. 创建只有执行步骤的工单模板
	// todo 临时模板创建的工单只能有执行节点
	w.StepTemplates = []*model.WorkflowStepTemplate{
		{
			Number:              1,
			Typ:                 model.WorkflowStepTypeSQLExecute,
			ExecuteByAuthorized: sql.NullBool{Bool: true, Valid: true},
		},
	}

	// 9. 创建工单
	if err := s.CreateWorkflowV2(req.Subject, w.WorkflowId, req.Desc, w.User, w.Tasks, w.StepTemplates, w.ProjectId, nil, nil, nil, w.GetOpExecUser); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 10. 获取创建的工单
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, w.WorkflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 11. 执行工单
	workFlowStatus, err := executeWorkflow(projectUid, workflow, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &AutoCreateAndExecuteWorkflowResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &AutoCreateAndExecuteWorkflowResV1Data{
			WorkflowID:     workflow.WorkflowId,
			WorkFlowStatus: workFlowStatus,
		},
	})
}

// validateSysUserPermission 验证当前用户是否为 sys 用户
func validateSysUserPermission(c echo.Context) (*model.User, error) {
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return nil, err
	}
	if user.Name != model.DefaultSysUser {
		return nil, errors.New(errors.DataInvalid, fmt.Errorf("only sys user can access this API"))
	}
	return user, nil
}

// parseAutoCreateWorkflowRequest 解析自动创建工单的请求参数
func parseAutoCreateWorkflowRequest(c echo.Context) (*AutoCreateAndExecuteWorkflowReqV1, error) {
	req := new(AutoCreateAndExecuteWorkflowReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return nil, err
	}

	// 解析 instances JSON 字符串
	if instancesStr := c.FormValue("instances"); instancesStr != "" {
		if err := json.Unmarshal([]byte(instancesStr), &req.Instances); err != nil {
			return nil, errors.New(errors.DataInvalid, fmt.Errorf("invalid instances JSON: %v", err))
		}
	}
	return req, nil
}

// prepareAndValidateInstances 准备和验证实例，返回实例名称到实例的映射
func prepareAndValidateInstances(ctx context.Context, projectUid string, reqInstances []*InstanceForCreatingTask) (map[string]*model.Instance, error) {
	if len(reqInstances) > MaximumDataSourceNum {
		return nil, ErrTooManyDataSource
	}

	instNames := make([]string, len(reqInstances))
	for i, instance := range reqInstances {
		instNames[i] = instance.InstanceName
	}

	distinctInstNames := utils.RemoveDuplicate(instNames)
	instances, err := dms.GetInstancesInProjectByNames(ctx, projectUid, distinctInstNames)
	if err != nil {
		return nil, err
	}

	nameInstanceMap := make(map[string]*model.Instance, len(reqInstances))
	for _, inst := range instances {
		inst, exist, err := dms.GetInstanceInProjectByName(ctx, projectUid, inst.Name)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, ErrInstanceNoAccess
		}

		if err := common.CheckInstanceIsConnectable(inst); err != nil {
			return nil, err
		}

		nameInstanceMap[inst.Name] = inst
	}

	return nameInstanceMap, nil
}

// createTaskGroup 创建任务组
func createTaskGroup(s *model.Storage, user *model.User, req *AutoCreateAndExecuteWorkflowReqV1, nameInstanceMap map[string]*model.Instance) ([]*model.Task, error) {
	now := time.Now()
	tasks := make([]*model.Task, len(req.Instances))
	for i, reqInstance := range req.Instances {
		instance := nameInstanceMap[reqInstance.InstanceName]
		task := &model.Task{
			Schema:          reqInstance.InstanceSchema,
			InstanceId:      instance.ID,
			CreateUserId:    uint64(user.ID),
			DBType:          instance.DbType,
			ExecMode:        req.ExecMode,
			FileOrderMethod: req.FileOrderMethod,
		}
		task.CreatedAt = now
		tasks[i] = task
	}

	taskGroup := model.TaskGroup{Tasks: tasks}
	if err := s.Save(&taskGroup); err != nil {
		return nil, err
	}

	return tasks, nil
}

// auditTaskGroup 审核任务组
func auditTaskGroup(c echo.Context, s *model.Storage, tasks []*model.Task, req *AutoCreateAndExecuteWorkflowReqV1) error {
	instanceIds := make([]uint64, 0, len(tasks))
	for _, task := range tasks {
		instanceIds = append(instanceIds, task.InstanceId)
	}

	instancesForAudit, err := dms.GetInstancesByIds(c.Request().Context(), instanceIds)
	if err != nil {
		return err
	}

	if len(instancesForAudit) == 0 {
		return errors.New(errors.DataNotExist, fmt.Errorf("no instances found for audit"))
	}

	projectId := instancesForAudit[0].ProjectId
	dbType := instancesForAudit[0].DbType

	// 获取 SQL 内容
	var sqls GetSQLFromFileResp
	var fileRecords []*model.AuditFile
	if req.Sql != "" {
		sqls = GetSQLFromFileResp{
			SourceType:       model.TaskSQLSourceFromFormData,
			SQLsFromFormData: req.Sql,
		}
	} else {
		sqls, err = GetSQLFromFile(c)
		if err != nil {
			return err
		}
		fileRecords, err = saveFileFromContext(c)
		if err != nil {
			return err
		}
	}

	// 处理任务 SQL 和备份配置
	l := log.NewEntry()
	plugin, err := common.NewDriverManagerWithoutCfg(l, dbType)
	if err != nil {
		return err
	}
	defer plugin.Close(context.TODO())

	instanceMap := make(map[uint64]*model.Instance)
	for _, instance := range instancesForAudit {
		instanceMap[instance.ID] = instance
	}

	backupService := server.BackupService{}
	for _, task := range tasks {
		task.SQLSource = sqls.SourceType

		instance, exist := instanceMap[task.InstanceId]
		if !exist {
			return fmt.Errorf("can not find instance in task")
		}

		// 使用数据源上的备份配置
		if instance.EnableBackup {
			if err := backupService.CheckBackupConflictWithExecMode(instance.EnableBackup, task.ExecMode); err != nil {
				return err
			}
			if err := backupService.CheckIsDbTypeSupportEnableBackup(task.DBType); err != nil {
				return err
			}
			task.EnableBackup = instance.EnableBackup
			// 使用数据源上的 BackupMaxRows 配置
			task.BackupMaxRows = backupService.AutoChooseBackupMaxRows(instance.EnableBackup, nil, *instance)
			task.InstanceEnableBackup = instance.EnableBackup
		}

		// 添加 SQL 到任务
		if err := addSQLsFromFileToTasks(sqls, task, plugin); err != nil {
			return errors.New(errors.GenericError, fmt.Errorf("add sqls from file to task failed: %v", err))
		}

		// 处理文件记录
		if len(fileRecords) > 0 {
			fileHeader, _, err := getFileHeaderFromContext(c)
			if err != nil {
				return err
			}
			if strings.HasSuffix(fileHeader.Filename, ZIPFileExtension) && task.FileOrderMethod != "" && task.ExecMode == model.ExecModeSqlFile {
				sortAuditFiles(fileRecords, task.FileOrderMethod)
			}

			if err := batchCreateFileRecords(s, fileRecords, task.ID); err != nil {
				return errors.New(errors.GenericError, fmt.Errorf("save sql file record failed: %v", err))
			}
		}
	}

	// 转换 SQL 编码
	for _, task := range tasks {
		if err := convertSQLSourceEncodingFromTask(task); err != nil {
			return err
		}
	}

	if err := s.Save(tasks); err != nil {
		return err
	}

	// 执行审核
	for i, task := range tasks {
		if task.Status != model.TaskStatusInit {
			continue
		}

		tasks[i], err = server.GetSqled().AddTaskWaitResult(projectId, fmt.Sprintf("%d", task.ID), server.ActionTypeAudit)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateBackupForNonDQLSQLs 验证非 DQL SQL 的备份配置
func validateBackupForNonDQLSQLs(s *model.Storage, tasks []*model.Task) error {
	backupService := server.BackupService{}
	for _, task := range tasks {
		if task.Instance == nil {
			return errors.New(errors.DataNotExist, fmt.Errorf("instance is nil for task %v", task.ID))
		}
		// 检查数据源是否开启备份
		if !task.Instance.EnableBackup {
			continue
		}
		// 检查数据源是否支持备份
		if err := backupService.CheckIsDbTypeSupportEnableBackup(task.DBType); err != nil {
			return errors.New(errors.FeatureNotImplemented, fmt.Errorf("%v instance %v does not support backup: %v", task.DBType, task.Instance.Name, err))
		}

		for _, executeSQL := range task.ExecuteSQLs {
			if executeSQL.SQLType != driverV2.SQLTypeDQL {
				// 检查备份策略
				backupTask, err := s.GetBackupTaskByExecuteSqlId(executeSQL.ID)
				if err != nil {
					return errors.New(errors.GenericError, fmt.Errorf("get backup task err: %v", err))
				}

				if backupTask.BackupStrategy == string(server.BackupStrategyNone) ||
					backupTask.BackupStrategy == string(server.BackupStrategyManually) {
					return errors.New(errors.FeatureNotImplemented, fmt.Errorf("SQL that does not support backup: %q", utils.TruncateStringByRunes(executeSQL.SqlFingerprint, 100)))
				}
			}
		}
	}
	return nil
}

// executeWorkflow 执行工单
func executeWorkflow(projectUid string, workflow *model.Workflow, user *model.User) (string, error) {
	needExecTaskIds, err := GetNeedExecTaskIds(workflow, user)
	if err != nil {
		return "", err
	}

	if len(needExecTaskIds) == 0 {
		return "", nil
	}

	ch, err := server.ExecuteTasksProcess(workflow.WorkflowId, projectUid, user)
	if err != nil {
		return "", err
	}

	return <-ch, nil
}
