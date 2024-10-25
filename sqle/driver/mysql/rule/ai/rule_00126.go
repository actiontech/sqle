package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00126 = "SQLE00126"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00126,
			Desc:       "对于MySQL的DML, 不建议对字段编号进行 GROUP BY",
			Annotation: "GROUP BY 1 表示按第一列进行GROUP BY；在GROUP BY子句中使用字段编号，而不是表达式或列名称，当查询列顺序改变时，会导致查询逻辑出现问题",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "对于MySQL的DML, 不建议对字段编号进行 GROUP BY",
		AllowOffline: true,
		Func:    RuleSQLE00126,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00126): "For dml, using position numbers in GROUP BY field is prohibited".
You should follow the following logic:
1. For "SELECT..." The statement,
  * Check if there is a GROUP by clause in the sentence, and if so, check it further.
  * Check if the content of Group By is a numeric number representing the location, if so, report a rule violation.
2. For "INSERT..." Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For "DELETE..." Statement, the SELECT clause in the DELETE statement is checked the same way as above.
4. For "UPDATE..." Statement, perform the same checks as above on the SELECT clause in the UPDATE statement.
5. For "UNION..." Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00126(input *rulepkg.RuleHandlerInput) error {

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.DeleteStmt, *ast.UnionStmt, *ast.UpdateStmt:
		// "SELECT..." "INSERT..." "DELETE..." "UNION..." "UPDATE..."
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// Apply the check to each select statement in a "union"
			groupByClause := selectStmt.GroupBy
			if groupByClause != nil {
				// Check if the GROUP BY clause is using numbers
				for _, item := range groupByClause.Items {
					if _, ok := item.Expr.(*ast.PositionExpr); ok {
						// the expr is a position number
						rulepkg.AddResult(input.Res, input.Rule, SQLE00126)
						return nil
					}
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
