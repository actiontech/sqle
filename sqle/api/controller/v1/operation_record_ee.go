//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v2 "github.com/actiontech/sqle/sqle/api/controller/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

type ApiInterfaceInfo struct {
	RouterPath               string
	Method                   string
	OperationType            string
	OperationAction          string
	GetProjectAndContentFunc func(c echo.Context) (projectName, objectName string, err error)
}

var ApiInterfaceInfoList = []ApiInterfaceInfo{
	// 项目
	{
		RouterPath:               "/v1/projects",
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeProject,
		OperationAction:          model.OperationRecordActionCreateProject,
		GetProjectAndContentFunc: getProjectAndContentFromCreateProject,
	},
	{
		RouterPath:      "/v1/projects/:project_name/",
		Method:          http.MethodPatch,
		OperationType:   model.OperationRecordTypeProject,
		OperationAction: model.OperationRecordActionUpdateProject,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			name := c.Param("project_name")
			return name, fmt.Sprintf("编辑项目，项目名：%v", name), nil
		},
	},
	{
		RouterPath:      "/v1/projects/:project_name/",
		Method:          http.MethodPost,
		OperationType:   model.OperationRecordTypeProject,
		OperationAction: model.OperationRecordActionDeleteProject,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			name := c.Param("project_name")
			return name, fmt.Sprintf("删除项目，项目名：%v", name), nil
		},
	},
	// 数据源
	{
		RouterPath:               "/v2/projects/:project_name/instances",
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeInstance,
		OperationAction:          model.OperationRecordActionCreateInstance,
		GetProjectAndContentFunc: getProjectAndContentFromCreatingInstance,
	},
	{
		RouterPath:      "/v1/projects/:project_name/instances/:instance_name/",
		Method:          http.MethodPatch,
		OperationType:   model.OperationRecordTypeInstance,
		OperationAction: model.OperationRecordActionUpdateInstance,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("编辑数据源，数据源名称：%v", c.Param("instance_name")), nil
		},
	},
	{
		RouterPath:      "/v1/projects/:project_name/instances/:instance_name/",
		Method:          http.MethodDelete,
		OperationType:   model.OperationRecordTypeInstance,
		OperationAction: model.OperationRecordActionDeleteInstance,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("删除数据源，数据源名称：%v", c.Param("instance_name")), nil
		},
	},
	// 项目规则模板
	{
		RouterPath:               "/v1/projects/:project_name/rule_templates",
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeProjectRuleTemplate,
		OperationAction:          model.OperationRecordActionCreateProjectRuleTemplate,
		GetProjectAndContentFunc: getProjectAndContentFromCreatingProjectRuleTemplate,
	},
	{
		RouterPath:      "/v1/projects/:project_name/rule_templates/:rule_template_name/",
		Method:          http.MethodDelete,
		OperationType:   model.OperationRecordTypeProjectRuleTemplate,
		OperationAction: model.OperationRecordActionDeleteProjectRuleTemplate,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("删除规则模板，模板名：%v", c.Param("rule_template_name")), nil
		},
	},
	{
		RouterPath:      "/v1/projects/:project_name/rule_templates/:rule_template_name/",
		Method:          http.MethodPatch,
		OperationType:   model.OperationRecordTypeProjectRuleTemplate,
		OperationAction: model.OperationRecordActionUpdateProjectRuleTemplate,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("编辑规则模板，模板名：%v", c.Param("rule_template_name")), nil
		},
	},
	// 流程模板
	{
		RouterPath:      "/v1/projects/:project_name/workflow_template",
		Method:          http.MethodPatch,
		OperationType:   model.OperationRecordTypeWorkflowTemplate,
		OperationAction: model.OperationRecordActionUpdateWorkflowTemplate,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), "编辑流程模板", nil
		},
	},
	// 智能扫描
	{
		RouterPath:               "/v1/projects/:project_name/audit_plans",
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeAuditPlan,
		OperationAction:          model.OperationRecordActionCreateAuditPlan,
		GetProjectAndContentFunc: getProjectAndContentFromCreatingAuditPlan,
	},
	{
		RouterPath:      "/v1/projects/:project_name/audit_plans/:audit_plan_name/",
		Method:          http.MethodDelete,
		OperationType:   model.OperationRecordTypeAuditPlan,
		OperationAction: model.OperationRecordActionDeleteAuditPlan,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("删除智能扫描任务，任务名：%v", c.Param("audit_plan_name")), nil
		},
	},
	{
		RouterPath:      "/v1/projects/:project_name/audit_plans/:audit_plan_name/",
		Method:          http.MethodPatch,
		OperationType:   model.OperationRecordTypeAuditPlan,
		OperationAction: model.OperationRecordActionUpdateAuditPlan,
		GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("编辑智能扫描任务，任务名：%v", c.Param("audit_plan_name")), nil
		},
	},
	{
		RouterPath:               "/v2/projects/:project_name/workflows",
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeWorkflow,
		OperationAction:          model.OperationRecordActionCreateWorkflow,
		GetProjectAndContentFunc: getProjectAndContentFromCreatingWorkflow,
	},
	{
		RouterPath:               "/v2/projects/:project_name/workflows/:workflow_id/cancel", // 取消工单
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeWorkflow,
		OperationAction:          model.OperationRecordActionCancelWorkflow,
		GetProjectAndContentFunc: getProjectAndContentFromCancelingWorkflow,
	},
	{
		RouterPath:               "/v2/projects/:project_name/workflows/cancel", // 批量取消工单
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeWorkflow,
		OperationAction:          model.OperationRecordActionCancelWorkflow,
		GetProjectAndContentFunc: getProjectAndContentFromBatchCancelingWorkflow,
	},
	{
		RouterPath:               "/v2/projects/:project_name/workflows/:workflow_id/steps/:workflow_step_id/approve",
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeWorkflow,
		OperationAction:          model.OperationRecordActionApproveWorkflow,
		GetProjectAndContentFunc: getProjectAndContentFromApprovingWorkflow,
	},
	{
		RouterPath:               "/v2/projects/:project_name/workflows/:workflow_id/steps/:workflow_step_id/reject",
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeWorkflow,
		OperationAction:          model.OperationRecordActionRejectWorkflow,
		GetProjectAndContentFunc: getProjectAndContentFromRejectingWorkflow,
	},
	{
		RouterPath:               "/v2/projects/:project_name/workflows/:workflow_id/tasks/:task_id/execute", // 上线单个数据源
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeWorkflow,
		OperationAction:          model.OperationRecordActionExecuteWorkflow,
		GetProjectAndContentFunc: getProjectAndContentFromExecutingWorkflow,
	},
	{
		RouterPath:               "/v2/projects/:project_name/workflows/:workflow_id/tasks/execute", //多数据源上线
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeWorkflow,
		OperationAction:          model.OperationRecordActionExecuteWorkflow,
		GetProjectAndContentFunc: getProjectAndContentFromBatchExecutingWorkflow,
	},
}

