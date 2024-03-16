package auditplan

import (
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"

	rdsCoreRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	rds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	rdsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	rdsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/region"
	"github.com/percona/go-mysql/query"
	"github.com/sirupsen/logrus"
)

const huaweiCloudRequestTimeFormat = "2006-01-02T15:04:05-0700"
const huaweiCloudResponseTimeFormat = "2006-01-02T15:04:05"

type SqlFromHuaweiCloud struct {
	sql                string
	executionStartTime time.Time
	schema             string
}

func NewHuaweiRdsMySQLSlowLogTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	sqlCollector := newSQLCollector(entry, ap)
	a := &HuaweiRdsMySQLSlowLogTask{}
	task := &huaweiRdsMySQLTask{
		sqlCollector: sqlCollector,
		lastEndTime:  nil,
		pullLogs:     a.pullLogs,
	}
	sqlCollector.do = task.collectorDo
	a.huaweiRdsMySQLTask = task
	return a
}

type HuaweiRdsMySQLSlowLogTask struct {
	*huaweiRdsMySQLTask
}

func (hr *HuaweiRdsMySQLSlowLogTask) pullLogs(client *rds.RdsClient, instanceId string, startTime, endTime time.Time, pageSize, pageNum int32) (sqls []SqlFromHuaweiCloud, err error) {
	request := hr.newSlowSqlsRequest(instanceId, startTime, endTime, pageSize, pageNum)
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

	sqls = make([]SqlFromHuaweiCloud, len(*response.SlowLogList))
	for i, slowRecord := range *response.SlowLogList {
		execStartTime, err := time.Parse(huaweiCloudResponseTimeFormat, utils.NvlString(&slowRecord.StartTime))
		if err != nil {
			return nil, fmt.Errorf("parse huawei cloud execution-start-time failed: %v", err)
		}
		sqls[i] = SqlFromHuaweiCloud{
			sql:                utils.NvlString(&slowRecord.QuerySample),
			executionStartTime: execStartTime,
			schema:             slowRecord.Database,
		}
	}
	return sqls, nil
}

func (at *HuaweiRdsMySQLSlowLogTask) Audit() (*AuditResultResp, error) {
	task := &model.Task{
		DBType: at.ap.DBType,
	}
	return at.baseTask.audit(task)
}

// huaweiRdsMySQLTask implement the Task interface.
//
// huaweiRdsMySQLTask is a loop task which collect slow log from ali rds MySQL instance.
type huaweiRdsMySQLTask struct {
	*sqlCollector
	lastEndTime *time.Time
	pullLogs    func(client *rds.RdsClient, instanceId string, startTime, endTime time.Time, pageSize, pageNum int32) (sqls []SqlFromHuaweiCloud, err error)
}

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

