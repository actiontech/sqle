package model

import (
	"gorm.io/gorm"
)

func init() {
	autoMigrateList = append(autoMigrateList, &BackupTask{})
	autoMigrateList = append(autoMigrateList, &ExecuteSqlRollbackWorkflows{})
	autoMigrateList = append(autoMigrateList, &RollbackWorkflowOriginalWorkflows{})
}

type BackupTask struct {
	gorm.Model
	TaskId            uint   `gorm:"index;column:task_id;not null"`                                                                   // 备份任务关联原始SQL的任务的id
	InstanceId        uint64 `gorm:"column:instance_id;not null"`                                                                     // 备份任务关联原始SQL对应的数据源
	ExecuteSqlId      uint   `gorm:"index;column:execute_sql_id;not null"`                                                            // 备份任务关联原始SQL的id
	BackupStrategy    string `gorm:"type:enum('none','reverse_sql','original_row','manual');column:backup_strategy;size:20;not null"` // 备份任务支持的备份策略
	BackupStrategyTip string `gorm:"column:backup_strategy_tip;size:255;not null;default:''"`                                         // 推荐备份任务的原因
	BackupStatus      string `gorm:"type:enum('waiting_for_execution','executing','failed','succeed');column:backup_status;size:20"`  // 备份任务的执行状态
	BackupExecResult  string `gorm:"column:backup_exec_result;size:255;not null;default:''"`                                          // 备份任务的执行结果
	SchemaName        string `gorm:"column:schema_name;size:50;not null;default:''"`                                                  // 备份的SQL对应的schema
	TableName         string `gorm:"column:table_name;size:50;not null;default:''"`                                                   // 备份的SQL对应的table
}

type ExecuteSqlRollbackWorkflows struct {
	gorm.Model
	TaskId             uint   `gorm:"column:task_id;index;not null"`        // 执行的SQL所属的执行任务id
	ExecuteSqlId       uint   `gorm:"column:execute_sql_id;not null"`       // 执行的SQL的id
	RollbackWorkflowId string `gorm:"column:rollback_workflow_id;not null"` // 回滚工单的id
}

type RollbackWorkflowOriginalWorkflows struct {
	gorm.Model
	OriginalWorkflowId string `gorm:"column:original_workflow_id;not null"` // 原始工单的id
	RollbackWorkflowId string `gorm:"column:rollback_workflow_id;not null"` // 回滚工单的id
}

type ExecuteSqlRollbackWorkflowsRelation struct {
	ExecuteSqlRollbackWorkflows
	RollbackWorkflowSubject string `gorm:"column:subject" json:"subject"`
	RollbackWorkflowStatus  string `gorm:"column:status"  json:"status"`
}

type RollbackWorkflowOriginalWorkflowsRelation struct {
	RollbackWorkflowOriginalWorkflows
	RollbackWorkflowSubject string `gorm:"column:subject" json:"subject"`
	RollbackWorkflowStatus  string `gorm:"column:status"  json:"status"`
}
