package model

import (
	"database/sql"

	"github.com/actiontech/sqle/sqle/pkg/params"
)

type AuditPlanListDetail struct {
	Name             string         `json:"name"`
	Cron             string         `json:"cron_expression"`
	DBType           string         `json:"db_type"`
	Token            string         `json:"token"`
	InstanceName     string         `json:"instance_name"`
	InstanceDatabase string         `json:"instance_database"`
	Type             sql.NullString `json:"type"`
	Params           params.Params  `json:"params"`
}

var auditPlanQueryTpl = `
SELECT audit_plans.name, audit_plans.cron_expression, audit_plans.db_type, audit_plans.token,
audit_plans.instance_name, audit_plans.instance_database, audit_plans.type, audit_plans.params

{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var auditPlanCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var auditPlanBodyTpl = `
{{ define "body" }}
FROM audit_plans
LEFT JOIN users ON audit_plans.create_user_id = users.id

WHERE audit_plans.deleted_at IS NULL

{{- if not .current_user_is_admin }}
AND users.login_name = :current_user_id
{{- end }}

{{- if .filter_audit_plan_db_type }}
AND audit_plans.db_type = :filter_audit_plan_db_type
{{- end }}

{{ end }}
`

func (s *Storage) GetAuditPlansByReq(data map[string]interface{}) (
	list []*AuditPlanListDetail, count uint64, err error) {

	err = s.getListResult(auditPlanBodyTpl, auditPlanQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(auditPlanBodyTpl, auditPlanCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}

type AuditPlanSQLListDetail struct {
	Fingerprint string `json:"fingerprint"`
	SQLContent  string `json:"sql_content"`
	Info        JSON   `json:"info"`
}

var auditPlanSQLQueryTpl = `
SELECT
audit_plan_sqls.fingerprint,
audit_plan_sqls.sql_content,
audit_plan_sqls.info

{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var auditPlanSQLCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var auditPlanSQLBodyTpl = `
{{ define "body" }}

FROM audit_plan_sqls_v2 AS audit_plan_sqls
JOIN audit_plans ON audit_plans.id = audit_plan_sqls.audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND audit_plans.deleted_at IS NULL
AND audit_plans.name = :audit_plan_name

{{ end }}
`

func (s *Storage) GetAuditPlanSQLsByReq(data map[string]interface{}) (
	list []*AuditPlanSQLListDetail, count uint64, err error) {

	err = s.getListResult(auditPlanSQLBodyTpl, auditPlanSQLQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(auditPlanSQLBodyTpl, auditPlanSQLCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}
