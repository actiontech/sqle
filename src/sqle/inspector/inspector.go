package inspector

import (
	"errors"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/parser"
	"sqle/storage"
)

type Inspector struct {
	Rules Rules
}

func NewInspector() *Inspector {
	return &Inspector{
		Rules: initRules(),
	}
}

func (s *Inspector) parseSql(dbType int, sql string) ([]ast.StmtNode, error) {
	switch dbType {
	case storage.DB_TYPE_MYSQL:
		p := parser.New()
		stmts, err := p.Parse(sql, "", "")
		if err != nil {
			fmt.Printf("parse error: %v\nsql: %v", err, sql)
			return nil, err
		}
		return stmts, nil
	default:
		return nil, errors.New("db type is invalid")
	}
}

//func (s *Inspector) Inspect(config Rules, task *storage.Task) ([]*storage.Sql, error) {
//	stmts, err := s.parseSql(task.Db.DbType, task.ReqSql)
//	if err != nil {
//		return nil, err
//	}
//	s.Rules.
//}

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
