package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====

func TestRuleSQLE00119(t *testing.T) {
	ruleName := ai.SQLE00119
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select... no group by, no order by", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select... with group by, no order by
	runSingleRuleInspectCase(rule, t, "select... with group by, no order by", DefaultMysqlInspect(), `
	SELECT id, COUNT(id) FROM exist_db.exist_tb_1 GROUP BY id;
	`, newTestResult().addResult(ruleName))

	//select... with group by, with order by
	runSingleRuleInspectCase(rule, t, "select... with group by, with order by", DefaultMysqlInspect(), `
	SELECT id, COUNT(id) FROM exist_db.exist_tb_1 GROUP BY id ORDER BY id;
	`, newTestResult())

	//insert... select... no group by, no order by
	runSingleRuleInspectCase(rule, t, "insert... select... no group by, no order by", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_2 SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult())

	//insert... select... with group by, no order by
	runSingleRuleInspectCase(rule, t, "insert... select... with group by, no order by", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_2 SELECT id, COUNT(id) FROM exist_db.exist_tb_1 GROUP BY id;
	`, newTestResult().addResult(ruleName))

	//insert... select... with group by, with order by
	runSingleRuleInspectCase(rule, t, "insert... select... with group by, with order by", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_2 SELECT id, COUNT(id) FROM exist_db.exist_tb_1 GROUP BY id ORDER BY id;
	`, newTestResult())

	//union... no group by, no order by
	runSingleRuleInspectCase(rule, t, "union... no group by, no order by", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union... with group by, no order by
	runSingleRuleInspectCase(rule, t, "union... with group by, no order by", DefaultMysqlInspect(), `
	SELECT id, COUNT(id) FROM exist_db.exist_tb_1 UNION ALL SELECT id, COUNT(id) FROM exist_db.exist_tb_2 GROUP BY id;
	`, newTestResult().addResult(ruleName))

	//union... with group by, with order by
	runSingleRuleInspectCase(rule, t, "union... with group by, with order by", DefaultMysqlInspect(), `
	SELECT id, COUNT(id) FROM exist_db.exist_tb_1 UNION ALL (SELECT id, COUNT(id) FROM exist_db.exist_tb_2 GROUP BY id ORDER BY id);
	`, newTestResult())

	//update... no group by, no order by
	runSingleRuleInspectCase(rule, t, "update... no group by, no order by", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 1;
	`, newTestResult())

	//update... with group by, no order by
	runSingleRuleInspectCase(rule, t, "update... with group by, no order by", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = (SELECT id FROM exist_db.exist_tb_2 GROUP BY id);
	`, newTestResult().addResult(ruleName))

	//update... with group by, with order by
	runSingleRuleInspectCase(rule, t, "update... with group by, with order by", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = (SELECT id FROM exist_db.exist_tb_2 GROUP BY id ORDER BY id);
	`, newTestResult())

	//delete... no group by, no order by
	runSingleRuleInspectCase(rule, t, "delete... no group by, no order by", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult())

	//delete... with group by, no order by
	runSingleRuleInspectCase(rule, t, "delete... with group by, no order by", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = (SELECT id FROM exist_db.exist_tb_2 GROUP BY id);
	`, newTestResult().addResult(ruleName))

	//delete... with group by, with order by
	runSingleRuleInspectCase(rule, t, "delete... with group by, with order by", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = (SELECT id FROM exist_db.exist_tb_2 GROUP BY id ORDER BY id);
	`, newTestResult())

}

// ==== Rule test code end ====
