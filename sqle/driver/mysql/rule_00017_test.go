package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00017(t *testing.T) {
	ruleName := ai.SQLE00017
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, no problem
	runSingleRuleInspectCase(rule, t, "create table, no problem", DefaultMysqlInspect(),
		`
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    );
    `, newTestResult())

	//create table, with problem (blob type)
	runSingleRuleInspectCase(rule, t, "create table, with problem (blob type)", DefaultMysqlInspect(),
		`
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    a int,
    b blob,
    PRIMARY KEY (id)
    );
    `, newTestResult().addResult(ruleName, "b"))

	//create table, with problem (text type)
	runSingleRuleInspectCase(rule, t, "create table, with problem (text type)", DefaultMysqlInspect(),
		`
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    a int,
    b text,
    PRIMARY KEY (id)
    );
    `, newTestResult().addResult(ruleName, "b"))

	//alter table add columns, no problem
	runSingleRuleInspectCase(rule, t, "alter table add columns, no problem", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a int;
    `, newTestResult())

	//alter table add columns, with problem (blob type)
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (blob type)", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a blob;
    `, newTestResult().addResult(ruleName, "a"))

	//alter table add columns, with problem (text type)
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (text type)", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a text;
    `, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, no problem
	runSingleRuleInspectCase(rule, t, "alter table modify column, no problem", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 MODIFY v1 int;
    `, newTestResult())

	//alter table modify column, with problem (blob type)
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (blob type)", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 MODIFY v1 blob;
    `, newTestResult().addResult(ruleName, "v1"))

	//alter table modify column, with problem (text type)
	runSingleRuleInspectCase(rule, t, "alter table modify column, with problem (text type)", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 MODIFY v1 text;
    `, newTestResult().addResult(ruleName, "v1"))

	//alter table change column, no problem
	runSingleRuleInspectCase(rule, t, "alter table change column, no problem", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v int;
    `, newTestResult())

	//alter table change column, with problem (blob type)
	runSingleRuleInspectCase(rule, t, "alter table change column, with problem (blob type)", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v blob;
    `, newTestResult().addResult(ruleName, "v"))

	//alter table change column, with problem (text type)
	runSingleRuleInspectCase(rule, t, "alter table change column, with problem (text type)", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v text;
    `, newTestResult().addResult(ruleName, "v"))
}

// ==== Rule test code end ====
