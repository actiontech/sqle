package mysql

import (
	"context"
	"testing"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
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
	assert.Equal(t, nodes[0].Type, driverV2.SQLTypeDQL)

	nodes, err = DefaultMysqlInspect().Parse(context.TODO(), "insert into tb1 values(1)")
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].Type, driverV2.SQLTypeDML)

	nodes, err = DefaultMysqlInspect().Parse(context.TODO(), `
INSERT INTO customers (customer_name, email)
SELECT first_name, email
FROM contacts
WHERE last_name = 'Smith';`)
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
func TestInspect_assertSQLType(t *testing.T) {
	args := []struct {
		Name string
		SQL  string
		Want string
	}{
		{
			"case 1",
			`select * from tb`,
			driverV2.SQLTypeDQL,
		},
		{
			"case 2",
			`
(SELECT a FROM t1 WHERE a=10 AND B=1 ORDER BY a LIMIT 10)
UNION
(SELECT a FROM t2 WHERE a=11 AND B=2 ORDER BY a LIMIT 10);`,
			driverV2.SQLTypeDQL,
		},
		{
			"case 3",
			`
DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3
WHERE t1.id=t2.id AND t2.id=t3.id;`,
			driverV2.SQLTypeDML,
		},
		{
			"case 4",
			`
INSERT INTO tbl_name (a,b,c) VALUES(1,2,3,4,5,6,7,8,9);`,
			driverV2.SQLTypeDML,
		},
		{
			"case 5",
			`
UPDATE t SET id = id + 1;`,
			driverV2.SQLTypeDML,
		},
		{
			"case 6",
			`
CREATE TABLE new_tbl AS SELECT * FROM orig_tbl;`,
			driverV2.SQLTypeDDL,
		},
		{
			"case 7",  // unparsed
			`
CREATEaa TABLE new_tbl AS SELECT * FROM orig_tbl;`,
			driverV2.SQLTypeDDL,
		},
	}
	i := &MysqlDriverImpl{}
	for _, arg := range args {
		t.Run(arg.Name, func(t *testing.T) {
			stmt, err := util.ParseSql(arg.SQL)
			assert.NoError(t, err)
			assert.Equal(t, arg.Want, i.assertSQLType(stmt[0]))
		})
	}
}
