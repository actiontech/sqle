//go:build enterprise
// +build enterprise

package auditplan

import "github.com/actiontech/sqle/sqle/pkg/params"

const (
	TypeOceanBaseForMySQLMybatis = "ocean_base_for_mysql_mybatis"
	TypeOceanBaseForMySQLTopSQL  = "ocean_base_for_mysql_top_sql"
	TypeObForOracleTopSQL        = "ob_for_oracle_top_sql"
	TypeDB2TopSQL                = "db2_top_sql"
	TypeDB2SchemaMeta            = "db2_schema_meta"
	TypeTDSQLSlowLog             = "tdsql_for_innodb_slow_log"
	TypeTDSQLSchemaMeta          = "tdsql_for_innodb_schema_meta"
	TypeDmTopSQL                 = "dm_top_sql"
	TypePostgreSQLTopSQL         = "postgresql_top_sql"
	TypePostgreSQLSchemaMeta     = "Postgresql_schema_meta"
	TypeTBaseSlowLog               = "TBase_slow_log"
)

const (
	InstanceTypeOceanBaseForMySQL = "OceanBase For MySQL"
	InstanceTypeObForOracle       = "OceanBase For Oracle"
	InstanceTypeDB2               = "DB2"
	InstanceTypeTDSQL             = "TDSQL For InnoDB"
	InstanceTypeDm                = "DM"
	InstanceTypePostgreSQL        = "PostgreSQL"
	InstanceTypeTBase             = "TBase"
)

const (
	paramKeyIndicator = "indicator"
	paramKeyTopN      = "top_n"
)

var EEMetas = []Meta{
	{
		Type:         TypeTDSQLSlowLog,
		Desc:         "慢日志",
		InstanceType: InstanceTypeTDSQL,
		Params: []*params.Param{
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
		},
		CreateTask: NewSlowLogTask,
	},
	{
		Type:         TypeOceanBaseForMySQLTopSQL,
		Desc:         "Top SQL",
		InstanceType: InstanceTypeOceanBaseForMySQL,
		CreateTask:   NewOBMySQLTopSQLTask,
		Params: []*params.Param{
			{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyTopN,
				Desc:  "Top N",
				Value: "3",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyIndicator,
				Desc:  "关注指标",
				Value: OBMySQLIndicatorElapsedTime,
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeDB2TopSQL,
		Desc:         "DB2 Top SQL",
		InstanceType: InstanceTypeDB2,
		CreateTask:   NewDB2TopSQLTask,
		Params: []*params.Param{
			{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyTopN,
				Desc:  "Top N",
				Value: "3",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyIndicator,
				Desc:  "关注指标",
				Value: DB2IndicatorAverageElapsedTime,
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeDB2SchemaMeta,
		Desc:         "库表元数据",
		InstanceType: InstanceTypeDB2,
		CreateTask:   NewDB2SchemaMetaTask,
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
		Type:         TypeTDSQLSchemaMeta,
		Desc:         "库表元数据",
		InstanceType: InstanceTypeTDSQL,
		CreateTask:   NewSchemaMetaTask,
		Params: []*params.Param{
			{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
		},
	},
	{
		Type:         TypePostgreSQLSchemaMeta,
		Desc:         "库表元数据",
		InstanceType: InstanceTypePostgreSQL,
		CreateTask:   NewPostgreSQLSchemaMetaTask,
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
		Type:         TypeDmTopSQL,
		Desc:         "DM TOP SQL",
		InstanceType: InstanceTypeDm,
		CreateTask:   NewDmTopSQLTask,
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
				Desc:  "排序字段",
				Value: DmTopSQLMetricTotalExecTime,
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeObForOracleTopSQL,
		Desc:         "OceanBase For Oracle TOP SQL",
		InstanceType: InstanceTypeObForOracle,
		CreateTask:   NewObForOracleTopSQLTask,
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
				Desc:  "排序字段",
				Value: DynPerformanceViewObForOracleColumnElapsedTime,
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypePostgreSQLTopSQL,
		Desc:         "TOP SQL",
		InstanceType: InstanceTypePostgreSQL,
		CreateTask:   NewPostgreSQLTopSQLTask,
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
				Desc:  "排序字段",
				Value: DynPerformanceViewPgSQLColumnElapsedTime,
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeTBaseSlowLog,
		Desc:         "慢日志",
		InstanceType: InstanceTypeTBase,
		Params: []*params.Param{
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
		},
		CreateTask: NewTBasePgLog,
	},
}

func init() {
	for _, meta := range EEMetas {
		Metas = append(Metas, meta)
		MetaMap[meta.Type] = meta
	}
}
