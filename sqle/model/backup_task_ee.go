//go:build enterprise
// +build enterprise

package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
)

/*
the default batch size is 100, if the length of input backup task slice is shorter than 100, the batch will set as the length of slice
*/
func (s *Storage) BatchCreateBackupTasks(backupTasks []*BackupTask) error {
	var batchSize = 100
	if len(backupTasks) < 100 {
		batchSize = len(backupTasks)
	}
	return s.db.CreateInBatches(backupTasks, batchSize).Error
}

/*
update backup status and backup execute information
*/
func (s *Storage) UpdateBackupExecuteResult(task *BackupTask) error {
	return s.db.Model(&BackupTask{}).
		Where("id = ?", task.ID).
		UpdateColumns(map[string]interface{}{
			"backup_status":      task.BackupStatus,
			"backup_exec_result": task.BackupExecResult,
		}).Error
}

func (s *Storage) GetBackupTaskByExecuteSqlId(executeSqlId uint) (*BackupTask, error) {
	var backupTask BackupTask
	err := s.db.Model(&BackupTask{}).Where("execute_sql_id = ?", executeSqlId).First(&backupTask).Error
	if err != nil {
		return nil, err
	}
	return &backupTask, nil
}

func (s *Storage) GetBackupTaskByTaskId(taskId uint) ([]*BackupTask, error) {
	var backupTasks []*BackupTask
	err := s.db.Model(&BackupTask{}).Where("task_id = ?", taskId).Find(&backupTasks).Error
	if err != nil {
		return nil, err
	}
	return backupTasks, nil
}

func (s *Storage) GetRollbackSqlByTaskId(taskId uint) ([]*RollbackSQL, error) {
	var rollbackSqls []*RollbackSQL
	err := s.db.Model(&RollbackSQL{}).Where("task_id = ?", taskId).Find(&rollbackSqls).Error
	if err != nil {
		return nil, err
	}
	return rollbackSqls, nil
}

type WorkflowOriginSql struct {
	Id             uint   `json:"id"`
	TaskId         uint   `json:"task_id"`
	Description    string `json:"description"`
	ExecuteOrder   uint   `json:"execute_order"`
	ExecuteSql     string `json:"execute_sql"`
	ExecStatus     string `json:"exec_status"`
	InstanceId     uint64 `json:"instance_id"`
	BackupStrategy string `json:"backup_strategy"`
}

func (s *Storage) GetWorkflowSqlsByReq(data map[string]interface{}) (
	result []*WorkflowOriginSql, count uint64, err error) {

	err = s.getListResult(workflowOriginSqlBodyTpl, workflowOriginSqlQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(workflowOriginSqlBodyTpl, workflowOriginSqlCountTpl, data)
	return result, count, err
}

var workflowOriginSqlQueryTpl = `
SELECT 
	e_sql.id,
	e_sql.task_id,
	e_sql.description,
	e_sql.number AS execute_order,
	e_sql.content AS execute_sql,
	e_sql.exec_status,
	w_instance.instance_id,
	IFNULL(b_task.backup_strategy,'none') AS backup_strategy

{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var workflowOriginSqlCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var workflowOriginSqlBodyTpl = `
{{ define "body" }}
FROM execute_sql_detail AS e_sql 
LEFT JOIN backup_tasks AS b_task ON b_task.execute_sql_id = e_sql.id
LEFT JOIN workflow_instance_records AS w_instance ON e_sql.task_id = w_instance.task_id
LEFT JOIN workflows ON w_instance.workflow_record_id = workflows.workflow_record_id

WHERE workflows.workflow_id = :filter_workflow_id

{{- if .filter_exec_status }}
AND e_sql.exec_status = :filter_exec_status
{{- end }}

{{- if .filter_instance_id }}
AND w_instance.instance_id = :filter_instance_id
{{- end }}

{{ end }}
`

type BackupSql struct {
	OriginSqlId uint   `json:"origin_sql_id"`
	BackupSqls  string `json:"backup_sql"`
}

func (s *Storage) GetRollbackSqlByFilters(data map[string]interface{}) (list []*BackupSql, err error) {
	err = s.getListResult(reverseBackupSqlBodyTpl, reverseBackupSqlListQueryTpl, data, &list)
	if err != nil {
		return nil, err
	}
	return
}

var reverseBackupSqlListQueryTpl = `
SELECT 
	backup_tasks.execute_sql_id AS origin_sql_id, 
	IFNULL(rollback_sql_detail.content,'') AS backup_sql

{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var reverseBackupSqlListCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var reverseBackupSqlBodyTpl = `
{{ define "body" }}

FROM backup_tasks 
LEFT JOIN rollback_sql_detail ON backup_tasks.execute_sql_id = rollback_sql_detail.execute_sql_id
WHERE backup_tasks.execute_sql_id IN ( {{range $index, $element := .filter_execute_sql_ids}}{{if $index}},{{end}}{{$element}}{{end}} )
AND backup_strategy = "reverse_sql"

{{ end }}
`


type ExecuteSqlRollbackWorkflows struct {
	gorm.Model
	TaskId             uint   `gorm:"column:task_id;index;not null"`
	ExecuteSqlId       uint   `gorm:"column:execute_sql_id;not null"`
	RollbackWorkflowId string `gorm:"column:rollback_workflow_id;not null"`
}

type RollbackWorkflowOriginalWorkflows struct {
	gorm.Model
	OriginalWorkflowId string `gorm:"column:original_workflow_id;not null"`
	RollbackWorkflowId string `gorm:"column:rollback_workflow_id;not null"`
}

func CreateRollbackWorkflowOriginalWorkflowRelation(txDB *gorm.DB, relation *RollbackWorkflowOriginalWorkflows) error {
	return txDB.Model(&RollbackWorkflowOriginalWorkflows{}).Create(&relation).Error
}

func CreateExecuteSqlRollbackWorkflowRelation(txDB *gorm.DB, relations []ExecuteSqlRollbackWorkflows) error {
	return txDB.Model(&ExecuteSqlRollbackWorkflows{}).Create(&relations).Error
}

func UpdateWorkflowByWorkflowId(txDB *gorm.DB, workflowRecordId string, workflowParam map[string]interface{}) error {
	err := txDB.Model(&WorkflowRecord{}).
		Where("id = ?", workflowRecordId).
		Updates(workflowParam).Error
	if err != nil {
		return err
	}
	return nil
}

func BatchUpdateExecuteSqlExecuteStatus(txDB *gorm.DB, executeSqlIds []uint, status string) error {
	return txDB.Model(&ExecuteSQL{}).Where("id IN ?", executeSqlIds).Update("exec_status", status).Error
}

func UpdateTaskStatusByIDsTx(txDB *gorm.DB, taskIDs []uint, attrs map[string]interface{}) error {
	err := txDB.Model(&Task{}).Where("id IN (?)", taskIDs).Updates(attrs).Error
	return errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetExecuteSQLByTaskId(taskId uint) ([]*ExecuteSQL, error) {
	var executeSqls []*ExecuteSQL
	err := s.db.Model(&ExecuteSQL{}).Where("task_id = ?", taskId).Find(&executeSqls).Error
	if err != nil {
		return nil, err
	}
	return executeSqls, nil
}
