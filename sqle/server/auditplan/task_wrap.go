package auditplan

import (
	"context"
	e "errors"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pkg/errors"

	"github.com/actiontech/sqle/sqle/dms"
	sqleErr "github.com/actiontech/sqle/sqle/errors"
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
		ProcessAuditPlanStatusAndLogError(at.logger, at.ap.ID, at.ap.InstanceID, at.ap.Type, err)
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
	}
	if at.ap.Type == TypeTDMySQLDistributedLock {
		err = at.pushSQLToDataLock(sqls, at.ap)
	} else {
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
}

func ProcessAuditPlanStatusAndLogError(l *logrus.Entry, auditPlanId uint, instanceID, auditPlnType string, err error) {
	s := model.GetStorage()
	status := model.LastCollectionNormal
	if err != nil {
		l.Error(sqleErr.NewAuditPlanExecuteExtractErr(err, instanceID, auditPlnType))
		status = model.LastCollectionAbnormal
	}
	updateErr := s.UpdateAuditPlanInfoByAPID(auditPlanId, map[string]interface{}{"last_collection_status": status})
	if updateErr != nil {
		l.Errorf("update audit plan task info collection status status failed, error : %v", updateErr)
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

	// 目前只支持MySQL
	if ap.DBType == driverV2.DriverTypeMySQL {
		for _, sqlQueue := range SqlQueueList {
			err = createSqlManageCostMetricRecord(sqlQueue, ap.Instance)
			if err != nil {
				log.Logger().Errorf("createSqlManageCostMetricRecord: %v", err)
			}
		}
	}
	err = at.persist.PushSQLToManagerSQLQueue(SqlQueueList)
	if err != nil {
		return err
	}

	return nil
}

func (at *TaskWrapper) pushSQLToDataLock(sqlList []*SQLV2, auditPlan *AuditPlan) error {
	dataLocks := make([]*model.DataLock, 0, len(sqlList))
	for _, sql := range sqlList {
		instanceAuditPlanId, err := strconv.ParseUint(sql.SourceId, 10, 64)
		if err != nil {
			log.Logger().Errorf("parse sql source id failed, error : %v", err)
			continue
		}
		auditPlanId, err := strconv.ParseUint(sql.AuditPlanId, 10, 64)
		if err != nil {
			log.Logger().Errorf("parse sql audit plan id failed, error : %v", err)
			continue
		}
		dataLocks = append(dataLocks, &model.DataLock{
			GrantedLockId:           sql.Info.Get(MetricNameGrantedLockId).String(),
			WaitingLockId:           sql.Info.Get(MetricNameWaitingLockId).String(),
			Engine:                  sql.Info.Get(MetricNameEngine).String(),
			AuditPlanId:             auditPlanId,
			InstanceAuditPlanId:     instanceAuditPlanId,
			DatabaseName:            sql.SchemaName,
			DbUser:                  sql.Info.Get(MetricNameDBUser).String(),
			Host:                    sql.Info.Get(MetricNameHost).String(),
			ObjectName:              sql.Info.Get(MetricNameObjectName).String(),
			IndexType:               sql.Info.Get(MetricNameIndexType).String(),
			LockType:                sql.Info.Get(MetricNameLockType).String(),
			LockMode:                sql.Info.Get(MetricNameLockMode).String(),
			GrantedLockConnectionId: sql.Info.Get(MetricNameGrantedLockConnectionId).Int(),
			WaitingLockConnectionId: sql.Info.Get(MetricNameWaitingLockConnectionId).Int(),
			GrantedLockSql:          sql.Info.Get(MetricNameGrantedLockSql).String(),
			WaitingLockSql:          sql.Info.Get(MetricNameWaitingLockSql).String(),
			GrantedLockTrxId:        sql.Info.Get(MetricNameGrantedLockTrxId).Int(),
			WaitingLockTrxId:        sql.Info.Get(MetricNameWaitingLockTrxId).Int(),
			TrxStarted:              sql.Info.Get(MetricNameGrantedLockTrxStarted).Time(),
			TrxWaitStarted:          sql.Info.Get(MetricNameWaitingLockTrxWaitStarted).Time(),
		})
	}
	return at.persist.PushSQLToDataLock(dataLocks, auditPlan.ID, auditPlan.InstanceAuditPlanId)
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
		dbUser := metrics.Get(MetricNameDBUser).String()
		matchedID, isInBlacklist := FilterSQLsByBlackList(metrics.Get(MetricNameEndpoints).StringArray(), sql.SqlText, sql.SqlFingerprint, instName, dbUser, blacklist)
		if isInBlacklist {
			matchedCount[matchedID]++
			continue
		}

		SqlQueueList = append(SqlQueueList, sql)
	}

	return matchedCount, SqlQueueList, nil
}

func FilterSQLsByBlackList(endpoint []string, sqlText, sqlFp, instName, dbUser string, blacklist []*model.BlackListAuditPlanSQL) (uint, bool) {
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
	if dbUser != "" {
		matchedID, isDbUserInBlackList := filter.IsDbUserInBlackList(dbUser)
		if isDbUserInBlackList {
			return matchedID, true
		}
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
		case model.FilterTypeDbUser:
			blackFilter.BlackDbUserList = append(blackFilter.BlackDbUserList, BlackDbUser{
				ID:         filter.ID,
				DbUserName: filter.FilterContent,
			})
		}
	}
	return &blackFilter
}

type BlackFilter struct {
	BlackSqlList    []BlackSqlList
	BlackFpList     []BlackFpList
	BlackIpList     []BlackIpList
	BlackHostList   []BlackHostList
	BlackCidrList   []BlackCidrList
	BlackInstList   []BlackInstList
	BlackDbUserList []BlackDbUser
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

type BlackDbUser struct {
	ID         uint
	DbUserName string
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

func (f BlackFilter) IsDbUserInBlackList(checkDbUser string) (uint, bool) {
	for _, blackDbUser := range f.BlackDbUserList {
		if blackDbUser.DbUserName == checkDbUser {
			return blackDbUser.ID, true
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
	cost, err := GetQueryCost(plugin, sqlManageQueue.SqlText)
	if err != nil {
		log.Logger().Errorf("createSqlManageCostMetricRecord: failed %v", err)
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

func GetQueryCost(plugin driver.Plugin, sql string) (cost float64, err error) {
	explainJSONResult, err := plugin.ExplainJSONFormat(context.TODO(), &driverV2.ExplainConf{Sql: sql})
	if err != nil {
		return 0, err
	}
	sqlNodes, err := plugin.Parse(context.TODO(), sql)
	if err != nil {
		return 0, err
	}
	if len(sqlNodes) == 0 {
		return 0, errors.Errorf("failed to get query cost invalid sql: %v", sql)
	}
	// 不是SELECT语句直接忽略
	if sqlNodes[0].Type != driverV2.SQLTypeDQL {
		return 0, errors.Errorf("failed to get query cost because it is not DQL: %v", sql)
	}
	costFloat, err := strconv.ParseFloat(explainJSONResult.QueryBlock.CostInfo.QueryCost, 64)
	if err != nil {
		return 0, err
	}
	return costFloat, nil
}

func Explain(dbType string, plugin driver.Plugin, sql string) (res *driverV2.ExplainResult, err error) {
	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleExplain) {
		return nil, driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleExplain)
	}

	return plugin.Explain(context.TODO(), &driverV2.ExplainConf{Sql: sql})
}
