package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00131(t *testing.T) {
	ruleName := ai.SQLE00131
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//select...order by...
	runSingleRuleInspectCase(rule, t, "select...order by...", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 ORDER BY v1;
	`, newTestResult())

	//select...order by rand()...
	runSingleRuleInspectCase(rule, t, "select...order by rand()...", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 ORDER BY RAND();
	`, newTestResult().addResult(ruleName))

	//select...order by rand(3)...
	runSingleRuleInspectCase(rule, t, "select...order by rand(3)...", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 ORDER BY RAND(3);
	`, newTestResult().addResult(ruleName))

	//insert...select...order by rand()...
	runSingleRuleInspectCase(rule, t, "insert...select...order by rand()...", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2 ORDER BY RAND();
	`, newTestResult().addResult(ruleName))

	//insert...select...order by...
	runSingleRuleInspectCase(rule, t, "insert...select...order by...", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2 ORDER BY v1;
	`, newTestResult())

	//union...order by rand()...
	runSingleRuleInspectCase(rule, t, "union...order by rand()...", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL (SELECT * FROM exist_db.exist_tb_2 ORDER BY RAND());
	`, newTestResult().addResult(ruleName))

	//union...order by...
	runSingleRuleInspectCase(rule, t, "union...order by...", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL (SELECT * FROM exist_db.exist_tb_2 ORDER BY v1);
	`, newTestResult())
}

// ==== Rule test code end ====
