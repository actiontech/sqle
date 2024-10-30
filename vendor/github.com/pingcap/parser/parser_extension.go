package parser

import "github.com/pingcap/parser/ast"

func (p Parser) Result() []ast.StmtNode {
	return p.result
}
