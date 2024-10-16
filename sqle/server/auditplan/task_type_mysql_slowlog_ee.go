//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/sirupsen/logrus"
)

type SlowLogTaskV2 struct{}

func NewSlowLogTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &SlowLogTaskV2{}
	}
}

func (at *SlowLogTaskV2) InstanceType() string {
	return InstanceTypeMySQL
}

func (at *SlowLogTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:      paramKeyCollectIntervalMinute,
			Value:    "60",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamCollectIntervalMinuteMySQL),
		},
		{
			Key:   paramKeySlowLogCollectInput,
			Value: "0",
			Type:  params.ParamTypeInt,
			Enums: []params.EnumsValue{
				{Value: "0", I18nDesc: locale.Bundle.LocalizeAll(locale.EnumSlowLogFileSource)},
				{Value: "1", I18nDesc: locale.Bundle.LocalizeAll(locale.EnumSlowLogTableSource)},
			},
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamSlowLogCollectInput),
		},
		{
			Key:      paramKeyFirstSqlsScrappedInLastPeriodHours,
			Value:    "24",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamFirstSqlsScrappedHours),
		},
	}
}

func (at *SlowLogTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{
		{
			Param: params.Param{
				Key:      MetricNameQueryTimeAvg,
				Value:    "10",
				Type:     params.ParamTypeFloat64,
				I18nDesc: locale.Bundle.LocalizeAll(locale.ApMetricQueryTimeAvg),
			},
			Operator: params.Operator{
				Value:      ">",
				EnumsValue: defaultOperatorEnums,
			},
		},
		{
			Param: params.Param{
				Key:      MetricNameRowExaminedAvg,
				Value:    "100",
				Type:     params.ParamTypeFloat64,
				I18nDesc: locale.Bundle.LocalizeAll(locale.ApMetricRowExaminedAvg),
			},
			Operator: params.Operator{
				Value:      ">",
				EnumsValue: defaultOperatorEnums,
			},
		},
		defaultAuditLevelOperateParams,
	}
}

func (at *SlowLogTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameLastReceiveTimestamp,
		MetricNameQueryTimeAvg,
		MetricNameQueryTimeMax,
		MetricNameRowExaminedAvg,
		MetricNameDBUser,
		MetricNameEndpoints,
		MetricNameStartTimeOfLastScrapedSQL,
	}
}

func (at *SlowLogTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}

	originSQL.SQLContent = mergedSQL.SQLContent

	// counter
	originCounter := originSQL.Info.Get(MetricNameCounter).Int()
	mergedCounter := mergedSQL.Info.Get(MetricNameCounter).Int()
	counter := originCounter + mergedCounter
	originSQL.Info.SetInt(MetricNameCounter, counter)

	// last_receive_timestamp
	originSQL.Info.SetString(MetricNameLastReceiveTimestamp, mergedSQL.Info.Get(MetricNameLastReceiveTimestamp).String())

	// query_time_avg
	queryTimeAvg := (originSQL.Info.Get(MetricNameQueryTimeAvg).Float()*float64(originCounter) +
		mergedSQL.Info.Get(MetricNameQueryTimeAvg).Float()*float64(mergedCounter)) /
		math.Max(float64(counter), 1)

	originSQL.Info.SetFloat(MetricNameQueryTimeAvg, queryTimeAvg)

	// query_time_max
	queryTimeMax := math.Max(originSQL.Info.Get(MetricNameQueryTimeMax).Float(), mergedSQL.Info.Get(MetricNameQueryTimeMax).Float())
	originSQL.Info.SetFloat(MetricNameQueryTimeMax, queryTimeMax)

	// row_examined_avg
	rowExaminedAvg := (originSQL.Info.Get(MetricNameRowExaminedAvg).Float()*float64(originCounter) +
		mergedSQL.Info.Get(MetricNameRowExaminedAvg).Float()*float64(mergedCounter)) /
		math.Max(float64(counter), 1)

	originSQL.Info.SetFloat(MetricNameRowExaminedAvg, rowExaminedAvg)

	// first_query_at //todo: 这个参数看起来没用上？

	// db_user
	originSQL.Info.SetString(MetricNameDBUser, mergedSQL.Info.Get(MetricNameDBUser).String())

	// endpoints
	originSQL.Info.SetString(MetricNameEndpoints, mergedSQL.Info.Get(MetricNameEndpoints).String())

	// start_time
	originSQL.Info.SetString(MetricNameStartTimeOfLastScrapedSQL, mergedSQL.Info.Get(MetricNameStartTimeOfLastScrapedSQL).String())

	return
}

