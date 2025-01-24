package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00089(t *testing.T) {
	ruleName := ai.SQLE00089
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: 使用 INSERT INTO 语句插入单个值，不违反规则",
		"INSERT INTO employees (id, name) VALUES (1, 'Alice');",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE employees (id INT, name VARCHAR(255));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 2: 使用 INSERT INTO...SELECT 语句，违反规则",
		"INSERT INTO employees (id, name) SELECT id, name FROM contractors;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE employees (id INT, name VARCHAR(255)); CREATE TABLE contractors (id INT, name VARCHAR(255));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	// runAIRuleCase(rule, t, "case 3: 使用 INSERT INTO 语句并包含 WITH 子句，违反规则",
	// 	"INSERT INTO employees (id, name) WITH temp AS (SELECT id, name FROM interns) SELECT id, name FROM temp;",
	// 	session.NewAIMockContext().WithSQL(
	// 		"CREATE TABLE employees (id INT, name VARCHAR(255)); CREATE TABLE interns (id INT, name VARCHAR(255));",
	// 	),
	// 	nil,
	// 	newTestResult().addResult(ruleName),
	// )

	runAIRuleCase(rule, t, "case 4: 使用 INSERT INTO 语句设置多个列值，不违反规则",
		"INSERT INTO employees SET id = 2, name = 'Bob';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE employees (id INT, name VARCHAR(255));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 5: 使用 INSERT INTO 语句的 VALUES 子句中包含子查询，不违反规则",
		"INSERT INTO employees (id, name) VALUES (3, (SELECT name FROM managers WHERE id = 10));",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE employees (id INT, name VARCHAR(255)); CREATE TABLE managers (id INT, name VARCHAR(255));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 6: 使用 INSERT INTO 语句插入多个值，不违反规则 (从xml中补充)",
		"INSERT INTO employees (id, name) VALUES (4, 'Charlie'), (5, 'David');",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE employees (id INT, name VARCHAR(255));",
		),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 7: 使用 INSERT INTO...SELECT 语句从同一表中选择数据，违反规则 (从xml中补充)",
		"INSERT INTO employees (id, name) SELECT id, name FROM employees WHERE id > 5;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE employees (id INT, name VARCHAR(255));",
		),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 9: 使用 INSERT INTO...SELECT 语句并包含 JOIN，违反规则 (从xml中补充)",
		"INSERT INTO employees (id, name) SELECT contractors.id, contractors.name FROM contractors JOIN projects ON contractors.project_id = projects.id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE employees (id INT, name VARCHAR(255)); CREATE TABLE contractors (id INT, name VARCHAR(255), project_id INT); CREATE TABLE projects (id INT);",
		),
		nil,
		newTestResult().addResult(ruleName),
	)
}

// ==== Rule test code end ====
