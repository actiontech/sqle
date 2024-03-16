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
	FirstAppearTimeStamp string            `json:"first_appear_timestamp"`
	LastReceiveTimeStamp string            `json:"last_receive_timestamp"`
	FpCount              uint64            `json:"fp_count"`
	Assignees            []string          `json:"assignees"`
	Status               string            `json:"status" enums:"unhandled,solved,ignored,manual_audited"`
	Remark               string            `json:"remark"`
	Endpoints            []string          `json:"endpoints"`
}

// GetSqlManageList
// @Summary 获取管控sql列表
// @Description get sql manage list
// @Tags SqlManage
// @Id GetSqlManageListV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param filter_assignee query string false "assignee"
// @Param filter_instance_name query string false "instance name"
// @Param filter_source query string false "source" Enums(audit_plan,sql_audit_record)
// @Param filter_audit_level query string false "audit level" Enums(normal,notice,warn,error)
// @Param filter_last_audit_start_time_from query string false "last audit start time from"
// @Param filter_last_audit_start_time_to query string false "last audit start time to"
// @Param filter_status query string false "status" Enums(unhandled,solved,ignored,manual_audited)
// @Param filter_rule_name query string false "rule name"
// @Param filter_db_type query string false "db type"
// @Param fuzzy_search_endpoint query string false "fuzzy search endpoint"
// @Param fuzzy_search_schema_name query string false "fuzzy search schema name"
// @Param sort_field query string false "sort field" Enums(first_appear_timestamp,last_receive_timestamp,fp_count)
// @Param sort_order query string false "sort order" Enums(asc,desc)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v2.GetSqlManageListResp
// @Router /v2/projects/{project_name}/sql_manages [get]
func GetSqlManageList(c echo.Context) error {
	return getSqlManageList(c)
}
