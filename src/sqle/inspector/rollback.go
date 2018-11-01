package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	_model "github.com/pingcap/tidb/model"
)

func (i *Inspector) GenerateRollbackSql() ([]string, error) {
	defer i.closeDbConn()

	for _, sql := range i.SqlArray {
		var node ast.StmtNode
		var err error

		node, err = parseOneSql(i.Db.DbType, sql.Sql)
		switch stmt := node.(type) {
		case *ast.AlterTableStmt:
			err = i.generateAlterTableRollbackSql(stmt)
		}
		if err != nil {
			return nil, err
		}
	}
	rollbackSqls := []string{}
	// Reverse order
	for n := len(i.rollbackSqls) - 1; n >= 0; n-- {
		rollbackSqls = append(rollbackSqls, i.rollbackSqls[n])
	}
	return rollbackSqls, nil
}

func (i *Inspector) generateAlterTableRollbackSql(stmt *ast.AlterTableStmt) error {
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

func (i *Inspector) generateCreateSchemaRollbackSql(stmt *ast.CreateDatabaseStmt) error {
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

func (i *Inspector) generateCreateTableRollbackSql(stmt *ast.CreateTableStmt) error {
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
		fmt.Sprintf("DROP TABLE IF EXISTS `%s`", getTableNameWithQuote(stmt.Table)))
	return nil
}

func (i *Inspector) generateDropSchemaRollbackSql(stmt *ast.DropDatabaseStmt) error {
	return nil
}

func (i *Inspector) generateDropTableRollbackSql(stmt *ast.DropTableStmt) error {
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

func (i *Inspector) generateCreateIndex(stmt *ast.CreateIndexStmt) error {
	i.rollbackSqls = append(i.rollbackSqls,
		fmt.Sprintf("DROP INDEX `%s` ON %s", stmt.IndexName, getTableNameWithQuote(stmt.Table)))
	return nil
}

func (i *Inspector) generateDropIndex(stmt *ast.CreateIndexStmt) error {
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
					indexName, getTableNameWithQuote(stmt.Table))
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
