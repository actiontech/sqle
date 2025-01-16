package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00099(t *testing.T) {
	ruleName := ai.SQLE00099
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select...without FOR UPDATE", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select...with FOR UPDATE
	runSingleRuleInspectCase(rule, t, "select...with FOR UPDATE", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 FOR UPDATE;
	`, newTestResult().addResult(ruleName))

	//union...
	runSingleRuleInspectCase(rule, t, "union...without FOR UPDATE", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union...with FOR UPDATE
	runSingleRuleInspectCase(rule, t, "union...with FOR UPDATE", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 FOR UPDATE UNION ALL SELECT * FROM exist_db.exist_tb_2 FOR UPDATE;
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
