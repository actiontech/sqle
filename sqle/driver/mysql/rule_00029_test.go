package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQLE00029(t *testing.T) {
	ruleName := ai.SQLE00029
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	for _, sql := range []string{
		`CREATE PROCEDURE my_procedure() BEGIN SELECT * FROM exist_db.exist_tb_1 ; END;`,
		`CREATE DEFINER='sqle_op'@'localhost' PROCEDURE my_procedure() BEGIN SELECT * FROM exist_db.exist_tb_1 ; END`,
		`CREATE DEFINER='sqle_op'@'localhost' PROCEDURE my_procedure() BEGIN END;`,
		`ALTER PROCEDURE my_procedure COMMENT 'Updated procedure'`,
	} {
		runSingleRuleInspectCase(rule, t, "",
			DefaultMysqlInspect(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").
				addResult(ruleName))
	}

	// ===== 语法错误
	for _, sql := range []string{
		`CREATE DEFINER='sqle_op'@'localhost'PROCEDURE my_procedure() BEGIN SELECT * FROM exist_db.exist_tb_1 ; END;`,
		`CREATEDEFINER='sqle_op'@'localhost' PROCEDURE my_procedure() BEGIN SELECT * FROM exist_db.exist_tb_1 ; END;`,
		`ALTER DEFINER='sqle_op'@'localhost' PROCEDURE my_procedure COMMENT 'Updated procedure'`,
	} {
		runSingleRuleInspectCase(rule, t, "",
			DefaultMysqlInspect(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))
	}
}

// ==== Rule test code end ====
