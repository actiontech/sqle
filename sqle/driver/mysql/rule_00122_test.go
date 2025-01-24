package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00122(t *testing.T) {
	ruleName := ai.SQLE00122
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT SUM on a column where all values are NULL", "SELECT SUM(age) FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, age INT); INSERT INTO customers (id, age) VALUES (1, NULL), (2, NULL);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT SUM(age) FROM customers;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
			{
				Query: "SELECT (SELECT COUNT(*) FROM customers) - (SELECT COUNT(*) FROM customers WHERE age IS NOT NULL) RESULT",
				Rows:  sqlmock.NewRows([]string{"RESULT"}).AddRow(0),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT COUNT on a column where all values are NULL", "SELECT COUNT(age) FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, age INT); INSERT INTO customers (id, age) VALUES (1, NULL), (2, NULL);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT COUNT(age) FROM customers;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
			{
				Query: "SELECT (SELECT COUNT(*) FROM customers) - (SELECT COUNT(*) FROM customers WHERE age IS NOT NULL) RESULT",
				Rows:  sqlmock.NewRows([]string{"RESULT"}).AddRow(0),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT SUM on a column with mixed NULL and non-NULL values", "SELECT SUM(salary) FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, salary INT); INSERT INTO employees (id, salary) VALUES (1, 1000), (2, NULL), (3, 2000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT SUM(salary) FROM employees;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
			{
				Query: "SELECT (SELECT COUNT(*) FROM employees) - (SELECT COUNT(*) FROM employees WHERE salary IS NOT NULL) RESULT;",
				Rows:  sqlmock.NewRows([]string{"RESULT"}).AddRow(2),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 4: SELECT COUNT on a column with mixed NULL and non-NULL values", "SELECT COUNT(salary) FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, salary INT); INSERT INTO employees (id, salary) VALUES (1, 1000), (2, NULL), (3, 2000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT COUNT(salary) FROM employees;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
			{
				Query: "SELECT (SELECT COUNT(*) FROM employees) - (SELECT COUNT(*) FROM employees WHERE salary IS NOT NULL) RESULT;",
				Rows:  sqlmock.NewRows([]string{"RESULT"}).AddRow(2),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 5: SELECT COUNT(*) from the table", "SELECT COUNT(*) FROM orders;",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (id INT, order_date DATE); INSERT INTO orders (id, order_date) VALUES (1, '2023-01-01'), (2, '2023-01-02');"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: SELECT statement without SUM or COUNT", "SELECT name, address FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100), address VARCHAR(100)); INSERT INTO customers (id, name, address) VALUES (1, 'John Doe', '123 Main St'), (2, 'Jane Doe', '456 Elm St');"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: SELECT with multiple aggregates, one violating the rule", "SELECT SUM(age), COUNT(salary) FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, age INT, salary INT); INSERT INTO employees (id, age, salary) VALUES (1, NULL, 1000), (2, NULL, NULL), (3, NULL, 2000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT (SELECT COUNT(*) FROM employees) - (SELECT COUNT(*) FROM employees WHERE salary IS NOT NULL) RESULT;",
				Rows:  sqlmock.NewRows([]string{"RESULT"}).AddRow(100),
			},
			{
				Query: "SELECT (SELECT COUNT(*) FROM employees) - (SELECT COUNT(*) FROM employees WHERE age IS NOT NULL) RESULT;",
				Rows:  sqlmock.NewRows([]string{"RESULT"}).AddRow(0),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: SELECT SUM on an expression instead of a column", "SELECT SUM(1) FROM table1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT); INSERT INTO table1 (id) VALUES (1), (2), (3);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: SELECT SUM and COUNT on a column where all values are NULL (从xml中补充)", "SELECT SUM(age) as sum_age, COUNT(age) as count_age FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, age INT); INSERT INTO customers (id, age) VALUES (1, NULL), (2, NULL);"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT (SELECT COUNT(*) FROM customers) - (SELECT COUNT(*) FROM customers WHERE age IS NOT NULL) RESULT",
				Rows:  sqlmock.NewRows([]string{"RESULT"}).AddRow(0),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: SELECT SUM and COUNT with IFNULL on a column where all values are NULL (从xml中补充)", "SELECT SUM(IFNULL(age, 0)) as sum_age, COUNT(IFNULL(age, 0)) as count_age FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, age INT); INSERT INTO customers (id, age) VALUES (1, NULL), (2, NULL);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT SUM(IFNULL(age, 0)) as sum_age, COUNT(IFNULL(age, 0)) as count_age FROM customers;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
			{
				Query: "SELECT (SELECT COUNT(*) FROM customers) - (SELECT COUNT(*) FROM customers WHERE age IS NOT NULL) RESULT",
				Rows:  sqlmock.NewRows([]string{"RESULT"}).AddRow(0),
			},
		}, newTestResult())
}

// ==== Rule test code end ====
