package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00055(t *testing.T) {
	ruleName := ai.SQLE00055
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, no index
	runSingleRuleInspectCase(rule, t, "create table, no index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, with index, no redundant
	runSingleRuleInspectCase(rule, t, "create table, with index, no redundant", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	INDEX idx_1 (v1),
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, with index, with repeat index
	runSingleRuleInspectCase(rule, t, "create table, with index, with repeat index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a int,
	INDEX idx_1 (id),
	INDEX idx_2 (a),
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, "[id]", "[id]"))

	//create table, with index, with no redundant index
	runSingleRuleInspectCase(rule, t, "create table, with index, with no redundant index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	INDEX idx_1 (v1,id),
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, with index, with redundant index
	runSingleRuleInspectCase(rule, t, "create table, with index, with redundant index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
  v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	INDEX idx_1 (id,v1,v2),
	PRIMARY KEY (id, v1)
	);
	`, newTestResult().addResult(ruleName, "[id v1 v2]", "[id v1]"))

	//create table, with index, with repeat index
	runSingleRuleInspectCase(rule, t, "create table, with index, with repeat index", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
		id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
		v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
		v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
		PRIMARY KEY (id),
		INDEX idx_2 (id)
		)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	`, newTestResult().addResult(ruleName, "[id]", "[id]"))

	//create index, with no redundant index
	runSingleRuleInspectCase(rule, t, "create index, with no redundant index", DefaultMysqlInspect(), `
	CREATE INDEX idx_2 on exist_db.exist_tb_9(v4);
	`, newTestResult())

	//create index, with redundant index
	runSingleRuleInspectCase(rule, t, "create index, with column, with redundant index", DefaultMysqlInspect(), `
	CREATE INDEX idx_2 on exist_db.exist_tb_9(v1, v2(10));
	`, newTestResult().addResult(ruleName, "[v1 v2 v3 v4]", "[v1 v2]"))

	//create index, with repeat index
	runSingleRuleInspectCase(rule, t, "create index, with repeat index", DefaultMysqlInspect(), `
	CREATE INDEX idx_2 on exist_db.exist_tb_9(v1,v2,v3, v4);
	`, newTestResult().addResult(ruleName, "[v1 v2 v3 v4]", "[v1 v2 v3 v4]"))

	//alter table, no index
	runSingleRuleInspectCase(rule, t, "alter table, no index", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHARSET=utf8mb4;
	`, newTestResult())

	// Alter table, add index, no redundant index
	runSingleRuleInspectCase(rule, t, "alter table, add index, no redundant index", DefaultMysqlInspect(), `
ALTER TABLE exist_db.exist_tb_9 
ADD INDEX idx_3 (v4);
`, newTestResult())

	// Alter table, add index, with redundant index
	runSingleRuleInspectCase(rule, t, "alter table, add index, with redundant index", DefaultMysqlInspect(), `
ALTER TABLE exist_db.exist_tb_9 
ADD INDEX idx_3 (v1);
`, newTestResult().addResult(ruleName, "[v1 v2 v3 v4]", "[v1]"))

	// Alter table, add index, with repeat index
	runSingleRuleInspectCase(rule, t, "alter table, add index, with repeat index", DefaultMysqlInspect(), `
ALTER TABLE exist_db.exist_tb_9 
ADD INDEX idx_2 (v2,v3),
ADD INDEX idx_3 (v3);
`, newTestResult().addResult(ruleName, "[v2 v3]", "[v2 v3]").addResult(ruleName, "[v3]", "[v3]"))

	// Alter table, drop index, no effect on redundancy
	runSingleRuleInspectCase(rule, t, "alter table, drop index, no effect on redundancy", DefaultMysqlInspect(), `
ALTER TABLE exist_db.exist_tb_1 
DROP INDEX idx_1;
`, newTestResult())
}

// ==== Rule test code end ====
