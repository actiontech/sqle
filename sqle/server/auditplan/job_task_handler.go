package auditplan

import (
	"database/sql"
	"fmt"
	"time"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
		return
	}

	go handlerSQLAudit(entry, sqlList)

}

// todo: 错误处理
func handlerSQLAudit(entry *logrus.Entry, sqlList []*model.SQLManageRecord) {
	s := model.GetStorage()
	sqlList, err := BatchAuditSQLs(sqlList, true)
	if err != nil {
		entry.Warnf("batch audit manager sql failed, error: %v", err)
	}
	// 设置高优先级
	sqlList, err = SetSQLPriority(sqlList)
	if err != nil {
		entry.Warnf("set sql priority sql failed, error: %v", err)
	}
	for _, sql := range sqlList {
		manageSqlParam := make(map[string]interface{}, 3)
		manageSqlParam["audit_level"] = sql.AuditLevel
		manageSqlParam["audit_results"] = sql.AuditResults
		manageSqlParam["priority"] = sql.Priority
		err = s.UpdateManagerSQLBySqlId(manageSqlParam, sql.SQLID)
		if err != nil {
			entry.Warnf("update manager sql failed, error: %v", err)
			continue
		}
	}
}

func BatchAuditSQLs(sqlList []*model.SQLManageRecord, isSkipAuditedSql bool) ([]*model.SQLManageRecord, error) {

	// SQL聚合
	sqlMap := make(map[string][]*model.SQLManageRecord)
	for _, sql := range sqlList {
		if isSkipAuditedSql && sql.AuditLevel != "" {
			continue
		}

		// 根据source id和schema name 聚合sqls，避免task内需要切换schema上下文审核
		key := fmt.Sprintf("%s:%s", sql.SourceId, sql.SchemaName)
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

func SetSQLPriority(sqlList []*model.SQLManageRecord) ([]*model.SQLManageRecord, error) {
	var err error
	s := model.GetStorage()
	// SQL聚合
	auditPlanMap := make(map[string]*model.AuditPlanV2, 0)

	for i, sql_ := range sqlList {
		sourceId := sql_.SourceId
		sourceType := sql_.Source
		auditPlan, ok := auditPlanMap[sourceId]
		if !ok {
			var exist bool
			auditPlan, exist, err = s.GetAuditPlanByInstanceIdAndType(sourceId, sourceType)
			if err != nil {
				return nil, err
			}
			if !exist {
				continue
			}
			auditPlanMap[sourceId] = auditPlan
		}
		priority, _, err := GetSingleSQLPriorityWithReasons(auditPlan, sql_)
		if err != nil {
			return nil, err
		}
		if priority == model.PriorityHigh {
			sqlList[i].Priority = sql.NullString{
				String: model.PriorityHigh,
				Valid:  true,
			}
		}
	}
	return sqlList, nil
}

func GetSingleSQLPriorityWithReasons(auditPlan *model.AuditPlanV2, sql *model.SQLManageRecord) (priority string, reasons []string, err error) {
	if auditPlan == nil || sql == nil {
		return "", reasons, nil
	}
	info, err := sql.Info.OriginValue()
	if err != nil {
		return "", nil, err
	}
	highPriorityConditions := auditPlan.HighPriorityParams
	// 遍历优先级条件
	for _, highPriorityCondition := range highPriorityConditions {
		var compareParamVale string
		// 特殊处理审核级别
		if highPriorityCondition.Key == OperationParamAuditLevel {
			switch sql.AuditLevel {
			case string(driverV2.RuleLevelNotice):
				compareParamVale = "1"
			case string(driverV2.RuleLevelWarn):
				compareParamVale = "2"
			case string(driverV2.RuleLevelError):
				compareParamVale = "3"
			default:
				compareParamVale = "0"
			}
		} else {
			// 获取信息中的相应字段值
			infoV, ok := info[highPriorityCondition.Key]
			if !ok {
				continue
			}
			compareParamVale = fmt.Sprintf("%v", infoV)
		}
		// 检查是否为高优先级条件
		if high, err := highPriorityConditions.CompareParamValue(highPriorityCondition.Key, compareParamVale); err == nil && high {
			// 添加匹配的条件作为原因
			if highPriorityCondition.Key == OperationParamAuditLevel {
				reasons = append(reasons, fmt.Sprintf("【%v %v %v，为：%s】", highPriorityCondition.Desc, highPriorityCondition.Operator.Value, toComparedValue(highPriorityCondition.Param.Value), toComparedValue(compareParamVale)))
			} else {
				reasons = append(reasons, fmt.Sprintf("【%v %v %v，为：%s】", highPriorityCondition.Desc, highPriorityCondition.Operator.Value, highPriorityCondition.Param.Value, toComparedValue(compareParamVale)))
			}
		}
	}
	if len(reasons) > 0 {
		return model.PriorityHigh, reasons, nil
	}
	return "", reasons, nil
}

func toComparedValue(compareParamVale string) string {
	switch compareParamVale {
	case "1":
		return "提示"
	case "2":
		return "警告"
	case "3":
		return "错误"
	default:
		return compareParamVale
	}
}
