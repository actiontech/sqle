package auditplan

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/sirupsen/logrus"
)

type MySQLSchemaMetaTaskV2 struct {
	BaseSchemaMetaTaskV2
}

func NewMySQLSchemaMetaTaskV2Fn() func() interface{} {
	return func() interface{} {
		return NewMySQLSchemaMetaTaskV2()
	}
}

func NewMySQLSchemaMetaTaskV2() *MySQLSchemaMetaTaskV2 {
	t := &MySQLSchemaMetaTaskV2{}
	t.BaseSchemaMetaTaskV2 = BaseSchemaMetaTaskV2{
		extract: t.extractSQL,
	}
	return t
}

func (at *MySQLSchemaMetaTaskV2) InstanceType() string {
	return InstanceTypeMySQL
}

func (at *BaseSchemaMetaTaskV2) extractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SchemaMetaSQL, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}

	sqls := []*SchemaMetaSQL{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	instance, _, err := dms.GetInstancesById(ctx, ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
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

	schemas, err := db.ShowDatabases(true)
	if err != nil {
		return nil, fmt.Errorf("get schema fail, error: %v", err)
	}

	for _, schema := range schemas {
		tables, err := db.ShowSchemaTables(schema)
		if err != nil {
			return nil, fmt.Errorf("get schema table fail, error: %v", err)
		}
		var views []string
		if ap.Params.GetParam("collect_view").Bool() {
			views, err = db.ShowSchemaViews(schema)
			if err != nil {
				return nil, fmt.Errorf("get schema view fail, error: %v", err)
			}
		}
		err = db.UseSchema(schema)
		if err != nil {
			return nil, fmt.Errorf("use schema fail, error: %v", err)
		}
		for _, table := range tables {
			sql, err := db.ShowCreateTable(utils.SupplementalQuotationMarks(schema), utils.SupplementalQuotationMarks(table))
			if err != nil {
				return nil, fmt.Errorf("show create table fail, error: %v", err)
			}
			sqls = append(sqls, &SchemaMetaSQL{
				SchemaName: schema,
				MetaName:   table,
				MetaType:   "table",
				SQLContent: sql,
			})
		}
		for _, view := range views {
			sql, err := db.ShowCreateView(utils.SupplementalQuotationMarks(view))
			if err != nil {
				return nil, fmt.Errorf("show create view fail, error: %v", err)
			}
			sqls = append(sqls, &SchemaMetaSQL{
				SchemaName: schema,
				MetaName:   view,
				MetaType:   "view",
				SQLContent: sql,
			})
		}
	}
	return sqls, nil
}

type BaseSchemaMetaTaskV2 struct {
	extract func(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SchemaMetaSQL, error)
}

type SchemaMetaSQL struct {
	SchemaName string
	MetaName   string
	MetaType   string
	SQLContent string
}

func (at *BaseSchemaMetaTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:   paramKeyCollectIntervalMinute,
			Desc:  "采集周期（分钟）",
			Value: "60",
			Type:  params.ParamTypeInt,
		},
		{
			Key:   "collect_view",
			Desc:  "是否采集视图信息",
			Value: "0",
			Type:  params.ParamTypeBool,
		},
	}
}

func (at *BaseSchemaMetaTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{
		{
			Param: params.Param{
				Key:   "audit_level",
				Value: "warn",
				Desc:  "告警级别",
				Type:  params.ParamTypeInt,
				Enums: []params.EnumsValue{
					{
						Value: "0",
						Desc:  "notice",
					}, {
						Value: "1",
						Desc:  "warn",
					}, {
						Value: "2",
						Desc:  "error",
					},
				},
			},
			BooleanOperatorParam: params.BooleanOperator{
				Value: ">",
				EnumsValue: []params.EnumsValue{
					{
						Value: ">",
						Desc:  "大于",
					}, {
						Value: "=",
						Desc:  "等于",
					}, {
						Value: "<",
						Desc:  "小于",
					},
				},
			},
		},
	}
}

func (at *BaseSchemaMetaTaskV2) Metrics() []string {
	return []string{
		MetricNameMetaType,
		MetricNameMetaName,
		MetricNameRecordDeleted,
	}
}

func (at *BaseSchemaMetaTaskV2) genSQLId(sql *SQLV2) string {
	md5Json, err := json.Marshal(
		struct {
			ProjectId string
			Schema    string
			InstName  string
			Source    string
			ApID      uint
			MetaName  string
			MetaType  string
		}{
			ProjectId: sql.ProjectId,
			Schema:    sql.SchemaName,
			InstName:  sql.InstanceID,
			Source:    sql.Source,
			ApID:      sql.SourceId,
			MetaName:  sql.Info.Get(MetricNameMetaName).String(),
			MetaType:  sql.Info.Get(MetricNameMetaType).String(),
		},
	)
	if err != nil { // todo: 处理错误
		return utils.Md5String(fmt.Sprintf("%d:%s:%s:%s", sql.SourceId, sql.SchemaName, sql.Info.Get(MetricNameMetaName).String(), sql.Info.Get(MetricNameMetaType).String()))
	} else {
		return utils.Md5String(string(md5Json))
	}
}

