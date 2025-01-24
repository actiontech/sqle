package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00053(t *testing.T) {
	ruleName := ai.SQLE00053
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT * 使用在简单查询中", "SELECT * FROM users;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT * 使用在带有WHERE条件的查询中", "SELECT * FROM orders WHERE order_id = 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT, product_name VARCHAR(100), quantity INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT * 使用在带有JOIN的查询中", "SELECT orders.* FROM customers JOIN orders ON customers.id = orders.customer_id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100)); CREATE TABLE orders (order_id INT, customer_id INT, product_name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: 指定列名的SELECT查询", "SELECT id, name FROM users;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: SELECT * 使用在子查询中", "SELECT name FROM (SELECT * FROM users) AS subquery;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: 指定列名的子查询", "SELECT name FROM (SELECT id, name FROM users) AS subquery;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: SELECT * 使用在简单查询中(从xml中补充)", "SELECT * FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: SELECT * 使用在带有WHERE条件的查询中(从xml中补充)", "SELECT * FROM customers WHERE age = 25;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: 指定列名的SELECT查询(从xml中补充)", "SELECT age FROM customers WHERE age = 25;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), age INT);"),
		nil, newTestResult())
}

// ==== Rule test code end ====
