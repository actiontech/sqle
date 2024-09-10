package auditplan

import (
	"context"
	"strconv"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

type DefaultTaskV2 struct{}

func NewDefaultTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &DefaultTaskV2{}
	}
}

func (at *DefaultTaskV2) InstanceType() string {
	return InstanceTypeAll
}

func (at *DefaultTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{}
}

var defaultOperatorEnums = []params.EnumsValue{
	{
		Value:    ">",
		I18nDesc: locale.ShouldLocalizeAll(locale.OperatorGreaterThan),
	}, {
		Value:    "=",
		I18nDesc: locale.ShouldLocalizeAll(locale.OperatorEqualTo),
	}, {
		Value:    "<",
		I18nDesc: locale.ShouldLocalizeAll(locale.OperatorLessThan),
	},
}

var defaultAuditLevelOperateParams = &params.ParamWithOperator{
	Param: params.Param{
		Key:   OperationParamAuditLevel,
		Value: "2",
		Type:  params.ParamTypeInt,
		Enums: []params.EnumsValue{
			{
				Value: "1",
				Desc:  string(driverV2.RuleLevelNotice),
			}, {
				Value: "2",
				Desc:  string(driverV2.RuleLevelWarn),
			},
			{
				Value: "3",
				Desc:  string(driverV2.RuleLevelError),
			},
		},
		I18nDesc: locale.ShouldLocalizeAll(locale.OperationParamAuditLevel),
	},
	Operator: params.Operator{
		Value:      ">",
		EnumsValue: defaultOperatorEnums,
	},
}

func (at *DefaultTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{
		defaultAuditLevelOperateParams,
	}
}

func (at *DefaultTaskV2) Metrics() []string {
	return []string{
		MetricNameCounter,
		MetricNameLastReceiveTimestamp,
	}
}

func (at *DefaultTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}
	originSQL.SQLContent = mergedSQL.SQLContent
	// counter
	originCounter := originSQL.Info.Get(MetricNameCounter).Int()
	mergedCounter := mergedSQL.Info.Get(MetricNameCounter).Int()
	counter := originCounter + mergedCounter
	originSQL.Info.SetInt(MetricNameCounter, counter)

	// last_receive_timestamp
	originSQL.Info.SetString(MetricNameLastReceiveTimestamp, mergedSQL.Info.Get(MetricNameLastReceiveTimestamp).String())
	return
}

func (at *DefaultTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
	originSQL, exist, err := cache.GetSQL(sql.SQLId)
	if err != nil {
		return err
	}
	if !exist {
		cache.CacheSQL(sql)
		return nil
	}
	at.mergeSQL(originSQL, sql)
	return nil
}

func (at *DefaultTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *DefaultTaskV2) Head(ap *AuditPlan) []Head {
	return []Head{
		{
			Name: "fingerprint",
			Desc: locale.ApSQLFingerprint,
			Type: "sql",
		},
		{
			Name: "sql",
			Desc: locale.ApLastSQL,
			Type: "sql",
		},
		{
			Name: "priority",
			Desc: locale.ApPriority,
		},
		{
			Name: model.AuditResultName,
			Desc: model.AuditResultDesc,
		},
		{
			Name: MetricNameCounter,
			Desc: locale.ApMetricNameCounter,
		},
		{
			Name: MetricNameLastReceiveTimestamp,
			Desc: locale.ApMetricNameLastReceiveTimestamp,
		},
	}
}

