package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/jinzhu/gorm"
)

const (
	SQLTypeDML = "dml"
	SQLTypeDDL = "ddl"
)

const (
	TaskStatusInit             = "initialized"
	TaskStatusAudited          = "audited"
	TaskStatusExecuting        = "executing"
	TaskStatusExecuteSucceeded = "exec_succeeded"
	TaskStatusExecuteFailed    = "exec_failed"
)

const (
	TaskSQLSourceFromFormData       = "form_data"
	TaskSQLSourceFromSQLFile        = "sql_file"
	TaskSQLSourceFromMyBatisXMLFile = "mybatis_xml_file"
	TaskSQLSourceFromAuditPlan      = "audit_plan"
)

const TaskExecResultOK = "OK"

type Task struct {
	Model
	InstanceId   uint    `json:"instance_id"`
	Schema       string  `json:"instance_schema" gorm:"column:instance_schema" example:"db1"`
	PassRate     float64 `json:"pass_rate"`
	SQLSource    string  `json:"sql_source" gorm:"column:sql_source"`
	DBType       string  `json:"db_type" gorm:"default:'mysql'" example:"mysql"`
	Status       string  `json:"status" gorm:"default:\"initialized\""`
	CreateUserId uint

	CreateUser   *User          `gorm:"foreignkey:CreateUserId"`
	Instance     *Instance      `json:"-" gorm:"foreignkey:InstanceId"`
	ExecuteSQLs  []*ExecuteSQL  `json:"-" gorm:"foreignkey:TaskId"`
	RollbackSQLs []*RollbackSQL `json:"-" gorm:"foreignkey:TaskId"`
}

func (t *Task) InstanceName() string {
	if t.Instance != nil {
		return t.Instance.Name
	}
	return ""
}

const (
	SQLAuditStatusInitialized = "initialized"
	SQLAuditStatusDoing       = "doing"
	SQLAuditStatusFinished    = "finished"
)

const (
	SQLExecuteStatusInitialized = "initialized"
	SQLExecuteStatusDoing       = "doing"
	SQLExecuteStatusFailed      = "failed"
	SQLExecuteStatusSucceeded   = "succeeded"
)

type BaseSQL struct {
	Model
	TaskId uint `json:"-" gorm:"index"`
	Number uint `json:"number"`

	// Content store single SQL or batch SQLs
	//
	// Content may store batch SQLs When BaseSQL embed to RollbackSQL.
	// Split Content to single SQL before execute RollbackSQL.
	Content         string `json:"sql" gorm:"type:text"`
	StartBinlogFile string `json:"start_binlog_file"`
	StartBinlogPos  int64  `json:"start_binlog_pos"`
	EndBinlogFile   string `json:"end_binlog_file"`
	EndBinlogPos    int64  `json:"end_binlog_pos"`
	RowAffects      int64  `json:"row_affects"`
	ExecStatus      string `json:"exec_status" gorm:"default:\"initialized\""`
	ExecResult      string `json:"exec_result"`
}

func (s *BaseSQL) GetExecStatusDesc() string {
	switch s.ExecStatus {
	case SQLExecuteStatusInitialized:
		return "准备执行"
	case SQLExecuteStatusDoing:
		return "正在执行"
	case SQLExecuteStatusFailed:
		return "执行失败"
	case SQLExecuteStatusSucceeded:
		return "执行成功"
	default:
		return "未知"
	}
}

type ExecuteSQL struct {
	BaseSQL
	AuditStatus string `json:"audit_status" gorm:"default:\"initialized\""`
	AuditResult string `json:"audit_result" gorm:"type:text"`
	// AuditFingerprint generate from SQL and SQL audit result use MD5 hash algorithm,
	// it used for deduplication in one audit task.
	AuditFingerprint string `json:"audit_fingerprint" gorm:"index;type:char(32)"`
	// AuditLevel has four level: error, warn, notice, normal.
	AuditLevel string `json:"audit_level"`
}

func (s ExecuteSQL) TableName() string {
	return "execute_sql_detail"
}

func (s *ExecuteSQL) GetAuditStatusDesc() string {
	switch s.AuditStatus {
	case SQLAuditStatusInitialized:
		return "未审核"
	case SQLAuditStatusDoing:
		return "正在审核"
	case SQLAuditStatusFinished:
		return "审核完成"
	default:
		return "未知状态"
	}
}

