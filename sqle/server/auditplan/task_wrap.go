package auditplan

import (
	"context"
	e "errors"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pkg/errors"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/sirupsen/logrus"
)

// todo: 补充单测
type SQLV2Cacher interface {
	GetSQLs() []*SQLV2
	GetSQL(sqlId string) (*SQLV2, bool, error)
	CacheSQL(sql *SQLV2)
}

type SQLV2Cache struct {
	sqls   []*SQLV2
	sqlMap map[string]*SQLV2
}

func NewSQLV2Cache() *SQLV2Cache {
	return &SQLV2Cache{
		sqls:   []*SQLV2{},
		sqlMap: map[string]*SQLV2{},
	}
}

func (c *SQLV2Cache) GetSQLs() []*SQLV2 {
	return c.sqls
}

func (c *SQLV2Cache) GetSQL(sqlId string) (*SQLV2, bool, error) {
	if sql, ok := c.sqlMap[sqlId]; ok {
		return sql, true, nil
	} else {
		return nil, false, nil
	}
}

func (c *SQLV2Cache) CacheSQL(sql *SQLV2) {
	c.sqls = append(c.sqls, sql)
	c.sqlMap[sql.SQLId] = sql
}

type SQLV2CacheWithPersist struct {
	SQLV2Cache
	persist *model.Storage
}

func NewSQLV2CacheWithPersist(persist *model.Storage) *SQLV2CacheWithPersist {
	return &SQLV2CacheWithPersist{
		SQLV2Cache: *NewSQLV2Cache(),
		persist:    persist,
	}
}

func (c *SQLV2CacheWithPersist) GetSQL(sqlId string) (*SQLV2, bool, error) {
	if sql, ok := c.sqlMap[sqlId]; ok {
		return sql, true, nil
	}
	originSQL, exist, err := c.persist.GetManageSQLBySQLId(sqlId)
	if err != nil {
		return nil, false, err
	}
	if exist {
		originSQLV2 := ConvertMangerSQLToSQLV2(originSQL)
		c.sqls = append(c.sqls, originSQLV2)
		c.sqlMap[sqlId] = originSQLV2
		return originSQLV2, true, nil
	}
	return nil, false, nil
}

type TaskWrapper struct {
	ap *AuditPlan
	// persist is a database handle which store AuditPlan.
	persist *model.Storage
	logger  *logrus.Entry

	// 回调
	collect AuditPlanCollector
	handler AuditPlanHandler

	// collect
	sync.WaitGroup
	isStarted    bool
	cancel       chan struct{}
	loopInterval func() time.Duration
}

func NewTaskWrap(taskFn func() interface{}) func(entry *logrus.Entry, ap *AuditPlan) Task {
	return func(entry *logrus.Entry, ap *AuditPlan) Task {
		tw := &TaskWrapper{
			ap:      ap,
			persist: model.GetStorage(),
			logger:  entry,
		}
		// task 必须是新初始化的实例，task 内部存了变量需要隔离
		task := taskFn()
		if collect, ok := task.(AuditPlanCollector); ok {
			tw.collect = collect
			tw.isStarted = false
			tw.cancel = make(chan struct{})
			tw.loopInterval = func() time.Duration {
				// todo: 由任务自己定义
				interval := ap.Params.GetParam(paramKeyCollectIntervalSecond).Int()
				if interval != 0 {
					return time.Second * time.Duration(interval)
				}
				interval = ap.Params.GetParam(paramKeyCollectIntervalMinute).Int()
				if interval != 0 {
					return time.Minute * time.Duration(interval)
				}
				return time.Minute * time.Duration(60)
			}
		}
		if handler, ok := task.(AuditPlanHandler); ok {
			tw.handler = handler
		}
		return tw
	}
}

func (at *TaskWrapper) Start() error {
	if at.collect != nil {
		return at.startCollect()
	}
	return nil
}

func (at *TaskWrapper) Stop() error {
	if at.collect != nil {
		return at.stopCollect()
	}
	return nil
}

func (at *TaskWrapper) Audit() (*AuditResultResp, error) {
	// TODO remove
	return at.handler.Audit(nil)
}

// func (at *TaskWrapper) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
// 	return at.handler.GetSQLs(at.ap, at.persist, args)
// }

