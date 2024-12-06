//go:build !enterprise
// +build !enterprise

package mysql

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
)

var ErrUnsupportedBackup error = fmt.Errorf("backup is unsupported for sqle community version")

func (i *MysqlDriverImpl) Backup(ctx context.Context, backupStrategy string, sql string, backupMaxRows uint64) (backupSqls []string, executeResult string, err error) {
	return nil, "", ErrUnsupportedBackup
}

func (i *MysqlDriverImpl) RecommendBackupStrategy(ctx context.Context, sql string) (*driver.RecommendBackupStrategyRes, error) {
	return nil, ErrUnsupportedBackup
}
