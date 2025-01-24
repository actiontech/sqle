package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00177(t *testing.T) {
	ruleName := ai.SQLE00177
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select..., no order by", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult())

	//select... with order by, less than default value
	runSingleRuleInspectCase(rule, t, "select..., with order by, less than default value", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 ORDER BY v1, v2, v3;
	`, newTestResult())

	//select... with order by, more than default value
	runSingleRuleInspectCase(rule, t, "select..., with order by, more than default value", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 ORDER BY v1, v2, v3, v4;
	`, newTestResult().addResult(ruleName, 3))

	//insert... select..., no order by
	runSingleRuleInspectCase(rule, t, "insert... select..., no order by", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2;
	`, newTestResult())

	//insert... select... with order by, less than default value
	runSingleRuleInspectCase(rule, t, "insert... select..., with order by, less than default value", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2 ORDER BY v1, v2, v3;
	`, newTestResult())

	//insert... select... with order by, more than default value
	runSingleRuleInspectCase(rule, t, "insert... select..., with order by, more than default value", DefaultMysqlInspect(), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2 ORDER BY v1, v2, v3, v4;
	`, newTestResult().addResult(ruleName, 3))

	//union... select..., no order by
	runSingleRuleInspectCase(rule, t, "union... select..., no order by", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union... select... with order by, less than default value
	runSingleRuleInspectCase(rule, t, "union... select..., with order by, less than default value", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL (SELECT * FROM exist_db.exist_tb_2 ORDER BY v1, v2, v3);
	`, newTestResult())

	//union... select... with order by, more than default value
	runSingleRuleInspectCase(rule, t, "union... select..., with order by, more than default value", DefaultMysqlInspect(), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL (SELECT * FROM exist_db.exist_tb_2 ORDER BY v1, v2, v3, v4);
	`, newTestResult().addResult(ruleName, 3))
}

// ==== Rule test code end ====
