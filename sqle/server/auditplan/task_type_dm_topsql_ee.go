//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/sirupsen/logrus"
)

type DmTopSQLTaskV2 struct{}

func NewDmTopSQLTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &DmTopSQLTaskV2{}
	}
}

func (at *DmTopSQLTaskV2) InstanceType() string {
	return InstanceTypeDm
}

func (at *DmTopSQLTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:   paramKeyCollectIntervalMinute,
			Desc:  "采集周期（分钟）",
			Value: "60",
			Type:  params.ParamTypeInt,
		},
		{
			Key:   "top_n",
			Desc:  "Top N",
			Value: "3",
			Type:  params.ParamTypeInt,
		},
		{
			Key:   "order_by_column",
			Desc:  "排序字段",
			Value: DmTopSQLMetricTotalExecTime,
			Type:  params.ParamTypeString,
		},
	}
}

func (at *DmTopSQLTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{}
}

func (at *DmTopSQLTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameQueryTimeTotal,
		MetricNameQueryTimeAvg,
		MetricNameCPUTimeTotal,
		MetricNamePhyReadPageTotal,
		MetricNameLogicReadPageTotal,
	}
}

func (at *DmTopSQLTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}
	// counter
	originSQL.Info.SetInt(MetricNameCounter, mergedSQL.Info.Get(MetricNameCounter).Int())

	// MetricNameQueryTimeTotal
	originSQL.Info.SetFloat(MetricNameQueryTimeTotal, mergedSQL.Info.Get(MetricNameQueryTimeTotal).Float())

	// MetricNameQueryTimeAvg
	originSQL.Info.SetFloat(MetricNameQueryTimeAvg, mergedSQL.Info.Get(MetricNameQueryTimeTotal).Float())

	// MetricNameCPUTimeTotal
	originSQL.Info.SetFloat(MetricNameCPUTimeTotal, mergedSQL.Info.Get(MetricNameCPUTimeTotal).Float())

	// MetricNamePhyReadPageTotal
	originSQL.Info.SetInt(MetricNamePhyReadPageTotal, mergedSQL.Info.Get(MetricNamePhyReadPageTotal).Int())

	// MetricNameLogicReadPageTotal
	originSQL.Info.SetInt(MetricNameLogicReadPageTotal, mergedSQL.Info.Get(MetricNameLogicReadPageTotal).Int())
	return
}

type DynPerformanceDmColumns struct {
	SQLFullText      string  `json:"sql_fulltext"`
	Executions       float64 `json:"executions"`
	TotalExecTime    float64 `json:"total_exec_time"`
	AverageExecTime  float64 `json:"average_exec_time"`
	CPUTime          float64 `json:"cpu_time"`
	PhyReadPageCnt   float64 `json:"phy_read_page_cnt"`
	LogicReadPageCnt float64 `json:"logic_read_page_cnt"`
}

// Dm Top SQL
const (
	DynPerformanceViewDmTpl = `
SELECT
    sql_fulltext,
    executions,
    total_exec_time,
    average_exec_time,
    cpu_time,
    phy_read_page_cnt,
    logic_read_page_cnt
FROM (
    SELECT
        SQL_TXT AS sql_fulltext,
        COUNT(*) AS executions,
        SUM(EXEC_TIME) AS total_exec_time,
        SUM(EXEC_TIME) / COUNT(*) OVER () AS average_exec_time,
        (SUM(EXEC_TIME) - SUM(PARSE_TIME) - SUM(IO_WAIT_TIME)) AS cpu_time,
        SUM(PHY_READ_CNT) AS phy_read_page_cnt,
        SUM(LOGIC_READ_CNT) AS logic_read_page_cnt,
        ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) AS row_num
    FROM V$SQL_STAT_HISTORY
    GROUP BY SQL_TXT
) t WHERE executions > 0 AND row_num <= %v ORDER BY %v DESC`
	DmTopSQLMetricExecutions       = "executions"
	DmTopSQLMetricTotalExecTime    = "total_exec_time"
	DmTopSQLMetricAverageExecTime  = "average_exec_time"
	DmTopSQLMetricCPUTime          = "cpu_time"
	DmTopSQLMetricPhyReadPageCnt   = "phy_read_page_cnt"
	DmTopSQLMetricLogicReadPageCnt = "logic_read_page_cnt"
)

