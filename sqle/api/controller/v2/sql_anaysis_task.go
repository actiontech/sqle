package v2

import (
	"context"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/fillsql"

	"github.com/labstack/echo/v4"
)

const (
	MybatisXMLCharDefaultValue  = "1"
	MybatisXMLIntDefaultValue   = 1
	MybatisXMLFloatDefaultValue = 1.0
	XMLFileExtension            = ".XML"
)

func getTaskAnalysisData(c echo.Context) error {

	taskID := c.Param("task_id")
	sqlNumber := c.Param("number")

	// 获取AffectRowsEnabled查询参数，默认为true以保持向后兼容
	affectRowsEnabled := true
	if affectRowsEnabledStr := c.QueryParam("affectRowsEnabled"); affectRowsEnabledStr != "" {
		if affectRowsEnabledStr == "false" {
			affectRowsEnabled = false
		}
	}

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

	sqlContent, err := fillsql.FillingSQLWithParamMarker(taskSql.Content, task)
	if err != nil {
		log.NewEntry().Errorf("fill param marker sql failed: %v", err)
		sqlContent = taskSql.Content
	}
	res, err := v1.GetSQLAnalysisResult(log.NewEntry(), task.Instance, task.Schema, sqlContent, affectRowsEnabled)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetTaskAnalysisDataResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertSQLAnalysisResultToRes(c.Request().Context(), res, taskSql.Content, affectRowsEnabled),
	})
}

func convertSQLAnalysisResultToRes(ctx context.Context, res *v1.AnalysisResult, rawSQL string, affectRowsEnabled bool) *TaskAnalysisDataV2 {

	data := &TaskAnalysisDataV2{}

	// explain
	{
		data.SQLExplain = &SQLExplain{SQL: rawSQL}
		if res.ExplainResultErr != nil {
			data.SQLExplain.ErrMessage = res.ExplainResultErr.Error()
		} else {
			classicResult := &v1.ExplainClassicResult{
				Head: make([]v1.TableMetaItemHeadResV1, len(res.ExplainResult.ClassicResult.Columns)),
				Rows: make([]map[string]string, len(res.ExplainResult.ClassicResult.Rows)),
			}

			// head
			for i := range res.ExplainResult.ClassicResult.Columns {
				col := res.ExplainResult.ClassicResult.Columns[i]
				classicResult.Head[i].FieldName = col.Name
				classicResult.Head[i].Desc = col.I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
			}

			// rows
			for i := range res.ExplainResult.ClassicResult.Rows {
				row := res.ExplainResult.ClassicResult.Rows[i]
				classicResult.Rows[i] = make(map[string]string, len(row))
				for k := range row {
					colName := res.ExplainResult.ClassicResult.Columns[k].Name
					classicResult.Rows[i][colName] = row[k]
				}
			}
			data.SQLExplain.ClassicResult = classicResult
		}
	}

	// table_metas
	{
		data.TableMetas = &TableMetas{}
		if res.TableMetaResultErr != nil {
			data.TableMetas.ErrMessage = res.TableMetaResultErr.Error()
		} else {
			tableMetaItemsData := make([]*v1.TableMeta, len(res.TableMetaResult.TableMetas))
			for i := range res.TableMetaResult.TableMetas {
				tableMeta := res.TableMetaResult.TableMetas[i]
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
					tableMetaColumnData.Head[j].Desc = col.I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
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
					tableMetaIndexData.Head[j].Desc = tableMetaIndexInfo.Columns[j].I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
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

		// 只有当AffectRowsEnabled为true时才处理影响行数
		if affectRowsEnabled {
			data.PerformanceStatistics.AffectRows = &AffectRows{}
			if res.AffectRowsResultErr != nil {
				data.PerformanceStatistics.AffectRows.ErrMessage = res.AffectRowsResultErr.Error()
			} else if res.AffectRowsResult != nil {
				data.PerformanceStatistics.AffectRows.ErrMessage = res.AffectRowsResult.ErrMessage
				data.PerformanceStatistics.AffectRows.Count = int(res.AffectRowsResult.Count)
			}
		}
	}

	return data
}
