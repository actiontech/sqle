package auditplan

import (
	"context"
	"strconv"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

type DefaultTaskV2 struct {
}

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

func (at *DefaultTaskV2) Audit(sqls []*model.OriginManageSQL) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *DefaultTaskV2) GetSQLs(ap *AuditPlan, persist *model.Storage, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := persist.GetInstanceAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head := []Head{
		{
			Name: "fingerprint",
			Desc: "SQL指纹",
			Type: "sql",
		},
		{
			Name: "sql",
			Desc: "最后一次匹配到该指纹的语句",
			Type: "sql",
		},
		{
			Name: model.AuditResultName,
			Desc: model.AuditResultDesc,
		},
		{
			Name: MetricNameCounter,
			Desc: "匹配到该指纹的语句数量",
		},
		{
			Name: MetricNameLastReceiveTimestamp,
			Desc: "最后一次匹配到该指纹的时间",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		data, err := sql.Info.OriginValue()
		if err != nil {
			return nil, nil, 0, err
		}
		info := LoadMetrics(data, at.Metrics())
		rows = append(rows, map[string]string{
			"sql":                          sql.SQLContent,
			"fingerprint":                  sql.Fingerprint,
			MetricNameCounter:              strconv.Itoa(int(info.Get(MetricNameCounter).Int())),
			MetricNameLastReceiveTimestamp: info.Get(MetricNameLastReceiveTimestamp).String(),
			model.AuditResultName:          sql.AuditResult.String,
		})
	}
	return head, rows, count, nil
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
