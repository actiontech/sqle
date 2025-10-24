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
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

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

	// 新增测试用例：验证修复效果
	runAIRuleCase(rule, t, "case 33: UPDATE with BIGINT UNSIGNED and int constant (should pass)",
		"UPDATE t1 SET name = 'jack' WHERE id = 2838923;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id BIGINT UNSIGNED NOT NULL, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 34: UPDATE with BIGINT and int constant (should pass)",
		"UPDATE t1 SET name = 'jack' WHERE id = 2838923;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id BIGINT NOT NULL, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 35: UPDATE with INT and int constant (should pass)",
		"UPDATE t1 SET name = 'jack' WHERE id = 100;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT NOT NULL, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 36: UPDATE with BIGINT UNSIGNED and string constant (should fail)",
		"UPDATE t1 SET name = 'jack' WHERE id = '2838923';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id BIGINT UNSIGNED NOT NULL, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 37: SELECT with BIGINT UNSIGNED and large constant (should pass)",
		"SELECT * FROM t1 WHERE id = 18446744073709551615;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id BIGINT UNSIGNED NOT NULL, name VARCHAR(100));"),
		nil, newTestResult())

	// ==== 测试用例分类：数值类型大类匹配 ====
	// 测试目标：验证数值类型之间的匹配，确保不会导致隐式转换
	// 测试范围：cases 38-43

	// 整数类型匹配测试
	runAIRuleCase(rule, t, "case 38: TINYINT with INT constant (should pass - 数值在范围内)",
		"SELECT * FROM t1 WHERE age = 25;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (age TINYINT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 39: SMALLINT with INT constant (should pass - 数值在范围内)",
		"SELECT * FROM t1 WHERE score = 1000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (score SMALLINT, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 40: INT with INT constant (should pass)",
		"SELECT * FROM t1 WHERE id = 2147483647;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, name VARCHAR(100));"),
		nil, newTestResult())

	// 浮点数与整数匹配测试
	runAIRuleCase(rule, t, "case 41: FLOAT with INT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE price = 100;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (price FLOAT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 42: DOUBLE with FLOAT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE rate = 3.14;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (rate DOUBLE, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 精确数值类型匹配测试
	runAIRuleCase(rule, t, "case 43: DECIMAL with INT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE amount = 100;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (amount DECIMAL(10,2), name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// ==== 测试用例分类：字符串类型大类匹配 ====
	// 测试目标：验证字符串类型之间的匹配，确保不会导致隐式转换
	// 测试范围：cases 44-47

	// 基础字符串类型匹配测试
	runAIRuleCase(rule, t, "case 44: CHAR with VARCHAR constant (should pass - 字符串类型兼容)",
		"SELECT * FROM t1 WHERE code = 'ABC';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (code CHAR(10), name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 45: VARCHAR with CHAR constant (should pass - 字符串类型兼容)",
		"SELECT * FROM t1 WHERE name = 'John';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (name VARCHAR(100), code CHAR(10));"),
		nil, newTestResult())

	// 二进制字符串类型匹配测试
	runAIRuleCase(rule, t, "case 46: BLOB with BLOB constant (should pass)",
		"SELECT * FROM t1 WHERE data = 0x48656C6C6F;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (data BLOB, name VARCHAR(100));"),
		nil, newTestResult())

	// 枚举类型匹配测试
	runAIRuleCase(rule, t, "case 47: ENUM with ENUM constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE status = 'active';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (status ENUM('active', 'inactive'), name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// ==== 测试用例分类：日期时间类型大类匹配 ====
	// 测试目标：验证日期时间类型的匹配，检测隐式转换风险
	// 测试范围：cases 48-50B
	// 注意：虽然 MySQL 会自动转换字符串为日期，但隐式转换可能导致索引失效
	// 正确的做法：应该使用显式转换，如 CAST('1990-01-01' AS DATE) 或 DATE('1990-01-01')

	// 日期时间类型与字符串常量不匹配测试（应该报错）
	runAIRuleCase(rule, t, "case 48: DATE with string constant (should fail - 隐式转换风险)",
		"SELECT * FROM t1 WHERE birth_date = '1990-01-01';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (birth_date DATE, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 49: DATETIME with string constant (should fail - 隐式转换风险)",
		"SELECT * FROM t1 WHERE created_at = '2023-01-01 12:00:00';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (created_at DATETIME, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 50: TIMESTAMP with string constant (should fail - 隐式转换风险)",
		"SELECT * FROM t1 WHERE updated_at = '2023-01-01 12:00:00';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (updated_at TIMESTAMP, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 日期时间类型与日期时间类型列匹配测试（应该通过）
	runAIRuleCase(rule, t, "case 50A: DATE with DATE column (should pass)",
		"SELECT * FROM t1 WHERE birth_date = created_date;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (birth_date DATE, created_date DATE, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 50B: DATETIME with DATETIME column (should pass)",
		"SELECT * FROM t1 WHERE created_at = updated_at;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (created_at DATETIME, updated_at DATETIME, name VARCHAR(100));"),
		nil, newTestResult())

	// ==== 测试用例分类：跨大类不匹配 ====
	// 测试目标：验证不同大类类型之间的不匹配，确保检测到隐式转换风险
	// 测试范围：cases 51-54

	// 数值与字符串类型不匹配测试（应该报错）
	runAIRuleCase(rule, t, "case 51: INT with VARCHAR constant (should fail)",
		"SELECT * FROM t1 WHERE id = '123';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 52: VARCHAR with INT constant (should fail)",
		"SELECT * FROM t1 WHERE name = 123;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (name VARCHAR(100), id INT);"),
		nil, newTestResult().addResult(ruleName))

	// 日期时间与字符串类型不匹配测试（应该报错）
	// 正确的做法：使用 CAST('1990-01-01' AS DATE) 或 DATE('1990-01-01')
	runAIRuleCase(rule, t, "case 53: DATE with VARCHAR constant (should fail - 隐式转换风险)",
		"SELECT * FROM t1 WHERE birth_date = '1990-01-01';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (birth_date DATE, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 日期时间与数值类型不匹配测试（应该报错）
	// 正确的做法：使用 CAST(19900101 AS DATE) 或 DATE(19900101)
	runAIRuleCase(rule, t, "case 54: DATE with INT constant (should fail - 隐式转换风险)",
		"SELECT * FROM t1 WHERE birth_date = 19900101;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (birth_date DATE, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// ==== 测试用例分类：边界情况 ====
	// 测试目标：验证特殊类型和边界情况的匹配
	// 测试范围：cases 55-57

	// 位类型匹配测试
	runAIRuleCase(rule, t, "case 55: BIT with INT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE flag = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (flag BIT(1), name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 年份类型匹配测试
	runAIRuleCase(rule, t, "case 56: YEAR with INT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE year_col = 2023;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (year_col YEAR, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 集合类型匹配测试
	runAIRuleCase(rule, t, "case 57: SET with SET constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE tags = 'tag1,tag2';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (tags SET('tag1', 'tag2', 'tag3'), name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// ==== 新增测试用例：时间转换函数 ====

	// 测试使用时间函数与时间列的比较（应该通过）
	runAIRuleCase(rule, t, "case 58: DATE with CURRENT_DATE function (should pass)",
		"SELECT * FROM t1 WHERE birth_date = CURRENT_DATE();",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (birth_date DATE, name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 59: DATETIME with NOW function (should pass)",
		"SELECT * FROM t1 WHERE created_at = NOW();",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (created_at DATETIME, name VARCHAR(100));"),
		nil, newTestResult())

	// ==== 测试用例分类：列与值大类相同但会隐式转换的情况 ====
	// 测试目标：验证列与值在大类相同但具体类型不同时的隐式转换检测
	// 这些情况虽然大类相同，但会导致隐式转换，影响性能

	// 数值类型内部的隐式转换测试（应该报错）
	runAIRuleCase(rule, t, "case 60: TINYINT column with BIGINT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE age = 9223372036854775807;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (age TINYINT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 61: INT column with BIGINT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE id = 9223372036854775807;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id INT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 62: FLOAT column with DOUBLE constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE price = 3.141592653589793;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (price FLOAT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 63: DECIMAL(5,2) column with DECIMAL(10,4) constant (should pass - DECIMAL类型兼容)",
		"SELECT * FROM t1 WHERE amount = 123456.7890;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (amount DECIMAL(5,2), name VARCHAR(100));"),
		nil, newTestResult())

	// 字符串类型内部的兼容性测试
	runAIRuleCase(rule, t, "case 64: CHAR(10) column with long VARCHAR constant (should pass - MySQL转换值)",
		"SELECT * FROM t1 WHERE code = 'This is a very long string that exceeds CHAR(10) limit';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (code CHAR(10), name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 65: VARCHAR(50) column with TEXT constant (should pass - 字符串类型兼容)",
		"SELECT * FROM t1 WHERE description = 'This is a very long text that would be stored in TEXT type';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (description VARCHAR(50), name VARCHAR(100));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 66: BLOB column with BLOB constant (should pass - BLOB类型兼容)",
		"SELECT * FROM t1 WHERE data = 0x48656C6C6F576F726C64;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (data BLOB, name VARCHAR(100));"),
		nil, newTestResult())

	// 日期时间类型内部的隐式转换测试（应该报错）
	runAIRuleCase(rule, t, "case 67: DATE column with DATETIME constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE birth_date = '2023-01-01 12:00:00';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (birth_date DATE, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 68: DATETIME column with TIMESTAMP constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE created_at = '2023-01-01 12:00:00';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (created_at DATETIME, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 69: TIMESTAMP column with DATE constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE updated_at = '2023-01-01';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (updated_at TIMESTAMP, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 特殊数值类型的隐式转换测试（应该报错）
	runAIRuleCase(rule, t, "case 70: TINYINT UNSIGNED column with overflow constant (should fail - 超出范围)",
		"SELECT * FROM t1 WHERE status = 256;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (status TINYINT UNSIGNED, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 71: INT UNSIGNED column with BIGINT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE count = 9223372036854775807;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (count INT UNSIGNED, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 枚举和集合类型的隐式转换测试（应该报错）
	runAIRuleCase(rule, t, "case 72: ENUM column with invalid ENUM constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE status = 'invalid_status';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (status ENUM('active', 'inactive'), name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 73: SET column with invalid SET constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE tags = 'invalid_tag,another_invalid';",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (tags SET('tag1', 'tag2', 'tag3'), name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	// 边界值隐式转换测试（应该报错）
	runAIRuleCase(rule, t, "case 74: SMALLINT column with INT overflow constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE score = 32768;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (score SMALLINT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 75: MEDIUMINT column with BIGINT constant (should fail - 隐式转换)",
		"SELECT * FROM t1 WHERE id = 2147483648;",
		session.NewAIMockContext().WithSQL("CREATE TABLE t1 (id MEDIUMINT, name VARCHAR(100));"),
		nil, newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
