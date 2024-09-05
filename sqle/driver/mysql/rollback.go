package mysql

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	"github.com/actiontech/sqle/sqle/pkg/i18nPkg"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/pingcap/parser/ast"
	_model "github.com/pingcap/parser/model"
	parserMysql "github.com/pingcap/parser/mysql"
)

func (i *MysqlDriverImpl) GenerateRollbackSql(node ast.Node) (string, i18nPkg.I18nStr, error) {
	switch node.(type) {
	case ast.DDLNode:
		return i.GenerateDDLStmtRollbackSql(node)
	case ast.DMLNode:
		return i.GenerateDMLStmtRollbackSql(node)
	}
	return "", nil, nil
}

func (i *MysqlDriverImpl) GenerateDDLStmtRollbackSql(node ast.Node) (rollbackSql string, unableRollbackReason i18nPkg.I18nStr, err error) {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		rollbackSql, unableRollbackReason, err = i.generateAlterTableRollbackSql(stmt)
	case *ast.CreateTableStmt:
		rollbackSql, unableRollbackReason, err = i.generateCreateTableRollbackSql(stmt)
	case *ast.CreateDatabaseStmt:
		rollbackSql, unableRollbackReason, err = i.generateCreateSchemaRollbackSql(stmt)
	case *ast.DropTableStmt:
		rollbackSql, unableRollbackReason, err = i.generateDropTableRollbackSql(stmt)
	case *ast.CreateIndexStmt:
		rollbackSql, unableRollbackReason, err = i.generateCreateIndexRollbackSql(stmt)
	case *ast.DropIndexStmt:
		rollbackSql, unableRollbackReason, err = i.generateDropIndexRollbackSql(stmt)
	}
	return rollbackSql, unableRollbackReason, err
}

func (i *MysqlDriverImpl) GenerateDMLStmtRollbackSql(node ast.Node) (rollbackSql string, unableRollbackReason i18nPkg.I18nStr, err error) {
	// MysqlDriverImpl may skip initialized cnf when Audited SQLs in whitelist.
	if i.cnf == nil || i.cnf.DMLRollbackMaxRows < 0 {
		return "", nil, nil
	}

	paramMarkerChecker := util.ParamMarkerChecker{}
	node.Accept(&paramMarkerChecker)
	if paramMarkerChecker.HasParamMarker {
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportParamMarkerStatementRollback), nil
	}

	hasVarChecker := util.HasVarChecker{}
	node.Accept(&hasVarChecker)
	if hasVarChecker.HasVar {
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportHasVariableRollback), nil
	}

	switch stmt := node.(type) {
	case *ast.InsertStmt:
		rollbackSql, unableRollbackReason, err = i.generateInsertRollbackSql(stmt)
	case *ast.DeleteStmt:
		rollbackSql, unableRollbackReason, err = i.generateDeleteRollbackSql(stmt)
	case *ast.UpdateStmt:
		rollbackSql, unableRollbackReason, err = i.generateUpdateRollbackSql(stmt)
	}
	return
}

