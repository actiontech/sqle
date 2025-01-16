package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00123 = "SQLE00123"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00123,
			Desc:       "禁止使用TRUNCATE操作",
			Annotation: "TRUNCATE是DDL，执行后数据默认隐式提交，无法回滚，在没有备份的场景下，谨慎使用TRUNCATE",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeIndexInvalidation,
		},
		Message: "禁止使用TRUNCATE操作.",
		AllowOffline: true,
		Func:    RuleSQLE00123,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00123): "For dml, using truncate is prohibited".
You should follow the following logic:
1. "For TRUNCATE..." Statement, report a rule violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00123(input *rulepkg.RuleHandlerInput) error {

	switch input.Node.(type) {
	case *ast.TruncateTableStmt:
		// "For TRUNCATE TABLE..."
		rulepkg.AddResult(input.Res, input.Rule, SQLE00123)
	}
	return nil
}

// ==== Rule code end ====
