package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00127(t *testing.T) {
	ruleName := ai.SQLE00127
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//select..., no problem
	runSingleRuleInspectCase(rule, t, "select..., no problem", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select..., with order by, no problem
	runSingleRuleInspectCase(rule, t, "select..., with order by, no problem", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 ORDER BY id;
	`, newTestResult())

	//select..., with order by, with function, violate the rule
	runSingleRuleInspectCase(rule, t, "select..., with order by, with function, violate the rule", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 ORDER BY concat(id, 1);
	`, newTestResult().addResult(ruleName))

	//select..., with order by, with arithmetic, violate the rule
	runSingleRuleInspectCase(rule, t, "select..., with order by, with arithmetic, violate the rule", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 ORDER BY id + 1;
	`, newTestResult().addResult(ruleName))

	//union..., no problem
	runSingleRuleInspectCase(rule, t, "union..., no problem", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 UNION SELECT id FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union..., with order by, no problem
	runSingleRuleInspectCase(rule, t, "union..., with order by, no problem", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 UNION SELECT id FROM exist_db.exist_tb_2 ORDER BY id;
	`, newTestResult())

	//union..., with order by, with function, violate the rule
	runSingleRuleInspectCase(rule, t, "union..., with order by, with function, violate the rule", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 UNION (SELECT id FROM exist_db.exist_tb_2 ORDER BY concat(id, 1));
	`, newTestResult().addResult(ruleName))

	//union..., with order by, with arithmetic, violate the rule
	runSingleRuleInspectCase(rule, t, "union..., with order by, with arithmetic, violate the rule", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 UNION (SELECT id FROM exist_db.exist_tb_2 ORDER BY id + 1);
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
