package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00072 = "SQLE00072"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00072,
			Desc:       "对于MySQL的DDL, 禁止进行删除外键的操作",
			Annotation: "删除已有约束会影响已有业务逻辑；开启该规则，SQLE将提醒删除外键为高危操作",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 禁止进行删除外键的操作",
		AllowOffline: true,
		Func:    RuleSQLE00072,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
