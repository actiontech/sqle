package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQL00047(t *testing.T) {
	ruleName := ai.SQLE00047
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// ===== CREATE
	// CREATE USER
	runSingleRuleInspectCase(rule, t, "CREATE USER", DefaultMysqlInspect(), `
	CREATE USER 'Test_useruseruseruseruseruseruseruseruseruseruseruseruseruseruseruser'@'localhost' REQUIRE NONE;
	`, newTestResult().addResult(ruleName, 64))

	// CREATE DATABASE
	runSingleRuleInspectCase(rule, t, "CREATE DATABASE ...", DefaultMysqlInspect(), `
	CREATE DATABASE no_exist_dbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdbdb;
	`, newTestResult().addResult(ruleName, 64))

	// CREATE TABLE
	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... ", DefaultMysqlInspect(),
		`CREATE TABLE exist_db.1no_exist_tb (idxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx INT);`, newTestResult().addResult(ruleName, 64))

	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... ", DefaultMysqlInspect(),
		`CREATE TABLE no_exist_tbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtbtb (id INT);`, newTestResult().addResult(ruleName, 64))

	// CREATE VIEW
	runSingleRuleInspectCase(rule, t, " CREATE VIEW ... ", DefaultMysqlInspect(),
		`CREATE VIEW order_viewviewviewviewviewviewviewviewviewviewviewviewviewviewview AS SELECT name FROM exist_db.exist_tb_1;`, newTestResult().addResult(ruleName, 64))

	// CREATE INDEX
	runSingleRuleInspectCase(rule, t, " CREATE INDEX ... ", DefaultMysqlInspect(),
		`CREATE INDEX 1ord_id_idxidxidxidxidxidxidxidxidxidxidxidxidxidxidxidxidxidxidxidx ON exist_db.exist_tb_1 (id);`, newTestResult().addResult(ruleName, 64))

	// CREATE EVENT
	runSingleRuleInspectCase(rule, t, " CREATE EVENT ... ", DefaultMysqlInspect(),
		`CREATE EVENT ename_eventeventeventeventeventeventeventeventeventeventeventevent ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName, 64))

	// ===== ALTER
	// ALTER TABLE ... ADD COLUMN
	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... ", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD c_columncolumncolumncolumncolumncolumncolumncolumncolumncolumncolumn INT;`, newTestResult().addResult(ruleName, 64))

	// ALTER TABLE ... ADD INDEX
	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... ", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD INDEX idx_2(id);`, newTestResult())

	// ====check2
	// ALTER TABLE ... RENAME TO
	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME .... ", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME TO new_tabletabletabletabletabletabletabletabletabletabletabletabletable`, newTestResult().addResult(ruleName, 64))

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
	runAIRuleCase(rule, t, " ALTER EVENT ... ", `ALTER EVENT ename RENAME TO newname_eventeventeventeventeventeventeventeventeventeventeventevent;`,
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName, 64))

	runAIRuleCase(rule, t, " ALTER EVENT ... ", `ALTER EVENT ename RENAME TO new_name;`,
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	//RENAME USER
	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		`RENAME USER 't1'@'localhost' TO 't2'@'%','user1'@'%' TO 'user2222222222222222222222222222222222222222222222222222222222222'@'%'`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName, 64))

	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		`RENAME USER user1 to user2`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

}

// ==== Rule test code end ====
