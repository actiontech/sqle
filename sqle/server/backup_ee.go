//go:build enterprise
// +build enterprise

package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

var ErrUnsupportedBackupInFileMode error = errors.New("enable backup in file mode is unsupported")

func initModelBackupTask(p driver.Plugin, task *model.Task, sql *model.ExecuteSQL) *model.BackupTask {
	var err error
	var strategyRes *driver.RecommendBackupStrategyRes
	// 为了兼容暂不支持backup接口但是支持rollback sql 接口的插件，这里将满足上述条件的插件的备份策略统一配置为生成反向SQL
	if driver.GetPluginManager().IsOptionalModuleEnabled(task.DBType, driverV2.OptionalBackup) {
		strategyRes, err = p.RecommendBackupStrategy(context.TODO(), sql.Content)
	} else if driver.GetPluginManager().IsOptionalModuleEnabled(task.DBType, driverV2.OptionalModuleGenRollbackSQL) {
		strategyRes = &driver.RecommendBackupStrategyRes{
			BackupStrategy:    driverV2.BackupStrategyReverseSql,
			BackupStrategyTip: fmt.Sprintf("数据源：%v,目前只支持根据反向SQL回滚", task.DBType),
		}
	}
	if err != nil {
		strategyRes.BackupStrategy = string(BackupStrategyManually)
		strategyRes.BackupStrategyTip = err.Error()
	}
	var schemaName string
	var tableName string
	if len(strategyRes.SchemasRefer) > 0 {
		schemaName = strategyRes.SchemasRefer[0]
	}
	if len(strategyRes.TablesRefer) > 0 {
		tableName = strategyRes.TablesRefer[0]
	}
	if strategyRes.BackupStrategy == "" {
		strategyRes.BackupStrategy = string(BackupStrategyManually)
		strategyRes.BackupStrategyTip = "暂不支持备份该SQL，请手工备份"
	}
	return &model.BackupTask{
		TaskId:            task.ID,
		InstanceId:        task.InstanceId,
		ExecuteSqlId:      sql.ID,
		BackupStatus:      string(BackupStatusWaitingForExecution),
		BackupStrategy:    strategyRes.BackupStrategy,
		BackupStrategyTip: strategyRes.BackupStrategyTip,
		SchemaName:        schemaName,
		TableName:         tableName,
	}
}

func getBackupManager(p driver.Plugin, sql *model.ExecuteSQL, dbType string, backupMaxRows uint64) (*BackupManager, error) {
	s := model.GetStorage()
	backupTask, err := s.GetBackupTaskByExecuteSqlId(sql.ID)
	if err != nil {
		return nil, err
	}
	return newBackupManager(p, backupTask, sql, dbType, backupMaxRows), nil
}

func newBackupManager(p driver.Plugin, modelBackupTask *model.BackupTask, sql *model.ExecuteSQL, dbType string, backupMaxRows uint64) *BackupManager {
	task := backupTask{
		ID:                modelBackupTask.ID,
		ExecTaskId:        sql.TaskId,
		ExecuteSqlId:      modelBackupTask.ExecuteSqlId,
		ExecuteSql:        sql.Content,
		SqlType:           sql.SQLType,
		BackupStatus:      BackupStatus(modelBackupTask.BackupStatus),
		InstanceId:        modelBackupTask.InstanceId,
		SchemaName:        modelBackupTask.SchemaName,
		TableName:         modelBackupTask.TableName,
		BackupStrategy:    BackupStrategy(modelBackupTask.BackupStrategy),
		BackupStrategyTip: modelBackupTask.BackupStrategyTip,
		BackupExecResult:  modelBackupTask.BackupExecResult,
		BackupMaxRows:     backupMaxRows,
	}
	var handler backupHandler
	switch modelBackupTask.BackupStrategy {
	case string(BackupStrategyManually):
		// 当用户选择手工备份时
		handler = &BackupManually{}
	case string(BackupStrategyOriginalRow):
		// 当用户选择备份行时
		handler = &BackupOriginalRow{baseBackupHandler: baseBackupHandler{plugin: p, task: task}}
	case string(BackupStrategyNone):
		// 当用户选择不备份时
		handler = &BackupNothing{}
	case string(BackupStrategyReverseSql):
		// 当用户不选择备份策略或选择了反向SQL
		handler = &BackupReverseSql{baseBackupHandler: baseBackupHandler{plugin: p, task: task}}
	default:
		handler = &BackupNothing{}
	}
	if !driver.GetPluginManager().IsOptionalModuleEnabled(dbType, driverV2.OptionalBackup) {
		if driver.GetPluginManager().IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleGenRollbackSQL) {
			handler = &BackupReverseSqlUseRollbackApi{baseBackupHandler: baseBackupHandler{plugin: p, task: task}}
		} else {
			handler = &BackupNothing{}
		}
	}
	return &BackupManager{
		backupTask:    task,
		backupHandler: handler,
	}
}

