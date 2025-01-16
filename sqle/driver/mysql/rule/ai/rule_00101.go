package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00101 = "SQLE00101"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00101,
			Desc:       "SELECT 语句不能有ORDER BY",
			Annotation: "ORDER BY 对查询性能影响较大，建议将排序部分放到业务处理。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "SELECT 语句不能有ORDER BY",
		AllowOffline: false,
		Func:    RuleSQLE00101,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00101): "For MySQL DQL, SELECT statements with ORDER BY is prohibited.".
You should follow the following logic:
1. For "select..." Statement, check if there is an ORDER BY clause in the statement, and if so, report a rule violation.
2. For "insert... "Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For "union..." Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00101(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt:
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// "select..." or "union..."
			if selectStmt.OrderBy != nil {
				// "select..." with "ORDER BY"
				rulepkg.AddResult(input.Res, input.Rule, SQLE00101)
				return nil
			}
		}

	}
	return nil
}

// ==== Rule code end ====
