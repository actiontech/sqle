//go:build enterprise
// +build enterprise

package v2

import (
	"context"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func getTaskAnalysisData(c echo.Context) error {

	taskID := c.Param("task_id")
	sqlNumber := c.Param("number")

	s := model.GetStorage()
	task, err := v1.GetTaskById(c.Request().Context(), taskID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := v1.CheckCurrentUserCanViewTask(c, task); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskSql, exist, err := s.GetTaskSQLByNumber(taskID, sqlNumber)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewDataNotExistErr("sql number not found"))
	}

	res, err := getSQLAnalysisResult(log.NewEntry(), task.Instance, task.Schema, taskSql.Content)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetTaskAnalysisDataResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertSQLAnalysisResultToRes(res, taskSql.Content),
	})
}

func convertSQLAnalysisResultToRes(res *analysisResult, rawSQL string) *TaskAnalysisDataV2 {

	data := &TaskAnalysisDataV2{}

	// explain
	{
		data.SQLExplain = &SQLExplain{SQL: rawSQL}
		if res.explainResultErr != nil {
			data.SQLExplain.ErrMessage = res.explainResultErr.Error()
		} else {
			classicResult := &v1.ExplainClassicResult{
				Head: make([]v1.TableMetaItemHeadResV1, len(res.explainResult.ClassicResult.Columns)),
				Rows: make([]map[string]string, len(res.explainResult.ClassicResult.Rows)),
			}

			// head
			for i := range res.explainResult.ClassicResult.Columns {
				col := res.explainResult.ClassicResult.Columns[i]
				classicResult.Head[i].FieldName = col.Name
				classicResult.Head[i].Desc = col.Desc
			}

			// rows
			for i := range res.explainResult.ClassicResult.Rows {
				row := res.explainResult.ClassicResult.Rows[i]
				classicResult.Rows[i] = make(map[string]string, len(row))
				for k := range row {
					colName := res.explainResult.ClassicResult.Columns[k].Name
					classicResult.Rows[i][colName] = row[k]
				}
			}
			data.SQLExplain.ClassicResult = classicResult
		}
	}

	// table_metas
	{
		data.TableMetas = &TableMetas{}
		if res.tableMetaResultErr != nil {
			data.TableMetas.ErrMessage = res.tableMetaResultErr.Error()
		} else {
			tableMetaItemsData := make([]*v1.TableMeta, len(res.tableMetaResult.TableMetas))
			for i := range res.tableMetaResult.TableMetas {
				tableMeta := res.tableMetaResult.TableMetas[i]
				tableMetaColumnsInfo := tableMeta.ColumnsInfo
				tableMetaIndexInfo := tableMeta.IndexesInfo
				tableMetaItemsData[i] = &v1.TableMeta{}
				tableMetaItemsData[i].Columns = v1.TableColumns{
					Rows: make([]map[string]string, len(tableMetaColumnsInfo.Rows)),
					Head: make([]v1.TableMetaItemHeadResV1, len(tableMetaColumnsInfo.Columns)),
				}

				tableMetaItemsData[i].Indexes = v1.TableIndexes{
					Rows: make([]map[string]string, len(tableMetaIndexInfo.Rows)),
					Head: make([]v1.TableMetaItemHeadResV1, len(tableMetaIndexInfo.Columns)),
				}

				tableMetaColumnData := tableMetaItemsData[i].Columns
				for j := range tableMetaColumnsInfo.Columns {
					col := tableMetaColumnsInfo.Columns[j]
					tableMetaColumnData.Head[j].FieldName = col.Name
					tableMetaColumnData.Head[j].Desc = col.Desc
				}

				for j := range tableMetaColumnsInfo.Rows {
					tableMetaColumnData.Rows[j] = make(map[string]string, len(tableMetaColumnsInfo.Rows[j]))
					for k := range tableMetaColumnsInfo.Rows[j] {
						colName := tableMetaColumnsInfo.Columns[k].Name
						tableMetaColumnData.Rows[j][colName] = tableMetaColumnsInfo.Rows[j][k]
					}
				}

				tableMetaIndexData := tableMetaItemsData[i].Indexes
				for j := range tableMetaIndexInfo.Columns {
					tableMetaIndexData.Head[j].FieldName = tableMetaIndexInfo.Columns[j].Name
					tableMetaIndexData.Head[j].Desc = tableMetaIndexInfo.Columns[j].Desc
				}

				for j := range tableMetaIndexInfo.Rows {
					tableMetaIndexData.Rows[j] = make(map[string]string, len(tableMetaIndexInfo.Rows[j]))
					for k := range tableMetaIndexInfo.Rows[j] {
						colName := tableMetaIndexInfo.Columns[k].Name
						tableMetaIndexData.Rows[j][colName] = tableMetaIndexInfo.Rows[j][k]
					}
				}

				tableMetaItemsData[i].Name = tableMeta.Name
				tableMetaItemsData[i].Schema = tableMeta.Schema
				tableMetaItemsData[i].CreateTableSQL = tableMeta.CreateTableSQL
				tableMetaItemsData[i].Message = tableMeta.Message
			}
			data.TableMetas.Items = tableMetaItemsData
		}
	}

	// performance_statistics
	{
		data.PerformanceStatistics = &PerformanceStatistics{}

		// affect_rows
		data.PerformanceStatistics.AffectRows = &AffectRows{}
		if res.affectRowsResultErr != nil {
			data.PerformanceStatistics.AffectRows.ErrMessage = res.affectRowsResultErr.Error()
		} else {
			data.PerformanceStatistics.AffectRows.ErrMessage = res.affectRowsResult.ErrMessage
			data.PerformanceStatistics.AffectRows.Count = int(res.affectRowsResult.Count)
		}

	}

	return data
}

