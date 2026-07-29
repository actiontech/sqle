package splitter

import (
	"testing"

	"github.com/pingcap/parser/ast"
	"github.com/stretchr/testify/assert"
)

func TestExtractOuterSQLAfterCTE(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantOuter string
		wantOK    bool
	}{
		{
			name: "minimal_dql_cte",
			sql: `WITH cte AS (
  SELECT 1 AS id
)
SELECT * FROM cte`,
			wantOuter: "SELECT * FROM cte",
			wantOK:    true,
		},
		{
			name: "with_recursive",
			sql: `WITH RECURSIVE t AS (
  SELECT 1 AS n
  UNION ALL
  SELECT n + 1 FROM t WHERE n < 3
)
SELECT * FROM t`,
			wantOuter: "SELECT * FROM t",
			wantOK:    true,
		},
		{
			name: "multi_cte",
			sql: `WITH cte1 AS (SELECT 1 AS id),
cte2 AS (SELECT 2 AS id)
SELECT * FROM cte1 JOIN cte2`,
			wantOuter: "SELECT * FROM cte1 JOIN cte2",
			wantOK:    true,
		},
		{
			name: "cte_with_column_list",
			sql: `WITH cte (id) AS (SELECT 1)
SELECT id FROM cte`,
			wantOuter: "SELECT id FROM cte",
			wantOK:    true,
		},
		{
			name:   "not_cte_with_grant_option",
			sql:    `GRANT SELECT ON t TO u WITH GRANT OPTION`,
			wantOK: false,
		},
		{
			name:   "plain_select",
			sql:    `SELECT 1`,
			wantOK: false,
		},
		{
			name:   "incomplete_cte_no_outer",
			sql:    `WITH cte AS (SELECT 1)`,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outer, ok := extractOuterSQLAfterCTE(tt.sql)
			assert.Equal(t, tt.wantOK, ok)
			if tt.wantOK {
				assert.Equal(t, tt.wantOuter, outer)
			}
		})
	}
}

// TestParseSqlText_MySQLCTE 覆盖 AC-007：合法 CTE 不因 CTE 语法本身落成 UnparsedStmt。
func TestParseSqlText_MySQLCTE(t *testing.T) {
	p := NewSplitter()

	tests := []struct {
		name     string
		sql      string
		wantType interface{} // concrete AST type expected
	}{
		{
			name:     "minimal_dql_cte",
			sql:      "WITH cte AS (SELECT 1 AS id) SELECT * FROM cte",
			wantType: &ast.SelectStmt{},
		},
		{
			name: "with_recursive",
			sql: `WITH RECURSIVE t AS (
  SELECT 1 AS n
  UNION ALL
  SELECT n + 1 FROM t WHERE n < 3
)
SELECT * FROM t`,
			wantType: &ast.SelectStmt{},
		},
		{
			name: "multi_cte",
			sql: `WITH cte1 AS (SELECT 1 AS id),
cte2 AS (SELECT id FROM cte1)
SELECT * FROM cte2`,
			wantType: &ast.SelectStmt{},
		},
		{
			name:     "cte_insert",
			sql:      "WITH cte AS (SELECT 1 AS id) INSERT INTO t (id) SELECT id FROM cte",
			wantType: &ast.InsertStmt{},
		},
		{
			name:     "cte_update",
			sql:      "WITH cte AS (SELECT 1 AS id) UPDATE t SET v = 1 WHERE id IN (SELECT id FROM cte)",
			wantType: &ast.UpdateStmt{},
		},
		{
			name:     "cte_delete",
			sql:      "WITH cte AS (SELECT id FROM exist_tb_1 WHERE id = 1) DELETE FROM exist_tb_1 WHERE id IN (SELECT id FROM cte)",
			wantType: &ast.DeleteStmt{},
		},
		{
			name:     "real_ddl_create_table",
			sql:      "CREATE TABLE t_cte_ctrl (id INT PRIMARY KEY)",
			wantType: &ast.CreateTableStmt{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := p.ParseSqlText(tt.sql)
			assert.NoError(t, err)
			if !assert.Len(t, stmts, 1) {
				return
			}
			_, isUnparsed := stmts[0].(*ast.UnparsedStmt)
			assert.False(t, isUnparsed, "合法 CTE/对照 SQL 不应落成 UnparsedStmt, got %T", stmts[0])
			assert.IsType(t, tt.wantType, stmts[0])
			assert.Equal(t, tt.sql, stmts[0].Text())
		})
	}

	t.Run("illegal_sql_still_unparsed", func(t *testing.T) {
		stmts, err := p.ParseSqlText("SELECTxxx FROM nowhere WHERE")
		assert.NoError(t, err)
		if !assert.Len(t, stmts, 1) {
			return
		}
		_, ok := stmts[0].(*ast.UnparsedStmt)
		assert.True(t, ok, "明显非法 SQL 应仍为 UnparsedStmt, got %T", stmts[0])
	})
}
