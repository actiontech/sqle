package mysql

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/stretchr/testify/assert"
)

// ==== Rule test code start ====
func NewMySQLInspectOnRuleSQLE00039(t *testing.T, colNames []string, mockDiscrimination ...driver.Value) *MysqlDriverImpl {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect := NewMockInspect(e)

	selectExpr, selectColumns := []string{}, []string{}

	for _, column := range colNames {
		selectExpr = append(
			selectExpr,
			fmt.Sprintf("COUNT( DISTINCT ( `%v` ) ) / COUNT( * ) * 100 AS '%v'", column, column),
		)
		selectColumns = append(selectColumns, "`"+column+"`")
	}

	handler.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(
		"SELECT %v FROM (SELECT %v FROM `exist_db`.`exist_tb_3` LIMIT 50000) t;",
		strings.Join(selectExpr, ","),
		strings.Join(selectColumns, ","),
	))).WillReturnRows(
		sqlmock.NewRows(colNames).AddRow(mockDiscrimination...),
	)

	return inspect
}

func TestRuleSQLE00039(t *testing.T) {
	ruleName := ai.SQLE00039
	rule := rulepkg.RuleHandlerMap[ruleName].Rule
	ruleParam := 70
	i := NewMySQLInspectOnRuleSQLE00039(t, []string{"v1"}, "100.0000")

	//create index, with index, no problem (index discrimination is greater than the threshold)
	runSingleRuleInspectCase(rule, t, "create index, with index, no problem (index discrimination is greater than the threshold)", i,
		`
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v1);
    `, newTestResult())

	i = NewMySQLInspectOnRuleSQLE00039(t, []string{"v3"}, "30.0000")
	//create index, with index, problem (index discrimination is less than the threshold)
	runSingleRuleInspectCase(rule, t, "create index, with index, problem (index discrimination is less than the threshold)", i,
		`
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v3);
    `, newTestResult().addResult(ruleName, "v3", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, []string{"v1", "v3"}, "100.0000", "30.0000")
	//create index, with index, problem (one of indexes discrimination is less than the threshold)
	runSingleRuleInspectCase(rule, t, "create index, with index, problem (one of indexes discrimination is less than the threshold)", i,
		`
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v1, v3);
    `, newTestResult().addResult(ruleName, "v3", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, []string{"v2"}, "70.0000")
	//create index, with index, no problem (index discrimination is equal with the threshold)
	runSingleRuleInspectCase(rule, t, "create index, with index, problem (index discrimination is equal with the threshold)", i,
		`
    CREATE INDEX idx_1 on exist_db.exist_tb_3(v2);
    `, newTestResult())

	i = NewMySQLInspectOnRuleSQLE00039(t, []string{"v1"}, "100.0000")
	//ALTER index, with index, no problem (index discrimination is greater than the threshold)
	runSingleRuleInspectCase(rule, t, "ALTER index, with index, no problem (index discrimination is greater than the threshold)", i,
		`
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX(v1);
    `, newTestResult())

	i = NewMySQLInspectOnRuleSQLE00039(t, []string{"v3"}, "30.0000")
	//ALTER index, with index, problem (index discrimination is less than the threshold)
	runSingleRuleInspectCase(rule, t, "ALTER index, with index, problem (index discrimination is less than the threshold)", i,
		`
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX(v3);
    `, newTestResult().addResult(ruleName, "v3", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, []string{"v1", "v3"}, "100.0000", "30.0000")
	//ALTER index, with index, problem (one of indexes discrimination is less than the threshold)
	runSingleRuleInspectCase(rule, t, "ALTER index, with index, problem (one of indexes discrimination is less than the threshold)", i,
		`
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX(v1, v3);
    `, newTestResult().addResult(ruleName, "v3", ruleParam))

	i = NewMySQLInspectOnRuleSQLE00039(t, []string{"v2"}, "70.0000")
	//ALTER index, with index, no problem (index discrimination is equal with the threshold)
	runSingleRuleInspectCase(rule, t, "ALTER index, with index, problem (index discrimination is equal with the threshold)", i,
		`
    ALTER TABLE exist_db.exist_tb_3 ADD INDEX(v2);
    `, newTestResult())

}

// ==== Rule test code end ====
