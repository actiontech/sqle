package v2

import (
	_err "errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

var ErrForbidMyBatisXMLTask = func(taskId uint) error {
	return errors.New(errors.DataConflict,
		fmt.Errorf("the task for audit mybatis xml file is not allow to create workflow. taskId=%v", taskId))
}
var ErrWorkflowExecuteTimeIncorrect = errors.New(errors.TaskActionInvalid, fmt.Errorf("please go online during instance operation and maintenance time"))

var errTaskHasBeenUsed = errors.New(errors.DataConflict, fmt.Errorf("task has been used in other workflow"))

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
	req := new(CreateWorkflowReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()

	_, exist, err := s.GetWorkflowBySubject(req.Subject)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("workflow is exist")))
	}

	taskIds := utils.RemoveDuplicateUint(req.TaskIds)
	if len(taskIds) > v1.MaximumDataSourceNum {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, fmt.Errorf("the max task count of a workflow is %v", v1.MaximumDataSourceNum)))
	}
	tasks, foundAllTasks, err := s.GetTasksByIds(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !foundAllTasks {
		return controller.JSONBaseErrorReq(c, v1.ErrTaskNoAccess)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowTemplateId := tasks[0].Instance.WorkflowTemplateId
	for _, task := range tasks {
		if task.Instance == nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist. taskId=%v", task.ID)))
		}

		count, err := s.GetTaskSQLCountByTaskID(task.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if count == 0 {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("workflow's execute sql is null. taskId=%v", task.ID)))
		}

		if task.CreateUserId != user.ID {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
				fmt.Errorf("the task is not created by yourself. taskId=%v", task.ID)))
		}

		if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
			return controller.JSONBaseErrorReq(c, ErrForbidMyBatisXMLTask(task.ID))
		}

		// all instances must use the same workflow template
		if task.Instance.WorkflowTemplateId != workflowTemplateId {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
				fmt.Errorf("all instances must use the same workflow template")))
		}
	}

	// check user role operations
	{
		err = checkCurrentUserCanCreateWorkflow(user, tasks)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	count, err := s.GetWorkflowRecordCountByTaskIds(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if count > 0 {
		return controller.JSONBaseErrorReq(c, errTaskHasBeenUsed)
	}

	template, exist, err := s.GetWorkflowTemplateById(workflowTemplateId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("the task instance is not bound workflow template")))
	}

	err = checkWorkflowCanCommit(template, tasks)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	stepTemplates, err := s.GetWorkflowStepsByTemplateId(template.ID)
	if err != nil {
		return err
	}
	err = s.CreateWorkflow(req.Subject, req.Desc, user, tasks, stepTemplates)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, exist, err := s.GetLastWorkflow()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("should exist at least one workflow after create workflow")))
	}
	go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeCreate)

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func checkWorkflowCanCommit(template *model.WorkflowTemplate, tasks []*model.Task) error {
	allowLevel := driver.RuleLevelError
	if template.AllowSubmitWhenLessAuditLevel != "" {
		allowLevel = driver.RuleLevel(template.AllowSubmitWhenLessAuditLevel)
	}
	for _, task := range tasks {
		if driver.RuleLevel(task.AuditLevel).More(allowLevel) {
			return errors.New(errors.DataInvalid,
				fmt.Errorf("there is an audit result with an error level higher than the allowable submission level(%v), please modify it before submitting. taskId=%v", allowLevel, task.ID))
		}
	}
	return nil
}

func checkCurrentUserCanCreateWorkflow(user *model.User, tasks []*model.Task) error {
	if model.IsDefaultAdminUser(user.Name) {
		return nil
	}

	instances := make([]*model.Instance, len(tasks))
	for i, task := range tasks {
		instances[i] = task.Instance
	}

	s := model.GetStorage()
	ok, err := s.CheckUserHasOpToInstances(user, instances, []uint{model.OP_WORKFLOW_SAVE})
	if err != nil {
		return err
	}
	if !ok {
		return errors.NewAccessDeniedErr("user has no access to create workflow for instance")
	}

	return nil
}

