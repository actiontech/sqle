package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQL00066(t *testing.T) {
	ruleName := ai.SQLE00066
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: ALTER TABLE drop COLUMN", "ALTER TABLE test_table DROP COLUMN test_column;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, test_column INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: ALTER TABLE drop INDEX", "ALTER TABLE test_table DROP INDEX test_index;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, test_column INT, INDEX test_index (test_column));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: ALTER TABLE drop FOREIGN KEY", "ALTER TABLE test_table DROP FOREIGN KEY fk_test;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, test_column INT, CONSTRAINT fk_test FOREIGN KEY (test_column) REFERENCES other_table (id));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: DROP TABLE ", "DROP TABLE test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, test_column INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: DROP INDEX ", "DROP INDEX test_index ON test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, test_column INT, INDEX test_index (test_column));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: DROP DATABASE ", "DROP DATABASE test_database;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE test_database;"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: ALTER TABLE drop PARTITION", "ALTER TABLE t1 DROP PARTITION p0;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE IF NOT EXISTS db_mysql; USE db_mysql; CREATE TABLE t1( id INT PRIMARY KEY, c1 INT DEFAULT 0, c2 INT DEFAULT 1, CONSTRAINT t1_check_global CHECK(c1 < c2) ENFORCED ) PARTITION BY RANGE(id) ( PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION p_max VALUES LESS THAN (MAXVALUE) );"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: ALTER TABLE drop CONSTRAINT", "ALTER TABLE t1 DROP INDEX test_index, DROP CONSTRAINT t1_check_global;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE IF NOT EXISTS db_mysql; USE db_mysql; CREATE TABLE t1( id INT PRIMARY KEY, c1 INT DEFAULT 0, c2 INT DEFAULT 1, CONSTRAINT t1_check_global CHECK(c1 < c2) ENFORCED ) PARTITION BY RANGE(id) ( PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION p_max VALUES LESS THAN (MAXVALUE) );"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: ALTER TABLE drop DEFAULT", "ALTER TABLE t1 ALTER c1 DROP DEFAULT;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE IF NOT EXISTS db_mysql; USE db_mysql; CREATE TABLE t1( id INT PRIMARY KEY, c1 INT DEFAULT 0, c2 INT DEFAULT 1, CONSTRAINT t1_check_global CHECK(c1 < c2) ENFORCED ) PARTITION BY RANGE(id) ( PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION p_max VALUES LESS THAN (MAXVALUE) );"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: DROP VIEW ", "DROP VIEW IF NOT EXISTS exist_db.view_1;",
		session.NewAIMockContext().WithSQL("CREATE view exist_db.view_1 as select * from exist_db.exist_tb_1;"),
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: DROP FUNCTION ", "DROP FUNCTION t1_func;",
		nil,
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: DROP PROCEDURE ", "DROP PROCEDURE t1_proc;",
		nil,
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: DROP TRIGGER ", "DROP TRIGGER t1_trigger;",
		nil,
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: DROP EVENT ", "DROP EVENT t1_event;",
		nil,
		nil, newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))

}

// ==== Rule test code end ====
