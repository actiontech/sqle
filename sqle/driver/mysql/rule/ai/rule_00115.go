package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00115 = "SQLE00115"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00115,
			Desc:       plocale.Rule00115Desc,
			Annotation: plocale.Rule00115Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID, plocale.RuleTagQuery.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00115Message,
		Func:    RuleSQLE00115,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00115): "For dml, using scalar subquery is prohibited".
You should follow the following logic:
1. For "SELECT..." Statement, checks if a SELECT clause exists in the sentence, and if so, then
  1. Check if a subquery is used to retrieve the fields in the clause, and if so, determine the number of fields that the subquery retrieves, and if there is only one, report a rule violation
  2. Check if there is a WHERE condition in the sentence, if there is a WHERE condition, then determine whether the subquery is used in the WHERE condition, if there is a use, then determine the number of fields obtained by the subquery, if there is only one, then report the rule violation
2. For "UNION..." Statement, does the same check as above for each SELECT clause in the statement.
3. For "DELETE..." Statement, the same checks as above are performed for the sub-queries in the statement.
   In addition, it checks whether the subquery is used in the WHERE condition, and if so, it determines the number of fields obtained by the subquery, and if the number is only one, it reports a violation of the rule
4. For "INSERT..." Statement, the same checks as above are performed for the sub-queries in the statement.
   In addition, it checks whether the subquery is used in the values part or in the SET clauses of the INSERT statement is not, and if so, it determines the number of fields obtained by the subquery, and if the number is only one, it reports a violation of the rule.
5. For "UPDATE..." Statement, the same checks as above are performed for the sub-queries in the statement.
   In addition, it checks whether the subquery is used in the SET clause, and if so, it determines the number of fields obtained by the subquery, and if the number is only one, it reports a violation of the rule
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00115(input *rulepkg.RuleHandlerInput) error {
	checkExprScalarSubquery := func(expr ast.ExprNode) {
		for _, subquery := range util.GetSubquery(expr) {
			// the expr is a subquery
			// Check if the number of fields fetched by the subquery is only one
			for _, s := range util.GetSelectStmt(subquery.Query) {
				if len(s.Fields.Fields) == 1 {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00115)
					return
				}
			}
		}
	}

	checkWhereScalarSubquery := func(where ast.ExprNode) {
		// check scalar subquery for the where condition
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			if whereExpr, ok := expr.(*ast.BinaryOperationExpr); ok {
				// In subquery is not a scalar subquery
				if whereExpr.Op != opcode.In {
					checkExprScalarSubquery(expr)
				}
			}
			return false
		}, where)

	}

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.DeleteStmt, *ast.UpdateStmt:
		// "SELECT...", "UNION...", "INSERT...", "DELETE...", "UPDATE..."
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// "SELECT..."

			// check scalar subquery for the field fetched in the select
			if selectStmt.Fields != nil {

				for _, field := range selectStmt.Fields.Fields {
					checkExprScalarSubquery(field.Expr)
				}
			}

			// check scalar subquery for the where condition
			if selectStmt.Where != nil {
				checkWhereScalarSubquery(selectStmt.Where)
			}
		}
	}

	switch stmt := input.Node.(type) {
	case *ast.InsertStmt:
		// "INSERT..."

		// checks whether the scalar subquery is used in the values
		for _, list := range stmt.Lists {
			for _, expr := range list {
				checkExprScalarSubquery(expr)
			}
		}

		// checks whether the scalar subquery is used in the SET clause
		for _, assignment := range stmt.Setlist {
			checkExprScalarSubquery(assignment.Expr)
		}

	case *ast.DeleteStmt:
		// "DELETE..."

		// check scalar subquery for the where condition
		if stmt.Where != nil {
			checkWhereScalarSubquery(stmt.Where)
		}
	case *ast.UpdateStmt:
		// "UPDATE..."

		// checks whether the scalar subquery is used in the SET clause
		for _, assignment := range stmt.List {
			checkExprScalarSubquery(assignment.Expr)
		}

	}
	return nil

}

// ==== Rule code end ====
