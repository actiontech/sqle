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
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
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

func (at *SlowLogTaskV2) Params() func(instanceId ...string) params.Params {
	return func(instanceId ...string) params.Params {
		return []*params.Param{
			{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟，仅对 mysql.slow_log 有效）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeySlowLogCollectInput,
				Desc:  "采集来源",
				Value: "0",
				Type:  params.ParamTypeInt,
				Enums: []params.EnumsValue{
					{Value: "0", Desc: "从slow.log 文件采集,需要适配scanner"}, {Value: "1", Desc: "从mysql.slow_log 表采集"},
				},
			},
		}
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

	return
}

func (at *SlowLogTaskV2) genSQLV2FromRow(ap *AuditPlan, row map[string]sql.NullString) (*SQLV2, error) {
	query := row["sql_text"].String
	sqlV2 := &SQLV2{
		Source:       ap.Type,
		SourceId:     ap.ID,
		ProjectId:    ap.ProjectId,
		InstanceName: ap.InstanceName,
		SQLContent:   query,
		SchemaName:   row["db"].String,
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
	if ap.InstanceName == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	instance, _, err := dms.GetInstanceInProjectByName(context.Background(), string(ap.ProjectId), ap.InstanceName)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
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

func (at *SlowLogTaskV2) Audit(sqls []*model.OriginManageSQL) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *SlowLogTaskV2) GetSQLs(ap *AuditPlan, persist *model.Storage, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head := []Head{
		{
			Name: "fingerprint",
			Desc: "SQL指纹",
			Type: "sql",
		},
		{
			Name: "sql",
			Desc: "SQL",
			Type: "sql",
		},
		{
			Name: model.AuditResultName,
			Desc: model.AuditResultDesc,
		},
		{
			Name: "counter",
			Desc: "数量",
		},
		{
			Name: "last_receive_timestamp",
			Desc: "最后匹配时间",
		},
		{
			Name: "average_query_time",
			Desc: "平均执行时间",
		},
		{
			Name: "max_query_time",
			Desc: "最长执行时间",
		},
		{
			Name: "row_examined_avg",
			Desc: "平均扫描行数",
		},
		{
			Name: "db_user",
			Desc: "用户",
		},
		{
			Name: "schema",
			Desc: "Schema",
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
			model.AuditResultName:    sql.AuditResult.String,
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
