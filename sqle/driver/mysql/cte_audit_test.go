package mysql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// unsupportedSQLWarnSnippet is the Chinese body of UnsupportedSyntaxError /
// server unsupportedSQLWarnMessage. AC-005 forbids this when triggered by legal CTE.
const unsupportedSQLWarnSnippet = "语法错误或者解析器不支持，请人工确认SQL正确性"

// TestCTEAuditNoUnsupportedSQLWarn covers AC-005: legal CTE audit must not emit
// unsupportedSQLWarnMessage (UnparsedStmt → checkUnparsedStmt path).
// Empty rule set isolates the parser-unsupported warning from unrelated rule hits.
func TestCTEAuditNoUnsupportedSQLWarn(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		{
			name: "non_recursive_dql_with_optional_cols",
			sql:  `WITH cte (id) AS (SELECT 1) SELECT * FROM cte;`,
		},
		{
			name: "with_recursive_dql",
			sql:  `WITH RECURSIVE t AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM t WHERE n < 3) SELECT * FROM t;`,
		},
		{
			name: "cte_plus_dml_delete",
			sql:  `WITH cte AS (SELECT id FROM exist_db.exist_tb_1 WHERE id = 1) DELETE FROM exist_db.exist_tb_1 WHERE id IN (SELECT id FROM cte);`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			runEmptyRuleInspectCase(t, c.name, DefaultMysqlInspectOffline(), c.sql, newTestResult())

			i := DefaultMysqlInspectOffline()
			i.rules = nil
			actual, err := i.Audit(nil, []string{c.sql})
			assert.NoError(t, err)
			if !assert.Len(t, actual, 1) {
				return
			}
			msg := actual[0].Message()
			assert.False(t, strings.Contains(msg, unsupportedSQLWarnSnippet),
				"legal CTE must not trigger unsupportedSQLWarnMessage; got %q", msg)
		})
	}
}
