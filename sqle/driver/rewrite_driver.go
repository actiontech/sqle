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

type SQLRewriteType string

const (
	SQLRewriteTypeQuery = "query" // rewrite query sql
)

const (
	SQLRewriteTypeKey = "sql-rewrite-type"
)

// rewrite query sql will perform limit processing on the basis of the original sql query result
func GenerateQueryRewriteParams(limit, offset uint32) params.Params {
	return params.Params{
		&params.Param{
			Key:   SQLRewriteTypeKey,
			Value: SQLRewriteTypeQuery,
		},
	}
}
