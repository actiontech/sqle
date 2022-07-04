//go:build !enterprise
// +build !enterprise

package mysql

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
)

func (d *DriverManager) getSQLQueryDriver() (driver.SQLQueryDriver, error) {
	return nil, fmt.Errorf("only support SQL query in enterprise edition")
}

func (d *DriverManager) getAnalysisDriver() (driver.AnalysisDriver, error) {
	return nil, fmt.Errorf("only support SQL analysis in enterprise edition")
}