package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00118(t *testing.T) {
	ruleName := ai.SQLE00118
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// case 1: DROP TABLE 语句用于删除表
	runAIRuleCase(rule, t, "case 1: DROP TABLE 语句用于删除表", "DROP TABLE test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	// case 2: TRUNCATE TABLE 语句用于清空表
	runAIRuleCase(rule, t, "case 2: TRUNCATE TABLE 语句用于清空表", "TRUNCATE TABLE test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	// case 3: SELECT 语句用于查询数据
	runAIRuleCase(rule, t, "case 3: SELECT 语句用于查询数据", "SELECT * FROM test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult())

	// case 4: INSERT 语句用于插入数据
	runAIRuleCase(rule, t, "case 4: INSERT 语句用于插入数据", "INSERT INTO test_table (id, name) VALUES (1, 'test');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult())

	// case 5: UPDATE 语句用于更新数据
	runAIRuleCase(rule, t, "case 5: UPDATE 语句用于更新数据", "UPDATE test_table SET name = 'updated' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult())

	// case 6: DELETE 语句用于删除数据
	runAIRuleCase(rule, t, "case 6: DELETE 语句用于删除数据", "DELETE FROM test_table WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult())

	// case 7: DROP TABLE 语句用于删除表 t1
	runAIRuleCase(rule, t, "case 7: DROP TABLE 语句用于删除表 t1", "DROP TABLE t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	// case 8: TRUNCATE TABLE 语句用于清空表 t1
	runAIRuleCase(rule, t, "case 8: TRUNCATE TABLE 语句用于清空表 t1", "TRUNCATE TABLE t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	// case 9: SELECT 语句用于查询表 t1 数据
	runAIRuleCase(rule, t, "case 9: SELECT 语句用于查询表 t1 数据", "SELECT * FROM t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, name VARCHAR(50));"),
		nil, newTestResult())
}

// ==== Rule test code end ====
