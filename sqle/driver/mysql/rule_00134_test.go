package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00134(t *testing.T) {
	ruleName := ai.SQLE00134
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: UPDATE 语句修改单一主键字段，违反规则",
		"UPDATE users SET id = 10 WHERE name = 'Alice';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50), email VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: UPDATE 语句修改非主键字段，通过规则",
		"UPDATE users SET email = 'alice@example.com' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50), email VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: UPDATE 语句修改复合主键中的一个字段，违反规则",
		"UPDATE orders SET order_id = 100 WHERE product_id = 'XYZ';",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT, product_id VARCHAR(50), quantity INT, PRIMARY KEY (order_id, product_id));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: UPDATE 语句修改复合主键外的字段，通过规则",
		"UPDATE orders SET quantity = 5 WHERE order_id = 100 AND product_id = 'XYZ';",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT, product_id VARCHAR(50), quantity INT, PRIMARY KEY (order_id, product_id));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: UPDATE 语句同时修改主键和非主键字段，违反规则",
		"UPDATE users SET id = 2, email = 'new@example.com' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50), email VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: UPDATE 语句使用函数修改主键字段，违反规则",
		"UPDATE users SET id = id + 1 WHERE name = 'Bob';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50), email VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: UPDATE 语句修改主键字段，违反规则(从xml中补充)",
		"UPDATE customers SET id = 100000000000 + id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, customer_no INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: UPDATE 语句修改非主键字段，通过规则(从xml中补充)",
		"UPDATE customers SET customer_no = 100000000000 + customer_no;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, customer_no INT, name VARCHAR(50));"),
		nil, newTestResult())
}

// ==== Rule test code end ====
