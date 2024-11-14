//go:build enterprise
// +build enterprise

package server

import (
	"errors"
	"fmt"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
)

var ErrUnsupportedBackupInFileMode error = errors.New("enable backup in file mode is unsupported")

type BackupService struct{}

// 文件模式不支持备份，仅支持SQL模式上线
func (BackupService) CheckBackupConflictWithExecMode(EnableBackup bool, ExecMode string) error {
	if EnableBackup && ExecMode == executeSqlFileMode {
		return ErrUnsupportedBackupInFileMode
	}
	return nil
}

// 检查数据源类型是否支持备份
func (BackupService) CheckIsDbTypeSupportEnableBackup(dbType string) error {
	if dbType != driverV2.DriverTypeMySQL {
		return fmt.Errorf("db type %v can not enable backup", dbType)
	}
	return nil
}

// TODO 不同数据库的备份推荐可能不同，后续考虑将推荐备份策略的推荐放到插件中
func initModelBackupTask(task *model.Task, sql *model.ExecuteSQL) *model.BackupTask {
	var tableName string
	var schemaName string
	var strategy BackupStrategy = BackupStrategyReverseSql
	var reason string = "default backup strategy is reverse sql in mvp1"
	if sql.RowAffects > int64(BackupRowsAffectedLimit) {
		strategy = BackupStrategyManually
		reason = fmt.Sprintf("the rows affected by this sql, is bigger than limit:%v", BackupRowsAffectedLimit)
	}
	// TODO 根据SQL的类型来推荐备份策略
	if sql.SQLType == driverV2.SQLTypeDQL {
		strategy = BackupStrategyNone
		reason = fmt.Sprintf("the type of sql is %v, has no need to backup", sql.SQLType)
	}
	// TODO 根据备份SQL所引用的schema和table的数量推荐备份策略
	// if len(sql.TableReferred) == 1 {
	// 	tableName = sql.TableReferred[0]
	// } else {
	// 	strategy = BackupStrategyNone
	// 	reason = "unsupported one sql refer to multi-table"
	// }
	// if len(sql.SchemaReferred) == 1 {
	// 	schemaName = sql.SchemaReferred[0]
	// } else {
	// 	strategy = BackupStrategyNone
	// 	reason = "unsupported one sql refer to multi-schema"
	// }

	return &model.BackupTask{
		TaskId:            task.ID,
		InstanceId:        task.InstanceId,
		ExecuteSqlId:      sql.ID,
		BackupStrategy:    string(strategy),
		BackupStrategyTip: reason,
		BackupStatus:      string(BackupStatusWaitingForExecution),
		SchemaName:        schemaName,
		TableName:         tableName,
	}
}

type BaseBackupTask struct {
	ID           uint
	ExecTaskId   uint
	InstanceId   uint64
	SchemaName   string
	TableName    string
	ExecuteSqlId uint
	ExecuteSql   string // load from
	SqlType      string // ddl dml dql

	BackupStatus      BackupStatus
	BackupStrategy    BackupStrategy
	BackupStrategyTip string

	BackupExecInfo string
}

func (t BaseBackupTask) Backup() error {
	return nil
}

type BackupNothing struct {
	BaseBackupTask
}

type BackupOriginRow struct {
	BaseBackupTask
}

type BackupManually struct {
	BaseBackupTask
}

type BackupReverseSql struct {
	BaseBackupTask
}
