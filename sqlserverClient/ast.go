package sqlserverClient

import (
	"github.com/pingcap/tidb/ast"
)

// Node implementations tidb ast Node interface
type Node struct {
	text string
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
}

type SqlServerStmt struct {
	Node
	isDDL bool
	isDML bool
}

func (s *SqlServerStmt) IsDDLStmt() bool {
	return s.isDDL
}

func (s *SqlServerStmt) IsDMLStmt() bool {
	return s.isDML
}

func NewSqlServerStmt(sql string, isDDL, isDML bool) *SqlServerStmt {
	return &SqlServerStmt{
		Node: Node{
			text: sql,
		},
		isDML: isDML,
		isDDL: isDDL,
	}
}
