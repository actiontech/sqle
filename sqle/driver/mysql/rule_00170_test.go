package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00170(t *testing.T) {
	ruleName := ai.SQLE00170
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	// case 1
	runAIRuleCase(rule, t, "case 1: 使用 ALTER TABLE ... MODIFY ... 语句缩短 VARCHAR 字段长度，且当前数据长度超过新长度",
		"ALTER TABLE test_table MODIFY name VARCHAR(50);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(100));").WithSQL("INSERT INTO test_table (id, name) VALUES (1, 'This is a very long name exceeding fifty characters.');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(name)) "max_length" FROM test_table`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(60),
			},
		}, newTestResult().addResult(ruleName))

	// case 2
	runAIRuleCase(rule, t, "case 2: 使用 ALTER TABLE ... MODIFY ... 语句不缩短 VARCHAR 字段长度，保持或增加长度",
		"ALTER TABLE test_table MODIFY name VARCHAR(100);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));").WithSQL("INSERT INTO test_table (id, name) VALUES (1, 'Short name');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(name)) "max_length" FROM test_table`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(10),
			},
		}, newTestResult())

	// case 3
	runAIRuleCase(rule, t, "case 3: 使用 ALTER TABLE ... CHANGE ... 语句缩短 VARCHAR 字段长度，且当前数据长度超过新长度",
		"ALTER TABLE test_table CHANGE name name VARCHAR(30);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(100));").WithSQL("INSERT INTO test_table (id, name) VALUES (1, 'This name is definitely longer than thirty characters.');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(name)) "max_length" FROM test_table`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(45),
			},
		}, newTestResult().addResult(ruleName))

	// case 4
	runAIRuleCase(rule, t, "case 4: 使用 ALTER TABLE ... CHANGE ... 语句不缩短 VARCHAR 字段长度，保持或增加长度",
		"ALTER TABLE test_table CHANGE name name VARCHAR(100);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, name VARCHAR(50));").WithSQL("INSERT INTO test_table (id, name) VALUES (1, 'Short name');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(name)) "max_length" FROM test_table`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(10),
			},
		}, newTestResult())

	// case 5
	runAIRuleCase(rule, t, "case 5: 使用 ALTER TABLE ... MODIFY ... 语句缩短 CHAR 字段长度，且当前数据长度超过新长度",
		"ALTER TABLE test_table MODIFY code CHAR(5);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, code CHAR(10));").WithSQL("INSERT INTO test_table (id, code) VALUES (1, '1234567890');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(code)) "max_length" FROM test_table`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(10),
			},
		}, newTestResult().addResult(ruleName))

	// case 6
	runAIRuleCase(rule, t, "case 6: 使用 ALTER TABLE ... MODIFY ... 语句不缩短 CHAR 字段长度，保持或增加长度",
		"ALTER TABLE test_table MODIFY code CHAR(10);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, code CHAR(5));").WithSQL("INSERT INTO test_table (id, code) VALUES (1, '12345');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(code)) "max_length" FROM test_table`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(5),
			},
		}, newTestResult())

	// case 7
	runAIRuleCase(rule, t, "case 7: 使用 ALTER TABLE ... CHANGE ... 语句缩短 CHAR 字段长度，且当前数据长度超过新长度",
		"ALTER TABLE test_table CHANGE code code CHAR(4);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, code CHAR(10));").WithSQL("INSERT INTO test_table (id, code) VALUES (1, '1234567890');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(code)) "max_length" FROM test_table`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(10),
			},
		}, newTestResult().addResult(ruleName))

	// case 8
	runAIRuleCase(rule, t, "case 8: 使用 ALTER TABLE ... CHANGE ... 语句不缩短 CHAR 字段长度，保持或增加长度",
		"ALTER TABLE test_table CHANGE code code CHAR(10);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT, code CHAR(5));").WithSQL("INSERT INTO test_table (id, code) VALUES (1, '1234');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(code)) "max_length" FROM test_table`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(4),
			},
		}, newTestResult())

	// case 9
	runAIRuleCase(rule, t, "case 9: 使用 ALTER TABLE ... MODIFY ... 语句缩短 VARCHAR 字段长度，且当前数据长度超过新长度 (从xml中补充)",
		"ALTER TABLE customers MODIFY city VARCHAR(4);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, city VARCHAR(10));").WithSQL("INSERT INTO customers (id, city) VALUES (1, 'New York');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(city)) "max_length" FROM customers`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(8),
			},
		}, newTestResult().addResult(ruleName))

	// case 10
	runAIRuleCase(rule, t, "case 10: 使用 ALTER TABLE ... MODIFY ... 语句不缩短 VARCHAR 字段长度，保持或增加长度 (从xml中补充)",
		"ALTER TABLE customers MODIFY city VARCHAR(10);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, city VARCHAR(5));").WithSQL("INSERT INTO customers (id, city) VALUES (1, 'LA');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(city)) "max_length" FROM customers`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(2),
			},
		}, newTestResult())

	// case 11
	runAIRuleCase(rule, t, "case 11: 使用 ALTER TABLE ... CHANGE ... 语句缩短 VARCHAR 字段长度，且当前数据长度超过新长度 (从xml中补充)",
		"ALTER TABLE customers CHANGE city city VARCHAR(4);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, city VARCHAR(10));").WithSQL("INSERT INTO customers (id, city) VALUES (1, 'Chicago');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(city)) "max_length" FROM customers`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(7),
			},
		}, newTestResult().addResult(ruleName))

	// case 12
	runAIRuleCase(rule, t, "case 12: 使用 ALTER TABLE ... CHANGE ... 语句不缩短 VARCHAR 字段长度，保持或增加长度 (从xml中补充)",
		"ALTER TABLE customers CHANGE city city VARCHAR(10);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT, city VARCHAR(5));").WithSQL("INSERT INTO customers (id, city) VALUES (1, 'SF');"),
		[]*AIMockSQLExpectation{
			{
				Query: `SELECT MAX(CHAR_LENGTH(city)) "max_length" FROM customers`,
				Rows:  sqlmock.NewRows([]string{"max_length"}).AddRow(2),
			},
		}, newTestResult())
}

// ==== Rule test code end ====
