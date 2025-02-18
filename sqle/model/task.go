package model

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/utils"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

const (
	TaskStatusInit             = "initialized"
	TaskStatusAudited          = "audited"
	TaskStatusExecuting        = "executing"
	TaskStatusManuallyExecuted = "manually_executed"
	TaskStatusExecuteSucceeded = "exec_succeeded"
	TaskStatusExecuteFailed    = "exec_failed"
	TaskStatusTerminating      = "terminating"
	TaskStatusTerminateFail    = "terminate_failed"
	TaskStatusTerminateSucc    = "terminate_succeeded"
)

const (
	TaskSQLSourceFromFormData       = "form_data"
	TaskSQLSourceFromSQLFile        = "sql_file"
	TaskSQLSourceFromMyBatisXMLFile = "mybatis_xml_file"
	TaskSQLSourceFromZipFile        = "zip_file"
	TaskSQLSourceFromGitRepository  = "git_repository"
	TaskSQLSourceFromAuditPlan      = "audit_plan"
)

const TaskExecResultOK = "OK"

const ExecModeSqlFile = "sql_file"
const ExecModeSqls = "sqls"

type Task struct {
	Model
	InstanceId           uint64  `json:"instance_id"`
	Schema               string  `json:"instance_schema" gorm:"column:instance_schema;type:varchar(255)" example:"db1"`
	PassRate             float64 `json:"pass_rate"`
	Score                int32   `json:"score"`
	AuditLevel           string  `json:"audit_level" gorm:"type:varchar(255)"`
	SQLSource            string  `json:"sql_source" gorm:"column:sql_source;type:varchar(255)"`
	DBType               string  `json:"db_type" gorm:"default:'mysql';type:varchar(255)" example:"mysql"`
	Status               string  `json:"status" gorm:"default:\"initialized\";type:varchar(255)"`
	GroupId              uint    `json:"group_id" gorm:"column:group_id"`
	CreateUserId         uint64
	RuleTemplateID       uint `json:"rule_template_id" gorm:"column:rule_template_id"`
	ExecStartAt          *time.Time
	ExecEndAt            *time.Time
	ExecMode             string         `json:"exec_mode" gorm:"default:'sqls';type:varchar(255)" example:"sqls"`
	EnableBackup         bool           `gorm:"column:enable_backup"`
	BackupMaxRows        uint64         `json:"backup_max_rows" gorm:"column:backup_max_rows;not null;default:0"`
	InstanceEnableBackup bool           `gorm:"column:instance_enable_backup"` // 用于记录创建task时，instance备份开关的状态
	FileOrderMethod      string         `json:"file_order_method" gorm:"column:file_order_method;type:varchar(255)"`
	Instance             *Instance      `json:"-" gorm:"-"`
	RuleTemplate         *RuleTemplate  `json:"-" gorm:"foreignkey:RuleTemplateID"`
	ExecuteSQLs          []*ExecuteSQL  `json:"-" gorm:"foreignkey:TaskId"`
	RollbackSQLs         []*RollbackSQL `json:"-" gorm:"foreignkey:TaskId"`
	AuditFiles           []*AuditFile   `json:"-" gorm:"foreignkey:TaskId"`
}

func (t *Task) RuleTemplateName() string {
	if t.RuleTemplate != nil {
		return t.RuleTemplate.Name
	}
	return ""
}

func (t *Task) InstanceName() string {
	if t.Instance != nil {
		return t.Instance.Name
	}
	return ""
}

func (t *Task) InstanceHost() string {
	if t.Instance != nil {
		return t.Instance.Host
	}
	return ""
}

func (t *Task) TaskExecStartAt() string {
	if t.ExecStartAt == nil {
		return ""
	}
	return t.ExecStartAt.Format("2006-01-02 15:04:05")
}

func (t *Task) TaskExecEndAt() string {
	if t.ExecEndAt == nil {
		return ""
	}
	return t.ExecEndAt.Format("2006-01-02 15:04:05")
}

const (
	SQLAuditStatusInitialized = "initialized"
	SQLAuditStatusDoing       = "doing"
	SQLAuditStatusFinished    = "finished"
)

const (
	SQLExecuteStatusInitialized      = "initialized"
	SQLExecuteStatusDoing            = "doing"
	SQLExecuteStatusFailed           = "failed"
	SQLExecuteStatusSucceeded        = "succeeded"
	SQLExecuteStatusManuallyExecuted = "manually_executed"
	SQLExecuteStatusTerminateSucc    = "terminate_succeeded"
	SQLExecuteStatusTerminateFailed  = "terminate_failed"
	SQLExecuteStatusExecuteRollback  = "execute_rollback" // 执行回滚
)

