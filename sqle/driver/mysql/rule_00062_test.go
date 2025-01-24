package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQL00062(t *testing.T) {
	ruleName := ai.SQLE00062
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SET TRANSACTION ISOLATION LEVEL set to READ UNCOMMITTED",
		"SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SET TRANSACTION ISOLATION LEVEL set to READ COMMITTED",
		"SET TRANSACTION ISOLATION LEVEL READ COMMITTED;",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 3: SET TRANSACTION ISOLATION LEVEL set to REPEATABLE READ",
		"SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: SET TRANSACTION ISOLATION LEVEL set to SERIALIZABLE",
		"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: SET SESSION TRANSACTION ISOLATION LEVEL set to READ COMMITTED",
		"SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED;",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 6: SET GLOBAL TRANSACTION ISOLATION LEVEL set to READ UNCOMMITTED",
		"SET GLOBAL TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: SET @@SESSION.TRansaction_isolation set to READ COMMITTED",
		"SET @@SESSION.transaction_isolation = 'READ-COMMITTED';",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 8: SET @@GLOBAL.transaction_isolation set to REPEATABLE READ",
		"SET @@GLOBAL.transaction_isolation = 'REPEATABLE-READ';",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: SET TRANSACTION ISOLATION LEVEL set to READ COMMITTED",
		"SET TRANSACTION ISOLATION LEVEL ReAd CoMmItTeD;",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 11: SET GLOBAL transaction_isolation set to READ COMMITTED",
		"SET GLOBAL transaction_isolation = 'READ-COMMITTED';",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 12: SET SESSION transaction_isolation set to READ UNCOMMITTED",
		"SET SESSION transaction_isolation = 'READ-UNCOMMITTED';",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SET @@GLOBAL.transaction_isolation set to READ UNCOMMITTED",
		"SET @@GLOBAL.transaction_isolation = 'READ-UNCOMMITTED';",
		nil,
		nil,
		newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
