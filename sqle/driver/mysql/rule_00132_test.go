package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00132(t *testing.T) {
	ruleName := ai.SQLE00132
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// select no subquery...
	runSingleRuleInspectCase(rule, t, "select..., no subquery", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1;
	`, newTestResult())

	// select with subquery...
	runSingleRuleInspectCase(rule, t, "select..., with subquery", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))

	// union no subquery...
	runSingleRuleInspectCase(rule, t, "union..., no subquery", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 UNION SELECT id FROM exist_db.exist_tb_2;
	`, newTestResult())

	// union with subquery...
	runSingleRuleInspectCase(rule, t, "union..., with subquery", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_2) UNION SELECT id FROM exist_db.exist_tb_3;
	`, newTestResult().addResult(ruleName))

	// update no subquery...
	runSingleRuleInspectCase(rule, t, "update..., no subquery", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET id = 2 WHERE id = 1;
	`, newTestResult())

	//update with subquery...
	runSingleRuleInspectCase(rule, t, "update..., with subquery", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET id = (SELECT id FROM exist_db.exist_tb_2) WHERE id = 1;
	`, newTestResult().addResult(ruleName))

	//insert no subquery...
	runSingleRuleInspectCase(rule, t, "insert..., no subquery", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 (id) VALUES (3);
	`, newTestResult())

	//insert with subquery...
	runSingleRuleInspectCase(rule, t, "insert..., with subquery", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT id FROM exist_db.exist_tb_2 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))

	//delete no subquery...
	runSingleRuleInspectCase(rule, t, "delete..., no subquery", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = 1;
	`, newTestResult())

	//delete with subquery...
	runSingleRuleInspectCase(rule, t, "delete..., with subquery", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
