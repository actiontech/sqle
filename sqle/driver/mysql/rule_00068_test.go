package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00068(t *testing.T) {
	ruleName := ai.SQLE00068
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, no timestamp column
	runSingleRuleInspectCase(rule, t, "create table, no timestamp column", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned DEFAULT 1000,
    a int,
    PRIMARY KEY (id)
    );
    `, newTestResult())

	//create table, with timestamp column
	runSingleRuleInspectCase(rule, t, "create table, with timestamp column", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned DEFAULT 1000,
    a int,
    ts timestamp,
    PRIMARY KEY (id)
    );
    `, newTestResult().addResult(ruleName, "ts"))

	//alter table add column, no timestamp column
	runSingleRuleInspectCase(rule, t, "alter table add column, no timestamp column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a int NOT NULL COMMENT "unit test";
    `, newTestResult())

	//alter table add column, with timestamp column
	runSingleRuleInspectCase(rule, t, "alter table add column, with timestamp column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN ts timestamp NOT NULL COMMENT "unit test";
    `, newTestResult().addResult(ruleName, "ts"))

	//alter table modify column, no timestamp column
	runSingleRuleInspectCase(rule, t, "alter table modify column, no timestamp column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 int;
    `, newTestResult())

	//alter table modify column, with timestamp column
	runSingleRuleInspectCase(rule, t, "alter table modify column, with timestamp column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 timestamp;
    `, newTestResult().addResult(ruleName, "v1"))

	//alter table change column, no timestamp column
	runSingleRuleInspectCase(rule, t, "alter table change column, no timestamp column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a int;
    `, newTestResult())

	//alter table change column, with timestamp column
	runSingleRuleInspectCase(rule, t, "alter table change column, with timestamp column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a timestamp;
    `, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
