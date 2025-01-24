package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQL00086(t *testing.T) {
	ruleName := ai.SQLE00086
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT 语句中存在后缀匹配模糊检索", "SELECT * FROM users WHERE username LIKE '%ab';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT 语句中存在子字符串匹配模糊检索", "SELECT * FROM users WHERE username LIKE '%ab%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT 语句中不存在模糊检索", "SELECT * FROM users WHERE username LIKE 'ab%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: INSERT...SELECT 语句中存在后缀匹配模糊检索", "INSERT INTO archive_users SELECT * FROM users WHERE username LIKE '%ab';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50)); CREATE TABLE archive_users (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: INSERT...SELECT 语句中存在子字符串匹配模糊检索", "INSERT INTO archive_users SELECT * FROM users WHERE username LIKE '%ab%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50)); CREATE TABLE archive_users (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: UPDATE...SET 语句中存在后缀匹配模糊检索", "UPDATE users SET status = 'inactive' WHERE username LIKE '%ab';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50), status VARCHAR(10))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: UPDATE...SET 语句中存在子字符串匹配模糊检索", "UPDATE users SET status = 'inactive' WHERE username LIKE '%ab%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50), status VARCHAR(10))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: DELETE 语句中存在后缀匹配模糊检索", "DELETE FROM users WHERE username LIKE '%ab';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: DELETE 语句中存在子字符串匹配模糊检索", "DELETE FROM users WHERE username LIKE '%ab%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: UNION ALL 语句中存在后缀匹配模糊检索", "SELECT * FROM users WHERE username LIKE '%ab' UNION ALL SELECT * FROM admins WHERE username LIKE '%cd';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50)); CREATE TABLE admins (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: UNION ALL 语句中存在子字符串匹配模糊检索", "SELECT * FROM users WHERE username LIKE '%ab%' UNION ALL SELECT * FROM admins WHERE username LIKE '%cd%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50)); CREATE TABLE admins (id INT PRIMARY KEY, username VARCHAR(50))"),
		nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 12: WITH 语句中存在后缀匹配模糊检索", "WITH user_data AS (SELECT * FROM users WHERE username LIKE '%ab') SELECT * FROM user_data;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50))"),
	// 	nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 13: WITH 语句中存在子字符串匹配模糊检索", "WITH user_data AS (SELECT * FROM users WHERE username LIKE '%ab%') SELECT * FROM user_data;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, username VARCHAR(50))"),
	// 	nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SELECT 语句中存在后缀匹配模糊检索", "SELECT * FROM t2 WHERE name LIKE '%t';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: SELECT 语句中存在子字符串匹配模糊检索", "SELECT * FROM t2 WHERE name LIKE '%t%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: INSERT...SELECT 语句中存在后缀匹配模糊检索", "INSERT INTO t2 SELECT id, name, type, addr FROM t2 WHERE name LIKE '%t';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: INSERT...SELECT 语句中存在子字符串匹配模糊检索", "INSERT INTO t2 SELECT id, name, type, addr FROM t2 WHERE name LIKE '%t%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: UPDATE...SET 语句中存在后缀匹配模糊检索", "UPDATE t2 SET type = '3' WHERE name LIKE '%t';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 19: UPDATE...SET 语句中存在子字符串匹配模糊检索", "UPDATE t2 SET type = '3' WHERE name LIKE '%t%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 20: DELETE 语句中存在后缀匹配模糊检索", "DELETE FROM t2 WHERE name LIKE '%t';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 21: DELETE 语句中存在子字符串匹配模糊检索", "DELETE FROM t2 WHERE name LIKE '%t%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 22: SELECT 语句中不存在模糊检索", "SELECT * FROM t2 WHERE name LIKE 't%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 23: INSERT...SELECT 语句中不存在模糊检索", "INSERT INTO t2 SELECT id, name, type, addr FROM t2 WHERE name LIKE 't%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 24: UPDATE...SET 语句中不存在模糊检索", "UPDATE t2 SET type = '3' WHERE name LIKE 't%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 25: DELETE 语句中不存在模糊检索", "DELETE FROM t2 WHERE name LIKE 't%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 26: SELECT 语句中存在后_检索", "SELECT * FROM t2 WHERE name LIKE '_t';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 27: SELECT 语句中存在子字符串_检索", "SELECT * FROM t2 WHERE name LIKE '_t_';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 28: INSERT...SELECT 语句中存在后缀_检索", "INSERT INTO t2 SELECT id, name, type, addr FROM t2 WHERE name LIKE '_t';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 29: INSERT...SELECT 语句中存在子字符串_检索", "INSERT INTO t2 SELECT id, name, type, addr FROM t2 WHERE name LIKE '_t%';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 30: UPDATE...SET 语句中存在后缀_检索", "UPDATE t2 SET type = '3' WHERE name LIKE '_t';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 31: UPDATE...SET 语句中存在子字符串_检索", "UPDATE t2 SET type = '3' WHERE name LIKE '_t_';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 32: DELETE 语句中存在后缀_检索", "DELETE FROM t2 WHERE name LIKE '_t';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 33: DELETE 语句中存在子字符串_检索", "DELETE FROM t2 WHERE name LIKE '_t_';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t2 (id INT PRIMARY KEY, name VARCHAR(50), type INT, addr VARCHAR(100))"),
		nil, newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
