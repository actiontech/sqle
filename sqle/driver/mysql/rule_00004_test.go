package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00004(t *testing.T) {
	ruleName := ai.SQLE00004
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, without auto_increment option
	runSingleRuleInspectCase(rule, t, "create table, without auto_increment option", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned DEFAULT 100 AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    );
    `, newTestResult())

	//create table, with auto_increment 0
	runSingleRuleInspectCase(rule, t, "create table, with auto_increment 0", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    ) auto_increment = 0;
    `, newTestResult())

	//create table, with auto_increment other than 0
	runSingleRuleInspectCase(rule, t, "create table, with auto_increment other than 0", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned DEFAULT 1 AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    ) auto_increment = 1;
    `, newTestResult().addResult(ruleName))

	//alter table, without auto_increment option
	runSingleRuleInspectCase(rule, t, "alter table, without auto_increment option", DefaultMysqlInspect(), `
    alter table exist_db.exist_tb_1 charset=utf8mb4;
    `, newTestResult())

	//alter table, with auto_increment 0
	runSingleRuleInspectCase(rule, t, "alter table, with auto_increment 0", DefaultMysqlInspect(), `
    alter table exist_db.exist_tb_1 auto_increment=0;
    `, newTestResult())

	//alter table, with auto_increment other than 0
	runSingleRuleInspectCase(rule, t, "alter table, with auto_increment other than 0", DefaultMysqlInspect(), `
    alter table exist_db.exist_tb_1 auto_increment=1;
    `, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
