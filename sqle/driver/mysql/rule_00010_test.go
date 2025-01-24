package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00010(t *testing.T) {
	ruleName := ai.SQLE00010
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//alter table, drop normal index
	runSingleRuleInspectCase(rule, t, "alter table, drop normal index", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1 DROP INDEX idx_1;
    `, newTestResult())

	//alter table, drop primary key index
	runSingleRuleInspectCase(rule, t, "alter table, drop primary key index", DefaultMysqlInspect(),
		"ALTER TABLE exist_db.exist_tb_1 DROP PRIMARY KEY", newTestResult().addResult(ruleName))

	//drop index, drop normal index
	runSingleRuleInspectCase(rule, t, "drop index, drop normal index", DefaultMysqlInspect(), `
    DROP INDEX idx_1 ON exist_tb_1;
    `, newTestResult())
}

// ==== Rule test code end ====
