package auditplan

import (
	"time"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuditPlanAggregateSQLJob struct {
	server.BaseJob
}

func NewAuditPlanAggregateSQLJob(entry *logrus.Entry) server.ServerJob {
	entry = entry.WithField("job", "audit_plan_aggregate_sql")
	j := &AuditPlanAggregateSQLJob{}
	j.BaseJob = *server.NewBaseJob(entry, 5*time.Second, j.AggregateSQL)
	return j
}

func (j *AuditPlanAggregateSQLJob) AggregateSQL(entry *logrus.Entry) {
	s := model.GetStorage()
	queues, err := s.PullSQLFromManagerSQLQueue()
	if err != nil {
		entry.Warnf("pull sql from queue failed, error: %v", err)
		return
	}
	if len(queues) == 0 {
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
	// todo: 错误处理
	if err = s.Tx(func(txDB *gorm.DB) error {
		for _, sql := range sqlList {
			err := s.SaveManagerSQL(txDB, sql)
			if err != nil {
				entry.Warnf("update manager sql failed, error: %v", err)
				return err
			}

			// 更新状态表
			err = s.UpdateManagerSQLStatus(txDB, sql)
			if err != nil {
				entry.Warnf("update manager sql status failed, error: %v", err)
				return err
			}
		}

		for _, sql := range queues {
			err := s.RemoveSQLFromQueue(txDB, sql)
			if err != nil {
				entry.Warnf("remove manager sql queue failed, error: %v", err)
				return err
			}
		}

		return nil

	}); err != nil {
		entry.Warnf("audit plan aggregate sql job failed, error: %v", err)
		return
	}
}
