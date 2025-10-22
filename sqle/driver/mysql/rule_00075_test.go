package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00075(t *testing.T) {
	ruleName := ai.SQLE00075
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, no table charset, no column charset - should pass
	runSingleRuleInspectCase(rule, t, "create table, no table charset, no column charset", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 100 AUTO_INCREMENT,
	a varchar(10),
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, with table charset utf8mb4, no column charset - should pass
	runSingleRuleInspectCase(rule, t, "create table, with table charset utf8mb4, no column charset", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 100 AUTO_INCREMENT,
	a varchar(10),
	PRIMARY KEY (id)
	) CHARSET utf8mb4;
	`, newTestResult())

	//create table, with table charset utf8mb4, column charset utf8mb4 - should pass
	runSingleRuleInspectCase(rule, t, "create table, with table charset utf8mb4, column charset utf8mb4", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 100 AUTO_INCREMENT,
	a varchar(10) CHARSET utf8mb4,
	PRIMARY KEY (id)
	) CHARSET utf8mb4;
	`, newTestResult())

	//create table, with table charset utf8mb4, column charset utf8 - should fail
	runSingleRuleInspectCase(rule, t, "create table, with table charset utf8mb4, column charset utf8", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 100 AUTO_INCREMENT,
	a varchar(10) CHARSET utf8,
	PRIMARY KEY (id)
	) CHARSET utf8mb4;
	`, newTestResult().addResult(ruleName, "a"))

	//create table, with table charset utf8, column charset utf8mb4 - should fail
	runSingleRuleInspectCase(rule, t, "create table, with table charset utf8, column charset utf8mb4", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 100 AUTO_INCREMENT,
	a varchar(10) CHARSET utf8mb4,
	PRIMARY KEY (id)
	) CHARSET utf8;
	`, newTestResult().addResult(ruleName, "a"))

	//alter table add column, no table charset change, no column charset - should pass
	runSingleRuleInspectCase(rule, t, "alter table add column, no table charset change, no column charset", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a varchar(10) COMMENT "unit test";
	`, newTestResult())

	//alter table add column, no table charset change, column charset utf8mb4 - should pass (assuming table charset is utf8mb4)
	runSingleRuleInspectCase(rule, t, "alter table add column, no table charset change, column charset utf8mb4", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a varchar(10) CHARSET utf8mb4 COMMENT "unit test";
	`, newTestResult())

	//alter table add column, change table charset to utf8, column charset utf8mb4 - should fail
	runSingleRuleInspectCase(rule, t, "alter table add column, change table charset to utf8, column charset utf8mb4", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a varchar(10) CHARSET utf8mb4 COMMENT "unit test", CONVERT TO CHARACTER SET utf8;
	`, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, no table charset change, no column charset - should pass
	runSingleRuleInspectCase(rule, t, "alter table modify column, no table charset change, no column charset", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY v1 varchar(10) COMMENT "unit test";
	`, newTestResult())

	//alter table modify column, no table charset change, column charset utf8mb4 - should pass (assuming table charset is utf8mb4)
	runSingleRuleInspectCase(rule, t, "alter table modify column, no table charset change, column charset utf8mb4", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY v1 varchar(10) CHARSET utf8mb4 COMMENT "unit test";
	`, newTestResult())

	//alter table modify column, change table charset to utf8, column charset utf8mb4 - should fail
	runSingleRuleInspectCase(rule, t, "alter table modify column, change table charset to utf8, column charset utf8mb4", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY v1 varchar(10) CHARSET utf8mb4 COMMENT "unit test", CONVERT TO CHARACTER SET utf8;
	`, newTestResult().addResult(ruleName, "v1"))

	//alter table change column, no table charset change, no column charset - should pass
	runSingleRuleInspectCase(rule, t, "alter table change column, no table charset change, no column charset", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a varchar(10) COMMENT "unit test";
	`, newTestResult())

	//alter table change column, no table charset change, column charset utf8mb4 - should pass (assuming table charset is utf8mb4)
	runSingleRuleInspectCase(rule, t, "alter table change column, no table charset change, column charset utf8mb4", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a varchar(10) CHARSET utf8mb4 COMMENT "unit test";
	`, newTestResult())

	//alter table change column, change table charset to utf8, column charset utf8mb4 - should fail
	runSingleRuleInspectCase(rule, t, "alter table change column, change table charset to utf8, column charset utf8mb4", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a varchar(10) CHARSET utf8mb4 COMMENT "unit test", CONVERT TO CHARACTER SET utf8;
	`, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
