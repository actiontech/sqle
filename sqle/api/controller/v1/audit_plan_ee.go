//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var errSQLAnalysisSQLIsNotSupport = errors.New(errors.SQLAnalysisSQLIsNotSupported, driver.ErrSQLIsNotSupported)

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

	explainResult, explainMessage, tableMetaResult, err := getSQLAnalysisResultFromDriver(log.NewEntry(), auditPlanReport.AuditPlan.InstanceDatabase, auditPlanReportSQLV2.SQL, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditPlanAnalysisDataResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    explainAndMetaDataToRes(explainResult, explainMessage, tableMetaResult, auditPlanReportSQLV2.SQL),
	})
}

func getSQLAnalysisResultFromDriver(l *logrus.Entry, database, sql string, instance *model.Instance) (explainResultInput *driver.ExplainResult, explainMessage string, metaDataResultInput *driver.GetTableMetaBySQLResult, err error) {
	dsn, err := common.NewDSN(instance, database)
	if err != nil {
		return nil, "", nil, err
	}
	drvMgr, err := driver.NewDriverManger(log.NewEntry(), instance.DbType, &driver.Config{DSN: dsn})
	if err != nil {
		return nil, "", nil, err
	}
	defer drvMgr.Close(context.TODO())

	analysisDriver, err := drvMgr.GetAnalysisDriver()
	if err != nil {
		return nil, "", nil, err
	}

	explainResult, err := analysisDriver.Explain(context.TODO(), &driver.ExplainConf{Sql: sql})
	if err != nil && err == driver.ErrSQLIsNotSupported {
		return nil, "", nil, errSQLAnalysisSQLIsNotSupport
	} else if err != nil {
		explainMessage = err.Error()
	}

	tableMetaResult, err := analysisDriver.GetTableMetaBySQL(context.TODO(), &driver.GetTableMetaBySQLConf{
		Sql: sql,
	})
	if err != nil && err == driver.ErrSQLIsNotSupported {
		return nil, "", nil, errSQLAnalysisSQLIsNotSupport
	} else if err != nil {
		l.Errorf("get table metadata failed: %v", err)
	}
	return explainResult, explainMessage, tableMetaResult, nil
}

func explainAndMetaDataToRes(explainResultInput *driver.ExplainResult, explainMessage string, metaDataResultInput *driver.GetTableMetaBySQLResult,
	rawSql string) GetSQLAnalysisDataResItemV1 {

	explainResult := explainResultInput
	if explainResult == nil {
		explainResult = &driver.ExplainResult{}
	}
	metaDataResult := metaDataResultInput
	if metaDataResult == nil {
		metaDataResult = &driver.GetTableMetaBySQLResult{}
	}

	analysisDataResItemV1 := GetSQLAnalysisDataResItemV1{
		SQLExplain: SQLExplain{
			ClassicResult: ExplainClassicResult{
				Rows: make([]map[string]string, len(explainResult.ClassicResult.Rows)),
				Head: make([]TableMetaItemHeadResV1, len(explainResult.ClassicResult.Columns)),
			},
		},
		TableMetas: make([]TableMeta, len(metaDataResult.TableMetas)),
	}

	explainResItemV1 := analysisDataResItemV1.SQLExplain.ClassicResult
	for i, column := range explainResult.ClassicResult.Columns {
		explainResItemV1.Head[i].FieldName = column.Name
		explainResItemV1.Head[i].Desc = column.Desc
	}

	for i, rows := range explainResult.ClassicResult.Rows {
		explainResItemV1.Rows[i] = make(map[string]string)
		for k, row := range rows {
			columnName := explainResult.ClassicResult.Columns[k].Name
			explainResItemV1.Rows[i][columnName] = row
		}
	}

	analysisDataResItemV1.SQLExplain.SQL = rawSql
	analysisDataResItemV1.SQLExplain.Message = explainMessage

	for i, tableMeta := range metaDataResult.TableMetas {
		tableMetaColumnsInfo := tableMeta.ColumnsInfo
		tableMetaIndexInfo := tableMeta.IndexesInfo

		analysisDataResItemV1.TableMetas[i].Columns = TableColumns{
			Rows: make([]map[string]string, len(tableMetaColumnsInfo.Rows)),
			Head: make([]TableMetaItemHeadResV1, len(tableMetaColumnsInfo.Columns)),
		}
		analysisDataResItemV1.TableMetas[i].Indexes = TableIndexes{
			Rows: make([]map[string]string, len(tableMetaIndexInfo.Rows)),
			Head: make([]TableMetaItemHeadResV1, len(tableMetaIndexInfo.Columns)),
		}

		tableMetaColumnRes := analysisDataResItemV1.TableMetas[i].Columns
		for i2, column := range tableMetaColumnsInfo.Columns {
			tableMetaColumnRes.Head[i2].FieldName = column.Name
			tableMetaColumnRes.Head[i2].Desc = column.Desc
		}

		for i2, rows := range tableMetaColumnsInfo.Rows {
			tableMetaColumnRes.Rows[i2] = make(map[string]string)
			for k, row := range rows {
				columnName := tableMetaColumnsInfo.Columns[k].Name
				tableMetaColumnRes.Rows[i2][columnName] = row
			}
		}

		tableMetaIndexRes := analysisDataResItemV1.TableMetas[i].Indexes
		for i2, column := range tableMetaIndexInfo.Columns {
			tableMetaIndexRes.Head[i2].FieldName = column.Name
			tableMetaIndexRes.Head[i2].Desc = column.Desc
		}

		for i2, rows := range tableMetaIndexInfo.Rows {
			tableMetaIndexRes.Rows[i2] = make(map[string]string)
			for k, row := range rows {
				columnName := tableMetaIndexInfo.Columns[k].Name
				tableMetaIndexRes.Rows[i2][columnName] = row
			}
		}

		analysisDataResItemV1.TableMetas[i].Name = tableMeta.Name
		analysisDataResItemV1.TableMetas[i].Schema = tableMeta.Schema
		analysisDataResItemV1.TableMetas[i].CreateTableSQL = tableMeta.CreateTableSQL
		analysisDataResItemV1.TableMetas[i].Message = tableMeta.Message
	}

	return analysisDataResItemV1
}
