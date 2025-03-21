//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

// TDSQL For MySQL Distribute Lock Collection
type TDMySQLDistributedLock struct {
	DefaultTaskV2
}

func NewTDMySQLDistributedLockTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &TDMySQLDistributedLock{}
	}
}

/*
	implement of AuditPlanMeta
*/

var _ AuditPlanMeta = &TDMySQLDistributedLock{}

func (at *TDMySQLDistributedLock) InstanceType() string {
	return InstanceTypeTDSQL
}

// Params定义了采集任务的用户可配置项
func (at *TDMySQLDistributedLock) Params(instanceId ...string) params.Params {
	return params.Params{
		{
			Key:      paramKeyCollectIntervalSecond,
			Value:    "60",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamCollectIntervalSecond),
		},
	}
}

func (at *TDMySQLDistributedLock) HighPriorityParams() params.ParamsWithOperator {
	return params.ParamsWithOperator{}
}

// Metrics定义了采集任务的采集项
func (at *TDMySQLDistributedLock) Metrics() []string {
	return []string{
		MetricNameGrantedLockTrxStarted,
		MetricNameWaitingLockTrxWaitStarted,
		MetricNameGrantedLockConnectionId,
		MetricNameWaitingLockConnectionId,
		MetricNameGrantedLockTrxId,
		MetricNameWaitingLockTrxId,
		MetricNameGrantedLockId,
		MetricNameWaitingLockId,
		MetricNameObjectName,
		MetricNameEngine,
		MetricNameDatabase,
		MetricNameLockType,
		MetricNameLockMode,
		MetricNameLockStatus,
		MetricNameDBUser,
		MetricNameHost,
		MetricNameIndexType,
		MetricNameGrantedLockSql,
		MetricNameWaitingLockSql,
	}
}

/*
	implement of AuditPlanCollector
*/

var _ AuditPlanCollector = &TDMySQLDistributedLock{}

func (at *TDMySQLDistributedLock) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	// 1. get executor
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	instance, err := getInstanceOfAuditPlan(ctx, ap)
	if err != nil {
		return nil, err
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

	// 2. execute SQL and get result
	resultRows, err := db.Db.Query(at.getLockSQL())
	if err != nil {
		logger.Errorf("execute SQL fail, error: %v", err)
		return nil, fmt.Errorf("execute SQL fail, error: %v", err)
	}

	// 3. convert result and return
	sqlV2s := at.convertToSQLV2(logger, ap, resultRows)
	return sqlV2s, nil
}

func getInstanceOfAuditPlan(ctx context.Context, ap *AuditPlan) (*model.Instance, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	instance, exist, err := dms.GetInstancesById(ctx, ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}
	if !exist {
		return nil, errors.NewInstanceNoExistErr()
	}
	return instance, nil
}

func (at *TDMySQLDistributedLock) getLockSQL() string {
	return `SELECT
    granted_locks.ENGINE_LOCK_ID granted_lock_id,
    waiting_locks.ENGINE_LOCK_ID waiting_lock_id,
    granted_locks.ENGINE AS engine,
    granted_locks.OBJECT_SCHEMA AS database_name,
    granted_locks.OBJECT_NAME AS object_name,
    granted_lock_threads.PROCESSLIST_USER AS db_user,
    granted_lock_threads.PROCESSLIST_HOST AS host,
    granted_locks.INDEX_NAME AS index_type,
    granted_locks.LOCK_TYPE AS lock_type,
    granted_locks.LOCK_MODE AS lock_mode,
    granted_lock_threads.PROCESSLIST_INFO AS granted_lock_sql,
    waiting_lock_threads.PROCESSLIST_INFO AS waiting_lock_sql,
    granted_lock_threads.PROCESSLIST_ID AS granted_lock_connection_id,
    waiting_lock_threads.PROCESSLIST_ID AS waiting_lock_connection_id,
    granted_locks.ENGINE_TRANSACTION_ID AS granted_lock_trx_id,
    waiting_locks.ENGINE_TRANSACTION_ID AS waiting_lock_trx_id,
    granted_trx.TRX_STARTED AS trx_started,
    waiting_trx.TRX_WAIT_STARTED AS trx_wait_started
FROM information_schema.innodb_trx granted_trx
JOIN performance_schema.data_locks granted_locks ON granted_trx.trx_id = granted_locks.ENGINE_TRANSACTION_ID
JOIN performance_schema.data_lock_waits lock_waits ON granted_locks.ENGINE_LOCK_ID = lock_waits.BLOCKING_ENGINE_LOCK_ID
JOIN performance_schema.data_locks waiting_locks ON lock_waits.REQUESTING_ENGINE_LOCK_ID = waiting_locks.ENGINE_LOCK_ID
JOIN information_schema.innodb_trx waiting_trx ON waiting_locks.ENGINE_TRANSACTION_ID = waiting_trx.trx_id
JOIN performance_schema.threads granted_lock_threads ON granted_locks.THREAD_ID = granted_lock_threads.THREAD_ID
JOIN performance_schema.threads waiting_lock_threads ON waiting_locks.THREAD_ID = waiting_lock_threads.THREAD_ID
WHERE granted_locks.LOCK_STATUS = 'GRANTED';`
}

