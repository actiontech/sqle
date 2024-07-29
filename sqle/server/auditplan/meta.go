package auditplan

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/sirupsen/logrus"
)

type Meta struct {
	Type         string        `json:"audit_plan_type"`
	Desc         string        `json:"audit_plan_type_desc"`
	InstanceType string        `json:"instance_type"`
	Params       params.Params `json:"audit_plan_params,omitempty"`
	Metrics      []string
	CreateTask   func(entry *logrus.Entry, ap *AuditPlan) Task `json:"-"`
	Handler      AuditPlanHandler
}

type MetaBuilder struct {
	Type          string
	Desc          string
	TaskHandlerFn func() interface{}
}

const (
	TypeDefault               = "default"
	TypeMySQLSlowLog          = "mysql_slow_log"
	TypeMySQLMybatis          = "mysql_mybatis"
	TypeMySQLSchemaMeta       = "mysql_schema_meta"
	TypeMySQLProcesslist      = "mysql_processlist"
	TypeAliRdsMySQLSlowLog    = "ali_rds_mysql_slow_log"
	TypeAliRdsMySQLAuditLog   = "ali_rds_mysql_audit_log"
	TypeHuaweiRdsMySQLSlowLog = "huawei_rds_mysql_slow_log"
	TypeOracleTopSQL          = "oracle_top_sql"
	TypeTiDBAuditLog          = "tidb_audit_log"
	TypeAllAppExtract         = "all_app_extract"
	TypeBaiduRdsMySQLSlowLog  = "baidu_rds_mysql_slow_log"
	TypeSQLFile               = "sql_file"
)

const (
	InstanceTypeAll    = ""
	InstanceTypeMySQL  = "MySQL"
	InstanceTypeOracle = "Oracle"
	InstanceTypeTiDB   = "TiDB"
)

const (
	paramKeyCollectIntervalSecond               = "collect_interval_second"
	paramKeyCollectIntervalMinute               = "collect_interval_minute"
	paramKeyAuditSQLsScrappedInLastPeriodMinute = "audit_sqls_scrapped_in_last_period_minute"
	paramKeySQLMinSecond                        = "sql_min_second"
	paramKeyDBInstanceId                        = "db_instance_id"
	paramKeyAccessKeyId                         = "access_key_id"
	paramKeyAccessKeySecret                     = "access_key_secret"
	paramKeyRdsPath                             = "rds_path"
	paramKeyFirstSqlsScrappedInLastPeriodHours  = "first_sqls_scrapped_in_last_period_hours"
	paramKeyProjectId                           = "project_id"
	paramKeyRegion                              = "region"
	paramKeySchema                              = "schema"
)

var MetaBuilderList = []MetaBuilder{
	{
		Type:          TypeDefault,
		Desc:          "自定义",
		TaskHandlerFn: NewDefaultTaskV2Fn(),
	},
	{
		Type:          TypeMySQLMybatis,
		Desc:          "Mybatis 扫描",
		TaskHandlerFn: NewDefaultTaskV2Fn(),
	},
	{
		Type:          TypeMySQLSchemaMeta,
		Desc:          "MySQL库表元数据",
		TaskHandlerFn: NewMySQLSchemaMetaTaskV2Fn(),
	},
	{
		Type:          TypeMySQLProcesslist,
		Desc:          "processlist 列表",
		TaskHandlerFn: NewMySQLProcessListTaskV2Fn(),
	},
	{
		Type:          TypeAliRdsMySQLSlowLog,
		Desc:          "阿里RDS MySQL慢日志",
		TaskHandlerFn: NewMySQLSlowLogAliTaskV2Fn(),
	},
	{
		Type:          TypeAliRdsMySQLAuditLog,
		Desc:          "阿里RDS MySQL审计日志",
		TaskHandlerFn: NewMySQLAuditLogAliTaskV2Fn(),
	},
	{
		Type:          TypeBaiduRdsMySQLSlowLog,
		Desc:          "百度云RDS MySQL慢日志",
		TaskHandlerFn: NewMySQLSlowLogBaiduTaskV2Fn(),
	},
	{
		Type:          TypeHuaweiRdsMySQLSlowLog,
		Desc:          "华为云RDS MySQL慢日志",
		TaskHandlerFn: NewMySQLSlowLogHuaweiTaskV2Fn(),
	},
	{
		Type:          TypeOracleTopSQL,
		Desc:          "Oracle TOP SQL",
		TaskHandlerFn: NewOracleTopSQLTaskV2Fn(),
	},
	{
		Type:          TypeAllAppExtract,
		Desc:          "应用程序SQL抓取",
		TaskHandlerFn: NewDefaultTaskV2Fn(),
	},
	{
		Type:          TypeTiDBAuditLog,
		Desc:          "TiDB审计日志",
		TaskHandlerFn: NewTiDBAuditLogTaskV2Fn(),
	},
	{
		Type:          TypeSQLFile,
		Desc:          "SQL文件",
		TaskHandlerFn: NewDefaultTaskV2Fn(),
	},
}

var MetaMap = map[string]Meta{}
var Metas = []Meta{}

func buildMeta(b MetaBuilder) Meta {
	task := b.TaskHandlerFn()

	handler, ok := task.(AuditPlanHandler)
	if !ok {
		panic(fmt.Sprintf("task %s don't implement audit plan handler interface, ", b.Type))
	}
	taskMeta, ok := task.(AuditPlanMeta)
	if !ok {
		panic(fmt.Sprintf("task %s don't implement audit plan meta interface, ", b.Type))
	}
	return Meta{
		Type:         b.Type,
		Desc:         b.Desc,
		InstanceType: taskMeta.InstanceType(),
		Params:       taskMeta.Params(),
		Metrics:      taskMeta.Metrics(),
		Handler:      handler,
		CreateTask:   NewTaskWrap(b.TaskHandlerFn),
	}
}

func init() {
	for _, b := range MetaBuilderList {
		meta := buildMeta(b)
		Metas = append(Metas, meta)
		MetaMap[b.Type] = meta
	}
}

func GetMeta(typ string) (Meta, error) {
	if typ == "" {
		typ = TypeDefault
	}
	meta, ok := MetaMap[typ]
	if !ok {
		return Meta{}, fmt.Errorf("audit plan type %s not found", typ)
	}
	return Meta{
		Type:         meta.Type,
		Desc:         meta.Desc,
		InstanceType: meta.InstanceType,
		Params:       meta.Params.Copy(),
		Metrics:      meta.Metrics,
		CreateTask:   meta.CreateTask,
		Handler:      meta.Handler,
	}, nil
}

func GetSupportedScannerAuditPlanType() map[string]struct{} {
	return map[string]struct{}{
		TypeMySQLSlowLog: {},
		TypeTiDBAuditLog: {},
	}
}

const (
	ParamsKeySchema = paramKeySchema
)

func GetEnumsByInstanceId(paramKey, instanceId string) (enumsValues []params.EnumsValue) {
	logger := log.NewEntry()

	switch paramKey {
	case paramKeySchema:
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
	default:
		return
	}
	return

}
