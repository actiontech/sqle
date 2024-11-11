package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00001(t *testing.T) {
	ruleName := ai.SQLE00001
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	// SELECT statements
	runAIRuleCase(rule, t, "case 1: SELECT statement without WHERE clause",
		"SELECT * FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: SELECT statement with WHERE condition always true",
		"SELECT * FROM employees WHERE 1=1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2_1: SELECT statement with WHERE condition always true",
		"SELECT * FROM employees WHERE  (1=1) or (name='JERYY');",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2_2: SELECT statement with WHERE condition always true",
		"SELECT * FROM employees WHERE 1=2 or (1=2 or name=name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2_3: SELECT statement with WHERE condition always true",
		"SELECT * FROM employees WHERE  (1=2 or name=name) or 1=2;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT statement with WHERE condition using OR with always true expression",
		"SELECT * FROM employees WHERE department = 'Sales' OR TRUE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: SELECT statement with WHERE condition using OR with always false expression",
		"SELECT * FROM employees WHERE department = 'Sales' OR false;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 3_1: SELECT statement with WHERE condition on non-nullable column using IS NOT NULL",
		"SELECT * FROM employees WHERE employee_id IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 3_2: SELECT statement with WHERE condition on nullable column using IS NOT NULL",
		"SELECT * FROM employees WHERE employee_id IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: SELECT statement with WHERE condition using OR with IS NOT NULL on non-nullable column",
		"SELECT * FROM employees WHERE department = 'HR' OR employee_id IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	// // INSERT statements
	runAIRuleCase(rule, t, "case 9: INSERT statement using SELECT without WHERE clause",
		"INSERT INTO employees_archive SELECT * FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));").WithSQL("CREATE TABLE employees_archive (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: INSERT statement using SELECT with WHERE condition always true",
		"INSERT INTO employees_archive SELECT * FROM employees WHERE 1=1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));").WithSQL("CREATE TABLE employees_archive (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: INSERT statement using SELECT with valid WHERE condition",
		"INSERT INTO employees_archive SELECT * FROM employees WHERE active = 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), active TINYINT);").WithSQL("CREATE TABLE employees_archive (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult())

	// UPDATE statements
	runAIRuleCase(rule, t, "case 12: UPDATE statement without WHERE clause",
		"UPDATE employees SET salary = salary * 1.1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), salary DECIMAL(10, 2));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: UPDATE statement with WHERE condition always true",
		"UPDATE employees SET salary = salary * 1.1 WHERE 1=1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), salary DECIMAL(10, 2));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: UPDATE statement with WHERE condition using OR with always true expression",
		"UPDATE employees SET salary = salary * 1.1 WHERE department = 'IT' OR TRUE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), salary DECIMAL(10, 2));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 15: UPDATE statement with valid WHERE condition",
		"UPDATE employees SET salary = salary * 1.1 WHERE performance_rating = 'A';",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), salary DECIMAL(10, 2), performance_rating CHAR(1));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 16: UPDATE statement with WHERE condition on non-nullable column using IS NOT NULL",
		"UPDATE employees SET bonus = 1000 WHERE employee_id IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), bonus DECIMAL(10, 2));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: UPDATE statement with WHERE condition on nullable column using IS NOT NULL",
		"UPDATE employees SET bonus = 1000 WHERE middle_name IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), bonus DECIMAL(10, 2));"),
		nil,
		newTestResult())

	// DELETE statements
	runAIRuleCase(rule, t, "case 18: DELETE statement without WHERE clause",
		"DELETE FROM employees;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 19: DELETE statement with WHERE condition always true",
		"DELETE FROM employees WHERE 1=1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 20: DELETE statement with WHERE condition using OR with always true expression",
		"DELETE FROM employees WHERE department = 'Marketing' OR TRUE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 21: DELETE statement with valid WHERE condition",
		"DELETE FROM employees WHERE last_login < '2023-01-01';",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), last_login DATE);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 22: DELETE statement with WHERE condition on non-nullable column using IS NOT NULL",
		"DELETE FROM employees WHERE employee_id IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 23: DELETE statement with WHERE condition on nullable column using IS NOT NULL",
		"DELETE FROM employees WHERE middle_name IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
		nil,
		newTestResult())

	// UNION statements
	runAIRuleCase(rule, t, "case 24: UNION statement with subSELECT without WHERE clause",
		"SELECT name FROM employees UNION SELECT name FROM contractors;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));").WithSQL("CREATE TABLE contractors (contractor_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), contract_end DATE);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 25: UNION statement with subSELECT with WHERE condition always true",
		"SELECT name FROM employees WHERE 1=1 UNION SELECT name FROM contractors WHERE TRUE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));").WithSQL("CREATE TABLE contractors (contractor_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), contract_end DATE);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 26: UNION statement with subSELECT with valid WHERE condition",
		"SELECT name FROM employees WHERE active = 1 UNION SELECT name FROM contractors WHERE contract_end > '2023-12-31';",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), active TINYINT);").WithSQL("CREATE TABLE contractors (contractor_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), contract_end DATE);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 27: UNION statement with subSELECT using OR with always true expression",
		"SELECT name FROM employees WHERE department = 'Finance' OR 1=1 UNION SELECT name FROM contractors WHERE department = 'Finance' OR TRUE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));").WithSQL("CREATE TABLE contractors (contractor_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), contract_end DATE);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 28: UNION statement with subSELECT on non-nullable column using IS NOT NULL",
		"SELECT name FROM employees WHERE employee_id IS NOT NULL UNION SELECT name FROM contractors WHERE contractor_id IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));").WithSQL("CREATE TABLE contractors (contractor_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), contract_end DATE);"),
		nil,
		newTestResult().addResult(ruleName))

	// // WITH statements
	// runAIRuleCase(rule, t, "case 29: WITH statement with SELECT without WHERE clause",
	// 	"WITH active_employees AS (SELECT * FROM employees) SELECT * FROM active_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
	// 	nil,
	// 	newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 30: WITH statement with SELECT with WHERE condition always true",
	// 	"WITH active_employees AS (SELECT * FROM employees WHERE 1=1) SELECT * FROM active_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
	// 	nil,
	// 	newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 31: WITH statement with SELECT with valid WHERE condition",
	// 	"WITH active_employees AS (SELECT * FROM employees WHERE status = 'active') SELECT * FROM active_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50), status VARCHAR(10));"),
	// 	nil,
	// 	newTestResult())

	// runAIRuleCase(rule, t, "case 32: WITH statement with SELECT using OR with always true expression",
	// 	"WITH active_employees AS (SELECT * FROM employees WHERE department = 'IT' OR TRUE) SELECT * FROM active_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
	// 	nil,
	// 	newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 33: WITH statement with SELECT on non-nullable column using IS NOT NULL",
	// 	"WITH valid_employees AS (SELECT * FROM employees WHERE employee_id IS NOT NULL) SELECT * FROM valid_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
	// 	nil,
	// 	newTestResult().addResult(ruleName))

	// runAIRuleCase(rule, t, "case 34: WITH statement with SELECT on nullable column using IS NOT NULL",
	// 	"WITH named_employees AS (SELECT * FROM employees WHERE middle_name IS NOT NULL) SELECT * FROM named_employees;",
	// 	session.NewAIMockContext().WithSQL("CREATE TABLE employees (employee_id INT NOT NULL, name VARCHAR(100), department VARCHAR(50), age INT, middle_name VARCHAR(50));"),
	// 	nil,
	// 	newTestResult())

	// * 新增示例
	runAIRuleCase(rule, t, "case 35: SELECT statement without WHERE clause on customers table (从xml中补充)",
		"SELECT * FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 36: SELECT statement with WHERE condition always true on customers table (从xml中补充)",
		"SELECT * FROM customers WHERE 1=1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 37: SELECT statement with valid WHERE condition on customers table (从xml中补充)",
		"SELECT * FROM customers WHERE age = 22;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 38: INSERT statement using SELECT without WHERE clause on customers table (从xml中补充)",
		"INSERT INTO customers_insert SELECT * FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);").WithSQL("CREATE TABLE customers_insert (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 39: UPDATE statement without WHERE clause on customers table (从xml中补充)",
		"UPDATE customers SET age = 30;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 40: DELETE statement without WHERE clause on customers table (从xml中补充)",
		"DELETE FROM customers;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 41: SELECT statement with WHERE EXISTS (haven't table)",
		"SELECT * FROM customers WHERE EXISTS (SELECT 1 FROM dual);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 42: SELECT statement with WHERE EXISTS (have table)",
		"SELECT * FROM customers WHERE EXISTS (SELECT 1 FROM customers where customer_id=1);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 43: SELECT statement with WHERE not EXISTS ",
		"SELECT * FROM customers WHERE not EXISTS (SELECT 1 FROM dual);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 44: SELECT statement with WHERE 1 in (1,2,3) ",
		"SELECT * FROM customers WHERE 1 in (1,2,3);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 45: SELECT statement with WHERE 1 in (SELECT 1 FROM dual) ",
		"SELECT * FROM customers WHERE 1 in (SELECT 1 FROM dual);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 46: SELECT statement with WHERE 1 in (SELECT 2 FROM dual) ",
		"SELECT * FROM customers WHERE 1 in (SELECT 2 FROM dual union SELECT 1 FROM dual );",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 47: SELECT statement with WHERE COALESCE(customer_id, 'default') IS NOT NULL ",
		"select count(*) from customers WHERE COALESCE(customer_id, 'default') IS NOT NULL;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT NOT NULL, name VARCHAR(100), age INT);"),
		nil,
		newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
