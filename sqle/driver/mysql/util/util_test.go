package util

import (
	"strings"
	"testing"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/stretchr/testify/assert"
)

func TestGetSelectNodeFromSelect(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"SELECT * FROM t1", "SELECT COUNT(1) FROM `t1`"},
		{"SELECT * FROM (SELECT * FROM t1) as t2", "SELECT COUNT(1) FROM (SELECT * FROM (`t1`)) AS `t2`"},
		{"SELECT * FROM t1 WHERE id = (SELECT id FROM t2 WHERE a = 1)", "SELECT COUNT(1) FROM `t1` WHERE `id`=(SELECT `id` FROM `t2` WHERE `a`=1)"},
		{"select t2.id from t2 where id = 1 order by id limit 1", "SELECT COUNT(1) FROM `t2` WHERE `id`=1 ORDER BY `id` LIMIT 1"},
		{"select t1.id,t2.id from t2 join t1 on t1.id = t2.id where id = 1 order by id limit 1, 1", "SELECT COUNT(1) FROM `t2` JOIN `t1` ON `t1`.`id`=`t2`.`id` WHERE `id`=1 ORDER BY `id` LIMIT 1,1"},
		{"delete from t1 where id = 1", "SELECT COUNT(1) FROM `t1` WHERE `id`=1"},
		{"DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id;", "SELECT COUNT(1) FROM (`t1` JOIN `t2`) JOIN `t3` WHERE `t1`.`id`=`t2`.`id` AND `t2`.`id`=`t3`.`id`"},
		{"DELETE FROM someLog WHERE user = jack ORDER BY timestamp_column LIMIT 1;", "SELECT COUNT(1) FROM `someLog` WHERE `user`=`jack` ORDER BY `timestamp_column` LIMIT 1"},
		{"DELETE t1 FROM t1 LEFT JOIN t2 ON t1.id=t2.id WHERE t2.id IS NULL;", "SELECT COUNT(1) FROM `t1` LEFT JOIN `t2` ON `t1`.`id`=`t2`.`id` WHERE `t2`.`id` IS NULL"},
		{"DELETE FROM a1, a2 USING t1 AS a1 INNER JOIN t2 AS a2 WHERE a1.id=a2.id;", "SELECT COUNT(1) FROM `t1` AS `a1` JOIN `t2` AS `a2` WHERE `a1`.`id`=`a2`.`id`"},
		{"UPDATE t1 SET col1 = col1 + 1;", "SELECT COUNT(1) FROM `t1`"},
		{"UPDATE t SET id = id + 1 ORDER BY id DESC limit 10;", "SELECT COUNT(1) FROM `t` ORDER BY `id` DESC LIMIT 10"},
		{"UPDATE items,month SET items.price=month.price WHERE items.id=month.id;", "SELECT COUNT(1) FROM (`items`) JOIN `month` WHERE `items`.`id`=`month`.`id`"},
	}

	for _, test := range tests {
		node, err := ParseOneSql(test.input)
		assert.NoError(t, err)

		var newNode ast.Node
		switch stmt := node.(type) {
		case *ast.SelectStmt:
			newNode = getSelectNodeFromSelect(stmt)
		case *ast.DeleteStmt:
			newNode = getSelectNodeFromDelete(stmt)
		case *ast.UpdateStmt:
			newNode = getSelectNodeFromUpdate(stmt)
		}

		sqlBuilder := new(strings.Builder)
		err = newNode.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, sqlBuilder))
		assert.NoError(t, err)

		assert.Equal(t, test.expect, sqlBuilder.String())
	}
}

func TestSaveFingerprint(t *testing.T) {
	type FpCase struct {
		input  string
		expect string
	}

	cases := []FpCase{
		{
			input:  `update  tb1 set a = "2" where a = "3" and b = 4`,
			expect: "update tb1 set a = ? where a = ? and b = ?",
		},
		// not sql
		{
			input:  "SELECT*FROM (SELECT * FROM tb values(1));",
			expect: "SELECT*FROM (SELECT * FROM tb values(1));",
		},
		{
			input:  "not sql",
			expect: "not sql",
		},
		// https://github.com/actiontech/sqle/issues/2603
		{
			input:  "insert into tb values(1)",
			expect: "insert into tb values(1)",
		},
	}

	for i := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			fp := SafeFingerprint(log.NewEntry(), cases[i].input)
			assert.Equal(t, cases[i].expect, fp)
		})
	}
}