func (at *TDMySQLDistributedLock) convertToSQLV2(logger *logrus.Entry, ap *AuditPlan, resultRows []map[string]sql.NullString) []*SQLV2 {
	sqlV2List := make([]*SQLV2, 0, len(resultRows))
	for _, row := range resultRows {
		sqlV2 := &SQLV2{
			Source:      ap.Type,
			SourceId:    strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
			ProjectId:   ap.ProjectId,
			InstanceID:  ap.InstanceID,
			AuditPlanId: strconv.FormatUint(uint64(ap.ID), 10),
			SchemaName:  row[MetricNameDatabase].String,
		}
		layout := "2006-01-02T15:04:05Z"
		metricNameTrxStartedStr := row[MetricNameGrantedLockTrxStarted].String
		metricNameTrxWaitStartedStr := row[MetricNameWaitingLockTrxWaitStarted].String
		metricNameTrxStarted, err := time.Parse(layout, metricNameTrxStartedStr)
		// 这里暂时将锁的信息全部存储到info中
		info := NewMetrics()
		if err == nil {
			info.SetTime(MetricNameGrantedLockTrxStarted, &metricNameTrxStarted)
		}
		metricNameTrxWaitStarted, err := time.Parse(layout, metricNameTrxWaitStartedStr)
		if err == nil {
			info.SetTime(MetricNameWaitingLockTrxWaitStarted, &metricNameTrxWaitStarted)
		}
		grantedLockConnectionId, err := strconv.ParseInt(row[MetricNameGrantedLockConnectionId].String, 10, 64)
		if err == nil {
			info.SetInt(MetricNameGrantedLockConnectionId, grantedLockConnectionId)
		}
		waitingLockConnectionId, err := strconv.ParseInt(row[MetricNameWaitingLockConnectionId].String, 10, 64)
		if err == nil {
			info.SetInt(MetricNameWaitingLockConnectionId, waitingLockConnectionId)
		}
		grantedLockTrxId, err := strconv.ParseInt(row[MetricNameGrantedLockTrxId].String, 10, 64)
		if err == nil {
			info.SetInt(MetricNameGrantedLockTrxId, grantedLockTrxId)
		}
		waitingLockTrxId, err := strconv.ParseInt(row[MetricNameWaitingLockTrxId].String, 10, 64)
		if err == nil {
			info.SetInt(MetricNameWaitingLockTrxId, waitingLockTrxId)
		}
		info.SetString(MetricNameGrantedLockId, row[MetricNameGrantedLockId].String)
		info.SetString(MetricNameWaitingLockId, row[MetricNameWaitingLockId].String)
		info.SetString(MetricNameObjectName, row[MetricNameObjectName].String)
		info.SetString(MetricNameEngine, row[MetricNameEngine].String)
		info.SetString(MetricNameDatabase, row[MetricNameDatabase].String)
		info.SetString(MetricNameLockType, row[MetricNameLockType].String)
		info.SetString(MetricNameLockMode, row[MetricNameLockMode].String)
		info.SetString(MetricNameLockStatus, row[MetricNameLockStatus].String)
		info.SetString(MetricNameDBUser, row[MetricNameDBUser].String)
		info.SetString(MetricNameHost, row[MetricNameHost].String)
		info.SetString(MetricNameIndexType, row[MetricNameIndexType].String)
		info.SetString(MetricNameGrantedLockSql, row[MetricNameGrantedLockSql].String)
		info.SetString(MetricNameWaitingLockSql, row[MetricNameWaitingLockSql].String)
		sqlV2.Info = info
		sqlV2List = append(sqlV2List, sqlV2)
	}
	return sqlV2List
}

/*
	implement of AuditPlanHandler
*/

var _ AuditPlanHandler = &TDMySQLDistributedLock{}

