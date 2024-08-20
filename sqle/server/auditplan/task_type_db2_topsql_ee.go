//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
	dry "github.com/ungerik/go-dry"
)

type DB2TopSQLTaskV2 struct{}

func NewDB2TopSQLTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &DB2TopSQLTaskV2{}
	}
}

func (at *DB2TopSQLTaskV2) InstanceType() string {
	return InstanceTypeDB2
}

func (at *DB2TopSQLTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:   paramKeyCollectIntervalMinute,
			Desc:  "采集周期（分钟）",
			Value: "60",
			Type:  params.ParamTypeInt,
		},
		{
			Key:   paramKeyTopN,
			Desc:  "Top N",
			Value: "3",
			Type:  params.ParamTypeInt,
		},
		{
			Key:   paramKeyIndicator,
			Desc:  "关注指标",
			Value: DB2IndicatorAverageElapsedTime,
			Type:  params.ParamTypeString,
		},
	}
}

func (at *DB2TopSQLTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{}
}

func (at *DB2TopSQLTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameQueryTimeTotal,
		MetricNameQueryTimeAvg,
		MetricNameCPUTimeAvg,
		MetricNameLockWaitTimeTotal,
		MetricNameLockWaitCounter,
		MetricNameActiveWaitTimeTotal,
		MetricNameActiveTimeTotal,
		MetricNameLastReceiveTimestamp,
	}
}

func (at *DB2TopSQLTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}

	// counter
	originSQL.Info.SetInt(MetricNameCounter, mergedSQL.Info.Get(MetricNameCounter).Int())

	// MetricNameQueryTimeTotal
	originSQL.Info.SetFloat(MetricNameQueryTimeTotal, mergedSQL.Info.Get(MetricNameQueryTimeTotal).Float())

	// MetricNameQueryTimeAvg
	originSQL.Info.SetFloat(MetricNameQueryTimeAvg, mergedSQL.Info.Get(MetricNameQueryTimeAvg).Float())

	// MetricNameCPUTimeAvg
	originSQL.Info.SetFloat(MetricNameCPUTimeAvg, mergedSQL.Info.Get(MetricNameCPUTimeAvg).Float())

	// MetricNameLockWaitTimeTotal
	originSQL.Info.SetFloat(MetricNameLockWaitTimeTotal, mergedSQL.Info.Get(MetricNameLockWaitTimeTotal).Float())

	// MetricNameLockWaitCounter
	originSQL.Info.SetInt(MetricNameLockWaitCounter, mergedSQL.Info.Get(MetricNameLockWaitCounter).Int())

	// MetricNameActiveWaitTimeTotal
	originSQL.Info.SetInt(MetricNameActiveWaitTimeTotal, mergedSQL.Info.Get(MetricNameActiveWaitTimeTotal).Int())

	//MetricNameActiveTimeTotal
	originSQL.Info.SetFloat(MetricNameActiveTimeTotal, mergedSQL.Info.Get(MetricNameActiveTimeTotal).Float())

	// last_receive_timestamp
	originSQL.Info.SetString(MetricNameLastReceiveTimestamp, mergedSQL.Info.Get(MetricNameLastReceiveTimestamp).String())
	return
}

const (
	DB2IndicatorNumExecutions      = "num_executions"      // 执行次数
	DB2IndicatorTotalElapsedTime   = "stmt_exec_time"      // 总执行时间
	DB2IndicatorAverageElapsedTime = "avg_elapsed_time_ms" // 平均执行时间
	DB2IndicatorAverageCPUTime     = "avg_cpu_time"        // 平均 CPU 时间
	DB2IndicatorLockWaitTime       = "lock_wait_time"      // 锁等待时间
	DB2IndicatorLockWaitNum        = "lock_waits"          // 锁等待次数
	DB2IndicatorSQLWaitTime        = "total_act_wait_time" // 活动等待总时间
	DB2IndicatorTotalActTime       = "total_act_time"      // 活动总时间
)

func (at *DB2TopSQLTaskV2) indicator(ap *AuditPlan) (string, error) {
	indicator := ap.Params.GetParam(paramKeyIndicator).String()
	if indicator == "" {
		return DB2IndicatorAverageElapsedTime, nil
	}

	if !dry.StringInSlice(indicator, []string{
		DB2IndicatorNumExecutions,
		DB2IndicatorTotalElapsedTime,
		DB2IndicatorAverageElapsedTime,
		DB2IndicatorAverageCPUTime,
		DB2IndicatorLockWaitTime,
		DB2IndicatorLockWaitNum,
		DB2IndicatorSQLWaitTime,
		DB2IndicatorTotalActTime,
	}) {
		return "", fmt.Errorf("invalid indicator: %v", indicator)
	}
	return indicator, nil
}

