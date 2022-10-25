package v2

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)


type GetAuditPlanReportSQLsReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanReportSQLsResV1 struct {
	controller.BaseRes
	Data      []AuditPlanReportSQLResV1 `json:"data"`
	TotalNums uint64                    `json:"total_nums"`
}

type AuditPlanReportSQLResV1 struct {
	SQL         string `json:"audit_plan_report_sql" example:"select * from t1 where id = 1"`
	AuditResult string `json:"audit_plan_report_sql_audit_result" example:"same format as task audit result"`
	Number      uint   `json:"number" example:"1"`
}

// GetAuditPlanReportSQLsV1 is to fix the irregular uri used by GetAuditPlanReportSQLs
// issue: https://github.com/actiontech/sqle/issues/429
// @Summary 获取指定扫描任务的SQL扫描详情
// @Description get audit plan report SQLs
// @Id getAuditPlanReportsSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_id path uint true "project id"
// @Param audit_plan_name path string true "audit plan name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanReportSQLsResV1
// @router /v1/projects/{project_id}/audit_plans/{audit_plan_name}/reports/{audit_plan_report_id}/sqls [get]
func GetAuditPlanReportSQLsV1(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlanReportSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	apName := c.Param("audit_plan_name")
	err := v1.CheckCurrentUserCanAccessAuditPlan(c, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"audit_plan_name":      apName,
		"audit_plan_report_id": c.Param("audit_plan_report_id"),
		"limit":                req.PageSize,
		"offset":               offset,
	}
	auditPlanReportSQLs, count, err := s.GetAuditPlanReportSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlanReportSQLsResV1 := make([]AuditPlanReportSQLResV1, len(auditPlanReportSQLs))
	for i, auditPlanReportSQL := range auditPlanReportSQLs {
		auditPlanReportSQLsResV1[i] = AuditPlanReportSQLResV1{
			SQL:         auditPlanReportSQL.SQL,
			AuditResult: auditPlanReportSQL.AuditResult,
			Number:      auditPlanReportSQL.Number,
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanReportSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlanReportSQLsResV1,
		TotalNums: count,
	})
}