func (at *TaskWrapper) FullSyncSQLs(sqls []*SQL) error {
	cache := NewSQLV2Cache()
	for _, sql := range sqls {
		sqlV2 := NewSQLV2FromSQL(at.ap, sql)
		meta, err := GetMeta(sqlV2.Source)
		if err != nil {
			at.logger.Warnf("get meta failed, error: %v", err)
			// todo: 有错误咋处理
			continue
		}
		if meta.Handler == nil {
			at.logger.Warnf("do not support this type (%s), error: %v", sqlV2.Source, err)
			// todo: 没有处理器咋办
			continue
		}
		err = meta.Handler.AggregateSQL(cache, sqlV2)
		if err != nil {
			at.logger.Warnf("aggregate sql failed, error: %v", err)
			// todo: 有错误咋处理
			continue
		}
	}
	sqlList := make([]*model.SQLManageQueue, 0, len(cache.GetSQLs()))
	for _, sql := range cache.GetSQLs() {
		sqlList = append(sqlList, ConvertSQLV2ToMangerSQLQueue(sql))
	}

	err := at.pushSQLToManagerSQLQueue(sqlList, at.ap)
	if err != nil {
		at.logger.Errorf("push sql to manager sql queue failed, error : %v", err)
	}

	return err
}

func (at *TaskWrapper) PartialSyncSQLs(sqls []*SQL) error {
	return at.FullSyncSQLs(sqls)
}

func (at *TaskWrapper) startCollect() error {
	if at.isStarted {
		return nil
	}
	interval := at.loopInterval()

	at.WaitGroup.Add(1)
	go func() {
		at.isStarted = true
		at.logger.Infof("start task")
		at.loop(at.cancel, interval)
		at.WaitGroup.Done()
	}()
	return nil
}

func (at *TaskWrapper) stopCollect() error {
	if !at.isStarted {
		return nil
	}
	at.cancel <- struct{}{}
	at.isStarted = false
	at.WaitGroup.Wait()
	at.logger.Infof("stop task")
	return nil
}

func (at *TaskWrapper) extractSQL() {
	var err error
	defer func() {
		status := model.LastCollectionNormal
		if err != nil {
			at.logger.Error(errors.NewAuditPlanExecuteExtractErr(err, at.ap.InstanceID, at.ap.Type))
			status = model.LastCollectionAbnormal
		}
		updateErr := at.persist.UpdateAuditPlanInfoByAPID(at.ap.ID, map[string]interface{}{"last_collection_status": status})
		if updateErr != nil {
			at.logger.Errorf("update audit plan task info collection status status failed, error : %v", updateErr)
		}
	}()
	collectionTime := time.Now()
	sqls, err := at.collect.ExtractSQL(at.logger, at.ap, at.persist)
	if err != nil {
		at.logger.Errorf("extract sql failed, %v", err)
		return
	}
	// todo: 对于mysql慢日志类型，采集来源是scannerd的任务的时间不应该在此处更新
	err = at.persist.UpdateAuditPlanLastCollectionTime(at.ap.ID, collectionTime)
	if err != nil {
		at.logger.Errorf("update audit plan last collection time failed, error : %v", err)
	}
	if len(sqls) == 0 {
		at.logger.Info("extract sql list is empty, skip")
		return
	}
	// 转换类型
	sqlQueues := make([]*model.SQLManageQueue, 0, len(sqls))
	for _, sql := range sqls {
		sqlQueues = append(sqlQueues, ConvertSQLV2ToMangerSQLQueue(sql))
	}

	err = at.pushSQLToManagerSQLQueue(sqlQueues, at.ap)
	if err != nil {
		at.logger.Errorf("push sql to manager sql queue failed, error : %v", err)
	}
}

func (at *TaskWrapper) loop(cancel chan struct{}, interval time.Duration) {
	if interval == 0 {
		at.logger.Warnf("task(%v) loop interval can not be zero", at.ap.Name)
		return
	}
	at.extractSQL()

	tk := time.NewTicker(interval)
	for {
		select {
		case <-cancel:
			tk.Stop()
			return
		case <-tk.C:
			at.logger.Infof("tick %s", at.ap.Name)
			at.extractSQL()
		}
	}
}

