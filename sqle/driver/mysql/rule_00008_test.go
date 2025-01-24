package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

// ==== Rule test code start ====
func TestRuleSQLE00008(t *testing.T) {
	ruleName := ai.SQLE00008
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// case 1: CREATE TABLE 定义包含列级别的主键
	runAIRuleCase(rule, t, "case 1: CREATE TABLE 定义包含列级别的主键",
		"CREATE TABLE users (user_id INT PRIMARY KEY, username VARCHAR(50));",
		nil,
		nil,
		newTestResult(),
	)

	// case 2: CREATE TABLE 定义包含表级别的主键
	runAIRuleCase(rule, t, "case 2: CREATE TABLE 定义包含表级别的主键",
		"CREATE TABLE orders (order_id INT, customer_id INT, PRIMARY KEY(order_id));",
		nil,
		nil,
		newTestResult(),
	)

	// case 3: CREATE TABLE 未包含主键定义
	runAIRuleCase(rule, t, "case 3: CREATE TABLE 未包含主键定义",
		"CREATE TABLE products (product_id INT, product_name VARCHAR(50));",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	// case 4: ALTER TABLE 添加主键定义
	runAIRuleCase(rule, t, "case 4: ALTER TABLE 添加主键定义",
		"ALTER TABLE products ADD PRIMARY KEY(product_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE products (product_id INT, product_name VARCHAR(50));"),
		nil,
		newTestResult(),
	)

	// case 5: ALTER TABLE 删除主键并添加新的主键定义
	runAIRuleCase(rule, t, "case 5: ALTER TABLE 删除主键并添加新的主键定义",
		"ALTER TABLE users DROP PRIMARY KEY, ADD PRIMARY KEY(username);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (user_id INT, username VARCHAR(50));"),
		nil,
		newTestResult().add(driverV2.RuleLevelError, "", "当前没有主键，不能执行删除"),
	)

	// case 6: ALTER TABLE 删除主键但未添加新的主键定义
	runAIRuleCase(rule, t, "case 6: ALTER TABLE 删除主键但未添加新的主键定义",
		"ALTER TABLE orders DROP PRIMARY KEY;",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT PRIMARY KEY, customer_id INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	// case 8: CREATE TABLE 无主键定义导致全表扫描(从xml中补充)
	runAIRuleCase(rule, t, "case 8: CREATE TABLE 无主键定义导致全表扫描(从xml中补充)",
		"CREATE TABLE no_primary_key (id INT, name VARCHAR(50));",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	// case 9: ALTER TABLE 为无主键表添加主键(从xml中补充)
	runAIRuleCase(rule, t, "case 9: ALTER TABLE 为无主键表添加主键(从xml中补充)",
		"ALTER TABLE no_primary_key ADD PRIMARY KEY(id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE no_primary_key (id INT, name VARCHAR(50));"),
		nil,
		newTestResult(),
	)
}

// ==== Rule test code end ====
