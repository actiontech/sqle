package mysql

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"

	dmsCommonSQLOp "github.com/actiontech/dms/pkg/dms-common/sql_op"
	"github.com/actiontech/sqle/sqle/log"
)

var (
	testSQLObjectOpT1Read = &dmsCommonSQLOp.SQLObjectOp{
		Op: dmsCommonSQLOp.SQLOpRead,
		Object: &dmsCommonSQLOp.SQLObject{
			Type:         dmsCommonSQLOp.SQLObjectTypeTable,
			DatabaseName: "s1",
			TableName:    "t1",
		},
	}
	testSQLObjectOpT1AddOrUpdate = &dmsCommonSQLOp.SQLObjectOp{
		Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
		Object: &dmsCommonSQLOp.SQLObject{
			Type:         dmsCommonSQLOp.SQLObjectTypeTable,
			DatabaseName: "s1",
			TableName:    "t1",
		},
	}
)

func TestSQLObjectOpsDuplicateRemoval(t *testing.T) {
	type args struct {
		ops []*dmsCommonSQLOp.SQLObjectOp
	}
	tests := []struct {
		name string
		args args
		want []*dmsCommonSQLOp.SQLObjectOp
	}{
		{
			name: "test1",
			args: args{ops: []*dmsCommonSQLOp.SQLObjectOp{testSQLObjectOpT1Read, testSQLObjectOpT1Read}},
			want: []*dmsCommonSQLOp.SQLObjectOp{testSQLObjectOpT1Read},
		},
		{
			name: "test1",
			args: args{ops: []*dmsCommonSQLOp.SQLObjectOp{testSQLObjectOpT1Read, testSQLObjectOpT1Read, testSQLObjectOpT1AddOrUpdate}},
			want: []*dmsCommonSQLOp.SQLObjectOp{testSQLObjectOpT1Read, testSQLObjectOpT1AddOrUpdate},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SQLObjectOpsDuplicateRemoval(tt.args.ops); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLObjectOpsDuplicateRemoval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMysqlDriverImpl_GetSQLOp(t *testing.T) {
	type args struct {
		ctx  context.Context
		sqls string
	}
	tests := []struct {
		name    string
		args    args
		want    []*dmsCommonSQLOp.SQLObjectOps
		wantErr error
	}{
		{
			name: "multi sql",
			args: args{ctx: context.Background(), sqls: "select * from s1.t1;select * from s1.t2"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "s1",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "select * from s1.t1;",
					},
				},
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "s1",
								TableName:    "t2",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "select * from s1.t2",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "select basic",
			args: args{ctx: context.Background(), sqls: "select * from s1.t1"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "s1",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "select * from s1.t1",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "select no db",
			args: args{ctx: context.Background(), sqls: "select * from t1"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "select * from t1",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "select 1",
			args: args{ctx: context.Background(), sqls: "select 1"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "select 1",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "select multi table",
			args: args{ctx: context.Background(), sqls: "select * from t1,t2"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t2",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "select * from t1,t2",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "select join table",
			args: args{ctx: context.Background(), sqls: "SELECT * FROM t1 INNER JOIN t2"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t2",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "SELECT * FROM t1 INNER JOIN t2",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "select into outfile",
			args: args{ctx: context.Background(),
				sqls: "SELECT * FROM t1 INTO OUTFILE '/tmp/select-values.txt'"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeServer,
								DatabaseName: "",
								TableName:    "",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "SELECT * FROM t1 INTO OUTFILE '/tmp/select-values.txt'",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "create table like",
			args: args{ctx: context.Background(), sqls: "CREATE TABLE new_tbl LIKE orig_tbl"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "new_tbl",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "orig_tbl",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "CREATE TABLE new_tbl LIKE orig_tbl",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "create table basic",
			args: args{ctx: context.Background(), sqls: "CREATE TABLE t (c CHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin);"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "CREATE TABLE t (c CHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin);",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "drop table basic",
			args: args{ctx: context.Background(), sqls: "DROP TABLE t1;"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "DROP TABLE t1;",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "alter table basic",
			args: args{ctx: context.Background(), sqls: "ALTER TABLE t2 DROP COLUMN c, DROP COLUMN d;"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t2",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "ALTER TABLE t2 DROP COLUMN c, DROP COLUMN d;",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "alter table rename",
			args: args{ctx: context.Background(), sqls: "ALTER TABLE old_table RENAME new_table;"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "new_table",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "old_table",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "old_table",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "ALTER TABLE old_table RENAME new_table;",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "rename table",
			args: args{ctx: context.Background(), sqls: "RENAME TABLE old_table TO new_table;"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "new_table",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "old_table",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "old_table",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "RENAME TABLE old_table TO new_table;",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "truncate table",
			args: args{ctx: context.Background(), sqls: "TRUNCATE TABLE t1"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "TRUNCATE TABLE t1",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "repair table",
			args:    args{ctx: context.Background(), sqls: "REPAIR TABLE t1"},
			wantErr: fmt.Errorf("there is unparsed stmt: REPAIR TABLE t1"),
		},
		{
			name: "alter databases",
			args: args{ctx: context.Background(), sqls: "ALTER DATABASE myDatabase CHARACTER SET= ascii;"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
								DatabaseName: "myDatabase",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "ALTER DATABASE myDatabase CHARACTER SET= ascii;",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "drop databases",
			args: args{ctx: context.Background(), sqls: "DROP DATABASE mydb"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
								DatabaseName: "mydb",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "DROP DATABASE mydb",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "create view",
			args: args{ctx: context.Background(), sqls: "CREATE VIEW test.v AS SELECT * FROM t;"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "test",
								TableName:    "v",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "CREATE VIEW test.v AS SELECT * FROM t;",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "create view select CURRENT_DATE",
			args: args{ctx: context.Background(), sqls: "CREATE VIEW v_today (today) AS SELECT CURRENT_DATE;"},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "v_today",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: "CREATE VIEW v_today (today) AS SELECT CURRENT_DATE;",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "create or replace view",
			args: args{ctx: context.Background(), sqls: `CREATE OR REPLACE VIEW view_name AS
SELECT column_name(s)
FROM table_name
WHERE condition`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "view_name",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "view_name",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "table_name",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `CREATE OR REPLACE VIEW view_name AS
SELECT column_name(s)
FROM table_name
WHERE condition`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "create view union",
			args: args{ctx: context.Background(), sqls: `CREATE VIEW netcheck.cpu_mp AS
(SELECT
 cpu.ID AS id,
 cpu.chanel_name AS chanel_name,
 cpu.first_channel AS first_channel,
 cpu.IMG_Url AS IMG_Url,
 cpu.lastModifyTime AS lastModifyTime,
 cpu.second_channel AS second_channel,
 cpu.SHOW_TIME AS SHOW_TIME,
 cpu.TASK_Id AS TASK_Id,
 cpu.TITLE AS TITLE,
 cpu.URL AS URL,
 cpu.checkSysTaskId AS checkSysTaskId,
 cpu.innerUUID AS innerUUID,
 cpu.isReject AS isReject,
 cpu.scanTime AS scanTime
FROM channel_page_update_result cpu
)
 UNION ALL
(SELECT
  mp.id AS id,
  '' AS chanel_name,
  '' AS first_channel,
  '' AS second_channel,
  '' AS TITLE,
  mp.imgUrl AS IMG_Url
  ,mp.lastModifyTime AS lastModifyTime,
  mp.showTime AS SHOW_TIME,
  mp.taskId AS TASK_Id,
  mp.url AS URL,
  mp.checkSysTaskId AS checkSysTaskId,
  mp.innerUUID AS innerUUID,
  mp.isReject AS isReject,
  mp.scanTime AS scanTime
FROM mainpageupdateresult mp
);`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "netcheck",
								TableName:    "cpu_mp",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "channel_page_update_result",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "mainpageupdateresult",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `CREATE VIEW netcheck.cpu_mp AS
(SELECT
 cpu.ID AS id,
 cpu.chanel_name AS chanel_name,
 cpu.first_channel AS first_channel,
 cpu.IMG_Url AS IMG_Url,
 cpu.lastModifyTime AS lastModifyTime,
 cpu.second_channel AS second_channel,
 cpu.SHOW_TIME AS SHOW_TIME,
 cpu.TASK_Id AS TASK_Id,
 cpu.TITLE AS TITLE,
 cpu.URL AS URL,
 cpu.checkSysTaskId AS checkSysTaskId,
 cpu.innerUUID AS innerUUID,
 cpu.isReject AS isReject,
 cpu.scanTime AS scanTime
FROM channel_page_update_result cpu
)
 UNION ALL
(SELECT
  mp.id AS id,
  '' AS chanel_name,
  '' AS first_channel,
  '' AS second_channel,
  '' AS TITLE,
  mp.imgUrl AS IMG_Url
  ,mp.lastModifyTime AS lastModifyTime,
  mp.showTime AS SHOW_TIME,
  mp.taskId AS TASK_Id,
  mp.url AS URL,
  mp.checkSysTaskId AS checkSysTaskId,
  mp.innerUUID AS innerUUID,
  mp.isReject AS isReject,
  mp.scanTime AS scanTime
FROM mainpageupdateresult mp
);`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "create index",
			args: args{ctx: context.Background(), sqls: `CREATE INDEX part_of_name ON customer (name(10));`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "customer",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `CREATE INDEX part_of_name ON customer (name(10));`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "drop index",
			args: args{ctx: context.Background(), sqls: `DROP INDEX i1 ON t;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `DROP INDEX i1 ON t;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "lock table",
			args: args{ctx: context.Background(), sqls: `LOCK TABLES t1 READ`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `LOCK TABLES t1 READ`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "unlock table",
			args: args{ctx: context.Background(), sqls: `UNLOCK TABLES;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `UNLOCK TABLES;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "union",
			args: args{ctx: context.Background(), sqls: `SELECT 1, 2 UNION SELECT 'a', 'b';`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SELECT 1, 2 UNION SELECT 'a', 'b';`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "union table",
			args: args{ctx: context.Background(), sqls: `SELECT city FROM customers
UNION
SELECT city FROM suppliers
ORDER BY city;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "customers",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "suppliers",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SELECT city FROM customers
UNION
SELECT city FROM suppliers
ORDER BY city;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "load data",
			args: args{ctx: context.Background(), sqls: `LOAD DATA INFILE 'data.txt' INTO TABLE db2.my_table;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeServer,
								DatabaseName: "",
								TableName:    "",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "db2",
								TableName:    "my_table",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `LOAD DATA INFILE 'data.txt' INTO TABLE db2.my_table;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "load data local",
			args: args{ctx: context.Background(), sqls: `LOAD DATA LOCAL INFILE 'data.txt' INTO TABLE db2.my_table;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "db2",
								TableName:    "my_table",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `LOAD DATA LOCAL INFILE 'data.txt' INTO TABLE db2.my_table;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "insert",
			args: args{ctx: context.Background(), sqls: `INSERT INTO tbl_name () VALUES();`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "tbl_name",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `INSERT INTO tbl_name () VALUES();`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "insert select from",
			args: args{ctx: context.Background(), sqls: `INSERT INTO tbl_temp2 (fld_id) 
SELECT tbl_temp1.fld_order_id FROM tbl_temp1 WHERE tbl_temp1.fld_order_id > 100;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "tbl_temp2",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "tbl_temp1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `INSERT INTO tbl_temp2 (fld_id) 
SELECT tbl_temp1.fld_order_id FROM tbl_temp1 WHERE tbl_temp1.fld_order_id > 100;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "delete where",
			args: args{ctx: context.Background(), sqls: `DELETE FROM somelog WHERE user = 'jcole' ORDER BY timestamp_column LIMIT 1;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "somelog",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "somelog",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `DELETE FROM somelog WHERE user = 'jcole' ORDER BY timestamp_column LIMIT 1;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "delete where use alias",
			args: args{ctx: context.Background(), sqls: `DELETE FROM somelog AS s WHERE s.user = 'jcole' ORDER BY timestamp_column LIMIT 1;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "somelog",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "somelog",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `DELETE FROM somelog AS s WHERE s.user = 'jcole' ORDER BY timestamp_column LIMIT 1;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "delete table",
			args: args{ctx: context.Background(), sqls: `DROP TABLE t_old;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t_old",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `DROP TABLE t_old;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "delete multi table",
			args: args{ctx: context.Background(), sqls: `DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpDelete,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t2",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t2",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t3",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "update table basic",
			args: args{ctx: context.Background(), sqls: `UPDATE t1 SET col1 = col1 + 1;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `UPDATE t1 SET col1 = col1 + 1;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "update table with two table",
			args: args{ctx: context.Background(), sqls: `UPDATE items,month SET items.price=month.price WHERE items.id=month.id;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAddOrUpdate,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "items",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "items",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "month",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `UPDATE items,month SET items.price=month.price WHERE items.id=month.id;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show columns",
			args: args{ctx: context.Background(), sqls: `SHOW COLUMNS FROM mytable FROM mydb;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "mydb",
								TableName:    "mytable",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW COLUMNS FROM mytable FROM mydb;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show create table",
			args: args{ctx: context.Background(), sqls: `SHOW CREATE TABLE t`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW CREATE TABLE t`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show create user",
			args: args{ctx: context.Background(), sqls: `SHOW CREATE USER 'u1'@'localhost'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
								DatabaseName: "mysql",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW CREATE USER 'u1'@'localhost'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show grants",
			args: args{ctx: context.Background(), sqls: `SHOW GRANTS FOR 'jeffrey'@'localhost';`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
								DatabaseName: "mysql",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW GRANTS FOR 'jeffrey'@'localhost';`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show procedure status",
			args: args{ctx: context.Background(), sqls: `SHOW PROCEDURE STATUS LIKE 'sp1'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW PROCEDURE STATUS LIKE 'sp1'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show index",
			args: args{ctx: context.Background(), sqls: `SHOW INDEX FROM City`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "city",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW INDEX FROM City`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show create databases",
			args: args{ctx: context.Background(), sqls: `SHOW CREATE DATABASE test`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
								DatabaseName: "test",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW CREATE DATABASE test`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show events",
			args: args{ctx: context.Background(), sqls: `SHOW EVENTS FROM test`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
								DatabaseName: "test",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW EVENTS FROM test`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "show master status",
			args: args{ctx: context.Background(), sqls: `SHOW MASTER STATUS`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHOW MASTER STATUS`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "explain for connection",
			args: args{ctx: context.Background(), sqls: `EXPLAIN FOR CONNECTION 4`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `EXPLAIN FOR CONNECTION 4`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "explain",
			args: args{ctx: context.Background(), sqls: `EXPLAIN ANALYZE SELECT * FROM t1 JOIN t2 ON (t1.c1 = t2.c2)`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t1",
							},
						},
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "t2",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `EXPLAIN ANALYZE SELECT * FROM t1 JOIN t2 ON (t1.c1 = t2.c2)`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "prepare",
			args: args{ctx: context.Background(), sqls: `PREPARE stmt1 FROM 'SELECT productCode, productName FROM products WHERE productCode = ?'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpRead,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "",
								TableName:    "products",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `PREPARE stmt1 FROM 'SELECT productCode, productName FROM products WHERE productCode = ?'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "binlog",
			args: args{ctx: context.Background(), sqls: `BINLOG 'str'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `BINLOG 'str'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "flush",
			args: args{ctx: context.Background(), sqls: `FLUSH BINARY LOGS`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `FLUSH BINARY LOGS`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "kill",
			args: args{ctx: context.Background(), sqls: `KILL 10`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `KILL 10`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "set password",
			args: args{ctx: context.Background(), sqls: `SET PASSWORD FOR 'jeffrey'@'localhost' = 'auth_string'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SET PASSWORD FOR 'jeffrey'@'localhost' = 'auth_string'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "create user",
			args: args{ctx: context.Background(), sqls: `CREATE USER 'jeffrey'@'localhost' IDENTIFIED BY 'password';`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `CREATE USER 'jeffrey'@'localhost' IDENTIFIED BY 'password';`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "alter user",
			args: args{ctx: context.Background(), sqls: `ALTER USER 'jeffrey'@'localhost' ACCOUNT LOCK`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `ALTER USER 'jeffrey'@'localhost' ACCOUNT LOCK`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "alter instance",
			args: args{ctx: context.Background(), sqls: `ALTER INSTANCE RELOAD TLS;`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `ALTER INSTANCE RELOAD TLS;`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "drop user",
			args: args{ctx: context.Background(), sqls: `DROP USER 'jeffrey'@'localhost';`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `DROP USER 'jeffrey'@'localhost';`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "revoke instance level",
			args: args{ctx: context.Background(), sqls: `REVOKE INSERT ON *.* FROM 'jeffrey'@'localhost'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpGrant,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `REVOKE INSERT ON *.* FROM 'jeffrey'@'localhost'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "revoke db level",
			args: args{ctx: context.Background(), sqls: `REVOKE INSERT ON db1.* FROM 'jeffrey'@'localhost'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpGrant,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
								DatabaseName: "db1",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `REVOKE INSERT ON db1.* FROM 'jeffrey'@'localhost'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "revoke table level",
			args: args{ctx: context.Background(), sqls: `REVOKE INSERT ON db1.t1 FROM 'jeffrey'@'localhost'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpGrant,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "db1",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `REVOKE INSERT ON db1.t1 FROM 'jeffrey'@'localhost'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "grant instance level",
			args: args{ctx: context.Background(), sqls: `GRANT ALL ON *.* TO 'jeffrey'@'localhost'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpGrant,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `GRANT ALL ON *.* TO 'jeffrey'@'localhost'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "grant db level",
			args: args{ctx: context.Background(), sqls: `GRANT ALL ON db1.* TO 'jeffrey'@'localhost'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpGrant,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
								DatabaseName: "db1",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `GRANT ALL ON db1.* TO 'jeffrey'@'localhost'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "grant table level",
			args: args{ctx: context.Background(), sqls: `GRANT ALL ON db1.t1 TO 'jeffrey'@'localhost'`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpGrant,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeTable,
								DatabaseName: "db1",
								TableName:    "t1",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `GRANT ALL ON db1.t1 TO 'jeffrey'@'localhost'`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "shutdown",
			args: args{ctx: context.Background(), sqls: `SHUTDOWN`},
			want: []*dmsCommonSQLOp.SQLObjectOps{
				{
					ObjectOps: []*dmsCommonSQLOp.SQLObjectOp{
						{
							Op: dmsCommonSQLOp.SQLOpAdmin,
							Object: &dmsCommonSQLOp.SQLObject{
								Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
								DatabaseName: "",
								TableName:    "",
							},
						},
					},
					Sql: dmsCommonSQLOp.SQLInfo{
						Sql: `SHUTDOWN`,
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "unparsed sql",
			args:    args{ctx: context.Background(), sqls: `SELECT * FROMa t1`},
			wantErr: fmt.Errorf("there is unparsed stmt: SELECT * FROMa t1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &MysqlDriverImpl{log: log.NewEntry()}
			got, err := i.GetSQLOp(tt.args.ctx, tt.args.sqls)
			if nil == err && nil == tt.wantErr {
				if !isResultEqual(got, tt.want) {
					t.Errorf("MysqlDriverImpl.GetSQLOp() = %v, want %v", SQLObjectOpsFingerPrint(got), SQLObjectOpsFingerPrint(tt.want))
				}
				return
			}
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.wantErr) {
				t.Errorf("MysqlDriverImpl.GetSQLOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func isResultEqual(a, b []*dmsCommonSQLOp.SQLObjectOps) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !isSQLObjectOpsEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func isSQLObjectOpsEqual(a, b *dmsCommonSQLOp.SQLObjectOps) bool {
	if len(a.ObjectOps) != len(b.ObjectOps) {
		return false
	}
	if a.Sql.Sql != b.Sql.Sql {
		return false
	}
	sort.Slice(a.ObjectOps, func(i, j int) bool {
		s1 := SQLObjectOpFingerPrint(a.ObjectOps[i])
		s2 := SQLObjectOpFingerPrint(a.ObjectOps[j])
		return s1 < s2
	})

	sort.Slice(b.ObjectOps, func(i, j int) bool {
		s1 := SQLObjectOpFingerPrint(b.ObjectOps[i])
		s2 := SQLObjectOpFingerPrint(b.ObjectOps[j])
		return s1 < s2
	})

	for i := range a.ObjectOps {
		if !isSQLObjectOpEqual(a.ObjectOps[i], b.ObjectOps[i]) {
			return false
		}
	}
	return true
}

func isSQLObjectOpEqual(a, b *dmsCommonSQLOp.SQLObjectOp) bool {
	if a.Op != b.Op {
		return false
	}
	return isSQLObjectEqual(a.Object, b.Object)
}

func isSQLObjectEqual(a, b *dmsCommonSQLOp.SQLObject) bool {
	if a.Type != b.Type {
		return false
	}
	if a.DatabaseName != b.DatabaseName {
		return false
	}
	if a.SchemaName != b.SchemaName {
		return false
	}
	if a.TableName != b.TableName {
		return false
	}
	return true
}
