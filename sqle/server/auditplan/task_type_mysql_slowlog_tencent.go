package auditplan

import (
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

// MySQLSlowLogTencentTaskV2 腾讯云 CDB MySQL 慢日志（OpenAPI）扫描类型。
type MySQLSlowLogTencentTaskV2 struct {
	DefaultTaskV2
	// 测试可注入；生产为 nil，走 OpenAPI 客户端。
	collectDeps *tencentCollectDeps
}

func NewMySQLSlowLogTencentTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &MySQLSlowLogTencentTaskV2{}
	}
}

func (at *MySQLSlowLogTencentTaskV2) InstanceType() string {
	return InstanceTypeMySQL
}

func (at *MySQLSlowLogTencentTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:      paramKeyAccessKeyId,
			Value:    "",
			Type:     params.ParamTypePassword,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamTencentSecretId),
		},
		{
			Key:      paramKeyAccessKeySecret,
			Value:    "",
			Type:     params.ParamTypePassword,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamTencentSecretKey),
		},
		{
			Key:      paramKeyRegion,
			Value:    "",
			Type:     params.ParamTypeString,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamTencentRegion),
		},
		{
			Key:      paramKeyDBInstanceId,
			Value:    "",
			Type:     params.ParamTypeString,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamTencentDBInstanceId),
		},
		{
			Key:      paramKeyFirstSqlsScrappedInLastPeriodHours,
			Value:    "24",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAllWithArgs(locale.ParamFirstCollectDurationWithMaxDays, 30),
		},
		{
			Key:      paramKeyCollectIntervalMinute,
			Value:    "60",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.Bundle.LocalizeAll(locale.ParamCollectIntervalMinute),
		},
	}
}

func (at *MySQLSlowLogTencentTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameLastReceiveTimestamp,
		MetricNameQueryTimeAvg,
		MetricNameQueryTimeMax,
		MetricNameLockTimeAvg,
		MetricNameLockTimeMax,
		MetricNameRowExaminedAvg,
		MetricNameDBUser,
		MetricNameHost,
		MetricNameStartTimeOfLastScrapedSQL,
	}
}

func (at *MySQLSlowLogTencentTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

// ExtractSQL 方案 A：拉取 + 入队成功后写 DB 游标，返回空切片避免 Wrapper 二次入队。
func (at *MySQLSlowLogTencentTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	deps := tencentCollectDeps{}
	if at.collectDeps != nil {
		deps = *at.collectDeps
	}
	return extractTencentMySQLSlowLog(logger, ap, persist, deps)
}
