package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00219(t *testing.T) {
	ruleName := ai.SQLE00219
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE 缺少创建时间字段",
		"CREATE TABLE test_table (id INT);",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: CREATE TABLE 创建时间字段类型错误",
		"CREATE TABLE test_table (id INT, create_time DATETIME);",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: CREATE TABLE 创建时间字段默认值错误",
		"CREATE TABLE test_table (id INT, create_time TIMESTAMP DEFAULT '2023-01-01 00:00:00');",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: CREATE TABLE 正确的创建时间字段",
		"CREATE TABLE test_table (id INT, create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP);",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 5: ALTER TABLE 新增创建时间字段但类型错误",
		"ALTER TABLE test_table ADD COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: ALTER TABLE 新增创建时间字段但默认值错误",
		"ALTER TABLE test_table ADD COLUMN create_time TIMESTAMP DEFAULT '2023-01-01 00:00:00';",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: ALTER TABLE 正确新增创建时间字段",
		"ALTER TABLE test_table ADD COLUMN create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 8: ALTER TABLE 修改创建时间字段但类型错误",
		"ALTER TABLE test_table MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: ALTER TABLE 修改创建时间字段但默认值错误",
		"ALTER TABLE test_table MODIFY COLUMN create_time TIMESTAMP DEFAULT '2023-01-01 00:00:00';",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: ALTER TABLE 正确修改创建时间字段",
		"ALTER TABLE test_table MODIFY COLUMN create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, create_time TIMESTAMP DEFAULT '2023-01-01 00:00:00');"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 11: CREATE TABLE 缺少创建时间字段 (从xml中补充)",
		"CREATE TABLE customers (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: CREATE TABLE 正确的创建时间字段 (从xml中补充)",
		"CREATE TABLE customers (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 13: ALTER TABLE 正确新增创建时间字段 (从xml中补充)",
		"ALTER TABLE customers ADD create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 14: ALTER TABLE 正确修改创建时间字段 (从xml中补充)",
		"ALTER TABLE customers MODIFY create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, create_time TIMESTAMP DEFAULT '2023-01-01 00:00:00', PRIMARY KEY (id));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 15: ALTER TABLE ... CHANGE正确修改创建时间字段 (从xml中补充)",
		"ALTER TABLE customers CHANGE create_time create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, create_time TIMESTAMP DEFAULT '2023-01-01 00:00:00', PRIMARY KEY (id));"),
		nil,
		newTestResult())
}

// ==== Rule test code end ====
