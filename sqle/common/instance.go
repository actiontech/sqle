package common

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
)

func CheckInstanceIsConnectable(instance *model.Instance) error {
	plugin, err := NewDriverManagerWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return err
	}
	defer plugin.Close(context.TODO())

	if err := plugin.Ping(context.TODO()); err != nil {
		return err
	}

	return nil
}

func CheckDeleteInstance(instanceId uint) error {
	s := model.GetStorage()

	isUnFinished, err := s.IsWorkflowUnFinishedByInstanceId(instanceId)
	if err != nil {
		return fmt.Errorf("check instance %d is finished failed: %v", instanceId, err)
	}
	if isUnFinished {
		return fmt.Errorf("instance %d is running,cannot be deleted", instanceId)
	}

	return nil
}
