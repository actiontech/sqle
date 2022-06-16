//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func getAuditPlanAnalysisData(c echo.Context) error {
	reportId := c.Param("audit_plan_report_id")
	number := c.Param("number")

	reportIdInt, err := strconv.Atoi(reportId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("parse audit plan report id failed: %v", err)))
	}

	numberInt, err := strconv.Atoi(number)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("parse number failed: %v", err)))
	}

	s := model.GetStorage()
	auditPlanReport, exist, err := s.GetAuditPlanReportByID(uint(reportIdInt))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("audit plan report not exist")))
	}

	err = CheckCurrentUserCanAccessAuditPlan(c, auditPlanReport.AuditPlan.Name, model.OP_AUDIT_PLAN_VIEW_OTHERS)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlanReportSQLV2, exist, err := s.GetAuditPlanReportSQLV2ByReportIDAndNumber(uint(reportIdInt), uint(numberInt))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("audit plan report sql v2 not exist")))
	}

	instance, exist, err := s.GetInstanceByName(auditPlanReport.AuditPlan.InstanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance not exist")))
	}

	dsn, err := newDSN(instance, auditPlanReport.AuditPlan.InstanceDatabase)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	analysisDriver, err := driver.NewAnalysisDriver(log.NewEntry(), instance.DbType, dsn)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	explainResult, err := analysisDriver.Explain(context.TODO(), &driver.ExplainConf{Sql: auditPlanReportSQLV2.SQL})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	//todo:remove NewAnalysisDriver function later
	analysisDriver01, err := driver.NewAnalysisDriver(log.NewEntry(), instance.DbType, dsn)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	tableMetaResult, err := analysisDriver01.GetTableMetaBySQL(context.TODO(), &driver.GetTableMetaBySQLConf{Sql: auditPlanReportSQLV2.SQL})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditPlanAnalysisDataResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    explainAndMetaDataToRes(explainResult, tableMetaResult, auditPlanReportSQLV2.SQL),
	})
}

func explainAndMetaDataToRes(explainResult *driver.ExplainResult, metaDataResult *driver.GetTableMetaBySQLResult,
	rawSql string) GetSQLAnalysisDataResItemV1 {
	analysisDataResItemV1 := GetSQLAnalysisDataResItemV1{
		SQLExplain: SQLExplain{
			ClassicResult: ExplainClassicResult{
				Rows: make([]map[string]string, len(explainResult.ClassicResult.Rows)),
				Head: make([]TableMetaItemHeadResV1, len(explainResult.ClassicResult.Column)),
			},
		},
		TableMetas: make([]TableMeta, len(metaDataResult.TableMetas)),
	}

	explainResItemV1 := analysisDataResItemV1.SQLExplain.ClassicResult
	for i, column := range explainResult.ClassicResult.Column {
		explainResItemV1.Head[i].FieldName = column.Name
		explainResItemV1.Head[i].Desc = column.Desc
	}

	for i, rows := range explainResult.ClassicResult.Rows {
		explainResItemV1.Rows[i] = make(map[string]string)
		for k, row := range rows {
			columnName := explainResult.ClassicResult.Column[k].Name
			explainResItemV1.Rows[i][columnName] = row
		}
	}

	analysisDataResItemV1.SQLExplain.SQL = rawSql

	for i, tableMeta := range metaDataResult.TableMetas {
		tableMetaColumnsInfo := tableMeta.ColumnsInfo
		tableMetaIndexInfo := tableMeta.IndexesInfo

		analysisDataResItemV1.TableMetas[i].Columns = TableColumns{
			Rows: make([]map[string]string, len(tableMetaColumnsInfo.Rows)),
			Head: make([]TableMetaItemHeadResV1, len(tableMetaColumnsInfo.Column)),
		}
		analysisDataResItemV1.TableMetas[i].Indexes = TableIndexes{
			Rows: make([]map[string]string, len(tableMetaIndexInfo.Rows)),
			Head: make([]TableMetaItemHeadResV1, len(tableMetaIndexInfo.Column)),
		}

		tableMetaColumnRes := analysisDataResItemV1.TableMetas[i].Columns
		for i2, column := range tableMetaColumnsInfo.Column {
			tableMetaColumnRes.Head[i2].FieldName = column.Name
			tableMetaColumnRes.Head[i2].Desc = column.Desc
		}

		for i2, rows := range tableMetaColumnsInfo.Rows {
			tableMetaColumnRes.Rows[i2] = make(map[string]string)
			for k, row := range rows {
				columnName := tableMetaColumnsInfo.Column[k].Name
				tableMetaColumnRes.Rows[i2][columnName] = row
			}
		}

		tableMetaIndexRes := analysisDataResItemV1.TableMetas[i].Indexes
		for i2, column := range tableMetaIndexInfo.Column {
			tableMetaIndexRes.Head[i2].FieldName = column.Name
			tableMetaIndexRes.Head[i2].Desc = column.Desc
		}

		for i2, rows := range tableMetaIndexInfo.Rows {
			tableMetaIndexRes.Rows[i2] = make(map[string]string)
			for k, row := range rows {
				columnName := tableMetaIndexInfo.Column[k].Name
				tableMetaIndexRes.Rows[i2][columnName] = row
			}
		}

		analysisDataResItemV1.TableMetas[i].Name = tableMeta.Name
		analysisDataResItemV1.TableMetas[i].Schema = tableMeta.Schema
		analysisDataResItemV1.TableMetas[i].CreateTableSQL = tableMeta.CreateTableSQL
	}

	return analysisDataResItemV1
}
