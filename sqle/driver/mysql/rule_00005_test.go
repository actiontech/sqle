package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00005(t *testing.T) {
	ruleName := ai.SQLE00005
	rule := rulepkg.RuleHandlerMap[ruleName].Rule
	ruleParams := []interface{}{3}

	//create table, no index
	runSingleRuleInspectCase(rule, t, "create table, no index", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
    );
    `, newTestResult())

	//create table, with index, no composite index
	runSingleRuleInspectCase(rule, t, "create table, with index, no composite index", DefaultMysqlInspect(), `
    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
    a int,
    b int,
    c int,
		INDEX idx_1 (a, b)
    ) ;
    `, newTestResult())

	//create table, with index,  with composite index (fewer than expected number of columns)
	runSingleRuleInspectCase(rule, t, "create table, with index, no composite index", DefaultMysqlInspect(), `
		CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
		id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
		a int,
		b int,
		c int,
		INDEX idx_1 (a, b, c)
		) ;
		`, newTestResult())

	//create table, with index, with composite index (more than expected number of columns)
	runSingleRuleInspectCase(rule, t, "create table, with index, with composite index (more than expected number of columns)", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	a int,
	b int,
	c int,
	d int,
	INDEX idx_1 (a, b, c, d)
	) ;
	`, newTestResult().addResult(ruleName, ruleParams...))

	//create index, no composite index
	runSingleRuleInspectCase(rule, t, "create table, with index, no composite index", DefaultMysqlInspect(), `
	CREATE INDEX idx_1 on exist_db.exist_tb_3(v1, v2);
    `, newTestResult())

	//create index, with composite index (fewer than expected number of columns)
	runSingleRuleInspectCase(rule, t, "create table, with index, with composite index (fewer than expected number of columns)", DefaultMysqlInspect(), `
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v1, v2, v3);
    `, newTestResult())

	//create index, with composite index (more than expected number of columns)
	runSingleRuleInspectCase(rule, t, "create table, with index, with composite index (more than expected number of columns)", DefaultMysqlInspect(), `
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v1, v2, v3, id);
    `, newTestResult().addResult(ruleName, ruleParams...))

	//alter table, no index
	runSingleRuleInspectCase(rule, t, "alter table, no index", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_3 ADD COLUMN a int NOT NULL COMMENT "unit test";
    `, newTestResult())

	//alter table, with index, no composite index
	runSingleRuleInspectCase(rule, t, "alter table, with index, no composite index", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX (v1, v2);
    `, newTestResult())

	//alter table, with index, with composite index (fewer than expected number of columns)
	runSingleRuleInspectCase(rule, t, "alter table, with index, with composite index (fewer than expected number of columns)", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX (v1, v2, v3);
    `, newTestResult())

	//alter table, with index, with composite index (more than expected number of columns)
	runSingleRuleInspectCase(rule, t, "alter table, with index, with composite index (more than expected number of columns)", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX (v1, v2, v3, id);
    `, newTestResult().addResult(ruleName, ruleParams...))
}

// ==== Rule test code end ====
