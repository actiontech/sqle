package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00178 = "SQLE00178"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00178,
			Desc:       plocale.Rule00178Desc,
			Annotation: plocale.Rule00178Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00178Message,
		Func:    RuleSQLE00178,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00178): "For dml, Full table sort are prohibited".
You should follow the following logic:
1. For SELECT... Statement, check If there is no WHERE condition or where condition is always True (such as where 1=1 or where True), and the SQL statement has an ORDER BY or GROUP BY or DISTINCT clause, report a violation.
2. For UNION ... statement, perform the same check as mentioned above for each SELECT statement within the UNION.
3. For INSERT... Statement, if the INSERT statement inserts data from a SELECT query that matches the SELECT statement firing rules, report a violation.
4. For DELETE... Statement, if the DELETE statement has no WHERE condition, or WHERE condition is always true, and involves an ORDER BY clause, report a violation.
5. For UPDATE... Statement, if the UPDATE statement has no WHERE condition, or WHERE condition is always true, and involves an ORDER BY clause, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00178(input *rulepkg.RuleHandlerInput) error {
	isSelectStmtViolation := func(stmt *ast.SelectStmt) bool {
		aliasInfo := util.GetTableAliasInfoFromJoin(stmt.From.TableRefs)
		if stmt.Where == nil || util.IsExprConstTrue(input.Ctx, stmt.Where, aliasInfo) {
			// where is nil or where is always true
			if stmt.OrderBy != nil || stmt.GroupBy != nil || stmt.Distinct {
				// with order by or group by or distinct
				return true
			}
		}
		return false
	}
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		// "select"
		if isSelectStmtViolation(stmt) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00178)
			return nil

		}
	case *ast.UnionStmt:
		// "union"
		for _, selectStmt := range util.GetSelectStmt(stmt) {

			if isSelectStmtViolation(selectStmt) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00178)
				return nil
			}
		}
	case *ast.InsertStmt:
		// "insert"
		// check if the INSERT statement inserts data from a SELECT query that matches the SELECT statement firing rules
		for _, selectStmt := range util.GetSelectStmt(stmt) {

			if isSelectStmtViolation(selectStmt) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00178)
				return nil
			}
		}
	case *ast.DeleteStmt:
		// "delete"
		aliasInfos := util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		if stmt.Where == nil || util.IsExprConstTrue(input.Ctx, stmt.Where, aliasInfos) {
			// where is nil or where is always true
			if stmt.Order != nil {
				// with order
				rulepkg.AddResult(input.Res, input.Rule, SQLE00178)
				return nil
			}
		}
	case *ast.UpdateStmt:
		// "update"
		aliasInfos := util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		if stmt.Where == nil || util.IsExprConstTrue(input.Ctx, stmt.Where, aliasInfos) {
			// where is nil or where is always true
			if stmt.Order != nil {
				// with order
				rulepkg.AddResult(input.Res, input.Rule, SQLE00178)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
