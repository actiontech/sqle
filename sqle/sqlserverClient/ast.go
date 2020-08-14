package sqlserverClient

import (
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
)

// Node implementations tidb ast Node interface
type Node struct {
	text string
}

func (n *Node) Restore(ctx *format.RestoreCtx) error {
	//TODO implement
	return nil
}

// SetText implements Node interface.
func (n *Node) SetText(text string) {
	n.text = text
}

// Text implements Node interface.
func (n *Node) Text() string {
	return n.text
}

func (n *Node) Accept(v ast.Visitor) (node ast.Node, ok bool) {
	return n, true
}

type SqlServerNode interface {
	ast.Node
	IsDDLStmt() bool
	IsDMLStmt() bool
	IsProcedureStmt() bool
	IsFunctionStmt() bool
}

type SqlServerStmt struct {
	Node
	isDDL       bool
	isDML       bool
	isProcedure bool
	isFunction  bool
}

func (s *SqlServerStmt) IsDDLStmt() bool {
	return s.isDDL
}

func (s *SqlServerStmt) IsDMLStmt() bool {
	return s.isDML
}

func (s *SqlServerStmt) IsProcedureStmt() bool {
	return s.isProcedure
}

func (s *SqlServerStmt) IsFunctionStmt() bool {
	return s.isFunction
}

func NewSqlServerStmt(sql string, isDDL, isDML, isProcedure, isFunction bool) *SqlServerStmt {
	return &SqlServerStmt{
		Node: Node{
			text: sql,
		},
		isDML:       isDML,
		isDDL:       isDDL,
		isProcedure: isProcedure,
		isFunction:  isFunction,
	}
}
