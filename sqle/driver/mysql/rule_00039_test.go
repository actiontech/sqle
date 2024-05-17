package mysql

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/stretchr/testify/assert"
)

// ==== Rule test code start ====
func NewMySQLInspectOnRuleSQLE00039(t *testing.T, showIndex bool, colNames []string, mockDiscrimination ...driver.Value) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)

	if showIndex {
		handler.ExpectQuery(regexp.QuoteMeta("SHOW INDEX FROM `exist_db`.`exist_tb_3`")).
			WillReturnRows(sqlmock.NewRows([]string{"Column_name"}).AddRow("id").AddRow("v1"))
	}

	handler.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) AS total FROM `exist_db`.`exist_tb_3` LIMIT 50000")).WillReturnRows(
		sqlmock.NewRows([]string{"total"}).AddRow(50000),
	)

	for i, column := range colNames {

		handler.ExpectQuery(regexp.QuoteMeta(
			fmt.Sprintf("SELECT COUNT(*) AS record_count FROM (SELECT `%s` FROM `exist_db`.`exist_tb_3` LIMIT 50000) AS limited GROUP BY `%s` ORDER BY record_count DESC LIMIT 1",
				column, column))).WillReturnRows(
			sqlmock.NewRows([]string{"record_count"}).AddRow(mockDiscrimination[i]),
		)
	}

	return inspect
}

func TestRuleSQLE00039(t *testing.T) {
	ruleName := ai.SQLE00039
	rule := rulepkg.RuleHandlerMap[ruleName].Rule
	ruleParam := 0.7
	i := NewMySQLInspectOnRuleSQLE00039(t, false, []string{"v1"}, "100")

	//create index, with index, no problem (index discrimination is greater than the threshold)
	runSingleRuleInspectCase(rule, t, "create index, with index, no problem (index discrimination is greater than the threshold)", i,
		`
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v1);
    `, newTestResult())

	i = NewMySQLInspectOnRuleSQLE00039(t, false, []string{"v3"}, "40000")
	//create index, with index, problem (index discrimination is less than the threshold)
	runSingleRuleInspectCase(rule, t, "create index, with index, problem (index discrimination is less than the threshold)", i,
		`
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v3);
    `, newTestResult().addResult(ruleName, "v3", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, false, []string{"v1", "v3"}, "100", "40000")
	//create index, with index, problem (one of indexes discrimination is less than the threshold)
	runSingleRuleInspectCase(rule, t, "create index, with index, problem (one of indexes discrimination is less than the threshold)", i,
		`
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v1, v3);
    `, newTestResult().addResult(ruleName, "v3", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, false, []string{"v2"}, "15000")
	//create index, with index, no problem (index discrimination is equal with the threshold)
	runSingleRuleInspectCase(rule, t, "create index, with index, problem (index discrimination is equal with the threshold)", i,
		`
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v2);
    `, newTestResult())

	i = NewMySQLInspectOnRuleSQLE00039(t, false, []string{"v1"}, "100")
	//ALTER index, with index, no problem (index discrimination is greater than the threshold)
	runSingleRuleInspectCase(rule, t, "ALTER index, with index, no problem (index discrimination is greater than the threshold)", i,
		`
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX(v1);
    `, newTestResult())

	i = NewMySQLInspectOnRuleSQLE00039(t, false, []string{"v3"}, "40000")
	//ALTER index, with index, problem (index discrimination is less than the threshold)
	runSingleRuleInspectCase(rule, t, "ALTER index, with index, problem (index discrimination is less than the threshold)", i,
		`
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX(v3);
    `, newTestResult().addResult(ruleName, "v3", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, false, []string{"v1", "v3"}, "100", "40000")
	//ALTER index, with index, problem (one of indexes discrimination is less than the threshold)
	runSingleRuleInspectCase(rule, t, "ALTER index, with index, problem (one of indexes discrimination is less than the threshold)", i,
		`
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX(v1, v3);
    `, newTestResult().addResult(ruleName, "v3", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, false, []string{"v2"}, "15000")
	//ALTER index, with index, no problem (index discrimination is equal with the threshold)
	runSingleRuleInspectCase(rule, t, "ALTER index, with index, problem (index discrimination is equal with the threshold)", i,
		`
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX(v2);
    `, newTestResult())

	i = NewMySQLInspectOnRuleSQLE00039(t, true, []string{"v2"}, "40000")
	//select...where, with no index
	runSingleRuleInspectCase(rule, t, "select...where, with no index", i, `
	SELECT * FROM exist_db.exist_tb_3 WHERE v2 != "1";
	`, newTestResult())

	i = NewMySQLInspectOnRuleSQLE00039(t, true, []string{"v1"}, "40000")
	// select...where, with problem (index field type is not the expected)  selectivity < threshold
	runSingleRuleInspectCase(rule, t, "select...where, with problem (index field type is not the expected)  selectivity < threshold", i, `
	SELECT * FROM exist_db.exist_tb_3 WHERE v1 != "1";
	`, newTestResult().addResult(ruleName, "v1", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, true, []string{"v1"}, "100")
	// select...where, with problem (index field type is not the expected)  selectivity >= threshold
	runSingleRuleInspectCase(rule, t, "select...where, with problem (index field type is not the expected)  selectivity >= threshold", i, `
	SELECT * FROM exist_db.exist_tb_3 WHERE v1!= "1";
	`, newTestResult())

}

// ==== Rule test code end ====
