package auditplan

import (
	"context"

	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

func GetSqlManagerRuleTips(ctx context.Context, logger *logrus.Entry, auditPlanId uint, persist *model.Storage) []FilterTip {
	var ruleFilterTips []FilterTip
	rules, err := persist.GetManagerSqlRuleTipsByAuditPlan(auditPlanId)
	if err != nil {
		logger.Warnf("get sql manager rule tips failed, error: %v", err)
		return []FilterTip{}
	} else {
		ruleFilterTips = make([]FilterTip, 0, len(rules))
		for _, rule := range rules {
			ruleFilterTips = append(ruleFilterTips, FilterTip{
				Value: rule.RuleName,
				Desc:  rule.I18nRuleInfo.GetRuleInfoByLangTag(locale.GetLangTagFromCtx(ctx)).Desc,
				Group: rule.DbType,
			})
		}
	}
	return ruleFilterTips
}

func GetSqlManagerSchemaNameTips(logger *logrus.Entry, auditPlanId uint, persist *model.Storage) []FilterTip {
	schemaNames, err := persist.GetManagerSqlSchemaNameByAuditPlan(auditPlanId)
	if err != nil {
		logger.Warnf("get sql manager schema name tips failed, error: %v", err)
		return []FilterTip{}
	} else {
		schemaNameTips := make([]FilterTip, 0, len(schemaNames))
		for _, schemaName := range schemaNames {
			schemaNameTips = append(schemaNameTips, FilterTip{
				Value: schemaName,
				Desc:  schemaName,
			})
		}
		return schemaNameTips
	}
}

func GetSqlManagerMetricTips(logger *logrus.Entry, auditPlanId uint, persist *model.Storage, metricName string) []FilterTip {
	metricNames, err := persist.GetManagerSqlMetricTipsByAuditPlan(auditPlanId, metricName)
	if err != nil {
		logger.Warnf("get sql manager metric tips failed, error: %v", err)
		return []FilterTip{}
	} else {
		metricNameTips := make([]FilterTip, 0, len(metricNames))
		for _, metricName := range metricNames {
			metricNameTips = append(metricNameTips, FilterTip{
				Value: metricName,
				Desc:  metricName,
			})
		}
		return metricNameTips
	}
}

func GetSqlManagerPriorityTips(ctx context.Context, logger *logrus.Entry) []FilterTip {
	return []FilterTip{
		{
			Value: model.PriorityHigh,
			Desc:  locale.ShouldLocalizeMsg(ctx, locale.ApPriorityHigh),
		},
	}
}
