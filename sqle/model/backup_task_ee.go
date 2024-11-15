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
			"backup_exec_info": task.BackupExecInfo,
		}).Error
}
