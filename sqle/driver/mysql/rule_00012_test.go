package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00012(t *testing.T) {
	ruleName := ai.SQLE00012
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, with BIGINT
	runSingleRuleInspectCase(rule, t, "create table, with BIGINT", DefaultMysqlInspect(),
		`
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    a int,
    PRIMARY KEY (id)
    );
    `, newTestResult())

	//create table, with DECIMAL
	runSingleRuleInspectCase(rule, t, "create table, with DECIMAL", DefaultMysqlInspect(),
		`
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    a decimal(10, 0),
    PRIMARY KEY (id)
    );
    `, newTestResult().addResult(ruleName, "a"))

	//alter table add columns, with BIGINT
	runSingleRuleInspectCase(rule, t, "alter table add columns, with BIGINT", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a bigint COMMENT "unit test";
    `, newTestResult())

	//alter table add columns, with DECIMAL
	runSingleRuleInspectCase(rule, t, "alter table add columns, with DECIMAL", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a decimal(10, 0) COMMENT "unit test";
    `, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, with BIGINT
	runSingleRuleInspectCase(rule, t, "alter table modify column, with BIGINT", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN a bigint COMMENT "unit test";
    `, newTestResult())

	//alter table modify column, with DECIMAL
	runSingleRuleInspectCase(rule, t, "alter table modify column, with DECIMAL", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN a decimal(10, 0) COMMENT "unit test";
    `, newTestResult().addResult(ruleName, "a"))

	//alter table change column, with BIGINT
	runSingleRuleInspectCase(rule, t, "alter table change column, with BIGINT", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a bigint COMMENT "unit test";
    `, newTestResult())

	//alter table change column, with DECIMAL
	runSingleRuleInspectCase(rule, t, "alter table change column, with DECIMAL", DefaultMysqlInspect(),
		`
    ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a decimal(10, 0) COMMENT "unit test";
    `, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
