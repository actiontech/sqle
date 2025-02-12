package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00039(t *testing.T) {
	ruleName := ai.SQLE00039
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE INDEX on table_users.username with discrimination above threshold",
		"CREATE INDEX idx_username ON table_users(username);",
		session.NewAIMockContext().WithSQL("CREATE TABLE table_users (id INT, username VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `username` ) ) / COUNT( * ) * 100 AS 'username' FROM (SELECT `username` FROM `exist_db`.`table_users` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"username"}).AddRow(1),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE INDEX on table_orders.order_id with discrimination below threshold",
		"CREATE INDEX idx_order_id ON table_orders(order_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE table_orders (order_id INT, order_date DATE);"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `order_id` ) ) / COUNT( * ) * 100 AS 'order_id' FROM (SELECT `order_id` FROM `exist_db`.`table_orders` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"order_id"}).AddRow(0.5),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: ALTER TABLE table_products ADD INDEX idx_product_code(product_code) with discrimination above threshold",
		"ALTER TABLE table_products ADD INDEX idx_product_code(product_code);",
		session.NewAIMockContext().WithSQL("CREATE TABLE table_products (product_id INT, product_code VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `product_code` ) ) / COUNT( * ) * 100 AS 'product_code' FROM (SELECT `product_code` FROM `exist_db`.`table_products` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"product_code"}).AddRow(1),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 4: ALTER TABLE table_customers ADD INDEX idx_email(email) with discrimination below threshold",
		"ALTER TABLE table_customers ADD INDEX idx_email(email);",
		session.NewAIMockContext().WithSQL("CREATE TABLE table_customers (customer_id INT, email VARCHAR(100));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `email` ) ) / COUNT( * ) * 100 AS 'email' FROM (SELECT `email` FROM `exist_db`.`table_customers` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"email"}).AddRow(0.5),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: SELECT query on table_employees with indexed field salary having discrimination above threshold",
		"SELECT * FROM table_employees WHERE salary > 50000;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table_employees (employee_id INT, salary DECIMAL(10, 2), INDEX idx_salary(salary));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `salary` ) ) / COUNT( * ) * 100 AS 'salary' FROM (SELECT `salary` FROM `exist_db`.`table_employees` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"salary"}).AddRow(1),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 6: SELECT query on table_inventory with indexed field sku having discrimination below threshold",
		"SELECT * FROM table_inventory WHERE sku = 'ABC123';",
		session.NewAIMockContext().WithSQL("CREATE TABLE table_inventory (inventory_id INT, sku VARCHAR(50), INDEX idx_sku(sku));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `sku` ) ) / COUNT( * ) * 100 AS 'sku' FROM (SELECT `sku` FROM `exist_db`.`table_inventory` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"sku"}).AddRow(0.6),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: UPDATE statement on table_accounts with indexed field account_id having discrimination above threshold",
		"UPDATE table_accounts SET status = 'active' WHERE account_id = 1001;",
		session.NewAIMockContext().WithSQL("CREATE TABLE table_accounts (account_id INT, status VARCHAR(20), INDEX idx_account_id(account_id));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `account_id` ) ) / COUNT( * ) * 100 AS 'account_id' FROM (SELECT `account_id` FROM `exist_db`.`table_accounts` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"account_id"}).AddRow(1),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 9: ALTER TABLE customers ADD INDEX idx_age_customers(age) with discrimination below threshold (从xml中补充)",
		"ALTER TABLE customers ADD INDEX idx_age_customers(age);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT, age INT, name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `age` ) ) / COUNT( * ) * 100 AS 'age' FROM (SELECT `age` FROM `exist_db`.`customers` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"age"}).AddRow(0.1),
			},
		}, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: CREATE INDEX on customers.name with discrimination above threshold (从xml中补充)",
		"CREATE INDEX idx_name_customers ON customers(name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT, name VARCHAR(50));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `name` ) ) / COUNT( * ) * 100 AS 'name' FROM (SELECT `name` FROM `exist_db`.`customers` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"name"}).AddRow(1),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 11: SELECT query on customers with indexed field name having discrimination above threshold (从xml中补充)",
		"SELECT * FROM customers WHERE name = '小王22222333' AND age < 50;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT, name VARCHAR(50), age INT, INDEX idx_name(name));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `name` ) ) / COUNT( * ) * 100 AS 'name' FROM (SELECT `name` FROM `exist_db`.`customers` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"name"}).AddRow(1),
			},
		}, newTestResult())

	runAIRuleCase(rule, t, "case 12: SELECT query on customers with indexed field age having discrimination below threshold (从xml中补充)",
		"SELECT /*+ index(customers idx_age_customers) */ * FROM customers WHERE name = '小王22222333' AND age < 50;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (customer_id INT, name VARCHAR(50), age INT, INDEX idx_age(age));"),
		[]*AIMockSQLExpectation{
			{
				Query: "SELECT COUNT( DISTINCT ( `age` ) ) / COUNT( * ) * 100 AS 'age' FROM (SELECT `age` FROM `exist_db`.`customers` LIMIT 50000) t;",
				Rows:  sqlmock.NewRows([]string{"age"}).AddRow(0.1),
			},
		}, newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
