package controller

import (
	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"github.com/labstack/echo/v4"
)

type UserRes struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type GetWorkflowTemplateRes struct {
	BaseRes
	Data WorkflowTemplateRes `json:"data"`
}

type WorkflowTemplateRes struct {
	TemplateName string          `json:"name"`
	TemplateDesc string          `json:"desc"`
	Steps        WorkFlowStepRes `json:"step_list"`
}

type WorkFlowStepRes struct {
	StepNumber string    `json:"number"`
	Typ        string    `json:"type"`
	Desc       string    `json:"desc"`
	Users      []UserRes `json:"assignee_user_list"`
}

// @Summary 获取审批流程模板详情
// @Description get workflow template detail
// @Param workflow_template_id path string true "workflow template ID"
// @Success 200 {object} GetWorkflowTemplateRes
// @router /v1/workflow_templates/{workflow_template_id}/ [get]
func GetWorkflowTemplate(c echo.Context) error {
	return nil
}

type UserReq struct {
	UserName string `json:"user_name" form:"user_name"`
}

type CreateWorkflowTemplateReq struct {
	TemplateName string          `json:"name" form:"name"`
	TemplateDesc string          `json:"desc" form:"desc"`
	Steps        WorkFlowStepReq `json:"step_list" form:"step_list"`
}

type WorkFlowStepReq struct {
	StepNumber string    `json:"number" form:"number"`
	Typ        string    `json:"type" form:"type"`
	Desc       string    `json:"desc" form:"desc"`
	Users      []UserReq `json:"assignee_user_list" form:"assignee_user_list"`
}

// @Summary 创建Sql审批流程模板
// @Description create a workflow template
// @Accept json
// @Produce json
// @Param instance body controller.CreateWorkflowTemplateReq true "create workflow template"
// @Success 200 {object} controller.GetWorkflowTemplateRes
// @router /v1/workflow_templates [post]
func CreateWorkflowTemplate(c echo.Context) error {
	return nil
}

// @Summary 更新Sql审批流程模板
// @Description update the workflow template
// @Accept json
// @Produce json
// @Param workflow_template_id path string true "workflow template ID"
// @Param instance body controller.CreateWorkflowTemplateReq true "create workflow template"
// @Success 200 {object} controller.GetWorkflowTemplateRes
// @router /v1/workflow_templates/{workflow_template_id}/ [put]
func UpdateWorkflowTemplate(c echo.Context) error {
	return nil
}

// @Summary 删除Sql审批流程模板
// @Description update the workflow template
// @Accept json
// @Produce json
// @Param workflow_template_id path string true "workflow template ID"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflow_templates/{workflow_template_id}/ [delete]
func DeleteWorkflowTemplate(c echo.Context) error {
	return nil
}


type GetAllWorkflowTemplatesRes struct {
	BaseRes
	Data []model.WorkflowTemplate `json:"data"`
}

// @Summary 获取审批流程模板列表
// @Description get workflow template list
// @Success 200 {object} controller.GetAllWorkflowTemplatesRes
// @router /v1/workflow_templates [get]
func GetWorkflowTemplates(c echo.Context) error {
	return nil
}

// @Summary 创建工单
// @Description create workflow
// @Param task_id path string true "Task ID"
// @Success 200 {object} controller.GetTaskRes
// @router /v1/workflows [post]
func CreateWorkflow(c echo.Context) error {
	return nil
}

type GetWorkflowRes struct {
	BaseRes
	Data WorkflowRes `json:"data"`
}

type WorkflowRes struct {
	Id                int              `json:"id"`
	TaskId            int              `json:"task_id"`
	CreateUser        string           `json:"create_user_name"`
	CreateTime        string           `json:"create_time"`
	CurrentStepNumber int              `json:"current_step_number"`
	State             int              `json:"state"`
	Steps             []WorkflowRecord `json:"step_record_list"`
}

type WorkflowRecord struct {
	Id            int       `json:"id"`
	StepNumber    int       `json:"step_number"`
	Typ           string    `json:"type"`
	Desc          string    `json:"desc"`
	OperationUser string    `json:"operation_user_name"`
	OperationTime string    `json:"operation_time"`
	State         int       `json:"state"`
	Reason        string    `json:"reason"`
	Users         []UserReq `json:"assignee_user_list"`
}

// @Summary 获取审批流程详情
// @Description get workflow detail
// @Param workflow_id path string true "workflow ID"
// @Success 200 {object} controller.GetWorkflowRes
// @router /v1/workflows/{workflow_id}/ [get]
func GetWorkflow(c echo.Context) error {
	return nil
}

// @Summary 审批通过
// @Description accept workflow
// @Param workflow_id path string true "workflow ID"
// @Param workflow_step_number path string true "workflow step number"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/{workflow_step_number}/accept [post]
func AcceptWorkflow(c echo.Context) error {
	return nil
}

// @Summary 审批驳回
// @Description reject workflow
// @Param workflow_id path string true "workflow ID"
// @Param workflow_step_number path string true "workflow step number"
// @Param reason query string false "reject reason"
// @Success 200 {object} controller.BaseRes
// @router /v1/workflows/{workflow_id}/{workflow_step_number}/reject [post]
func RejectWorkflow(c echo.Context) error {
	return nil
}
