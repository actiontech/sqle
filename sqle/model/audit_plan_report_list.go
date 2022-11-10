package model

import "database/sql"

type AuditPlanReportListDetail struct {
	ID         string          `json:"id"`
	AuditLevel sql.NullString  `json:"audit_level"`
	Score      sql.NullInt32   `json:"score"`
	PassRate   sql.NullFloat64 `json:"pass_rate"`
	CreateAt   string          `json:"created_at"`
}

var auditPlanReportQueryTpl = `
SELECT reports.id, reports.score , reports.pass_rate, reports.audit_level, reports.created_at

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

FROM audit_plan_reports_v2 AS reports
JOIN audit_plans ON audit_plans.id = reports.audit_plan_id
JOIN projects ON projects.id = audit_plans.project_id

WHERE reports.deleted_at IS NULL
AND audit_plans.deleted_at IS NULL
AND audit_plans.name = :audit_plan_name
AND projects.name = :project_name

ORDER BY reports.created_at DESC , reports.id DESC 

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
	SQL         string `json:"sql"`
	AuditResult string `json:"audit_result"`
	Number      uint   `json:"number"`
}

var auditPlanReportSQLQueryTpl = `
SELECT report_sqls.sql, report_sqls.audit_result, report_sqls.number

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

FROM audit_plan_report_sqls_v2 AS report_sqls
WHERE deleted_at IS NULL
AND audit_plan_report_id = :audit_plan_report_id

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