func (at *DmTopSQLTaskV2) queryTopSQLsForDm(inst *model.Instance, database string, orderBy string, topN int) ([]*DynPerformanceDmColumns, error) {
	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, database)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	sql := fmt.Sprintf(DynPerformanceViewDmTpl, topN, orderBy)
	result, err := plugin.Query(ctx, sql, &driverV2.QueryConf{TimeOutSecond: 120})
	if err != nil {
		return nil, err
	}
	var ret []*DynPerformanceDmColumns
	rows := result.Rows
	for _, row := range rows {
		values := row.Values
		if len(values) < 7 {
			continue
		}
		executions, err := strconv.ParseFloat(values[1].Value, 64)
		if err != nil {
			return nil, err
		}
		totalExecTime, err := strconv.ParseFloat(values[2].Value, 64)
		if err != nil {
			return nil, err
		}
		averageExecTime, err := strconv.ParseFloat(values[3].Value, 64)
		if err != nil {
			return nil, err
		}
		cpuTime, err := strconv.ParseFloat(values[4].Value, 64)
		if err != nil {
			return nil, err
		}
		phyReadPageCnt, err := strconv.ParseFloat(values[5].Value, 64)
		if err != nil {
			return nil, err
		}
		logicReadPageCnt, err := strconv.ParseFloat(values[6].Value, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &DynPerformanceDmColumns{
			SQLFullText:      values[0].Value,
			Executions:       executions,
			TotalExecTime:    totalExecTime,
			AverageExecTime:  averageExecTime,
			CPUTime:          cpuTime,
			PhyReadPageCnt:   phyReadPageCnt,
			LogicReadPageCnt: logicReadPageCnt,
		})
	}
	return ret, nil
}

func (at *DmTopSQLTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	inst, _, err := dms.GetInstancesById(ctx, ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}

	sqls, err := at.queryTopSQLsForDm(inst, "", ap.Params.GetParam("order_by_column").String(),
		ap.Params.GetParam("top_n").Int())
	if err != nil {
		return nil, fmt.Errorf("query top sql fail, error: %v", err)
	}

	if len(sqls) == 0 {
		logger.Info("sql result count is 0")
		return nil, nil
	}
	cache := NewSQLV2Cache()

	for _, sql := range sqls {
		info := NewMetrics()
		sqlV2 := &SQLV2{
			Source:      ap.Type,
			SourceId:    ap.ID,
			ProjectId:   ap.ProjectId,
			InstanceID:  ap.InstanceID,
			SchemaName:  "", // todo: top sql 未采集schema, 需要填充
			Info:        info,
			SQLContent:  sql.SQLFullText,
			Fingerprint: sql.SQLFullText,
		}
		info.SetInt(MetricNameCounter, int64(sql.Executions))
		info.SetFloat(MetricNameQueryTimeTotal, sql.TotalExecTime)
		info.SetFloat(MetricNameQueryTimeAvg, sql.AverageExecTime)
		info.SetFloat(MetricNameCPUTimeTotal, sql.CPUTime)
		info.SetInt(MetricNamePhyReadPageTotal, int64(sql.PhyReadPageCnt))
		info.SetInt(MetricNameLogicReadPageTotal, int64(sql.LogicReadPageCnt))
		sqlV2.GenSQLId()
		at.AggregateSQL(cache, sqlV2)
	}
	return cache.GetSQLs(), nil
}

func (at *DmTopSQLTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *DmTopSQLTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *DmTopSQLTaskV2) Head(ap *AuditPlan) []Head {
	return []Head{
		{
			Name: "sql",
			Desc: "SQL语句",
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
			Name: MetricNameCounter,
			Desc: "总执行次数",
		},
		{
			Name: MetricNameQueryTimeTotal,
			Desc: "总执行时间(s)",
		},
		{
			Name: MetricNameQueryTimeAvg,
			Desc: "平均执行时间(s)",
		},
		{
			Name: MetricNameCPUTimeTotal,
			Desc: "CPU时间占用(s)",
		},
		{
			Name: MetricNamePhyReadPageTotal,
			Desc: "物理读页数",
		},
		{
			Name: MetricNameLogicReadPageTotal,
			Desc: "逻辑读页数",
		},
	}
}

func (at *DmTopSQLTaskV2) Filters(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta {
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

func (at *DmTopSQLTaskV2) GetSQLData(ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
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
		data, err := sql.Info.OriginValue()
		if err != nil {
			return nil, 0, err
		}
		info := LoadMetrics(data, at.Metrics())
		rows = append(rows, map[string]string{
			"sql":                        sql.SQLContent,
			"id":                         sql.AuditPlanSqlId,
			"priority":                   sql.Priority.String,
			MetricNameCounter:            strconv.Itoa(int(info.Get(MetricNameCounter).Int())),
			MetricNameQueryTimeTotal:     fmt.Sprintf("%v", utils.Round(float64(info.Get(MetricNameQueryTimeTotal).Float())/1000, 3)), //视图中时间单位是毫秒，所以除以1000得到秒
			MetricNameQueryTimeAvg:       fmt.Sprintf("%v", utils.Round(float64(info.Get(MetricNameQueryTimeAvg).Float())/1000, 3)),   //视图中时间单位是毫秒，所以除以1000得到秒
			MetricNameCPUTimeTotal:       fmt.Sprintf("%v", utils.Round(float64(info.Get(MetricNameCPUTimeTotal).Float())/1000, 3)),   //视图中时间单位是毫秒，所以除以1000得到秒
			MetricNamePhyReadPageTotal:   strconv.Itoa(int(info.Get(MetricNamePhyReadPageTotal).Int())),
			MetricNameLogicReadPageTotal: strconv.Itoa(int(info.Get(MetricNameLogicReadPageTotal).Int())),
			model.AuditResultName:        sql.AuditResult.String,
		})
	}
	return rows, count, nil
}
