//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	sqleMiddleware "github.com/actiontech/sqle/sqle/api/middleware"
	dms "github.com/actiontech/sqle/sqle/dms"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

func init() {
	sqleMiddleware.ApiInterfaceInfoList = append(sqleMiddleware.ApiInterfaceInfoList, []sqleMiddleware.ApiInterfaceInfo{
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
			RouterPath:      "/v1/configurations/system_variables",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateSystemVariables,
			GetProjectAndContentFunc: func(c echo.Context) (string, string, error) {
				return "", "修改全局配置", nil
			},
		},
		// 工单
		{
			RouterPath:               "/v1/projects/:project_name/workflows/:workflow_id/tasks/:task_id/order_file",
			Method:                   http.MethodPost,
			OperationType:            model.OperationRecordTypeWorkflow,
			OperationAction:          model.OperationRecordActionUpdateWorkflow,
			GetProjectAndContentFunc: getProjectAndContentFromUpdatingFilesOrder,
		},
	}...)
}

func getProjectAndContentFromCreateRuleTemplate(c echo.Context) (string, string, error) {
	req := new(CreateRuleTemplateReqV1)
	if err := marshalRequestBody(c, req); err != nil {
		return "", "", err
	}
	return "", fmt.Sprintf("创建全局规则模板，模板名：%v", req.Name), nil
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

func getProjectAndContentFromUpdatingFilesOrder(c echo.Context) (string, string, error) {
	req := new(UpdateSqlFileOrderV1Req)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}

	s := model.GetStorage()
	contents := []string{}
	fileIds := []uint{}
	idIndexMap := make(map[uint]uint)
	for _, updateFile := range req.FilesToSort {
		fileIds = append(fileIds, updateFile.FileID)
		idIndexMap[updateFile.FileID] = updateFile.NewIndex
	}

	auditFiles, err := s.GetFileByIds(fileIds)
	if err != nil {
		return "", "", err
	}

	for _, file := range auditFiles {
		newIndex := idIndexMap[file.ID]
		contents = append(contents, fmt.Sprintf("将%s->%d", file.FileName, newIndex))
	}

	projectName := c.Param("project_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), projectName)
	if err != nil {
		return "", "", err
	}
	id := c.Param("workflow_id")
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", ErrWorkflowNoAccess
	}
	content := "文件上线顺序调整：" + strings.Join(contents, "，") + fmt.Sprintf("，工单名称：%s", workflow.Subject)
	return projectName, content, nil
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
	model.OperationRecordTypeProjectMember:       "项目成员",
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
	model.OperationRecordActionArchiveProject:              "冻结项目",
	model.OperationRecordActionUnarchiveProject:            "取消冻结项目",
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
	model.OperationRecordActionCreateMember:                "添加成员",
	model.OperationRecordActionCreateMemberGroup:           "添加成员组",
	model.OperationRecordActionDeleteMember:                "删除成员",
	model.OperationRecordActionDeleteMemberGroup:           "删除成员组",
	model.OperationRecordActionUpdateMember:                "编辑成员",
	model.OperationRecordActionUpdateMemberGroup:           "编辑成员组",
}

func getOperationActionList(c echo.Context) error {
	type action struct {
		OperationType   string
		OperationAction string
	}
	var operationActionList []action
	removeDuplicate := make(map[string]struct{})
	for _, info := range sqleMiddleware.ApiInterfaceInfoList {
		if _, ok := removeDuplicate[info.OperationAction]; ok {
			continue
		}
		removeDuplicate[info.OperationAction] = struct{}{}
		operationActionList = append(operationActionList, action{
			info.OperationType,
			info.OperationAction,
		})
	}

	var operationActionNameList []OperationActionList
	for _, operationAction := range operationActionList {
		operationActionNameList = append(operationActionNameList, OperationActionList{
			OperationType:   operationAction.OperationType,
			OperationAction: operationAction.OperationAction,
			Desc:            actionNameDescMap[operationAction.OperationAction],
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

var operationRecordStatusMap = map[string]string{
	model.OperationRecordStatusSucceeded: "成功",
	model.OperationRecordStatusFailed:    "失败",
}

func exportOperationRecordList(c echo.Context) error {
	req := new(GetExportOperationRecordListReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := map[string]interface{}{
		"filter_operate_time_from":       req.FilterOperateTimeFrom,
		"filter_operate_time_to":         req.FilterOperateTimeTo,
		"fuzzy_search_operate_user_name": req.FuzzySearchOperateUserName,
		"filter_operate_type_name":       req.FilterOperateTypeName,
		"filter_operate_action":          req.FilterOperateAction,
	}
	if req.FilterOperateProjectName != nil {
		data["filter_operate_project_name"] = req.FilterOperateProjectName
	}

	s := model.GetStorage()
	exportList, err := s.GetOperationRecordExportList(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	buff := new(bytes.Buffer)
	buff.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，为了兼容 windows 系统

	csvWriter := csv.NewWriter(buff)

	csvColumnNameList := []string{"操作时间", "项目", "操作人", "操作对象", "操作内容", "状态"}
	err = csvWriter.Write(csvColumnNameList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	for _, record := range exportList {
		csvLine := []string{
			record.OperationTime.Format("2006-01-02 15:04:05"),
			record.OperationProjectName,
			record.OperationUserName,
			actionNameDescMap[record.OperationAction],
			record.OperationContent,
			operationRecordStatusMap[record.OperationStatus],
		}
		err = csvWriter.Write(csvLine)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	csvWriter.Flush()

	fileName := fmt.Sprintf("%s_操作记录.csv", time.Now().Format("20060102150405"))
	c.Response().Header().Set(echo.HeaderContentDisposition, mime.FormatMediaType("attachment", map[string]string{
		"filename": fileName,
	}))

	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}
