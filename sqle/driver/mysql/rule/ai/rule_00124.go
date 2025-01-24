package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00124 = "SQLE00124"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00124,
			Desc:         plocale.Rule00124Desc,
			Annotation:   plocale.Rule00124Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00124Message,
		Func:    RuleSQLE00124,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00124): "For dml, using DELETE to delete an entire table is prohibited".
You should follow the following logic:
1. For "DELETE... Statement, checks if there is no WHERE condition or where condition is always True (for example, where 1=1 or where True), reports a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00124(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.DeleteStmt:
		// "delete"
		aliasInfos := util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		if stmt.Where == nil || util.IsExprConstTrue(input.Ctx, stmt.Where, aliasInfos) {
			// "delete...where..."
			rulepkg.AddResult(input.Res, input.Rule, SQLE00124)
			return nil
		}
	}
	return nil

}

// ==== Rule code end ====
