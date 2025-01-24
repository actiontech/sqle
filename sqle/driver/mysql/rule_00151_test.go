package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00151(t *testing.T) {
	ruleName := ai.SQLE00151
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, without tablespace
	runSingleRuleInspectCase(rule, t, "create table, without tablespace", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    );
    `, newTestResult())

	//create table, with tablespace
	runSingleRuleInspectCase(rule, t, "create table, with tablespace", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    PRIMARY KEY (id)
    ) TABLESPACE innodb_system;
    `, newTestResult().addResult(ruleName))

	//alter table, without tablespace
	runSingleRuleInspectCase(rule, t, "alter table, without tablespace", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 CHARACTER SET 'utf8mb4';
    `, newTestResult())

	//alter table, with tablespace
	runSingleRuleInspectCase(rule, t, "alter table, with tablespace", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 TABLESPACE innodb_system;
    `, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
