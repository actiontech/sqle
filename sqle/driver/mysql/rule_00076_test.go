package mysql

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/stretchr/testify/assert"
)

// ==== Rule test code start ====
// For rule involving online information, use NewMockExecutor to simulate sql statements.
func NewMySQLInspectOnRuleSQLE0076(t *testing.T, sql string, planType string, rows int) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)

	handler.ExpectQuery(regexp.QuoteMeta("EXPLAIN " + sql)).
		WillReturnRows(sqlmock.NewRows([]string{"select_type", "rows"}).AddRow(planType, rows))
	handler.ExpectQuery(regexp.QuoteMeta("SHOW WARNINGS")).WillReturnRows(sqlmock.NewRows(nil))

	return inspect
}

func TestRuleSQLE00076(t *testing.T) {
	ruleName := ai.SQLE00076
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runSingleRuleInspectCase(rule, t, "case 1", NewMySQLInspectOnRuleSQLE0076(t, "delete from exist_db.exist_tb_1", "DELETE", 15000),
		`delete from exist_db.exist_tb_1`, newTestResult().addResult(ruleName))
	runSingleRuleInspectCase(rule, t, "case 1", NewMySQLInspectOnRuleSQLE0076(t, "delete from exist_db.exist_tb_1", "DELETE", 100),
		`delete from exist_db.exist_tb_1`, newTestResult())

	runSingleRuleInspectCase(rule, t, "case 1", NewMySQLInspectOnRuleSQLE0076(t, "UPDATE exist_db.exist_tb_1 set id = 0", "UPDATE", 15000),
		`UPDATE exist_db.exist_tb_1 set id = 0`, newTestResult().addResult(ruleName))
	runSingleRuleInspectCase(rule, t, "case 1", NewMySQLInspectOnRuleSQLE0076(t, "UPDATE exist_db.exist_tb_1 set id = 0", "UPDATE", 100),
		`UPDATE exist_db.exist_tb_1 set id = 0`, newTestResult())

}

// ==== Rule test code end ====