// ref: https://www.ibm.com/docs/zh/db2/11.1?topic=views-snap-get-dyn-sql-dynsql-snapshot
func (at *DB2TopSQLTaskV2) collectSQL(ap *AuditPlan) (string, error) {
	// SET TENANT ?是用于设置当前会话的租户标识的语句
	// MON_GET_PKG_CACHE_STMT表函数中可能会存在两次`SET TENANT ?`。因为指纹唯一性导致存表失败，所以过滤对`SET TENANT ?`语句的采集
	sql := `
	SELECT
	num_executions,
	total_act_time,
	total_act_wait_time,
	lock_wait_time,
	lock_waits,
	stmt_exec_time,
	total_cpu_time / NUM_EXEC_WITH_METRICS AS avg_cpu_time,
	STMT_EXEC_TIME/NUM_EXECUTIONS AS avg_elapsed_time_ms,
	substr(stmt_text, 1, 300) AS stmt_text
	FROM
	TABLE(MON_GET_PKG_CACHE_STMT(NULL, NULL, NULL, -2)) T
	WHERE
	upper(stmt_text) NOT LIKE '%%MON_GET_PKG_CACHE_STMT%%' AND NUM_EXEC_WITH_METRICS <> 0 AND upper(stmt_text) <> 'SET TENANT ?'
ORDER BY %s DESC   
`
	indicator, err := at.indicator(ap)
	if err != nil {
		return "", err
	}

	sql = fmt.Sprintf(sql, indicator)

	// limit top N
	{
		topN := ap.Params.GetParam(paramKeyTopN).Int()
		if topN == 0 {
			topN = 10
		}
		sql = fmt.Sprintf(`%v FETCH FIRST %d ROWS ONLY `, sql, topN)
	}

	return sql, nil
}

