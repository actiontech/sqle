package auditplan

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/oracle"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/sirupsen/logrus"
)

type OracleTopSQLTaskV2 struct{ DefaultTaskV2 }

func NewOracleTopSQLTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &OracleTopSQLTaskV2{
			DefaultTaskV2: DefaultTaskV2{},
		}
	}
}

func (at *OracleTopSQLTaskV2) InstanceType() string {
	return InstanceTypeOracle
}

func (at *OracleTopSQLTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:      paramKeyCollectIntervalMinute,
			Value:    "60",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamCollectIntervalMinuteOracle),
		},
	}
}

func (at *OracleTopSQLTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameLastReceiveTimestamp,
		MetricNameQueryTimeAvg,
		MetricNameQueryTimeTotal,
		MetricNameUserIOWaitTimeTotal,
		MetricNameCPUTimeTotal,
		MetricNameDiskReadTotal,
		MetricNameBufferGetCounter,
		MetricNameDBUser,
	}
}

func (at *OracleTopSQLTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}
	// counter
	originSQL.Info.SetInt(MetricNameCounter, mergedSQL.Info.Get(MetricNameCounter).Int())

	// MetricNameQueryTimeTotal
	originSQL.Info.SetFloat(MetricNameQueryTimeTotal, mergedSQL.Info.Get(MetricNameQueryTimeTotal).Float())

	// MetricNameUserIOWaitTimeTotal
	originSQL.Info.SetFloat(MetricNameUserIOWaitTimeTotal, mergedSQL.Info.Get(MetricNameUserIOWaitTimeTotal).Float())

	// MetricNameCPUTimeTotal
	originSQL.Info.SetFloat(MetricNameCPUTimeTotal, mergedSQL.Info.Get(MetricNameCPUTimeTotal).Float())

	// MetricNameDiskReadTotal
	originSQL.Info.SetInt(MetricNameDiskReadTotal, mergedSQL.Info.Get(MetricNameDiskReadTotal).Int())

	// MetricNameBufferGetCounter
	originSQL.Info.SetInt(MetricNameBufferGetCounter, mergedSQL.Info.Get(MetricNameBufferGetCounter).Int())
	// MetricNameDBUser
	originSQL.Info.SetString(MetricNameDBUser, mergedSQL.Info.Get(MetricNameDBUser).String())
}

func (at *OracleTopSQLTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	inst, exist, err := dms.GetInstancesById(ctx, ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}
	if !exist {
		return nil, errors.NewInstanceNoExistErr()
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
	// get db user blacklist
	dbUserBlacklists, err := model.GetStorage().
		GetBlacklistByProjectIDAndFilterType(model.ProjectUID(ap.ProjectId), model.FilterTypeDbUser)
	if err != nil {
		return nil, fmt.Errorf("get blacklist fail, error: %v", err)
	}
	// convert to string slice
	notInUser := make([]string, 0, len(dbUserBlacklists))
	for _, blacklist := range dbUserBlacklists {
		notInUser = append(notInUser, blacklist.FilterContent)
	}
	// NOTE: top_n and order_by_column are not defined in Params(), GetParam returns zero values.
	// top_n=0 defaults to 10 in QueryTopSQLs, order_by_column="" iterates multiple metrics.
	// This is a known issue, kept as-is to avoid introducing risk. See design doc risk item 8.
	sqls, err := db.QueryTopSQLs(ctx, ap.Params.GetParam("collect_interval_minute").String(), ap.Params.GetParam("top_n").Int(), notInUser, ap.Params.GetParam("order_by_column").String())
	if err != nil {
		return nil, fmt.Errorf("query top sql fail, error: %v", err)
	}

	cache := NewSQLV2Cache()
	rawSQLs := make([]*model.SQLManageRawSQL, 0, len(sqls))
	for _, sql := range sqls {
		info := NewMetrics()
		sqlV2 := &SQLV2{
			Source:      ap.Type,
			SourceId:    strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
			AuditPlanId: strconv.FormatUint(uint64(ap.ID), 10),
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
		info.SetInt(MetricNameDiskReadTotal, sql.DiskReads)
		info.SetInt(MetricNameBufferGetCounter, sql.BufferGets)
		info.SetString(MetricNameDBUser, sql.UserName)
		info.SetString(MetricNameLastReceiveTimestamp, time.Now().Format(time.RFC3339))
		if sql.Executions > 0 {
			avgQueryTime := float64(sql.ElapsedTime) / float64(sql.Executions) / 1e6 // 微秒转秒
			info.SetFloat(MetricNameQueryTimeAvg, avgQueryTime)
		}
		sqlV2.GenSQLId()
		rawSQLs = append(rawSQLs, ConvertSQLV2ToMangerRawSQL(sqlV2))
		err = at.AggregateSQL(cache, sqlV2)
		if err != nil {
			logger.Warnf("aggregate sql failed,error : %v", err)
			continue
		}
	}

	if err := persist.CreateSqlManageRawSQLs(rawSQLs); err != nil {
		logger.Errorf("OracleTopSQLTaskV2 create sql manage raw sql failed, error: %v", err)
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
			Name:     MetricNameQueryTimeTotal,
			Desc:     locale.ApMetricNameQueryTimeTotal,
			Sortable: true,
		},
		{
			Name:     MetricNameCPUTimeTotal,
			Desc:     locale.ApMetricNameCPUTimeTotal,
			Sortable: true,
		},
		{
			Name:     MetricNameDiskReadTotal,
			Desc:     locale.ApMetricNameDiskReadTotal,
			Sortable: true,
		},
		{
			Name:     MetricNameBufferGetCounter,
			Desc:     locale.ApMetricNameBufferGetCounter,
			Sortable: true,
		},
		{
			Name:     MetricNameUserIOWaitTimeTotal,
			Desc:     locale.ApMetricNameUserIOWaitTimeTotal,
			Sortable: true,
		},
		{
			Name: MetricNameDBUser,
			Desc: locale.ApMetricNameDBUser,
		},
	}
}

func (at *OracleTopSQLTaskV2) GetSQLData(ctx context.Context, ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
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
			MetricNameQueryTimeTotal:      fmt.Sprintf("%v", utils.Round(info.Get(MetricNameQueryTimeTotal).Float()/1000/1000, 3)),
			MetricNameCPUTimeTotal:        fmt.Sprintf("%v", utils.Round(info.Get(MetricNameCPUTimeTotal).Float()/1000/1000, 3)),
			MetricNameDiskReadTotal:       strconv.Itoa(int(info.Get(MetricNameDiskReadTotal).Int())),
			MetricNameBufferGetCounter:    strconv.Itoa(int(info.Get(MetricNameBufferGetCounter).Int())),
			MetricNameUserIOWaitTimeTotal: fmt.Sprintf("%v", utils.Round(info.Get(MetricNameUserIOWaitTimeTotal).Float()/1000, 3)),
			model.AuditResultName:         sql.AuditResult.GetAuditJsonStrByLangTag(locale.Bundle.GetLangTagFromCtx(ctx)),
			model.AuditStatus:             sql.AuditStatus,
			MetricNameDBUser:              info.Get(MetricNameDBUser).String(),
		})
	}
	return rows, count, nil
}