func (at *TaskWrapper) pushSQLToManagerSQLQueue(sqlList []*model.SQLManageQueue, ap *AuditPlan) error {
	if len(sqlList) == 0 {
		return nil
	}

	matchedCount, SqlQueueList, err := at.filterSqlManageQueue(sqlList)
	if err != nil {
		return err
	}

	lastMatchedTime := time.Now()
	err = at.persist.BatchUpdateBlackListCount(matchedCount, lastMatchedTime)
	if err != nil {
		return err
	}

	for _, sqlQueue := range SqlQueueList {
		err = createSqlManageCostMetricRecord(sqlQueue, ap.Instance)
		if err != nil {
			log.Logger().Errorf("createSqlManageCostMetricRecord: %v", err)
		}
	}
	err = at.persist.PushSQLToManagerSQLQueue(SqlQueueList)
	if err != nil {
		return err
	}

	return nil
}

func (at *TaskWrapper) filterSqlManageQueue(sqlList []*model.SQLManageQueue) (map[uint]uint, []*model.SQLManageQueue, error) {
	blackListMap := make(map[string][]*model.BlackListAuditPlanSQL)
	instanceMap := make(map[string]string)
	var err error
	matchedCount := make(map[uint]uint)
	SqlQueueList := make([]*model.SQLManageQueue, 0)
	for _, sql := range sqlList {
		blacklist, ok := blackListMap[sql.ProjectId]
		if !ok {
			blacklist, err = at.persist.GetBlackListByProjectID(model.ProjectUID(sql.ProjectId))
			if err != nil {
				return nil, nil, err
			}
			blackListMap[sql.ProjectId] = blacklist
		}

		instName, ok := instanceMap[sql.InstanceID]
		if !ok {
			instance, exist, err := dms.GetInstancesById(context.TODO(), sql.InstanceID)
			if err != nil {
				return nil, nil, err
			}
			if !exist {
				return nil, nil, e.New("instance not exist")
			}
			instName = instance.Name
			instanceMap[sql.InstanceID] = instName
		}

		value, err := sql.Info.OriginValue()
		if err != nil {
			return nil, nil, err
		}
		metrics := LoadMetrics(value, []string{MetricNameEndpoints})
		matchedID, isInBlacklist := FilterSQLsByBlackList(metrics.Get(MetricNameEndpoints).StringArray(), sql.SqlText, sql.SqlFingerprint, instName, blacklist)
		if isInBlacklist {
			matchedCount[matchedID]++
			continue
		}

		SqlQueueList = append(SqlQueueList, sql)
	}

	return matchedCount, SqlQueueList, nil
}

func FilterSQLsByBlackList(endpoint []string, sqlText, sqlFp, instName string, blacklist []*model.BlackListAuditPlanSQL) (uint, bool) {
	if len(blacklist) == 0 {
		return 0, false
	}

	filter := ConvertToBlackFilter(blacklist)

	matchedID, hasEndpointInBlacklist := filter.HasEndpointInBlackList(endpoint)
	if hasEndpointInBlacklist {
		return matchedID, true
	}

	matchedID, isSqlInBlackList := filter.IsSqlInBlackList(sqlText)
	if isSqlInBlackList {
		return matchedID, true
	}

	matchedID, isFpInBlackList := filter.IsFpInBlackList(sqlFp)
	if isFpInBlackList {
		return matchedID, true
	}

	matchedID, isInstNameInBlackList := filter.IsInstNameInBlackList(instName)
	if isInstNameInBlackList {
		return matchedID, true
	}

	return 0, false
}

