package driver

import (
	"context"
	"database/sql"
	"os"

	driverV1 "github.com/actiontech/sqle/sqle/driver/v1"

	hclog "github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

type ListTablesInSchemaFunc func(ctx context.Context, conf *driverV1.ListTablesInSchemaConf, dbConf DbConf) (*driverV1.ListTablesInSchemaResult, error)
type GetTableMetaByTableNameFunc func(ctx context.Context, conf *driverV1.GetTableMetaByTableNameConf, dbConf DbConf) (*driverV1.GetTableMetaByTableNameResult, error)
type GetTableMetaBySQLFunc func(ctx context.Context, conf *driverV1.GetTableMetaBySQLConf, dbConf DbConf) (*driverV1.GetTableMetaBySQLResult, error)
type ExplainFunc func(ctx context.Context, conf *driverV1.ExplainConf, dbConf DbConf) (*driverV1.ExplainResult, error)

type DbConf struct {
	Db   *sql.DB
	Conn *sql.Conn
}

type AnalysisAdaptor struct {
	l hclog.Logger

	dt  Dialector
	dsn *driverV1.DSN

	listTablesInSchemaFunc      ListTablesInSchemaFunc
	getTableMetaByTableNameFunc GetTableMetaByTableNameFunc
	getTableMetaBySQLFunc       GetTableMetaBySQLFunc
	explainFunc                 ExplainFunc
}

func NewAnalysisAdaptor(dt Dialector) *AnalysisAdaptor {
	return &AnalysisAdaptor{
		dt: dt,
		l: hclog.New(&hclog.LoggerOptions{
			JSONFormat: true,
			Output:     os.Stderr,
			Level:      hclog.Trace,
		}),
	}
}

func (a *AnalysisAdaptor) AddListTablesInSchemaFunc(f ListTablesInSchemaFunc) {
	a.listTablesInSchemaFunc = f
}

func (a *AnalysisAdaptor) AddGetTableMetaByTableNameFunc(f GetTableMetaByTableNameFunc) {
	a.getTableMetaByTableNameFunc = f
}

func (a *AnalysisAdaptor) AddGetTableMetaBySQLFunc(f GetTableMetaBySQLFunc) {
	a.getTableMetaBySQLFunc = f
}

func (a *AnalysisAdaptor) AddExplainFunc(f ExplainFunc) {
	a.explainFunc = f
}

func (a *AnalysisAdaptor) GeneratePlugin() goPlugin.Plugin {
	defer func() {
		if err := recover(); err != nil {
			a.l.Error("panic", "err", err)
		}
	}()
	newDriver := func(dsn *driverV1.DSN) driverV1.AnalysisDriver {
		if p, exist := pluginImpls[driverV1.PluginNameAnalysisDriver]; exist {
			return p
		}

		a.dsn = dsn
		di := &pluginImpl{
			analysisAdaptor: a,
		}
		if a.dsn == nil {
			pluginImpls[driverV1.PluginNameAnalysisDriver] = di
			return di
		}
		driverName, dsnDetail := a.dt.Dialect(dsn)
		db, conn := getDbConn(driverName, dsnDetail)
		di.db = db
		di.conn = conn
		pluginImpls[driverV1.PluginNameAnalysisDriver] = di
		return di
	}
	return driverV1.NewAnalysisDriverPlugin(newDriver)
}

func getDbConn(driverName, dsnDetail string) (db *sql.DB, conn *sql.Conn) {
	var err error
	db, err = sql.Open(driverName, dsnDetail)
	if err != nil {
		panic(errors.Wrap(err, "open database failed when new driver"))
	}
	conn, err = db.Conn(context.TODO())
	if err != nil {
		panic(errors.Wrap(err, "get database connection failed when new driver"))
	}
	if err := conn.PingContext(context.TODO()); err != nil {
		panic(errors.Wrap(err, "ping database connection failed when new driver"))
	}
	return
}