// generateAlterTableRollbackSql generate alter table SQL for alter table.
func (i *MysqlDriverImpl) generateAlterTableRollbackSql(stmt *ast.AlterTableStmt) (string, i18nPkg.I18nStr, error) {
	schemaName := i.Ctx.GetSchemaName(stmt.Table)
	tableName := stmt.Table.Name.String()

	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(stmt.Table)
	if err != nil || !exist {
		return "", nil, err
	}
	rollbackStmt := &ast.AlterTableStmt{
		Table: util.NewTableName(schemaName, tableName),
		Specs: []*ast.AlterTableSpec{},
	}
	// rename table
	if specs := util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableRenameTable); len(specs) > 0 {
		spec := specs[len(specs)-1]
		rollbackStmt.Table = util.NewTableName(schemaName, spec.NewTable.Name.String())
		rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
			Tp:       ast.AlterTableRenameTable,
			NewTable: util.NewTableName(schemaName, tableName),
		})
	}
	// Add columns need drop columns
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns) {
		if spec.NewColumns == nil {
			continue
		}
		for _, col := range spec.NewColumns {
			rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
				Tp:            ast.AlterTableDropColumn,
				OldColumnName: &ast.ColumnName{Name: _model.NewCIStr(col.Name.String())},
			})
		}
	}
	// drop columns need add columns
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropColumn) {
		colName := spec.OldColumnName.String()
		for _, col := range createTableStmt.Cols {
			if col.Name.String() == colName {
				rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
					Tp:         ast.AlterTableAddColumns,
					NewColumns: []*ast.ColumnDef{col},
				})
			}
		}
	}
	// change column need change
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableChangeColumn) {
		if spec.NewColumns == nil {
			continue
		}
		for _, col := range createTableStmt.Cols {
			if col.Name.String() == spec.OldColumnName.String() {
				rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
					Tp:            ast.AlterTableChangeColumn,
					OldColumnName: spec.NewColumns[0].Name,
					NewColumns:    []*ast.ColumnDef{col},
				})
			}
		}
	}

	// modify column need modify
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableModifyColumn) {
		if spec.NewColumns == nil {
			continue
		}
		for _, col := range createTableStmt.Cols {
			if col.Name.String() == spec.NewColumns[0].Name.String() {
				rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
					Tp:         ast.AlterTableModifyColumn,
					NewColumns: []*ast.ColumnDef{col},
				})
			}
		}
	}

	/*
		+----------------------------------- alter column -----------------------------------+
		v1 varchar(20) NOT NULL  DEFAULT "test",
			1. alter column v1 set default "TEST" -> alter column v1 set default "test",
			2. alter column v1 drop default -> alter column v1 set default "test",

		v2 varchar(20) NOT NULL,
			1. alter column v1 set default "TEST", -> alter column v1 DROP DEFAULT,
			2. alter column v1 DROP DEFAULT, -> no nothing,
		+------------------------------------------------------------------------------------+
	*/
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAlterColumn) {
		if spec.NewColumns == nil {
			continue
		}
		newColumn := spec.NewColumns[0]
		newSpec := &ast.AlterTableSpec{
			Tp: ast.AlterTableAlterColumn,
			NewColumns: []*ast.ColumnDef{
				{
					Name: newColumn.Name,
				},
			},
		}
		for _, col := range createTableStmt.Cols {
			if col.Name.String() == newColumn.Name.String() {
				if util.HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
					for _, op := range col.Options {
						if op.Tp == ast.ColumnOptionDefaultValue {
							newSpec.NewColumns[0].Options = []*ast.ColumnOption{
								{
									Expr: op.Expr,
								},
							}
							rollbackStmt.Specs = append(rollbackStmt.Specs, newSpec)
						}
					}
				} else {
					// if *ast.ColumnDef.Options is nil, it is "DROP DEFAULT",
					if newColumn.Options != nil {
						rollbackStmt.Specs = append(rollbackStmt.Specs, newSpec)
					}
				}
			}
		}
	}
	// drop index need add
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropIndex) {
		for _, constraint := range createTableStmt.Constraints {
			if constraint.Name == spec.Name {
				rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
					Tp:         ast.AlterTableAddConstraint,
					Constraint: constraint,
				})
			}
		}
	}
	// drop primary key need add
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropPrimaryKey) {
		_ = spec
		for _, constraint := range createTableStmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
					Tp:         ast.AlterTableAddConstraint,
					Constraint: constraint,
				})
			}
		}
	}

	// drop foreign key need add
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropForeignKey) {
		for _, constraint := range createTableStmt.Constraints {
			if constraint.Name == spec.Name {
				rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
					Tp:         ast.AlterTableAddConstraint,
					Constraint: constraint,
				})
			}
		}
	}

	// rename index
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableRenameIndex) {
		spec.FromKey, spec.ToKey = spec.ToKey, spec.FromKey
		rollbackStmt.Specs = append(rollbackStmt.Specs, spec)
	}

	// Add constraint (index key, primary key ...) need drop
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
		switch spec.Constraint.Tp {
		case ast.ConstraintIndex, ast.ConstraintUniq:
			// Add index without index name, index name will be created by db
			if spec.Constraint.Name == "" {
				continue
			}
			rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
				Tp:   ast.AlterTableDropIndex,
				Name: spec.Constraint.Name,
			})
		case ast.ConstraintPrimaryKey:
			rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
				Tp: ast.AlterTableDropPrimaryKey,
			})
		case ast.ConstraintForeignKey:
			rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
				Tp:   ast.AlterTableDropForeignKey,
				Name: spec.Constraint.Name,
			})
		}
	}

	rollbackSql := util.AlterTableStmtFormat(rollbackStmt)
	return rollbackSql, nil, nil
}

