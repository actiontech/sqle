package common

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
)

func CheckInstanceIsConnectable(instance *model.Instance) error {
	drvMgr, err := NewDriverManagerWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return err
	}
	defer drvMgr.Close(context.TODO())

	d, err := drvMgr.GetAuditDriver()
	if err != nil {
		return err
	}

	if err := d.Ping(context.TODO()); err != nil {
		return err
	}

	return nil
}

func CheckDeleteInstance(instanceId uint) error {
	s := model.GetStorage()

	tasks, err := s.GetTaskByInstanceId(instanceId)
	if err != nil {
		return err
	}
	taskIds := make([]uint, 0, len(tasks))
	for _, task := range tasks {
		taskIds = append(taskIds, task.ID)
	}
	isRunning, err := s.TaskWorkflowIsUnfinished(taskIds)
	if err != nil {
		return err
	}
	if isRunning {
		return fmt.Errorf("instance %d is running,cannot be deleted", instanceId)
	}

	return nil
}
