package driver

import (
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

// SQLRewriteDriver is a SQL rewrite driver
type SQLRewriteDriver interface {
	Rewrite(sql string, params params.Params) (newSql string, err error)
}

// NewSQLRewriteDriver return a new instantiated SQLRewriteDriver.
func NewSQLRewriteDriver(log *logrus.Entry, dbType string) (SQLRewriteDriver, error) {
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
