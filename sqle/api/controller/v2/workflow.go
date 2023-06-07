package v2

import (
	_err "errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/server"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/model"

	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"

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
	projectName := c.Param("project_name")
	workflowId := c.Param("workflow_id")

	s := model.GetStorage()
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrProjectNotExist(projectName))
	}

	userName := controller.GetUserName(c)
	if err := v1.CheckIsProjectMember(userName, project.Name); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, project, workflow, []uint{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	stepIdStr := c.Param("workflow_step_id")
	stepId, err := v1.FormatStringToInt(stepIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowIdStr := strconv.Itoa(int(workflow.ID))
	workflow, exist, err = s.GetWorkflowDetailById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	nextStep := workflow.NextStep()

	err = v1.CheckUserCanOperateStep(user, workflow, stepId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	if err := server.ApproveWorkflowProcess(workflow, user, s); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	go im.UpdateApprove(workflow.ID, user.Phone, model.ApproveStatusAgree, "")

	if nextStep.Template.Typ != model.WorkflowStepTypeSQLExecute {
		go im.CreateApprove(strconv.Itoa(int(workflow.ID)))
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

	projectName := c.Param("project_name")
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrProjectNotExist(projectName))
	}

	workflowID := c.Param("workflow_id")
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(project.Name, workflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	// RejectWorkflow no need extra operation code for now.
	err = v1.CheckCurrentUserCanOperateWorkflow(c, project, workflow, []uint{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	stepIdStr := c.Param("workflow_step_id")
	stepId, err := v1.FormatStringToInt(stepIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowIdStr := strconv.Itoa(int(workflow.ID))
	workflow, exist, err = s.GetWorkflowDetailById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	err = v1.CheckUserCanOperateStep(user, workflow, stepId)
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

	go im.UpdateApprove(workflow.ID, user.Phone, model.ApproveStatusRefuse, req.Reason)

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

	projectName := c.Param("project_name")
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrProjectNotExist(projectName))
	}

	workflowID := c.Param("workflow_id")
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(project.Name, workflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, project, workflow, []uint{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, err = checkCancelWorkflow(project.Name, workflow.WorkflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowStatus := workflow.Record.Status

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	isManager, err := s.IsProjectManager(user.Name, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !(user.ID == workflow.CreateUserId || user.Name == model.DefaultAdminUser || isManager) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	workflow.Record.Status = model.WorkflowStatusCancel
	workflow.Record.CurrentWorkflowStepId = 0

	err = s.UpdateWorkflowStatus(workflow)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if workflowStatus == model.WorkflowStatusWaitForAudit {
		go im.CancelApprove(workflow.ID)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

func checkCancelWorkflow(projectName, workflowID string) (*model.Workflow, error) {
	workflow, exist, err := model.GetStorage().GetWorkflowDetailByWorkflowID(projectName, workflowID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, v1.ErrWorkflowNoAccess
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

	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	if err := v1.CheckIsProjectManager(userName, projectName); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflows := make([]*model.Workflow, len(req.WorkflowIDList))
	for i, workflowID := range req.WorkflowIDList {
		workflow, err := checkCancelWorkflow(projectName, workflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		workflows[i] = workflow

		workflow.Record.Status = model.WorkflowStatusCancel
		workflow.Record.CurrentWorkflowStepId = 0
	}

	if err := model.GetStorage().BatchUpdateWorkflowStatus(workflows); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type BatchCompleteWorkflowsReqV2 struct {
	WorkflowIDList []string `json:"workflow_id_list" form:"workflow_id_list"`
}

// BatchCompleteWorkflowsV2 complete workflows.
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

	projectName := c.Param("project_name")
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	isManager, err := s.IsProjectManager(user.Name, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflows := make([]*model.Workflow, len(req.WorkflowIDList))
	for i, workflowID := range req.WorkflowIDList {
		workflow, err := checkCanCompleteWorkflow(projectName, workflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		// 执行上线的人可以决定真的上线这个工单还是直接标记完成
		lastStep := workflow.Record.Steps[len(workflow.Record.Steps)-1]
		canFinishWorkflow := isManager
		if !canFinishWorkflow {
			for _, assignee := range lastStep.Assignees {
				if assignee.Name == user.Name {
					canFinishWorkflow = true
					break
				}
			}
		}

		if !canFinishWorkflow {
			return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("the current user does not have permission to end these work orders")))
		}

		lastStep.State = model.WorkflowStepStateApprove
		lastStep.OperationUserId = user.ID
		workflows[i] = workflow
		workflow.Record.Status = model.WorkflowStatusFinish
		workflow.Record.CurrentWorkflowStepId = 0

		needExecInstanceRecords := []*model.WorkflowInstanceRecord{}
		for _, inst := range workflow.Record.InstanceRecords {
			if !inst.IsSQLExecuted {
				inst.ExecutionUserId = user.ID
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

func checkCanCompleteWorkflow(projectName, workflowID string) (*model.Workflow, error) {
	workflow, exist, err := model.GetStorage().GetWorkflowDetailByWorkflowID(projectName, workflowID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, v1.ErrWorkflowNoAccess
	}
	if !(workflow.Record.Status == model.WorkflowStatusWaitForExecution) {
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
	projectName := c.Param("project_name")
	workflowID := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, workflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	workflowId := fmt.Sprintf("%v", workflow.ID)

	taskIdStr := c.Param("task_id")
	taskId, err := v1.FormatStringToInt(taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, exist, err = s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = v1.PrepareForWorkflowExecution(c, projectName, workflow, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	isCan, err := v1.IsTaskCanExecute(s, taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !isCan {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("task has no need to be executed. taskId=%v workflowId=%v", taskId, workflowId))
	}

	err = server.ExecuteWorkflow(workflow, map[uint]uint{uint(taskId): user.ID})
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
	projectName := c.Param("project_name")
	workflowId := c.Param("workflow_id")

	if err := CheckCurrentUserCanViewWorkflow(c, workflowId, projectName); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}

	// TODO: Code logic optimization
	instances, err := s.GetUserCanOpInstancesFromProject(user, projectName, []uint{model.OP_WORKFLOW_EXECUTE})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	instanceMap := make(map[string]string)
	for i := range instances {
		inst := instances[i]
		instanceMap[inst.Name] = user.Name
	}

	queryData := map[string]interface{}{
		"workflow_id":  workflowId,
		"project_name": projectName,
	}

	taskDetails, err := s.GetWorkflowTasksSummaryByReqV2(queryData)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetWorkflowTasksResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowToTasksSummaryRes(taskDetails, instanceMap),
	})
}

func convertWorkflowToTasksSummaryRes(taskDetails []*model.WorkflowTasksSummaryDetail, instanceMap map[string]string) []*GetWorkflowTasksItemV2 {
	res := make([]*GetWorkflowTasksItemV2, len(taskDetails))

	for i, taskDetail := range taskDetails {

		res[i] = &GetWorkflowTasksItemV2{
			TaskId:                   taskDetail.TaskId,
			InstanceName:             utils.AddDelTag(taskDetail.InstanceDeletedAt, taskDetail.InstanceName),
			Status:                   v1.GetTaskStatusRes(taskDetail.WorkflowRecordStatus, taskDetail.TaskStatus, taskDetail.InstanceScheduledAt),
			ExecStartTime:            taskDetail.TaskExecStartAt,
			ExecEndTime:              taskDetail.TaskExecEndAt,
			ScheduleTime:             taskDetail.InstanceScheduledAt,
			CurrentStepAssigneeUser:  taskDetail.CurrentStepAssigneeUsers,
			TaskPassRate:             taskDetail.TaskPassRate,
			TaskScore:                taskDetail.TaskScore,
			InstanceMaintenanceTimes: v1.ConvertPeriodToMaintenanceTimeResV1(taskDetail.InstanceMaintenancePeriod),
			ExecutionUserName:        utils.AddDelTag(taskDetail.ExecutionUserDeletedAt, taskDetail.ExecutionUserName),
		}

		// NOTE: 当 SQL 处于上线中时，CurrentStepAssigneeUser 可能为空。此处需要「拥有上线权限的用户」
		if taskDetail.TaskStatus == model.TaskStatusExecuting {
			res[i].CurrentStepAssigneeUser = []string{instanceMap[taskDetail.InstanceName]}
		}
	}
	return res
}

type CreateWorkflowReqV2 struct {
	Subject string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc    string `json:"desc" form:"desc"`
	TaskIds []uint `json:"task_ids" form:"task_ids" valid:"required"`
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

	projectName := c.Param("project_name")

	s := model.GetStorage()
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrProjectNotExist(projectName))
	}
	if project.IsArchived() {
		return controller.JSONBaseErrorReq(c, v1.ErrProjectArchived)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := v1.CheckIsProjectMember(user.Name, project.Name); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowId, err := utils.GenUid()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	_, exist, err = s.GetWorkflowByProjectNameAndWorkflowId(project.Name, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("workflow is exist")))
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

	insIdtMap := make(map[uint] /* project instance id */ struct{}, len(project.Instances))
	for _, instance := range project.Instances {
		insIdtMap[instance.ID] = struct{}{}
	}

	workflowTemplateId := tasks[0].Instance.WorkflowTemplateId
	for _, task := range tasks {
		if task.Instance == nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist. taskId=%v", task.ID)))
		}

		if _, ok := insIdtMap[task.InstanceId]; !ok {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not in project. taskId=%v", task.ID)))
		}

		count, err := s.GetTaskSQLCountByTaskID(task.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if count == 0 {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("workflow's execute sql is null. taskId=%v", task.ID)))
		}

		if task.CreateUserId != user.ID {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
				fmt.Errorf("the task is not created by yourself. taskId=%v", task.ID)))
		}

		if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
			return controller.JSONBaseErrorReq(c, v1.ErrForbidMyBatisXMLTask(task.ID))
		}

		// all instances must use the same workflow template
		if task.Instance.WorkflowTemplateId != workflowTemplateId {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict,
				fmt.Errorf("all instances must use the same workflow template")))
		}
	}

	// check user role operations
	{
		err = v1.CheckCurrentUserCanCreateWorkflow(user, tasks, projectName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	count, err := s.GetWorkflowRecordCountByTaskIds(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if count > 0 {
		return controller.JSONBaseErrorReq(c, errTaskHasBeenUsed)
	}

	template, exist, err := s.GetWorkflowTemplateById(workflowTemplateId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("the task instance is not bound workflow template")))
	}

	err = v1.CheckWorkflowCanCommit(template, tasks)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	stepTemplates, err := s.GetWorkflowStepsByTemplateId(template.ID)
	if err != nil {
		return err
	}
	err = s.CreateWorkflowV2(req.Subject, workflowId, req.Desc, user, tasks, stepTemplates, project.ID)
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

	workFlowId := strconv.Itoa(int(workflow.ID))
	go notification.NotifyWorkflow(workFlowId, notification.WorkflowNotifyTypeCreate)

	go im.CreateApprove(workFlowId)

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
// @Summary 更新工单（驳回后才可更新）
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

	projectName := c.Param("project_name")
	workflowId := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewDataNotExistErr("workflow not exist"))
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, &model.Project{Name: projectName}, workflow, []uint{})
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

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskIds := make([]uint, len(tasks))
	for i, task := range tasks {
		taskIds[i] = task.ID

		count, err := s.GetTaskSQLCountByTaskID(task.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if count == 0 {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("task's execute sql is null. taskId=%v", task.ID)))
		}

		err = v1.CheckCurrentUserCanViewTask(c, task)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		if task.Instance == nil {
			return controller.JSONBaseErrorReq(c, v1.ErrInstanceNotExist)
		}

		if user.ID != task.CreateUserId {
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

	workflowIdStr := fmt.Sprintf("%v", workflow.ID)
	workflow, exist, err = s.GetWorkflowDetailById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	if workflow.Record.Status != model.WorkflowStatusReject {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)))
	}

	if user.ID != workflow.CreateUserId {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("you are not allow to operate the workflow")))
	}

	template, exist, err := s.GetWorkflowTemplateById(tasks[0].Instance.WorkflowTemplateId)
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
	go notification.NotifyWorkflow(workflowIdStr, notification.WorkflowNotifyTypeCreate)

	workFlowId := strconv.Itoa(int(workflow.ID))
	go im.CreateApprove(workFlowId)

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateWorkflowScheduleReqV2 struct {
	ScheduleTime *time.Time `json:"schedule_time"`
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
	projectName := c.Param("project_name")
	workflowId := c.Param("workflow_id")

	s := model.GetStorage()

	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	err = v1.CheckCurrentUserCanOperateWorkflow(c, &model.Project{Name: projectName}, workflow, []uint{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflowId = strconv.Itoa(int(workflow.ID))

	taskId := c.Param("task_id")
	taskIdUint, err := v1.FormatStringToUint64(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	req := new(UpdateWorkflowScheduleReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflow, exist, err = s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}
	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, _err.New("workflow current step not found")))
	}

	if workflow.Record.Status != model.WorkflowStatusWaitForExecution {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
			fmt.Errorf("workflow need to be approved first")))
	}

	err = v1.CheckUserCanOperateStep(user, workflow, int(currentStep.ID))
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

	instance, exist, err := s.GetInstanceById(fmt.Sprintf("%v", curTaskRecord.InstanceId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrInstanceNotExist)
	}

	if req.ScheduleTime != nil && len(instance.MaintenancePeriod) != 0 && !instance.MaintenancePeriod.IsWithinScope(*req.ScheduleTime) {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowExecuteTimeIncorrect)
	}

	err = s.UpdateInstanceRecordSchedule(curTaskRecord, user.ID, req.ScheduleTime)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
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
	projectName := c.Param("project_name")
	workflowId := c.Param("workflow_id")

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	workflowId = fmt.Sprintf("%v", workflow.ID)

	workflow, exist, err = s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err := v1.PrepareForWorkflowExecution(c, projectName, workflow, user); err != nil {
		return err
	}

	needExecTaskIds, err := v1.GetNeedExecTaskIds(s, workflow, user)
	if err != nil {
		return err
	}

	err = server.ExecuteWorkflow(workflow, needExecTaskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

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
	Steps             []*WorkflowStepResV2 `json:"workflow_step_list,omitempty"`
}

