//go:build enterprise
// +build enterprise

package auditplan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestMergeSlowlogSQLsByFingerprint(t *testing.T) {
// 	log := logrus.WithField("test", "test")
// 	log.Level = logrus.DebugLevel
// 	cases := []struct {
// 		sqls     []*sqlFromSlowLog
// 		expected []sqlInfo
// 	}{
// 		{
// 			sqls: []*sqlFromSlowLog{
// 				{sql: "set names utf8", schema: "", queryTimeSeconds: 2},
// 				{sql: "set names utf8", schema: "", queryTimeSeconds: 1},
// 				{sql: "set names utf8", schema: "", queryTimeSeconds: 3},
// 			},
// 			expected: []sqlInfo{
// 				{counter: 3, fingerprint: "SET NAMES ?", sql: "set names utf8", schema: "", queryTimeSeconds: 2},
// 			},
// 		},
// 		{
// 			sqls: []*sqlFromSlowLog{
// 				{sql: "select sleep(2)", schema: "", queryTimeSeconds: 2},
// 				{sql: "select sleep(3)", schema: "", queryTimeSeconds: 3},
// 				{sql: "select sleep(4)", schema: "", queryTimeSeconds: 4},
// 			},
// 			expected: []sqlInfo{
// 				{counter: 3, fingerprint: "SELECT SLEEP(?)", sql: "select sleep(4)", schema: "", queryTimeSeconds: 3},
// 			},
// 		},
// 		{
// 			sqls: []*sqlFromSlowLog{
// 				{sql: "select * from tb1 where a=1", schema: "tb1", queryTimeSeconds: 1},
// 				{sql: "select sleep(2)", schema: "", queryTimeSeconds: 2},
// 				{sql: "select sleep(4)", schema: "", queryTimeSeconds: 4},
// 				{sql: "select * from tb1 where a=3", schema: "tb1", queryTimeSeconds: 3},
// 			},
// 			expected: []sqlInfo{
// 				{counter: 2, fingerprint: "SELECT * FROM `tb1` WHERE `a`=?", sql: "select * from tb1 where a=3", schema: "tb1", queryTimeSeconds: 2},
// 				{counter: 2, fingerprint: "SELECT SLEEP(?)", sql: "select sleep(4)", schema: "", queryTimeSeconds: 3},
// 			},
// 		},
// 	}

// 	for i := range cases {
// 		c := cases[i]
// 		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
// 			actual, err := sqlFromSlowLogs(c.sqls).mergeByFingerprint(log)
// 			assert.NoError(t, err)
// 			assert.EqualValues(t, c.expected, actual)
// 		})
// 	}
// }

func TestSlowSQLMerge(t *testing.T) {
	newSQLTest := func() *SQLV2 {
		return &SQLV2{
			SQLId:      "test",
			SQLContent: "select 1",
			Info: LoadMetrics(map[string]interface{}{
				MetricNameCounter:              1,
				MetricNameLastReceiveTimestamp: "2000-1-1 1:00:00",
				MetricNameQueryTimeAvg:         1,
				MetricNameQueryTimeMax:         1,
				MetricNameRowExaminedAvg:       1,
				MetricNameDBUser:               "root",
				MetricNameEndpoints:            "1",
			}, []string{
				MetricNameCounter,
				MetricNameLastReceiveTimestamp,
				MetricNameQueryTimeAvg,
				MetricNameQueryTimeMax,
				MetricNameRowExaminedAvg,
				MetricNameDBUser,
				MetricNameEndpoints,
			}),
		}
	}
	task := &SlowLogTaskV2{}

	// sql content
	o := newSQLTest()
	m := newSQLTest()
	m.SQLContent = "select 2"
	task.mergeSQL(o, m)
	assert.Equal(t, "select 2", o.SQLContent)

	// MetricNameCounter
	o = newSQLTest()
	m = newSQLTest()
	m.Info.SetInt(MetricNameCounter, 1)
	task.mergeSQL(o, m)
	assert.Equal(t, int64(2), o.Info.Get(MetricNameCounter).Value())

	o = newSQLTest()
	m = newSQLTest()
	m.Info.SetInt(MetricNameCounter, 1000)
	task.mergeSQL(o, m)
	assert.Equal(t, int64(1001), o.Info.Get(MetricNameCounter).Value())

	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetInt(MetricNameCounter, 1000)
	task.mergeSQL(o, m)
	assert.Equal(t, int64(1001), o.Info.Get(MetricNameCounter).Value())

	// MetricNameLastReceiveTimestamp
	o = newSQLTest()
	m = newSQLTest()
	m.Info.SetString(MetricNameLastReceiveTimestamp, "2000-1-2 1:00:00")
	task.mergeSQL(o, m)
	assert.Equal(t, "2000-1-2 1:00:00", o.Info.Get(MetricNameLastReceiveTimestamp).Value())

	// MetricNameQueryTimeAvg
	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetFloat(MetricNameQueryTimeAvg, 10.0)
	o.Info.SetInt(MetricNameCounter, 100)

	m.Info.SetFloat(MetricNameQueryTimeAvg, 20.0)
	m.Info.SetInt(MetricNameCounter, 100)

	task.mergeSQL(o, m)
	assert.Equal(t, float64((10*100+20*100)/200), o.Info.Get(MetricNameQueryTimeAvg).Float())
	assert.Equal(t, float64(15), o.Info.Get(MetricNameQueryTimeAvg).Float())

	// MetricNameQueryTimeMax
	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetFloat(MetricNameQueryTimeMax, 10.0)

	task.mergeSQL(o, m)
	assert.Equal(t, float64(10), o.Info.Get(MetricNameQueryTimeMax).Float())

	o = newSQLTest()
	m = newSQLTest()
	m.Info.SetFloat(MetricNameQueryTimeMax, 11.0)

	task.mergeSQL(o, m)
	assert.Equal(t, float64(11), o.Info.Get(MetricNameQueryTimeMax).Float())

	// MetricNameRowExaminedAvg
	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetFloat(MetricNameRowExaminedAvg, 1000.0)
	o.Info.SetInt(MetricNameCounter, 100)

	m.Info.SetFloat(MetricNameRowExaminedAvg, 3000.0)
	m.Info.SetInt(MetricNameCounter, 100)

	task.mergeSQL(o, m)
	assert.Equal(t, float64((1000*100+3000*100)/200), o.Info.Get(MetricNameRowExaminedAvg).Float())
	assert.Equal(t, float64(2000), o.Info.Get(MetricNameRowExaminedAvg).Float())

	// MetricNameDBUser
	o = newSQLTest()
	m = newSQLTest()
	m.Info.SetString(MetricNameDBUser, "root2")
	task.mergeSQL(o, m)
	assert.Equal(t, "root2", o.Info.Get(MetricNameDBUser).Value())

	// MetricNameEndpoints
	o = newSQLTest()
	m = newSQLTest()
	m.Info.SetString(MetricNameEndpoints, "2")
	task.mergeSQL(o, m)
	assert.Equal(t, "2", o.Info.Get(MetricNameEndpoints).Value())
}

func TestSchemaMetaMerge(t *testing.T) {
	newSQLTest := func() *SQLV2 {
		return &SQLV2{
			SQLId:       "test",
			SQLContent:  "select 1",
			Fingerprint: "select 1",
			Info: LoadMetrics(map[string]interface{}{}, []string{
				MetricNameRecordDeleted,
				MetricNameMetaName,
				MetricNameMetaType,
			}),
		}
	}
	task := &BaseSchemaMetaTaskV2{}

	// sql content
	o := newSQLTest()
	m := newSQLTest()
	m.SQLContent = "select 2"
	task.mergeSQL(o, m)
	assert.Equal(t, "select 2", o.SQLContent)

	// sql fp
	o = newSQLTest()
	m = newSQLTest()
	m.Fingerprint = "select 2"
	task.mergeSQL(o, m)
	assert.Equal(t, "select 2", o.Fingerprint)

	// record deleted
	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetBool(MetricNameRecordDeleted, true)
	m.Info.SetBool(MetricNameRecordDeleted, false)
	task.mergeSQL(o, m)
	assert.Equal(t, false, o.Info.Get(MetricNameRecordDeleted).Value())

	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetBool(MetricNameRecordDeleted, false)
	m.Info.SetBool(MetricNameRecordDeleted, true)
	task.mergeSQL(o, m)
	assert.Equal(t, true, o.Info.Get(MetricNameRecordDeleted).Value())

	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetBool(MetricNameRecordDeleted, false)
	m.Info.SetBool(MetricNameRecordDeleted, false)
	task.mergeSQL(o, m)
	assert.Equal(t, false, o.Info.Get(MetricNameRecordDeleted).Value())

	// meta name
	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetString(MetricNameMetaName, "test1")
	m.Info.SetString(MetricNameMetaName, "test2")
	task.mergeSQL(o, m)
	assert.Equal(t, "test2", o.Info.Get(MetricNameMetaName).Value())

	// meta type
	o = newSQLTest()
	m = newSQLTest()
	o.Info.SetString(MetricNameMetaType, "view")
	m.Info.SetString(MetricNameMetaType, "table")
	task.mergeSQL(o, m)
	assert.Equal(t, "table", o.Info.Get(MetricNameMetaType).Value())
}
