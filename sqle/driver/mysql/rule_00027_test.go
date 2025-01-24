package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00027(t *testing.T) {
	ruleName := ai.SQLE00027
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, no problem
	runSingleRuleInspectCase(rule, t, "create table, no problem", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test comment",
	v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test comment",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test comment",
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, with problem (no comment for column)
	runSingleRuleInspectCase(rule, t, "create table, with problem (no comment for column)", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test comment",
	v1 varchar(255) NOT NULL DEFAULT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test",
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, "v1,v2"))

	//create table, with problem (only space comment for column)
	runSingleRuleInspectCase(rule, t, "create table, with problem (only space comment for column)", DefaultMysqlInspect(), `
		CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
		id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test comment",
		v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "   ",
		v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "   ",
		PRIMARY KEY (id)
		);
		`, newTestResult().addResult(ruleName, "v1,v2"))

	//create table, with problem (empty comment for column)
	runSingleRuleInspectCase(rule, t, "create table, with problem (empty comment for column)", DefaultMysqlInspect(), `
				CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
				id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test comment",
				v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "",
				v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "",
				PRIMARY KEY (id)
				);
				`, newTestResult().addResult(ruleName, "v1,v2"))

	//alter table add column, no problem
	runSingleRuleInspectCase(rule, t, "alter table add column, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int COMMENT "unit test comment";
	`, newTestResult())

	//alter table add column, with problem (no comment for column)
	runSingleRuleInspectCase(rule, t, "alter table add column, with problem (no comment for column)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int;
	`, newTestResult().addResult(ruleName, "v"))

	//alter table add column, with problem (only space comment for column)
	runSingleRuleInspectCase(rule, t, "alter table add column, with problem (only space comment for column)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int COMMENT "   ";
	`, newTestResult().addResult(ruleName, "v"))

	//alter table add column, with problem (empty comment for column)
	runSingleRuleInspectCase(rule, t, "alter table add column, with problem (empty comment for column)", DefaultMysqlInspect(), `
		ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int COMMENT "";
		`, newTestResult().addResult(ruleName, "v"))

	//alter table modify column, no problem
	runSingleRuleInspectCase(rule, t, "alter table modify column, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 int COMMENT "unit test comment";
	`, newTestResult())

	//alter table modify column, with problem (no comment for column)
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (no comment for column)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 int;
	`, newTestResult().addResult(ruleName, "v1"))

	//alter table modify column, with problem (only space for column)
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (only space for column)", DefaultMysqlInspect(), `
		ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 int COMMENT "   ";
		`, newTestResult().addResult(ruleName, "v1"))

	//alter table modify column, with problem (empty for column)
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (empty for column)", DefaultMysqlInspect(), `
				ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 int COMMENT "";
				`, newTestResult().addResult(ruleName, "v1"))

	//alter table change column, no problem
	runSingleRuleInspectCase(rule, t, "alter table modify column, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a int COMMENT "unit test comment";
	`, newTestResult())

	//alter table modify column, with problem (no comment for column)
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (no comment for column)", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a int;
	`, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, with problem (only space comment for column)
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (only space comment for column)", DefaultMysqlInspect(), `
		ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a int COMMENT "   ";
		`, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, with problem (empty comment for column)
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (empty comment for column)", DefaultMysqlInspect(), `
				ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a int COMMENT "";
				`, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
