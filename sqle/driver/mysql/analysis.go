package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/pingcap/parser/ast"
	"github.com/pkg/errors"
)

// ListTablesInSchema list tables in specified schema.
// func (i *MysqlDriverImpl) ListTablesInSchema(ctx context.Context, conf *driver.ListTablesInSchemaConf) (*driver.ListTablesInSchemaResult, error) {
// 	conn, err := i.getDbConn()
// 	if err != nil {
// 		return nil, err
// 	}

// 	schema := conf.Schema
// 	if schema == "" {
// 		schema = i.Ctx.CurrentSchema()
// 	}
// 	tables, err := conn.ShowSchemaTables(schema)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resItems := make([]driver.Table, len(tables))
// 	for i, t := range tables {
// 		resItems[i].Name = t
// 	}
// 	return &driver.ListTablesInSchemaResult{Tables: resItems}, nil
// }

// GetTableMetaByTableName get table's metadata by table name.
// func (i *MysqlDriverImpl) GetTableMetaByTableName(ctx context.Context, conf *driver.GetTableMetaByTableNameConf) (*driver.GetTableMetaByTableNameResult, error) {
// 	schema := conf.Schema
// 	if schema == "" {
// 		schema = i.Ctx.CurrentSchema()
// 	}
// 	columnsInfo, indexesInfo, sql, err := i.getTableMetaByTableName(ctx, schema, conf.Table)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &driver.GetTableMetaByTableNameResult{TableMeta: driver.TableMetaItem{
// 		Name:           conf.Table,
// 		Schema:         schema,
// 		ColumnsInfo:    columnsInfo,
// 		IndexesInfo:    indexesInfo,
// 		CreateTableSQL: sql,
// 	}}, nil
// }

func (i *MysqlDriverImpl) getTableMetaByTableName(ctx context.Context, schema, table string) (driverV2.ColumnsInfo, driverV2.IndexesInfo, string, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return driverV2.ColumnsInfo{}, driverV2.IndexesInfo{}, "", err
	}

	columnsInfo, err := i.getTableColumnsInfo(conn, schema, table)
	if err != nil {
		return driverV2.ColumnsInfo{}, driverV2.IndexesInfo{}, "", err
	}

	indexesInfo, err := i.getTableIndexesInfo(conn, schema, table)
	if err != nil {
		return driverV2.ColumnsInfo{}, driverV2.IndexesInfo{}, "", err
	}

	sql, err := conn.ShowCreateTable(utils.SupplementalQuotationMarks(schema), utils.SupplementalQuotationMarks(table))
	if err != nil {
		return driverV2.ColumnsInfo{}, driverV2.IndexesInfo{}, "", err
	}

	return columnsInfo, indexesInfo, sql, nil
}

