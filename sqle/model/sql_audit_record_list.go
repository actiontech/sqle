package model

import (
	"database/sql"
	"time"
)

type SQLAuditRecordListItem struct {
	AuditRecordId   string          `json:"audit_record_id"`
	RecordCreatedAt *time.Time      `json:"record_created_at"`
	CreatorName     string          `json:"creator_name"`
	Tags            sql.NullString  `json:"tags"`
	InstanceName    sql.NullString  `json:"instance_name"`
	InstanceHost    sql.NullString  `json:"instance_host"`
	InstancePort    sql.NullString  `json:"instance_port"`
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
       create_user.login_name            AS creator_name,
       instances.name                    AS instance_name,
       instances.db_host                 AS instance_host,
       instances.db_port                 AS instance_port,
       tasks.id                          AS task_id,
       tasks.db_type                     AS db_type,
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

var sqlAuditRecordQueryBodyTpl = `
{{ define "body" }}
FROM sql_audit_records
LEFT JOIN projects AS p ON sql_audit_records.project_id = p.id
LEFT JOIN users AS create_user ON sql_audit_records.creator_id = create_user.id
LEFT JOIN tasks ON sql_audit_records.task_id = tasks.id
LEFT JOIN instances ON tasks.instance_id = instances.id

WHERE
sql_audit_records.deleted_at IS NULL
AND p.name = :filter_project_name

{{- if .check_user_can_access }}
AND create_user.id = :filter_creator_id
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

{{- if .filter_instance_name }}
AND instances.name = :filter_instance_name
{{- end }}

{{- if .filter_create_time_from }}
AND sql_audit_records.created_at > :filter_create_time_from
{{- end }}

{{- if .filter_create_time_to }}
AND sql_audit_records.created_at < :filter_create_time_to
{{- end }}

{{ end }}

`

func (s *Storage) GetSQLAuditRecordsByReq(data map[string]interface{}) (
	result []*SQLAuditRecordListItem, err error) {

	err = s.getListResult(sqlAuditRecordQueryBodyTpl, sqlAuditRecordQueryTpl, data, &result)
	if err != nil {
		return result, err
	}

	return result, err
}
