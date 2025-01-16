package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00083(t *testing.T) {
	ruleName := ai.SQLE00083
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT语句不包含GROUP BY或DISTINCT，且涉及一张表，违反索引跳跃扫描规则",
		"SELECT col1 FROM test_table WHERE col2 = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col1 INT, col2 INT, col3 INT, PRIMARY KEY (col1, col2));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT col1 FROM test_table WHERE col2 = 1",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using index for skip scan"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	// 语法不支持
	// runAIRuleCase(rule, t, "case 5: WITH语句不包含GROUP BY或DISTINCT，且涉及一张表，违反索引跳跃扫描规则",
	// 	"WITH cte AS (SELECT col1 FROM test_table WHERE col2 = 1) SELECT * FROM cte;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col1 INT, col2 INT, col3 INT, PRIMARY KEY (col1, col2));"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "EXPLAIN WITH cte AS (SELECT col1 FROM test_table WHERE col2 = 1) SELECT * FROM cte",
	// 			Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using index for skip scan"),
	// 		},
	// 		{
	// 			Query: "SHOW WARNINGS",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: UNION语句中的SELECT子句不包含GROUP BY或DISTINCT，且涉及一张表，违反索引跳跃扫描规则",
		"SELECT col1 FROM test_table WHERE col2 = 1 UNION SELECT col1 FROM test_table WHERE col3 = 2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col1 INT, col2 INT, col3 INT, PRIMARY KEY (col1, col2));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT col1 FROM test_table WHERE col2 = 1 UNION SELECT col1 FROM test_table WHERE col3 = 2",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using index for skip scan"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: SELECT语句在联合索引上进行跳跃扫描，违反索引跳跃扫描规则(从xml中补充)",
		"SELECT sex, age FROM customers WHERE age < 22;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, sex VARCHAR(10), age INT, PRIMARY KEY (id, sex));"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT sex, age FROM customers WHERE age < 22",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using index for skip scan"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: INSERT SELECT语句在联合索引上进行跳跃扫描，违反索引跳跃扫描规则(从xml中补充)",
		"INSERT INTO customers_sub SELECT sex, age FROM customers WHERE age < 22;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, sex VARCHAR(10), age INT, PRIMARY KEY (id, sex)); CREATE TABLE customers_sub (sex VARCHAR(10), age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO customers_sub SELECT sex, age FROM customers WHERE age < 22",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using index for skip scan"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: INSERT SELECT语句在联合索引上进行跳跃扫描，不违反索引跳跃扫描规则(从xml中补充)",
		"INSERT INTO customers_sub SELECT sex, age FROM customers WHERE age < 22;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, sex VARCHAR(10), age INT, PRIMARY KEY (id, sex)); CREATE TABLE customers_sub (sex VARCHAR(10), age INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO customers_sub SELECT sex, age FROM customers WHERE age < 22",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using index"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult())
}

// ==== Rule test code end ====
