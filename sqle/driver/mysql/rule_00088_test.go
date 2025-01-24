package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQL00088(t *testing.T) {
	ruleName := ai.SQLE00088
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: INSERT语句未指定列名", "INSERT INTO test_table VALUES (1, 'example');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: INSERT语句指定列名", "INSERT INTO test_table (id, name) VALUES (1, 'example');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: INSERT语句部分指定列名", "INSERT INTO test_table (id) VALUES (1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: INSERT语句指定所有列名", "INSERT INTO test_table (id, name, age) VALUES (1, 'example', 25);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: INSERT语句未指定列名并使用SELECT", "INSERT INTO t1 SELECT 3, cid, name FROM t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, cid INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: INSERT语句指定列名并使用SELECT", "INSERT INTO t1 (id, name, cid) SELECT id, name, cid FROM t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, cid INT, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: INSERT语句未指定列名并使用SELECT", "INSERT INTO t1 SELECT id, name, cid FROM t1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, cid INT, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
