package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00127 = "SQLE00127"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00127,
			Desc:         plocale.Rule00127Desc,
			Annotation:   plocale.Rule00127Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00127Message,
		Func:    RuleSQLE00127,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
