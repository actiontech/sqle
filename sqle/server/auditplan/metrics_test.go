package auditplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricsInt(t *testing.T) {
	ms := NewMetrics()
	ms.SetInt("test_int", 10)
	assert.Equal(t, int64(10), ms["test_int"].i)
	assert.Equal(t, MetricTypeInt, ms["test_int"].typ)
	assert.Equal(t, int64(10), ms.Get("test_int").i)
	assert.Equal(t, int64(10), ms.Get("test_int").Int())
	assert.Equal(t, int64(10), ms.Get("test_int").Value())

	ms.SetInt("test_int", 0)
	assert.Equal(t, int64(0), ms["test_int"].i)
	assert.Equal(t, MetricTypeInt, ms["test_int"].typ)
	assert.Equal(t, int64(0), ms.Get("test_int").i)
	assert.Equal(t, int64(0), ms.Get("test_int").Int())
	assert.Equal(t, int64(0), ms.Get("test_int").Value())

	assert.Equal(t, int64(0), ms.Get("test_int_not_exist").Int())
}

func TestMetricsFloat(t *testing.T) {
	ms := NewMetrics()
	ms.SetFloat("test_float", 10.0)
	assert.Equal(t, 10.0, ms["test_float"].f)
	assert.Equal(t, MetricTypeFloat, ms["test_float"].typ)
	assert.Equal(t, 10.0, ms.Get("test_float").f)
	assert.Equal(t, 10.0, ms.Get("test_float").Float())
	assert.Equal(t, 10.0, ms.Get("test_float").Value())

	ms.SetFloat("test_float", 0)
	assert.Equal(t, float64(0), ms["test_float"].f)
	assert.Equal(t, MetricTypeFloat, ms["test_float"].typ)
	assert.Equal(t, float64(0), ms.Get("test_float").f)
	assert.Equal(t, float64(0), ms.Get("test_float").Float())
	assert.Equal(t, float64(0), ms.Get("test_float").Value())

	assert.Equal(t, float64(0), ms.Get("test_float_not_exist").Float())
}

func TestMetricsString(t *testing.T) {
	ms := NewMetrics()
	ms.SetString("test_string", "test")
	assert.Equal(t, "test", ms["test_string"].s)
	assert.Equal(t, MetricTypeString, ms["test_string"].typ)
	assert.Equal(t, "test", ms.Get("test_string").s)
	assert.Equal(t, "test", ms.Get("test_string").String())
	assert.Equal(t, "test", ms.Get("test_string").Value())

	ms.SetString("test_string", "")
	assert.Equal(t, "", ms["test_string"].s)
	assert.Equal(t, MetricTypeString, ms["test_string"].typ)
	assert.Equal(t, "", ms.Get("test_string").s)
	assert.Equal(t, "", ms.Get("test_string").String())
	assert.Equal(t, "", ms.Get("test_string").Value())

	assert.Equal(t, "", ms.Get("test_string_not_exist").String())
}

func TestMetricsTime(t *testing.T) {
	ms := NewMetrics()
	now1 := time.Now()
	ms.SetTime("test_time", &now1)
	assert.Equal(t, &now1, ms["test_time"].t)
	assert.Equal(t, MetricTypeTime, ms["test_time"].typ)
	assert.Equal(t, &now1, ms.Get("test_time").t)
	assert.Equal(t, &now1, ms.Get("test_time").Time())
	assert.Equal(t, &now1, ms.Get("test_time").Value())

	now2 := time.Now()
	ms.SetTime("test_time", &now2)
	assert.Equal(t, &now2, ms["test_time"].t)
	assert.Equal(t, MetricTypeTime, ms["test_time"].typ)
	assert.Equal(t, &now2, ms.Get("test_time").t)
	assert.Equal(t, &now2, ms.Get("test_time").Time())
	assert.Equal(t, &now2, ms.Get("test_time").Value())
	assert.NotEqual(t, &now1, ms.Get("test_time").Value())

	assert.Nil(t, ms.Get("test_time_not_exist").Time())
}

func TestMetricsBool(t *testing.T) {
	ms := NewMetrics()

	ms.SetBool("test_bool", true)
	assert.Equal(t, true, ms["test_bool"].b)
	assert.Equal(t, MetricTypeBool, ms["test_bool"].typ)
	assert.Equal(t, true, ms.Get("test_bool").b)
	assert.Equal(t, true, ms.Get("test_bool").Bool())
	assert.Equal(t, true, ms.Get("test_bool").Value())

	ms.SetBool("test_bool", false)
	assert.Equal(t, false, ms["test_bool"].b)
	assert.Equal(t, MetricTypeBool, ms["test_bool"].typ)
	assert.Equal(t, false, ms.Get("test_bool").b)
	assert.Equal(t, false, ms.Get("test_bool").Bool())
	assert.Equal(t, false, ms.Get("test_bool").Value())
	assert.NotEqual(t, true, ms.Get("test_bool").Value())

	assert.Equal(t, false, ms.Get("test_bool_not_exist").Bool())
}

func TestMetricsToMap(t *testing.T) {
	ms := NewMetrics()
	now1 := time.Now()

	ms.SetInt("test_int", 10)
	ms.SetInt("test_int", 0)
	ms.SetFloat("test_float", 10.0)
	ms.SetString("test_string", "test")
	ms.SetTime("test_time", &now1)
	ms.SetBool("test_bool", true)
	m := ms.ToMap()

	assert.Equal(t, len(m), 5)
	assert.Equal(t, int64(0), m["test_int"])
	assert.Equal(t, float64(10.0), m["test_float"])
	assert.Equal(t, "test", m["test_string"])
	assert.Equal(t, &now1, m["test_time"])
	assert.Equal(t, true, m["test_bool"])
}

