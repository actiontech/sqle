package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00072 = "SQLE00072"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00072,
			Desc:         plocale.Rule00072Desc,
			Annotation:   plocale.Rule00072Annotation,
			Category:     plocale.RuleTypeDDLConvention,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00072Message,
		Func:    RuleSQLE00072,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00072): "In DDL, deleting foreign keys is prohibited".
You should follow the following logic:
1. For "alter table ... drop foreign key ..." statement, report a violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00072(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		for range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableDropForeignKey) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00072)
		}
	default:
		return nil
	}
	return nil
}

// ==== Rule code end ====
