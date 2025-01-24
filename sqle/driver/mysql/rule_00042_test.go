package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQL00042(t *testing.T) {
	ruleName := ai.SQLE00042
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TEMPORARY TABLE without prefix", "CREATE TEMPORARY TABLE test_table (id INT);",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: CREATE TEMPORARY TABLE with correct prefix", "CREATE TEMPORARY TABLE tmp_test_table (id INT);",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: CREATE TABLE without TEMPORARY keyword", "CREATE TABLE test_table (id INT);",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: ALTER TABLE RENAME TEMPORARY TABLE without prefix", "ALTER TABLE test_table RENAME TO new_test_table;",
		session.NewAIMockContext().WithSQL("CREATE TEMPORARY TABLE test_table (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: ALTER TABLE RENAME TEMPORARY TABLE with correct prefix", "ALTER TABLE tmp_test_table RENAME TO tmp_new_test_table;",
		session.NewAIMockContext().WithSQL("CREATE TEMPORARY TABLE tmp_test_table (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: ALTER TABLE RENAME non-temporary table", "ALTER TABLE test_table RENAME TO new_test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: CREATE TEMPORARY TABLE with prefix and additional columns", "CREATE TEMPORARY TABLE tmp_order_his(id BIGINT, name varchar(64) DEFAULT '');",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: ALTER TABLE RENAME TEMPORARY TABLE with new prefix", "ALTER TABLE tmp_order_his RENAME TO tmp_new_order_his;",
		session.NewAIMockContext().WithSQL("CREATE TEMPORARY TABLE tmp_order_his(id BIGINT, name varchar(64) DEFAULT '');"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: ALTER TABLE RENAME TEMPORARY TABLE without new prefix", "ALTER TABLE tmp_order_his RENAME TO new_order_his;",
		session.NewAIMockContext().WithSQL("CREATE TEMPORARY TABLE tmp_order_his(id BIGINT, name varchar(64) DEFAULT '');"),
		nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
