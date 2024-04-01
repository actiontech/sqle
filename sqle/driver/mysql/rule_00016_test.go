package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00016(t *testing.T) {
	ruleName := ai.SQLE00016
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, blob/text column is with DEFAULT value NULL, no problem
	runSingleRuleInspectCase(rule, t, "create table, no problem", DefaultMysqlInspect(), `
	 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	 id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	 a blob DEFAULT NULL,
	 b text DEFAULT NULL,
	 PRIMARY KEY (id)
	 );
	 `, newTestResult())

	//create table, with problem (blob type, with default value other than NULL)
	runSingleRuleInspectCase(rule, t, "create table, with problem (blob type, with default value)", DefaultMysqlInspect(), `
	 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	 id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	 a blob DEFAULT 'a',
	 PRIMARY KEY (id)
	 );
	 `, newTestResult().addResult(ruleName, "a"))

	//create table, with problem (text type, with default value other than NULL)
	runSingleRuleInspectCase(rule, t, "create table, with problem (text type, with default value)", DefaultMysqlInspect(), `
	 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	 id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	 a text DEFAULT 'a',
	 PRIMARY KEY (id)
	 );
	 `, newTestResult().addResult(ruleName, "a"))

	//alter table add columns, blob/text column is without DEFAULT value, no problem
	runSingleRuleInspectCase(rule, t, "alter table add columns, no problem", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a blob DEFAULT NULL;
	 `, newTestResult())

	//alter table add columns, with problem (blob type, with default value other than NULL)
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (blob type, no default value)", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a blob DEFAULT 'b';
	 `, newTestResult().addResult(ruleName, "a"))

	//alter table add columns, with problem (text type, with default value other than NULL)
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (text type, no default value)", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a text DEFAULT 'b';
	 `, newTestResult().addResult(ruleName, "a"))

	//alter table modify columns, blob/text column is without DEFAULT value, no problem
	runSingleRuleInspectCase(rule, t, "alter table add columns, no problem", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 blob DEFAULT NULL;
	 `, newTestResult())

	//alter table modify columns, with problem (blob type, with default value other than NULL)
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (blob type, no default value)", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 blob DEFAULT 'b';
	 `, newTestResult().addResult(ruleName, "v1"))

	//alter table modify columns, with problem (text type, with default value other than NULL)
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (text type, no default value)", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 text DEFAULT 'b';
	 `, newTestResult().addResult(ruleName, "v1"))

	//alter table change columns, blob/text column is without DEFAULT value, no problem
	runSingleRuleInspectCase(rule, t, "alter table add columns, no problem", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a blob DEFAULT NULL;
	 `, newTestResult())

	//alter table change columns, with problem (blob type, with default value other than NULL)
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (blob type, no default value)", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a blob DEFAULT 'b';
	 `, newTestResult().addResult(ruleName, "a"))

	//alter table change columns, with problem (text type, with default value other than NULL)
	runSingleRuleInspectCase(rule, t, "alter table add columns, with problem (text type, no default value)", DefaultMysqlInspect(), `
	 ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a text DEFAULT 'b';
	 `, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
