package driver

import (
	"context"
	"database/sql/driver"

	v2 "github.com/actiontech/sqle/sqle/driver/v2"
)

type Plugin interface {
	Close(ctx context.Context)

	// Parse parse sqlText to Node array. sqlText may be single SQL or batch SQLs.
	Parse(ctx context.Context, sqlText string) ([]v2.Node, error)

	// Audit sql with rules. sql is single SQL text or multi audit.
	Audit(ctx context.Context, sqls []string) ([]*v2.AuditResults, error)

	// GenRollbackSQL generate sql's rollback SQL.
	GenRollbackSQL(ctx context.Context, sql string) (string, string, error)

	Ping(ctx context.Context) error
	Exec(ctx context.Context, query string) (driver.Result, error)
	Tx(ctx context.Context, queries ...string) ([]driver.Result, error)
	Query(ctx context.Context, sql string, conf *v2.QueryConf) (*v2.QueryResult, error)
	Explain(ctx context.Context, conf *v2.ExplainConf) (*v2.ExplainResult, error)

	// Schemas export all supported schemas.
	//
	// For example, performance_schema/performance_schema... which in MySQL is not allowed for auditing.
	Schemas(ctx context.Context) ([]string, error)

	// in v2, this is a virtual api, it is a combination of [ExtractTableFromSQL, GetTableMeta]
	GetTableMetaBySQL(ctx context.Context, conf *GetTableMetaBySQLConf) (*GetTableMetaBySQLResult, error)
}

type PluginBoot interface {
	Register() (*v2.DriverMetas, error)
	Open(*v2.Config) (Plugin, error)
	Stop() error
}

type GetTableMetaBySQLConf struct {
	// this SQL should be a single SQL
	Sql string
}

type GetTableMetaBySQLResult struct {
	TableMetas []*TableMeta
}

type TableMeta struct {
	v2.Table
	v2.TableMeta
}