func (at *huaweiRdsMySQLTask) collectorDo() {
	//1. Load Configuration
	// if at.ap.InstanceName == "" {
	// 	at.logger.Warnf("instance is not configured")
	// 	return
	// }
	accessKeyId := at.ap.Params.GetParam(paramKeyAccessKeyId).String()
	if accessKeyId == "" {
		at.logger.Warnf("huawei cloud access key id is not configured")
		return
	}
	secretAccessKey := at.ap.Params.GetParam(paramKeyAccessKeySecret).String()
	if secretAccessKey == "" {
		at.logger.Warnf("huawei cloud secret access key is not configured")
		return
	}
	projectId := at.ap.Params.GetParam(paramKeyProjectId).String()
	if projectId == "" {
		at.logger.Warnf("huawei cloud project id is not configured")
		return
	}
	instanceId := at.ap.Params.GetParam(paramKeyDBInstanceId).String()
	if instanceId == "" {
		at.logger.Warnf("huawei cloud instance id is not configured")
		return
	}
	region := at.ap.Params.GetParam(paramKeyRegion).String()
	if region == "" {
		at.logger.Warnf("huawei cloud region is not configured")
		return
	}
	if !isHuaweiRdsRegionExist(region) {
		at.logger.Warnf("huawei cloud region is not exist")
		return
	}
	periodHours := at.ap.Params.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).String()
	if periodHours == "" {
		at.logger.Warnf("huawei cloud period hours is not configured")
		return
	}
	firstScrapInLastHours, err := strconv.Atoi(periodHours)
	if err != nil {
		at.sqlCollector.logger.Warnf("convert first sqls scrapped in last period hours failed: %v", err)
		return
	}

	if firstScrapInLastHours == 0 {
		firstScrapInLastHours = 24
	}
	theMaxSupportedDays := 30 // 支持往前查看慢日志的最大天数
	hoursDuringADay := 24
	if firstScrapInLastHours > theMaxSupportedDays*hoursDuringADay {
		at.logger.Warnf("Can not get slow logs from so early time. firstScrapInLastHours=%v", firstScrapInLastHours)
		return
	}
	//2. Init Client
	client := at.CreateClient(accessKeyId, secretAccessKey, projectId, region)
	if client == nil {
		at.logger.Warnf("create client for huawei rdb mysql failed")
		return
	}

	// 3. Request for slow logs

	now := time.Now()
	slowSqls := []SqlFromHuaweiCloud{}
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
			at.logger.Warnf("pull rds logs failed: %v", err)
			return
		}
		filteredNewSlowSqls := at.filterSlowSqlsByExecutionTime(newSlowSqls, startTime)
		slowSqls = append(slowSqls, filteredNewSlowSqls...)

		if len(newSlowSqls) < int(pageSize) {
			break
		}
		pageNum++
	}

	//4. Merge SQL
	mergedSlowSqls := mergeSQLsFromHuaweiCloud(slowSqls)
	if len(mergedSlowSqls) > 0 {
		if at.isFirstScrap() {
			err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, at.convertSQLInfosToModelSQLs(mergedSlowSqls, now))
			if err != nil {
				at.logger.Errorf("save sqls to storage fail, error: %v", err)
				return
			}
		} else {
			err = at.persist.UpdateDefaultAuditPlanSQLs(at.ap.ID, at.convertSQLInfosToModelSQLs(mergedSlowSqls, now))
			if err != nil {
				at.logger.Errorf("save sqls to storage fail, error: %v", err)
				return
			}
		}
	}
	// update lastEndTime
	// 查询的起始时间为上一次查询到的最后一条慢语句的开始执行时间
	if len(slowSqls) > 0 {
		lastSlowSql := slowSqls[len(slowSqls)-1]
		at.lastEndTime = &lastSlowSql.executionStartTime
	}

}

func (at *huaweiRdsMySQLTask) convertSQLInfosToModelSQLs(sqls []sqlInfo, now time.Time) []*model.AuditPlanSQLV2 {
	return convertRawSlowSQLWitchFromSqlInfo(sqls, now)
}

func mergeSQLsFromHuaweiCloud(sqls []SqlFromHuaweiCloud) []sqlInfo {
	sqlInfos := []sqlInfo{}

	counter := map[string]int /*slice subscript*/ {}
	for _, sql := range sqls {
		fp := query.Fingerprint(sql.sql)
		if index, exist := counter[fp]; exist {
			sqlInfos[index].counter += 1
			sqlInfos[index].fingerprint = fp
			sqlInfos[index].sql = sql.sql
			sqlInfos[index].schema = sql.schema
		} else {
			sqlInfos = append(sqlInfos, sqlInfo{
				counter:     1,
				fingerprint: fp,
				sql:         sql.sql,
				schema:      sql.schema,
			})
			counter[fp] = len(sqlInfos) - 1
		}

	}
	return sqlInfos
}
func (at *huaweiRdsMySQLTask) CreateClient(accessKeyId, secretAccessKey, projectId, region string) *rds.RdsClient {
	defer func() {
		if err := recover(); err != nil {
			at.logger.Warnf("in huawei rds of audit plan, recovered from panic: %+v", err)
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

func (at *huaweiRdsMySQLTask) newSlowSqlsRequest(instanceId string, startTime, endTime time.Time, pageSize, pageNum int32) *rdsModel.ListSlowLogsRequest {
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

func (at *huaweiRdsMySQLTask) isFirstScrap() bool {
	return at.lastEndTime == nil
}

// 因为查询的起始时间为上一次查询到的最后一条慢语句的executionStartTime（精确到秒），而查询起始时间只能精确到分钟，所以有可能还是会查询到上一次查询过的慢语句，需要将其过滤掉
func (at *huaweiRdsMySQLTask) filterSlowSqlsByExecutionTime(slowSqls []SqlFromHuaweiCloud, executionTime time.Time) (res []SqlFromHuaweiCloud) {
	for _, sql := range slowSqls {
		if !sql.executionStartTime.After(executionTime) {
			continue
		}
		res = append(res, sql)
	}
	return
}
