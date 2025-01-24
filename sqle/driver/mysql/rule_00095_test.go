package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQL00095(t *testing.T) {
	ruleName := ai.SQLE00095
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: DELETE 语句中使用 != 操作符", "DELETE FROM employees WHERE salary != 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 4000); INSERT INTO employees VALUES (2, 'Jane Doe', 6000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM employees WHERE salary != 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: DELETE 语句中使用 <> 操作符", "DELETE FROM employees WHERE salary <> 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 4000); INSERT INTO employees VALUES (2, 'Jane Doe', 6000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM employees WHERE salary <> 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 3: DELETE 语句中没有使用不等于操作符", "DELETE FROM employees WHERE salary = 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 5000); INSERT INTO employees VALUES (2, 'Jane Doe', 5000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM employees WHERE salary = 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 4: INSERT 语句中使用 != 操作符", "INSERT INTO employees (id, name, salary) SELECT id, name, salary FROM temp_employees WHERE salary != 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); CREATE TABLE temp_employees (id INT, name VARCHAR(50), salary INT); INSERT INTO temp_employees VALUES (1, 'John Doe', 4000); INSERT INTO temp_employees VALUES (2, 'Jane Doe', 6000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO employees (id, name, salary) SELECT id, name, salary FROM temp_employees WHERE salary != 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: INSERT 语句中使用 <> 操作符", "INSERT INTO employees (id, name, salary) SELECT id, name, salary FROM temp_employees WHERE salary <> 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); CREATE TABLE temp_employees (id INT, name VARCHAR(50), salary INT); INSERT INTO temp_employees VALUES (1, 'John Doe', 4000); INSERT INTO temp_employees VALUES (2, 'Jane Doe', 6000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO employees (id, name, salary) SELECT id, name, salary FROM temp_employees WHERE salary <> 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 6: INSERT 语句中没有使用不等于操作符", "INSERT INTO employees (id, name, salary) SELECT id, name, salary FROM temp_employees WHERE salary = 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); CREATE TABLE temp_employees (id INT, name VARCHAR(50), salary INT); INSERT INTO temp_employees VALUES (1, 'John Doe', 5000); INSERT INTO temp_employees VALUES (2, 'Jane Doe', 5000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN INSERT INTO employees (id, name, salary) SELECT id, name, salary FROM temp_employees WHERE salary = 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 7: UPDATE 语句中使用 != 操作符", "UPDATE employees SET salary = salary * 1.1 WHERE salary != 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 4000); INSERT INTO employees VALUES (2, 'Jane Doe', 6000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE employees SET salary = salary * 1.1 WHERE salary != 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: UPDATE 语句中使用 <> 操作符", "UPDATE employees SET salary = salary * 1.1 WHERE salary <> 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 4000); INSERT INTO employees VALUES (2, 'Jane Doe', 6000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE employees SET salary = salary * 1.1 WHERE salary <> 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 9: UPDATE 语句中没有使用不等于操作符", "UPDATE employees SET salary = salary * 1.1 WHERE salary = 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 5000); INSERT INTO employees VALUES (2, 'Jane Doe', 5000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE employees SET salary = salary * 1.1 WHERE salary = 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 10: SELECT 语句中使用 != 操作符", "SELECT * FROM employees WHERE salary != 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 4000); INSERT INTO employees VALUES (2, 'Jane Doe', 6000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM employees WHERE salary != 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: SELECT 语句中使用 <> 操作符", "SELECT * FROM employees WHERE salary <> 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 4000); INSERT INTO employees VALUES (2, 'Jane Doe', 6000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM employees WHERE salary <> 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 12: SELECT 语句中没有使用不等于操作符", "SELECT * FROM employees WHERE salary = 5000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (id INT, name VARCHAR(50), salary INT); INSERT INTO employees VALUES (1, 'John Doe', 5000); INSERT INTO employees VALUES (2, 'Jane Doe', 5000);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM employees WHERE salary = 5000;",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 13: SELECT 语句中使用 != 操作符在 customers 表", "SELECT * FROM customers WHERE name != '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM customers WHERE name != '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: SELECT 语句中使用 <> 操作符在 customers 表", "SELECT * FROM customers WHERE name <> '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM customers WHERE name <> '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 15: SELECT 语句中没有使用不等于操作符在 customers 表", "SELECT * FROM customers WHERE name = '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN SELECT * FROM customers WHERE name = '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 16: UPDATE 语句中使用 != 操作符在 customers 表", "UPDATE customers SET age = 29 WHERE name != '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE customers SET age = 29 WHERE name != '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: UPDATE 语句中使用 <> 操作符在 customers 表", "UPDATE customers SET age = 29 WHERE name <> '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE customers SET age = 29 WHERE name <> '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 18: UPDATE 语句中没有使用不等于操作符在 customers 表", "UPDATE customers SET age = 29 WHERE name = '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN UPDATE customers SET age = 29 WHERE name = '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 19: DELETE 语句中使用 != 操作符在 customers 表", "DELETE FROM customers WHERE name != '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM customers WHERE name != '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 20: DELETE 语句中使用 <> 操作符在 customers 表", "DELETE FROM customers WHERE name <> '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM customers WHERE name <> '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 21: DELETE 语句中没有使用不等于操作符在 customers 表", "DELETE FROM customers WHERE name = '小青';",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers(id INT NOT NULL, name VARCHAR(32) DEFAULT 'lucy', sex int NOT NULL default 0, city VARCHAR(32) NOT NULL default 'beijing', age INT NOT NULL default 0, PRIMARY KEY (id)); INSERT INTO customers VALUES (1,'小青',0,'杭州',25); INSERT INTO customers VALUES (2,'小白',0,'杭州',25);"),
		[]*AIMockSQLExpectation{
			{
				Query: "EXPLAIN DELETE FROM customers WHERE name = '小青';",
				Rows:  sqlmock.NewRows(nil),
			},
		}, newTestResult())
}

// ==== Rule test code end ====
