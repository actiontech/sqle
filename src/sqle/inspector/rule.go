package inspector

import "github.com/pingcap/tidb/ast"

type Rule struct {
	Code    int
	Desc    string
	Level   int
	Disable bool
	CheckFn func(node ast.StmtNode) (string, error)
}

type Rules []*Rule

func(r Rules)Do(sql string) {

}

func (r Rules) add(code int, desc string, fn func(node ast.StmtNode) (string, error)) {
	r = append(r, &Rule{
		Code:    code,
		Desc:    desc,
		CheckFn: fn,
	})
}

const (
	SELECT_STMT_TABLE_MUST_EXIST = iota
)

func initRules() Rules {
	r := Rules{}
	r.add(SELECT_STMT_TABLE_MUST_EXIST, "", SelectStmtTableMustExist)
	return r
}

func SelectStmtTableMustExist(node ast.StmtNode) (string, error) {
	selectStmt, ok := node.(*ast.SelectStmt)
	if !ok {
		return "", nil
	}
	_ = selectStmt
	return "", nil
}