func (at *BaseSchemaMetaTaskV2) mergeSQL(originSQL, mergedSQL *SQLV2) {
	if originSQL.SQLId != mergedSQL.SQLId {
		return
	}

	originSQL.SQLContent = mergedSQL.SQLContent
	originSQL.Fingerprint = mergedSQL.Fingerprint

	// meta type
	originSQL.Info.SetString(MetricNameMetaType, mergedSQL.Info.Get(MetricNameMetaType).String())

	// meta name
	originSQL.Info.SetString(MetricNameMetaName, mergedSQL.Info.Get(MetricNameMetaName).String())

	// record deleted
	originSQL.Info.SetBool(MetricNameRecordDeleted, mergedSQL.Info.Get(MetricNameRecordDeleted).Bool())

	return
}

func (at *BaseSchemaMetaTaskV2) ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) {
	if ap.InstanceID == "" {
		logger.Warnf("instance is not configured")
		return nil, nil
	}

	sqls, err := at.extract(logger, ap, persist)
	if err != nil {
		return nil, err
	}
	cache := NewSQLV2Cache()
	// 1. 取出当期任务的存量数据，先将所有存量数据都标记为 MetricNameRecordDeleted = true(表示在扫描任务界面删除、但SQL管控内还可见)
	// 2. 再将采集的数据和存量数据进行合并，所有还存在的数据还会被标记为 MetricNameRecordDeleted = false
	// 例子：
	// 1. 假设存在表 db1.t1, db1.t2
	// 2. 进行采集数据：db1.t1, db1.t2 被采集到了扫描任务，此时它们的 MetricNameRecordDeleted 指标都为 false
	// 3. 在下一次采集之前将 db1.t1 表删除
	// 4. 进行采集数据：最终 db1.t1 会被标记为 MetricNameRecordDeleted = false, db1.t2 会被标记为 MetricNameRecordDeleted = true
	originSQLV2, err := persist.GetManagerSQLListByAuditPlanId(ap.ID)
	if err != nil {
		logger.Errorf("get manager sql failed, error: %v", err)
		return nil, err
	}
	for _, sql := range originSQLV2 {
		sqlV2 := ConvertMangerSQLToSQLV2(sql)
		if sqlV2.Info.Get(MetricNameRecordDeleted).Bool() == false {
			sqlV2.Info.SetBool(MetricNameRecordDeleted, true)
			cache.CacheSQL(sqlV2)
		}
	}
	for _, sql := range sqls {
		sqlV2 := &SQLV2{
			Source:      ap.Type,
			SourceId:    ap.ID,
			ProjectId:   ap.ProjectId,
			InstanceID:  ap.InstanceID,
			SchemaName:  sql.SchemaName,
			SQLContent:  sql.SQLContent,
			Fingerprint: sql.SQLContent, // DDL不能参数化
			Info:        NewMetrics(),
		}
		sqlV2.Info.SetBool(MetricNameRecordDeleted, false)
		sqlV2.Info.SetString(MetricNameMetaName, sql.MetaName)
		sqlV2.Info.SetString(MetricNameMetaType, sql.MetaType)
		sqlV2.SQLId = at.genSQLId(sqlV2)
		if err := at.AggregateSQL(cache, sqlV2); err != nil {
			logger.Errorf("aggregate sql failed, error: %v", err)
			continue
		}
	}
	return cache.GetSQLs(), nil
}

func (at *BaseSchemaMetaTaskV2) AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error {
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

func (at *BaseSchemaMetaTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return auditSQLs(sqls)
}

func (at *BaseSchemaMetaTaskV2) Head(ap *AuditPlan) []Head {
	return []Head{
		{
			Name: "sql",
			Desc: "SQL语句",
			Type: "sql",
		},
		{
			Name: model.AuditResultName,
			Desc: model.AuditResultDesc,
		},
		{
			Name: "schema_name",
			Desc: "schema",
		},
		{
			Name: MetricNameMetaName,
			Desc: "对象名称",
		},
		{
			Name: MetricNameMetaType,
			Desc: "对象类型",
		},
	}
}

func (at *BaseSchemaMetaTaskV2) Filters(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta {
	return []FilterMeta{
		{
			Name:            "sql", // 模糊筛选
			Desc:            "SQL",
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
		},
		{
			Name:            "rule_name",
			Desc:            "审核规则",
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerRuleTips(logger, ap.ID, persist),
		},
		{
			Name:            "schema_name",
			Desc:            "schema",
			FilterInputType: FilterInputTypeString,
			FilterOpType:    FilterOpTypeEqual,
			FilterTips:      GetSqlManagerSchemaNameTips(logger, ap.ID, persist),
		},
	}
}

func (at *BaseSchemaMetaTaskV2) GetSQLData(ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
	// todo: 需要过滤掉	MetricNameRecordDeleted = true 的记录，因为有分页所以需要在db里过滤，还要考虑概览界面统计的问题
	args := make(map[model.FilterName]interface{}, len(filters))
	for _, filter := range filters {
		switch filter.Name {
		case "sql":
			args[model.FilterSQL] = filter.FilterComparisonValue

		case "schema_name":
			args[model.FilterSchemaName] = filter.FilterComparisonValue

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
			"sql":                 sql.SQLContent,
			"schema_name":         sql.Schema,
			MetricNameMetaName:    info.Get(MetricNameMetaName).String(),
			MetricNameMetaType:    info.Get(MetricNameMetaType).String(),
			model.AuditResultName: sql.AuditResult.String,
		})
	}
	return rows, count, nil
}
