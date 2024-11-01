package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00019(t *testing.T) {
	ruleName := ai.SQLE00019
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE 使用 ENUM 类型字段",
		"CREATE TABLE sample_table (id INT, status ENUM('active', 'inactive'));",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: CREATE TABLE 使用 SET 类型字段",
		"CREATE TABLE sample_table (id INT, tags SET('tag1', 'tag2', 'tag3'));",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: CREATE TABLE 使用标准数据类型字段",
		"CREATE TABLE sample_table (id INT, name VARCHAR(255));",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: CREATE TABLE 包含多个字段，其中一个字段使用 ENUM 类型",
		"CREATE TABLE sample_table (id INT, name VARCHAR(255), status ENUM('active', 'inactive'));",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: ALTER TABLE 新增字段使用 ENUM 类型",
		"ALTER TABLE sample_table ADD COLUMN role ENUM('admin', 'user');",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(255));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: ALTER TABLE 新增字段使用 SET 类型",
		"ALTER TABLE sample_table ADD COLUMN permissions SET('read', 'write', 'execute');",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(255));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: ALTER TABLE 新增字段使用标准数据类型",
		"ALTER TABLE sample_table ADD COLUMN age INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(255));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: ALTER TABLE 修改字段为 ENUM 类型",
		"ALTER TABLE sample_table MODIFY COLUMN status ENUM('active', 'inactive', 'pending');",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, status VARCHAR(255));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: ALTER TABLE 修改字段为标准数据类型",
		"ALTER TABLE sample_table MODIFY COLUMN name TEXT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, name VARCHAR(255));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: ALTER TABLE 不涉及 ENUM 或 SET 类型的字段变更",
		"ALTER TABLE sample_table DROP COLUMN age;",
		session.NewAIMockContext().WithSQL("CREATE TABLE sample_table (id INT, age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 11: CREATE TABLE 使用 ENUM 类型字段，测试 ENUM 的值插入和排序问题",
		"CREATE TABLE t1 (a INT PRIMARY KEY AUTO_INCREMENT, b ENUM('A','3','2','1') DEFAULT '3');",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
