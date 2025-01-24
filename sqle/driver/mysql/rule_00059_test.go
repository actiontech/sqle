package mysql

import (
	"testing"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/stretchr/testify/assert"
)

// ==== Rule test code start ====

// For rule involving online information, use NewMockExecutor to simulate sql statements.
func NewMySQLInspectOnRuleSQLE00059(t *testing.T, tableSize int /*table size MB*/) *MysqlDriverImpl {
	e, _, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)
	inspect.Ctx = session.NewMockContextForTestTableSize(e, map[string]int{
		"exist_tb_1": tableSize,
	})

	return inspect
}

func TestRuleSQLE00059(t *testing.T) {
	ruleName := ai.SQLE00059
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule
	ruleParams := []interface{}{"5"}

	// Alter table, no modify/change column
	i := NewMySQLInspectOnRuleSQLE00059(t, 6000)
	runSingleRuleInspectCase(rule, t, "alter table, no modify/change column", i, `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int NOT NULL COMMENT "unit test";
	`, newTestResult())

	// Alter table, with modify column, rows less than threshold
	i = NewMySQLInspectOnRuleSQLE00059(t, 50)
	runSingleRuleInspectCase(rule, t, "alter table, with modify column, rows less than threshold", i, `
	ALTER TABLE exist_db.exist_tb_1 MODIFY v int NOT NULL COMMENT "unit test";
	`, newTestResult())

	// Alter table, with modify column, rows greater than threshold
	i = NewMySQLInspectOnRuleSQLE00059(t, 6000)
	runSingleRuleInspectCase(rule, t, "alter table, with modify column, rows greater than threshold", i, `
	ALTER TABLE exist_db.exist_tb_1 MODIFY v int NOT NULL COMMENT "unit test";`, newTestResult().addResult(ruleName, ruleParams...))

	// Alter table, with change column, rows less than threshold
	i = NewMySQLInspectOnRuleSQLE00059(t, 50)
	runSingleRuleInspectCase(rule, t, "alter table, with change column, rows less than threshold", i, `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v int NOT NULL COMMENT "unit test";
	`, newTestResult())

	// Alter table, with change column, rows greater than threshold
	i = NewMySQLInspectOnRuleSQLE00059(t, 6000)
	runSingleRuleInspectCase(rule, t, "alter table, with change column, rows greater than threshold", i, `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v int NOT NULL COMMENT "unit test";`, newTestResult().addResult(ruleName, ruleParams...))
}

// ==== Rule test code end ====