// generateCreateSchemaRollbackSql generate drop database SQL for create database.
func (i *MysqlDriverImpl) generateCreateSchemaRollbackSql(stmt *ast.CreateDatabaseStmt) (string, i18nPkg.I18nStr, error) {
	schemaName := stmt.Name
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return "", nil, err
	}
	if schemaExist {
		return "", nil, err
	}
	rollbackSql := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", schemaName)
	return rollbackSql, nil, nil
}

// generateCreateTableRollbackSql generate drop table SQL for create table.
func (i *MysqlDriverImpl) generateCreateTableRollbackSql(stmt *ast.CreateTableStmt) (string, i18nPkg.I18nStr, error) {
	schemaExist, err := i.Ctx.IsSchemaExist(i.Ctx.GetSchemaName(stmt.Table))
	if err != nil {
		return "", nil, err
	}
	// if schema not exist, create table will be failed. don't rollback
	if !schemaExist {
		return "", nil, nil
	}

	tableExist, err := i.Ctx.IsTableExist(stmt.Table)
	if err != nil {
		return "", nil, err
	}

	if tableExist {
		return "", nil, nil
	}
	rollbackSql := fmt.Sprintf("DROP TABLE IF EXISTS %s", i.getTableNameWithQuote(stmt.Table))
	return rollbackSql, nil, nil
}

// generateDropTableRollbackSql generate create table SQL for drop table.
func (i *MysqlDriverImpl) generateDropTableRollbackSql(stmt *ast.DropTableStmt) (string, i18nPkg.I18nStr, error) {
	rollbackSql := ""
	for _, table := range stmt.Tables {
		stmt, tableExist, err := i.Ctx.GetCreateTableStmt(table)
		if err != nil {
			return "", nil, err
		}
		// if table not exist, can not rollback it.
		if !tableExist {
			continue
		}
		rollbackSql += stmt.Text() + ";\n"
	}
	return rollbackSql, nil, nil
}

// generateCreateIndexRollbackSql generate drop index SQL for create index.
func (i *MysqlDriverImpl) generateCreateIndexRollbackSql(stmt *ast.CreateIndexStmt) (string, i18nPkg.I18nStr, error) {
	return fmt.Sprintf("DROP INDEX `%s` ON %s", stmt.IndexName, i.getTableNameWithQuote(stmt.Table)), nil, nil
}

// generateDropIndexRollbackSql generate create index SQL for drop index.
func (i *MysqlDriverImpl) generateDropIndexRollbackSql(stmt *ast.DropIndexStmt) (string, i18nPkg.I18nStr, error) {
	indexName := stmt.IndexName
	createTableStmt, tableExist, err := i.Ctx.GetCreateTableStmt(stmt.Table)
	if err != nil {
		return "", nil, err
	}
	// if table not exist, don't rollback
	if !tableExist {
		return "", nil, nil
	}
	rollbackSql := ""
	for _, constraint := range createTableStmt.Constraints {
		if constraint.Name == indexName {
			sql := ""
			switch constraint.Tp {
			case ast.ConstraintIndex:
				sql = fmt.Sprintf("CREATE INDEX `%s` ON %s",
					indexName, i.getTableNameWithQuote(stmt.Table))
			case ast.ConstraintUniq:
				sql = fmt.Sprintf("CREATE UNIQUE INDEX `%s` ON %s",
					indexName, i.getTableNameWithQuote(stmt.Table))
			default:
				return "", plocale.ShouldLocalizeAll(plocale.NotSupportStatementRollback), nil
			}
			if constraint.Option != nil {
				sql = fmt.Sprintf("%s %s", sql, util.IndexOptionFormat(constraint.Option))
			}
			rollbackSql = sql
		}
	}
	return rollbackSql, nil, nil
}

