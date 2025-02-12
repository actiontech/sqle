package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00131 = "SQLE00131"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00131,
			Desc:       plocale.Rule00131Desc,
			Annotation: plocale.Rule00131Annotation,
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
		Message: plocale.Rule00131Message,
		Func:    RuleSQLE00131,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00131): "For dml, using ORDER BY RAND() is prohibited".
You should follow the following logic:
1. For SELECT... Order BY... Statement, checks whether the Order By clause contains a RAND function, and if so, report a violation.
2. For INSERT... Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For UNION... Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00131(input *rulepkg.RuleHandlerInput) error {

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UnionStmt:
		// "SELECT...", "INSERT...", "UNION..."
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// Apply the check to each order by items
			orderBy := selectStmt.OrderBy
			if orderBy != nil {
				for _, item := range orderBy.Items {
					if item.Expr == nil {
						continue
					}

					// Check for function calls in the expr.
					funcNames := util.GetFuncName(item.Expr)
					for _, name := range funcNames {
						// Check if the function is the one we want to check.
						if name == "rand" {
							// Add a rule violation result if an issue is found.
							rulepkg.AddResult(input.Res, input.Rule, SQLE00131)
						}
					}
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
