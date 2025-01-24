package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00120(t *testing.T) {
	ruleName := ai.SQLE00120
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: DELETE 语句中使用 IN (NULL)",
		"DELETE FROM users WHERE id IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: DELETE 语句中使用 NOT IN (NULL)",
		"DELETE FROM users WHERE id NOT IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: DELETE 语句中不使用 IN 或 NOT IN (NULL)",
		"DELETE FROM users WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: INSERT...SELECT 语句中使用 IN (NULL)",
		"INSERT INTO archive_users SELECT * FROM users WHERE id IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100)); CREATE TABLE archive_users (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: INSERT...SELECT 语句中使用 NOT IN (NULL)",
		"INSERT INTO archive_users SELECT * FROM users WHERE id NOT IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100)); CREATE TABLE archive_users (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: INSERT...SELECT 语句中不使用 IN 或 NOT IN (NULL)",
		"INSERT INTO archive_users SELECT * FROM users WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100)); CREATE TABLE archive_users (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: SELECT 语句中使用 IN (NULL)",
		"SELECT * FROM users WHERE id IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: SELECT 语句中使用 NOT IN (NULL)",
		"SELECT * FROM users WHERE id NOT IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: SELECT 语句中不使用 IN 或 NOT IN (NULL)",
		"SELECT * FROM users WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: UPDATE 语句中使用 IN (NULL)",
		"UPDATE users SET name = 'John' WHERE id IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: UPDATE 语句中使用 NOT IN (NULL)",
		"UPDATE users SET name = 'John' WHERE id NOT IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: UPDATE 语句中不使用 IN 或 NOT IN (NULL)",
		"UPDATE users SET name = 'John' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	// runAIRuleCase(rule, t, "case 13: WITH 语句中使用 IN (NULL)",
	// 	"WITH user_data AS (SELECT * FROM users WHERE id IN (NULL)) SELECT * FROM user_data;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
	// 	nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 14: WITH 语句中使用 NOT IN (NULL)",
	// 	"WITH user_data AS (SELECT * FROM users WHERE id NOT IN (NULL)) SELECT * FROM user_data;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
	// 	nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 15: WITH 语句中不使用 IN 或 NOT IN (NULL)",
	// 	"WITH user_data AS (SELECT * FROM users WHERE id = 1) SELECT * FROM user_data;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT, name VARCHAR(100));"),
	// 	nil, newTestResult())

	runAIRuleCase(rule, t, "case 16: SELECT 语句中使用 IN (NULL) (从xml中补充)",
		"SELECT * FROM t1 WHERE mark1 IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (mark1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: SELECT 语句中使用 NOT IN (NULL) (从xml中补充)",
		"SELECT * FROM t1 WHERE mark1 NOT IN (NULL);",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (mark1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: SELECT 语句中使用 IS NULL (从xml中补充)",
		"SELECT * FROM t1 WHERE mark1 IS NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (mark1 INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 19: SELECT 语句中使用 IS NOT NULL (从xml中补充)",
		"SELECT * FROM t1 WHERE mark1 IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (mark1 INT);"),
		nil, newTestResult())
}

// ==== Rule test code end ====