func (s *ExecuteSQL) GetAuditResultDesc() string {
	if s.AuditResult == "" {
		return "审核通过"
	}
	return s.AuditResult
}

type RollbackSQL struct {
	BaseSQL
	ExecuteSQLId uint `gorm:"index;column:execute_sql_id"`
}

func (s RollbackSQL) TableName() string {
	return "rollback_sql_detail"
}

func (t *Task) HasDoingAudit() bool {
	if t.ExecuteSQLs != nil {
		for _, commitSQL := range t.ExecuteSQLs {
			if commitSQL.AuditStatus != SQLAuditStatusInitialized {
				return true
			}
		}
	}
	return false
}

func (t *Task) HasDoingExecute() bool {
	if t.ExecuteSQLs != nil {
		for _, commitSQL := range t.ExecuteSQLs {
			if commitSQL.ExecStatus != SQLExecuteStatusInitialized {
				return true
			}
		}
	}
	return false
}

func (t *Task) IsExecuteFailed() bool {
	if t.ExecuteSQLs != nil {
		for _, commitSQL := range t.ExecuteSQLs {
			if commitSQL.ExecStatus == SQLExecuteStatusFailed {
				return true
			}
		}
	}
	return false
}

func (t *Task) HasDoingRollback() bool {
	if t.RollbackSQLs != nil {
		for _, rollbackSQL := range t.RollbackSQLs {
			if rollbackSQL.ExecStatus != SQLExecuteStatusInitialized {
				return true
			}
		}
	}
	return false
}

