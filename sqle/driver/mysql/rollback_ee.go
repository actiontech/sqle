//go:build enterprise
// +build enterprise

package mysql

import (
	"fmt"
	"strings"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"

	"github.com/pingcap/parser/ast"
	_model "github.com/pingcap/parser/model"
	parserMysql "github.com/pingcap/parser/mysql"
)

func (i *MysqlDriverImpl) GenerateRollbackSqls(node ast.Node) ([]string, i18nPkg.I18nStr, error) {
	switch node.(type) {
	case ast.DDLNode:
		return i.GenerateDDLStmtRollbackSqls(node)
	case ast.DMLNode:
		return i.GenerateDMLStmtRollbackSqls(node)
	// other
	case *ast.UnparsedStmt:
		return []string{}, i18nPkg.ConvertStr2I18nAsDefaultLang("无法正常解析该SQL，无法进行备份"), nil
	}
	return []string{}, i18nPkg.ConvertStr2I18nAsDefaultLang("暂不支持，该SQL的行备份"), nil
}

func (i *MysqlDriverImpl) GenerateDDLStmtRollbackSqls(node ast.Node) (rollbackSql []string, unableRollbackReason i18nPkg.I18nStr, err error) {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		rollbackSql, unableRollbackReason, err = i.generateAlterTableRollbackSqls(stmt)
	case *ast.CreateTableStmt:
		rollbackSql, unableRollbackReason, err = i.generateCreateTableRollbackSqls(stmt)
	case *ast.CreateDatabaseStmt:
		rollbackSql, unableRollbackReason, err = i.generateCreateSchemaRollbackSqls(stmt)
	case *ast.DropDatabaseStmt:
		return i.generateDropDatabaseRollbackSqls(stmt)
	case *ast.DropTableStmt:
		rollbackSql, unableRollbackReason, err = i.generateDropTableRollbackSqls(stmt)
	case *ast.CreateIndexStmt:
		rollbackSql, unableRollbackReason, err = i.generateCreateIndexRollbackSqls(stmt)
	case *ast.DropIndexStmt:
		rollbackSql, unableRollbackReason, err = i.generateDropIndexRollbackSqls(stmt)
	}
	return rollbackSql, unableRollbackReason, err
}

func (i *MysqlDriverImpl) GenerateDMLStmtRollbackSqls(node ast.Node) (rollbackSql []string, unableRollbackReason i18nPkg.I18nStr, err error) {

	paramMarkerChecker := util.ParamMarkerChecker{}
	node.Accept(&paramMarkerChecker)
	if paramMarkerChecker.HasParamMarker {
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportParamMarkerStatementRollback), nil
	}

	hasVarChecker := util.HasVarChecker{}
	node.Accept(&hasVarChecker)
	if hasVarChecker.HasVar {
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportHasVariableRollback), nil
	}

	switch stmt := node.(type) {
	case *ast.InsertStmt:
		rollbackSql, unableRollbackReason, err = i.generateInsertRollbackSqls(stmt)
	case *ast.DeleteStmt:
		rollbackSql, unableRollbackReason, err = i.generateDeleteRollbackSqls(stmt)
	case *ast.UpdateStmt:
		rollbackSql, unableRollbackReason, err = i.generateUpdateRollbackSqls(stmt)
	}
	return
}

// generateDeleteRollbackSql generate insert SQL for delete.
func (i *MysqlDriverImpl) generateDeleteRollbackSqls(stmt *ast.DeleteStmt) ([]string, i18nPkg.I18nStr, error) {
	// not support multi-table syntax
	if stmt.IsMultiTable {
		i.Logger().Infof("not support generate rollback sql with multi-delete statement")
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportMultiTableStatementRollback), nil
	}
	// sub query statement
	if util.WhereStmtHasSubQuery(stmt.Where) {
		i.Logger().Infof("not support generate rollback sql with sub query")
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportSubQueryStatementRollback), nil
	}
	var err error
	tables := util.GetTables(stmt.TableRefs.TableRefs)
	table := tables[0]
	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil || !exist {
		return []string{}, nil, err
	}
	_, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return []string{}, nil, err
	}
	if !hasPk {
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportNoPrimaryKeyTableRollback), nil
	}

	var max = TemporaryMaxRows
	limit, err := util.GetLimitCount(stmt.Limit, max+1)
	if err != nil {
		return []string{}, nil, err
	}
	if limit > max {
		count, err := i.getRecordCount(table, "", stmt.Where, stmt.Order, limit)
		if err != nil {
			return []string{}, nil, err
		}
		if count > max {
			return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportExceedMaxRowsRollback), nil
		}
	}
	records, err := i.getRecords(table, "", stmt.Where, stmt.Order, limit)
	if err != nil {
		return []string{}, nil, err
	}

	columnsName := []string{}
	colNameDefMap := make(map[string]*ast.ColumnDef)
	for _, col := range createTableStmt.Cols {
		columnsName = append(columnsName, col.Name.Name.String())
		colNameDefMap[col.Name.Name.String()] = col
	}
	rollbackSqls := make([]string, 0, len(records))
	for _, record := range records {
		if len(record) != len(columnsName) {
			return []string{}, nil, nil
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
		rollbackSqls = append(rollbackSqls,
			fmt.Sprintf("INSERT INTO %s (`%s`) VALUES %s;", i.getTableNameWithQuote(table), strings.Join(columnsName, "`, `"), fmt.Sprintf("(%s)", strings.Join(vs, ", "))))
	}
	return rollbackSqls, nil, nil
}

