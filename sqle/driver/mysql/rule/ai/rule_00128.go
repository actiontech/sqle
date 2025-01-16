package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00128 = "SQLE00128"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00128,
			Desc:       "不建议使用 HAVING 子句",
			Annotation: "对于索引字段，放在HAVING子句中时不会走索引；建议将HAVING子句改写为WHERE中的查询条件，可以在查询处理期间使用索引，提高SQL的执行效率",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "不建议使用 HAVING 子句.",
		AllowOffline: true,
		Func:    RuleSQLE00128,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00128): "For dml, Using having clause are prohibited".
You should follow the following logic:
1. For "SELECT..." Statement, checks for a HAVING clause in the sentence and, if present, reports a rule violation.
2. For "INSERT..." Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3, for "DELETE..." Statement, the SELECT clause in the DELETE statement is checked the same way as above.
4. For "UPDATE..." Statement, perform the same checks as above on the SELECT clause in the UPDATE statement.
5. For "UNION..." Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00128(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.DeleteStmt, *ast.UpdateStmt:
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// "select..."
			if selectStmt.Having != nil {
				// "select..." with "having"
				rulepkg.AddResult(input.Res, input.Rule, SQLE00128)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
