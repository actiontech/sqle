package inspector

import (
	"database/sql"
	"fmt"
	"github.com/pingcap/tidb/ast"
	_model "github.com/pingcap/tidb/model"
	"sqle/model"
	"strings"
)

func (i *Inspect) GenerateAllRollbackSql() ([]*model.RollbackSql, error) {
	defer i.closeDbConn()
	i.Logger().Info("start generate rollback sql")
	for _, commitSql := range i.Task.CommitSqls {
		err := i.Add(&commitSql.Sql, func(sql *model.Sql) error {
			return i.GenerateRollbackSql(sql)
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
	return i.GetAllRollbackSql(), nil
}

func (i *Inspect) GetAllRollbackSql() []*model.RollbackSql {
	rollbackSqls := []*model.RollbackSql{}
	// Reverse order
	var number uint = 1
	for n := len(i.rollbackSqls) - 1; n >= 0; n-- {
		rollbackSqls = append(rollbackSqls, &model.RollbackSql{
			Sql: model.Sql{
				Number:  number,
				Content: i.rollbackSqls[n],
			},
		})
		number += 1
	}
	return rollbackSqls
}

func (i *Inspect) GenerateRollbackSql(sql *model.Sql) error {
	node := sql.Stmts[0]
	switch node.(type) {
	case ast.DDLNode:
		return i.GenerateDDLStmtRollbackSql(node)
	case ast.DMLNode:
		return i.GenerateDMLStmtRollbackSql(node)
	}
	return nil
}

func (i *Inspect) GenerateDDLStmtRollbackSql(node ast.StmtNode) error {
	var err error
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		err = i.generateAlterTableRollbackSql(stmt)
	case *ast.CreateTableStmt:
		err = i.generateCreateTableRollbackSql(stmt)
	case *ast.CreateDatabaseStmt:
		err = i.generateCreateSchemaRollbackSql(stmt)
	case *ast.DropTableStmt:
		err = i.generateDropTableRollbackSql(stmt)
	}
	if err != nil {
		return err
	}
	return nil
}

func (i *Inspect) GenerateDMLStmtRollbackSql(node ast.StmtNode) error {
	var err error
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		err = i.generateInsertRollbackSql(stmt)
	case *ast.DeleteStmt:
		err = i.generateDeleteRollbackSql(stmt)
	case *ast.UpdateStmt:
		err = i.generateUpdateRollbackSql(stmt)
	}
	if err != nil {
		return err
	}
	return nil
}

func (i *Inspect) generateAlterTableRollbackSql(stmt *ast.AlterTableStmt) error {
	schemaName := i.getSchemaName(stmt.Table)
	tableName := stmt.Table.Name.String()

	createTableStmt, exist, err := i.getCreateTableStmt(fmt.Sprintf("%s.%s", schemaName, tableName))
	if err != nil || !exist {
		return err
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
	// add columns need drop columns
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
				&ast.ColumnDef{
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
								&ast.ColumnOption{
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
	// add constraint (index key, primary key ...)
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
		switch spec.Constraint.Tp {
		case ast.ConstraintIndex, ast.ConstraintUniqKey:
			// add index without index name, index name will be created by db
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
		}
	}
	// drop index
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
	// drop primary
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
	rollbackSql := alterTableStmtFormat(rollbackStmt)
	if rollbackSql != "" {
		i.rollbackSqls = append(i.rollbackSqls, rollbackSql)
	}
	return nil
}

func (i *Inspect) generateCreateSchemaRollbackSql(stmt *ast.CreateDatabaseStmt) error {
	schemaName := stmt.Name
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if schemaExist && stmt.IfNotExists {
		return err
	}
	i.rollbackSqls = append(i.rollbackSqls, fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", schemaName))
	return nil
}

func (i *Inspect) generateCreateTableRollbackSql(stmt *ast.CreateTableStmt) error {
	schemaName := i.getSchemaName(stmt.Table)
	tableName := i.getTableName(stmt.Table)

	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	// if schema not exist, create table will be failed. don't rollback
	if !schemaExist {
		return nil
	}

	tableExist, err := i.isTableExist(tableName)
	if err != nil {
		return err
	}

	if tableExist && stmt.IfNotExists {
		return nil
	}
	i.rollbackSqls = append(i.rollbackSqls,
		fmt.Sprintf("DROP TABLE IF EXISTS %s", i.getTableNameWithQuote(stmt.Table)))
	return nil
}

func (i *Inspect) generateDropTableRollbackSql(stmt *ast.DropTableStmt) error {
	for _, table := range stmt.Tables {
		tableName := i.getTableName(table)
		stmt, tableExist, err := i.getCreateTableStmt(tableName)
		if err != nil {
			return err
		}
		// if table not exist, don't rollback
		if !tableExist {
			return nil
		}
		i.rollbackSqls = append(i.rollbackSqls, stmt.Text())
	}
	return nil
}

func (i *Inspect) generateCreateIndexRollbackSql(stmt *ast.CreateIndexStmt) error {
	i.rollbackSqls = append(i.rollbackSqls,
		fmt.Sprintf("DROP INDEX `%s` ON %s", stmt.IndexName, i.getTableNameWithQuote(stmt.Table)))
	return nil
}

func (i *Inspect) generateDropIndexRollbackSql(stmt *ast.CreateIndexStmt) error {
	indexName := stmt.IndexName
	createTableStmt, tableExist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
	if err != nil {
		return err
	}
	// if table not exist, don't rollback
	if !tableExist {
		return nil
	}
	for _, constraint := range createTableStmt.Constraints {
		if constraint.Name == indexName {
			switch constraint.Tp {
			case ast.ConstraintIndex:
				sql := fmt.Sprintf("CREATE INDEX `%s` ON %s",
					indexName, i.getTableNameWithQuote(stmt.Table))
				if constraint.Option != nil {
					sql = fmt.Sprintf("%s %s", sql, indexOptionFormat(constraint.Option))
				}
				i.rollbackSqls = append(i.rollbackSqls, sql)
			case ast.ConstraintUniq:

			}
		}
	}
	return nil
}

func (i *Inspect) generateInsertRollbackSql(stmt *ast.InsertStmt) error {
	table := getTables(stmt.Table.TableRefs)
	// table just has one in insert stmt.
	if len(table) != 1 {
		return nil
	}
	tableName := i.getTableName(table[0])

	if stmt.OnDuplicate != nil {
		return nil
	}

	createTableStmt, exist, err := i.getCreateTableStmt(tableName)
	if err != nil {
		return err
	}
	// if table not exist, insert will failed.
	if !exist {
		return nil
	}
	pkColumnsName, hasPk := getPrimaryKey(createTableStmt)
	if !hasPk {
		return nil
	}

	rollbackSql := ""

	// match "insert into table_name (column_name,...) value (v1,...)"
	// match "insert into table_name value (v1,...)"
	if stmt.Lists != nil {
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
				return nil
			}
			for n, name := range columnsName {
				_, isPk := pkColumnsName[name]
				if isPk {
					where = append(where, fmt.Sprintf("%s = '%s'", name, exprFormat(value[n])))
				}
			}
			if len(where) != len(pkColumnsName) {
				return nil
			}
			rollbackSql += fmt.Sprintf("DELETE FROM %s WHERE %s;\n",
				i.getTableNameWithQuote(table[0]), strings.Join(where, ", "))
		}
		i.rollbackSqls = append(i.rollbackSqls, rollbackSql)
		return nil
	}

	// match "insert into table_name set col_name = value1, ..."
	if stmt.Setlist != nil {
		where := []string{}
		for _, setExpr := range stmt.Setlist {
			name := setExpr.Column.Name.String()
			_, isPk := pkColumnsName[name]
			if isPk {
				where = append(where, fmt.Sprintf("%s = '%s'", name, exprFormat(setExpr.Expr)))
			}
		}
		if len(where) != len(pkColumnsName) {
			return nil
		}
		i.rollbackSqls = append(i.rollbackSqls, fmt.Sprintf("DELETE FROM %s WHERE %s;\n",
			i.getTableNameWithQuote(table[0]), strings.Join(where, " AND ")))
	}
	return nil
}

func (i *Inspect) generateDeleteRollbackSql(stmt *ast.DeleteStmt) error {
	// not support multi-table syntax
	if stmt.IsMultiTable {
		return nil
	}
	tables := getTables(stmt.TableRefs.TableRefs)
	table := tables[0]
	createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(table))
	if err != nil || !exist {
		return err
	}
	_, hasPk := getPrimaryKey(createTableStmt)
	if !hasPk {
		return nil
	}

	records, err := i.getRecords(table, stmt.Where, stmt.Order, stmt.Limit)

	values := []string{}

	columnsName := []string{}
	for _, col := range createTableStmt.Cols {
		columnsName = append(columnsName, col.Name.Name.String())
	}
	for _, record := range records {
		if len(record) != len(columnsName) {
			return nil
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
	if rollbackSql != "" {
		i.rollbackSqls = append(i.rollbackSqls, rollbackSql)
	}
	return nil
}

func (i *Inspect) generateUpdateRollbackSql(stmt *ast.UpdateStmt) error {
	tables := getTables(stmt.TableRefs.TableRefs)
	// multi table syntax
	if len(tables) != 1 {
		return nil
	}
	table := tables[0]
	createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(table))
	if err != nil || !exist {
		return err
	}
	pkColumnsName, hasPk := getPrimaryKey(createTableStmt)
	if !hasPk {
		return nil
	}

	records, err := i.getRecords(table, stmt.Where, stmt.Order, stmt.Limit)

	columnsName := []string{}
	rollbackSql := ""
	for _, col := range createTableStmt.Cols {
		columnsName = append(columnsName, col.Name.Name.String())
	}
	for _, record := range records {
		if len(record) != len(columnsName) {
			return nil
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
	if rollbackSql != "" {
		i.rollbackSqls = append(i.rollbackSqls, rollbackSql)
	}
	return nil
}

func (i *Inspect) getRecords(tableName *ast.TableName, where ast.ExprNode,
	order *ast.OrderByClause, limit *ast.Limit) ([]map[string]sql.NullString, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	recordSql := fmt.Sprintf("SELECT * FROM %s", i.getTableNameWithQuote(tableName))
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
	if limit != nil {
		recordSql = fmt.Sprintf("%s LIMIT %s", recordSql, exprFormat(limit.Count))
	} else {
		count := GetConfigInt(CONFIG_DML_ROLLBACK_MAX_ROWS)
		recordSql = fmt.Sprintf("%s LIMIT %d", recordSql, count)
	}
	recordSql += ";"
	return conn.Db.Query(recordSql)
}
