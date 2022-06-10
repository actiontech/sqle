//go:build enterprise
// +build enterprise

package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	driver.RegisterAnalysisDriver(driver.DriverTypeMySQL, newAnalysisDriverInspect)
}

func newAnalysisDriverInspect(log *logrus.Entry, dsn *driver.DSN) (driver.AnalysisDriver, error) {
	var inspect = &Inspect{}

	if dsn != nil {
		conn, err := executor.NewExecutor(log, dsn, dsn.DatabaseName)
		if err != nil {
			return nil, errors.Wrap(err, "new executor in inspect")
		}
		inspect.isConnected = true
		inspect.dbConn = conn
		inspect.inst = dsn

		ctx := session.NewContext(nil, session.WithExecutor(conn))
		ctx.SetCurrentSchema(dsn.DatabaseName)

		inspect.Ctx = ctx
	} else {
		ctx := session.NewContext(nil)
		inspect.Ctx = ctx
	}

	inspect.log = log
	inspect.result = driver.NewInspectResults()
	inspect.isOfflineAudit = dsn == nil

	inspect.cnf = &Config{
		DMLRollbackMaxRows: -1,
		DDLOSCMinSize:      -1,
		DDLGhostMinSize:    -1,
	}

	return inspect, nil
}

// ListTablesInSchema list tables in specified schema.
func (i *Inspect) ListTablesInSchema(ctx context.Context, conf *driver.ListTablesInSchemaConf) (*driver.ListTablesInSchemaResult, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	tables, err := conn.ShowSchemaTables(conf.Schema)
	if err != nil {
		return nil, err
	}

	resItems := make([]driver.Table, len(tables))
	for i, t := range tables {
		resItems[i].Name = t
	}
	return &driver.ListTablesInSchemaResult{Tables: resItems}, nil
}

// GetTableMetaByTableName get table's metadata by table name.
func (i *Inspect) GetTableMetaByTableName(ctx context.Context, conf *driver.GetTableMetaByTableNameConf) (*driver.GetTableMetaByTableNameResult, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	columnsInfo, err := i.getTableColumnsInfo(conn, conf.Schema, conf.Table)
	if err != nil {
		return nil, err
	}

	indexesInfo, err := i.getTableIndexesInfo(conn, conf.Schema, conf.Table)
	if err != nil {
		return nil, err
	}

	sql, err := conn.ShowCreateTable(conf.Schema, conf.Table)
	if err != nil {
		return nil, err
	}

	return &driver.GetTableMetaByTableNameResult{TableMeta: driver.TableMetaItem{
		ColumnsInfo:    columnsInfo,
		IndexesInfo:    indexesInfo,
		CreateTableSQL: sql,
	}}, nil
}

func (i *Inspect) getTableColumnsInfo(conn *executor.Executor, schema, tableName string) (driver.ColumnsInfo, error) {
	columns := []*params.Param{
		{
			Key:   "COLUMN_NAME",
			Value: "COLUMN_NAME",
			Desc:  "列名",
		}, {
			Key:   "COLUMN_TYPE",
			Value: "COLUMN_TYPE",
			Desc:  "列类型",
		}, {
			Key:   "CHARACTER_SET_NAME",
			Value: "CHARACTER_SET_NAME",
			Desc:  "列字符集",
		}, {
			Key:   "IS_NULLABLE",
			Value: "IS_NULLABLE",
			Desc:  "是否可以为空",
		}, {
			Key:   "COLUMN_KEY",
			Value: "COLUMN_KEY",
			Desc:  "列索引",
		}, {
			Key:   "COLUMN_DEFAULT",
			Value: "COLUMN_DEFAULT",
			Desc:  "默认值",
		}, {
			Key:   "EXTRA",
			Value: "EXTRA",
			Desc:  "拓展信息",
		}, {
			Key:   "COLUMN_COMMENT",
			Value: "COLUMN_COMMENT",
			Desc:  "列说明",
		},
	}

	queryColumns := make([]string, len(columns))
	for i, c := range columns {
		queryColumns[i] = c.Value
	}
	records, err := conn.GetTableColumnsInfo(schema, tableName)
	if err != nil {
		return driver.ColumnsInfo{}, err
	}

	rows := make([][]string, len(records))
	for i, record := range records {
		row := []string{
			record.ColumnName,
			record.ColumnType,
			record.CharacterSetName,
			record.IsNullable,
			record.ColumnKey,
			record.ColumnDefault,
			record.Extra,
			record.ColumnComment,
		}
		rows[i] = row
	}

	ret := driver.ColumnsInfo{}
	ret.Column = columns
	ret.Rows = rows
	return ret, nil
}

