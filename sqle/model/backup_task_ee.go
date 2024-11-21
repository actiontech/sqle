//go:build enterprise
// +build enterprise

package model

/*
the default batch size is 100, if the length of input backup task slice is shorter than 100, the batch will set as the length of slice
*/
func (s *Storage) BatchCreateBackupTasks(backupTasks []*BackupTask) error {
	var batchSize = 100
	if len(backupTasks) < 100 {
		batchSize = len(backupTasks)
	}
	return s.db.CreateInBatches(backupTasks, batchSize).Error
}

/*
update backup status and backup execute information
*/
func (s *Storage) UpdateBackupExecuteResult(task *BackupTask) error {
	return s.db.Model(&BackupTask{}).
		Where("id = ?", task.ID).
		UpdateColumns(map[string]interface{}{
			"backup_status":    task.BackupStatus,
			"backup_exec_result": task.BackupExecResult,
		}).Error
}

func (s *Storage) GetBackupTaskByExecuteSqlId(executeSqlId uint) (*BackupTask, error) {
	var backupTask BackupTask
	err := s.db.Model(&BackupTask{}).Where("execute_sql_id = ?", executeSqlId).First(&backupTask).Error
	if err != nil {
		return nil, err
	}
	return &backupTask, nil
}


func (s *Storage) GetBackupTaskByTaskId(taskId uint) ([]*BackupTask, error) {
	var backupTasks []*BackupTask
	err := s.db.Model(&BackupTask{}).Where("task_id = ?", taskId).Find(&backupTasks).Error
	if err != nil {
		return nil, err
	}
	return backupTasks, nil
}

func (s *Storage) GetRollbackSqlByTaskId(taskId uint) ([]*RollbackSQL, error) {
	var rollbackSqls []*RollbackSQL
	err := s.db.Model(&RollbackSQL{}).Where("task_id = ?", taskId).Find(&rollbackSqls).Error
	if err != nil {
		return nil, err
	}
	return rollbackSqls, nil
}
