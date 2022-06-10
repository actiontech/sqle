package util

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/format"

	"github.com/stretchr/testify/assert"
)

func TestSelectStmtExtractor(t *testing.T) {
	tests := []struct {
		input  string
		output []string
	}{
		{"SELECT * FROM t1", []string{"SELECT * FROM t1"}},
		{"SELECT * FROM (SELECT * FROM t1) as t2", []string{"SELECT * FROM (SELECT * FROM t1) AS t2", "SELECT * FROM t1"}},
		{"SELECT * FROM t1 WHERE id = (SELECT id FROM t2 WHERE a = 1)", []string{"SELECT * FROM t1 WHERE id=(SELECT id FROM t2 WHERE a=1)", "SELECT id FROM t2 WHERE a=1"}},
		{"SELECT * FROM t1 UNION SELECT * FROM t2", []string{"SELECT * FROM t1", "SELECT * FROM t2"}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tt.input, "", "")
			assert.NoError(t, err)

			visitor := &SelectStmtExtractor{}
			stmt.Accept(visitor)

			for i, ss := range visitor.SelectStmts {
				var buf strings.Builder
				assert.NoError(t, ss.Restore(format.NewRestoreCtx(0, &buf)))
				assert.Equal(t, tt.output[i], buf.String())
			}
		})
	}
}
