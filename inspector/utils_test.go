package inspector

import (
	"github.com/pingcap/tidb/ast"
	"github.com/stretchr/testify/assert"
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
