package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00101 = "SQLE00101"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00101,
			Desc:       plocale.Rule00101Desc,
			Annotation: plocale.Rule00101Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00101Message,
		Func:    RuleSQLE00101,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
