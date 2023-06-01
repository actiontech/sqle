package server

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"

	"github.com/sirupsen/logrus"
)

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
	now := time.Now()
	for _, workflow := range workflows {
		w, exist, err := st.GetWorkflowDetailById(strconv.Itoa(int(workflow.ID)))
		if err != nil {
			entry.Errorf("get workflow from storage error: %v", err)
			return
		}
		if !exist {
			entry.Errorf("workflow %s not found", workflow.Subject)
			return
		}

		currentStep := w.CurrentStep()
		if currentStep == nil {
			entry.Errorf("workflow %s not found", w.Subject)
			return
		}
		if currentStep.Template.Typ != model.WorkflowStepTypeSQLExecute {
			entry.Errorf("workflow %s need to be approved first", w.Subject)
			return
		}

		entry.Infof("start to execute scheduled workflow %s", w.Subject)
		needExecuteTaskIds := map[uint]uint{}
		for _, ir := range w.Record.InstanceRecords {
			if !ir.IsSQLExecuted && ir.ScheduledAt != nil && ir.ScheduledAt.Before(now) {
				needExecuteTaskIds[ir.TaskId] = ir.ScheduleUserId
			}
		}
		if len(needExecuteTaskIds) == 0 {
			entry.Warnf("workflow %s need to execute scheduled, but no task find", w.Subject)
		}
		err = ExecuteWorkflow(w, needExecuteTaskIds)
		if err != nil {
			entry.Errorf("execute scheduled workflow %s error: %v", w.Subject, err)
		} else {
			entry.Infof("execute scheduled workflow %s success", w.Subject)
		}
	}
}

func ExecuteWorkflow(workflow *model.Workflow, needExecTaskIdToUserId map[uint]uint) error {
	s := model.GetStorage()

	// get task and check connection before to execute it.
	for taskId := range needExecTaskIdToUserId {
		taskId := fmt.Sprintf("%d", taskId)
		task, exist, err := s.GetTaskDetailById(taskId)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New(errors.DataNotExist, fmt.Errorf("task is not exist. taskID=%v", taskId))
		}
		if task.Instance == nil {
			return errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))
		}

		// if instance is not connectable, exec sql must be failed;
		// commit action unable to retry, so don't to exec it.
		if err = common.CheckInstanceIsConnectable(task.Instance); err != nil {
			return errors.New(errors.ConnectRemoteDatabaseError, err)
		}
	}

	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return fmt.Errorf("workflow current step not found")
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

	err := s.UpdateWorkflowExecInstanceRecord(workflow, operateStep, needExecTaskRecords)
	if err != nil {
		return err
	}

	l := log.NewEntry()
	var lock sync.Mutex
	for taskId := range needExecTaskIdToUserId {
		id := taskId
		go func() {
			sqledServer := GetSqled()
			task, err := sqledServer.AddTaskWaitResult(strconv.Itoa(int(id)), ActionTypeExecute)

			{ // NOTE: Update the workflow status before sending notifications to ensure that the notification content reflects the latest information.
				lock.Lock()
				updateStatus(s, workflow, l)
				lock.Unlock()
			}

			if err != nil || task.Status == model.TaskStatusExecuteFailed {
				go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeExecuteFail)
			} else {
				go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeExecuteSuccess)
			}

		}()
	}

	return nil
}

func updateStatus(s *model.Storage, workflow *model.Workflow, l *logrus.Entry) {
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
	currentStep.OperationUserId = user.ID
	nextStep := workflow.NextStep()
	workflow.Record.CurrentWorkflowStepId = nextStep.ID
	if nextStep.Template.Typ == model.WorkflowStepTypeSQLExecute {
		workflow.Record.Status = model.WorkflowStatusWaitForExecution
	}

	err := s.UpdateWorkflowStep(workflow, currentStep)
	if err != nil {
		return fmt.Errorf("update workflow status failed, %v", err)
	}

	go notification.NotifyWorkflow(strconv.Itoa(int(workflow.ID)), notification.WorkflowNotifyTypeApprove)

	return nil
}

func RejectWorkflowProcess(workflow *model.Workflow, reason string, user *model.User, s *model.Storage) error {
	currentStep := workflow.CurrentStep()
	currentStep.State = model.WorkflowStepStateReject
	currentStep.Reason = reason
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.ID

	workflow.Record.Status = model.WorkflowStatusReject
	workflow.Record.CurrentWorkflowStepId = 0

	if err := s.UpdateWorkflowStep(workflow, currentStep); err != nil {
		return fmt.Errorf("update workflow status failed, %v", err)
	}

	go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeReject)

	return nil
}
