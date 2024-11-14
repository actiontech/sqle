//go:build !enterprise
// +build !enterprise

package server

type BackupService struct{}

func (BackupService) CheckBackupConflictWithExecMode(EnableBackup bool, ExecMode string) error {
	return nil
}

func (BackupService) CheckIsDbTypeSupportEnableBackup(dbType string) error {
	return nil
}
