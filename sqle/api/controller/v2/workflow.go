package v2

import (
	"context"
	_err "errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/sqlversion"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

var errTaskHasBeenUsed = errors.New(errors.DataConflict, fmt.Errorf("task has been used in other workflow"))

type WorkflowStepResV2 struct {
	Id            uint       `json:"workflow_step_id,omitempty"`
	Number        uint       `json:"number"`
	Type          string     `json:"type" enums:"create_workflow,update_workflow,sql_review,sql_execute"`
	Desc          string     `json:"desc,omitempty"`
	Users         []string   `json:"assignee_user_name_list,omitempty"`
	OperationUser string     `json:"operation_user_name,omitempty"`
	OperationTime *time.Time `json:"operation_time,omitempty"`
	State         string     `json:"state,omitempty" enums:"initialized,approved,rejected"`
	Reason        string     `json:"reason,omitempty"`
}

// @Summary 审批通过
// @Description approve workflow
// @Tags workflow
// @Id approveWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param workflow_step_id path string true "workflow step id"
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/steps/{workflow_step_id}/approve [post]
func ApproveWorkflowV2(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowId := c.Param("workflow_id")

	s := model.GetStorage()

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, projectUid, workflow, []dmsV1.OpPermissionType{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	stepIdStr := c.Param("workflow_step_id")
	stepId, err := v1.FormatStringToInt(stepIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	nextStep := workflow.NextStep()

	err = server.CheckUserCanOperateStep(user, workflow, stepId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	if err := server.ApproveWorkflowProcess(workflow, user, s); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	go im.UpdateApprove(workflow.WorkflowId, user, model.ApproveStatusAgree, "")

	if nextStep != nil {
		go im.CreateApprove(string(workflow.ProjectId), workflow.WorkflowId)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type RejectWorkflowReqV2 struct {
	Reason string `json:"reason" form:"reason"`
}

// @Summary 审批驳回
// @Description reject workflow
// @Tags workflow
// @Id rejectWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Param workflow_step_id path string true "workflow step id"
// @param workflow_approve body v2.RejectWorkflowReqV2 true "workflow approve request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/steps/{workflow_step_id}/reject [post]
func RejectWorkflowV2(c echo.Context) error {
	req := new(RejectWorkflowReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowID := c.Param("workflow_id")
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// RejectWorkflow no need extra operation code for now.
	err = v1.CheckCurrentUserCanOperateWorkflow(c, projectUid, workflow, []dmsV1.OpPermissionType{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	stepIdStr := c.Param("workflow_step_id")
	stepId, err := v1.FormatStringToInt(stepIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = server.CheckUserCanOperateStep(user, workflow, stepId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	for _, inst := range workflow.Record.InstanceRecords {
		if inst.IsSQLExecuted {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("can not reject workflow, cause there is any task is executed"))
		}
		if inst.ScheduledAt != nil {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("can not reject workflow, cause there is any task is scheduled to be executed"))
		}
	}

	if err := server.RejectWorkflowProcess(workflow, req.Reason, user, s); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	go im.UpdateApprove(workflow.WorkflowId, user, model.ApproveStatusRefuse, req.Reason)

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// @Summary 审批关闭（中止）
// @Description cancel workflow
// @Tags workflow
// @Id cancelWorkflowV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param workflow_id path string true "workflow id"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/cancel [post]
func CancelWorkflowV2(c echo.Context) error {
	s := model.GetStorage()

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowID := c.Param("workflow_id")
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, workflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, projectUid, workflow, []dmsV1.OpPermissionType{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, err = checkCancelWorkflow(projectUid, workflow.WorkflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	up, err := dms.NewUserPermission(controller.GetUserID(c), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !(controller.GetUserID(c) == workflow.CreateUserId || up.CanOpProject()) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	workflow.Record.Status = model.WorkflowStatusCancel
	workflow.Record.CurrentWorkflowStepId = 0

	err = s.UpdateWorkflowStatus(workflow)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflowStageParam := make(map[string]interface{}, 1)
	workflowStageParam["workflow_release_status"] = model.WorkflowReleaseStatusNotNeedReleased
	err = s.UpdateStageWorkflowIfNeed(workflow.WorkflowId, workflowStageParam)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	go im.BatchCancelApprove([]string{workflow.WorkflowId}, user)

	return controller.JSONBaseErrorReq(c, nil)
}

func checkCancelWorkflow(projectId, workflowID string) (*model.Workflow, error) {
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectId, workflowID, model.GetStorage().GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return nil, err
	}

	if !(workflow.Record.Status == model.WorkflowStatusWaitForAudit ||
		workflow.Record.Status == model.WorkflowStatusWaitForExecution ||
		workflow.Record.Status == model.WorkflowStatusReject) {
		return nil, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status))
	}
	return workflow, nil
}

type BatchCancelWorkflowsReqV2 struct {
	WorkflowIDList []string `json:"workflow_id_list" form:"workflow_id_list"`
}

// BatchCancelWorkflowsV2 batch cancel workflows.
// @Summary 批量取消工单
// @Description batch cancel workflows
// @Tags workflow
// @Id batchCancelWorkflowsV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param BatchCancelWorkflowsReqV2 body v2.BatchCancelWorkflowsReqV2 true "batch cancel workflows request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/cancel [post]
func BatchCancelWorkflowsV2(c echo.Context) error {
	req := new(BatchCancelWorkflowsReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflows := make([]*model.Workflow, len(req.WorkflowIDList))
	workflowIds := make([]string, 0, len(req.WorkflowIDList))
	for i, workflowID := range req.WorkflowIDList {
		workflow, err := checkCancelWorkflow(projectUid, workflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		workflows[i] = workflow
		workflowIds = append(workflowIds, workflow.WorkflowId)
		workflow.Record.Status = model.WorkflowStatusCancel
		workflow.Record.CurrentWorkflowStepId = 0
	}

	if err := model.GetStorage().BatchUpdateWorkflowStatus(workflows); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	go im.BatchCancelApprove(workflowIds, user)

	return controller.JSONBaseErrorReq(c, nil)
}

type BatchCompleteWorkflowsReqV2 struct {
	WorkflowIDList []string `json:"workflow_id_list" form:"workflow_id_list"`
}

// BatchCompleteWorkflowsV2 complete workflows.
// @Deprecated
// @Summary 批量完成工单
// @Description this api will directly change the work order status to finished without real online operation
// @Tags workflow
// @Id batchCompleteWorkflowsV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param data body v2.BatchCompleteWorkflowsReqV2 true "batch complete workflows request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/complete [post]
func BatchCompleteWorkflowsV2(c echo.Context) error {
	req := new(BatchCancelWorkflowsReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
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

	workflows := make([]*model.Workflow, len(req.WorkflowIDList))
	for i, workflowID := range req.WorkflowIDList {
		workflow, err := CheckCanCompleteWorkflow(projectUid, workflowID)
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
			if !inst.IsSQLExecuted {
				inst.ExecutionUserId = user.GetIDStr()
				inst.IsSQLExecuted = true
				needExecInstanceRecords = append(needExecInstanceRecords, inst)
			}
		}
		if err := model.GetStorage().CompletionWorkflow(workflow, lastStep, needExecInstanceRecords); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	return controller.JSONBaseErrorReq(c, nil)
}

func CheckCanCompleteWorkflow(projectId, workflowID string) (*model.Workflow, error) {
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectId, workflowID, model.GetStorage().GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return nil, err
	}

	if workflow.Record.Status != model.WorkflowStatusWaitForExecution && workflow.Record.Status != model.WorkflowStatusExecFailed {
		return nil, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status))
	}
	return workflow, nil
}

// ExecuteOneTaskOnWorkflowV2
// @Summary 工单提交单个数据源上线
// @Description execute one task on workflow
// @Tags workflow
// @Id executeOneTaskOnWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Param task_id path string true "task id"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/execute [post]
func ExecuteOneTaskOnWorkflowV2(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflowID := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskIdStr := c.Param("task_id")
	taskId, err := v1.FormatStringToInt(taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = v1.PrepareForTaskExecution(c, projectUid, workflow, user, taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	executable, reason, err := sqlversion.CheckWorkflowExecutable(c.Request().Context(), projectUid, workflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !executable {
		return controller.JSONBaseErrorReq(c, fmt.Errorf(reason))
	}

	isCan, err := v1.IsTaskCanExecute(s, taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !isCan {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("task has no need to be executed. taskId=%v workflowId=%v", taskId, workflow.WorkflowId))
	}

	_, err = server.ExecuteWorkflow(workflow, map[uint]string{uint(taskId): user.GetIDStr()})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetWorkflowTasksResV2 struct {
	controller.BaseRes
	Data []*GetWorkflowTasksItemV2 `json:"data"`
}

type GetWorkflowTasksItemV2 struct {
	TaskId                   uint                       `json:"task_id"`
	InstanceName             string                     `json:"instance_name"`
	Status                   string                     `json:"status" enums:"wait_for_audit,wait_for_execution,exec_scheduled,exec_failed,exec_succeeded,executing,manually_executed,terminating,terminate_succeeded,terminate_failed"`
	ExecStartTime            *time.Time                 `json:"exec_start_time,omitempty"`
	ExecEndTime              *time.Time                 `json:"exec_end_time,omitempty"`
	ScheduleTime             *time.Time                 `json:"schedule_time,omitempty"`
	CurrentStepAssigneeUser  []string                   `json:"current_step_assignee_user_name_list,omitempty"`
	TaskPassRate             float64                    `json:"task_pass_rate"`
	TaskScore                int32                      `json:"task_score"`
	InstanceMaintenanceTimes []*v1.MaintenanceTimeResV1 `json:"instance_maintenance_times"`
	ExecutionUserName        string                     `json:"execution_user_name"`
}

// GetSummaryOfWorkflowTasksV2
// @Summary 获取工单数据源任务概览
// @Description get summary of workflow instance tasks
// @Tags workflow
// @Id getSummaryOfInstanceTasksV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Success 200 {object} v2.GetWorkflowTasksResV2
// @router /v2/projects/{project_name}/workflows/{workflow_id}/tasks [get]
func GetSummaryOfWorkflowTasksV2(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowId := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	err = v1.CheckCurrentUserCanViewWorkflow(c, projectUid, workflow, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeViewOthersWorkflow})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	queryData := map[string]interface{}{
		"workflow_id": workflowId,
		"project_id":  projectUid,
	}

	var taskDetails []*model.WorkflowTasksSummaryDetail
	workflowStatus := workflow.Record.Status
	// 当工单处于工作流程模板的审核阶段时，工单概览应该显示每个task的待操作人对应的审核模板步骤待操作人
	// 当工单处于工作流程模板的上线阶段时，工单概览应该分别显示每个task的待操作人，而不是审核模板步骤的待操作人
	if workflowStatus == model.WorkflowStatusExecuting || workflowStatus == model.WorkflowStatusWaitForExecution {
		taskDetails, err = s.GetWorkflowTaskSummaryByReq(queryData)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	} else {
		taskDetails, err = s.GetWorkflowStepSummaryByReqV2(queryData)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	for i, detail := range taskDetails {
		instance, exist, err := dms.GetInstanceInProjectById(c.Request().Context(), projectUid, detail.InstanceId)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if exist {
			taskDetails[i].InstanceName = instance.Name
			taskDetails[i].InstanceMaintenancePeriod = instance.MaintenancePeriod
		}
	}

	return c.JSON(http.StatusOK, &GetWorkflowTasksResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowToTasksSummaryRes(taskDetails),
	})
}

func convertWorkflowToTasksSummaryRes(taskDetails []*model.WorkflowTasksSummaryDetail) []*GetWorkflowTasksItemV2 {
	res := make([]*GetWorkflowTasksItemV2, len(taskDetails))

	for i, taskDetail := range taskDetails {
		userNames := make([]string, 0)
		for _, userId := range strings.Split(taskDetail.CurrentStepAssigneeUserIds.String, ",") {
			if userId == "" {
				continue
			}
			userNames = append(userNames, dms.GetUserNameWithDelTag(userId))
		}

		res[i] = &GetWorkflowTasksItemV2{
			TaskId:                   taskDetail.TaskId,
			InstanceName:             utils.AddDelTag(taskDetail.InstanceDeletedAt, taskDetail.InstanceName),
			Status:                   v1.GetTaskStatusRes(taskDetail.WorkflowRecordStatus, taskDetail.TaskStatus, taskDetail.InstanceScheduledAt),
			ExecStartTime:            taskDetail.TaskExecStartAt,
			ExecEndTime:              taskDetail.TaskExecEndAt,
			ScheduleTime:             taskDetail.InstanceScheduledAt,
			CurrentStepAssigneeUser:  userNames,
			TaskPassRate:             taskDetail.TaskPassRate,
			TaskScore:                taskDetail.TaskScore,
			InstanceMaintenanceTimes: v1.ConvertPeriodToMaintenanceTimeResV1(taskDetail.InstanceMaintenancePeriod),
			ExecutionUserName:        dms.GetUserNameWithDelTag(taskDetail.ExecutionUserId),
		}
	}

	return res
}

type CreateWorkflowReqV2 struct {
	Subject              string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc                 string `json:"desc" form:"desc"`
	SqlVersionID         *uint  `json:"sql_version_id" form:"sql_version_id"`
	TaskIds              []uint `json:"task_ids" form:"task_ids" valid:"required"`
}

type CreateWorkflowResV2 struct {
	controller.BaseRes
	Data *CreateWorkflowResV2Data `json:"data"`
}

type CreateWorkflowResV2Data struct {
	WorkflowID string `json:"workflow_id"`
}

// CreateWorkflowV2
// @Summary 创建工单
// @Description create workflow
// @Accept json
// @Produce json
// @Tags workflow
// @Id createWorkflowV2
// @Security ApiKeyAuth
// @Param instance body v2.CreateWorkflowReqV2 true "create workflow request"
// @Param project_name path string true "project name"
// @Success 200 {object} CreateWorkflowResV2
// @router /v2/projects/{project_name}/workflows [post]
func CreateWorkflowV2(c echo.Context) error {
	req := new(CreateWorkflowReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// dms-todo: 与 dms 生成uid保持一致
	workflowId, err := utils.GenUid()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	_, exist, err := s.GetWorkflowByProjectAndWorkflowId(projectUid, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("workflow[%v] is exist", workflowId)))
	}

	taskIds := utils.RemoveDuplicateUint(req.TaskIds)
	if len(taskIds) > v1.MaximumDataSourceNum {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, fmt.Errorf("the max task count of a workflow is %v", v1.MaximumDataSourceNum)))
	}
	tasks, foundAllTasks, err := s.GetTasksByIds(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !foundAllTasks {
		return controller.JSONBaseErrorReq(c, errors.NewTaskNoExistOrNoAccessErr())
	}

	instanceIdsOfWorkflowTasks := make([]uint64, 0, len(tasks))
	for _, task := range tasks {
		instanceIdsOfWorkflowTasks = append(instanceIdsOfWorkflowTasks, task.InstanceId)
	}

	instancesOfWorkflowInProject, err := dms.GetInstancesInProjectByIds(c.Request().Context(), projectUid, instanceIdsOfWorkflowTasks)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectInstanceMap := map[uint64]*model.Instance{}
	for _, instance := range instancesOfWorkflowInProject {
		projectInstanceMap[instance.ID] = instance
	}

	workflowTemplate, exist, err := s.GetWorkflowTemplateByProjectId(model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("the task instance is not bound workflow template")))
	}

	for _, task := range tasks {
		if instance, ok := projectInstanceMap[task.InstanceId]; ok {
			task.Instance = instance
		}

		if task.Instance == nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist. taskId=%v", task.ID)))
		}

		if task.Instance.ProjectId != projectUid {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not in project. taskId=%v", task.ID)))
		}

		count, err := s.GetTaskSQLCountByTaskID(task.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if count == 0 {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("workflow's execute sql is null. taskId=%v", task.ID)))
		}

		if task.CreateUserId != uint64(user.ID) {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
				fmt.Errorf("the task is not created by yourself. taskId=%v", task.ID)))
		}

		if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
			return controller.JSONBaseErrorReq(c, v1.ErrForbidMyBatisXMLTask(task.ID))
		}
	}

	// check user role operations
	{

		canOperationInstance, err := v1.CheckCurrentUserCanCreateWorkflow(c.Request().Context(), projectUid, user, tasks)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !canOperationInstance {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("can't operation instance"))
		}

	}

	count, err := s.GetWorkflowRecordCountByTaskIds(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if count > 0 {
		return controller.JSONBaseErrorReq(c, errTaskHasBeenUsed)
	}

	stepTemplates, err := s.GetWorkflowStepsByTemplateId(workflowTemplate.ID)
	if err != nil {
		return err
	}

	memberWithPermissions, _, err := dmsobject.ListMembersInProject(c.Request().Context(), controller.GetDMSServerAddress(), dmsV1.ListMembersForInternalReq{
		ProjectUid: projectUid,
		PageSize:   999,
		PageIndex:  1,
	})
	if err != nil {
		return err
	}

	if req.SqlVersionID != nil {
		err = sqlversion.CheckInstanceInWorkflowCanAssociateToTheFirstStageOfVersion(*req.SqlVersionID, instanceIdsOfWorkflowTasks)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	err = s.CreateWorkflowV2(req.Subject, workflowId, req.Desc, user, tasks, stepTemplates, model.ProjectUID(projectUid), req.SqlVersionID, nil, nil, func(tasks []*model.Task) (auditWorkflowUsers, canExecUser [][]*model.User) {
		auditWorkflowUsers = make([][]*model.User, len(tasks))
		executorWorkflowUsers := make([][]*model.User, len(tasks))
		for i, task := range tasks {
			auditWorkflowUsers[i], err = v1.GetCanOpInstanceUsers(memberWithPermissions, task.Instance, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeAuditWorkflow})
			if err != nil {
				return
			}
			executorWorkflowUsers[i], err = v1.GetCanOpInstanceUsers(memberWithPermissions, task.Instance, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeExecuteWorkflow})
			if err != nil {
				return
			}
		}
		return auditWorkflowUsers, executorWorkflowUsers
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, exist, err := s.GetLastWorkflow()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("should exist at least one workflow after create workflow")))
	}

	go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeCreate)

	go im.CreateApprove(string(workflow.ProjectId), workflow.WorkflowId)

	return c.JSON(http.StatusOK, &CreateWorkflowResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data: &CreateWorkflowResV2Data{
			WorkflowID: workflow.WorkflowId,
		},
	})
}

