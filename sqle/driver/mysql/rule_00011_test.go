package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00011(t *testing.T) {
	ruleName := ai.SQLE00011
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// case 1: 单个ALTER TABLE语句，针对customers表
	runAIRuleCase(rule, t, "case 1: 单个ALTER TABLE语句，针对customers表",
		"ALTER TABLE customers ADD COLUMN age INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil,
		newTestResult())

	// case 2: 单个ALTER TABLE语句，针对orders表
	runAIRuleCase(rule, t, "case 2: 单个ALTER TABLE语句，针对orders表",
		"ALTER TABLE orders MODIFY COLUMN status VARCHAR(20);",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT PRIMARY KEY, status VARCHAR(10));"),
		nil,
		newTestResult())

	runSingleRuleInspectCase(rule, t, "alter_table: alter table need merge", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 Add column v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
		ALTER TABLE exist_db.exist_tb_1 Add column v6 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
		ALTER TABLE exist_db.exist_tb_1 Add column v7 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
		ALTER TABLE exist_db.exist_tb_1 DROP COLUMN v7;`,
		newTestResult(),
		newTestResult().addResult(ruleName),
		newTestResult().addResult(ruleName),
		newTestResult().addResult(ruleName))

	// case 3: 多个ALTER TABLE语句，针对不同表customers和orders
	runAIRuleCase(rule, t, "case 3: 多个ALTER TABLE语句，针对不同表customers和orders",
		`ALTER TABLE customers ADD COLUMN email VARCHAR(255); ALTER TABLE orders DROP COLUMN order_date;`,
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50));").
			WithSQL("CREATE TABLE orders (order_id INT PRIMARY KEY, order_date DATE, status VARCHAR(10));"),
		nil,
		newTestResult(),
		newTestResult())

	// case 4: 多个ALTER TABLE语句，针对同一个表customers
	runAIRuleCase(rule, t, "case 4: 多个ALTER TABLE语句，针对同一个表customers",
		"ALTER TABLE customers ADD COLUMN phone VARCHAR(20); ALTER TABLE customers DROP COLUMN address;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), address VARCHAR(255));"),
		nil,
		newTestResult(),
		newTestResult().addResult(ruleName))

	// case 5: 多个ALTER TABLE语句，针对同一个表customers，包含不同操作
	runAIRuleCase(rule, t, "case 5: 多个ALTER TABLE语句，针对同一个表customers，包含不同操作",
		"ALTER TABLE customers ADD COLUMN email VARCHAR(255); ALTER TABLE customers ADD COLUMN phone VARCHAR(20);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil,
		newTestResult(),
		newTestResult().addResult(ruleName))

	// case 6: 多个ALTER TABLE语句，针对不同表customers和products
	runAIRuleCase(rule, t, "case 6: 多个ALTER TABLE语句，针对不同表customers和products",
		"ALTER TABLE customers ADD COLUMN loyalty_points INT; ALTER TABLE products ADD COLUMN stock INT;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50));").
			WithSQL("CREATE TABLE products (product_id INT PRIMARY KEY, name VARCHAR(50));"),
		nil,
		newTestResult(), newTestResult())

	// case 7: 多个ALTER TABLE语句，针对同一个表customers，包含添加和修改操作(从xml中补充)
	runAIRuleCase(rule, t, "case 7: 多个ALTER TABLE语句，针对同一个表customers，包含添加和修改操作(从xml中补充)",
		"ALTER TABLE customers ADD COLUMN type VARCHAR(4) NOT NULL DEFAULT '1'; ALTER TABLE customers ADD COLUMN sex SMALLINT(2) NOT NULL DEFAULT 0; ALTER TABLE customers CHANGE COLUMN name name VARCHAR(64) COMMENT '名称';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil,
		newTestResult(),
		newTestResult().addResult(ruleName),
		newTestResult().addResult(ruleName))

	// case 8: 多个ALTER TABLE语句，针对同一个表customers，包含添加和删除索引操作(从xml中补充)
	runAIRuleCase(rule, t, "case 8: 多个ALTER TABLE语句，针对同一个表customers，包含添加和删除索引操作(从xml中补充)",
		"ALTER TABLE customers ADD INDEX idx_customers_column(type); ALTER TABLE customers DROP INDEX idx_customers_column;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50),type VARCHAR(1));"),
		nil,
		newTestResult(),
		newTestResult().addResult(ruleName))

	// case 9: 单个ALTER TABLE语句，包含多个操作，针对customers表(从xml中补充)
	runAIRuleCase(rule, t, "case 9: 单个ALTER TABLE语句，包含多个操作，针对customers表(从xml中补充)",
		"ALTER TABLE customers ADD COLUMN type VARCHAR(4) NOT NULL DEFAULT '1', ADD COLUMN sex SMALLINT(2) NOT NULL DEFAULT 0, CHANGE COLUMN name name VARCHAR(64) COMMENT '名称', CHANGE COLUMN id id INT(11) COMMENT '编号', ADD INDEX idx_customers_column(type);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil,
		newTestResult())
}

// ==== Rule test code end ====
