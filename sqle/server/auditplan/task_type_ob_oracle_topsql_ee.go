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

type ObForOracleTopSQLTaskV2 struct {
	obVersion string
	DefaultTaskV2
}

func NewObForOracleTopSQLTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &ObForOracleTopSQLTaskV2{DefaultTaskV2: DefaultTaskV2{}}
	}
}

func (at *ObForOracleTopSQLTaskV2) InstanceType() string {
	return InstanceTypeObForOracle
}

func (at *ObForOracleTopSQLTaskV2) Params(instanceId ...string) params.Params {
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
			Value: DynPerformanceViewObForOracleColumnElapsedTime,
			Type:  params.ParamTypeString,
		},
	}
}

func (at *ObForOracleTopSQLTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameQueryTimeTotal,
		MetricNameCPUTimeTotal,
		MetricNameDiskReadTotal,
		MetricNameBufferGetCounter,
		MetricNameUserIOWaitTimeTotal,
	}
}

func (at *ObForOracleTopSQLTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}
	// counter
	originSQL.Info.SetInt(MetricNameCounter, mergedSQL.Info.Get(MetricNameCounter).Int())

	// MetricNameQueryTimeTotal
	originSQL.Info.SetFloat(MetricNameQueryTimeTotal, mergedSQL.Info.Get(MetricNameQueryTimeTotal).Float())

	// MetricNameCPUTimeTotal
	originSQL.Info.SetFloat(MetricNameCPUTimeTotal, mergedSQL.Info.Get(MetricNameCPUTimeTotal).Float())

	// MetricNameDiskReadTotal
	originSQL.Info.SetInt(MetricNameDiskReadTotal, mergedSQL.Info.Get(MetricNameDiskReadTotal).Int())

	// MetricNameBufferGetCounter
	originSQL.Info.SetInt(MetricNameBufferGetCounter, mergedSQL.Info.Get(MetricNameBufferGetCounter).Int())

	// MetricNameUserIOWaitTimeTotal
	originSQL.Info.SetFloat(MetricNameUserIOWaitTimeTotal, mergedSQL.Info.Get(MetricNameUserIOWaitTimeTotal).Float())
	return
}

type DynPerformanceObForOracleColumns struct {
	SQLFullText    string  `json:"sql_fulltext"`
	Executions     float64 `json:"executions"`
	ElapsedTime    float64 `json:"elapsed_time"`
	CPUTime        float64 `json:"cpu_time"`
	DiskReads      float64 `json:"disk_reads"`
	BufferGets     float64 `json:"buffer_gets"`
	UserIOWaitTime float64 `json:"user_io_wait_time"`
}

const (
	DynPerformanceViewObForOracleTpl = `
SELECT
    t1.sql_fulltext as sql_fulltext,
    sum(t1.EXECUTIONS) as executions,
    sum(t1.ELAPSED_TIME) as elapsed_time,
    sum(t1.CPU_TIME) as cpu_time,
    sum(t1.DISK_READS) as disk_reads,
    sum(t1.BUFFERS_GETS) as buffer_gets,
    sum(t1.USER_IO_WAIT_TIME) as user_io_wait_time
FROM 
	(
		SELECT
			to_char(QUERY_SQL) as sql_fulltext,
			EXECUTIONS,
			ELAPSED_TIME,
			CPU_TIME,
			DISK_READS,
			BUFFERS_GETS,
			USER_IO_WAIT_TIME
		FROM 
			%v
		WHERE
			to_char(QUERY_SQL) != 'null'
	) t1 
GROUP BY 
    t1.sql_fulltext
ORDER BY 
    %v DESC
FETCH FIRST %v ROWS ONLY
`
	DynPerformanceViewObForOracleColumnExecutions     = "executions"
	DynPerformanceViewObForOracleColumnElapsedTime    = "elapsed_time"
	DynPerformanceViewObForOracleColumnCPUTime        = "cpu_time"
	DynPerformanceViewObForOracleColumnDiskReads      = "disk_reads"
	DynPerformanceViewObForOracleColumnBufferGets     = "buffer_gets"
	DynPerformanceViewObForOracleColumnUserIOWaitTime = "user_io_wait_time"
)

