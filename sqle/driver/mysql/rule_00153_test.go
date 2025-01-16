package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00153(t *testing.T) {
	ruleName := ai.SQLE00153
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, no secondary index
	runSingleRuleInspectCase(rule, t, "create table, no secondary index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
	);
	`, newTestResult().addResult(ruleName))

	//create table, with secondary index
	runSingleRuleInspectCase(rule, t, "create table, with secondary index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	INDEX idx_1 (id)
	);
	`, newTestResult())
}

// ==== Rule test code end ====
