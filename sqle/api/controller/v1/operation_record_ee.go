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

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	sqleMiddleware "github.com/actiontech/sqle/sqle/api/middleware"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"

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
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return c.Param("project_name"), locale.Bundle.LocalizeAllWithArgs(locale.OprDelRuleTemplateWithName, c.Param("rule_template_name")), nil
			},
		},
		{
			RouterPath:      "/v1/projects/:project_name/rule_templates/:rule_template_name/",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeProjectRuleTemplate,
			OperationAction: model.OperationRecordActionUpdateProjectRuleTemplate,
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return c.Param("project_name"), locale.Bundle.LocalizeAllWithArgs(locale.OprEditRuleTemplateWithName, c.Param("rule_template_name")), nil
			},
		},
		// 流程模板
		{
			RouterPath:      "/v1/projects/:project_name/workflow_template",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeWorkflowTemplate,
			OperationAction: model.OperationRecordActionUpdateWorkflowTemplate,
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return c.Param("project_name"), locale.Bundle.LocalizeAll(locale.OprEditProcedureTemplate), nil
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
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return c.Param("project_name"), locale.Bundle.LocalizeAllWithArgs(locale.OprDelAuditPlanWithName, c.Param("audit_plan_name")), nil
			},
		},
		{
			RouterPath:      "/v1/projects/:project_name/audit_plans/:audit_plan_name/",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeAuditPlan,
			OperationAction: model.OperationRecordActionUpdateAuditPlan,
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return c.Param("project_name"), locale.Bundle.LocalizeAllWithArgs(locale.OprEditAuditPlanWithName, c.Param("audit_plan_name")), nil
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
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return "", locale.Bundle.LocalizeAllWithArgs(locale.OprEditGlobalRuleTemplateWithName, c.Param("rule_template_name")), nil
			},
		},
		{
			RouterPath:      "/v1/rule_templates/:rule_template_name/",
			Method:          http.MethodDelete,
			OperationType:   model.OperationRecordTypeGlobalRuleTemplate,
			OperationAction: model.OperationRecordActionDeleteGlobalRuleTemplate,
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return "", locale.Bundle.LocalizeAllWithArgs(locale.OprDelGlobalRuleTemplateWithName, c.Param("rule_template_name")), nil
			},
		},
		// 系统配置
		{
			RouterPath:      "/v1/configurations/ding_talk",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateDingTalkConfiguration,
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return "", locale.Bundle.LocalizeAll(locale.OprEditDingConfig), nil
			},
		},
		{
			RouterPath:      "/v1/configurations/system_variables",
			Method:          http.MethodPatch,
			OperationType:   model.OperationRecordTypeSystemConfiguration,
			OperationAction: model.OperationRecordActionUpdateSystemVariables,
			GetProjectAndContentFunc: func(c echo.Context) (string, i18nPkg.I18nStr, error) {
				return "", locale.Bundle.LocalizeAll(locale.OprEditGlobalConfig), nil
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

func getProjectAndContentFromCreateRuleTemplate(c echo.Context) (string, i18nPkg.I18nStr, error) {
	req := new(CreateRuleTemplateReqV1)
	if err := marshalRequestBody(c, req); err != nil {
		return "", nil, err
	}
	return "", locale.Bundle.LocalizeAllWithArgs(locale.OprAddGlobalRuleTemplateWithName, req.Name), nil
}

func getProjectAndContentFromCreatingAuditPlan(c echo.Context) (string, i18nPkg.I18nStr, error) {
	req := new(CreateAuditPlanReqV1)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", nil, err
	}
	return c.Param("project_name"), locale.Bundle.LocalizeAllWithArgs(locale.OprAddAuditPlanWithName, req.Name), nil
}

func getProjectAndContentFromCreatingProjectRuleTemplate(c echo.Context) (string, i18nPkg.I18nStr, error) {
	req := new(CreateProjectRuleTemplateReqV1)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", nil, err
	}
	return c.Param("project_name"), locale.Bundle.LocalizeAllWithArgs(locale.OprAddRuleTemplateWithName, req.Name), nil
}

