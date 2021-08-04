package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"actiontech.cloud/sqle/sqle/sqle/driver"
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/model"

	"github.com/pingcap/parser/ast"
	_model "github.com/pingcap/parser/model"
)

func (i *Inspect) GenerateAllRollbackSql(executeSQLs []*model.ExecuteSQL) ([]*model.RollbackSQL, error) {
	i.Logger().Info("start generate rollback sql")

	rollbackSqls := []*model.RollbackSQL{}
	for _, executeSQL := range executeSQLs {
		currentSql := executeSQL
		err := i.Add(&currentSql.BaseSQL, func(node ast.Node) error {
			rollbackSql, reason, err := i.GenerateRollbackSql(node)
			if rollbackSql != "" {
				rollbackSqls = append(rollbackSqls, &model.RollbackSQL{
					BaseSQL: model.BaseSQL{
						TaskId:  currentSql.TaskId,
						Content: rollbackSql,
					},
					ExecuteSQLId: currentSql.ID,
				})
			}
			if reason != "" {
				result := driver.NewInspectResults()
				if currentSql.AuditResult != "" {
					result.Add(currentSql.AuditLevel, currentSql.AuditResult)
				}
				result.Add(model.RuleLevelNotice, reason)
				currentSql.AuditLevel = result.Level()
				currentSql.AuditResult = result.Message()
			}
			return err
		})
		if err != nil {
			i.Logger().Error("add rollback sql failed")
			return nil, err
		}
	}
	if err := i.Do(); err != nil {
		i.Logger().Errorf("generate rollback sql failed")
		return nil, err
	}
	i.Logger().Info("generate rollback sql finish")
	return i.GetAllRollbackSqlReversed(rollbackSqls), nil
}

func (i *Inspect) GetAllRollbackSqlReversed(sqls []*model.RollbackSQL) []*model.RollbackSQL {
	rollbackSqls := []*model.RollbackSQL{}
	// Reverse order
	var number uint = 1
	for n := len(sqls) - 1; n >= 0; n-- {
		rollbackSql := sqls[n]
		if rollbackSql != nil {
			rollbackSql.Number = number
		}
		rollbackSqls = append(rollbackSqls, rollbackSql)
		number += 1
	}
	return rollbackSqls
}

func (i *Inspect) GenerateRollbackSql(node ast.Node) (string, string, error) {
	switch node.(type) {
	case ast.DDLNode:
		return i.GenerateDDLStmtRollbackSql(node)
	case ast.DMLNode:
		return i.GenerateDMLStmtRollbackSql(node)
	}
	return "", "", nil
}

func (i *Inspect) GenerateDDLStmtRollbackSql(node ast.Node) (rollbackSql, unableRollbackReason string, err error) {
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

func (i *Inspect) GenerateDMLStmtRollbackSql(node ast.Node) (rollbackSql, unableRollbackReason string, err error) {
	// Inspect may skip initialized cnf when Audited SQLs in whitelist.
	if i.cnf == nil || i.cnf.DMLRollbackMaxRows < 0 {
		return "", "", nil
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

const (
	NotSupportStatementRollback               = "暂不支持回滚该类型的语句"
	NotSupportMultiTableStatementRollback     = "暂不支持回滚多表的 DML 语句"
	NotSupportOnDuplicatStatementRollback     = "暂不支持回滚 ON DUPLICATE 语句"
	NotSupportSubQueryStatementRollback       = "暂不支持回滚带子查询的语句"
	NotSupportNoPrimaryKeyTableRollback       = "不支持回滚没有主键的表的DML语句"
	NotSupportInsertWithoutPrimaryKeyRollback = "不支持回滚 INSERT 没有指定主键的语句"
	NotSupportExceedMaxRowsRollback           = "预计影响行数超过配置的最大值，不生成回滚语句"
)

// generateAlterTableRollbackSql generate alter table SQL for alter table.
func (i *Inspect) generateAlterTableRollbackSql(stmt *ast.AlterTableStmt) (string, string, error) {
	schemaName := i.getSchemaName(stmt.Table)
	tableName := stmt.Table.Name.String()

	createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
	if err != nil || !exist {
		return "", "", err
	}
	rollbackStmt := &ast.AlterTableStmt{
		Table: newTableName(schemaName, tableName),
		Specs: []*ast.AlterTableSpec{},
	}
	// rename table
	if specs := getAlterTableSpecByTp(stmt.Specs, ast.AlterTableRenameTable); len(specs) > 0 {
		spec := specs[len(specs)-1]
		rollbackStmt.Table = newTableName(schemaName, spec.NewTable.Name.String())
		rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
			Tp:       ast.AlterTableRenameTable,
			NewTable: newTableName(schemaName, tableName),
		})
	}
	// Add columns need drop columns
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropColumn) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableChangeColumn) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableModifyColumn) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAlterColumn) {
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
				if HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
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
					} else {
						// do nothing
					}
				}
			}
		}
	}
	// drop index need add
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropIndex) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropPrimaryKey) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropForeignKey) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableRenameIndex) {
		spec.FromKey, spec.ToKey = spec.ToKey, spec.FromKey
		rollbackStmt.Specs = append(rollbackStmt.Specs, spec)
	}

	// Add constraint (index key, primary key ...) need drop
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
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

	rollbackSql := alterTableStmtFormat(rollbackStmt)
	return rollbackSql, "", nil
}

