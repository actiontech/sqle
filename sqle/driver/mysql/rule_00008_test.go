package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00008(t *testing.T) {
	ruleName := ai.SQLE00008
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, primary key not specified
	runSingleRuleInspectCase(rule, t, "create table, primary key not specified", DefaultMysqlInspect(),
		`
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
    );
    `, newTestResult().addResult(ruleName))

	//create table, primary key specified in table constraint
	runSingleRuleInspectCase(rule, t, "create table, primary key specified", DefaultMysqlInspect(),
		`
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    );
    `, newTestResult())

	//create table, primary key specified in column definition
	runSingleRuleInspectCase(rule, t, "create table, primary key specified", DefaultMysqlInspect(),
	`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test"
	);
	`, newTestResult())
}

// ==== Rule test code end ====
