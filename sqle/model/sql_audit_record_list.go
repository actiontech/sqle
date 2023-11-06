package model

import (
	"database/sql"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
)

type SQLAuditRecordListItem struct {
	AuditRecordId   string          `json:"audit_record_id"`
	RecordCreatedAt *time.Time      `json:"record_created_at"`
	CreatorId       string          `json:"creator_id"`
	Tags            sql.NullString  `json:"tags"`
	InstanceId      uint64          `json:"instance_id"`
	TaskId          uint            `json:"task_id"`
	DbType          string          `json:"db_type"`
	InstanceSchema  string          `json:"instance_schema"`
	AuditLevel      sql.NullString  `json:"audit_level"`
	TaskStatus      string          `json:"task_status"`
	AuditScore      sql.NullInt32   `json:"audit_score"`
	AuditPassRate   sql.NullFloat64 `json:"audit_pass_rate"`
	SQLSource       string          `json:"sql_source"`
}

var sqlAuditRecordQueryTpl = `
SELECT sql_audit_records.audit_record_id AS audit_record_id,
       sql_audit_records.created_at      AS record_created_at,
       sql_audit_records.tags            AS tags,
			 sql_audit_records.creator_id      AS creator_id,
       tasks.id                          AS task_id,
       tasks.db_type                     AS db_type,
			 tasks.instance_id      		       AS instance_id,
       tasks.instance_schema             AS instance_schema,
       tasks.audit_level                 AS audit_level,
       tasks.status                      AS task_status,
       tasks.score                       AS audit_score,
       tasks.pass_rate                   AS audit_pass_rate,
       tasks.sql_source                  AS sql_source

{{- template "body" . -}}
ORDER BY sql_audit_records.created_at DESC
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var sqlAuditRecordCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var sqlAuditRecordQueryBodyTpl = `
{{ define "body" }}
FROM sql_audit_records
LEFT JOIN tasks ON sql_audit_records.task_id = tasks.id

WHERE
sql_audit_records.deleted_at IS NULL
AND sql_audit_records.project_id = :filter_project_id

{{- if .check_user_can_access }}
AND sql_audit_records.creator_id = :filter_creator_id
{{- end }}

{{- if .fuzzy_search_tags }}
AND sql_audit_records.tags LIKE '%{{ .fuzzy_search_tags }}%'
{{- end }}

{{- if .filter_task_status }}
AND tasks.status = :filter_task_status
{{- end }}

{{- if .filter_task_status_exclude }}
AND tasks.status <> :filter_task_status_exclude
{{- end }}

{{- if .filter_instance_id }}
AND  tasks.instance_id = :filter_instance_id
{{- end }}

{{- if .filter_create_time_from }}
AND sql_audit_records.created_at > :filter_create_time_from
{{- end }}

{{- if .filter_create_time_to }}
AND sql_audit_records.created_at < :filter_create_time_to
{{- end }}

{{- if .filter_audit_record_ids }}
AND sql_audit_records.audit_record_id IN ( {{ .filter_audit_record_ids }} )
{{- end }}

{{ end }}

`

func (s *Storage) GetSQLAuditRecordsByReq(data map[string]interface{}) (
	result []*SQLAuditRecordListItem, count uint64, err error) {

	err = s.getListResult(sqlAuditRecordQueryBodyTpl, sqlAuditRecordQueryTpl, data, &result)
	if err != nil {
		return result, 0, errors.New(errors.ConnectStorageError, err)
	}
	count, err = s.getCountResult(sqlAuditRecordQueryBodyTpl, sqlAuditRecordCountTpl, data)
	if err != nil {
		return result, 0, errors.New(errors.ConnectStorageError, err)
	}
	return result, count, err
}
