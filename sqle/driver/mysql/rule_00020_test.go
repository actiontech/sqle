package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00020(t *testing.T) {
	ruleName := ai.SQLE00020
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE 定义39列，符合规则",
		"CREATE TABLE test_table_39 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT);",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 2: CREATE TABLE 定义40列，符合规则",
		"CREATE TABLE test_table_40 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT, col40 INT);",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 3: CREATE TABLE 定义41列，违反规则",
		"CREATE TABLE test_table_41 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT, col40 INT, col41 INT);",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 4: ALTER TABLE 向已有39列的表添加1列，符合规则",
		"ALTER TABLE test_table_39 ADD COLUMN col40 INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table_39 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT);"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 5: ALTER TABLE 向已有39列的表添加2列，违反规则",
		"ALTER TABLE test_table_39 ADD COLUMN col40 INT, ADD COLUMN col41 INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table_39 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 6: ALTER TABLE 向已有40列的表添加1列且删除1列，符合规则",
		"ALTER TABLE test_table_40 ADD COLUMN col41 INT, DROP COLUMN col1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table_40 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT, col40 INT);"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 7: ALTER TABLE 向已有35列的表删除5列并添加10列，违反规则",
		"ALTER TABLE test_table_35 DROP COLUMN col1, DROP COLUMN col2, DROP COLUMN col3, DROP COLUMN col4, DROP COLUMN col5, ADD COLUMN col36 INT, ADD COLUMN col37 INT, ADD COLUMN col38 INT, ADD COLUMN col39 INT, ADD COLUMN col40 INT, ADD COLUMN col41 INT, ADD COLUMN col42 INT, ADD COLUMN col43 INT, ADD COLUMN col44 INT, ADD COLUMN col45 INT, ADD COLUMN col46 INT, ADD COLUMN col47 INT, ADD COLUMN col48 INT, ADD COLUMN col49 INT, ADD COLUMN col50 INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table_35 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 8: ALTER TABLE 向已有36列的表添加5列，违反规则",
		"ALTER TABLE test_table_36 ADD COLUMN col37 INT, ADD COLUMN col38 INT, ADD COLUMN col39 INT, ADD COLUMN col40 INT, ADD COLUMN col41 INT;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table_36 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT);"),
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 9: ALTER TABLE 向已有40列的表删除1列，符合规则",
		"ALTER TABLE test_table_40 DROP COLUMN col1;",
		session.NewAIMockContext().WithSQL("CREATE TABLE test_table_40 (col1 INT, col2 INT, col3 INT, col4 INT, col5 INT, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT, col40 INT);"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 10: CREATE TABLE 定义40列，符合规则 (从xml中补充)",
		"CREATE TABLE order_table_oltp (order_id INT, order_name VARCHAR(255), sales_id INT, order_create_time DATETIME, order_end_time DATETIME, col6 INT, col7 INT, col8 INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT, col40 INT);",
		nil,
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 11: CREATE TABLE 定义41列，违反规则 (从xml中补充)",
		"CREATE TABLE order_table_olap (order_id INT, order_name VARCHAR(255), sales_id INT, order_create_time DATETIME, order_end_time DATETIME, department_name VARCHAR(255), sales_name VARCHAR(255), sales_department_id INT, col9 INT, col10 INT, col11 INT, col12 INT, col13 INT, col14 INT, col15 INT, col16 INT, col17 INT, col18 INT, col19 INT, col20 INT, col21 INT, col22 INT, col23 INT, col24 INT, col25 INT, col26 INT, col27 INT, col28 INT, col29 INT, col30 INT, col31 INT, col32 INT, col33 INT, col34 INT, col35 INT, col36 INT, col37 INT, col38 INT, col39 INT, col40 INT, col41 INT);",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)
}

// ==== Rule test code end ====
