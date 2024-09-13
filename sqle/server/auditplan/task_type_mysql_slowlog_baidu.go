package auditplan

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/sirupsen/logrus"
)

type MySQLSlowLogBaiduTaskV2 struct {
	DefaultTaskV2
	lastEndTime *time.Time
}

func NewMySQLSlowLogBaiduTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &MySQLSlowLogBaiduTaskV2{}
	}
}

func (at *MySQLSlowLogBaiduTaskV2) InstanceType() string {
	return InstanceTypeMySQL
}

func (at *MySQLSlowLogBaiduTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:      paramKeyDBInstanceId,
			Value:    "",
			Type:     params.ParamTypeString,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamDBInstanceId),
		},
		{
			Key:      paramKeyAccessKeyId,
			Value:    "",
			Type:     params.ParamTypePassword,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamAccessKeyId),
		},
		{
			Key:      paramKeyAccessKeySecret,
			Value:    "",
			Type:     params.ParamTypePassword,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamAccessKeySecret),
		},
		{
			Key: paramKeyFirstSqlsScrappedInLastPeriodHours,
			// 百度云RDS慢日志只能拉取最近7天的数据
			// https://cloud.baidu.com/doc/RDS/s/Tjwvz046g
			Value:    "24",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAllWithArgs(locale.ParamFirstCollectDurationWithMaxDays, 7),
		},
		{
			Key:      paramKeyRdsPath,
			Value:    "rds.bj.baidubce.com",
			Type:     params.ParamTypeString,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamRdsPath),
		},
	}
}

func (at *MySQLSlowLogBaiduTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *MySQLSlowLogBaiduTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	rdsDBInstanceId := ap.Params.GetParam(paramKeyDBInstanceId).String()
	if rdsDBInstanceId == "" {
		return nil, fmt.Errorf("db instance id is not configured")
	}

	accessKeyID := ap.Params.GetParam(paramKeyAccessKeyId).String()
	if accessKeyID == "" {
		return nil, fmt.Errorf("access key id is not configured")
	}

	accessKeySecret := ap.Params.GetParam(paramKeyAccessKeySecret).String()
	if accessKeySecret == "" {
		return nil, fmt.Errorf("access key secret is not configured")
	}

	rdsPath := ap.Params.GetParam(paramKeyRdsPath).String()
	if rdsPath == "" {
		return nil, fmt.Errorf("rds path is not configured")
	}

	firstScrapInLastHours := ap.Params.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).Int()
	if firstScrapInLastHours == 0 {
		firstScrapInLastHours = 24
	}

	theMaxSupportedDays := 7 // 支持往前查看慢日志的最大天数
	hoursDuringADay := 24
	if firstScrapInLastHours > theMaxSupportedDays*hoursDuringADay {
		logger.Warnf("Can not get slow logs from so early time. firstScrapInLastHours=%v", firstScrapInLastHours)
		return nil, nil
	}

	client, err := at.CreateClient(rdsPath, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("create client failed, error: %v", err)
	}

	// use UTC time
	// https://cloud.baidu.com/doc/RDS/s/Ejwvz0uoq
	now := time.Now().UTC()
	var startTime time.Time
	if at.isFirstScrap() {
		startTime = now.Add(time.Duration(-1*firstScrapInLastHours) * time.Hour)
	} else {
		startTime = *at.lastEndTime
	}

	var pageNum int32 = 1
	pageSize := 100
	//nolint:staticcheck
	slowSqls := make([]BaiduSlowSqlRes, 0)
	for {
		// 每次取100条
		slowSqlListFromBaiduRds, err := at.pullSlowLogs(client, rdsDBInstanceId, startTime, now, int32(pageSize), pageNum)
		if err != nil {
			return nil, fmt.Errorf("pull slow logs failed, error: %v", err)
		}

		slowSqls = append(slowSqls, slowSqlListFromBaiduRds...)
		if len(slowSqlListFromBaiduRds) < pageSize {
			break
		}
		pageNum++
	}

	// 缓存采集时间点
	if len(slowSqls) > 0 {
		lastSlowSql := slowSqls[len(slowSqls)-1]
		at.lastEndTime = &lastSlowSql.executionStartTime
	}

	cache := NewSQLV2Cache()
	for _, sql := range slowSqls {
		sqlV2 := &SQLV2{
			Source:     ap.Type,
			SourceId:   strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
			ProjectId:  ap.ProjectId,
			InstanceID: ap.InstanceID,
			SchemaName: sql.schema,
			SQLContent: sql.sql,
		}
		fp, err := util.Fingerprint(sql.sql, true)
		if err != nil {
			logger.Warnf("get sql finger print failed, err: %v, sql: %s", err, sql.sql)
			fp = sql.sql
		} else if fp == "" {
			logger.Warn("get sql finger print failed, fp is empty")
			fp = sql.sql
		}
		sqlV2.Fingerprint = fp

		info := NewMetrics()
		// counter
		info.SetInt(MetricNameCounter, 1)

		// latest query time, todo: 是否可以从数据库取
		info.SetString(MetricNameLastReceiveTimestamp, time.Now().Format(time.RFC3339))

		sqlV2.Info = info
		sqlV2.GenSQLId()
		err = at.AggregateSQL(cache, sqlV2)
		if err != nil {
			logger.Warnf("aggregate sql failed, error: %v", err)
			// todo: 有错误咋处理
			continue
		}
	}

	return cache.GetSQLs(), nil
}