// generateInsertRollbackSql generate delete SQL for insert.
func (i *MysqlDriverImpl) generateInsertRollbackSql(stmt *ast.InsertStmt) (string, i18nPkg.I18nStr, error) {
	tables := util.GetTables(stmt.Table.TableRefs)
	// table just has one in insert stmt.
	if len(tables) != 1 {
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportMultiTableStatementRollback), nil
	}
	if stmt.OnDuplicate != nil {
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportOnDuplicatStatementRollback), nil
	}
	table := tables[0]
	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil {
		return "", nil, err
	}
	// if table not exist, insert will failed.
	if !exist {
		return "", nil, nil
	}
	pkColumnsName, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return "", nil, err
	}
	if !hasPk {
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportNoPrimaryKeyTableRollback), nil
	}

	rollbackSql := ""

	// match "insert into table_name (column_name,...) value (v1,...)"
	// match "insert into table_name value (v1,...)"
	if stmt.Lists != nil {
		if int64(len(stmt.Lists)) > i.cnf.DMLRollbackMaxRows {
			return "", plocale.ShouldLocalizeAll(plocale.NotSupportExceedMaxRowsRollback), nil
		}
		columnsName := []string{}
		if stmt.Columns != nil {
			for _, col := range stmt.Columns {
				columnsName = append(columnsName, col.Name.String())
			}
		} else {
			for _, col := range createTableStmt.Cols {
				columnsName = append(columnsName, col.Name.String())
			}
		}
		for _, value := range stmt.Lists {
			where := []string{}
			// mysql will throw error: 1136 (21S01): Column count doesn't match value count
			if len(columnsName) != len(value) {
				return "", nil, nil
			}
			for n, name := range columnsName {
				_, isPk := pkColumnsName[name]
				if isPk {
					where = append(where, fmt.Sprintf("%s = '%s'", name, util.ExprFormat(value[n])))
				}
			}
			if len(where) != len(pkColumnsName) {
				return "", plocale.ShouldLocalizeAll(plocale.NotSupportInsertWithoutPrimaryKeyRollback), nil
			}
			rollbackSql += fmt.Sprintf("DELETE FROM %s WHERE %s;\n",
				i.getTableNameWithQuote(table), strings.Join(where, " AND "))
		}
		return rollbackSql, nil, nil
	}

	// match "insert into table_name set col_name = value1, ..."
	if stmt.Setlist != nil {
		if 1 > i.cnf.DMLRollbackMaxRows {
			return "", plocale.ShouldLocalizeAll(plocale.NotSupportExceedMaxRowsRollback), nil
		}
		where := []string{}
		for _, setExpr := range stmt.Setlist {
			name := setExpr.Column.Name.String()
			_, isPk := pkColumnsName[name]
			if isPk {
				where = append(where, fmt.Sprintf("%s = '%s'", name, util.ExprFormat(setExpr.Expr)))
			}
		}
		if len(where) != len(pkColumnsName) {
			return "", nil, nil
		}
		rollbackSql = fmt.Sprintf("DELETE FROM %s WHERE %s;\n",
			i.getTableNameWithQuote(table), strings.Join(where, " AND "))
	}
	return rollbackSql, nil, nil
}

// 将二进制字段转化为十六进制字段
func getHexStrFromBytesStr(byteStr string) string {
	encode := []byte(byteStr)
	return hex.EncodeToString(encode)
}