/*
查询OceanBase版本的SQL语句

 1. 根据文档查询到的支持范围包括：3.2.3-4.3.1(最新)
 2. 测试时使用OceanBase版本3.1.5，3.1.5也支持该用法

参考链接：https://www.oceanbase.com/quicksearch?q=OB_VERSION
*/
const GetOceanbaseVersionSQL string = "SELECT OB_VERSION() FROM DUAL"
const OceanBaseVersion4_0_0 string = "4.0.0"

// ob for oracle 是从4.1.0开始适配的，因此默认版本设定为4.1.0
const DefaultOBForOracleVersion string = "4.1.0"

func (at *ObForOracleTopSQLTaskV2) getOceanBaseVersion(ctx context.Context, inst *model.Instance, database string) (string, error) {
	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, database)
	if err != nil {
		return "", err
	}
	defer plugin.Close(context.TODO())
	versionResult, err := plugin.Query(ctx, GetOceanbaseVersionSQL, &driverV2.QueryConf{TimeOutSecond: 20})
	if err != nil {
		return "", err
	}
	if len(versionResult.Column) == 1 && len(versionResult.Rows) == 1 {
		return versionResult.Rows[0].Values[0].Value, nil
	}
	return "", fmt.Errorf("unexpected result of ob version")
}

// 从OceanBase4.0.0版本开始，GV$PLAN_CACHE_PLAN_STAT视图名称调整为GV$OB_PLAN_CACHE_PLAN_STAT
// 参考：https://www.oceanbase.com/docs/common-oceanbase-database-cn-1000000000820360
const PlanCacheViewNameOBV4After string = "GV$OB_PLAN_CACHE_PLAN_STAT"
const PlanCacheViewNameOBV4Before string = "GV$PLAN_CACHE_PLAN_STAT"

func (at *ObForOracleTopSQLTaskV2) getPlanCacheViewNameByObVersion(obVersion string) string {
	var viewName string = PlanCacheViewNameOBV4After
	isLess, err := utils.IsVersionLessThan(obVersion, OceanBaseVersion4_0_0)
	if err != nil {
		log.Logger().Errorf("compare ocean base version failed, use default view name for top sql %v, error is %v", viewName, err)
		isLess = false
	}
	if isLess {
		viewName = PlanCacheViewNameOBV4Before
	}
	return viewName
}

func (at *ObForOracleTopSQLTaskV2) queryTopSQLs(inst *model.Instance, database string, obVersion string, orderBy string, topN int) ([]*DynPerformanceObForOracleColumns, error) {
	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, database)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	sql := fmt.Sprintf(
		DynPerformanceViewObForOracleTpl,
		// 通过视图PlanCacheView查询OecanBase for Oracle的Top SQL
		at.getPlanCacheViewNameByObVersion(obVersion), orderBy, topN,
	)
	result, err := plugin.Query(ctx, sql, &driverV2.QueryConf{TimeOutSecond: 20})
	if err != nil {
		return nil, err
	}
	var ret []*DynPerformanceObForOracleColumns
	rows := result.Rows
	for _, row := range rows {
		values := row.Values
		if len(values) == 0 {
			continue
		}
		executions, err := strconv.ParseFloat(values[1].Value, 64)
		if err != nil {
			return nil, err
		}
		elapsedTime, err := strconv.ParseFloat(values[2].Value, 64)
		if err != nil {
			return nil, err
		}
		cpuTime, err := strconv.ParseFloat(values[3].Value, 64)
		if err != nil {
			return nil, err
		}
		diskReads, err := strconv.ParseFloat(values[4].Value, 64)
		if err != nil {
			return nil, err
		}
		bufferGets, err := strconv.ParseFloat(values[5].Value, 64)
		if err != nil {
			return nil, err
		}
		userIoWaitTime, err := strconv.ParseFloat(values[6].Value, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &DynPerformanceObForOracleColumns{
			SQLFullText:    values[0].Value,
			Executions:     executions,
			ElapsedTime:    elapsedTime,
			CPUTime:        cpuTime,
			DiskReads:      diskReads,
			BufferGets:     bufferGets,
			UserIOWaitTime: userIoWaitTime,
		})
	}
	return ret, nil
}

