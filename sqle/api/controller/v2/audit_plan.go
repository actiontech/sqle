package v2

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetAuditPlanSQLsResV2 struct {
	controller.BaseRes
	Data      AuditPlanSQLResV2 `json:"data"`
	TotalNums uint64            `json:"total_nums"`
}

type AuditPlanSQLResV2 struct {
	Head []AuditPlanSQLHeadV2                 `json:"head"`
	Rows []map[string] /* head name */ string `json:"rows"`
}

type AuditPlanSQLHeadV2 struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Type string `json:"type" enums:"sql"`
}

// @Summary 获取指定审核计划的SQLs信息(不包括审核结果)
// @Description get audit plan SQLs
// @Id getAuditPlanSQLsV2
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v2.GetAuditPlanSQLsResV2
// @router /v2/audit_plans/{audit_plan_name}/sqls [get]
func GetAuditPlanSQLs(c echo.Context) error {
	return c.JSON(http.StatusOK, &GetAuditPlanSQLsResV2{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      AuditPlanSQLResV2{},
		TotalNums: 0,
	})
}

type GetAuditPlanReportSQLsReqV2 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanReportSQLsResV2 struct {
	controller.BaseRes
	Data      []AuditPlanReportSQLResV2 `json:"data"`
	TotalNums uint64                    `json:"total_nums"`
}

type AuditPlanReportSQLResV2 struct {
	SQL         string `json:"audit_plan_report_sql" example:"select * from t1 where id = 1"`
	AuditResult string `json:"audit_plan_report_sql_audit_result" example:"same format as task audit result"`
}

// @Summary 获取指定审核计划的SQL审核详情
// @Description get audit plan report SQLs
// @Id getAuditPlanReportSQLsV2
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v2.GetAuditPlanReportSQLsResV2
// @router /v2/audit_plans/{audit_plan_name}/report/{audit_plan_report_id}/ [get]
func GetAuditPlanReportSQLs(c echo.Context) error {
	return c.JSON(http.StatusOK, &GetAuditPlanReportSQLsResV2{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      []AuditPlanReportSQLResV2{},
		TotalNums: 0,
	})
}