type BaseSQL struct {
	Model
	TaskId uint `json:"-" gorm:"index"`
	Number uint `json:"number"`

	// Content store single SQL or batch SQLs
	//
	// Content may store batch SQLs When BaseSQL embed to RollbackSQL.
	// Split Content to single SQL before execute RollbackSQL.
	Content         string `json:"sql" gorm:"type:longtext"`
	Description     string `json:"description" gorm:"type:text"`
	StartBinlogFile string `json:"start_binlog_file" gorm:"type:varchar(255)"`
	StartBinlogPos  int64  `json:"start_binlog_pos"`
	EndBinlogFile   string `json:"end_binlog_file" gorm:"type:varchar(255)"`
	EndBinlogPos    int64  `json:"end_binlog_pos"`
	RowAffects      int64  `json:"row_affects"`
	ExecStatus      string `json:"exec_status" gorm:"default:\"initialized\""`
	ExecResult      string `json:"exec_result" gorm:"type:text"`
	ExecBatchId     uint64 `json:"exec_batch_id"`
	Schema          string `json:"schema" gorm:"type:varchar(255)"`
	SourceFile      string `json:"source_file" gorm:"type:varchar(255)"`
	StartLine       uint64 `json:"start_line" gorm:"not null"`
	SQLType         string `json:"sql_type" gorm:"type:varchar(255)"` // such as DDL,DML,DQL...
}

func (s *BaseSQL) GetExecStatusDesc(ctx context.Context) string {
	switch s.ExecStatus {
	case SQLExecuteStatusInitialized:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLExecuteStatusInitialized)
	case SQLExecuteStatusDoing:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLExecuteStatusDoing)
	case SQLExecuteStatusFailed:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLExecuteStatusFailed)
	case SQLExecuteStatusSucceeded:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLExecuteStatusSucceeded)
	case SQLExecuteStatusManuallyExecuted:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLExecuteStatusManuallyExecuted)
	default:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLExecuteStatusUnknown)
	}
}

type AuditResult struct {
	Level               string              `json:"level"`
	RuleName            string              `json:"rule_name"`
	ExecutionFailed     bool                `json:"execution_failed"`
	I18nAuditResultInfo I18nAuditResultInfo `json:"i18n_audit_result_info"`
}

func (ar *AuditResult) GetAuditMsgByLangTag(lang language.Tag) string {
	return ar.I18nAuditResultInfo.GetAuditResultInfoByLangTag(lang).Message
}

func (ar *AuditResult) GetAuditErrorMsgByLangTag(lang language.Tag) string {
	return ar.I18nAuditResultInfo.GetAuditResultInfoByLangTag(lang).ErrorInfo
}

type AuditResultInfo struct {
	Message   string
	ErrorInfo string
}

type I18nAuditResultInfo map[language.Tag]AuditResultInfo

func (i *I18nAuditResultInfo) GetAuditResultInfoByLangTag(lang language.Tag) *AuditResultInfo {
	if i == nil {
		return &AuditResultInfo{}
	}

	if info, ok := (*i)[lang]; ok {
		return &info
	}

	info := (*i)[i18nPkg.DefaultLang]
	return &info
}

func (i I18nAuditResultInfo) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	return string(b), err
}

func (i *I18nAuditResultInfo) Scan(input interface{}) error {
	if input == nil {
		return nil
	}
	if data, ok := input.([]byte); !ok {
		return fmt.Errorf("I18nAuditResultInfo Scan input is not bytes")
	} else {
		return json.Unmarshal(data, i)
	}
}

func ConvertI18NAuditResultInfoMapToI18nStr(m I18nAuditResultInfo) i18nPkg.I18nStr {
	s := make(i18nPkg.I18nStr, len(m))
	for lang, v := range m {
		s[lang] = v.Message
	}
	return s
}

func ConvertI18nStrToI18NAuditResultInfoMap(s i18nPkg.I18nStr) I18nAuditResultInfo {
	m := make(I18nAuditResultInfo, len(s))
	for lang, v := range s {
		m[lang] = AuditResultInfo{
			Message: v,
		}
	}
	return m
}

type AuditResults []AuditResult

