package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQL00056(t *testing.T) {
	ruleName := ai.SQLE00056
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE specify character set as latin1", "CREATE TABLE exist_db.no_exist_tb_1 (id INT) DEFAULT CHARSET=latin1;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName, "UTF8MB4"))

	runAIRuleCase(rule, t, "case 2: CREATE TABLE specify character set as utf8mb4", "CREATE TABLE exist_db.no_exist_tb_1 (id INT) DEFAULT CHARSET=utf8mb4;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: CREATE TABLE unspecified character set", "CREATE TABLE exist_db.no_exist_tb_1 (id INT);",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: ALTER TABLE specify character set as latin1", "ALTER TABLE exist_db.exist_tb_1 DEFAULT CHARSET=latin1;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName, "UTF8MB4"))

	runAIRuleCase(rule, t, "case 5: ALTER TABLE specify character set as utf8mb4", "ALTER TABLE exist_db.exist_tb_1 DEFAULT CHARSET=utf8mb4;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: ALTER TABLE unspecified character set", "ALTER TABLE exist_db.exist_tb_1 ADD COLUMN name VARCHAR(255);",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: CREATE DATABASE specify character set as latin1", "CREATE DATABASE no_exist_db DEFAULT CHARACTER SET latin1;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName, "UTF8MB4"))

	runAIRuleCase(rule, t, "case 8: CREATE DATABASE specify character set as utf8mb4", "CREATE DATABASE no_exist_db DEFAULT CHARACTER SET utf8mb4;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: CREATE DATABASE unspecified character set", "CREATE DATABASE no_exist_db;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: ALTER DATABASE specify character set as latin1", "ALTER DATABASE exist_db DEFAULT CHARACTER SET latin1;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName, "UTF8MB4"))

	runAIRuleCase(rule, t, "case 11: ALTER DATABASE specify character set as utf8mb4", "ALTER DATABASE exist_db DEFAULT CHARACTER SET utf8mb4;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 13: CREATE TABLE specify character set as KOI8R", "CREATE TABLE exist_db.no_exist_tb_1 (c1 VARCHAR(100)) CHARSET=KOI8R;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName, "UTF8MB4"))

	runAIRuleCase(rule, t, "case 14: CREATE TABLE specify character set as UTF8MB4", "CREATE TABLE exist_db.no_exist_tb_1 (c1 VARCHAR(100)) CHARSET=UTF8MB4;",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 15: CREATE TABLE unspecified character set", "CREATE TABLE exist_db.no_exist_tb_1 (c1 VARCHAR(100));",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 16: ALTER TABLE convert to character setutf8mb4", "ALTER TABLE exist_db.exist_tb_1 CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;",
		nil, /*mock context*/
		nil, newTestResult())
}

// ==== Rule test code end ====