func (at *SlowLogTaskV2) genSQLV2FromRow(ap *AuditPlan, row map[string]sql.NullString) (*SQLV2, error) {
	query := row["sql_text"].String
	sqlV2 := &SQLV2{
		Source:     ap.Type,
		SourceId:   strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
		ProjectId:  ap.ProjectId,
		InstanceID: ap.InstanceID,
		SQLContent: query,
		SchemaName: row["db"].String,
	}

	fp, err := util.Fingerprint(query, true)
	if err != nil {
		return nil, fmt.Errorf("get sql finger print failed, err: %v", err)
	}
	if fp == "" {
		return nil, fmt.Errorf("get sql finger print failed, fp is empty")
	}
	sqlV2.Fingerprint = fp

	info := NewMetrics()
	// counter
	info.SetInt(MetricNameCounter, 1)

	// latest query time, todo: 是否可以从数据库取
	info.SetString(MetricNameLastReceiveTimestamp, time.Now().Format(time.RFC3339))

	// start time
	info.SetString(MetricNameStartTimeOfLastScrapedSQL, row["start_time"].String)

	// query time avg and max
	queryTime, err := strconv.Atoi(row["query_time"].String)
	if err != nil {
		return nil, fmt.Errorf("unexpected data format: %v, ", row["query_time"].String)
	}
	info.SetFloat(MetricNameQueryTimeAvg, float64(queryTime))
	info.SetFloat(MetricNameQueryTimeMax, float64(queryTime))

	// row examined avg
	rowExamined, err := strconv.ParseFloat(row["rows_examined"].String, 64)
	if err != nil {
		return nil, fmt.Errorf("unexpected data format: %v, ", row["rows_examined"].String)
	}
	info.SetFloat(MetricNameRowExaminedAvg, rowExamined)

	sqlV2.Info = info
	sqlV2.GenSQLId()
	return sqlV2, nil
}

func (at *SlowLogTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	// 如果配置为表采集，则开启采集任务
	if ap.Params.GetParam(paramKeySlowLogCollectInput).Int() != slowlogCollectInputTable {
		return nil, nil
	}
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	instance, exist, err := dms.GetInstancesById(context.Background(), ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}
	if !exist {
		return nil, fmt.Errorf("instance: %v is not exist", ap.InstanceID)
	}

	db, err := executor.NewExecutor(logger, &driverV2.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
	}, "")
	if err != nil {
		return nil, fmt.Errorf("connect to instance fail, error: %v", err)
	}
	defer db.Db.Close()

	queryStartTime, err := persist.GetLatestStartTimeAuditPlanSQLV2(ap.ID)
	if err != nil {
		return nil, fmt.Errorf("get start time failed, error: %v", err)
	}

	if queryStartTime == "" {
		firstScrapInLastHours := ap.Params.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).Int()
		queryStartTime = time.Now().Add(time.Duration(-1*firstScrapInLastHours) * time.Hour).Format("2006-01-02 15:04:05")
	}

	querySQL := `
	SELECT sql_text,db,TIME_TO_SEC(query_time) AS query_time, start_time, rows_examined
	FROM mysql.slow_log
	WHERE sql_text != ''
		AND db NOT IN ('information_schema','performance_schema','mysql','sys')
	`

	cache := NewSQLV2Cache()
	for {
		extraCondition := fmt.Sprintf(" AND start_time>'%s' ORDER BY start_time LIMIT %d", queryStartTime, SlowLogQueryNums)
		execQuerySQL := querySQL + extraCondition

		res, err := db.Db.Query(execQuerySQL)
		if err != nil {
			return nil, fmt.Errorf("query slow log failed, error: %v", err)
		}
		for _, row := range res {
			sqlV2, err := at.genSQLV2FromRow(ap, row)
			if err != nil {
				logger.Warnf("skip collect, error: %v, sql is %s", err, row["sql_text"].String)
			}
			at.AggregateSQL(cache, sqlV2)
		}

		if len(res) < SlowLogQueryNums {
			break
		}

		queryStartTime = res[len(res)-1]["start_time"].String

		time.Sleep(500 * time.Millisecond)
	}
	return cache.GetSQLs(), nil
}

