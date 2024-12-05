package server

import "github.com/actiontech/sqle/sqle/model"

type backupHandler interface {
	backup() (backupResult string, err error)
}

type BackupStrategy string

const (
	BackupStrategyNone        BackupStrategy = "none"         // 不备份(不支持备份、无需备份、选择不备份)
	BackupStrategyReverseSql  BackupStrategy = "reverse_sql"  // 备份为反向SQL
	BackupStrategyOriginalRow BackupStrategy = "original_row" // 备份为原始行
	BackupStrategyManually    BackupStrategy = "manual"       // 标记为人工备份
	BackupRowsAffectedLimit   int            = 1000           // SQL影响行数上限，超过该上限的SQL不进行备份
)

type BackupStatus string

const (
	BackupStatusWaitingForExecution BackupStatus = "waiting_for_execution" // 等待备份
	BackupStatusExecuting           BackupStatus = "executing"             // 备份中
	BackupStatusFailed              BackupStatus = "failed"                // 备份失败
	BackupStatusSucceed             BackupStatus = "succeed"               // 备份成功
)

/* backupTaskMap mapping origin sql id to backup task */
type backupTaskMap map[uint]*model.BackupTask

func (m backupTaskMap) GetBackupStrategy(sqlId uint) string {
	if task, exist := m[sqlId]; exist {
		return task.BackupStrategy
	}
	return ""
}

func (m backupTaskMap) GetBackupStatus(sqlId uint) string {
	if task, exist := m[sqlId]; exist {
		return task.BackupStatus
	}
	return ""
}

func (m backupTaskMap) GetBackupResult(sqlId uint) string {
	if task, exist := m[sqlId]; exist {
		return task.BackupExecResult
	}
	return ""
}

func (m backupTaskMap) GetBackupStrategyTip(sqlId uint) string {
	if task, exist := m[sqlId]; exist {
		return task.BackupStrategyTip
	}
	return ""
}

func (m backupTaskMap) AddBackupTask(backupTask *model.BackupTask) {
	if backupTask == nil {
		return
	}
	if _, exist := m[backupTask.ExecuteSqlId]; !exist {
		m[backupTask.ExecuteSqlId] = backupTask
	}
}
