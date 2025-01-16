package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====

func TestRuleSQLE00123(t *testing.T) {
	ruleName := ai.SQLE00123
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//for truncate, with problem
	runSingleRuleInspectCase(rule, t, "for truncate, no problem", DefaultMysqlInspect(), `
	TRUNCATE TABLE exist_db.exist_tb_1;
	`, newTestResult().addResult(ruleName))

	//for delete, no problem
	runSingleRuleInspectCase(rule, t, "for delete, no problem", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = 1;
	`, newTestResult())
}

// ==== Rule test code end ====
