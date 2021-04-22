package sqlfmt

import (
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	_ "github.com/pingcap/tidb/types/parser_driver"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

func FormatSQL(sql string) string {
	node, err := ParseOneSql(sql)
	if err != nil {
		return sql
	}
	node.Accept(&FormatVisitor{})
	f, err := RestoreToSqlWithFlag(format.RestoreNameBackQuotes|format.RestoreStringDoubleQuotes, node)
	if err != nil {
		return sql
	}
	return f
}

type FormatVisitor struct{}

func (f *FormatVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if v, ok := n.(*driver.ValueExpr); ok {
		v.Type.Charset = ""
	}
	return n, false
}

func (f *FormatVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func RestoreToSqlWithFlag(flag format.RestoreFlags, stmt ast.Node) (string, error) {
	var sb strings.Builder
	ctx := format.NewRestoreCtx(flag, &sb)
	err := stmt.Restore(ctx)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func ParseOneSql(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmt, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}
	return stmt, nil
}
