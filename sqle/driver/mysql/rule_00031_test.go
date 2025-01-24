package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00031(t *testing.T) {
	ruleName := ai.SQLE00031
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	for _, sql := range []string{
		`CREATE VIEW test_view AS SELECT * FROM exist_db.exist_tb_1;`,
		`CREATE VIEW complex_view AS SELECT a.id, b.name FROM exist_db.exist_tb_1 a JOIN  exist_db.exist_tb_2 b ON a.id = b.id;`,
		`CREATE VIEW subquery_view AS SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_2);`,
		`CREATE VIEW aggregate_view AS SELECT COUNT(*) FROM exist_db.exist_tb_1;`,
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(ruleName))
	}

	for _, sql := range []string{
		`CREATE TABLE exist_db.not_exist_tb_1 (id INT PRIMARY KEY);`,
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}
}

// ==== Rule test code end ====
