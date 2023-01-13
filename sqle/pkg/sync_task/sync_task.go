package sync_task

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/model"
)

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
