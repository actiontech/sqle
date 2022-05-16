package driver

import (
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

// SQLQueryDriver is a SQL rewrite and execute driver
type SQLQueryDriver interface {
	Rewrite(sql string, params params.Params) (newSql string, err error)
	Query(sql string) ([]map[string]string, error)
}

// NewSQLQueryDriver return a new instantiated SQLQueryDriver.
func NewSQLQueryDriver(log *logrus.Entry, dbType string, cfg *DSN) (SQLQueryDriver, error) {
	return nil, nil
}

type SQLQueryType string

const (
	SQLQueryTypeQuery = "query" // rewrite query sql
)

const (
	SQLQueryTypeKey = "sql-rewrite-type"
)

// rewrite query sql will perform limit processing on the basis of the original sql query result
func GenerateQueryQueryParams(limit, offset uint32) params.Params {
	return params.Params{
		&params.Param{
			Key:   SQLQueryTypeKey,
			Value: SQLQueryTypeQuery,
		},
	}
}
