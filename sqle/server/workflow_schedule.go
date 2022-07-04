package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/notification"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

func (s *Sqled) workflowScheduleLoop() {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	entry := log.NewEntry().WithField("type", "schedule_workflow")
	for {
		select {
		case <-s.exit:
			return
		case <-tick.C:
			s.WorkflowSchedule(entry)
		}
	}
}

func (s *Sqled) WorkflowSchedule(entry *logrus.Entry) {
	st := model.GetStorage()
	workflows, err := st.GetNeedScheduledWorkflows()
	if err != nil {
		entry.Errorf("get need scheduled workflows from storage error: %v", err)
		return
	}
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
		err = ExecuteWorkflow(w, w.Record.ScheduleUserId)
		if err != nil {
			entry.Errorf("execute scheduled workflow %s error: %v", w.Subject, err)
		} else {
			entry.Infof("execute scheduled workflow %s success", w.Subject)
		}
	}
}

func ExecuteWorkflow(workflow *model.Workflow, userId uint) error {
	s := model.GetStorage()

	// get task and check connection before to execute it.
	taskId := fmt.Sprintf("%d", workflow.Record.TaskId)
	task, exist, err := s.GetTaskDetailById(taskId)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(errors.DataNotExist, fmt.Errorf("task is not exist"))
	}
	if task.Instance == nil {
		return errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))
	}

	// if instance is not connectable, exec sql must be failed;
	// commit action unable to retry, so don't to exec it.
	dsn := &driver.DSN{
		Host:             task.Instance.Host,
		Port:             task.Instance.Port,
		User:             task.Instance.User,
		Password:         task.Instance.Password,
		AdditionalParams: task.Instance.AdditionalParams,
		DatabaseName:     task.Schema,
	}

	cfg, err := driver.NewConfig(dsn, nil)
	if err != nil {
		return errors.New(errors.LoadDriverFail, err)
	}

	drvMgr, err := driver.NewDriverManger(log.NewEntry(), task.DBType, cfg)
	if err != nil {
		return errors.New(errors.LoadDriverFail, err)
	}
	defer drvMgr.Close(context.TODO())
	d, err := drvMgr.GetAuditDriver()
	if err != nil {
		return errors.New(errors.LoadDriverFail, err)
	}
	if err := d.Ping(context.TODO()); err != nil {
		return errors.New(errors.ConnectRemoteDatabaseError, err)
	}

	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return fmt.Errorf("workflow current step not found")
	}
	// update workflow
	currentStep.State = model.WorkflowStepStateApprove
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = userId
	workflow.Record.Status = model.WorkflowStatusFinish
	workflow.Record.CurrentWorkflowStepId = 0

	err = s.UpdateWorkflowStatus(workflow, currentStep)
	if err != nil {
		return err
	}
	go func() {
		sqledServer := GetSqled()
		task, err := sqledServer.AddTaskWaitResult(taskId, ActionTypeExecute)
		if err != nil || task.Status == model.TaskStatusExecuteFailed {
			go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeExecuteFail)
		} else {
			go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeExecuteSuccess)
		}
	}()
	return nil
}