type GetWorkflowsReqV2 struct {
	FilterSubject                     string `json:"filter_subject" query:"filter_subject"`
	FilterCreateTimeFrom              string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo                string `json:"filter_create_time_to" query:"filter_create_time_to"`
	FilterCreateUserName              string `json:"filter_create_user_name" query:"filter_create_user_name"`
	FilterStatus                      string `json:"filter_status" query:"filter_status" valid:"omitempty,oneof=wait_for_audit wait_for_execution rejected canceled exec_failed finished"`
	FilterCurrentStepAssigneeUserName string `json:"filter_current_step_assignee_user_name" query:"filter_current_step_assignee_user_name"`
	FilterTaskInstanceName            string `json:"filter_task_instance_name" query:"filter_task_instance_name"`
	PageIndex                         uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize                          uint32 `json:"page_size" query:"page_size" valid:"required"`
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
// @Param filter_status query string false "filter workflow status" Enums(wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,finished)
// @Param filter_current_step_assignee_user_name query string false "filter current step assignee user name"
// @Param filter_task_instance_name query string false "filter instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v2.GetWorkflowsResV2
// @router /v2/workflows [get]
func GetWorkflowsV2(c echo.Context) error {
	req := new(GetWorkflowsReqV2)
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

	var workflowStatus string
	var taskStatus string

	// task status
	switch req.FilterStatus {
	case model.WorkflowStatusExecFailed:
		taskStatus = model.TaskStatusExecuteFailed
	case model.WorkflowStatusFinish:
		taskStatus = model.TaskStatusExecuteSucceeded
	}

	// workflow status
	switch req.FilterStatus {
	case model.WorkflowStatusWaitForAudit, model.WorkflowStatusWaitForExecution, model.WorkflowStatusCancel,
		model.WorkflowStatusReject:

		workflowStatus = req.FilterStatus
	}

	data := map[string]interface{}{
		"filter_subject":                         req.FilterSubject,
		"filter_create_time_from":                req.FilterCreateTimeFrom,
		"filter_create_time_to":                  req.FilterCreateTimeTo,
		"filter_create_user_name":                req.FilterCreateUserName,
		"filter_status":                          workflowStatus,
		"filter_task_status":                     taskStatus,
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

	workflowsReq := make([]*WorkflowDetailResV2, 0, len(workflows))
	for _, workflow := range workflows {
		workflowReq := &WorkflowDetailResV2{
			Id:                      workflow.Id,
			Subject:                 workflow.Subject,
			Desc:                    workflow.Desc,
			CreateUser:              utils.AddDelTag(workflow.CreateUserDeletedAt, workflow.CreateUser.String),
			CreateTime:              workflow.CreateTime,
			CurrentStepType:         workflow.CurrentStepType.String,
			CurrentStepAssigneeUser: workflow.CurrentStepAssigneeUser,
			Status:                  convertWorkflowStatusToRes(workflow.Status, workflow.TaskStatus),
		}
		workflowsReq = append(workflowsReq, workflowReq)
	}

	return c.JSON(http.StatusOK, GetWorkflowsResV2{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      workflowsReq,
		TotalNums: count,
	})
}

func convertWorkflowStatusToRes(workflowStatus string, taskStatus []string) string {
	var status = workflowStatus

	if workflowStatus == model.WorkflowStatusFinish {
		var hasExecuteFailTask bool
		for _, taskStat := range taskStatus {
			if taskStat == model.TaskStatusExecuteFailed {
				hasExecuteFailTask = true
			}
		}

		if hasExecuteFailTask {
			status = model.WorkflowStatusExecFailed
		} else {
			status = model.WorkflowStatusFinish
		}
	}

	return status
}

type GetWorkflowResV2 struct {
	controller.BaseRes
	Data *WorkflowResV2 `json:"data"`
}

type WorkflowTaskItem struct {
	Id uint `json:"task_id"`
}

type WorkflowRecordResV2 struct {
	Tasks             []*WorkflowTaskItem     `json:"tasks"`
	CurrentStepNumber uint                    `json:"current_step_number,omitempty"`
	Status            string                  `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,finished"`
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
	workflowIdStr := c.Param("workflow_id")
	s := model.GetStorage()

	workflowId, err := v1.FormatStringToInt(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = v1.CheckCurrentUserCanViewWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(workflowId)}})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflow, exist, err := s.GetWorkflowDetailById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}
	history, err := s.GetWorkflowHistoryById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflow.RecordHistory = history

	taskIds := workflow.GetTaskIds()
	tasks, _, err := s.GetTasksByIds(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetWorkflowResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowToRes(workflow, tasks),
	})
}

