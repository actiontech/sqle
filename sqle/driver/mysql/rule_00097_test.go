package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00097(t *testing.T) {
	ruleName := ai.SQLE00097
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT语句中ORDER BY使用VARCHAR字段长度超过100",
		"SELECT name FROM users ORDER BY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='users' and column_name='name';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 255),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT语句中ORDER BY使用VARCHAR字段长度未超过100",
		"SELECT name FROM users ORDER BY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='users' and column_name='name';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 50),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 3: SELECT语句中ORDER BY使用TEXT字段",
		"SELECT description FROM products ORDER BY description;",
		session.NewAIMockContext().WithSQL("CREATE TABLE products (description TINYTEXT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='products' and column_name='description';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("TINYTEXT", 65535),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: SELECT语句中GROUP BY使用BLOB字段",
		"SELECT image FROM gallery GROUP BY image;",
		session.NewAIMockContext().WithSQL("CREATE TABLE gallery (image BLOB);"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='gallery' and column_name='image';",
				Rows:  sqlmock.NewRows([]string{"data_type"}).AddRow("blob"),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: INSERT...SELECT语句中ORDER BY使用VARCHAR字段长度超过100",
		"INSERT INTO archive SELECT name FROM users ORDER BY name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (name VARCHAR(255)); CREATE TABLE archive (name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='users' and column_name='name';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 255),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: UNION语句中一个SELECT使用VARCHAR字段长度没超过100",
		"SELECT name FROM users1 UNION (SELECT name FROM users2 ORDER BY name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users1 (name VARCHAR(255)); CREATE TABLE users2 (name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='users1' and column_name='name';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 255),
			},
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='users2' and column_name='name';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 50),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 7: DISTINCT语句中使用VARCHAR字段长度超过100",
		"SELECT DISTINCT name FROM users;",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='users' and column_name='name';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 255),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: SELECT语句中ORDER BY使用VARCHAR字段长度超过512(从xml中补充)",
		"SELECT mark1 FROM customers ORDER BY mark1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (mark1 VARCHAR(2000));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='customers' and column_name='mark1';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 2000),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: SELECT语句中ORDER BY使用VARCHAR字段长度未超过512(从xml中补充)",
		"SELECT mark2 FROM customers ORDER BY mark2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (mark2 VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='customers' and column_name='mark2';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 100),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 10: INSERT...SELECT union ...语句中ORDER BY使用VARCHAR字段长度超过100",
		"INSERT INTO archive SELECT name FROM users union (SELECT name FROM archive ORDER BY name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE users (name VARCHAR(255)); CREATE TABLE archive (name VARCHAR(255));"),
		[]*AIMockSQLExpectation{
			{
				Query: "select data_type, character_maximum_length from information_schema.columns where table_name='users' and column_name='name';",
				Rows:  sqlmock.NewRows([]string{"data_type", "character_maximum_length"}).AddRow("VARCHAR", 255),
			},
		}, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
