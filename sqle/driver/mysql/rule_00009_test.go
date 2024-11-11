package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00009(t *testing.T) {
	ruleName := ai.SQLE00009
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// runAIRuleCase(rule, t, "case 1: DELETE 语句的 WHERE 子句中对字段应用函数但没有函数索引",
	// 	"DELETE FROM employees WHERE substr(name,2,8) = 'JOHN';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME, INDEX idx_trim_name (name));"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "SHOW INDEX FROM `exist_db`.`employees`",
	// 			Rows:  sqlmock.NewRows([]string{"Expression"}).AddRow("substr(`name`, 2, 8)"),
	// 		},
	// 	},
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 1: DELETE 语句的 WHERE 子句中对字段应用函数但没有函数索引",
	// 	"DELETE FROM employees WHERE name = 'JOHN';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME, INDEX idx_trim_name (name));"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 4: INSERT...SELECT 语句中 SELECT 子句的 WHERE 条件应用函数且没有函数索引",
	// 	"INSERT INTO archived_employees SELECT * FROM employees WHERE LOWER(status) = 'active';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
	// 		WithSQL("CREATE TABLE archived_employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "SHOW INDEX FROM `exist_db`.`employees`",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 5: INSERT...SELECT 语句中 SELECT 子句的 WHERE 条件应用函数且存在函数索引",
	// 	"INSERT INTO archived_employees SELECT * FROM employees WHERE YEAR(created_at) = 2022;",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME, INDEX idx_year_created_at (YEAR(created_at)));").
	// 		WithSQL("CREATE TABLE archived_employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "SHOW INDEX FROM `exist_db`.`employees`",
	// 			Rows:  sqlmock.NewRows([]string{"Expression"}).AddRow("year(`created_at`)"),
	// 		},
	// 	},
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 6: INSERT...SELECT 语句中 SELECT 子句的 WHERE 条件未应用函数",
	// 	"INSERT INTO archived_employees SELECT * FROM employees WHERE status = 'active';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
	// 		WithSQL("CREATE TABLE archived_employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 7: SELECT 语句的 WHERE 子句中对字段应用函数且没有函数索引",
	// 	"SELECT * FROM employees WHERE LENGTH(name) > 5;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "SHOW INDEX FROM `exist_db`.`employees`",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 8: SELECT 语句的 WHERE 子句中对字段应用函数且存在函数索引",
	// 	"SELECT * FROM employees WHERE TRIM(name) = 'John';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "SHOW INDEX FROM `exist_db`.`employees`",
	// 			Rows:  sqlmock.NewRows([]string{"Expression"}).AddRow("trim(`name`)"),
	// 		},
	// 	},
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 9: SELECT 语句的 WHERE 子句中没有对字段应用函数",
	// 	"SELECT * FROM employees WHERE name = 'John';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 10: UPDATE 语句的 WHERE 子句中对字段应用函数且没有函数索引",
	// 	"UPDATE employees SET status = 'inactive' WHERE LOWER(name) = 'john';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "SHOW INDEX FROM `exist_db`.`employees`",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 11: UPDATE 语句的 WHERE 子句中对字段应用函数且存在函数索引",
	// 	"UPDATE employees SET status = 'inactive' WHERE YEAR(created_at) = 2021;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "SHOW INDEX FROM `exist_db`.`employees`",
	// 			Rows:  sqlmock.NewRows([]string{"Expression"}).AddRow("year(`created_at`)"),
	// 		},
	// 	},
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 12: UPDATE 语句的 WHERE 子句中没有对字段应用函数",
	// 	"UPDATE employees SET status = 'inactive' WHERE name = 'John';",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	// runAIRuleCase(rule, t, "case 13: UNION 语句中一个 SELECT 子句的 WHERE 条件应用函数且没有函数索引",
	// 	"SELECT name FROM employees WHERE UPPER(status) = 'ACTIVE' UNION SELECT name FROM contractors WHERE name = 'Jane';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
	// 		WithSQL("CREATE TABLE contractors (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	[]*AIMockSQLExpectation{
	// 		{
	// 			Query: "SHOW INDEX FROM `exist_db`.`employees`",
	// 			Rows:  sqlmock.NewRows(nil),
	// 		},
	// 	},
	// 	newTestResult().addResult(ruleName),
	// )

	// runAIRuleCase(rule, t, "case 14: UNION 语句中所有 SELECT 子句的 WHERE 条件未应用函数",
	// 	"SELECT name FROM employees WHERE status = 'active' UNION SELECT name FROM contractors WHERE name = 'Jane';",
	// 	session.NewAIMockContext().
	// 		WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
	// 		WithSQL("CREATE TABLE contractors (id INT, name VARCHAR(50), created_at DATETIME);"),
	// 	nil,
	// 	newTestResult(),
	// )

	runAIRuleCase(rule, t, "case 15: UNION 语句中所有 SELECT 子句的 WHERE 条件应用函数且存在函数索引",
		"SELECT name FROM employees WHERE DATE(created_at) = '2023-01-01' UNION SELECT name FROM contractors WHERE LOWER(name) = 'jane';",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), status VARCHAR(20), created_at DATETIME);").
			WithSQL("CREATE TABLE contractors (id INT, name VARCHAR(50), created_at DATETIME);"),
		[]*AIMockSQLExpectation{
			{
				Query: "SHOW INDEX FROM `exist_db`.`employees`",
				Rows:  sqlmock.NewRows([]string{"Expression"}).AddRow("date(`created_at`)"),
			},
			{
				Query: "SHOW INDEX FROM `exist_db`.`contractors`",
				Rows:  sqlmock.NewRows([]string{"Expression"}).AddRow("lower(`name`)"),
			},
		},
		newTestResult(),
	)
}

// ==== Rule test code end ====
