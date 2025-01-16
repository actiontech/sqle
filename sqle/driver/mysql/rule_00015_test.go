package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00015(t *testing.T) {
	ruleName := ai.SQLE00015
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE 未指定表级 COLLATION，列未指定 COLLATION",
		"CREATE TABLE employees (id INT, name VARCHAR(100), description TEXT);",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 2: CREATE TABLE 指定表级 COLLATION，与数据库默认 COLLATION 一致",
		"CREATE TABLE departments (id INT, dept_name VARCHAR(100)) COLLATE utf8_general_ci;",
		nil,
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 3: CREATE TABLE 指定表级 COLLATION，与数据库默认 COLLATION 不一致",
		"CREATE TABLE projects (id INT, project_name VARCHAR(100)) COLLATE utf8_unicode_ci;",
		nil,
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 4: CREATE TABLE 未指定表级 COLLATION，但列级 COLLATION 与数据库默认 COLLATION 一致",
		"CREATE TABLE salaries (id INT, amount DECIMAL(10,2) COLLATE utf8_general_ci);",
		nil,
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 5: CREATE TABLE 未指定表级 COLLATION，但列级 COLLATION 与数据库默认 COLLATION 不一致",
		"CREATE TABLE bonuses (id INT, bonus_amount DECIMAL(10,2) COLLATE utf8_unicode_ci);",
		nil,
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 6: CREATE TABLE 指定表级 COLLATION 与默认一致，所有列级 COLLATION 与默认一致",
		"CREATE TABLE benefits (id INT, benefit_name VARCHAR(100) COLLATE utf8_general_ci) COLLATE utf8_general_ci;",
		nil,
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 7: CREATE TABLE 指定表级 COLLATION 与默认一致，但存在列级 COLLATION 与默认不一致",
		"CREATE TABLE leaves (id INT, leave_type VARCHAR(100) COLLATE utf8_unicode_ci) COLLATE utf8_general_ci;",
		nil,
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 8: CREATE TABLE 未指定表级 COLLATION，有多列，部分列指定 COLLATION 与默认一致，部分不一致",
		"CREATE TABLE attendance (id INT, status VARCHAR(50) COLLATE utf8_general_ci, remarks TEXT COLLATE utf8_unicode_ci);",
		nil,
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 9: ALTER TABLE 不包含 CONVERT TO CHARACTER SET，未修改字符类型列",
		"ALTER TABLE employees ADD COLUMN hire_date DATE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(100), description TEXT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 10: ALTER TABLE 包含 CONVERT TO CHARACTER SET，COLLATION 与默认一致",
		"ALTER TABLE departments CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci;",
		session.NewAIMockContext().WithSQL("CREATE TABLE departments (id INT, dept_name VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 11: ALTER TABLE 包含 CONVERT TO CHARACTER SET，COLLATION 与默认不一致",
		"ALTER TABLE projects CONVERT TO CHARACTER SET utf8 COLLATE utf8_unicode_ci;",
		session.NewAIMockContext().WithSQL("CREATE TABLE projects (id INT, project_name VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 12: ALTER TABLE 添加字符类型列，未指定 COLLATION",
		"ALTER TABLE salaries ADD COLUMN currency VARCHAR(10);",
		session.NewAIMockContext().WithSQL("CREATE TABLE salaries (id INT, amount DECIMAL(10,2));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 13: ALTER TABLE 添加字符类型列，指定 COLLATION 与默认一致",
		"ALTER TABLE bonuses ADD COLUMN bonus_type VARCHAR(50) COLLATE utf8_general_ci;",
		session.NewAIMockContext().WithSQL("CREATE TABLE bonuses (id INT, bonus_amount DECIMAL(10,2));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 14: ALTER TABLE 添加字符类型列，指定 COLLATION 与默认不一致",
		"ALTER TABLE benefits ADD COLUMN benefit_description TEXT COLLATE utf8_unicode_ci;",
		session.NewAIMockContext().WithSQL("CREATE TABLE benefits (id INT, benefit_name VARCHAR(100) COLLATE utf8_general_ci) COLLATE utf8_general_ci;"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 15: ALTER TABLE 修改字符类型列，指定 COLLATION 与默认一致",
		"ALTER TABLE leaves MODIFY COLUMN leave_type VARCHAR(100) COLLATE utf8_general_ci;",
		session.NewAIMockContext().WithSQL("CREATE TABLE leaves (id INT, leave_type VARCHAR(100) COLLATE utf8_unicode_ci) COLLATE utf8_general_ci;"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 16: ALTER TABLE 修改字符类型列，指定 COLLATION 与默认不一致",
		"ALTER TABLE attendance MODIFY COLUMN remarks TEXT COLLATE utf8_unicode_ci;",
		session.NewAIMockContext().WithSQL("CREATE TABLE attendance (id INT, status VARCHAR(50) COLLATE utf8_general_ci, remarks TEXT COLLATE utf8_unicode_ci);"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 17: CREATE TABLE 指定表级 COLLATION 与默认一致，但存在列级 COLLATION 与默认不一致 (从xml中补充)",
		"CREATE TABLE table_a (id INT, col_1 VARCHAR(50) COLLATE utf8mb4_bin, col_2 VARCHAR(50)) COLLATE utf8mb4_bin;",
		nil,
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8_general_ci"),
			},
		},
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 18: CREATE TABLE 未指定表级 COLLATION，列未指定 COLLATION (从xml中补充)",
		"CREATE TABLE table_b (id INT, col_1 VARCHAR(50), col_2 VARCHAR(50));",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 19: ALTER TABLE CONVERT TO CHARACTER SET 与默认一致 (从xml中补充)",
		"ALTER TABLE table_b CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table_b (id INT, col_1 VARCHAR(50), col_2 VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@collation_database",
				Rows:  sqlmock.NewRows([]string{"@@collation_database"}).AddRow("utf8mb4_bin"),
			},
		},
		newTestResult(),
	)
}

// ==== Rule test code end ====
