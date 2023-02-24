package driver

import (
	"context"
	"os"

	driverV1 "github.com/actiontech/sqle/sqle/driver/v1"

	hclog "github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"
)

type SQLQueryPrepareFunc func(ctx context.Context, sql string, conf *driverV1.QueryPrepareConf, dbConf DbConf) (*driverV1.QueryPrepareResult, error)
type SQLQueryFunc func(ctx context.Context, sql string, conf *driverV1.QueryConf, dbConf DbConf) (*driverV1.QueryResult, error)

type QueryAdaptor struct {
	l hclog.Logger

	dt  Dialector
	dsn *driverV1.DSN

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
	newDriver := func(dsn *driverV1.DSN) driverV1.SQLQueryDriver {
		if p, exist := pluginImpls[driverV1.PluginNameQueryDriver]; exist {
			return p
		}

		q.dsn = dsn
		di := &pluginImpl{
			queryAdaptor: q,
		}
		if q.dsn == nil {
			pluginImpls[driverV1.PluginNameQueryDriver] = di
			return di
		}
		driverName, dsnDetail := q.dt.Dialect(dsn)
		db, conn := getDbConn(driverName, dsnDetail)
		di.db = db
		di.conn = conn
		pluginImpls[driverV1.PluginNameQueryDriver] = di
		return di
	}
	return driverV1.NewQueryDriverPlugin(newDriver)
}
