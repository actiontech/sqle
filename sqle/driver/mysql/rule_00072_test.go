package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00072(t *testing.T) {
	ruleName := ai.SQLE00072
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//alter table, not drop foreign key
	runSingleRuleInspectCase(rule, t, "alter table, not drop foreign key", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v int NOT NULL COMMENT "unit test";
    `, newTestResult())

	//alter table, drop foreign key
	runSingleRuleInspectCase(rule, t, "alter table, drop foreign key", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 DROP FOREIGN KEY v1;
    `, newTestResult().addResult(ruleName))
}
// ==== Rule test code end ====