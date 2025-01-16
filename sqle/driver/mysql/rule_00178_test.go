package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00178(t *testing.T) {
	ruleName := ai.SQLE00178
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//select... order by
	runSingleRuleInspectCase(rule, t, "select... order by", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//select... group by
	runSingleRuleInspectCase(rule, t, "select... group by", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 GROUP BY id;
	`, newTestResult().addResult(ruleName))

	//select... distinct
	runSingleRuleInspectCase(rule, t, "select... distinct", DefaultMysqlInspect(), `
	SELECT DISTINCT * FROM exist_db.exist_tb_1;
	`, newTestResult().addResult(ruleName))

	//select... order by, with where condition is always true(TRUE)
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition is always true(TRUE)", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE TRUE ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//select... order by, with where condition is always true(1=1)
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition is always true(1=1)", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE 1=1 ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//select... order by, with where condition is always true(0<1)
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition is always true(0<1)", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE 0<1 ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//select... order by, with where condition is always true(NOT FALSE)
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition is always true(NOT FALSE)", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE NOT FALSE ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//select... order by, with where condition is always true(TRUE OR FALSE)
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition is always true(TRUE OR FALSE)", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE TRUE OR FALSE ORDER BY id;
	`, newTestResult().addResult(ruleName))

	// select... order by, with where condition is always true(TRUE AND TRUE)
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition is always true(TRUE AND TRUE)", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE TRUE AND TRUE ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//select... order by, with where condition is always true((id = id) AND (v1 = v1))
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition is always true((id = id) AND (v1 = v1))", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE (id = id) AND (v1 = v1) ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//select... order by, with where condition is always true((id = id) OR (false))
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition is always true((id = id) OR (false))", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE (id = id) OR (false) ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//select... group by, with where condition is always true
	runSingleRuleInspectCase(rule, t, "select... group by, with where condition is always true", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE TRUE GROUP BY id;
	`, newTestResult().addResult(ruleName))

	//select... distinct, with where condition is always true
	runSingleRuleInspectCase(rule, t, "select... distinct, with where condition is always true", DefaultMysqlInspect(), `
	SELECT DISTINCT * FROM exist_db.exist_tb_1 WHERE TRUE;
	`, newTestResult().addResult(ruleName))

	//select... order by, with where condition
	runSingleRuleInspectCase(rule, t, "select... order by, with where condition", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE id > 0 ORDER BY id;
	`, newTestResult())

	//select... group by, with where condition
	runSingleRuleInspectCase(rule, t, "select... group by, with where condition", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 WHERE id > 0 GROUP BY id;
	`, newTestResult())

	//select... distinct, with where condition
	runSingleRuleInspectCase(rule, t, "select... distinct, with where condition", DefaultMysqlInspect(), `
	SELECT DISTINCT * FROM exist_db.exist_tb_1 WHERE id > 0;
	`, newTestResult())

	//insert... values
	runSingleRuleInspectCase(rule, t, "insert... values,", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 VALUES (1, "2", "3");
	`, newTestResult())

	//insert... select, with order by
	runSingleRuleInspectCase(rule, t, "insert... select, with order by", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_1 ORDER BY id;
	`, newTestResult().addResult(ruleName))

	//insert... select, with group by
	runSingleRuleInspectCase(rule, t, "insert... select, with group by", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_1 GROUP BY id;
	`, newTestResult().addResult(ruleName))

	//delete... where, with order by
	runSingleRuleInspectCase(rule, t, "delete... where, with order by", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 WHERE id = 1 ORDER BY id LIMIT 1;
	`, newTestResult())

	//delete... where, with order by
	runSingleRuleInspectCase(rule, t, "delete..., with order by", DefaultMysqlInspect(), `
	DELETE FROM exist_db.exist_tb_1 ORDER BY id LIMIT 1;
	`, newTestResult().addResult(ruleName))

	//update... set, with order by
	runSingleRuleInspectCase(rule, t, "update... set, with order by", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = "1" ORDER BY id LIMIT 1;
	`, newTestResult().addResult(ruleName))

	//update... where, with order by
	runSingleRuleInspectCase(rule, t, "update... where, with order by", DefaultMysqlInspect(), `
	UPDATE exist_db.exist_tb_1 SET v1 = "1" WHERE id = 1 ORDER BY id LIMIT 1;
	`, newTestResult())
}

// ==== Rule test code end ====
