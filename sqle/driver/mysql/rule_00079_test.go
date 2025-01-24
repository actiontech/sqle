package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00079(t *testing.T) {
	ruleName := ai.SQLE00079
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT 语句中表别名与表名相同", "SELECT t1.column1 FROM table1 AS table1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.column1 FROM table1 AS table1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT 语句中列别名与列名相同", "SELECT column1 AS column1 FROM table1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT column1 AS column1 FROM table1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT 语句中表别名与列名相同", "SELECT t1.column1 FROM table1 AS column1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.column1 FROM table1 AS column1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: SELECT 语句中列别名与表名相同", "SELECT column1 AS table1 FROM table1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT column1 AS table1 FROM table1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: SELECT 语句中表别名与其他表名相同", "SELECT t1.column1 FROM table1 AS table2, table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT); CREATE TABLE table2 (column2 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.column1 FROM table1 AS table2, table2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: UNION 语句中子查询表别名与表名相同", "SELECT t1.column1 FROM table1 AS table1 UNION SELECT t2.column2 FROM table2 AS table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT); CREATE TABLE table2 (column2 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.column1 FROM table1 AS table1 UNION SELECT t2.column2 FROM table2 AS table2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: UNION 语句中子查询列别名与列名相同", "SELECT column1 AS column1 FROM table1 UNION SELECT column2 AS column2 FROM table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT); CREATE TABLE table2 (column2 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT column1 AS column1 FROM table1 UNION SELECT column2 AS column2 FROM table2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: UNION 语句中子查询表别名与列名相同", "SELECT t1.column1 FROM table1 AS column1 UNION SELECT t2.column2 FROM table2 AS column2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT); CREATE TABLE table2 (column2 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.column1 FROM table1 AS column1 UNION SELECT t2.column2 FROM table2 AS column2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: UNION 语句中子查询列别名与表名相同", "SELECT column1 AS table1 FROM table1 UNION SELECT column2 AS table2 FROM table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table1 (column1 INT); CREATE TABLE table2 (column2 INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT column1 AS table1 FROM table1 UNION SELECT column2 AS table2 FROM table2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: SELECT 语句中表别名与表名相同", "SELECT t1.id AS id FROM t1 AS t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.id AS id FROM t1 AS t1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: SELECT 语句中列别名与列名相同", "SELECT id AS id FROM t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id AS id FROM t1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: SELECT 语句中表别名与列名相同", "SELECT t1.id FROM t1 AS id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.id FROM t1 AS id;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: SELECT 语句中列别名与表名相同", "SELECT id AS t1 FROM t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id AS t1 FROM t1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SELECT 语句中表别名与其他表名相同", "SELECT t1.id FROM t1 AS t2, t2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.id FROM t1 AS t2, t2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: UNION 语句中子查询表别名与表名相同", "SELECT t1.id FROM t1 AS t1 UNION SELECT t2.id FROM t2 AS t2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.id FROM t1 AS t1 UNION SELECT t2.id FROM t2 AS t2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: UNION 语句中子查询列别名与列名相同", "SELECT id AS id FROM t1 UNION SELECT id AS id FROM t2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id AS id FROM t1 UNION SELECT id AS id FROM t2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: UNION 语句中子查询表别名与列名相同", "SELECT t1.id FROM t1 AS id UNION SELECT t2.id FROM t2 AS id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT t1.id FROM t1 AS id UNION SELECT t2.id FROM t2 AS id;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: UNION 语句中子查询列别名与表名相同", "SELECT id AS t1 FROM t1 UNION SELECT id AS t2 FROM t2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id AS t1 FROM t1 UNION SELECT id AS t2 FROM t2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: join select  ", "select t1.id id, t1.r1 r1 ,t2.id id,t2.r1 r1 from t1 inner join t2 on t1.id = t2.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, r1 int); CREATE TABLE t2 (id INT, r1 int);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: join select   ", "select t1.id t1_id, t1.r1 t1_r1 ,t2.id t2_id,t2.r1 t2_r1 from t1 inner join t2 on t1.id = t2.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, r1 int); CREATE TABLE t2 (id INT, r1 int);CREATE TABLE t3 (id INT, r1 int,id2 INT, r2 int);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 18: join select.   ", "select a.id t1_id, a.r1 t1_r1 ,b.id t2_id,b.r1 t2_r1 from t1 a inner join t2 as b on a.id = b.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, r1 int); CREATE TABLE t2 (id INT, r1 int);CREATE TABLE t3 (id INT, r1 int,id2 INT, r2 int);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 19: join select  ", "select t1.id id, t1.r1 r1 ,t2.id id,t2.r1 r1 from t1 t2 inner join t2 t1 on t1.id = t2.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, r1 int); CREATE TABLE t2 (id INT, r1 int);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 20: insert .... select ...   ", "insert into t3 (select t1.id id, t1.r1 r1 ,t2.id id,t2.r1 r1 from t1 t2 inner join t2 t1 on t1.id = t2.id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, r1 int); CREATE TABLE t2 (id INT, r1 int);CREATE TABLE t3 (id INT, r1 int,id2 INT, r2 int);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 21: update ....where id select ...   ", "UPDATE t3 SET id = 2 where id in (select t1.id id from t1 t2 inner join t2 t1 on t1.id = t2.id limit 1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, r1 int); CREATE TABLE t2 (id INT, r1 int);CREATE TABLE t3 (id INT, r1 int,id2 INT, r2 int);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 22: delete ....where id select ...   ", "DELETE FROM t3 WHERE id in (select t1.id id from t1 t2 inner join t2 t1 on t1.id = t2.id limit 1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, r1 int); CREATE TABLE t2 (id INT, r1 int);CREATE TABLE t3 (id INT, r1 int,id2 INT, r2 int);"),
		nil, newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
