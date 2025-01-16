package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00084(t *testing.T) {
	ruleName := ai.SQLE00084
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TEMPORARY TABLE 使用临时表",
		"CREATE TEMPORARY TABLE temp_table (id INT);",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: CREATE TABLE 不使用临时表",
		"CREATE TABLE regular_table (id INT);",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 3: INSERT INTO SELECT 使用临时表",
		"INSERT INTO target_table SELECT * FROM source_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE target_table (id INT); CREATE TABLE source_table (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO target_table SELECT * FROM source_table",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using temporary"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: SELECT 使用临时表",
		"SELECT * FROM source_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE source_table (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM source_table",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using temporary"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: UPDATE 使用临时表",
		"UPDATE target_table SET id = 1 WHERE id IN (SELECT id FROM source_table);",
		session.NewAIMockContext().WithSQL("CREATE TABLE target_table (id INT); CREATE TABLE source_table (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE target_table SET id = 1 WHERE id IN (SELECT id FROM source_table)",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using temporary"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: DELETE 使用临时表",
		"DELETE FROM target_table WHERE id IN (SELECT id FROM source_table);",
		session.NewAIMockContext().WithSQL("CREATE TABLE target_table (id INT); CREATE TABLE source_table (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM target_table WHERE id IN (SELECT id FROM source_table)",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using temporary"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: UNION 使用临时表 (从xml中补充)",
		"SELECT * FROM customers WHERE id = 1 UNION SELECT * FROM customers WHERE id = 2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM customers WHERE id = 1 UNION SELECT * FROM customers WHERE id = 2",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using temporary"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: JOIN ORDER BY 使用临时表 (从xml中补充)",
		"SELECT * FROM customers a JOIN customers_records b USING(id) ORDER BY b.log_date DESC;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT); CREATE TABLE customers_records (id INT, log_date DATE);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM customers a JOIN customers_records b USING(id) ORDER BY b.log_date DESC",
				Rows:  sqlmock.NewRows([]string{"Extra"}).AddRow("Using temporary"),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
