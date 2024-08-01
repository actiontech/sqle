package model

import (
	"database/sql"
	"time"
)

type InstanceAuditPlanListDetail struct {
	Id           uint           `json:"id"`
	DBType       string         `json:"db_type"`
	InstanceID   string         `json:"instance_id"`
	Token        string         `json:"token"`
	ActiveStatus string         `json:"active_status"`
	CreateUserId string         `json:"create_user_id"`
	CreateTime   string         `json:"created_at"`
	Types        sql.NullString `json:"types"`
}

var instanceAuditPlanQueryTpl = `
SELECT 
    instance_audit_plans.id,
    instance_audit_plans.db_type,
    instance_audit_plans.instance_id,
    instance_audit_plans.token,
    instance_audit_plans.active_status,
    instance_audit_plans.create_user_id,
    instance_audit_plans.created_at,
    audit_plans.types

 
{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var instanceAuditPlanCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var instanceAuditPlanBodyTpl = `
{{ define "body" }}

FROM 
    instance_audit_plans 
LEFT JOIN 
    (SELECT instance_audit_plan_id, GROUP_CONCAT(audit_plans_v2.type) AS types 
     FROM audit_plans_v2 
	 WHERE audit_plans_v2.deleted_at IS NULL
     GROUP BY instance_audit_plan_id) AS audit_plans 
ON 
    audit_plans.instance_audit_plan_id = instance_audit_plans.id
WHERE 
    instance_audit_plans.deleted_at IS NULL

{{- if not .current_user_is_admin }}
AND (
    instance_audit_plans.create_user_id = :current_user_id
    {{- if .accessible_instances_id }}
    OR instance_audit_plans.instance_id IN ( {{ .accessible_instances_id }} )
    {{- end }}
)
{{- end }}

{{- if .filter_instance_audit_plan_db_type }}
AND instance_audit_plans.db_type = :filter_instance_audit_plan_db_type
{{- end }}

{{- if .filter_audit_plan_type }}
AND FIND_IN_SET(:filter_audit_plan_type, audit_plans.types)
{{- end }}

{{- if .filter_audit_plan_instance_id }}
AND instance_audit_plans.instance_id = :filter_audit_plan_instance_id
{{- end }}

{{- if .filter_project_id }}
AND instance_audit_plans.project_id = :filter_project_id
{{- end }}

{{- if .filter_by_active_status }}
AND instance_audit_plans.active_status = :filter_by_active_status
{{- end }}
{{ end }}
`

func (s *Storage) GetInstanceAuditPlansByReq(data map[string]interface{}) (
	list []*InstanceAuditPlanListDetail, count uint64, err error) {
	err = s.getListResult(instanceAuditPlanBodyTpl, instanceAuditPlanQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(instanceAuditPlanBodyTpl, instanceAuditPlanCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}

type InstanceAuditPlanSQLListDetail struct {
	Fingerprint string         `json:"sql_fingerprint"`
	SQLContent  string         `json:"sql_text"`
	Schema      string         `json:"schema_name"`
	Info        JSON           `json:"info"`
	AuditResult sql.NullString `json:"audit_results"`
}

const (
	AuditResultName = "audit_results"
	AuditResultDesc = "审核结果"
)

var instanceAuditPlanSQLQueryTpl = `
SELECT
audit_plan_sqls.sql_fingerprint,
audit_plan_sqls.sql_text,
audit_plan_sqls.schema_name,
audit_plan_sqls.info,
audit_plan_sqls.audit_results

{{- template "body" . -}} 

order by audit_plan_sqls.id

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var instanceAuditPlanSQLCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var instanceAuditPlanSQLBodyTpl = `
{{ define "body" }}

FROM sql_manage_records AS audit_plan_sqls
JOIN audit_plans_v2 ON audit_plans_v2.id = audit_plan_sqls.source_id
JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND instance_audit_plans.deleted_at IS NULL
AND audit_plans_v2.id = :audit_plan_id

{{ end }}
`

func (s *Storage) GetInstanceAuditPlanSQLsByReq(data map[string]interface{}) (
	list []*InstanceAuditPlanSQLListDetail, count uint64, err error) {

	err = s.getListResult(instanceAuditPlanSQLBodyTpl, instanceAuditPlanSQLQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(instanceAuditPlanSQLBodyTpl, instanceAuditPlanSQLCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}

type SQLManagerList struct {
	Model

	SQLID                string       `json:"sql_id" gorm:"unique_index:sql_id;not null"`
	Source               string       `json:"source"`
	SourceId             string       `json:"source_id"`
	ProjectId            string       `json:"project_id"`
	InstanceID           string       `json:"instance_id"`
	SchemaName           string       `json:"schema_name"`
	SqlFingerprint       string       `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	SqlText              string       `json:"sql_text" gorm:"type:mediumtext;not null"`
	Info                 JSON         `gorm:"type:json"`
	AuditLevel           string       `json:"audit_level"`
	AuditResults         AuditResults `json:"audit_results" gorm:"type:json"`
	FpCount              uint64       `json:"fp_count"`
	FirstAppearTimestamp *time.Time   `json:"first_appear_timestamp"`
	LastReceiveTimestamp *time.Time   `json:"last_receive_timestamp"`

	// 任务属性字段
	Assignees string `json:"assignees"`
	Status    string `json:"status" gorm:"default:\"unhandled\""`
	Remark    string `json:"remark" gorm:"type:varchar(4000)"`
}
