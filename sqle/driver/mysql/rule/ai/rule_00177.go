package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00177 = "SQLE00177"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00177,
			Desc:       plocale.Rule00177Desc,
			Annotation: plocale.Rule00177Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "3",
				Desc:  plocale.Rule00177Params1,
				Type:  params.ParamTypeInt,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00177Message,
		Func:    RuleSQLE00177,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00177): "For dml, the number of Order By fields should not exceed the specified threshold".the threshold should be a parameter whose default value is 3.
You should follow the following logic:
1. For "select..." The statement,
* Check if the ORDER BY keyword is present in the SQL statement, if it is, proceed to the next step
* Count the number of fields in the ORDER BY clause and report a rule violation if the number of fields exceeds a threshold
2. For INSERT... Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For UNION... Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00177(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxOrderByCount := param.Int()

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UnionStmt:
		// "select", "insert", "union"
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			if selectStmt.OrderBy == nil {
				continue
			}

			if len(selectStmt.OrderBy.Items) > maxOrderByCount {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00177, maxOrderByCount)
			}
		}
	}
	return nil
}

// ==== Rule code end ====
