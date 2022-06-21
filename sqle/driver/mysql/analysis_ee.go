//go:build enterprise
// +build enterprise

package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/pingcap/parser/ast"
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
	defer conn.Db.Close() //todo do not close connect here, but expose common close() to gracefully close

	schema := conf.Schema
	if schema == "" {
		schema = i.Ctx.CurrentSchema()
	}
	tables, err := conn.ShowSchemaTables(schema)
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
	schema := conf.Schema
	if schema == "" {
		schema = i.Ctx.CurrentSchema()
	}
	columnsInfo, indexesInfo, sql, err := i.getTableMetaByTableName(ctx, schema, conf.Table)
	if err != nil {
		return nil, err
	}

	return &driver.GetTableMetaByTableNameResult{TableMeta: driver.TableMetaItem{
		Name:           conf.Table,
		Schema:         schema,
		ColumnsInfo:    columnsInfo,
		IndexesInfo:    indexesInfo,
		CreateTableSQL: sql,
	}}, nil
}

func (i *Inspect) getTableMetaByTableName(ctx context.Context, schema, table string) (driver.ColumnsInfo, driver.IndexesInfo, string, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return driver.ColumnsInfo{}, driver.IndexesInfo{}, "", err
	}
	defer i.Close(ctx) //todo do not close connect here, but expose common close() to gracefully close

	columnsInfo, err := i.getTableColumnsInfo(conn, schema, table)
	if err != nil {
		return driver.ColumnsInfo{}, driver.IndexesInfo{}, "", err
	}

	indexesInfo, err := i.getTableIndexesInfo(conn, schema, table)
	if err != nil {
		return driver.ColumnsInfo{}, driver.IndexesInfo{}, "", err
	}

	sql, err := conn.ShowCreateTable(utils.SupplementalQuotationMarks(schema), utils.SupplementalQuotationMarks(table))
	if err != nil {
		return driver.ColumnsInfo{}, driver.IndexesInfo{}, "", err
	}

	return columnsInfo, indexesInfo, sql, nil
}

func (i *Inspect) getTableColumnsInfo(conn *executor.Executor, schema, tableName string) (driver.ColumnsInfo, error) {
	columns := []driver.AnalysisInfoHead{
		{
			Name: "COLUMN_NAME",
			Desc: "列名",
		}, {
			Name: "COLUMN_TYPE",
			Desc: "列类型",
		}, {
			Name: "CHARACTER_SET_NAME",
			Desc: "列字符集",
		}, {
			Name: "IS_NULLABLE",
			Desc: "是否可以为空",
		}, {
			Name: "COLUMN_KEY",
			Desc: "列索引",
		}, {
			Name: "COLUMN_DEFAULT",
			Desc: "默认值",
		}, {
			Name: "EXTRA",
			Desc: "拓展信息",
		}, {
			Name: "COLUMN_COMMENT",
			Desc: "列说明",
		},
	}

	queryColumns := make([]string, len(columns))
	for i, c := range columns {
		queryColumns[i] = c.Name
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
	columns := []driver.AnalysisInfoHead{
		{
			Name: "Column_name",
			Desc: "列名",
		}, {
			Name: "Key_name",
			Desc: "索引名",
		}, {
			// set the row's value as Yes if Non_unique is 0 and No if Non_unique is 1
			Name: "Unique",
			Desc: "唯一性",
		}, {
			Name: "Seq_in_index",
			Desc: "列序列",
		}, {
			Name: "Cardinality",
			Desc: "基数",
		}, {
			// set the row's value as Yes if the column may contain NULL values and No if not
			Name: "Null",
			Desc: "是否可以为空",
		}, {
			Name: "Index_type",
			Desc: "索引类型",
		}, {
			Name: "Comment",
			Desc: "备注",
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
	// check sql
	if conf.Sql == "" {
		return nil, errors.New("the SQL should not be empty")
	}
	// only support dml
	if isDML, err := i.isDML(conf.Sql); err != nil {
		return nil, err
	} else if !isDML {
		return nil, fmt.Errorf("the sql is `%v`, but we only support DML", conf.Sql)
	}

	node, err := util.ParseOneSql(conf.Sql)
	if err != nil {
		return nil, err
	}

	type schemaTable struct {
		Schema string
		Table  string
	}

	schemaTables := []schemaTable{}
	addTable := func(t *ast.TableName) {
		schema := t.Schema.String()
		if schema == "" {
			schema = i.Ctx.CurrentSchema()
		}
		schemaTables = append(schemaTables, schemaTable{
			Schema: schema,
			Table:  t.Name.String(),
		})
	}
	getMultiTables := func(stmt *ast.Join) {
		tables := util.GetTables(stmt)
		for _, t := range tables {
			addTable(t)
		}
	}

	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil {
			break
		}
		getMultiTables(stmt.From.TableRefs)
	case *ast.UnionStmt:
		for _, selectStmt := range stmt.SelectList.Selects {
			if selectStmt.From == nil {
				continue
			}
			getMultiTables(selectStmt.From.TableRefs)
		}
	case *ast.UpdateStmt:
		getMultiTables(stmt.TableRefs.TableRefs)
	case *ast.InsertStmt:
		getMultiTables(stmt.Table.TableRefs)
		if stmt.Select != nil {
			getMultiTables(stmt.Select.(*ast.SelectStmt).From.TableRefs)
		}
	case *ast.DeleteStmt:
		getMultiTables(stmt.TableRefs.TableRefs)
	case *ast.LoadDataStmt:
		addTable(stmt.Table)
	case *ast.ShowStmt:
		addTable(stmt.Table)
	default:
		return nil, fmt.Errorf("the sql is `%v`, we don't support analysing this sql", conf.Sql)
	}

	tableMetas := make([]driver.TableMetaItem, len(schemaTables))
	for j, schemaTable := range schemaTables {
		columnsInfo, indexesInfo, sql, err := i.getTableMetaByTableName(ctx, schemaTable.Schema, schemaTable.Table)
		if err != nil {
			return nil, err
		}
		tableMetas[j] = driver.TableMetaItem{
			Name:           schemaTable.Table,
			Schema:         schemaTable.Schema,
			ColumnsInfo:    columnsInfo,
			IndexesInfo:    indexesInfo,
			CreateTableSQL: sql,
		}
	}
	return &driver.GetTableMetaBySQLResult{
		TableMetas: tableMetas,
	}, nil
}

func (i *Inspect) isDML(sql string) (bool, error) {
	//get tables from sql
	node, err := util.ParseOneSql(sql)
	if err != nil {
		return false, err
	}
	switch node.(type) {
	case ast.DMLNode:
		return true, nil
	default:
		return false, nil
	}
}

// Explain get explain result for SQL.
func (i *Inspect) Explain(ctx context.Context, conf *driver.ExplainConf) (*driver.ExplainResult, error) {
	// check sql
	// only support dml
	if isDML, err := i.isDML(conf.Sql); err != nil {
		return nil, err
	} else if !isDML {
		return nil, fmt.Errorf("the sql is `%v`, but we only support DML", conf.Sql)
	}

	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close() //todo do not close connect here, but expose common Close() to gracefully close
	columns, rows, err := conn.Explain(conf.Sql)
	if err != nil {
		return nil, err
	}

	resColumn := make([]driver.AnalysisInfoHead, len(columns))
	for i, column := range columns {
		resColumn[i] = driver.AnalysisInfoHead{Name: column}
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
