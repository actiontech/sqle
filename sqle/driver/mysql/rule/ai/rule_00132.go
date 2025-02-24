package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00132 = "SQLE00132"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00132,
			Desc:       plocale.Rule00132Desc,
			Annotation: plocale.Rule00132Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID, plocale.RuleTagQuery.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00132Message,
		Func:    RuleSQLE00132,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00132): "For dml, Using subqueries are prohibited".
You should follow the following logic:
1. For "select..." The statement, checks if a SELECT subquery exists in the sentence, and if so, reports a rule violation
2. For "union..." Statement, perform the same checking process as above
3. For "update..." Statement, perform the same checking process as above
4. For "insert..." Statement, perform the same checking process as above
5. For "delete..." Statement, perform the same checking process as above
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00132(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
		if len(util.GetSubquery(stmt)) > 0 {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00132)
		}
		return nil
	}
	return nil
}

// ==== Rule code end ====
