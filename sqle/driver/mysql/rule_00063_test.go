package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00063(t *testing.T) {
	ruleName := ai.SQLE00063
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// case 1: CREATE语句中唯一索引名不符合格式
	runAIRuleCase(rule, t, "case 1: CREATE语句中唯一索引名不符合格式",
		"CREATE TABLE test_table (id INT, UNIQUE INDEX idx_id (id));",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	// case 2: CREATE语句中唯一索引名符合格式
	runAIRuleCase(rule, t, "case 2: CREATE语句中唯一索引名符合格式",
		"CREATE TABLE test_table (id INT, UNIQUE INDEX IDX_UK_TEST_TABLE_ID (id));",
		nil, /*mock context*/
		nil, newTestResult())

	// case 3: ALTER语句中添加唯一索引名不符合格式
	runAIRuleCase(rule, t, "case 3: ALTER语句中添加唯一索引名不符合格式",
		"ALTER TABLE test_table ADD UNIQUE INDEX idx_id (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult().addResult(ruleName))

	// case 4: ALTER语句中添加唯一索引名符合格式
	runAIRuleCase(rule, t, "case 4: ALTER语句中添加唯一索引名符合格式",
		"ALTER TABLE test_table ADD UNIQUE INDEX IDX_UK_TEST_TABLE_ID (id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult())

	// case 5: ALTER语句中重命名唯一索引名不符合格式
	runAIRuleCase(rule, t, "case 5: ALTER语句中重命名唯一索引名不符合格式",
		"ALTER TABLE test_table RENAME INDEX IDX_UK_OLD_NAME TO idx_new_name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT,UNIQUE INDEX IDX_UK_OLD_NAME(id));"),
		nil, newTestResult().addResult(ruleName))

	// case 6: ALTER语句中重命名唯一索引名符合格式
	runAIRuleCase(rule, t, "case 6: ALTER语句中重命名唯一索引名符合格式",
		"ALTER TABLE test_table RENAME INDEX IDX_UK_OLD_NAME TO IDX_UK_TEST_TABLE_ID;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT,UNIQUE INDEX IDX_UK_OLD_NAME(id));"),
		nil, newTestResult())

	// case 7: CREATE语句中唯一索引名符合格式(从xml中补充)
	runAIRuleCase(rule, t, "case 7: CREATE语句中唯一索引名符合格式(从xml中补充)",
		"CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', UNIQUE INDEX IDX_UK_ORDER_HIS_NAME(name));",
		nil, /*mock context*/
		nil, newTestResult())

	// case 8: CREATE唯一索引名符合格式(从xml中补充)
	runAIRuleCase(rule, t, "case 8: CREATE唯一索引名符合格式(从xml中补充)",
		"CREATE UNIQUE INDEX IDX_UK_ORDER_HIS_NAME ON order_his (name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '');"),
		nil, newTestResult())

	// case 9: ALTER语句中添加唯一索引名符合格式(从xml中补充)
	runAIRuleCase(rule, t, "case 9: ALTER语句中添加唯一索引名符合格式(从xml中补充)",
		"ALTER TABLE order_his ADD UNIQUE INDEX IDX_UK_ORDER_HIS_NAME(id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '');"),
		nil, newTestResult().addResult(ruleName))

	// case 10: ALTER语句中重命名唯一索引名符合格式(从xml中补充)
	runAIRuleCase(rule, t, "case 10: ALTER语句中重命名唯一索引名符合格式(从xml中补充)",
		"ALTER TABLE order_his RENAME INDEX name_idx TO IDX_UK_ORDER_HIS_NAME;",
		session.NewAIMockContext().WithSQL("CREATE TABLE order_his (id BIGINT, name varchar(64) DEFAULT '', UNIQUE INDEX name_idx(name)); "),
		nil, newTestResult())
}

// ==== Rule test code end ====
