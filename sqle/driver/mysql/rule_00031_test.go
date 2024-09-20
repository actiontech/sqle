package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
)

// ==== Rule test code start ====
func TestRuleSQLE00031(t *testing.T) {
	ruleName := "SQLE00031"
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	for _, sql := range []string{
		`CREATE VIEW test_view AS SELECT * FROM test_table;`,
		`CREATE VIEW complex_view AS SELECT a.id, b.name FROM table_a a JOIN table_b b ON a.id = b.id;`,
		`CREATE VIEW subquery_view AS SELECT id FROM test_table WHERE id IN (SELECT id FROM another_table);`,
		`CREATE VIEW aggregate_view AS SELECT COUNT(*) FROM test_table;`,
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(ruleName))
	}

	for _, sql := range []string{
		`CREATE TABLE test_table (id INT PRIMARY KEY);`,
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}
}

// ==== Rule test code end ====
