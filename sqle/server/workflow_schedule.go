package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/sirupsen/logrus"
)

var ErrWorkflowNoAccess = errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist or you can't access it"))

type WorkflowScheduleJob struct {
	BaseJob
}

func NewWorkflowScheduleJob(entry *logrus.Entry) ServerJob {
	entry = entry.WithField("job", "schedule_workflow")
	j := &WorkflowScheduleJob{}
	j.BaseJob = *NewBaseJob(entry, 5*time.Second, j.WorkflowSchedule)
	return j
}

func (j *WorkflowScheduleJob) WorkflowSchedule(entry *logrus.Entry) {
	st := model.GetStorage()
	workflows, err := st.GetNeedScheduledWorkflows()
	if err != nil {
		entry.Errorf("get need scheduled workflows from storage error: %v", err)
		return
	}

	for _, workflow := range workflows {
		w, err := dms.GetWorkflowDetailByWorkflowId(string(workflow.ProjectId), workflow.WorkflowId, st.GetWorkflowDetailWithoutInstancesByWorkflowID)
		if err != nil {
			entry.Errorf("get workflow from storage error: %v", err)
			return
		}
		needExecuteTaskIds, err := w.GetNeedScheduledTaskIds(entry)
		if err != nil {
			entry.Errorf("get need scheduled taskIds error: %v", err)
			return
		}

		entry.Infof("start to execute scheduled workflow %s", w.Subject)
		if len(needExecuteTaskIds) == 0 {
			entry.Warnf("workflow %s need to execute scheduled, but no task find", w.Subject)
		}
		_, err = ExecuteWorkflow(w, needExecuteTaskIds)
		if err != nil {
			entry.Errorf("execute scheduled workflow %s error: %v", w.Subject, err)
		} else {
			entry.Infof("execute scheduled workflow %s success", w.Subject)
		}
	}
}

