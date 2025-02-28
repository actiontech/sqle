package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00123 = "SQLE00123"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00123,
			Desc:       plocale.Rule00123Desc,
			Annotation: plocale.Rule00123Annotation,
			Category:   plocale.RuleTypeIndexInvalidation,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID, plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagSecurity.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00123Message,
		Func:    RuleSQLE00123,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00123): "For dml, using truncate is prohibited".
You should follow the following logic:
1. "For TRUNCATE..." Statement, report a rule violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00123(input *rulepkg.RuleHandlerInput) error {

	switch input.Node.(type) {
	case *ast.TruncateTableStmt:
		// "For TRUNCATE TABLE..."
		rulepkg.AddResult(input.Res, input.Rule, SQLE00123)
	}
	return nil
}

// ==== Rule code end ====
