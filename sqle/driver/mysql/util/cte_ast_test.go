package util

import (
	"testing"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/stretchr/testify/assert"
)

// TestCTEParseASTPositiveSupport covers AC-004: Parse succeeds and AST carries CTE
// definition fields (name / optional columns / AS subquery). Must not only assert
// "not Unparsed".
func TestCTEParseASTPositiveSupport(t *testing.T) {
	p := parser.New()

	t.Run("non_recursive_dql_with_optional_cols", func(t *testing.T) {
		sql := "WITH cte (id) AS (SELECT 1) SELECT * FROM cte"
		stmt, err := p.ParseOneStmt(sql, "", "")
		assert.NoError(t, err)
		assert.NotNil(t, stmt)
		_, isUnparsed := stmt.(*ast.UnparsedStmt)
		assert.False(t, isUnparsed, "expect parsed CTE SelectStmt, got UnparsedStmt")

		sel, ok := stmt.(*ast.SelectStmt)
		if !assert.True(t, ok, "expect *ast.SelectStmt, got %T", stmt) {
			return
		}
		if !assert.NotNil(t, sel.With) {
			return
		}
		assert.False(t, sel.With.IsRecursive)
		if !assert.Len(t, sel.With.CTEs, 1) {
			return
		}

		cte := sel.With.CTEs[0]
		assert.Equal(t, "cte", cte.Name.L)
		if assert.Len(t, cte.ColNameList, 1) {
			assert.Equal(t, "id", cte.ColNameList[0].L)
		}
		if assert.NotNil(t, cte.Query, "AS subquery must be present") {
			assert.NotNil(t, cte.Query.Query, "AS subquery Query node must be present")
		}
	})

	t.Run("with_recursive_dql", func(t *testing.T) {
		sql := `WITH RECURSIVE t AS (
  SELECT 1 AS n UNION ALL SELECT n + 1 FROM t WHERE n < 3
) SELECT * FROM t`
		stmt, err := p.ParseOneStmt(sql, "", "")
		assert.NoError(t, err)
		assert.NotNil(t, stmt)
		_, isUnparsed := stmt.(*ast.UnparsedStmt)
		assert.False(t, isUnparsed, "expect parsed CTE SelectStmt, got UnparsedStmt")

		sel, ok := stmt.(*ast.SelectStmt)
		if !assert.True(t, ok, "expect *ast.SelectStmt, got %T", stmt) {
			return
		}
		if !assert.NotNil(t, sel.With) {
			return
		}
		assert.True(t, sel.With.IsRecursive)
		if !assert.Len(t, sel.With.CTEs, 1) {
			return
		}

		cte := sel.With.CTEs[0]
		assert.Equal(t, "t", cte.Name.L)
		assert.Empty(t, cte.ColNameList)
		if assert.NotNil(t, cte.Query, "AS subquery must be present") {
			assert.NotNil(t, cte.Query.Query, "AS subquery Query node must be present")
		}
	})

	t.Run("cte_plus_dml_delete", func(t *testing.T) {
		sql := `WITH cte AS (SELECT id FROM exist_tb_1 WHERE id = 1)
DELETE FROM exist_tb_1 WHERE id IN (SELECT id FROM cte)`
		stmt, err := p.ParseOneStmt(sql, "", "")
		assert.NoError(t, err)
		assert.NotNil(t, stmt)
		_, isUnparsed := stmt.(*ast.UnparsedStmt)
		assert.False(t, isUnparsed, "expect parsed CTE DeleteStmt, got UnparsedStmt")

		del, ok := stmt.(*ast.DeleteStmt)
		if !assert.True(t, ok, "expect *ast.DeleteStmt, got %T", stmt) {
			return
		}
		if !assert.NotNil(t, del.With) {
			return
		}
		assert.False(t, del.With.IsRecursive)
		if !assert.Len(t, del.With.CTEs, 1) {
			return
		}

		cte := del.With.CTEs[0]
		assert.Equal(t, "cte", cte.Name.L)
		assert.Empty(t, cte.ColNameList)
		if assert.NotNil(t, cte.Query, "AS subquery must be present") {
			assert.NotNil(t, cte.Query.Query, "AS subquery Query node must be present")
		}
	})
}
