package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQL00094(t *testing.T) {
	ruleName := ai.SQLE00094
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: DELETE 语句中使用 GROUP_CONCAT 函数", "DELETE FROM test_table WHERE id IN (SELECT GROUP_CONCAT(id) FROM another_table);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data VARCHAR(255)); CREATE TABLE another_table (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM test_table WHERE id IN (SELECT GROUP_CONCAT(id) FROM another_table);",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 2: DELETE 语句中没有使用内置函数", "DELETE FROM test_table WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data VARCHAR(255));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: DELETE 语句中使用 CONCAT_WS 函数", "DELETE FROM test_table WHERE id = CONCAT_WS('-', '1', '2');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM test_table WHERE id = CONCAT_WS('-', '1', '2');",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 4: INSERT 语句中使用 JSON_ARRAY 函数", "INSERT INTO test_table (id, data) VALUES (1, JSON_ARRAY('a', 'b', 'c'));",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data JSON);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO test_table (id, data) VALUES (1, JSON_ARRAY('a', 'b', 'c'));",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 5: INSERT 语句中没有使用内置函数", "INSERT INTO test_table (id, data) VALUES (1, 'test_data');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data VARCHAR(255));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: INSERT 语句中使用 GROUP_CONCAT 函数", "INSERT INTO test_table (id, data) SELECT id, GROUP_CONCAT(data) FROM another_table GROUP BY id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data VARCHAR(255)); CREATE TABLE another_table (id INT, data VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO test_table (id, data) SELECT id, GROUP_CONCAT(data) FROM another_table GROUP BY id;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 7: SELECT 语句中使用 JSON_ARRAY 函数", "SELECT JSON_ARRAY(id, name) FROM test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT JSON_ARRAY(id, name) FROM test_table;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 8: SELECT 语句中没有使用内置函数", "SELECT id, name FROM test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(255));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: SELECT 语句中使用 GROUP_CONCAT 函数", "SELECT GROUP_CONCAT(name) FROM test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT GROUP_CONCAT(name) FROM test_table;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 10: UPDATE 语句中使用 CONCAT_WS 函数", "UPDATE test_table SET name = CONCAT_WS('-', first_name, last_name) WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, first_name VARCHAR(255), last_name VARCHAR(255), name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE test_table SET name = CONCAT_WS('-', first_name, last_name) WHERE id = 1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 11: UPDATE 语句中没有使用内置函数", "UPDATE test_table SET name = 'new_name' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(255));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 12: UPDATE 语句中使用 JSON_ARRAY 函数", "UPDATE test_table SET data = JSON_ARRAY('a', 'b') WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data JSON);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE test_table SET data = JSON_ARRAY('a', 'b') WHERE id = 1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 12_1: UPDATE 语句中使用 JSON_ARRAY 函数", "UPDATE test_table SET data = 'hahha' WHERE data = JSON_ARRAY('a', 'b');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data JSON);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE test_table SET data = JSON_ARRAY('a', 'b') WHERE id = 1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 12_2: UPDATE 语句中使用 JSON_ARRAY 函数", "UPDATE test_table SET data = 'hahha' WHERE data in (select JSON_ARRAY(id, data) from exist_db.exist_tb_1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, data JSON);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE test_table SET data = JSON_ARRAY('a', 'b') WHERE id = 1;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 13: SELECT 语句中使用 JSON_ARRAY 函数并按结果分组", "SELECT JSON_ARRAY(id, name, age) result FROM customers GROUP BY result;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255), age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT JSON_ARRAY(id, name, age) result FROM customers GROUP BY result;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 14: SELECT 语句中使用 GROUP_CONCAT 函数并按 age 分组", "SELECT GROUP_CONCAT(name) AS name_list, age FROM customers GROUP BY age;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255), age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT GROUP_CONCAT(name) AS name_list, age FROM customers GROUP BY age;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 15: SELECT 语句中使用 CONCAT_WS 函数并限制行数", "SELECT CONCAT_ws(',', id, name, sex, city, age) AS row_list FROM customers LIMIT 2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255), sex VARCHAR(10), city VARCHAR(255), age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT CONCAT_WS(',', id, name, sex, city, age) AS row_list FROM customers LIMIT 2;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 16: SELECT 语句中没有使用内置函数", "SELECT id, name, age FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 17: SELECT 语句中没有使用内置函数并按 age 分组", "SELECT name, age FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 18: SELECT 语句中没有使用内置函数并限制行数", "SELECT id, name, sex, city, age FROM customers LIMIT 2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255), sex VARCHAR(10), city VARCHAR(255), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 19: SELECT 语句中没有使用内置函数并通过索引过滤", "SELECT COUNT(*) FROM customers WHERE name = '小王1';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 20: SELECT 语句中使用 FIND_IN_SET 函数进行过滤", "SELECT COUNT(*) FROM customers WHERE FIND_IN_SET('小王1', name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT COUNT(*) FROM customers WHERE FIND_IN_SET('小王1', name);",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 21: SELECT 语句中使用 JSON_ARRAY 函数进行group by", "SELECT COUNT(*) FROM customers GROUP BY JSON_ARRAY(id,name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT COUNT(*) FROM customers GROUP BY JSON_ARRAY(id,name);",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 22: SELECT 语句中使用 CONCAT_WS 函数进行order by", "SELECT id, name FROM customers ORDER BY CONCAT_WS(',', id,name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id, name FROM customers ORDER BY CONCAT_WS(',', id,name);",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 24: UNION select 语句中使用 CONCAT_WS 函数进行order by", "SELECT id, name FROM customers ORDER BY CONCAT_WS(',', id,name) union all SELECT id, name FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id, name FROM customers ORDER BY CONCAT_WS(',', id,name);",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))

	runAIRuleCase(rule, t, "case 25: select ...join 语句中使用 CONCAT_WS 函数进行order by", "SELECT CONCAT_WS(',', id,name) FROM customers as a join customers as b on a.id=b.id",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT CONCAT_WS(',', id,name) FROM customers as a join customers as b on a.id=b.id;",
				Rows:  sqlmock.NewRows([]string{"type"}).AddRow("index"),
			},
		}, newTestResult().addResult(ruleName, "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"))
}

// ==== Rule test code end ====
