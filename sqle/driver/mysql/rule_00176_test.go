package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00176(t *testing.T) {
	ruleName := ai.SQLE00176
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: DELETE 语句包含 hint 指令 /*+ ... */",
		"DELETE /*+ MAX_EXECUTION_TIME(1000) */ FROM test_table WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: DELETE 语句包含 FORCE INDEX",
		"DELETE FROM test_table FORCE INDEX (idx_test) WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: DELETE 语句不包含 hint 指令",
		"DELETE FROM test_table WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: INSERT 语句包含 hint 指令 /*+ ... */",
		"INSERT /*+ MAX_EXECUTION_TIME(1000) */ INTO test_table (id, name) VALUES (1, 'test');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: INSERT 语句不包含 hint 指令",
		"INSERT INTO test_table (id, name) VALUES (1, 'test');",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: UPDATE 语句包含 hint 指令 /*+ ... */",
		"UPDATE /*+ MAX_EXECUTION_TIME(1000) */ test_table SET name = 'test' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: UPDATE 语句包含 USE INDEX",
		"UPDATE test_table USE INDEX (idx_test) SET name = 'test' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: UPDATE 语句不包含 hint 指令",
		"UPDATE test_table SET name = 'test' WHERE id = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 11: SELECT 语句包含 hint 指令 /*+ ... */",
		"SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: SELECT 语句包含 STRAIGHT_JOIN",
		"SELECT * FROM test_table STRAIGHT_JOIN another_table ON test_table.id = another_table.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50)); CREATE TABLE another_table (id INT PRIMARY KEY, description VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: SELECT 语句不包含 hint 指令",
		"SELECT * FROM test_table;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult())

	// runAIRuleCase(rule, t, "case 14: SELECT 语句包含 hint 指令 /*+ ... */ (从xml中补充)",
	// 	"SELECT /*+ index (customers idx_sex_customers) */ * FROM customers WHERE age = 20;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
	// 	nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: SELECT 语句包含 FORCE INDEX (从xml中补充)",
		"SELECT age FROM customers FORCE INDEX FOR GROUP BY (idx_sex_customers) GROUP BY age;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: SELECT 语句包含 STRAIGHT_JOIN (从xml中补充)",
		"SELECT STRAIGHT_JOIN COUNT(*) FROM customers a JOIN customers_small b ON a.id = b.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT); CREATE TABLE customers_small (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 17: SELECT 语句包含 set_var hint 指令 /*+ ... */ (从xml中补充)",
	// 	"SELECT /*+ set_var(max_execution_time=1) */ * FROM customers ORDER BY name DESC LIMIT 1;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
	// 	nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: UPDATE 语句包含 FORCE INDEX (从xml中补充)",
		"UPDATE customers FORCE INDEX (idx_sex_customers) SET sex = 1 WHERE age = 20;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 19: DELETE 语句包含 hint 指令 /*+ ... */ (从xml中补充)",
	// 	"DELETE /*+ index(customers idx_sex_customers) */ FROM customers WHERE age = 20;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
	// 	nil, newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 20: INSERT 语句包含 hint 指令 /*+ ... */ (从xml中补充)",
	// 	"INSERT /*+ index(customers idx_sex_customers) */ INTO customers VALUES (9999999,'xiaozhangxiaowangxiaofeng',0,'shanghai',90);",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
	// 	nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 22: SELECT 语句不包含 hint 指令 (从xml中补充)",
		"SELECT * FROM customers WHERE age = 20;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 23: SELECT 语句不包含 hint 指令 (从xml中补充)",
		"SELECT age FROM customers GROUP BY age;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 24: SELECT 语句不包含 hint 指令 (从xml中补充)",
		"SELECT age FROM customers ORDER BY age;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 25: SELECT 语句不包含 hint 指令 (从xml中补充)",
		"SELECT COUNT(*) FROM customers a JOIN customers_small b ON a.id = b.id;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT); CREATE TABLE customers_small (id INT PRIMARY KEY, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 26: SELECT 语句不包含 hint 指令 (从xml中补充)",
		"SELECT * FROM customers ORDER BY name DESC LIMIT 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 27: UPDATE 语句不包含 hint 指令 (从xml中补充)",
		"UPDATE customers SET sex = 1 WHERE age = 20;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 28: DELETE 语句不包含 hint 指令 (从xml中补充)",
		"DELETE FROM customers WHERE age = 20;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 29: INSERT 语句不包含 hint 指令 (从xml中补充)",
		"INSERT INTO customers VALUES (9999999,'xiaozhangxiaowangxiaofeng',0,'shanghai',90);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), sex INT, city VARCHAR(50), age INT);"),
		nil, newTestResult())

}

// ==== Rule test code end ====
