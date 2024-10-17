package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00040(t *testing.T) {
	ruleName := ai.SQLE00040
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE INDEX 没有使用固定前缀", "CREATE INDEX my_index ON test_table (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (column1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 1_1: CREATE UNIQUE INDEX 不需要固定前缀", "CREATE UNIQUE INDEX my_index ON test_table (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (column1 INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE INDEX 使用固定前缀", "CREATE INDEX idx_my_index ON test_table (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (column1 INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: ALTER TABLE ADD INDEX 没有使用固定前缀", "ALTER TABLE test_table ADD INDEX my_index (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (column1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3_1: ALTER TABLE ADD FULLTEXT INDEX 不需要固定前缀", "ALTER TABLE test_table ADD FULLTEXT INDEX my_index (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (column1 INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: ALTER TABLE ADD INDEX 使用固定前缀", "ALTER TABLE test_table ADD INDEX idx_my_index (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (column1 INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: ALTER TABLE RENAME INDEX 没有使用固定前缀", "ALTER TABLE test_table RENAME INDEX old_index TO new_index;",
		session.NewAIMockContext().WithSQL("use exist_db;CREATE TABLE test_table (column1 INT, INDEX old_index(column1));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: ALTER TABLE RENAME INDEX 使用固定前缀", "ALTER TABLE test_table RENAME INDEX old_index TO idx_new_index;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (column1 INT, INDEX old_index(column1));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: CREATE TABLE WITH INDEX 使用固定前缀(从xml中补充)", "CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', INDEX idx_name(name));",
		session.NewAIMockContext(),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7_1: CREATE TABLE WITH PRIMARY KEY 不需要固定前缀(从xml中补充)", "CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', PRIMARY KEY idx_name(name));",
		session.NewAIMockContext(),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: CREATE INDEX 使用固定前缀(从xml中补充)", "CREATE INDEX idx_name_idx ON order_his (name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '');"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: ALTER TABLE ADD INDEX 使用固定前缀(从xml中补充)", "ALTER TABLE order_his ADD INDEX idx_id(id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '');"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: ALTER TABLE RENAME INDEX 使用固定前缀(从xml中补充)", "ALTER TABLE order_his RENAME INDEX order_his TO idx_name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', INDEX order_his(name)); "),
		nil, newTestResult())
}

// ==== Rule test code end ====
