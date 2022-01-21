package ast

import (
	"github.com/pingcap/parser/format"
)

type UnparsedStmt struct {
	node
}

func (n *UnparsedStmt) statement() {}

// Restore returns the sql text from ast tree
func (n *UnparsedStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteString(n.Text())
	return nil
}

func (n *UnparsedStmt) Accept(v Visitor) (node Node, ok bool) {
	return nil, false
}
