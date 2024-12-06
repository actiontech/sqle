//go:build enterprise
// +build enterprise

package mysql

import (
	"context"

	"testing"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/stretchr/testify/assert"
)

func testBackupOneSql(t *testing.T, caseName, querySql string, expectedBackupResults []string, expectedQueryResult []*queryResult) {
	mockExecutor, mockHandler, err := executor.NewMockExecutor()
	// 预期无错误产生
	if err != nil {
		t.Errorf("caseName: %s\n, err: %s\n", caseName, err.Error())
	}

	// 在执行过程中预期会做的请求和响应
	for _, expect := range expectedQueryResult {
		mockHandler.ExpectQuery(expect.query).WillReturnRows(expect.result)
	}
	mysqlDriverImpl := NewMockInspectWithIsExecutedSQL(mockExecutor)

	// 执行备份
	backupResults, _, err := mysqlDriverImpl.Backup(context.TODO(), BackupStrategyOriginalRow, querySql, 1000)
	// 预期无错误产生
	if err != nil {
		t.Errorf("caseName: %s\n, err: %s\n", caseName, err.Error())
	}

	// 预期备份结果和预期结果一致
	if assert.Equal(t, len(expectedBackupResults), len(backupResults), caseName) {
		for idx := range backupResults {
			if expectedBackupResults[idx] != backupResults[idx] {
				t.Errorf("caseName: %s\nexpectedBackupResults[%d]: %s\nactualBackupResults[%d]: %s\n", caseName, idx, expectedBackupResults[idx], idx, backupResults[idx])
			}
		}
	}
}

type BackupTestCase struct {
	caseName              string         // 测试用例名称
	querySql              string         // 执行的SQL语句
	expectedBackupResults []string       // 预期的备份结果
	expectedQueryResult   []*queryResult // 预期的查询结果
}

// 测试使用行备份备份Update语句
func TestBackupOriginalRowForUpdateClause(t *testing.T) {
	var selectAllRegex string = `^SELECT \* FROM`
	var allColumnsOfExistTB1 = []string{"id", "v1", "v2"}
	var allColumnsOfExistTB2 = []string{"id", "v1", "v2", "user_id"}
	var allColumnsOfExistTB3 = []string{"id", "v1", "v2", "v3"}
	testCases := []BackupTestCase{
		/* !!UNSUPPORTED CASES!!
		{
			caseName: "update with join clause in exist_tb_2",
			querySql: "UPDATE exist_tb_2 et2 JOIN exist_db.exist_tb_1 et1 ON et2.user_id = et1.id SET et2.v2 = 'joined_value' WHERE et1.v1 = 'v1';",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows([]string{"id", "v1", "v2", "user_id"}).AddRow(2, "v1", "joined_value", 1),
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_2` (`id`, `v1`, `v2`, `user_id`) VALUES ('2', 'v1', 'joined_value', '1');",
			},
		},

		{
			caseName: "multi-table update in exist_tb_1 and exist_tb_2",
			querySql: "UPDATE exist_tb_1 et1 JOIN exist_tb_2 et2 ON et1.id = et2.id SET et1.v2 = 'updated_value', et2.v2 = 'updated_value' WHERE et1.id = 1;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows([]string{"id", "v1", "v2"}).AddRow(1, "v1", "updated_value"), // exist_tb_1 更新
				},
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows([]string{"id", "v1", "v2", "user_id"}).AddRow(1, "v1", "updated_value", 1), // exist_tb_2 更新
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', 'v1', 'updated_value');",
				"REPLACE INTO `exist_db`.`exist_tb_2` (`id`, `v1`, `v2`, `user_id`) VALUES ('1', 'v1', 'updated_value', '1');",
			},
		},
		*/
		{
			caseName: "update with conditional WHERE clause in exist_tb_1",
			querySql: "UPDATE exist_tb_1 SET v2 = 'conditional_update' WHERE v1 = 'v1' AND id > 1;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(2, "v1", "conditional_update"), // 满足条件的行被更新
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('2', 'v1', 'conditional_update');",
			},
		},
		{
			caseName: "update column v1 of exist_tb_1",
			querySql: "update exist_tb_1 set v1 = 1 where id = 1;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,                                            // 该字段固定
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(1, "2", "3"), //请求的响应，填写某表的所有列，以及对应的返回值
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', '2', '3');", // 希望把这一部分提取出一个模板，填入表名、列名和值
			},
		},
		{
			caseName: "update column v2 of exist_tb_2 with set",
			querySql: "update exist_tb_2 set v2 = 'changed' where id = 2;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB2).AddRow(2, "v1", "changed", 1),
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_2` (`id`, `v1`, `v2`, `user_id`) VALUES ('2', 'v1', 'changed', '1');",
			},
		},
		{
			caseName: "update multiple columns in exist_tb_3",
			querySql: "UPDATE exist_tb_3 SET v1 = 'new_value', v3 = 10 WHERE id = 1;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB3).AddRow(1, "new_value", "old_value", 10),
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_3` (`id`, `v1`, `v2`, `v3`) VALUES ('1', 'new_value', 'old_value', '10');",
			},
		},
		{
			caseName: "update multiple rows in exist_tb_1",
			querySql: "UPDATE exist_tb_1 SET v2 = 'updated_value' WHERE id IN (1, 2);",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(1, "v1", "updated_value").AddRow(2, "v1", "updated_value"),
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', 'v1', 'updated_value');",
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('2', 'v1', 'updated_value');",
			},
		},
		{
			caseName: "update no rows in exist_tb_1",
			querySql: "UPDATE exist_tb_1 SET v2 = 'updated_value' WHERE id = 999;", // id 999 不存在，更新0行
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1), // 没有匹配的行返回空
				},
			},
			expectedBackupResults: []string{
				// 无更新，备份也不需要任何数据
			},
		},
		{
			caseName: "update using subquery in exist_tb_1",
			querySql: "UPDATE exist_tb_1 SET v2 = (SELECT v2 FROM exist_db.exist_tb_2 WHERE exist_db.exist_tb_2.id = exist_tb_1.id) WHERE id = 1;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(1, "v1", "updated_value"),
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', 'v1', 'updated_value');",
			},
		},
		{
			caseName: "update with LIMIT in exist_tb_1",
			querySql: "UPDATE exist_tb_1 SET v2 = 'limited_update' WHERE v1 = 'v1' LIMIT 1;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(1, "v1", "limited_update"), // 仅更新了第一行
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', 'v1', 'limited_update');",
			},
		},
	}
	for _, testCase := range testCases {
		testBackupOneSql(t, testCase.caseName, testCase.querySql, testCase.expectedBackupResults, testCase.expectedQueryResult)
	}
}

