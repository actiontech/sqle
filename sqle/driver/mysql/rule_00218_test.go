package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00218(t *testing.T) {
	ruleName := ai.SQLE00218
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT语句中的WHERE条件包含联合索引的最左字段user_id，应通过",
		"SELECT * FROM user_table WHERE user_id = 123;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 2: SELECT语句中的WHERE条件仅包含联合索引的非最左字段username，应违反",
		"SELECT * FROM user_table WHERE username = 'john_doe';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 3: SELECT语句无WHERE条件，但ORDER BY包含联合索引的最左字段user_id，应通过",
		"SELECT * FROM user_table ORDER BY user_id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 4: SELECT语句无WHERE条件，但ORDER BY包含联合索引的非最左字段username，应违反",
		"SELECT * FROM user_table ORDER BY username;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 5: INSERT语句中的SELECT包含WHERE条件且包含联合索引的最左字段user_id，应通过",
		"INSERT INTO user_archive (user_id, username) SELECT user_id, username FROM user_table WHERE user_id = 456;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username)); CREATE TABLE user_archive (user_id INT, username VARCHAR(50));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 6: INSERT语句中的SELECT仅包含联合索引的非最左字段username，应违反",
		"INSERT INTO user_archive (user_id, username) SELECT user_id, username FROM user_table WHERE username = 'jane_doe';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username)); CREATE TABLE user_archive (user_id INT, username VARCHAR(50));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 7: UNION语句的所有SELECT子句均包含WHERE条件且包含联合索引的最左字段user_id，应通过",
		"SELECT * FROM user_table WHERE user_id = 1 UNION SELECT * FROM user_table WHERE user_id = 2;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 8: UNION语句的一个SELECT子句包含WHERE条件的非最左字段username，应违反",
		"SELECT * FROM user_table WHERE user_id = 1 UNION SELECT * FROM user_table WHERE username = 'alice';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 11: SELECT语句的WHERE条件总为真且GROUP BY包含联合索引的最左字段user_id，应通过",
		"SELECT user_id, COUNT(*) FROM user_table WHERE 1=1 GROUP BY user_id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 12: SELECT语句的WHERE条件总为真且GROUP BY包含联合索引的非最左字段username，应违反",
		"SELECT username, COUNT(*) FROM user_table WHERE 1=1 GROUP BY username;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE user_table (user_id INT, username VARCHAR(50), INDEX idx_user (user_id, username));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 13: SELECT语句中的WHERE条件包含联合索引的最左字段name，应通过(从xml中补充)",
		"SELECT * FROM customers WHERE name = '小王1';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE customers (id INT, name VARCHAR(50), age INT, INDEX idx_customers (name, age));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 14: SELECT语句中的WHERE条件仅包含联合索引的非最左字段age，应违反(从xml中补充)",
		"SELECT * FROM customers WHERE age < 30;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE customers (id INT, name VARCHAR(50), age INT, INDEX idx_customers (name, age));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 15: SELECT语句无WHERE条件，但ORDER BY包含联合索引的最左字段name，应通过(从xml中补充)",
		"SELECT * FROM customers ORDER BY name;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE customers (id INT, name VARCHAR(50), age INT, INDEX idx_customers (name, age));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 16: SELECT语句无WHERE条件，但ORDER BY包含联合索引的非最左字段age，应违反(从xml中补充)",
		"SELECT * FROM customers ORDER BY age;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE customers (id INT, name VARCHAR(50), age INT, INDEX idx_customers (name, age));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 17: INSERT语句中的SELECT包含WHERE条件且包含联合索引的最左字段name，应通过(从xml中补充)",
		"INSERT INTO customers_archive (id, name, sex, city, age) SELECT id, name, sex, city, age FROM customers WHERE name = '小王2';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT, INDEX idx_customers (name, age)); CREATE TABLE customers_archive (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT);",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 18: INSERT语句中的SELECT仅包含联合索引的非最左字段age，应违反(从xml中补充)",
		"INSERT INTO customers_archive (id, name, sex, city, age) SELECT id, name, sex, city, age FROM customers WHERE age < 25;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT, INDEX idx_customers (name, age)); CREATE TABLE customers_archive (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT);",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 19: UNION语句的所有SELECT子句均包含WHERE条件且包含联合索引的最左字段name，应通过(从xml中补充)",
		"SELECT * FROM customers WHERE name = '小王1' UNION SELECT * FROM customers WHERE name = '小王2';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE customers (id INT, name VARCHAR(50), age INT, INDEX idx_customers (name, age));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 20: UNION语句的一个SELECT子句包含WHERE条件的非最左字段age，应违反(从xml中补充)",
		"SELECT * FROM customers WHERE name = '小王1' UNION SELECT * FROM customers WHERE age < 25;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE customers (id INT, name VARCHAR(50), age INT, INDEX idx_customers (name, age));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

}

// ==== Rule test code end ====
