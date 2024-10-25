package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00073 = "SQLE00073"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00073,
			Desc:       "对于MySQL的DDL, 不建议修改表的默认字符集",
			Annotation: "修改表的默认字符集，只会影响后续新增的字段，不会修表已有字段的字符集；如需修改整张表所有字段的字符集建议开启此规则",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 不建议修改表的默认字符集",
		AllowOffline: true,
		Func:    RuleSQLE00073,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
