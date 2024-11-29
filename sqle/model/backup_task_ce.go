//go:build !enterprise
// +build !enterprise

package model

func (s *Storage) BatchCreateBackupTasks(backupTasks []*BackupTask) error {
	return nil
}

func (s *Storage) GetBackupTaskByExecuteSqlId(executeSqlId uint) (*BackupTask, error){
	return &BackupTask{},nil
}

func (s *Storage) GetExecuteSqlRollbackWorkflowRelationByTaskId(taskId uint) ([]*ExecuteSqlRollbackWorkflowsRelation, error) {
	return []*ExecuteSqlRollbackWorkflowsRelation{}, nil
}

func (s *Storage) GetRollbackWorkflowByOriginalWorkflowId(workflowId string) ([]*RollbackWorkflowOriginalWorkflowsRelation, error) {
	return []*RollbackWorkflowOriginalWorkflowsRelation{}, nil
}