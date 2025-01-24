package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00080(t *testing.T) {
	ruleName := ai.SQLE00080
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: INSERT...VALUES 语句超过阈值的行数",
		"INSERT INTO test_table (id) VALUES (1),(2),(3),(4),(5),(6),(7),(8),(9),(10),(11),(12),(13),(14),(15),(16),(17),(18),(19),(20),(21),(22),(23),(24),(25),(26),(27),(28),(29),(30),(31),(32),(33),(34),(35),(36),(37),(38),(39),(40),(41),(42),(43),(44),(45),(46),(47),(48),(49),(50),(51),(52),(53),(54),(55),(56),(57),(58),(59),(60),(61),(62),(63),(64),(65),(66),(67),(68),(69),(70),(71),(72),(73),(74),(75),(76),(77),(78),(79),(80),(81),(82),(83),(84),(85),(86),(87),(88),(89),(90),(91),(92),(93),(94),(95),(96),(97),(98),(99),(100), (101);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: INSERT...VALUES 语句未超过阈值的行数",
		"INSERT INTO test_table (id) VALUES (1), (2), (3), (100);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: REPLACE...VALUES 语句超过阈值的行数",
		"REPLACE INTO test_table (id) VALUES (1),(2),(3),(4),(5),(6),(7),(8),(9),(10),(11),(12),(13),(14),(15),(16),(17),(18),(19),(20),(21),(22),(23),(24),(25),(26),(27),(28),(29),(30),(31),(32),(33),(34),(35),(36),(37),(38),(39),(40),(41),(42),(43),(44),(45),(46),(47),(48),(49),(50),(51),(52),(53),(54),(55),(56),(57),(58),(59),(60),(61),(62),(63),(64),(65),(66),(67),(68),(69),(70),(71),(72),(73),(74),(75),(76),(77),(78),(79),(80),(81),(82),(83),(84),(85),(86),(87),(88),(89),(90),(91),(92),(93),(94),(95),(96),(97),(98),(99),(100), (101);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: REPLACE...VALUES 语句未超过阈值的行数",
		"REPLACE INTO test_table (id) VALUES (1), (2), (3), (100);",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: INSERT...SELECT 语句估算行数超过阈值",
		"INSERT INTO test_table SELECT * FROM large_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE large_table (id INT); CREATE TABLE test_table (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO test_table SELECT * FROM large_table",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(90).AddRow(400),
			},
			{
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: INSERT...SELECT 语句估算行数未超过阈值",
		"INSERT INTO test_table SELECT id FROM small_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE small_table (id INT); CREATE TABLE test_table (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO test_table SELECT id FROM small_table;",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(100),
			}, {
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 7: REPLACE...SELECT 语句估算行数超过阈值",
		"REPLACE INTO test_table (id) SELECT id FROM large_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE large_table (id INT); CREATE TABLE test_table (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN REPLACE INTO test_table (id) SELECT id FROM large_table;",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(101),
			}, {
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: REPLACE...SELECT 语句估算行数未超过阈值",
		"REPLACE INTO test_table (id) SELECT id FROM small_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE small_table (id INT); CREATE TABLE test_table (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN REPLACE INTO test_table (id) SELECT id FROM small_table;",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(100),
			}, {
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 9: UNION 语句中一个 SELECT 子句估算行数超过阈值",
		"INSERT INTO test_table (id) SELECT id FROM large_table UNION SELECT id FROM small_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE large_table (id INT); CREATE TABLE test_table (id INT); CREATE TABLE small_table (id INT); "),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO test_table (id) SELECT id FROM large_table UNION SELECT id FROM small_table;",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(101).AddRow(50),
			}, {
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: UNION 语句中所有 SELECT 子句估算行数未超过阈值",
		"INSERT INTO test_table (id) SELECT id FROM small_table1 UNION SELECT id FROM small_table2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE small_table1 (id INT); CREATE TABLE test_table (id INT); CREATE TABLE small_table2 (id INT);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO test_table (id) SELECT id FROM small_table1 UNION SELECT id FROM small_table2;",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(50).AddRow(50),
			}, {
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 12: INSERT...SELECT 语句估算行数超过阈值(从xml中补充)",
		"INSERT INTO customers (id, cname, sex, city, age) SELECT id, cname, sex, city, age FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, cname VARCHAR(32), sex INT, city VARCHAR(32), age INT); INSERT INTO customers SELECT * FROM generate_series(1, 101);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO customers (id, cname, sex, city, age) SELECT id, cname, sex, city, age FROM customers;",
				Rows:  sqlmock.NewRows([]string{"rows"}).AddRow(101),
			}, {
				Query: "SHOW WARNINGS",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
