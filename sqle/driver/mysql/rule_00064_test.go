package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00064(t *testing.T) {
	ruleName := ai.SQLE00064
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, no index
	runSingleRuleInspectCase(rule, t, "create table, no index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
	);
	`, newTestResult())

	//create table, with index, no varchar index
	runSingleRuleInspectCase(rule, t, "create table, with index, no varchar index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, with index, with varchar index less than expected length
	runSingleRuleInspectCase(rule, t, "create table, with index, with varchar index less than expected length", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 varchar(100),
	PRIMARY KEY (id),
	INDEX idx_1 (v1)
	);
	`, newTestResult())

	//create table, with index, with varchar index greater than expected length
	runSingleRuleInspectCase(rule, t, "create table, with index, with varchar index greater than expected length", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 varchar(1000),
	PRIMARY KEY (id),
	INDEX idx_1 (v1)
	);
	`, newTestResult().addResult(ruleName, "v1"))

	//create index, no problem
	runSingleRuleInspectCase(rule, t, "create index, no problem", DefaultMysqlInspect(), `
	CREATE INDEX idx_1 on exist_db.exist_tb_3(v3);
	`, newTestResult())

	//create index, with problem (varchar index greater than expected length)
	runSingleRuleInspectCase(rule, t, "create index, with problem (varchar index greater than expected length)", DefaultMysqlInspect(), `
	CREATE INDEX idx_1 on exist_db.exist_tb_12(v3);
	`, newTestResult().addResult(ruleName, "v3"))

	//alter table, no index
	runSingleRuleInspectCase(rule, t, "alter table, no index", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHARACTER SET 'utf8mb4';
	`, newTestResult())

	//alter table, with index, no varchar index
	runSingleRuleInspectCase(rule, t, "alter table, with index, no varchar index", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD INDEX idx_2(v1);
	`, newTestResult())

	//alter table, with index, with varchar index less than expected length
	runSingleRuleInspectCase(rule, t, "alter table, with index, with varchar index less than expected length", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD INDEX idx_2(v1);
	`, newTestResult())

	//alter table, with index, with varchar index greater than expected length
	runSingleRuleInspectCase(rule, t, "alter table, with index, with varchar index greater than expected length", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_12 ADD INDEX idx_2(v3);
	`, newTestResult().addResult(ruleName, "v3"))
}

// ==== Rule test code end ====