func convertWorkflowToRes(workflow *model.Workflow, tasks []*model.Task) *WorkflowResV2 {
	workflowRes := &WorkflowResV2{
		Id:         workflow.ID,
		Subject:    workflow.Subject,
		Desc:       workflow.Desc,
		Mode:       workflow.Mode,
		CreateUser: workflow.CreateUser.Name,
		CreateTime: &workflow.CreatedAt,
	}

	// convert workflow record
	workflowRecordRes := convertWorkflowRecordToRes(workflow, workflow.Record, tasks)

	// convert workflow record history
	recordHistory := make([]*WorkflowRecordResV2, 0, len(workflow.RecordHistory))
	for _, record := range workflow.RecordHistory {
		recordRes := convertWorkflowRecordToRes(workflow, record, tasks)
		recordHistory = append(recordHistory, recordRes)
	}
	workflowRes.RecordHistory = recordHistory
	workflowRes.Record = workflowRecordRes
	return workflowRes
}

func convertWorkflowRecordToRes(workflow *model.Workflow,
	record *model.WorkflowRecord, tasks []*model.Task) *WorkflowRecordResV2 {

	steps := make([]*v1.WorkflowStepResV1, 0, len(record.Steps)+1)
	// It is filled by create user and create time;
	// and tell others that this is a creating or updating operation.
	var stepType string
	if workflow.IsFirstRecord(record) {
		stepType = model.WorkflowStepTypeCreateWorkflow
	} else {
		stepType = model.WorkflowStepTypeUpdateWorkflow
	}

	firstVirtualStep := &v1.WorkflowStepResV1{
		Type:          stepType,
		OperationTime: &record.CreatedAt,
		OperationUser: workflow.CreateUserName(),
	}
	steps = append(steps, firstVirtualStep)

	// convert workflow actual step
	for _, step := range record.Steps {
		stepRes := convertWorkflowStepToRes(step)
		steps = append(steps, stepRes)
	}
	// fill step number
	var currentStepNum uint
	for i, step := range steps {
		number := uint(i + 1)
		step.Number = number
		if step.Id != 0 && step.Id == record.CurrentWorkflowStepId {
			currentStepNum = number
		}
	}

	tasksRes := make([]*WorkflowTaskItem, len(record.InstanceRecords))
	for i, inst := range record.InstanceRecords {
		tasksRes[i] = &WorkflowTaskItem{Id: inst.TaskId}
	}

	workflowStatus := record.Status
	if record.Status == model.WorkflowStatusFinish {
		for _, task := range tasks {
			if task.Status == model.TaskStatusExecuteFailed {
				workflowStatus = model.WorkflowStatusExecFailed
			}
		}
	}

	return &WorkflowRecordResV2{
		Tasks:             tasksRes,
		CurrentStepNumber: currentStepNum,
		Status:            workflowStatus,
		Steps:             steps,
	}
}

