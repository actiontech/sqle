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
func NewMySQLInspectOnRuleSQLE00111(t *testing.T) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)

	for i := 0; i < 100; i++ {
		handler.ExpectQuery(regexp.QuoteMeta("SHOW INDEX FROM `exist_db`.`exist_tb_1`")).
			WillReturnRows(sqlmock.NewRows([]string{"Expression"}).AddRow("lower(`v1`)"))
	}

	return inspect
}

func TestRuleSQLE00111(t *testing.T) {
	ruleName := ai.SQLE00111
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule
	i := NewMySQLInspectOnRuleSQLE00111(t)

	// without function call or math operation in where clause
	runSingleRuleInspectCase(rule, t, "without function call or math operation in where clause", i, `
    SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id > 0
    `, newTestResult())

	// with function call in where clause
	for _, sql := range []string{
		"select id,v1,v2 from exist_db.exist_tb_1 where year(id) = 1",
		"select id,v1,v2 from exist_db.exist_tb_1 where CONCAT(id, v1) = 1",
		"update exist_db.exist_tb_1 set id = 1 where year(id) = 1",
		"update exist_db.exist_tb_1 set id = 1 where CONCAT(id, v1) = 1",
		"delete from exist_db.exist_tb_1 where CONCAT(id, v1) = 1",
		"delete from exist_db.exist_tb_1 where year(id) = 1",
		"delete from exist_db.exist_tb_1 where LENGTH(v1) + 1 = 10",
		"delete from exist_db.exist_tb_1 where DAY(v1_date) - 1 = 15",
		"delete from exist_db.exist_tb_1 where CONCAT(v1, '_', v2) = 'value_combined'",
	} {
		runSingleRuleInspectCase(rule, t, "with function call in where clause", i, sql, newTestResult().addResult(ruleName))
	}

	// with math operation in where clause
	for _, sql := range []string{
		"select id,v1,v2 from exist_db.exist_tb_1 where id + 1 = 1",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where id - 1 = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where id * 1 = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where id / 1 = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where id % 1 = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where id MOD 1 = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where id DIV 1 = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where -id = 12",

		"select id,v1,v2 from exist_db.exist_tb_1 where 1 + id = 1",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where 1- id = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where 1 * id = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where 1 / id = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where 1 % id = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where 1 MOD id = 12",
		"SELECT id,v1,v2 from exist_db.exist_tb_1 where 1 DIV id = 12",

		"select id,v1,v2 from exist_db.exist_tb_1 where (SELECT id from exist_db.exist_tb_1 where -id = 12 limit 1) = 1",
		"select id,v1,v2 from exist_db.exist_tb_1 where (SELECT id from exist_db.exist_tb_1 where id + 1 = 12 limit 1) = 1",
		"select id,v1,v2 from exist_db.exist_tb_1 where (SELECT id from exist_db.exist_tb_1 where id * 1 = 12 limit 1) = 1",
		"select id,v1,v2 from exist_db.exist_tb_1 where (SELECT id from exist_db.exist_tb_1 where id / 1 = 12 limit 1) = 1",
		"select id,v1,v2 from exist_db.exist_tb_1 where (SELECT id from exist_db.exist_tb_1 where id % 1 = 12 limit 1) = 1",
		"select id,v1,v2 from exist_db.exist_tb_1 where (SELECT id from exist_db.exist_tb_1 where id DIV 1 = 12 limit 1) = 1",
		"select (SELECT id from exist_db.exist_tb_1 where id DIV 1 = 12 limit 1),v1,v2 from exist_db.exist_tb_1",
		"select (SELECT id from exist_db.exist_tb_1 where id / 1 = 12 limit 1),v1,v2 from exist_db.exist_tb_1",
		"select (SELECT id from exist_db.exist_tb_1 where id * 1 = 12 limit 1),v1,v2 from exist_db.exist_tb_1",
		"select (SELECT (SELECT id from exist_db.exist_tb_1 where id * 1 = 12 limit 1) from exist_db.exist_tb_1 limit 1),v1,v2 from exist_db.exist_tb_1",
		"update exist_db.exist_tb_1 set id = 1 where id + 1 = 1",
		"update exist_db.exist_tb_1 set id = 1 where id - 1 = 1",
		"update exist_db.exist_tb_1 set id = 1 where id * 1 = 1",
		"update exist_db.exist_tb_1 set id = 1 where id / 1 = 1",
		"update exist_db.exist_tb_1 set id = 1 where id % 1 = 1",
		"update exist_db.exist_tb_1 set id = 1 where id MOD 1 = 1",
		"update exist_db.exist_tb_1 set id = 1 where id DIV 1 = 1",
		"update exist_db.exist_tb_1 set id = 1 where -id = 1",

		"update exist_db.exist_tb_1 set id = 1 where 1 + id = 1",
		"update exist_db.exist_tb_1 set id = 1 where 1 - id = 1",
		"update exist_db.exist_tb_1 set id = 1 where 1 * id = 1",
		"update exist_db.exist_tb_1 set id = 1 where 1 / id = 1",
		"update exist_db.exist_tb_1 set id = 1 where 1 % id = 1",
		"update exist_db.exist_tb_1 set id = 1 where 1 MOD id = 1",
		"update exist_db.exist_tb_1 set id = 1 where 1 DIV id = 1",
		"update exist_db.exist_tb_1 set id = 1 where -id = 1",

		"delete from exist_db.exist_tb_1 where id + 1 = 1",
		"delete from exist_db.exist_tb_1 where id - 1 = 1",
		"delete from exist_db.exist_tb_1 where id * 1 = 1",
		"delete from exist_db.exist_tb_1 where id / 1 = 1",
		"delete from exist_db.exist_tb_1 where id % 1 = 1",
		"delete from exist_db.exist_tb_1 where id MOD 1 = 1",
		"delete from exist_db.exist_tb_1 where id DIV 1 = 1",
		"delete from exist_db.exist_tb_1 where -id = 1",

		"delete from exist_db.exist_tb_1 where 1 + id = 1",
		"delete from exist_db.exist_tb_1 where 1 - id = 1",
		"delete from exist_db.exist_tb_1 where 1 * id = 1",
		"delete from exist_db.exist_tb_1 where 1 / id = 1",
		"delete from exist_db.exist_tb_1 where 1 % id = 1",
		"delete from exist_db.exist_tb_1 where 1 MOD id = 1",
		"delete from exist_db.exist_tb_1 where 1 DIV id = 1",
	} {
		runSingleRuleInspectCase(rule, t, "with math operation in where clause", i, sql, newTestResult().addResult(ruleName))
	}

	//with union, without function call or math operation in where clause
	runSingleRuleInspectCase(rule, t, "without function call or math operation in union clause", i, `
    (SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id > 0)
    UNION ALL
    (SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id < 100)
    `, newTestResult())

	//with union, with function call in where clause
	runSingleRuleInspectCase(rule, t, "with function call in union clause", i, `
    (SELECT id, v1 FROM exist_db.exist_tb_1 WHERE CONCAT(id, v1) = 1)
    UNION ALL
    (SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id < 100)
    `, newTestResult().addResult(ruleName))

	//with union, with math operation in where clause
	runSingleRuleInspectCase(rule, t, "with math operation in union clause", i, `
    (SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id + 1 = 1)
    UNION ALL
    (SELECT id, v1 FROM exist_db.exist_tb_1 WHERE id < 100)
    `, newTestResult().addResult(ruleName))

	// with function call in where clause and exist function index
	runSingleRuleInspectCase(rule, t, " with function call in where clause and exist function index", i, `
	SELECT id, v1 FROM exist_db.exist_tb_1 WHERE lower(v1) = '1'
	`, newTestResult())

	// with function call in where clause and not exist function index
	runSingleRuleInspectCase(rule, t, " with function call in where clause and not exist function index", i, `
		SELECT id, v1 FROM exist_db.exist_tb_1 WHERE lower(v1) = '1' AND UPPER(v1) = '1'
		`, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
