package inspector

import (
	"errors"
	"github.com/pingcap/tidb/ast"
	"sqle/executor"
	"sqle/storage"
)

// createRollbackSql create rollback sql for input sql; this sql is single sql.
func createRollbackSql(conn *executor.Conn, query string) (string, error) {
	stmts, err := parseSql(storage.DB_TYPE_MYSQL, query)
	if err != nil {
		return "", err
	}
	stmt := stmts[0]
	switch n := stmt.(type) {
	case *ast.AlterTableStmt:
		tableName := getTableName(n.Table)
		createQuery, err := conn.ShowCreateDatabase(tableName)
		if err != nil {
			return "", err
		}
		t, err := parseSql(storage.DB_TYPE_MYSQL, createQuery)

		createStmt, ok := t[0].(*ast.CreateTableStmt)
		if !ok {
			return "", errors.New("")
		}
		return alterTableRollbackSql(createStmt, n)
	default:
		return "", nil
	}
}

func alterTableRollbackSql(t1 *ast.CreateTableStmt, t2 *ast.AlterTableStmt) (string, error) {
	return "", nil
}
