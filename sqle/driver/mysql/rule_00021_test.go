package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00021(t *testing.T) {
	ruleName := ai.SQLE00021
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, without NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "create table, without NOT NULL constraint", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a int,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, "a"))

	//create table, with NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "create table, with NOT NULL constraint", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a int NOT NULL,
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//alter table add columns, without NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "alter table add columns, without NOT NULL constraint", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a int COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table add columns, with NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "alter table add columns, with NOT NULL constraint", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a int NOT NULL COMMENT "unit test";
	`, newTestResult())

	//alter table modify column, without NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "alter table modify column, without NOT NULL constraint", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN a int COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, with NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "alter table modify column, with NOT NULL constraint", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN a int NOT NULL COMMENT "unit test";
	`, newTestResult())

	//alter table change column, without NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "alter table change column, without NOT NULL constraint", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a int COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table change column, with NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "alter table change column, with NOT NULL constraint", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a int NOT NULL COMMENT "unit test";
	`, newTestResult())
}
// ==== Rule test code end ====