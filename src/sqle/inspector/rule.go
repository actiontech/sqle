package inspector

import "github.com/pingcap/tidb/ast"

type Rule struct {
	Code    int
	Desc    string
	Level   int
	Disable bool
	CheckFn func(node *ast.StmtNode) (string, string, error)
}

const (
	SELECT_STMT_TABLE_MUST_EXIST = iota
)

var RuleMap = []*Rule{
	&Rule{
		Code:    SELECT_STMT_TABLE_MUST_EXIST,
		Desc:    "test",
		Level:   0,
		Disable: false,
	},
}

func CheckSql() {

}
//
//func SelectStmtTableMustExist(node ast.StmtNode) (string, string, error) {
//	selectStmt, ok := node.(*ast.SelectStmt)
//	if !ok {
//		return "", "", nil
//	}
//	selectStmt.TableHints
//}
