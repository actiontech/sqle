package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00073 = "SQLE00073"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00073,
			Desc:         plocale.Rule00073Desc,
			Annotation:   plocale.Rule00073Annotation,
			Category:     plocale.RuleTypeDDLConvention,
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00073Message,
		Func:    RuleSQLE00073,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00073): "In DDL, Changing default character set of a table is prohibited".
You should follow the following logic:
1. For "ALTER TABLE ... CONVERT TO CHARACTER SET ..." statement, report a violation.
2. For "ALTER TABLE ... CHARACTER SET ..." statement, report a violation.
3. For "ALTER TABLE ... COLLATE ..." statement, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00073(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// "alter table"
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableOption) {
			// "alter table option"
			if util.IsAlterTableCommandAlterOption(spec, ast.TableOptionCharset) {
				// "ALTER TABLE ... CONVERT TO CHARACTER SET ..." or "ALTER TABLE ... CHARACTER SET ..."
				rulepkg.AddResult(input.Res, input.Rule, input.Rule.Name)
			}

			if util.IsAlterTableCommandAlterOption(spec, ast.TableOptionCollate) {
				// "ALTER TABLE ... COLLATE ..."
				rulepkg.AddResult(input.Res, input.Rule, input.Rule.Name)
			}
		}
	}
	return nil
}

// ==== Rule code end ====
