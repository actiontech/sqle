package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQL00078(t *testing.T) {
	ruleName := ai.SQLE00078
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT语句包含聚合函数COUNT", "SELECT COUNT(*) FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT COUNT(*) FROM employees;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT语句包含聚合函数SUM", "SELECT SUM(salary) FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary DECIMAL(10,2));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT SUM(salary) FROM employees;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT语句包含聚合函数AVG", "SELECT AVG(age) FROM users;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(50), age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT AVG(age) FROM users;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: SELECT语句包含聚合函数MAX", "SELECT MAX(score) FROM results;",
		session.NewAIMockContext().WithSQL("CREATE TABLE results (id INT, name VARCHAR(50), score INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT MAX(score) FROM results;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: SELECT语句包含聚合函数MIN", "SELECT MIN(price) FROM products;",
		session.NewAIMockContext().WithSQL("CREATE TABLE products (id INT, name VARCHAR(50), price DECIMAL(10,2));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT MIN(price) FROM products;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: UNION语句中的一个SELECT包含聚合函数", "SELECT name FROM employees UNION SELECT COUNT(*) FROM departments;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50)); CREATE TABLE departments (id INT, name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT name FROM employees UNION SELECT COUNT(*) FROM departments;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: UNION语句中的所有SELECT都不包含聚合函数", "SELECT name FROM employees UNION SELECT department FROM departments;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50)); CREATE TABLE departments (id INT, department VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT name FROM employees UNION SELECT department FROM departments;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 8: 简单的SELECT语句不包含聚合函数", "SELECT name, age FROM users;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(50), age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT name, age FROM users;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 9: SELECT语句包含子查询，子查询中有聚合函数", "SELECT name FROM employees WHERE id IN (SELECT MAX(id) FROM employees);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT name FROM employees WHERE id IN (SELECT MAX(id) FROM employees);",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: SELECT语句包含子查询，子查询中没有聚合函数", "SELECT name FROM employees WHERE id IN (SELECT id FROM departments);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50)); CREATE TABLE departments (id INT, name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT name FROM employees WHERE id IN (SELECT id FROM departments);",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 11: SELECT语句包含聚合函数COUNT并配合GROUP BY", "SELECT name, COUNT(1) as cn FROM t2 GROUP BY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(50), score INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT name, COUNT(1) as cn FROM t2 GROUP BY name;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: SELECT语句包含聚合函数AVG并配合GROUP BY", "SELECT name, AVG(score) as avgscore FROM t2 GROUP BY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(50), score INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT name, AVG(score) as avgscore FROM t2 GROUP BY name;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: SELECT语句包含聚合函数AVG并配合GROUP BY和HAVING", "SELECT name, AVG(score) as avgscore FROM t2 GROUP BY name HAVING AVG(score) > 0;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(50), score INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT name, AVG(score) as avgscore FROM t2 GROUP BY name HAVING AVG(score) > 0;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SELECT语句包含聚合函数SUM并配合WHERE", "SELECT SUM(score) FROM t2 WHERE score IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(50), score INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT SUM(score) FROM t2 WHERE score IS NOT NULL;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: SELECT语句包含聚合函数AVG并配合WHERE", "SELECT AVG(score) as avgs FROM t2 WHERE score IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(50), score INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT AVG(score) as avgs FROM t2 WHERE score IS NOT NULL;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
