//go:build enterprise
// +build enterprise

package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type apiInterfaceInfo struct {
	routerPath               string
	method                   string
	operationType            string
	operationAction          string
	getProjectAndContentFunc func(c echo.Context) (projectName, content string, err error)
}

var apiInterfaceInfoList = []apiInterfaceInfo{
	// 项目
	{
		routerPath:               "/v1/projects",
		method:                   http.MethodPost,
		operationType:            model.OperationRecordTypeProject,
		operationAction:          model.OperationRecordActionCreateProject,
		getProjectAndContentFunc: getProjectAndContentFromCreateProject,
	},
	// 项目规则模板
	{
		routerPath:               "/v1/projects/:project_name/rule_templates",
		method:                   http.MethodPost,
		operationType:            model.OperationRecordTypeProjectRuleTemplate,
		operationAction:          model.OperationRecordActionCreateProjectRuleTemplate,
		getProjectAndContentFunc: getProjectAndContentFromCreatingProjectRuleTemplate,
	},
	{
		routerPath:      "/v1/projects/:project_name/rule_templates/:rule_template_name/",
		method:          http.MethodDelete,
		operationType:   model.OperationRecordTypeProjectRuleTemplate,
		operationAction: model.OperationRecordActionDeleteProjectRuleTemplate,
		getProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("删除规则模板，模板名：%v", c.Param("rule_template_name")), nil
		},
	},
	{
		routerPath:      "/v1/projects/:project_name/rule_templates/:rule_template_name/",
		method:          http.MethodPatch,
		operationType:   model.OperationRecordTypeProjectRuleTemplate,
		operationAction: model.OperationRecordActionUpdateProjectRuleTemplate,
		getProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("编辑规则模板，模板名：%v", c.Param("rule_template_name")), nil
		},
	},
	// 流程模板
	{
		routerPath:      "/v1/projects/:project_name/workflow_template",
		method:          http.MethodPatch,
		operationType:   model.OperationRecordTypeWorkflowTemplate,
		operationAction: model.OperationRecordActionUpdateWorkflowTemplate,
		getProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), "编辑流程模板", nil
		},
	},
	// 智能扫描
	{
		routerPath:               "/v1/projects/:project_name/audit_plans",
		method:                   http.MethodPost,
		operationType:            model.OperationRecordTypeAuditPlan,
		operationAction:          model.OperationRecordActionCreateAuditPlan,
		getProjectAndContentFunc: getProjectAndContentFromCreatingAuditPlan,
	},
	{
		routerPath:      "/v1/projects/:project_name/audit_plans/:audit_plan_name/",
		method:          http.MethodDelete,
		operationType:   model.OperationRecordTypeAuditPlan,
		operationAction: model.OperationRecordActionDeleteAuditPlan,
		getProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("删除智能扫描任务，任务名：%v", c.Param("audit_plan_name")), nil
		},
	},
	{
		routerPath:      "/v1/projects/:project_name/audit_plans/:audit_plan_name/",
		method:          http.MethodPatch,
		operationType:   model.OperationRecordTypeAuditPlan,
		operationAction: model.OperationRecordActionUpdateAuditPlan,
		getProjectAndContentFunc: func(c echo.Context) (string, string, error) {
			return c.Param("project_name"), fmt.Sprintf("编辑智能扫描任务，任务名：%v", c.Param("audit_plan_name")), nil
		},
	},
}

func getProjectAndContentFromCreatingAuditPlan(c echo.Context) (string, string, error) {
	req := new(v1.CreateAuditPlanReqV1)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}
	return c.Param("project_name"), fmt.Sprintf("创建智能扫描任务，任务名：%v", req.Name), nil
}

func getProjectAndContentFromCreatingProjectRuleTemplate(c echo.Context) (string, string, error) {
	req := new(v1.CreateProjectRuleTemplateReqV1)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}
	return c.Param("project_name"), fmt.Sprintf("添加规则模板，模板名：%v", req.Name), nil
}

func getProjectAndContentFromCreateProject(c echo.Context) (string, string, error) {
	req := new(v1.CreateProjectReqV1)
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

type ResponseBodyWrite struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseBodyWrite) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ResponseBodyWrite) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.Write([]byte(s))
}

func OperationLogRecord() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			reqIP := c.Request().Host
			path := c.Path()
			newLog := log.NewEntry()
			for _, interfaceInfo := range apiInterfaceInfoList {
				if c.Request().Method == interfaceInfo.method && interfaceInfo.routerPath == path {
					userName := controller.GetUserName(c)

					operationRecord := &model.OperationRecord{
						OperationTime:     time.Now(),
						OperationUserName: userName,
						OperationReqIP:    reqIP,
						OperationTypeName: interfaceInfo.operationType,
						OperationAction:   interfaceInfo.operationAction,
					}

					projectName, content, err := interfaceInfo.getProjectAndContentFunc(c)
					if err != nil {
						newLog.Errorf("get content and project name error: %s", err)
					}

					operationRecord.OperationProjectName = projectName
					operationRecord.OperationContent = content

					respBodyWrite := &ResponseBodyWrite{body: new(bytes.Buffer), ResponseWriter: c.Response().Writer}

					c.Response().Writer = respBodyWrite

					if err = next(c); err != nil {
						c.Error(err)
					}

					resp := respBodyWrite.body.Bytes()
					var respBody map[string]interface{}
					if err := json.Unmarshal(resp, &respBody); err == nil {
						if code, ok := respBody["code"]; ok {
							codeInt := int(code.(float64))
							if codeInt != 0 {
								operationRecord.OperationStatus = model.OperationRecordStatusFail
							} else {
								operationRecord.OperationStatus = model.OperationRecordStatusSuccess
							}
						}
					} else {
						operationRecord.OperationStatus = model.OperationRecordStatusFail
					}

					s := model.GetStorage()
					if err := s.Save(&operationRecord); err != nil {
						newLog.Errorf("save operation record error: %s", err)
						return nil
					}

					return nil
				}
			}

			return next(c)
		}
	}
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
