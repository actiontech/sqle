//go:build !enterprise
// +build !enterprise

package model

func (s *Storage) UpdateStageWorkflowExecTimeIfNeed(workflowId string) error {
	return nil
}

func (s *Storage) UpdateStageWorkflowIfNeed(workflowId string, workflowStage map[string]interface{}) error {
	return nil
}

func (stage SqlVersionStage) InitialStatusOfWorkflow() string {
	return ""
}

func (s *Storage) GetAssociatedStageWorkflows(workflowId string) ([]*AssociatedStageWorkflow, error) {
	return nil, nil
}

func (s *Storage) GetSQLVersionByWorkflowId(workflowId string) (*SqlVersion, error) {
	return &SqlVersion{}, nil
}
