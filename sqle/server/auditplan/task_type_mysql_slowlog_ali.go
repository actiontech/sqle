package auditplan

import (
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	rds20140815 "github.com/alibabacloud-go/rds-20140815/v2/client"
	_util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
)

type MySQLSlowLogAliTaskV2 struct {
	DefaultTaskV2
	lastEndTime *time.Time
}

func NewMySQLSlowLogAliTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &MySQLSlowLogAliTaskV2{}
	}
}

func (at *MySQLSlowLogAliTaskV2) InstanceType() string {
	return InstanceTypeMySQL
}

func (at *MySQLSlowLogAliTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:   paramKeyDBInstanceId,
			Desc:  "实例ID",
			Value: "",
			Type:  params.ParamTypeString,
		},
		{
			Key:   paramKeyAccessKeyId,
			Desc:  "Access Key ID",
			Value: "",
			Type:  params.ParamTypePassword,
		},
		{
			Key:   paramKeyAccessKeySecret,
			Desc:  "Access Key Secret",
			Value: "",
			Type:  params.ParamTypePassword,
		},
		{
			Key:   paramKeyFirstSqlsScrappedInLastPeriodHours,
			Desc:  "启动任务时拉取慢日志时间范围(单位:小时,最大31天)",
			Value: "",
			Type:  params.ParamTypeInt,
		},
		{
			Key:   paramKeyRdsPath,
			Desc:  "RDS Open API地址",
			Value: "rds.aliyuncs.com",
			Type:  params.ParamTypeString,
		},
	}
}

func (at *MySQLSlowLogAliTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{
		defaultAuditLevelOperateParams,
	}
}

func (at *MySQLSlowLogAliTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameLastReceiveTimestamp,
	}
}

func (at *MySQLSlowLogAliTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *MySQLSlowLogAliTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	rdsDBInstanceId := ap.Params.GetParam(paramKeyDBInstanceId).String()
	if rdsDBInstanceId == "" {
		return nil, fmt.Errorf("rds DB instance ID is not configured")
	}

	accessKeyId := ap.Params.GetParam(paramKeyAccessKeyId).String()
	if accessKeyId == "" {
		return nil, fmt.Errorf("Access Key ID is not configured")
	}

	accessKeySecret := ap.Params.GetParam(paramKeyAccessKeySecret).String()
	if accessKeySecret == "" {
		return nil, fmt.Errorf("Access Key Secret is not configured")
	}

	rdsPath := ap.Params.GetParam(paramKeyRdsPath).String()
	if rdsPath == "" {
		return nil, fmt.Errorf("RDS Path is not configured")
	}

	firstScrapInLastHours := ap.Params.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).Int()
	if firstScrapInLastHours == 0 {
		firstScrapInLastHours = 24
	}
	theMaxSupportedDays := 31 // 支持往前查看慢日志的最大天数
	hoursDuringADay := 24
	if firstScrapInLastHours > theMaxSupportedDays*hoursDuringADay {
		logger.Warnf("Can not get slow logs from so early time. firstScrapInLastHours=%v", firstScrapInLastHours)
		return nil, nil
	}

	client, err := at.CreateClient(rdsPath, tea.String(accessKeyId), tea.String(accessKeySecret))
	if err != nil {
		return nil, fmt.Errorf("create client for polardb mysql failed: %v", err)
	}

	pageSize := 100
	now := time.Now().UTC()
	var startTime time.Time
	if at.isFirstScrap() {
		startTime = now.Add(time.Duration(-1*firstScrapInLastHours) * time.Hour)
	} else {
		startTime = *at.lastEndTime
	}
	var pageNum int32 = 1
	slowSqls := []AliSlowSqlRes{}
	for {
		newSlowSqls, err := at.pullSlowLogs(client, rdsDBInstanceId, startTime, now, int32(pageSize), pageNum)
		if err != nil {
			return nil, fmt.Errorf("pull rds logs failed: %v", err)
		}
		filteredNewSlowSqls := at.filterSlowSqlsByExecutionTime(newSlowSqls, startTime)
		slowSqls = append(slowSqls, filteredNewSlowSqls...)

		if len(newSlowSqls) < pageSize {
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
			SourceId:   ap.ID,
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
			logger.Warn("get sql finger print failed, fp is empty", err, sql.sql)
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

type AliSlowSqlRes struct {
	sql                string
	executionStartTime time.Time
	schema             string
}

func (at *MySQLSlowLogAliTaskV2) CreateClient(rdsPath string, accessKeyId *string, accessKeySecret *string) (_result *rds20140815.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(rdsPath)
	_result, _err = rds20140815.NewClient(config)
	return _result, _err
}

// 查询内容范围是开始时间的0s到设置结束时间的0s，所以结束时间点的慢日志是查询不到的
// startTime和endTime对应的是慢语句的开始执行时间
func (at *MySQLSlowLogAliTaskV2) pullSlowLogs(client *rds20140815.Client, DBInstanId string, startTime, endTime time.Time, pageSize, pageNum int32) (sqls []AliSlowSqlRes, err error) {
	describeSlowLogRecordsRequest := &rds20140815.DescribeSlowLogRecordsRequest{
		DBInstanceId: tea.String(DBInstanId),
		StartTime:    tea.String(startTime.Format("2006-01-02T15:04Z")),
		EndTime:      tea.String(endTime.Format("2006-01-02T15:04Z")),
		PageSize:     tea.Int32(pageSize),
		PageNumber:   tea.Int32(pageNum),
	}

	runtime := &_util.RuntimeOptions{}
	response := &rds20140815.DescribeSlowLogRecordsResponse{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		var err error
		response, err = client.DescribeSlowLogRecordsWithOptions(describeSlowLogRecordsRequest, runtime)
		if err != nil {
			return err
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		errMsg := _util.AssertAsString(error.Message)
		return nil, fmt.Errorf("get slow log failed: %v", *errMsg)
	}

	sqls = make([]AliSlowSqlRes, len(response.Body.Items.SQLSlowRecord))
	for i, slowRecord := range response.Body.Items.SQLSlowRecord {
		execStartTime, err := time.Parse("2006-01-02T15:04:05Z", utils.NvlString(slowRecord.ExecutionStartTime))
		if err != nil {
			return nil, fmt.Errorf("parse execution-start-time failed: %v", err)
		}
		sqls[i] = AliSlowSqlRes{
			sql:                utils.NvlString(slowRecord.SQLText),
			executionStartTime: execStartTime,
			schema:             utils.NvlString(slowRecord.DBName),
		}
	}
	return sqls, nil
}

// 因为查询的起始时间为上一次查询到的最后一条慢语句的executionStartTime（精确到秒），而查询起始时间只能精确到分钟，所以有可能还是会查询到上一次查询过的慢语句，需要将其过滤掉
func (at *MySQLSlowLogAliTaskV2) filterSlowSqlsByExecutionTime(slowSqls []AliSlowSqlRes, executionTime time.Time) (res []AliSlowSqlRes) {
	for _, sql := range slowSqls {
		if !sql.executionStartTime.After(executionTime) {
			continue
		}
		res = append(res, sql)
	}
	return
}

func (at *MySQLSlowLogAliTaskV2) isFirstScrap() bool {
	return at.lastEndTime == nil
}
