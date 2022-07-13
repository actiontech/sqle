package driver

import (
	"context"
	"os"

	"github.com/actiontech/sqle/sqle/driver"

	hclog "github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"
)

type SQLQueryPrepareFunc func(ctx context.Context, sql string, conf *driver.QueryPrepareConf, dbConf DbConf) (*driver.QueryPrepareResult, error)
type SQLQueryFunc func(ctx context.Context, sql string, conf *driver.QueryConf, dbConf DbConf) (*driver.QueryResult, error)

type QueryAdaptor struct {
	l hclog.Logger

	dt  Dialector
	dsn *driver.DSN

	queryPrepare SQLQueryPrepareFunc
	query        SQLQueryFunc
}

func NewQueryAdaptor(dt Dialector) *QueryAdaptor {
	return &QueryAdaptor{
		dt: dt,
		l: hclog.New(&hclog.LoggerOptions{
			JSONFormat: true,
			Output:     os.Stderr,
			Level:      hclog.Trace,
		}),
	}
}

func (q *QueryAdaptor) AddSQLQueryPrepareFunc(f SQLQueryPrepareFunc) {
	q.queryPrepare = f
}

func (q *QueryAdaptor) AddSQLQueryFunc(f SQLQueryFunc) {
	q.query = f
}

func (q *QueryAdaptor) GeneratePlugin() goPlugin.Plugin {
	defer func() {
		if err := recover(); err != nil {
			q.l.Error("panic", "err", err)
		}
	}()
	newDriver := func(dsn *driver.DSN) driver.SQLQueryDriver {
		if p, exist := pluginImpls[driver.PluginNameQueryDriver]; exist {
			return p
		}

		q.dsn = dsn
		di := &pluginImpl{
			queryAdaptor: q,
		}
		if q.dsn == nil {
			pluginImpls[driver.PluginNameQueryDriver] = di
			return di
		}
		driverName, dsnDetail := q.dt.Dialect(dsn)
		db, conn := getDbConn(driverName, dsnDetail)
		di.db = db
		di.conn = conn
		pluginImpls[driver.PluginNameQueryDriver] = di
		return di
	}
	return driver.NewQueryDriverPlugin(newDriver)
}