func ConvertToBlackFilter(blackList []*model.BlackListAuditPlanSQL) *BlackFilter {
	var blackFilter BlackFilter
	for _, filter := range blackList {
		switch filter.FilterType {
		case model.FilterTypeSQL:
			blackFilter.BlackSqlList = append(blackFilter.BlackSqlList, BlackSqlList{
				ID:     filter.ID,
				Regexp: utils.FullFuzzySearchRegexp(filter.FilterContent),
			})
		case model.FilterTypeFpSQL:
			blackFilter.BlackFpList = append(blackFilter.BlackFpList, BlackFpList{
				ID:     filter.ID,
				Regexp: utils.FullFuzzySearchRegexp(filter.FilterContent),
			})
		case model.FilterTypeHost:
			blackFilter.BlackHostList = append(blackFilter.BlackHostList, BlackHostList{
				ID:     filter.ID,
				Regexp: utils.FullFuzzySearchRegexp(filter.FilterContent),
			})
		case model.FilterTypeIP:
			ip := net.ParseIP(filter.FilterContent)
			if ip == nil {
				log.Logger().Errorf("wrong ip in black list,ip:%s", filter.FilterContent)
				continue
			}
			blackFilter.BlackIpList = append(blackFilter.BlackIpList, BlackIpList{
				ID: filter.ID,
				Ip: ip,
			})
		case model.FilterTypeCIDR:
			_, cidr, err := net.ParseCIDR(filter.FilterContent)
			if err != nil {
				log.Logger().Errorf("wrong cidr in black list,cidr:%s,err:%v", filter.FilterContent, err)
				continue
			}
			blackFilter.BlackCidrList = append(blackFilter.BlackCidrList, BlackCidrList{
				ID:   filter.ID,
				Cidr: cidr,
			})
		case model.FilterTypeInstance:
			blackFilter.BlackInstList = append(blackFilter.BlackInstList, BlackInstList{
				ID:       filter.ID,
				InstName: filter.FilterContent,
			})
		}
	}
	return &blackFilter
}

type BlackFilter struct {
	BlackSqlList  []BlackSqlList
	BlackFpList   []BlackFpList
	BlackIpList   []BlackIpList
	BlackHostList []BlackHostList
	BlackCidrList []BlackCidrList
	BlackInstList []BlackInstList
}

type BlackSqlList struct {
	ID     uint
	Regexp *regexp.Regexp
}

type BlackFpList struct {
	ID     uint
	Regexp *regexp.Regexp
}

type BlackIpList struct {
	ID uint
	Ip net.IP
}

type BlackHostList struct {
	ID     uint
	Regexp *regexp.Regexp
}

type BlackCidrList struct {
	ID   uint
	Cidr *net.IPNet
}

type BlackInstList struct {
	ID       uint
	InstName string
}

func (f BlackFilter) IsSqlInBlackList(checkSql string) (uint, bool) {
	for _, blackSql := range f.BlackSqlList {
		if blackSql.Regexp.MatchString(checkSql) {
			return blackSql.ID, true
		}
	}
	return 0, false
}

func (f BlackFilter) IsFpInBlackList(fp string) (uint, bool) {
	for _, blackFp := range f.BlackFpList {
		if blackFp.Regexp.MatchString(fp) {
			return blackFp.ID, true
		}
	}
	return 0, false
}

func (f BlackFilter) IsInstNameInBlackList(instName string) (uint, bool) {
	for _, blackInstName := range f.BlackInstList {
		if blackInstName.InstName == instName {
			return blackInstName.ID, true
		}
	}
	return 0, false
}

// 输入一组ip若其中有一个ip在黑名单中则返回true
func (f BlackFilter) HasEndpointInBlackList(checkIps []string) (uint, bool) {
	var checkNetIp net.IP
	for _, checkIp := range checkIps {
		checkNetIp = net.ParseIP(checkIp)
		if checkNetIp == nil {
			// 无法解析IP，可能是域名，需要正则匹配
			for _, blackHost := range f.BlackHostList {
				if blackHost.Regexp.MatchString(checkIp) {
					return blackHost.ID, true
				}
			}
		} else {
			for _, blackIp := range f.BlackIpList {
				if blackIp.Ip.Equal(checkNetIp) {
					return blackIp.ID, true
				}
			}
			for _, blackCidr := range f.BlackCidrList {
				if blackCidr.Cidr.Contains(checkNetIp) {
					return blackCidr.ID, true
				}
			}
		}
	}

	return 0, false
}

