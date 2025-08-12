package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00108(t *testing.T) {
	ruleName := ai.SQLE00108
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: DELETE语句where中包含6层嵌套子查询",
		"DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 = 'value'))))))",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: DELETE语句where中包含4层嵌套子查询",
		"DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id = 'value'))))",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 3: INSERT...select语句中包含5层嵌套子查询",
		"INSERT INTO exist_db.exist_tb_1 (id) SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 = 'value')))))",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: INSERT语句where中包含4层嵌套子查询",
		"INSERT INTO exist_db.exist_tb_1 (id) SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id = 'value')))",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: SELECT语句where中包含6层嵌套子查询",
		"SELECT * FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 = 'value'))))))",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: SELECT语句where中包含4层嵌套子查询",
		"SELECT * FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id = 'value'))))",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: UPDATE语句where中包含6层嵌套子查询",
		"UPDATE exist_db.exist_tb_1 SET id = 'value' WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id0 IN (SELECT id1 FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 = 'value'))))))",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: UPDATE语句where中包含4层嵌套子查询",
		"UPDATE exist_db.exist_tb_1 SET id = 'value' WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id0 = 'value'))))",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: UNION语句where中包含6层嵌套子查询",
		"SELECT * FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 = 'value')))))) UNION SELECT * FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 = 'value')",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: UNION语句where中包含4层嵌套子查询",
		"SELECT * FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id = 'value')))) UNION SELECT * FROM exist_db.exist_tb_1 WHERE id2 = 'value'",
		nil, nil, newTestResult())

	// runAIRuleCase(rule, t, "case 11: WITH语句where中包含6层嵌套子查询",
	// 	"WITH CTE AS (SELECT * FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 = 'value')))))) SELECT * FROM CTE",
	// 	nil, nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 12: WITH语句where中包含5层嵌套子查询",
	// 	"WITH CTE AS (SELECT * FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM exist_db.exist_tb_1 WHERE id = 'value'))))) SELECT * FROM CTE",
	// 	nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 13_tes: SELECT语句where中包含6层嵌套子查询，使用示例中的表结构",
		"SELECT AVG(subquery_middle.subquery_grade) AS subquery_middle_avg FROM (SELECT grade AS subquery_grade FROM st1 WHERE st1.cid IN (SELECT cid FROM st_class WHERE cname = 'class2')) subquery_middle;",
		session.NewAIMockContext().WithSQL("CREATE TABLE st1 (id bigint, name VARCHAR(32), cid bigint, grade NUMERIC); CREATE TABLE st_class (cid bigint, cname VARCHAR(32));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 14: SELECT语句where中包含1层子查询，使用示例中的表结构",
		"SELECT count(1) cn FROM st1 CROSS JOIN (SELECT AVG(st1.grade) AS avg_grade FROM st1 WHERE st1.cid IN (SELECT cid FROM st_class WHERE cname = 'class2')) avg_grades WHERE st1.grade > avg_grades.avg_grade",
		session.NewAIMockContext().WithSQL("CREATE TABLE st1 (id bigint, name VARCHAR(32), cid bigint, grade NUMERIC); CREATE TABLE st_class (cid bigint, cname VARCHAR(32));"),
		nil, newTestResult())

	// 这个子查询中实际扫描表的子查询 只有1个，因此不算违规
	runAIRuleCase(rule, t, "case 15: SELECT语句中使用JOIN的ON条件中嵌套子查询5层",
		"SELECT st1.id FROM st1 JOIN st_class ON st1.cid in (SELECT cid FROM (SELECT cid FROM (SELECT cid FROM (SELECT cid FROM (SELECT cid FROM st_class WHERE cname = 'class2') AS sub1) AS sub2) AS sub3) AS sub4);",
		session.NewAIMockContext().WithSQL("CREATE TABLE st1 (id INT, cid INT); CREATE TABLE st_class (cid INT, cname VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 16: SELECT语句中使用JOIN的ON条件中嵌套子查询6层, 违规",
		"SELECT st1.id FROM st1 JOIN st_class ON st1.cid in (SELECT cid FROM (SELECT cid FROM (SELECT cid FROM (SELECT cid FROM (SELECT cid FROM st_class WHERE cname in (SELECT cname FROM st_class WHERE cname = 'class2')) AS sub1) AS sub2) AS sub3) AS sub4);",
		session.NewAIMockContext().WithSQL("CREATE TABLE st1 (id INT, cid INT); CREATE TABLE st_class (cid INT, cname VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: SELECT语句中 查询列中, 嵌套子查询2层",
		"SELECT 1, st1.id, (SELECT (SELECT id0 FROM exist_db.exist_tb_1 WHERE id1 = 'value') xx2 FROM exist_db.exist_tb_1 WHERE id1 = 'value') xxx FROM st1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE st1 (id INT, cid INT); CREATE TABLE st_class (cid INT, cname VARCHAR(50));"),
		nil, newTestResult())

}

// ==== Rule test code end ====
