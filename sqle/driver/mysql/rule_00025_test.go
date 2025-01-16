package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00025(t *testing.T) {
	ruleName := ai.SQLE00025
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, timestamp column with DEFAULT value
	runSingleRuleInspectCase(rule, t, "create table, no timestamp column", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	ts timestamp default current_timestamp,
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, timestamp column without DEFAULT value
	runSingleRuleInspectCase(rule, t, "create table, timestamp column without DEFAULT value", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	ts timestamp,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, "ts"))

	//alter table add column, timestamp column with DEFAULT value
	runSingleRuleInspectCase(rule, t, "alter table add column, timestamp column with DEFAULT value", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN ts timestamp DEFAULT current_timestamp NOT NULL COMMENT "unit test";
	`, newTestResult())

	//alter table add column, timestamp column without DEFAULT value
	runSingleRuleInspectCase(rule, t, "alter table add column, timestamp column without DEFAULT value", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN ts timestamp NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "ts"))

	//alter table modify column, timestamp column with DEFAULT value
	runSingleRuleInspectCase(rule, t, "alter table modify column, timestamp column with DEFAULT value", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 timestamp DEFAULT current_timestamp NOT NULL COMMENT "unit test";
	`, newTestResult())

	//alter table modify column, timestamp column without DEFAULT value
	runSingleRuleInspectCase(rule, t, "alter table modify column, timestamp column without DEFAULT value", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 timestamp NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "v1"))

	//alter table change column, timestamp column with DEFAULT value
	runSingleRuleInspectCase(rule, t, "alter table change column, timestamp column with DEFAULT value", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a timestamp DEFAULT current_timestamp NOT NULL COMMENT "unit test";
	`, newTestResult())

	//alter table change column, timestamp column without DEFAULT value
	runSingleRuleInspectCase(rule, t, "alter table change column, timestamp column without DEFAULT value", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a timestamp NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
