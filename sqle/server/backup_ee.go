//go:build enterprise
// +build enterprise

package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

var ErrUnsupportedBackupInFileMode error = errors.New("enable backup in file mode is unsupported")

type BackupService struct{}

func (BackupService) GetBackupTasksMap(taskId uint) (backupTaskMap, error) {
	backupTasks, err := model.GetStorage().GetBackupTaskByTaskId(taskId)
	if err != nil {
		return nil, err
	}
	backupTaskMap := make(backupTaskMap)
	for _, task := range backupTasks {
		backupTaskMap.AddBackupTask(task)
	}
	return backupTaskMap, nil
}

/* rollbackSqlMap mapping origin sql id to rollback sqls */
func (BackupService) GetRollbackSqlsMap(taskId uint) (map[uint][]string, error) {
	rollbackSqls, err := model.GetStorage().GetRollbackSqlByTaskId(taskId)
	if err != nil {
		return nil, err
	}
	rollbackSqlMap := make(map[uint][]string)
	for _, sql := range rollbackSqls {
		rollbackSqlMap[sql.ExecuteSQLId] = append(rollbackSqlMap[sql.ExecuteSQLId], sql.Content)
	}
	return rollbackSqlMap, nil
}

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
	// if sql.RowAffects > int64(BackupRowsAffectedLimit) {
	// 	strategy = BackupStrategyManually
	// 	reason = fmt.Sprintf("the rows affected by this sql, is bigger than limit:%v", BackupRowsAffectedLimit)
	// }
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
			plugin: a.plugin,
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
				BackupExecResult:  backupTask.BackupExecResult,
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
		BackupExecResult:  task.BackupExecResult,
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
	BackupExecResult  string         // 备份执行详情信息
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
	plugin driver.Plugin
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
	rollbackSQL, info, updateStatusErr := backup.plugin.GenRollbackSQL(context.TODO(), backup.ExecuteSql)
	if updateStatusErr != nil {
		return updateStatusErr
	}
	// set backup execute result
	backup.BackupExecResult = info.GetStrInLang(language.Chinese)
	if backup.BaseBackupTask.BackupExecResult == "" {
		backup.BaseBackupTask.BackupExecResult = string(BackupStatusSucceed)
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

type BackupSqlData struct {
	ExecOrder      uint     `json:"exec_order"`
	ExecSqlID      uint     `json:"exec_sql_id"`
	OriginSQL      string   `json:"origin_sql"`
	OriginTaskId   uint     `json:"origin_task_id"`
	BackupSqls     []string `json:"backup_sqls"`
	BackupStrategy string   `json:"backup_strategy" enum:"none,manual,reverse_sql,original_row"`
	InstanceName   string   `json:"instance_name"`
	InstanceId     uint64   `json:"instance_id "`
	ExecStatus     string   `json:"exec_status"`
	Description    string   `json:"description"`
}

func (BackupService) GetBackupSqlList(ctx context.Context, workflowId, filterInstanceId, filterExecStatus string, limit, offset uint32) (backupSqlList []*BackupSqlData, count uint64, err error) {
	s := model.GetStorage()
	// 1. get origin sql filter by filters and limit and offset
	data := map[string]interface{}{
		"filter_workflow_id": workflowId,
		"filter_exec_status": filterExecStatus,
		"filter_instance_id": filterInstanceId,
		"limit":              limit,
		"offset":             offset,
	}
	sqlsOfWorkflow, count, err := s.GetWorkflowSqlsByReq(data)
	if err != nil {
		return nil, 0, err
	}
	if len(sqlsOfWorkflow) == 0 {
		return []*BackupSqlData{}, 0, nil
	}

	// 2. get instance from dms
	instanceIds := []uint64{}
	instanceIdMap := make(map[uint64]struct{})
	originSqlIds := make([]uint, 0, len(sqlsOfWorkflow))
	for _, sql := range sqlsOfWorkflow {
		if _, exist := instanceIdMap[sql.InstanceId]; !exist {
			instanceIdMap[sql.InstanceId] = struct{}{}
			instanceIds = append(instanceIds, sql.InstanceId)
		}
		originSqlIds = append(originSqlIds, sql.Id)
	}
	instances, err := dms.GetInstancesByIds(ctx, instanceIds)
	if err != nil {
		return nil, 0, err
	}
	instanceMap := make(map[uint64]*model.Instance, len(instances))
	for _, instance := range instances {
		instanceMap[instance.ID] = instance
	}

	// 3. get rollback sql
	backupSqlMap := make(map[uint][]string)
	if len(originSqlIds) > 0 {
		backupSqls, err := s.GetRollbackSqlByFilters(map[string]interface{}{
			"filter_execute_sql_ids": originSqlIds,
		})
		if err != nil {
			return nil, 0, err
		}
		for _, backupSql := range backupSqls {
			backupSqlMap[backupSql.OriginSqlId] = append(backupSqlMap[backupSql.OriginSqlId], backupSql.BackupSqls)
		}
	}

	// 4. fill sqls with instance and rollback sqls content
	backupSqlList = make([]*BackupSqlData, 0, len(sqlsOfWorkflow))
	for _, originSql := range sqlsOfWorkflow {
		var instanceName string
		if inst, exist := instanceMap[originSql.InstanceId]; exist {
			instanceName = inst.Name
		}
		backupSqlList = append(backupSqlList, &BackupSqlData{
			ExecSqlID:      originSql.Id,
			OriginTaskId:   originSql.TaskId,
			ExecOrder:      originSql.ExecuteOrder,
			OriginSQL:      originSql.ExecuteSql,
			BackupSqls:     backupSqlMap[originSql.Id],
			BackupStrategy: originSql.BackupStrategy,
			InstanceName:   instanceName,
			InstanceId:     originSql.InstanceId,
			ExecStatus:     originSql.ExecStatus,
			Description:    originSql.Description,
		})
	}
	return backupSqlList, count, nil
}

func (BackupService) CheckOriginWorkflowCanRollback(workflowId string) (*model.Workflow, error) {
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByWorkflowId(workflowId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("workflow does not exist")
	}
	switch workflow.Record.Status {
	case model.WorkflowStatusCancel, model.WorkflowStatusExecFailed, model.WorkflowStatusReject, model.SQLAuditStatusFinished:
		return workflow, nil
	default:
		return nil, fmt.Errorf("can not create rollback workflow, because the status of origin workflow is %v", workflow.Record.Status)
	}
}

func (BackupService) CheckSqlsTasksMappingRelationship(taskIds, sqlIds []uint) (map[uint][]uint, error) {
	s := model.GetStorage()
	// check origin sql belongs to origin task
	originSqlMap := make(map[uint] /* sql id  */ uint /* task id */)
	for _, taskId := range taskIds {
		executeSQLs, err := s.GetExecuteSQLByTaskId(taskId)
		if err != nil {
			return nil, err
		}
		for _, sql := range executeSQLs {
			originSqlMap[sql.ID] = taskId
		}
	}
	originTaskSqlMap := make(map[uint] /* task id  */ []uint /* sql id */)
	for _, execSqlId := range sqlIds {
		if taskId, exist := originSqlMap[execSqlId]; exist {
			originTaskSqlMap[taskId] = append(originTaskSqlMap[taskId], execSqlId)
		} else {
			return nil, fmt.Errorf("sql id %v not exist in task %v", execSqlId, taskId)
		}
	}
	return originTaskSqlMap, nil
}

/*
备份服务更新原有工单:
 1. 建立原有工单、原有工单的任务以及原有工单的执行SQL和回滚工单的关联关系
 2. 更新原有工单的状态至执行失败
 3. 更新开启回滚工单的原有SQL对应的任务的任务状态至执行失败
 4. 更新开启回滚工单的SQL对应其原有工单的SQL状态为执行回滚
*/
func (svc BackupService) UpdateOriginWorkflow(rollbackWorkflowId, originWorkflowId, originWorkflowRecordId string, originTaskSqlsMap map[uint] /* task id  */ []uint /* sql id */) error {
	var sqlsNeedToRollback []uint
	var tasksNeedToRollback []uint
	for taskId, sqlId := range originTaskSqlsMap {
		sqlsNeedToRollback = append(sqlsNeedToRollback, sqlId...)
		tasksNeedToRollback = append(tasksNeedToRollback, taskId)
	}
	return model.GetStorage().Tx(
		func(txDB *gorm.DB) error {
			err := svc.AssociateRollbackWorkflowWithOriginTaskSqls(txDB, rollbackWorkflowId, originTaskSqlsMap)
			if err != nil {
				return err
			}
			err = svc.AssociateRollbackWorkflowWithOriginalWorkflow(txDB, rollbackWorkflowId, originWorkflowId)
			if err != nil {
				return err
			}
			err = svc.UpdateOriginSqlExecuteStatusToExecuteRollback(txDB, sqlsNeedToRollback)
			if err != nil {
				return err
			}
			err = svc.UpdateOriginTaskStatusToExecuteFailed(txDB, tasksNeedToRollback)
			if err != nil {
				return err
			}
			err = svc.UpdateOriginWorkflowStatusToExecuteFailed(txDB, originWorkflowRecordId)
			if err != nil {
				return err
			}
			return nil

		},
	)
}

func (BackupService) AssociateRollbackWorkflowWithOriginTaskSqls(txDB *gorm.DB, rollbackWorkflowId string, originTaskSqlMap map[uint][]uint) error {
	var relations []model.ExecuteSqlRollbackWorkflows
	for taskId, rollbackSqlIds := range originTaskSqlMap {
		for _, sqlId := range rollbackSqlIds {
			relations = append(relations, model.ExecuteSqlRollbackWorkflows{
				TaskId:             taskId,
				ExecuteSqlId:       sqlId,
				RollbackWorkflowId: rollbackWorkflowId,
			})
		}
	}
	err := model.CreateExecuteSqlRollbackWorkflowRelation(txDB, relations)
	if err != nil {
		return err
	}
	return nil
}

func (BackupService) AssociateRollbackWorkflowWithOriginalWorkflow(txDB *gorm.DB, rollbackWorkflowId, originWorkflowId string) error {
	err := model.CreateRollbackWorkflowOriginalWorkflowRelation(txDB, &model.RollbackWorkflowOriginalWorkflows{
		RollbackWorkflowId: rollbackWorkflowId,
		OriginalWorkflowId: originWorkflowId,
	})
	if err != nil {
		return err
	}
	return nil
}

// 更新原始工单的状态为执行失败
func (BackupService) UpdateOriginWorkflowStatusToExecuteFailed(txDB *gorm.DB, originWorkflowId string) error {
	return model.UpdateWorkflowByWorkflowId(txDB, originWorkflowId, map[string]interface{}{
		"status": model.WorkflowStatusExecFailed,
	})
}

// 在原始工单执行失败的SQL处，可以看到已回滚/正在回滚的标记
func (BackupService) UpdateOriginSqlExecuteStatusToExecuteRollback(txDB *gorm.DB, sqlsNeedToRollback []uint) error {
	return model.BatchUpdateExecuteSqlExecuteStatus(txDB, sqlsNeedToRollback, model.SQLExecuteStatusExecuteRollback)
}

// 在原始工单执行失败的SQL对应的task的状态变更为失败
func (BackupService) UpdateOriginTaskStatusToExecuteFailed(txDB *gorm.DB, tasksNeedToRollback []uint) error {
	return model.UpdateTaskStatusByIDsTx(txDB, tasksNeedToRollback, map[string]interface{}{
		"status": model.TaskStatusExecuteFailed,
	})
}