func (a *AuditResults) GetAuditJsonStrByLangTag(lang language.Tag) string {
	type AuditResultRes struct {
		Level           string `json:"level"`
		Message         string `json:"message"`
		RuleName        string `json:"rule_name"`
		ErrorInfo       string `json:"error_info"`
		ExecutionFailed bool   `json:"execution_failed"`
	}
	results := make([]AuditResultRes, len(*a))
	for k, v := range *a {
		results[k] = AuditResultRes{
			Level:           v.Level,
			Message:         v.GetAuditMsgByLangTag(lang),
			ErrorInfo:       v.GetAuditErrorMsgByLangTag(lang),
			RuleName:        v.RuleName,
			ExecutionFailed: v.ExecutionFailed,
		}
	}
	data, _ := json.Marshal(results)
	return string(data)
}

func (a AuditResults) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	return string(b), err
}

func (a *AuditResults) Scan(input interface{}) error {
	if input == nil {
		return nil
	}
	if data, ok := input.([]byte); !ok {
		return fmt.Errorf("AuditResults Scan input is not bytes")
	} else {
		return json.Unmarshal(data, a)
	}
}

// todo check somewhere fmt Sprint AuditResults to frontend?
func (a *AuditResults) String(ctx context.Context) string {
	lang := locale.Bundle.GetLangTagFromCtx(ctx)
	msgs := make([]string, len(*a))
	for i := range *a {
		res := (*a)[i]
		msgs[i] = res.GetAuditMsgByLangTag(lang)
	}
	return strings.Join(msgs, "\n")
}

func (a *AuditResults) Append(dar *driverV2.AuditResult) {
	for i := range *a {
		ar := (*a)[i]
		if ar.Level == string(dar.Level) && ar.RuleName == dar.RuleName {
			return
		}
	}
	newAr := ConvertAuditResultFromDriverToModel(dar)
	*a = append(*a, *newAr)
}

func ConvertAuditResultFromDriverToModel(dar *driverV2.AuditResult) *AuditResult {
	newAr := &AuditResult{
		Level:               string(dar.Level),
		RuleName:            dar.RuleName,
		ExecutionFailed:     dar.ExecutionFailed,
		I18nAuditResultInfo: make(map[language.Tag]AuditResultInfo, len(dar.I18nAuditResultInfo)),
	}
	for langTag, info := range dar.I18nAuditResultInfo {
		newAr.I18nAuditResultInfo[langTag] = AuditResultInfo{
			Message:   info.Message,
			ErrorInfo: info.ErrorInfo,
		}
	}
	return newAr
}

func ConvertAuditResultFromModelToDriver(mar *AuditResult) *driverV2.AuditResult {
	newAr := &driverV2.AuditResult{
		Level:               driverV2.RuleLevel(mar.Level),
		RuleName:            mar.RuleName,
		ExecutionFailed:     mar.ExecutionFailed,
		I18nAuditResultInfo: make(map[language.Tag]driverV2.AuditResultInfo, len(mar.I18nAuditResultInfo)),
	}
	for langTag, info := range mar.I18nAuditResultInfo {
		newAr.I18nAuditResultInfo[langTag] = driverV2.AuditResultInfo{
			Message:   info.Message,
			ErrorInfo: info.ErrorInfo,
		}
	}
	return newAr
}

type ExecuteSQL struct {
	BaseSQL
	AuditStatus  string       `json:"audit_status" gorm:"default:\"initialized\""`
	AuditResults AuditResults `json:"audit_results" gorm:"type:json"`
	// AuditFingerprint generate from SQL and SQL audit result use MD5 hash algorithm,
	// it used for deduplication in one audit task.
	AuditFingerprint string `json:"audit_fingerprint" gorm:"index;type:char(32)"`
	// AuditLevel has four level: error, warn, notice, normal.
	AuditLevel string      `json:"audit_level" gorm:"type:varchar(255)"`
	BackupTask *BackupTask `json:"-" gorm:"foreignkey:execute_sql_id"`
}

func (s ExecuteSQL) TableName() string {
	return "execute_sql_detail"
}

func (s *ExecuteSQL) GetAuditStatusDesc(ctx context.Context) string {
	switch s.AuditStatus {
	case SQLAuditStatusInitialized:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLAuditStatusInitialized)
	case SQLAuditStatusDoing:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLAuditStatusDoing)
	case SQLAuditStatusFinished:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLAuditStatusFinished)
	default:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLAuditStatusUnknown)
	}
}

func (s *ExecuteSQL) GetAuditResults(ctx context.Context) string {
	if len(s.AuditResults) == 0 {
		return ""
	}

	return s.AuditResults.String(ctx)
}

func (s *ExecuteSQL) GetAuditResultDesc(ctx context.Context) string {
	if len(s.AuditResults) == 0 {
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.SQLAuditResultDescPass)
	}

	return s.AuditResults.String(ctx)
}

