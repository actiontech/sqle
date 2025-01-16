package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00005(t *testing.T) {
	ruleName := ai.SQLE00005
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// Test case 1: CREATE TABLE 时创建包含6个字段的复合索引
	runAIRuleCase(rule, t, "case 1: CREATE TABLE 时创建包含6个字段的复合索引",
		"CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT, address VARCHAR(100), email VARCHAR(100), phone VARCHAR(20), INDEX idx_sample (id, name, age, address, email, phone));",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	// Test case 2: CREATE TABLE 时创建包含5个字段的复合索引
	runAIRuleCase(rule, t, "case 2: CREATE TABLE 时创建包含5个字段的复合索引",
		"CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT, address VARCHAR(100), email VARCHAR(100), INDEX idx_sample (id, name, age, address, email));",
		nil,
		nil,
		newTestResult(),
	)

	// Test case 3: CREATE TABLE 时创建包含3个字段的复合索引
	runAIRuleCase(rule, t, "case 3: CREATE TABLE 时创建包含3个字段的复合索引",
		"CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT, INDEX idx_sample (id, name, age));",
		nil,
		nil,
		newTestResult(),
	)

	// Test case 4: ALTER TABLE 时添加包含6个字段的复合索引
	runAIRuleCase(rule, t, "case 4: ALTER TABLE 时添加包含6个字段的复合索引",
		"ALTER TABLE sample_table ADD INDEX idx_new (id, name, age, address, email, phone);",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT, address VARCHAR(100), email VARCHAR(100), phone VARCHAR(20));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	// Test case 5: ALTER TABLE 时添加包含5个字段的复合索引
	runAIRuleCase(rule, t, "case 5: ALTER TABLE ... ADD FULLTEXT KEY 时添加包含5个字段的复合索引",
		"ALTER TABLE sample_table ADD FULLTEXT KEY idx_new (id, name, age, address, email);",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT, address VARCHAR(100), email VARCHAR(100));"),
		nil,
		newTestResult(),
	)

	// Test case 6: ALTER TABLE 时添加包含3个字段的复合索引
	runAIRuleCase(rule, t, "case 6: ALTER TABLE 时添加包含3个字段的复合索引",
		"ALTER TABLE sample_table ADD INDEX idx_new (id, name, age);",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT);"),
		nil,
		newTestResult(),
	)

	// Test case 7: CREATE INDEX 时创建包含6个字段的复合索引
	runAIRuleCase(rule, t, "case 7: CREATE INDEX 时创建包含6个字段的复合索引",
		"CREATE INDEX idx_sample ON sample_table (id, name, age, address, email, phone);",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT, address VARCHAR(100), email VARCHAR(100), phone VARCHAR(20));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	// Test case 8: CREATE INDEX 时创建包含5个字段的复合索引
	runAIRuleCase(rule, t, "case 8: CREATE INDEX 时创建包含5个字段的复合索引",
		"CREATE INDEX idx_sample ON sample_table (id, name, age, address, email);",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT, address VARCHAR(100), email VARCHAR(100));"),
		nil,
		newTestResult(),
	)

	// Test case 9: CREATE INDEX 时创建包含3个字段的复合索引
	runAIRuleCase(rule, t, "case 9: CREATE INDEX 时创建包含3个字段的复合索引",
		"CREATE INDEX idx_sample ON sample_table (id, name, age);",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(50), age INT);"),
		nil,
		newTestResult(),
	)

	// Test case 10: CREATE TABLE 时创建包含6个字段的复合索引(从xml中补充)
	runAIRuleCase(rule, t, "case 10: CREATE TABLE 时创建包含6个字段的复合索引(从xml中补充)",
		"CREATE TABLE customers (id INT, name VARCHAR(32), sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100), INDEX idx_customers (id, name, sex, age, mark1, mark2));",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	// Test case 11: CREATE TABLE 时创建包含5个字段的复合索引(从xml中补充)
	runAIRuleCase(rule, t, "case 11: CREATE TABLE 时创建包含5个字段的复合索引(从xml中补充)",
		"CREATE TABLE customers (id INT, name VARCHAR(32), sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100), INDEX idx_customers (id, name, sex, age, mark1));",
		nil,
		nil,
		newTestResult(),
	)

	// Test case 12: ALTER TABLE 时添加包含6个字段的复合索引(从xml中补充)
	runAIRuleCase(rule, t, "case 12: ALTER TABLE 时添加包含6个字段的复合索引(从xml中补充)",
		"ALTER TABLE customers ADD INDEX idx_customers (id, name, sex, age, mark1, mark2);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(32), sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	// Test case 13: ALTER TABLE 时添加包含5个字段的复合索引(从xml中补充)
	runAIRuleCase(rule, t, "case 13: ALTER TABLE 时添加包含5个字段的复合索引(从xml中补充)",
		"ALTER TABLE customers ADD INDEX idx_customers (id, name, sex, age, mark1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(32), sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100));"),
		nil,
		newTestResult(),
	)

	// Test case 14: CREATE INDEX 时创建包含6个字段的复合索引(从xml中补充)
	runAIRuleCase(rule, t, "case 14: CREATE INDEX 时创建包含6个字段的复合索引(从xml中补充)",
		"CREATE INDEX idx_customers ON customers (id, name, sex, age, mark1, mark2);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(32), sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	// Test case 15: CREATE UNIQUE INDEX 时创建包含5个字段的复合索引(从xml中补充)
	runAIRuleCase(rule, t, "case 15: CREATE INDEX 时创建包含5个字段的复合索引(从xml中补充)",
		"CREATE UNIQUE INDEX idx_customers ON customers (id, name, sex, age, mark1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(32), sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100));"),
		nil,
		newTestResult(),
	)
}

// ==== Rule test code end ====
