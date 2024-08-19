package auditplan

import (
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/model"
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

	err := at.persist.PushSQLToManagerSQLQueue(sqlList)
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
		return
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
	err = at.persist.PushSQLToManagerSQLQueue(sqlQueues)
	if err != nil {
		at.logger.Errorf("push sql to manager sql queue failed, error : %v", err)
		return
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
