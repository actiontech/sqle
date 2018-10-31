package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	_model "github.com/pingcap/tidb/model"
)

//func CreateRollbackSql(task *model.Task, sql string) (string, error) {
//	return "", nil
//	conn, err := executor.OpenDbWithTask(task)
//	if err != nil {
//		return "", err
//	}
//	switch task.Instance.DbType {
//	case model.DB_TYPE_MYSQL:
//		return generateRollbackSql(conn, sql)
//	default:
//		return "", errors.New("db type is invalid")
//	}
//}
//
//// createRollbackSql create rollback sql for input sql; this sql is single sql.
//func generateRollbackSql(conn *executor.Conn, query string) (string, error) {
//	stmts, err := parseSql(model.DB_TYPE_MYSQL, query)
//	if err != nil {
//		return "", err
//	}
//	stmt := stmts[0]
//	switch n := stmt.(type) {
//	case *ast.AlterTableStmt:
//		tableName := getTableName(n.Table)
//		createQuery, err := conn.ShowCreateTable(tableName)
//		if err != nil {
//			return "", err
//		}
//		t, err := parseOneSql(model.DB_TYPE_MYSQL, createQuery)
//		if err != nil {
//			return "", err
//		}
//		createStmt, ok := t.(*ast.CreateTableStmt)
//		if !ok {
//			return "", errors.New("")
//		}
//		return alterTableRollbackSql(createStmt, n)
//	default:
//		return "", nil
//	}
//}

func (i *Inspector) alterTableRollbackSql(stmt *ast.AlterTableStmt) (string, error) {
	schemaName := i.getSchemaName(stmt.Table)
	tableName := stmt.Table.Name.String()

	createTableStmt, exist, err := i.getCreateTableStmt(fmt.Sprintf("%s.%s", schemaName, tableName))
	if err != nil || !exist {
		return "", err
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
		rollbackStmt.Specs = append(rollbackStmt.Specs, &ast.AlterTableSpec{
			Tp:            ast.AlterTableDropColumn,
			OldColumnName: &ast.ColumnName{Name: _model.NewCIStr(spec.NewColumns[0].Name.String())},
		})
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

	// add index

	// drop index

	// add primary key

	return alterTableStmtFormat(rollbackStmt), nil
}
