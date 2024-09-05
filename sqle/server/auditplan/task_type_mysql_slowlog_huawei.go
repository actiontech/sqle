package auditplan

import (
	"fmt"
	"html"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"

	rdsCoreRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	rds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	rdsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	rdsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/region"
	"github.com/sirupsen/logrus"
)

type MySQLSlowLogHuaweiTaskV2 struct {
	DefaultTaskV2
	lastEndTime *time.Time
}

func NewMySQLSlowLogHuaweiTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &MySQLSlowLogHuaweiTaskV2{}
	}
}

func (at *MySQLSlowLogHuaweiTaskV2) InstanceType() string {
	return InstanceTypeMySQL
}

func (at *MySQLSlowLogHuaweiTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:      paramKeyProjectId,
			Value:    "",
			Type:     params.ParamTypeString,
			I18nDesc: locale.ShouldLocalizeAll(locale.ParamProjectId),
		},
		{
			Key:      paramKeyDBInstanceId,
			Value:    "",
			Type:     params.ParamTypeString,
			I18nDesc: locale.ShouldLocalizeAll(locale.ParamDBInstanceId),
		},
		{
			Key:      paramKeyAccessKeyId,
			Value:    "",
			Type:     params.ParamTypePassword,
			I18nDesc: locale.ShouldLocalizeAll(locale.ParamAccessKeyId),
		},
		{
			Key:      paramKeyAccessKeySecret,
			Value:    "",
			Type:     params.ParamTypePassword,
			I18nDesc: locale.ShouldLocalizeAll(locale.ParamAccessKeySecret),
		},
		{
			Key:      paramKeyFirstSqlsScrappedInLastPeriodHours,
			Value:    "24",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.ShouldLocalizeAllWithArgs(locale.ParamFirstCollectDurationWithMaxDays, 30),
		},
		{
			Key:      paramKeyRegion,
			Value:    "",
			Type:     params.ParamTypeString,
			I18nDesc: locale.ShouldLocalizeAll(locale.ParamRegion),
		},
	}
}

func (at *MySQLSlowLogHuaweiTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *MySQLSlowLogHuaweiTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	//1. Load Configuration
	// if at.ap.InstanceName == "" {
	// 	at.logger.Warnf("instance is not configured")
	// 	return
	// }
	accessKeyId := ap.Params.GetParam(paramKeyAccessKeyId).String()
	if accessKeyId == "" {
		return nil, fmt.Errorf("huawei cloud access key id is not configured")
	}
	secretAccessKey := ap.Params.GetParam(paramKeyAccessKeySecret).String()
	if secretAccessKey == "" {
		return nil, fmt.Errorf("huawei cloud secret access key is not configured")
	}
	projectId := ap.Params.GetParam(paramKeyProjectId).String()
	if projectId == "" {
		return nil, fmt.Errorf("huawei cloud project id is not configured")
	}
	instanceId := ap.Params.GetParam(paramKeyDBInstanceId).String()
	if instanceId == "" {
		return nil, fmt.Errorf("huawei cloud instance id is not configured")
	}
	region := ap.Params.GetParam(paramKeyRegion).String()
	if region == "" {
		return nil, fmt.Errorf("huawei cloud region is not configured")
	}
	if !isHuaweiRdsRegionExist(region) {
		return nil, fmt.Errorf("huawei cloud region is not exist")
	}
	periodHours := ap.Params.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).String()
	if periodHours == "" {
		return nil, fmt.Errorf("huawei cloud period hours is not configured")
	}
	firstScrapInLastHours, err := strconv.Atoi(periodHours)
	if err != nil {
		return nil, fmt.Errorf("convert first sqls scrapped in last period hours failed: %v", err)
	}

	if firstScrapInLastHours == 0 {
		firstScrapInLastHours = 24
	}
	theMaxSupportedDays := 30 // 支持往前查看慢日志的最大天数
	hoursDuringADay := 24
	if firstScrapInLastHours > theMaxSupportedDays*hoursDuringADay {
		logger.Warnf("Can not get slow logs from so early time. firstScrapInLastHours=%v", firstScrapInLastHours)
		return nil, nil
	}
	//2. Init Client
	client := at.CreateClient(logger, accessKeyId, secretAccessKey, projectId, region)
	if client == nil {
		return nil, fmt.Errorf("create client for huawei rdb mysql failed")
	}

	// 3. Request for slow logs

	now := time.Now()
	slowSqls := []HuaweiSlowSqlRes{}
	var pageNum int32 = 1
	var pageSize int32 = 100
	var startTime time.Time
	if at.isFirstScrap() {
		startTime = now.Add(time.Duration(-1*firstScrapInLastHours) * time.Hour)
	} else {
		startTime = *at.lastEndTime
	}

	for {
		newSlowSqls, err := at.pullLogs(client, instanceId, startTime, now, pageSize, pageNum)
		if err != nil {
			return nil, fmt.Errorf("pull rds logs failed: %v", err)
		}
		filteredNewSlowSqls := at.filterSlowSqlsByExecutionTime(newSlowSqls, startTime)
		slowSqls = append(slowSqls, filteredNewSlowSqls...)

		if len(newSlowSqls) < int(pageSize) {
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

type HuaweiSlowSqlRes struct {
	sql                string
	executionStartTime time.Time
	schema             string
}

const huaweiCloudRequestTimeFormat = "2006-01-02T15:04:05-0700"
const huaweiCloudResponseTimeFormat = "2006-01-02T15:04:05"

var staticFields = map[string]struct{}{
	"af-south-1":     {},
	"cn-north-4":     {},
	"cn-north-1":     {},
	"cn-east-2":      {},
	"cn-east-3":      {},
	"cn-south-1":     {},
	"cn-southwest-2": {},
	"ap-southeast-2": {},
	"ap-southeast-1": {},
	"ap-southeast-3": {},
	"ru-northwest-2": {},
	"sa-brazil-1":    {},
	"la-north-2":     {},
	"cn-south-2":     {},
	"na-mexico-1":    {},
	"la-south-2":     {},
	"cn-north-9":     {},
	"cn-north-2":     {},
	"tr-west-1":      {},
	"ap-southeast-4": {},
	"ae-ad-1":        {},
	"eu-west-101":    {},
}

func isHuaweiRdsRegionExist(region string) bool {
	var provider = rdsCoreRegion.DefaultProviderChain("RDS")
	reg := provider.GetRegion(region)
	if reg == nil {
		_, exist := staticFields[region]
		return exist
	}
	return true
}

func (at *MySQLSlowLogHuaweiTaskV2) newSlowSqlsRequest(instanceId string, startTime, endTime time.Time, pageSize, pageNum int32) *rdsModel.ListSlowLogsRequest {
	startDate := startTime.Format(huaweiCloudRequestTimeFormat)
	endDate := endTime.Format(huaweiCloudRequestTimeFormat)

	return &rdsModel.ListSlowLogsRequest{
		InstanceId: instanceId,
		StartDate:  startDate,
		EndDate:    endDate,
		Offset:     &pageNum,
		Limit:      &pageSize,
	}
}

func (at *MySQLSlowLogHuaweiTaskV2) CreateClient(logger *logrus.Entry, accessKeyId, secretAccessKey, projectId, region string) *rds.RdsClient {
	defer func() {
		if err := recover(); err != nil {
			logger.Warnf("in huawei rds of audit plan, recovered from panic: %+v", err)
		}
	}()
	credential := basic.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		WithProjectId(projectId).
		Build()
	hcClient := rds.RdsClientBuilder().
		WithRegion(rdsRegion.ValueOf(region)).
		WithCredential(credential).
		Build()
	return rds.NewRdsClient(hcClient)
}

func (at *MySQLSlowLogHuaweiTaskV2) pullLogs(client *rds.RdsClient, instanceId string, startTime, endTime time.Time, pageSize, pageNum int32) (sqls []HuaweiSlowSqlRes, err error) {
	request := at.newSlowSqlsRequest(instanceId, startTime, endTime, pageSize, pageNum)
	response := &rdsModel.ListSlowLogsResponse{}

	tryErr := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("panic: %v", r)
				}
			}
		}()
		response, err = client.ListSlowLogs(request)
		if err != nil {
			return err
		}
		return nil
	}()

	if tryErr != nil {
		return nil, fmt.Errorf("get huawei cloud slow log failed: %+v", tryErr)
	}

	sqls = make([]HuaweiSlowSqlRes, len(*response.SlowLogList))
	for i, slowRecord := range *response.SlowLogList {
		execStartTime, err := time.Parse(huaweiCloudResponseTimeFormat, utils.NvlString(&slowRecord.StartTime))
		if err != nil {
			return nil, fmt.Errorf("parse huawei cloud execution-start-time failed: %v", err)
		}
		sqls[i] = HuaweiSlowSqlRes{
			sql:                html.UnescapeString(utils.NvlString(&slowRecord.QuerySample)),
			executionStartTime: execStartTime,
			schema:             slowRecord.Database,
		}
	}
	return sqls, nil
}

func (at *MySQLSlowLogHuaweiTaskV2) isFirstScrap() bool {
	return at.lastEndTime == nil
}

// 因为查询的起始时间为上一次查询到的最后一条慢语句的executionStartTime（精确到秒），而查询起始时间只能精确到分钟，所以有可能还是会查询到上一次查询过的慢语句，需要将其过滤掉
func (at *MySQLSlowLogHuaweiTaskV2) filterSlowSqlsByExecutionTime(slowSqls []HuaweiSlowSqlRes, executionTime time.Time) (res []HuaweiSlowSqlRes) {
	for _, sql := range slowSqls {
		if !sql.executionStartTime.After(executionTime) {
			continue
		}
		res = append(res, sql)
	}
	return
}
