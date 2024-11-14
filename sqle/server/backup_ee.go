//go:build enterprise
// +build enterprise

package server

import (
	"errors"
	"fmt"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

var ErrUnsupportedBackupInFileMode error = errors.New("enable backup in file mode is unsupported")

type BackupService struct{}

// 文件模式不支持备份，仅支持SQL模式上线
func (BackupService) CheckBackupConflictWithExecMode(EnableBackup bool, ExecMode string) error {
	if EnableBackup && ExecMode == executeSqlFileMode {
		return ErrUnsupportedBackupInFileMode
	}
	return nil
}

// 检查数据源类型是否支持备份
func (BackupService) CheckIsDbTypeSupportEnableBackup(dbType string) error {
	if dbType != driverV2.DriverTypeMySQL {
		return fmt.Errorf("db type %v can not enable backup", dbType)
	}
	return nil
}