// ExecuteWorkflow 执行工单中的任务
// 参数说明:
//   - workflow: 要执行的工单对象
//   - needExecTaskIdToUserId: 需要执行的任务ID到用户ID的映射
//   - isAutoCreated: 可变参数，用于标识工单是否为自动创建
//     * 使用场景:
//       - 当工单是通过 AutoCreateAndExecuteWorkflowV1 等自动创建接口创建时，应传递 true
//       - 当工单是通过普通创建流程（如 CreateWorkflowV2）创建时，不传递此参数或传递 false
//     * 默认行为:
//       - 如果不传递此参数（即 len(isAutoCreated) == 0），则默认为 false，使用普通工单的通知类型
//       - 如果传递 true，则使用自动创建工单的特殊通知类型（auto_exec_success/auto_exec_failed）
//       - 如果传递 false，则使用普通工单的通知类型（exec_success/exec_failed）
//     * 通知类型差异:
//       - 自动创建工单: WorkflowNotifyTypeAutoExecuteSuccess/AutoExecuteFail
//       - 普通工单: WorkflowNotifyTypeExecuteSuccess/ExecuteFail
//     * 注意: 此参数仅影响通知的 action 类型，不影响工单的实际执行逻辑
func ExecuteWorkflow(workflow *model.Workflow, needExecTaskIdToUserId map[uint]string, isAutoCreated ...bool) (chan string, error) {
	s := model.GetStorage()
	l := log.NewEntry()
	err := s.UpdateStageWorkflowExecTimeIfNeed(workflow.WorkflowId)
	if err != nil {
		l.Errorf("update workflow execute time for version stage error: %v", err)
	}
	// get task and check connection before to execute it.
	for taskId := range needExecTaskIdToUserId {
		taskId := fmt.Sprintf("%d", taskId)
		task, exist, err := s.GetTaskDetailById(taskId)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, errors.New(errors.DataNotExist, fmt.Errorf("task is not exist. taskID=%v", taskId))
		}
		instance, exist, err := dms.GetInstancesById(context.Background(), fmt.Sprintf("%d", task.InstanceId))
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist. instanceId=%v", task.InstanceId))
		}
		task.Instance = instance
		if task.Instance == nil {
			return nil, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))
		}

		// if instance is not connectable, exec sql must be failed;
		// commit action unable to retry, so don't to exec it.
		if err = common.CheckInstanceIsConnectable(task.Instance); err != nil {
			return nil, errors.New(errors.ConnectRemoteDatabaseError, err)
		}
	}

	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return nil, fmt.Errorf("workflow current step not found")
	}

	// update workflow
	needExecTaskRecords := make([]*model.WorkflowInstanceRecord, 0, len(needExecTaskIdToUserId))
	for _, inst := range workflow.Record.InstanceRecords {
		if userId, ok := needExecTaskIdToUserId[inst.TaskId]; ok {
			inst.IsSQLExecuted = true
			inst.ExecutionUserId = userId
			needExecTaskRecords = append(needExecTaskRecords, inst)
		}
	}

	var operateStep *model.WorkflowStep
	// 只有当所有数据源都执行上线操作时，current step状态才改为"approved"
	allTaskHasExecuted := true
	for _, inst := range workflow.Record.InstanceRecords {
		if !inst.IsSQLExecuted {
			allTaskHasExecuted = false
		}
	}
	if allTaskHasExecuted {
		currentStep.State = model.WorkflowStepStateApprove
		workflow.Record.Status = model.WorkflowStatusExecuting
		workflow.Record.CurrentWorkflowStepId = 0
		operateStep = currentStep
	} else {
		operateStep = nil
	}

	err = s.UpdateWorkflowExecInstanceRecord(workflow, operateStep, needExecTaskRecords)
	if err != nil {
		return nil, err
	}
	workflowStatusChan := make(chan string, 1)
	var lock sync.Mutex
	for taskId := range needExecTaskIdToUserId {
		id := taskId
		go func() {
			sqledServer := GetSqled()
			task, err := sqledServer.AddTaskWaitResult(string(workflow.ProjectId), strconv.Itoa(int(id)), ActionTypeExecute)

			{ // NOTE: Update the workflow status before sending notifications to ensure that the notification content reflects the latest information.
				lock.Lock()
				updateStatus(s, workflow, l, workflowStatusChan)
				lock.Unlock()
			}

			// 判断是否为自动创建的工单
			// 逻辑说明:
			//   - 如果 isAutoCreated 参数未传递（len == 0），则 isAuto = false，使用普通工单通知类型
			//   - 如果 isAutoCreated[0] == true，则 isAuto = true，使用自动创建工单的特殊通知类型
			//   - 如果 isAutoCreated[0] == false，则 isAuto = false，使用普通工单通知类型
			// 通知类型说明:
			//   - 自动创建工单成功: WorkflowNotifyTypeAutoExecuteSuccess -> action: "auto_exec_success"
			//   - 自动创建工单失败: WorkflowNotifyTypeAutoExecuteFail -> action: "auto_exec_failed"
			//   - 普通工单成功: WorkflowNotifyTypeExecuteSuccess -> action: "exec_success"
			//   - 普通工单失败: WorkflowNotifyTypeExecuteFail -> action: "exec_failed"
			isAuto := len(isAutoCreated) > 0 && isAutoCreated[0]
			if err != nil || task.Status == model.TaskStatusExecuteFailed {
				if isAuto {
					go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeAutoExecuteFail)
				} else {
					go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeExecuteFail)
				}
			} else {
				if isAuto {
					go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeAutoExecuteSuccess)
				} else {
					go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeExecuteSuccess)
				}
			}

		}()
	}

	return workflowStatusChan, nil
}

