package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00052(t *testing.T) {
	ruleName := ai.SQLE00052
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// Case 1: CREATE TABLE 主键未使用 AUTO_INCREMENT
	runAIRuleCase(rule, t, "case 1: CREATE TABLE 主键未使用 AUTO_INCREMENT",
		"CREATE TABLE test_table (id INT PRIMARY KEY);",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	// Case 2: CREATE TABLE 主键使用 AUTO_INCREMENT
	runAIRuleCase(rule, t, "case 2: CREATE TABLE 主键使用 AUTO_INCREMENT",
		"CREATE TABLE test_table (id INT PRIMARY KEY AUTO_INCREMENT);",
		nil, /*mock context*/
		nil, newTestResult())

	// Case 3: ALTER TABLE 添加主键未使用 AUTO_INCREMENT
	runAIRuleCase(rule, t, "case 3: ALTER TABLE 添加主键未使用 AUTO_INCREMENT",
		"ALTER TABLE test_table ADD COLUMN id INT PRIMARY KEY;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// Case 4: ALTER TABLE 添加主键使用 AUTO_INCREMENT
	runAIRuleCase(rule, t, "case 4: ALTER TABLE 添加主键使用 AUTO_INCREMENT",
		"ALTER TABLE test_table ADD COLUMN id INT PRIMARY KEY AUTO_INCREMENT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (name VARCHAR(100));"),
		nil, newTestResult())

	// Case 5: ALTER TABLE 修改主键未使用 AUTO_INCREMENT
	runAIRuleCase(rule, t, "case 5: ALTER TABLE 修改主键未使用 AUTO_INCREMENT",
		"ALTER TABLE test_table MODIFY COLUMN id INT PRIMARY KEY;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult())

	// Case 6: ALTER TABLE 修改主键使用 AUTO_INCREMENT
	runAIRuleCase(rule, t, "case 6: ALTER TABLE 修改主键使用 AUTO_INCREMENT",
		"ALTER TABLE test_table MODIFY COLUMN id INT PRIMARY KEY AUTO_INCREMENT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult())

	// Case 7: CREATE TABLE 非自增主键表(从xml中补充)
	runAIRuleCase(rule, t, "case 7: CREATE TABLE 非自增主键表(从xml中补充)",
		"CREATE TABLE customers_no_auto_increment (id CHAR(36) PRIMARY KEY, name VARCHAR(32) DEFAULT '' NOT NULL, sex INT DEFAULT 0, age INT DEFAULT 0, mark1 VARCHAR(20) NOT NULL, mark2 VARCHAR(30) NOT NULL, KEY idx_name_customers_no_auto_increment (name));",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	// Case 8: CREATE TABLE 自增主键表(从xml中补充)
	runAIRuleCase(rule, t, "case 8: CREATE TABLE 自增主键表(从xml中补充)",
		"CREATE TABLE customers_auto_increment (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(32) DEFAULT '' NOT NULL, sex INT DEFAULT 0, age INT DEFAULT 0, mark1 VARCHAR(20) NOT NULL, mark2 VARCHAR(30) NOT NULL, KEY idx_name_customers_auto_increment (name));",
		nil, /*mock context*/
		nil, newTestResult())

	// Case 9: ALTER TABLE 添加非自增主键(从xml中补充)
	runAIRuleCase(rule, t, "case 9: ALTER TABLE 添加非自增主键(从xml中补充)",
		"ALTER TABLE customers ADD PRIMARY KEY(id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
