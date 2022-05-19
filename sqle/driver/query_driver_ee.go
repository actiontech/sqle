//go:build enterprise
// +build enterprise

package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/sirupsen/logrus"
)

// SQLQueryDriver is a SQL rewrite and execute driver
type SQLQueryDriver interface {
	QueryPrepare(ctx context.Context, sql string, conf *QueryPrepareConf) (*QueryPrepareResult, error)
	Query(ctx context.Context, sql string, conf *QueryConf) (*QueryResult, error)
}

var queryDriverMu = &sync.RWMutex{}
var queryDrivers = make(map[string]queryHandler)

// queryHandler is a template which SQLQueryDriver plugin should provide such function signature.
type queryHandler func(log *logrus.Entry, c *DSN) (SQLQueryDriver, error)

// NewSQLQueryDriver return a new instantiated SQLQueryDriver.
func NewSQLQueryDriver(log *logrus.Entry, dbType string, cfg *DSN) (SQLQueryDriver, error) {
	queryDriverMu.RLock()
	defer queryDriverMu.RUnlock()
	d, exist := queryDrivers[dbType]
	if !exist {
		return nil, fmt.Errorf("driver type %v is not supported", dbType)
	}
	return d(log, cfg)
}

// RegisterSQLQueryDriver like sql.Register.
//
// RegisterSQLQueryDriver makes a database driver available by the provided driver name.
// SQLQueryDriver's initialize handler and audit rules register by RegisterSQLQueryDriver.
func RegisterSQLQueryDriver(name string, h queryHandler) {
	queryDriverMu.RLock()
	_, exist := queryDrivers[name]
	queryDriverMu.RUnlock()
	if exist {
		panic("duplicated driver name")
	}

	queryDriverMu.Lock()
	queryDrivers[name] = h
	queryDriverMu.Unlock()
}

type ErrorType string

const (
	ErrorTypeNotQuery = "not query"
	ErrorTypeNotError = "not error"
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
	Rows   []*QueryResultRow
}

type QueryResultRow struct {
	Values []*QueryResultValue
}

type QueryResultValue struct {
	Value string
}
