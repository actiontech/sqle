package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00128 = "SQLE00128"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00128,
			Desc:       plocale.Rule00128Desc,
			Annotation: plocale.Rule00128Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00128Message,
		Func:    RuleSQLE00128,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