// generateCreateSchemaRollbackSql generate drop database SQL for create database.
func (i *Inspect) generateCreateSchemaRollbackSql(stmt *ast.CreateDatabaseStmt) (string, string, error) {
	schemaName := stmt.Name
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return "", "", err
	}
	if schemaExist {
		return "", "", err
	}
	rollbackSql := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", schemaName)
	return rollbackSql, "", nil
}

// generateCreateTableRollbackSql generate drop table SQL for create table.
func (i *Inspect) generateCreateTableRollbackSql(stmt *ast.CreateTableStmt) (string, string, error) {
	schemaExist, err := i.isSchemaExist(i.getSchemaName(stmt.Table))
	if err != nil {
		return "", "", err
	}
	// if schema not exist, create table will be failed. don't rollback
	if !schemaExist {
		return "", "", nil
	}

	tableExist, err := i.isTableExist(stmt.Table)
	if err != nil {
		return "", "", err
	}

	if tableExist {
		return "", "", nil
	}
	rollbackSql := fmt.Sprintf("DROP TABLE IF EXISTS %s", i.getTableNameWithQuote(stmt.Table))
	return rollbackSql, "", nil
}

// generateDropTableRollbackSql generate create table SQL for drop table.
func (i *Inspect) generateDropTableRollbackSql(stmt *ast.DropTableStmt) (string, string, error) {
	rollbackSql := ""
	for _, table := range stmt.Tables {
		stmt, tableExist, err := i.getCreateTableStmt(table)
		if err != nil {
			return "", "", err
		}
		// if table not exist, can not rollback it.
		if !tableExist {
			continue
		}
		rollbackSql += stmt.Text() + ";\n"
	}
	return rollbackSql, "", nil
}

// generateCreateIndexRollbackSql generate drop index SQL for create index.
func (i *Inspect) generateCreateIndexRollbackSql(stmt *ast.CreateIndexStmt) (string, string, error) {
	return fmt.Sprintf("DROP INDEX `%s` ON %s", stmt.IndexName, i.getTableNameWithQuote(stmt.Table)), "", nil
}

// generateDropIndexRollbackSql generate create index SQL for drop index.
func (i *Inspect) generateDropIndexRollbackSql(stmt *ast.DropIndexStmt) (string, string, error) {
	indexName := stmt.IndexName
	createTableStmt, tableExist, err := i.getCreateTableStmt(stmt.Table)
	if err != nil {
		return "", "", err
	}
	// if table not exist, don't rollback
	if !tableExist {
		return "", "", nil
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
				return "", NotSupportStatementRollback, nil
			}
			if constraint.Option != nil {
				sql = fmt.Sprintf("%s %s", sql, indexOptionFormat(constraint.Option))
			}
			rollbackSql = sql
		}
	}
	return rollbackSql, "", nil
}

