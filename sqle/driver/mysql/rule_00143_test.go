package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====

func TestRuleSQLE00143(t *testing.T) {
	ruleName := ai.SQLE00143
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select... no problem", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select... with join, no problem
	runSingleRuleInspectCase(rule, t, "select... with join, no problem", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 JOIN exist_db.exist_tb_2 ON id = id;
	`, newTestResult())

	//select... with join, problem (on)
	runSingleRuleInspectCase(rule, t, "select... with join, problem (on)", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 t1 JOIN exist_db.exist_tb_2 t2 ON t1.id = t2.id OR t1.v1 = t2.v1;
	`, newTestResult().addResult(ruleName))

	//select... with join, problem (where)
	runSingleRuleInspectCase(rule, t, "select... with join, problem (where)", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 t1, exist_db.exist_tb_2 t2 WHERE t1.id = t2.id OR t1.v1 = t2.v1;
	`, newTestResult().addResult(ruleName))

	//insert... select... no problem
	runSingleRuleInspectCase(rule, t, "insert... select... no problem", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT id FROM exist_db.exist_tb_2;
	`, newTestResult())

	//insert... select... with join, no problem
	runSingleRuleInspectCase(rule, t, "insert... select... with join, no problem", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT id FROM exist_db.exist_tb_2 WHERE id = 1 OR id2 = 2;
	`, newTestResult())

	// insert... select... with join, problem
	runSingleRuleInspectCase(rule, t, "insert... select... with join, problem", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_1 t1 JOIN exist_db.exist_tb_2 t2 ON t1.id = t2.id OR t1.v1 = t2.v1;
	`, newTestResult().addResult(ruleName))

	//update... with join, no problem
	runSingleRuleInspectCase(rule, t, "update... with join, no problem", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET id = (SELECT id FROM exist_db.exist_tb_2 WHERE id = 1);
	`, newTestResult())

	//update... with join, problem
	runSingleRuleInspectCase(rule, t, "update... with join, problem", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET id = 1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 t1 JOIN exist_db.exist_tb_2 t2 ON t1.id = t2.id OR t1.v1 = t2.v1);
	`, newTestResult().addResult(ruleName))

	//delete... with join, no problem
	runSingleRuleInspectCase(rule, t, "delete... with join, no problem", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult())

	//delete... with join, problem
	runSingleRuleInspectCase(rule, t, "delete... with join, problem", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT * FROM exist_db.exist_tb_1 t1 JOIN exist_db.exist_tb_2 t2 ON t1.id = t2.id OR t1.v1 = t2.v1);
	`, newTestResult().addResult(ruleName))

	//union... with join, no problem
	runSingleRuleInspectCase(rule, t, "union... with join, no problem", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 UNION ALL SELECT id FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union... with join, problem
	runSingleRuleInspectCase(rule, t, "union... with join, problem", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 WHERE id = 1 UNION ALL (SELECT * FROM exist_db.exist_tb_1 t1 JOIN exist_db.exist_tb_2 t2 ON t1.id = t2.id OR t1.v1 = t2.v1);
	`, newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
