package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00091(t *testing.T) {
	ruleName := ai.SQLE00091
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 0: SELECT语句中显式JOIN 包含unsing连接条件",
		"SELECT * FROM table1 JOIN table2 using(id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 1: SELECT语句中显式JOIN缺少连接条件",
		"SELECT * FROM table1 JOIN table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT语句中显式JOIN包含连接条件",
		"SELECT * FROM table1 JOIN table2 ON table1.id = table2.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: SELECT语句中隐式JOIN缺少连接条件",
		"SELECT * FROM table1, table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: SELECT语句中隐式JOIN包含连接条件",
		"SELECT * FROM table1, table2 WHERE table1.id <> table2.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: UNION语句中SELECT子句缺少连接条件",
		"SELECT * FROM table1 JOIN table2 UNION SELECT * FROM table3;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100)); CREATE TABLE table3 (id INT, value VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 6: WITH语句中SELECT子句缺少连接条件",
	// 	"WITH cte AS (SELECT * FROM table1 JOIN table2) SELECT * FROM cte;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
	// 	nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: INSERT...SELECT语句中JOIN缺少连接条件",
		"INSERT INTO table3 SELECT * FROM table1 JOIN table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100)); CREATE TABLE table3 (id INT, name VARCHAR(50), description VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: INSERT...SELECT语句中JOIN包含连接条件",
		"INSERT INTO table3 SELECT * FROM table1 JOIN table2 ON table1.id = table2.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100)); CREATE TABLE table3 (id INT, name VARCHAR(50), description VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: UPDATE语句中JOIN缺少连接条件",
		"UPDATE table1 JOIN table2 SET table1.name = 'test';",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: UPDATE语句中JOIN包含连接条件",
		"UPDATE table1 JOIN table2 ON table1.id != table2.id SET table1.name = 'test';",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 11: DELETE语句中JOIN缺少连接条件",
		"DELETE FROM table1 USING table1 JOIN table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: DELETE语句中JOIN包含连接条件",
		"DELETE FROM table1 USING table1 JOIN table2 ON table1.id = table2.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50)); CREATE TABLE table2 (id INT, description VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 13: SELECT语句中隐式JOIN缺少连接条件（使用示例表）",
		"SELECT a.id, a.name, b.id as ord_id, b.amount FROM customers a, orders b;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT); CREATE TABLE orders (id INT, c_id INT, amount DECIMAL(10,2));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SELECT语句中隐式JOIN包含连接条件（使用示例表）",
		"SELECT a.id, a.name, b.id as ord_id, b.amount FROM customers a, orders b WHERE a.id = b.c_id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT); CREATE TABLE orders (id INT, c_id INT, amount DECIMAL(10,2));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 15: INSERT...SELECT语句中JOIN缺少连接条件（使用示例表）",
		"INSERT INTO customers(name, sex, city, age) SELECT a.name, a.sex, a.city, a.age FROM customers a, orders b;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT); CREATE TABLE orders (id INT, c_id INT, amount DECIMAL(10,2));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: INSERT...SELECT语句中JOIN包含连接条件（使用示例表）",
		"INSERT INTO customers(name, sex, city, age) SELECT a.name, a.sex, a.city, a.age FROM customers a, orders b WHERE a.id = b.c_id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT); CREATE TABLE orders (id INT, c_id INT, amount DECIMAL(10,2));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 17: UPDATE语句中JOIN缺少连接条件（使用示例表）",
		"UPDATE customers SET age = age + 10 WHERE id IN (SELECT a.id FROM customers a, orders b);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT); CREATE TABLE orders (id INT, c_id INT, amount DECIMAL(10,2));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: UPDATE语句中JOIN包含连接条件（使用示例表）",
		"UPDATE customers SET age = age + 10 WHERE id IN (SELECT a.id FROM customers a, orders b WHERE a.id = b.c_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT); CREATE TABLE orders (id INT, c_id INT, amount DECIMAL(10,2));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 19: DELETE语句中JOIN缺少连接条件（使用示例表）",
		"DELETE FROM customers WHERE id IN (SELECT a.id FROM customers a, orders b);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT); CREATE TABLE orders (id INT, c_id INT, amount DECIMAL(10,2));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 20: DELETE语句中JOIN包含连接条件（使用示例表）",
		"DELETE FROM customers WHERE id IN (SELECT a.id FROM customers a, orders b WHERE a.id = b.c_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), sex VARCHAR(10), city VARCHAR(50), age INT); CREATE TABLE orders (id INT, c_id INT, amount DECIMAL(10,2));"),
		nil, newTestResult())
}

// ==== Rule test code end ====
