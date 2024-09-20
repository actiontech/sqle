package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQLE00030(t *testing.T) {
	ruleName := ai.SQLE00030
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	for _, sql := range []string{
		`CREATE TRIGGER my_trigger BEFORE INSERT ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.name = UPPER(NEW.name);`,
		`CREATE DEFINER='sqle_op'@'localhost' TRIGGER my_trigger BEFORE INSERT ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.name = UPPER(NEW.name);`,
		`CREATE TRIGGER my_trigger AFTER UPDATE ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.updated_at = NOW();`,
		`CREATE TRIGGER ins_check BEFORE INSERT ON customers FOR EACH ROW BEGIN IF NEW.age < 18 THEN SET NEW.mark1 = '未满18岁'; ELSEIF NEW.age >= 18 THEN SET NEW.mark1 = '满18岁，已经成年了'; END IF; END;`,
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))
	}

	for _, sql := range []string{
		`DROP TRIGGER IF EXISTS ins_check;`,
		`CREATEDEFINER='sqle_op'@'localhost' TRIGGER my_trigger BEFORE INSERT ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.name = UPPER(NEW.name);`,
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))
	}

	// // Test case 1
	// i := NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 1", i, `
	// CREATE TRIGGER my_trigger BEFORE INSERT ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.name = UPPER(NEW.name);
	// `, newTestResult().addResult(ruleName))

	// // Test case 2
	// i = NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 2", i, `
	// CREATE TABLE exist_db.exist_tb_1 (id INT, name VARCHAR(100));
	// `, newTestResult())

	// // Test case 3
	// i = NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 3", i, `
	// CREATE TABLE exist_db.exist_tb_1 (id INT, name VARCHAR(100)); CREATE TRIGGER my_trigger BEFORE INSERT ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.name = UPPER(NEW.name);
	// `, newTestResult().addResult(ruleName))

	// // Test case 4
	// i = NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 4", i, `
	// CREATE TABLE orders (order_id INT, customer_id INT, PRIMARY KEY (order_id), FOREIGN KEY (customer_id) REFERENCES customers(customer_id));
	// `, newTestResult())

	// // Test case 5
	// i = NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 5", i, `
	// CREATE TRIGGER my_trigger AFTER UPDATE ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.updated_at = NOW();
	// `, newTestResult().addResult(ruleName))

	// // Test case 6
	// i = NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 6", i, `
	// CREATE TRIGGER ins_check BEFORE INSERT ON customers FOR EACH ROW BEGIN IF NEW.age < 18 THEN SET NEW.mark1 = '未满18岁'; ELSEIF NEW.age >= 18 THEN SET NEW.mark1 = '满18岁，已经成年了'; END IF; END;
	// `, newTestResult().addResult(ruleName))

	// // Test case 7
	// i = NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 7", i, `
	// DROP TRIGGER IF EXISTS ins_check;
	// `, newTestResult())

	// // Test case 8
	// i = NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 8", i, `
	// CREATE TABLE customers (id INT NOT NULL PRIMARY KEY, name VARCHAR(32) DEFAULT '' NOT NULL, sex INT DEFAULT 0, age INT DEFAULT 0, mark1 VARCHAR(200));
	// `, newTestResult())

	// // Test case 9
	// i = NewMockInspect()
	// runSingleRuleInspectCase(rule, t, "Test case 9", i, `
	// INSERT INTO customers (id, name, sex, age) VALUES (1, '小季', 0, 20);
	// `, newTestResult())
}

// ==== Rule test code end ====
