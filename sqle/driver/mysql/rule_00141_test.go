package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00141(t *testing.T) {
	ruleName := ai.SQLE00141
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT 语句中包含超过3个表的嵌套 JOIN",
		"SELECT * FROM table1 JOIN table2 ON table1.id = table2.id JOIN table3 ON table2.id = table3.id JOIN table4 ON table3.id = table4.id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table4 (id INT, name VARCHAR(50));CREATE TABLE table3 (id INT, name VARCHAR(50));CREATE TABLE table2 (id INT, name VARCHAR(50));CREATE TABLE table1 (id INT, name VARCHAR(50));",
		),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT 语句中包含3个表的嵌套 JOIN",
		"SELECT * FROM table1 JOIN table2 ON table1.id = table2.id JOIN table3 ON table2.id = table3.id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table3 (id INT, name VARCHAR(50));CREATE TABLE table2 (id INT, name VARCHAR(50));CREATE TABLE table1 (id INT, name VARCHAR(50));",
		),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: UPDATE 语句中包含超过3个表的嵌套 JOIN",
		"UPDATE table1 JOIN table2 ON table1.id = table2.id JOIN table3 ON table2.id = table3.id JOIN table4 ON table3.id = table4.id SET table1.name = 'Test';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table4 (id INT, name VARCHAR(50));CREATE TABLE table3 (id INT, name VARCHAR(50));CREATE TABLE table2 (id INT, name VARCHAR(50));CREATE TABLE table1 (id INT, name VARCHAR(50));",
		),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: UPDATE 语句中包含3个表的嵌套 JOIN",
		"UPDATE table1 JOIN table2 ON table1.id = table2.id JOIN table3 ON table2.id = table3.id SET table1.name = 'Test';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table3 (id INT, name VARCHAR(50));CREATE TABLE table2 (id INT, name VARCHAR(50));CREATE TABLE table1 (id INT, name VARCHAR(50));",
		),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: SELECT 语句中包含超过3个表的嵌套 JOIN (从xml中补充)",
		"SELECT a.id, a.name, b.post, c.addr, d.dp_name FROM st1 a JOIN st_ps b ON a.pid = b.id JOIN st_addr c ON a.addr_id = c.id JOIN st_dp d ON a.dp_id = d.id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE st_dp (id INT, dp_name VARCHAR(50));CREATE TABLE st_addr (id INT, addr VARCHAR(100));CREATE TABLE st_ps (id INT, post VARCHAR(50));CREATE TABLE st1 (id INT, name VARCHAR(50), pid INT, addr_id INT, dp_id INT);",
		),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: INSERT 语句中包含超过3个表的嵌套 JOIN (从xml中补充)",
		"INSERT INTO st_detail(id, name, post, addr, dp_name, grade) SELECT a.id, a.name, b.post, c.addr, d.dp_name, a.grade FROM st1 a JOIN st_ps b ON a.pid = b.id JOIN st_addr c ON a.addr_id = c.id JOIN st_dp d ON a.dp_id = d.id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE st_dp (id INT, dp_name VARCHAR(50));CREATE TABLE st_addr (id INT, addr VARCHAR(100));CREATE TABLE st_ps (id INT, post VARCHAR(50))CREATE TABLE st1 (id INT, name VARCHAR(50), pid INT, addr_id INT, dp_id INT, grade INT);CREATE TABLE st_detail (id INT, name VARCHAR(50), post VARCHAR(50), addr VARCHAR(100), dp_name VARCHAR(50), grade INT);",
		),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: UPDATE 语句中包含超过3个表的嵌套 JOIN (从xml中补充)",
		"UPDATE st1 JOIN st_ps b ON st1.pid = b.id JOIN st_addr c ON st1.pid = c.id JOIN st_dp d ON st1.pid = d.id SET grade = grade + 10;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE st_dp (id INT);CREATE TABLE st_addr (id INT);CREATE TABLE st_ps (id INT);CREATE TABLE st1 (id INT, pid INT, grade INT);",
		),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: DELETE 语句中包含超过3个表的嵌套 JOIN (从xml中补充)",
		"DELETE st1 FROM st1 JOIN st_ps b ON st1.pid = b.id JOIN st_addr c ON st1.pid = c.id JOIN st_dp d ON st1.pid = d.id WHERE b.id = 2;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE st_dp (id INT);CREATE TABLE st_addr (id INT);CREATE TABLE st_ps (id INT);CREATE TABLE st1 (id INT, pid INT);",
		),
		nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
