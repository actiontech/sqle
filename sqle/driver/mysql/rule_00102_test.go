package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00102(t *testing.T) {
	ruleName := ai.SQLE00102
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// DELETE语句测试用例
	runAIRuleCase(rule, t, "case 1: DELETE语句包含ORDER BY子句",
		"DELETE FROM my_table WHERE id = 1 ORDER BY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: DELETE语句不包含ORDER BY子句",
		"DELETE FROM my_table WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: DELETE语句包含ORDER BY和LIMIT子句",
		"DELETE FROM my_table WHERE age > 30 ORDER BY age LIMIT 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: 带有不同大小写的DELETE语句包含ORDER BY子句",
		"delete from my_table where id = 2 OrDeR bY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: 带有多余空格的DELETE语句包含ORDER BY子句",
		"DELETE  FROM  my_table  WHERE  id = 3  ORDER BY  name ;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: DELETE语句包含子查询，但主语句不包含ORDER BY",
		"DELETE FROM my_table WHERE id IN (SELECT id FROM other_table WHERE active = 1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT); CREATE TABLE other_table (id INT, active INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: DELETE语句包含ORDER BY子句(从xml中补充)",
		"DELETE FROM customers ORDER BY age;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: DELETE语句不包含ORDER BY子句(从xml中补充)",
		"DELETE FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: DELETE语句包含子查询且子查询包含ORDER BY子句(从xml中补充)",
		"DELETE FROM customers WHERE id IN (SELECT c_id FROM orders ORDER BY id DESC);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT); CREATE TABLE orders (c_id INT, id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: DELETE语句包含子查询但主语句不包含ORDER BY子句(从xml中补充)",
		"DELETE FROM customers WHERE id IN (SELECT c_id FROM orders);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT); CREATE TABLE orders (c_id INT, id INT);"),
		nil, newTestResult())

	// UPDATE语句测试用例
	runAIRuleCase(rule, t, "case 11: UPDATE语句包含ORDER BY子句",
		"UPDATE my_table SET name = 'John' WHERE id = 1 ORDER BY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: UPDATE语句不包含ORDER BY子句",
		"UPDATE my_table SET name = 'John' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 13: UPDATE语句包含ORDER BY和LIMIT子句",
		"UPDATE my_table SET name = 'Jane' WHERE age > 30 ORDER BY age LIMIT 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: 带有不同大小写的UPDATE语句包含ORDER BY子句",
		"update my_table set name = 'Alice' where id = 2 OrDeR bY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: 带有多余空格的UPDATE语句包含ORDER BY子句",
		"UPDATE  my_table  SET  name = 'Bob'  WHERE  id = 3  ORDER BY  name ;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: UPDATE语句包含子查询，但主语句不包含ORDER BY",
		"UPDATE my_table SET name = (SELECT name FROM other_table WHERE id = my_table.id) WHERE id = 4;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(100), age INT); CREATE TABLE other_table (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 17: UPDATE语句包含ORDER BY子句(从xml中补充)",
		"UPDATE customers SET city = '北京' ORDER BY age;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT, city VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: UPDATE语句不包含ORDER BY子句(从xml中补充)",
		"UPDATE customers SET city = '北京';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT, city VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 19: UPDATE语句包含子查询且子查询包含ORDER BY子句(从xml中补充)",
		"UPDATE customers SET city = '北京' WHERE id IN (SELECT c_id FROM orders ORDER BY id DESC);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT, city VARCHAR(100)); CREATE TABLE orders (c_id INT, id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 20: UPDATE语句包含子查询但主语句不包含ORDER BY子句(从xml中补充)",
		"UPDATE customers SET city = '北京' WHERE id IN (SELECT c_id FROM orders);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT, city VARCHAR(100)); CREATE TABLE orders (c_id INT, id INT);"),
		nil, newTestResult())
}

// ==== Rule test code end ====
