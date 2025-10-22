package server

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
)

func ReExecuteTaskSQLs(workflow *model.Workflow, task *model.Task, execSqlIds []uint, user *model.User) error {
	s := model.GetStorage()
	l := log.NewEntry()

	instance, exist, err := dms.GetInstancesById(context.Background(), fmt.Sprintf("%d", task.InstanceId))
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist. instanceId=%v", task.InstanceId))
	}
	task.Instance = instance
	if task.Instance == nil {
		return errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))
	}

	// if instance is not connectable, exec sql must be failed;
	// commit action unable to retry, so don't to exec it.
	if err = common.CheckInstanceIsConnectable(task.Instance); err != nil {
		return errors.New(errors.ConnectRemoteDatabaseError, err)
	}

	needExecTaskRecords := make([]*model.WorkflowInstanceRecord, 0, len(workflow.Record.InstanceRecords))
	// update workflow
	for _, inst := range workflow.Record.InstanceRecords {
		if inst.TaskId != task.ID {
			continue
		}
		inst.IsSQLExecuted = true
		inst.ExecutionUserId = user.GetIDStr()
		needExecTaskRecords = append(needExecTaskRecords, inst)
	}

	workflow.Record.Status = model.WorkflowStatusExecuting
	workflow.Record.CurrentWorkflowStepId = 0

	err = s.UpdateWorkflowExecInstanceRecordForReExecute(workflow, needExecTaskRecords)
	if err != nil {
		return err
	}
	var lock sync.Mutex
	go func() {
		sqledServer := GetSqled()
		task, err := sqledServer.AddTaskWaitResultWithSQLIds(string(workflow.ProjectId), strconv.Itoa(int(task.ID)), execSqlIds, ActionTypeExecute)
		{
			lock.Lock()
			updateStatus(s, workflow, l, nil)
			lock.Unlock()
		}
		if err != nil || task.Status == model.TaskStatusExecuteFailed {
			go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeExecuteFail)
		} else {
			go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeExecuteSuccess)
		}

	}()

	return nil
}