type BaiduSlowSqlRes struct {
	sql                string
	executionStartTime time.Time
	schema             string
}

type slowLogReq struct {
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
	PageNo    int      `json:"pageNo"`
	PageSize  int      `json:"pageSize"`
	DbName    []string `json:"dbName"`
}

type slowLogDetail struct {
	InstanceId   string  `json:"instanceId"`
	UserName     string  `json:"userName"`
	DbName       string  `json:"dbName"`
	HostIp       string  `json:"hostIp"`
	QueryTime    float64 `json:"queryTime"`
	LockTime     float64 `json:"lockTime"`
	RowsExamined int     `json:"rowsExamined"`
	RowsSent     int     `json:"rowsSent"`
	Sql          string  `json:"sql"`
	ExecuteTime  string  `json:"executeTime"`
}

type slowLogDetailResp struct {
	Count    int             `json:"count"`
	SlowLogs []slowLogDetail `json:"slowLogs"`
}

func getRdsUriWithInstanceId(instanceId string) string {
	return rds.URI_PREFIX + rds.REQUEST_RDS_URL + "/" + instanceId + "/slowlogs/details"
}

func (at *MySQLSlowLogBaiduTaskV2) CreateClient(rdsPath string, accessKeyID string, accessKeySecret string) (*rds.Client, error) {
	client, err := rds.NewClient(accessKeyID, accessKeySecret, rdsPath)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (at *MySQLSlowLogBaiduTaskV2) pullSlowLogs(client *rds.Client, instanceID string, startTime time.Time, endTime time.Time, size int32, num int32) (sqlList []BaiduSlowSqlRes, err error) {
	req := slowLogReq{
		StartTime: startTime.Format("2006-01-02T15:04:05Z"),
		EndTime:   endTime.Format("2006-01-02T15:04:05Z"),
		PageNo:    int(num),
		PageSize:  int(size),
		DbName:    nil,
	}

	uri := getRdsUriWithInstanceId(instanceID)

	resp := new(slowLogDetailResp)
	err = bce.NewRequestBuilder(client).
		WithMethod(http.MethodPost).
		WithURL(uri).
		WithBody(&req).
		WithResult(&resp).
		Do()
	if err != nil {
		return nil, err
	}

	sqlList = make([]BaiduSlowSqlRes, len(resp.SlowLogs))
	for i, slowLog := range resp.SlowLogs {
		execStartTime, err := time.Parse("2006-01-02 15:04:05", slowLog.ExecuteTime)
		if err != nil {
			return nil, err
		}

		sqlList[i] = BaiduSlowSqlRes{
			sql:                slowLog.Sql,
			executionStartTime: execStartTime,
			schema:             slowLog.DbName,
		}
	}

	return sqlList, nil
}

// 因为查询的起始时间为上一次查询到的最后一条慢语句的executionStartTime（精确到秒），而查询起始时间只能精确到分钟，所以有可能还是会查询到上一次查询过的慢语句，需要将其过滤掉
func (at *MySQLSlowLogBaiduTaskV2) filterSlowSqlsByExecutionTime(slowSqls []AliSlowSqlRes, executionTime time.Time) (res []AliSlowSqlRes) {
	for _, sql := range slowSqls {
		if !sql.executionStartTime.After(executionTime) {
			continue
		}
		res = append(res, sql)
	}
	return
}

func (at *MySQLSlowLogBaiduTaskV2) isFirstScrap() bool {
	return at.lastEndTime == nil
}