// generateUpdateRollbackSql generate update SQL for update.
func (i *MysqlDriverImpl) generateUpdateRollbackSqls(stmt *ast.UpdateStmt) ([]string, i18nPkg.I18nStr, error) {
	tableSources := util.GetTableSources(stmt.TableRefs.TableRefs)
	// multi table syntax
	if len(tableSources) != 1 {
		i.Logger().Infof("not support generate rollback sql with multi-update statement")
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportMultiTableStatementRollback), nil
	}
	// sub query statement
	if util.WhereStmtHasSubQuery(stmt.Where) {
		i.Logger().Infof("not support generate rollback sql with sub query")
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportSubQueryStatementRollback), nil
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
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportSubQueryStatementRollback), nil
	default:
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportStatementRollback), nil
	}
	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil || !exist {
		return []string{}, nil, err
	}
	pkColumnsName, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return []string{}, nil, err
	}
	if !hasPk {
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportNoPrimaryKeyTableRollback), nil
	}

	var max = TemporaryMaxRows
	limit, err := util.GetLimitCount(stmt.Limit, max+1)
	if err != nil {
		return []string{}, nil, err
	}
	if limit > max {
		count, err := i.getRecordCount(table, tableAlias, stmt.Where, stmt.Order, limit)
		if err != nil {
			return []string{}, nil, err
		}
		if count > max {
			return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportExceedMaxRowsRollback), nil
		}
	}
	records, err := i.getRecords(table, tableAlias, stmt.Where, stmt.Order, limit)
	if err != nil {
		return []string{}, nil, err
	}
	colNameDefMap := make(map[string]*ast.ColumnDef)
	for _, col := range createTableStmt.Cols {
		colNameDefMap[col.Name.Name.String()] = col
	}
	rollbackSqls := make([]string, 0, len(records))
	for _, record := range records {
		if len(record) != len(colNameDefMap) {
			return []string{}, nil, nil
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
						pkValue = restore(l.Expr)
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
		rollbackSqls = append(rollbackSqls,
			fmt.Sprintf("UPDATE %s SET %s WHERE %s;", i.getTableNameWithQuote(table), strings.Join(value, ", "), strings.Join(where, " AND ")))
	}
	return rollbackSqls, nil, nil
}

// generateAlterTableRollbackSql generate alter table SQL for alter table.
func (i *MysqlDriverImpl) generateAlterTableRollbackSqls(stmt *ast.AlterTableStmt) ([]string, i18nPkg.I18nStr, error) {
	schemaName := i.Ctx.GetSchemaName(stmt.Table)
	tableName := stmt.Table.Name.String()

	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(stmt.Table)
	if err != nil || !exist {
		return []string{}, nil, err
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
	return []string{rollbackSql}, nil, nil
}

// generateCreateSchemaRollbackSql generate drop database SQL for create database.
func (i *MysqlDriverImpl) generateCreateSchemaRollbackSqls(stmt *ast.CreateDatabaseStmt) ([]string, i18nPkg.I18nStr, error) {
	schemaName := stmt.Name
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return []string{}, nil, err
	}
	if schemaExist {
		return []string{}, nil, err
	}
	rollbackSql := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", schemaName)
	return []string{rollbackSql}, nil, nil
}

