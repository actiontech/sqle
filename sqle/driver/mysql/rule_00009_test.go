package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00009(t *testing.T) {
	ruleName := ai.SQLE00009
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: DELETE 语句的 WHERE 子句中对字段应用函数但没有函数索引",
		"DELETE FROM employees WHERE substr(name,2,8) = 'JOHN';",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME, INDEX idx_trim_name (name));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 1: DELETE 语句的 WHERE 子句中对字段应用函数但没有函数索引",
		"DELETE FROM employees WHERE name = 'JOHN';",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME, INDEX idx_trim_name (name));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 4: INSERT...SELECT 语句中 SELECT 子句的 WHERE 条件应用函数且没有函数索引",
		"INSERT INTO archived_employees SELECT * FROM employees WHERE LOWER(status) = 'active';",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
			WithSQL("CREATE TABLE archived_employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	// runAIRuleCase(rule, t, "case 1: DELETE 语句的 WHERE 子句中对字段应用函数且没有函数索引",
	// 	"DELETE FROM employees WHERE UPPER(name) = 'JOHN';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN DELETE FROM employees WHERE UPPER(name) = 'JOHN'",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 2: DELETE 语句的 WHERE 子句中对字段应用函数且存在函数索引",
	// 	"DELETE FROM employees WHERE DATE(created_at) = '2023-01-01';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME, INDEX idx_date_created_at (DATE(created_at)));"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 3: DELETE 语句的 WHERE 子句中没有对字段应用函数",
	// 	"DELETE FROM employees WHERE name = 'JOHN';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 4: INSERT...SELECT 语句中 SELECT 子句的 WHERE 条件应用函数且没有函数索引",
	// 	"INSERT INTO archived_employees SELECT * FROM employees WHERE LOWER(status) = 'active';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
	// 		WithSQL("CREATE TABLE archived_employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN INSERT INTO archived_employees SELECT * FROM employees WHERE LOWER(status) = 'active'",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 5: INSERT...SELECT 语句中 SELECT 子句的 WHERE 条件应用函数且存在函数索引",
	// 	"INSERT INTO archived_employees SELECT * FROM employees WHERE YEAR(created_at) = 2022;",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME, INDEX idx_year_created_at (YEAR(created_at)));").
	// 		WithSQL("CREATE TABLE archived_employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 6: INSERT...SELECT 语句中 SELECT 子句的 WHERE 条件未应用函数",
	// 	"INSERT INTO archived_employees SELECT * FROM employees WHERE status = 'active';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
	// 		WithSQL("CREATE TABLE archived_employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 7: SELECT 语句的 WHERE 子句中对字段应用函数且没有函数索引",
	// 	"SELECT * FROM employees WHERE LENGTH(name) > 5;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN SELECT * FROM employees WHERE LENGTH(name) > 5",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 8: SELECT 语句的 WHERE 子句中对字段应用函数且存在函数索引",
	// 	"SELECT * FROM employees WHERE TRIM(name) = 'John';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME, INDEX idx_trim_name (TRIM(name)));"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 9: SELECT 语句的 WHERE 子句中没有对字段应用函数",
	// 	"SELECT * FROM employees WHERE name = 'John';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 10: UPDATE 语句的 WHERE 子句中对字段应用函数且没有函数索引",
	// 	"UPDATE employees SET status = 'inactive' WHERE LOWER(name) = 'john';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN UPDATE employees SET status = 'inactive' WHERE LOWER(name) = 'john'",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 11: UPDATE 语句的 WHERE 子句中对字段应用函数且存在函数索引",
	// 	"UPDATE employees SET status = 'inactive' WHERE YEAR(created_at) = 2021;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME, INDEX idx_year_created_at (YEAR(created_at)));"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 12: UPDATE 语句的 WHERE 子句中没有对字段应用函数",
	// 	"UPDATE employees SET status = 'inactive' WHERE name = 'John';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 13: UNION 语句中一个 SELECT 子句的 WHERE 条件应用函数且没有函数索引",
	// 	"SELECT name FROM employees WHERE UPPER(status) = 'ACTIVE' UNION SELECT name FROM contractors WHERE name = 'Jane';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
	// 		WithSQL("CREATE TABLE contractors (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN SELECT name FROM employees WHERE UPPER(status) = 'ACTIVE'",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 14: UNION 语句中所有 SELECT 子句的 WHERE 条件未应用函数",
	// 	"SELECT name FROM employees WHERE status = 'active' UNION SELECT name FROM contractors WHERE name = 'Jane';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
	// 		WithSQL("CREATE TABLE contractors (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 15: UNION 语句中所有 SELECT 子句的 WHERE 条件应用函数且存在函数索引",
	// 	"SELECT name FROM employees WHERE DATE(created_at) = '2023-01-01' UNION SELECT name FROM contractors WHERE LOWER(name) = 'jane';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME, INDEX idx_date_created_at (DATE(created_at)));").
	// 		WithSQL("CREATE TABLE contractors (id INT, name VARCHAR(50), created_at DATETIME, INDEX idx_lower_name (LOWER(name)));"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 16: WITH 语句中一个 CTE 的 WHERE 条件应用函数且没有函数索引",
	// 	"WITH active_employees AS (SELECT * FROM employees WHERE UPPER(status) = 'ACTIVE') SELECT * FROM active_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN WITH active_employees AS (SELECT * FROM employees WHERE UPPER(status) = 'ACTIVE') SELECT * FROM active_employees",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 17: WITH 语句中所有 CTE 的 WHERE 条件未应用函数",
	// 	"WITH active_employees AS (SELECT * FROM employees WHERE status = 'active') SELECT * FROM active_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 18: WITH 语句中所有 CTE 的 WHERE 条件应用函数且存在函数索引",
	// 	"WITH recent_employees AS (SELECT * FROM employees WHERE YEAR(created_at) = 2022) SELECT * FROM recent_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME, INDEX idx_year_created_at (YEAR(created_at)));"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 19: SELECT 语句中 WHERE 条件对字段应用函数且没有函数索引(从xml中补充)",
	// 	"SELECT count(*) FROM customers WHERE SUBSTR(log_date, 2, 6) = '02';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, log_date VARCHAR(20));"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN SELECT count(*) FROM customers WHERE SUBSTR(log_date, 2, 6) = '02'",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 20: UPDATE 语句中 WHERE 条件对字段应用函数且没有函数索引(从xml中补充)",
	// 	"UPDATE customers SET mark = '10' WHERE SUBSTR(log_date, 2, 6) = '02';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, log_date VARCHAR(20), mark VARCHAR(10));"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN UPDATE customers SET mark = '10' WHERE SUBSTR(log_date, 2, 6) = '02'",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 21: DELETE 语句中 WHERE 条件对字段应用函数且没有函数索引(从xml中补充)",
	// 	"DELETE FROM customers WHERE SUBSTR(log_date, 2, 6) = '02';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, log_date VARCHAR(20));"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN DELETE FROM customers WHERE SUBSTR(log_date, 2, 6) = '02'",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )
}

// ==== Rule test code end ====