func (at *DefaultTaskV2) Filters(ctx context.Context, logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta {
	return []FilterMeta{
		{
			Name:            "sql", // 模糊筛选
			Desc:            locale.ApSQLStatement,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
		},
		{
			Name:            "rule_name",
			Desc:            locale.ApRuleName,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerRuleTips(ctx, logger, ap.ID, persist),
		},
		{
			Name:            "priority",
			Desc:            locale.ApPriority,
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerPriorityTips(ctx, logger),
		}}
}

func genArgsByFilters(filters []Filter) map[model.FilterName]interface{} {
	args := make(map[model.FilterName]interface{}, len(filters))
	for _, filter := range filters {
		switch filter.Name {
		case "sql":
			args[model.FilterSQL] = filter.FilterComparisonValue

		case "priority":
			args[model.FilterPriority] = filter.FilterComparisonValue

		case "rule_name":
			args[model.FilterRuleName] = filter.FilterComparisonValue

		case "schema_name":
			args[model.FilterSchemaName] = filter.FilterComparisonValue

		case MetricNameDBUser:
			args[model.FilterNameDBUser] = filter.FilterComparisonValue

		case MetricNameCounter:
			counter, err := strconv.Atoi(filter.FilterComparisonValue)
			if err != nil {
				continue
			}
			args[model.FilterCounter] = counter

		case MetricNameQueryTimeAvg:
			value, err := strconv.Atoi(filter.FilterComparisonValue)
			if err != nil {
				continue
			}
			args[model.FilterQueryTimeAvg] = value

		case MetricNameRowExaminedAvg:
			value, err := strconv.Atoi(filter.FilterComparisonValue)
			if err != nil {
				continue
			}
			args[model.FilterRowExaminedAvg] = value

		case MetricNameLastReceiveTimestamp:
			args[model.FilterLastReceiveTimestampFrom] = filter.FilterBetweenValue.From
			args[model.FilterLastReceiveTimestampTo] = filter.FilterBetweenValue.To
		}
	}
	return args
}

func (at *DefaultTaskV2) GetSQLData(ctx context.Context, ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
	args := make(map[model.FilterName]interface{}, len(filters))
	for _, filter := range filters {
		switch filter.Name {
		case "sql":
			args[model.FilterSQL] = filter.FilterComparisonValue

		case "priority":
			args[model.FilterPriority] = filter.FilterComparisonValue

		case "rule_name":
			args[model.FilterRuleName] = filter.FilterComparisonValue
		}
	}
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReqV2(ap.ID, ap.Type, limit, offset, checkAndGetOrderByName(at.Head(ap), orderBy), isAsc, args)
	if err != nil {
		return nil, count, err
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		data, err := sql.Info.OriginValue()
		if err != nil {
			return nil, 0, err
		}
		info := LoadMetrics(data, at.Metrics())
		rows = append(rows, map[string]string{
			"sql":                          sql.SQLContent,
			"fingerprint":                  sql.Fingerprint,
			"id":                           sql.AuditPlanSqlId,
			"priority":                     sql.Priority.String,
			MetricNameCounter:              strconv.Itoa(int(info.Get(MetricNameCounter).Int())),
			MetricNameLastReceiveTimestamp: info.Get(MetricNameLastReceiveTimestamp).String(),
			model.AuditResultName:          sql.AuditResult.GetAuditJsonStrByLangTag(locale.GetLangTagFromCtx(ctx)),
		})
	}
	return rows, count, nil
}

func ShowSchemaEnumsByInstanceId(instanceId string) (enumsValues []params.EnumsValue) {
	logger := log.NewEntry()
	if instanceId == "" {
		return
	}
	inst, exist, err := dms.GetInstancesById(context.Background(), instanceId)
	if err != nil {
		logger.Errorf("can't find instance by id, %v", instanceId)
		return
	}
	if !exist {
		return
	}
	if !driver.GetPluginManager().IsOptionalModuleEnabled(inst.DbType, driverV2.OptionalModuleQuery) {
		logger.Errorf("can not do this task, %v", driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleQuery))
		return
	}

	plugin, err := driver.GetPluginManager().OpenPlugin(logger, inst.DbType, &driverV2.Config{
		DSN: &driverV2.DSN{
			Host:             inst.Host,
			Port:             inst.Port,
			User:             inst.User,
			Password:         inst.Password,
			AdditionalParams: inst.AdditionalParams,
		},
	})
	if err != nil {
		logger.Errorf("get plugin failed, error: %v", err)
		return
	}
	defer plugin.Close(context.Background())

	schemas, err := plugin.Schemas(context.Background())
	if err != nil {
		logger.Errorf("show schema failed, error: %v", err)
		return
	}

	for _, schema := range schemas {
		enumsValues = append(enumsValues, params.EnumsValue{
			Value: schema,
			Desc:  schema,
		})
	}
	return
}
