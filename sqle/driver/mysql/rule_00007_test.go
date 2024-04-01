package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00007(t *testing.T) {
	ruleName := ai.SQLE00007
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, no auto-increment column
	runSingleRuleInspectCase(rule, t, "create table, no auto-increment column", DefaultMysqlInspect(), `
		 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
		 id bigint unsigned NOT NULL COMMENT "unit test"
		 );
		 `, newTestResult())

	//create table, with one auto-increment column
	runSingleRuleInspectCase(rule, t, "create table, with auto-increment column", DefaultMysqlInspect(), `
		 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
		 id bigint unsigned NOT NULL COMMENT "unit test",
		 id2 bigint unsigned AUTO_INCREMENT COMMENT "unit test"
		 );
		 `, newTestResult())

	//create table, with multiple auto-increment column
	runSingleRuleInspectCase(rule, t, "create table, with multiple auto-increment column", DefaultMysqlInspect(), `
		 CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
		 id bigint unsigned NOT NULL COMMENT "unit test",
		 id2 bigint unsigned AUTO_INCREMENT COMMENT "unit test",
		 id3 bigint unsigned AUTO_INCREMENT
		 );
		 `, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
