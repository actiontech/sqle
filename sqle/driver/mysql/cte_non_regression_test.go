package mysql

import (
	"context"
	"strings"
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/stretchr/testify/assert"
)

// TestCTENonRegressionIllegalAndNonCTEWith covers AC-009:
// illegal SQL still fails / warns reasonably; GRANT … WITH GRANT OPTION must not be
// treated as CTE success path and must keep non-"parser unsupported" rule hits.
func TestCTENonRegressionIllegalAndNonCTEWith(t *testing.T) {
	t.Run("illegal_sql_not_silent_legal_cte", func(t *testing.T) {
		sql := `SELECT FROM WHERE;`

		stmts, err := util.ParseSql(sql)
		assert.NoError(t, err)
		if !assert.Len(t, stmts, 1) {
			return
		}
		_, isUnparsed := stmts[0].(*ast.UnparsedStmt)
		assert.True(t, isUnparsed, "illegal SQL must be UnparsedStmt, got %T", stmts[0])
		if sel, ok := stmts[0].(*ast.SelectStmt); ok {
			assert.Nil(t, sel.With, "illegal SQL must not become SelectStmt with CTE With")
		}

		i := DefaultMysqlInspectOffline()
		i.rules = nil
		actual, err := i.Audit(context.TODO(), []string{sql})
		assert.NoError(t, err)
		if !assert.Len(t, actual, 1) {
			return
		}
		msg := actual[0].Message()
		assert.True(t, strings.Contains(msg, unsupportedSQLWarnSnippet),
			"illegal SQL must keep reasonable unsupported/parse warn; got %q", msg)

		assert.NotEqual(t, driverV2.SQLTypeDQL, (&MysqlDriverImpl{}).assertSQLType(stmts[0]),
			"illegal UnparsedStmt must not be classified as CTE dql success path")
	})

	t.Run("grant_with_grant_option_not_cte", func(t *testing.T) {
		sql := `GRANT SELECT ON database.* TO 'user'@'localhost' WITH GRANT OPTION;`

		stmts, err := util.ParseSql(sql)
		assert.NoError(t, err)
		if !assert.Len(t, stmts, 1) {
			return
		}
		grant, ok := stmts[0].(*ast.GrantStmt)
		if !assert.True(t, ok, "expect *ast.GrantStmt, got %T", stmts[0]) {
			return
		}
		assert.True(t, grant.WithGrant, "WITH GRANT OPTION must set GrantStmt.WithGrant")
		_, isSelect := stmts[0].(*ast.SelectStmt)
		assert.False(t, isSelect, "GRANT must not be misparsed as CTE SelectStmt")

		sqlType := (&MysqlDriverImpl{}).assertSQLType(stmts[0])
		assert.NotEqual(t, driverV2.SQLTypeDQL, sqlType,
			"GRANT WITH GRANT OPTION must not follow CTE dql success path; got %s", sqlType)
		assert.Equal(t, driverV2.SQLTypeDDL, sqlType)

		// Non-"parser unsupported" reasonable alert must remain (SQLE00174).
		rule := rulepkg.AIRuleHandlerMap[ai.SQLE00174].Rule
		runAIRuleCase(rule, t, "grant_with_option_keeps_sqle00174", sql,
			nil, nil, newTestResult().addResult(ai.SQLE00174))
	})
}