func (s *Storage) GetTaskById(taskId string) (*Task, bool, error) {
	task := &Task{}
	err := s.db.Where("id = ?", taskId).Preload("Instance").First(task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return task, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetTaskDetailById(taskId string) (*Task, bool, error) {
	task := &Task{}
	err := s.db.Where("id = ?", taskId).Preload("Instance").
		Preload("ExecuteSQLs").Preload("RollbackSQLs").First(task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return task, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetTaskExecuteSQLContent(taskId string) ([]byte, error) {
	rows, err := s.db.Model(&ExecuteSQL{}).Select("content").
		Where("task_id = ?", taskId).Rows()
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	defer rows.Close()
	buff := &bytes.Buffer{}
	for rows.Next() {
		var content string
		rows.Scan(&content)
		buff.WriteString(strings.TrimRight(content, ";"))
		buff.WriteString(";\n")
	}
	return buff.Bytes(), nil
}

func (s *Storage) UpdateTask(task *Task, attrs ...interface{}) error {
	err := s.db.Table("tasks").Where("id = ?", task.ID).Update(attrs...).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateExecuteSQLs(ExecuteSQLs []*ExecuteSQL) error {
	tx := s.db.Begin()
	for _, executeSQL := range ExecuteSQLs {
		currentSql := executeSQL
		if err := tx.Save(currentSql).Error; err != nil {
			tx.Rollback()
			return errors.New(errors.ConnectStorageError, err)
		}
	}
	return errors.New(errors.ConnectStorageError, tx.Commit().Error)
}

func (s *Storage) UpdateRollbackSQLs(rollbackSQLs []*RollbackSQL) error {
	tx := s.db.Begin()
	for _, rollbackSQL := range rollbackSQLs {
		currentSql := rollbackSQL
		if err := tx.Save(currentSql).Error; err != nil {
			tx.Rollback()
			return errors.New(errors.ConnectStorageError, err)
		}
	}
	return errors.New(errors.ConnectStorageError, tx.Commit().Error)
}

func (s *Storage) UpdateTaskStatusById(taskId uint, status string) error {
	err := s.db.Model(&Task{}).Where("id = ?", taskId).Update(map[string]string{
		"status": status,
	}).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateExecuteSQLStatusByTaskId(task *Task, status string) error {
	query := "UPDATE execute_sql_detail SET exec_status=? WHERE task_id=?"

	tx := s.db.Begin()
	if err := tx.Exec(query, status, task.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *Storage) UpdateExecuteSqlStatus(baseSQL *BaseSQL, status, result string) error {
	attr := map[string]interface{}{}
	if status != "" {
		baseSQL.ExecStatus = status
		attr["exec_status"] = status
	}
	if result != "" {
		baseSQL.ExecResult = result
		attr["exec_result"] = result
	}
	return s.UpdateExecuteSQLById(fmt.Sprintf("%v", baseSQL.ID), attr)
}

func (s *Storage) UpdateExecuteSQLById(executeSQLId string, attrs ...interface{}) error {
	err := s.db.Table(ExecuteSQL{}.TableName()).Where("id = ?", executeSQLId).Update(attrs...).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateRollbackSqlStatus(baseSQL *BaseSQL, status, result string) error {
	attr := map[string]interface{}{}
	if status != "" {
		baseSQL.ExecStatus = status
		attr["exec_status"] = status
	}
	if result != "" {
		baseSQL.ExecResult = result
		attr["exec_result"] = result
	}
	return s.UpdateRollbackSQLById(fmt.Sprintf("%v", baseSQL.ID), attr)
}

func (s *Storage) UpdateRollbackSQLById(rollbackSQLId string, attrs ...interface{}) error {
	err := s.db.Table(RollbackSQL{}.TableName()).Where("id = ?", rollbackSQLId).Update(attrs...).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRelatedDDLTask(task *Task) ([]Task, error) {
	tasks := []Task{}
	err := s.db.Where(Task{
		InstanceId: task.InstanceId,
		Schema:     task.Schema,
		PassRate:   1,
		Status:     TaskStatusAudited,
	}).Preload("Instance").Preload("ExecuteSQLs").Find(&tasks).Error
	return tasks, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetTaskByInstanceId(instanceId uint) ([]Task, error) {
	tasks := []Task{}
	err := s.db.Where(&Task{InstanceId: instanceId}).Find(&tasks).Error
	return tasks, errors.New(errors.ConnectStorageError, err)
}

type TaskSQLDetail struct {
	Number      uint           `json:"number"`
	ExecSQL     string         `json:"exec_sql"`
	AuditResult string         `json:"audit_result"`
	AuditLevel  string         `json:"audit_level"`
	AuditStatus string         `json:"audit_status"`
	ExecResult  string         `json:"exec_result"`
	ExecStatus  string         `json:"exec_status"`
	RollbackSQL sql.NullString `json:"rollback_sql"`
}

var taskSQLsQueryTpl = `SELECT e_sql.number, e_sql.content AS exec_sql, r_sql.content AS rollback_sql,
e_sql.audit_result, e_sql.audit_level, e_sql.audit_status, e_sql.exec_result, e_sql.exec_status

{{- template "body" . -}}

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var taskSQLsCountTpl = `SELECT COUNT(*)

{{- template "body" . -}}
`

var taskSQLsQueryBodyTpl = `
{{ define "body" }}
FROM execute_sql_detail AS e_sql
LEFT JOIN rollback_sql_detail AS r_sql ON e_sql.id = r_sql.execute_sql_id
WHERE
e_sql.deleted_at IS NULL
AND e_sql.task_id = :task_id

{{- if .filter_exec_status }}
AND e_sql.exec_status = :filter_exec_status
{{- end }}

{{- if .filter_audit_status }}
AND e_sql.audit_status = :filter_audit_status
{{- end }}

{{- if .no_duplicate }}
AND e_sql.id IN (
SELECT SQL_BIG_RESULT MIN(id) AS id FROM execute_sql_detail WHERE task_id = :task_id 
GROUP BY audit_fingerprint, IFNULL(audit_fingerprint, id) ORDER BY null
)
{{- end }}
ORDER BY e_sql.id
{{- end }}
`

func (s *Storage) GetTaskSQLsByReq(data map[string]interface{}) (
	result []*TaskSQLDetail, count uint64, err error) {

	err = s.getListResult(taskSQLsQueryBodyTpl, taskSQLsQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(taskSQLsQueryBodyTpl, taskSQLsCountTpl, data)
	return result, count, err
}

func (s *Storage) DeleteTask(task *Task) error {
	return s.TxExec(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM tasks WHERE id = ?", task.ID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM execute_sql_detail WHERE task_id = ?", task.ID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM rollback_sql_detail WHERE task_id = ?", task.ID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) GetExpiredTasks(start time.Time) ([]*Task, error) {
	tasks := []*Task{}
	err := s.db.Model(&Task{}).Select("tasks.id").
		Joins("LEFT JOIN workflow_records ON tasks.id = workflow_records.task_id").
		Where("tasks.created_at < ?", start).
		Where("workflow_records.id is NULL").
		Scan(&tasks).Error

	return tasks, errors.New(errors.ConnectStorageError, err)
}
