package mysql

import (
	"testing"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
	"github.com/stretchr/testify/assert"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

func TestCTE_ParseAndSQLType(t *testing.T) {
	i := &MysqlDriverImpl{}

	t.Run("select", func(t *testing.T) {
		sql := `WITH cte AS (SELECT 1 AS id) SELECT * FROM cte;`
		stmts, err := util.ParseSql(sql)
		assert.NoError(t, err)
		assert.IsType(t, &ast.SelectStmt{}, stmts[0])
		s := stmts[0].(*ast.SelectStmt)
		assert.NotNil(t, s.With)
		assert.Equal(t, "cte", s.With.CTEs[0].Name.O)
		assert.NotNil(t, s.With.CTEs[0].Query)
		assert.False(t, s.With.IsRecursive)
		assert.Equal(t, driverV2.SQLTypeDQL, i.assertSQLType(stmts[0]))
	})

	t.Run("recursive", func(t *testing.T) {
		sql := `WITH RECURSIVE cte AS (SELECT 1 AS n UNION ALL SELECT n+1 FROM cte WHERE n < 3) SELECT * FROM cte;`
		stmts, err := util.ParseSql(sql)
		assert.NoError(t, err)
		s, ok := stmts[0].(*ast.SelectStmt)
		assert.True(t, ok)
		assert.NotNil(t, s.With)
		assert.True(t, s.With.IsRecursive)
		assert.Equal(t, "cte", s.With.CTEs[0].Name.O)
		assert.Equal(t, driverV2.SQLTypeDQL, i.assertSQLType(stmts[0]))
	})

	t.Run("delete", func(t *testing.T) {
		sql := `WITH cte AS (SELECT 1 AS id) DELETE FROM t WHERE id IN (SELECT id FROM cte);`
		stmts, err := util.ParseSql(sql)
		assert.NoError(t, err)
		s, ok := stmts[0].(*ast.DeleteStmt)
		assert.True(t, ok)
		assert.NotNil(t, s.With)
		assert.Equal(t, "cte", s.With.CTEs[0].Name.O)
		assert.Equal(t, driverV2.SQLTypeDML, i.assertSQLType(stmts[0]))
	})

	t.Run("ddl still ddl", func(t *testing.T) {
		sql := `CREATE TABLE new_tbl AS SELECT * FROM orig_tbl;`
		stmts, err := util.ParseSql(sql)
		assert.NoError(t, err)
		assert.Equal(t, driverV2.SQLTypeDDL, i.assertSQLType(stmts[0]))
	})
}
