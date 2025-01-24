package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00058(t *testing.T) {
	ruleName := ai.SQLE00058
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE 使用 PARTITION BY RANGE 定义分区",
		"CREATE TABLE partitioned_table (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20));",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 2: CREATE TABLE 不使用分区相关功能",
		"CREATE TABLE simple_table (id INT, name VARCHAR(50));",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 3: ALTER TABLE 添加分区",
		"ALTER TABLE partitioned_table ADD PARTITION (PARTITION p2 VALUES LESS THAN (30));",
		session.NewAIMockContext().WithSQL("CREATE TABLE partitioned_table (name VARCHAR(255));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 4: ALTER TABLE 不修改分区",
		"ALTER TABLE simple_table ADD COLUMN age INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE simple_table (name VARCHAR(255));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 7: CREATE TABLE 使用 PARTITION BY RANGE 定义分区 (从xml中补充)",
		"CREATE TABLE customers (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id)) PARTITION BY RANGE (id-1000) (PARTITION p0 VALUES LESS THAN (100000), PARTITION p1 VALUES LESS THAN (300000), PARTITION p2 VALUES LESS THAN (500000), PARTITION p3 VALUES LESS THAN (700000), PARTITION p4 VALUES LESS THAN (900000), PARTITION p5 VALUES LESS THAN (MAXVALUE));",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 9: CREATE TABLE 不使用分区相关功能 (从xml中补充)",
		"CREATE TABLE customers0 (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 10: CREATE TABLE 不使用分区相关功能 (从xml中补充)",
		"CREATE TABLE customers1 (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 11: CREATE TABLE 不使用分区相关功能 (从xml中补充)",
		"CREATE TABLE customers2 (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 12: CREATE TABLE 不使用分区相关功能 (从xml中补充)",
		"CREATE TABLE customers3 (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 13: CREATE TABLE 不使用分区相关功能 (从xml中补充)",
		"CREATE TABLE customers4 (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 14: CREATE TABLE 不使用分区相关功能 (从xml中补充)",
		"CREATE TABLE customers5 (id INT NOT NULL, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult(),
	)
}

// ==== Rule test code end ====
