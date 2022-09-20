package v1

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
)

var ErrWorkflowNoAccess = errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist or you can't access it"))

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
	Users                []string `json:"assignee_user_name_list"`
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

type CreateWorkflowTemplateReqV1 struct {
	Name                          string                       `json:"workflow_template_name" form:"workflow_template_name" valid:"required,name"`
	Desc                          string                       `json:"desc" form:"desc"`
	AllowSubmitWhenLessAuditLevel string                       `json:"allow_submit_when_less_audit_level" enums:"normal,notice,warn,error"`
	Steps                         []*WorkFlowStepTemplateReqV1 `json:"workflow_step_template_list" form:"workflow_step_template_list" valid:"required,dive,required"`
	Instances                     []string                     `json:"instance_name_list" form:"instance_name_list"`
}

type WorkFlowStepTemplateReqV1 struct {
	Type                 string   `json:"type" form:"type" valid:"oneof=sql_review sql_execute" enums:"sql_review,sql_execute"`
	Desc                 string   `json:"desc" form:"desc"`
	ApprovedByAuthorized bool     `json:"approved_by_authorized"`
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
		if len(step.Users) == 0 && !step.ApprovedByAuthorized {
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
	allowSubmitWhenLessAuditLevel := string(driver.RuleLevelWarn)
	if req.AllowSubmitWhenLessAuditLevel != "" {
		allowSubmitWhenLessAuditLevel = req.AllowSubmitWhenLessAuditLevel
	}
	workflowTemplate := &model.WorkflowTemplate{
		Name:                          req.Name,
		Desc:                          req.Desc,
		AllowSubmitWhenLessAuditLevel: allowSubmitWhenLessAuditLevel,
	}
	steps := make([]*model.WorkflowStepTemplate, 0, len(req.Steps))
	for i, step := range req.Steps {
		s := &model.WorkflowStepTemplate{
			Number: uint(i + 1),
			ApprovedByAuthorized: sql.NullBool{
				Bool:  step.ApprovedByAuthorized,
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
				ApprovedByAuthorized: sql.NullBool{
					Bool:  step.ApprovedByAuthorized,
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
// @Deprecated
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
	return nil
}

type GetWorkflowResV1 struct {
	controller.BaseRes
	Data *WorkflowResV1 `json:"data"`
}

type WorkflowResV1 struct {
	Id                       uint                    `json:"workflow_id"`
	Subject                  string                  `json:"subject"`
	Desc                     string                  `json:"desc,omitempty"`
	CreateUser               string                  `json:"create_user_name"`
	CreateTime               *time.Time              `json:"create_time"`
	InstanceMaintenanceTimes []*MaintenanceTimeResV1 `json:"instance_maintenance_times"`
	Record                   *WorkflowRecordResV1    `json:"record"`
	RecordHistory            []*WorkflowRecordResV1  `json:"record_history_list,omitempty"`
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

func CheckCurrentUserCanOperateWorkflow(c echo.Context, workflow *model.Workflow, ops []uint) error {
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
	if access {
		return nil
	}
	if len(ops) > 0 {
		instances, err := s.GetInstancesByWorkflowID(workflow.ID)
		if err != nil {
			return err
		}
		ok, err := s.CheckUserHasOpToInstances(user, instances, ops)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return ErrWorkflowNoAccess
}

// @Summary 获取审批流程详情
// @Deprecated
// @Description get workflow detail
// @Tags workflow
// @Id getWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path integer true "workflow id"
// @Success 200 {object} v1.GetWorkflowResV1
// @router /v1/workflows/{workflow_id}/ [get]
func GetWorkflow(c echo.Context) error {
	return nil
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
	Id                      uint       `json:"workflow_id"`
	Subject                 string     `json:"subject"`
	Desc                    string     `json:"desc"`
	TaskPassRate            float64    `json:"task_pass_rate"`
	TaskScore               int32      `json:"task_score"`
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
// @Deprecated
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
// @Param filter_task_execute_start_time_from query string false "filter task execute start time from"
// @Param filter_task_execute_start_time_to query string false "filter task execute start time to"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetWorkflowsResV1
// @router /v1/workflows [get]
func GetWorkflows(c echo.Context) error {
	return nil
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
	err = CheckCurrentUserCanOperateWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	}, []uint{})
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

	err = CheckUserCanOperateStep(user, workflow, stepId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	currentStep := workflow.CurrentStep()

	if workflow.Record.Status == model.WorkflowStatusWaitForExecution {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow has been approved, you should to execute it")))
	}

	currentStep.State = model.WorkflowStepStateApprove
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.ID
	nextStep := workflow.NextStep()
	workflow.Record.CurrentWorkflowStepId = nextStep.ID
	if nextStep.Template.Typ == model.WorkflowStepTypeSQLExecute {
		workflow.Record.Status = model.WorkflowStatusWaitForExecution
	}

	err = s.UpdateWorkflowStatus(workflow, currentStep, nil)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	go notification.NotifyWorkflow(workflowId, notification.WorkflowNotifyTypeApprove)

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

	// RejectWorkflow no need extra operation code for now.
	err = CheckCurrentUserCanOperateWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	}, []uint{})
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

	err = CheckUserCanOperateStep(user, workflow, stepId)
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

	err = s.UpdateWorkflowStatus(workflow, currentStep, nil)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeReject)

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
	err = CheckCurrentUserCanOperateWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(id)},
	}, []uint{})
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

	err = model.GetStorage().UpdateWorkflowStatus(workflow, nil, nil)
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
		if err := model.GetStorage().UpdateWorkflowStatus(workflow, nil, nil); err != nil {
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
// @Deprecated
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

type UpdateWorkflowScheduleV1 struct {
	ScheduleTime *time.Time `json:"schedule_time"`
}

// @Summary 设置工单定时上线时间（设置为空则代表取消定时时间，需要SQL审核流程都通过后才可以设置）
// @Description update workflow schedule.
// @Deprecated
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
	return nil
}

// @Summary 工单提交 SQL 上线
// @Description execute task on workflow
// @Deprecated
// @Tags workflow
// @Id executeTaskOnWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/task/execute [post]
func ExecuteTaskOnWorkflow(c echo.Context) error {
	return nil
}

// ExecuteOneTaskOnWorkflowV1
// @Summary 工单提交单个数据源上线
// @Description execute one task on workflow
// @Tags workflow
// @Id executeOneTaskOnWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param task_id path string true "task id"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/tasks/{task_id}/execute [post]
func ExecuteOneTaskOnWorkflowV1(c echo.Context) error {
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := FormatStringToInt(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	taskIdStr := c.Param("task_id")
	taskId, err := FormatStringToInt(taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(strconv.Itoa(workflowId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = PrepareForWorkflowExecution(c, workflow, user, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	needExecTaskIds, err := GetNeedExecTaskIds(s, workflow)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if _, ok := needExecTaskIds[uint(taskId)]; !ok {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("task has no need to be executed. taskId=%v workflowId=%v", taskId, workflowId))
	}

	err = server.ExecuteWorkflow(workflow, map[uint]struct{}{uint(taskId): {}}, user.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func GetNeedExecTaskIds(s *model.Storage, workflow *model.Workflow) (taskIds map[uint] /*task id*/ struct{}, err error) {
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
	needExecTaskIds := make(map[uint]struct{})
	for _, instRecord := range workflow.Record.InstanceRecords {
		if instRecord.ScheduledAt != nil || instRecord.IsSQLExecuted {
			continue
		}
		needExecTaskIds[instRecord.TaskId] = struct{}{}
	}
	return needExecTaskIds, nil
}

func PrepareForWorkflowExecution(c echo.Context, workflow *model.Workflow, user *model.User, workflowId int) error {
	err := CheckCurrentUserCanOperateWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(workflowId)},
	}, []uint{})
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
	Status                   string                  `json:"status" enums:"wait_for_audit,wait_for_execution,exec_scheduled,exec_failed,exec_succeeded,executing"`
	ExecStartTime            *time.Time              `json:"exec_start_time,omitempty"`
	ExecEndTime              *time.Time              `json:"exec_end_time,omitempty"`
	ScheduleTime             *time.Time              `json:"schedule_time,omitempty"`
	CurrentStepAssigneeUser  []string                `json:"current_step_assignee_user_name_list,omitempty"`
	TaskPassRate             float64                 `json:"task_pass_rate"`
	TaskScore                int32                   `json:"task_score"`
	InstanceMaintenanceTimes []*MaintenanceTimeResV1 `json:"instance_maintenance_times"`
}

// GetSummaryOfWorkflowTasksV1
// @Summary 获取工单数据源任务概览
// @Description get summary of workflow instance tasks
// @Tags workflow
// @Id getSummaryOfInstanceTasksV1
// @Security ApiKeyAuth
// @Param workflow_id path integer true "workflow id"
// @Success 200 {object} v1.GetWorkflowTasksResV1
// @router /v1/workflows/{workflow_id}/tasks [get]
func GetSummaryOfWorkflowTasksV1(c echo.Context) error {
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := FormatStringToInt(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanViewWorkflow(c, &model.Workflow{
		Model: model.Model{ID: uint(workflowId)}})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrWorkflowNoAccess)
	}

	data, err := convertWorkflowToTasksSummaryRes(s, workflow)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetWorkflowTasksResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    data,
	})
}

func convertWorkflowToTasksSummaryRes(s *model.Storage, workflow *model.Workflow) ([]*GetWorkflowTasksItemV1, error) {
	res := make([]*GetWorkflowTasksItemV1, len(workflow.Record.InstanceRecords))
	taskIds := workflow.GetTaskIds()
	taskIds = utils.RemoveDuplicateUint(taskIds)
	tasks, err := s.GetTasksByIds(taskIds)
	if err != nil {
		return nil, err
	}
	if len(tasks) != len(taskIds) {
		return nil, ErrWorkflowNoAccess
	}
	taskIdToTask := make(map[uint]*model.Task, len(tasks))
	for _, task := range tasks {
		taskIdToTask[task.ID] = task
	}

	for i, inst := range workflow.Record.InstanceRecords {
		// convert assignees
		var assignees []string
		// current step is nil if workflow is finished
		if workflow.Record.CurrentStep != nil {
			assignees = make([]string, len(workflow.Record.CurrentStep.Assignees))
			for i, user := range workflow.Record.CurrentStep.Assignees {
				assignees[i] = user.Name
			}
		}

		res[i] = &GetWorkflowTasksItemV1{
			TaskId:                   inst.TaskId,
			InstanceName:             taskIdToTask[inst.TaskId].Instance.Name,
			Status:                   getTaskStatusRes(workflow, taskIdToTask[inst.TaskId], inst.ScheduledAt),
			ExecStartTime:            taskIdToTask[inst.TaskId].ExecStartAt,
			ExecEndTime:              taskIdToTask[inst.TaskId].ExecEndAt,
			ScheduleTime:             inst.ScheduledAt,
			CurrentStepAssigneeUser:  assignees,
			TaskPassRate:             taskIdToTask[inst.TaskId].PassRate,
			TaskScore:                taskIdToTask[inst.TaskId].Score,
			InstanceMaintenanceTimes: convertPeriodToMaintenanceTimeResV1(taskIdToTask[inst.TaskId].Instance.MaintenancePeriod),
		}
	}
	return res, nil
}

const (
	taskDisplayStatusWaitForAudit     = "wait_for_audit"
	taskDisplayStatusWaitForExecution = "wait_for_execution"
	taskDisplayStatusExecFailed       = "exec_failed"
	taskDisplayStatusExecSucceeded    = "exec_succeeded"
	taskDisplayStatusExecuting        = "executing"
	taskDisplayStatusScheduled        = "exec_scheduled"
)

func getTaskStatusRes(workflow *model.Workflow, task *model.Task, scheduleAt *time.Time) (status string) {
	if workflow.Record.Status == model.WorkflowStatusWaitForAudit {
		return taskDisplayStatusWaitForAudit
	}

	if scheduleAt != nil && task.Status == model.TaskStatusAudited {
		return taskDisplayStatusScheduled
	}

	switch task.Status {
	case model.TaskStatusAudited:
		return taskDisplayStatusWaitForExecution
	case model.TaskStatusExecuteSucceeded:
		return taskDisplayStatusExecSucceeded
	case model.TaskStatusExecuteFailed:
		return taskDisplayStatusExecFailed
	case model.TaskStatusExecuting:
		return taskDisplayStatusExecuting
	}
	return ""
}

func CheckCurrentUserCanViewWorkflow(c echo.Context, workflow *model.Workflow) error {
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
	if access {
		return nil
	}
	instances, err := s.GetInstancesByWorkflowID(workflow.ID)
	if err != nil {
		return err
	}
	ok, err := s.CheckUserHasOpToAnyInstance(user, instances, []uint{model.OP_WORKFLOW_VIEW_OTHERS})
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return ErrWorkflowNoAccess
}
