package auditplan

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

type MySQLProcessListTaskV2 struct {
	DefaultTaskV2
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
			Key:      paramKeyCollectIntervalSecond,
			Value:    "60",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamCollectIntervalSecond),
		},
		{
			Key:      paramKeySQLMinSecond,
			Value:    "0",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamSQLMinSecond),
		},
	}

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

	instance, exist, err := dms.GetInstancesById(ctx, ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}
	if !exist {
		return nil, errors.NewInstanceNoExistErr()
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

	// 查询 SHOW FULL PROCESSLIST
	res, err := db.Db.Query("SHOW FULL PROCESSLIST")
	if err != nil {
		return nil, fmt.Errorf("SHOW FULL PROCESSLIST failed, error: %v", err)
	}

	if len(res) <= 1 { // 仅有自己的连接
		return nil, nil
	}
	cache := NewSQLV2Cache()
	sqlMinSecond := ap.Params.GetParam(paramKeySQLMinSecond).Int()
	for i := range res {
		if sqlMinSecond > 0 {
			queryTime, _ := strconv.Atoi(res[i]["Time"].String)
			if sqlMinSecond > queryTime {
				continue
			}
		}
		if res[i]["Info"].String == "" ||
			res[i]["Id"].String == db.Db.GetConnectionID() ||
			res[i]["db"].String == "information_schema" ||
			res[i]["db"].String == "performance_schema" ||
			res[i]["db"].String == "mysql" ||
			res[i]["db"].String == "sys" {
			continue
		}
		query := res[i]["Info"].String
		sqlV2 := &SQLV2{
			Source:      ap.Type,
			SourceId:    strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
			AuditPlanId: strconv.FormatUint(uint64(ap.ID), 10),
			ProjectId:   ap.ProjectId,
			InstanceID:  ap.InstanceID,
			SchemaName:  res[i]["db"].String,
			SQLContent:  query,
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
