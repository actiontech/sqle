//go:build !enterprise
// +build !enterprise

package mysql

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

func (i *MysqlDriverImpl) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	return nil, fmt.Errorf("only support Explain in enterprise edition")
}

func (i *MysqlDriverImpl) GetTableMetaBySQL(ctx context.Context, conf *driver.GetTableMetaBySQLConf) (*driver.GetTableMetaBySQLResult, error) {
	return nil, fmt.Errorf("only support GetTableMetaBySQL in enterprise edition")
}

func (i *MysqlDriverImpl) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	return nil, fmt.Errorf("only support Query in enterprise edition")
}
