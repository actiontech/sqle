//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sqleMiddleware "github.com/actiontech/sqle/sqle/api/middleware"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

func init() {
	sqleMiddleware.ApiInterfaceInfoList = append(sqleMiddleware.ApiInterfaceInfoList, []sqleMiddleware.ApiInterfaceInfo{
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
		// 平台用户
		{
			RouterPath:               "/v1/users",
			Method:                   http.MethodPost,
			OperationType:            model.OperationRecordTypeGlobalUser,
			OperationAction:          model.OperationRecordActionCreateUser,
			GetProjectAndContentFunc: getProjectAndContentFromCreateUser,
		},
		{
			RouterPath:      "/v1/users/:user_name/",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeGlobalUser,
			OperationAction: model.OperationRecordActionUpdateUser,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", fmt.Sprintf("编辑用户，用户名：%v", c.Param("user_name")), nil
			},
		},
		{
			RouterPath:      "/v1/users/:user_name/",
			Method:          http.MethodDelete,
			OperationType:   model.OperationRecordTypeGlobalUser,
			OperationAction: model.OperationRecordActionDeleteUser,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", fmt.Sprintf("删除用户，用户名：%v", c.Param("user_name")), nil
			},
		},
		// 全局规则模板
		{
			RouterPath:               "/v1/rule_templates",
			Method:                   http.MethodPost,
			OperationType:            model.OperationRecordTypeGlobalRuleTemplate,
			OperationAction:          model.OperationRecordActionCreateGlobalRuleTemplate,
			GetProjectAndContentFunc: getProjectAndContentFromCreateRuleTemplate,
		},
		{
			RouterPath:      "/v1/rule_templates/:rule_template_name/",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeGlobalRuleTemplate,
			OperationAction: model.OperationRecordActionUpdateGlobalRuleTemplate,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", fmt.Sprintf("编辑全局规则模板，模板名：%v", c.Param("rule_template_name")), nil
			},
		},
		{
			RouterPath:      "/v1/rule_templates/:rule_template_name/",
			Method:          http.MethodDelete,
			OperationType:   model.OperationRecordTypeGlobalRuleTemplate,
			OperationAction: model.OperationRecordActionDeleteGlobalRuleTemplate,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", fmt.Sprintf("删除全局规则模板，模板名：%v", c.Param("rule_template_name")), nil
			},
		},
		// 系统配置
		{
			RouterPath:      "/v1/configurations/ding_talk",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateDingTalkConfiguration,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", "修改钉钉配置", nil
			},
		},
		{
			RouterPath:      "/v1/configurations/smtp",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateSMTPConfiguration,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", "修改SMTP配置", nil
			},
		},
		{
			RouterPath:      "/v1/configurations/wechat",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateWechatConfiguration,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", "修改微信配置", nil
			},
		},
		{
			RouterPath:      "/v1/configurations/system_variables",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateSystemVariables,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", "修改全局配置", nil
			},
		},
		{
			RouterPath:      "/v1/configurations/ldap",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateLDAPConfiguration,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", "修改LDAP配置", nil
			},
		},
		{
			RouterPath:      "/v1/configurations/oauth2",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateOAuth2Configuration,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", "修改OAuth2配置", nil
			},
		},
		// 成员
		{
			RouterPath:               "/v1/projects/:project_name/members",
			Method:                   http.MethodPost,
			OperationType:            model.OperationRecordTypeMember,
			OperationAction:          model.OperationRecordActionCreateMember,
			GetProjectAndContentFunc: getProjectAndContentFromCreateMember,
		},
		{
			RouterPath:               "/v1/projects/:project_name/member_groups",
			Method:                   http.MethodPost,
			OperationType:            model.OperationRecordTypeMember,
			OperationAction:          model.OperationRecordActionCreateMemberGroup,
			GetProjectAndContentFunc: getProjectAndContentFromCreateMemberGroup,
		},
		{
			RouterPath:      "/v1/projects/:project_name/members/:user_name/",
			Method:          http.MethodDelete,
			OperationType:   model.OperationRecordTypeMember,
			OperationAction: model.OperationRecordActionDeleteMember,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return c.Param("project_name"), fmt.Sprintf("删除成员，用户名：%v", c.Param("user_name")), nil
			},
		},
		{
			RouterPath:      "/v1/projects/:project_name/member_groups/:user_group_name/",
			Method:          http.MethodDelete,
			OperationType:   model.OperationRecordTypeMember,
			OperationAction: model.OperationRecordActionDeleteMemberGroup,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return c.Param("project_name"), fmt.Sprintf("删除成员组，组名：%v", c.Param("user_group_name")), nil
			},
		},
		{
			RouterPath:      "/v1/projects/:project_name/members/:user_name/",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeMember,
			OperationAction: model.OperationRecordActionUpdateMember,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return c.Param("project_name"), fmt.Sprintf("编辑成员，用户名：%v", c.Param("user_name")), nil
			},
		},
		{
			RouterPath:      "/v1/projects/:project_name/member_groups/:user_group_name/",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeMember,
			OperationAction: model.OperationRecordActionUpdateMemberGroup,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return c.Param("project_name"), fmt.Sprintf("编辑成员组，组名：%v", c.Param("user_group_name")), nil
			},
		},
	}...)
}

