//go:build enterprise
// +build enterprise

package auditplan

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMergeSlowlogSQLsByFingerprint(t *testing.T) {
	log := logrus.WithField("test", "test")
	log.Level = logrus.DebugLevel
	cases := []struct {
		sqls     []*sqlFromSlowLog
		expected []sqlInfo
	}{
		{
			sqls: []*sqlFromSlowLog{
				{sql: "set names utf8", schema: "", queryTimeSeconds: 2},
				{sql: "set names utf8", schema: "", queryTimeSeconds: 1},
				{sql: "set names utf8", schema: "", queryTimeSeconds: 3},
			},
			expected: []sqlInfo{
				{counter: 3, fingerprint: "SET NAMES ?", sql: "set names utf8", schema: "", queryTimeSeconds: 2},
			},
		},
		{
			sqls: []*sqlFromSlowLog{
				{sql: "select sleep(2)", schema: "", queryTimeSeconds: 2},
				{sql: "select sleep(3)", schema: "", queryTimeSeconds: 3},
				{sql: "select sleep(4)", schema: "", queryTimeSeconds: 4},
			},
			expected: []sqlInfo{
				{counter: 3, fingerprint: "SELECT SLEEP(?)", sql: "select sleep(4)", schema: "", queryTimeSeconds: 3},
			},
		},
		{
			sqls: []*sqlFromSlowLog{
				{sql: "select * from tb1 where a=1", schema: "tb1", queryTimeSeconds: 1},
				{sql: "select sleep(2)", schema: "", queryTimeSeconds: 2},
				{sql: "select sleep(4)", schema: "", queryTimeSeconds: 4},
				{sql: "select * from tb1 where a=3", schema: "tb1", queryTimeSeconds: 3},
			},
			expected: []sqlInfo{
				{counter: 2, fingerprint: "SELECT * FROM `tb1` WHERE `a`=?", sql: "select * from tb1 where a=3", schema: "tb1", queryTimeSeconds: 2},
				{counter: 2, fingerprint: "SELECT SLEEP(?)", sql: "select sleep(4)", schema: "", queryTimeSeconds: 3},
			},
		},
	}

	for i := range cases {
		c := cases[i]
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			actual, _ := sqlFromSlowLogs(c.sqls).mergeByFingerprint(log)
			assert.EqualValues(t, c.expected, actual)
		})
	}
}