func (at *ObForOracleTopSQLTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	inst, _, err := dms.GetInstancesById(ctx, ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}

	if at.obVersion == "" {
		at.obVersion, err = at.getOceanBaseVersion(ctx, inst, "")
		if err != nil {
			log.Logger().Errorf("get ocean base version failed, use default version %v, error is %v", DefaultOBForOracleVersion, err)
			at.obVersion = DefaultOBForOracleVersion
		}
	}

	sqls, err := at.queryTopSQLs(inst, "", at.obVersion,
		ap.Params.GetParam("order_by_column").String(),
		ap.Params.GetParam("top_n").Int())
	if err != nil {
		logger.Errorf("query top sql fail, error: %v", err)
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
		info.SetFloat(MetricNameQueryTimeTotal, sql.ElapsedTime)
		info.SetFloat(MetricNameCPUTimeTotal, sql.CPUTime)
		info.SetInt(MetricNameDiskReadTotal, int64(sql.DiskReads))
		info.SetInt(MetricNameBufferGetCounter, int64(sql.BufferGets))
		info.SetFloat(MetricNameUserIOWaitTimeTotal, sql.UserIOWaitTime)
		sqlV2.GenSQLId()
		at.AggregateSQL(cache, sqlV2)
	}
	return cache.GetSQLs(), nil
}

func (at *ObForOracleTopSQLTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *ObForOracleTopSQLTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *ObForOracleTopSQLTaskV2) Head(ap *AuditPlan) []Head {
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
			Desc: "执行时间(s)",
		},
		{
			Name: MetricNameCPUTimeTotal,
			Desc: "CPU消耗时间(s)",
		},
		{
			Name: MetricNameDiskReadTotal,
			Desc: "物理读次数",
		},
		{
			Name: MetricNameBufferGetCounter,
			Desc: "逻辑读次数",
		},
		{
			Name: MetricNameUserIOWaitTimeTotal,
			Desc: "I/O等待时间(s)",
		},
	}
}

func (at *ObForOracleTopSQLTaskV2) GetSQLData(ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
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
			"sql":                         sql.SQLContent,
			"id":                          sql.AuditPlanSqlId,
			MetricNameCounter:             strconv.Itoa(int(info.Get(MetricNameCounter).Int())),
			"priority":                    sql.Priority.String,
			MetricNameQueryTimeTotal:      fmt.Sprintf("%v", utils.Round(float64(info.Get(MetricNameQueryTimeTotal).Float())/1000, 3)), //视图中时间单位是毫秒，所以除以1000得到秒
			MetricNameCPUTimeTotal:        fmt.Sprintf("%v", utils.Round(float64(info.Get(MetricNameCPUTimeTotal).Float())/1000, 3)),   //视图中时间单位是毫秒，所以除以1000得到秒
			MetricNameDiskReadTotal:       strconv.Itoa(int(info.Get(MetricNameDiskReadTotal).Int())),
			MetricNameBufferGetCounter:    strconv.Itoa(int(info.Get(MetricNameBufferGetCounter).Int())),
			MetricNameUserIOWaitTimeTotal: fmt.Sprintf("%v", utils.Round(float64(info.Get(MetricNameUserIOWaitTimeTotal).Float())/1000, 3)), //视图中时间单位是毫秒，所以除以1000得到秒
			model.AuditResultName:         sql.AuditResult.String,
		})
	}
	return rows, count, nil
}