func (i *Inspect) getTableIndexesInfo(conn *executor.Executor, schema, tableName string) (driver.IndexesInfo, error) {
	columns := []*params.Param{

		{
			Key:   "Column_name",
			Value: "Column_name",
			Desc:  "列名",
		}, {
			Key:   "Key_name",
			Value: "Key_name",
			Desc:  "索引名",
		}, {
			// set the row's value as Yes if Non_unique is 0 and No if Non_unique is 1
			Key:   "Unique",
			Value: "Unique",
			Desc:  "唯一性",
		}, {
			Key:   "Seq_in_index",
			Value: "Seq_in_index",
			Desc:  "列序列",
		}, {
			Key:   "Cardinality",
			Value: "Cardinality",
			Desc:  "基数",
		}, {
			// set the row's value as Yes if the column may contain NULL values and No if not
			Key:   "Null",
			Value: "Null",
			Desc:  "是否可以为空",
		}, {
			Key:   "Index_type",
			Value: "Index_type",
			Desc:  "索引类型",
		}, {
			Key:   "Comment",
			Value: "Comment",
			Desc:  "备注",
		},
	}

	indexRecords, err := conn.GetTableIndexesInfo(schema, tableName)
	if err != nil {
		return driver.IndexesInfo{}, err
	}

	rows := make([][]string, len(indexRecords))
	for i, record := range indexRecords {
		nullable := strings.ToUpper(record.Null)
		if nullable != "YES" {
			nullable = "NO"
		}
		unique := "YES"
		if record.NonUnique == "1" {
			unique = "NO"
		}

		row := []string{
			record.ColumnName,
			record.KeyName,
			unique,
			record.SeqInIndex,
			record.Cardinality,
			nullable,
			record.IndexType,
			record.Comment,
		}
		rows[i] = row
	}

	ret := driver.IndexesInfo{}
	ret.Column = columns
	ret.Rows = rows
	return ret, nil
}

// GetTableMetaBySQL get table's metadata by SQL.
func (i *Inspect) GetTableMetaBySQL(ctx context.Context, conf *driver.GetTableMetaBySQLConf) (*driver.GetTableMetaBySQLResult, error) {
	return nil, nil
}

// Explain get explain result for SQL.
func (i *Inspect) Explain(ctx context.Context, conf *driver.ExplainConf) (*driver.ExplainResult, error) {
	// check sql
	// only support dml
	nodes, err := i.ParseSql(conf.Sql)
	if err != nil {
		return nil, err
	}
	switch nodes[0].(type) {
	case ast.DMLNode:
	default:
		return nil, fmt.Errorf("the sql is `%v`, but we only support DML", conf.Sql)
	}

	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	columns, rows, err := conn.Explain(conf.Sql)
	if err != nil {
		return nil, err
	}

	resColumn := params.Params{}
	for _, column := range columns {
		resColumn = append(resColumn, &params.Param{
			Key:   column,
			Value: column,
		})
	}

	resRows := make([][]string, len(rows))
	for i, row := range rows {
		for _, s := range row {
			resRows[i] = append(resRows[i], s.String)
		}
	}
	res := driver.ExplainClassicResult{
		AnalysisInfoInTableFormat: driver.AnalysisInfoInTableFormat{
			Column: resColumn,
			Rows:   resRows,
		},
	}

	return &driver.ExplainResult{
		ClassicResult: res,
	}, nil
}
