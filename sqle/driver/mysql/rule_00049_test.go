package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQL00049(t *testing.T) {
	ruleName := ai.SQLE00049
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// ===== CREATE
	// CREATE USER
	runSingleRuleInspectCase(rule, t, "CREATE USER", DefaultMysqlInspect(), `
	CREATE USER 'TABLE'@'localhost' REQUIRE NONE;
	`, newTestResult().addResult(ruleName))

	// CREATE DATABASE
	runSingleRuleInspectCase(rule, t, "CREATE DATABASE ...", DefaultMysqlInspect(),
		"CREATE DATABASE `INT`;", newTestResult().addResult(ruleName))

	// CREATE TABLE
	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... ", DefaultMysqlInspect(),
		"CREATE TABLE exist_db.1no_exist_tb (`COLUMN` INT);", newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... ", DefaultMysqlInspect(),
		"CREATE TABLE `TABLE` (id INT);", newTestResult().addResult(ruleName))

	// CREATE VIEW
	runSingleRuleInspectCase(rule, t, " CREATE VIEW ... ", DefaultMysqlInspect(),
		"CREATE VIEW `VIEW` AS SELECT name FROM exist_db.exist_tb_1;", newTestResult().addResult(ruleName))

	// CREATE INDEX
	runSingleRuleInspectCase(rule, t, " CREATE INDEX ... ", DefaultMysqlInspect(),
		"CREATE INDEX `INDEX` ON exist_db.exist_tb_1 (id);", newTestResult().addResult(ruleName))

	// CREATE EVENT
	runSingleRuleInspectCase(rule, t, " CREATE EVENT ... ", DefaultMysqlInspect(),
		"CREATE EVENT CREATE ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;", newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	// ===== ALTER
	// ALTER TABLE ... ADD COLUMN
	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... ", DefaultMysqlInspect(),
		"ALTER TABLE exist_db.exist_tb_1 ADD `INT` INT;", newTestResult().addResult(ruleName))

	// ALTER TABLE ... ADD INDEX
	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... ", DefaultMysqlInspect(),
		"ALTER TABLE exist_db.exist_tb_1 ADD INDEX `SELECT`(id);", newTestResult().addResult(ruleName))

	// ====check2
	// ALTER TABLE ... RENAME TO
	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME .... ", DefaultMysqlInspect(),
		"ALTER TABLE exist_db.exist_tb_1 RENAME TO `RENAME`", newTestResult().addResult(ruleName))

	// ALTER TABLE RENAME INDEX
	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME INDEX ... ", DefaultMysqlInspect(),
		"ALTER TABLE exist_db.exist_tb_1 RENAME INDEX idx_1 TO `ALTER`;", newTestResult().addResult(ruleName))

	// RENAME TABLE ... TO ...
	runSingleRuleInspectCase(rule, t, " RENAME TABLE ....", DefaultMysqlInspect(),
		"RENAME TABLE exist_db.exist_tb_1 TO exist_db.`exists`;", newTestResult().addResult(ruleName))

	// ALTER TABLE ... CHANGE ...
	runSingleRuleInspectCase(rule, t, " ALTER TABLE CHANGE .... column names start with a number", DefaultMysqlInspect(),
		"ALTER TABLE exist_db.exist_tb_1 CHANGE id `EXPLAIN` INT;", newTestResult().addResult(ruleName))

	// ALTER EVENT ...  TO ...
	runAIRuleCase(rule, t, " ALTER EVENT ... ", "ALTER EVENT ename RENAME TO HAVING;",
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runAIRuleCase(rule, t, " ALTER EVENT ... ", "ALTER EVENT ename RENAME TO LIMIT;",
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	//RENAME USER
	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		"RENAME USER 't1'@'localhost' TO 't2'@'%','user1'@'%' TO 'CREATE'@'%'", newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		"RENAME USER user1 to BLOB", newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

}

// ==== Rule test code end ====
