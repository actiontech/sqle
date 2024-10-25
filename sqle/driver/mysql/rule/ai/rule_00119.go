package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00119 = "SQLE00119"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00119,
			Desc:       "对于MySQL的DML, 建议为GROUP BY语句添加ORDER BY条件",
			Annotation: "在5.7中，MySQL默认会对’GROUP BY col1, …’按如下顺序’ORDER BY col1,…’隐式排序，导致产生无谓的排序，带来额外的开销，影响SQL执行效率；在8.0中，则不会出现这种情况。如果不需要排序建议显示添加’ORDER BY NULL’",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "对于MySQL的DML, 建议为GROUP BY语句添加ORDER BY条件.",
		AllowOffline: true,
		Func:    RuleSQLE00119,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00119): "For dml, It is recommended to add an ORDER BY condition to the GROUP BY statement".
You should follow the following logic:
1. For "SELECT..." The statement,
  1. Check if the group by keyword is present in the sentence, and if so, proceed to the next step.
  2. Check if there is an order by keyword in the sentence, if not, report a rule violation.
2. For INSERT... Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For UNION... Statement, does the same check as above for each SELECT clause in the statement.
4. For UPDATE... Statement, the same checks as above are performed for the sub-queries in the statement.
5. For DELETE... Statement, the same checks as above are performed for the sub-queries in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00119(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UnionStmt, *ast.UpdateStmt, *ast.DeleteStmt:
		// "SELECT..." "INSERT..." "UNION..." "UPDATE..." "DELETE..."
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// Apply the check to each select statement in a "union"
			groupByClause := selectStmt.GroupBy
			if groupByClause != nil && len(groupByClause.Items) > 0 {
				// "GROUP BY" is present in the SQL statement
				orderByClause := selectStmt.OrderBy
				if orderByClause != nil && len(orderByClause.Items) > 0 {
					// "ORDER BY" is present in the SQL statement
					return nil
				}
				// "ORDER BY" is not present in the SQL statement
				rulepkg.AddResult(input.Res, input.Rule, SQLE00119)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