func (at *SlowLogTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *SlowLogTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *SlowLogTaskV2) GetSQLs(ctx context.Context, ap *AuditPlan, persist *model.Storage, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head := []Head{
		{
			Name: "fingerprint",
			Desc: locale.ApSQLFingerprint,
			Type: "sql",
		},
		{
			Name: "sql",
			Desc: locale.ApSQLStatement,
			Type: "sql",
		},
		{
			Name: model.AuditResultName,
			Desc: model.AuditResultDesc,
		},
		{
			Name: "counter",
			Desc: locale.ApNum,
		},
		{
			Name: "last_receive_timestamp",
			Desc: locale.ApLastMatchTime,
		},
		{
			Name: "average_query_time",
			Desc: locale.ApQueryTimeAvg,
		},
		{
			Name: "max_query_time",
			Desc: locale.ApQueryTimeMax,
		},
		{
			Name: "row_examined_avg",
			Desc: locale.ApRowExaminedAvg,
		},
		{
			Name: "db_user",
			Desc: locale.ApMetricNameDBUser,
		},
		{
			Name: "schema",
			Desc: locale.ApSchema,
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		var info = struct {
			Counter              uint64   `json:"counter"`
			LastReceiveTimestamp string   `json:"last_receive_timestamp"`
			AverageQueryTime     *float64 `json:"query_time_avg"`
			MaxQueryTime         *float64 `json:"query_time_max"`
			RowExaminedAvg       *float64 `json:"row_examined_avg"`
			DBUser               string   `json:"db_user"`
		}{}
		err := json.Unmarshal(sql.Info, &info)
		if err != nil {
			return nil, nil, 0, err
		}
		row := map[string]string{
			"sql":                    sql.SQLContent,
			"fingerprint":            sql.Fingerprint,
			"counter":                strconv.FormatUint(info.Counter, 10),
			"last_receive_timestamp": info.LastReceiveTimestamp,
			"db_user":                info.DBUser,
			"schema":                 sql.Schema,
			model.AuditResultName:    sql.AuditResult.GetAuditJsonStrByLangTag(locale.Bundle.GetLangTagFromCtx(ctx)),
		}

		if info.RowExaminedAvg != nil {
			row["row_examined_avg"] = fmt.Sprintf("%.6f", *info.RowExaminedAvg)
		}
		// 兼容之前没有平均执行时间和最长执行时间的数据，没有数据的时候不会在前端显示0.00000导致误解
		if info.AverageQueryTime != nil {
			row["average_query_time"] = fmt.Sprintf("%.6f", *info.AverageQueryTime)
		}
		if info.MaxQueryTime != nil {
			row["max_query_time"] = fmt.Sprintf("%.6f", *info.MaxQueryTime)
		}
		rows = append(rows, row)
	}
	return head, rows, count, nil
}

func (at *SlowLogTaskV2) Head(ap *AuditPlan) []Head {
	return []Head{
		{
			Name: "fingerprint",
			Desc: locale.ApSQLFingerprint,
			Type: "sql",
		},
		{
			Name: "sql",
			Desc: locale.ApSQLStatement,
			Type: "sql",
		},
		{
			Name: "priority",
			Desc: locale.ApPriority,
		},
		{
			Name: model.AuditResultName,
			Desc: model.AuditResultDesc,
		},
		{
			Name:     MetricNameCounter,
			Desc:     locale.ApMetricNameCounter,
			Sortable: true,
		},
		{
			Name:     MetricNameLastReceiveTimestamp,
			Desc:     locale.ApLastMatchTime,
			Sortable: true,
		},
		{
			Name:     MetricNameQueryTimeAvg,
			Desc:     locale.ApMetricNameQueryTimeAvg,
			Sortable: true,
		},
		{
			Name:     MetricNameQueryTimeMax,
			Desc:     locale.ApMetricNameQueryTimeMax,
			Sortable: true,
		},
		{
			Name:     MetricNameRowExaminedAvg,
			Desc:     locale.ApMetricNameRowExaminedAvg,
			Sortable: true,
		},
		{
			Name: MetricNameDBUser,
			Desc: locale.ApMetricNameDBUser,
		},
		{
			Name: "schema_name",
			Desc: locale.ApSchema,
		},
	}
}

func (at *SlowLogTaskV2) Filters(ctx context.Context, logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta {
	return []FilterMeta{
		{
			Name:            "sql", // 模糊筛选
			Desc:            locale.ApSQLStatement,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
		},
		{
			Name:            "rule_name",
			Desc:            locale.ApRuleName,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerRuleTips(ctx, logger, ap.ID, persist),
		},
		{
			Name:            "priority",
			Desc:            locale.ApPriority,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerPriorityTips(ctx, logger),
		},
		{
			Name:            MetricNameDBUser,
			Desc:            locale.ApMetricNameDBUser,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerMetricTips(logger, ap.ID, persist, MetricNameDBUser),
		},
		{
			Name:            "schema_name",
			Desc:            locale.ApSchema,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerSchemaNameTips(logger, ap.ID, persist),
		},
		{
			Name:            MetricNameCounter, // 阈值查询
			Desc:            locale.ApMetricNameCounterMoreThan,
			FilterInputType: FilterInputTypeInt,
			FilterOpType:    FilterOpTypeEqual,
		},
		{
			Name:            MetricNameQueryTimeAvg, // 阈值查询
			Desc:            locale.ApMetricNameQueryTimeAvgMoreThan,
			FilterInputType: FilterInputTypeInt,
			FilterOpType:    FilterOpTypeEqual,
		},
		{
			Name:            MetricNameRowExaminedAvg, // 阈值查询
			Desc:            locale.ApMetricNameRowExaminedAvgMoreThan,
			FilterInputType: FilterInputTypeInt,
			FilterOpType:    FilterOpTypeEqual,
		},
		{
			Name:            MetricNameLastReceiveTimestamp,
			Desc:            locale.ApLastMatchTime,
			FilterInputType: FilterInputTypeDateTime,
			FilterOpType:    FilterOpTypeBetween,
		},
	}
}

func (at *SlowLogTaskV2) GetSQLData(ctx context.Context, ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
	// todo: 需要过滤掉	MetricNameRecordDeleted = true 的记录，因为有分页所以需要在db里过滤，还要考虑概览界面统计的问题
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReqV2(ap.ID, ap.Type, limit, offset, checkAndGetOrderByName(at.Head(ap), orderBy), isAsc, genArgsByFilters(filters))
	if err != nil {
		return nil, count, err
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		data, err := sql.Info.OriginValue()
		if err != nil {
			return nil, 0, err
		}
		info := LoadMetrics(data, at.Metrics())
		rows = append(rows, map[string]string{
			"fingerprint":                  sql.Fingerprint,
			"sql":                          sql.SQLContent,
			"id":                           sql.AuditPlanSqlId,
			"priority":                     sql.Priority.String,
			model.AuditResultName:          sql.AuditResult.GetAuditJsonStrByLangTag(locale.Bundle.GetLangTagFromCtx(ctx)),
			MetricNameCounter:              fmt.Sprint(info.Get(MetricNameCounter).Int()),
			MetricNameLastReceiveTimestamp: info.Get(MetricNameLastReceiveTimestamp).String(),
			MetricNameQueryTimeAvg:         fmt.Sprint(utils.Round(info.Get(MetricNameQueryTimeAvg).Float(), 2)),
			MetricNameQueryTimeMax:         fmt.Sprint(utils.Round(info.Get(MetricNameQueryTimeMax).Float(), 2)),
			MetricNameRowExaminedAvg:       fmt.Sprint(utils.Round(info.Get(MetricNameRowExaminedAvg).Float(), 2)),
			MetricNameDBUser:               info.Get(MetricNameDBUser).String(),
			"schema_name":                  sql.Schema,
		})
	}
	return rows, count, nil
}
