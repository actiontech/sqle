package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00013(t *testing.T) {
	ruleName := ai.SQLE00013
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, no problem
	runSingleRuleInspectCase(rule, t, "create table, no problem", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, with problem (float type))
	runSingleRuleInspectCase(rule, t, "create table, with problem (float type)", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a float,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, "a"))

	//create table, with problem (double type))
	runSingleRuleInspectCase(rule, t, "create table, with problem (double type)", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a double,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, "a"))

	//alter table add columns, no problem
	runSingleRuleInspectCase(rule, t, "alter table add columns, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a int NOT NULL COMMENT "unit test";
	`, newTestResult())

	//alter table add columns, with problem (float type))
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (float type)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a float NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table add columns, with problem (double type))
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (double type)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a double NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, no problem
	runSingleRuleInspectCase(rule, t, "alter table modify column, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN a int NOT NULL COMMENT "unit test";
	`, newTestResult())

	//alter table modify column, with problem (float type))
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (float type)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN a float NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, with problem (double type))
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (double type)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN a double NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table change column, no problem
	runSingleRuleInspectCase(rule, t, "alter table change column, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a int NOT NULL COMMENT "unit test";
	`, newTestResult())

	//alter table change column, with problem (float type))
	runSingleRuleInspectCase(rule, t, "alter table change column, with problem (float type)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a float NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table change column, with problem (double type))
	runSingleRuleInspectCase(rule, t, "alter table change column, with problem (double type)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a double NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
