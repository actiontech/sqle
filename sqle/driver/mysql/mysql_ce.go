//go:build !enterprise
// +build !enterprise

package mysql

import (
	"context"
	"fmt"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

func (i *MysqlDriverImpl) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	return nil, fmt.Errorf("only support Query in enterprise edition")
}
