package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00016_1(t *testing.T) {
	ruleName := ai.SQLE00016_1
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// create table, with blob/text column, with NOT NULL
	runSingleRuleInspectCase(rule, t, "create table, with expected column type, with NOT NULL", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
a blob NOT NULL,
b text NOT NULL,
PRIMARY KEY (id)
);
`, newTestResult().addResult(ruleName, "a,b"))

}

// ==== Rule test code end ====
