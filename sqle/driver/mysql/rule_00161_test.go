package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00161(t *testing.T) {
	ruleName := ai.SQLE00161
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SET auto_increment_increment 设置步长为 2，违反规则",
		"SET @@auto_increment_increment = 1,@@read_only=true;",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 1: SET auto_increment_increment 设置步长为 2，违反规则",
		"SET @@read_only=true,@@auto_increment_increment = 2;",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)
}

// ==== Rule test code end ====
