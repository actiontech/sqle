package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00071(t *testing.T) {
	ruleName := ai.SQLE00071
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//alter table add column, no problem
	runSingleRuleInspectCase(rule, t, "alter table add column, no problem", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int NOT NULL COMMENT "unit test";
    `, newTestResult())

	//alter table drop column, violate the rule
	runSingleRuleInspectCase(rule, t, "alter table drop column, violate the rule", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 DROP COLUMN v1;
    `, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
