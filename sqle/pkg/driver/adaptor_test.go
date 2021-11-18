package driver

import (
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/sijms/go-ora/v2"
	"github.com/stretchr/testify/assert"
)

func Test_splitSQL(t *testing.T) {
	tests := []struct {
		sqls    string
		want    []string
		wantErr bool
	}{
		{sqls: "select * from t1;select * from t1;", want: []string{"select * from t1;", "select * from t1;"}},
		{sqls: "select * from t1;select * from t1;", want: []string{"select * from t1;", "select * from t1;"}},
		{sqls: "select * from t1;select * from t1", want: []string{"select * from t1;", "select * from t1"}},
		{sqls: "select * from `t;1`;select * from `t;1`", want: []string{"select * from `t;1`;", "select * from `t;1`"}},
		{sqls: "insert into t1(name) values ('hello;world');insert into t1(name) values ('hello;world')", want: []string{"insert into t1(name) values ('hello;world');", "insert into t1(name) values ('hello;world')"}},
		{sqls: "insert into t1(name) values (\"hello;world\");insert into t1(name) values (\"hello;world\")", want: []string{"insert into t1(name) values (\"hello;world\");", "insert into t1(name) values (\"hello;world\")"}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := splitSQL(tt.sqls)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