// generateDropDatabaseRollbackSqls generate create database SQL for dropping database.
func (i *MysqlDriverImpl) generateDropDatabaseRollbackSqls(stmt *ast.DropDatabaseStmt) ([]string, i18nPkg.I18nStr, error) {
	return []string{fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", stmt.Name)}, nil, nil
}

// generateCreateTableRollbackSql generate drop table SQL for create table.
func (i *MysqlDriverImpl) generateCreateTableRollbackSqls(stmt *ast.CreateTableStmt) ([]string, i18nPkg.I18nStr, error) {
	schemaExist, err := i.Ctx.IsSchemaExist(i.Ctx.GetSchemaName(stmt.Table))
	if err != nil {
		return []string{}, nil, err
	}
	// if schema not exist, create table will be failed. don't rollback
	if !schemaExist {
		return []string{}, nil, nil
	}

	tableExist, err := i.Ctx.IsTableExist(stmt.Table)
	if err != nil {
		return []string{}, nil, err
	}

	if tableExist {
		return []string{}, nil, nil
	}
	rollbackSql := fmt.Sprintf("DROP TABLE IF EXISTS %s", i.getTableNameWithQuote(stmt.Table))
	return []string{rollbackSql}, nil, nil
}

// generateDropTableRollbackSql generate create table SQL for drop table.
func (i *MysqlDriverImpl) generateDropTableRollbackSqls(stmt *ast.DropTableStmt) ([]string, i18nPkg.I18nStr, error) {
	rollbackSql := ""
	for _, table := range stmt.Tables {
		stmt, tableExist, err := i.Ctx.GetCreateTableStmt(table)
		if err != nil {
			return []string{}, nil, err
		}
		// if table not exist, can not rollback it.
		if !tableExist {
			continue
		}
		rollbackSql += stmt.Text() + ";\n"
	}
	return []string{rollbackSql}, nil, nil
}

// generateCreateIndexRollbackSql generate drop index SQL for create index.
func (i *MysqlDriverImpl) generateCreateIndexRollbackSqls(stmt *ast.CreateIndexStmt) ([]string, i18nPkg.I18nStr, error) {
	return []string{fmt.Sprintf("DROP INDEX `%s` ON %s", stmt.IndexName, i.getTableNameWithQuote(stmt.Table))}, nil, nil
}

// generateDropIndexRollbackSql generate create index SQL for drop index.
func (i *MysqlDriverImpl) generateDropIndexRollbackSqls(stmt *ast.DropIndexStmt) ([]string, i18nPkg.I18nStr, error) {
	indexName := stmt.IndexName
	createTableStmt, tableExist, err := i.Ctx.GetCreateTableStmt(stmt.Table)
	if err != nil {
		return []string{}, nil, err
	}
	// if table not exist, don't rollback
	if !tableExist {
		return []string{}, nil, nil
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
				return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportStatementRollback), nil
			}
			if constraint.Option != nil {
				sql = fmt.Sprintf("%s %s", sql, util.IndexOptionFormat(constraint.Option))
			}
			rollbackSql = sql
		}
	}
	return []string{rollbackSql}, nil, nil
}

// generateInsertRollbackSql generate delete SQL for insert.
func (i *MysqlDriverImpl) generateInsertRollbackSqls(stmt *ast.InsertStmt) ([]string, i18nPkg.I18nStr, error) {
	tables := util.GetTables(stmt.Table.TableRefs)
	// table just has one in insert stmt.
	if len(tables) != 1 {
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportMultiTableStatementRollback), nil
	}
	if stmt.OnDuplicate != nil {
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportOnDuplicatStatementRollback), nil
	}
	table := tables[0]
	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil {
		return []string{}, nil, err
	}
	// if table not exist, insert will failed.
	if !exist {
		return []string{}, nil, nil
	}
	pkColumnsName, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return []string{}, nil, err
	}
	if !hasPk {
		return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportNoPrimaryKeyTableRollback), nil
	}

	rollbackSqls := []string{}

	// match "insert into table_name (column_name,...) value (v1,...)"
	// match "insert into table_name value (v1,...)"
	if stmt.Lists != nil {
		if int64(len(stmt.Lists)) > TemporaryMaxRows {
			return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportExceedMaxRowsRollback), nil
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
				return []string{}, nil, nil
			}
			for n, name := range columnsName {
				_, isPk := pkColumnsName[name]
				if isPk {
					where = append(where, fmt.Sprintf("%s = %s", name, restore(value[n])))
				}
			}
			if len(where) != len(pkColumnsName) {
				return []string{}, plocale.Bundle.LocalizeAll(plocale.NotSupportInsertWithoutPrimaryKeyRollback), nil
			}
			rollbackSqls = append(rollbackSqls,
				fmt.Sprintf("DELETE FROM %s WHERE %s;\n", i.getTableNameWithQuote(table), strings.Join(where, " AND ")))
		}
		return rollbackSqls, nil, nil
	}

	// match "insert into table_name set col_name = value1, ..."
	if stmt.Setlist != nil {
		where := []string{}
		for _, setExpr := range stmt.Setlist {
			name := setExpr.Column.Name.String()
			_, isPk := pkColumnsName[name]
			if isPk {
				where = append(where, fmt.Sprintf("%s = '%s'", name, restore(setExpr.Expr)))
			}
		}
		if len(where) != len(pkColumnsName) {
			return []string{}, nil, nil
		}
		rollbackSqls = append(rollbackSqls, fmt.Sprintf("DELETE FROM %s WHERE %s;\n",
			i.getTableNameWithQuote(table), strings.Join(where, " AND ")))
	}
	return rollbackSqls, nil, nil
}
