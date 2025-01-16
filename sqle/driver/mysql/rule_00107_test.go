package mysql

import (
	"strings"
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00107(t *testing.T) {
	ruleName := ai.SQLE00107
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT语句长度小于1024，符合规则",
		"SELECT id, name FROM users WHERE status = 'active';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, status VARCHAR(50), age INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 2: SELECT语句长度大于等于1024，违反规则",
		"SELECT "+strings.Repeat("column", 200)+" FROM users WHERE "+strings.Repeat("status = 'active' AND ", 50)+"1=1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, status VARCHAR(50), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: UPDATE语句长度小于1024，符合规则",
		"UPDATE orders SET status = 'shipped' WHERE order_id = 12345;",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT, status VARCHAR(100), quantity INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 4: UPDATE语句长度大于等于1024，违反规则",
		"UPDATE orders SET order_id = 7 WHERE  1=1 "+strings.Repeat("AND 1=1 ", 200)+";",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT, status VARCHAR(100), quantity INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: DELETE语句长度小于1024，符合规则",
		"DELETE FROM orders WHERE order_id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT, status VARCHAR(100), quantity INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 7: INSERT...SELECT语句长度小于1024，符合规则",
		"INSERT INTO archive_orders (order_id, customer_id) SELECT order_id, customer_id FROM orders WHERE status = 'completed';",
		session.NewAIMockContext().WithSQL("CREATE TABLE archive_orders (order_id INT, status VARCHAR(100), customer_id INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 9: UNION语句长度小于1024，符合规则",
		"SELECT id, name FROM customers UNION SELECT id, name FROM suppliers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100));CREATE TABLE suppliers (id INT, name VARCHAR(100));"),
		nil,
		newTestResult())

}

// ==== Rule test code end ====
