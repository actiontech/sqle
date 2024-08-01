package auditplan

import (
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/sirupsen/logrus"
)

type AuditPlanHandlerJob struct {
	server.BaseJob
}

func NewAuditPlanHandlerJob(entry *logrus.Entry) server.ServerJob {
	entry = entry.WithField("job", "audit_plan_handler")
	j := &AuditPlanHandlerJob{}
	j.BaseJob = *server.NewBaseJob(entry, 5*time.Second, j.HandlerSQL)
	return j
}

func (j *AuditPlanHandlerJob) HandlerSQL(entry *logrus.Entry) {
	s := model.GetStorage()
	queues, err := s.PullSQLFromManagerSQLQueue()
	if err != nil {
		entry.Warnf("pull sql from queue failed, error: %v", err)
		return
	}
	cache := NewSQLV2CacheWithPersist(s)
	for _, sql := range queues {
		sqlV2 := ConvertMangerSQLQueueToSQLV2(sql)
		meta, err := GetMeta(sqlV2.Source)
		if err != nil {
			entry.Warnf("get meta failed, error: %v", err)
			// todo: 有错误咋处理
			continue
		}
		if meta.Handler == nil {
			entry.Warnf("do not support this type (%s), error: %v", sqlV2.Source, err)
			// todo: 没有处理器咋办
			continue
		}
		err = meta.Handler.AggregateSQL(cache, sqlV2)
		if err != nil {
			entry.Warnf("aggregate sql failed, error: %v", err)
			// todo: 有错误咋处理
			continue
		}

	}

	sqlList := make([]*model.SQLManageRecord, 0, len(cache.GetSQLs()))
	for _, sql := range cache.GetSQLs() {
		sqlList = append(sqlList, ConvertSQLV2ToMangerSQL(sql))
	}

	if len(sqlList) == 0 {
		return
	}
	// 审核
	auditSQLs, err := batchAuditSQLs(sqlList)
	if err != nil {
		entry.Warnf("batch audit origin manager sql failed, error: %v", err)
		return
	}

	// todo: 保证事务和错误处理
	for _, sql := range auditSQLs {
		err := s.UpdateManagerSQL(sql)
		if err != nil {
			entry.Warnf("update manager sql failed, error: %v", err)
			return
		}

		// 同时更新状态表
		err = s.UpdateManagerSQLStatus(sql)
		if err != nil {
			entry.Warnf("update manager sql status failed, error: %v", err)
			return
		}

	}
	for _, sql := range queues {
		err := s.RemoveSQLFromQueue(sql)
		if err != nil {
			entry.Warnf("remove manager sql queue failed, error: %v", err)
			return
		}
	}
}

func batchAuditSQLs(sqlList []*model.SQLManageRecord) ([]*model.SQLManageRecord, error) {

	// SQL聚合
	sqlMap := make(map[string][]*model.SQLManageRecord)
	for _, sql := range sqlList {
		// skip audited sql
		if sql.AuditLevel != "" {
			continue
		}
		// 根据source id和schema name 聚合sqls，避免task内需要切换schema上下文审核
		key := fmt.Sprintf("%d:%s", sql.SourceId, sql.SchemaName)
		_, ok := sqlMap[key]
		if !ok {
			sqlMap[key] = make([]*model.SQLManageRecord, 0)
		}
		sqlMap[key] = append(sqlMap[key], sql)

	}

	auditSQLs := make([]*model.SQLManageRecord, 0)
	// 聚合的SQL批量审核
	for _, sqls := range sqlMap {
		// get audit plan by source id
		auditPlanType := sqls[0].Source

		meta, err := GetMeta(auditPlanType)
		if err != nil {
			return nil, err
		}
		resp, err := meta.Handler.Audit(sqls)
		if err != nil {
			log.NewEntry().Errorf("audit sqls in task fail %v,ignore audit result", err)
			auditSQLs = append(auditSQLs, sqls...)
			continue
		}

		// 更新原值
		auditSQLs = append(auditSQLs, resp.AuditedSqls...)

	}
	return auditSQLs, nil
}
