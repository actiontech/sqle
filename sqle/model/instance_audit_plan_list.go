package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/locale"
)

type InstanceAuditPlanListDetail struct {
	Id           uint           `json:"id"`
	DBType       string         `json:"db_type"`
	InstanceID   string         `json:"instance_id"`
	Token        string         `json:"token"`
	ActiveStatus string         `json:"active_status"`
	CreateUserId string         `json:"create_user_id"`
	CreateTime   string         `json:"created_at"`
	AuditPlanIds sql.NullString `json:"audit_plan_ids"`
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
	instance_audit_plans.token,
    audit_plans.audit_plan_ids

 
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
    (SELECT instance_audit_plan_id, 
		GROUP_CONCAT(audit_plans_v2.id) AS audit_plan_ids,
		GROUP_CONCAT(audit_plans_v2.type) AS types
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

{{- if .filter_instance_ids_by_env }}
AND instance_audit_plans.instance_id in ( {{ .filter_instance_ids_by_env}} )
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
	AuditPlanSqlId string         `json:"id"`
	Fingerprint    string         `json:"sql_fingerprint"`
	SQLContent     string         `json:"sql_text"`
	Schema         string         `json:"schema_name"`
	Info           JSON           `json:"info"`
	AuditResult    AuditResults   `json:"audit_results"`
	AuditStatus    string         `json:"audit_status"`
	Priority       sql.NullString `json:"priority"`
}

const (
	AuditResultName = "audit_results"
	AuditStatus     = "audit_status"
)

var AuditResultDesc = locale.ApAuditResult

type FilterName string

const (
	FilterNameDBUser               FilterName = "db_user"
	FilterLastReceiveTimestampFrom FilterName = "last_receive_timestamp_from"
	FilterLastReceiveTimestampTo   FilterName = "last_receive_timestamp_to"
	FilterRuleName                 FilterName = "rule_name"
	FilterSchemaName               FilterName = "schema_name"
	FilterSQL                      FilterName = "sql"
	FilterCounter                  FilterName = "counter"
	FilterQueryTimeAvg             FilterName = "query_time_avg"
	FilterRowExaminedAvg           FilterName = "row_examined_avg"
	FilterPriority                 FilterName = "priority"
)

type FilterType string

const (
	FilterTypeCommon FilterType = "common"
	FilterTypeLike   FilterType = "like"
)

var FilterMap = map[FilterName]FilterType{
	FilterNameDBUser:               FilterTypeCommon,
	FilterLastReceiveTimestampFrom: FilterTypeCommon,
	FilterLastReceiveTimestampTo:   FilterTypeCommon,
	FilterRuleName:                 FilterTypeCommon,
	FilterSchemaName:               FilterTypeCommon,
	FilterCounter:                  FilterTypeCommon,
	FilterSQL:                      FilterTypeLike,
	FilterQueryTimeAvg:             FilterTypeCommon,
	FilterRowExaminedAvg:           FilterTypeCommon,
	FilterPriority:                 FilterTypeCommon,
}

var OrderByMap = map[string] /* field name */ string /* field name with table*/ {
	"counter":                 "audit_plan_sqls.info->'$.counter'",
	"query_time_total":        "audit_plan_sqls.info->'$.query_time_total'",
	"cpu_time_total":          "audit_plan_sqls.info->'$.cpu_time_total'",
	"disk_read_total":         "audit_plan_sqls.info->'$.disk_read_total'",
	"buffer_read_total":       "audit_plan_sqls.info->'$.buffer_read_total'",
	"user_io_wait_time_total": "audit_plan_sqls.info->'$.user_io_wait_time_total'",
	"last_receive_timestamp":  "audit_plan_sqls.info->'$.last_receive_timestamp'",
	"query_time_avg":          "audit_plan_sqls.info->'$.query_time_avg'",
	"query_time_max":          "audit_plan_sqls.info->'$.query_time_max'",
	"row_examined_avg":        "audit_plan_sqls.info->'$.row_examined_avg'",
	"cpu_time_avg":            "audit_plan_sqls.info->'$.cpu_time_avg'",
	"lock_wait_time_total":    "audit_plan_sqls.info->'$.lock_wait_time_total'",
	"lock_wait_counter":       "audit_plan_sqls.info->'$.lock_wait_counter'",
	"act_wait_time_total":     "audit_plan_sqls.info->'$.act_wait_time_total'",
	"act_time_total":          "audit_plan_sqls.info->'$.act_time_total'",
	"phy_read_page_total":     "audit_plan_sqls.info->'$.phy_read_page_total'",
	"logic_read_page_total":   "audit_plan_sqls.info->'$.logic_read_page_total'",
}

var instanceAuditPlanSQLQueryTpl = `
SELECT
audit_plan_sqls.id,
audit_plan_sqls.sql_fingerprint,
audit_plan_sqls.sql_text,
audit_plan_sqls.schema_name,
audit_plan_sqls.info,
audit_plan_sqls.audit_results,
IF(audit_plan_sqls.audit_results IS NULL,'being_audited','') AS audit_status,
audit_plan_sqls.priority

{{- template "body" . -}} 

{{- if .order_by -}}
ORDER BY {{.order_by}}
{{- if .is_asc }}
ASC
{{- else}}
DESC
{{- end -}}
{{else}}
ORDER BY audit_plan_sqls.id
{{- end -}}

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
JOIN audit_plans_v2 ON CONCAT(audit_plans_v2.instance_audit_plan_id, '') = audit_plan_sqls.source_id AND audit_plans_v2.type = audit_plan_sqls.source
JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id

WHERE audit_plan_sqls.deleted_at IS NULL
AND instance_audit_plans.deleted_at IS NULL
AND audit_plans_v2.id = :audit_plan_id

{{- if .schema_name }}
AND audit_plan_sqls.schema_name = :schema_name
{{- end }}

{{- if .last_receive_timestamp_from }}
AND JSON_EXTRACT(audit_plan_sqls.info, '$.last_receive_timestamp') >= :last_receive_timestamp_from
{{- end }}

{{- if .last_receive_timestamp_to }}
AND JSON_EXTRACT(audit_plan_sqls.info, '$.last_receive_timestamp') <= :last_receive_timestamp_to
{{- end }}

{{- if .rule_name }}
AND JSON_CONTAINS(JSON_EXTRACT(audit_plan_sqls.audit_results,'$[*].rule_name'), '"{{ .rule_name }}"') > 0
{{- end }}

{{- if .db_user }}
AND JSON_EXTRACT(audit_plan_sqls.info, '$.db_user') = :db_user
{{- end }}

{{- if .schema_meta_name }}
AND JSON_EXTRACT(audit_plan_sqls.info, '$.schema_meta_name') = :schema_meta_name
{{- end }}

{{- if .sql }}
AND audit_plan_sqls.sql_fingerprint LIKE :sql
{{- end}}

{{- if .counter }}
AND JSON_EXTRACT(audit_plan_sqls.info, '$.counter') >= :counter
{{- end}}

{{- if .query_time_avg }}
AND JSON_EXTRACT(audit_plan_sqls.info, '$.query_time_avg') >= :query_time_avg
{{- end}}

{{- if .row_examined_avg }}
AND JSON_EXTRACT(audit_plan_sqls.info, '$.row_examined_avg') >= :row_examined_avg
{{- end}}


{{- if .priority }}
AND audit_plan_sqls.priority = :priority
{{- end}}

{{ end }}
`

// todo: 删除
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

// todo: 等 GetInstanceAuditPlanSQLsByReq 方法删除后替换它
func (s *Storage) GetInstanceAuditPlanSQLsByReqV2(apId uint, apType string, limit, offset int, orderBy string, isAsc bool, filters map[FilterName]interface{}) (
	list []*InstanceAuditPlanSQLListDetail, count uint64, err error) {
	args := map[string]interface{}{}

	args["audit_plan_id"] = apId
	args["audit_plan_type"] = apType
	args["is_asc"] = isAsc
	if limit > 0 {
		args["limit"] = limit
		args["offset"] = offset
	}

	// order by 无法预编译，使用map定义预期的排序字段，非预期值则跳过防止SQL注入
	if v, ok := OrderByMap[orderBy]; ok {
		args["order_by"] = v
	}

	for filterName, filterValue := range filters {
		v, ok := FilterMap[filterName]
		// 非预期的筛选条件跳过
		if !ok {
			continue
		}
		switch v {
		case FilterTypeLike:
			args[string(filterName)] = fmt.Sprintf("%%%s%%", filterValue)
		default:
			args[string(filterName)] = filterValue
		}

	}
	err = s.getListResult(instanceAuditPlanSQLBodyTpl, instanceAuditPlanSQLQueryTpl, args, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(instanceAuditPlanSQLBodyTpl, instanceAuditPlanSQLCountTpl, args)
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
