package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00071 = "SQLE00071"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00071,
			Desc:         plocale.Rule00071Desc,
			Annotation:   plocale.Rule00071Annotation,
			Category:     plocale.RuleTypeDDLConvention,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00071Message,
		Func:    RuleSQLE00071,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00071): "In DDL, Deleting columns is prohibited".
You should follow the following logic:
1. For "alter table ... drop column ..." statement, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00071(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// "alter table ..."
		for _, spec := range stmt.Specs {
			if util.IsAlterTableCommand(spec, ast.AlterTableDropColumn) {
				// "alter table ... drop column ..."
				rulepkg.AddResult(input.Res, input.Rule, SQLE00071)
				return nil
			}
		}
	default:
		return nil
	}
	return nil
}

// ==== Rule code end ====