type WorkflowResV2 struct {
	Name          string                 `json:"workflow_name"`
	WorkflowID    string                 `json:"workflow_id"`
	Desc          string                 `json:"desc,omitempty"`
	Mode          string                 `json:"mode" enums:"same_sqls,different_sqls"`
	CreateUser    string                 `json:"create_user_name"`
	CreateTime    *time.Time             `json:"create_time"`
	Record        *WorkflowRecordResV2   `json:"record"`
	RecordHistory []*WorkflowRecordResV2 `json:"record_history_list,omitempty"`
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
	projectName := c.Param("project_name")
	workflowID := c.Param("workflow_id")

	s := model.GetStorage()

	err := CheckCurrentUserCanViewWorkflow(c, workflowID, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, exist, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, workflowID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	workflowIdStr := strconv.Itoa(int(workflow.ID))

	workflow, exist, err = s.GetWorkflowDetailById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrWorkflowNoAccess)
	}

	history, err := s.GetWorkflowHistoryById(workflowIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflow.RecordHistory = history

	return c.JSON(http.StatusOK, &GetWorkflowResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowToRes(workflow),
	})
}

func convertWorkflowToRes(workflow *model.Workflow) *WorkflowResV2 {
	workflowRes := &WorkflowResV2{
		Name:       workflow.Subject,
		WorkflowID: workflow.WorkflowId,
		Desc:       workflow.Desc,
		Mode:       workflow.Mode,
		CreateUser: workflow.CreateUser.Name,
		CreateTime: &workflow.CreatedAt,
	}

	// convert workflow record
	workflowRecordRes := convertWorkflowRecordToRes(workflow, workflow.Record)

	// convert workflow record history
	recordHistory := make([]*WorkflowRecordResV2, 0, len(workflow.RecordHistory))
	for _, record := range workflow.RecordHistory {
		recordRes := convertWorkflowRecordToRes(workflow, record)
		recordHistory = append(recordHistory, recordRes)
	}
	workflowRes.RecordHistory = recordHistory
	workflowRes.Record = workflowRecordRes

	return workflowRes
}

