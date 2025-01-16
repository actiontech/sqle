package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00016_1(t *testing.T) {
	ruleName := ai.SQLE00016_1
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, with blob/text column, no NOT NULL
	runSingleRuleInspectCase(rule, t, "create table, with blob/text column, no NOT NULL", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
a blob,
b text,
PRIMARY KEY (id)
);
`, newTestResult())

	// create table, with blob/text column, with NOT NULL
	runSingleRuleInspectCase(rule, t, "create table, with expected column type, with NOT NULL", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
a blob NOT NULL,
b text NOT NULL,
PRIMARY KEY (id)
);
`, newTestResult().addResult(ruleName, "a,b"))

	// alter table add columns, with blob/text column, no NOT NULL
	runSingleRuleInspectCase(rule, t, "alter table add columns, with blob/text column, no NOT NULL", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a blob;
`, newTestResult())

	// alter table add columns, with blob/text column, with NOT NULL
	runSingleRuleInspectCase(rule, t, "alter table add columns, with expected column type, with NOT NULL", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a blob NOT NULL;
`, newTestResult().addResult(ruleName, "a"))

	// alter table modify column, with blob/text column, no NOT NULL
	runSingleRuleInspectCase(rule, t, "alter table modify column, with blob/text column, no NOT NULL", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 blob;
`, newTestResult())

	// alter table modify column, with blob/text column, with NOT NULL
	runSingleRuleInspectCase(rule, t, "alter table modify column, with expected column type, with NOT NULL", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 blob NOT NULL;
`, newTestResult().addResult(ruleName, "v1"))

	// alter table change column, with blob/text column, no NOT NULL
	runSingleRuleInspectCase(rule, t, "alter table change column, with blob/text column, no NOT NULL", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a blob;
`, newTestResult())

	// alter table change column, with blob/text column, with NOT NULL
	runSingleRuleInspectCase(rule, t, "alter table change column, with expected column type, with NOT NULL", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a blob NOT NULL;
`, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