func (at *TDMySQLDistributedLock) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *TDMySQLDistributedLock) Head(ap *AuditPlan) []Head {
	return []Head{
		{
			Name: MetricNameGrantedLockId,
			Desc: locale.ApMetricNameGrantedLockId,
		},
		{
			Name: MetricNameGrantedLockTrxId,
			Desc: locale.ApMetricNameGrantedLockTrxId,
		},
		{
			Name: MetricNameGrantedLockTrxStarted,
			Desc: locale.ApMetricNameTrxStarted,
		},
		{
			Name: MetricNameGrantedLockConnectionId,
			Desc: locale.ApMetricNameGrantedLockConnectionId,
		},
		{
			Name: MetricNameGrantedLockSql,
			Desc: locale.ApMetricNameGrantedLockSql,
		},
		{
			Name: MetricNameWaitingLockId,
			Desc: locale.ApMetricNameWaitingLockId,
		},
		{
			Name: MetricNameWaitingLockTrxId,
			Desc: locale.ApMetricNameWaitingLockTrxId,
		},
		{
			Name: MetricNameWaitingLockTrxWaitStarted,
			Desc: locale.ApMetricNameTrxWaitStarted,
		},
		{
			Name: MetricNameWaitingLockConnectionId,
			Desc: locale.ApMetricNameWaitingLockConnectionId,
		},
		{
			Name: MetricNameWaitingLockSql,
			Desc: locale.ApMetricNameWaitingLockSql,
		},
		{
			Name: MetricNameLockType,
			Desc: locale.ApMetricNameLockType,
		},
		{
			Name: MetricNameLockMode,
			Desc: locale.ApMetricNameLockMode,
		},
		{
			Name: MetricNameDatabase,
			Desc: locale.ApDatabase,
		},
		{
			Name: MetricNameObjectName,
			Desc: locale.ApMetricNameTable,
		},
		{
			Name: MetricNameHost,
			Desc: locale.ApMetricNameHost,
		},
		{
			Name: MetricNameDBUser,
			Desc: locale.ApMetricNameDBUser,
		},
		{
			Name: MetricNameEngine,
			Desc: locale.ApMetricEngine,
		},
	}
}

func (at *TDMySQLDistributedLock) Filters(ctx context.Context, logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta {
	return []FilterMeta{
		{
			Name:            MetricNameLockType,
			Desc:            locale.ApMetricNameLockType,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      retrieveFilterTips(MetricNameLockType, persist, ap, logger),
		},
		{
			Name:            MetricNameDatabase,
			Desc:            locale.ApDatabase,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      retrieveFilterTips(MetricNameDatabase, persist, ap, logger),
		},
		{
			Name:            MetricNameObjectName,
			Desc:            locale.ApMetricNameTable,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      retrieveFilterTips(MetricNameObjectName, persist, ap, logger),
		},
	}
}

func retrieveFilterTips(filterType string, persist *model.Storage, ap *AuditPlan, logger *logrus.Entry) []FilterTip {
	filterTypes, err := persist.SelectDistinctColumn(ap.ID, ap.InstanceAuditPlanId, filterType)
	if err != nil {
		logger.Warnf("retrieve filter tips failed, err: %v", err)
		return []FilterTip{}
	}
	filterTips := make([]FilterTip, 0, len(filterTypes))
	for _, lockType := range filterTypes {
		filterTips = append(filterTips, FilterTip{
			Value: lockType,
		})
	}
	return filterTips
}

// GetSQLData 获取分布式锁结果集接口
func (at *TDMySQLDistributedLock) GetSQLData(ctx context.Context, ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
	dataLocks, err := persist.GetDataLockList(genDataLockFilters(filters), limit, offset, ap.ID, ap.InstanceAuditPlanId)
	if err != nil {
		return nil, 0, err
	}
	count, err := persist.CountDataLock(genDataLockFilters(filters), ap.ID, ap.InstanceAuditPlanId)
	if err != nil {
		return nil, 0, err
	}
	rows := make([]map[string]string, 0, len(dataLocks))
	layout := "2006-01-02T15:04:05Z"
	for _, lock := range dataLocks {
		trxStarted := ""
		if lock.TrxStarted != nil {
			trxStarted = lock.TrxStarted.Format(layout)
		}
		trxWaitStarted := ""
		if lock.TrxWaitStarted != nil {
			trxWaitStarted = lock.TrxWaitStarted.Format(layout)
		}
		rows = append(rows, map[string]string{
			MetricNameGrantedLockId:             lock.GrantedLockId,
			MetricNameWaitingLockId:             lock.WaitingLockId,
			MetricNameEngine:                    lock.Engine,
			MetricNameDBUser:                    lock.DbUser,
			MetricNameHost:                      lock.Host,
			MetricNameDatabase:                  lock.DatabaseName,
			MetricNameObjectName:                lock.ObjectName,
			MetricNameLockType:                  lock.LockType,
			MetricNameLockMode:                  lock.LockMode,
			MetricNameGrantedLockConnectionId:   strconv.FormatInt(lock.GrantedLockConnectionId, 10),
			MetricNameWaitingLockConnectionId:   strconv.FormatInt(lock.WaitingLockConnectionId, 10),
			MetricNameGrantedLockTrxId:          strconv.FormatInt(lock.GrantedLockTrxId, 10),
			MetricNameWaitingLockTrxId:          strconv.FormatInt(lock.WaitingLockTrxId, 10),
			MetricNameGrantedLockSql:            lock.GrantedLockSql,
			MetricNameWaitingLockSql:            lock.WaitingLockSql,
			MetricNameGrantedLockTrxStarted:     trxStarted,
			MetricNameWaitingLockTrxWaitStarted: trxWaitStarted,
		})
	}
	return rows, uint64(count), nil
}

func genDataLockFilters(filters []Filter) map[string]string {
	dataLockFilters := make(map[string]string, len(filters))
	for _, filter := range filters {
		dataLockFilters[filter.Name] = filter.FilterComparisonValue
	}
	return dataLockFilters
}