type backupHandler interface {
	backup() (backupResult string, err error)
}

type BackupManager struct {
	backupHandler
	backupTask
}

type backupTask struct {
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
	BackupMaxRows     uint64         // 备份的最大行数
}

func (backup backupTask) toModel() *model.BackupTask {
	return &model.BackupTask{
		Model: gorm.Model{
			ID: backup.ID,
		},
		TaskId:            backup.ExecTaskId,
		InstanceId:        backup.InstanceId,
		ExecuteSqlId:      backup.ExecuteSqlId,
		BackupStrategy:    string(backup.BackupStrategy),
		BackupStrategyTip: backup.BackupStrategyTip,
		BackupStatus:      string(backup.BackupStatus),
		BackupExecResult:  backup.BackupExecResult,
		SchemaName:        backup.SchemaName,
		TableName:         backup.TableName,
	}
}

func (mgr *BackupManager) Backup() (backupErr error) {
	s := model.GetStorage()
	var backupResult string
	defer func() {
		// update status to database according to backup error
		var status BackupStatus
		if backupErr != nil {
			status = BackupStatusFailed
		} else {
			status = BackupStatusSucceed
		}
		if updateStatusErr := mgr.UpdateStatusForBackupTaskTo(status); updateStatusErr != nil {
			backupErr = fmt.Errorf("in backup task %v, when UpdateStatusForBackupTaskTo %v failed %v %w", mgr.backupTask.ID, status, backupErr, updateStatusErr)
		}
		mgr.backupTask.BackupExecResult = backupResult
		updateTaskErr := s.UpdateBackupExecuteResult(mgr.backupTask.toModel())
		if updateTaskErr != nil {
			backupErr = fmt.Errorf("in backup task %v, when UpdateBackupExecuteResult failed %v %w", mgr.backupTask.ID, backupErr, updateTaskErr)
		}
	}()
	// update status in memory
	if err := mgr.UpdateStatusForBackupTaskTo(BackupStatusExecuting); err != nil {
		return err
	}

	// 执行备份操作
	backupResult, backupErr = mgr.backupHandler.backup()
	if backupErr != nil {
		return backupErr
	}

	return nil
}

