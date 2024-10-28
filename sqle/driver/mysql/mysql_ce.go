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

func (i *MysqlDriverImpl) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	return nil, fmt.Errorf("only support Query in enterprise edition")
}

func (i *MysqlDriverImpl) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabasSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {
	return nil, fmt.Errorf("only support Query in enterprise edition")
}
