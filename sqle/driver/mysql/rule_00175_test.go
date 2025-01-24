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
func NewMySQLInspectOnRuleSQLE00175(t *testing.T, sql string, planType string) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)

	handler.ExpectQuery(regexp.QuoteMeta("EXPLAIN " + sql)).
		WillReturnRows(sqlmock.NewRows([]string{"type"}).AddRow(planType))
	handler.ExpectQuery(regexp.QuoteMeta("SHOW WARNINGS")).WillReturnRows(sqlmock.NewRows(nil))

	return inspect
}

func TestRuleSQLE00175(t *testing.T) {
	ruleName := ai.SQLE00175
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//select...
	runSingleRuleInspectCase(rule, t, "select..., no problem", NewMySQLInspectOnRuleSQLE00175(t, "SELECT * FROM exist_db.exist_tb_1", "index"), `
	SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult())

	//insert..., no problem
	runSingleRuleInspectCase(rule, t, "insert..., no problem", NewMySQLInspectOnRuleSQLE00175(t, "INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2", "index"), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_2;
	`, newTestResult())

	//insert..., with problem in select clause
	runSingleRuleInspectCase(rule, t, "insert..., with problem in select clause", NewMySQLInspectOnRuleSQLE00175(t, "INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_1", "index_merge"), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult().addResult(ruleName))

	//union..., no problem
	runSingleRuleInspectCase(rule, t, "union..., no problem", NewMySQLInspectOnRuleSQLE00175(t, "SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_2", "index"), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_2;
	`, newTestResult())

	//union..., with problem in first select clause
	runSingleRuleInspectCase(rule, t, "union..., with problem in first select clause", NewMySQLInspectOnRuleSQLE00175(t, "SELECT a FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_1", "index_merge"), `
	SELECT a FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult().addResult(ruleName))

	//union..., with problem in second select clause
	runSingleRuleInspectCase(rule, t, "union..., with problem in second select clause", NewMySQLInspectOnRuleSQLE00175(t, "SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_1", "index_merge"), `
	SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT * FROM exist_db.exist_tb_1;
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
