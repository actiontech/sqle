package inspector

import (
	"github.com/pingcap/tidb/ast"
	"github.com/stretchr/testify/assert"
	"sqle/model"
	"testing"
)

func TestRemoveArrayRepeat(t *testing.T) {
	input := []string{"a", "b", "c", "c", "a"}
	expect := []string{"a", "b", "c"}
	actual := RemoveArrayRepeat(input)
	assert.Equal(t, expect, actual)
}

var columnOptionsForTest = []*ast.ColumnOption{
	&ast.ColumnOption{
		Tp: ast.ColumnOptionAutoIncrement,
	},
	&ast.ColumnOption{
		Tp: ast.ColumnOptionPrimaryKey,
	},
	&ast.ColumnOption{
		Tp: ast.ColumnOptionNotNull,
	},
}

func TestIsAllInOptions(t *testing.T) {
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionAutoIncrement), true)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionNotNull), true)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionNull), false)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionAutoIncrement, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionNull, ast.ColumnOptionPrimaryKey), false)
}

func TestHasOneInOptions(t *testing.T) {
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionAutoIncrement), true)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionNotNull), true)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionNull), false)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionAutoIncrement, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionNull, ast.ColumnOptionPrimaryKey), true)
}

func TestReplaceSchemaName(t *testing.T) {
	input := "alter table `db1`.tb1 drop column a1"
	output := "alter table `tb1` drop column a1"
	assert.Equal(t, replaceTableName(input, "db1", "tb1"), output)
}

func TestGetDuplicate(t *testing.T) {
	assert.Equal(t, []string{}, getDuplicate([]string{"1", "2", "3"}))
	assert.Equal(t, []string{"2"}, getDuplicate([]string{"1", "2", "2"}))
	assert.Equal(t, []string{"2", "3"}, getDuplicate([]string{"1", "2", "2", "3", "3", "3"}))
}

func TestRemoveDuplicate(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3"}, removeDuplicate([]string{"1", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3"}, removeDuplicate([]string{"1", "2", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3"}, removeDuplicate([]string{"1", "2", "2", "3", "3", "3"}))
}

func TestInspectResults(t *testing.T) {
	results := newInspectResults()
	handler := RuleHandlerMap[DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST]
	results.add(handler.Rule.Level, handler.Message)
	assert.Equal(t, "error", results.level())
	assert.Equal(t, "[error]新建表必须加入if not exists create，保证重复执行不报错", results.message())

	results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "not_exist_tb")
	assert.Equal(t, "error", results.level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在`, results.message())

	results2 := newInspectResults()
	results2.add(results.level(), results.message())
	results2.add("notice", "test")
	assert.Equal(t, "error", results2.level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test`, results2.message())

	results3 := newInspectResults()
	results3.add(results2.level(), results2.message())
	results3.add("notice", "[osc]test")
	assert.Equal(t, "error", results3.level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test
[osc]test`, results3.message())

	results4 := newInspectResults()
	results4.add("notice", "[notice]test")
	results4.add("error", "[osc]test")
	assert.Equal(t, "error", results4.level())
	assert.Equal(t,
		`[notice]test
[osc]test`, results4.message())

	results5 := newInspectResults()
	results5.add("warn", "[warn]test")
	results5.add("notice", "[osc]test")
	assert.Equal(t, "warn", results5.level())
	assert.Equal(t,
		`[warn]test
[osc]test`, results5.message())
}
