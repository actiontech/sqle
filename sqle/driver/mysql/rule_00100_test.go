package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00100(t *testing.T) {
	ruleName := ai.SQLE00100
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: 简单SELECT语句，返回行数小于1000", "SELECT * FROM my_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(50)); INSERT INTO my_table VALUES (1, 'Alice'), (2, 'Bob');"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM my_table",
				Rows: sqlmock.NewRows([]string{"operation", "object_name", "object_owner", "rows"}).
					AddRow("SELECT STATEMENT", "my_table", "SYSTEM", 2).
					AddRow("TABLE ACCESS", "my_table", "SYSTEM", 2),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 2: 简单SELECT语句，返回行数超过1000", "SELECT * FROM my_large_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_large_table (id INT, name VARCHAR(50)); INSERT INTO my_large_table SELECT id, CONCAT('Name', id) FROM (SELECT @row := @row + 1 AS id FROM (SELECT 0 UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3) t1, (SELECT 0 UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3) t2, (SELECT @row := 0) t3 LIMIT 1001) t;"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM my_large_table",
				Rows: sqlmock.NewRows([]string{"operation", "object_name", "object_owner", "rows"}).
					AddRow("SELECT STATEMENT", "my_large_table", "SYSTEM", 1001).
					AddRow("TABLE ACCESS", "my_large_table", "SYSTEM", 1001),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT语句带LIMIT，限制行数小于等于1000，但实际大于1000", "SELECT * FROM my_table LIMIT 1000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(50)); INSERT INTO my_table VALUES (1, 'Alice'), (2, 'Bob');"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM my_table LIMIT 1000",
				Rows: sqlmock.NewRows([]string{"operation", "object_name", "object_owner", "rows"}).
					AddRow("SELECT STATEMENT", "my_table", "SYSTEM", 1001).
					AddRow("TABLE ACCESS", "my_table", "SYSTEM", 2),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: SELECT语句带LIMIT，限制行数超过1000", "SELECT * FROM my_table LIMIT 1500;",
		session.NewAIMockContext().WithSQL("CREATE TABLE my_table (id INT, name VARCHAR(50)); INSERT INTO my_table SELECT id, CONCAT('Name', id) FROM (SELECT @row := @row + 1 AS id FROM (SELECT 0 UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3) t1, (SELECT 0 UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3) t2, (SELECT @row := 0) t3 LIMIT 1001) t;"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM my_table LIMIT 1500",
				Rows: sqlmock.NewRows([]string{"operation", "object_name", "object_owner", "rows"}).
					AddRow("SELECT STATEMENT", "my_table", "SYSTEM", 1001).
					AddRow("TABLE ACCESS", "my_table", "SYSTEM", 1001),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: SELECT语句返回行数小于1000，带过滤条件", "SELECT * FROM customers WHERE id < 900;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(50)); INSERT INTO customers SELECT id, CONCAT('Customer', id) FROM (SELECT @row := @row + 1 AS id FROM (SELECT 0 UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3) t1, (SELECT 0 UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3) t2, (SELECT @row := 0) t3 LIMIT 900) t;"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM customers WHERE id < 900",
				Rows: sqlmock.NewRows([]string{"operation", "object_name", "object_owner", "rows"}).
					AddRow("SELECT STATEMENT", "customers", "SYSTEM", 900).
					AddRow("TABLE ACCESS", "customers", "SYSTEM", 900),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 9: UNION语句中一个SELECT子句带LIMIT超过1000,但实际大于1000", "SELECT id FROM table1 UNION SELECT id FROM table2 LIMIT 500;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT); INSERT INTO table1 VALUES (1), (2); CREATE TABLE table2 (id INT); INSERT INTO table2 VALUES (3), (4);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id FROM table1 UNION SELECT id FROM table2 LIMIT 500",
				Rows: sqlmock.NewRows([]string{"operation", "object_name", "object_owner", "rows"}).
					AddRow("SELECT STATEMENT", "table1", "SYSTEM", 1500).
					AddRow("TABLE ACCESS", "table1", "SYSTEM", 2).
					AddRow("SELECT STATEMENT", "table2", "SYSTEM", 2).
					AddRow("TABLE ACCESS", "table2", "SYSTEM", 2),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: INSER...select子句带LIMIT没有超过1000,但实际大于1000", "INSERT INTO table2 SELECT id FROM table1 UNION SELECT id FROM table2 WHERE 1=1 LIMIT 500;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT); INSERT INTO table1 VALUES (1), (2); CREATE TABLE table2 (id INT); INSERT INTO table2 VALUES (3), (4);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id FROM table1 UNION SELECT id FROM table2 WHERE 1=1 LIMIT 500",
				Rows: sqlmock.NewRows([]string{"operation", "object_name", "object_owner", "rows"}).
					AddRow("SELECT STATEMENT", "table1", "SYSTEM", 1500).
					AddRow("TABLE ACCESS", "table1", "SYSTEM", 2).
					AddRow("SELECT STATEMENT", "table2", "SYSTEM", 2).
					AddRow("TABLE ACCESS", "table2", "SYSTEM", 2),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: Delete中 select子句带LIMIT没有超过1000,但实际大于1000", "DELETE from table1 where id in (SELECT id FROM table1 UNION SELECT id FROM table2 LIMIT 500);",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT); INSERT INTO table1 VALUES (1), (2); CREATE TABLE table2 (id INT); INSERT INTO table2 VALUES (3), (4);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id FROM table1 UNION SELECT id FROM table2 LIMIT 500",
				Rows: sqlmock.NewRows([]string{"operation", "object_name", "object_owner", "rows"}).
					AddRow("SELECT STATEMENT", "table1", "SYSTEM", 1500).
					AddRow("TABLE ACCESS", "table1", "SYSTEM", 2).
					AddRow("SELECT STATEMENT", "table2", "SYSTEM", 2).
					AddRow("TABLE ACCESS", "table2", "SYSTEM", 2),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
