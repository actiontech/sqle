package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00016(t *testing.T) {
	ruleName := ai.SQLE00016
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE包含BLOB类型字段并设置为NOT NULL，违反规则SQLE00016",
		"CREATE TABLE sample_create_table (id INT, data BLOB NOT NULL);",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 2: CREATE TABLE包含BLOB类型字段并设置为NULL，符合规则SQLE00016",
		"CREATE TABLE sample_create_table (id INT, data BLOB NULL);",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 3: CREATE TABLE包含TEXT类型字段并设置为NOT NULL，违反规则SQLE00016",
		"CREATE TABLE sample_create_table (id INT, description TEXT NOT NULL);",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 4: CREATE TABLE包含INT类型字段并设置为NOT NULL，符合规则SQLE00016",
		"CREATE TABLE sample_create_table (id INT NOT NULL, name VARCHAR(100));",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 5: ALTER TABLE添加TEXT类型字段并设置为NOT NULL，违反规则SQLE00016",
		"ALTER TABLE sample_alter_table ADD COLUMN description TEXT NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_alter_table (id INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 6: ALTER TABLE添加JSON类型字段并设置为NULL，符合规则SQLE00016",
		"ALTER TABLE sample_alter_table ADD COLUMN config JSON NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_alter_table (id INT);"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 7: ALTER TABLE修改GEOMETRY类型字段并设置为NOT NULL，违反规则SQLE00016",
		"ALTER TABLE sample_alter_table MODIFY COLUMN location GEOMETRY NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_alter_table (id INT, location GEOMETRY NULL);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 8: ALTER TABLE添加INT类型字段并设置为NOT NULL，符合规则SQLE00016",
		"ALTER TABLE sample_alter_table ADD COLUMN quantity INT NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_alter_table (id INT);"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 9: CREATE TABLE包含TEXT类型字段并设置为NOT NULL，违反规则SQLE00016(从xml中补充)",
		"CREATE TABLE tb_text (id INT NOT NULL AUTO_INCREMENT, a TINYTEXT, b TEXT NOT NULL, c VARCHAR(255) DEFAULT NULL, d BLOB NOT NULL, PRIMARY KEY (id));",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 10: ALTER TABLE添加BLOB类型字段并设置为NOT NULL，违反规则SQLE00016(从xml中补充)",
		"ALTER TABLE tb_text ADD COLUMN d2 BLOB NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE tb_text (id INT NOT NULL AUTO_INCREMENT, a TINYTEXT, b TEXT, c VARCHAR(255) DEFAULT NULL, d BLOB, PRIMARY KEY (id));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 11: ALTER TABLE修改BLOB类型字段并设置为NOT NULL，违反规则SQLE00016(从xml中补充)",
		"ALTER TABLE tb_text CHANGE COLUMN c c BLOB NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE tb_text (id INT NOT NULL AUTO_INCREMENT, a TINYTEXT, b TEXT, c VARCHAR(255) DEFAULT NULL, d BLOB, PRIMARY KEY (id));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 12: ALTER TABLE修改GEOMETRY类型字段并设置为NOT NULL，违反规则SQLE00016(从xml中补充)",
		"ALTER TABLE tb_text MODIFY COLUMN c GEOMETRY NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE tb_text (id INT NOT NULL AUTO_INCREMENT, a TINYTEXT, b TEXT, c VARCHAR(255) DEFAULT NULL, d BLOB, PRIMARY KEY (id));"),
		nil,
		newTestResult().addResult(ruleName),
	)
}

// ==== Rule test code end ====