type analysisResult struct {
	explainResult    *driverV2.ExplainResult
	explainResultErr error

	tableMetaResult    *driver.GetTableMetaBySQLResult
	tableMetaResultErr error

	affectRowsResult    *driverV2.EstimatedAffectRows
	affectRowsResultErr error
}

func getSQLAnalysisResult(l *logrus.Entry, instance *model.Instance, schema, sql string) (
	res *analysisResult, err error) {

	dsn, err := common.NewDSN(instance, schema)
	if err != nil {
		return nil, err
	}

	plugin, err := driver.GetPluginManager().
		OpenPlugin(l, instance.DbType, &driverV2.Config{DSN: dsn})
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	res = &analysisResult{}
	res.explainResult, res.explainResultErr = explain(instance.DbType, plugin, sql)
	res.tableMetaResult, res.tableMetaResultErr = getTableMetas(instance.DbType, plugin, sql)
	res.affectRowsResult, res.affectRowsResultErr = getRowsAffected(instance.DbType, plugin, sql)

	return res, nil
}

func explain(dbType string, plugin driver.Plugin, sql string) (
	res *driverV2.ExplainResult, err error) {

	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleExplain) {
		return nil, driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleExplain)
	}

	return plugin.Explain(context.TODO(), &driverV2.ExplainConf{Sql: sql})
}

func getTableMetas(dbType string, plugin driver.Plugin, sql string) (
	res *driver.GetTableMetaBySQLResult, err error) {

	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleGetTableMeta) {
		return nil, driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleGetTableMeta)
	}

	return plugin.GetTableMetaBySQL(context.TODO(), &driver.GetTableMetaBySQLConf{Sql: sql})
}

func getRowsAffected(dbType string, plugin driver.Plugin, sql string) (
	res *driverV2.EstimatedAffectRows, err error) {

	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleEstimateSQLAffectRows) {
		return nil, driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleEstimateSQLAffectRows)
	}

	return plugin.EstimateSQLAffectRows(context.TODO(), sql)
}
