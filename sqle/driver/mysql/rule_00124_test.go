package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00124(t *testing.T) {
	ruleName := ai.SQLE00124
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// DELETE... without WHERE
	runSingleRuleInspectCase(rule, t, "delete without where", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1;
	`, newTestResult().addResult(ruleName))

	// DELETE... with WHERE, with where condition always true
	runSingleRuleInspectCase(rule, t, "delete with where, with where condition always true", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE 1=1;
	`, newTestResult().addResult(ruleName))

	// DELETE... with WHERE, with where condition
	runSingleRuleInspectCase(rule, t, "delete with where, with where condition", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = 1;
	`, newTestResult())
}

// ==== Rule test code end ====