func getProjectAndContentFromUpdatingFilesOrder(c echo.Context) (string, i18nPkg.I18nStr, error) {
	req := new(UpdateSqlFileOrderV1Req)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", nil, err
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
		return "", nil, err
	}

	for _, file := range auditFiles {
		newIndex := idIndexMap[file.ID]
		contents = append(contents, fmt.Sprintf("%s->%d", file.FileName, newIndex))
	}

	projectName := c.Param("project_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), projectName)
	if err != nil {
		return "", nil, err
	}
	id := c.Param("workflow_id")
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, id)
	if err != nil {
		return "", nil, fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", nil, ErrWorkflowNoAccess
	}
	content := locale.Bundle.LocalizeAllWithArgs(locale.OprUpdateFilesOrderWithOrderAndName, strings.Join(contents, "，"), workflow.Subject)
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

var typeNameDescMap = map[string]*i18n.Message{
	model.OperationRecordTypeProject:             locale.OprTypeProject,
	model.OperationRecordTypeInstance:            locale.OprTypeInstance,
	model.OperationRecordTypeProjectRuleTemplate: locale.OprTypeProjectRuleTemplate,
	model.OperationRecordTypeWorkflowTemplate:    locale.OprTypeWorkflowTemplate,
	model.OperationRecordTypeAuditPlan:           locale.OprTypeAuditPlan,
	model.OperationRecordTypeWorkflow:            locale.OprTypeWorkflow,
	model.OperationRecordTypeGlobalUser:          locale.OprTypeGlobalUser,
	model.OperationRecordTypeGlobalRuleTemplate:  locale.OprTypeGlobalRuleTemplate,
	model.OperationRecordTypeSystemConfiguration: locale.OprTypeSystemConfiguration,
	model.OperationRecordTypeProjectMember:       locale.OprTypeProjectMember,
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
			Desc:              locale.Bundle.LocalizeMsgByCtx(c.Request().Context(), typeNameDescMap[operationType]),
		})
	}

	return c.JSON(http.StatusOK, GetOperationTypeNamesListResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    operationTypeNameList,
	})
}

