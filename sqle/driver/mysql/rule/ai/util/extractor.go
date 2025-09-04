package util

import (
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
)

type FuncInfo struct {
	FuncName string
	Columns  []*ast.ColumnName
	Expr     string
}

type funcExtractor struct {
	funcs []*FuncInfo
}

func (fe *funcExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.FuncCallExpr:
		fi := &FuncInfo{}
		for _, columnNameExpr := range n.Args {
			col, ok := columnNameExpr.(*ast.ColumnNameExpr)
			if !ok {
				continue
			}
			fi.Columns = append(fi.Columns, col.Name)
		}
		fi.FuncName = n.FnName.L
		fi.Expr = ExprFormat(n)
		fe.funcs = append(fe.funcs, fi)
		return in, true
	case *ast.AggregateFuncExpr:
		fi := &FuncInfo{}
		for _, columnNameExpr := range n.Args {
			col, ok := columnNameExpr.(*ast.ColumnNameExpr)
			if !ok {
				continue
			}
			fi.Columns = append(fi.Columns, col.Name)
		}
		// TODO: print aggregate function
		// fi.expr = ExprFormat(n)
		fi.FuncName = n.F
		fe.funcs = append(fe.funcs, fi)
	}
	return in, false
}

func (fe *funcExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type mathOpExtractor struct {
	columnList []*ast.ColumnName
	expr       []string
}

func (me *mathOpExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.BinaryOperationExpr:
		// https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html
		if !isMathComputation(n) {
			return n, false
		}

		if col, ok := n.L.(*ast.ColumnNameExpr); ok {
			me.columnList = append(me.columnList, col.Name)
		}

		if col, ok := n.R.(*ast.ColumnNameExpr); ok {
			me.columnList = append(me.columnList, col.Name)
		}
		me.expr = append(me.expr, ExprFormat(n))
		return in, true

	case *ast.UnaryOperationExpr:
		if n.Op == opcode.Minus {
			col, ok := n.V.(*ast.ColumnNameExpr)
			if !ok {
				return n, false
			}
			me.columnList = append(me.columnList, col.Name)
		}
		me.expr = append(me.expr, ExprFormat(n))
		return in, true
	}

	return in, false
}

func (me *mathOpExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

func isMathComputation(stmt *ast.BinaryOperationExpr) bool {
	return stmt.Op == opcode.Plus || stmt.Op == opcode.Minus || stmt.Op == opcode.Mul || stmt.Op == opcode.Div || stmt.Op == opcode.IntDiv || stmt.Op == opcode.Mod
}

type selectStmtExtractor struct {
	SelectStmts []*ast.SelectStmt
}

func (se *selectStmtExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		se.SelectStmts = append(se.SelectStmts, stmt)
	}
	return in, false
}

func (se *selectStmtExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

// tableNameExtractor implements ast.Visitor interface.
type tableNameExtractor struct {
	tableNames []*ast.TableName
}

func (te *tableNameExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableName:
		te.tableNames = append(te.tableNames, stmt)
	}
	return in, false
}

func (te *tableNameExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

// columnNameExprExtractor is a visitor that extracts column names from SQL expressions.
type columnNameExprExtractor struct {
	columnExpr []*ast.ColumnNameExpr // List to store extracted column exprs
}

// Enter is the method called when a node is visited.
// It checks if the node is of a type that contains a column name and extracts it.
func (ce *columnNameExprExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch expr := in.(type) {
	case *ast.ColumnNameExpr:
		// Append the column if the node is a ColumnNameExpr
		ce.columnExpr = append(ce.columnExpr, expr)
	}
	// Continue traversing to child nodes
	return in, false
}

// Leave is the method called after a node's children have been visited.
func (ce *columnNameExprExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	// Nothing specific happens when leaving a node, just return the node and true
	return in, true
}

// SubqueryExprExtractor implements ast.Visitor interface.
type SubqueryExprExtractor struct {
	expr []*ast.SubqueryExpr
}

func (te *SubqueryExprExtractor) Enter(in ast.Node) (ast.Node, bool) {
	// 情况 1: WHERE、IN、EXISTS 中的子查询
	if sub, ok := in.(*ast.SubqueryExpr); ok {
		te.expr = append(te.expr, sub)
		return in, false
	}

	// 情况 2: FROM 中的子查询（即 SelectStmt 嵌套在 TableSource 中）
	if ts, ok := in.(*ast.TableSource); ok {
		if sel, ok := ts.Source.(*ast.SelectStmt); ok {
			// 包装成 SubqueryExpr 形式以便统一处理
			te.expr = append(te.expr, &ast.SubqueryExpr{
				Query: sel,
			})
		}
	}

	return in, false
}

func (te *SubqueryExprExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type JoinExtractor struct {
	joins []*ast.Join
}

func (je *JoinExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.Join:
		je.joins = append(je.joins, stmt)
	}
	return in, false
}

func (je *JoinExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}
