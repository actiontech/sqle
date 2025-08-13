package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00109(t *testing.T) {
	ruleName := ai.SQLE00109
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: DELETE语句中的子查询使用LIMIT",
		"DELETE FROM orders WHERE id IN (SELECT id FROM orders_archive LIMIT 10);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE orders (id INT, amount DECIMAL(10,2));").
			WithSQL("CREATE TABLE orders_archive (id INT, amount DECIMAL(10,2));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: DELETE语句中的子查询不使用LIMIT",
		"DELETE FROM orders WHERE id IN (SELECT id FROM orders_archive);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE orders (id INT, amount DECIMAL(10,2));").
			WithSQL("CREATE TABLE orders_archive (id INT, amount DECIMAL(10,2));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: UPDATE语句中的子查询使用LIMIT",
		"UPDATE products SET price = price * 0.9 WHERE id IN (SELECT id FROM products_archive LIMIT 5);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE products (id INT, price DECIMAL(10,2));").
			WithSQL("CREATE TABLE products_archive (id INT, price DECIMAL(10,2));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: UPDATE语句中的子查询不使用LIMIT",
		"UPDATE products SET price = price * 0.9 WHERE id IN (SELECT id FROM products_archive);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE products (id INT, price DECIMAL(10,2));").
			WithSQL("CREATE TABLE products_archive (id INT, price DECIMAL(10,2));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: SELECT语句中的子查询使用LIMIT",
		"SELECT * FROM customers WHERE id IN (SELECT id FROM customers_archive LIMIT 3);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100));").
			WithSQL("CREATE TABLE customers_archive (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: SELECT语句中的子查询不使用LIMIT",
		"SELECT * FROM customers WHERE id IN (SELECT id FROM customers_archive);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100));").
			WithSQL("CREATE TABLE customers_archive (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: INSERT ... SELECT语句中的子查询使用LIMIT",
		"INSERT INTO new_orders (id, amount) SELECT id, amount FROM old_orders WHERE id IN (SELECT id FROM orders_archive LIMIT 2);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE new_orders (id INT, amount DECIMAL(10,2));").
			WithSQL("CREATE TABLE old_orders (id INT, amount DECIMAL(10,2));").
			WithSQL("CREATE TABLE orders_archive (id INT, amount DECIMAL(10,2));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: INSERT ... SELECT语句中的子查询不使用LIMIT",
		"INSERT INTO new_orders (id, amount) SELECT id, amount FROM old_orders WHERE id IN (SELECT id FROM orders_archive);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE new_orders (id INT, amount DECIMAL(10,2));").
			WithSQL("CREATE TABLE old_orders (id INT, amount DECIMAL(10,2));").
			WithSQL("CREATE TABLE orders_archive (id INT, amount DECIMAL(10,2));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: SELECT ... UNION ALL SELECT语句中的子查询使用LIMIT",
		"SELECT * FROM (SELECT id FROM products LIMIT 1) AS sub UNION ALL SELECT * FROM products;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE products (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: SELECT ... UNION ALL SELECT语句中的子查询不使用LIMIT",
		"SELECT * FROM (SELECT id FROM products) AS sub UNION ALL SELECT * FROM products;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE products (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 11: SELECT语句中的子查询使用LIMIT (反例)",
		"SELECT name, city, age FROM customers WHERE id IN (SELECT id FROM customers LIMIT 2);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), city VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: UPDATE语句中的子查询使用LIMIT (反例)",
		"UPDATE customers SET city='北京' WHERE id IN (SELECT id FROM customers LIMIT 2);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), city VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: DELETE语句中的子查询使用LIMIT (反例)",
		"DELETE FROM customers WHERE id IN (SELECT id FROM customers LIMIT 2);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), city VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SELECT语句使用JOIN代替LIMIT子查询 (正例)",
		"SELECT a.name, a.city, a.age FROM customers a JOIN (SELECT id FROM customers LIMIT 2) b USING(id);",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), city VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