type UpdateWorkflowReqV2 struct {
	TaskIds []uint `json:"task_ids" form:"task_ids" valid:"required"`
}

// UpdateWorkflowV2
// @Summary 更新工单（工单被驳回、工单被关闭、执行成功、执行失败后才可更新）
// @Description update workflow when it is rejected to creator.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Param instance body v2.UpdateWorkflowReqV2 true "update workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/ [patch]
func UpdateWorkflowV2(c echo.Context) error {
	req := new(UpdateWorkflowReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowId := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, projectUid, workflow, []dmsV1.OpPermissionType{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	tasks, _, err := s.GetTasksByIds(req.TaskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if len(tasks) <= 0 {
		return controller.JSONBaseErrorReq(c, errors.NewTaskNoExistOrNoAccessErr())
	}

	instanceIds := make([]uint64, 0, len(tasks))
	for _, task := range tasks {
		instanceIds = append(instanceIds, task.InstanceId)
	}

	instances, err := dms.GetInstancesInProjectByIds(c.Request().Context(), projectUid, instanceIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceMap := map[uint64]*model.Instance{}
	for _, instance := range instances {
		instanceMap[instance.ID] = instance
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskIds := make([]uint, len(tasks))
	for i, task := range tasks {
		taskIds[i] = task.ID

		if instance, ok := instanceMap[task.InstanceId]; ok {
			task.Instance = instance
		}

		count, err := s.GetTaskSQLCountByTaskID(task.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if count == 0 {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("task's execute sql is null. taskId=%v", task.ID)))
		}

		err = v1.CheckCurrentUserCanOpTask(c, task)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		if task.Instance == nil {
			return controller.JSONBaseErrorReq(c, v1.ErrInstanceNotExist)
		}

		if uint64(user.ID) != task.CreateUserId {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
				fmt.Errorf("the task is not created by yourself. taskId=%v", task.ID)))
		}

		if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
			return controller.JSONBaseErrorReq(c, v1.ErrForbidMyBatisXMLTask(task.ID))
		}
	}

	count, err := s.GetWorkflowRecordCountByTaskIds(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if count > 0 {
		return controller.JSONBaseErrorReq(c, errTaskHasBeenUsed)
	}
	// When workflow status is rejected or exec failed, the user can recommit workflow. And the workflow becomes waiting for the audit process.
	// TODO 重试工单 驳回、成功、失败、关闭的工单可以重试，应抽离出一个函数
	if workflow.Record.Status != model.WorkflowStatusReject && workflow.Record.Status != model.WorkflowStatusExecFailed {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)))
	}

	if user.GetIDStr() != workflow.CreateUserId {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	template, exist, err := s.GetWorkflowTemplateByProjectId(workflow.ProjectId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
			fmt.Errorf("failed to find the corresponding workflow template based on the task id")))
	}

	err = v1.CheckWorkflowCanCommit(template, tasks)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateWorkflowRecord(workflow, tasks)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeCreate)

	go im.CreateApprove(string(workflow.ProjectId), workflow.WorkflowId)

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateWorkflowScheduleReqV2 struct {
	ScheduleTime *time.Time `json:"schedule_time"`
	IsNotify     *bool      `json:"is_notify"`
	NotifyType   *string    `json:"notify_type" enums:"wechat,feishu"`
}

