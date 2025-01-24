package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====

func TestRuleSQLE00115(t *testing.T) {
	ruleName := ai.SQLE00115
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select..., no subquery", DefaultMysqlInspect(), `
	SELECT id, v1, v2 FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select..., with scalar subquery in select clause
	runSingleRuleInspectCase(rule, t, "select..., with scalar subquery in select clause", DefaultMysqlInspect(), `
	SELECT id, (SELECT id FROM exist_db.exist_tb_1), v2 FROM exist_db.exist_tb_1;
	`, newTestResult().addResult(ruleName))

	//select..., with scalar subquery in where clause
	runSingleRuleInspectCase(rule, t, "select..., with scalar subquery in where clause", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 WHERE id = (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult().addResult(ruleName))

	//select..., with in subquery in select clause
	runSingleRuleInspectCase(rule, t, "select..., with in subquery in select clause", DefaultMysqlInspect(), `
	SELECT id, v1, v2 FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult())

	//select..., with exists subquery in select clause
	runSingleRuleInspectCase(rule, t, "select..., with exists subquery in select clause", DefaultMysqlInspect(), `
	SELECT id, v1, v2 FROM exist_db.exist_tb_1 WHERE EXISTS (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult())

	//union..., no subquery
	runSingleRuleInspectCase(rule, t, "union..., no subquery", DefaultMysqlInspect(), `
	SELECT id, v1, v2 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1, v2 FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union..., with scalar subquery in select clause
	runSingleRuleInspectCase(rule, t, "union..., with scalar subquery in select clause", DefaultMysqlInspect(), `
	SELECT id, (SELECT id FROM exist_db.exist_tb_1), v2 FROM exist_db.exist_tb_1 UNION ALL SELECT id, (SELECT id FROM exist_db.exist_tb_2), v2 FROM exist_db.exist_tb_2;
	`, newTestResult().addResult(ruleName))

	//union..., with in subquery in select clause
	runSingleRuleInspectCase(rule, t, "union..., with in subquery in select clause", DefaultMysqlInspect(), `
	SELECT id, v1, v2 FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1) UNION ALL SELECT id, v1, v2 FROM exist_db.exist_tb_2 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult())

	//union..., with exists subquery in select clause
	runSingleRuleInspectCase(rule, t, "union..., with exists subquery in select clause", DefaultMysqlInspect(), `
	SELECT id, v1, v2 FROM exist_db.exist_tb_1 WHERE EXISTS (SELECT id FROM exist_db.exist_tb_1) UNION ALL SELECT id, v1, v2 FROM exist_db.exist_tb_2 WHERE EXISTS (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult())

	//union..., with scalar subquery in where clause
	runSingleRuleInspectCase(rule, t, "union..., with scalar subquery in where clause", DefaultMysqlInspect(), `
	SELECT id FROM exist_db.exist_tb_1 WHERE id = (SELECT id FROM exist_db.exist_tb_1) UNION ALL SELECT id FROM exist_db.exist_tb_2 WHERE id = (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))

	//delete..., no subquery
	runSingleRuleInspectCase(rule, t, "delete..., no subquery", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = 1;
	`, newTestResult())

	//delete..., with scalar subquery in where clause
	runSingleRuleInspectCase(rule, t, "delete..., with scalar subquery in where clause", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult().addResult(ruleName))

	//delete..., with in subquery in where clause
	runSingleRuleInspectCase(rule, t, "delete..., with in subquery in where clause", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult())

	//delete..., with exists subquery in where clause
	runSingleRuleInspectCase(rule, t, "delete..., with exists subquery in where clause", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE EXISTS (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult())

	//update..., no subquery
	runSingleRuleInspectCase(rule, t, "update..., no subquery", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 1, v2 = 2 WHERE id = 1;
	`, newTestResult())

	//update..., with scalar subquery in SET clause
	runSingleRuleInspectCase(rule, t, "update..., with scalar subquery in SET clause", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = (SELECT id FROM exist_db.exist_tb_1), v2 = 2 WHERE id = 1;
	`, newTestResult().addResult(ruleName))

	//update..., with in subquery
	runSingleRuleInspectCase(rule, t, "update..., with in subquery", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 1, v2 = 2 WHERE id IN (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult())

	//update..., with exists subquery
	runSingleRuleInspectCase(rule, t, "update..., with exists subquery", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 1, v2 = 2 WHERE EXISTS (SELECT id FROM exist_db.exist_tb_1);
	`, newTestResult())

	//insert..., no subquery
	runSingleRuleInspectCase(rule, t, "insert..., no subquery", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 (id, v1, v2) VALUES (1, 2, 3);
	`, newTestResult())

	//insert..., with scalar subquery in VALUES clause
	runSingleRuleInspectCase(rule, t, "insert..., with scalar subquery in VALUES clause", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 (id, v1, v2) VALUES ((SELECT id FROM exist_db.exist_tb_1), 2, 3);
	`, newTestResult().addResult(ruleName))

	//insert..., with exists subquery in SET clause
	runSingleRuleInspectCase(rule, t, "insert..., with exists subquery in VALUES clause", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SET id = (SELECT id FROM exist_db.exist_tb_1), v1 = 2, v2 = 3;
	`, newTestResult().addResult(ruleName))

	//insert..., with scalar subquery in column definition
	runSingleRuleInspectCase(rule, t, "insert..., with scalar subquery in column definition", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 (id, v1, v2) VALUES (1, 2, (SELECT id FROM exist_db.exist_tb_1));
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
