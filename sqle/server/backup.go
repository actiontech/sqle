package server

type BackupTask interface {
	Backup() error
}

type BackupStrategy string

const (
	BackupStrategyNone       BackupStrategy = "none"        // 不备份(不支持备份、无需备份、选择不备份)
	BackupStrategyReverseSql BackupStrategy = "reverse_sql" // 备份为反向SQL
	BackupStrategyOriginRow  BackupStrategy = "origin_row"  // 备份为原始行
	BackupStrategyManually   BackupStrategy = "manually"    // 标记为人工备份
	BackupRowsAffectedLimit  int            = 1000          // SQL影响行数上限，超过该上限的SQL不进行备份
)

type BackupStatus string

const (
	BackupStatusWaitingForExecution BackupStatus = "waiting_for_execution" // 等待备份
	BackupStatusExecuting           BackupStatus = "executing"             // 备份中
	BackupStatusFailed              BackupStatus = "failed"                // 备份失败
	BackupStatusSucceed             BackupStatus = "succeed"               // 备份成功
)