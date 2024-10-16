//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"fmt"
	"strconv"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

type ObForMysqlTopSQLTaskV2 struct {
	DefaultTaskV2
}

func NewObForMysqlTopSQLTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &ObForMysqlTopSQLTaskV2{DefaultTaskV2: DefaultTaskV2{}}
	}
}

func (at *ObForMysqlTopSQLTaskV2) InstanceType() string {
	return InstanceTypeOceanBaseForMySQL
}

func (at *ObForMysqlTopSQLTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:      paramKeyCollectIntervalMinute,
			Value:    "60",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamCollectIntervalMinute),
		},
		{
			Key:   paramKeyTopN,
			Desc:  "Top N",
			Value: "3",
			Type:  params.ParamTypeInt,
		},
		{
			Key:      paramKeyIndicator,
			Value:    DB2IndicatorAverageElapsedTime,
			Type:     params.ParamTypeString,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamIndicator),
		},
	}
}

func (at *ObForMysqlTopSQLTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameFirstQueryAt,
		MetricNameLastQueryAt,
		MetricNameQueryTimeAvg,
		MetricNameQueryTimeMax,
		MetricNameIoWaitTimeAvg,
		MetricNameCPUTimeAvg,
		MetricNameDiskReadAvg,
		MetricNameBufferReadAvg,
	}
}

func (at *ObForMysqlTopSQLTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}
	// counter
	originSQL.Info.SetInt(MetricNameCounter, mergedSQL.Info.Get(MetricNameCounter).Int())

	// MetricNameFirstQueryAt
	originSQL.Info.SetString(MetricNameFirstQueryAt, mergedSQL.Info.Get(MetricNameFirstQueryAt).String())

	// MetricNameLastQueryAt
	originSQL.Info.SetString(MetricNameLastQueryAt, mergedSQL.Info.Get(MetricNameLastQueryAt).String())

	// MetricNameQueryTimeAvg
	originSQL.Info.SetFloat(MetricNameQueryTimeAvg, mergedSQL.Info.Get(MetricNameQueryTimeAvg).Float())

	// MetricNameQueryTimeMax
	originSQL.Info.SetFloat(MetricNameQueryTimeMax, mergedSQL.Info.Get(MetricNameQueryTimeMax).Float())

	// MetricNameIoWaitTimeAvg
	originSQL.Info.SetFloat(MetricNameIoWaitTimeAvg, mergedSQL.Info.Get(MetricNameIoWaitTimeAvg).Float())

	// MetricNameCPUTimeAvg
	originSQL.Info.SetFloat(MetricNameCPUTimeAvg, mergedSQL.Info.Get(MetricNameCPUTimeAvg).Float())

	// MetricNameDiskReadAvg
	originSQL.Info.SetFloat(MetricNameDiskReadAvg, mergedSQL.Info.Get(MetricNameDiskReadAvg).Float())

	// MetricNameBufferReadAvg
	originSQL.Info.SetFloat(MetricNameBufferReadAvg, mergedSQL.Info.Get(MetricNameBufferReadAvg).Float())
	return
}

const (
	// 通用采集项
	OBMySQLSQLKeySQLText            = "sql_text"
	OBMySQLSQLInfoKeyFirstRequest   = "first_request"
	OBMySQLSQLInfoKeyExecCount      = "exec_count"
	OBMySQLSQLInfoKeyLastRequest    = "last_request"
	OBMySQLSQLInfoKeyAverageElapsed = "average_elapsed"

	// OBMySQLIndicatorElapsedTime 对应采集项
	OBMySQLSQLInfoKeyMaxElapsed = "max_elapsed"

	// OBMySQLIndicatorCPUTime 对应采集项
	OBMySQLSQLInfoKeyAverageCPU = "average_cpu"

	// OBMySQLIndicatorIOWait 对应采集项
	OBMySQLSQLInfoKeyAverageIOWait = "average_io_wait"
	OBMySQLSQLInfoKeyDiskRead      = "disk_read"
	OBMySQLSQLInfoKeyBufferRead    = "buffer_read"
)

