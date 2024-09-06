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

type PGTopSQLTaskV2 struct {
	DefaultTaskV2
}

func NewPGTopSQLTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &PGTopSQLTaskV2{DefaultTaskV2: DefaultTaskV2{}}
	}
}

func (at *PGTopSQLTaskV2) InstanceType() string {
	return InstanceTypePostgreSQL
}

func (at *PGTopSQLTaskV2) Params(instanceId ...string) params.Params {
	id := ""
	if len(instanceId) != 0 {
		id = instanceId[0]
	}
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
			Value: DynPerformanceViewPgSQLColumnElapsedTime,
			Type:  params.ParamTypeString,
		},
		{
			Key:   paramKeySchema,
			Desc:  "schema",
			Value: "postgres",
			Type:  params.ParamTypeString,
			Enums: ShowSchemaEnumsByInstanceId(id),
		},
	}
}

func (at *PGTopSQLTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameQueryTimeTotal,
		MetricNameUserIOWaitTimeTotal,
		MetricNameDiskReadTotal,
		MetricNameBufferGetCounter,
	}
}

func (at *PGTopSQLTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}
	// counter
	originSQL.Info.SetInt(MetricNameCounter, mergedSQL.Info.Get(MetricNameCounter).Int())

	// MetricNameQueryTimeTotal
	originSQL.Info.SetFloat(MetricNameQueryTimeTotal, mergedSQL.Info.Get(MetricNameQueryTimeTotal).Float())

	// MetricNameUserIOWaitTimeTotal
	originSQL.Info.SetFloat(MetricNameUserIOWaitTimeTotal, mergedSQL.Info.Get(MetricNameUserIOWaitTimeTotal).Float())

	// MetricNameDiskReadTotal
	originSQL.Info.SetFloat(MetricNameDiskReadTotal, mergedSQL.Info.Get(MetricNameDiskReadTotal).Float())

	// MetricNameBufferGetCounter
	originSQL.Info.SetFloat(MetricNameBufferGetCounter, mergedSQL.Info.Get(MetricNameBufferGetCounter).Float())
	return
}

type DynPerformancePgColumns struct {
	SQLFullText    string  `json:"sql_fulltext"`
	Executions     float64 `json:"executions"`
	ElapsedTime    float64 `json:"elapsed_time"`
	DiskReads      float64 `json:"disk_reads"`
	BufferGets     float64 `json:"buffer_gets"`
	UserIOWaitTime float64 `json:"user_io_wait_time"`
}

const (
	DynPerformanceViewPgTpl = `
	select
		pss.query as sql_fulltext,
		sum ( pss.calls ) as executions,
		sum ( pss.total_exec_time ) as elapsed_time,
		sum ( pss.shared_blks_read ) as disk_reads,
		sum ( pss.shared_blks_hit ) as buffer_gets,
		sum ( pss.blk_read_time ) as user_io_wait_time 
	from
		pg_stat_statements pss
		join pg_database pd on pss.dbid = pd.oid 
	where
		pss.calls > 0 
		and query <> '<insufficient privilege>'
		and pd.datname = '%v' 
	group by pss.query 
	order by %v desc limit '%v'`
	DynPerformanceViewPgSQLColumnExecutions     = "executions"
	DynPerformanceViewPgSQLColumnElapsedTime    = "elapsed_time"
	DynPerformanceViewPgSQLColumnDiskReads      = "disk_reads"
	DynPerformanceViewPgSQLColumnBufferGets     = "buffer_gets"
	DynPerformanceViewPgSQLColumnUserIOWaitTime = "user_io_wait_time"
)

