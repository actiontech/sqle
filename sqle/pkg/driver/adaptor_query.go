package driver

import (
	"context"
	"database/sql"
	"os"

	"github.com/actiontech/sqle/sqle/driver"

	"github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

type SQLQueryPrepareFunc func(ctx context.Context, sql string, conf *driver.QueryPrepareConf) (*driver.QueryPrepareResult, error)
type SQLQueryFunc func(ctx context.Context, sql string, conf *driver.QueryConf) (*driver.QueryResult, error)

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
			q: q,
		}
		if q.dsn == nil {
			pluginImpls[driver.PluginNameQueryDriver] = di
			return di
		}
		driverName, dsnDetail := q.dt.Dialect(dsn)
		db, err := sql.Open(driverName, dsnDetail)
		if err != nil {
			panic(errors.Wrap(err, "open database failed when new driver"))
		}
		conn, err := db.Conn(context.TODO())
		if err != nil {
			panic(errors.Wrap(err, "get database connection failed when new driver"))
		}
		if err := conn.PingContext(context.TODO()); err != nil {
			panic(errors.Wrap(err, "ping database connection failed when new driver"))
		}

		di.db = db
		di.conn = conn
		pluginImpls[driver.PluginNameQueryDriver] = di
		return di
	}
	return driver.NewQueryDriverPlugin(newDriver)
}
