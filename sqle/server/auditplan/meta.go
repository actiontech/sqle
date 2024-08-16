package auditplan

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/sirupsen/logrus"
)

type Meta struct {
	Type         string `json:"audit_plan_type"`
	Desc         string `json:"audit_plan_type_desc"`
	InstanceType string `json:"instance_type"`
	// instanceId means gen `enums` by db conn, default is a constant definition
	Params             func(instanceId ...string) params.Params `json:"audit_plan_params,omitempty"`
	HighPriorityParams params.ParamsWithOperator                `json:"high_priority_params,omitempty"`
	Metrics            []string
	CreateTask         func(entry *logrus.Entry, ap *AuditPlan) Task `json:"-"`
	Handler            AuditPlanHandler
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
	paramKeyCollectIntervalSecond              = "collect_interval_second"
	paramKeyCollectIntervalMinute              = "collect_interval_minute"
	paramKeySQLMinSecond                       = "sql_min_second"
	paramKeyDBInstanceId                       = "db_instance_id"
	paramKeyAccessKeyId                        = "access_key_id"
	paramKeyAccessKeySecret                    = "access_key_secret"
	paramKeyRdsPath                            = "rds_path"
	paramKeyFirstSqlsScrappedInLastPeriodHours = "first_sqls_scrapped_in_last_period_hours"
	paramKeyProjectId                          = "project_id"
	paramKeyRegion                             = "region"
)

const (
	OperationParamAuditLevel     = "audit_level"
	OperationParamQueryTimeAvg   = MetricNameQueryTimeAvg
	OperationParamRowExaminedAvg = MetricNameRowExaminedAvg
)

var MetaBuilderList = []MetaBuilder{
	{
		Type:          TypeDefault,
		Desc:          "自定义",
		TaskHandlerFn: NewDefaultTaskV2Fn(),
	},
	{
		Type:          TypeMySQLSchemaMeta,
		Desc:          "库表元数据",
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
		Params: func(instanceId ...string) params.Params {
			return taskMeta.Params(instanceId...)
		},
		HighPriorityParams: taskMeta.HighPriorityParams(),
		Metrics:            taskMeta.Metrics(),
		Handler:            handler,
		CreateTask:         NewTaskWrap(b.TaskHandlerFn),
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
		Type:               meta.Type,
		Desc:               meta.Desc,
		InstanceType:       meta.InstanceType,
		Params:             meta.Params,
		HighPriorityParams: meta.HighPriorityParams,
		Metrics:            meta.Metrics,
		CreateTask:         meta.CreateTask,
		Handler:            meta.Handler,
	}, nil
}

var supportedCmdTypeList = map[string]struct{}{
	TypeMySQLSlowLog: {},
	TypeTiDBAuditLog: {},
}

func GetSupportedScannerAuditPlanType() map[string]struct{} {
	return supportedCmdTypeList
}
