package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00051(t *testing.T) {
	ruleName := ai.SQLE00051
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//ALTER TABLE customers ADD COLUMN id2 INT AUTO_INCREMENT PRIMARY KEY; -- 添加列，同时追加 AUTO_INCREMENT、PRIMARY KEY
	runAIRuleCase(rule, t, "case 0: ALTER TABLE add PRIMARY KEY without AUTO_INCREMENT to existing column",
		"ALTER TABLE customers ADD COLUMN id2 INT AUTO_INCREMENT PRIMARY KEY",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(32), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	//ALTER TABLE customers ADD PRIMARY KEY (id);
	runAIRuleCase(rule, t, "case 0: ALTER TABLE add PRIMARY KEY without AUTO_INCREMENT to existing column",
		"ALTER TABLE customers ADD PRIMARY KEY (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(32), age INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 1: CREATE TABLE with PRIMARY KEY and AUTO_INCREMENT",
		"CREATE TABLE test_table (id INT PRIMARY KEY AUTO_INCREMENT);",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: CREATE TABLE with PRIMARY KEY without AUTO_INCREMENT",
		"CREATE TABLE test_table (id INT PRIMARY KEY);",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 3: ALTER TABLE to add PRIMARY KEY with AUTO_INCREMENT",
		"ALTER TABLE test_table ADD COLUMN id INT PRIMARY KEY AUTO_INCREMENT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (name VARCHAR(32));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: ALTER TABLE to add PRIMARY KEY without AUTO_INCREMENT",
		"ALTER TABLE test_table ADD COLUMN id INT PRIMARY KEY;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (name VARCHAR(32));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 5: ALTER TABLE existing table with AUTO_INCREMENT primary key",
		"ALTER TABLE test_table MODIFY id INT AUTO_INCREMENT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY );"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: ALTER TABLE existing table without AUTO_INCREMENT primary key",
		"ALTER TABLE test_table MODIFY id INT PRIMARY KEY;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 7: CREATE TABLE with multiple columns, one PRIMARY KEY with AUTO_INCREMENT",
		"CREATE TABLE customers (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(32), age INT);",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: CREATE TABLE with multiple columns, PRIMARY KEY without AUTO_INCREMENT",
		"CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(32), age INT);",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 9: CREATE TABLE with multiple columns, PRIMARY KEY without AUTO_INCREMENT",
		"CREATE TABLE customers (id INT AUTO_INCREMENT, name VARCHAR(32), age INT, PRIMARY KEY(id));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: ALTER TABLE add PRIMARY KEY without AUTO_INCREMENT to existing column",
		"ALTER TABLE customers ADD PRIMARY KEY (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(32), age INT);"),
		nil,
		newTestResult())
}

// ==== Rule test code end ====