func convertWorkflowRecordToRes(workflow *model.Workflow, record *model.WorkflowRecord) *WorkflowRecordResV2 {
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
		OperationUser: workflow.CreateUserName(),
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

	return &WorkflowRecordResV2{
		Tasks:             tasksRes,
		CurrentStepNumber: currentStepNum,
		Status:            record.Status,
		Steps:             steps,
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
	}
	if step.OperationUser != nil {
		stepRes.OperationUser = step.OperationUser.Name
	}
	if step.Assignees != nil {
		for _, user := range step.Assignees {
			stepRes.Users = append(stepRes.Users, user.Name)
		}
	}
	return stepRes
}

func CheckCurrentUserCanViewWorkflow(c echo.Context, workflowID, projectName string) error {
	userName := controller.GetUserName(c)
	s := model.GetStorage()
	isManager, err := s.IsProjectManager(userName, projectName)
	if err != nil {
		return err
	}
	if userName == model.DefaultAdminUser || isManager {
		return nil
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	workflow, _, err := s.GetWorkflowByProjectNameAndWorkflowId(projectName, workflowID)
	if err != nil {
		return err
	}

	access, err := s.UserCanAccessWorkflow(user, workflow)
	if err != nil {
		return err
	}
	if access {
		return nil
	}
	instances, err := s.GetInstancesByWorkflowID(workflow.ID)
	if err != nil {
		return err
	}
	ok, err := s.CheckUserHasOpToAnyInstance(user, instances, []uint{model.OP_WORKFLOW_VIEW_OTHERS})
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return v1.ErrWorkflowNoAccess
}
