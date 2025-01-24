package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00035(t *testing.T) {
	ruleName := ai.SQLE00035
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//without full-width Chinese quotation marks
	runSingleRuleInspectCase(rule, t, "without full-width Chinese quotation marks", DefaultMysqlInspect(), `
			alter table exist_tb_1 add column a int comment 'a'
			`, newTestResult())

	//with full-width Chinese quotation marks in comment
	runSingleRuleInspectCase(rule, t, "with full-width Chinese quotation marks in comment", DefaultMysqlInspect(), `
			alter table exist_tb_1 add column a int comment '”a“'
			`, newTestResult().addResult(ruleName))

	//with full-width Chinese quotation marks in column name
	runSingleRuleInspectCase(rule, t, "with full-width Chinese quotation marks in column name", DefaultMysqlInspect(),
		"alter table exist_tb_1 add column `”a“` int comment 'a'",
		newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
