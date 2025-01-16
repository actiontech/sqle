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
func NewMySQLInspectOnRuleSQLE00179(t *testing.T, sql string, withWarning bool) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)

	handler.ExpectQuery(regexp.QuoteMeta("EXPLAIN " + sql)).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("")) // The results of explain won't be used in this use case, so return empty here
	if withWarning {
		handler.ExpectQuery(regexp.QuoteMeta("SHOW WARNINGS")).
			WillReturnRows(sqlmock.NewRows([]string{"Message"}).AddRow("Cannot use ref access on index 'idx_log_date_customers' due to type or collation conversion on field 'log_date'"))
	} else {
		handler.ExpectQuery(regexp.QuoteMeta("SHOW WARNINGS")).WillReturnRows(sqlmock.NewRows(nil))
	}

	return inspect
}

func TestRuleSQLE00179(t *testing.T) {
	ruleName := ai.SQLE00179
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// SELECT statements testing
	runSingleRuleInspectCase(rule, t, "select..., no problem", NewMySQLInspectOnRuleSQLE00179(t, "SELECT * FROM exist_db.exist_tb_1 where v1 = concat(now(),'')", false), `
	SELECT * FROM exist_db.exist_tb_1 where v1 = concat(now(),'');
	`, newTestResult())

	runSingleRuleInspectCase(rule, t, "select..., with problem(due to type or collation conversion on field)", NewMySQLInspectOnRuleSQLE00179(t, "SELECT * FROM exist_db.exist_tb_1 where v1 = now()", true), `
	SELECT * FROM exist_db.exist_tb_1 where v1 = now();
	`, newTestResult().addResult(ruleName))

	// UPDATE statements testing
	// No problem scenario
	runSingleRuleInspectCase(rule, t, "update..., no problem", NewMySQLInspectOnRuleSQLE00179(t, "UPDATE exist_db.exist_tb_1 SET v1 = concat(now(),'') WHERE id = 1", false), `
	UPDATE exist_db.exist_tb_1 SET v1 = concat(now(),'') WHERE id = 1;
	`, newTestResult())

	// With problem scenario
	runSingleRuleInspectCase(rule, t, "update..., with problem(due to type or collation conversion on field)", NewMySQLInspectOnRuleSQLE00179(t, "UPDATE exist_db.exist_tb_1 SET v1 = now() WHERE id = 1", true), `
	UPDATE exist_db.exist_tb_1 SET v1 = now() WHERE id = 1;
	`, newTestResult().addResult(ruleName))

	// DELETE statements testing
	runSingleRuleInspectCase(rule, t, "delete..., no problem", NewMySQLInspectOnRuleSQLE00179(t, "DELETE FROM exist_db.exist_tb_1 WHERE v1 = concat(now(),'')", false), `
	DELETE FROM exist_db.exist_tb_1 WHERE v1 = concat(now(),'');
	`, newTestResult())

	// INSERT statements testing
	// No problem scenario
	runSingleRuleInspectCase(rule, t, "insert..., no problem", NewMySQLInspectOnRuleSQLE00179(t, "INSERT INTO exist_db.exist_tb_1 (v1) VALUES (concat(now(),''))", false), `
	INSERT INTO exist_db.exist_tb_1 (v1) VALUES (concat(now(),''));
	`, newTestResult())

	// With problem scenario (if applicable)
	runSingleRuleInspectCase(rule, t, "insert..., with problem(due to type or collation conversion on field)", NewMySQLInspectOnRuleSQLE00179(t, "INSERT INTO exist_db.exist_tb_1 (v1) VALUES (now())", true), `
	INSERT INTO exist_db.exist_tb_1 (v1) VALUES (now());
	`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