// createSqlManageCostMetricRecord 针对SELECT语句的SQL下发执行计划，并且生成对应的实体类
func createSqlManageCostMetricRecord(sqlManageQueue *model.SQLManageQueue, instance *model.Instance) error {
	// 建表语句直接忽略 mysql_schema_meta
	if sqlManageQueue.Source == TypeMySQLSchemaMeta {
		return nil
	}
	dsn, err := common.NewDSN(instance, sqlManageQueue.SchemaName)
	if err != nil {
		return err
	}

	plugin, err := driver.GetPluginManager().OpenPlugin(log.NewEntry(), instance.DbType, &driverV2.Config{DSN: dsn})
	if err != nil {
		return err
	}
	defer plugin.Close(context.TODO())
	sqlNodes, err := plugin.Parse(context.TODO(), sqlManageQueue.SqlText)
	if err != nil {
		return err
	}
	if len(sqlNodes) == 0 {
		return errors.Errorf("invalid sql: %v", sqlManageQueue.SqlText)
	}
	// 不是SELECT语句直接忽略
	if sqlNodes[0].Type != driverV2.SQLTypeDQL {
		return nil
	}
	// EXPLAIN
	explainResult, err := plugin.Explain(context.TODO(), &driverV2.ExplainConf{Sql: sqlManageQueue.SqlText})
	if err != nil {
		return err
	}
	// EXPLAIN FORMAT
	explainJSONResult, err := plugin.ExplainJSONFormat(context.TODO(), &driverV2.ExplainConf{Sql: sqlManageQueue.SqlText})
	if err != nil {
		return err
	}
	cost, err := strconv.ParseFloat(explainJSONResult.QueryBlock.CostInfo.QueryCost, 64)
	if err != nil {
		log.Logger().Errorf("createSqlManageCostMetricRecord: parse explain cost to float64 failed %v", err)
		return err
	}
	storage := model.GetStorage()
	nowTime := time.Now()
	sqlManageMetricRecord := &model.SqlManageMetricRecord{
		SQLID:          sqlManageQueue.SQLID,
		ExecutionCount: 1,
		RecordBeginAt:  nowTime,
		RecordEndAt:    nowTime,
	}
	if err = storage.Create(sqlManageMetricRecord); err != nil {
		log.Logger().Errorf("createSqlManageCostMetricRecord: create SqlManageMetricRecord error sqlId: %v", sqlManageQueue.SQLID)
		return err
	}
	sqlManageMetricValue := &model.SqlManageMetricValue{
		SqlManageMetricRecordID: sqlManageMetricRecord.ID,
		MetricName:              MetricNameExplainCost,
		MetricValue:             cost,
	}
	if err = storage.Create(sqlManageMetricValue); err != nil {
		log.Logger().Errorf("createSqlManageCostMetricRecord: create sqlManageMetricValue error sqlId: %v", sqlManageQueue.SQLID)
		return err
	}
	sqlManageMetricExecutePlanRecords := buildSqlManageMetricExecutePlanRecord(explainResult.ClassicResult.Rows, sqlManageMetricRecord.ID)
	if err = storage.Create(sqlManageMetricExecutePlanRecords); err != nil {
		log.Logger().Errorf("createSqlManageCostMetricRecord: create sqlManageMetricExecutePlanRecord error sqlId: %v", sqlManageQueue.SQLID)
		return err
	}
	return nil
}

// buildSqlManageMetricExecutePlanRecord 构建执行计划记录-批量
func buildSqlManageMetricExecutePlanRecord(explainResultRows [][]string, sqlManageMetricRecordID uint) []model.SqlManageMetricExecutePlanRecord {
	sqlManageMetricExecutePlanRecords := make([]model.SqlManageMetricExecutePlanRecord, 0)
	for _, row := range explainResultRows {
		selectId, _ := strconv.Atoi(row[0])
		keyLength, _ := strconv.Atoi(row[7])
		rows, _ := strconv.Atoi(row[9])
		filtered, _ := strconv.ParseFloat(row[10], 64)
		sqlManageMetricExecutePlanRecords = append(sqlManageMetricExecutePlanRecords, model.SqlManageMetricExecutePlanRecord{
			SqlManageMetricRecordID: sqlManageMetricRecordID,
			SelectId:                selectId,
			SelectType:              row[1],
			Table:                   row[2],
			Partitions:              row[3],
			Type:                    row[4],
			PossibleKeys:            row[5],
			Key:                     row[6],
			KeyLen:                  keyLength,
			Ref:                     row[8],
			Rows:                    rows,
			Filtered:                filtered,
			Extra:                   row[11],
		})
	}
	return sqlManageMetricExecutePlanRecords
}
