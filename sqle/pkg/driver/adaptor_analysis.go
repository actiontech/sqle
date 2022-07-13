package driver

import (
	"context"
	"database/sql"
	"os"

	"github.com/actiontech/sqle/sqle/driver"

	hclog "github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

type ListTablesInSchemaFunc func(ctx context.Context, conf *driver.ListTablesInSchemaConf, dbConf DbConf) (*driver.ListTablesInSchemaResult, error)
type GetTableMetaByTableNameFunc func(ctx context.Context, conf *driver.GetTableMetaByTableNameConf, dbConf DbConf) (*driver.GetTableMetaByTableNameResult, error)
type GetTableMetaBySQLFunc func(ctx context.Context, conf *driver.GetTableMetaBySQLConf, dbConf DbConf) (*driver.GetTableMetaBySQLResult, error)
type ExplainFunc func(ctx context.Context, conf *driver.ExplainConf, dbConf DbConf) (*driver.ExplainResult, error)

type DbConf struct {
	Db   *sql.DB
	Conn *sql.Conn
}

type AnalysisAdaptor struct {
	l hclog.Logger

	dt  Dialector
	dsn *driver.DSN

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
	newDriver := func(dsn *driver.DSN) driver.AnalysisDriver {
		if p, exist := pluginImpls[driver.PluginNameAnalysisDriver]; exist {
			return p
		}

		a.dsn = dsn
		di := &pluginImpl{
			analysisAdaptor: a,
		}
		if a.dsn == nil {
			pluginImpls[driver.PluginNameAnalysisDriver] = di
			return di
		}
		driverName, dsnDetail := a.dt.Dialect(dsn)
		db, conn := getDbConn(driverName, dsnDetail)
		di.db = db
		di.conn = conn
		pluginImpls[driver.PluginNameAnalysisDriver] = di
		return di
	}
	return driver.NewAnalysisDriverPlugin(newDriver)
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
