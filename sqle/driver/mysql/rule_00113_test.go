package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====

func TestRuleSQLE00113(t *testing.T) {
	ruleName := ai.SQLE00113
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select... no problem", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select... is null
	runSingleRuleInspectCase(rule, t, "select... is null", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id IS NULL;
	`, newTestResult())

	//select... is not null
	runSingleRuleInspectCase(rule, t, "select... is not null", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id IS NOT NULL;
	`, newTestResult())

	//select... is not 1
	runSingleRuleInspectCase(rule, t, "select... is not 1", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 WHERE NOT id = 1;
	`, newTestResult().addResult(ruleName))

	//select... with negative where condition, not in
	runSingleRuleInspectCase(rule, t, "select... with negative where condition, not in", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id NOT IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))

	//select... with negative where condition, not like
	runSingleRuleInspectCase(rule, t, "select... with negative where condition, not like", DefaultMysqlInspect(), `
  SELECT id, v1 FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE 'prefix%';
  `, newTestResult().addResult(ruleName))

	//select... with negative where condition, not exists
	runSingleRuleInspectCase(rule, t, "select... with negative where condition, not exists", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 WHERE NOT EXISTS (SELECT id FROM exist_db.exist_tb_2 WHERE exist_tb_1.id = exist_tb_2.id);
	`, newTestResult().addResult(ruleName))

	//select... with negative where condition, not equal to
	runSingleRuleInspectCase(rule, t, "select... with negative where condition, not equal to", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id <> 1;
	`, newTestResult().addResult(ruleName))

	// select... with negative where condition, not between
	runSingleRuleInspectCase(rule, t, "select... with negative where condition, not between", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id NOT BETWEEN 1 AND 10;
	`, newTestResult().addResult(ruleName))

	//delete...
	runSingleRuleInspectCase(rule, t, "delete... no problem", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult())

	//delete... with negative where condition
	runSingleRuleInspectCase(rule, t, "delete... with negative where condition", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id NOT IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))

	// DELETE with negative where condition, not like
	runSingleRuleInspectCase(rule, t, "delete... with negative where condition, not like", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE 'prefix%';
	`, newTestResult().addResult(ruleName))

	// DELETE with negative where condition, not exists
	runSingleRuleInspectCase(rule, t, "delete... with negative where condition, not exists", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE NOT EXISTS (SELECT id FROM exist_db.exist_tb_2 WHERE exist_tb_1.id = exist_tb_2.id);
	`, newTestResult().addResult(ruleName))

	// DELETE with negative where condition, not equal to
	runSingleRuleInspectCase(rule, t, "delete... with negative where condition, not equal to", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id <> 1;
	`, newTestResult().addResult(ruleName))

	//insert...
	runSingleRuleInspectCase(rule, t, "insert... no problem", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT id, v1 FROM exist_db.exist_tb_2 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult())

	//insert... with negative where condition
	runSingleRuleInspectCase(rule, t, "insert... with negative where condition", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT id, v1 FROM exist_db.exist_tb_2 WHERE id NOT IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))

	// INSERT with negative where condition, not like
	runSingleRuleInspectCase(rule, t, "insert... with negative where condition, not like", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT id, v1 FROM exist_db.exist_tb_2 WHERE v1 NOT LIKE 'prefix%';
	`, newTestResult().addResult(ruleName))

	// INSERT with negative where condition, not exists
	runSingleRuleInspectCase(rule, t, "insert... with negative where condition, not exists", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT id, v1 FROM exist_db.exist_tb_2 WHERE NOT EXISTS (SELECT id FROM exist_db.exist_tb_3 WHERE exist_tb_2.id = exist_tb_3.id);
	`, newTestResult().addResult(ruleName))

	// INSERT with negative where condition, not equal to
	runSingleRuleInspectCase(rule, t, "insert... with negative where condition, not equal to", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT id, v1 FROM exist_db.exist_tb_2 WHERE id <> 1;
	`, newTestResult().addResult(ruleName))

	//update...
	runSingleRuleInspectCase(rule, t, "update... no problem", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = (SELECT v1 FROM exist_db.exist_tb_2 WHERE id IN (SELECT id FROM exist_db.exist_tb_2)) WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult())

	//update... with negative where condition
	runSingleRuleInspectCase(rule, t, "update... with negative where condition", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = (SELECT v1 FROM exist_db.exist_tb_2 WHERE id NOT IN (SELECT id FROM exist_db.exist_tb_2)) WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))

	// UPDATE with negative where condition, not like
	runSingleRuleInspectCase(rule, t, "update... with negative where condition, not like", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 'new_value' WHERE v1 NOT LIKE 'prefix%';
	`, newTestResult().addResult(ruleName))

	// UPDATE with negative where condition, not exists
	runSingleRuleInspectCase(rule, t, "update... with negative where condition, not exists", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 'new_value' WHERE NOT EXISTS (SELECT id FROM exist_db.exist_tb_2 WHERE exist_tb_1.id = exist_tb_2.id);
	`, newTestResult().addResult(ruleName))

	// UPDATE with negative where condition, not equal to
	runSingleRuleInspectCase(rule, t, "update... with negative where condition, not equal to", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = 'new_value' WHERE id <> 1;
	`, newTestResult().addResult(ruleName))

	//union...
	runSingleRuleInspectCase(rule, t, "union... no problem", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult())

	//union... with negative where condition
	runSingleRuleInspectCase(rule, t, "union... with negative where condition", DefaultMysqlInspect(), `
	SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2 WHERE id NOT IN (SELECT id FROM exist_db.exist_tb_2);
	`, newTestResult().addResult(ruleName))

	// UNION with negative where condition, not like
	runSingleRuleInspectCase(rule, t, "union... with negative where condition, not like", DefaultMysqlInspect(), `
		SELECT id, v1 FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2 WHERE v1 NOT LIKE 'prefix%';
		`, newTestResult().addResult(ruleName))

	// UNION with negative where condition, not exists
	runSingleRuleInspectCase(rule, t, "union... with negative where condition, not exists", DefaultMysqlInspect(), `
		SELECT id, v1 FROM exist_db.exist_tb_1 WHERE NOT EXISTS (SELECT id FROM exist_db.exist_tb_2 WHERE exist_tb_1.id = exist_tb_2.id) UNION ALL SELECT id, v1 FROM exist_db.exist_tb_3;
		`, newTestResult().addResult(ruleName))

	// UNION with negative where condition, not equal to
	runSingleRuleInspectCase(rule, t, "union... with negative where condition, not equal to", DefaultMysqlInspect(), `
		SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id <> 1 UNION ALL SELECT id, v1 FROM exist_db.exist_tb_2;
		`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
