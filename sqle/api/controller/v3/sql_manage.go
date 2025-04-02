package v3

import (
	v2 "github.com/actiontech/sqle/sqle/api/controller/v2"
	"github.com/labstack/echo/v4"
)

// GetSqlManageList
// @Summary 获取管控sql列表
// @Description get sql manage list
// @Tags SqlManage
// @Id GetSqlManageListV3
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
// @Param filter_by_environment_tag query string false "filter by environment tag"
// @Param filter_priority query string false "priority" Enums(high,low)
// @Param fuzzy_search_endpoint query string false "fuzzy search endpoint"
// @Param fuzzy_search_schema_name query string false "fuzzy search schema name"
// @Param sort_field query string false "sort field" Enums(first_appear_timestamp,last_receive_timestamp,fp_count)
// @Param sort_order query string false "sort order" Enums(asc,desc)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v2.GetSqlManageListResp
// @Router /v3/projects/{project_name}/sql_manages [get]
func GetSqlManageList(c echo.Context) error {
	return v2.GetSqlManageList(c)
}