func updateStatus(s *model.Storage, workflow *model.Workflow, l *logrus.Entry, workflowStatusChan chan string) {
	tasks, err := s.GetTasksByWorkFlowRecordID(workflow.Record.ID)
	if err != nil {
		l.Errorf("get tasks by workflow record id error: %v", err)
	}

	var workFlowStatus string

	var hasExecuting bool
	var hasExecuteFailed bool
	var hasWaitExecute bool

	for _, task := range tasks {
		if task.Status == model.TaskStatusExecuting {
			hasExecuting = true
		}
		if task.Status == model.TaskStatusExecuteFailed {
			hasExecuteFailed = true
		}
		if task.Status == model.TaskStatusAudited {
			hasWaitExecute = true
		}

		// termination by user
		if task.Status == model.TaskStatusTerminating ||
			task.Status == model.TaskStatusTerminateSucc ||
			task.Status == model.TaskStatusTerminateFail {
			hasExecuteFailed = true
		}
	}

	if hasWaitExecute {
		workFlowStatus = model.WorkflowStatusWaitForExecution
	} else if hasExecuting {
		workFlowStatus = model.WorkflowStatusExecuting
	} else if hasExecuteFailed {
		workFlowStatus = model.WorkflowStatusExecFailed
	} else {
		workFlowStatus = model.WorkflowStatusFinish
	}

	if workFlowStatus != "" {
		err = s.UpdateWorkflowRecordByID(workflow.Record.ID, map[string]interface{}{
			"status": workFlowStatus,
		})
		if err != nil {
			l.Errorf("update workflow record status failed: %v", err)
		}
		if workflowStatusChan != nil && workFlowStatus != model.WorkflowStatusExecuting && workFlowStatus != model.WorkflowStatusWaitForExecution {
			workflowStatusChan <- workFlowStatus
		}
	}
}

func ApproveWorkflowProcess(workflow *model.Workflow, user *model.User, s *model.Storage) error {
	currentStep := workflow.CurrentStep()

	if workflow.Record.Status == model.WorkflowStatusWaitForExecution {
		return errors.New(errors.DataInvalid,
			fmt.Errorf("workflow has been approved, you should to execute it"))
	}

	currentStep.State = model.WorkflowStepStateApprove
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.GetIDStr()
	nextStep := workflow.NextStep()
	workflow.Record.CurrentWorkflowStepId = nextStep.ID
	if nextStep.Template.Typ == model.WorkflowStepTypeSQLExecute {
		workflow.Record.Status = model.WorkflowStatusWaitForExecution
	}

	err := s.UpdateWorkflowStep(workflow, currentStep)
	if err != nil {
		return fmt.Errorf("update workflow status failed, %v", err)
	}

	go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeApprove)

	return nil
}

func RejectWorkflowProcess(workflow *model.Workflow, reason string, user *model.User, s *model.Storage) error {
	currentStep := workflow.CurrentStep()
	currentStep.State = model.WorkflowStepStateReject
	currentStep.Reason = reason
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.GetIDStr()

	workflow.Record.Status = model.WorkflowStatusReject
	workflow.Record.CurrentWorkflowStepId = 0

	if err := s.UpdateWorkflowStep(workflow, currentStep); err != nil {
		return fmt.Errorf("update workflow status failed, %v", err)
	}

	go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeReject)

	return nil
}

// ExecuteTasksProcess 执行工单任务处理流程
// 参数说明:
//   - workflowId: 工单ID
//   - projectUid: 项目UID
//   - user: 执行用户
//   - isAutoCreated: 可变参数，用于标识工单是否为自动创建
//     * 使用场景: 与 ExecuteWorkflow 的 isAutoCreated 参数相同
//     * 默认行为: 如果不传递此参数，则默认为 false，表示普通工单
//     * 传递方式: 此参数会透传给 ExecuteWorkflow 函数，用于控制通知类型
//     * 示例:
//       - 自动创建工单: ExecuteTasksProcess(workflowId, projectUid, user, true)
//       - 普通工单: ExecuteTasksProcess(workflowId, projectUid, user) 或 ExecuteTasksProcess(workflowId, projectUid, user, false)
func ExecuteTasksProcess(workflowId string, projectUid string, user *model.User, isAutoCreated ...bool) (chan string, error) {
	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		return nil, err
	}

	if err = PrepareForWorkflowExecution(projectUid, workflow, user); err != nil {
		return nil, err
	}

	needExecTaskIds, err := GetNeedExecTaskIds(workflow, user)
	if err != nil {
		return nil, err
	}

	workflowExecResultChan, err := ExecuteWorkflow(workflow, needExecTaskIds, isAutoCreated...)
	if err != nil {
		return nil, err
	}

	return workflowExecResultChan, nil
}

