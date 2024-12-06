package auditplan

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTidbCompletionSchema(t *testing.T) {
	// https://github.com/actiontech/sqle-ee/issues/395
	sql := "INSERT INTO t1(a1,a2,a3,a4) VALUES('','','Y',CURRENT_DATE)"
	newSQL, err := tidbCompletionSchema(sql, "test")
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `test`.`t1` (`a1`,`a2`,`a3`,`a4`) VALUES ('','','Y',CURRENT_DATE())", newSQL)
}

func TestMergeSQLsByFingerprint(t *testing.T) {
	tests := []struct {
		sqls []SqlFromAliCloud
		want []sqlInfo
	}{
		{
			sqls: []SqlFromAliCloud{{sql: `select sleep(2)`}},
			want: []sqlInfo{
				{
					counter:     1,
					fingerprint: "select sleep(?)",
					sql:         "select sleep(2)",
				},
			},
		}, {
			sqls: []SqlFromAliCloud{
				{
					sql: `select sleep(2)`,
				},
				{
					sql: `select sleep(4)`,
				},
				{
					sql: `select sleep(3)`,
				},
			},
			want: []sqlInfo{
				{
					counter:     3,
					fingerprint: "select sleep(?)",
					sql:         "select sleep(3)",
				},
			},
		}, {
			sqls: []SqlFromAliCloud{{sql: `select * from tb1 where a=1`}},
			want: []sqlInfo{
				{
					counter:     1,
					fingerprint: "select * from tb1 where a=?",
					sql:         "select * from tb1 where a=1",
				},
			},
		}, {
			sqls: []SqlFromAliCloud{
				{
					sql: `select * from tb1 where a=1`,
				},
				{
					sql: `select * from tb1 where a=2`,
				},
				{
					sql: `select * from tb1 where a=3`,
				},
			},
			want: []sqlInfo{
				{
					counter:     3,
					fingerprint: "select * from tb1 where a=?",
					sql:         "select * from tb1 where a=3",
				},
			},
		}, {
			sqls: []SqlFromAliCloud{
				{
					sql: `select * from tb1 where a=1`,
				},
				{
					sql: `select sleep(2)`,
				},
				{
					sql: `select sleep(4)`,
				},
				{
					sql: `select * from tb1 where a=3`,
				},
			},
			want: []sqlInfo{
				{
					counter:     2,
					fingerprint: "select * from tb1 where a=?",
					sql:         "select * from tb1 where a=3",
				},
				{
					counter:     2,
					fingerprint: "select sleep(?)",
					sql:         "select sleep(4)",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test sqls NO.%v", i), func(t *testing.T) {
			actual := mergeSQLsByFingerprint(nil, tt.sqls)
			if !assert.EqualValues(t, actual, tt.want) {
				t.Fatal("unexpected result")
			}
		})
	}
}
