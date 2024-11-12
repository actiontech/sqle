package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00071 = "SQLE00071"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00071,
			Desc:       "禁止进行删除列的操作",
			Annotation: "业务逻辑与删除列依赖未完全消除，列被删除后可能导致程序异常（无法正常读写）的情况；开启该规则，SQLE将提醒删除列为高危操作",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "禁止进行删除列的操作",
		AllowOffline: true,
		Func:    RuleSQLE00071,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