/*
备份任务的备份状态机:

	[BackupStatusWaitingForExecution] --> [BackupStatusExecuting] --> [BackupStatusSucceed/BackupStatusFailed]
*/
func (task *BackupManager) UpdateStatusForBackupTaskTo(targetStatus BackupStatus) error {
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

type baseBackupHandler struct {
	plugin driver.Plugin
	task   backupTask
	svc    BackupService
}

func (backup *baseBackupHandler) backup() (backupResult string, backupErr error) {
	// generate reverse sql
	backupSqls, executeInfo, backupErr := backup.plugin.Backup(context.TODO(), string(BackupStrategyOriginalRow), backup.task.ExecuteSql, backup.task.BackupMaxRows)
	if backupErr != nil {
		return executeInfo, backupErr
	}
	if backupErr = backup.svc.saveBackupResultToRollbackSQLs(backup.task, backupSqls, executeInfo); backupErr != nil {
		return executeInfo, backupErr
	}
	return executeInfo, nil
}

type BackupNothing struct{}

func (BackupNothing) backup() (backupResult string, backupErr error) {
	return "", nil
}

type BackupManually struct{}

func (BackupManually) backup() (backupResult string, backupErr error) {
	return "", nil
}

type BackupOriginalRow struct {
	baseBackupHandler
}

type BackupReverseSql struct {
	baseBackupHandler
}

// 为了兼容暂不支持backup接口但是支持rollback sql 接口的插件，这里新增一个用于兼容的备份handler
type BackupReverseSqlUseRollbackApi struct {
	baseBackupHandler
}

func (backup *BackupReverseSqlUseRollbackApi) backup() (backupResult string, backupErr error) {
	// generate reverse sql
	backupSqlText, executeInfo, backupErr := backup.plugin.GenRollbackSQL(context.TODO(), backup.task.ExecuteSql)
	executeInfoInChinese := executeInfo.GetStrInLang(language.Chinese)
	if backupErr != nil {
		return executeInfoInChinese, backupErr
	}
	// adapter to backup sqls
	backupSqlNodes, backupErr := backup.plugin.Parse(context.TODO(), backupSqlText)
	if backupErr != nil {
		return executeInfoInChinese, backupErr
	}
	backupSqls := make([]string, 0, len(backupSqlNodes))
	for _, sql := range backupSqlNodes {
		if !strings.HasSuffix(sql.Text, ";") {
			sql.Text = sql.Text + ";"
		}
		backupSqls = append(backupSqls, sql.Text)
	}
	// save backup result into database
	if backupErr = backup.svc.saveBackupResultToRollbackSQLs(backup.task, backupSqls, executeInfoInChinese); backupErr != nil {
		return executeInfoInChinese, backupErr
	}
	return executeInfoInChinese, nil
}

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
	if driver.GetPluginManager().IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleGenRollbackSQL) {
		return nil
	}
	if !driver.GetPluginManager().IsOptionalModuleEnabled(dbType, driverV2.OptionalBackup) {
		return fmt.Errorf("db type %v can not enable backup", dbType)
	}
	return nil
}

func (svc BackupService) CheckCanTaskBackup(task *model.Task) bool {
	return task.EnableBackup && svc.CheckIsDbTypeSupportEnableBackup(task.DBType) == nil
}

func (BackupService) saveBackupResultToRollbackSQLs(task backupTask, backupSqls []string, executeInfo string) error {
	s := model.GetStorage()
	rollbackSqls := make([]*model.RollbackSQL, 0, len(backupSqls))
	for idx, rollbackSql := range backupSqls {
		rollbackSqls = append(rollbackSqls, &model.RollbackSQL{
			BaseSQL: model.BaseSQL{
				Number:      uint(idx) + 1,
				TaskId:      task.ExecTaskId,
				Content:     rollbackSql,
				Description: executeInfo,
			},
			ExecuteSQLId: task.ExecuteSqlId,
		})
	}
	// save backup result into database
	return s.UpdateRollbackSQLs(rollbackSqls)
}

type BackupSqlData struct {
	ExecOrder      uint     `json:"exec_order"`
	ExecSqlID      uint     `json:"exec_sql_id"`
	OriginSQL      string   `json:"origin_sql"`
	OriginTaskId   uint     `json:"origin_task_id"`
	BackupSqls     []string `json:"backup_sqls"`
	BackupStrategy string   `json:"backup_strategy" enum:"none,manual,reverse_sql,original_row"`
	BackupResult   string   `json:"backup_result"`
	BackupStatus   string   `json:"backup_status" enum:"waiting_for_execution,executing,succeed,failed"`
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
	type backupSqlInfo struct {
		backupSqls   []string
		backupResult string
		backupStatus string
	}
	// 3. get rollback sql
	backupSqlMap := make(map[uint]*backupSqlInfo)
	if len(originSqlIds) > 0 {
		backupInfos, err := s.GetBackupInfoFilterBy(map[string]interface{}{
			"filter_execute_sql_ids": originSqlIds,
		})
		if err != nil {
			return nil, 0, err
		}
		for _, backupInfo := range backupInfos {
			if _, exist := backupSqlMap[backupInfo.OriginSqlId]; !exist {
				backupSqlMap[backupInfo.OriginSqlId] = &backupSqlInfo{}
				backupSqlMap[backupInfo.OriginSqlId].backupResult = backupInfo.BackupExecResult
				backupSqlMap[backupInfo.OriginSqlId].backupStatus = backupInfo.BackupStatus
			}
			backupSqlMap[backupInfo.OriginSqlId].backupSqls = append(backupSqlMap[backupInfo.OriginSqlId].backupSqls, backupInfo.BackupSqls)
		}
	}

	// 4. fill sqls with instance and rollback sqls content
	backupSqlList = make([]*BackupSqlData, 0, len(sqlsOfWorkflow))
	for _, originSql := range sqlsOfWorkflow {
		var instanceName string
		if inst, exist := instanceMap[originSql.InstanceId]; exist {
			instanceName = inst.Name
		}
		backupSqlData := &BackupSqlData{
			ExecSqlID:      originSql.Id,
			OriginTaskId:   originSql.TaskId,
			ExecOrder:      originSql.ExecuteOrder,
			OriginSQL:      originSql.ExecuteSql,
			BackupStrategy: originSql.BackupStrategy,
			InstanceName:   instanceName,
			InstanceId:     originSql.InstanceId,
			ExecStatus:     originSql.ExecStatus,
			Description:    originSql.Description,
		}
		if sql, exist := backupSqlMap[originSql.Id]; exist {
			backupSqlData.BackupSqls = sql.backupSqls
			backupSqlData.BackupStatus = sql.backupStatus
			backupSqlData.BackupResult = sql.backupResult
		}
		backupSqlList = append(backupSqlList, backupSqlData)
	}
	return backupSqlList, count, nil
}

