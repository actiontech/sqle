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
	paramKeySlowLogCollectInput                 = "slow_log_collect_input"
	paramKeyAuditSQLsScrappedInLastPeriodMinute = "audit_sqls_scrapped_in_last_period_minute"
	paramKeySQLMinSecond                        = "sql_min_second"
	paramKeyDBInstanceId                        = "db_instance_id"
	paramKeyAccessKeyId                         = "access_key_id"
	paramKeyAccessKeySecret                     = "access_key_secret"
	paramKeyRdsPath                             = "rds_path"
	paramKeyFirstSqlsScrappedInLastPeriodHours  = "first_sqls_scrapped_in_last_period_hours"
	paramKeyProjectId                           = "project_id"
	paramKeyRegion                              = "region"
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
		Params: []*params.Param{
			{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟，仅对 mysql.slow_log 有效）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeySlowLogCollectInput,
				Desc:  "采集来源。0：mysql-slow.log 文件；1：mysql.slow_log 表",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
		},
		CreateTask: NewSlowLogTask,
	},
	{
		Type:         TypeMySQLMybatis,
		Desc:         "Mybatis 扫描",
		InstanceType: InstanceTypeAll,
		CreateTask:   NewDefaultTask,
	},
	{
		Type:         TypeMySQLSchemaMeta,
		Desc:         "库表元数据",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewSchemaMetaTask,
		Params: []*params.Param{
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
		},
	},
	{
		Type:         TypeMySQLProcesslist,
		Desc:         "processlist 列表",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewMySQLProcesslistTask,
		Params: []*params.Param{
			{
				Key:   paramKeyCollectIntervalSecond,
				Desc:  "采集周期（秒）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeySQLMinSecond,
				Desc:  "SQL 最小执行时间（秒）",
				Value: "0",
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
			{
				Key:   paramKeyRdsPath,
				Desc:  "RDS Open API地址",
				Value: "rds.aliyuncs.com",
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeAliRdsMySQLAuditLog,
		Desc:         "阿里RDS MySQL审计日志",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewAliRdsMySQLAuditLogTask,
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
				Desc:  "启动任务时拉取日志时间范围(单位:小时,最大31天)",
				Value: "",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyRdsPath,
				Desc:  "RDS Open API地址",
				Value: "rds.aliyuncs.com",
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeBaiduRdsMySQLSlowLog,
		Desc:         "百度云RDS MySQL慢日志",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewBaiduRdsMySQLSlowLogTask,
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
				Key: paramKeyFirstSqlsScrappedInLastPeriodHours,
				// 百度云RDS慢日志只能拉取最近7天的数据
				// https://cloud.baidu.com/doc/RDS/s/Tjwvz046g
				Desc:  "启动任务时拉取慢日志时间范围(单位:小时,最大7天)",
				Value: "",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyRdsPath,
				Desc:  "RDS Open API地址",
				Value: "rds.bj.baidubce.com",
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeHuaweiRdsMySQLSlowLog,
		Desc:         "华为云RDS MySQL慢日志",
		InstanceType: InstanceTypeMySQL,
		CreateTask:   NewHuaweiRdsMySQLSlowLogTask,
		Params: []*params.Param{
			{
				Key:   paramKeyProjectId,
				Desc:  "项目ID",
				Value: "",
				Type:  params.ParamTypeString,
			},
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
				Desc:  "启动任务时拉取慢日志的时间范围（单位：小时，最大30天）",
				Value: "",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyRegion,
				Desc:  "当前RDS实例所在的地区（示例：cn-east-2）",
				Value: "",
				Type:  params.ParamTypeString,
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
	{
		Type:         TypeSQLFile,
		Desc:         "SQL文件",
		InstanceType: InstanceTypeAll,
		CreateTask:   NewDefaultTask,
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
		CreateTask:   meta.CreateTask,
	}, nil
}
