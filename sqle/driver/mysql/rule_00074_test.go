package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00074(t *testing.T) {
	ruleName := ai.SQLE00074
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// alter table, no rename
	runSingleRuleInspectCase(rule, t, "alter table, normal", DefaultMysqlInspect(), `
ALTER TABLE exist_db.exist_tb_1 add column v4 int;
`, newTestResult())

	// alter table, rename table
	runSingleRuleInspectCase(rule, t, "alter table, rename table", DefaultMysqlInspect(), `
ALTER TABLE exist_db.exist_tb_1 RENAME TO new_db.new_tb_1;
`, newTestResult().addResult(ruleName))

	// rename table
	runSingleRuleInspectCase(rule, t, "rename table", DefaultMysqlInspect(), `
RENAME TABLE exist_db.exist_tb_1 TO new_db.new_tb_1;
`, newTestResult().addResult(ruleName))

	// alter table, rename column
	runSingleRuleInspectCase(rule, t, "alter table, rename column", DefaultMysqlInspect(), `
ALTER TABLE exist_db.exist_tb_1 RENAME COLUMN v1 TO v2;
`, newTestResult().addResult(ruleName))

	// alter table, change column
	runSingleRuleInspectCase(rule, t, "alter table, change column", DefaultMysqlInspect(), `
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v int NOT NULL COMMENT "unit test";
`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
