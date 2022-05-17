package driver

import (
	"context"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

// SQLQueryDriver is a SQL rewrite and execute driver
type SQLQueryDriver interface {
	QueryPrepare(ctx context.Context, sql string, conf *QueryPrepareConf) (*QueryPrepareResult, error)
	Query(ctx context.Context, sql string, conf *QueryConf) (*QueryResult, error)
}

// NewSQLQueryDriver return a new instantiated SQLQueryDriver.
func NewSQLQueryDriver(log *logrus.Entry, dbType string, cfg *DSN) (SQLQueryDriver, error) {
	return nil, nil
}

type ErrorType string

const (
	ErrorTypeNotQuery = "not query"
)

type QueryPrepareConf struct {
	Limit  uint32
	Offset uint32
}

type QueryPrepareResult struct {
	NewSQL    string
	ErrorType ErrorType
	Error     string
}

type QueryConf struct {
	TimeOutSecond uint32
}

// The data location in Values should be consistent with that in Column
type QueryResult struct {
	Column params.Params
	Rows   []*QueryResultValue
}

type QueryResultRow struct {
	Values []*QueryResultValue
}

type QueryResultValue struct {
	Value string
}
