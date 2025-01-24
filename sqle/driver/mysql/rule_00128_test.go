package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00128(t *testing.T) {
	ruleName := ai.SQLE00128
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select... no having", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select... with having
	runSingleRuleInspectCase(rule, t, "select... with having", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 GROUP BY id HAVING SUM(id) > 1;
	`, newTestResult().addResult(ruleName))

	//insert... select...
	runSingleRuleInspectCase(rule, t, "insert... select... no having", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2;
	`, newTestResult())

	//insert... select... with having
	runSingleRuleInspectCase(rule, t, "insert... select... with having", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2 GROUP BY id HAVING SUM(id) > 1;
	`, newTestResult().addResult(ruleName))

	//update...
	runSingleRuleInspectCase(rule, t, "update... no where", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 1;
	`, newTestResult())

	//update... with where
	runSingleRuleInspectCase(rule, t, "update... with where", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 1 WHERE id > 1;
	`, newTestResult())

	//update... with having
	runSingleRuleInspectCase(rule, t, "update... with having", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 1 WHERE id IN (SELECT id FROM exist_db.exist_tb_2 GROUP BY id HAVING SUM(id) > 1);
	`, newTestResult().addResult(ruleName))

	//delete...
	runSingleRuleInspectCase(rule, t, "delete... no where", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1;
	`, newTestResult())

	//delete... with where
	runSingleRuleInspectCase(rule, t, "delete... with where", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id > 1;
	`, newTestResult())

	//delete... with having
	runSingleRuleInspectCase(rule, t, "delete... with having", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_2 GROUP BY id HAVING SUM(id) > 1);
	`, newTestResult().addResult(ruleName))

	//union...
	runSingleRuleInspectCase(rule, t, "union... no having", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union... with having
	runSingleRuleInspectCase(rule, t, "union... with having", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_2 GROUP BY id HAVING SUM(id) > 1;
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