func TestLoadMetrics(t *testing.T) {
	// load int
	ms := LoadMetrics(map[string]interface{}{
		MetricNameCounter: 1,
	}, []string{
		MetricNameCounter,
	})
	assert.Equal(t, int64(1), ms.Get(MetricNameCounter).Int())
	assert.Equal(t, 1, len(ms))

	ms = LoadMetrics(map[string]interface{}{
		MetricNameCounter: int64(2),
	}, []string{
		MetricNameCounter,
	})
	assert.Equal(t, int64(2), ms.Get(MetricNameCounter).Int())
	assert.Equal(t, 1, len(ms))

	ms = LoadMetrics(map[string]interface{}{
		MetricNameCounter: int32(3),
	}, []string{
		MetricNameCounter,
	})
	assert.Equal(t, int64(3), ms.Get(MetricNameCounter).Int())
	assert.Equal(t, 1, len(ms))

	ms = LoadMetrics(map[string]interface{}{
		MetricNameCounter: float64(4),
	}, []string{
		MetricNameCounter,
	})
	assert.Equal(t, int64(4), ms.Get(MetricNameCounter).Int())
	assert.Equal(t, 1, len(ms))

	ms = LoadMetrics(map[string]interface{}{
		MetricNameQueryTimeAvg: 1,
	}, []string{
		MetricNameQueryTimeAvg,
	})
	assert.Equal(t, float64(1), ms.Get(MetricNameQueryTimeAvg).Float())
	assert.Equal(t, 1, len(ms))

	// load float
	ms = LoadMetrics(map[string]interface{}{
		MetricNameQueryTimeAvg: 1.2,
	}, []string{
		MetricNameQueryTimeAvg,
	})
	assert.Equal(t, 1.2, ms.Get(MetricNameQueryTimeAvg).Float())
	assert.Equal(t, 1, len(ms))

	ms = LoadMetrics(map[string]interface{}{
		MetricNameQueryTimeAvg: int64(1),
	}, []string{
		MetricNameQueryTimeAvg,
	})
	assert.Equal(t, float64(1), ms.Get(MetricNameQueryTimeAvg).Float())
	assert.Equal(t, 1, len(ms))

	ms = LoadMetrics(map[string]interface{}{
		MetricNameQueryTimeAvg: int(2),
	}, []string{
		MetricNameQueryTimeAvg,
	})
	assert.Equal(t, float64(2), ms.Get(MetricNameQueryTimeAvg).Float())
	assert.Equal(t, 1, len(ms))

	// load string
	ms = LoadMetrics(map[string]interface{}{
		MetricNameDBUser: "root",
	}, []string{
		MetricNameDBUser,
	})
	assert.Equal(t, "root", ms.Get(MetricNameDBUser).String())
	assert.Equal(t, 1, len(ms))

	// load bool
	assert.Equal(t, false, ms.Get(MetricNameRecordDeleted).Bool())
	ms = LoadMetrics(map[string]interface{}{
		MetricNameRecordDeleted: true,
	}, []string{
		MetricNameRecordDeleted,
	})
	assert.Equal(t, true, ms.Get(MetricNameRecordDeleted).Bool())
	assert.Equal(t, 1, len(ms))

	// test load all
	info := map[string]interface{}{
		MetricNameCounter:              1,
		MetricNameLastReceiveTimestamp: "2024-7-7 10:00:00",
		MetricNameQueryTimeAvg:         2.1,
		MetricNameQueryTimeMax:         10,
		MetricNameRowExaminedAvg:       1000,
		MetricNameDBUser:               "root",
		MetricNameEndpoints:            "1",
		MetricNameRecordDeleted:        true,
	}
	ms = LoadMetrics(info, []string{
		MetricNameCounter,
		MetricNameLastReceiveTimestamp,
		MetricNameQueryTimeAvg,
		MetricNameQueryTimeMax,
		MetricNameRowExaminedAvg,
		MetricNameDBUser,
		MetricNameEndpoints,
		MetricNameRecordDeleted,
	})
	assert.Equal(t, 8, len(ms))

	// 测试没有提供 metrics
	ms = LoadMetrics(info, []string{})
	assert.Equal(t, 0, len(ms))

	// 测试提供少量的metrics
	ms = LoadMetrics(info, []string{
		MetricNameCounter,
	})
	assert.Equal(t, 1, len(ms))
	assert.Equal(t, int64(1), ms.Get(MetricNameCounter).Int())
	assert.Equal(t, float64(0), ms.Get(MetricNameRowExaminedAvg).Float())

	// 测试提供少量的metrics
	ms = LoadMetrics(info, []string{
		MetricNameCounter,
		MetricNameRowExaminedAvg,
	})
	assert.Equal(t, 2, len(ms))
	assert.Equal(t, int64(1), ms.Get(MetricNameCounter).Int())
	assert.Equal(t, float64(1000), ms.Get(MetricNameRowExaminedAvg).Float())

	// 测试提供不存在的metrics
	ms = LoadMetrics(info, []string{
		"not_exist_metrics",
	})
	assert.Equal(t, 0, len(ms))

	// 测试提供不存在的metrics
	ms = LoadMetrics(info, []string{
		"not_exist_metrics",
		MetricNameCounter,
	})
	assert.Equal(t, 1, len(ms))
	assert.Equal(t, int64(1), ms.Get(MetricNameCounter).Int())
}
