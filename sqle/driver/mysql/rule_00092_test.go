package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQL00092(t *testing.T) {
	ruleName := ai.SQLE00092
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// DELETE语句未使用LIMIT子句
	runAIRuleCase(rule, t, "case 1: DELETE语句未使用LIMIT子句", "DELETE FROM users WHERE age > 30;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult().addResult(ruleName))

	// DELETE语句使用LIMIT子句
	runAIRuleCase(rule, t, "case 2: DELETE语句使用LIMIT子句", "DELETE FROM users WHERE age > 30 LIMIT 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult())

	// UPDATE语句未使用LIMIT子句
	runAIRuleCase(rule, t, "case 3: UPDATE语句未使用LIMIT子句", "UPDATE users SET age = age + 1 WHERE age > 30;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult().addResult(ruleName))

	// UPDATE语句使用LIMIT子句
	runAIRuleCase(rule, t, "case 4: UPDATE语句使用LIMIT子句", "UPDATE users SET age = age + 1 WHERE age > 30 LIMIT 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult())

	// DELETE语句未使用LIMIT子句 - customers表
	runAIRuleCase(rule, t, "case 5: DELETE语句未使用LIMIT子句 - customers表", "DELETE FROM customers WHERE age > 30;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult().addResult(ruleName))

	// DELETE语句使用LIMIT子句 - customers表
	runAIRuleCase(rule, t, "case 6: DELETE语句使用LIMIT子句 - customers表", "DELETE FROM customers WHERE age > 30 LIMIT 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult())

	// UPDATE语句未使用LIMIT子句 - customers表
	runAIRuleCase(rule, t, "case 7: UPDATE语句未使用LIMIT子句 - customers表", "UPDATE customers SET age = 100 WHERE age > 30;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult().addResult(ruleName))

	// UPDATE语句使用LIMIT子句 - customers表
	runAIRuleCase(rule, t, "case 8: UPDATE语句使用LIMIT子句 - customers表", "UPDATE customers SET age = 100 WHERE age > 30 LIMIT 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult())
}

// ==== Rule test code end ====
