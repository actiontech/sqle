//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/api/controller"
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
	{
		RouterPath:               "/v1/projects",
		Method:                   http.MethodPost,
		OperationType:            model.OperationRecordTypeProject,
		OperationAction:          model.OperationRecordActionCreateProject,
		GetProjectAndContentFunc: getProjectAndContentFromCreateProject,
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
		"filter_operate_project_name":    req.FilterOperateProjectName,
		"fuzzy_search_operate_user_name": req.FuzzySearchOperateUserName,
		"filter_operate_type_name":       req.FilterOperateTypeName,
		"filter_operate_action":          req.FilterOperateAction,
		"limit":                          req.PageSize,
		"offset":                         offset,
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
