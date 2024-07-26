//go:build enterprise
// +build enterprise

package auditplan

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
	TypeTBaseSlowLog             = "TBase_slow_log"
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
	paramKeyIndicator           = "indicator"
	paramKeyTopN                = "top_n"
	paramKeySlowLogCollectInput = "slow_log_collect_input"
)

var EEMetaBuilderList = []MetaBuilder{
	{
		Type:          TypeMySQLSlowLog,
		Desc:          "慢日志",
		TaskHandlerFn: NewSlowLogTaskV2Fn(),
	},
	{
		Type:          TypeTDSQLSlowLog,
		Desc:          "慢日志",
		TaskHandlerFn: NewTDMySQLSlowLogTaskV2Fn(),
	},
	{
		Type:          TypeOceanBaseForMySQLTopSQL,
		Desc:          "Top SQL",
		TaskHandlerFn: NewObForMysqlTopSQLTaskV2Fn(),
	},
	{
		Type:          TypeDB2TopSQL,
		Desc:          "DB2 Top SQL",
		TaskHandlerFn: NewDB2TopSQLTaskV2Fn(),
	},
	{
		Type:          TypeDB2SchemaMeta,
		Desc:          "DB2库表元数据",
		TaskHandlerFn: NewDB2SchemaMetaTaskV2Fn(),
	},
	{
		Type:          TypeTDSQLSchemaMeta,
		Desc:          "TDSQL库表元数据",
		TaskHandlerFn: NewTDMySQLSchemaMetaTaskV2Fn(),
	},
	{
		Type:          TypePostgreSQLSchemaMeta,
		Desc:          "PostgreSQL库表元数据",
		TaskHandlerFn: NewPGSchemaMetaTaskV2Fn(),
	},
	{
		Type:          TypeDmTopSQL,
		Desc:          "DM TOP SQL",
		TaskHandlerFn: NewDmTopSQLTaskV2Fn(),
	},
	{
		Type:          TypeObForOracleTopSQL,
		Desc:          "OceanBase For Oracle TOP SQL",
		TaskHandlerFn: NewObForOracleTopSQLTaskV2Fn(),
	},
	{
		Type: TypePostgreSQLTopSQL,
		Desc: "TOP SQL",

		TaskHandlerFn: NewPGTopSQLTaskV2Fn()},
	{
		Type:          TypeTBaseSlowLog,
		Desc:          "慢日志",
		TaskHandlerFn: NewTBaseSlowLogTaskV2Fn(),
	},
}

func init() {
	for _, b := range EEMetaBuilderList {
		meta := buildMeta(b)
		Metas = append(Metas, meta)
		MetaMap[b.Type] = meta
	}
}
