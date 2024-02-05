package v2

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver/mysql"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/pingcap/parser/format"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"

	"github.com/labstack/echo/v4"
	"github.com/pingcap/parser/ast"
	pingcapMysql "github.com/pingcap/parser/mysql"
	"github.com/sirupsen/logrus"
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

	sqlContent, err := fillingSQLWithParamMarker(taskSql.Content, task)
	if err != nil {
		log.NewEntry().Errorf("fill param marker sql failed: %v", err)
		sqlContent = taskSql.Content
	}
	res, err := v1.GetSQLAnalysisResult(log.NewEntry(), task.Instance, task.Schema, sqlContent)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetTaskAnalysisDataResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertSQLAnalysisResultToRes(res, taskSql.Content),
	})
}

func convertSQLAnalysisResultToRes(res *v1.AnalysisResult, rawSQL string) *TaskAnalysisDataV2 {

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
				classicResult.Head[i].Desc = col.Desc
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
		if res.AffectRowsResultErr != nil {
			data.PerformanceStatistics.AffectRows.ErrMessage = res.AffectRowsResultErr.Error()
		} else {
			data.PerformanceStatistics.AffectRows.ErrMessage = res.AffectRowsResult.ErrMessage
			data.PerformanceStatistics.AffectRows.Count = int(res.AffectRowsResult.Count)
		}

	}

	return data
}
func fillingSQLWithParamMarker(sqlContent string, task *model.Task) (string, error) {
	l := log.NewEntry()
	driver := mysql.MysqlDriverImpl{}
	nodes, err := driver.ParseSql(sqlContent)
	if err != nil {
		return sqlContent, err
	}

	if task.Instance == nil {
		return sqlContent, nil
	}

	// sql分析是单条sql分析
	if len(nodes) == 0 {
		return sqlContent, nil
	}
	node := nodes[0]
	var tableRefs *ast.Join
	var where ast.ExprNode
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		tableRefs = stmt.From.TableRefs
		where = stmt.Where
	case *ast.UpdateStmt:
		tableRefs = stmt.TableRefs.TableRefs
		where = stmt.Where
	case *ast.DeleteStmt:
		tableRefs = stmt.TableRefs.TableRefs
		where = stmt.Where
	default:
		return sqlContent, nil
	}

	if where == nil {
		return sqlContent, nil
	}
	schema := task.Schema
	if task.Schema == "" {
		schema = getSchemaFromTableRefs(tableRefs)
	}

	dsn, err := common.NewDSN(task.Instance, schema)
	if err != nil {
		return sqlContent, err
	}
	conn, err := executor.NewExecutor(log.NewEntry(), dsn, schema)
	if err != nil {
		return sqlContent, err
	}
	ctx := session.NewContext(nil, session.WithExecutor(conn))
	ctx.SetCurrentSchema(schema)

	tableNameCreateTableStmtMap := ctx.GetTableNameCreateTableStmtMap(tableRefs)
	fillParamMarker(l, where, tableNameCreateTableStmtMap)
	return restore(node)
}

func fillParamMarker(l *logrus.Entry, where ast.ExprNode, tableCreateStmtMap map[string]*ast.CreateTableStmt) {
	util.ScanWhereStmt(func(expr ast.ExprNode) bool {
		switch stmt := expr.(type) {
		case *ast.BinaryOperationExpr:
			// where name=?, 解析器会将'?'解析为ParamMarkerExpr
			if column, ok := stmt.L.(*ast.ColumnNameExpr); ok {
				if _, ok := stmt.R.(*parserdriver.ParamMarkerExpr); !ok {
					return true
				}
				defaultValue, err := fillColumnDefaultValue(column, tableCreateStmtMap)
				if err != nil {
					l.Error(err)
				}
				if defaultValue == nil {
					return false
				}
				stmt.R = defaultValue
			} else if column, ok := stmt.R.(*ast.ColumnNameExpr); ok {
				// 存在列名在比较符号左侧的情况 `where ?=name`
				if _, ok := stmt.L.(*parserdriver.ParamMarkerExpr); !ok {
					return true
				}
				defaultValue, err := fillColumnDefaultValue(column, tableCreateStmtMap)
				if err != nil {
					l.Error(err)
				}
				if defaultValue == nil {
					return false
				}
				stmt.L = defaultValue
			}
		}
		return false
	}, where)
}

func fillColumnDefaultValue(column *ast.ColumnNameExpr, tableCreateStmtMap map[string]*ast.CreateTableStmt) (ast.ExprNode, error) {
	tableName := column.Name.Table.L
	columnName := column.Name.Name.L
	// table name为空，代表没有进行连表查询也没有使用别名 e.g: select * from users where id=?
	if tableName == "" {
		for k := range tableCreateStmtMap {
			tableName = k
		}
	}
	createTableStmt, ok := tableCreateStmtMap[tableName]
	if !ok {
		return nil, fmt.Errorf("fillXmlSql get create table sql failed, table:%v", tableName)
	}
	currentTime := time.Now()
	for _, col := range createTableStmt.Cols {
		if col.Name.Name.L != columnName {
			continue
		}
		var value interface{}
		switch col.Tp.Tp {
		case pingcapMysql.TypeVarchar, pingcapMysql.TypeString:
			value = MybatisXMLCharDefaultValue
		// int类型
		case pingcapMysql.TypeLong, pingcapMysql.TypeTiny, pingcapMysql.TypeShort,
			pingcapMysql.TypeInt24, pingcapMysql.TypeLonglong:
			value = MybatisXMLIntDefaultValue
		// 浮点类型
		case pingcapMysql.TypeNewDecimal, pingcapMysql.TypeFloat, pingcapMysql.TypeDouble:
			value = MybatisXMLFloatDefaultValue
		case pingcapMysql.TypeDatetime:
			value = currentTime.Format("2006-01-02 15:04:05")
		}
		if value != nil {
			defaultValue := &parserdriver.ValueExpr{}
			defaultValue.SetValue(value)
			return defaultValue, nil
		}
	}
	return nil, nil
}

// 还原抽象语法树节点至SQL
func restore(node ast.Node) (string, error) {
	var buf strings.Builder
	sql := ""
	rc := format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)

	if err := node.Restore(rc); err != nil {
		return "", err
	}
	sql = buf.String()
	return sql, nil
}

func getSchemaFromTableRefs(stmt *ast.Join) string {
	schema := ""
	if stmt == nil {
		return schema
	}
	if n := stmt.Left; n != nil {
		switch t := n.(type) {
		case *ast.TableSource:
			if tableName, ok := t.Source.(*ast.TableName); ok {
				schema = tableName.Schema.L
			}
		}
	}
	return schema
}