// 测试使用行备份备份Delete语句
func TestBackupOriginalRowForDeleteClause(t *testing.T) {
	var selectAllRegex string = `^SELECT \* FROM`
	var allColumnsOfExistTB1 = []string{"id", "v1", "v2"}
	var allColumnsOfExistTB2 = []string{"id", "v1", "v2", "user_id"}

	testCases := []BackupTestCase{
		/*  !!UNSUPPORTED CASES!!
		{
			caseName: "delete rows from exist_tb_1 and exist_tb_2 using JOIN",
			querySql: "DELETE et1, et2 FROM exist_tb_1 et1 JOIN exist_tb_2 et2 ON et1.id = et2.id WHERE et1.v1 = 'v1';",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(1, "v1", "old_value"), // 删除前查询到的结果，满足条件的行
				},
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB2).AddRow(1, "v1", "old_value", 1), // 删除前查询到的结果，满足条件的行
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', 'v1', 'old_value');", // 删除前备份的结果
				"REPLACE INTO `exist_db`.`exist_tb_2` (`id`, `v1`, `v2`, `user_id`) VALUES ('1', 'v1', 'old_value', '1');",
			},
		},

		*/
		{
			caseName: "delete multiple rows with LIMIT from exist_tb_1",
			querySql: "DELETE FROM exist_tb_1 WHERE v1 = 'v1' LIMIT 2;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(1, "v1", "old_value").AddRow(2, "v1", "old_value"), // 删除前查询到的结果，满足条件的行
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', 'v1', 'old_value');", // 删除前的备份数据
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('2', 'v1', 'old_value');",
			},
		},
		{
			caseName: "delete rows using subquery from exist_tb_1",
			querySql: "DELETE FROM exist_tb_1 WHERE id IN (SELECT id FROM exist_tb_2 WHERE v1 = 'v1');",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(1, "v1", "old_value"), // 删除前查询到的结果，满足条件的行
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', 'v1', 'old_value');", // 删除前的备份数据
			},
		},
		{
			caseName: "delete from exist_tb_1 with condition",
			querySql: "DELETE FROM exist_tb_1 WHERE v1 = 'v1' AND id = 1 ORDER BY v1 DESC;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(2, "v1", "old_value"), // 剩余的行
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('2', 'v1', 'old_value');", // 删除前的备份
			},
		},
		{
			caseName: "delete multiple rows in exist_tb_2",
			querySql: "DELETE FROM exist_tb_2 WHERE user_id = 1;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB2).AddRow(3, "v1", "value", 2), // 剩余的行
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_2` (`id`, `v1`, `v2`, `user_id`) VALUES ('3', 'v1', 'value', '2');", // 被删除行的备份
			},
		},
		{
			caseName: "delete no rows in exist_tb_1",
			querySql: "DELETE FROM exist_tb_1 WHERE id = 999;", // id 999 不存在，删除0行
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1), // 没有匹配的行返回空
				},
			},
			expectedBackupResults: []string{
				// 无删除，备份也不需要任何数据
			},
		},
		{
			caseName: "delete using subquery in exist_tb_1",
			querySql: "DELETE FROM exist_tb_1 WHERE id = (SELECT id FROM exist_db.exist_tb_2 WHERE v1 = 'v1' LIMIT 1);",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(2, "v1", "value"),
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('2', 'v1', 'value');", // 删除的行备份
			},
		},
		{
			caseName: "delete all rows from exist_tb_1",
			querySql: "DELETE FROM exist_tb_1;",
			expectedQueryResult: []*queryResult{
				{
					query:  selectAllRegex,
					result: sqlMock.NewRows(allColumnsOfExistTB1).AddRow(1, "v1", "old_value").AddRow(2, "v1", "old_value"), // 删除前查询到的结果
				},
			},
			expectedBackupResults: []string{
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('1', 'v1', 'old_value');", // 删除前的备份
				"REPLACE INTO `exist_db`.`exist_tb_1` (`id`, `v1`, `v2`) VALUES ('2', 'v1', 'old_value');",
				// 所有行都会备份
			},
		},
	}
	for _, testCase := range testCases {
		testBackupOneSql(t, testCase.caseName, testCase.querySql, testCase.expectedBackupResults, testCase.expectedQueryResult)
	}
}
