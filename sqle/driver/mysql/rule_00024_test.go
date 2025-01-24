package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00024(t *testing.T) {
	ruleName := ai.SQLE00024
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, no SET column
	runSingleRuleInspectCase(rule, t, "create table, no SET column", DefaultMysqlInspect(), `
 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
 id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
 );
 `, newTestResult())

	//create table, with SET column
	runSingleRuleInspectCase(rule, t, "create table, with SET column", DefaultMysqlInspect(), `
 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
 id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
 a SET('reading', 'writing', 'painting', 'cooking')
 );
 `, newTestResult().addResult(ruleName, "a"))

	//alter table add column, no SET column
	runSingleRuleInspectCase(rule, t, "alter table add column, no SET column", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a int NOT NULL COMMENT "unit test";
 `, newTestResult())

	//alter table add column, with SET column
	runSingleRuleInspectCase(rule, t, "alter table add column, with SET column", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a set('reading', 'writing', 'painting', 'cooking') NOT NULL COMMENT "unit test";
 `, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, no SET column
	runSingleRuleInspectCase(rule, t, "alter table modify column, no SET column", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 Modify COLUMN v1 int NOT NULL COMMENT "unit test";
 `, newTestResult())

	//alter table modify column, with SET column
	runSingleRuleInspectCase(rule, t, "alter table modify column, with SET column", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 Modify COLUMN v1 set('reading', 'writing', 'painting', 'cooking') NOT NULL COMMENT "unit test";
 `, newTestResult().addResult(ruleName, "v1"))

	//alter table change column, no SET column
	runSingleRuleInspectCase(rule, t, "alter table modify column, no SET column", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v int NOT NULL COMMENT "unit test";
 `, newTestResult())

	//alter table modify column, with SET column
	runSingleRuleInspectCase(rule, t, "alter table modify column, with SET column", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v set('reading', 'writing', 'painting', 'cooking') NOT NULL COMMENT "unit test";
 `, newTestResult().addResult(ruleName, "v"))
}

// ==== Rule test code end ====