func convertWorkflowStepToRes(step *model.WorkflowStep) *v1.WorkflowStepResV1 {
	stepRes := &v1.WorkflowStepResV1{
		Id:            step.ID,
		Type:          step.Template.Typ,
		Desc:          step.Template.Desc,
		OperationTime: step.OperateAt,
		State:         step.State,
		Reason:        step.Reason,
		Users:         []string{},
	}
	if step.OperationUser != nil {
		stepRes.OperationUser = step.OperationUser.Name
	}
	if step.Assignees != nil {
		for _, user := range step.Assignees {
			stepRes.Users = append(stepRes.Users, user.Name)
		}
	}
	return stepRes
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
	req := new(UpdateWorkflowReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := v1.FormatStringToInt(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = v1.CheckCurrentUserCanOperateWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(workflowId)},
	}, []uint{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	tasks, _, err := s.GetTasksByIds(req.TaskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if len(tasks) <= 0 {
		return controller.JSONBaseErrorReq(c, v1.ErrTaskNoAccess)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskIds := make([]uint, len(tasks))
	for i, task := range tasks {
		taskIds[i] = task.ID

		count, err := s.GetTaskSQLCountByTaskID(task.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if count == 0 {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("task's execute sql is null. taskId=%v", task.ID)))
		}

		err = v1.CheckCurrentUserCanViewTask(c, task)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		if task.Instance == nil {
			return controller.JSONBaseErrorReq(c, v1.ErrInstanceNotExist)
		}

		if user.ID != task.CreateUserId {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
				fmt.Errorf("the task is not created by yourself. taskId=%v", task.ID)))
		}

		if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
			return controller.JSONBaseErrorReq(c, ErrForbidMyBatisXMLTask(task.ID))
		}
	}

	count, err := s.GetWorkflowRecordCountByTaskIds(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if count > 0 {
		return controller.JSONBaseErrorReq(c, errTaskHasBeenUsed)
	}

	workflow, exist, err := s.GetWorkflowDetailById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	if workflow.Record.Status != model.WorkflowStatusReject {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)))
	}

	if user.ID != workflow.CreateUserId {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	template, exist, err := s.GetWorkflowTemplateById(tasks[0].Instance.WorkflowTemplateId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
			fmt.Errorf("failed to find the corresponding workflow template based on the task id")))
	}

	err = checkWorkflowCanCommit(template, tasks)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateWorkflowRecord(workflow, tasks)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	go notification.NotifyWorkflow(workflowIdStr, notification.WorkflowNotifyTypeCreate)

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
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
// @Param instance body v1.UpdateWorkflowScheduleV1 true "update workflow schedule request"
// @Success 200 {object} controller.BaseRes
// @router /v2/workflows/{workflow_id}/tasks/{task_id}/schedule [put]
func UpdateWorkflowScheduleV2(c echo.Context) error {
	workflowId := c.Param("workflow_id")
	workflowIdInt, err := v1.FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	taskId := c.Param("task_id")
	taskIdUint, err := v1.FormatStringToUint64(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	req := new(v1.UpdateWorkflowScheduleV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(workflowIdInt)},
	}, []uint{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}
	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, _err.New("workflow current step not found")))
	}

	if workflow.Record.Status != model.WorkflowStatusWaitForExecution {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow need to be approved first")))
	}

	err = v1.CheckUserCanOperateStep(user, workflow, int(currentStep.ID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}
	var curTaskRecord *model.WorkflowInstanceRecord
	for _, ir := range workflow.Record.InstanceRecords {
		if uint64(ir.TaskId) == taskIdUint {
			curTaskRecord = ir
		}
	}
	if curTaskRecord == nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, _err.New("task is not found in workflow")))
	}

	if req.ScheduleTime != nil && req.ScheduleTime.Before(time.Now()) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf(
			"request schedule time is too early")))
	}

	if curTaskRecord.IsSQLExecuted {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf(
			"task has been executed")))
	}

	instance, exist, err := s.GetInstanceById(fmt.Sprintf("%v", curTaskRecord.InstanceId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrInstanceNotExist)
	}

	if req.ScheduleTime != nil && len(instance.MaintenancePeriod) != 0 && !instance.MaintenancePeriod.IsWithinScope(*req.ScheduleTime) {
		return controller.JSONBaseErrorReq(c, ErrWorkflowExecuteTimeIncorrect)
	}

	err = s.UpdateInstanceRecordSchedule(curTaskRecord, user.ID, req.ScheduleTime)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// ExecuteTasksOnWorkflow
// @Summary 多数据源批量上线
// @Description execute tasks on workflow
// @Tags workflow
// @Id executeTasksOnWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Success 200 {object} controller.BaseRes
// @router /v2/workflows/{workflow_id}/tasks/execute [post]
func ExecuteTasksOnWorkflow(c echo.Context) error {
	workflowId := c.Param("workflow_id")
	id, err := v1.FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(strconv.Itoa(id))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err := v1.PrepareForWorkflowExecution(c, workflow, user, id); err != nil {
		return err
	}

	needExecTaskIds, err := v1.GetNeedExecTaskIds(s, workflow)
	if err != nil {
		return err
	}

	err = server.ExecuteWorkflow(workflow, needExecTaskIds, user.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}
