package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00126(t *testing.T) {
	ruleName := ai.SQLE00126
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select..., no group by", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select... with group by position
	runSingleRuleInspectCase(rule, t, "select..., with group by position", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 GROUP BY 1;
	`, newTestResult().addResult(ruleName))

	//select... with group by, with group by option
	runSingleRuleInspectCase(rule, t, "select..., with group by, with group by option", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 GROUP BY id, v1;
	`, newTestResult())

	//union...
	runSingleRuleInspectCase(rule, t, "union..., no group by", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union... with group by position
	runSingleRuleInspectCase(rule, t, "union..., with group by position", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2 GROUP BY 1;
	`, newTestResult().addResult(ruleName))

	//union... with group by, with group by option
	runSingleRuleInspectCase(rule, t, "union..., with group by, with group by option", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2 GROUP BY id, v1;
	`, newTestResult())

	//insert... select..., no group by
	runSingleRuleInspectCase(rule, t, "insert... select..., no group by", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_3 SELECT id, v1 FROM exist_db.exist_tb_1;
	`, newTestResult())

	//insert... select..., with group by position
	runSingleRuleInspectCase(rule, t, "insert... select..., with group by position", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_3 SELECT id, v1 FROM exist_db.exist_tb_1 GROUP BY 1;
	`, newTestResult().addResult(ruleName))

	//insert... select..., with group by, with group by option
	runSingleRuleInspectCase(rule, t, "insert... select..., with group by, with group by option", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_3 SELECT id, v1 FROM exist_db.exist_tb_1 GROUP BY id, v1;
	`, newTestResult())

	//update... select..., no group by
	runSingleRuleInspectCase(rule, t, "update... select..., no group by", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_3 SET v1 = (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult())

	//update... select..., with group by position
	runSingleRuleInspectCase(rule, t, "update... select..., with group by position", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_3 SET v1 = (SELECT id FROM exist_db.exist_tb_1 GROUP BY 1);
	`, newTestResult().addResult(ruleName))

	//update... select..., with group by, with group by option
	runSingleRuleInspectCase(rule, t, "update... select..., with group by, with group by option", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_3 SET v1 = (SELECT id FROM exist_db.exist_tb_1 GROUP BY id, v1);
	`, newTestResult())

	//delete... select..., no group by
	runSingleRuleInspectCase(rule, t, "delete... select..., no group by", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_3 WHERE EXISTS (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult())

	//delete... select..., with group by position
	runSingleRuleInspectCase(rule, t, "delete... select..., with group by position", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_3 WHERE EXISTS (SELECT id FROM exist_db.exist_tb_1 GROUP BY 1);
	`, newTestResult().addResult(ruleName))

	//delete... select..., with group by, with group by option
	runSingleRuleInspectCase(rule, t, "delete... select..., with group by, with group by option", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_3 WHERE EXISTS (SELECT id FROM exist_db.exist_tb_1 GROUP BY id, v1);
	`, newTestResult())
}

// ==== Rule test code end ====
