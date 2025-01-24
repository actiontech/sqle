package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00073(t *testing.T) {
	ruleName := ai.SQLE00073
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//alter table, normal
	runSingleRuleInspectCase(rule, t, "alter table, normal", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int NOT NULL COMMENT "unit test";
    `, newTestResult())

	//ALTER TABLE ... CONVERT TO CHARACTER SET ...
	runSingleRuleInspectCase(rule, t, "ALTER TABLE ... CONVERT TO CHARACTER SET ...", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 CONVERT TO CHARACTER SET utf8;
    `, newTestResult().addResult(ruleName))

	//ALTER TABLE ... CHARACTER SET ...
	runSingleRuleInspectCase(rule, t, "ALTER TABLE ... CHARACTER SET ...t", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 CHARACTER SET 'utf8';
    `, newTestResult().addResult(ruleName))

	//ALTER TABLE ... COLLATE ...
	runSingleRuleInspectCase(rule, t, "ALTER TABLE ... COLLATE ...", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 default collate = utf8_unicode_ci;
    `, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
