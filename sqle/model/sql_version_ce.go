//go:build !enterprise
// +build !enterprise

package model

func (s *Storage) UpdateStageWorkflowExecTimeIfNeed(workflowId string) error {
	return nil
}

func (stage SqlVersionStage) InitialStatusOfWorkflow() string {
	return ""
}