// generateInsertRollbackSql generate delete SQL for insert.
func (i *Inspect) generateInsertRollbackSql(stmt *ast.InsertStmt) (string, string, error) {
	tables := getTables(stmt.Table.TableRefs)
	// table just has one in insert stmt.
	if len(tables) != 1 {
		return "", NotSupportMultiTableStatementRollback, nil
	}
	if stmt.OnDuplicate != nil {
		return "", NotSupportOnDuplicatStatementRollback, nil
	}
	table := tables[0]
	createTableStmt, exist, err := i.getCreateTableStmt(table)
	if err != nil {
		return "", "", err
	}
	// if table not exist, insert will failed.
	if !exist {
		return "", "", nil
	}
	pkColumnsName, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return "", "", err
	}
	if !hasPk {
		return "", NotSupportNoPrimaryKeyTableRollback, nil
	}

	rollbackSql := ""

	// match "insert into table_name (column_name,...) value (v1,...)"
	// match "insert into table_name value (v1,...)"
	if stmt.Lists != nil {
		if int64(len(stmt.Lists)) > i.cnf.DMLRollbackMaxRows {
			return "", NotSupportExceedMaxRowsRollback, nil
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
				return "", "", nil
			}
			for n, name := range columnsName {
				_, isPk := pkColumnsName[name]
				if isPk {
					where = append(where, fmt.Sprintf("%s = '%s'", name, exprFormat(value[n])))
				}
			}
			if len(where) != len(pkColumnsName) {
				return "", NotSupportInsertWithoutPrimaryKeyRollback, nil
			}
			rollbackSql += fmt.Sprintf("DELETE FROM %s WHERE %s;\n",
				i.getTableNameWithQuote(table), strings.Join(where, " AND "))
		}
		return rollbackSql, "", nil
	}

	// match "insert into table_name set col_name = value1, ..."
	if stmt.Setlist != nil {
		if 1 > i.cnf.DMLRollbackMaxRows {
			return "", NotSupportExceedMaxRowsRollback, nil
		}
		where := []string{}
		for _, setExpr := range stmt.Setlist {
			name := setExpr.Column.Name.String()
			_, isPk := pkColumnsName[name]
			if isPk {
				where = append(where, fmt.Sprintf("%s = '%s'", name, exprFormat(setExpr.Expr)))
			}
		}
		if len(where) != len(pkColumnsName) {
			return "", "", nil
		}
		rollbackSql = fmt.Sprintf("DELETE FROM %s WHERE %s;\n",
			i.getTableNameWithQuote(table), strings.Join(where, " AND "))
	}
	return rollbackSql, "", nil
}

