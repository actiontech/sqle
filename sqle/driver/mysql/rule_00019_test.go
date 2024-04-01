package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00019(t *testing.T) {
	ruleName := ai.SQLE00019
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, no ENUM column
	runSingleRuleInspectCase(rule, t, "create table, no ENUM column", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
    );
    `, newTestResult())

	//create table, with ENUM column
	runSingleRuleInspectCase(rule, t, "create table, with ENUM column", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    a enum('reading', 'writing', 'painting', 'cooking')
    );
    `, newTestResult().addResult(ruleName, "a"))

	//alter table add column, no ENUM column
	runSingleRuleInspectCase(rule, t, "alter table add column, no ENUM column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a int NOT NULL COMMENT "unit test";
    `, newTestResult())

	//alter table add column, with ENUM column
	runSingleRuleInspectCase(rule, t, "alter table add column, with ENUM column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a enum('reading', 'writing', 'painting', 'cooking') NOT NULL COMMENT "unit test";
    `, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, no ENUM column
	runSingleRuleInspectCase(rule, t, "alter table modify column, no ENUM column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 Modify COLUMN v1 int NOT NULL COMMENT "unit test";
    `, newTestResult())

	//alter table modify column, with ENUM column
	runSingleRuleInspectCase(rule, t, "alter table modify column, with ENUM column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 Modify COLUMN v1 enum('reading', 'writing', 'painting', 'cooking') NOT NULL COMMENT "unit test";
    `, newTestResult().addResult(ruleName, "v1"))

	//alter table change column, no ENUM column
	runSingleRuleInspectCase(rule, t, "alter table modify column, no ENUM column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v int NOT NULL COMMENT "unit test";
    `, newTestResult())

	//alter table modify column, with ENUM column
	runSingleRuleInspectCase(rule, t, "alter table modify column, with ENUM column", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v enum('reading', 'writing', 'painting', 'cooking') NOT NULL COMMENT "unit test";
    `, newTestResult().addResult(ruleName, "v"))
}

// ==== Rule test code end ====
