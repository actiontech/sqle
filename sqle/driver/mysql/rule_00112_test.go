package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00112(t *testing.T) {
	ruleName := ai.SQLE00112
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// SELECT语句测试用例
	runAIRuleCase(rule, t, "case 1: SELECT语句中WHERE子句比较Table1.column1 (INT)与Table1.column2 (INT)，预期通过",
		"SELECT * FROM Table1 WHERE column1 = column2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: SELECT语句中WHERE子句比较Table1.column1 (INT)与Table1.column3 (VARCHAR)，预期违规",
		"SELECT * FROM Table1 WHERE column1 = column3;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT语句中WHERE子句比较Table1.column1 (INT)与常量100 (INT)，预期通过",
		"SELECT * FROM Table1 WHERE column1 = 100;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: SELECT语句中WHERE子句比较Table1.column1 (INT)与常量'text' (VARCHAR)，预期违规",
		"SELECT * FROM Table1 WHERE column1 = 'text';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: SELECT语句中WHERE子句比较常量'100' (VARCHAR)与常量'200' (VARCHAR)，预期通过",
		"SELECT * FROM Table1 WHERE '100' = '200';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"), nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: SELECT语句中WHERE子句比较常量100 (INT)与常量'text' (VARCHAR)，预期违规",
		"SELECT * FROM Table1 WHERE 100 = 'text';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: SELECT语句中JOIN USING比较Table1.column1 (INT)与Table2.columnA (INT)，预期通过",
		"SELECT * FROM Table1 JOIN Table2 USING (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table2 (columnA INT, columnB VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: SELECT语句中JOIN USING比较Table1.column1 (INT)与Table2.column1 (VARCHAR)，预期违规",
		"SELECT * FROM Table1 JOIN Table2 USING (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table2 (column1 VARCHAR(100), columnB VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8.1: SELECT语句中复杂JOIN USING比较Table1.column1 (INT)与Table2.column1 (VARCHAR)，预期违规",
		"SELECT * FROM Table1 JOIN Table2 USING (column1) JOIN Table2 USING (column3);",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table2 (column1 VARCHAR(100), column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// UPDATE语句测试用例
	runAIRuleCase(rule, t, "case 9: UPDATE语句中WHERE子句比较Table1.column1 (INT)与Table1.column2 (INT)，预期通过",
		"UPDATE Table1 SET column3 = 'new_value' WHERE column1 = column2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: UPDATE语句中WHERE子句比较Table1.column1 (INT)与Table1.column3 (VARCHAR)，预期违规",
		"UPDATE Table1 SET column3 = 'new_value' WHERE column1 = column3;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: UPDATE语句中WHERE子句比较Table1.column1 (INT)与常量200 (INT)，预期通过",
		"UPDATE Table1 SET column3 = 'new_value' WHERE column1 = 200;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 12: UPDATE语句中WHERE子句比较Table1.column1 (INT)与常量'text' (VARCHAR)，预期违规",
		"UPDATE Table1 SET column3 = 'new_value' WHERE column1 = 'text';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: UPDATE语句中WHERE子句比较常量'300' (VARCHAR)与常量'400' (VARCHAR)，预期通过",
		"UPDATE Table1 SET column3 = 'new_value' WHERE '300' = '400';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"), nil, newTestResult())

	runAIRuleCase(rule, t, "case 14: UPDATE语句中WHERE子句比较常量300 (INT)与常量'text' (VARCHAR)，预期违规",
		"UPDATE Table1 SET column3 = 'new_value' WHERE 300 = 'text';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: UPDATE语句中JOIN USING比较Table1.column1 (INT)与Table2.columnA (INT)，预期通过",
		"UPDATE Table1 JOIN Table2 USING (column1) SET Table1.column3 = 'updated';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table2 (columnA INT, columnB VARCHAR(100));"),
		nil, newTestResult())

	// DELETE语句测试用例
	runAIRuleCase(rule, t, "case 17: DELETE语句中WHERE子句比较Table1.column1 (INT)与Table1.column2 (INT)，预期通过",
		"DELETE FROM Table1 WHERE column1 = column2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 18: DELETE语句中WHERE子句比较Table1.column1 (INT)与Table1.column3 (VARCHAR)，预期违规",
		"DELETE FROM Table1 WHERE column1 = column3;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 19: DELETE语句中WHERE子句比较Table1.column1 (INT)与常量300 (INT)，预期通过",
		"DELETE FROM Table1 WHERE column1 = 300;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 20: DELETE语句中WHERE子句比较Table1.column1 (INT)与常量'text' (VARCHAR)，预期违规",
		"DELETE FROM Table1 WHERE column1 = 'text';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 21: DELETE语句中WHERE子句比较常量'400' (VARCHAR)与常量'500' (VARCHAR)，预期通过",
		"DELETE FROM Table1 WHERE '400' = '500';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"), nil, newTestResult())

	runAIRuleCase(rule, t, "case 22: DELETE语句中WHERE子句比较常量400 (INT)与常量'text' (VARCHAR)，预期违规",
		"DELETE FROM Table1 WHERE 400 = 'text';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));"), nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 23: DELETE语句中JOIN USING比较Table1.column1 (INT)与Table2.columnA (INT)，预期通过",
		"DELETE Table1 FROM Table1 JOIN Table2 USING (column1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table2 (columnA INT, columnB VARCHAR(100));"),
		nil, newTestResult())

	// INSERT语句测试用例
	runAIRuleCase(rule, t, "case 25: INSERT INTO SELECT语句中WHERE子句比较Table1.column1 (INT)与Table2.columnA (INT)，预期通过",
		"INSERT INTO Table3 (column1, column3) SELECT column1, column3 FROM Table1 JOIN Table2 USING (column1) WHERE Table1.column1 = Table2.columnA;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table2 (columnA INT, columnB VARCHAR(100));").
			WithSQL("CREATE TABLE Table3 (column1 INT, column3 VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 26: INSERT INTO SELECT语句中WHERE子句比较Table1.column1 (INT)与Table2.column3 (VARCHAR)，预期违规",
		"INSERT INTO Table3 (column1, column3) SELECT column1, column3 FROM Table1 JOIN Table2 USING (column1) WHERE Table1.column1 = Table2.column3;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table2 (columnA INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table3 (column1 INT, column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 27: INSERT INTO SELECT语句中WHERE子句比较Table1.column1 (INT)与常量500 (INT)，预期通过",
		"INSERT INTO Table3 (column1, column3) SELECT column1, column3 FROM Table1 WHERE column1 = 500;",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table3 (column1 INT, column3 VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 28: INSERT INTO SELECT语句中WHERE子句比较Table1.column1 (INT)与常量'text' (VARCHAR)，预期违规",
		"INSERT INTO Table3 (column1, column3) SELECT column1, column3 FROM Table1 WHERE column1 = 'text';",
		session.NewAIMockContext().WithSQL("CREATE TABLE Table1 (column1 INT, column2 INT, column3 VARCHAR(100));").
			WithSQL("CREATE TABLE Table3 (column1 INT, column3 VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 新增示例
	runAIRuleCase(rule, t, "case 29: SELECT语句中WHERE子句比较customers.c_id (INT)与常量'123' (VARCHAR)，预期违规",
		"SELECT * FROM customers WHERE c_id = '123';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (c_id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 30: SELECT语句中JOIN ON比较customers.c_id (INT)与orders.c_id (VARCHAR)，预期违规",
		"SELECT * FROM customers a JOIN orders b ON a.c_id = b.c_id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (c_id INT, name VARCHAR(100));").
			WithSQL("CREATE TABLE orders (c_id VARCHAR(100), order_date DATE);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 30.1: SELECT语句中复杂JOIN ON比较customers.c_id (INT)与orders.c_id (VARCHAR)，预期违规",
		"SELECT * FROM customers a JOIN orders b ON a.c_id = b.c_id JOIN orders c ON c.c_id = b.c_id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (c_id INT, name VARCHAR(100));").
			WithSQL("CREATE TABLE orders (c_id VARCHAR(100), order_date DATE);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 31: UPDATE语句中WHERE子句比较customers.log_date (VARCHAR)与常量CURRENT_DATE() (DATE)，预期违规",
		"UPDATE customers SET name = 'updated' WHERE log_date = CURRENT_DATE();",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (c_id INT, name VARCHAR(100), log_date VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 32: DELETE语句中JOIN ON比较orders.c_id (VARCHAR)与customers.c_id (INT)，预期违规",
		"DELETE FROM orders WHERE c_id IN (SELECT c_id FROM customers WHERE c_id = orders.c_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (c_id VARCHAR(100), order_date DATE);").
			WithSQL("CREATE TABLE customers (c_id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
