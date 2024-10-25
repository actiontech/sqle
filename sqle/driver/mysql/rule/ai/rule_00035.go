package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00035 = "SQLE00035"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00035,
			Desc:       "对于MySQL的DDL, DDL语句中不建议使用中文全角引号",
			Annotation: "建议开启此规则，可避免MySQL会将中文全角引号识别为命名的一部分，执行结果与业务预期不符",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, DDL语句中不建议使用中文全角引号",
		AllowOffline: true,
		Func:    RuleSQLE00035,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00035): "In DDL, using full-width Chinese quotation marks in DDL statements is prohibited".
You should follow the following logic:
1. For DDL statement, check the sql text, report a violation if it has full-width Chinese quotation marks.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00035(input *rulepkg.RuleHandlerInput) error {
	switch input.Node.(type) {
	case ast.DDLNode:
		if strings.Contains(input.Node.Text(), "“") {
			rulepkg.AddResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

// ==== Rule code end ====