func (at *DB2TopSQLTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	inst, _, err := dms.GetInstancesById(context.Background(), ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}

	if !driver.GetPluginManager().IsOptionalModuleEnabled(inst.DbType, driverV2.OptionalModuleQuery) {
		return nil, fmt.Errorf("can not do this task, %v", driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleQuery))
	}

	plugin, err := driver.GetPluginManager().OpenPlugin(
		logger, inst.DbType, &driverV2.Config{
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

	sql, err := at.collectSQL(ap)
	if err != nil {
		return nil, fmt.Errorf("generate collect sql failed, error: %v", err)
	}

	result, err := plugin.Query(context.Background(), sql,
		&driverV2.QueryConf{TimeOutSecond: 10})
	if err != nil {
		return nil, fmt.Errorf("collect failed, error: %v", err)
	}

	if len(result.Rows) == 0 {
		logger.Infof("db2 top sql audit_plan(%v) collected no statement", ap.ID)
		return nil, nil
	}

	logger.Infof("db2 top sql audit_plan(%v) collected %v statements", ap.ID, len(result.Rows))

	cache := NewSQLV2Cache()

	for i := range result.Rows {
		row := result.Rows[i]
		info := NewMetrics()
		sqlV2 := &SQLV2{
			Source:     ap.Type,
			SourceId:   ap.ID,
			ProjectId:  ap.ProjectId,
			InstanceID: ap.InstanceID,
			SchemaName: "", // todo: top sql 未采集schema, 需要填充
			Info:       info,
		}

		info.SetString(MetricNameLastReceiveTimestamp, time.Now().Format(time.RFC3339)) // todo: 没啥大用
		for j := range row.Values {
			switch result.Column[j].Key {
			case "stmt_test":
				sqlV2.SQLContent = row.Values[j].Value
				sqlV2.Fingerprint = row.Values[j].Value

			case DB2IndicatorNumExecutions:
				var counter int64
				counter, err = strconv.ParseInt(row.Values[j].Value, 10, 64)
				if err != nil {
					return nil, err
				}
				info.SetInt(MetricNameCounter, counter)

			case DB2IndicatorTotalElapsedTime:
				var queryTimeTotal float64
				queryTimeTotal, err = strconv.ParseFloat(row.Values[j].Value, 10)
				if err != nil {
					return nil, err
				}
				info.SetFloat(MetricNameQueryTimeTotal, queryTimeTotal)

			case DB2IndicatorAverageElapsedTime:
				var queryTimeAvg float64
				queryTimeAvg, err = strconv.ParseFloat(row.Values[j].Value, 10)
				if err != nil {
					return nil, err
				}
				info.SetFloat(MetricNameQueryTimeAvg, queryTimeAvg)

			case DB2IndicatorAverageCPUTime:
				var cpuTimeAvg float64
				cpuTimeAvg, err = strconv.ParseFloat(row.Values[j].Value, 10)
				if err != nil {
					return nil, err
				}
				info.SetFloat(MetricNameCPUTimeAvg, cpuTimeAvg)
			case DB2IndicatorLockWaitTime:
				var lockWaitTimeTotal float64
				lockWaitTimeTotal, err = strconv.ParseFloat(row.Values[j].Value, 10)
				if err != nil {
					return nil, err
				}
				info.SetFloat(MetricNameLockWaitTimeTotal, lockWaitTimeTotal)

			case DB2IndicatorLockWaitNum:
				var lockWaitCounter int64
				lockWaitCounter, err = strconv.ParseInt(row.Values[j].Value, 10, 64)
				if err != nil {
					return nil, err
				}
				info.SetInt(MetricNameLockWaitCounter, lockWaitCounter)
			case DB2IndicatorSQLWaitTime:
				var actWaitTimeTotal float64
				actWaitTimeTotal, err = strconv.ParseFloat(row.Values[j].Value, 10)
				if err != nil {
					return nil, err
				}
				info.SetFloat(MetricNameActiveWaitTimeTotal, actWaitTimeTotal)
			case DB2IndicatorTotalActTime:
				var actTimeTotal float64
				actTimeTotal, err = strconv.ParseFloat(row.Values[j].Value, 10)
				if err != nil {
					return nil, err
				}
				info.SetFloat(MetricNameActiveTimeTotal, actTimeTotal)
			}
		}
		sqlV2.GenSQLId()
		at.AggregateSQL(cache, sqlV2)
	}
	return cache.GetSQLs(), nil
}

func (at *DB2TopSQLTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *DB2TopSQLTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *DB2TopSQLTaskV2) Head(ap *AuditPlan) []Head {
	return []Head{
		{
			Name: "sql",
			Desc: "SQL语句",
			Type: "sql",
		},
		{
			Name: model.AuditResultName,
			Desc: model.AuditResultDesc,
		},
		{
			Name: "priority",
			Desc: "优先级",
		},
		{
			Name: MetricNameQueryTimeTotal,
			Desc: "总执行时间(ms)",
		},
		{
			Name: MetricNameQueryTimeAvg,
			Desc: "平均执行时间(ms)",
		},
		{
			Name: MetricNameCounter,
			Desc: "执行次数",
		},
		{
			Name: MetricNameCPUTimeAvg,
			Desc: "平均 CPU 时间(μs)",
		},
		{
			Name: MetricNameLockWaitTimeTotal,
			Desc: "锁等待时间(ms)",
		},
		{
			Name: MetricNameLockWaitCounter,
			Desc: "锁等待次数",
		},
		{
			Name: MetricNameActiveWaitTimeTotal,
			Desc: "活动等待总时间(ms)",
		},
		{
			Name: MetricNameActiveTimeTotal,
			Desc: "活动总时间(ms)",
		},
	}
}

func (at *DB2TopSQLTaskV2) Filters(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta {
	return []FilterMeta{
		{
			Name:            "sql", // 模糊筛选
			Desc:            "SQL",
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
		},
		{
			Name:            "rule_name",
			Desc:            "审核规则",
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerRuleTips(logger, ap.ID, persist),
		},
		{
			Name:            "priority",
			Desc:            "SQL优先级",
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerPriorityTips(logger),
		}}
}

func (at *DB2TopSQLTaskV2) GetSQLData(ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
	args := make(map[model.FilterName]interface{}, len(filters))
	for _, filter := range filters {
		switch filter.Name {
		case "sql":
			args[model.FilterSQL] = filter.FilterComparisonValue

		case "priority":
			args[model.FilterPriority] = filter.FilterComparisonValue

		case "rule_name":
			args[model.FilterRuleName] = filter.FilterComparisonValue
		}
	}
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReqV2(ap.ID, ap.Type, limit, offset, checkAndGetOrderByName(at.Head(ap), orderBy), isAsc, args)
	if err != nil {
		return nil, count, err
	}

	result := []map[string]string{}

	for i := range auditPlanSQLs {
		mp := map[string]string{
			"sql":                 auditPlanSQLs[i].SQLContent,
			"id":                  auditPlanSQLs[i].AuditPlanSqlId,
			"priority":            auditPlanSQLs[i].Priority.String,
			model.AuditResultName: auditPlanSQLs[i].AuditResult.String,
		}

		origin, err := auditPlanSQLs[i].Info.OriginValue()
		if err != nil {
			return nil, 0, err
		}
		for k := range origin {
			mp[k] = fmt.Sprintf("%v", origin[k])
		}
		result = append(result, mp)
	}
	return result, count, nil
}
