package common

import (
	"context"
	"errors"
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

func CheckDeleteInstance(instanceId int64) error {
	s := model.GetStorage()

	isUnFinished, err := s.IsWorkflowUnFinishedByInstanceId(instanceId)
	if err != nil {
		return fmt.Errorf("check if all workflows are finished failed: %v", err)
	}
	if isUnFinished {
		return errors.New("has unfinished workflows")
	}
	_, isBoundToAuditPlan, err := s.GetInstanceAuditPlanByInstanceID(instanceId)
	if err != nil {
		return fmt.Errorf("check that all audit plans are unbound failed: %v", err)
	}
	if isBoundToAuditPlan {
		return errors.New("there is an unbound audit plan")
	}
	nodes, err := s.GetPipelineNodesByInstanceId(uint64(instanceId))
	if err != nil {
		return fmt.Errorf("check that all pipeline node are unbound failed: %v", err)
	}
	if len(nodes) > 0 {
		return errors.New("there is an unbound pipeline node")
	}
	return nil
}
