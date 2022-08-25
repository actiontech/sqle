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

func TestDeduplicateSQLsByFingerprint(t *testing.T) {
	tests := []struct {
		sqls []string
		want []*sqlInfo
	}{
		{
			sqls: []string{`select sleep(2)`},
			want: []*sqlInfo{
				{
					counter:     1,
					fingerprint: "select sleep(?)",
					sql:         "select sleep(2)",
				},
			},
		}, {
			sqls: []string{`select sleep(2)`, `select sleep(4)`, `select sleep(3)`},
			want: []*sqlInfo{
				{
					counter:     3,
					fingerprint: "select sleep(?)",
					sql:         "select sleep(3)",
				},
			},
		}, {
			sqls: []string{`select * from tb1 where a=1`},
			want: []*sqlInfo{
				{
					counter:     1,
					fingerprint: "select * from tb1 where a=?",
					sql:         "select * from tb1 where a=1",
				},
			},
		}, {
			sqls: []string{`select * from tb1 where a=1`, `select * from tb1 where a=2`, `select * from tb1 where a=3`},
			want: []*sqlInfo{
				{
					counter:     3,
					fingerprint: "select * from tb1 where a=?",
					sql:         "select * from tb1 where a=3",
				},
			},
		}, {
			sqls: []string{`select * from tb1 where a=1`, `select sleep(2)`, `select sleep(4)`, `select * from tb1 where a=3`},
			want: []*sqlInfo{
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
			actual := deduplicateSQLsByFingerprint(tt.sqls)
			if !assert.EqualValues(t, actual, tt.want) {
				t.Fatal("unexpected result")
			}
		})
	}
}
