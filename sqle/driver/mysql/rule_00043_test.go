package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00043(t *testing.T) {
	ruleName := ai.SQLE00043
	rule := rulepkg.RuleHandlerMap[ruleName].Rule
	ruleParams := []interface{}{2}

	//create table, no index
	runSingleRuleInspectCase(rule, t, "create table, no index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
	);
	`, newTestResult())

	//create table, with index, no column repeat
	runSingleRuleInspectCase(rule, t, "create table, with index, no column repeat", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	PRIMARY KEY (id),
	INDEX index_1 (id)
	);
	`, newTestResult())

	//create table, with index, column repeat
	runSingleRuleInspectCase(rule, t, "create table, with index, column repeat", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a int,
	b int,
	PRIMARY KEY (id),
	INDEX index_1 (id, a),
	INDEX index_2 (a),
	INDEX index_3 (b, a)	
	);
	`, newTestResult().addResult(ruleName, "a", ruleParams[0]))

	//create table, with index in column definition, column repeat
	runSingleRuleInspectCase(rule, t, "create table, with index in column definition, column repeat", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a int UNIQUE,
	b int,
	PRIMARY KEY (id),
	INDEX index_1 (id, a),
	INDEX index_2 (b, a)
	);
	`, newTestResult().addResult(ruleName, "a", ruleParams[0]))

	//create index, no column repeat
	runSingleRuleInspectCase(rule, t, "create index, no column repeat", DefaultMysqlInspect(), `
	CREATE INDEX index_1 on exist_db.exist_tb_3(v2);
	`, newTestResult())

	//create index, with column repeat
	runSingleRuleInspectCase(rule, t, "create index, with column repeat", DefaultMysqlInspect(), `
	CREATE INDEX index_1 on exist_db.exist_tb_1(id, v1);
	`, newTestResult().addResult(ruleName, "v1", ruleParams[0]))

	//alter table, no index
	runSingleRuleInspectCase(rule, t, "alter table, no index", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHARSET=utf8mb4;
	`, newTestResult())

	//alter table, with index, no column repeat
	runSingleRuleInspectCase(rule, t, "alter table, with index, no column repeat", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD INDEX (v2);
	`, newTestResult())

	//alter table, with index, column repeat
	runSingleRuleInspectCase(rule, t, "alter table, with index, column repeat", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD INDEX (id, v1);
	`, newTestResult().addResult(ruleName, "v1", ruleParams[0]))

}

// ==== Rule test code end ====
