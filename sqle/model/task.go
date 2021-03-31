package model

import (
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pingcap/parser/ast"
)

// task action
const (
	TASK_ACTION_AUDIT = iota + 1
	TASK_ACTION_EXECUTE
	TASK_ACTION_ROLLBACK
)

const (
	SQL_TYPE_DML                      = "dml"
	SQL_TYPE_DDL                      = "ddl"
	SQL_TYPE_MULTI                    = "dml&ddl"
	SQL_TYPE_PROCEDURE_FUNCTION       = "procedure&function"
	SQL_TYPE_PROCEDURE_FUNCTION_MULTI = "procedure&function&dml&ddl"
)

const (
	TaskStatusInit             = "initialized"
	TaskStatusAudited          = "audited"
	TaskStatusExecuting        = "executing"
	TaskStatusExecuteSucceeded = "exec_succeeded"
	TaskStatusExecuteFailed    = "exec_failed"
)

type Task struct {
	Model
	InstanceId   uint    `json:"instance_id"`
	Schema       string  `json:"instance_schema" gorm:"column:instance_schema" example:"db1"`
	PassRate     float64 `json:"pass_rate"`
	SQLType      string  `json:"sql_type" gorm:"column:sql_type"`
	Status       string  `json:"status" gorm:"default:\"initialized\""`
	CreateUserId uint

	CreateUser   *User          `gorm:"foreignkey:CreateUserId"`
	Instance     *Instance      `json:"-" gorm:"foreignkey:InstanceId"`
	ExecuteSQLs  []*ExecuteSQL  `json:"-" gorm:"foreignkey:TaskId"`
	RollbackSQLs []*RollbackSQL `json:"-" gorm:"foreignkey:TaskId"`
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
	TaskId          uint       `json:"-"`
	Number          uint       `json:"number"`
	Content         string     `json:"sql" gorm:"type:text"`
	StartBinlogFile string     `json:"start_binlog_file"`
	StartBinlogPos  int64      `json:"start_binlog_pos"`
	EndBinlogFile   string     `json:"end_binlog_file"`
	EndBinlogPos    int64      `json:"end_binlog_pos"`
	RowAffects      int64      `json:"row_affects"`
	ExecStatus      string     `json:"exec_status" gorm:"default:\"initialized\""`
	ExecResult      string     `json:"exec_result"`
	Stmts           []ast.Node `json:"-" gorm:"-"`
	Fingerprint     string     `json:"fingerprint" gorm:"type:text"`
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
	AuditLevel  string `json:"audit_level"` // level: error, warn, notice, normal
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
	ExecuteSQLId uint `gorm:"column:execute_sql_id"`
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

func (t *Task) ValidAction(typ int) error {
	switch typ {
	case TASK_ACTION_AUDIT:
		// audit sql allowed at all times
		return nil
	case TASK_ACTION_EXECUTE:
		if t.HasDoingExecute() {
			return errors.New(errors.TASK_ACTION_DONE, fmt.Errorf("task has been executed"))
		}
	case TASK_ACTION_ROLLBACK:
		if t.HasDoingRollback() {
			return errors.New(errors.TASK_ACTION_DONE, fmt.Errorf("task has been rolled back"))
		}
		if t.IsExecuteFailed() {
			return errors.New(errors.TASK_ACTION_INVALID, fmt.Errorf("task is commit failed, not allow rollback"))
		}
		if !t.HasDoingExecute() {
			return errors.New(errors.TASK_ACTION_INVALID, fmt.Errorf("task need commit first"))
		}
	}
	return nil
}

func (s *Storage) GetTaskById(taskId string) (*Task, bool, error) {
	task := &Task{}
	err := s.db.Where("id = ?", taskId).Preload("Instance").First(task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return task, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetTaskDetailById(taskId string) (*Task, bool, error) {
	task := &Task{}
	err := s.db.Where("id = ?", taskId).Preload("Instance").
		Preload("ExecuteSQLs").Preload("RollbackSQLs").First(task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return task, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetTaskExecuteSQLContent(taskId string) ([]string, error) {
	var SQLContents []string
	rows, err := s.db.Model(&ExecuteSQL{}).Select("content").
		Where("task_id = ?", taskId).Rows()
	if err != nil {
		return []string{}, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	defer rows.Close()
	for rows.Next() {
		var content string
		rows.Scan(&content)
		SQLContents = append(SQLContents, content)
	}
	return SQLContents, nil
}

func (s *Storage) GetTasksInstanceName() ([]string, error) {
	query := `SELECT instances.name AS name FROM tasks 
JOIN instances ON tasks.instance_id = instances.id 
GROUP BY instances.name`

	var instancesName []string
	rows, err := s.db.Raw(query).Rows()
	if err != nil {
		return []string{}, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		instancesName = append(instancesName, name)
	}
	return instancesName, nil
}

func (s *Storage) GetTasks() ([]Task, error) {
	tasks := []Task{}
	err := s.db.Find(&tasks).Error
	return tasks, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

//func (s *Storage) GetExpiredTaskIds() ([]string, error) {
//	now := time.Now().Format(time.RFC3339)
//	query := fmt.Sprintf(`SELECT id FROM tasks WHERE expired_at < %s`, now)
//
//	var ids []string
//	rows, err := s.db.Raw(query).Rows()
//	if err != nil {
//		return []string{}, errors.New(errors.CONNECT_STORAGE_ERROR, err)
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var id string
//		rows.Scan(&id)
//		ids = append(ids, id)
//	}
//	return ids, nil
//}

func (s *Storage) GetTasksByIds(ids []string) ([]Task, error) {
	tasks := []Task{}
	err := s.db.Where("id IN (?)", ids).Find(&tasks).Error
	return tasks, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateTask(task *Task, attrs ...interface{}) error {
	err := s.db.Table("tasks").Where("id = ?", task.ID).Update(attrs...).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateExecuteSQLs(task *Task, ExecuteSQLs []*ExecuteSQL) error {
	tx := s.db.Begin()
	//if err := tx.Unscoped().Where("task_id=?", task.ID).Delete(ExecuteSQL{}).Error; err != nil {
	//	return err
	//}

	for _, executeSQL := range ExecuteSQLs {
		currentSql := executeSQL
		if err := tx.Save(currentSql).Error; err != nil {
			tx.Rollback()
			return errors.New(errors.CONNECT_STORAGE_ERROR, err)
		}
		//if err := tx.Exec("INSERT execute_sql_detail(task_id, number, content, "+
		//	"start_binlog_file, start_binlog_pos, end_binlog_file, end_binlog_pos, row_affects, "+
		//	"exec_status, exec_result, fingerprint, audit_status, audit_result, audit_level) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		//	task.ID, executeSQL.Number, executeSQL.Content, executeSQL.StartBinlogFile, executeSQL.StartBinlogPos,
		//	executeSQL.EndBinlogFile, executeSQL.EndBinlogPos, executeSQL.RowAffects, executeSQL.ExecStatus,
		//	executeSQL.ExecResult, executeSQL.Fingerprint, executeSQL.AuditStatus, executeSQL.AuditResult,
		//	executeSQL.AuditLevel).Error; err != nil {
		//	tx.Rollback()
		//	return err
		//}
	}
	return errors.New(errors.CONNECT_STORAGE_ERROR, tx.Commit().Error)
}

func (s *Storage) UpdateRollbackSQLs(task *Task, rollbackSQLs []*RollbackSQL) error {
	tx := s.db.Begin()
	//if err := tx.Unscoped().Where("task_id=?", task.ID).Delete(RollbackSQL{}).Error; err != nil {
	//	return err
	//}

	for _, rollbackSQL := range rollbackSQLs {
		currentSql := rollbackSQL
		if err := tx.Save(currentSql).Error; err != nil {
			tx.Rollback()
			return errors.New(errors.CONNECT_STORAGE_ERROR, err)
		}
		//if err := tx.Exec("INSERT INTO rollback_sql_detail(task_id, number, content, "+
		//	"start_binlog_file, start_binlog_pos, end_binlog_file, end_binlog_pos, row_affects, "+
		//	"exec_status, exec_result, fingerprint, execute_sql_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)",
		//	task.ID, rollbackSQL.Number, rollbackSQL.Content, rollbackSQL.StartBinlogFile,
		//	rollbackSQL.StartBinlogPos, rollbackSQL.EndBinlogFile, rollbackSQL.EndBinlogPos,
		//	rollbackSQL.RowAffects, rollbackSQL.ExecStatus, rollbackSQL.ExecResult,
		//	rollbackSQL.Fingerprint, rollbackSQL.ExecuteSQLId).Error; err != nil {
		//	tx.Rollback()
		//	return err
		//}
	}
	return errors.New(errors.CONNECT_STORAGE_ERROR, tx.Commit().Error)
}

func (s *Storage) UpdateTaskStatusById(taskId uint, status string) error {
	err := s.db.Model(&Task{}).Where("id = ?", taskId).Update(map[string]string{
		"status": status,
	}).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
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
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
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
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) HardDeleteRollbackSQLByTaskIds(taskIds []string) error {
	rollbackSQL := RollbackSQL{}
	err := s.db.Unscoped().Where("task_id IN (?)", taskIds).Delete(rollbackSQL).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

//func (s *Storage) GetRollbackSqlByTaskId(taskId string, commitSqlNum []string) ([]RollbackSQL, error) {
//	rollbackSqls := []RollbackSQL{}
//	querySql := "task_id=?"
//	queryArgs := make([]interface{}, 0)
//	queryArgs = append(queryArgs, taskId)
//	if len(commitSqlNum) > 0 {
//		querySql += " AND COMMIT_SQL_NUMBER IN (?)"
//		queryArgs = append(queryArgs, commitSqlNum)
//	}
//	err := s.db.Where(querySql, queryArgs...).Find(&rollbackSqls).Error
//	return rollbackSqls, errors.New(errors.CONNECT_STORAGE_ERROR, err)
//}

func (s *Storage) GetRelatedDDLTask(task *Task) ([]Task, error) {
	tasks := []Task{}
	err := s.db.Where(Task{
		InstanceId: task.InstanceId,
		Schema:     task.Schema,
		PassRate:   1,
		SQLType:    SQL_TYPE_DDL,
		Status:     TaskStatusAudited,
	}).Preload("Instance").Preload("ExecuteSQLs").Find(&tasks).Error
	return tasks, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) HardDeleteExecuteSqlResultByTaskIds(ids []string) error {
	executeSQL := ExecuteSQL{}
	err := s.db.Unscoped().Where("task_id IN (?)", ids).Delete(executeSQL).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetExecErrorExecuteSQLsByTaskId(taskId string) ([]ExecuteSQL, error) {
	executeSQLs := []ExecuteSQL{}
	err := s.db.Not("exec_result", []string{"ok", ""}).Where("task_id=? ", taskId).Find(&executeSQLs).Error
	return executeSQLs, errors.New(errors.CONNECT_STORAGE_ERROR, err)
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
GROUP BY audit_result, IFNULL(audit_result, id), fingerprint, IFNULL(fingerprint, id) ORDER BY null
)
{{- end }}

{{- if .filter_next_step_assignee_user_name }}
AND ass_user.login_name = :filter_next_step_assignee_user_name
{{- end }}
ORDER BY e_sql.id
{{- end }}
`

func (s *Storage) GetTaskSQLsByReq(data map[string]interface{}) ([]*TaskSQLDetail, uint64, error) {
	result := []*TaskSQLDetail{}
	count, err := s.getListResult(taskSQLsQueryBodyTpl, taskSQLsQueryTpl, taskSQLsCountTpl, data, &result)
	return result, count, err
}

func (s *Storage) DeleteTasksByIds(ids []string) error {
	tasks := Task{}
	err := s.db.Where("id IN (?)", ids).Delete(tasks).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) HardDeleteTasksByIds(ids []string) error {
	s.db.Begin()
	tx, err := s.db.DB().Begin()
	if err != nil {
		return errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	if _, err := tx.Exec("DELETE FROM task WHERE id in (?)", ids); err != nil {
		tx.Rollback()
		return errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	if _, err := tx.Exec("DELETE FROM commit_sql_detail WHERE task_id in (?)", ids); err != nil {
		tx.Rollback()
		return errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}

	if _, err := tx.Exec("DELETE FROM rollback_sql_detail WHERE task_id in (?)", ids); err != nil {
		tx.Rollback()
		return errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}

	err = tx.Commit()
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
