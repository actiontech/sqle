package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00179 = "SQLE00179"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00179,
			Desc:       plocale.Rule00179Desc,
			Annotation: plocale.Rule00179Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID, plocale.RuleTagCorrection.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00179Message,
		Func:    RuleSQLE00179,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00179): "For dml, SQL queries with implicit data type conversions are prohibited".
You should follow the following logic:
1. For "select..." Statement, check the warning information of the execution plan of the SQL statement, and report the violation of the rule if the keyword "due to type or collation conversion on field" appears in the warning information. The warning information of the execution plan should be the information obtained online.
2. For "update..." Statement, performs the same checking process as above.
3. For "delete..." Statement, performs the same checking process as above.
4. For "insert..." Statement, performs the same checking process as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00179(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
		// "select...", "insert...", "update...", "delete..."
		plan, err := util.GetExecutionPlan(input.Ctx, stmt.Text())
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)
			return nil
		}
		for _, warning := range plan.Warnings {
			// "due to type or collation conversion on field"
			if strings.Contains(warning.Message, "due to type or collation conversion on field") {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00179)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
