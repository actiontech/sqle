package inspector

import (
	"errors"
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/ast"
	"sqle/storage"
)

func Inspect(task *storage.Task) ([]*storage.Sql, error) {
	switch task.Db.DbType {
	case storage.DB_TYPE_MYSQL:
		sqls, err := inspectMysql(task.ReqSql)
		if err != nil {
			return nil, err
		}
		return sqls, nil
	default:
		return nil, errors.New("db type is invalid")
	}
}

// InspectMysql support multi-sql, split by ";".
func inspectMysql(sql string) ([]*storage.Sql, error) {
	sqls := []*storage.Sql{}
	p := parser.New()
	stmts, err := p.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("parse error: %v\nsql: %v", err, sql)
		return nil, err
	}
	for _, stmt := range stmts {
		sql := &storage.Sql{}
		sql.CommitSql = stmt.Text()
		sqls = append(sqls, sql)
		switch stmt.(type) {
		case *ast.SelectStmt:

		}
	}
	return sqls, nil
}
