package auditplan

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/oracle"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/sirupsen/logrus"
)

type Meta struct {
	Type         string                                              `json:"audit_plan_type"`
	Desc         string                                              `json:"audit_plan_type_desc"`
	InstanceType string                                              `json:"instance_type"`
	Params       params.Params                                       `json:"audit_plan_params,omitempty"`
	CreateTask   func(entry *logrus.Entry, ap *model.AuditPlan) Task `json:"-"`
}

const (
	TypeDefault            = "default"
	TypeMySQLSlowLog       = "mysql_slow_log"
	TypeMySQLMybatis       = "mysql_mybatis"
	TypeMySQLSchemaMeta    = "mysql_schema_meta"
	TypeAliRdsMySQLSlowLog = "ali_rds_mysql_slow_log"
	TypeOracleTopSQL       = "oracle_top_sql"
	TypeTiDBAuditLog       = "tidb_audit_log"
	TypeAllAppExtract      = "all_app_extract"
)

const (
	InstanceTypeAll    = ""
	InstanceTypeMySQL  = "MySQL"
	InstanceTypeOracle = "Oracle"
	InstanceTypeTiDB   = "TiDB"
)

const (
	paramKeyCollectIntervalMinute               = "collect_interval_minute"
	paramKeyAuditSQLsScrappedInLastPeriodMinute = "audit_sqls_scrapped_in_last_period_minute"
	paramKeyDBInstanceId                        = "db_instance_id"
	paramKeyAccessKeyId                         = "access_key_id"
	paramKeyAccessKeySecret                     = "access_key_secret"
	paramKeyFirstSqlsScrappedInLastPeriodHours  = "first_sqls_scrapped_in_last_period_hours"
)

var Metas = []Meta{
	{
		Type:         TypeDefault,
		Desc:         "自定义",
		InstanceType: InstanceTypeAll,
		CreateTask:   NewDefaultTask,
	},
	{
		Type:         TypeMySQLSlowLog,
		Desc:         "慢日志",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewDefaultTask,
		Params: []*params.Param{
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
		},
	},
	{
		Type:         TypeMySQLMybatis,
		Desc:         "Mybatis 扫描",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewDefaultTask,
	},
	{
		Type:         TypeMySQLSchemaMeta,
		Desc:         "库表元数据",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewSchemaMetaTask,
		Params: []*params.Param{
			&params.Param{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			&params.Param{
				Key:   "collect_view",
				Desc:  "是否采集视图信息",
				Value: "0",
				Type:  params.ParamTypeBool,
			},
		},
	},
	{
		Type:         TypeAliRdsMySQLSlowLog,
		Desc:         "阿里RDS MySQL慢日志",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewAliRdsMySQLSlowLogTask,
		Params: []*params.Param{
			{
				Key:   paramKeyDBInstanceId,
				Desc:  "实例ID",
				Value: "",
				Type:  params.ParamTypeString,
			},
			{
				Key:   paramKeyAccessKeyId,
				Desc:  "Access Key ID",
				Value: "",
				Type:  params.ParamTypeString,
			},
			{
				Key:   paramKeyAccessKeySecret,
				Desc:  "Access Key Secret",
				Value: "",
				Type:  params.ParamTypeString,
			},
			{
				Key:   paramKeyFirstSqlsScrappedInLastPeriodHours,
				Desc:  "启动任务时拉取慢日志时间范围(单位:小时,最大31天)",
				Value: "",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
		},
	},
	{
		Type:         TypeOracleTopSQL,
		Desc:         "Oracle TOP SQL",
		InstanceType: InstanceTypeOracle,
		CreateTask:   NewOracleTopSQLTask,
		Params: []*params.Param{
			{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   "top_n",
				Desc:  "Top N",
				Value: "3",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   "order_by_column",
				Desc:  "V$SQLAREA中的排序字段",
				Value: oracle.DynPerformanceViewSQLAreaColumnElapsedTime,
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeAllAppExtract,
		Desc:         "应用程序SQL抓取",
		InstanceType: InstanceTypeAll,
		CreateTask:   NewDefaultTask,
	},
	{
		Type:         TypeTiDBAuditLog,
		Desc:         "TiDB审计日志",
		InstanceType: InstanceTypeTiDB,
		CreateTask:   NewTiDBAuditLogTask,
		Params: []*params.Param{
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
		},
	},
}

var MetaMap = map[string]Meta{}

func init() {
	for _, meta := range Metas {
		MetaMap[meta.Type] = meta
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
	}, nil
}
