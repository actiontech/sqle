package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00034(t *testing.T) {
	ruleName := ai.SQLE00034
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// Case 1: CREATE TABLE with NOT NULL constraint but no DEFAULT clause
	runSingleRuleInspectCase(rule, t, "Case 1", DefaultMysqlInspect(), `
	CREATE TABLE exist_db.not_exist_tb_1  (id INT NOT NULL);
	`, newTestResult().addResult(ruleName))

	// Case 2: CREATE TABLE with NOT NULL constraint and DEFAULT clause
	runSingleRuleInspectCase(rule, t, "Case 2", DefaultMysqlInspect(), `
	CREATE TABLE exist_db.not_exist_tb_1  (id INT NOT NULL DEFAULT 0);
	`, newTestResult())

	// Case 3: CREATE TABLE without NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "Case 3", DefaultMysqlInspect(), `
	CREATE TABLE exist_db.not_exist_tb_1  (id INT);
	`, newTestResult())

	// Case 4: ALTER TABLE with NOT NULL constraint but no DEFAULT clause
	runSingleRuleInspectCase(rule, t, "Case 4", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN name VARCHAR(255) NOT NULL;
    `, newTestResult().addResult(ruleName))

	// // Case 5: ALTER TABLE with NOT NULL constraint and DEFAULT clause
	runSingleRuleInspectCase(rule, t, "Case 5", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT 'unknown';
	`, newTestResult())

	// Case 6: ALTER TABLE modifying column to NOT NULL but no DEFAULT clause
	runSingleRuleInspectCase(rule, t, "Case 6", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN name VARCHAR(255) NOT NULL;
	`, newTestResult().addResult(ruleName))

	// // Case 7: ALTER TABLE modifying column to NOT NULL and DEFAULT clause
	runSingleRuleInspectCase(rule, t, "Case 7", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN name VARCHAR(255) NOT NULL DEFAULT 'unknown';
	`, newTestResult())

	// Case 8: ALTER TABLE without NOT NULL constraint
	runSingleRuleInspectCase(rule, t, "Case 8", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN age INT;
	`, newTestResult())

	// Case 9: CREATE TABLE with NOT NULL constraint but no DEFAULT clause (negative example)
	runSingleRuleInspectCase(rule, t, "Case 9", DefaultMysqlInspect(), `
	CREATE TABLE exist_db.not_exist_tb_1 (id INT NOT NULL default 0, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age SMALLINT NOT NULL, PRIMARY KEY (id));
	`, newTestResult().addResult(ruleName))

	// Case 10: CREATE TABLE with NOT NULL constraint and DEFAULT clause (positive example)
	runSingleRuleInspectCase(rule, t, "Case 10", DefaultMysqlInspect(), `
	CREATE TABLE exist_db.not_exist_tb_1 (id INT NOT NULL default 0, name VARCHAR(32) DEFAULT '', sex TINYINT NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age SMALLINT NOT NULL default 0, PRIMARY KEY (id));
	`, newTestResult())

	// Case 11: ALTER TABLE modifying column to NOT NULL but no DEFAULT clause (negative example)
	runSingleRuleInspectCase(rule, t, "Case 11", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN age INT NOT NULL;
	`, newTestResult().addResult(ruleName))

	// Case 12: ALTER TABLE modifying column to NOT NULL and DEFAULT clause (positive example)
	runSingleRuleInspectCase(rule, t, "Case 12", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN age INT NOT NULL DEFAULT 0;
	`, newTestResult())

	// Case 13: ALTER TABLE change column to NOT NULL and DEFAULT clause (positive example)
	runSingleRuleInspectCase(rule, t, "Case 13", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE id id2 INT NOT NULL DEFAULT 0;
	`, newTestResult())

	// Case 14: ALTER TABLE change column to NOT NULL but no DEFAULT clause (positive example)
	runSingleRuleInspectCase(rule, t, "Case 14", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE id id2 INT NOT NULL;
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
