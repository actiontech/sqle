package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00028(t *testing.T) {
	ruleName := ai.SQLE00028
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, all column has DEFAULT value
	runSingleRuleInspectCase(rule, t, "create table, auto increment column has no DEFAULT value", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 1,
	a int DEFAULT 1
	);
	`, newTestResult())

	//create table, column has no DEFAULT value
	runSingleRuleInspectCase(rule, t, "create table, auto increment column has no DEFAULT value", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 1,
	a int
	);
	`, newTestResult().addResult(ruleName, "a"))

	//create table, auto increment column has no DEFAULT value
	runSingleRuleInspectCase(rule, t, "create table, auto increment column has no DEFAULT value", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 1,
	a int AUTO_INCREMENT
	);
	`, newTestResult())

	//create table, blob/text column has no DEFAULT value
	runSingleRuleInspectCase(rule, t, "create table, auto increment column has no DEFAULT value", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 1,
	a blob AUTO_INCREMENT
	);
	`, newTestResult())

	//alter table add column, column has no DEFAULT value
	runSingleRuleInspectCase(rule, t, "alter table add column, no column has DEFAULT value", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int NOT NULL COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "v"))

	//alter table add column, column has DEFAULT value
	runSingleRuleInspectCase(rule, t, "alter table add column, column has DEFAULT value", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int DEFAULT 100 COMMENT "unit test";
	`, newTestResult())
}

// ==== Rule test code end ====
