package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00007 = "SQLE00007"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00007,
			Desc:       "对于MySQL的DDL, 自增字段只能设置一个",
			Annotation: "MySQL InnoDB，MyISAM 引擎不允许存在多个自增字段，设置多个自增字段会导致上线失败。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 自增字段只能设置一个",
		AllowOffline: true,
		Func:    RuleSQLE00007,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00007): "In DDL, when creating table, You can only define a maximum of one auto-increment column".
You should follow the following logic:
1. For "create table ..." statement, check the count of auto-increment columns, if the count is more than 1, report violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00007(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table"
		found := 0

		for _, col := range stmt.Cols {
			if util.IsColumnAutoIncrement(col) {
				found++
			}
		}

		if found > 1 {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00007)
			return nil
		}
	}
	return nil
}

// ==== Rule code end ====
