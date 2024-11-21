//go:build enterprise
// +build enterprise

package server

import (
	"context"
	"errors"
	"fmt"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"golang.org/x/text/language"
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

func toBackupTask(a *action, sql *model.ExecuteSQL) (BackupTask, error) {
	s := model.GetStorage()
	backupTask, err := s.GetBackupTaskByExecuteSqlId(sql.ID)
	if err != nil {
		return nil, err
	}
	switch backupTask.BackupStrategy {
	case string(BackupStrategyManually):
		// 当用户选择手工备份时
		return &BackupManually{}, nil
	case string(BackupStrategyOriginalRow):
		// 当用户选择备份行时
		return &BackupOriginalRow{}, nil
	case string(BackupStrategyNone):
		// 当用户选择不备份时
		return &BackupNothing{}, nil
	case string(BackupStrategyReverseSql):
		// 当用户不选择备份策略或选择了反向SQL
		return &BackupReverseSql{
			action: a,
			BaseBackupTask: BaseBackupTask{
				ID:                backupTask.ID,
				ExecTaskId:        sql.TaskId,
				ExecuteSqlId:      backupTask.ExecuteSqlId,
				ExecuteSql:        sql.Content,
				SqlType:           sql.SQLType,
				BackupStatus:      BackupStatus(backupTask.BackupStatus),
				InstanceId:        backupTask.InstanceId,
				SchemaName:        backupTask.SchemaName,
				TableName:         backupTask.TableName,
				BackupStrategy:    BackupStrategy(backupTask.BackupStrategy),
				BackupStrategyTip: backupTask.BackupStrategyTip,
				BackupExecInfo:    backupTask.BackupExecInfo,
			},
		}, nil
	default:
		return &BackupNothing{}, nil
	}
}

func (task BaseBackupTask) toModel() *model.BackupTask {
	return &model.BackupTask{
		TaskId:            task.ExecTaskId,
		InstanceId:        task.InstanceId,
		ExecuteSqlId:      task.ExecuteSqlId,
		BackupStrategy:    string(task.BackupStrategy),
		BackupStrategyTip: task.BackupStrategyTip,
		BackupStatus:      string(task.BackupStatus),
		BackupExecInfo:    task.BackupExecInfo,
		SchemaName:        task.SchemaName,
		TableName:         task.TableName,
	}
}

type BaseBackupTask struct {
	ID         uint   // 备份任务id
	ExecTaskId uint   // 备份任务对应的执行任务id
	InstanceId uint64 // 备份任务对应的数据源id

	ExecuteSqlId uint   // 备份的原始SQL的id
	ExecuteSql   string // 备份的原始SQL
	SchemaName   string // 备份的原始SQL对应的schema
	TableName    string // 备份的原始SQL对应的table
	SqlType      string // 备份的原始SQL类型 ddl dml dql

	BackupStrategy    BackupStrategy // 备份策略
	BackupStrategyTip string         // 备份策略推荐原因
	BackupStatus      BackupStatus   // 备份执行状态
	BackupExecInfo    string         // 备份执行详情信息
}

func (t BaseBackupTask) Backup() error {
	return nil
}

/*
备份任务的备份状态机:

	[BackupStatusWaitingForExecution] --> [BackupStatusExecuting] --> [BackupStatusSucceed/BackupStatusFailed]
*/
func (task *BaseBackupTask) UpdateStatusTo(targetStatus BackupStatus) error {
	// 定义状态流转规则
	validTransitions := map[BackupStatus][]BackupStatus{
		BackupStatusWaitingForExecution: {BackupStatusExecuting},
		BackupStatusExecuting:           {BackupStatusSucceed, BackupStatusFailed},
	}

	// 检查目标状态是否是允许的流转状态
	allowedStatuses, ok := validTransitions[task.BackupStatus]
	if !ok {
		return fmt.Errorf("current status %s does not allow any transitions", task.BackupStatus)
	}

	for _, status := range allowedStatuses {
		if status == targetStatus {
			task.BackupStatus = targetStatus
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", task.BackupStatus, targetStatus)
}

type BackupNothing struct {
	BaseBackupTask
}

type BackupOriginalRow struct {
	BaseBackupTask
}

type BackupManually struct {
	BaseBackupTask
}

type BackupReverseSql struct {
	BaseBackupTask
	action *action
}

// TODO 不同数据库的备份方式可能不同,备份动作，应该放到插件里面
func (backup *BackupReverseSql) Backup() (backupErr error) {
	s := model.GetStorage()
	var modelBackupTask *model.BackupTask = backup.toModel()
	defer func() {
		// update status to database according to backup error
		var status BackupStatus
		if backupErr != nil {
			status = BackupStatusFailed
		} else {
			status = BackupStatusSucceed
		}
		if updateStatusErr := backup.UpdateStatusTo(status); updateStatusErr != nil {
			backupErr = fmt.Errorf("%v%w", backupErr, updateStatusErr)
		}

		updateTaskErr := s.UpdateBackupExecuteResult(backup.toModel())
		if updateTaskErr != nil {
			backupErr = fmt.Errorf("%v%w", backupErr, updateTaskErr)
		}
	}()

	// update status in memory
	if err := backup.UpdateStatusTo(BackupStatusExecuting); err != nil {
		return err
	}
	// generate reverse sql
	rollbackSQL, info, updateStatusErr := backup.action.plugin.GenRollbackSQL(context.TODO(), backup.ExecuteSql)
	if updateStatusErr != nil {
		return updateStatusErr
	}
	// set backup info
	backup.BaseBackupTask.BackupExecInfo = info.GetStrInLang(language.Chinese)
	if backup.BaseBackupTask.BackupExecInfo == "" {
		backup.BaseBackupTask.BackupExecInfo = string(BackupStatusSucceed)
	}
	// save backup result into database
	updateStatusErr = s.UpdateRollbackSQLs([]*model.RollbackSQL{
		{
			BaseSQL: model.BaseSQL{
				TaskId:  modelBackupTask.TaskId,
				Content: rollbackSQL,
			},
			ExecuteSQLId: modelBackupTask.ExecuteSqlId,
		},
	})
	if updateStatusErr != nil {
		return updateStatusErr
	}
	return nil
}
