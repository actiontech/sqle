//go:build enterprise
// +build enterprise

package auditplan

import (
	"github.com/actiontech/sqle/sqle/cmd/scannerd/command"
	"github.com/actiontech/sqle/sqle/locale"
)

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
	TypeTBaseSlowLog             = command.TypeTBaseSlowLog
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
	paramKeySchema              = "schema"
)

var EEMetaBuilderList = []MetaBuilder{
	{
		Type:          TypeMySQLSlowLog,
		Desc:          locale.ApMetaSlowLog,
		TaskHandlerFn: NewSlowLogTaskV2Fn(),
	},
	{
		Type:          TypeTDSQLSlowLog,
		Desc:          locale.ApMetaSlowLog,
		TaskHandlerFn: NewTDMySQLSlowLogTaskV2Fn(),
	},
	{
		Type:          TypeOceanBaseForMySQLTopSQL,
		Desc:          locale.ApMetaTopSQL,
		TaskHandlerFn: NewObForMysqlTopSQLTaskV2Fn(),
	},
	{
		Type:          TypeDB2TopSQL,
		Desc:          locale.ApMetaDB2TopSQL,
		TaskHandlerFn: NewDB2TopSQLTaskV2Fn(),
	},
	{
		Type:          TypeDB2SchemaMeta,
		Desc:          locale.ApMetaSchemaMeta,
		TaskHandlerFn: NewDB2SchemaMetaTaskV2Fn(),
	},
	{
		Type:          TypeTDSQLSchemaMeta,
		Desc:          locale.ApMetaSchemaMeta,
		TaskHandlerFn: NewTDMySQLSchemaMetaTaskV2Fn(),
	},
	{
		Type:          TypePostgreSQLSchemaMeta,
		Desc:          locale.ApMetaSchemaMeta,
		TaskHandlerFn: NewPGSchemaMetaTaskV2Fn(),
	},
	{
		Type:          TypeDmTopSQL,
		Desc:          locale.ApMetaDmTopSQL,
		TaskHandlerFn: NewDmTopSQLTaskV2Fn(),
	},
	{
		Type:          TypeObForOracleTopSQL,
		Desc:          locale.ApMetaObForOracleTopSQL,
		TaskHandlerFn: NewObForOracleTopSQLTaskV2Fn(),
	},
	{
		Type:          TypePostgreSQLTopSQL,
		Desc:          locale.ApMetaPostgreSQLTopSQL,
		TaskHandlerFn: NewPGTopSQLTaskV2Fn(),
	},
	{
		Type:          TypeTBaseSlowLog,
		Desc:          locale.ApMetaSlowLog,
		TaskHandlerFn: NewTBaseSlowLogTaskV2Fn(),
	},
}

func init() {
	for _, b := range EEMetaBuilderList {
		meta := buildMeta(b)
		Metas = append(Metas, meta)
		MetaMap[b.Type] = meta
	}

	supportedCmdTypeList[TypeTBaseSlowLog] = struct{}{}
	supportedCmdTypeList[TypeTDSQLSlowLog] = struct{}{}
}
