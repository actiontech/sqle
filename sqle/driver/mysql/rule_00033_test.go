package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00033(t *testing.T) {
	ruleName := ai.SQLE00033
	rule := rulepkg.RuleHandlerMap[ruleName].Rule
	ruleParams := []interface{}{"UPDATE_TIME"}

	//create table, without update_time column
	runSingleRuleInspectCase(rule, t, "create table, without update_time column", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create table, with update_time column, without DEFAULT
	runSingleRuleInspectCase(rule, t, "create table, with update_time column, without DEFAULT", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	UPDATE_TIME datetime,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create table, with update_time column, with DEFAULT value not current_timestamp
	runSingleRuleInspectCase(rule, t, "create table, with update_time column, with DEFAULT value not current_timestamp", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	UPDATE_TIME datetime DEFAULT 0,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create table, with update_time column, with DEFAULT value current_timestamp, without ON UPDATE
	runSingleRuleInspectCase(rule, t, "create table, with update_time column, with DEFAULT value current_timestamp", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	UPDATE_TIME datetime DEFAULT current_timestamp,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create table, with update_time column, with DEFAULT value current_timestamp, with ON UPDATE value current_timestamp
	runSingleRuleInspectCase(rule, t, "create table, with update_time column, with DEFAULT value current_timestamp", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	UPDATE_TIME datetime DEFAULT current_timestamp ON UPDATE current_timestamp,
	PRIMARY KEY (id)
	);
	`, newTestResult())
}

// ==== Rule test code end ====
