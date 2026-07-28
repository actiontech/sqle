package splitter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTrailingCommentOnlySkipped(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		wantLen  int
		wantText string
	}{
		{"semicolon then line comment", "select 1;\n-- ces", 1, "select 1;"},
		{"semicolon then hash comment", "select 1;\n# ces", 1, "select 1;"},
		{"semicolon then block comment", "select 1;\n/* ces */", 1, "select 1;"},
		{"inline after semicolon", "select 1; -- ces", 1, "select 1;"},
		{"comment before sql kept in text", "-- ces\nselect 1;", 1, "-- ces\nselect 1;"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes, err := NewSplitter().ParseSqlText(tc.sql)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantLen, len(nodes))
			assert.Equal(t, tc.wantText, nodes[0].Text())
		})
	}
}
