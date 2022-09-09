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
		{"SELECT * FROM (SELECT * FROM t1) as t2", []string{"SELECT * FROM (SELECT * FROM (t1)) AS t2", "SELECT * FROM t1"}},
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

func TestSubQueryMaxNestNumExtractor(t *testing.T) {
	tests := []struct {
		input  string
		expect int
	}{
		{`select (select count(*) from users) as a
from users
where (select count(*) from users) > 1
  and (select (select count(*) from users limit 1) from users where (select id from users where (select id from users limit 1) = 1 limit 1) = 1) > 1`, 3},
		{`select (select count(*) from users) as a
from exist_db.exist_tb_1
where (select count(*) from exist_db.exist_tb_2) > 1
  and (select count(*)
       from exist_db.exist_tb_1
       where (select id
              from exist_db.exist_tb_1
              where (select count(*) from exist_db.exist_tb_2 where (select count(*) from exist_db.exist_tb_2) = 1) =
                    1) = 1) = 1`, 4},
		{`update exist_db.exist_tb_1,exist_db.exist_tb_2
set exist_tb_1.v1 = exist_tb_2.v1
where (select count(*) from exist_db.exist_tb_2) > 1
  and (select count(*)
       from exist_db.exist_tb_1
       where exist_tb_1.id = 1
         and (select id from exist_db.exist_tb_1 limit 1) = 1) > 1`, 2},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tt.input, "", "")
			assert.NoError(t, err)

			var maxNestNum int
			visitor := &SubQueryMaxNestNumExtractor{MaxNestNum: &maxNestNum, CurrentNestNum: 1}
			stmt.Accept(visitor)

			assert.Equal(t, tt.expect, *visitor.MaxNestNum)
		})
	}
}