func (at *ObForMysqlTopSQLTaskV2) getCollectSQL(ap *AuditPlan) string {
	topN := ap.Params.GetParam(paramKeyTopN).Int()

	switch ap.Params.GetParam(paramKeyIndicator).String() {
	case OBMySQLIndicatorElapsedTime:
		return fmt.Sprintf(`
SELECT
    SQL_TEXT AS %v, 
    EXECUTIONS AS %v, 
    CEIL(AVG_EXE_USEC/1000) AS %v, 
    CEIL(SLOWEST_EXE_USEC/1000) AS %v, 
    FROM_UNIXTIME(TIME_TO_USEC(FIRST_LOAD_TIME)/1000000) AS %v,
    FROM_UNIXTIME(TIME_TO_USEC(LAST_ACTIVE_TIME)/1000000) AS %v
FROM
    OCEANBASE.GV$SQL
WHERE %v != ''
GROUP BY
    SQL_ID
ORDER BY
    %v
DESC
LIMIT %v
`, OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyExecCount,
			OBMySQLSQLInfoKeyAverageElapsed,
			OBMySQLSQLInfoKeyMaxElapsed,
			OBMySQLSQLInfoKeyFirstRequest,
			OBMySQLSQLInfoKeyLastRequest,
			OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyMaxElapsed,
			topN)

	case OBMySQLIndicatorCPUTime:
		return fmt.Sprintf(`
SELECT
    SQL_TEXT AS %v, 
    EXECUTIONS AS %v, 
    CEIL(AVG_EXE_USEC/1000) AS %v,
    CEIL(CPU_TIME/EXECUTIONS/1000) AS %v, 
    FROM_UNIXTIME(TIME_TO_USEC(FIRST_LOAD_TIME)/1000000) AS %v,
    FROM_UNIXTIME(TIME_TO_USEC(LAST_ACTIVE_TIME)/1000000) AS %v
FROM
    OCEANBASE.GV$SQL
WHERE %v != ''
GROUP BY
    SQL_ID
ORDER BY
    %v
DESC
LIMIT %v
`, OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyExecCount,
			OBMySQLSQLInfoKeyAverageElapsed,
			OBMySQLSQLInfoKeyAverageCPU,
			OBMySQLSQLInfoKeyFirstRequest,
			OBMySQLSQLInfoKeyLastRequest,
			OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyAverageCPU,
			topN,
		)

	case OBMySQLIndicatorIOWait:
		return fmt.Sprintf(`
SELECT
    SQL_TEXT AS %v, 
    EXECUTIONS AS %v, 
    CEIL(USER_IO_WAIT_TIME/EXECUTIONS/1000) AS %v, 
    CEIL(BUFFER_GETS/EXECUTIONS) AS %v,
    CEIL(DISK_READS/EXECUTIONS) AS %v,
    FROM_UNIXTIME(TIME_TO_USEC(FIRST_LOAD_TIME)/1000000) AS %v,
    FROM_UNIXTIME(TIME_TO_USEC(LAST_ACTIVE_TIME)/1000000) AS %v
FROM
    OCEANBASE.GV$SQL
WHERE %v != ''
GROUP BY
    SQL_ID
ORDER BY
    %v
DESC
LIMIT %v
`, OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyExecCount,
			OBMySQLSQLInfoKeyAverageIOWait,
			OBMySQLSQLInfoKeyBufferRead,
			OBMySQLSQLInfoKeyDiskRead,
			OBMySQLSQLInfoKeyFirstRequest,
			OBMySQLSQLInfoKeyLastRequest,
			OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyAverageIOWait,
			topN,
		)

	default:
		return ""
	}
}

type DynPerformanceObForMySQLColumns struct {
	SQLText       string  `json:"sql_text"`
	Executions    float64 `json:"exec_count"`
	FirstQueryAt  string  `json:"first_request"`
	LastQueryAt   string  `json:"last_request"`
	QueryTimeAvg  float64 `json:"average_elapsed"`
	QueryTimeMax  float64 `json:"max_elapsed"`
	IoWaitTimeAvg float64 `json:"average_io_wait"`
	BufferReadAvg float64 `json:"buffer_read"`
	DiskReadAvg   float64 `json:"disk_read"`
	CPUTimeAvg    float64 `json:"average_cpu"`
}

