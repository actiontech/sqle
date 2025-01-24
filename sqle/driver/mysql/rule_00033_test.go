package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00033(t *testing.T) {
	ruleName := ai.SQLE00033
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule
	ruleParams := []interface{}{"UPDATE_TIME"}

	//create table, without update_time column
	runSingleRuleInspectCase(rule, t, "create table, without update_time column", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create table, with update_time column, without DEFAULT
	runSingleRuleInspectCase(rule, t, "create table, with update_time column, without DEFAULT", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	UPDATE_TIME datetime,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create table, with update_time column, with DEFAULT value not current_timestamp
	runSingleRuleInspectCase(rule, t, "create table, with update_time column, with DEFAULT value not current_timestamp", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	UPDATE_TIME datetime DEFAULT 0,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create table, with update_time column, with DEFAULT value current_timestamp, without ON UPDATE
	runSingleRuleInspectCase(rule, t, "create table, with update_time column, with DEFAULT value current_timestamp", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	UPDATE_TIME datetime DEFAULT current_timestamp,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create table, with update_time column, with DEFAULT value current_timestamp, with ON UPDATE value current_timestamp
	runSingleRuleInspectCase(rule, t, "create table, with update_time column, with DEFAULT value current_timestamp", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	UPDATE_TIME datetime DEFAULT current_timestamp ON UPDATE current_timestamp,
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//alter table, add not update_time column
	runSingleRuleInspectCase(rule, t, "alter table, add not update_time column", DefaultMysqlInspect(), `
		ALTER TABLE exist_db.exist_tb_1 ADD UPDATE_TIME2 int;
		`, newTestResult())

	//alter table, add update_time column, without DEFAULT
	runSingleRuleInspectCase(rule, t, "alter table, add update_time column, without DEFAULT", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD UPDATE_TIME datetime COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, add update_time column, with DEFAULT value not CURRENT_TIMESTAMP
	runSingleRuleInspectCase(rule, t, "alter table, add update_time column, with DEFAULT value not CURRENT_TIMESTAMP", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD UPDATE_TIME datetime DEFAULT 0 COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, add update_time column, with DEFAULT value CURRENT_TIMESTAMP, not ON UPDATE
	runSingleRuleInspectCase(rule, t, "alter table, add update_time column, with DEFAULT value CURRENT_TIMESTAMP, not ON UPDATE", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD UPDATE_TIME datetime DEFAULT CURRENT_TIMESTAMP COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, add update_time column, with DEFAULT value CURRENT_TIMESTAMP
	runSingleRuleInspectCase(rule, t, "alter table, add update_time column, with DEFAULT value CURRENT_TIMESTAMP", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD UPDATE_TIME datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE current_timestamp COMMENT "unit test";
	`, newTestResult())

	//alter table, modify update_time column, without DEFAULT
	runSingleRuleInspectCase(rule, t, "alter table, modify update_time column, without DEFAULT", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY UPDATE_TIME datetime COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, modify update_time column, with DEFAULT value not CURRENT_TIMESTAMP
	runSingleRuleInspectCase(rule, t, "alter table, modify update_time column, with DEFAULT value not CURRENT_TIMESTAMP", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY UPDATE_TIME datetime DEFAULT 0 COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, modify update_time column, with DEFAULT value CURRENT_TIMESTAMP, not ON UPDATE
	runSingleRuleInspectCase(rule, t, "alter table, modify update_time column, with DEFAULT value CURRENT_TIMESTAMP, not ON UPDATE", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY UPDATE_TIME datetime DEFAULT CURRENT_TIMESTAMP COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, modify update_time column, with DEFAULT value CURRENT_TIMESTAMP
	runSingleRuleInspectCase(rule, t, "alter table, modify update_time column, with DEFAULT value CURRENT_TIMESTAMP", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY UPDATE_TIME datetime DEFAULT CURRENT_TIMESTAMP  ON UPDATE current_timestamp COMMENT "unit test";
	`, newTestResult())

	//alter table, change update_time column, without DEFAULT
	runSingleRuleInspectCase(rule, t, "alter table, change update_time column, without DEFAULT", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE v1 UPDATE_TIME datetime COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, change update_time column, with DEFAULT value not CURRENT_TIMESTAMP
	runSingleRuleInspectCase(rule, t, "alter table, change update_time column, with DEFAULT value not CURRENT_TIMESTAMP", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE v1 UPDATE_TIME datetime DEFAULT 0 COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, change update_time column, with DEFAULT value CURRENT_TIMESTAMP, not ON UPDATE
	runSingleRuleInspectCase(rule, t, "alter table, change update_time column, with DEFAULT value CURRENT_TIMESTAMP, not ON UPDATE", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE v1 UPDATE_TIME datetime DEFAULT CURRENT_TIMESTAMP  COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "UPDATE_TIME"))

	//alter table, change update_time column, with DEFAULT value CURRENT_TIMESTAMP
	runSingleRuleInspectCase(rule, t, "alter table, change update_time column, with DEFAULT value CURRENT_TIMESTAMP", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE v1 UPDATE_TIME datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE current_timestamp COMMENT "unit test";
	`, newTestResult())

}

// ==== Rule test code end ====
