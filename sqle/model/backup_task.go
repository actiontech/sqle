package model

import (
	"gorm.io/gorm"
)

func init() {
	autoMigrateList = append(autoMigrateList, &BackupTask{})
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
