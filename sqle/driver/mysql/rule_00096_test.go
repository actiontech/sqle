package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00096(t *testing.T) {
	ruleName := ai.SQLE00096
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: SELECT 语句涉及3个表，未超过阈值",
		"SELECT a.*, b.*, c.* FROM table_a a JOIN table_b b ON a.id = b.a_id JOIN table_c c ON b.id = c.b_id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));",
		), nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: SELECT 语句涉及4个表，超过阈值",
		"SELECT a.*, b.*, c.*, d.* FROM table_a a JOIN table_b b ON a.id = b.a_id JOIN table_c c ON b.id = c.b_id JOIN table_d d ON c.id = d.c_id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));CREATE TABLE table_d (id INT, c_id INT, info VARCHAR(100));",
		), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: UPDATE 语句涉及3个表，未超过阈值",
		"UPDATE table_a a JOIN table_b b ON a.id = b.a_id SET a.value = b.value WHERE a.id = 1;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, value VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));",
		), nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: UPDATE 语句涉及4个表，超过阈值",
		"UPDATE table_a a JOIN table_b b ON a.id = b.a_id JOIN table_c c ON b.id = c.b_id JOIN table_d d ON c.id = d.c_id SET a.value = d.value WHERE a.id = 1;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, value VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));CREATE TABLE table_d (id INT, c_id INT, info VARCHAR(100));",
		), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: DELETE 语句涉及3个表，未超过阈值",
		"DELETE a FROM table_a a JOIN table_b b ON a.id = b.a_id JOIN table_c c ON b.id = c.b_id WHERE a.id = 1;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));",
		), nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: DELETE 语句涉及4个表，超过阈值",
		"DELETE a FROM table_a a JOIN table_b b ON a.id = b.a_id JOIN table_c c ON b.id = c.b_id JOIN table_d d ON c.id = d.c_id WHERE a.id = 1;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));CREATE TABLE table_d (id INT, c_id INT, info VARCHAR(100));",
		), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: INSERT ... SELECT 语句涉及3个表，未超过阈值",
		"INSERT INTO table_e (id, value) SELECT a.id, b.value FROM table_a a JOIN table_b b ON a.id = b.a_id JOIN table_c c ON b.id = c.b_id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));CREATE TABLE table_e (id INT, value VARCHAR(100));",
		), nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: INSERT ... SELECT 语句涉及4个表，超过阈值",
		"INSERT INTO table_e (id, value) SELECT a.id, d.value FROM table_a a JOIN table_b b ON a.id = b.a_id JOIN table_c c ON b.id = c.b_id JOIN table_d d ON c.id = d.c_id;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));CREATE TABLE table_d (id INT, c_id INT, info VARCHAR(100));CREATE TABLE table_e (id INT, value VARCHAR(100));",
		), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: UNION 语句涉及3个表，未超过阈值",
		"SELECT id FROM table_a UNION SELECT id FROM table_b UNION SELECT id FROM table_c;",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, name VARCHAR(100));CREATE TABLE table_c (id INT, name VARCHAR(100));",
		), nil, newTestResult())

	// runAIRuleCase(rule, t, "case 10: UNION 语句涉及4个表，超过阈值",
	// 	"SELECT id FROM table_a UNION SELECT id FROM table_b UNION SELECT id FROM table_c UNION SELECT id FROM table_d;",
	// 	session.NewAIMockContext().WithSQL(
	// 		"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, name VARCHAR(100));CREATE TABLE table_c (id INT, name VARCHAR(100));CREATE TABLE table_d (id INT, name VARCHAR(100));",
	// 	), nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 11: WITH 语句涉及3个表，未超过阈值",
	// 	"WITH cte1 AS (SELECT * FROM table_a), cte2 AS (SELECT * FROM table_b) SELECT cte1.id, cte2.value FROM cte1 JOIN cte2 ON cte1.id = cte2.a_id JOIN table_c ON cte2.id = table_c.b_id;",
	// 	session.NewAIMockContext().WithSQL(
	// 		"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));",
	// 	), nil, newTestResult())

	// runAIRuleCase(rule, t, "case 12: WITH 语句涉及4个表，超过阈值",
	// 	"WITH cte1 AS (SELECT * FROM table_a), cte2 AS (SELECT * FROM table_b), cte3 AS (SELECT * FROM table_c) SELECT cte1.id, cte2.value, table_d.info FROM cte1 JOIN cte2 ON cte1.id = cte2.a_id JOIN cte3 ON cte2.id = cte3.b_id JOIN table_d ON cte3.id = table_d.c_id;",
	// 	session.NewAIMockContext().WithSQL(
	// 		"CREATE TABLE table_a (id INT, name VARCHAR(100));CREATE TABLE table_b (id INT, a_id INT, value VARCHAR(100));CREATE TABLE table_c (id INT, b_id INT, description VARCHAR(100));CREATE TABLE table_d (id INT, c_id INT, info VARCHAR(100));",
	// 	), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: SELECT 语句涉及6个表，超过阈值(从xml中补充)",
		"SELECT a.id, a.name, c.name, hr.is_health, sc.score, t.name FROM student a, sc, course c, health_archives ha, health_report hr, teacher t WHERE a.id = sc.student_id AND sc.course_id = c.id AND ha.student_id = a.id AND ha.health_report_id = hr.id AND t.id = c.teacher_id AND sc.score > 80 AND c.name='课程911' AND hr.is_health='强壮';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE student (id INT, name VARCHAR(100));CREATE TABLE sc (student_id INT, course_id INT, score INT);CREATE TABLE course (id INT, name VARCHAR(100), teacher_id INT);CREATE TABLE health_archives (student_id INT, health_report_id INT);CREATE TABLE health_report (id INT, is_health VARCHAR(100));CREATE TABLE teacher (id INT, name VARCHAR(100));",
		), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: UPDATE 语句涉及6个表，超过阈值(从xml中补充)",
		"UPDATE student a, sc, course c, health_archives ha, health_report hr, teacher t SET a.update_time = NOW() WHERE a.id = sc.student_id AND sc.course_id = c.id AND ha.student_id = a.id AND ha.health_report_id = hr.id AND t.id = c.teacher_id AND sc.score > 80 AND c.name='课程911' AND hr.is_health='强壮';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE student (id INT, name VARCHAR(100),update_time date);CREATE TABLE sc (student_id INT, course_id INT, score INT);CREATE TABLE course (id INT, name VARCHAR(100), teacher_id INT);CREATE TABLE health_archives (student_id INT, health_report_id INT);CREATE TABLE health_report (id INT, is_health VARCHAR(100));CREATE TABLE teacher (id INT, name VARCHAR(100));",
		), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: DELETE 语句涉及6个表，超过阈值(从xml中补充)",
		"DELETE a FROM student a INNER JOIN sc INNER JOIN course c INNER JOIN health_archives ha INNER JOIN health_report hr INNER JOIN teacher t WHERE a.id = sc.student_id AND sc.course_id = c.id AND ha.student_id = a.id AND ha.health_report_id = hr.id AND t.id = c.teacher_id AND sc.score > 80 AND c.name='课程911' AND hr.is_health='强壮';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE student (id INT, name VARCHAR(100));CREATE TABLE sc (student_id INT, course_id INT, score INT);CREATE TABLE course (id INT, name VARCHAR(100), teacher_id INT);CREATE TABLE health_archives (student_id INT, health_report_id INT);CREATE TABLE health_report (id INT, is_health VARCHAR(100));CREATE TABLE teacher (id INT, name VARCHAR(100));",
		), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: INSERT ... SELECT 语句涉及6个表，超过阈值(从xml中补充)",
		"INSERT INTO student_ids SELECT a.id FROM student a, sc, course c, health_archives ha, health_report hr, teacher t WHERE a.id = sc.student_id AND sc.course_id = c.id AND ha.student_id = a.id AND ha.health_report_id = hr.id AND t.id = c.teacher_id AND sc.score > 80 AND c.name='课程911' AND hr.is_health='强壮';",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE student (id INT, name VARCHAR(100));CREATE TABLE sc (student_id INT, course_id INT, score INT);CREATE TABLE course (id INT, name VARCHAR(100), teacher_id INT);CREATE TABLE health_archives (student_id INT, health_report_id INT);CREATE TABLE health_report (id INT, is_health VARCHAR(100));CREATE TABLE teacher (id INT, name VARCHAR(100));CREATE TABLE student_ids (id INT);",
		), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: SELECT 子查询语句涉及6个表，超过阈值(从xml中补充)",
		"SELECT * from student where name in (SELECT a.name FROM student a, sc, course c, health_archives ha, health_report hr, teacher t WHERE a.id = sc.student_id AND sc.course_id = c.id AND ha.student_id = a.id AND ha.health_report_id = hr.id AND t.id = c.teacher_id AND sc.score > 80 AND c.name='课程911' AND hr.is_health='强壮')",
		session.NewAIMockContext().WithSQL(
			"CREATE TABLE student (id INT, name VARCHAR(100));CREATE TABLE sc (student_id INT, course_id INT, score INT);CREATE TABLE course (id INT, name VARCHAR(100), teacher_id INT);CREATE TABLE health_archives (student_id INT, health_report_id INT);CREATE TABLE health_report (id INT, is_health VARCHAR(100));CREATE TABLE teacher (id INT, name VARCHAR(100));",
		), nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
