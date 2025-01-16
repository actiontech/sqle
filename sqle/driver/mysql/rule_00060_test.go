package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQL00060(t *testing.T) {
	ruleName := ai.SQLE00060
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE without COMMENT", "CREATE TABLE exist_db.no_exist_tb_1 (id INT);",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: CREATE TABLE with COMMENT", "CREATE TABLE exist_db.no_exist_tb_1 (id INT) COMMENT='This is a test table';",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: CREATE TABLE with COMMENT and multiple columns", "CREATE TABLE exist_db.no_exist_tb_1 (id INT, name VARCHAR(50)) COMMENT='Table with multiple columns';",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: CREATE TABLE with COMMENT and constraints", "CREATE TABLE exist_db.no_exist_tb_1 (id INT PRIMARY KEY, name VARCHAR(50) NOT NULL) COMMENT='Table with constraints';",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: CREATE TABLE without COMMENT but with constraints", "CREATE TABLE exist_db.no_exist_tb_1 (id INT PRIMARY KEY, name VARCHAR(50) NOT NULL);",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: CREATE TABLE with COMMENT and foreign key", "CREATE TABLE exist_db.no_exist_tb_1 (id INT PRIMARY KEY, ref_id INT, CONSTRAINT fk_ref FOREIGN KEY (ref_id) REFERENCES other_table(id)) COMMENT='Table with foreign key';",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: CREATE TABLE without COMMENT but with foreign key", "CREATE TABLE exist_db.no_exist_tb_1 (id INT PRIMARY KEY, ref_id INT, CONSTRAINT fk_ref FOREIGN KEY (ref_id) REFERENCES other_table(id));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: CREATE TABLE with COMMENT and unique constraint", "CREATE TABLE exist_db.no_exist_tb_1 (id INT PRIMARY KEY, name VARCHAR(50) UNIQUE) COMMENT='Table with unique constraint';",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: CREATE TABLE without COMMENT but with unique constraint", "CREATE TABLE exist_db.no_exist_tb_1 (id INT PRIMARY KEY, name VARCHAR(50) UNIQUE);",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: CREATE TABLE with COMMENT and check constraint", "CREATE TABLE exist_db.no_exist_tb_1 (id INT PRIMARY KEY, age INT CHECK (age > 0)) COMMENT='Table with check constraint';",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 11: CREATE TABLE without COMMENT but with check constraint", "CREATE TABLE exist_db.no_exist_tb_1 (id INT PRIMARY KEY, age INT CHECK (age > 0));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: CREATE TABLE with COMMENT and IF NOT EXISTS", "CREATE TABLE IF NOT EXISTS exist_db.no_exist_tb_1 (id BIGINT, name VARCHAR(32) DEFAULT '', sex SMALLINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id)) COMMENT='客户';",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 13: CREATE TABLE without COMMENT and IF NOT EXISTS", "CREATE TABLE IF NOT EXISTS exist_db.no_exist_tb_1 (id BIGINT, name VARCHAR(32) DEFAULT '', sex SMALLINT NOT NULL, city VARCHAR(32) NOT NULL, age SMALLINT NOT NULL, PRIMARY KEY (id));",
		nil, nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
