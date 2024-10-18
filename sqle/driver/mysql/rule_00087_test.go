package mysql

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func generateLongInClauseFor00087(count int) string {
	var ids []string
	for i := 1; i <= count; i++ {
		ids = append(ids, strconv.Itoa(i))
	}
	return strings.Join(ids, ", ")
}

func TestRuleSQLE00087(t *testing.T) {
	ruleName := ai.SQLE00087
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT语句中WHERE条件的IN列表超过500个元素",
		fmt.Sprintf("SELECT * FROM users WHERE id IN (%s);", generateLongInClauseFor00087(501)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT语句中WHERE条件的IN列表少于500个元素",
		fmt.Sprintf("SELECT * FROM users WHERE id IN (%s);", generateLongInClauseFor00087(499)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: INSERT ... SELECT语句中WHERE条件的IN列表超过500个元素",
		fmt.Sprintf("INSERT INTO archive_users SELECT * FROM users WHERE id IN (%s);", generateLongInClauseFor00087(501)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100)); CREATE TABLE archive_users (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: INSERT ... SELECT语句中WHERE条件的IN列表少于500个元素",
		fmt.Sprintf("INSERT INTO archive_users SELECT * FROM users WHERE id IN (%s);", generateLongInClauseFor00087(499)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100)); CREATE TABLE archive_users (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: UPDATE语句中WHERE条件的IN列表超过500个元素",
		fmt.Sprintf("UPDATE users SET status = 'active' WHERE id IN (%s);", generateLongInClauseFor00087(501)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), status VARCHAR(20));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: UPDATE语句中WHERE条件的IN列表少于500个元素",
		fmt.Sprintf("UPDATE users SET status = 'active' WHERE id IN (%s);", generateLongInClauseFor00087(499)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), status VARCHAR(20));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: DELETE语句中WHERE条件的IN列表超过500个元素",
		fmt.Sprintf("DELETE FROM users WHERE id IN (%s);", generateLongInClauseFor00087(501)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: DELETE语句中WHERE条件的IN列表少于500个元素",
		fmt.Sprintf("DELETE FROM users WHERE id IN (%s);", generateLongInClauseFor00087(499)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: UNION语句中第一个SELECT的WHERE条件的IN列表超过500个元素",
		fmt.Sprintf("SELECT * FROM users WHERE id IN (%s) UNION SELECT * FROM admins;", generateLongInClauseFor00087(501)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100)); CREATE TABLE admins (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: UNION语句中第二个SELECT的WHERE条件的IN列表超过500个元素",
		fmt.Sprintf("SELECT * FROM users UNION SELECT * FROM admins WHERE id IN (%s);", generateLongInClauseFor00087(501)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100)); CREATE TABLE admins (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: UNION语句中所有SELECT的WHERE条件的IN列表少于500个元素",
		fmt.Sprintf("SELECT * FROM users WHERE id IN (%s) UNION SELECT * FROM admins WHERE id IN (%s);", generateLongInClauseFor00087(499), generateLongInClauseFor00087(499)),
		session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100)); CREATE TABLE admins (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil, newTestResult())

	// 不支持 with语法
	// // runAIRuleCase(rule, t, "case 12: WITH语句中子查询的WHERE条件的IN列表超过500个元素",
	// // 	fmt.Sprintf("WITH temp AS (SELECT * FROM users WHERE id IN (%s)) SELECT * FROM temp;", generateLongInClauseFor00087(501)),
	// // 	session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));"),
	// // 	nil, newTestResult().addResult(ruleName))

	// // runAIRuleCase(rule, t, "case 13: WITH语句中子查询的WHERE条件的IN列表少于500个元素",
	// // 	fmt.Sprintf("WITH temp AS (SELECT * FROM users WHERE id IN (%s)) SELECT * FROM temp;", generateLongInClauseFor00087(499)),
	// // 	session.NewAIMockContext().WithSQL("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));"),
	// // 	nil, newTestResult())

	runAIRuleCase(rule, t, "case 14: SELECT语句中WHERE条件的IN子查询扫描行数超过500个(从xml中补充)",
		"SELECT * FROM customers WHERE id IN (SELECT id FROM customers_ids1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(100)); CREATE TABLE customers_ids1 (id INT PRIMARY KEY);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id FROM customers_ids1",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(1000),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SELECT语句中WHERE条件的IN子查询扫描行数少于500个(从xml中补充)",
		"SELECT * FROM customers WHERE id IN (SELECT id FROM customers_ids1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(100)); CREATE TABLE customers_ids1 (id INT PRIMARY KEY);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id FROM customers_ids1",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(488),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult())

	runAIRuleCase(rule, t, "case 15: delete语句中WHERE条件的NOT IN子查询扫描行数少于500个(从xml中补充)",
		"delete from customers WHERE id NOT IN (SELECT id FROM customers_ids1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(100)); CREATE TABLE customers_ids1 (id INT PRIMARY KEY);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT id FROM customers_ids1",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(600),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		},
		newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
