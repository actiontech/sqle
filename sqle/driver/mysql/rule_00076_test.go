package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00076(t *testing.T) {
	ruleName := ai.SQLE00076
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: UPDATE语句影响行数超过阈值且没有WHERE子句",
		"UPDATE employees SET salary = salary * 1.1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT, salary DECIMAL(10,2), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE employees SET salary = salary * 1.1",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("UPDATE", 15000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: UPDATE语句影响行数超过阈值但有WHERE子句",
		"UPDATE employees SET salary = salary * 1.1 WHERE department_id = 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT, salary DECIMAL(10,2), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE employees SET salary = salary * 1.1 WHERE department_id = 10",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("UPDATE", 12000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: UPDATE语句影响行数不超过阈值且有WHERE子句",
		"UPDATE employees SET salary = salary * 1.1 WHERE employee_id = 5;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT, salary DECIMAL(10,2), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE employees SET salary = salary * 1.1 WHERE employee_id = 5",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("UPDATE", 1),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 4: DELETE语句影响行数超过阈值且没有WHERE子句",
		"DELETE FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT, salary DECIMAL(10,2), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM employees",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("DELETE", 15000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: DELETE语句影响行数超过阈值但有WHERE子句",
		"DELETE FROM employees WHERE department_id = 10;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT, salary DECIMAL(10,2), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM employees WHERE department_id = 10",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("DELETE", 12000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: DELETE语句影响行数不超过阈值且有WHERE子句",
		"DELETE FROM employees WHERE employee_id = 5;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT, salary DECIMAL(10,2), department_id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM employees WHERE employee_id = 5",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("DELETE", 1),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 7: UPDATE语句在没有WHERE子句时影响行数超过阈值(从xml中补充)",
		"UPDATE t2 SET name = concat(name,'1');",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE t2 SET name = concat(name,'1')",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("UPDATE", 15000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: UPDATE语句使用索引字段的WHERE子句影响行数超过阈值(从xml中补充)",
		"UPDATE t2 SET name = concat(name,'1') WHERE name LIKE '%t%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE t2 SET name = concat(name,'1') WHERE name LIKE '%t%'",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("UPDATE", 12000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: DELETE语句在没有WHERE子句时影响行数超过阈值(从xml中补充)",
		"DELETE FROM t2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM t2",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("DELETE", 15000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: DELETE语句使用子字符串匹配的WHERE子句影响行数超过阈值(从xml中补充)",
		"DELETE FROM t2 WHERE name LIKE '%t%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM t2 WHERE name LIKE '%t%'",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("DELETE", 12000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: UPDATE语句使用合适的条件影响行数不超过阈值(从xml中补充)",
		"UPDATE t2 SET name = concat(name,'1') WHERE (id < 101) AND name LIKE 't%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE t2 SET name = concat(name,'1') WHERE (id < 101) AND name LIKE 't%'",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("UPDATE", 50),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 12: DELETE语句使用精确条件影响行数不超过阈值(从xml中补充)",
		"DELETE FROM t2 WHERE id <= 100;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM t2 WHERE id <= 100",
				Rows:  sqlmock.NewRows([]string{"operation", "rows"}).AddRow("DELETE", 100),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())
}

// ==== Rule test code end ====