func (at *ObForMysqlTopSQLTaskV2) collect(ap *AuditPlan, persist *model.Storage, p driver.Plugin, sql string) ([]*DynPerformanceObForMySQLColumns, error) {
	result, err := p.Query(context.Background(), sql, &driverV2.QueryConf{TimeOutSecond: 5})
	if err != nil {
		return nil, err
	}
	if len(result.Column) <= 0 {
		return nil, nil
	}
	sqls := []*DynPerformanceObForMySQLColumns{}

	for _, row := range result.Rows {
		s := &DynPerformanceObForMySQLColumns{}
		for i, value := range row.Values {
			switch result.Column[i].Key {
			case OBMySQLSQLKeySQLText:
				s.SQLText = value.Value
			case OBMySQLSQLInfoKeyFirstRequest:
				s.FirstQueryAt = value.Value
			case OBMySQLSQLInfoKeyLastRequest:
				s.LastQueryAt = value.Value
			case OBMySQLSQLInfoKeyExecCount:
				executions, err := strconv.ParseInt(value.Value, 10, 64)
				if err != nil {
					return nil, err
				}
				s.Executions = float64(executions)
			case OBMySQLSQLInfoKeyAverageElapsed:
				elapsedTime, err := strconv.ParseFloat(value.Value, 64)
				if err != nil {
					return nil, err
				}
				s.QueryTimeAvg = elapsedTime
			case OBMySQLSQLInfoKeyMaxElapsed:
				elapsedTimeMax, err := strconv.ParseFloat(value.Value, 64)
				if err != nil {
					return nil, err
				}
				s.QueryTimeMax = elapsedTimeMax
			case OBMySQLSQLInfoKeyAverageCPU:
				cpuTimeAvg, err := strconv.ParseFloat(value.Value, 64)
				if err != nil {
					return nil, err
				}
				s.CPUTimeAvg = cpuTimeAvg
			case OBMySQLSQLInfoKeyAverageIOWait:
				ioWaitTime, err := strconv.ParseFloat(value.Value, 64)
				if err != nil {
					return nil, err
				}
				s.IoWaitTimeAvg = ioWaitTime
			case OBMySQLSQLInfoKeyDiskRead:
				diskReadAvg, err := strconv.ParseFloat(value.Value, 64)
				if err != nil {
					return nil, err
				}
				s.DiskReadAvg = diskReadAvg
			case OBMySQLSQLInfoKeyBufferRead:
				bufferReadAvg, err := strconv.ParseFloat(value.Value, 64)
				if err != nil {
					return nil, err
				}
				s.BufferReadAvg = bufferReadAvg
			}

		}
		sqls = append(sqls, s)
	}
	return sqls, nil
}

func (at *ObForMysqlTopSQLTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	inst, exist, err := dms.GetInstancesById(context.Background(), ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}
	if !exist {
		return nil, fmt.Errorf("instance: %v is not exist", ap.InstanceID)
	}

	if !driver.GetPluginManager().IsOptionalModuleEnabled(inst.DbType, driverV2.OptionalModuleQuery) {
		return nil, fmt.Errorf("can not do this task, %v", driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleQuery))
	}

	plugin, err := driver.GetPluginManager().OpenPlugin(logger, inst.DbType, &driverV2.Config{
		DSN: &driverV2.DSN{
			Host:             inst.Host,
			Port:             inst.Port,
			User:             inst.User,
			Password:         inst.Password,
			AdditionalParams: inst.AdditionalParams,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get plugin failed, error: %v", err)
	}
	defer plugin.Close(context.Background())

	sql := at.getCollectSQL(ap)
	if sql == "" {
		return nil, fmt.Errorf("unknown metric of interest")
	}
	sqls, err := at.collect(ap, persist, plugin, sql)
	if err != nil {
		return nil, fmt.Errorf("collect failed, error: %v", err)
	}

	cache := NewSQLV2Cache()
	for _, sql := range sqls {
		info := NewMetrics()
		sqlV2 := &SQLV2{
			Source:      ap.Type,
			SourceId:    strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
			ProjectId:   ap.ProjectId,
			InstanceID:  ap.InstanceID,
			SchemaName:  "", // todo: top sql 未采集schema, 需要填充
			Info:        info,
			SQLContent:  sql.SQLText,
			Fingerprint: sql.SQLText,
		}
		info.SetInt(MetricNameCounter, int64(sql.Executions))
		info.SetString(MetricNameFirstQueryAt, sql.FirstQueryAt)
		info.SetString(MetricNameLastQueryAt, sql.LastQueryAt)
		info.SetFloat(MetricNameQueryTimeAvg, sql.QueryTimeAvg)
		info.SetFloat(MetricNameQueryTimeMax, sql.QueryTimeMax)
		info.SetFloat(MetricNameIoWaitTimeAvg, sql.IoWaitTimeAvg)

		info.SetFloat(MetricNameCPUTimeAvg, sql.CPUTimeAvg)
		info.SetInt(MetricNameDiskReadAvg, int64(sql.DiskReadAvg))
		info.SetInt(MetricNameBufferReadAvg, int64(sql.BufferReadAvg))

		sqlV2.GenSQLId()
		at.AggregateSQL(cache, sqlV2)
	}
	return cache.GetSQLs(), nil
}

func (at *ObForMysqlTopSQLTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
	originSQL, exist, err := cache.GetSQL(sql.SQLId)
	if err != nil {
		return err
	}
	if !exist {
		cache.CacheSQL(sql)
		return nil
	}
	at.mergeSQL(originSQL, sql)
	return nil
}

func (at *ObForMysqlTopSQLTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *ObForMysqlTopSQLTaskV2) Head(ap *AuditPlan) []Head {
	switch ap.Params.GetParam(paramKeyIndicator).String() {
	case OBMySQLIndicatorElapsedTime:
		return []Head{
			{
				Name: "sql",
				Desc: locale.ApSQLFingerprint,
				Type: "sql",
			},
			{
				Name: "priority",
				Desc: locale.ApPriority,
			}, {
				Name: model.AuditResultName,
				Desc: model.AuditResultDesc,
			}, {
				Name: MetricNameCounter,
				Desc: locale.ApMetricNameCounter,
			}, {
				Name: MetricNameQueryTimeAvg,
				Desc: locale.ApMetricNameQueryTimeAvg,
			}, {
				Name: MetricNameQueryTimeMax,
				Desc: locale.ApMetricNameQueryTimeMax,
			}, {
				Name: MetricNameFirstQueryAt,
				Desc: locale.ApMetricNameFirstQueryAt,
			}, {
				Name: MetricNameLastQueryAt,
				Desc: locale.ApMetricNameLastQueryAt,
			},
		}
	case OBMySQLIndicatorIOWait:
		return []Head{
			{
				Name: "sql",
				Desc: locale.ApSQLFingerprint,
				Type: "sql",
			}, {
				Name: model.AuditResultName,
				Desc: model.AuditResultDesc,
			}, {
				Name: MetricNameCounter,
				Desc: locale.ApMetricNameCounter,
			}, {
				Name: MetricNameIoWaitTimeAvg,
				Desc: locale.ApMetricNameIoWaitTimeAvg,
			}, {
				Name: MetricNameBufferReadAvg,
				Desc: locale.ApMetricNameBufferReadAvg,
			}, {
				Name: MetricNameDiskReadAvg,
				Desc: locale.ApMetricNameDiskReadAvg,
			}, {
				Name: MetricNameFirstQueryAt,
				Desc: locale.ApMetricNameFirstQueryAt,
			}, {
				Name: MetricNameLastQueryAt,
				Desc: locale.ApMetricNameLastQueryAt,
			},
		}
	case OBMySQLIndicatorCPUTime:
		return []Head{
			{
				Name: "sql",
				Desc: locale.ApSQLFingerprint,
				Type: "sql",
			}, {
				Name: model.AuditResultName,
				Desc: model.AuditResultDesc,
			}, {
				Name: MetricNameCounter,
				Desc: locale.ApMetricNameCounter,
			}, {
				Name: MetricNameCPUTimeAvg,
				Desc: locale.ApMetricNameCPUTimeAvg,
			}, {
				Name: MetricNameQueryTimeAvg,
				Desc: locale.ApMetricNameQueryTimeAvg,
			}, {
				Name: MetricNameFirstQueryAt,
				Desc: locale.ApMetricNameFirstQueryAt,
			}, {
				Name: MetricNameLastQueryAt,
				Desc: locale.ApMetricNameLastQueryAt,
			},
		}
	}
	return []Head{}
}

func (at *ObForMysqlTopSQLTaskV2) GetSQLData(ctx context.Context, ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReqV2(ap.ID, ap.Type, limit, offset, checkAndGetOrderByName(at.Head(ap), orderBy), isAsc, genArgsByFilters(filters))
	if err != nil {
		return nil, count, err
	}
	result := []map[string]string{}
	for _, planSQL := range auditPlanSQLs {
		mp := map[string]string{
			"sql":                 planSQL.SQLContent,
			"id":                  planSQL.AuditPlanSqlId,
			"priority":            planSQL.Priority.String,
			model.AuditResultName: planSQL.AuditResult.GetAuditJsonStrByLangTag(locale.Bundle.GetLangTagFromCtx(ctx)),
		}

		origin, err := planSQL.Info.OriginValue()
		if err != nil {
			return nil, 0, err
		}
		for k, v := range origin {
			mp[k] = fmt.Sprintf("%v", v)
		}
		result = append(result, mp)
	}
	return result, count, nil
}
