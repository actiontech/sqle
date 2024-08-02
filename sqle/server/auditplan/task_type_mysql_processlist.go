package auditplan

import (
	"context"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

type MySQLProcessListTaskV2 struct {
}

func NewMySQLProcessListTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &MySQLProcessListTaskV2{}
	}
}

func (at *MySQLProcessListTaskV2) InstanceType() string {
	return InstanceTypeMySQL
}

func (at *MySQLProcessListTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:   paramKeyCollectIntervalSecond,
			Desc:  "采集周期（秒）",
			Value: "60",
			Type:  params.ParamTypeInt,
		},
		{
			Key:   paramKeySQLMinSecond,
			Desc:  "SQL 最小执行时间（秒）",
			Value: "0",
			Type:  params.ParamTypeInt,
		},
	}
}

func (at *MySQLProcessListTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameLastReceiveTimestamp,
	}
}

func (at *MySQLProcessListTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
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
	return
}

func (at *MySQLProcessListTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *MySQLProcessListTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *MySQLProcessListTaskV2) processListSQL(ap *AuditPlan) string {
	sql := `
SELECT DISTINCT db,time,info
FROM information_schema.processlist
WHERE ID != connection_id() AND info != '' AND db NOT IN ('information_schema','performance_schema','mysql','sys')
%v
`
	whereSqlMinSecond := ""
	{
		sqlMinSecond := ap.Params.GetParam(paramKeySQLMinSecond).Int()
		if sqlMinSecond > 0 {
			whereSqlMinSecond = fmt.Sprintf("AND TIME > %d", sqlMinSecond)
		}
	}
	sql = fmt.Sprintf(sql, whereSqlMinSecond)
	return sql
}

func (at *MySQLProcessListTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	instance, _, err := dms.GetInstancesById(ctx, ap.InstanceID)
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

	res, err := db.Db.Query(at.processListSQL(ap))
	if err != nil {
		return nil, fmt.Errorf("show processlist failed, error: %v", err)
	}

	if len(res) == 0 {
		return nil, nil
	}
	cache := NewSQLV2Cache()
	for i := range res {
		query := res[i]["info"].String
		sqlV2 := &SQLV2{
			Source:     ap.Type,
			SourceId:   ap.ID,
			ProjectId:  ap.ProjectId,
			InstanceID: ap.InstanceID,
			SchemaName: res[i]["db"].String,
			SQLContent: query,
		}
		fp, err := util.Fingerprint(query, true)
		if err != nil {
			logger.Warnf("get sql finger print failed, err: %v, sql: %s", err, query)
			fp = query
		} else if fp == "" {
			logger.Warn("get sql finger print failed, fp is empty")
			fp = query
		}
		sqlV2.Fingerprint = fp

		info := NewMetrics()
		// counter
		info.SetInt(MetricNameCounter, 1)

		// latest query time, todo: 是否可以从数据库取
		info.SetString(MetricNameLastReceiveTimestamp, time.Now().Format(time.RFC3339))

		sqlV2.Info = info
		sqlV2.GenSQLId()
		if err = at.AggregateSQL(cache, sqlV2); err != nil {
			logger.Warnf("aggregate sql failed error : %v", err)
			continue
		}
	}
	return cache.GetSQLs(), nil
}

func (at *MySQLProcessListTaskV2) GetSQLs(ap *AuditPlan, persist *model.Storage, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	return baseTaskGetSQLs(args, persist)
}
