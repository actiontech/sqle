package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00036(t *testing.T) {
	ruleName := ai.SQLE00036
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, without blob column in index
	runSingleRuleInspectCase(rule, t, "create table, without blob column in index", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    );
    `, newTestResult())

	//create table, with blob column in index
	runSingleRuleInspectCase(rule, t, "create table, with blob column in index", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    a int,
    b blob UNIQUE,
    PRIMARY KEY (id)
    );
    `, newTestResult().addResult(ruleName))

	//create index, without blob column
	runSingleRuleInspectCase(rule, t, "create index, without blob column", DefaultMysqlInspect(), `
    CREATE INDEX idx_1 on exist_db.exist_tb_12 (v3, v2);
    `, newTestResult())

	//create index, with blob column
	runSingleRuleInspectCase(rule, t, "create index, with blob column", DefaultMysqlInspect(), `
    CREATE INDEX idx_1 on exist_db.exist_tb_12 (v1);
    `, newTestResult().addResult(ruleName))

	//alter table, without blob column in index
	runSingleRuleInspectCase(rule, t, "alter table, without blob column in index", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_12 ADD INDEX idx_3 (v2, v3) COMMENT "unit test";
    `, newTestResult())

	//alter table, with blob column in index
	runSingleRuleInspectCase(rule, t, "alter table, with blob column in index", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_12 ADD INDEX idx_3 (v1, v2) COMMENT "unit test";
    `, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
