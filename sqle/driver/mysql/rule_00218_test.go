package mysql

import (
	"fmt"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/stretchr/testify/assert"
)

// ==== Rule test code start ====

func NewMySQLInspectOnRuleSQLE00218(t *testing.T, tableNames ...string) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)
	handler.MatchExpectationsInOrder(false)

	inspect := NewMockInspect(e)

	for _, tableName := range tableNames {
		handler.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SHOW INDEX FROM `exist_db`.`%v`", tableName))).
			WillReturnRows(sqlmock.NewRows([]string{"Column_name", "Seq_in_index"}).AddRow("id", "1").AddRow("v2", "2"))
	}

	return inspect
}

func TestRuleSQLE00218(t *testing.T) {
	ruleName := ai.SQLE00218
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//select...where use leftmost index
	runSingleRuleInspectCase(rule, t, "select...where use leftmost index", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE id = 1;
	`, newTestResult())

	//select...where not use leftmost index
	runSingleRuleInspectCase(rule, t, "select...where not use leftmost index", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE v2 = 1;
	`, newTestResult().addResult(ruleName, "v2"))

	//select...where 1=1
	runSingleRuleInspectCase(rule, t, "select...where 1=1", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE 1=1;
	`, newTestResult())

	//select...where True
	runSingleRuleInspectCase(rule, t, "select...where True", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE True;
	`, newTestResult())

	//select...where use leftmost index, with group by
	runSingleRuleInspectCase(rule, t, "select...where use leftmost index, with group by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE id = 1 GROUP BY id;
	`, newTestResult())

	//select...where not use leftmost index, with group by
	runSingleRuleInspectCase(rule, t, "select...where not use leftmost index, with group by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE v2 = 1 GROUP BY id;
	`, newTestResult().addResult(ruleName, "v2"))

	//select...where use leftmost index, with order by
	runSingleRuleInspectCase(rule, t, "select...where use leftmost index, with order by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE id = 1 ORDER BY id;
	`, newTestResult())

	//select...where not use leftmost index, with order by
	runSingleRuleInspectCase(rule, t, "select...where not use leftmost index, with order by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE v2 = 1 ORDER BY id;
	`, newTestResult().addResult(ruleName, "v2"))

	//select...where 1=1, with group by
	runSingleRuleInspectCase(rule, t, "select...where 1=1, with group by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE 1=1 GROUP BY id;
	`, newTestResult())

	//select...where 1=1, with group by, not use leftmost index
	runSingleRuleInspectCase(rule, t, "select...where 1=1, with group by, not use leftmost index", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE 1=1 GROUP BY v2;
	`, newTestResult().addResult(ruleName, "v2"))

	//select...where 1=1, with order by
	runSingleRuleInspectCase(rule, t, "select...where 1=1, with order by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE 1=1 ORDER BY id;
	`, newTestResult())

	//select...where 1=1, with order by, not use leftmost index
	runSingleRuleInspectCase(rule, t, "select...where 1=1, with order by, not use leftmost index", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1"), `
	SELECT * FROM exist_db.exist_tb_1 WHERE 1=1 ORDER BY v2;
	`, newTestResult().addResult(ruleName, "v2"))

	//select...where normal, with join, on condition not use leftmost index, with group by
	runSingleRuleInspectCase(rule, t, "select...where normal, with join, on condition not use leftmost index, with group by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2"), `
	SELECT * FROM exist_db.exist_tb_1 a LEFT JOIN exist_db.exist_tb_2 b ON a.v2 = b.v2 WHERE a.id = 1 GROUP BY a.id;
	`, newTestResult().addResult(ruleName, "a.v2").addResult(ruleName, "b.v2"))

	//select...where normal, with join, on condition not use leftmost index, with order by
	runSingleRuleInspectCase(rule, t, "select...where normal, with join, on condition not use leftmost index, with order by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2"), `
	SELECT * FROM exist_db.exist_tb_1 a LEFT JOIN exist_db.exist_tb_2 b ON a.v2 = b.v2 WHERE a.id = 1 ORDER BY a.id;
	`, newTestResult().addResult(ruleName, "a.v2").addResult(ruleName, "b.v2"))

	//select...where 1=1, with join, on condition not use leftmost index, with group by
	runSingleRuleInspectCase(rule, t, "select...where 1=1, with join, on condition not use leftmost index, with group by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2"), `
	SELECT * FROM exist_db.exist_tb_1 a LEFT JOIN exist_db.exist_tb_2 b ON a.v2 = b.v2 WHERE 1=1 GROUP BY a.id;
	`, newTestResult().addResult(ruleName, "a.v2").addResult(ruleName, "b.v2"))

	//select...where 1=1, with join, on condition not use leftmost index, with order by
	runSingleRuleInspectCase(rule, t, "select...where 1=1, with join, on condition not use leftmost index, with order by", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2"), `
	SELECT * FROM exist_db.exist_tb_1 a LEFT JOIN exist_db.exist_tb_2 b ON a.v2 = b.v2 WHERE 1=1 ORDER BY a.id;
	`, newTestResult().addResult(ruleName, "a.v2").addResult(ruleName, "b.v2"))

	//insert...select...where use leftmost index
	runSingleRuleInspectCase(rule, t, "insert...select...where use leftmost index", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_3"), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_3 WHERE id = 1;
	`, newTestResult())

	//insert...select...where not use leftmost index
	runSingleRuleInspectCase(rule, t, "insert...select...where not use leftmost index", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_3"), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_3 WHERE v2 = 1;
	`, newTestResult().addResult(ruleName, "v2"))

	//insert...select...where 1=1
	runSingleRuleInspectCase(rule, t, "insert...select...where 1=1", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_3"), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_3 WHERE 1=1;
	`, newTestResult())

	//insert...select...where use leftmost index, with join
	runSingleRuleInspectCase(rule, t, "insert...select...where use leftmost index, with join", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_3", "exist_tb_4"), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_3 a LEFT JOIN exist_db.exist_tb_4 b ON a.id = b.id WHERE a.id = 1;
	`, newTestResult())

	//insert...select...where not use leftmost index, with join
	runSingleRuleInspectCase(rule, t, "insert...select...where not use leftmost index, with join", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_3", "exist_tb_4"), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_3 a LEFT JOIN exist_db.exist_tb_4 b ON a.v2 = b.v2 WHERE a.id = 1;
	`, newTestResult().addResult(ruleName, "a.v2").addResult(ruleName, "b.v2"))

	//insert...select...where 1=1, with join
	runSingleRuleInspectCase(rule, t, "insert...select...where 1=1, with join", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_3", "exist_tb_4"), `
	INSERT INTO exist_db.exist_tb_1 SELECT * FROM exist_db.exist_tb_3 a LEFT JOIN exist_db.exist_tb_4 b ON a.v2 = b.v2 WHERE 1=1;
	`, newTestResult().addResult(ruleName, "a.v2").addResult(ruleName, "b.v2"))

	//union...where use leftmost index
	runSingleRuleInspectCase(rule, t, "union...where use leftmost index", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2"), `
	SELECT * FROM (SELECT id FROM exist_db.exist_tb_1 WHERE id = 1 UNION ALL SELECT id FROM exist_db.exist_tb_2 WHERE id = 1) a;
	`, newTestResult())

	//union...where 1=1
	runSingleRuleInspectCase(rule, t, "union...where 1=1", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2"), `
	SELECT * FROM (SELECT id FROM exist_db.exist_tb_1 WHERE 1=1 UNION ALL SELECT id FROM exist_db.exist_tb_2 WHERE 1=1) a;
	`, newTestResult())

	//union...where use leftmost index, with join
	runSingleRuleInspectCase(rule, t, "union...where use leftmost index, with join", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2", "exist_tb_3", "exist_tb_4"), `
	SELECT * FROM (SELECT id FROM exist_db.exist_tb_1 a JOIN exist_db.exist_tb_2 b ON a.id = b.id WHERE a.id = 1 UNION ALL SELECT id FROM exist_db.exist_tb_3 a JOIN exist_db.exist_tb_4 b ON a.id = b.id WHERE a.id = 1) a;
	`, newTestResult())

	//union...where not use leftmost index, with join
	runSingleRuleInspectCase(rule, t, "union...where not use leftmost index, with join", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2", "exist_tb_3", "exist_tb_4"), `
	SELECT * FROM (SELECT id FROM exist_db.exist_tb_1 a JOIN exist_db.exist_tb_2 b ON a.v2 = b.v2 WHERE a.id = 1 UNION ALL SELECT id FROM exist_db.exist_tb_3 a JOIN exist_db.exist_tb_4 b ON a.v2 = b.v2 WHERE a.id = 1) a;
	`, newTestResult().addResult(ruleName, "a.v2").addResult(ruleName, "b.v2"))

	//union...where 1=1, with join
	runSingleRuleInspectCase(rule, t, "union...where 1=1, with join", NewMySQLInspectOnRuleSQLE00218(t, "exist_tb_1", "exist_tb_2", "exist_tb_3", "exist_tb_4"), `
	SELECT * FROM (SELECT id FROM exist_db.exist_tb_1 a JOIN exist_db.exist_tb_2 b ON a.v2 = b.v2 WHERE 1=1 UNION ALL SELECT id FROM exist_db.exist_tb_3 a JOIN exist_db.exist_tb_4 b ON a.v2 = b.v2 WHERE 1=1) a;
	`, newTestResult().addResult(ruleName, "a.v2").addResult(ruleName, "b.v2"))

}

// ==== Rule test code end ====