// UpdateWorkflowScheduleV2
// @Summary 设置工单数据源定时上线时间（设置为空则代表取消定时时间，需要SQL审核流程都通过后才可以设置）
// @Description update workflow schedule.
// @Tags workflow
// @Accept json
// @Produce json
// @Id updateWorkflowScheduleV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param task_id path string true "task id"
// @Param project_name path string true "project name"
// @Param instance body v2.UpdateWorkflowScheduleReqV2 true "update workflow schedule request"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/schedule [put]
func UpdateWorkflowScheduleV2(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowId := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, projectUid, workflow, []dmsV1.OpPermissionType{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskId := c.Param("task_id")
	taskIdUint, err := v1.FormatStringToUint64(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	req := new(UpdateWorkflowScheduleReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, _err.New("workflow current step not found")))
	}

	if workflow.Record.Status != model.WorkflowStatusWaitForExecution {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow need to be approved first")))
	}

	err = server.CheckUserCanOperateStep(user, workflow, int(currentStep.ID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}
	var curTaskRecord *model.WorkflowInstanceRecord
	for _, ir := range workflow.Record.InstanceRecords {
		if uint64(ir.TaskId) == taskIdUint {
			curTaskRecord = ir
		}
	}
	if curTaskRecord == nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, _err.New("task is not found in workflow")))
	}

	if req.ScheduleTime != nil && req.ScheduleTime.Before(time.Now()) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf(
			"request schedule time is too early")))
	}

	if curTaskRecord.IsSQLExecuted {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf(
			"task has been executed")))
	}

	instance, exist, err := dms.GetInstanceInProjectById(c.Request().Context(), projectUid, curTaskRecord.InstanceId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrInstanceNotExist)
	}

	if req.ScheduleTime != nil && len(instance.MaintenancePeriod) != 0 && !instance.MaintenancePeriod.IsWithinScope(*req.ScheduleTime) {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowExecuteTimeIncorrect)
	}

	executable, reason, err := sqlversion.CheckWorkflowExecutable(c.Request().Context(), projectUid, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !executable {
		return controller.JSONBaseErrorReq(c, fmt.Errorf(reason))
	}

	err = s.UpdateInstanceRecordSchedule(curTaskRecord, user.GetIDStr(), req.ScheduleTime)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.IsNotify != nil && *req.IsNotify && req.NotifyType != nil && req.ScheduleTime != nil {
		if err := s.CreateNotifyRecord(*req.NotifyType, curTaskRecord); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.ScheduleTime == nil {
		if err := s.CancelNotify(uint(taskIdUint)); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// ExecuteTasksOnWorkflowV2
// @Summary 多数据源批量上线
// @Description execute tasks on workflow
// @Tags workflow
// @Id executeTasksOnWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/workflows/{workflow_id}/tasks/execute [post]
func ExecuteTasksOnWorkflowV2(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflowId := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	executable, reason, err := sqlversion.CheckWorkflowExecutable(c.Request().Context(), projectUid, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !executable {
		return controller.JSONBaseErrorReq(c, fmt.Errorf(reason))
	}

	if err := v1.PrepareForWorkflowExecution(c, projectUid, workflow, user); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	_, err = server.ExecuteTasksProcess(workflow.WorkflowId, projectUid, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	im.UpdateApprove(workflow.WorkflowId, user, model.ApproveStatusAgree, "")

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetWorkflowResV2 struct {
	controller.BaseRes
	Data *WorkflowResV2 `json:"data"`
}

type WorkflowTaskItem struct {
	Id uint `json:"task_id"`
}

type WorkflowRecordResV2 struct {
	Tasks             []*WorkflowTaskItem  `json:"tasks"`
	CurrentStepNumber uint                 `json:"current_step_number,omitempty"`
	Status            string               `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,executing,finished"`
	Executable        bool                 `json:"executable"`
	ExecutableReason  string               `json:"executable_reason"`
	Steps             []*WorkflowStepResV2 `json:"workflow_step_list,omitempty"`
}

type WorkflowResV2 struct {
	Name                        string                        `json:"workflow_name"`
	WorkflowID                  string                        `json:"workflow_id"`
	Desc                        string                        `json:"desc,omitempty"`
	Mode                        string                        `json:"mode" enums:"same_sqls,different_sqls"`
	ExecMode                    string                        `json:"exec_mode" enums:"sql_file,sqls"`
	CreateUser                  string                        `json:"create_user_name"`
	CreateTime                  *time.Time                    `json:"create_time"`
	SqlVersion                  *SqlVersion                   `json:"sql_version,omitempty"`
	Record                      *WorkflowRecordResV2          `json:"record"`
	RecordHistory               []*WorkflowRecordResV2        `json:"record_history_list,omitempty"`
	AssociatedStageWorkflows    []*AssociatedStageWorkflows   `json:"associated_stage_workflows,omitempty"`
	AssociatedRollbackWorkflows []*AssociatedRollbackWorkflow `json:"associated_rollback_workflows"`
}

type AssociatedStageWorkflows struct {
	WorkflowID        string `json:"workflow_id"`
	WorkflowName      string `json:"workflow_name"`
	Status            string `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,executing,finished"`
	SqlVersionStageID uint   `json:"sql_version_stage_id"`
	StageSequence     int    `json:"stage_sequence"`
}

type SqlVersion struct {
	SqlVersionName string `json:"sql_version_name"`
	SqlVersionId   uint   `json:"sql_version_id"`
}

// GetWorkflowV2
// @Summary 获取工单详情
// @Description get workflow detail
// @Tags workflow
// @Id getWorkflowV2
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Success 200 {object} GetWorkflowResV2
// @router /v2/projects/{project_name}/workflows/{workflow_id}/ [get]
func GetWorkflowV2(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflowID := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = v1.CheckCurrentUserCanViewWorkflow(c, projectUid, workflow, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeViewOthersWorkflow})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// TODO 优化为一次批量用户查询,history 记录也许一并处理
	for i := range workflow.Record.Steps {
		step := workflow.Record.Steps[i]
		AssigneesUserNames := make([]string, 0)
		for _, id := range strings.Split(step.Assignees, ",") {
			if id == "" {
				continue
			}
			AssigneesUserNames = append(AssigneesUserNames, dms.GetUserNameWithDelTag(id))
		}
		step.Assignees = strings.Join(AssigneesUserNames, ",")
		if workflow.CurrentStep() != nil && step.ID == workflow.CurrentStep().ID {
			workflow.Record.CurrentStep = step
		}
		workflow.Record.Steps[i] = step
	}

	history, err := s.GetWorkflowHistoryById(workflow.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflow.RecordHistory = history
	associatedWorkflows, err := s.GetAssociatedStageWorkflows(workflow.WorkflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	sqlVersion, err := s.GetSQLVersionByWorkflowId(workflow.WorkflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetWorkflowResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowToRes(c.Request().Context(), workflow, sqlVersion, associatedWorkflows),
	})
}

func convertWorkflowToRes(ctx context.Context, workflow *model.Workflow, sqlVersion *model.SqlVersion, associatedWorkflows []*model.AssociatedStageWorkflow) *WorkflowResV2 {
	workflowRes := &WorkflowResV2{
		Name:                     workflow.Subject,
		WorkflowID:               workflow.WorkflowId,
		Desc:                     workflow.Desc,
		Mode:                     workflow.Mode,
		ExecMode:                 workflow.ExecMode,
		CreateUser:               dms.GetUserNameWithDelTag(workflow.CreateUserId),
		CreateTime:               &workflow.CreatedAt,
		AssociatedStageWorkflows: convertAssociatedWorkflowToRes(associatedWorkflows),
	}
	sqlVersionRes := &SqlVersion{
		SqlVersionId:   sqlVersion.ID,
		SqlVersionName: sqlVersion.Version,
	}
	workflowRes.SqlVersion = sqlVersionRes
	// convert workflow record
	workflowRecordRes := convertWorkflowRecordToRes(ctx, workflow, workflow.Record)

	// convert workflow record history
	recordHistory := make([]*WorkflowRecordResV2, 0, len(workflow.RecordHistory))
	for _, record := range workflow.RecordHistory {
		recordRes := convertWorkflowRecordToRes(ctx, workflow, record)
		recordHistory = append(recordHistory, recordRes)
	}
	workflowRes.RecordHistory = recordHistory
	workflowRes.Record = workflowRecordRes

	return workflowRes
}

func convertAssociatedWorkflowToRes(associatedWorkflows []*model.AssociatedStageWorkflow) []*AssociatedStageWorkflows {
	associatedWorkflowsRes := make([]*AssociatedStageWorkflows, 0, len(associatedWorkflows))
	for _, associatedWorkflow := range associatedWorkflows {
		associatedWorkflowsRes = append(associatedWorkflowsRes, &AssociatedStageWorkflows{
			WorkflowID:        associatedWorkflow.WorkflowID.String,
			WorkflowName:      associatedWorkflow.WorkflowName.String,
			Status:            associatedWorkflow.Status.String,
			SqlVersionStageID: associatedWorkflow.SqlVersionStageID,
			StageSequence:     associatedWorkflow.StageSequence,
		})
	}

	return associatedWorkflowsRes
}

func convertWorkflowRecordToRes(ctx context.Context, workflow *model.Workflow, record *model.WorkflowRecord) *WorkflowRecordResV2 {
	steps := make([]*WorkflowStepResV2, 0, len(record.Steps)+1)
	// It is filled by create user and create time;
	// and tell others that this is a creating or updating operation.
	var stepType string
	if workflow.IsFirstRecord(record) {
		stepType = model.WorkflowStepTypeCreateWorkflow
	} else {
		stepType = model.WorkflowStepTypeUpdateWorkflow
	}

	firstVirtualStep := &WorkflowStepResV2{
		Type:          stepType,
		OperationTime: &record.CreatedAt,
		OperationUser: dms.GetUserNameWithDelTag(workflow.CreateUserId),
	}
	steps = append(steps, firstVirtualStep)

	// convert workflow actual step
	for _, step := range record.Steps {
		stepRes := convertWorkflowStepToRes(step)
		steps = append(steps, stepRes)
	}
	// fill step number
	var currentStepNum uint
	for i, step := range steps {
		number := uint(i + 1)
		step.Number = number
		if step.Id != 0 && step.Id == record.CurrentWorkflowStepId {
			currentStepNum = number
		}
	}

	tasksRes := make([]*WorkflowTaskItem, len(record.InstanceRecords))
	for i, inst := range record.InstanceRecords {
		tasksRes[i] = &WorkflowTaskItem{Id: inst.TaskId}
	}

	var err error
	var executable bool
	var reason string = fmt.Sprintf("the status of workflow is %v", record.Status)
	if record.Status == model.WorkflowStatusWaitForExecution {
		executable, reason, err = sqlversion.CheckWorkflowExecutable(ctx, string(workflow.ProjectId), workflow.WorkflowId)
		if err != nil {
			reason = err.Error()
		}
	}

	return &WorkflowRecordResV2{
		Tasks:             tasksRes,
		CurrentStepNumber: currentStepNum,
		Status:            record.Status,
		Steps:             steps,
		Executable:        executable,
		ExecutableReason:  reason,
	}
}

func convertWorkflowStepToRes(step *model.WorkflowStep) *WorkflowStepResV2 {
	stepRes := &WorkflowStepResV2{
		Id:            step.ID,
		Type:          step.Template.Typ,
		Desc:          step.Template.Desc,
		OperationTime: step.OperateAt,
		State:         step.State,
		Reason:        step.Reason,
		Users:         []string{},
		OperationUser: dms.GetUserNameWithDelTag(step.OperationUserId),
	}
	stepRes.Users = append(stepRes.Users, strings.Split(step.Assignees, ",")...)
	return stepRes
}

type GetAuditFileListReq struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditFileListRes struct {
	controller.BaseRes
	Data      []AuditFileStatistic `json:"data"`
	TotalNums uint64               `json:"total_nums"`
}

type AuditFileStatistic struct {
	FileID           string            `json:"file_id"`
	FileName         string            `json:"file_name"`
	ExecOrder        uint              `json:"exec_order"`
	ExecStatus       string            `json:"exec_status"`
	AuditResultCount *AuditResultCount `json:"audit_result_count"`
}

type AuditResultCount struct {
	ErrorSQLCount   uint `json:"error_sql_count"`
	WarningSQLCount uint `json:"warning_sql_count"`
	NormalSQLCount  uint `json:"normal_sql_count"`
	NoticeSQLCount  uint `json:"notice_sql_count"`
}

// GetAuditFileList
// @Summary 获取审核任务文件概览列表
// @Description get audit task file list
// @Tags task
// @Id getAuditFileList
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} GetAuditFileListRes
// @router /v2/tasks/audits/{task_id}/files [get]
func GetAuditFileList(c echo.Context) error {
	req := new(GetAuditFileListReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	taskId := c.Param("task_id")

	task, err := v1.GetTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = v1.CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	data := map[string]interface{}{
		"task_id": taskId,
		"limit":   limit,
		"offset":  offset,
	}
	result, count, err := s.GetAuditStatisticByTaskId(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditFileListRes{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      convertToAuditFileList(result),
		TotalNums: count,
	})
}

func convertToAuditFileList(input []*model.AuditResultStatistic) (output []AuditFileStatistic) {
	output = make([]AuditFileStatistic, 0, len(input))
	for _, file := range input {
		output = append(output, AuditFileStatistic{
			FileID:     file.ExecFileID,
			FileName:   file.ExecFileName,
			ExecOrder:  file.ExecOrder,
			ExecStatus: file.FileExecStatus(),
			AuditResultCount: &AuditResultCount{
				ErrorSQLCount:   file.ErrorCount,
				WarningSQLCount: file.WarningCount,
				NormalSQLCount:  file.NormalCount,
				NoticeSQLCount:  file.NoticeCount,
			},
		})
	}
	return
}

type GetAuditFileExecStatisticRes struct {
	controller.BaseRes
	Data *AuditFileExecStatistic `json:"data"`
}

type AuditFileExecStatistic struct {
	FileID          string           `json:"file_id"`
	FileName        string           `json:"file_name"`
	ExecResultCount *ExecResultCount `json:"exec_result_count"`
}

type ExecResultCount struct {
	FailedCount             uint `json:"failed_count"`
	SucceededCount          uint `json:"succeeded_count"`
	InitializedCount        uint `json:"initialized_count"`
	DoingCount              uint `json:"doing_count"`
	ManuallyExecutedCount   uint `json:"manually_executed_count"`
	TerminateSucceededCount uint `json:"terminate_succeeded_count"`
	TerminateFailedCount    uint `json:"terminate_failed_count"`
}

// GetAuditFileExecStatistic
// @Summary 获取审核任务文件执行概览
// @Description get audit task file execute statistic
// @Tags task
// @Id getAuditFileExecStatistic
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param file_id path string true "file id"
// @Success 200 {object} GetAuditFileExecStatisticRes
// @router /v2/tasks/audits/{task_id}/files/{file_id}/ [get]
func GetAuditFileExecStatistic(c echo.Context) error {
	taskId := c.Param("task_id")
	fileId := c.Param("file_id")

	task, err := v1.GetTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = v1.CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	fileIdNum, err := strconv.Atoi(fileId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("the file id should be a number %v", err))
	}
	if fileIdNum <= 0 {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("file id cannot be less than or equal with zero"))
	}

	s := model.GetStorage()
	data := map[string]interface{}{
		"task_id": taskId,
		"file_id": fileId,
	}
	result, err := s.GetAuditFileExecStatisticByFileId(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditFileExecStatisticRes{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertToAuditFileExecStatistic(result),
	})
}

func convertToAuditFileExecStatistic(file *model.AuditFileExecStatistic) *AuditFileExecStatistic {
	return &AuditFileExecStatistic{
		FileID:   file.ExecFileID,
		FileName: file.ExecFileName,
		ExecResultCount: &ExecResultCount{
			InitializedCount:        file.InitializedCount,
			SucceededCount:          file.SucceededCount,
			FailedCount:             file.FailedCount,
			DoingCount:              file.DoingCount,
			ManuallyExecutedCount:   file.ManuallyExecutedCount,
			TerminateSucceededCount: file.TerminateSucceededCount,
			TerminateFailedCount:    file.TerminateFailedCount,
		},
	}
}
