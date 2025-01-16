package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00004(t *testing.T) {
	ruleName := ai.SQLE00004
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE with AUTO_INCREMENT column starting at 0",
		"CREATE TABLE test_table (id INT AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT=0;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE TABLE with AUTO_INCREMENT column starting at 1",
		"CREATE TABLE test_table (id INT AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT=1;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: CREATE TABLE without AUTO_INCREMENT column",
		"CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: CREATE TABLE with multiple columns, one AUTO_INCREMENT starting at 0",
		"CREATE TABLE test_table (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(50)) AUTO_INCREMENT=0;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: CREATE TABLE with multiple AUTO_INCREMENT columns, one starting at 1",
		"CREATE TABLE test_table (id INT AUTO_INCREMENT PRIMARY KEY, seq INT AUTO_INCREMENT) AUTO_INCREMENT=1;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: CREATE TABLE with AUTO_INCREMENT column without specifying AUTO_INCREMENT option",
		"CREATE TABLE test_table (id INT AUTO_INCREMENT PRIMARY KEY);",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: SET auto_increment_offset to 0",
		"SET auto_increment_offset = 0;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: SET auto_increment_offset to 1",
		"SET auto_increment_offset = 1;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: SET auto_increment_increment to 1",
		"SET auto_increment_increment = 1;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: SET some_other_variable to 10",
		"SET some_other_variable = 10;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 11: CREATE TABLE with AUTO_INCREMENT column starting at 3",
		"CREATE TABLE t1 (id INT AUTO_INCREMENT PRIMARY KEY, c1 INT) AUTO_INCREMENT = 3;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: CREATE TABLE with AUTO_INCREMENT column starting at 100",
		"CREATE TABLE t1 (id INT AUTO_INCREMENT PRIMARY KEY, c1 INT) AUTO_INCREMENT = 100;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: CREATE TABLE with AUTO_INCREMENT column starting at INT_MAX",
		"CREATE TABLE t2 (id INT AUTO_INCREMENT PRIMARY KEY, c1 INT) AUTO_INCREMENT = 2147483647;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SET auto_increment_offset to 100",
		"SET @@auto_increment_offset = 100;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: CREATE TABLE with default AUTO_INCREMENT",
		"CREATE TABLE t1 (id INT AUTO_INCREMENT PRIMARY KEY, c1 INT);",
		nil, /*mock context*/
		nil, newTestResult())
}

// ==== Rule test code end ====
