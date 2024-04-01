package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00018(t *testing.T) {
	ruleName := ai.SQLE00018
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, no problem
	runSingleRuleInspectCase(rule, t, "create table, no problem", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a char(10)
	);
	`, newTestResult())

	//create table, with problem (string type, greater than expected length))
	runSingleRuleInspectCase(rule, t, "create table, with problem (string type, greater than expected length))", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a char(21)
	);
	`, newTestResult().addResult(ruleName, "a"))

	//alter table add columns, no problem
	runSingleRuleInspectCase(rule, t, "alter table add columns, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a char(10) COMMENT "unit test";
	`, newTestResult())

	//alter table add columns, with problem (string type, greater than expected length))
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (string type, greater than expected length))", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a char(21) COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table modify columns, no problem
	runSingleRuleInspectCase(rule, t, "alter table modify columns, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 char(10) COMMENT "unit test";
	`, newTestResult())

	//alter table modify columns, with problem (string type, greater than expected length))
	runSingleRuleInspectCase(rule, t, "alter table modify columns, with problem (string type, greater than expected length))", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 char(21) COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "v1"))

	//alter table change columns, no problem
	runSingleRuleInspectCase(rule, t, "alter table change columns, no problem", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a char(10) COMMENT "unit test";
	`, newTestResult())

	//alter table change columns, with problem (string type, greater than expected length))
	runSingleRuleInspectCase(rule, t, "alter table change columns, with problem (string type, greater than expected length))", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a char(21) COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
