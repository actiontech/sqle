//go:build enterprise
// +build enterprise

package model


func (s *Storage) BatchCreateBackupTasks(backupTasks []*BackupTask) error {
	var batchSize = 100
	if len(backupTasks) < 100 {
		batchSize = len(backupTasks)
	}
	return s.db.CreateInBatches(backupTasks, batchSize).Error
}