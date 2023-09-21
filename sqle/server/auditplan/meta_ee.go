//go:build enterprise
// +build enterprise

package auditplan

import "github.com/actiontech/sqle/sqle/pkg/params"

const (
	TypeOceanBaseForMySQLMybatis = "ocean_base_for_mysql_mybatis"
	TypeOceanBaseForMySQLTopSQL  = "ocean_base_for_mysql_top_sql"
	TypeDB2TopSQL                = "db2_top_sql"
	TypeDB2SchemaMeta            = "db2_schema_meta"
	TypeTDSQLSlowLog             = "tdsql_for_innodb_slow_log"
	TypeTDSQLSchemaMeta          = "tdsql_for_innodb_schema_meta"
)

const (
	InstanceTypeOceanBaseForMySQL = "OceanBase For MySQL"
	InstanceTypeDB2               = "DB2"
	InstanceTypeTDSQL             = "TDSQL For InnoDB"
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
}

func init() {
	for _, meta := range EEMetas {
		Metas = append(Metas, meta)
		MetaMap[meta.Type] = meta
	}
}