func (BackupService) CheckOriginWorkflowCanRollback(workflowId string) (*model.Workflow, error) {
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByWorkflowId(workflowId)
	if err != nil {
		log.Logger().Errorf("in CheckOriginWorkflowCanRollback when GetWorkflowByWorkflowId %v failed, error is %v", workflowId, err)
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("in CheckOriginWorkflowCanRollback when GetWorkflowByWorkflowId %v workflow does not exist", workflowId)
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
			log.Logger().Errorf("in CheckSqlsTasksMappingRelationship when get execute sql by task id %v failed, error is: %v", taskId, err)
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
			return nil, fmt.Errorf("in CheckSqlsTasksMappingRelationship sql id %v not exist in task %v", execSqlId, taskId)
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

func (BackupService) CanUpdateStrategyForTask(task *model.Task) error {
	if task.EnableBackup {
		return nil
	}
	return fmt.Errorf("can not update strategy for task which did not enable backup, task id %v", task.ID)
}

func (svc BackupService) UpdateBackupStrategyForSql(sqlId, backupStrategy string) error {
	if err := svc.checkStrategyIsSupported(BackupStrategy(backupStrategy)); err != nil {
		return err
	}
	return model.GetStorage().UpdateBackupStrategyForSql(sqlId, backupStrategy, "该备份策略由人工手动修改")
}

func (svc BackupService) BatchUpdateBackupStrategyForTask(taskId, backupStrategy string) error {
	if err := svc.checkStrategyIsSupported(BackupStrategy(backupStrategy)); err != nil {
		return err
	}
	return model.GetStorage().BatchUpdateBackupStrategyForTask(taskId, backupStrategy, "该备份策略由人工手动批量修改")
}

func (BackupService) checkStrategyIsSupported(backupStrategy BackupStrategy) error {
	switch backupStrategy {
	case BackupStrategyManually, BackupStrategyNone, BackupStrategyOriginalRow, BackupStrategyReverseSql:
	default:
		return fmt.Errorf("strategy %v is unsupported", backupStrategy)
	}
	return nil
}

// 如果数据源的备份开启，工单的备份关闭，则返回true
func (BackupService) IsBackupConflictWithInstance(taskEnableBackup, instanceEnableBackup bool) bool {
	return instanceEnableBackup && !taskEnableBackup
}

const DefaultBackupMaxRows uint64 = 1000

// AutoChooseBackupMaxRows 方法根据是否启用备份以及备份最大行数的设置来确定最终的备份最大行数。
func (BackupService) AutoChooseBackupMaxRows(enableBackup bool, backupMaxRows *uint64, instance model.Instance) uint64 {
	// 如果未启用备份，则返回 0。
	if !enableBackup {
		return 0
	}
	// 如果 backupMaxRows 不为 nil 并且其值大于 0，则返回 backupMaxRows 的值。
	if backupMaxRows != nil && *backupMaxRows >= 0 {
		return *backupMaxRows
	}
	// 如果实例启用了备份并且实例的备份最大行数大于等于 0，则返回实例的备份最大行数。
	if instance.EnableBackup && instance.BackupMaxRows >= 0 {
		return instance.BackupMaxRows
	}
	// 如果以上条件都不满足，则返回默认的备份最大行数 DefaultBackupMaxRows。
	return DefaultBackupMaxRows
}

func (BackupService) SupportedBackupStrategy(dbType string) []string {
	if driver.GetPluginManager().IsOptionalModuleEnabled(dbType, driverV2.OptionalBackup) {
		return []string{
			string(BackupStrategyManually),
			string(BackupStrategyNone),
			string(BackupStrategyOriginalRow),
			string(BackupStrategyReverseSql),
		}
	}
	if driver.GetPluginManager().IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleGenRollbackSQL) {
		return []string{
			string(BackupStrategyNone),
			string(BackupStrategyManually),
			string(BackupStrategyReverseSql),
		}
	}
	return []string{}
}

func modifyRulesWithBackupMaxRows(rules []*model.Rule, dbType string, backupMaxRows uint64) []*model.Rule {
	if backupMaxRows == 0 {
		return rules
	}
	svc := BackupService{}
	if err := svc.CheckIsDbTypeSupportEnableBackup(dbType); err != nil {
		return rules
	}
	switch dbType {
	case driverV2.DriverTypeTBase, driverV2.DriverTypePostgreSQL, "GaussDB for MySQL":
		return modifyRulesForPgLikeDriver(rules, backupMaxRows)
	case driverV2.DriverTypeOceanBase, "GoldenDB", "TiDB", driverV2.DriverTypeSQLServer, driverV2.DriverTypeTDSQLForInnoDB:
		return modifyRulesForMySQLLikeDriver(rules, backupMaxRows)
	case driverV2.DriverTypeOracle, "DM", "OceanBase For Oracle":
		return modifyRulesForOracleLikeDriver(rules, backupMaxRows)
	case driverV2.DriverTypeDB2:
		// 没有限制行数，找不到对应规则
	}
	return rules
}
func modifyRulesForOracleLikeDriver(rules []*model.Rule, backupMaxRows uint64) []*model.Rule {
	// [{"key": "first_key", "desc": "影响行数", "type": "string", "enums": null, "value": "1000", "i18n_desc": null}] Oracle_084
	var Oracle84 bool
	var Oracle85 bool
	for i := range rules {
		if Oracle84 && Oracle85 {
			break
		}
		if rules[i].Name == "Oracle_084" {
			rules[i].Params = params.Params{
				&params.Param{
					Key:   "first_key",
					Desc:  "影响行数",
					Value: fmt.Sprint(backupMaxRows),
					Type:  "int",
				},
			}
			Oracle84 = true
			continue
		}
		if rules[i].Name == "Oracle_085" {
			Oracle85 = true
			continue
		}
	}
	if !Oracle84 {
		rules = append(rules, &model.Rule{
			Name: "Oracle_084",
			Params: params.Params{
				&params.Param{
					Key:   "first_key",
					Desc:  "影响行数",
					Value: fmt.Sprint(backupMaxRows),
					Type:  "int",
				},
			},
			Level:         model.NoticeAuditLevel,
			DBType:        driverV2.DriverTypeOracle,
			HasAuditPower: true,
			I18nRuleInfo: driverV2.I18nRuleInfo{
				language.Tag{}: &driverV2.RuleInfo{
					Desc:       "在 DML 语句中预计影响行数超过指定值时不生成回滚语句",
					Annotation: "大事务回滚，容易影响数据库性能，使得业务发生波动；具体规则阈值可以根据业务需求调整，默认值：1000",
					Category:   "全局配置",
					Knowledge:  driverV2.RuleKnowledge{},
				},
			},
		})
	}
	if !Oracle85 {
		rules = append(rules, &model.Rule{
			Name:          "Oracle_085",
			Level:         model.NoticeAuditLevel,
			DBType:        driverV2.DriverTypeOracle,
			HasAuditPower: true,
			I18nRuleInfo: driverV2.I18nRuleInfo{
				language.Tag{}: &driverV2.RuleInfo{
					Desc:       "开启审核时生成回滚语句",
					Annotation: "回滚语句可以挽回错误的sql执行,提供容错机制",
					Category:   "全局配置",
					Knowledge:  driverV2.RuleKnowledge{},
				},
			},
		})
	}
	return rules
}

func modifyRulesForMySQLLikeDriver(rules []*model.Rule, backupMaxRows uint64) []*model.Rule {
	var ruleRollBackMaxRows bool
	for i := range rules {
		if rules[i].Name == "dml_rollback_max_rows" {
			rules[i].Params = params.Params{
				&params.Param{
					Key:   "first_key",
					Desc:  "最大影响行数",
					Value: fmt.Sprint(backupMaxRows),
					Type:  "int",
				},
			}
			ruleRollBackMaxRows = true
			break
		}
	}
	if !ruleRollBackMaxRows {
		rules = append(rules, &model.Rule{
			Name:   "dml_rollback_max_rows",
			DBType: driverV2.DriverTypeOceanBase,
			Level:  model.NoticeAuditLevel,
			Params: params.Params{
				&params.Param{
					Key:   "first_key",
					Desc:  "最大影响行数",
					Value: fmt.Sprint(backupMaxRows),
					Type:  "int",
				},
			},
			HasAuditPower: true,
			I18nRuleInfo: driverV2.I18nRuleInfo{
				language.Tag{}: &driverV2.RuleInfo{
					Desc:       "在 DML 语句中预计影响行数超过指定值则不回滚",
					Annotation: "大事务回滚，容易影响数据库性能，使得业务发生波动；具体规则阈值可以根据业务需求调整，默认值：1000",
					Category:   "全局配置",
					Knowledge:  driverV2.RuleKnowledge{},
				},
			},
		})
	}
	return rules
}

func modifyRulesForPgLikeDriver(rules []*model.Rule, backupMaxRows uint64) []*model.Rule {
	var rule24 bool
	var rule25 bool

	for i := range rules {
		if rule24 && rule25 {
			break
		}
		// [{"key": "max_affected_rows", "desc": "回滚语句影响行数", "type": "int", "enums": null, "value": "1000", "i18n_desc": null}]
		if !rule24 && rules[i].Name == "pg_024" {
			rules[i].Params = params.Params{
				&params.Param{
					Key:   "max_affected_rows",
					Desc:  "回滚语句影响行数",
					Value: fmt.Sprint(backupMaxRows),
					Type:  "int",
				},
			}
			rule24 = true
			continue
		}
		// 这条规则无param
		if !rule25 && rules[i].Name == "pg_025" {
			rule25 = true
			rules[i].HasAuditPower = true
			continue
		}
	}
	if !rule24 {
		rules = append(rules, &model.Rule{
			Name:   "pg_024",
			DBType: driverV2.DriverTypePostgreSQL,
			Level:  model.NoticeAuditLevel,
			Params: params.Params{
				&params.Param{
					Key:   "max_affected_rows",
					Desc:  "回滚语句影响行数",
					Value: fmt.Sprint(backupMaxRows),
					Type:  "int",
				},
			},
			HasAuditPower: true,
			I18nRuleInfo: driverV2.I18nRuleInfo{
				language.Tag{}: &driverV2.RuleInfo{
					Desc:       "在 DML 语句中预计影响行数超过指定值则不回滚",
					Annotation: "大事务回滚，容易影响数据库性能，使得业务发生波动；具体规则阈值可以根据业务需求调整，默认值：1000",
					Category:   "全局配置",
					Knowledge:  driverV2.RuleKnowledge{},
				},
			},
		})
	}
	if !rule25 {
		rules = append(rules, &model.Rule{
			Name:          "pg_025",
			DBType:        driverV2.DriverTypePostgreSQL,
			Level:         model.NoticeAuditLevel,
			HasAuditPower: true,
			I18nRuleInfo: driverV2.I18nRuleInfo{
				language.Tag{}: &driverV2.RuleInfo{
					Desc:       "使用sql语句回滚功能",
					Annotation: "回滚语句可以挽回错误的sql执行,提供容错机制",
					Category:   "全局配置",
					Knowledge:  driverV2.RuleKnowledge{},
				},
			},
		})
	}
	return rules
}
