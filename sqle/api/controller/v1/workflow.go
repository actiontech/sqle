package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/misc"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

var ErrWorkflowNoAccess = errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist or you can't access it"))
var ErrForbidMyBatisXMLTask = errors.New(errors.DataConflict,
	fmt.Errorf("the task for audit mybatis xml file is not allow to create workflow"))

type GetWorkflowTemplateResV1 struct {
	controller.BaseRes
	Data *WorkflowTemplateDetailResV1 `json:"data"`
}

type WorkflowTemplateDetailResV1 struct {
	Name                          string                       `json:"workflow_template_name"`
	Desc                          string                       `json:"desc,omitempty"`
	AllowSubmitWhenLessAuditLevel *string                      `json:"allow_submit_when_less_audit_level" enums:"normal,notice,warn,error"`
	Steps                         []*WorkFlowStepTemplateResV1 `json:"workflow_step_template_list"`
	Instances                     []string                     `json:"instance_name_list,omitempty"`
}

type WorkFlowStepTemplateResV1 struct {
	Number int      `json:"number"`
	Typ    string   `json:"type"`
	Desc   string   `json:"desc,omitempty"`
	Users  []string `json:"assignee_user_name_list"`
}

// @Summary 获取审批流程模板详情
// @Description get workflow template detail
// @Tags workflow
// @Id getWorkflowTemplateV1
// @Security ApiKeyAuth
// @Param workflow_template_name path string true "workflow template name"
// @Success 200 {object} v1.GetWorkflowTemplateResV1
// @router /v1/workflow_templates/{workflow_template_name}/ [get]
func GetWorkflowTemplate(c echo.Context) error {
	s := model.GetStorage()
	templateName := c.Param("workflow_template_name")
	template, exist, err := s.GetWorkflowTemplateByName(templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("workflow template is not exist")))
	}
	steps, err := s.GetWorkflowStepsDetailByTemplateId(template.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	template.Steps = steps
	res := &WorkflowTemplateDetailResV1{
		Name: template.Name,
		Desc: template.Desc,
	}
	stepsRes := make([]*WorkFlowStepTemplateResV1, 0, len(steps))
	for _, step := range steps {
		stepRes := &WorkFlowStepTemplateResV1{
			Number: int(step.Number),
			Typ:    step.Typ,
			Desc:   step.Desc,
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
		return controller.JSONBaseErrorReq(c, err)
	}
	res.Instances = instanceNames

	return c.JSON(http.StatusOK, &GetWorkflowTemplateResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    res,
	})
}

type CreateWorkflowTemplateReqV1 struct {
	Name                          string                       `json:"workflow_template_name" form:"workflow_template_name" valid:"required,name"`
	Desc                          string                       `json:"desc" form:"desc"`
	AllowSubmitWhenLessAuditLevel *string                      `json:"allow_submit_when_less_audit_level" enums:"normal,notice,warn,error"`
	Steps                         []*WorkFlowStepTemplateReqV1 `json:"workflow_step_template_list" form:"workflow_step_template_list" valid:"required,dive,required"`
	Instances                     []string                     `json:"instance_name_list" form:"instance_name_list"`
}

type WorkFlowStepTemplateReqV1 struct {
	Type  string   `json:"type" form:"type" valid:"oneof=sql_review sql_execute" enums:"sql_review,sql_execute"`
	Desc  string   `json:"desc" form:"desc"`
	Users []string `json:"assignee_user_name_list" form:"assignee_user_name_list" valid:"required"`
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
		if len(step.Users) == 0 {
			return fmt.Errorf("the assignee is empty for step %s", step.Desc)
		}
		if len(step.Users) > 3 {
			return fmt.Errorf("the assignee for step cannot be more than 3")
		}
	}
	return nil
}