func (at *PGTopSQLTaskV2) queryTopSQLsForPg(inst *model.Instance, database string, orderBy string, topN int) ([]*DynPerformancePgColumns, error) {
	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, database)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	ctxForCreatePgExtension, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 执行创建pg_stat_statements扩展
	_, err = plugin.Exec(ctxForCreatePgExtension, `CREATE EXTENSION IF NOT EXISTS pg_stat_statements`)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	sql := fmt.Sprintf(DynPerformanceViewPgTpl, database, orderBy, topN)
	result, err := plugin.Query(ctx, sql, &driverV2.QueryConf{TimeOutSecond: 120})
	if err != nil {
		return nil, err
	}
	var ret []*DynPerformancePgColumns
	rows := result.Rows
	for _, row := range rows {
		values := row.Values
		if len(values) < 6 {
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
		diskReads, err := strconv.ParseFloat(values[3].Value, 64)
		if err != nil {
			return nil, err
		}
		bufferGets, err := strconv.ParseFloat(values[4].Value, 64)
		if err != nil {
			return nil, err
		}
		userIoWaitTime, err := strconv.ParseFloat(values[5].Value, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &DynPerformancePgColumns{
			SQLFullText:    values[0].Value,
			Executions:     executions,
			ElapsedTime:    elapsedTime,
			DiskReads:      diskReads,
			BufferGets:     bufferGets,
			UserIOWaitTime: userIoWaitTime,
		})
	}
	return ret, nil
}

func (at *PGTopSQLTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	schema := ap.Params.GetParam(paramKeySchema).String()

	// 设置默认数据库为：postgres，因为连接PG必须指定数据库
	if len(schema) == 0 {
		schema = "postgres"
	}
	// 超时2分钟
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	inst, _, err := dms.GetInstancesById(ctx, ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}

	sqls, err := at.queryTopSQLsForPg(inst, schema, ap.Params.GetParam("order_by_column").String(),
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
			SourceId:    strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
			ProjectId:   ap.ProjectId,
			InstanceID:  ap.InstanceID,
			SchemaName:  "", // todo: top sql 未采集schema, 需要填充
			Info:        info,
			SQLContent:  sql.SQLFullText,
			Fingerprint: sql.SQLFullText,
		}
		info.SetInt(MetricNameCounter, int64(sql.Executions))
		info.SetFloat(MetricNameQueryTimeTotal, float64(sql.ElapsedTime))
		info.SetFloat(MetricNameUserIOWaitTimeTotal, float64(sql.UserIOWaitTime))
		info.SetFloat(MetricNameDiskReadTotal, float64(sql.DiskReads))
		info.SetFloat(MetricNameBufferGetCounter, float64(sql.BufferGets))
		sqlV2.GenSQLId()
		at.AggregateSQL(cache, sqlV2)
	}
	return cache.GetSQLs(), nil
}

func (at *PGTopSQLTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *PGTopSQLTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *PGTopSQLTaskV2) Head(ap *AuditPlan) []Head {
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
			Name: MetricNameDiskReadTotal,
			Desc: "物理读块数",
		},
		{
			Name: MetricNameBufferGetCounter,
			Desc: "逻辑读块数",
		},
		{
			Name: MetricNameUserIOWaitTimeTotal,
			Desc: "I/O等待时间(s)",
		},
	}
}

func (at *PGTopSQLTaskV2) GetSQLData(ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
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
			"priority":                    sql.Priority.String,
			MetricNameCounter:             strconv.Itoa(int(info.Get(MetricNameCounter).Int())),
			MetricNameQueryTimeTotal:      fmt.Sprintf("%v", utils.Round(float64(info.Get(MetricNameQueryTimeTotal).Float())/1000, 3)), //视图中时间单位是毫秒，所以除以1000得到秒
			MetricNameDiskReadTotal:       strconv.Itoa(int(info.Get(MetricNameDiskReadTotal).Int())),
			MetricNameBufferGetCounter:    strconv.Itoa(int(info.Get(MetricNameBufferGetCounter).Int())),
			MetricNameUserIOWaitTimeTotal: fmt.Sprintf("%v", utils.Round(float64(info.Get(MetricNameUserIOWaitTimeTotal).Float())/1000, 3)), //视图中时间单位是毫秒，所以除以1000得到秒
			model.AuditResultName:         sql.AuditResult.String,
		})
	}
	return rows, count, nil
}
