package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====

func TestRuleSQLE00121(t *testing.T) {
	ruleName := ai.SQLE00121
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// Test: SELECT without ORDER BY and without LIMIT
	runSingleRuleInspectCase(rule, t, "select... without order by, without limit", DefaultMysqlInspect(), `
    SELECT * FROM exist_db.exist_tb_1;
    `, newTestResult())

	// Test: SELECT with LIMIT but without ORDER BY
	runSingleRuleInspectCase(rule, t, "select... with limit, without order by", DefaultMysqlInspect(), `
    SELECT * FROM exist_db.exist_tb_1 LIMIT 1;
    `, newTestResult().addResult(ruleName))

	// Test: SELECT with both LIMIT and ORDER BY
	runSingleRuleInspectCase(rule, t, "select... with limit, with order by", DefaultMysqlInspect(), `
    SELECT * FROM exist_db.exist_tb_1 ORDER BY v1 LIMIT 1;
    `, newTestResult())

	// Test: INSERT with a SELECT subquery that lacks both ORDER BY and LIMIT
	runSingleRuleInspectCase(rule, t, "insert... select... without order by, without limit", DefaultMysqlInspect(), `
    INSERT INTO exist_db.exist_tb_2 SELECT * FROM exist_db.exist_tb_1;
    `, newTestResult())

	// Test: INSERT with a SELECT subquery that has LIMIT but lacks ORDER BY
	runSingleRuleInspectCase(rule, t, "insert... select... with limit, without order by", DefaultMysqlInspect(), `
    INSERT INTO exist_db.exist_tb_2 SELECT * FROM exist_db.exist_tb_1 LIMIT 1;
    `, newTestResult().addResult(ruleName))

	// Test: INSERT with a SELECT subquery that includes both LIMIT and ORDER BY
	runSingleRuleInspectCase(rule, t, "insert... select... with limit, with order by", DefaultMysqlInspect(), `
    INSERT INTO exist_db.exist_tb_2 SELECT * FROM exist_db.exist_tb_1 ORDER BY v1 LIMIT 1;
    `, newTestResult())

	// Test: UNION of two SELECT queries without ORDER BY and without LIMIT
	runSingleRuleInspectCase(rule, t, "union... select... without order by, without limit", DefaultMysqlInspect(), `
    SELECT * FROM exist_db.exist_tb_1 UNION SELECT * FROM exist_db.exist_tb_2;
    `, newTestResult())

	// Test: UNION of two SELECT queries with LIMIT but without ORDER BY
	runSingleRuleInspectCase(rule, t, "union... select... with limit, without order by", DefaultMysqlInspect(), `
    SELECT * FROM exist_db.exist_tb_1 UNION SELECT * FROM exist_db.exist_tb_2 LIMIT 1;
    `, newTestResult())

	// Test: UNION of two SELECT queries with both LIMIT and ORDER BY
	runSingleRuleInspectCase(rule, t, "union... select... with limit, with order by", DefaultMysqlInspect(), `
    SELECT * FROM exist_db.exist_tb_1 ORDER BY v1 UNION SELECT * FROM exist_db.exist_tb_2 ORDER BY v2 LIMIT 1;
    `, newTestResult())

	// Test: UPDATE using a SELECT subquery without ORDER BY and without LIMIT
	runSingleRuleInspectCase(rule, t, "update... select... without order by, without limit", DefaultMysqlInspect(), `
    UPDATE exist_db.exist_tb_1 SET v1 = (SELECT MAX(v2) FROM exist_db.exist_tb_2);
    `, newTestResult())

	// Test: UPDATE using a SELECT subquery with LIMIT but without ORDER BY
	runSingleRuleInspectCase(rule, t, "update... select... with limit, without order by", DefaultMysqlInspect(), `
    UPDATE exist_db.exist_tb_1 SET v1 = (SELECT MAX(v2) FROM exist_db.exist_tb_2 LIMIT 1);
    `, newTestResult().addResult(ruleName))

	// Test: UPDATE using a SELECT subquery with both LIMIT and ORDER BY
	runSingleRuleInspectCase(rule, t, "update... select... with limit, with order by", DefaultMysqlInspect(), `
    UPDATE exist_db.exist_tb_1 SET v1 = (SELECT v2 FROM exist_db.exist_tb_2 ORDER BY v2 DESC LIMIT 1);
    `, newTestResult())

	// Test: DELETE using a SELECT subquery without ORDER BY and without LIMIT
	runSingleRuleInspectCase(rule, t, "delete... select...without order by, without limit", DefaultMysqlInspect(), `
    DELETE FROM exist_db.exist_tb_1 WHERE v1 IN (SELECT v2 FROM exist_db.exist_tb_2);
    `, newTestResult())

	// Test: DELETE using a SELECT subquery with LIMIT but without ORDER BY
	runSingleRuleInspectCase(rule, t, "delete... select... with limit, without order by", DefaultMysqlInspect(), `
    DELETE FROM exist_db.exist_tb_1 WHERE v1 IN (SELECT v2 FROM exist_db.exist_tb_2 LIMIT 1);
    `, newTestResult().addResult(ruleName))

	// Test: DELETE using a SELECT subquery with both LIMIT and ORDER BY
	runSingleRuleInspectCase(rule, t, "delete... select... with limit, with order by", DefaultMysqlInspect(), `
    DELETE FROM exist_db.exist_tb_1 WHERE v1 IN (SELECT v2 FROM exist_db.exist_tb_2 ORDER BY v2 DESC LIMIT 1);
    `, newTestResult())
}

// ==== Rule test code end ====
