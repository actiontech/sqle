package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00099 = "SQLE00099"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00099,
			Desc:       plocale.Rule00099Desc,
			Annotation: plocale.Rule00099Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
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
		Message: plocale.Rule00099Message,
		Func:    RuleSQLE00099,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00099): "For MySQL DQL, SELECT FOR UPDATE is prohibited.".
You should follow the following logic:
1. For "select..." Statement, checks FOR the presence of an FOR UPDATE clause in the statement, If it does, report a violation.
2. For "insert... "Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For "union..." Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00099(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt:
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// "select..."
			if selectStmt.LockTp == ast.SelectLockForUpdate {
				//"select..." with "FOR UPDATE"
				rulepkg.AddResult(input.Res, input.Rule, SQLE00099)
				return nil
			}
		}

	}
	return nil
}

// ==== Rule code end ====