// generateDeleteRollbackSql generate insert SQL for delete.
func (i *MysqlDriverImpl) generateDeleteRollbackSql(stmt *ast.DeleteStmt) (string, i18nPkg.I18nStr, error) {
	// not support multi-table syntax
	if stmt.IsMultiTable {
		i.Logger().Infof("not support generate rollback sql with multi-delete statement")
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportMultiTableStatementRollback), nil
	}
	// sub query statement
	if util.WhereStmtHasSubQuery(stmt.Where) {
		i.Logger().Infof("not support generate rollback sql with sub query")
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportSubQueryStatementRollback), nil
	}
	var err error
	tables := util.GetTables(stmt.TableRefs.TableRefs)
	table := tables[0]
	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil || !exist {
		return "", nil, err
	}
	_, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return "", nil, err
	}
	if !hasPk {
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportNoPrimaryKeyTableRollback), nil
	}

	var max = i.cnf.DMLRollbackMaxRows
	limit, err := util.GetLimitCount(stmt.Limit, max+1)
	if err != nil {
		return "", nil, err
	}
	if limit > max {
		count, err := i.getRecordCount(table, "", stmt.Where, stmt.Order, limit)
		if err != nil {
			return "", nil, err
		}
		if count > max {
			return "", plocale.ShouldLocalizeAll(plocale.NotSupportExceedMaxRowsRollback), nil
		}
	}
	records, err := i.getRecords(table, "", stmt.Where, stmt.Order, limit)
	if err != nil {
		return "", nil, err
	}
	values := []string{}

	columnsName := []string{}
	colNameDefMap := make(map[string]*ast.ColumnDef)
	for _, col := range createTableStmt.Cols {
		columnsName = append(columnsName, col.Name.Name.String())
		colNameDefMap[col.Name.Name.String()] = col
	}
	for _, record := range records {
		if len(record) != len(columnsName) {
			return "", nil, nil
		}
		vs := []string{}
		for _, name := range columnsName {
			v := "NULL"
			if record[name].Valid {
				colDef := colNameDefMap[name]
				if parserMysql.HasBinaryFlag(colDef.Tp.Flag) {
					hexStr := getHexStrFromBytesStr(record[name].String)
					v = fmt.Sprintf("X'%s'", hexStr)
				} else {
					v = fmt.Sprintf("'%s'", record[name].String)
				}
			}
			vs = append(vs, v)
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(vs, ", ")))
	}
	rollbackSql := ""
	if len(values) > 0 {
		rollbackSql = fmt.Sprintf("INSERT INTO %s (`%s`) VALUES %s;",
			i.getTableNameWithQuote(table), strings.Join(columnsName, "`, `"),
			strings.Join(values, ", "))
	}
	return rollbackSql, nil, nil
}

// generateUpdateRollbackSql generate update SQL for update.
func (i *MysqlDriverImpl) generateUpdateRollbackSql(stmt *ast.UpdateStmt) (string, i18nPkg.I18nStr, error) {
	tableSources := util.GetTableSources(stmt.TableRefs.TableRefs)
	// multi table syntax
	if len(tableSources) != 1 {
		i.Logger().Infof("not support generate rollback sql with multi-update statement")
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportMultiTableStatementRollback), nil
	}
	// sub query statement
	if util.WhereStmtHasSubQuery(stmt.Where) {
		i.Logger().Infof("not support generate rollback sql with sub query")
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportSubQueryStatementRollback), nil
	}
	var (
		table      *ast.TableName
		tableAlias string
	)
	tableSource := tableSources[0]
	switch source := tableSource.Source.(type) {
	case *ast.TableName:
		table = source
		tableAlias = tableSource.AsName.String()
	case *ast.SelectStmt, *ast.UnionStmt:
		i.Logger().Infof("not support generate rollback sql with update-select statement")
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportSubQueryStatementRollback), nil
	default:
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportStatementRollback), nil
	}
	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil || !exist {
		return "", nil, err
	}
	pkColumnsName, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return "", nil, err
	}
	if !hasPk {
		return "", plocale.ShouldLocalizeAll(plocale.NotSupportNoPrimaryKeyTableRollback), nil
	}

	var max = i.cnf.DMLRollbackMaxRows
	limit, err := util.GetLimitCount(stmt.Limit, max+1)
	if err != nil {
		return "", nil, err
	}
	if limit > max {
		count, err := i.getRecordCount(table, tableAlias, stmt.Where, stmt.Order, limit)
		if err != nil {
			return "", nil, err
		}
		if count > max {
			return "", plocale.ShouldLocalizeAll(plocale.NotSupportExceedMaxRowsRollback), nil
		}
	}
	records, err := i.getRecords(table, tableAlias, stmt.Where, stmt.Order, limit)
	if err != nil {
		return "", nil, err
	}
	rollbackSql := ""
	colNameDefMap := make(map[string]*ast.ColumnDef)
	for _, col := range createTableStmt.Cols {
		colNameDefMap[col.Name.Name.String()] = col
	}
	for _, record := range records {
		if len(record) != len(colNameDefMap) {
			return "", nil, nil
		}
		where := []string{}
		value := []string{}
		for _, col := range createTableStmt.Cols {
			colChanged := false
			_, isPk := pkColumnsName[col.Name.Name.L]
			isPkChanged := false
			pkValue := ""

			for _, l := range stmt.List {
				if col.Name.Name.L == l.Column.Name.L {
					colChanged = true
					if isPk {
						isPkChanged = true
						pkValue = util.ExprFormat(l.Expr)
					}
				}
			}
			name := col.Name.Name.O
			v := "NULL"
			if record[name].Valid {
				colDef := colNameDefMap[name]
				if parserMysql.HasBinaryFlag(colDef.Tp.Flag) {
					hexStr := getHexStrFromBytesStr(record[name].String)
					v = fmt.Sprintf("X'%s'", hexStr)
				} else {
					v = fmt.Sprintf("'%s'", record[name].String)
				}
			}

			if colChanged {
				value = append(value, fmt.Sprintf("%s = %s", name, v))
			}
			if isPk {
				if isPkChanged {
					where = append(where, fmt.Sprintf("%s = %s", name, pkValue))
				} else {
					where = append(where, fmt.Sprintf("%s = %s", name, v))

				}
			}
		}
		rollbackSql += fmt.Sprintf("UPDATE %s SET %s WHERE %s;", i.getTableNameWithQuote(table),
			strings.Join(value, ", "), strings.Join(where, " AND "))
	}
	return rollbackSql, nil, nil
}

