package mysql

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/stretchr/testify/assert"
)

// ==== Rule test code start ====
// For rule involving online information, use NewMockExecutor to simulate sql statements.
func NewMySQLInspectOnRuleSQLE00139(t *testing.T, sql string, planTypes, tableNames []string, tableSize map[string] /*table name*/ int /*table size MB*/) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)
	inspect.Ctx = session.NewMockContextForTestTableSize(e, tableSize)

	assert.Equal(t, len(planTypes), len(tableNames))
	r := sqlmock.NewRows([]string{"type", "table"})
	for i, t := range planTypes {
		r.AddRow(t, tableNames[i])
	}
	handler.ExpectQuery(regexp.QuoteMeta("EXPLAIN " + sql)).WillReturnRows(r)

	handler.ExpectQuery(regexp.QuoteMeta("SHOW WARNINGS")).WillReturnRows(sqlmock.NewRows(nil))

	return inspect
}

func TestRuleSQLE00139(t *testing.T) {
	ruleName := ai.SQLE00139
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule
	ruleParams := []interface{}{5}

	//select, no problem
	sql := "SELECT * FROM exist_db.exist_tb_1 where id =1"
	runSingleRuleInspectCase(rule, t, "select, no problem",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"index"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 1}),
		sql, newTestResult())

	//select, full table scan, table size is less than threshold
	sql = "SELECT * FROM exist_db.exist_tb_1 where id =1"
	runSingleRuleInspectCase(rule, t, "select, full table scan, table size is less than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 1}),
		sql, newTestResult())

	//select, full table scan, table size is greater than threshold
	sql = "SELECT * FROM exist_db.exist_tb_1 where id =1"
	runSingleRuleInspectCase(rule, t, "select, full table scan, table size is greater than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 6144}),
		sql, newTestResult().addResult(ruleName, ruleParams...))

	//union, no problem
	sql = "SELECT * FROM exist_db.exist_tb_1 where id =1 UNION ALL SELECT id, v1, v2 FROM exist_db.exist_tb_2 where id =1"
	runSingleRuleInspectCase(rule, t, "union, no problem",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"index", "index"}, []string{"exist_tb_1", "exist_tb_2"}, map[string]int{"exist_tb_1": 100, "exist_tb_2": 100}),
		sql, newTestResult())

	// union, full table scan, table size is less than threshold
	sql = "SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1, v2 FROM exist_db.exist_tb_2"
	runSingleRuleInspectCase(rule, t, "union, full table scan, table size is less than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL", "ALL"}, []string{"exist_tb_1", "exist_tb_2"}, map[string]int{"exist_tb_1": 4, "exist_tb_2": 4}),
		sql, newTestResult())

	// union, full table scan, one of the table size is less than threshold
	sql = "SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1, v2 FROM exist_db.exist_tb_2"
	runSingleRuleInspectCase(rule, t, "union, full table scan, one of the table size is less than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"index", "ALL"}, []string{"exist_tb_1", "exist_tb_2"}, map[string]int{"exist_tb_1": 6144, "exist_tb_2": 4}),
		sql, newTestResult())

	// union, full table scan, table size is greater than threshold
	sql = "SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1, v2 FROM exist_db.exist_tb_2"
	runSingleRuleInspectCase(rule, t, "union, full table scan, table size is greater than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL", "ALL"}, []string{"exist_tb_1", "exist_tb_2"}, map[string]int{"exist_tb_1": 6144, "exist_tb_2": 6144}),
		sql, newTestResult().addResult(ruleName, ruleParams...))

	// union, full table scan, one of the table size is greater than threshold
	sql = "SELECT * FROM exist_db.exist_tb_1 UNION ALL SELECT id, v1, v2 FROM exist_db.exist_tb_2"
	runSingleRuleInspectCase(rule, t, "union, full table scan, one of the table size is greater than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL", "ALL"}, []string{"exist_tb_1", "exist_tb_2"}, map[string]int{"exist_tb_1": 6144, "exist_tb_2": 100}),
		sql, newTestResult().addResult(ruleName, ruleParams...))

	//update, no problem
	sql = "UPDATE exist_db.exist_tb_1 SET v1 = 'value' WHERE id = 1"
	runSingleRuleInspectCase(rule, t, "update, no problem",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"index"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 100}),
		sql, newTestResult())

	//update, full table scan, table size is less than threshold
	sql = "UPDATE exist_db.exist_tb_1 SET v1 = 'value' WHERE v2 > 50"
	runSingleRuleInspectCase(rule, t, "update, full table scan, table size is less than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 4}),
		sql, newTestResult())

	//update, full table scan, table size is greater than threshold
	sql = "UPDATE exist_db.exist_tb_1 SET v1 = 'value' WHERE v2 > 50"
	runSingleRuleInspectCase(rule, t, "update, full table scan, table size is greater than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 6144}),
		sql, newTestResult().addResult(ruleName, ruleParams...))

	//delete, no problem
	sql = "DELETE FROM exist_db.exist_tb_1 WHERE id = 1"
	runSingleRuleInspectCase(rule, t, "delete, no problem",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"index"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 100}),
		sql, newTestResult())

	//delete, full table scan, table size is less than threshold
	sql = "DELETE FROM exist_db.exist_tb_1 WHERE v2 > 50"
	runSingleRuleInspectCase(rule, t, "delete, full table scan, table size is less than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 4}),
		sql, newTestResult())

	//delete, full table scan, table size is greater than threshold
	sql = "DELETE FROM exist_db.exist_tb_1 WHERE v2 > 50"
	runSingleRuleInspectCase(rule, t, "delete, full table scan, table size is greater than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 6144}),
		sql, newTestResult().addResult(ruleName, ruleParams...))

	// insert, no problem (subquery uses index)
	sql = "INSERT INTO exist_db.exist_tb_1 (id, v1) SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id = 1"
	runSingleRuleInspectCase(rule, t, "insert, no problem",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"index"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 100}),
		sql, newTestResult())

	// insert, subquery full table scan, table size is less than threshold
	sql = "INSERT INTO exist_db.exist_tb_1 (id, v1) SELECT id, v1 FROM exist_db.exist_tb_1 WHERE v2 > 50"
	runSingleRuleInspectCase(rule, t, "insert, full table scan, table size is less than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 4}),
		sql, newTestResult())

	// insert, subquery full table scan, table size is greater than threshold
	sql = "INSERT INTO exist_db.exist_tb_1 (id, v1) SELECT id, v1 FROM exist_db.exist_tb_1 WHERE v2 > 50"
	runSingleRuleInspectCase(rule, t, "insert, full table scan, table size is greater than threshold",
		NewMySQLInspectOnRuleSQLE00139(t, sql, []string{"ALL"}, []string{"exist_tb_1"}, map[string]int{"exist_tb_1": 6144}),
		sql, newTestResult().addResult(ruleName, ruleParams...))

}

// ==== Rule test code end ====
