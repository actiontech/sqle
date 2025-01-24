package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// "testing"

// rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
// "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"

// ==== Rule test code start ====
func TestRuleSQL00046(t *testing.T) {
	ruleName := ai.SQLE00046
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// ===== CREATE
	// CREATE USER
	runSingleRuleInspectCase(rule, t, "CREATE USER", DefaultMysqlInspect(), `
	CREATE USER 'Test'@'localhost' REQUIRE NONE;
	`, newTestResult().addResult(ruleName))

	// CREATE DATABASE
	runSingleRuleInspectCase(rule, t, "CREATE DATABASE ...", DefaultMysqlInspect(), `
	CREATE DATABASE No_exist_db;
	`, newTestResult().addResult(ruleName))

	// CREATE TABLE
	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... ", DefaultMysqlInspect(),
		`CREATE TABLE exist_db.1no_exist_tb (Id INT);`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... ", DefaultMysqlInspect(),
		`CREATE TABLE no_Exist_TB (id INT);`, newTestResult().addResult(ruleName))

	// CREATE VIEW
	runSingleRuleInspectCase(rule, t, " CREATE VIEW ... ", DefaultMysqlInspect(),
		`CREATE VIEW order_View AS SELECT name FROM exist_db.exist_tb_1;`, newTestResult().addResult(ruleName))

	// CREATE INDEX
	runSingleRuleInspectCase(rule, t, " CREATE INDEX ... ", DefaultMysqlInspect(),
		`CREATE INDEX 1ord_id_IDX ON exist_db.exist_tb_1 (id);`, newTestResult().addResult(ruleName))

	// CREATE EVENT
	runSingleRuleInspectCase(rule, t, " CREATE EVENT ... ", DefaultMysqlInspect(),
		`CREATE EVENT Ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	// ===== ALTER
	// ALTER TABLE ... ADD COLUMN
	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... ", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD c_CC INT;`, newTestResult().addResult(ruleName))

	// ALTER TABLE ... ADD INDEX
	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... ", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD INDEX idx_2(id);`, newTestResult())

	// ====check2
	// ALTER TABLE ... RENAME TO
	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME .... ", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME TO new_TABLE`, newTestResult().addResult(ruleName))

	// ALTER TABLE RENAME INDEX
	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME INDEX ... ", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME INDEX idx_1 TO new_idx_1;`, newTestResult())

	// RENAME TABLE ... TO ...
	runSingleRuleInspectCase(rule, t, " RENAME TABLE ....", DefaultMysqlInspect(),
		`RENAME TABLE exist_db.exist_tb_1 TO exist_db.exist_tb_xxx;`, newTestResult())

	// ALTER TABLE ... CHANGE ...
	runSingleRuleInspectCase(rule, t, " ALTER TABLE CHANGE .... column names start with a number", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 CHANGE id new_id INT;`, newTestResult())

	// ALTER EVENT ...  TO ...
	runAIRuleCase(rule, t, " ALTER EVENT ... ", `ALTER EVENT ename RENAME TO ENAME;`,
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runAIRuleCase(rule, t, " ALTER EVENT ... ", `ALTER EVENT ename RENAME TO NEW_NAME;`,
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	//RENAME USER
	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		`RENAME USER 't1'@'localhost' TO 't2'@'%','user1'@'%' TO 'User11'@'%'`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		`RENAME USER user1 to user2`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

}

// ==== Rule test code end ====
