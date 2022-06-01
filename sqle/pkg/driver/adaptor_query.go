package driver

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/pkg/params"
	goPlugin "github.com/hashicorp/go-plugin"

	"github.com/hashicorp/go-hclog"
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
		q.dsn = dsn
		di := &queryDriverImpl{
			q: q,
		}
		if q.dsn == nil {
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
		return di
	}
	return driver.NewQueryDriverPlugin(newDriver)
}

type queryDriverImpl struct {
	q    *QueryAdaptor
	db   *sql.DB
	conn *sql.Conn
}

func (q *queryDriverImpl) QueryPrepare(ctx context.Context, sql string, conf *driver.QueryPrepareConf) (*driver.QueryPrepareResult, error) {
	if q.q.queryPrepare != nil {
		return q.q.queryPrepare(ctx, sql, conf)
	}
	return &driver.QueryPrepareResult{
		NewSQL:    sql,
		ErrorType: driver.ErrorTypeNotError,
		Error:     "",
	}, nil
}

func (q *queryDriverImpl) Query(ctx context.Context, query string, conf *driver.QueryConf) (*driver.QueryResult, error) {
	if q.q.query != nil {
		return q.q.query(ctx, query, conf)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(conf.TimeOutSecond)*time.Second)
	defer cancel()
	rows, err := q.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := &driver.QueryResult{
		Column: params.Params{},
		Rows:   []*driver.QueryResultRow{},
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for _, column := range columns {
		result.Column = append(result.Column, &params.Param{
			Key:   column,
			Value: column,
			Desc:  column,
		})
	}

	for rows.Next() {
		buf := make([]interface{}, len(columns))
		data := make([]sql.NullString, len(columns))
		for i := range buf {
			buf[i] = &data[i]
		}
		if err := rows.Scan(buf...); err != nil {
			return nil, err
		}
		value := &driver.QueryResultRow{
			Values: []*driver.QueryResultValue{},
		}
		for i := 0; i < len(columns); i++ {
			value.Values = append(value.Values, &driver.QueryResultValue{Value: data[i].String})
		}
		result.Rows = append(result.Rows, value)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
