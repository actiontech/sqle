package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00045(t *testing.T) {
	ruleName := ai.SQLE00045
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// Case 1: SELECT ..., LIMIT offset exceeds 10000
	runSingleRuleInspectCase(rule, t, "Case 1: SELECT ..., LIMIT offset exceeds 10000", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 LIMIT 5 OFFSET 50000;
	`, newTestResult().addResult(ruleName, 10000))

	// Case 2: SELECT ..., LIMIT offset is equal to 10000
	runSingleRuleInspectCase(rule, t, "Case 2: SELECT ..., LIMIT offset is equal to 10000", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 LIMIT 10000, 10;
	`, newTestResult())

	// Case 3: INSERT ... SELECT ..., LIMIT offset exceeds 10000
	runSingleRuleInspectCase(rule, t, "Case 3: INSERT ... SELECT ..., LIMIT offset exceeds 10000", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_1 LIMIT 10001, 10;
	`, newTestResult().addResult(ruleName, 10000))

	// Case 4: INSERT ... SELECT ..., LIMIT offset is equal to 10000
	runSingleRuleInspectCase(rule, t, "Case 4: INSERT ... SELECT ..., LIMIT offset is equal to 10000", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_1 LIMIT 10 OFFSET 10000;
	`, newTestResult())

	// Case 5: UNION ALL ..., first SELECT clause, LIMIT offset exceeds 10000
	runSingleRuleInspectCase(rule, t, "Case 5: UNION ALL ...,first SELECT clause, LIMIT offset exceeds 10000", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 LIMIT 10001, 10 UNION ALL SELECT * FROM exist_db.exist_tb_1 LIMIT 10;
	`, newTestResult().addResult(ruleName, 10000))

	// Case 6: UNION ALL ...,second SELECT clause, LIMIT offset exceeds 10000
	runSingleRuleInspectCase(rule, t, "Case 6: UNION ALL ...,second SELECT clause, LIMIT offset exceeds 10000", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 LIMIT 10 UNION ALL SELECT * FROM exist_db.exist_tb_1 LIMIT 10001, 10;
	`, newTestResult().addResult(ruleName, 10000))

	// Case 7: UNION ...,first SELECT clause, LIMIT offset exceeds 10000
	runSingleRuleInspectCase(rule, t, "Case 7: UNION ...,first SELECT clause, LIMIT offset exceeds 10000", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 LIMIT 10001, 10 UNION SELECT * FROM exist_db.exist_tb_1 LIMIT 10;
	`, newTestResult().addResult(ruleName, 10000))

	// Case 8: UNION ..., second SELECT clause LIMIT offset exceeds 10000
	runSingleRuleInspectCase(rule, t, "Case 8: UNION ...,second SELECT clause, LIMIT offset exceeds 10000", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 LIMIT 10 UNION SELECT * FROM exist_db.exist_tb_1 LIMIT 10001, 10;
	`, newTestResult().addResult(ruleName, 10000))

}

// ==== Rule test code end ====