func getProjectAndContentFromCreatingInstance(c echo.Context) (string, string, error) {
	req := new(v2.CreateInstanceReqV2)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}
	projectName := c.Param("project_name")
	return projectName, fmt.Sprintf("添加数据源，数据源名称：%v", req.Name), nil
}

func getProjectAndContentFromBatchExecutingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("上线工单，工单名称：%v", workflow.Subject), nil
}

func getProjectAndContentFromExecutingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("上线工单的单个数据源，工单名称：%v", workflow.Subject), nil // todo issue1281 添加数据源名称到记录里
}

func getProjectAndContentFromBatchCancelingWorkflow(c echo.Context) (string, string, error) {
	req := new(v2.BatchCancelWorkflowsReqV2)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}
	projectName := c.Param("project_name")
	workflowNames, err := model.GetStorage().GetWorkflowNamesByIDs(req.WorkflowIDList)
	if err != nil {
		return "", "", err
	}
	return projectName, fmt.Sprintf("批量取消工单，工单名称：%v", workflowNames), nil
}

func getProjectAndContentFromCancelingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("取消工单，工单名称：%v", workflow.Subject), nil
}

func getProjectAndContentFromApprovingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("审核通过工单，工单名称：%v", workflow.Subject), nil
}

func getProjectAndContentFromRejectingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("驳回工单，工单名称：%v", workflow.Subject), nil
}

func getProjectAndContentFromCreatingWorkflow(c echo.Context) (string, string, error) {
	req := new(v2.CreateWorkflowReqV2)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}
	return c.Param("project_name"), fmt.Sprintf("创建工单，工单名：%v", req.Subject), nil
}

func getProjectAndContentFromCreatingAuditPlan(c echo.Context) (string, string, error) {
	req := new(CreateAuditPlanReqV1)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}
	return c.Param("project_name"), fmt.Sprintf("创建智能扫描任务，任务名：%v", req.Name), nil
}

func getProjectAndContentFromCreatingProjectRuleTemplate(c echo.Context) (string, string, error) {
	req := new(CreateProjectRuleTemplateReqV1)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}
	return c.Param("project_name"), fmt.Sprintf("添加规则模板，模板名：%v", req.Name), nil
}

func getProjectAndContentFromCreateProject(c echo.Context) (string, string, error) {
	req := new(CreateProjectReqV1)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}

	return req.Name, fmt.Sprintf("创建项目，项目名：%v", req.Name), nil
}

