package auditplan

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/oracle"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/sirupsen/logrus"
)

type OracleTopSQLTaskV2 struct{}

func NewOracleTopSQLTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &OracleTopSQLTaskV2{}
	}
}

func (at *OracleTopSQLTaskV2) InstanceType() string {
	return InstanceTypeOracle
}

func (at *OracleTopSQLTaskV2) Params(instanceId ...string) params.Params {
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
			Desc:  "V$SQLAREA中的排序字段",
			Value: oracle.DynPerformanceViewSQLAreaColumnElapsedTime,
			Type:  params.ParamTypeString,
		},
	}
}

func (at *OracleTopSQLTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{}
}

func (at *OracleTopSQLTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameQueryTimeTotal,
		MetricNameUserIOWaitTimeTotal,
		MetricNameCPUTimeTotal,
		MetricNameDiskReadTotal,
		MetricNameBufferGetCounter,
	}
}

func (at *OracleTopSQLTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}
	// counter
	originSQL.Info.SetInt(MetricNameCounter, mergedSQL.Info.Get(MetricNameCounter).Int())

	// // MetricNameQueryTimeTotal
	originSQL.Info.SetFloat(MetricNameQueryTimeTotal, mergedSQL.Info.Get(MetricNameQueryTimeTotal).Float())

	// // MetricNameUserIOWaitTimeTotal
	originSQL.Info.SetFloat(MetricNameUserIOWaitTimeTotal, mergedSQL.Info.Get(MetricNameUserIOWaitTimeTotal).Float())

	// // MetricNameCPUTimeTotal
	originSQL.Info.SetFloat(MetricNameCPUTimeTotal, mergedSQL.Info.Get(MetricNameCPUTimeTotal).Float())

	// // MetricNameDiskReadTotal
	originSQL.Info.SetFloat(MetricNameDiskReadTotal, mergedSQL.Info.Get(MetricNameDiskReadTotal).Float())

	// MetricNameBufferGetCounter
	originSQL.Info.SetFloat(MetricNameBufferGetCounter, mergedSQL.Info.Get(MetricNameBufferGetCounter).Float())
	return
}

func (at *OracleTopSQLTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	inst, _, err := dms.GetInstancesById(ctx, ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}
	// This depends on: https://github.com/actiontech/sqle-oracle-plugin.
	// If your Oracle db plugin does not implement the parameter `service_name`,
	// you can only use the default service name `XE`.
	// TODO: using DB plugin to query SQL.
	serviceName := inst.AdditionalParams.GetParam("service_name").String()
	dsn := &oracle.DSN{
		Host:        inst.Host,
		Port:        inst.Port,
		User:        inst.User,
		Password:    inst.Password,
		ServiceName: serviceName,
	}
	db, err := oracle.NewDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to instance fail, error: %v", err)
	}
	defer db.Close()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	sqls, err := db.QueryTopSQLs(ctx, ap.Params.GetParam("top_n").Int(), ap.Params.GetParam("order_by_column").String())
	if err != nil {
		return nil, fmt.Errorf("query top sql fail, error: %v", err)
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
		info.SetInt(MetricNameCounter, sql.Executions)
		info.SetFloat(MetricNameQueryTimeTotal, float64(sql.ElapsedTime))
		info.SetFloat(MetricNameUserIOWaitTimeTotal, float64(sql.UserIOWaitTime))
		info.SetFloat(MetricNameCPUTimeTotal, float64(sql.CPUTime))
		info.SetFloat(MetricNameDiskReadTotal, float64(sql.DiskReads))
		info.SetFloat(MetricNameBufferGetCounter, float64(sql.BufferGets))
		sqlV2.GenSQLId()
		err = at.AggregateSQL(cache, sqlV2)
		if err != nil {
			logger.Warnf("aggregate sql failed,error : %v", err)
			continue
		}
	}
	return cache.GetSQLs(), nil
}

func (at *OracleTopSQLTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *OracleTopSQLTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *OracleTopSQLTaskV2) Head(ap *AuditPlan) []Head {
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
			Desc: "物理读",
		},
		{
			Name: MetricNameBufferGetCounter,
			Desc: "逻辑读",
		},
		{
			Name: MetricNameUserIOWaitTimeTotal,
			Desc: "I/O等待时间(s)",
		},
	}
}

func (at *OracleTopSQLTaskV2) Filters(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta {
	return []FilterMeta{}
}

func (at *OracleTopSQLTaskV2) GetSQLData(ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReqV2(ap.ID, ap.Type, limit, offset, checkAndGetOrderByName(at.Head(ap), orderBy), isAsc, map[model.FilterName]interface{}{})
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
			MetricNameQueryTimeTotal:      fmt.Sprintf("%v", utils.Round(info.Get(MetricNameQueryTimeTotal).Float()/1000/1000, 3)),
			MetricNameCPUTimeTotal:        fmt.Sprintf("%v", utils.Round(info.Get(MetricNameCPUTimeTotal).Float()/1000/1000, 3)),
			MetricNameDiskReadTotal:       strconv.Itoa(int(info.Get(MetricNameDiskReadTotal).Int())),
			MetricNameBufferGetCounter:    strconv.Itoa(int(info.Get(MetricNameBufferGetCounter).Int())),
			MetricNameUserIOWaitTimeTotal: fmt.Sprintf("%v", utils.Round(info.Get(MetricNameUserIOWaitTimeTotal).Float()/1000, 3)),
			model.AuditResultName:         sql.AuditResult.String,
		})
	}
	return rows, count, nil
}
