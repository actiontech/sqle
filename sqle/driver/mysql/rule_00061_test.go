package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00061(t *testing.T) {
	ruleName := ai.SQLE00061
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE not included IF NOT EXISTS", "CREATE TABLE test_table (id INT);",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: CREATE TABLE included IF NOT EXISTS", "CREATE TABLE IF NOT EXISTS test_table (id INT);",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: CREATE TEMPORARY TABLE included IF NOT EXISTS", "CREATE TEMPORARY TABLE IF NOT EXISTS temp_table (id INT);",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: CREATE TEMPORARY TABLE not included IF NOT EXISTS", "CREATE TEMPORARY TABLE temp_table (id INT);",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
