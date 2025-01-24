package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00180(t *testing.T) {
	ruleName := ai.SQLE00180
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT without WHERE clause", "SELECT * FROM table1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE SELECT * FROM table1;",
				Rows:  sqlmock.NewRows([]string{"EXPLAIN"}).AddRow(`-> Table scan on table1  (cost=10.2 rows=100)`),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 2: SELECT with one WHERE clause", "SELECT * FROM table1 WHERE column1 = 'value';",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE SELECT * FROM table1 WHERE column1 = 'value';",
				Rows:  sqlmock.NewRows([]string{"EXPLAIN"}).AddRow(`-> Filter on column1  (cost=5.1 rows=50)`),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 3: SELECT with multiple WHERE clauses", "SELECT * FROM table1 WHERE column1 = 'value' AND column2 = 'value2' AND column3 = 'value3';",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE SELECT * FROM table1 WHERE column1 = 'value' AND column2 = 'value2' AND column3 = 'value3';",
				Rows: sqlmock.NewRows([]string{"EXPLAIN"}).
					AddRow(`-> Filter on column1  (cost=5.1 rows=50)
									-> Filter on column2  (cost=5.1 rows=50)
									-> Filter on column3  (cost=5.1 rows=50)`),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: SELECT with subquery having low NDV field", "SELECT count(*) FROM customers tb_outer WHERE name IN (SELECT name FROM customers tb_inner WHERE name='lily75');",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE SELECT count(*) FROM customers tb_outer WHERE name IN (SELECT name FROM customers tb_inner WHERE name='lily75');",
				Rows: sqlmock.NewRows([]string{"EXPLAIN"}).
					AddRow(`-> Subquery on tb_inner  (cost=8.5 rows=85)
									-> Filter on name  (cost=3.2 rows=32)
									-> Filter on name  (cost=3.2 rows=32)`),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 5: SELECT with subquery having high NDV field", "SELECT count(*) FROM customers tb_outer WHERE id IN (SELECT id FROM customers tb_inner WHERE name='lily75');",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE SELECT count(*) FROM customers tb_outer WHERE id IN (SELECT id FROM customers tb_inner WHERE name='lily75');",
				Rows:  sqlmock.NewRows([]string{"EXPLAIN"}).AddRow(`-> Subquery on tb_inner  (cost=8.5 rows=85)`),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 6: INSERT INTO...SELECT without WHERE", "INSERT INTO table2 SELECT * FROM table1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100)); CREATE TABLE table2 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE INSERT INTO table2 SELECT * FROM table1;",
				Rows:  sqlmock.NewRows([]string{"EXPLAIN"}).AddRow(`-> Table scan on table1  (cost=12.3 rows=123)`),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 7: INSERT INTO...SELECT with one WHERE clause", "INSERT INTO table2 SELECT * FROM table1 WHERE column1 = 'value';",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100)); CREATE TABLE table2 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE INSERT INTO table2 SELECT * FROM table1 WHERE column1 = 'value';",
				Rows:  sqlmock.NewRows([]string{"EXPLAIN"}).AddRow(`-> Filter on column1  (cost=6.4 rows=64)`),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 8: INSERT INTO...SELECT with multiple WHERE clauses", "INSERT INTO table2 SELECT * FROM table1 WHERE column1 = 'value' AND column2 = 'value2' AND column3 = 'value3';",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100)); CREATE TABLE table2 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE INSERT INTO table2 SELECT * FROM table1 WHERE column1 = 'value' AND column2 = 'value2' AND column3 = 'value3';",
				Rows: sqlmock.NewRows([]string{"EXPLAIN"}).
					AddRow(`-> Filter on column1  (cost=6.4 rows=64)
									-> Filter on column2  (cost=6.4 rows=64)
									-> Filter on column3  (cost=6.4 rows=64)`),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: UNION without WHERE clauses", "SELECT * FROM table1 UNION SELECT * FROM table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100)); CREATE TABLE table2 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE SELECT * FROM table1 UNION SELECT * FROM table2;",
				Rows: sqlmock.NewRows([]string{"EXPLAIN"}).
					AddRow(`-> Union on table1  (cost=15.0 rows=150)
									-> Union on table2  (cost=15.0 rows=150)`),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 10: UNION with one subquery having WHERE clause", "SELECT * FROM table1 WHERE column1 = 'value' UNION SELECT * FROM table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100)); CREATE TABLE table2 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE SELECT * FROM table1 WHERE column1 = 'value' UNION SELECT * FROM table2;",
				Rows: sqlmock.NewRows([]string{"EXPLAIN"}).
					AddRow(`-> Filter on table1  (cost=7.7 rows=77)
									-> Union on table2  (cost=15.0 rows=150)`),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 11: UNION with multiple subqueries having WHERE clauses", "SELECT * FROM table1 WHERE column1 = 'value' UNION SELECT * FROM table2 WHERE column2 = 'value2' UNION SELECT * FROM table3 WHERE column3 = 'value3';",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100)); CREATE TABLE table2 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100)); CREATE TABLE table3 (id INT, column1 VARCHAR(100), column2 VARCHAR(100), column3 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN FORMAT=TREE SELECT * FROM table1 WHERE column1 = 'value' UNION SELECT * FROM table2 WHERE column2 = 'value2' UNION SELECT * FROM table3 WHERE column3 = 'value3';",
				Rows: sqlmock.NewRows([]string{"EXPLAIN"}).
					AddRow(`-> Filter on table1  (cost=7.7 rows=77)
									-> Filter on table2  (cost=7.7 rows=77)
									-> Filter on table3  (cost=7.7 rows=77)`),
			},
		}, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
