package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00041(t *testing.T) {
	ruleName := ai.SQLE00041
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// case 7: CREATE TABLE with UNIQUE KEY with correct prefix (从xml中补充)
	runAIRuleCase(rule, t, "case 7_1: CREATE TABLE with UNIQUE KEY with correct prefix (从xml中补充)",
		"CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', UNIQUE KEY uniq_name(name));",
		nil,
		nil,
		newTestResult())

	// case 1: CREATE TABLE with UNIQUE INDEX without prefix
	runAIRuleCase(rule, t, "case 1: CREATE TABLE with UNIQUE INDEX without prefix",
		"CREATE TABLE test_table (id INT, UNIQUE INDEX idx_id (id));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	// case 2: CREATE TABLE with UNIQUE INDEX with correct prefix
	runAIRuleCase(rule, t, "case 2: CREATE TABLE with UNIQUE INDEX with correct prefix",
		"CREATE TABLE test_table (id INT, UNIQUE INDEX uniq_idx_id (id));",
		nil,
		nil,
		newTestResult())

	// case 3: ALTER TABLE ADD UNIQUE INDEX without prefix
	runAIRuleCase(rule, t, "case 3: ALTER TABLE ADD UNIQUE INDEX without prefix",
		"ALTER TABLE test_table ADD UNIQUE INDEX idx_id (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil,
		newTestResult().addResult(ruleName))

	// case 4: ALTER TABLE ADD UNIQUE INDEX with correct prefix
	runAIRuleCase(rule, t, "case 4: ALTER TABLE ADD UNIQUE INDEX with correct prefix",
		"ALTER TABLE test_table ADD UNIQUE INDEX uniq_idx_id (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil,
		newTestResult())

	// case 5: ALTER TABLE RENAME INDEX to UNIQUE INDEX without prefix
	runAIRuleCase(rule, t, "case 5: ALTER TABLE RENAME INDEX to UNIQUE INDEX without prefix",
		"ALTER TABLE test_table RENAME INDEX old_idx TO idx_id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, UNIQUE INDEX old_idx (id));"),
		nil,
		newTestResult().addResult(ruleName))

	// case 6: ALTER TABLE RENAME INDEX to UNIQUE INDEX with correct prefix
	runAIRuleCase(rule, t, "case 6: ALTER TABLE RENAME INDEX to UNIQUE INDEX with correct prefix",
		"ALTER TABLE test_table RENAME INDEX old_idx TO uniq_idx_id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, UNIQUE INDEX old_idx (id));"),
		nil,
		newTestResult())

	// case 7: CREATE TABLE with UNIQUE INDEX with correct prefix (从xml中补充)
	runAIRuleCase(rule, t, "case 7: CREATE TABLE with UNIQUE INDEX with correct prefix (从xml中补充)",
		"CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', UNIQUE INDEX uniq_name(name));",
		nil,
		nil,
		newTestResult())

	// case 7: CREATE TABLE with UNIQUE KEY with correct prefix (从xml中补充)
	runAIRuleCase(rule, t, "case 7_1: CREATE TABLE with UNIQUE KEY with correct prefix (从xml中补充)",
		"CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', UNIQUE KEY uniq_name(name));",
		nil,
		nil,
		newTestResult())

	// case 8: CREATE UNIQUE INDEX on existing table with correct prefix (从xml中补充)
	runAIRuleCase(rule, t, "case 8: CREATE UNIQUE INDEX on existing table with correct prefix (从xml中补充)",
		"CREATE UNIQUE INDEX uniq_id_idx ON order_his (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '');"),
		nil,
		newTestResult())

	// case 9: ALTER TABLE ADD UNIQUE INDEX with correct prefix (从xml中补充)
	runAIRuleCase(rule, t, "case 9: ALTER TABLE ADD UNIQUE INDEX with correct prefix (从xml中补充)",
		"ALTER TABLE order_his ADD UNIQUE INDEX uniq_id(id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '');"),
		nil,
		newTestResult())

	// case 10: ALTER TABLE RENAME INDEX with correct prefix (从xml中补充)
	runAIRuleCase(rule, t, "case 10: ALTER TABLE RENAME INDEX with correct prefix (从xml中补充)",
		"ALTER TABLE order_his RENAME INDEX name TO uniq_name2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', UNIQUE INDEX name(name));"),
		nil,
		newTestResult())
}

// ==== Rule test code end ====
