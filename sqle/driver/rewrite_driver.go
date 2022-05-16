package driver

import (
	"github.com/sirupsen/logrus"
)

// SQLQueryDriver is a SQL rewrite and execute driver
type SQLQueryDriver interface {
	IsQuerySQL(sql string) (bool, error)
	Repage(sql string, limit, offset uint32) (newSql string, err error)
	// The data location in values should be consistent with that in column
	Query(sql string, timeOutSecond uint32) (column []string, values [][]string, err error)
}

// NewSQLQueryDriver return a new instantiated SQLQueryDriver.
func NewSQLQueryDriver(log *logrus.Entry, dbType string, cfg *DSN) (SQLQueryDriver, error) {
	return nil, nil
}
