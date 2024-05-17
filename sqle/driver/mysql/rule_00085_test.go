package mysql

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/stretchr/testify/assert"
)

// ==== Rule test code start ====
// For rule involving online information, use NewMockExecutor to simulate sql statements.
func NewMySQLInspectOnRuleSQLE00085(t *testing.T, sql string, planType string) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)

	handler.ExpectQuery(regexp.QuoteMeta("EXPLAIN " + sql)).
		WillReturnRows(sqlmock.NewRows([]string{"type"}).AddRow(planType))
	handler.ExpectQuery(regexp.QuoteMeta("SHOW WARNINGS")).WillReturnRows(sqlmock.NewRows(nil))

	return inspect
}

func TestRuleSQLE00085(t *testing.T) {
	ruleName := ai.SQLE00085
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// select, no index
	i := NewMySQLInspectOnRuleSQLE00085(t, "SELECT * FROM exist_db.exist_tb_1", "ALL")
	runSingleRuleInspectCase(rule, t, "select, no index", i, `
	SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult())

	// select, with index
	i = NewMySQLInspectOnRuleSQLE00085(t, "SELECT v1 FROM exist_db.exist_tb_1 ORDER BY v1", "index")
	runSingleRuleInspectCase(rule, t, "select, with full index scan", i, `
	SELECT v1 FROM exist_db.exist_tb_1 ORDER BY v1;
	`, newTestResult().addResult(ruleName))

	// union, no index
	i = NewMySQLInspectOnRuleSQLE00085(t, "SELECT v1 FROM exist_db.exist_tb_1 UNION SELECT v1 FROM exist_db.exist_tb_2", "ALL")
	runSingleRuleInspectCase(rule, t, "union, no index", i, `
	SELECT v1 FROM exist_db.exist_tb_1 UNION SELECT v1 FROM exist_db.exist_tb_2;
	`, newTestResult())

	// union, with index
	i = NewMySQLInspectOnRuleSQLE00085(t, " (SELECT v1 FROM exist_db.exist_tb_1 ORDER BY v1) UNION SELECT v1 FROM exist_db.exist_tb_2", "index")
	runSingleRuleInspectCase(rule, t, "union, with index", i, `
	 (SELECT v1 FROM exist_db.exist_tb_1 ORDER BY v1) UNION SELECT v1 FROM exist_db.exist_tb_2;
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
