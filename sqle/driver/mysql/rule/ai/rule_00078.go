package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00078 = "SQLE00078"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00078,
			Desc:       "禁止使用聚合函数",
			Annotation: "禁止使用SQL聚合函数是为了确保查询的简单性、高性能和数据一致性。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "禁止使用聚合函数",
		AllowOffline: true,
		Func:         RuleSQLE00078,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00078): "在 MySQL 中，禁止使用聚合函数."
您应遵循以下逻辑：
1、检查句子中是否包含SELECT语法节点，存在则进一步检查。
2、检查句子中是否存在聚合函数语法节点，存在报告违反规则。

1. 对于UNION...语句, 对于其中的所有SELECT子句进行与SELECT语句相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00078(input *rulepkg.RuleHandlerInput) error {
	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}
	selectVisitor := &util.SelectVisitor{}
	input.Node.Accept(selectVisitor)
	for _, selectNode := range selectVisitor.SelectList {
		if selectNode.Having != nil {
			isHavingUseFunc := false
			util.ScanWhereStmt(func(expr ast.ExprNode) bool {
				switch expr.(type) {
				case *ast.AggregateFuncExpr:
					isHavingUseFunc = true
					return true
				}
				return false
			}, selectNode.Having.Expr)

			if isHavingUseFunc {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00078)
				return nil
			}
		}
		for _, field := range selectNode.Fields.Fields {
			if _, ok := field.Expr.(*ast.AggregateFuncExpr); ok {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00078)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
