//go:build !enterprise
// +build !enterprise

package server

import (
	"github.com/actiontech/sqle/sqle/driver"
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

func initModelBackupTask(p driver.Plugin, task *model.Task, sql *model.ExecuteSQL) *model.BackupTask {
	return &model.BackupTask{}
}

func toBackupTask(p driver.Plugin, sql *model.ExecuteSQL) (*BaseBackupTask, error) {
	return &BaseBackupTask{}, nil
}

func (BackupService) GetRollbackSqlsMap(taskId uint) (map[uint][]string, error) {
	return make(map[uint][]string), nil
}

func (BackupService) GetBackupTasksMap(taskId uint) (backupTaskMap, error) {
	return make(backupTaskMap), nil
}

func (BackupService) IsBackupConflictWithInstance(taskEnableBackup, instanceEnableBackup bool) bool {
	return false
}

func (BackupService) CheckCanTaskBackup(task *model.Task) bool {
	return false
}


func (BackupService) SupportedBackupStrategy(dbType string) []string {
	return []string{}
}