package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00082(t *testing.T) {
	ruleName := ai.SQLE00082
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// case 1: SELECT 语句使用 ORDER BY 可能触发 filesort
	runAIRuleCase(rule, t, "case 1: SELECT 语句使用 ORDER BY 可能触发 filesort",
		"SELECT * FROM employees ORDER BY last_name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, first_name VARCHAR(50), last_name VARCHAR(50), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM employees ORDER BY last_name;",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using filesort"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	// case 2: SELECT 语句使用 WHERE 子句不触发 filesort
	runAIRuleCase(rule, t, "case 2: SELECT 语句使用 WHERE 子句不触发 filesort",
		"SELECT * FROM employees WHERE department_id = 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, first_name VARCHAR(50), last_name VARCHAR(50), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM employees WHERE department_id = 10;",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow(""),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	// case 3: SELECT 语句使用 GROUP BY 可能触发 filesort
	runAIRuleCase(rule, t, "case 3: SELECT 语句使用 GROUP BY 可能触发 filesort",
		"SELECT department_id, COUNT(*) FROM employees GROUP BY department_id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, first_name VARCHAR(50), last_name VARCHAR(50), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT department_id, COUNT(*) FROM employees GROUP BY department_id;",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using temporary; Using filesort"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	// case 4: SELECT 语句使用 LIMIT 不触发 filesort
	runAIRuleCase(rule, t, "case 4: SELECT 语句使用 LIMIT 不触发 filesort",
		"SELECT * FROM employees LIMIT 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, first_name VARCHAR(50), last_name VARCHAR(50), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM employees LIMIT 10;",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow(""),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	// case 5: SELECT 语句使用窗口函数可能触发 filesort // TODO 当前解析器解析不了此语法
	// runAIRuleCase(rule, t, "case 5: SELECT 语句使用窗口函数可能触发 filesort",
	// 	"SELECT age FROM (SELECT age, ROW_NUMBER() OVER (PARTITION BY age ORDER BY age DESC) gn FROM customers) T WHERE gn=1 LIMIT 1;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, age INT);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN SELECT age FROM (SELECT age, ROW_NUMBER() OVER (PARTITION BY age ORDER BY age DESC) gn FROM customers) T WHERE gn=1 LIMIT 1;",
	// 			Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using filesort"),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	}, newTestResult().addResult(ruleName))

	// case 6: SELECT 语句使用 GROUP BY 和 ORDER BY 不一致可能触发 filesort
	runAIRuleCase(rule, t, "case 6: SELECT 语句使用 GROUP BY 和 ORDER BY 不一致可能触发 filesort",
		"SELECT COUNT(*) FROM (SELECT age FROM customers GROUP BY age ORDER BY name DESC) T;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, age INT, name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT COUNT(*) FROM (SELECT age FROM customers GROUP BY age ORDER BY name DESC) T;",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using filesort"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	// case 7: SELECT 语句求最小值不触发 filesort
	runAIRuleCase(rule, t, "case 7: SELECT 语句求最小值不触发 filesort",
		"SELECT MIN(age) FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT MIN(age) FROM customers;",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow(""),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	// case 8: SELECT 语句使用 GROUP BY 无排序不触发 filesort
	runAIRuleCase(rule, t, "case 8: SELECT 语句使用 GROUP BY 无排序不触发 filesort",
		"SELECT COUNT(*) FROM (SELECT age FROM customers GROUP BY age) T;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT COUNT(*) FROM (SELECT age FROM customers GROUP BY age) T;",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow(""),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())
}

// ==== Rule test code end ====
