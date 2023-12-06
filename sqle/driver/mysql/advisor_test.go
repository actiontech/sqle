package mysql

import (
	"fmt"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var explainColumns []string = []string{"type", "table"}

const explainTypeAll string = "ALL"
const drivingTable string = "exist_tb_1"
const explainFormat string = "EXPLAIN %s"

func mockPrefixIndexOptimizeResult(caseName string, c optimizerTestContent, t *testing.T) []*OptimizeResult {
	return mockOptimizeResultWithAdvisor(c.sql, c.maxColumn, c.queryResults, caseName, t, newPrefixIndexAdvisor)
}

func mockThreeStarOptimizeResult(caseName string, c optimizerTestContent, t *testing.T) []*OptimizeResult {
	return mockOptimizeResultWithAdvisor(c.sql, c.maxColumn, c.queryResults, caseName, t, newThreeStarIndexAdvisor)
}

func mockFunctionOptimizeResult(caseName string, c optimizerTestContent, t *testing.T) []*OptimizeResult {
	return mockOptimizeResultWithAdvisor(c.sql, c.maxColumn, c.queryResults, caseName, t, newFunctionIndexAdvisor)
}

func mockExtremalOptimizeResult(caseName string, c optimizerTestContent, t *testing.T) []*OptimizeResult {
	return mockOptimizeResultWithAdvisor(c.sql, c.maxColumn, c.queryResults, caseName, t, newExtremalIndexAdvisor)
}

func mockJoinOptimizeResult(caseName string, c optimizerTestContent, t *testing.T) []*OptimizeResult {
	return mockOptimizeResultWithAdvisor(c.sql, c.maxColumn, c.queryResults, caseName, t, newJoinIndexAdvisor)
}

func mockOptimizeResultWithAdvisor(sql string, maxColumn int, queryResults []*queryResult, caseName string, t *testing.T, f func(ctx *session.Context, log *logrus.Entry, originNode ast.Node, params params.Params) CreateIndexAdvisor) []*OptimizeResult {
	e, handler, err := executor.NewMockExecutor()
	assert.NoErrorf(t, err, caseName)
	for _, expect := range queryResults {
		handler.ExpectQuery(expect.query).WillReturnRows(expect.result)
	}

	impl := NewMockInspectWithIsExecutedSQL(e)
	assert.NoErrorf(t, err, caseName)

	node, err := util.ParseOneSql(sql)
	assert.NoErrorf(t, err, caseName)
	assert.Truef(t, canOptimize(impl.log, impl.Ctx, node), caseName)

	advisor := f(impl.Ctx, impl.log, node, params.Params{
		{
			Key:   MAX_INDEX_COLUMN,
			Value: fmt.Sprint(maxColumn),
			Type:  params.ParamTypeInt,
		},
	})
	return advisor.GiveAdvices()
}

type optimizerTestContent struct {
	sql           string
	maxColumn     int
	queryResults  []*queryResult
	expectResults []*OptimizeResult
}

type queryResult struct {
	query  string
	result *sqlmock.Rows
}

type optimizerTestCaseMap map[string] /*case name*/ optimizerTestContent

func (testCases optimizerTestCaseMap) testAll(testFunc func(caseName string, c optimizerTestContent, t *testing.T) []*OptimizeResult, t *testing.T) {
	for caseName := range testCases {
		testCases.testOne(caseName, testFunc, t)
	}
}

// 当需要单独测试一个测试用例时使用
func (testCases optimizerTestCaseMap) testOne(caseName string, testFunc func(caseName string, c optimizerTestContent, t *testing.T) []*OptimizeResult, t *testing.T) {
	c := testCases[caseName]
	result := testFunc(caseName, c, t)
	assert.Equalf(t, c.expectResults, result, caseName)
}

func TestPrefixIndexOptimize(t *testing.T) {
	testCases := make(optimizerTestCaseMap)
	testCases["test1"] = optimizerTestContent{
		sql: `SELECT * FROM exist_tb_1 WHERE v1 LIKE "%_set"`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT * FROM exist_tb_1 WHERE v1 LIKE "%_set"`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：`v1` LIKE '%_set' 中，使用了前缀模式匹配，在数据量大的时候，可以建立翻转函数索引",
				IndexedColumns: []string{"v1"},
				TableName:      "exist_tb_1",
			},
		},
		maxColumn: 1,
	}
	testCases["test2"] = optimizerTestContent{
		sql: `SELECT * FROM exist_tb_1 WHERE v1 LIKE upper("_set")`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT * FROM exist_tb_1 WHERE v1 LIKE upper("_set")`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：`v1` LIKE UPPER('_set') 中，使用了前缀模式匹配，在数据量大的时候，可以建立翻转函数索引",
				IndexedColumns: []string{"v1"},
				TableName:      "exist_tb_1",
			},
		},
		maxColumn: 1,
	}
	testCases["test3"] = optimizerTestContent{
		sql: `SELECT * FROM exist_tb_1 WHERE v1 = '_set'`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT * FROM exist_tb_1 WHERE v1 = '_set'`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			},
		},
		maxColumn: 1,
	}
	testCases.testAll(mockPrefixIndexOptimizeResult, t)
}

func TestFunctionIndexOptimize(t *testing.T) {
	testCases := make(optimizerTestCaseMap)
	testCases["test1"] = optimizerTestContent{
		sql: `SELECT v1 FROM exist_tb_1 WHERE LOWER(v1) = "s"`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT v1 FROM exist_tb_1 WHERE LOWER(v1) = "s"`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			}, {
				query:  regexp.QuoteMeta("SHOW GLOBAL VARIABLES LIKE 'version'"),
				result: sqlmock.NewRows([]string{"Value"}).AddRow("8.0.13"),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：LOWER(`v1`) 中，使用了函数作为查询条件，在MySQL8.0.13以上的版本，可以创建函数索引",
				IndexedColumns: []string{"v1"},
				TableName:      "exist_tb_1",
			},
		},
		maxColumn: 4,
	}
	testCases["test2"] = optimizerTestContent{
		sql: `SELECT v1 FROM exist_tb_1 WHERE LOWER(v1) = "s"`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT v1 FROM exist_tb_1 WHERE LOWER(v1) = "s"`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			}, {
				query:  regexp.QuoteMeta("SHOW GLOBAL VARIABLES LIKE 'version'"),
				result: sqlmock.NewRows([]string{"Value"}).AddRow("5.7.1"),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：LOWER(`v1`) 中，使用了函数作为查询条件，在MySQL5.7以上的版本，可以在虚拟列上创建索引",
				IndexedColumns: []string{"v1"},
				TableName:      "exist_tb_1",
			},
		},
		maxColumn: 4,
	}
	testCases["test3"] = optimizerTestContent{
		sql: `SELECT v1 FROM exist_tb_1 WHERE LOWER(v1) = "s"`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT v1 FROM exist_tb_1 WHERE LOWER(v1) = "s"`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			}, {
				query:  regexp.QuoteMeta("SHOW GLOBAL VARIABLES LIKE 'version'"),
				result: sqlmock.NewRows([]string{"Value"}).AddRow("5.2.1"),
			},
		},
		maxColumn: 4,
	}
	testCases["test4"] = optimizerTestContent{
		sql: `SELECT v1 FROM exist_tb_1 WHERE v1 = "s"`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT v1 FROM exist_tb_1 WHERE v1 = "s"`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			}, {
				query:  regexp.QuoteMeta("SHOW GLOBAL VARIABLES LIKE 'version'"),
				result: sqlmock.NewRows([]string{"Value"}).AddRow("5.2.1"),
			},
		},
		maxColumn: 4,
	}
	testCases.testAll(mockFunctionOptimizeResult, t)
}

func TestThreeStarOptimize(t *testing.T) {
	testCases := make(optimizerTestCaseMap)
	testCases["test1"] = optimizerTestContent{
		sql: `SELECT v1,v2 FROM exist_tb_3 WHERE v1 = "s" ORDER BY v3`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT v1,v2 FROM exist_tb_3 WHERE v1 = "s" ORDER BY v3`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "exist_tb_3"),
			}, {
				query:  regexp.QuoteMeta(`SELECT COUNT`), // 组成请求的来源包含map，会导致query的形式随机
				result: sqlmock.NewRows([]string{"v1", "v2", "v3"}).AddRow(70.12, 80.98, 34.2),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：SELECT `v1`,`v2` FROM `exist_tb_3` WHERE `v1`='s' ORDER BY `v3` 中，根据三星索引设计规范",
				IndexedColumns: []string{"v1", "v3", "v2"},
				TableName:      "exist_tb_3",
			},
		},
		maxColumn: 4,
	}
	testCases["test2-无法给出覆盖索引"] = optimizerTestContent{
		sql: `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" ORDER BY v3`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" ORDER BY v3`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "exist_tb_3"),
			}, {
				query:  regexp.QuoteMeta(`SELECT COUNT`),
				result: sqlmock.NewRows([]string{"v1", "v3", "id", "v2"}).AddRow(100.00, 23.56, 70.12, 80.98),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：SELECT `id`,`v1`,`v2`,`v3` FROM `exist_tb_3` WHERE `v1`='s' ORDER BY `v3` 中，根据三星索引设计规范",
				IndexedColumns: []string{"v1", "v3"},
				TableName:      "exist_tb_3",
			},
		},
		maxColumn: 3,
	}
	testCases["test3-给出覆盖索引"] = optimizerTestContent{
		sql: `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" ORDER BY v3`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" ORDER BY v3`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "exist_tb_3"),
			}, {
				query:  regexp.QuoteMeta(`SELECT COUNT`),
				result: sqlmock.NewRows([]string{"id", "v1", "v2", "v3"}).AddRow(100.00, 23.56, 70.12, 80.98),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：SELECT `id`,`v1`,`v2`,`v3` FROM `exist_tb_3` WHERE `v1`='s' ORDER BY `v3` 中，根据三星索引设计规范",
				IndexedColumns: []string{"v1", "v3", "id", "v2"},
				TableName:      "exist_tb_3",
			},
		},
		maxColumn: 4,
	}
	testCases["test4-排序列抢占"] = optimizerTestContent{
		sql: `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" AND v2 = "s" AND id = 20 ORDER BY v3`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" AND v2 = "s" AND id = 20 ORDER BY v3`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "exist_tb_3"),
			}, {
				query:  regexp.QuoteMeta(`SELECT COUNT`),
				result: sqlmock.NewRows([]string{"id", "v1", "v2", "v3"}).AddRow(100.00, 23.56, 70.12, 80.98),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：SELECT `id`,`v1`,`v2`,`v3` FROM `exist_tb_3` WHERE `v1`='s' AND `v2`='s' AND `id`=20 ORDER BY `v3` 中，根据三星索引设计规范",
				IndexedColumns: []string{"id", "v2", "v3"},
				TableName:      "exist_tb_3",
			},
		},
		maxColumn: 3,
	}
	testCases["test5-排序列在WHERE中"] = optimizerTestContent{
		sql: `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" AND v2 = "s" AND id = 20 ORDER BY v2`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" AND v2 = "s" AND id = 20 ORDER BY v2`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "exist_tb_3"),
			}, {
				query:  regexp.QuoteMeta(`SELECT COUNT`),
				result: sqlmock.NewRows([]string{"id", "v1", "v2", "v3"}).AddRow(100.00, 23.56, 70.12, 80.98),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：SELECT `id`,`v1`,`v2`,`v3` FROM `exist_tb_3` WHERE `v1`='s' AND `v2`='s' AND `id`=20 ORDER BY `v2` 中，根据三星索引设计规范",
				IndexedColumns: []string{"id", "v2", "v1"},
				TableName:      "exist_tb_3",
			},
		},
		maxColumn: 3,
	}
	testCases["test6-非等值列放在最后"] = optimizerTestContent{
		sql: `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" AND v2 = "s" AND id <= 20 ORDER BY v3`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT id,v1,v2,v3 FROM exist_tb_3 WHERE v1 = "s" AND v2 = "s" AND id <= 20 ORDER BY v3`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "exist_tb_3"),
			}, {
				query:  regexp.QuoteMeta(`SELECT COUNT`),
				result: sqlmock.NewRows([]string{"id", "v1", "v2", "v3"}).AddRow(100.00, 23.56, 70.12, 80.98),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：SELECT `id`,`v1`,`v2`,`v3` FROM `exist_tb_3` WHERE `v1`='s' AND `v2`='s' AND `id`<=20 ORDER BY `v3` 中，根据三星索引设计规范",
				IndexedColumns: []string{"v2", "v1", "v3", "id"},
				TableName:      "exist_tb_3",
			},
		},
		maxColumn: 4,
	}
	testCases.testAll(mockThreeStarOptimizeResult, t)
}

func TestExtremalOptimize(t *testing.T) {
	testCases := make(optimizerTestCaseMap)
	testCases["test1-v3v2都无索引"] = optimizerTestContent{
		sql: `SELECT MIN(v3),MAX(v2) FROM exist_tb_1 WHERE v1 = "s" GROUP BY v2`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT MIN(v3),MAX(v2) FROM exist_tb_1 WHERE v1 = "s" GROUP BY v2`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：MIN(`v3`) 中，使用了最值函数，可以利用索引有序的性质快速找到最值",
				IndexedColumns: []string{"v3"},
				TableName:      "exist_tb_1",
			},
			{
				Reason:         "索引建议 | SQL：MAX(`v2`) 中，使用了最值函数，可以利用索引有序的性质快速找到最值",
				IndexedColumns: []string{"v2"},
				TableName:      "exist_tb_1",
			},
		},
		maxColumn: 4,
	}
	testCases["test2-v2无索引"] = optimizerTestContent{
		sql: `SELECT v1,MIN(v2) FROM exist_tb_1 WHERE v1 = "s" GROUP BY v1`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT v1,MIN(v2) FROM exist_tb_1 WHERE v1 = "s" GROUP BY v1`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：MIN(`v2`) 中，使用了最值函数，可以利用索引有序的性质快速找到最值",
				IndexedColumns: []string{"v2"},
				TableName:      "exist_tb_1",
			},
		},
		maxColumn: 4,
	}
	testCases["test3-v1有索引v2无索引"] = optimizerTestContent{
		sql: `SELECT MAX(v1),MIN(v2) FROM exist_tb_1 WHERE v1 = "s" GROUP BY v1`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT MAX(v1),MIN(v2) FROM exist_tb_1 WHERE v1 = "s" GROUP BY v1`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			},
		},
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：MIN(`v2`) 中，使用了最值函数，可以利用索引有序的性质快速找到最值",
				IndexedColumns: []string{"v2"},
				TableName:      "exist_tb_1",
			},
		},
		maxColumn: 4,
	}
	testCases["test4-无最值函数"] = optimizerTestContent{
		sql: `SELECT v1,v2 FROM exist_tb_1 WHERE v1 = "s" GROUP BY v1`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT v1,v2 FROM exist_tb_1 WHERE v1 = "s" GROUP BY v1`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			},
		},
		maxColumn: 4,
	}
	testCases["test5-v1有索引"] = optimizerTestContent{
		sql: `SELECT MAX(v1) FROM exist_tb_1 WHERE v1 = "s" GROUP BY v1`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT MAX(v1) FROM exist_tb_1 WHERE v1 = "s" GROUP BY v1`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, drivingTable),
			},
		},
		maxColumn: 1,
	}
	testCases.testAll(mockExtremalOptimizeResult, t)
}