func (i *MysqlDriverImpl) getTableColumnsInfo(conn *executor.Executor, schema, tableName string) (driverV2.ColumnsInfo, error) {
	columns := []driverV2.TabularDataHead{
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
		return driverV2.ColumnsInfo{}, err
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

	ret := driverV2.ColumnsInfo{}
	ret.Columns = columns
	ret.Rows = rows
	return ret, nil
}

func (i *MysqlDriverImpl) getTableIndexesInfo(conn *executor.Executor, schema, tableName string) (driverV2.IndexesInfo, error) {
	columns := []driverV2.TabularDataHead{
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
		return driverV2.IndexesInfo{}, err
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

	ret := driverV2.IndexesInfo{}
	ret.Columns = columns
	ret.Rows = rows
	return ret, nil
}

// GetTableMetaBySQL get table's metadata by SQL.
func (i *MysqlDriverImpl) GetTableMetaBySQL(ctx context.Context, conf *driver.GetTableMetaBySQLConf) (*driver.GetTableMetaBySQLResult, error) {
	schemaTableList, err := i.ExtractSchemaTableList(conf.Sql)
	if err != nil {
		return nil, err
	}

	tableMetas := make([]*driver.TableMeta, 0, len(schemaTableList))
	for _, schemaTable := range schemaTableList {
		tableMeta := i.GetTableMeta(ctx, schemaTable)
		tableMetas = append(tableMetas, tableMeta)
	}

	return &driver.GetTableMetaBySQLResult{
		TableMetas: tableMetas,
	}, nil
}

func (i *MysqlDriverImpl) GetTableMeta(ctx context.Context, schemaTable SchemaTable) *driver.TableMeta {
	msg := ""
	columnsInfo, indexesInfo, sql, err := i.getTableMetaByTableName(ctx, schemaTable.Schema, schemaTable.Table)
	if err != nil {
		msg = err.Error()
	}

	tm := &driver.TableMeta{}
	tm.Name = schemaTable.Table
	tm.Schema = schemaTable.Schema
	tm.ColumnsInfo = columnsInfo
	tm.IndexesInfo = indexesInfo
	tm.CreateTableSQL = sql
	tm.Message = msg

	return tm
}

type SchemaTable struct {
	Schema string
	Table  string
}

func (i *MysqlDriverImpl) ExtractSchemaTableList(sql string) ([]SchemaTable, error) {
	// check sql
	if sql == "" {
		return nil, errors.New("the SQL should not be empty")
	}
	// only support dml
	if isDML, err := i.isDML(sql); err != nil {
		return nil, err
	} else if !isDML {
		return nil, driverV2.ErrSQLIsNotSupported
	}

	node, err := util.ParseOneSql(sql)
	if err != nil {
		return nil, err
	}

	var schemaTables []SchemaTable
	schemaTableMap := make(map[string]struct{}, 0)
	addTable := func(t *ast.TableName) {
		schema := t.Schema.String()
		if schema == "" {
			schema = i.Ctx.CurrentSchema()
		}

		schemaTableKey := fmt.Sprintf("%s.%s", schema, t.Name.String())
		if _, ok := schemaTableMap[schemaTableKey]; !ok {
			schemaTableMap[schemaTableKey] = struct{}{}
			schemaTables = append(schemaTables, SchemaTable{
				Schema: schema,
				Table:  t.Name.String(),
			})
		}
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
		if stmt.Table != nil {
			addTable(stmt.Table)
		}
	default:
		return nil, fmt.Errorf("the sql is `%v`, we don't support analysing this sql", sql)
	}

	return schemaTables, nil
}

func (i *MysqlDriverImpl) isDML(sql string) (bool, error) {
	//get tables from sql
	node, err := util.ParseOneSql(sql)
	if err != nil {
		return false, err
	}
	switch node.(type) {
	// pingcap将show语句归为DML语句，应该判断为非DML语句
	// DML文档 https://dev.mysql.com/doc/refman/5.7/en/sql-data-manipulation-statements.html
	case *ast.ShowStmt:
		return false, nil
	case ast.DMLNode:
		return true, nil
	default:
		return false, nil
	}
}

// Explain get explain result for SQL.
func (i *MysqlDriverImpl) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	// check sql
	// only support dml
	if isDML, err := i.isDML(conf.Sql); err != nil {
		return nil, err
	} else if !isDML {
		return nil, driverV2.ErrSQLIsNotSupported
	}

	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	columns, rows, err := conn.Explain(conf.Sql)
	if err != nil {
		return nil, err
	}

	resColumn := make([]driverV2.TabularDataHead, len(columns))
	for i, column := range columns {
		resColumn[i] = driverV2.TabularDataHead{Name: column}
	}

	resRows := make([][]string, len(rows))
	for i, row := range rows {
		for _, s := range row {
			resRows[i] = append(resRows[i], s.String)
		}
	}
	res := driverV2.ExplainClassicResult{
		TabularData: driverV2.TabularData{
			Columns: resColumn,
			Rows:    resRows,
		},
	}

	return &driverV2.ExplainResult{
		ClassicResult: res,
	}, nil
}
