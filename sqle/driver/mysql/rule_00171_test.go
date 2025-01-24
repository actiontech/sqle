package mysql

//import (
//	"testing"
//
//	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
//	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
//)
//
//// ==== Rule test code start ====
//func TestRuleSQLE00171(t *testing.T) {
//	ruleName := ai.SQLE00171
//	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule
//
//	ruleParams := []interface{}{"CREATE_TIME"}
//
//	//create table, has create_time column
//	runSingleRuleInspectCase(rule, t, "create table, has create_time column", DefaultMysqlInspect(), `
//	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
//	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
//	CREATE_TIME datetime DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
//	PRIMARY KEY (id)
//	);
//	`, newTestResult())
//
//	//create table, has create_time column, with wrong DEFAULT value
//	runSingleRuleInspectCase(rule, t, "create table, has create_time column", DefaultMysqlInspect(), `
//	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
//	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
//	CREATE_TIME datetime DEFAULT NULL COMMENT "unit test",
//	PRIMARY KEY (id)
//	);
//	`, newTestResult().addResult(ruleName, ruleParams...))
//
//	//create table, has no create_time column
//	runSingleRuleInspectCase(rule, t, "create table, has no create_time column", DefaultMysqlInspect(), `
//	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
//	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
//	PRIMARY KEY (id)
//	);
//	`, newTestResult().addResult(ruleName, ruleParams...))
//}
//
//// ==== Rule test code end ====
