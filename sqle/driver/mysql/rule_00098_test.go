package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====

func TestRuleSQLE00098(t *testing.T) {
	ruleName := ai.SQLE00098
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//select, no join
	runSingleRuleInspectCase(rule, t, "select, no join", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select, with join, no repeat
	runSingleRuleInspectCase(rule, t, "select, with join, no repeat", DefaultMysqlInspect(), `
	SELECT t1.id, t2.id FROM exist_db.exist_tb_1 t1, exist_db.exist_tb_2 t2 WHERE t1.id = t2.id;
	`, newTestResult())

	//select, with join, with repeat
	runSingleRuleInspectCase(rule, t, "select, with join, with repeat", DefaultMysqlInspect(), `
	SELECT t1.id, t1.id FROM exist_db.exist_tb_1 t1, exist_db.exist_tb_1 t2, exist_db.exist_tb_1 t3 WHERE t1.id = t2.id;
	`, newTestResult().addResult(ruleName, "exist_db.exist_tb_1"))

	//select, with from subquery
	runSingleRuleInspectCase(rule, t, "select, with from subquery", DefaultMysqlInspect(), `
	SELECT count(*) FROM (
		SELECT * FROM exist_db.exist_tb_2 WHERE v1 < "1" 
		union  
		SELECT * FROM exist_db.exist_tb_2 WHERE v1 < "1" 
		union
		SELECT * FROM exist_db.exist_tb_2 WHERE v1 < "1" 
		union
		SELECT * FROM exist_db.exist_tb_2 WHERE v1 < "1"
		) T;;
	`, newTestResult().addResult(ruleName, "exist_db.exist_tb_2"))

	//union, no join
	runSingleRuleInspectCase(rule, t, "union, no join", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION SELECT id, v1 FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union, with join, no repeat
	runSingleRuleInspectCase(rule, t, "union, with join, no repeat", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION SELECT id, v1 FROM exist_db.exist_tb_2 WHERE id = 1;
	`, newTestResult())

	//union, with join, with repeat
	runSingleRuleInspectCase(rule, t, "union, with join, with repeat", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id = 1 UNION SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id = 1;
	`, newTestResult())

	//union all, no join
	runSingleRuleInspectCase(rule, t, "union all, no join", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union all, with join, no repeat
	runSingleRuleInspectCase(rule, t, "union all, with join, no repeat", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2 WHERE id = 1;
	`, newTestResult())

	//union all, with join, with repeat
	runSingleRuleInspectCase(rule, t, "union all, with join, with repeat", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id = 1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id = 1;
	`, newTestResult())

	//insert select, no join
	runSingleRuleInspectCase(rule, t, "insert select, no join", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_3 SELECT id, v1 FROM exist_db.exist_tb_1;
	`, newTestResult())

	//insert select, with join, no repeat
	runSingleRuleInspectCase(rule, t, "insert select, with join, no repeat", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_3 SELECT id, v1 FROM exist_db.exist_tb_1 t1, exist_db.exist_tb_2 t2 WHERE t1.id = t2.id;
	`, newTestResult())

	//insert select, with join, with repeat
	runSingleRuleInspectCase(rule, t, "insert select, with join, with repeat", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_3 SELECT id, v1 FROM exist_db.exist_tb_1 t1, exist_db.exist_tb_1 t2, exist_db.exist_tb_1 t3 WHERE t1.id = t2.id;
	`, newTestResult().addResult(ruleName, "exist_db.exist_tb_1"))

	//update, no join
	runSingleRuleInspectCase(rule, t, "update, no join", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 1 WHERE id = 1;
	`, newTestResult())

	//update, with join, no repeat
	runSingleRuleInspectCase(rule, t, "update, with join, no repeat", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 t1, exist_db.exist_tb_2 t2 SET t1.v1 = t2.v1 WHERE t1.id = t2.id;
	`, newTestResult())

	//update, with join, with repeat
	runSingleRuleInspectCase(rule, t, "update, with join, with repeat", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 t1, exist_db.exist_tb_1 t2, exist_db.exist_tb_1 t3 SET t1.v1 = t2.v1 WHERE t1.id = t2.id;
	`, newTestResult())

	//delete, no join
	runSingleRuleInspectCase(rule, t, "delete, no join", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = 1;
	`, newTestResult())

	//delete, with join, no repeat
	runSingleRuleInspectCase(rule, t, "delete, with join, no repeat", DefaultMysqlInspect(), `
	DELETE t1, t2 FROM exist_db.exist_tb_1 AS t1 JOIN exist_db.exist_tb_2 AS t2 ON t1.id = t2.id;
	`, newTestResult())

	//delete, with join, with repeat
	runSingleRuleInspectCase(rule, t, "delete, with join, with repeat", DefaultMysqlInspect(), `
	DELETE t1, t2, t3 FROM exist_db.exist_tb_1 AS t1 JOIN exist_db.exist_tb_2 AS t2 JOIN exist_db.exist_tb_3 AS t3 ON t1.id = t2.id;
	`, newTestResult())
}

// ==== Rule test code end ====
