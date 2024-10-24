package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00054(t *testing.T) {
	ruleName := ai.SQLE00054
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// CREATE语句
	runAIRuleCase(rule, t, "case 1: CREATE TABLE with BIGINT UNSIGNED as primary key",
		"CREATE TABLE user_table (id BIGINT UNSIGNED PRIMARY KEY, name VARCHAR(100));",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE TABLE with BIGINT signed as primary key",
		"CREATE TABLE user_table (id BIGINT PRIMARY KEY, name VARCHAR(100));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: CREATE TABLE with INT as primary key",
		"CREATE TABLE user_table (id INT PRIMARY KEY, name VARCHAR(100));",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: CREATE TABLE without primary key",
		"CREATE TABLE user_table (id BIGINT UNSIGNED, name VARCHAR(100));",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: CREATE TABLE with composite primary key without BIGINT signed",
		"CREATE TABLE user_table (id BIGINT UNSIGNED, code VARCHAR(50), PRIMARY KEY (id, code));",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: CREATE TABLE with composite primary key including BIGINT signed",
		"CREATE TABLE user_table (id BIGINT, code VARCHAR(50), PRIMARY KEY (id, code));",
		nil, nil, newTestResult().addResult(ruleName))

	// ALTER语句
	runAIRuleCase(rule, t, "case 7: ALTER TABLE to add primary key with BIGINT UNSIGNED",
		"ALTER TABLE user_table ADD PRIMARY KEY (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_table (id BIGINT UNSIGNED, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: ALTER TABLE to modify primary key to BIGINT UNSIGNED",
		"ALTER TABLE user_table MODIFY COLUMN id BIGINT UNSIGNED PRIMARY KEY;",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_table (id BIGINT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: ALTER TABLE to modify primary key to BIGINT signed",
		"ALTER TABLE user_table MODIFY COLUMN id BIGINT PRIMARY KEY;",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_table (id BIGINT UNSIGNED, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: ALTER TABLE to add primary key with INT",
		"ALTER TABLE user_table ADD PRIMARY KEY (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_table (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 12: ALTER TABLE without changing primary key",
		"ALTER TABLE user_table ADD COLUMN email VARCHAR(100);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_table (id BIGINT UNSIGNED PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult())

	// * 新增示例
	runAIRuleCase(rule, t, "case 13: CREATE TABLE with BIGINT signed as primary key (从xml中补充)",
		"CREATE TABLE customers (id BIGINT NOT NULL PRIMARY KEY, name VARCHAR(32) DEFAULT '' NOT NULL, sex INT DEFAULT 0, age INT DEFAULT 0);",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: ALTER TABLE to modify primary key to BIGINT UNSIGNED (从xml中补充)",
		"ALTER TABLE customers MODIFY id BIGINT UNSIGNED NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id BIGINT NOT NULL PRIMARY KEY, name VARCHAR(32) DEFAULT '' NOT NULL, sex INT DEFAULT 0, age INT DEFAULT 0);"),
		nil, newTestResult())
}

// ==== Rule test code end ====
