//go:build enterprise
// +build enterprise

package mysql

import (
	"context"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"vitess.io/vitess/go/vt/sqlparser"
)

func init() {
	driver.RegisterSQLQueryDriver(driver.DriverTypeMySQL, newQueryDriverInspect)
}

func newQueryDriverInspect(log *logrus.Entry, dsn *driver.DSN) (driver.SQLQueryDriver, error) {
	var inspect = &Inspect{}

	if dsn != nil {
		conn, err := executor.NewExecutor(log, dsn, dsn.DatabaseName)
		if err != nil {
			return nil, errors.Wrap(err, "new executor in inspect")
		}
		inspect.isConnected = true
		inspect.dbConn = conn
		inspect.inst = dsn

		ctx := session.NewContext(nil, session.WithExecutor(conn))
		ctx.SetCurrentSchema(dsn.DatabaseName)

		inspect.Ctx = ctx
	} else {
		ctx := session.NewContext(nil)
		inspect.Ctx = ctx
	}

	inspect.log = log
	inspect.result = driver.NewInspectResults()
	inspect.isOfflineAudit = dsn == nil

	inspect.cnf = &Config{
		DMLRollbackMaxRows: -1,
		DDLOSCMinSize:      -1,
		DDLGhostMinSize:    -1,
	}

	return inspect, nil
}

func (*Inspect) QueryPrepare(ctx context.Context, sql string, conf *driver.QueryPrepareConf) (*driver.QueryPrepareResult, error) {
	return QueryPrepare(ctx, sql, conf)
}

func QueryPrepare(ctx context.Context, sql string, conf *driver.QueryPrepareConf) (*driver.QueryPrepareResult, error) {
	node, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	// check is query sql
	stmt, ok := node.(*sqlparser.Select)
	if !ok {
		return &driver.QueryPrepareResult{
			ErrorType: driver.ErrorTypeNotQuery,
			Error:     driver.ErrorTypeNotQuery,
		}, nil
	}

	// Generate new limit
	limit, offset := -1, -1
	if stmt.Limit != nil {
		if stmt.Limit.Rowcount != nil {
			limit, _ = strconv.Atoi(stmt.Limit.Rowcount.(*sqlparser.Literal).Val)
		}
		if stmt.Limit.Rowcount != nil {
			offset, _ = strconv.Atoi(stmt.Limit.Offset.(*sqlparser.Literal).Val)
		} else if limit != -1 {
			offset = 0
		}
	}
	appendLimit, appendOffset := -1, -1
	if conf != nil {
		appendLimit, appendOffset = int(conf.Limit), int(conf.Offset)
	}
	if appendLimit != -1 && appendOffset == -1 {
		appendLimit = 0
	}

	newLimit, newOffset := CalculateOffset(limit, offset, appendLimit, appendOffset)

	if newLimit != -1 {
		l := &sqlparser.Limit{
			Offset: &sqlparser.Literal{
				Type: sqlparser.IntVal,
				Val:  strconv.Itoa(newOffset),
			},
			Rowcount: &sqlparser.Literal{
				Type: sqlparser.IntVal,
				Val:  strconv.Itoa(newLimit),
			},
		}
		stmt.SetLimit(l)
	}

	// rewrite
	return &driver.QueryPrepareResult{
		NewSQL:    sqlparser.String(stmt),
		ErrorType: driver.ErrorTypeNotError,
	}, nil
}

// 1 means this item has no value or no limit
func CalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) (newLimit, newOffset int) {
	if checkIsInvalidCalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset) {
		return oldLimit, oldOffset
	}
	return calculateOffset(oldLimit, oldOffset, appendLimit, appendOffset)
}

func calculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) (newLimit, newOffset int) {
	if oldLimit == -1 {
		return appendLimit, appendOffset
	}
	newOffset = oldOffset + appendOffset
	newLimit = appendLimit
	if newOffset+newLimit > oldLimit+oldOffset {
		newLimit = oldLimit - appendOffset
	}

	return newLimit, newOffset
}

func checkIsInvalidCalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) bool {
	if appendLimit == -1 {
		return true
	}
	if oldLimit != -1 && appendOffset > oldLimit+oldOffset {
		return true
	}

	return false
}

func (i *Inspect) Query(ctx context.Context, sql string, conf *driver.QueryConf) (*driver.QueryResult, error) {
	// check sql
	prepareRes, err := i.QueryPrepare(ctx, sql, &driver.QueryPrepareConf{
		Limit:  1,
		Offset: 1,
	})
	if err != nil {
		return nil, err
	}
	if prepareRes.ErrorType != "" && prepareRes.ErrorType != driver.ErrorTypeNotError {
		return nil, errors.New(prepareRes.Error)
	}

	// add timeout
	cancel := func() {}
	if conf != nil {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(conf.TimeOutSecond)*time.Second)
	}
	defer cancel()

	// execute sql
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	result, err := conn.Db.QueryWithContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	// generate result
	res := &driver.QueryResult{
		Column: params.Params{},
		Rows:   []*driver.QueryResultRow{},
	}
	for i, row := range result {
		r := &driver.QueryResultRow{
			Values: []*driver.QueryResultValue{},
		}
		for key, value := range row {
			if i == 0 {
				res.Column = append(res.Column, &params.Param{
					Key:   key,
					Value: key,
				})
			}
			r.Values = append(r.Values, &driver.QueryResultValue{
				Value: value.String,
			})
		}
	}
	return res, nil
}