func marshalRequestBody(c echo.Context, pattern interface{}) error {
	reqBody, err := getReqBodyBytes(c)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(reqBody, pattern); err != nil {
		return err
	}

	if err := controller.Validate(pattern); err != nil {
		return err
	}
	return nil
}

func getReqBodyBytes(c echo.Context) ([]byte, error) {
	var bodyBytes []byte
	var err error

	if c.Request().Body != nil {
		bodyBytes, err = ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return nil, err
		}

		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		return bodyBytes, nil
	}

	return nil, fmt.Errorf("request body is nil")
}

var typeNameDescMap = map[string]string{
	model.OperationRecordTypeProject:             "项目",
	model.OperationRecordTypeInstance:            "数据源",
	model.OperationRecordTypeProjectRuleTemplate: "项目规则模板",
	model.OperationRecordTypeWorkflowTemplate:    "流程模板",
	model.OperationRecordTypeAuditPlan:           "智能扫描任务",
	model.OperationRecordTypeWorkflow:            "工单",
}

func getOperationTypeNameList(c echo.Context) error {
	var operationTypeList []string
	for _, info := range ApiInterfaceInfoList {
		operationTypeList = append(operationTypeList, info.OperationType)
	}

	distinctOperationTypeList := utils.RemoveDuplicate(operationTypeList)

	var operationTypeNameList []OperationTypeNameList
	for _, operationType := range distinctOperationTypeList {
		operationTypeNameList = append(operationTypeNameList, OperationTypeNameList{
			OperationTypeName: operationType,
			Desc:              typeNameDescMap[operationType],
		})
	}

	return c.JSON(http.StatusOK, GetOperationTypeNamesListResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    operationTypeNameList,
	})
}

var actionNameDescMap = map[string]string{
	model.OperationRecordActionCreateProject:             "创建项目",
	model.OperationRecordActionCreateProjectRuleTemplate: "添加规则模版",
	model.OperationRecordActionDeleteProjectRuleTemplate: "删除规则模版",
	model.OperationRecordActionUpdateProjectRuleTemplate: "编辑规则模版",
	model.OperationRecordActionUpdateWorkflowTemplate:    "编辑流程模版",
	model.OperationRecordActionCreateAuditPlan:           "创建智能扫描任务",
	model.OperationRecordActionDeleteAuditPlan:           "删除智能扫描任务",
	model.OperationRecordActionUpdateAuditPlan:           "编辑智能扫描任务",
}

func getOperationActionList(c echo.Context) error {
	var operationActionList []string
	for _, info := range ApiInterfaceInfoList {
		operationActionList = append(operationActionList, info.OperationAction)
	}

	distinctOperationActionList := utils.RemoveDuplicate(operationActionList)

	var operationActionNameList []OperationActionList
	for _, operationAction := range distinctOperationActionList {
		operationActionNameList = append(operationActionNameList, OperationActionList{
			OperationAction: operationAction,
			Desc:            actionNameDescMap[operationAction],
		})
	}

	return c.JSON(http.StatusOK, GetOperationActionListResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    operationActionNameList,
	})
}

func getOperationRecordList(c echo.Context) error {
	req := new(GetOperationRecordListReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex > 0 {
		offset = (req.PageIndex - 1) * req.PageSize
	}

	data := map[string]interface{}{
		"filter_operate_time_from":       req.FilterOperateTimeFrom,
		"filter_operate_time_to":         req.FilterOperateTimeTo,
		"fuzzy_search_operate_user_name": req.FuzzySearchOperateUserName,
		"filter_operate_type_name":       req.FilterOperateTypeName,
		"filter_operate_action":          req.FilterOperateAction,
		"limit":                          req.PageSize,
		"offset":                         offset,
	}

	if req.FilterOperateProjectName != nil {
		data["filter_operate_project_name"] = req.FilterOperateProjectName
	}

	s := model.GetStorage()
	operationRecordList, count, err := s.GetOperationRecordList(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var operationRecordListRes []OperationRecordList
	for _, operationRecord := range operationRecordList {
		operationRecordListRes = append(operationRecordListRes, OperationRecordList{
			ID:            uint64(operationRecord.ID),
			OperationTime: &operationRecord.OperationTime,
			OperationUser: OperationUser{
				UserName: operationRecord.OperationUserName,
				IP:       operationRecord.OperationReqIP,
			},
			OperationTypeName: typeNameDescMap[operationRecord.OperationTypeName],
			OperationAction:   actionNameDescMap[operationRecord.OperationAction],
			OperationContent:  operationRecord.OperationContent,
			ProjectName:       operationRecord.OperationProjectName,
			Status:            operationRecord.OperationStatus,
		})
	}

	return c.JSON(http.StatusOK, GetOperationRecordListResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      operationRecordListRes,
		TotalNums: count,
	})
}
