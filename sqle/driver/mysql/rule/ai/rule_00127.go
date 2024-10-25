package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00127 = "SQLE00127"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00127,
			Desc:       "对于MySQL的DML, 不建议在ORDER BY中使用表达式或函数",
			Annotation: "在ORDER BY子句中使用表达式或函数会导致无法有效利用索引，从而可能涉及到全表扫描和使用临时表进行数据排序。这样的操作在处理大数据量时会显著降低查询性能。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "对于MySQL的DML, 不建议在ORDER BY中使用表达式或函数",
		AllowOffline: false,
		Func:    RuleSQLE00127,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00127): "For dml, Mathematical operations and the use of functions on order by is prohibited".
You should follow the following logic:
1. For "select..." statement, checks if the field in the ORDER BY is an function call or math operation, if so, report a violation.
2. For "union ..." statement, perform the same check as mentioned above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00127(input *rulepkg.RuleHandlerInput) error {

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt:
		// "select..." or "union..."
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// Apply the check to each order by items
			orderBy := selectStmt.OrderBy
			if orderBy != nil {
				for _, item := range orderBy.Items {
					if item.Expr == nil {
						continue
					}
					// Check for mathematical operations in the expr. If any are found, return true to indicate a potential issue.
					mathExpr := util.GetMathOpExpr(item.Expr)
					if len(mathExpr) > 0 {
						// Add a rule violation result if an issue is found.
						rulepkg.AddResult(input.Res, input.Rule, SQLE00127)
					}

					// Check for function calls in the expr.
					funcExpr := util.GetFuncExpr(item.Expr)
					// If there are no function expressions, then there's no issue.
					if len(funcExpr) > 0 {
						// Add a rule violation result if an issue is found.
						rulepkg.AddResult(input.Res, input.Rule, SQLE00127)
					}
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
