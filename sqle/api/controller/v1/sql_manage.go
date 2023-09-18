package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetSqlManageListReq struct {
	FuzzySearchSqlFingerprint    *string `query:"fuzzy_search_sql_fingerprint" json:"fuzzy_search_sql_fingerprint,omitempty"`
	FilterAssignee               *string `query:"filter_assignee" json:"filter_assignee,omitempty"`
	FilterInstanceName           *string `query:"filter_instance_name" json:"filter_instance_name,omitempty"`
	FilterSource                 *string `query:"filter_source" json:"filter_source,omitempty"`
	FilterAuditLevel             *string `query:"filter_audit_level" json:"filter_audit_level,omitempty"`
	FilterLastAuditStartTimeFrom *string `query:"filter_last_audit_start_time_from" json:"filter_last_audit_start_time_from,omitempty"`
	FilterLastAuditStartTimeTo   *string `query:"filter_last_audit_start_time_to" json:"filter_last_audit_start_time_to,omitempty"`
	FilterStatus                 *string `query:"filter_status" json:"filter_status,omitempty"`
	PageIndex                    uint32  `query:"page_index" valid:"required" json:"page_index"`
	PageSize                     uint32  `query:"page_size" valid:"required" json:"page_size"`
}

type GetSqlManageListResp struct {
	controller.BaseRes
	Data                  []*SqlManage `json:"data"`
	SqlManageTotalNum     uint64       `json:"sql_manage_total_num"`
	SqlManageBadNum       uint64       `json:"sql_manage_bad_num"`
	SqlManageOptimizedNum uint64       `json:"sql_manage_optimized_num"`
}

type SqlManage struct {
	Id              uint64         `json:"id"`
	SqlFingerprint  string         `json:"sql_fingerprint"`
	Sql             string         `json:"sql"`
	Source          *Source        `json:"source"`
	InstanceName    string         `json:"instance_name"`
	SchemaName      string         `json:"schema_name"`
	AuditResult     []*AuditResult `json:"audit_result"`
	FirstAppearTime string         `json:"first_appear_time"`
	LastAppearTime  string         `json:"last_appear_time"`
	AppearNum       uint64         `json:"appear_num"`
	Assignees       []string       `json:"assignees"`
	Status          string         `json:"status" enums:"unhandled,solved,ignored"`
	Remark          string         `json:"remark"`
}

// AuditResult 用于SQL全生命周期展示的AuditResult
type AuditResult struct {
	Level    string `json:"level" example:"warn"`
	Message  string `json:"message" example:"避免使用不必要的内置函数md5()"`
	RuleName string `json:"rule_name"`
}

type Source struct {
	Type             string `json:"type" enums:"audit_plan,sql_audit_record"`
	AuditPlanName    string `json:"audit_plan_name"`
	SqlAuditRecordId string `json:"sql_audit_record_id"`
}

// GetSqlManageList
// @Summary 获取管控sql列表
// @Description get sql manage list
// @Tags SqlManage
// @Id GetSqlManageList
// @Security ApiKeyAuth
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param filter_assignee query string false "assignee"
// @Param filter_instance_name query string false "instance name"
// @Param filter_source query string false "source" Enums(audit_plan,api_audit)
// @Param filter_audit_level query string false "audit level" Enums(normal,notice,warn,error)
// @Param filter_last_audit_start_time_from query string false "last audit start time from"
// @Param filter_last_audit_start_time_to query string false "last audit start time to"
// @Param filter_status query string false "status" Enums(unhandled,solved,ignored)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetSqlManageListResp
// @Router /v1/projects/{project_name}/sql_manages [get]
func GetSqlManageList(c echo.Context) error {
	return nil
}