// @Summary 创建Sql审批流程模板
// @Description create a workflow template
// @Accept json
// @Produce json
// @Tags workflow
// @Id createWorkflowTemplateV1
// @Security ApiKeyAuth
// @Param instance body v1.CreateWorkflowTemplateReqV1 true "create workflow template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflow_templates [post]
func CreateWorkflowTemplate(c echo.Context) error {
	req := new(CreateWorkflowTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	_, exist, err := s.GetWorkflowTemplateByName(req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("workflow template is exist")))
	}

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

	instances, err := s.GetAndCheckInstanceExist(req.Instances)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowTemplate := &model.WorkflowTemplate{
		Name: req.Name,
		Desc: req.Desc,
	}
	steps := make([]*model.WorkflowStepTemplate, 0, len(req.Steps))
	for i, step := range req.Steps {
		s := &model.WorkflowStepTemplate{
			Number: uint(i + 1),
			Typ:    step.Type,
			Desc:   step.Desc,
		}
		stepUsers := make([]*model.User, 0, len(step.Users))
		for _, userName := range step.Users {
			stepUsers = append(stepUsers, userMap[userName])
		}
		s.Users = stepUsers
		steps = append(steps, s)
	}
	workflowTemplate.Steps = steps

	err = s.SaveWorkflowTemplate(workflowTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateWorkflowTemplateInstances(workflowTemplate, instances...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
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
// @Param workflow_template_name path string true "workflow template name"
// @Param instance body v1.UpdateWorkflowTemplateReqV1 true "create workflow template"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflow_templates/{workflow_template_name}/ [patch]
func UpdateWorkflowTemplate(c echo.Context) error {
	req := new(UpdateWorkflowTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	templateName := c.Param("workflow_template_name")
	workflowTemplate, exist, err := s.GetWorkflowTemplateByName(templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("workflow template is not exist")))
	}
	var instances []*model.Instance
	if req.Instances != nil && len(req.Instances) > 0 {
		instances, err = s.GetAndCheckInstanceExist(req.Instances)
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
				Typ:    step.Type,
				Desc:   step.Desc,
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
		err = s.Save(workflowTemplate)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Instances != nil {
		err = s.UpdateWorkflowTemplateInstances(workflowTemplate, instances...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// @Summary 删除Sql审批流程模板
// @Description update the workflow template
// @Tags workflow
// @Id deleteWorkflowTemplateV1
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param workflow_template_name path string true "workflow template name"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflow_templates/{workflow_template_name}/ [delete]
func DeleteWorkflowTemplate(c echo.Context) error {
	s := model.GetStorage()
	templateName := c.Param("workflow_template_name")
	workflowTemplate, exist, err := s.GetWorkflowTemplateByName(templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("workflow template is not exist")))
	}
	err = s.Delete(workflowTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetWorkflowTemplatesReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetWorkflowTemplatesResV1 struct {
	controller.BaseRes
	Data      []*WorkflowTemplateResV1 `json:"data"`
	TotalNums uint64                   `json:"total_nums"`
}

type WorkflowTemplateResV1 struct {
	Name string `json:"workflow_template_name"`
	Desc string `json:"desc"`
}

// @Summary 获取审批流程模板列表
// @Description get workflow template list
// @Tags workflow
// @Id getWorkflowTemplateListV1
// @Security ApiKeyAuth
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetWorkflowTemplatesResV1
// @router /v1/workflow_templates [get]
func GetWorkflowTemplates(c echo.Context) error {
	req := new(GetWorkflowTemplatesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"limit":  req.PageSize,
		"offset": offset,
	}
	s := model.GetStorage()
	workflowTemplates, count, err := s.GetWorkflowTemplatesByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflowTemplatesReq := make([]*WorkflowTemplateResV1, 0, len(workflowTemplates))
	for _, template := range workflowTemplates {
		workflowTemplateReq := &WorkflowTemplateResV1{
			Name: template.Name,
			Desc: template.Desc,
		}
		workflowTemplatesReq = append(workflowTemplatesReq, workflowTemplateReq)
	}
	return c.JSON(http.StatusOK, &GetWorkflowTemplatesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      workflowTemplatesReq,
		TotalNums: count,
	})
}

type GetWorkflowTemplateTipResV1 struct {
	controller.BaseRes
	Data []*WorkflowTemplateTipResV1 `json:"data"`
}

type WorkflowTemplateTipResV1 struct {
	Name string `json:"workflow_template_name"`
}

// @Summary 获取审批流程模板提示信息
// @Description get workflow template tips
// @Tags workflow
// @Id getWorkflowTemplateTipsV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowTemplateTipResV1
// @router /v1/workflow_template_tips [get]
func GetWorkflowTemplateTips(c echo.Context) error {
	s := model.GetStorage()
	templates, err := s.GetWorkflowTemplateTip()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	templateTipsResV1 := make([]*WorkflowTemplateTipResV1, 0, len(templates))
	for _, template := range templates {
		instanceTipRes := &WorkflowTemplateTipResV1{
			Name: template.Name,
		}
		templateTipsResV1 = append(templateTipsResV1, instanceTipRes)
	}
	return c.JSON(http.StatusOK, &GetWorkflowTemplateTipResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    templateTipsResV1,
	})
}

type CreateWorkflowReqV1 struct {
	Subject string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc    string `json:"desc" form:"desc"`
	TaskId  string `json:"task_id" form:"task_id" valid:"required"`
}

// @Summary 创建工单
// @Description create workflow
// @Accept json
// @Produce json
// @Tags workflow
// @Id createWorkflowV1
// @Security ApiKeyAuth
// @Param instance body v1.CreateWorkflowReqV1 true "create workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows [post]
func CreateWorkflow(c echo.Context) error {
	req := new(CreateWorkflowReqV1)
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

	task, exist, err := s.GetTaskById(req.TaskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}
	err = checkCurrentUserCanAccessTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if task.Instance == nil {
		return controller.JSONBaseErrorReq(c, errInstanceNotExist)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if task.CreateUserId != user.ID {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
			fmt.Errorf("the task is not created by yourself")))
	}

	if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
		return controller.JSONBaseErrorReq(c, ErrForbidMyBatisXMLTask)
	}

	_, exist, err = s.GetWorkflowRecordByTaskId(req.TaskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
			fmt.Errorf("task has been used in other workflow")))
	}

	template, exist, err := s.GetWorkflowTemplateById(task.Instance.WorkflowTemplateId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("the task instance is not bound workflow template")))
	}
	stepTemplates, err := s.GetWorkflowStepsByTemplateId(template.ID)
	if err != nil {
		return err
	}
	err = s.CreateWorkflow(req.Subject, req.Desc, user, task, stepTemplates)
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
	if err := misc.SendEmailIfConfigureSMTP(fmt.Sprintf("%v", workflow.ID)); err != nil {
		log.Logger().Errorf("after create workflow, send email error: %v", err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetWorkflowResV1 struct {
	controller.BaseRes
	Data *WorkflowResV1 `json:"data"`
}

type WorkflowResV1 struct {
	Id            uint                   `json:"workflow_id"`
	Subject       string                 `json:"subject"`
	Desc          string                 `json:"desc,omitempty"`
	CreateUser    string                 `json:"create_user_name"`
	CreateTime    *time.Time             `json:"create_time"`
	Record        *WorkflowRecordResV1   `json:"record"`
	RecordHistory []*WorkflowRecordResV1 `json:"record_history_list,omitempty"`
}

type WorkflowRecordResV1 struct {
	TaskId            uint                 `json:"task_id"`
	CurrentStepNumber uint                 `json:"current_step_number,omitempty"`
	Status            string               `json:"status" enums:"on_process,rejected,canceled,exec_scheduled,executing,exec_failed,finished"`
	ScheduleTime      *time.Time           `json:"schedule_time,omitempty"`
	ScheduleUser      string               `json:"schedule_user,omitempty"`
	Steps             []*WorkflowStepResV1 `json:"workflow_step_list,omitempty"`
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

func checkCurrentUserCanAccessWorkflow(c echo.Context, workflow *model.Workflow) error {
	if controller.GetUserName(c) == model.DefaultAdminUser {
		return nil
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	s := model.GetStorage()
	access, err := s.UserCanAccessWorkflow(user, workflow)
	if err != nil {
		return err
	}
	if !access {
		return ErrWorkflowNoAccess
	}
	return nil
}

func convertWorkflowToRes(workflow *model.Workflow, task *model.Task) *WorkflowResV1 {
	workflowRes := &WorkflowResV1{
		Id:         workflow.ID,
		Subject:    workflow.Subject,
		Desc:       workflow.Desc,
		CreateTime: &workflow.CreatedAt,
	}

	workflowRes.CreateUser = utils.AddDelTag(workflow.CreateUser.DeletedAt, workflow.CreateUserName())

	// convert workflow record
	recordRes := convertWorkflowRecordToRes(workflow, workflow.Record)
	// fill current step number
	for _, step := range recordRes.Steps {
		if step.Id != 0 && step.Id == workflow.Record.CurrentWorkflowStepId {
			recordRes.CurrentStepNumber = step.Number
		}
	}

	// find schedule user name by id in final step(sql execute step),
	// only the person specified in the final step can set the schedule time.
	finalStep := workflow.FinalStep()
	if workflow.Record.ScheduledAt != nil && finalStep.Template.Users != nil {
		for _, user := range finalStep.Template.Users {
			if user.ID == workflow.Record.ScheduleUserId {
				recordRes.ScheduleUser = user.Name
			}
		}
	}
	recordRes.Status = convertWorkflowStatusToRes(workflow.Record.Status, task.Status, workflow.Record.ScheduledAt)
	workflowRes.Record = recordRes

	// convert workflow record history
	recordHistory := make([]*WorkflowRecordResV1, 0, len(workflow.RecordHistory))
	for _, record := range workflow.RecordHistory {
		recordRes := convertWorkflowRecordToRes(workflow, record)
		recordHistory = append(recordHistory, recordRes)
	}
	workflowRes.RecordHistory = recordHistory
	return workflowRes
}

func convertWorkflowRecordToRes(workflow *model.Workflow,
	record *model.WorkflowRecord) *WorkflowRecordResV1 {

	steps := make([]*WorkflowStepResV1, 0, len(record.Steps)+1)
	// It is filled by create user and create time;
	// and tell others that this is a creating or updating operation.
	var stepType string
	if workflow.IsFirstRecord(record) {
		stepType = model.WorkflowStepTypeCreateWorkflow
	} else {
		stepType = model.WorkflowStepTypeUpdateWorkflow
	}

	firstVirtualStep := &WorkflowStepResV1{
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
	for i, step := range steps {
		number := uint(i + 1)
		step.Number = number
	}
	return &WorkflowRecordResV1{
		TaskId:       record.TaskId,
		Status:       record.Status,
		ScheduleTime: record.ScheduledAt,
		Steps:        steps,
	}
}

func convertWorkflowStepToRes(step *model.WorkflowStep) *WorkflowStepResV1 {
	stepRes := &WorkflowStepResV1{
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
	if step.Template.Users != nil {
		for _, user := range step.Template.Users {
			stepRes.Users = append(stepRes.Users, user.Name)
		}
	}
	return stepRes
}

func convertWorkflowStatusToRes(workflowStatus, taskStatus string, scheduleTime *time.Time) string {
	var status = workflowStatus
	switch taskStatus {
	case model.TaskStatusExecuting:
		status = model.WorkflowStatusExecuting
	case model.TaskStatusExecuteSucceeded:
		status = model.WorkflowStatusFinish
	case model.TaskStatusExecuteFailed:
		status = model.WorkflowStatusExecFailed
	}
	if status == model.WorkflowStatusRunning && scheduleTime != nil {
		status = model.WorkflowStatusExecScheduled
	}
	return status
}

// @Summary 获取审批流程详情
// @Description get workflow detail
// @Tags workflow
// @Id getWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path integer true "workflow id"
// @Success 200 {object} v1.GetWorkflowResV1
// @router /v1/workflows/{workflow_id}/ [get]
func GetWorkflow(c echo.Context) error {
	workflowId := c.Param("workflow_id")
	s := model.GetStorage()

	id, err := FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = checkCurrentUserCanAccessWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
	}
	history, err := s.GetWorkflowHistoryById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflow.RecordHistory = history

	task, exist, err := s.GetTaskById(strconv.Itoa(int(workflow.Record.TaskId)))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}

	return c.JSON(http.StatusOK, &GetWorkflowResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowToRes(workflow, task),
	})
}

type GetWorkflowsReqV1 struct {
	FilterSubject                     string `json:"filter_subject" query:"filter_subject"`
	FilterCreateTimeFrom              string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo                string `json:"filter_create_time_to" query:"filter_create_time_to"`
	FilterCreateUserName              string `json:"filter_create_user_name" query:"filter_create_user_name"`
	FilterCurrentStepType             string `json:"filter_current_step_type" query:"filter_current_step_type" valid:"omitempty,oneof=sql_review sql_execute"`
	FilterStatus                      string `json:"filter_status" query:"filter_status" valid:"omitempty,oneof=on_process rejected canceled exec_scheduled executing exec_failed finished"`
	FilterCurrentStepAssigneeUserName string `json:"filter_current_step_assignee_user_name" query:"filter_current_step_assignee_user_name"`
	FilterTaskInstanceName            string `json:"filter_task_instance_name" query:"filter_task_instance_name"`
	PageIndex                         uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize                          uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetWorkflowsResV1 struct {
	controller.BaseRes
	Data      []*WorkflowDetailResV1 `json:"data"`
	TotalNums uint64                 `json:"total_nums"`
}

type WorkflowDetailResV1 struct {
	Id                      uint       `json:"workflow_id"`
	Subject                 string     `json:"subject"`
	Desc                    string     `json:"desc"`
	TaskPassRate            float64    `json:"task_pass_rate"`
	TaskInstance            string     `json:"task_instance_name"`
	TaskInstanceSchema      string     `json:"task_instance_schema"`
	CreateUser              string     `json:"create_user_name"`
	CreateTime              *time.Time `json:"create_time"`
	CurrentStepType         string     `json:"current_step_type,omitempty" enums:"sql_review,sql_execute"`
	CurrentStepAssigneeUser []string   `json:"current_step_assignee_user_name_list,omitempty"`
	Status                  string     `json:"status" enums:"on_process,rejected,canceled,exec_scheduled,executing,exec_failed,finished"`
	ScheduleTime            *time.Time `json:"schedule_time,omitempty"`
}

// @Summary 获取审批流程列表
// @Description get workflow list
// @Tags workflow
// @Id getWorkflowListV1
// @Security ApiKeyAuth
// @Param filter_subject query string false "filter subject"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param filter_create_user_name query string false "filter create user name"
// @Param filter_current_step_type query string false "filter current step type" Enums(sql_review, sql_execute)
// @Param filter_status query string false "filter workflow status" Enums(on_process, rejected, canceled, exec_scheduled, executing, exec_failed, finished)
// @Param filter_current_step_assignee_user_name query string false "filter current step assignee user name"
// @Param filter_task_instance_name query string false "filter instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetWorkflowsResV1
// @router /v1/workflows [get]
func GetWorkflows(c echo.Context) error {
	req := new(GetWorkflowsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	var workflowStatus string
	var taskStatus string
	var isScheduled bool
	var notScheduled bool
	// filter task status
	switch req.FilterStatus {
	case model.WorkflowStatusExecuting:
		taskStatus = model.TaskStatusExecuting
	case model.WorkflowStatusExecFailed:
		taskStatus = model.TaskStatusExecuteFailed
	case model.WorkflowStatusFinish:
		taskStatus = model.TaskStatusExecuteSucceeded
	}
	// filter workflow status
	switch req.FilterStatus {
	case model.WorkflowStatusRunning:
		workflowStatus = model.WorkflowStatusRunning
		notScheduled = true
	case model.WorkflowStatusExecScheduled:
		workflowStatus = model.WorkflowStatusRunning
		isScheduled = true
	case model.WorkflowStatusCancel, model.WorkflowStatusReject:
		workflowStatus = req.FilterStatus
	}

	data := map[string]interface{}{
		"filter_subject":                         req.FilterSubject,
		"filter_create_time_from":                req.FilterCreateTimeFrom,
		"filter_create_time_to":                  req.FilterCreateTimeTo,
		"filter_create_user_name":                req.FilterCreateUserName,
		"filter_status":                          workflowStatus,
		"filter_task_status":                     taskStatus,
		"is_scheduled":                           isScheduled,
		"not_scheduled":                          notScheduled,
		"filter_current_step_type":               req.FilterCurrentStepType,
		"filter_current_step_assignee_user_name": req.FilterCurrentStepAssigneeUserName,
		"filter_task_instance_name":              req.FilterTaskInstanceName,
		"current_user_id":                        user.ID,
		"check_user_can_access":                  user.Name != model.DefaultAdminUser,
		"limit":                                  req.PageSize,
		"offset":                                 offset,
	}
	s := model.GetStorage()
	workflows, count, err := s.GetWorkflowsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowsReq := make([]*WorkflowDetailResV1, 0, len(workflows))
	for _, workflow := range workflows {
		workflowReq := &WorkflowDetailResV1{
			Id:                      workflow.Id,
			Subject:                 workflow.Subject,
			Desc:                    workflow.Desc,
			TaskPassRate:            workflow.TaskPassRate,
			TaskInstance:            utils.AddDelTag(workflow.TaskInstanceDeletedAt, workflow.TaskInstance.String),
			TaskInstanceSchema:      workflow.TaskInstanceSchema,
			CreateUser:              utils.AddDelTag(workflow.CreateUserDeletedAt, workflow.CreateUser.String),
			CreateTime:              workflow.CreateTime,
			CurrentStepType:         workflow.CurrentStepType.String,
			CurrentStepAssigneeUser: workflow.CurrentStepAssigneeUser,
			Status:                  convertWorkflowStatusToRes(workflow.Status, workflow.TaskStatus, workflow.ScheduleTime),
			ScheduleTime:            workflow.ScheduleTime,
		}
		workflowsReq = append(workflowsReq, workflowReq)
	}
	return c.JSON(http.StatusOK, &GetWorkflowsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      workflowsReq,
		TotalNums: count,
	})
}

func checkUserCanOperateStep(user *model.User, workflow *model.Workflow, stepId int) error {
	if workflow.Record.Status != model.WorkflowStatusRunning {
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

// @Summary 审批通过
// @Description approve workflow
// @Tags workflow
// @Id approveWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param workflow_step_id path string true "workflow step id"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/steps/{workflow_step_id}/approve [post]
func ApproveWorkflow(c echo.Context) error {
	workflowId := c.Param("workflow_id")
	id, err := FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = checkCurrentUserCanAccessWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	//TODO: try to using struct tag valid.
	stepIdStr := c.Param("workflow_step_id")
	stepId, err := FormatStringToInt(stepIdStr)
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
		return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
	}

	err = checkUserCanOperateStep(user, workflow, stepId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	currentStep := workflow.CurrentStep()

	if currentStep.Template.Typ == model.WorkflowStepTypeSQLExecute {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow has been approved, you should to execute it")))
	}

	currentStep.State = model.WorkflowStepStateApprove
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.ID
	nextStep := workflow.NextStep()
	workflow.Record.CurrentWorkflowStepId = nextStep.ID

	err = s.UpdateWorkflowStatus(workflow, currentStep)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}

	if err := misc.SendEmailIfConfigureSMTP(workflowId); err != nil {
		log.Logger().Errorf("after approve workflow, send email error: %v", err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type RejectWorkflowReqV1 struct {
	Reason string `json:"reason" form:"reason"`
}

// @Summary 审批驳回
// @Description reject workflow
// @Tags workflow
// @Id rejectWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param workflow_step_id path string true "workflow step id"
// @param workflow_approve body v1.RejectWorkflowReqV1 true "workflow approve request"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/steps/{workflow_step_id}/reject [post]
func RejectWorkflow(c echo.Context) error {
	req := new(RejectWorkflowReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	workflowId := c.Param("workflow_id")
	id, err := FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = checkCurrentUserCanAccessWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	//TODO: try to using struct tag valid.
	stepIdStr := c.Param("workflow_step_id")
	stepId, err := FormatStringToInt(stepIdStr)
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
		return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
	}

	err = checkUserCanOperateStep(user, workflow, stepId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	currentStep := workflow.CurrentStep()
	currentStep.State = model.WorkflowStepStateReject
	currentStep.Reason = req.Reason
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.ID

	workflow.Record.Status = model.WorkflowStatusReject
	workflow.Record.CurrentWorkflowStepId = 0

	err = s.UpdateWorkflowStatus(workflow, currentStep)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// @Summary 审批关闭（中止）
// @Description cancel workflow
// @Tags workflow
// @Id cancelWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/cancel [post]
func CancelWorkflow(c echo.Context) error {
	workflowId := c.Param("workflow_id")
	id, err := FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = checkCurrentUserCanAccessWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, err := checkCancelWorkflow(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !(user.ID == workflow.CreateUserId || user.Name == model.DefaultAdminUser) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	workflow.Record.Status = model.WorkflowStatusCancel
	workflow.Record.CurrentWorkflowStepId = 0

	err = model.GetStorage().UpdateWorkflowStatus(workflow, nil)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type BatchCancelWorkflowsReqV1 struct {
	WorkflowIds []string `json:"workflow_ids" form:"workflow_ids"`
}

// BatchCancelWorkflows batch cancel workflows.
// @Summary 批量取消工单
// @Description batch cancel workflows
// @Tags workflow
// @Id batchCancelWorkflowsV1
// @Security ApiKeyAuth
// @Param BatchCancelWorkflowsReqV1 body v1.BatchCancelWorkflowsReqV1 true "batch cancel workflows request"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/cancel [post]
func BatchCancelWorkflows(c echo.Context) error {
	req := new(BatchCancelWorkflowsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	workflows := make([]*model.Workflow, len(req.WorkflowIds))
	for i, workflowId := range req.WorkflowIds {
		workflow, err := checkCancelWorkflow(workflowId)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		workflows[i] = workflow

		workflow.Record.Status = model.WorkflowStatusCancel
		workflow.Record.CurrentWorkflowStepId = 0
	}

	for _, workflow := range workflows {
		if err := model.GetStorage().UpdateWorkflowStatus(workflow, nil); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	return controller.JSONBaseErrorReq(c, nil)
}

func checkCancelWorkflow(id string) (*model.Workflow, error) {
	workflow, exist, err := model.GetStorage().GetWorkflowDetailById(id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrWorkflowNoAccess
	}
	if !(workflow.Record.Status == model.WorkflowStatusRunning ||
		workflow.Record.Status == model.WorkflowStatusReject) {
		return nil, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status))
	}
	return workflow, nil
}

type UpdateWorkflowReqV1 struct {
	TaskId string `json:"task_id" form:"task_id" valid:"required"`
}

// @Summary 更新审批流程（驳回后才可更新）
// @Description update workflow when it is rejected to creator.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param instance body v1.UpdateWorkflowReqV1 true "update workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/ [patch]
func UpdateWorkflow(c echo.Context) error {
	req := new(UpdateWorkflowReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	workflowId := c.Param("workflow_id")
	id, err := FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = checkCurrentUserCanAccessWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	task, exist, err := s.GetTaskById(req.TaskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrTaskNoAccess)
	}
	err = checkCurrentUserCanAccessTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if task.Instance == nil {
		return controller.JSONBaseErrorReq(c, errInstanceNotExist)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if user.ID != task.CreateUserId {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
			fmt.Errorf("the task is not created by yourself")))
	}

	if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
		return controller.JSONBaseErrorReq(c, ErrForbidMyBatisXMLTask)
	}

	_, exist, err = s.GetWorkflowRecordByTaskId(req.TaskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
			fmt.Errorf("task has been used in other workflow")))
	}

	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
	}

	if workflow.Record.Status != model.WorkflowStatusReject {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)))
	}

	if user.ID != workflow.CreateUserId {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	err = s.UpdateWorkflowRecord(workflow, task)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}

	if err := misc.SendEmailIfConfigureSMTP(workflowId); err != nil {
		log.Logger().Errorf("after update workflow, send email error: %v", err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
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

type UpdateWorkflowScheduleV1 struct {
	ScheduleTime *time.Time `json:"schedule_time"`
}

// @Summary 设置工单定时上线时间（设置为空则代表取消定时时间，需要SQL审核流程都通过后才可以设置）
// @Description update workflow schedule.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowScheduleV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param instance body v1.UpdateWorkflowScheduleV1 true "update workflow schedule request"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/schedule [put]
func UpdateWorkflowSchedule(c echo.Context) error {
	req := new(UpdateWorkflowScheduleV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	workflowId := c.Param("workflow_id")
	id, err := FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = checkCurrentUserCanAccessWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	})
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
		return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
	}
	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return fmt.Errorf("workflow current step not found")
	}

	if currentStep.Template.Typ != model.WorkflowStepTypeSQLExecute {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow need to be approved first")))
	}

	err = checkUserCanOperateStep(user, workflow, int(currentStep.ID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	if req.ScheduleTime != nil && req.ScheduleTime.Before(time.Now()) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf(
			"request schedule time is too early")))
	}

	err = s.UpdateWorkflowSchedule(workflow, user.ID, req.ScheduleTime)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// @Summary 工单提交 SQL 上线
// @Description execute task on workflow
// @Tags workflow
// @Id executeTaskOnWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/task/execute [post]
func ExecuteTaskOnWorkflow(c echo.Context) error {
	workflowId := c.Param("workflow_id")
	id, err := FormatStringToInt(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = checkCurrentUserCanAccessWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	})
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
		return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
	}
	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return fmt.Errorf("workflow current step not found")
	}

	if currentStep.Template.Typ != model.WorkflowStepTypeSQLExecute {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow need to be approved first")))
	}

	err = checkUserCanOperateStep(user, workflow, int(currentStep.ID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	if workflow.Record.ScheduledAt != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow has been set to scheduled execution, not allowed to be executed")))
	}

	err = server.ExecuteWorkflow(workflow, user.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}
