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

func TestSelectCountSqlExtractor(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"SELECT * FROM t1", "SELECT count(1) FROM t1"},
		{"SELECT * FROM (SELECT * FROM t1) as t2", "SELECT count(1) FROM (SELECT * FROM (t1)) AS t2"},
		{"SELECT * FROM t1 WHERE id = (SELECT id FROM t2 WHERE a = 1)", "SELECT count(1) FROM t1 WHERE id=(SELECT id FROM t2 WHERE a=1)"},
		{"select t2.id from t2 where id = 1 order by id limit 1", "SELECT count(1) FROM t2 WHERE id=1 ORDER BY id LIMIT 1"},
		{"select t1.id,t2.id from t2 join t1 on t1.id = t2.id where id = 1 order by id limit 1, 1", "SELECT count(1) FROM t2 JOIN t1 ON t1.id=t2.id WHERE id=1 ORDER BY id LIMIT 1,1"},
		{"delete from t1 where id = 1", "SELECT count(1) FROM t1 WHERE id=1"},
		{"DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id;", "SELECT count(1) FROM (t1 JOIN t2) JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id"},
		{"DELETE FROM somelog WHERE user = jcole ORDER BY timestamp_column LIMIT 1;", "SELECT count(1) FROM somelog WHERE user=jcole ORDER BY timestamp_column LIMIT 1"},
		{"DELETE t1 FROM t1 LEFT JOIN t2 ON t1.id=t2.id WHERE t2.id IS NULL;", "SELECT count(1) FROM t1 LEFT JOIN t2 ON t1.id=t2.id WHERE t2.id IS NULL"},
		{"DELETE FROM a1, a2 USING t1 AS a1 INNER JOIN t2 AS a2 WHERE a1.id=a2.id;", "SELECT count(1) FROM t1 AS a1 JOIN t2 AS a2 WHERE a1.id=a2.id"},
		{"UPDATE t1 SET col1 = col1 + 1;", "SELECT count(1) FROM t1"},
		{"UPDATE t SET id = id + 1 ORDER BY id DESC limit 10;", "SELECT count(1) FROM t ORDER BY id DESC LIMIT 10"},
		{"UPDATE items,month SET items.price=month.price WHERE items.id=month.id;", "SELECT count(1) FROM (items) JOIN month WHERE items.id=month.id"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tt.input, "", "")
			assert.NoError(t, err)

			visitor := &SelectCountNodeExtractor{}
			node, _ := stmt.Accept(visitor)

			var buf strings.Builder
			if err := node.Restore(format.NewRestoreCtx(0, &buf)); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expect, buf.String())
		})
	}
}

func TestSelectFieldExtractor(t *testing.T) {
	tests := []struct {
		input              string
		isOnlyIncludeCount bool
	}{
		{"SELECT * FROM t1", false},
		{"SELECT COUNT(1) FROM `test`.`test`", true},
		{"SELECT count(*) FROM (SELECT * FROM t1) as t2", true},
		{"SELECT count(*),count(id) FROM (SELECT * FROM t1) as t2", false},
		{"SELECT count(1) FROM (SELECT * FROM t1) as t2", true},
		{"SELECT count(1),id FROM (SELECT * FROM t1) as t2", false},
		{"(SELECT count(1),id FROM (SELECT * FROM t1) as t2)", false},
		{"(SELECT count(1) FROM (SELECT * FROM t1) as t2)", true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tt.input, "", "")
			assert.NoError(t, err)

			visitor := &SelectFieldExtractor{}
			stmt.Accept(visitor)

			assert.Equal(t, tt.isOnlyIncludeCount, visitor.IsOnlyIncludeCountFunc)
		})
	}
}
