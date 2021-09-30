package model

type AuditPlanListDetail struct {
	Name             string `json:"name"`
	Cron             string `json:"cron_expression"`
	DBType           string `json:"db_type"`
	Token            string `json:"token"`
	InstanceName     string `json:"instance_name"`
	InstanceDatabase string `json:"instance_database"`
}

var auditPlanQueryTpl = `
SELECT audit_plans.name, audit_plans.cron_expression, audit_plans.db_type, audit_plans.token,
audit_plans.instance_name, audit_plans.instance_database

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

WHERE audit_plans.deleted_at IS NULL

{{- if not .current_user_is_admin }}
AND audit_plans.create_user_id = :current_user_id
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
	Fingerprint          string `json:"fingerprint"`
	Counter              string `json:"counter"`
	LastReceiveText      string `json:"last_sql"`
	LastReceiveTimestamp string `json:"last_receive_timestamp"`
}

var auditPlanSQLQueryTpl = `
SELECT audit_plan_sqls.fingerprint, audit_plan_sqls.counter, audit_plan_sqls.last_sql, audit_plan_sqls.last_receive_timestamp

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

FROM audit_plan_sqls
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

type AuditPlanReportListDetail struct {
	ID       string `json:"id"`
	CreateAt string `json:"created_at"`
}

var auditPlanReportQueryTpl = `
SELECT audit_plan_reports.id, audit_plan_reports.created_at

{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var auditPlanReportCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var auditPlanReportBodyTpl = `
{{ define "body" }}

FROM audit_plan_reports
JOIN audit_plans ON audit_plans.id = audit_plan_reports.audit_plan_id

WHERE audit_plan_reports.deleted_at IS NULL
AND audit_plans.deleted_at IS NULL
AND audit_plans.name = :audit_plan_name

{{ end }}
`

func (s *Storage) GetAuditPlanReportsByReq(data map[string]interface{}) (
	list []*AuditPlanReportListDetail, count uint64, err error) {

	err = s.getListResult(auditPlanReportBodyTpl, auditPlanReportQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(auditPlanReportBodyTpl, auditPlanReportCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}

type AuditPlanReportSQLListDetail struct {
	AuditResult string `json:"audit_result"`

	Fingerprint          string `json:"fingerprint"`
	LastReceiveText      string `json:"last_sql"`
	LastReceiveTimestamp string `json:"last_receive_timestamp"`
}

var auditPlanReportSQLQueryTpl = `
SELECT audit_plan_report_sqls.audit_result, 
audit_plan_sqls.fingerprint, audit_plan_sqls.last_sql, audit_plan_sqls.last_receive_timestamp

{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var auditPlanReportSQLCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var auditPlanReportSQLBodyTpl = `
{{ define "body" }}

FROM audit_plan_report_sqls
JOIN audit_plan_reports ON audit_plan_report_sqls.audit_plan_report_id = audit_plan_reports.id
JOIN audit_plans ON audit_plan_reports.audit_plan_id = audit_plans.id
JOIN audit_plan_sqls ON audit_plan_sqls.id = audit_plan_report_sqls.audit_plan_sql_id

WHERE audit_plan_report_sqls.deleted_at IS NULL
AND audit_plans.name = :audit_plan_name
AND audit_plan_reports.id = :audit_plan_report_id

{{ end }}
`

func (s *Storage) GetAuditPlanReportSQLsByReq(data map[string]interface{}) (
	list []*AuditPlanReportSQLListDetail, count uint64, err error) {

	err = s.getListResult(auditPlanReportSQLBodyTpl, auditPlanReportSQLQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(auditPlanReportSQLBodyTpl, auditPlanReportSQLCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}
