package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00067(t *testing.T) {
	ruleName := ai.SQLE00067
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE with foreign key constraints",
		"CREATE TABLE orders (order_id INT, customer_id INT, PRIMARY KEY(order_id), FOREIGN KEY (customer_id) REFERENCES customers(customer_id));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: CREATE TABLE without foreign key constraints",
		"CREATE TABLE customers (customer_id INT PRIMARY KEY, customer_name VARCHAR(100));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 3: CREATE TABLE ...",
		"CREATE TABLE order_items (order_id INT, item_id INT, PRIMARY KEY(order_id, item_id));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 4: ALTER TABLE adding a foreign key constraint",
		"ALTER TABLE orders ADD CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers(customer_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT, customer_id INT, PRIMARY KEY(order_id));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: ALTER TABLE add columns without foreign keys",
		"ALTER TABLE orders ADD COLUMN order_date DATE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE orders (order_id INT, customer_id INT, PRIMARY KEY(order_id));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 6: CREATE TABLE ... using unique constraint and no foreign keys",
		"CREATE TABLE products (product_id INT PRIMARY KEY, product_name VARCHAR(100) UNIQUE);",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 7: CREATE TABLE .... using CHECK constraints and no foreign keys",
		"CREATE TABLE employees (employee_id INT PRIMARY KEY, age INT CHECK (age >= 18));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 8: ALTER TABLE ... add unique constraint and no foreign key",
		"ALTER TABLE products ADD CONSTRAINT unique_product_name UNIQUE (product_name);",
		session.NewAIMockContext().WithSQL("CREATE TABLE products (product_id INT PRIMARY KEY, product_name VARCHAR(100));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 9: CREATE TABLE with a foreign key constraint",
		"CREATE TABLE customers_fk1 (id INT, customer_id INT, FOREIGN KEY (customer_id) REFERENCES customers(id));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: CREATE TABLE .... with a foreign key constraint referencing the customers table and cascading deletes",
		"CREATE TABLE customers_fk2 (id INT, customer_id INT, FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE);",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: CREATE TABLE .... with a foreign key constraint referencing the customers table and cascading deletes",
		"CREATE TABLE customers_fk3 (id INT, customer_id INT, FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE);",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: ALTER TABLE customers_fk1 adding a foreign key constraint",
		"ALTER TABLE customers_fk1 ADD CONSTRAINT fk1_id_customers FOREIGN KEY (id) REFERENCES customers(id) ON DELETE CASCADE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers_fk1 (id INT, customer_id INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: ALTER TABLE customers_fk2 adding a foreign key constraint",
		"ALTER TABLE customers_fk2 ADD CONSTRAINT fk2_id_customers FOREIGN KEY (id) REFERENCES customers(id) ON DELETE CASCADE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers_fk2 (id INT, customer_id INT);"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: ALTER TABLE customers_fk3 adding a foreign key constraint",
		"ALTER TABLE customers_fk3 ADD CONSTRAINT fk3_id_customers FOREIGN KEY (id) REFERENCES customers(id) ON DELETE CASCADE;",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers_fk3 (id INT, customer_id INT);"),
		nil,
		newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
