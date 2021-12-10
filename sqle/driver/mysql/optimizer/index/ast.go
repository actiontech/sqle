package index

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
	"github.com/pkg/errors"
)

type selectAST struct {
	selectStmt *ast.SelectStmt
}

func newSelectAST(sql string) (*selectAST, error) {
	stmt, err := parser.New().ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "parse sql failed")
	}

	s, ok := stmt.(*ast.SelectStmt)
	if !ok {
		return nil, errors.New("not select stmt")
	}

	return &selectAST{
		selectStmt: s,
	}, nil
}

func (sa *selectAST) EqualPredicateColumnsInWhere() []string {
	var columns []string

	util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
		if boExpr, ok := expr.(*ast.BinaryOperationExpr); ok {
			if boExpr.Op == opcode.EQ {
				if col, ok := boExpr.L.(*ast.ColumnNameExpr); ok {
					columns = append(columns, col.Name.Name.L)
				}
			}
		}

		return false
	}, sa.selectStmt.Where)

	return columns
}

func (sa *selectAST) ColumnsInOrderBy() []string {
	var columns []string

	if sa.selectStmt.OrderBy != nil {
		for _, item := range sa.selectStmt.OrderBy.Items {
			if col, ok := item.Expr.(*ast.ColumnNameExpr); ok {
				// support descending indexes after MySQL 8.0
				// ref: https://dev.mysql.com/doc/refman/8.0/en/descending-indexes.html
				// Before MySQL 8.0, they are parsed but ignored; index values are always
				// stored in ascending order. It's OK to add sequence to column.
				if item.Desc {
					columns = append(columns, fmt.Sprintf("%s desc", col.Name.Name.L))
				} else {
					columns = append(columns, fmt.Sprintf("%s", col.Name.Name.L))
				}
			}
		}
	}

	return columns
}

func (sa *selectAST) ColumnsInProjection() []string {
	var columns []string

	for _, field := range sa.selectStmt.Fields.Fields {
		if field.WildCard == nil {
			if colExpr, ok := field.Expr.(*ast.ColumnNameExpr); ok {
				columns = append(columns, colExpr.Name.Name.L)
			}
		} else {
			return nil
		}
	}

	return columns
}
