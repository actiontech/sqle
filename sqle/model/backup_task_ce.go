//go:build !enterprise
// +build !enterprise

package model

func (s *Storage) BatchCreateBackupTasks(backupTasks []*BackupTask) error {
	return nil
}
