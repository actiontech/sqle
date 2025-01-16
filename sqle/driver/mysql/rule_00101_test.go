package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00101(t *testing.T) {
	ruleName := ai.SQLE00101
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// select...
	runSingleRuleInspectCase(rule, t, "select... without ORDER BY", DefaultMysqlInspect(), `
SELECT id FROM exist_db.exist_tb_1;
`, newTestResult())

	// select... with ORDER BY
	runSingleRuleInspectCase(rule, t, "select... with ORDER BY", DefaultMysqlInspect(), `
SELECT id FROM exist_db.exist_tb_1 ORDER BY id;
`, newTestResult().addResult(ruleName))

	// union...
	runSingleRuleInspectCase(rule, t, "union... without ORDER BY", DefaultMysqlInspect(), `
SELECT id FROM exist_db.exist_tb_1
UNION ALL
SELECT id FROM exist_db.exist_tb_2;
`, newTestResult())

	// union... with ORDER BY
	runSingleRuleInspectCase(rule, t, "union... with ORDER BY", DefaultMysqlInspect(), `
SELECT id FROM exist_db.exist_tb_1
UNION ALL
(SELECT id FROM exist_db.exist_tb_2 ORDER BY id);
`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
