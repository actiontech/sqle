package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetSqlManageListReq struct {
	FilterDataSource             *string `form:"filter_data_source"`
	FilterSource                 *string `form:"filter_source"`
	FilterAuditLevel             *string `form:"filter_audit_level"`
	FilterOperationPerson        *string `form:"filter_operation_person"`
	FilterLastAuditStartTimeFrom *string `form:"filter_last_audit_start_time_from"`
	FilterLastAuditStartTimeTo   *string `form:"filter_last_audit_start_time_to"`
	FilterIsRelatedToMe          *bool   `form:"filter_is_related_to_me"`
	FilterIsSolved               *bool   `form:"filter_is_solved"`
	PageIndex                    uint32  `form:"page_index" valid:"required"`
	PageSize                     uint32  `form:"page_size" valid:"required"`
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
	DataSource            string `json:"data_source"`
	AuditResult           string `json:"audit_result"`
	FirstAppearTime       string `json:"first_appear_time"`
	LastAppearTime        string `json:"last_appear_time"`
	AppearNum             uint64 `json:"appear_num"`
	AssignPerson          string `json:"assign_person"`
	SolveStatus           string `json:"solve_status"`
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
// @Param filter_data_source query string false "data source"
// @Param filter_source query string false "source" Enums(audit_plan,api_audit)
// @Param filter_audit_level query string false "audit level" Enums(normal,notice,warn,error)
// @Param filter_operation_person query string false "operation person"
// @Param filter_last_audit_start_time_from query string false "last audit start time from"
// @Param filter_last_audit_start_time_to query string false "last audit start time to"
// @Param filter_is_related_to_me query boolean false "is related to me"
// @Param filter_is_solved query boolean false "is solved"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetSqlManageListResp
// @Router /v1/projects/{project_name}/sql_manages [get]
func GetSqlManageList(c echo.Context) error {
	return nil
}
