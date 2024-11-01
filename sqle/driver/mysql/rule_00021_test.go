package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00021(t *testing.T) {
	ruleName := ai.SQLE00021
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// CREATE语句测试用例
	runAIRuleCase(rule, t, "case 1: 创建表employees，所有字段包含NOT NULL约束",
		"CREATE TABLE employees (id INT NOT NULL, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: 创建表employees，字段id缺少NOT NULL约束",
		"CREATE TABLE employees (id INT, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: 创建表employees，字段id为主键且包含NOT NULL约束",
		"CREATE TABLE employees (id INT NOT NULL PRIMARY KEY, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: 创建表employees，主键在表级定义，id字段未显式指定NOT NULL",
		"CREATE TABLE employees (id INT, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL, PRIMARY KEY (id));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: 创建表employees，主键定义在多个字段上，salary字段缺少NOT NULL约束",
		"CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), PRIMARY KEY (id, name));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: 创建表employees，所有主键字段通过主键约束定义且其他字段包含NOT NULL约束",
		"CREATE TABLE employees (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);",
		nil, nil, newTestResult().addResult(ruleName))

	// ALTER语句测试用例
	runAIRuleCase(rule, t, "case 7: 修改表employees，添加字段age并包含NOT NULL约束",
		"ALTER TABLE employees ADD COLUMN age INT NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT NOT NULL, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: 修改表employees，添加字段age但缺少NOT NULL约束",
		"ALTER TABLE employees ADD COLUMN age INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT NOT NULL, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: 修改表employees，修改字段name以包含NOT NULL约束",
		"ALTER TABLE employees MODIFY COLUMN name VARCHAR(100) NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT NOT NULL, name VARCHAR(100), salary DECIMAL(10,2) NOT NULL);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: 修改表employees，修改字段name以移除NOT NULL约束",
		"ALTER TABLE employees MODIFY COLUMN name VARCHAR(100);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT NOT NULL, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: 修改表employees，修改字段name且保持NOT NULL约束不变",
		"ALTER TABLE employees MODIFY COLUMN name VARCHAR(100) NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT NOT NULL, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 12: 修改表employees，同时添加字段department_id并包含NOT NULL约束，同时删除字段salary",
		"ALTER TABLE employees ADD COLUMN department_id INT NOT NULL, DROP COLUMN salary;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT NOT NULL, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 13: 修改表employees，同时添加字段department_id但缺少NOT NULL约束，同时删除字段salary",
		"ALTER TABLE employees ADD COLUMN department_id INT, DROP COLUMN salary;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT NOT NULL, name VARCHAR(100) NOT NULL, salary DECIMAL(10,2) NOT NULL);"),
		nil, newTestResult().addResult(ruleName))

	// 新增示例测试用例
	runAIRuleCase(rule, t, "case 14: 创建表employees，所有字段缺少NOT NULL约束(从xml中补充)",
		"CREATE TABLE employees (id INT, name VARCHAR(50), salary DECIMAL(10,2));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: 创建表employees，所有字段包含NOT NULL约束(从xml中补充)",
		"CREATE TABLE employees (id INT NOT NULL, name VARCHAR(50) NOT NULL, salary DECIMAL(10,2) NOT NULL);",
		nil, nil, newTestResult())
}

// ==== Rule test code end ====
