//go:build enterprise
// +build enterprise

package v2

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/labstack/echo/v4"
)

func getAuditPlanAnalysisData(c echo.Context) error {
	reportId := c.Param("audit_plan_report_id")
	sqlNumber := c.Param("number")
	apName := c.Param("audit_plan_name")
	projectName := c.Param("project_name")

	reportIdInt, err := strconv.Atoi(reportId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("parse audit plan report id failed: %v", err)))
	}

	sqlNumberInt, err := strconv.Atoi(sqlNumber)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("parse number failed: %v", err)))
	}

	auditPlanReport, auditPlanReportSQLV2, instance, err := v1.GetAuditPlantReportAndInstance(c, projectName, apName, reportIdInt, sqlNumberInt)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	res, err := getSQLAnalysisResult(log.NewEntry(), instance, auditPlanReport.AuditPlan.InstanceDatabase, auditPlanReportSQLV2.SQL)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditPlanAnalysisDataResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertSQLAnalysisResultToRes(res, auditPlanReportSQLV2.SQL),
	})
}