func PrepareForWorkflowExecution(projectUid string, workflow *model.Workflow, user *model.User) error {
	err := CheckCurrentUserCanOperateWorkflowByUser(user, projectUid, workflow, []dmsV1.OpPermissionType{})
	if err != nil {
		return err
	}

	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return errors.New(errors.DataInvalid, fmt.Errorf("workflow current step not found"))
	}

	if workflow.Record.Status != model.WorkflowStatusWaitForExecution {
		return errors.New(errors.DataInvalid,
			fmt.Errorf("workflow need to be approved first"))
	}

	err = CheckUserCanOperateStep(user, workflow, int(currentStep.ID))
	if err != nil {
		return errors.New(errors.DataInvalid, err)
	}
	return nil
}

func GetNeedExecTaskIds(workflow *model.Workflow, user *model.User) (taskIds map[uint] /*task id*/ string /*user id*/, err error) {
	instances := make([]*model.Instance, 0, len(workflow.Record.InstanceRecords))
	for _, item := range workflow.Record.InstanceRecords {
		instances = append(instances, item.Instance)
	}
	// 有不在运维时间内的instances报错
	var cannotExecuteInstanceNames []string
	for _, inst := range instances {
		if len(inst.MaintenancePeriod) != 0 && !inst.MaintenancePeriod.IsWithinScope(time.Now()) {
			cannotExecuteInstanceNames = append(cannotExecuteInstanceNames, inst.Name)
		}
	}
	if len(cannotExecuteInstanceNames) > 0 {
		return nil, errors.New(errors.TaskActionInvalid,
			fmt.Errorf("please go online during instance operation and maintenance time. these instances are not in maintenance time[%v]", strings.Join(cannotExecuteInstanceNames, ",")))
	}

	// 定时的instances和已上线的跳过
	needExecTaskIds := make(map[uint]string)
	for _, instRecord := range workflow.Record.InstanceRecords {
		if instRecord.ScheduledAt != nil || instRecord.IsSQLExecuted {
			continue
		}
		needExecTaskIds[instRecord.TaskId] = user.GetIDStr()
	}
	return needExecTaskIds, nil
}

func CheckCurrentUserCanOperateWorkflowByUser(user *model.User, projectUid string, workflow *model.Workflow, ops []dmsV1.OpPermissionType) error {
	if user.Name == model.DefaultAdminUser || user.Name == model.DefaultSysUser {
		return nil
	}

	s := model.GetStorage()
	up, err := dms.NewUserPermission(user.GetIDStr(), projectUid)
	if err != nil {
		return err
	}
	canOpProject := up.CanOpProject()
	if canOpProject {
		return nil
	}

	access, err := s.UserCanAccessWorkflow(user.GetIDStr(), workflow.WorkflowId)
	if err != nil {
		return err
	}
	if access {
		return nil
	}
	if len(ops) > 0 {
		for _, item := range workflow.Record.InstanceRecords {
			if !up.CanOpInstanceNoAdmin(item.Instance.GetIDStr(), ops...) {
				return ErrWorkflowNoAccess
			}
		}
	}

	return ErrWorkflowNoAccess
}

func CheckUserCanOperateStep(user *model.User, workflow *model.Workflow, stepId int) error {
	if workflow.Record.Status != model.WorkflowStatusWaitForAudit && workflow.Record.Status != model.WorkflowStatusWaitForExecution {
		return fmt.Errorf("workflow status is %s, not allow operate it", workflow.Record.Status)
	}

	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return fmt.Errorf("workflow current step not found")
	}
	if uint(stepId) != workflow.CurrentStep().ID {
		return fmt.Errorf("workflow current step is not %d", stepId)
	}

	// sys 用户和 admin 用户可以直接操作工单
	if user.Name == model.DefaultSysUser || user.Name == model.DefaultAdminUser {
		return nil
	}

	if !workflow.IsOperationUser(user) {
		return fmt.Errorf("you are not allow to operate the workflow")
	}

	return nil
}
