package util

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

// FingerprintVisitor implements ast.Visitor interface.
type FingerprintVisitor struct{}

func (f *FingerprintVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if v, ok := n.(*driver.ValueExpr); ok {
		v.Type.Charset = ""
		v.SetValue([]byte("?"))
	}
	return n, false
}

func (f *FingerprintVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

type ParamMarkerChecker struct {
	HasParamMarker bool
}

func (p *ParamMarkerChecker) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	if _, ok := in.(*driver.ParamMarkerExpr); ok {
		p.HasParamMarker = true
		return in, true
	}
	return in, false
}

func (p *ParamMarkerChecker) Leave(in ast.Node) (node ast.Node, skipChildren bool) {
	return in, true
}

func ParseCreateTableStmt(sql string) (*ast.CreateTableStmt, error) {
	t, err := ParseOneSql(sql)
	if err != nil {
		return nil, err
	}
	createStmt, ok := t.(*ast.CreateTableStmt)
	if !ok {
		return nil, fmt.Errorf("stmt not support")
	}
	return createStmt, nil
}

// CapitalizeProcessor implements ast.Visitor interface.
//
// CapitalizeProcessor capitalize identifiers as needed.
//
// format.RestoreNameUppercase can not control name comparisons accurate.
// CASE:
// Database/Table/Table-alias names are case-insensitive when lower_case_table_names equals 1.
// Some identifiers, such as Tablespace names are case-sensitive which not affected by lower_case_table_names.
// ref: https://dev.mysql.com/doc/refman/5.7/en/identifier-case-sensitivity.html
type CapitalizeProcessor struct {
	capitalizeTableName      bool
	capitalizeTableAliasName bool
	capitalizeDatabaseName   bool
}

func (cp *CapitalizeProcessor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableSource:
		if cp.capitalizeTableAliasName {
			stmt.AsName.O = strings.ToUpper(stmt.AsName.O)
		}
	case *ast.TableName:
		if cp.capitalizeTableName {
			stmt.Name.O = strings.ToUpper(stmt.Name.O)
		}
		if cp.capitalizeDatabaseName {
			stmt.Schema.O = strings.ToUpper(stmt.Schema.O)
		}
	}

	if cp.capitalizeDatabaseName {
		switch stmt := in.(type) {
		case *ast.DropDatabaseStmt:
			stmt.Name = strings.ToUpper(stmt.Name)
		case *ast.CreateDatabaseStmt:
			stmt.Name = strings.ToUpper(stmt.Name)
		case *ast.AlterDatabaseStmt:
			stmt.Name = strings.ToUpper(stmt.Name)
		}
	}
	return in, false
}

func (cp *CapitalizeProcessor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

// TableNameExtractor implements ast.Visitor interface.
type TableNameExtractor struct {
	TableNames map[string] /*origin table name without database name*/ *ast.TableName
}

func (te *TableNameExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableName:
		te.TableNames[stmt.Name.O] = stmt
	}
	return in, false
}

func (te *TableNameExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type SelectStmtExtractor struct {
	SelectStmts []*ast.SelectStmt
}

func (se *SelectStmtExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		se.SelectStmts = append(se.SelectStmts, stmt)
	}
	return in, false
}

func (se *SelectStmtExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type SubQueryMaxNestNumExtractor struct {
	MaxNestNum     *int
	CurrentNestNum int
}

func (se *SubQueryMaxNestNumExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	stmt, ok := in.(*ast.SubqueryExpr)
	if !ok {
		return in, false
	}

	if *se.MaxNestNum < se.CurrentNestNum {
		*se.MaxNestNum = se.CurrentNestNum
	}

	numExtractor := SubQueryMaxNestNumExtractor{MaxNestNum: se.MaxNestNum, CurrentNestNum: se.CurrentNestNum + 1}
	stmt.Query.Accept(&numExtractor)

	return in, true
}

func (se *SubQueryMaxNestNumExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type TableSourceExtractor struct {
	TableSources map[string] /*origin table name without database name*/ *ast.TableSource
}

func (ts *TableSourceExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableSource:
		ts.TableSources[(stmt.Source).(*ast.TableName).Name.O] = stmt
	}
	return in, false
}

func (ts *TableSourceExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

// SelectFieldExtractor
// 检测select的字段是否只包含count(*)函数
type SelectFieldExtractor struct {
	IsSelectOnlyIncludeCountFunc bool
}

func (se *SelectFieldExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		isOneFiled := len(stmt.Fields.Fields) == 1
		if !isOneFiled {
			return in, true
		}

		if aggregateFuncExpr, ok := stmt.Fields.Fields[0].Expr.(*ast.AggregateFuncExpr); ok {
			isOneArg := len(aggregateFuncExpr.Args) == 1
			if !isOneArg {
				return in, true
			}

			var arg interface{}
			if expr, ok := aggregateFuncExpr.Args[0].(ast.ValueExpr); ok {
				arg = expr.GetValue()
			}

			isDigitOne := arg.(int64) == 1
			isCountFunc := strings.ToLower(aggregateFuncExpr.F) == ast.AggFuncCount
			if isCountFunc && isDigitOne {
				se.IsSelectOnlyIncludeCountFunc = true
				return in, true
			}
		}
	}
	return in, true
}

func (se *SelectFieldExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type SelectVisitor struct {
	SelectList []*ast.SelectStmt
}

func (v *SelectVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		v.SelectList = append(v.SelectList, stmt)
	}
	return in, false
}

func (v *SelectVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type ColumnNameVisitor struct {
	ColumnNameList []*ast.ColumnNameExpr
}

func (v *ColumnNameVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.ColumnNameExpr:
		v.ColumnNameList = append(v.ColumnNameList, stmt)
	}
	return in, false
}

func (v *ColumnNameVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type EqualConditionVisitor struct {
	ConditionList []*ast.BinaryOperationExpr
}

func (v *EqualConditionVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.BinaryOperationExpr:
		if stmt.Op == opcode.EQ {
			v.ConditionList = append(v.ConditionList, stmt)
		}
	}
	return in, false
}

func (v *EqualConditionVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}
