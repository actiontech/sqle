package auditplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/locale"
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
	sqlList, err := s.GetSQLsToAuditFromManage()
	if err != nil {
		entry.Warnf("get sqls to audit failed, error: %v", err)
		return
	}
	if len(sqlList) == 0 {
		return
	}
	sqlList, err = BatchAuditSQLs(entry, sqlList)
	if err != nil {
		entry.Warnf("batch audit manager sql failed, error: %v", err)
	}
	// 设置高优先级
	sqlList, err = SetSQLPriority(sqlList)
	if err != nil {
		entry.Warnf("set sql priority sql failed, error: %v", err)
	}
	// 更新审核结果和优先级
	recordIds := make([]uint, len(sqlList))
	for i, sql := range sqlList {
		recordIds[i] = sql.ID
		err = s.UpdateManagerSQLBySqlId(sql.SQLID, map[string]interface{}{"audit_level": sql.AuditLevel, "audit_results": sql.AuditResults, "priority": sql.Priority})
		if err != nil {
			entry.Warnf("update manager sql failed, error: %v", err)
			continue
		}
	}
	// 更新最后审核时间
	err = s.UpdateManageSQLProcessByManageIDs(recordIds, map[string]interface{}{"last_audit_time": time.Now()})
	if err != nil {
		entry.Warnf("update manage record process failed, error: %v", err)
	}
}

func BatchAuditSQLs(l *logrus.Entry, sqlList []*model.SQLManageRecord) ([]*model.SQLManageRecord, error) {
	s := model.GetStorage()
	sqlMap := make(map[string][]*model.SQLManageRecord)

	for _, sql := range sqlList {
		// 根据扫描任务和 schema name 聚合 sqls，避免 task 内需要切换 schema 上下文审核
		key := fmt.Sprintf("%s:%s:%s", sql.SourceId, sql.Source, sql.SchemaName)
		sqlMap[key] = append(sqlMap[key], sql)
	}

	var (
		auditedSQLs []*model.SQLManageRecord
		mu          sync.Mutex
		wg          sync.WaitGroup
	)

	for _, sqls := range sqlMap {
		wg.Add(1)
		go func(sqls []*model.SQLManageRecord) {
			defer wg.Done()

			sourceType := sqls[0].Source
			sourceId := sqls[0].SourceId
			l.Infof("processing audit for source id %s, type %s", sourceId, sourceType)

			var resp *AuditResultResp
			meta, err := GetMeta(sourceType)
			if err != nil {
				l.Errorf("get meta to audit sql fail %v", err)
			} else {
				resp, err = meta.Handler.Audit(sqls)
				if err != nil {
					if errors.Is(err, model.ErrAuditPlanNotFound) {
						l.Warnf("audit sqls in task fail %v, can't find source", err)
						// TODO 调整至clean中清理未关联扫描任务的sql
						// 扫描任务已被删除的sql不需要save到管控中
						if err := s.DeleteSQLManageRecordBySourceId(sourceId); err != nil {
							l.Errorf("delete sql manage record fail %v", err)
						}
						return
					}
					l.Errorf("audit sqls in task fail %v, ignore audit result", err)
				}
			}
			mu.Lock()
			// 当审核结果不为nil时，从审核结果中获取sql，其他错误情况直接返回原始sql
			if resp != nil {
				auditedSQLs = append(auditedSQLs, resp.AuditedSqls...)
			} else {
				auditedSQLs = append(auditedSQLs, sqls...)
			}
			mu.Unlock()
		}(sqls)
	}

	wg.Wait()

	return auditedSQLs, nil
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
		priority, _, err := GetSingleSQLPriorityWithReasons(context.TODO(), auditPlan, sql_)
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

// 获取SQL的优先级以及优先级触发的原因，只有高优先级或者无优先级，若是高优先级，则返回model.PriorityHigh=high,如果无优先级则返回空字符串
func GetSingleSQLPriorityWithReasons(ctx context.Context, auditPlan *model.AuditPlanV2, sql *model.SQLManageRecord) (priority string, reasons []string, err error) {
	if auditPlan == nil || sql == nil {
		return "", reasons, nil
	}
	info, err := sql.Info.OriginValue()
	if err != nil {
		return "", nil, err
	}
	toAuditLevel := func(valueToBeCompared string) string {
		switch valueToBeCompared {
		case "0":
			return locale.Bundle.LocalizeMsgByCtx(ctx, locale.RuleLevelNormal)
		case "1":
			return locale.Bundle.LocalizeMsgByCtx(ctx, locale.RuleLevelNotice)
		case "2":
			return locale.Bundle.LocalizeMsgByCtx(ctx, locale.RuleLevelWarn)
		case "3":
			return locale.Bundle.LocalizeMsgByCtx(ctx, locale.RuleLevelError)
		default:
			return valueToBeCompared
		}
	}
	highPriorityConditions := auditPlan.HighPriorityParams
	// 遍历优先级条件
	for _, highPriorityCondition := range highPriorityConditions {
		var valueToBeCompared string
		// 特殊处理审核级别
		if highPriorityCondition.Key == OperationParamAuditLevel {
			switch sql.AuditLevel {
			case string(driverV2.RuleLevelNotice):
				valueToBeCompared = "1"
			case string(driverV2.RuleLevelWarn):
				valueToBeCompared = "2"
			case string(driverV2.RuleLevelError):
				valueToBeCompared = "3"
			case string(driverV2.RuleLevelNormal):
				valueToBeCompared = "0"
			default:
				valueToBeCompared = "0"
			}
		} else {
			// 获取信息中的相应字段值
			infoV, ok := info[highPriorityCondition.Key]
			if !ok {
				continue
			}
			valueToBeCompared = fmt.Sprintf("%v", infoV)
		}
		// 检查是否为高优先级条件
		if high, err := highPriorityConditions.CompareParamValue(highPriorityCondition.Key, valueToBeCompared); err == nil && high {
			// 添加匹配的条件作为原因
			comparedValue := highPriorityCondition.Param.Value
			if highPriorityCondition.Key == OperationParamAuditLevel {
				valueToBeCompared = toAuditLevel(valueToBeCompared)
				comparedValue = toAuditLevel(comparedValue)
			}
			reasons = append(reasons, fmt.Sprintf("【%v %v %v，%s %s】",
				highPriorityCondition.I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx)),
				highPriorityCondition.Operator.Value,
				comparedValue,
				locale.Bundle.LocalizeMsgByCtx(ctx, locale.WordIs), valueToBeCompared),
			)
		}
	}
	if len(reasons) > 0 {
		return model.PriorityHigh, reasons, nil
	}
	return "", reasons, nil
}
