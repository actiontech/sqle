package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00022(t *testing.T) {
	ruleName := ai.SQLE00022
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule
	ruleParams := []interface{}{40}

	//create table, no problem
	runSingleRuleInspectCase(rule, t, "create table, no problem", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    );
    `, newTestResult())

	//create table, with too many columns
	runSingleRuleInspectCase(rule, t, "create table, with too many columns", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    aaaa1 int,
    aaaa2 int,
    aaaa3 int,
    aaaa4 int,
    aaaa5 int,
    aaaa6 int,
    aaaa7 int,
    aaaa8 int,
    aaaa9 int,
    aaaa10 int,
    aaaa11 int,
    aaaa12 int,
    aaaa13 int,
    aaaa14 int,
    aaaa15 int,
    aaaa16 int,
    aaaa17 int,
    aaaa18 int,
    aaaa19 int,
    aaaa20 int,
    aaaa21 int,
    aaaa22 int,
    aaaa23 int,
    aaaa24 int,
    aaaa25 int,
    aaaa26 int,
    aaaa27 int,
    aaaa28 int,
    aaaa29 int,
    aaaa30 int,
    aaaa31 int,
    aaaa32 int,
    aaaa33 int,
    aaaa34 int,
    aaaa35 int,
    aaaa36 int,
    aaaa37 int,
    aaaa38 int,
    aaaa39 int,
    aaaa40 int
    );
    `, newTestResult().addResult(ruleName, ruleParams...))
}

// ==== Rule test code end ====