// generateDeleteRollbackSql generate insert SQL for delete.
func (i *Inspect) generateDeleteRollbackSql(stmt *ast.DeleteStmt) (string, string, error) {
	// not support multi-table syntax
	if stmt.IsMultiTable {
		i.Logger().Infof("not support generate rollback sql with multi-delete statement")
		return "", NotSupportMultiTableStatementRollback, nil
	}
	// sub query statement
	if whereStmtHasSubQuery(stmt.Where) {
		i.Logger().Infof("not support generate rollback sql with sub query")
		return "", NotSupportSubQueryStatementRollback, nil
	}
	var err error
	tables := getTables(stmt.TableRefs.TableRefs)
	table := tables[0]
	createTableStmt, exist, err := i.getCreateTableStmt(table)
	if err != nil || !exist {
		return "", "", err
	}
	_, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return "", "", err
	}
	if !hasPk {
		return "", NotSupportNoPrimaryKeyTableRollback, nil
	}

	var max = i.cnf.DMLRollbackMaxRows
	limit, err := getLimitCount(stmt.Limit, max+1)
	if err != nil {
		return "", "", err
	}
	if limit > max {
		count, err := i.getRecordCount(table, "", stmt.Where, stmt.Order, limit)
		if err != nil {
			return "", "", err
		}
		if count > max {
			return "", NotSupportExceedMaxRowsRollback, nil
		}
	}
	records, err := i.getRecords(table, "", stmt.Where, stmt.Order, limit)

	values := []string{}

	columnsName := []string{}
	for _, col := range createTableStmt.Cols {
		columnsName = append(columnsName, col.Name.Name.String())
	}
	for _, record := range records {
		if len(record) != len(columnsName) {
			return "", "", nil
		}
		vs := []string{}
		for _, name := range columnsName {
			v := "NULL"
			if record[name].Valid {
				v = fmt.Sprintf("'%s'", record[name].String)
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
	return rollbackSql, "", nil
}

// generateUpdateRollbackSql generate update SQL for update.
func (i *Inspect) generateUpdateRollbackSql(stmt *ast.UpdateStmt) (string, string, error) {
	tableSources := getTableSources(stmt.TableRefs.TableRefs)
	// multi table syntax
	if len(tableSources) != 1 {
		i.Logger().Infof("not support generate rollback sql with multi-update statement")
		return "", NotSupportMultiTableStatementRollback, nil
	}
	// sub query statement
	if whereStmtHasSubQuery(stmt.Where) {
		i.Logger().Infof("not support generate rollback sql with sub query")
		return "", NotSupportSubQueryStatementRollback, nil
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
		return "", NotSupportSubQueryStatementRollback, nil
	default:
		return "", NotSupportStatementRollback, nil
	}
	createTableStmt, exist, err := i.getCreateTableStmt(table)
	if err != nil || !exist {
		return "", "", err
	}
	pkColumnsName, hasPk, err := i.getPrimaryKey(createTableStmt)
	if err != nil {
		return "", "", err
	}
	if !hasPk {
		return "", NotSupportNoPrimaryKeyTableRollback, nil
	}

	var max = i.cnf.DMLRollbackMaxRows
	limit, err := getLimitCount(stmt.Limit, max+1)
	if err != nil {
		return "", "", err
	}
	if limit > max {
		count, err := i.getRecordCount(table, tableAlias, stmt.Where, stmt.Order, limit)
		if err != nil {
			return "", "", err
		}
		if count > max {
			return "", NotSupportExceedMaxRowsRollback, nil
		}
	}
	records, err := i.getRecords(table, tableAlias, stmt.Where, stmt.Order, limit)

	columnsName := []string{}
	rollbackSql := ""
	for _, col := range createTableStmt.Cols {
		columnsName = append(columnsName, col.Name.Name.String())
	}
	for _, record := range records {
		if len(record) != len(columnsName) {
			return "", "", nil
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
						pkValue = exprFormat(l.Expr)
					}
				}
			}
			name := col.Name.String()
			v := "NULL"
			if record[name].Valid {
				v = fmt.Sprintf("'%s'", record[name].String)
			}

			if colChanged {
				value = append(value, fmt.Sprintf("%s = %s", name, v))
			}
			if isPk {
				if isPkChanged {
					where = append(where, fmt.Sprintf("%s = '%s'", name, pkValue))
				} else {
					where = append(where, fmt.Sprintf("%s = %s", name, v))

				}
			}
		}
		rollbackSql += fmt.Sprintf("UPDATE %s SET %s WHERE %s;", i.getTableNameWithQuote(table),
			strings.Join(value, ", "), strings.Join(where, " AND "))
	}
	return rollbackSql, "", nil
}

// getRecords select all data which will be update or delete.
func (i *Inspect) getRecords(tableName *ast.TableName, tableAlias string, where ast.ExprNode,
	order *ast.OrderByClause, limit int64) ([]map[string]sql.NullString, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	sql := i.generateGetRecordsSql("*", tableName, tableAlias, where, order, limit)
	return conn.Db.Query(sql)
}

// getRecordCount select all data count which will be update or delete.
func (i *Inspect) getRecordCount(tableName *ast.TableName, tableAlias string, where ast.ExprNode,
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
func (i *Inspect) generateGetRecordsSql(expr string, tableName *ast.TableName, tableAlias string, where ast.ExprNode,
	order *ast.OrderByClause, limit int64) string {
	recordSql := fmt.Sprintf("SELECT %s FROM %s", expr, getTableNameWithQuote(tableName))
	if tableAlias != "" {
		recordSql = fmt.Sprintf("%s AS %s", recordSql, tableAlias)
	}
	if where != nil {
		recordSql = fmt.Sprintf("%s WHERE %s", recordSql, exprFormat(where))
	}
	if order != nil {
		recordSql = fmt.Sprintf("%s ORDER BY", recordSql)
		for _, item := range order.Items {
			recordSql = fmt.Sprintf("%s %s", recordSql, exprFormat(item.Expr))
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
