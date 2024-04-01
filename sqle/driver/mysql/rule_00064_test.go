package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00064(t *testing.T) {
	ruleName := ai.SQLE00064
	rule := rulepkg.RuleHandlerMap[ruleName].Rule
	ruleParams := []interface{}{1024}

	//create table, no varchar column
	runSingleRuleInspectCase(rule, t, "create table, no varchar column", DefaultMysqlInspect(), `
 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
 id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
 PRIMARY KEY (id)
 );
 `, newTestResult())

	//create table, with varchar column, less than expected length
	runSingleRuleInspectCase(rule, t, "create table, with varchar column, less than expected length", DefaultMysqlInspect(), `
 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
 id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
 v1 varchar(100),
 PRIMARY KEY (id)
 );
 `, newTestResult())

	//create table, with varchar column, greater than expected length
	runSingleRuleInspectCase(rule, t, "create table, with varchar column, greater than expected length", DefaultMysqlInspect(), `
 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
 id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
 v1 varchar(1025),
 PRIMARY KEY (id)
 );
 `, newTestResult().addResult(ruleName, ruleParams...))

	//alter table, no varchar column
	runSingleRuleInspectCase(rule, t, "alter table, no varchar column", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 ADD COLUMN new_col int NOT NULL COMMENT "unit test";
 `, newTestResult())

	//alter table, with varchar column, less than expected length
	runSingleRuleInspectCase(rule, t, "alter table, with varchar column, less than expected length", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 ADD COLUMN new_col varchar(100) NOT NULL COMMENT "unit test";
 `, newTestResult())

	//alter table, with varchar column, greater than expected length
	runSingleRuleInspectCase(rule, t, "alter table, with varchar column, greater than expected length", DefaultMysqlInspect(), `
 ALTER TABLE exist_db.exist_tb_1 ADD COLUMN new_col varchar(1025) NOT NULL COMMENT "unit test";
 `, newTestResult().addResult(ruleName, ruleParams...))
}

// ==== Rule test code end ====
