package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00090(t *testing.T) {
	ruleName := ai.SQLE00090
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT 语句使用 UNION 而非 UNION ALL",
		"SELECT * FROM table1 UNION SELECT * FROM table2;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50));").
			WithSQL("CREATE TABLE table2 (id INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT 语句使用 UNION ALL",
		"SELECT * FROM table1 UNION ALL SELECT * FROM table2;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50));").
			WithSQL("CREATE TABLE table2 (id INT, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: INSERT ... SELECT 语句使用 UNION 而非 UNION ALL",
		"INSERT INTO table3 (col1) SELECT col1 FROM table1 UNION SELECT col1 FROM table2;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE table1 (col1 INT);").
			WithSQL("CREATE TABLE table2 (col1 INT);").
			WithSQL("CREATE TABLE table3 (col1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: INSERT ... SELECT 语句使用 UNION ALL",
		"INSERT INTO table3 (col1) SELECT col1 FROM table1 UNION ALL SELECT col1 FROM table2;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE table1 (col1 INT);").
			WithSQL("CREATE TABLE table2 (col1 INT);").
			WithSQL("CREATE TABLE table3 (col1 INT);"),
		nil, newTestResult())

	// runAIRuleCase(rule, t, "case 5: WITH 语句使用 UNION 而非 UNION ALL",
	// 	"WITH cte AS (SELECT * FROM table1 UNION SELECT * FROM table2) SELECT * FROM cte;",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50));").
	// 		WithSQL("CREATE TABLE table2 (id INT, name VARCHAR(50));"),
	// 	nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 6: WITH 语句使用 UNION ALL",
	// 	"WITH cte AS (SELECT * FROM table1 UNION ALL SELECT * FROM table2) SELECT * FROM cte;",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE table1 (id INT, name VARCHAR(50));").
	// 		WithSQL("CREATE TABLE table2 (id INT, name VARCHAR(50));"),
	// 	nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: SELECT 语句使用 UNION 而非 UNION ALL (从xml中补充)",
		"SELECT name, city FROM customers UNION SELECT name, city FROM suppliers;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (name VARCHAR(50), city VARCHAR(50));").
			WithSQL("CREATE TABLE suppliers (name VARCHAR(50), city VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: SELECT 语句使用 UNION ALL (从xml中补充)",
		"SELECT name, city FROM customers UNION ALL SELECT name, city FROM suppliers;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (name VARCHAR(50), city VARCHAR(50));").
			WithSQL("CREATE TABLE suppliers (name VARCHAR(50), city VARCHAR(50));"),
		nil, newTestResult())
}

// ==== Rule test code end ====
