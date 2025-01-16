package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00110(t *testing.T) {
	ruleName := ai.SQLE00110
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: DELETE statement with indexed field in WHERE clause",
		"DELETE FROM employee WHERE emp_id = 123;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), age INT);"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 2: DELETE statement with non-indexed field in WHERE clause",
		"DELETE FROM employee WHERE name = 'John';",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), age INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 3: DELETE statement with mixed indexed and non-indexed fields in WHERE clause",
		"DELETE FROM employee WHERE emp_id = 123 AND age = 30;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), age INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 4: INSERT statement with SELECT using indexed fields in WHERE clause",
		"INSERT INTO employee_archive (emp_id, name) SELECT emp_id, name FROM employee WHERE dept_id = 10;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT, INDEX idx_dept_id (dept_id));").
			WithSQL("CREATE TABLE employee_archive (emp_id INT, name VARCHAR(50));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 5: INSERT statement with SELECT using non-indexed fields in WHERE clause",
		"INSERT INTO employee_archive (emp_id, name) SELECT emp_id, name FROM employee WHERE age > 30;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT);").
			WithSQL("CREATE TABLE employee_archive (emp_id INT, name VARCHAR(50));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 6: SELECT statement with indexed fields in WHERE, GROUP BY, ORDER BY clauses",
		"SELECT emp_id, COUNT(*) FROM employee WHERE dept_id = 10 GROUP BY emp_id ORDER BY COUNT(*) DESC;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT, INDEX idx_dept_id (dept_id), INDEX idx_emp_id (emp_id));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 7: SELECT statement with non-indexed field in WHERE clause",
		"SELECT emp_id, name FROM employee WHERE name = 'Alice';",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 8: SELECT statement with non-indexed field in GROUP BY clause",
		"SELECT dept_id, COUNT(*) FROM employee GROUP BY dept_id, name;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 9: SELECT statement with non-indexed field in ORDER BY clause",
		"SELECT emp_id, name FROM employee ORDER BY age ASC;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 10: SELECT statement with mixed indexed and non-indexed fields in WHERE, GROUP BY, ORDER BY clauses",
		"SELECT emp_id, name FROM employee WHERE dept_id = 10 AND age > 30 GROUP BY emp_id ORDER BY salary DESC;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT, salary DECIMAL(10, 2), INDEX idx_dept_id (dept_id));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 11: UPDATE statement with indexed field in WHERE clause",
		"UPDATE employee SET salary = salary * 1.1 WHERE emp_id = 123;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), age INT, salary DECIMAL(10, 2));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 12: UPDATE statement with non-indexed field in WHERE clause",
		"UPDATE employee SET salary = salary * 1.1 WHERE name = 'Bob';",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), age INT, salary DECIMAL(10, 2));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 13: UPDATE statement with mixed indexed and non-indexed fields in WHERE clause",
		"UPDATE employee SET salary = salary * 1.1 WHERE emp_id = 123 AND age > 30;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), age INT, salary DECIMAL(10, 2));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 14: UNION statement with sub-select using indexed fields in WHERE clause",
		"SELECT emp_id, name FROM employee WHERE dept_id = 10 UNION SELECT emp_id, name FROM employee_archive WHERE dept_id = 20;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT, INDEX idx_dept_id (dept_id));").
			WithSQL("CREATE TABLE employee_archive (emp_id INT, name VARCHAR(50), dept_id INT, INDEX idx_dept_id_archive (dept_id));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 15: UNION statement with sub-select using non-indexed fields in WHERE clause",
		"SELECT emp_id, name FROM employee WHERE dept_id = 10 UNION SELECT emp_id, name FROM employee_archive WHERE age > 30;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT, INDEX idx_dept_id (dept_id));").
			WithSQL("CREATE TABLE employee_archive (emp_id INT, name VARCHAR(50), age INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 16: UNION statement with sub-select using mixed indexed and non-indexed fields in WHERE clause",
		"SELECT emp_id, name FROM employee WHERE dept_id = 10 UNION SELECT emp_id, name FROM employee_archive WHERE emp_id = 456 AND salary > 50000;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE employee (emp_id INT PRIMARY KEY, name VARCHAR(50), dept_id INT, age INT, INDEX idx_dept_id (dept_id));").
			WithSQL("CREATE TABLE employee_archive (emp_id INT, name VARCHAR(50), salary DECIMAL(10, 2), INDEX idx_emp_id_archive (emp_id));"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 17: SELECT statement with indexed field in WHERE clause (从xml中补充)",
		"SELECT COUNT(*) FROM customers WHERE age = 20;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), age INT, INDEX idx_age (age));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 18: SELECT statement with non-indexed field in WHERE clause (从xml中补充)",
		"SELECT * FROM customers WHERE name LIKE '小王1';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), age INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 19: UPDATE statement with indexed field in WHERE clause (从xml中补充)",
		"UPDATE customers SET city='beijing' WHERE age=20;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), age INT, city VARCHAR(50), INDEX idx_age (age));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 20: DELETE statement with indexed field in WHERE clause (从xml中补充)",
		"DELETE FROM customers WHERE age = 20;",
		session.NewAIMockContext().
			WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), age INT, INDEX idx_age (age));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 21: SELECT statement with non-indexed field in ORDER BY clause (从xml中补充)",
		"SELECT * FROM customers ORDER BY age DESC LIMIT 1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(50), age INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)
}

// ==== Rule test code end ====
