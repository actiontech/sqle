package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00220(t *testing.T) {
	ruleName := ai.SQLE00220
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	testCases := []struct {
		desc           string
		sqlToTest      string
		expectedResult string
		mockContextSql string
	}{
		{"case 0: SELECT count(id) without WHERE clause", "SELECT count(id) FROM users;", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 1: SELECT count(*) without WHERE clause", "SELECT count(*) FROM users;", "违规", "CREATE TABLE users (id INT, age INT);"},
		{"case 2: SELECT count(1) without WHERE clause", "SELECT count(1) FROM users;", "违规", "CREATE TABLE users (id INT, age INT);"},
		{"case 3: SELECT count(*) with WHERE clause", "SELECT count(*) FROM users WHERE age > 30;", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 4: SELECT count(1) with WHERE clause", "SELECT count(1) FROM users WHERE age > 30;", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 5: SELECT count(*) in subquery without WHERE clause", "SELECT * FROM (SELECT count(*) FROM users) AS sub;", "违规", "CREATE TABLE users (id INT, age INT);"},
		{"case 6: SELECT count(1) in subquery without WHERE clause", "SELECT * FROM (SELECT count(1) FROM users) AS sub;", "违规", "CREATE TABLE users (id INT, age INT);"},
		{"case 7: SELECT count(*) in subquery with WHERE clause", "SELECT * FROM (SELECT count(*) FROM users WHERE age > 30) AS sub;", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 8: SELECT count(1) in subquery with WHERE clause", "SELECT * FROM (SELECT count(1) FROM users WHERE age > 30) AS sub;", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 9: UPDATE with count(*) without WHERE clause", "UPDATE users SET age = (SELECT count(*) FROM users);", "违规", "CREATE TABLE users (id INT, age INT);"},
		{"case 10: UPDATE with count(1) without WHERE clause", "UPDATE users SET age = (SELECT count(1) FROM users);", "违规", "CREATE TABLE users (id INT, age INT);"},
		{"case 11: UPDATE with count(*) with WHERE clause", "UPDATE users SET age = (SELECT count(*) FROM users WHERE age > 30);", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 12: UPDATE with count(1) with WHERE clause", "UPDATE users SET age = (SELECT count(1) FROM users WHERE age > 30);", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 13: DELETE with count(*) without WHERE clause", "DELETE FROM users WHERE id = (SELECT count(*) FROM users);", "违规", "CREATE TABLE users (id INT, age INT);"},
		{"case 14: DELETE with count(1) without WHERE clause", "DELETE FROM users WHERE id = (SELECT count(1) FROM users);", "违规", "CREATE TABLE users (id INT, age INT);"},
		{"case 15: DELETE with count(*) with WHERE clause", "DELETE FROM users WHERE id = (SELECT count(*) FROM users WHERE age > 30);", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 16: DELETE with count(1) with WHERE clause", "DELETE FROM users WHERE id = (SELECT count(1) FROM users WHERE age > 30);", "通过", "CREATE TABLE users (id INT, age INT);"},
		{"case 17: SELECT count(*) without WHERE clause on customers table", "SELECT count(*) FROM customers;", "违规", "CREATE TABLE customers (id INT, age INT);"},
		{"case 18: SELECT count(1) without WHERE clause on customers table", "SELECT count(1) FROM customers;", "违规", "CREATE TABLE customers (id INT, age INT);"},
		{"case 19: SELECT count(*) with WHERE clause on customers table", "SELECT count(*) FROM customers WHERE age > 25;", "通过", "CREATE TABLE customers (id INT, age INT);"},
		{"case 20: SELECT count(1) with WHERE clause on customers table", "SELECT count(1) FROM customers WHERE age > 25;", "通过", "CREATE TABLE customers (id INT, age INT);"},
		{"case 21: UPDATE with count(*) without WHERE clause on customers table", "UPDATE customers SET age = (SELECT count(*) FROM customers);", "违规", "CREATE TABLE customers (id INT, age INT);"},
		{"case 22: UPDATE with count(1) without WHERE clause on customers table", "UPDATE customers SET age = (SELECT count(1) FROM customers);", "违规", "CREATE TABLE customers (id INT, age INT);"},
		{"case 23: UPDATE with count(*) with WHERE clause on customers table", "UPDATE customers SET age = (SELECT count(*) FROM customers WHERE age > 25);", "通过", "CREATE TABLE customers (id INT, age INT);"},
		{"case 24: UPDATE with count(1) with WHERE clause on customers table", "UPDATE customers SET age = (SELECT count(1) FROM customers WHERE age > 25);", "通过", "CREATE TABLE customers (id INT, age INT);"},
		{"case 25: DELETE with count(*) without WHERE clause on customers table", "DELETE FROM customers WHERE id = (SELECT count(*) FROM customers);", "违规", "CREATE TABLE customers (id INT, age INT);"},
		{"case 26: DELETE with count(1) without WHERE clause on customers table", "DELETE FROM customers WHERE id = (SELECT count(1) FROM customers);", "违规", "CREATE TABLE customers (id INT, age INT);"},
		{"case 27: DELETE with count(*) with WHERE clause on customers table", "DELETE FROM customers WHERE id = (SELECT count(*) FROM customers WHERE age > 25);", "通过", "CREATE TABLE customers (id INT, age INT);"},
		{"case 28: DELETE with count(1) with WHERE clause on customers table", "DELETE FROM customers WHERE id = (SELECT count(1) FROM customers WHERE age > 25);", "通过", "CREATE TABLE customers (id INT, age INT);"},
	}

	for _, tc := range testCases {
		if tc.expectedResult == "违规" {
			runAIRuleCase(rule, t, tc.desc, tc.sqlToTest,
				session.NewAIMockContext().WithSQL(tc.mockContextSql),
				nil, newTestResult().addResult(ruleName))
		} else {
			runAIRuleCase(rule, t, tc.desc, tc.sqlToTest,
				session.NewAIMockContext().WithSQL(tc.mockContextSql),
				nil, newTestResult())
		}
	}
}

// ==== Rule test code end ====
