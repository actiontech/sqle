package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// "testing"

// rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
// "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"

// ==== Rule test code start ====
func TestRuleSQL00065(t *testing.T) {
	ruleName := ai.SQLE00065
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: ALTER TABLE MODIFY column and specify FIRST",
		"ALTER TABLE test_table MODIFY COLUMN col1 INT FIRST;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col2 INT, col1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: ALTER TABLE MODIFY column and specify AFTER",
		"ALTER TABLE test_table MODIFY COLUMN col1 INT AFTER col2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col2 INT, col1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: ALTER TABLE MODIFY column and not specify",
		"ALTER TABLE test_table MODIFY COLUMN col1 INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col1 INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: ALTER TABLE CHANGE column and specify FIRST",
		"ALTER TABLE test_table CHANGE COLUMN col1 col1_new INT FIRST;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col2 INT, col1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: ALTER TABLE CHANGE column and specify AFTER",
		"ALTER TABLE test_table CHANGE COLUMN col1 col1_new INT AFTER col2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col2 INT, col1 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: ALTER TABLE CHANGE column and not specify",
		"ALTER TABLE test_table CHANGE COLUMN col1 col1_new INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col1 INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: ALTER TABLE ADD column and specify FIRST",
		"ALTER TABLE test_table ADD COLUMN col1 INT FIRST;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col2 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: ALTER TABLE ADD column and specify AFTER",
		"ALTER TABLE test_table ADD COLUMN col1 INT AFTER col2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col2 INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: ALTER TABLE ADD column and not specify",
		"ALTER TABLE test_table ADD COLUMN col1 INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (col2 INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: ALTER TABLE CHANGE age column and specify FIRST",
		"ALTER TABLE customers CHANGE age age INT NOT NULL FIRST;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: ALTER TABLE CHANGE age column and specify AFTER",
		"ALTER TABLE customers CHANGE age age INT NOT NULL AFTER name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: ALTER TABLE CHANGE age column and not specify",
		"ALTER TABLE customers CHANGE age age INT NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 13: ALTER TABLE MODIFY age column and specify FIRST",
		"ALTER TABLE customers MODIFY age INT NOT NULL FIRST;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: ALTER TABLE MODIFY age column and specify AFTER",
		"ALTER TABLE customers MODIFY age INT NOT NULL AFTER name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: ALTER TABLE MODIFY age column and not specify",
		"ALTER TABLE customers MODIFY age INT NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 16: ALTER TABLE ADD addr column and specify FIRST",
		"ALTER TABLE customers ADD COLUMN addr VARCHAR(2000) NULL FIRST;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: ALTER TABLE ADD addr column and specify AFTER",
		"ALTER TABLE customers ADD COLUMN addr VARCHAR(2000) NULL AFTER name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: ALTER TABLE ADD addr column and not specify",
		"ALTER TABLE customers ADD COLUMN addr VARCHAR(2000) NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (name VARCHAR(255), age INT);"),
		nil, newTestResult())
}

// ==== Rule test code end ====
