package v3

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/api/controller"
	v2 "github.com/actiontech/sqle/sqle/api/controller/v2"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/sqlversion"
	"github.com/labstack/echo/v4"
)

type BatchCompleteWorkflowsReqV3 struct {
	WorkflowList []*CompleteWorkflowReq `json:"workflow_list" form:"workflow_list"`
}

type CompleteWorkflowReq struct {
	WorkflowID string  `json:"workflow_id" form:"workflow_id"`
	Desc       *string `json:"desc" form:"desc"`
}

// BatchCompleteWorkflowsV3 complete workflows.
// @Summary 批量完成工单
// @Description this api will directly change the work order status to finished without real online operation
// @Tags workflow
// @Id batchCompleteWorkflowsV3
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param data body v3.BatchCompleteWorkflowsReqV3 true "batch complete workflows request"
// @Success 200 {object} controller.BaseRes
// @router /v3/projects/{project_name}/workflows/complete [post]
func BatchCompleteWorkflowsV3(c echo.Context) error {
	req := new(BatchCompleteWorkflowsReqV3)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	up, err := dms.NewUserPermission(user.GetIDStr(), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	workflows := make([]*model.Workflow, len(req.WorkflowList))
	for i, completeWorkflow := range req.WorkflowList {
		executable, reason, err := sqlversion.CheckWorkflowExecutable(c.Request().Context(), projectUid, completeWorkflow.WorkflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !executable {
			return controller.JSONBaseErrorReq(c, fmt.Errorf(reason))
		}
		workflow, err := v2.CheckCanCompleteWorkflow(projectUid, completeWorkflow.WorkflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		// 执行上线的人可以决定真的上线这个工单还是直接标记完成
		lastStep := workflow.Record.Steps[len(workflow.Record.Steps)-1]
		canFinishWorkflow := up.CanOpProject()
		if !canFinishWorkflow {
			for _, assignee := range strings.Split(lastStep.Assignees, ",") {
				if assignee == user.GetIDStr() {
					canFinishWorkflow = true
					break
				}
			}
		}

		if !canFinishWorkflow {
			return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("the current user does not have permission to end these work orders")))
		}

		lastStep.State = model.WorkflowStepStateApprove
		lastStep.OperationUserId = user.GetIDStr()
		workflows[i] = workflow
		workflow.Record.Status = model.WorkflowStatusFinish
		workflow.Record.CurrentWorkflowStepId = 0

		needExecInstanceRecords := []*model.WorkflowInstanceRecord{}
		for _, inst := range workflow.Record.InstanceRecords {
			if inst.Task.Status == model.TaskStatusExecuteFailed || inst.Task.Status == model.TaskStatusAudited {
				inst.ExecutionUserId = user.GetIDStr()
				inst.IsSQLExecuted = true
				needExecInstanceRecords = append(needExecInstanceRecords, inst)
			}
		}
		if err := s.CompletionWorkflow(workflow, lastStep, needExecInstanceRecords); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if completeWorkflow.Desc != nil {
			workflowParam := make(map[string]interface{}, 1)
			carriageReturnNewline := "\r\n"
			workflowParam["desc"] = workflow.Desc + carriageReturnNewline + *completeWorkflow.Desc
			if err := s.UpdateWorkflowById(workflow.ID, workflowParam); err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
		}
	}

	return controller.JSONBaseErrorReq(c, nil)
}
