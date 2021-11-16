package mysql

import (
	"context"
	"testing"

	"github.com/actiontech/sqle/sqle/driver"
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
		assert.Equal(t, node.Type, driver.SQLTypeDDL)
	}

	nodes, err = DefaultMysqlInspect().Parse(context.TODO(), "select * from t1")
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].Type, driver.SQLTypeDML)
}

func TestInspect_onlineddlWithGhost(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		setUp   func(*Inspect) *Inspect
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "alter stmt(true); config onlineddl(true); table size enough(true)",
			setUp: func(i *Inspect) *Inspect {
				i.Ctx.schemas["exist_db"].Tables["exist_tb_1"].Size = 17
				return i
			},
			args:    args{query: "alter table exist_db.exist_tb_1 add column col1 varchar(100);"},
			want:    true,
			wantErr: false,
		},

		{
			name: "alter stmt(true); config onlineddl(true); table size enough(false)",
			setUp: func(i *Inspect) *Inspect {
				i.Ctx.schemas["exist_db"].Tables["exist_tb_1"].Size = 15
				return i
			},
			args:    args{query: "alter table exist_db.exist_tb_1 add column col1 varchar(100);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(true); config onlineddl(false); table size enough(true)",
			setUp: func(i *Inspect) *Inspect {
				i.cnf.DDLGhostMinSize = -1
				i.Ctx.schemas["exist_db"].Tables["exist_tb_1"].Size = 17
				return i
			},
			args:    args{query: "alter table exist_db.exist_tb_1 add column col1 varchar(100);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(false); config onlineddl(true); table size enough(true)",
			setUp: func(i *Inspect) *Inspect {
				i.Ctx.schemas["exist_db"].Tables["exist_tb_1"].Size = 17
				return i
			},
			args:    args{query: "create index idx_exist_db_exist_tb_1 on exist_db(v2);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(false); config onlineddl(false); table size enough(true)",
			setUp: func(i *Inspect) *Inspect {
				i.cnf.DDLGhostMinSize = -1
				i.Ctx.schemas["exist_db"].Tables["exist_tb_1"].Size = 17
				return i
			},
			args:    args{query: "create index idx_exist_db_exist_tb_1 on exist_db(v2);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(false); config onlineddl(true); table size enough(false)",
			setUp: func(i *Inspect) *Inspect {
				i.Ctx.schemas["exist_db"].Tables["exist_tb_1"].Size = 15
				return i
			},
			args:    args{query: "create index idx_exist_db_exist_tb_1 on exist_db(v2);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(true); config onlineddl(false); table size enough(false)",
			setUp: func(i *Inspect) *Inspect {
				i.cnf.DDLGhostMinSize = -1
				i.Ctx.schemas["exist_db"].Tables["exist_tb_1"].Size = 15
				return i
			},
			args:    args{query: "alter table exist_db.exist_tb_1 add column col1 varchar(100);"},
			want:    false,
			wantErr: false,
		},

		{
			name: "alter stmt(false); config onlineddl(false); table size enough(false)",
			setUp: func(i *Inspect) *Inspect {
				i.cnf.DDLGhostMinSize = -1
				i.Ctx.schemas["exist_db"].Tables["exist_tb_1"].Size = 15
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
			got, err := tt.setUp(i).onlineddlWithGhost(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Inspect.onlineddlWithGhost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Inspect.onlineddlWithGhost() = %v, want %v", got, tt.want)
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