func (s *Storage) BatchSaveExecuteSqls(models []*ExecuteSQL) error {
	return s.Tx(func(txDB *gorm.DB) error {
		for _, model := range models {
			if err := txDB.Save(&model).Error; err != nil {
				txDB.Rollback()
				return err
			}
		}
		return nil
	})
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

func (s *Storage) GetTaskStatusByID(id string) (string, error) {
	task := &Task{}
	err := s.db.Select("status").Where("id = (?)", id).First(task).Error
	if err != nil {
		return "", err
	}
	return task.Status, nil
}

func (s *Storage) GetTaskDbTypeByID(id string) (string, error) {
	task := &Task{}
	err := s.db.Select("db_type").Where("id = (?)", id).First(task).Error
	if err != nil {
		return "", err
	}
	return task.DBType, nil
}

func (s *Storage) GetTaskById(taskId string) (*Task, bool, error) {
	task := &Task{}
	err := s.db.Where("id = ?", taskId).First(task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return task, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetTasksByIds(taskIds []uint) (tasks []*Task, foundAllIds bool, err error) {
	taskIds = utils.RemoveDuplicateUint(taskIds)
	err = s.db.Where("id IN (?)", taskIds).Find(&tasks).Error
	if err != nil {
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	if len(tasks) == len(taskIds) {
		foundAllIds = true
	}
	return
}

func (s *Storage) GetTaskDetailById(taskId string) (*Task, bool, error) {
	task := &Task{}
	err := s.db.Where("id = ?", taskId).
		Preload("RuleTemplate").Preload("ExecuteSQLs").Preload("RollbackSQLs").First(task).Error
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
		if err := rows.Scan(&content); err != nil {
			return nil, errors.New(errors.DataInvalid, err)
		}
		buff.WriteString(strings.TrimRight(content, ";"))
		buff.WriteString(";\n")
	}
	if rows.Err() != nil {
		return nil, errors.New(errors.DataParseFail, rows.Err())
	}

	return buff.Bytes(), nil
}

func (s *Storage) UpdateTask(task *Task, attrs interface{}) error {
	err := s.db.Table("tasks").Where("id = ?", task.ID).Updates(attrs).Error

	// .Update(attrs...).Error
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
	err := updateTaskStatusById(s.db, taskId, status)
	return errors.New(errors.ConnectStorageError, err)
}

func updateTaskStatusById(tx *gorm.DB, taskId uint, status string) error {
	return tx.Model(&Task{}).Where("id = ?", taskId).Updates(map[string]interface{}{
		"status": status,
	}).Error
}

func (s *Storage) UpdateTaskStatusByIDs(taskIDs []uint, attrs interface{}) error {
	err := s.db.Model(&Task{}).Where("id IN (?)", taskIDs).Updates(attrs).Error
	return errors.ConnectStorageErrWrapper(err)
}

func updateExecuteSQLStatusByTaskIdAndStatus(tx *gorm.DB, taskId uint, status []string, updateToStatus string) error {
	query := "UPDATE execute_sql_detail SET exec_status=? WHERE task_id=? AND exec_status IN (?)"
	return tx.Exec(query, updateToStatus, taskId, status).Error
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

func (s *Storage) UpdateExecuteSQLById(executeSQLId string, attrs interface{}) error {
	err := s.db.Table(ExecuteSQL{}.TableName()).Where("id = ?", executeSQLId).Updates(attrs).Error
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

func (s *Storage) UpdateRollbackSQLById(rollbackSQLId string, attrs interface{}) error {
	err := s.db.Table(RollbackSQL{}.TableName()).Where("id = ?", rollbackSQLId).Updates(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRelatedDDLTask(task *Task) ([]Task, error) {
	tasks := []Task{}
	err := s.db.Where(Task{
		InstanceId: task.InstanceId,
		Schema:     task.Schema,
		PassRate:   1,
		Status:     TaskStatusAudited,
	}).Preload("ExecuteSQLs").Find(&tasks).Error
	return tasks, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetTaskByInstanceId(instanceId uint64) ([]Task, error) {
	tasks := []Task{}
	err := s.db.Where(&Task{InstanceId: instanceId}).Find(&tasks).Error
	return tasks, errors.New(errors.ConnectStorageError, err)
}

type TaskSQLDetail struct {
	Id            uint           `json:"id"`
	Number        uint           `json:"number"`
	Description   string         `json:"description"`
	ExecSQL       string         `json:"exec_sql"`
	SQLSourceFile sql.NullString `json:"sql_source_file"`
	SQLStartLine  uint64         `json:"sql_start_line"`
	AuditResults  AuditResults   `json:"audit_results"`
	AuditLevel    string         `json:"audit_level"`
	AuditStatus   string         `json:"audit_status"`
	ExecResult    string         `json:"exec_result"`
	ExecStatus    string         `json:"exec_status"`
	RollbackSQL   sql.NullString `json:"rollback_sql"`
	SQLType       sql.NullString `json:"sql_type"`
}

func (t *TaskSQLDetail) GetAuditResults(ctx context.Context) string {
	if len(t.AuditResults) == 0 {
		return ""
	}

	return t.AuditResults.String(ctx)
}

var taskSQLsQueryTpl = `SELECT e_sql.id,e_sql.number, e_sql.description, e_sql.content AS exec_sql,  e_sql.source_file AS sql_source_file, e_sql.start_line AS sql_start_line, e_sql.sql_type,
e_sql.audit_results, e_sql.audit_level, e_sql.audit_status, e_sql.exec_result, e_sql.exec_status

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
{{- if .filter_audit_file_id }}
LEFT JOIN audit_files ON audit_files.task_id = e_sql.task_id
AND audit_files.file_name = e_sql.source_file
{{- end }}
WHERE
e_sql.deleted_at IS NULL
AND e_sql.task_id = :task_id

{{- if .filter_exec_status }}
AND e_sql.exec_status = :filter_exec_status
{{- end }}

{{- if .filter_audit_status }}
AND e_sql.audit_status = :filter_audit_status
{{- end }}

{{- if .filter_audit_level }}
AND e_sql.audit_level = :filter_audit_level
{{- end }}

{{- if .filter_audit_file_id }}
AND audit_files.id = :filter_audit_file_id
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
		Joins("LEFT JOIN workflow_instance_records ON tasks.id = workflow_instance_records.task_id").
		Joins("LEFT JOIN sql_audit_records ON tasks.id = sql_audit_records.task_id").
		Where("tasks.created_at < ?", start).
		Where("workflow_instance_records.id is NULL").
		Where("sql_audit_records.id is NULL").
		Scan(&tasks).Error

	return tasks, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetTaskSQLByNumber(taskId, number string) (*ExecuteSQL, bool, error) {
	e := &ExecuteSQL{}
	err := s.db.Where("task_id = ?", taskId).Where("number = ?", number).First(e).Error
	if err == gorm.ErrRecordNotFound {
		return e, false, nil
	}
	return e, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetTaskSQLCountByTaskID(taskId uint) (int64, error) {
	var count int64
	return count, s.db.Model(&ExecuteSQL{}).Where("task_id = ?", taskId).Count(&count).Error
}

type TaskGroup struct {
	Model
	Tasks []*Task `json:"tasks" gorm:"foreignkey:GroupId"`
}

func (s *Storage) GetTaskGroupByGroupId(groupId uint) (*TaskGroup, error) {
	taskGroup := &TaskGroup{}
	err := s.db.Preload("Tasks").
		Where("id = ?", groupId).Find(&taskGroup).Error
	return taskGroup, errors.New(errors.ConnectStorageError, err)
}

type SqlExecuteStatistic struct {
	InstanceID       uint `json:"instance_id"`
	AvgExecutionTime uint `json:"avg_execution_time"`
	MaxExecutionTime uint `json:"max_execution_time"`
	MinExecutionTime uint `json:"min_execution_time"`
}

func (s *Storage) GetSqlAvgExecutionTimeStatistic(limit uint) ([]*SqlExecuteStatistic, error) {
	var sqlExecuteStatistics []*SqlExecuteStatistic
	err := s.db.Model(&Workflow{}).Select("t.instance_id,"+
		"round(avg(timestampdiff(second, t.exec_start_at, t.exec_end_at))) avg_execution_time,"+
		"max(timestampdiff(second, t.exec_start_at, t.exec_end_at)) max_execution_time,"+
		"min(timestampdiff(second, t.exec_start_at, t.exec_end_at)) min_execution_time").
		Joins("left join workflow_records wr on workflows.workflow_record_id = wr.id").
		Joins("left join workflow_instance_records wir on wr.id = wir.workflow_record_id").
		Joins("left join tasks t on wir.task_id = t.id").
		Where("t.status = ?", TaskStatusExecuteSucceeded).
		Group("t.instance_id").Order("avg_execution_time desc").Limit(int(limit)).
		Scan(&sqlExecuteStatistics).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	return sqlExecuteStatistics, nil
}

type SqlExecutionCount struct {
	Count      uint   `json:"count"`
	InstanceId uint64 `json:"instance_id"`
}
