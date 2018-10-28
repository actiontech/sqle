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
