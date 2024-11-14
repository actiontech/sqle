//go:build !enterprise
// +build !enterprise

package server

import (
	"github.com/actiontech/sqle/sqle/model"
)

type BackupService struct{}

func (BackupService) CheckBackupConflictWithExecMode(EnableBackup bool, ExecMode string) error {
	return nil
}

func (BackupService) CheckIsDbTypeSupportEnableBackup(dbType string) error {
	return nil
}

type BaseBackupTask struct{}

func (t BaseBackupTask) Backup() error {
	return nil
}

func initModelBackupTask(task *model.Task, sql *model.ExecuteSQL) *model.BackupTask {
	return &model.BackupTask{}
}
