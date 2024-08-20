//go:build enterprise
// +build enterprise

package auditplan

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/sirupsen/logrus"
)

type TBaseSlowLogTaskV2 struct {
	lastEndTime *time.Time
}

func NewTBaseSlowLogTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &TBaseSlowLogTaskV2{}
	}
}

func (at *TBaseSlowLogTaskV2) InstanceType() string {
	return InstanceTypeTBase
}

func (at *TBaseSlowLogTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{}
}

func (at *TBaseSlowLogTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{}
}

func (at *TBaseSlowLogTaskV2) Metrics() []string {
	return []string{}
}

func (at *TBaseSlowLogTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
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

func (at *TBaseSlowLogTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *TBaseSlowLogTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *TBaseSlowLogTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	return nil, nil
}

func (at *TBaseSlowLogTaskV2) Head(ap *AuditPlan) []Head {
	return []Head{
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
			Name: "priority",
			Desc: "优先级",
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
}

func (at *TBaseSlowLogTaskV2) Filters(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta {
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

func (at *TBaseSlowLogTaskV2) GetSQLData(ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
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
			return nil, 0, err
		}
		row := map[string]string{
			"sql":                    sql.SQLContent,
			"id":                     sql.AuditPlanSqlId,
			"fingerprint":            sql.Fingerprint,
			"priority":               sql.Priority.String,
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
	return rows, count, nil
}
