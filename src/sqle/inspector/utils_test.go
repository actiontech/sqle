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

func TestHasSpecialOption(t *testing.T) {
	options := []*ast.ColumnOption{
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
	assert.Equal(t, HasSpecialOption(options, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, HasSpecialOption(options, ast.ColumnOptionAutoIncrement), true)
	assert.Equal(t, HasSpecialOption(options, ast.ColumnOptionNotNull), true)
	assert.Equal(t, HasSpecialOption(options, ast.ColumnOptionNull), false)
	assert.Equal(t, HasSpecialOption(options, ast.ColumnOptionAutoIncrement, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, HasSpecialOption(options, ast.ColumnOptionNull, ast.ColumnOptionPrimaryKey), false)
}

func TestInspectResults(t *testing.T) {
	results := newInspectResults()
	results.add(model.RULE_LEVEL_ERROR, model.DDL_CREATE_TABLE_NOT_EXIST)
	assert.Equal(t, "error", results.level())
	assert.Equal(t, "[error]新建表必须加入if not exists create，保证重复执行不报错", results.message())

	results.add(model.RULE_LEVEL_NOTICE, model.TABLE_NOT_EXIST, "not_exist_tb")
	assert.Equal(t, "error", results.level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[notice]表 not_exist_tb 不存在`, results.message())
}