func TestJoinOptimize(t *testing.T) {
	testCases := make(optimizerTestCaseMap)
	testCases["test1"] = optimizerTestContent{
		sql: `SELECT t1.v1, t1.v2 FROM exist_tb_1 t1 JOIN exist_tb_2 t2 ON t1.v1 = t2.v1 WHERE t1.v1 = "s"`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT t1.v1, t1.v2 FROM exist_tb_1 t1 JOIN exist_tb_2 t2 ON t1.v1 = t2.v1 WHERE t1.v1 = "s"`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "t1").AddRow(explainTypeAll, "t2"),
			},
		},
		maxColumn: 1,
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：`exist_tb_1` AS `t1` JOIN `exist_tb_2` AS `t2` ON `t1`.`v1`=`t2`.`v1` 中，字段 v1 为被驱动表 t2 上的关联字段",
				IndexedColumns: []string{"v1"},
				TableName:      "t2",
			},
		},
	}
	testCases["test2"] = optimizerTestContent{
		sql: `SELECT t1.v1, t1.v2 FROM exist_tb_1 t1 JOIN exist_tb_2 t2 USING(v1) WHERE t1.v1 = "s"`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT t1.v1, t1.v2 FROM exist_tb_1 t1 JOIN exist_tb_2 t2 USING(v1) WHERE t1.v1 = "s"`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "t1").AddRow(explainTypeAll, "t2"),
			},
		},
		maxColumn: 1,
		expectResults: []*OptimizeResult{
			{
				Reason:         "索引建议 | SQL：`exist_tb_1` AS `t1` JOIN `exist_tb_2` AS `t2` USING (`v1`) 中，字段 v1 为被驱动表 t2 上的关联字段",
				IndexedColumns: []string{"v1"},
				TableName:      "t2",
			},
		},
	}
	testCases["test3"] = optimizerTestContent{
		sql: `SELECT t1.v1, t1.v2 FROM exist_tb_1 t1 JOIN exist_tb_2 t2 WHERE t1.v1 = "s"`,
		queryResults: []*queryResult{
			{
				query:  regexp.QuoteMeta(fmt.Sprintf(explainFormat, `SELECT t1.v1, t1.v2 FROM exist_tb_1 t1 JOIN exist_tb_2 t2 WHERE t1.v1 = "s"`)),
				result: sqlmock.NewRows(explainColumns).AddRow(explainTypeAll, "t1").AddRow(explainTypeAll, "t2"),
			},
		},
		maxColumn: 1,
	}
	testCases.testAll(mockJoinOptimizeResult, t)
}
