package mysql

//import (
//	"testing"
//
//	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
//	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
//)
//
//// ==== Rule test code start ====
//func TestRuleSQLE00023(t *testing.T) {
//	ruleName := ai.SQLE00023
//	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule
//	ruleParams := []interface{}{2}
//	//create table, no primary key
//	runSingleRuleInspectCase(rule, t, "create table, no primary key", DefaultMysqlInspect(), `
//    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
//    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
//    );
//    `, newTestResult())
//
//	//create table, with primary key, fewer than expected number of columns
//	runSingleRuleInspectCase(rule, t, "create table, with primary key, fewer than expected number of columns", DefaultMysqlInspect(), `
//    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
//    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
//    id2 bigint,
//    PRIMARY KEY (id, id2)
//    );
//    `, newTestResult())
//
//	//create table, with primary key, more than expected number of columns
//	runSingleRuleInspectCase(rule, t, "create table, with primary key, fewer than expected number of columns", DefaultMysqlInspect(), `
//    CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
//    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
//    a int,
//    b int,
//    PRIMARY KEY (id, a, b)
//    );
//    `, newTestResult().addResult(ruleName, ruleParams...))
//
//	//alter table, no primary key
//	runSingleRuleInspectCase(rule, t, "alter table, no primary key", DefaultMysqlInspect(), `
//    ALTER TABLE exist_db.exist_tb_3 ADD COLUMN a int NOT NULL COMMENT "unit test";
//    `, newTestResult())
//
//	//alter table, with primary key, fewer than expected number of columns
//	runSingleRuleInspectCase(rule, t, "alter table, with primary key, fewer than expected number of columns", DefaultMysqlInspect(), `
//    ALTER TABLE exist_db.exist_tb_3 ADD CONSTRAINT PK_c PRIMARY KEY (id, v1);
//    `, newTestResult())
//
//	//alter table, with primary key, more than expected number of columns
//	runSingleRuleInspectCase(rule, t, "alter table, with primary key, more than expected number of columns", DefaultMysqlInspect(), `
//    ALTER TABLE exist_db.exist_tb_3 ADD CONSTRAINT PK_c PRIMARY KEY (id, v1, v2);
//    `, newTestResult().addResult(ruleName, ruleParams...))
//}
//
//// ==== Rule test code end ====
