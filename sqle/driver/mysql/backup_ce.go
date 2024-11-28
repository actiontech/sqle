//go:build !enterprise
// +build !enterprise

package mysql

import (
	"context"
	"fmt"
	"github.com/actiontech/sqle/sqle/driver"
)

var ErrUnsupportedBackup error = fmt.Errorf("backup is unsupported for sqle community version")

func (i *MysqlDriverImpl) Backup(ctx context.Context, backupStrategy string, sql string) (BackupSql []string, ExecuteInfo string, err error) {
	return nil, "", ErrUnsupportedBackup
}

func (i *MysqlDriverImpl) GetBackupStrategy(ctx context.Context, sql string) (*driver.GetBackupStrategyRes, error) {
	return nil, ErrUnsupportedBackup
}