func getProjectAndContentFromCreateMemberGroup(c echo.Context) (string, string, error) {
	req := new(CreateMemberGroupReqV1)
	if err := marshalRequestBody(c, req); err != nil {
		return "", "", err
	}
	return c.Param("project_name"), fmt.Sprintf("创建成员组，组名：%v", req.UserGroupName), nil
}

func getProjectAndContentFromCreateMember(c echo.Context) (string, string, error) {
	req := new(CreateMemberReqV1)
	if err := marshalRequestBody(c, req); err != nil {
		return "", "", err
	}
	return c.Param("project_name"), fmt.Sprintf("添加成员，用户名：%v", req.UserName), nil
}

func getProjectAndContentFromCreateRuleTemplate(c echo.Context) (string, string, error) {
	req := new(CreateRuleTemplateReqV1)
	if err := marshalRequestBody(c, req); err != nil {
		return "", "", err
	}
	return "", fmt.Sprintf("创建全局规则模板，模板名：%v", req.Name), nil
}

func getProjectAndContentFromCreateUser(c echo.Context) (string, string, error) {
	req := new(CreateUserReqV1)
	if err := marshalRequestBody(c, req); err != nil {
		return "", "", err
	}
	return "", fmt.Sprintf("创建用户，用户名：%v", req.Name), nil
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
	model.OperationRecordTypeGlobalUser:          "平台用户",
	model.OperationRecordTypeGlobalRuleTemplate:  "全局规则模板",
	model.OperationRecordTypeSystemConfiguration: "系统配置",
	model.OperationRecordTypeMember:              "成员",
}

func getOperationTypeNameList(c echo.Context) error {
	var operationTypeList []string
	for _, info := range sqleMiddleware.ApiInterfaceInfoList {
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
	model.OperationRecordActionCreateProject:               "创建项目",
	model.OperationRecordActionDeleteProject:               "删除项目",
	model.OperationRecordActionUpdateProject:               "编辑项目",
	model.OperationRecordActionCreateInstance:              "创建数据源",
	model.OperationRecordActionUpdateInstance:              "编辑数据源",
	model.OperationRecordActionDeleteInstance:              "删除数据源",
	model.OperationRecordActionCreateProjectRuleTemplate:   "添加规则模版",
	model.OperationRecordActionDeleteProjectRuleTemplate:   "删除规则模版",
	model.OperationRecordActionUpdateProjectRuleTemplate:   "编辑规则模版",
	model.OperationRecordActionUpdateWorkflowTemplate:      "编辑流程模版",
	model.OperationRecordActionCreateAuditPlan:             "创建智能扫描任务",
	model.OperationRecordActionDeleteAuditPlan:             "删除智能扫描任务",
	model.OperationRecordActionUpdateAuditPlan:             "编辑智能扫描任务",
	model.OperationRecordActionCreateWorkflow:              "创建工单",
	model.OperationRecordActionCancelWorkflow:              "关闭工单",
	model.OperationRecordActionApproveWorkflow:             "审核通过工单",
	model.OperationRecordActionRejectWorkflow:              "驳回工单",
	model.OperationRecordActionExecuteWorkflow:             "上线工单",
	model.OperationRecordActionScheduleWorkflow:            "定时上线",
	model.OperationRecordActionCreateUser:                  "创建用户",
	model.OperationRecordActionUpdateUser:                  "编辑用户",
	model.OperationRecordActionDeleteUser:                  "删除用户",
	model.OperationRecordActionCreateGlobalRuleTemplate:    "创建全局规则模版",
	model.OperationRecordActionUpdateGlobalRuleTemplate:    "编辑全局规则模版",
	model.OperationRecordActionDeleteGlobalRuleTemplate:    "删除全局规则模版",
	model.OperationRecordActionUpdateDingTalkConfiguration: "修改钉钉配置",
	model.OperationRecordActionUpdateSMTPConfiguration:     "修改SMTP配置",
	model.OperationRecordActionUpdateWechatConfiguration:   "修改微信配置",
	model.OperationRecordActionUpdateSystemVariables:       "修改系统变量",
	model.OperationRecordActionUpdateLDAPConfiguration:     "修改LDAP配置",
	model.OperationRecordActionUpdateOAuth2Configuration:   "修改OAuth2配置",
	model.OperationRecordActionCreateMember:                "创建成员",
	model.OperationRecordActionCreateMemberGroup:           "创建成员组",
	model.OperationRecordActionDeleteMember:                "删除成员",
	model.OperationRecordActionDeleteMemberGroup:           "删除成员组",
	model.OperationRecordActionUpdateMember:                "编辑成员",
	model.OperationRecordActionUpdateMemberGroup:           "编辑成员组",
}

func getOperationActionList(c echo.Context) error {
	var operationActionList []string
	for _, info := range sqleMiddleware.ApiInterfaceInfoList {
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