// getRecords select all data which will be update or delete.
func (i *MysqlDriverImpl) getRecords(tableName *ast.TableName, tableAlias string, where ast.ExprNode,
	order *ast.OrderByClause, limit int64) ([]map[string]sql.NullString, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	sql := i.generateGetRecordsSql("*", tableName, tableAlias, where, order, limit)
	return conn.Db.Query(sql)
}

// getRecordCount select all data count which will be update or delete.
func (i *MysqlDriverImpl) getRecordCount(tableName *ast.TableName, tableAlias string, where ast.ExprNode,
	order *ast.OrderByClause, limit int64) (int64, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return 0, err
	}
	sql := i.generateGetRecordsSql("count(*) as count", tableName, tableAlias, where, order, limit)

	var count int64
	var ok bool
	records, err := conn.Db.Query(sql)
	if err != nil {
		return 0, err
	}
	if len(records) != 1 {
		goto ERROR
	}
	_, ok = records[0]["count"]
	if !ok {
		goto ERROR
	}
	count, err = strconv.ParseInt(records[0]["count"].String, 10, 64)
	if err != nil {
		goto ERROR
	}
	return count, nil

ERROR:
	return 0, errors.New(errors.ConnectRemoteDatabaseError, fmt.Errorf("do not match records for select count(*)"))
}

// generateGetRecordsSql generate select SQL.
func (i *MysqlDriverImpl) generateGetRecordsSql(expr string, tableName *ast.TableName, tableAlias string, where ast.ExprNode,
	order *ast.OrderByClause, limit int64) string {
	recordSql := fmt.Sprintf("SELECT %s FROM %s", expr, i.getTableNameWithQuote(tableName))
	if tableAlias != "" {
		recordSql = fmt.Sprintf("%s AS %s", recordSql, tableAlias)
	}
	if where != nil {
		recordSql = fmt.Sprintf("%s WHERE %s", recordSql, util.ExprFormat(where))
	}
	if order != nil {
		recordSql = fmt.Sprintf("%s ORDER BY", recordSql)
		for _, item := range order.Items {
			recordSql = fmt.Sprintf("%s %s", recordSql, util.ExprFormat(item.Expr))
			if item.Desc {
				recordSql = fmt.Sprintf("%s DESC", recordSql)
			}
		}
	}
	if limit > 0 {
		recordSql = fmt.Sprintf("%s LIMIT %d", recordSql, limit)
	}
	recordSql += ";"
	return recordSql
}
