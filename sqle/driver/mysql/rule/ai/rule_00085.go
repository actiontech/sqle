package ai

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00085 = "SQLE00085"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00085,
			Desc:       plocale.Rule00085Desc,
			Annotation: plocale.Rule00085Annotation,
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
			Version:      2,
		},
		Message: plocale.Rule00085Message,
		Func:    RuleSQLE00085,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00085): "For dml, Full index scans of tables are prohibited".
You should follow the following logic:
1. For "select..." statement, check the output of the execution plan. If type is index, report a violation. The execution plan should be the information obtained online.
2. For "union ..." statement, perform the same check as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00085(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt:
		// "select..."
		executionPlan, err := util.GetExecutionPlan(input.Ctx, stmt.Text())
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)
			return nil
		}

		for _, record := range executionPlan.Plan {
			if record.Type == executor.ExplainRecordAccessTypeIndex {
				// full index scan
				rulepkg.AddResult(input.Res, input.Rule, SQLE00085)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
