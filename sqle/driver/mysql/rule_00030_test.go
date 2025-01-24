package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQLE00030(t *testing.T) {
	ruleName := ai.SQLE00030
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	for _, sql := range []string{
		`CREATE TRIGGER my_trigger BEFORE INSERT ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.name = UPPER(NEW.name);`,
		`CREATE DEFINER='sqle_op'@'localhost' TRIGGER my_trigger BEFORE INSERT ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.name = UPPER(NEW.name);`,
		`CREATE TRIGGER my_trigger AFTER UPDATE ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.updated_at = NOW();`,
		`CREATE TRIGGER ins_check BEFORE INSERT ON customers FOR EACH ROW BEGIN IF NEW.age < 18 THEN SET NEW.mark1 = '未满18岁'; ELSEIF NEW.age >= 18 THEN SET NEW.mark1 = '满18岁，已经成年了'; END IF; END;`,
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(ruleName))
	}

	for _, sql := range []string{
		`DROP TRIGGER IF EXISTS ins_check;`,
		`CREATEDEFINER='sqle_op'@'localhost' TRIGGER my_trigger BEFORE INSERT ON exist_db.exist_tb_1 FOR EACH ROW SET NEW.name = UPPER(NEW.name);`,
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))
	}
}

// ==== Rule test code end ====
