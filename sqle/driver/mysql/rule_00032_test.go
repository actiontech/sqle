package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00032(t *testing.T) {
	ruleName := ai.SQLE00032
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runSingleRuleInspectCase(rule, t, "Create database name without the '_DB' fixed suffix", DefaultMysqlInspect(), `
    CREATE DATABASE no_exist;
    `, newTestResult().addResult(ruleName, "_DB"))

	runSingleRuleInspectCase(rule, t, "Create database name with the fixed suffix '_DB'", DefaultMysqlInspect(), `
    CREATE DATABASE no_exist_DB;
    `, newTestResult())

	runSingleRuleInspectCase(rule, t, "Create database name without the '_DB' fixed suffix", DefaultMysqlInspect(), `
    CREATE DATABASE no_exist_db;
    `, newTestResult().addResult(ruleName, "_DB"))
}

// ==== Rule test code end ====
