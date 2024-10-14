package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00140(t *testing.T) {
	ruleName := ai.SQLE00140
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE 时指定库名",
		"CREATE TABLE testdb.test_table (id INT);",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb;"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE TABLE 时未指定库名",
		"CREATE TABLE test_table (id INT);",
		nil,
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: ALTER TABLE 时指定库名",
		"ALTER TABLE testdb.test_table ADD (name VARCHAR(100));",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb; CREATE TABLE testdb.test_table (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: ALTER TABLE 时未指定库名",
		"ALTER TABLE test_table ADD (name VARCHAR(100));",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: DROP TABLE 时指定库名",
		"DROP TABLE testdb.test_table;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb; CREATE TABLE testdb.test_table (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: DROP TABLE 时未指定库名",
		"DROP TABLE test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: SELECT 时指定库名",
		"SELECT * FROM testdb.test_table;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb; CREATE TABLE testdb.test_table (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: SELECT 时未指定库名",
		"SELECT * FROM test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: SELECT UNION 时所有子句均指定库名",
		"SELECT * FROM testdb.table1 UNION SELECT * FROM testdb.table2;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb; CREATE TABLE testdb.table1 (id INT); CREATE TABLE testdb.table2 (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: SELECT UNION 时部分子句未指定库名",
		"SELECT * FROM testdb.table1 UNION SELECT * FROM table2;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb; CREATE TABLE testdb.table1 (id INT); CREATE TABLE table2 (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: INSERT INTO 时指定库名",
		"INSERT INTO testdb.test_table (id, name) VALUES (1, 'Alice');",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb; CREATE TABLE testdb.test_table (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 12: INSERT INTO 时未指定库名",
		"INSERT INTO test_table (id, name) VALUES (1, 'Alice');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: UPDATE 时指定库名",
		"UPDATE testdb.test_table SET name = 'Bob' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb; CREATE TABLE testdb.test_table (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 14: UPDATE 时未指定库名",
		"UPDATE test_table SET name = 'Bob' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: DELETE 时指定库名",
		"DELETE FROM testdb.test_table WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE testdb; CREATE TABLE testdb.test_table (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 16: DELETE 时未指定库名",
		"DELETE FROM test_table WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 19: CREATE VIEW 时指定库名(从xml中补充)",
		"CREATE VIEW private.v1 AS SELECT * FROM private.t1;",
		session.NewAIMockContext().WithSQL("CREATE DATABASE private; CREATE TABLE private.t1 (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 20: CREATE VIEW 时未指定库名(从xml中补充)",
		"CREATE VIEW v1 AS SELECT * FROM t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 29: CREATE INDEX 时指定库名(从xml中补充)",
		"CREATE INDEX idx_t1_id ON private.t1(id);",
		session.NewAIMockContext().WithSQL("CREATE DATABASE private; CREATE TABLE private.t1 (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 30: CREATE INDEX 时未指定库名(从xml中补充)",
		"CREATE INDEX idx_t1_id ON t1(id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT);"),
		nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
