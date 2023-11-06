package util

import (
	"strings"

	"github.com/pingcap/parser/ast"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

/*
FuncCallStringResultGenerator:

 1. Process the *ast.FuncCallExpr node and return the corresponding result
 2. Support functions in MySQL document: string functions
 3. Support function nesting

https://dev.mysql.com/doc/refman/8.0/en/string-functions.html
*/
type FuncCallStringResultGenerator interface {
	GenerateResult() string // returns an empty string when encounter an unsupported function
}

func NewFuncCallStringResultGenerator(funcCall *ast.FuncCallExpr) FuncCallStringResultGenerator {
	switch funcCall.FnName.L {
	case "concat":
		return &ConCatFunc{ConCat: funcCall}
	case "upper":
		return &UpperFunc{Upper: funcCall}
	default:
		return &UnsupportFunc{}
	}
}

type UnsupportFunc struct {
	Unsupport *ast.FuncCallExpr
}

func (f *UnsupportFunc) GenerateResult() string {
	return ""
}

// https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_concat
type ConCatFunc struct {
	ConCat *ast.FuncCallExpr
}

func (f *ConCatFunc) GenerateResult() string {
	var result string
	for _, arg := range f.ConCat.Args {
		switch pattern := arg.(type) {
		case *driver.ValueExpr:
			result += pattern.Datum.GetString()
		case *ast.FuncCallExpr:
			result += NewFuncCallStringResultGenerator(pattern).GenerateResult()
		}
	}
	return result
}

// https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_upper
type UpperFunc struct {
	Upper *ast.FuncCallExpr
}

func (f *UpperFunc) GenerateResult() string {
	if len(f.Upper.Args) > 0 {
		switch pattern := f.Upper.Args[0].(type) {
		case *driver.ValueExpr:
			return strings.ToUpper(pattern.Datum.GetString())
		case *ast.FuncCallExpr:
			return strings.ToUpper(NewFuncCallStringResultGenerator(pattern).GenerateResult())
		}
	}
	return ""
}
