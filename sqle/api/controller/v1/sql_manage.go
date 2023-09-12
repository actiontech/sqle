package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetSqlManageListReq struct {
	FuzzySearchSqlFingerprint    *string `query:"fuzzy_search_sql_fingerprint" json:"fuzzy_search_sql_fingerprint,omitempty"`
	FuzzySearchAssignee          *string `query:"fuzzy_search_assignee" json:"fuzzy_search_assignee,omitempty"`
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
	Data      []*SqlManage `json:"data"`
	TotalNums uint64       `json:"total_nums"`
}

type SqlManage struct {
	Id                    uint64 `json:"id"`
	SqlFingerprint        string `json:"sql_fingerprint"`
	Sql                   string `json:"sql"`
	Source                string `json:"source"`
	Instance              string `json:"instance"`
	AuditResult           string `json:"audit_result"`
	FirstAppearTime       string `json:"first_appear_time"`
	LastAppearTime        string `json:"last_appear_time"`
	AppearNum             uint64 `json:"appear_num"`
	Assignee              string `json:"assignee"`
	Status                string `json:"status"`
	Remark                string `json:"remark"`
	SqlManageTotalNum     uint64 `json:"sql_manage_total_num"`
	SqlManageBadNum       uint64 `json:"sql_manage_bad_num"`
	SqlManageOptimizedNum uint64 `json:"sql_manage_optimized_num"`
}

// GetSqlManageList
// @Summary 获取管控sql列表
// @Description get sql manage list
// @Tags SqlManage
// @Security ApiKeyAuth
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param fuzzy_search_assignee query string false "fuzzy search assignee"
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
