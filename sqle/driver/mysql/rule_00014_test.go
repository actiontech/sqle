package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQLE00014(t *testing.T) {
	ruleName := ai.SQLE00014
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: 使用 CREATE FUNCTION 定义自定义函数",
		"CREATE FUNCTION my_func (a INT) RETURNS INT DETERMINISTIC RETURN a + 1;",
		nil, /*mock context*/
		nil,
		newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 2: 使用 CREATE FUNCTION 定义自定义函数（不同大小写）",
		"create function my_func (a INT) returns INT deterministic return a + 1;",
		nil, /*mock context*/
		nil,
		newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName),
	)

}

// ==== Rule test code end ====
