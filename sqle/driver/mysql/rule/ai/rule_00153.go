package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00153 = "SQLE00153"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00153,
			Desc:       plocale.Rule00153Desc,
			Annotation: plocale.Rule00153Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00153Message,
		Func:    RuleSQLE00153,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00153): "In table definition, secondary index must be used when create table".
You should follow the following logic:
1. For "create table ..." statement, check if there is any secondary index definition, and if not, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00153(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table..."
		found := false

		// check if secondary index is defined
		constraints := util.GetTableConstraints(stmt.Constraints,
			ast.ConstraintIndex,
			ast.ConstraintUniqIndex,
			ast.ConstraintKey,
			ast.ConstraintUniq,
			ast.ConstraintUniqKey,
		)
		if len(constraints) > 0 {
			found = true
		}

		// check if index is defined in column definition
		for _, col := range stmt.Cols {
			if util.IsColumnHasOption(col, ast.ColumnOptionUniqKey) {
				found = true
			}
		}

		if !found {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00153)
		}
	}
	return nil
}

// ==== Rule code end ====
