package v2

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/labstack/echo/v4"
)

type GetSqlManageListResp struct {
	controller.BaseRes
	Data                  []*SqlManage `json:"data"`
	SqlManageTotalNum     uint64       `json:"sql_manage_total_num"`
	SqlManageBadNum       uint64       `json:"sql_manage_bad_num"`
	SqlManageOptimizedNum uint64       `json:"sql_manage_optimized_num"`
}

type SqlManage struct {
	Id                   uint64            `json:"id"`
	SqlFingerprint       string            `json:"sql_fingerprint"`
	Sql                  string            `json:"sql"`
	Source               *v1.Source        `json:"source"`
	InstanceName         string            `json:"instance_name"`
	SchemaName           string            `json:"schema_name"`
	AuditResult          []*v1.AuditResult `json:"audit_result"`
	AuditStatus          string            `json:"audit_status" enums:"being_audited"`
	FirstAppearTimeStamp string            `json:"first_appear_timestamp"`
	LastReceiveTimeStamp string            `json:"last_receive_timestamp"`
	FpCount              uint64            `json:"fp_count"`
	Assignees            []string          `json:"assignees"`
	Status               string            `json:"status" enums:"unhandled,solved,ignored,manual_audited,sent"`
	Remark               string            `json:"remark"`
	Endpoints            []string          `json:"endpoints"`
	Priority             string            `json:"priority"`
}

// @Deprecated
// GetSqlManageList
// @Summary 获取管控sql列表
// @Description get sql manage list
// @Tags SqlManage
// @Id GetSqlManageListV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param filter_assignee query string false "assignee"
// @Param filter_instance_id query string false "instance id"
// @Param filter_source query string false "source" Enums(audit_plan,sql_audit_record)
// @Param filter_audit_level query string false "audit level" Enums(normal,notice,warn,error)
// @Param filter_last_audit_start_time_from query string false "last audit start time from"
// @Param filter_last_audit_start_time_to query string false "last audit start time to"
// @Param filter_status query string false "status" Enums(unhandled,solved,ignored,manual_audited,sent)
// @Param filter_rule_name query string false "rule name"
// @Param filter_db_type query string false "db type"
// @Param filter_business query string false "filter by business" // This parameter is deprecated
// @Param filter_by_environment_tag query string false "filter by environment tag"
// @Param filter_priority query string false "priority" Enums(high,low)
// @Param fuzzy_search_endpoint query string false "fuzzy search endpoint"
// @Param fuzzy_search_schema_name query string false "fuzzy search schema name"
// @Param sort_field query string false "sort field" Enums(first_appear_timestamp,last_receive_timestamp,fp_count)
// @Param sort_order query string false "sort order" Enums(asc,desc)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v2.GetSqlManageListResp
// @Router /v2/projects/{project_name}/sql_manages [get]
func GetSqlManageList(c echo.Context) error {
	return nil
}

type ExportSqlManagesReq struct {
	FuzzySearchSqlFingerprint *string `query:"fuzzy_search_sql_fingerprint" json:"fuzzy_search_sql_fingerprint,omitempty"`
	FilterAssignee            *string `query:"filter_assignee" json:"filter_assignee,omitempty"`
	FilterByEnvironmentTag       *string `query:"filter_by_environment_tag" json:"filter_by_environment_tag,omitempty"`
	FilterInstanceID             *string `query:"filter_instance_id" json:"filter_instance_id,omitempty"`
	FilterSource                 *string `query:"filter_source" json:"filter_source,omitempty"`
	FilterAuditLevel             *string `query:"filter_audit_level" json:"filter_audit_level,omitempty"`
	FilterLastAuditStartTimeFrom *string `query:"filter_last_audit_start_time_from" json:"filter_last_audit_start_time_from,omitempty"`
	FilterLastAuditStartTimeTo   *string `query:"filter_last_audit_start_time_to" json:"filter_last_audit_start_time_to,omitempty"`
	FilterStatus                 *string `query:"filter_status" json:"filter_status,omitempty"`
	FilterDbType                 *string `query:"filter_db_type" json:"filter_db_type,omitempty"`
	FilterRuleName               *string `query:"filter_rule_name" json:"filter_rule_name,omitempty"`
	FilterPriority               *string `query:"filter_priority" json:"filter_priority,omitempty" enums:"high,low"`
	FuzzySearchEndpoint          *string `query:"fuzzy_search_endpoint" json:"fuzzy_search_endpoint,omitempty"`
	FuzzySearchSchemaName        *string `query:"fuzzy_search_schema_name" json:"fuzzy_search_schema_name,omitempty"`
	SortField                    *string `query:"sort_field" json:"sort_field,omitempty" valid:"omitempty,oneof=first_appear_timestamp last_receive_timestamp fp_count" enums:"first_appear_timestamp,last_receive_timestamp,fp_count"`
	SortOrder                    *string `query:"sort_order" json:"sort_order,omitempty" valid:"omitempty,oneof=asc desc" enums:"asc,desc"`
}

// ExportSqlManagesV2
// @Summary 导出SQL管控
// @Description export sql manage
// @Id exportSqlManageV2
// @Tags SqlManage
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param filter_assignee query string false "assignee"
// @Param filter_by_environment_tag query string false "filter by environment tag"
// @Param filter_priority query string false "priority" Enums(high,low)
// @Param filter_instance_id query string false "instance id"
// @Param filter_source query string false "source" Enums(audit_plan,sql_audit_record)
// @Param filter_audit_level query string false "audit level" Enums(normal,notice,warn,error)
// @Param filter_last_audit_start_time_from query string false "last audit start time from"
// @Param filter_last_audit_start_time_to query string false "last audit start time to"
// @Param filter_status query string false "status" Enums(unhandled,solved,ignored,manual_audited)
// @Param filter_db_type query string false "db type"
// @Param filter_rule_name query string false "rule name"
// @Param fuzzy_search_endpoint query string false "fuzzy search endpoint"
// @Param fuzzy_search_schema_name query string false "fuzzy search schema name"
// @Param sort_field query string false "sort field" Enums(first_appear_timestamp,last_receive_timestamp,fp_count)
// @Param sort_order query string false "sort order" Enums(asc,desc)
// @Success 200 {file} file "export sql manage"
// @Router /v2/projects/{project_name}/sql_manages/exports [get]
func ExportSqlManagesV2(c echo.Context) error {
	return exportSqlManagesV2(c)
}
