package mysql

import (
	"context"
	"strings"
	"testing"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/format"
	"github.com/stretchr/testify/assert"
)

func TestInspect_Parse(t *testing.T) {
	nodes, err := DefaultMysqlInspect().Parse(context.TODO(), `
use test_db;
create trigger my_trigger before insert on t1 for each row insert into t2(id, c1) values(1, '2');
create table t1(id int);
	`)
	assert.NoError(t, err)
	for _, node := range nodes {
		assert.Equal(t, node.Type, driverV2.SQLTypeDDL)
	}

	nodes, err = DefaultMysqlInspect().Parse(context.TODO(), "select * from t1")
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].Type, driverV2.SQLTypeDML)
}

func TestInspect_onlineddlWithGhost(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		setUp   func(*MysqlDriverImpl) *MysqlDriverImpl
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "alter stmt(true); config onlineddl(true); table size enough(true)",
			setUp: func(i *MysqlDriverImpl) *MysqlDriverImpl {
				i.Ctx.Schemas()["exist_db"].Tables["exist_tb_1"].Size = 17
				return i
			},
			args:    args{query: "alter table exist_db.exist_tb_1 add column col1 varchar(100);"},
			want:    true,
			wantErr: false,
		},

		{
			name: "alter stmt(true); config onlineddl(true); table size enough(false)",
			setUp: func(i *MysqlDriverImpl) *MysqlDriverImpl {
				i.Ctx.Schemas()["exist_db"].Tables["exist_tb_1"].Size = 15
				return i
			},
			args:    args{query: "alter table exist_db.exist_tb_1 add column col1 varchar(100);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(true); config onlineddl(false); table size enough(true)",
			setUp: func(i *MysqlDriverImpl) *MysqlDriverImpl {
				i.cnf.DDLGhostMinSize = -1
				i.Ctx.Schemas()["exist_db"].Tables["exist_tb_1"].Size = 17
				return i
			},
			args:    args{query: "alter table exist_db.exist_tb_1 add column col1 varchar(100);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(false); config onlineddl(true); table size enough(true)",
			setUp: func(i *MysqlDriverImpl) *MysqlDriverImpl {
				i.Ctx.Schemas()["exist_db"].Tables["exist_tb_1"].Size = 17
				return i
			},
			args:    args{query: "create index idx_exist_db_exist_tb_1 on exist_db(v2);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(false); config onlineddl(false); table size enough(true)",
			setUp: func(i *MysqlDriverImpl) *MysqlDriverImpl {
				i.cnf.DDLGhostMinSize = -1
				i.Ctx.Schemas()["exist_db"].Tables["exist_tb_1"].Size = 17
				return i
			},
			args:    args{query: "create index idx_exist_db_exist_tb_1 on exist_db(v2);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(false); config onlineddl(true); table size enough(false)",
			setUp: func(i *MysqlDriverImpl) *MysqlDriverImpl {
				i.Ctx.Schemas()["exist_db"].Tables["exist_tb_1"].Size = 15
				return i
			},
			args:    args{query: "create index idx_exist_db_exist_tb_1 on exist_db(v2);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(true); config onlineddl(false); table size enough(false)",
			setUp: func(i *MysqlDriverImpl) *MysqlDriverImpl {
				i.cnf.DDLGhostMinSize = -1
				i.Ctx.Schemas()["exist_db"].Tables["exist_tb_1"].Size = 15
				return i
			},
			args:    args{query: "alter table exist_db.exist_tb_1 add column col1 varchar(100);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(false); config onlineddl(false); table size enough(false)",
			setUp: func(i *MysqlDriverImpl) *MysqlDriverImpl {
				i.cnf.DDLGhostMinSize = -1
				i.Ctx.Schemas()["exist_db"].Tables["exist_tb_1"].Size = 15
				return i
			},
			args:    args{query: "create index idx_exist_db_exist_tb_1 on exist_db(v2);"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := DefaultMysqlInspect()
			i.cnf.DDLGhostMinSize = 16
			got, err := tt.setUp(i).onlineddlWithGhost(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("MysqlDriverImpl.onlineddlWithGhost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MysqlDriverImpl.onlineddlWithGhost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInspect_GenRollbackSQL(t *testing.T) {
	i := DefaultMysqlInspect()

	rollback, reason, err := i.GenRollbackSQL(context.TODO(), "create table t1(id int, c1 int)")
	assert.NoError(t, err)
	assert.Equal(t, "", reason)
	assert.Equal(t, "DROP TABLE IF EXISTS `exist_db`.`t1`", rollback)

	rollback, reason, err = i.GenRollbackSQL(context.TODO(), "alter table t1 drop column c1")
	assert.NoError(t, err)
	assert.Equal(t, "", reason)
	assert.Equal(t, "ALTER TABLE `exist_db`.`t1`\nADD COLUMN `c1` int(11);", rollback)

	rollback, reason, err = i.GenRollbackSQL(context.TODO(), "alter table t1 add column c1 int")
	assert.NoError(t, err)
	assert.Equal(t, "", reason)
	assert.Equal(t, "ALTER TABLE `exist_db`.`t1`\nDROP COLUMN `c1`;", rollback)
}

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
		{"DELETE FROM somelog WHERE user = jcole ORDER BY timestamp_column LIMIT 1;", "SELECT COUNT(1) FROM `somelog` WHERE `user`=`jcole` ORDER BY `timestamp_column` LIMIT 1"},
		{"DELETE t1 FROM t1 LEFT JOIN t2 ON t1.id=t2.id WHERE t2.id IS NULL;", "SELECT COUNT(1) FROM `t1` LEFT JOIN `t2` ON `t1`.`id`=`t2`.`id` WHERE `t2`.`id` IS NULL"},
		{"DELETE FROM a1, a2 USING t1 AS a1 INNER JOIN t2 AS a2 WHERE a1.id=a2.id;", "SELECT COUNT(1) FROM `t1` AS `a1` JOIN `t2` AS `a2` WHERE `a1`.`id`=`a2`.`id`"},
		{"UPDATE t1 SET col1 = col1 + 1;", "SELECT COUNT(1) FROM `t1`"},
		{"UPDATE t SET id = id + 1 ORDER BY id DESC limit 10;", "SELECT COUNT(1) FROM `t` ORDER BY `id` DESC LIMIT 10"},
		{"UPDATE items,month SET items.price=month.price WHERE items.id=month.id;", "SELECT COUNT(1) FROM (`items`) JOIN `month` WHERE `items`.`id`=`month`.`id`"},
	}

	for _, test := range tests {
		node, err := util.ParseOneSql(test.input)
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
