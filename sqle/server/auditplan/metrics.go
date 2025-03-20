package auditplan

import (
	"time"
)

const MetricNameCounter string = "counter" // 总次数
const MetricNameLastReceiveTimestamp string = "last_receive_timestamp"
const MetricNameQueryTimeAvg string = "query_time_avg"     // 平均执行时间
const MetricNameQueryTimeMax string = "query_time_max"     // 最大执行时间
const MetricNameRowExaminedAvg string = "row_examined_avg" // 平均扫描行数
const MetricNameFirstQueryAt string = "first_query_at"
const MetricNameDBUser string = "db_user"
const MetricNameEndpoints string = "endpoints"
const MetricNameStartTimeOfLastScrapedSQL string = "start_time_of_last_scraped_sql" // 抓取sql的开始时间

const MetricNameMetaName string = "schema_meta_name"    // 表或者视图的名字
const MetricNameMetaType string = "schema_meta_type"    // 表或者视图等等
const MetricNameRecordDeleted string = "record_deleted" // 标记记录是否被删除掉

const MetricNameQueryTimeTotal string = "query_time_total" // 总执行时间
const MetricNameCPUTimeAvg string = "cpu_time_avg"
const MetricNameLockWaitTimeTotal string = "lock_wait_time_total"
const MetricNameLockWaitCounter string = "lock_wait_counter"
const MetricNameActiveWaitTimeTotal string = "act_wait_time_total"
const MetricNameActiveTimeTotal string = "act_time_total"

const MetricNameCPUTimeTotal string = "cpu_time_total"
const MetricNamePhyReadPageTotal string = "phy_read_page_total"
const MetricNameLogicReadPageTotal string = "logic_read_page_total"
const MetricNameUserIOWaitTimeTotal string = "user_io_wait_time_total"
const MetricNameDiskReadTotal string = "disk_read_total"
const MetricNameBufferGetCounter string = "buffer_read_total"

const MetricNameLastQueryAt = "last_query_at"
const MetricNameIoWaitTimeAvg = "io_wait_time_avg"
const MetricNameDiskReadAvg = "disk_read_avg"
const MetricNameBufferReadAvg = "buffer_read_avg"
const MetricNameExplainCost = "explain_cost"

// Lock
const MetricNameLockType string = "lock_type"
const MetricNameLockMode string = "lock_mode"
const MetricNameLockStatus string = "lock_status"
const MetricNameTrxStarted string = "trx_started"
const MetricNameTrxWaitStarted string = "trx_wait_started"
const MetricNameEngine string = "engine"
const MetricNameTable string = "table_name"
const MetricNameIndexName string = "index_name"
const MetricNameObjectName string = "object_name"
const MetricNameGrantedLockSql string = "granted_lock_sql"
const MetricNameWaitingLockSql string = "waiting_lock_sql"
const MetricNameGrantedLockConnectionId string = "granted_lock_connection_id"
const MetricNameWaitingLockConnectionId string = "waiting_lock_connection_id"
const MetricNameGrantedLockTrxId string = "granted_lock_trx_id"
const MetricNameWaitingLockTrxId string = "waiting_lock_trx_id"

var ALLMetric = map[string]MetricType{
	MetricNameCounter:                   MetricTypeInt,    // MySQL slow log
	MetricNameLastReceiveTimestamp:      MetricTypeString, // MySQL slow log
	MetricNameQueryTimeAvg:              MetricTypeFloat,  // MySQL slow log
	MetricNameQueryTimeMax:              MetricTypeFloat,  // MySQL slow log
	MetricNameRowExaminedAvg:            MetricTypeFloat,  // MySQL slow log
	MetricNameFirstQueryAt:              MetricTypeString, // MySQL slow log, 好像没用上 | OB MySQL TOP SQL
	MetricNameDBUser:                    MetricTypeString, // MySQL slow log
	MetricNameEndpoints:                 MetricTypeArray,  // MySQL slow log
	MetricNameStartTimeOfLastScrapedSQL: MetricTypeString, // MySQL slow log
	MetricNameMetaName:                  MetricTypeString, // MySQL schema meta
	MetricNameMetaType:                  MetricTypeString, // MySQL schema meta
	MetricNameRecordDeleted:             MetricTypeBool,   // MySQL schema meta

	MetricNameQueryTimeTotal:      MetricTypeFloat, // DB2 TOP SQL | OB Oracle TOP SQL
	MetricNameCPUTimeAvg:          MetricTypeFloat, // DB2 TOP SQL | OB MySQL TOP SQL
	MetricNameLockWaitTimeTotal:   MetricTypeFloat, // DB2 TOP SQL
	MetricNameLockWaitCounter:     MetricTypeInt,   // DB2 TOP SQL
	MetricNameActiveWaitTimeTotal: MetricTypeFloat, // DB2 TOP SQL
	MetricNameActiveTimeTotal:     MetricTypeFloat, // DB2 TOP SQL
	MetricNameCPUTimeTotal:        MetricTypeFloat, // DM TOP SQL  | OB Oracle TOP SQL
	MetricNamePhyReadPageTotal:    MetricTypeInt,   // DM TOP SQL | OB Oracle TOP SQL
	MetricNameLogicReadPageTotal:  MetricTypeInt,   // DM TOP SQL | OB Oracle TOP SQL

	MetricNameUserIOWaitTimeTotal: MetricTypeFloat, // OB Oracle TOP SQL
	MetricNameBufferGetCounter:    MetricTypeInt,   // OB Oracle TOP SQL
	MetricNameDiskReadTotal:       MetricTypeInt,   // OB Oracle TOP SQL

	MetricNameLastQueryAt:   MetricTypeString, // OB MySQL TOP SQL
	MetricNameIoWaitTimeAvg: MetricTypeFloat,  // OB MySQL TOP SQL

	MetricNameLockType:   MetricTypeString, // Lock
	MetricNameLockMode:   MetricTypeString, // Lock
	MetricNameLockStatus: MetricTypeString, // Lock
	MetricNameEngine:     MetricTypeString, // Lock
	MetricNameTable:      MetricTypeString, // Lock
}