var actionNameDescMap = map[string]*i18n.Message{
	model.OperationRecordActionCreateProject:               locale.OprActionCreateProject,
	model.OperationRecordActionDeleteProject:               locale.OprActionDeleteProject,
	model.OperationRecordActionUpdateProject:               locale.OprActionUpdateProject,
	model.OperationRecordActionArchiveProject:              locale.OprActionArchiveProject,
	model.OperationRecordActionUnarchiveProject:            locale.OprActionUnarchiveProject,
	model.OperationRecordActionCreateInstance:              locale.OprActionCreateInstance,
	model.OperationRecordActionUpdateInstance:              locale.OprActionUpdateInstance,
	model.OperationRecordActionDeleteInstance:              locale.OprActionDeleteInstance,
	model.OperationRecordActionCreateProjectRuleTemplate:   locale.OprActionCreateProjectRuleTemplate,
	model.OperationRecordActionDeleteProjectRuleTemplate:   locale.OprActionDeleteProjectRuleTemplate,
	model.OperationRecordActionUpdateProjectRuleTemplate:   locale.OprActionUpdateProjectRuleTemplate,
	model.OperationRecordActionUpdateWorkflowTemplate:      locale.OprActionUpdateWorkflowTemplate,
	model.OperationRecordActionCreateAuditPlan:             locale.OprActionCreateAuditPlan,
	model.OperationRecordActionDeleteAuditPlan:             locale.OprActionDeleteAuditPlan,
	model.OperationRecordActionUpdateAuditPlan:             locale.OprActionUpdateAuditPlan,
	model.OperationRecordActionCreateWorkflow:              locale.OprActionCreateWorkflow,
	model.OperationRecordActionCancelWorkflow:              locale.OprActionCancelWorkflow,
	model.OperationRecordActionApproveWorkflow:             locale.OprActionApproveWorkflow,
	model.OperationRecordActionRejectWorkflow:              locale.OprActionRejectWorkflow,
	model.OperationRecordActionExecuteWorkflow:             locale.OprActionExecuteWorkflow,
	model.OperationRecordActionScheduleWorkflow:            locale.OprActionScheduleWorkflow,
	model.OperationRecordActionCreateUser:                  locale.OprActionCreateUser,
	model.OperationRecordActionUpdateUser:                  locale.OprActionUpdateUser,
	model.OperationRecordActionDeleteUser:                  locale.OprActionDeleteUser,
	model.OperationRecordActionCreateGlobalRuleTemplate:    locale.OprActionCreateGlobalRuleTemplate,
	model.OperationRecordActionUpdateGlobalRuleTemplate:    locale.OprActionUpdateGlobalRuleTemplate,
	model.OperationRecordActionDeleteGlobalRuleTemplate:    locale.OprActionDeleteGlobalRuleTemplate,
	model.OperationRecordActionUpdateDingTalkConfiguration: locale.OprActionUpdateDingTalkConfiguration,
	model.OperationRecordActionUpdateSMTPConfiguration:     locale.OprActionUpdateSMTPConfiguration,
	model.OperationRecordActionUpdateWechatConfiguration:   locale.OprActionUpdateWechatConfiguration,
	model.OperationRecordActionUpdateSystemVariables:       locale.OprActionUpdateSystemVariables,
	model.OperationRecordActionUpdateLDAPConfiguration:     locale.OprActionUpdateLDAPConfiguration,
	model.OperationRecordActionUpdateOAuth2Configuration:   locale.OprActionUpdateOAuth2Configuration,
	model.OperationRecordActionCreateMember:                locale.OprActionCreateMember,
	model.OperationRecordActionCreateMemberGroup:           locale.OprActionCreateMemberGroup,
	model.OperationRecordActionDeleteMember:                locale.OprActionDeleteMember,
	model.OperationRecordActionDeleteMemberGroup:           locale.OprActionDeleteMemberGroup,
	model.OperationRecordActionUpdateMember:                locale.OprActionUpdateMember,
	model.OperationRecordActionUpdateMemberGroup:           locale.OprActionUpdateMemberGroup,
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
			Desc:            locale.Bundle.LocalizeMsgByCtx(c.Request().Context(), actionNameDescMap[operationAction.OperationAction]),
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

	ctx := c.Request().Context()
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
			OperationTypeName: locale.Bundle.LocalizeMsgByCtx(ctx, typeNameDescMap[operationRecord.OperationTypeName]),
			OperationAction:   locale.Bundle.LocalizeMsgByCtx(ctx, actionNameDescMap[operationRecord.OperationAction]),
			OperationContent:  operationRecord.GetOperationContentByLangTag(locale.Bundle.GetLangTagFromCtx(ctx)),
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

var operationRecordStatusMap = map[string]*i18n.Message{
	model.OperationRecordStatusSucceeded: locale.OprStatusSucceeded,
	model.OperationRecordStatusFailed:    locale.OprStatusFailed,
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

	ctx := c.Request().Context()
	s := model.GetStorage()
	exportList, err := s.GetOperationRecordExportList(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	buff := new(bytes.Buffer)
	buff.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，为了兼容 windows 系统

	csvWriter := csv.NewWriter(buff)

	csvColumnNameList := []string{
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.OprOperationTime),        //"操作时间",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.OprOperationProjectName), //"项目",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.OprOperationUserName),    //"操作人",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.OprOperationAction),      //"操作对象",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.OprOperationContent),     //"操作内容",
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.OprOperationStatus),      //"状态",
	}
	err = csvWriter.Write(csvColumnNameList)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	for _, record := range exportList {
		csvLine := []string{
			record.OperationTime.Format("2006-01-02 15:04:05"),
			record.OperationProjectName,
			record.OperationUserName,
			locale.Bundle.LocalizeMsgByCtx(ctx, actionNameDescMap[record.OperationAction]),
			record.GetOperationContentByLangTag(locale.Bundle.GetLangTagFromCtx(ctx)),
			locale.Bundle.LocalizeMsgByCtx(ctx, operationRecordStatusMap[record.OperationStatus]),
		}
		err = csvWriter.Write(csvLine)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	csvWriter.Flush()

	fileName := fmt.Sprintf("%s_operation_record.csv", time.Now().Format("20060102150405"))
	c.Response().Header().Set(echo.HeaderContentDisposition, mime.FormatMediaType("attachment", map[string]string{
		"filename": fileName,
	}))

	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}
