package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQL00048(t *testing.T) {
	ruleName := ai.SQLE00048
	rule := rulepkg.RuleHandlerMap[ruleName].Rule
	// ===== CREATE
	// CREATE DATABASE
	runSingleRuleInspectCase(rule, t, "CREATE DATABASE no_exist_db;", DefaultMysqlInspect(), `
	CREATE DATABASE No_exist_db;
	`, newTestResult())

	runSingleRuleInspectCase(rule, t, "CREATE DATABASE 2no_exist_db;", DefaultMysqlInspect(), `
	CREATE DATABASE 2no_exist_db;
	`, newTestResult().addResult(ruleName))
	// CREATE TABLE
	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... Name starts with a number", DefaultMysqlInspect(),
		`CREATE TABLE exist_db.1no_exist_tb_1 (id INT);`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... name contains illegal characters", DefaultMysqlInspect(),
		`CREATE TABLE exist_db.1no-exist-tb_1 (id INT);`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " CREATE TABLE ... name contains illegal characters", DefaultMysqlInspect(),
		`CREATE TABLE exist_db.no_exist_tb_1 (id INT);`, newTestResult())

	runSingleRuleInspectCase(rule, t, " CREATE TABLE .... column... name contains illegal characters", DefaultMysqlInspect(),
		`CREATE TABLE exist_db.no_exist_tb_1 (_id INT);`, newTestResult().addResult(ruleName))
	// CREATE VIEW
	runSingleRuleInspectCase(rule, t, " CREATE VIEW ... Name starts with a number", DefaultMysqlInspect(),
		`CREATE VIEW 1order_view AS SELECT name FROM exist_db.exist_tb_1;`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " CREATE VIEW ... name contains illegal characters", DefaultMysqlInspect(),
		`CREATE VIEW order-view AS SELECT name FROM order*his;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " CREATE VIEW ... name contains illegal characters", DefaultMysqlInspect(),
		`CREATE VIEW order_view AS SELECT name FROM exist_db.exist_tb_1;`, newTestResult())
	// CREATE INDEX
	runSingleRuleInspectCase(rule, t, " CREATE INDEX ... Name starts with a number", DefaultMysqlInspect(),
		`CREATE INDEX 1ord_id_idx ON exist_db.exist_tb_1 (id);`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " CREATE INDEX ... name contains illegal characters", DefaultMysqlInspect(),
		`CREATE INDEX ord_id_idx ON exist_db.exist_tb_1 (id);`, newTestResult())

	// CREATE EVENT
	runSingleRuleInspectCase(rule, t, " CREATE EVENT ... Name starts with a number", DefaultMysqlInspect(),
		`CREATE EVENT 1ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " CREATE EVENT ... name contains illegal characters", DefaultMysqlInspect(),
		`CREATE EVENT e-name ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " CREATE EVENT ... name contains illegal characters", DefaultMysqlInspect(),
		`CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	// ===== ALTER
	// ALTER TABLE ... ADD COLUMN
	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... column name start with a number", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD 1column INT;`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... column name contains illegal characters", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD column-name INT;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... column names conform to the rules", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD column_name INT;`, newTestResult())

	// ALTER TABLE ... ADD INDEX
	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... column names start with a number", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD INDEX 1idx_2(id);`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... column name contains illegal characters", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD INDEX idx-2(id);`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE ADD .... column names conform to the rules", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 ADD INDEX idx_2(id);`, newTestResult())

	// ====check2
	// ALTER TABLE ... RENAME TO
	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME .... table name start with a number", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME TO 2new_table`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME .... table name contains illegal characters", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME TO new-table;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME .... table names conform to the rules", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME TO new_table;`, newTestResult())

	// ALTER TABLE RENAME INDEX
	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME INDEX ... Name starts with a number", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME INDEX idx_1 TO 1uniq_id;`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME INDEX ... name contains illegal characters", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME INDEX idx_1 TO uniq-id;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE RENAME INDEX ... name contains illegal characters", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 RENAME INDEX idx_1 TO uniq_id;`, newTestResult())

	// RENAME TABLE ... TO ...
	runSingleRuleInspectCase(rule, t, " RENAME TABLE ....", DefaultMysqlInspect(),
		`RENAME TABLE exist_db.exist_tb_1 TO exist_db.123exist_tb_xxx;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE CHANGE .... column name contains illegal characters", DefaultMysqlInspect(),
		`RENAME TABLE exist_db.exist_tb_1 TO exist_db.exist-tb*xxx;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE CHANGE .... column names conform to the rules", DefaultMysqlInspect(),
		`RENAME TABLE exist_db.exist_tb_1 TO exist_db.exist_tb_xxx1, exist_db.exist_tb_2 TO exist_db.exist_tb_xxx2;`, newTestResult())

	// ALTER TABLE ... CHANGE ...
	runSingleRuleInspectCase(rule, t, " ALTER TABLE CHANGE .... column names start with a number", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 CHANGE id 1new_id INT;`, newTestResult().addResult(ruleName))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE CHANGE .... column name contains illegal characters", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 CHANGE id new-id INT;`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, " ALTER TABLE CHANGE .... column names conform to the rules", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1 CHANGE id new_id INT;`, newTestResult())

	// ALTER EVENT ...  TO ...
	runAIRuleCase(rule, t, " ALTER EVENT ... Name starts with a number", `ALTER EVENT ename RENAME TO 1newname;`,
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runAIRuleCase(rule, t, " ALTER EVENT ... name contains illegal characters", `ALTER EVENT ename RENAME TO new-name;`,
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runAIRuleCase(rule, t, " ALTER EVENT ... name contains illegal characters", `ALTER EVENT ename RENAME TO newname;`,
		session.NewAIMockContext().WithSQL("CREATE EVENT ename ON SCHEDULE EVERY 10 SECOND DO DELETE FROM exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	//RENAME USER
	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		`RENAME USER 't1'@'localhost' TO 't2'@'%','user1'@'%' TO 'user2222'@'%'`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))

	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		`RENAME USER user1 to _user1`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runSingleRuleInspectCase(rule, t, "RENAME USER ...", DefaultMysqlInspect(),
		`RENAME USER user1 to 1user1`, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

}

// ==== Rule test code end ====
