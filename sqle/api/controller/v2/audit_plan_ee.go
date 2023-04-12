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
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

func getAuditPlanAnalysisData(c echo.Context) error {
	reportId := c.Param("audit_plan_report_id")
	sqlNumber := c.Param("number")
	apName := c.Param("audit_plan_name")
	projectName := c.Param("project_name")

	if err := v1.CheckIsProjectMember(controller.GetUserName(c), projectName); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ap, exist, err := v1.GetAuditPlanIfCurrentUserCanAccess(c, projectName, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewAuditPlanNotExistErr())
	}

	reportIdInt, err := strconv.Atoi(reportId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("parse audit plan report id failed: %v", err)))
	}

	sqlNumberInt, err := strconv.Atoi(sqlNumber)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("parse number failed: %v", err)))
	}

	s := model.GetStorage()
	auditPlanReport, exist, err := s.GetAuditPlanReportByID(ap.ID, uint(reportIdInt))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("audit plan report not exist")))
	}

	auditPlanReportSQLV2, exist, err := s.GetAuditPlanReportSQLV2ByReportIDAndNumber(uint(reportIdInt), uint(sqlNumberInt))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("audit plan report sql v2 not exist")))
	}

	instance, exist, err := s.GetInstanceByNameAndProjectID(auditPlanReport.AuditPlan.InstanceName, auditPlanReport.AuditPlan.ProjectId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance not exist")))
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