func LoadMetrics(info map[string]interface{}, metrics []string) Metrics {
	ms := NewMetrics()
	for _, metric := range metrics {
		typ, ok := ALLMetric[metric]
		// 指标不存在，说明task上定义的指标并没有预先申明
		if !ok {
			continue
		}
		// 判断info内是否有该指标，说明定义的指标在info里未存储
		v, ok := info[metric]
		if !ok {
			continue
		}
		// todo: 需要对类型断言进行判断
		switch typ {
		case MetricTypeInt:
			loadInt(ms, metric, v)

		case MetricTypeFloat:
			loadFloat(ms, metric, v)

		case MetricTypeString:
			if s, ok := v.(string); ok {
				ms.SetString(metric, s)
			} else {
				ms.SetString(metric, "")
			}
		case MetricTypeArray:
			if ss, ok := v.([]string); ok {
				ms.SetStringArray(metric, ss)
			} else if ss, ok := v.([]interface{}); ok {
				var valList []string
				for _, s := range ss {
					if val, ok := s.(string); ok {
						valList = append(valList, val)
					}
				}
				ms.SetStringArray(metric, valList)
			} else {
				ms.SetStringArray(metric, nil)
			}
		case MetricTypeTime:
			if t, ok := v.(*time.Time); ok {
				ms.SetTime(metric, t)
			} else {
				ms.SetTime(metric, nil)
			}

		case MetricTypeBool:
			if b, ok := v.(bool); ok {
				ms.SetBool(metric, b)
			} else {
				ms.SetBool(metric, false)
			}
		default:
		}
	}
	return ms
}

func loadInt(ms Metrics, name string, v interface{}) {
	// 如果值为float则强制转成 int
	switch tv := v.(type) {
	case int:
		ms.SetInt(name, int64(tv))
	case uint64:
		ms.SetInt(name, int64(tv))
	case int64:
		ms.SetInt(name, tv)
	case int32:
		ms.SetInt(name, int64(tv))
	case float64:
		ms.SetInt(name, int64(tv))
	}
}

func loadFloat(ms Metrics, name string, v interface{}) {
	// 如果值为整形，则强制转float
	switch tv := v.(type) {
	case int:
		ms.SetFloat(name, float64(tv))
	case int64:
		ms.SetFloat(name, float64(tv))
	case int32:
		ms.SetFloat(name, float64(tv))
	case float64:
		ms.SetFloat(name, tv)
	}
}

type MetricType int

const (
	MetricTypeInt    MetricType = 1
	MetricTypeFloat  MetricType = 2
	MetricTypeString MetricType = 3
	MetricTypeTime   MetricType = 4
	MetricTypeBool   MetricType = 5
	MetricTypeArray  MetricType = 6
)

type Metric struct {
	name string
	i    int64
	s    string
	f    float64
	t    *time.Time
	b    bool
	ss   []string
	typ  MetricType
}

func (m *Metric) Int() int64 {
	if m == nil || m.typ != MetricTypeInt {
		return 0
	}
	return m.i
}

func (m *Metric) Float() float64 {
	if m == nil || m.typ != MetricTypeFloat {
		return 0
	}
	return m.f
}

func (m *Metric) String() string {
	if m == nil || m.typ != MetricTypeString {
		return ""
	}
	return m.s
}

func (m *Metric) StringArray() []string {
	if m == nil || m.typ != MetricTypeArray {
		return nil
	}
	return m.ss
}

func (m *Metric) Time() *time.Time {
	if m == nil || m.typ != MetricTypeTime {
		return nil
	}
	return m.t
}

func (m *Metric) Bool() bool {
	if m == nil || m.typ != MetricTypeBool {
		return false
	}
	return m.b
}

func (m *Metric) Value() interface{} {
	if m == nil {
		return nil
	}
	switch m.typ {
	case MetricTypeInt:
		return m.i
	case MetricTypeFloat:
		return m.f
	case MetricTypeString:
		return m.s
	case MetricTypeArray:
		return m.ss
	case MetricTypeTime:
		return m.t
	case MetricTypeBool:
		return m.b
	default:
		return nil
	}
}

func NewMetric() *Metric {
	return &Metric{}
}

type Metrics map[string]*Metric

func NewMetrics() Metrics {
	return Metrics{}
}

func (m Metrics) SetInt(name string, i int64) {
	m[name] = &Metric{
		typ: MetricTypeInt,
		i:   i,
	}
}

func (m Metrics) SetFloat(name string, f float64) {
	m[name] = &Metric{
		typ: MetricTypeFloat,
		f:   f,
	}
}

func (m Metrics) SetString(name string, s string) {
	m[name] = &Metric{
		typ: MetricTypeString,
		s:   s,
	}
}

func (m Metrics) SetStringArray(name string, ss []string) {
	m[name] = &Metric{
		typ: MetricTypeArray,
		ss:  ss,
	}
}

func (m Metrics) SetTime(name string, t *time.Time) {
	m[name] = &Metric{
		typ: MetricTypeTime,
		t:   t,
	}
}

func (m Metrics) SetBool(name string, b bool) {
	m[name] = &Metric{
		typ: MetricTypeBool,
		b:   b,
	}
}

func (m Metrics) Get(name string) *Metric {
	return m[name]
}

func (m Metrics) ToMap() map[string]interface{} {
	data := map[string]interface{}{}
	for k, v := range m {
		data[k] = v.Value()
	}
	return data
}
