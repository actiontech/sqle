package v1

import (
	"actiontech.cloud/universe/sqle/v4/sqle/api/controller"
	"actiontech.cloud/universe/sqle/v4/sqle/api/server"
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"actiontech.cloud/universe/sqle/v4/sqle/executor"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type GetWorkflowTemplateResV1 struct {
	controller.BaseRes
	Data *WorkflowTemplateDetailResV1 `json:"data"`
}

type WorkflowTemplateDetailResV1 struct {
	Name      string                       `json:"workflow_template_name"`
	Desc      string                       `json:"desc"`
	Steps     []*WorkFlowStepTemplateResV1 `json:"workflow_step_template_list"`
	Instances []string                     `json:"instance_name_list"`
}

type WorkFlowStepTemplateResV1 struct {
	Number int      `json:"number"`
	Typ    string   `json:"type"`
	Desc   string   `json:"desc"`
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
	Name      string                       `json:"workflow_template_name" form:"workflow_template_name" valid:"required"`
	Desc      string                       `json:"desc" form:"desc"`
	Steps     []*WorkFlowStepTemplateReqV1 `json:"workflow_step_template_list" form:"workflow_step_template_list" valid:"required"`
	Instances []string                     `json:"instance_name_list" form:"instance_name_list"`
}

type WorkFlowStepTemplateReqV1 struct {
	Type  string   `json:"type" form:"type" valid:"in(sql_review|sql_execute)" enums:"sql_review, sql_execute"`
	Desc  string   `json:"desc" form:"desc"`
	Users []string `json:"assignee_user_name_list" form:"assignee_user_name_list" valid:"required"`
}

func validWorkflowTemplateReq(steps []*WorkFlowStepTemplateReqV1) error {
	if steps == nil || len(steps) == 0 {
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
	Desc      *string                      `json:"desc" form:"desc"`
	Steps     []*WorkFlowStepTemplateReqV1 `json:"workflow_step_template_list" form:"workflow_step_template_list"`
	Instances []string                     `json:"instance_name_list" form:"instance_name_list"`
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
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required,int"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required,int"`
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
	Subject string `json:"workflow_subject" form:"workflow_subject"`
	Desc    string `json:"desc" form:"desc"`
	TaskId  string `json:"task_id" form:"task_id"`
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

	task, exist, err := s.GetTaskById(req.TaskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("task is not exist")))
	}

	_, exist, err = s.GetWorkflowByTaskId(req.TaskId)
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
	user, exist, err := s.GetUserByName(controller.GetUserName(c))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("current user is not exist")))
	}
	workflow := &model.Workflow{
		Subject:      req.Subject,
		Desc:         req.Desc,
		CreateUserId: user.ID,
		TaskId:       task.ID,
	}
	workflow.InitWorkflowStepByTemplate(stepTemplates)
	err = s.SaveWorkflow(workflow)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetWorkflowResV1 struct {
	controller.BaseRes
	Data *WorkflowResV1 `json:"data"`
}

type WorkflowResV1 struct {
	Id                uint                 `json:"workflow_id"`
	Subject           string               `json:"subject"`
	Desc              string               `json:"desc"`
	TaskId            uint                 `json:"task_id"`
	CreateUser        string               `json:"create_user_name"`
	CreateTime        *time.Time           `json:"create_time"`
	CurrentStepNumber uint                 `json:"current_step_number,omitempty"`
	Status            string               `json:"status" enums:"on_process, finished, rejected, canceled"`
	Steps             []*WorkflowStepResV1 `json:"workflow_step_list,omitempty"`
}

type WorkflowStepResV1 struct {
	Number        uint       `json:"number"`
	Type          string     `json:"type"`
	Desc          string     `json:"desc"`
	Users         []string   `json:"assignee_user_name_list,omitempty"`
	OperationUser string     `json:"operation_user_name,omitempty"`
	OperationTime *time.Time `json:"operation_time,omitempty"`
	State         string     `json:"state" enums:"initialized, approved, rejected"`
	Reason        string     `json:"reason"`
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
	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist")))
	}
	workflowRes := &WorkflowResV1{
		Id:         workflow.ID,
		Subject:    workflow.Subject,
		Desc:       workflow.Desc,
		TaskId:     workflow.TaskId,
		CreateTime: &workflow.CreatedAt,
		Status:     workflow.Record.Status,
	}
	if workflow.CreateUser != nil {
		workflowRes.CreateUser = workflow.CreateUser.Name
	}

	stepsRes := make([]*WorkflowStepResV1, 0, len(workflow.Record.Steps))
	for _, step := range workflow.Record.Steps {
		stepRes := &WorkflowStepResV1{
			Number:        step.Template.Number,
			Type:          step.Template.Typ,
			Desc:          step.Template.Desc,
			OperationTime: step.OperateAt,
			State:         step.State,
			Reason:        step.Reason,
		}
		if step.OperationUser != nil {
			stepRes.OperationUser = step.OperationUser.Name
		}
		userNames := []string{}
		if step.Template.Users != nil {
			for _, user := range step.Template.Users {
				userNames = append(userNames, user.Name)
			}
		}
		stepRes.Users = userNames
		stepsRes = append(stepsRes, stepRes)

		if workflow.CurrentStep() != nil {
			workflowRes.CurrentStepNumber = workflow.CurrentStep().Template.Number
		}
	}
	workflowRes.Steps = stepsRes
	return c.JSON(http.StatusOK, &GetWorkflowResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    workflowRes,
	})
}

type GetWorkflowsReqV1 struct {
	FilterCreateUserName              string `json:"filter_create_user_name" query:"filter_create_user_name"`
	FilterCurrentStepType             string `json:"filter_current_step_type" query:"filter_current_step_type" enums:"sql_review, sql_execute"`
	FilterStatus                      string `json:"filter_status" query:"filter_status" enums:"on_process, finished, rejected, canceled"`
	FilterCurrentStepAssigneeUserName string `json:"filter_current_step_assignee_user_name" query:"filter_current_step_assignee_user_name"`
	FilterTaskStatus                  string `json:"filter_task_status" query:"filter_task_status" enums:"initialized, audited, executing, exec_success, exec_failed"`
	FilterTaskInstanceName            string `json:"filter_task_instance_name" query:"filter_task_instance_name"`
	PageIndex                         uint32 `json:"page_index" query:"page_index" valid:"required,int"`
	PageSize                          uint32 `json:"page_size" query:"page_size" valid:"required,int"`
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
	TaskStatus              string     `json:"task_status" enums:"initialized, audited, executing, exec_success, exec_failed"`
	TaskPassRate            float64    `json:"task_pass_rate"`
	TaskInstance            string     `json:"task_instance_name"`
	TaskInstanceSchema      string     `json:"task_instance_schema"`
	CreateUser              string     `json:"create_user_name"`
	CreateTime              *time.Time `json:"create_time"`
	CurrentStepType         string     `json:"current_step_type,omitempty" enums:"sql_review, sql_execute"`
	CurrentStepAssigneeUser []string   `json:"current_step_assignee_user_name_list,omitempty"`
	Status                  string     `json:"status" enums:"on_process, finished, rejected, canceled"`
}

// @Summary 获取审批流程列表
// @Description get workflow list
// @Tags workflow
// @Id getWorkflowListV1
// @Security ApiKeyAuth
// @Param filter_create_user_name query string false "filter create user name"
// @Param filter_current_step_type query string false "filter current step type" Enums(sql_review, sql_execute)
// @Param filter_status query string false "filter workflow status" Enums(on_process, finished, rejected, canceled)
// @Param filter_current_step_assignee_user_name query string false "filter current step assignee user name"
// @Param filter_task_status query string false "filter task status" Enums(initialized, audited, executing, exec_success, exec_failed)
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

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"filter_create_user_name":                req.FilterCreateUserName,
		"filter_status":                          req.FilterStatus,
		"filter_current_step_type":               req.FilterCurrentStepType,
		"filter_current_step_assignee_user_name": req.FilterCurrentStepAssigneeUserName,
		"filter_task_status":                     req.FilterTaskStatus,
		"filter_task_instance_name":              req.FilterTaskInstanceName,
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
			TaskStatus:              workflow.TaskStatus,
			TaskPassRate:            workflow.TaskPassRate,
			TaskInstance:            workflow.TaskInstance,
			TaskInstanceSchema:      workflow.TaskInstanceSchema,
			CreateUser:              workflow.CreateUser,
			CreateTime:              workflow.CreateTime,
			CurrentStepType:         workflow.CurrentStepType.String,
			CurrentStepAssigneeUser: workflow.CurrentStepAssigneeUser,
			Status:                  workflow.Status,
		}
		workflowsReq = append(workflowsReq, workflowReq)
	}
	return c.JSON(http.StatusOK, &GetWorkflowsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      workflowsReq,
		TotalNums: count,
	})
}

// @Summary 审批通过
// @Description approve workflow
// @Tags workflow
// @Id approveWorkflowV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param workflow_step_number path string true "workflow step number"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/steps/{workflow_step_number}/approve [post]
func ApproveWorkflow(c echo.Context) error {
	workflowId := c.Param("workflow_id")
	//TODO: try to using struct tag valid.
	stepNumberStr := c.Param("workflow_step_number")
	stepNumber, err := FormatStringToInt(stepNumberStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist")))
	}

	if workflow.Record.Status != model.WorkflowStatusRunning {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)))
	}

	if workflow.CurrentStep() == nil {
		return controller.JSONBaseErrorReq(c, errors.New(
			errors.DataInvalid, fmt.Errorf("workflow current step not found")))
	}

	if uint(stepNumber) != workflow.CurrentStep().Template.Number {
		return controller.JSONBaseErrorReq(c, errors.New(
			errors.DataInvalid, fmt.Errorf("workflow current step is not %d", stepNumber)))
	}

	user, exist, err := s.GetUserByName(controller.GetUserName(c))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("current user is not exist")))
	}

	if !workflow.IsOperationUser(user) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	currentStep := workflow.CurrentStep()
	currentStep.State = model.WorkflowStepStateApprove
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.ID

	if currentStep == workflow.FinalStep() {
		workflow.Record.Status = model.WorkflowStatusFinish
		workflow.Record.CurrentWorkflowStepId = 0
	} else {
		nextStep := workflow.NextStep()
		workflow.Record.CurrentWorkflowStepId = nextStep.ID
	}

	err = s.UpdateWorkflowStatus(workflow, currentStep)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	if currentStep.Template.Typ == model.WorkflowStepTypeSQLExecute {
		taskId := fmt.Sprintf("%d", workflow.TaskId)
		task, exist, err := s.GetTaskDetailById(taskId)
		if err != nil {
			return c.JSON(http.StatusOK, controller.NewBaseReq(err))
		}
		if !exist {
			return c.JSON(http.StatusOK, controller.NewBaseReq(
				errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))))
		}
		if task.Instance == nil {
			return c.JSON(http.StatusOK, controller.NewBaseReq(
				errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))))
		}

		// if instance is not connectable, exec sql must be failed;
		// commit action unable to retry, so don't to exec it.
		if err := executor.Ping(log.NewEntry(), task.Instance); err != nil {
			return c.JSON(http.StatusOK, controller.NewBaseReq(err))
		}

		sqledServer := server.GetSqled()
		err = sqledServer.AddTask(taskId, model.TASK_ACTION_EXECUTE)
		if err != nil {
			return c.JSON(http.StatusOK, controller.NewBaseReq(err))
		}
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
// @Param workflow_step_number path string true "workflow step number"
// @param workflow_approve body v1.RejectWorkflowReqV1 true "workflow approve request"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/steps/{workflow_step_number}/reject [post]
func RejectWorkflow(c echo.Context) error {
	req := new(RejectWorkflowReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	workflowId := c.Param("workflow_id")
	//TODO: try to using struct tag valid.
	stepNumberStr := c.Param("workflow_step_number")
	stepNumber, err := FormatStringToInt(stepNumberStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist")))
	}

	if workflow.Record.Status != model.WorkflowStatusRunning {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)))
	}

	if workflow.CurrentStep() == nil {
		return controller.JSONBaseErrorReq(c, errors.New(
			errors.DataInvalid, fmt.Errorf("workflow current step not found")))
	}

	if uint(stepNumber) != workflow.CurrentStep().Template.Number {
		return controller.JSONBaseErrorReq(c, errors.New(
			errors.DataInvalid, fmt.Errorf("workflow current step is not %d", stepNumber)))
	}

	user, exist, err := s.GetUserByName(controller.GetUserName(c))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("current user is not exist")))
	}

	if !workflow.IsOperationUser(user) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
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

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist")))
	}

	if !(workflow.Record.Status == model.WorkflowStatusRunning ||
		workflow.Record.Status == model.WorkflowStatusReject) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)))
	}

	user, exist, err := s.GetUserByName(controller.GetUserName(c))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("current user is not exist")))
	}

	if !(user.ID == workflow.CreateUserId || user.Name == defaultAdminUser) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	workflow.Record.Status = model.WorkflowStatusCancel
	workflow.Record.CurrentWorkflowStepId = 0

	err = s.UpdateWorkflowStatus(workflow, nil)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
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
