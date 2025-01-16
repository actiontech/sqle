package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00026(t *testing.T) {
	ruleName := ai.SQLE00026
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// case 1: CREATE TABLE 中 INT 类型字段指定了长度且不包含 zerofill
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with INT field with specified length but without zerofill", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id INT(11));
    `, newTestResult().addResult(ruleName))

	// Test case 2
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with TINYINT field with specified length but without zerofill", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id TINYINT(4));
    `, newTestResult().addResult(ruleName))

	// Test case 3
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with SMALLINT field with specified length but without zerofill", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id SMALLINT(6));
    `, newTestResult().addResult(ruleName))

	// Test case 4
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with MEDIUMINT field with specified length but without zerofill", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id MEDIUMINT(9));
    `, newTestResult().addResult(ruleName))

	// Test case 5
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with BIGINT field with specified length but without zerofill", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id BIGINT(20));
    `, newTestResult().addResult(ruleName))

	// Test case 6
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with INT field with specified length and zerofill", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id INT(11) ZEROFILL);
    `, newTestResult())

	// Test case 7
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with INT field not specified length", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id INT);
    `, newTestResult())

	// Test case 8
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to INT with specified length but without zerofill", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age INT(11);
    `, newTestResult().addResult(ruleName))

	// Test case 9
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to TINYINT with specified length but without zerofill", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age TINYINT(4);
    `, newTestResult().addResult(ruleName))

	// Test case 10
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to SMALLINT with specified length but without zerofill", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age SMALLINT(6);
    `, newTestResult().addResult(ruleName))

	// Test case 11
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to MEDIUMINT with specified length but without zerofill", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age MEDIUMINT(9);
    `, newTestResult().addResult(ruleName))

	// Test case 12
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to BIGINT with specified length but without zerofill", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age BIGINT(20);
    `, newTestResult().addResult(ruleName))

	// Test case 13
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to INT with specified length and zerofill", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age INT(11) ZEROFILL;
    `, newTestResult())

	// Test case 14
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to INT not specified length", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age INT;
    `, newTestResult())

	// Test case 15
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with INTEGER field with specified length but without zerofill", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id INTEGER(11));
    `, newTestResult().addResult(ruleName))

	// Test case 16
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with INTEGER field with specified length and zerofill", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id INTEGER(11) ZEROFILL);
    `, newTestResult())

	// Test case 17
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to INTEGER with specified length but without zerofill", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age INTEGER(11);
    `, newTestResult().addResult(ruleName))

	// Test case 18
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to INTEGER with specified length and zerofill", DefaultMysqlInspect(), `
    ALTER TABLE exist_db.exist_tb_1  ADD COLUMN age INTEGER(11) ZEROFILL;
    `, newTestResult())

	// Test case 19
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with INTEGER field not specified length", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id TINYINT);
    `, newTestResult())

	// Test case 20
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with SMALLINT field not specified length", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id SMALLINT);
    `, newTestResult())

	// Test case 21
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with MEDIUMINT field not specified length", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id MEDIUMINT);
    `, newTestResult())

	// Test case 22
	runSingleRuleInspectCase(rule, t, "CREATE TABLE with BIGINT field not specified length", DefaultMysqlInspect(), `
    CREATE TABLE exist_db.not_exist_tb_1  (id BIGINT);
    `, newTestResult())

	// Test case 23
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to INTEGER field with specified length and zerofill", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE id id2 INTEGER(11) ZEROFILL;
    `, newTestResult())

	// Test case 24
	runSingleRuleInspectCase(rule, t, "ALTER TABLE add field to INTEGER field with specified length but without zerofill", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE id id2 INTEGER(11);
    `, newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
