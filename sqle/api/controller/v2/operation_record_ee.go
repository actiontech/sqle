//go:build enterprise
// +build enterprise

package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	sqleMiddleware "github.com/actiontech/sqle/sqle/api/middleware"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func init() {
	sqleMiddleware.ApiInterfaceInfoList = append(sqleMiddleware.ApiInterfaceInfoList, []sqleMiddleware.ApiInterfaceInfo{
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
		{
			RouterPath:               "/v2/projects/:project_name/workflows/:workflow_id/tasks/:task_id/schedule", // 设置定时上线
			Method:                   http.MethodPut,
			OperationType:            model.OperationRecordTypeWorkflow,
			OperationAction:          model.OperationRecordActionScheduleWorkflow,
			GetProjectAndContentFunc: getProjectAndContentFromSchedulingWorkflow,
		},
	}...)
}

func getProjectAndContentFromSchedulingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), projectName)
	if err != nil {
		return "", "", err
	}
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", v1.ErrWorkflowNoAccess
	}

	taskId := c.Param("task_id")
	task, err := v1.GetTaskById(c.Request().Context(), taskId)
	if err != nil {
		return "", "", fmt.Errorf("get task failed: %v", err)
	}

	req := new(UpdateWorkflowScheduleReqV2)
	err = marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}

	if req.ScheduleTime != nil {
		return projectName, fmt.Sprintf("设置定时上线，工单名称：%v, 数据源名: %v", workflow.Subject, task.InstanceName()), nil
	} else {
		return projectName, fmt.Sprintf("取消定时上线，工单名称：%v, 数据源名: %v", workflow.Subject, task.InstanceName()), nil
	}
}

func getProjectAndContentFromBatchExecutingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), projectName)
	if err != nil {
		return "", "", err
	}

	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", v1.ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("上线工单，工单名称：%v", workflow.Subject), nil
}

func getProjectAndContentFromExecutingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), projectName)
	if err != nil {
		return "", "", err
	}

	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", v1.ErrWorkflowNoAccess
	}

	taskId := c.Param("task_id")
	task, err := v1.GetTaskById(context.Background(), taskId)
	if err != nil {
		return "", "", fmt.Errorf("get task failed: %v", err)
	}

	return projectName, fmt.Sprintf("上线工单的单个数据源, 工单名称：%v, 数据源名: %v", workflow.Subject, task.InstanceName()), nil
}

func getProjectAndContentFromBatchCancelingWorkflow(c echo.Context) (string, string, error) {
	req := new(BatchCancelWorkflowsReqV2)
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
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), projectName)
	if err != nil {
		return "", "", err
	}
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", v1.ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("取消工单，工单名称：%v", workflow.Subject), nil
}

func getProjectAndContentFromApprovingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), projectName)
	if err != nil {
		return "", "", err
	}
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", v1.ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("审核通过工单，工单名称：%v", workflow.Subject), nil
}

func getProjectAndContentFromRejectingWorkflow(c echo.Context) (string, string, error) {
	projectName := c.Param("project_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), projectName)
	if err != nil {
		return "", "", err
	}
	id := c.Param("workflow_id")
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, id)
	if err != nil {
		return "", "", fmt.Errorf("get workflow failed: %v", err)
	}
	if !exist {
		return "", "", v1.ErrWorkflowNoAccess
	}
	return projectName, fmt.Sprintf("驳回工单，工单名称：%v", workflow.Subject), nil
}

func getProjectAndContentFromCreatingWorkflow(c echo.Context) (string, string, error) {
	req := new(CreateWorkflowReqV2)
	err := marshalRequestBody(c, req)
	if err != nil {
		return "", "", err
	}
	return c.Param("project_name"), fmt.Sprintf("创建工单，工单名：%v", req.Subject), nil
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
