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

			assert.Equal(t, tt.isOnlyIncludeCount, visitor.IsSelectOnlyIncludeCountFunc)
		})
	}
}

func TestColumnNameVisitor(t *testing.T) {
	tests := []struct {
		input       string
		columnCount uint
	}{
		{"SELECT * FROM t1", 0},                                                         //不包含列
		{"SELECT a,b,c FROM t1 WHERE id > 1", 4},                                        //使用不等号
		{"SELECT COUNT(*) FROM t1", 0},                                                  //使用函数并不包含列
		{"SELECT a,COUNT(*) FROM t1 GROUP BY a", 2},                                     //使用函数包含列
		{"SELECT * FROM table1 INNER JOIN table2 ON table1.id = table2.table1_id", 2},   //使用JOIN
		{"SELECT * FROM table1 WHERE id IN ( SELECT id FROM table2 WHERE age > 30)", 3}, //使用子查询
		{"SELECT UPPER(name), LENGTH(comments) FROM table1", 2},                         //使用函数
		{"SELECT CAST(price AS DECIMAL(10,2))FROM products", 1},                         //使用类型转换
		{"SELECT * FROM table1 INNER JOIN table2 ON table1.id = table2.table1_id INNER JOIN table3 ON table2.id = table3.table2_id", 4}, //使用JOIN嵌套
		{"SELECT column1 AS alias1, column2 AS alias2 FROM table1", 2},                                                                  //使用列别名
		{"SELECT column1 + column2 AS sum_columns FROM table1", 2},
		{"SELECT t1.column1 AS t1_col1, t2.column2 AS t2_col2 FROM table1 t1 INNER JOIN table2 t2 ON t1.id = t2.t1_id", 4}, //不带AS的表别名
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tt.input, "", "")
			assert.NoError(t, err)

			visitor := &ColumnNameVisitor{}
			stmt.Accept(visitor)

			assert.Equal(t, tt.columnCount, uint(len(visitor.ColumnNameList)))
		})
	}
}

func TestEqualConditionVisitor(t *testing.T) {
	tests := []struct {
		input          string
		conditionCount int
	}{
		{"SELECT * FROM t1 WHERE t1.id1 = t2.id2", 1},
		{"SELECT * FROM t1 WHERE t1.id1 = t3.id3 OR t2.id2 = t1.id1", 2},
		{"SELECT * FROM t JOIN t2 ON t.id = t2.id2 WHERE t.name = ? AND t2.age = ?", 1},
		{"DELETE FROM t1 WHERE t1.id1 = t1.id2", 0},
		{"UPDATE t1 SET id2 = 2 WHERE id1 > 1", 0},
		{"INSERT INTO t1 (id1, id2) VALUES (1, 2)", 0},
		{"DELETE FROM t1 WHERE id2 > 2", 0},
		{"UPDATE t1 SET id1 = 2 WHERE t2.id2 = t3.id3", 1}, //SET id1 = 2不是BinaryOperation，而是Assignment
		{"INSERT INTO t1 (id1, id2) VALUES (2, 1)", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tt.input, "", "")
			assert.NoError(t, err)

			visitor := &EqualConditionVisitor{}
			stmt.Accept(visitor)

			assert.Equal(t, tt.conditionCount, len(visitor.ConditionList))
		})
	}
}

func TestFuncCallVisitor(t *testing.T) {
	tests := []struct {
		input          string
		conditionCount int
	}{
		{"SELECT * FROM t1 WHERE t1.id1 = t3.id3 OR t2.id2 = t1.id1", 0},
		{"SELECT UPPER(CONCAT(CONCAT('a_',UPPER('b'),'_c'),'_','a_',UPPER('b'),'_c'));", 5},
		{"SELECT UPPER('a');", 1},
		{"SELECT UPPER(CONCAT('a_',UPPER('b'),'_c'));", 3},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tt.input, "", "")
			assert.NoError(t, err)

			visitor := &FuncCallExprVisitor{}
			stmt.Accept(visitor)

			assert.Equal(t, tt.conditionCount, len(visitor.FuncCallList))
		})
	}
}
