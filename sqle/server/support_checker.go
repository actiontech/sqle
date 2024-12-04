package server

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	optimization "github.com/actiontech/sqle/sqle/server/optimization/rule"
)

const (
	executeSqlFileMode string = "execute_sql_file_mode"
	sqlOptimization    string = "sql_optimization"
	backup             string = "backup"
)

type StatusChecker interface {
	CheckIsSupport() bool
}

func NewModuleStatusChecker(driverType string, moduleName string) (StatusChecker, error) {
	switch moduleName {
	case executeSqlFileMode:
		return executeSqlFileChecker{driverType: driverType, moduleName: moduleName}, nil
	case sqlOptimization:
		return sqlOptimizationChecker{}, nil
	case backup:
		return sqlBackupChecker{driverType: driverType}, nil
	}
	return nil, fmt.Errorf("no checker matched")
}

type executeSqlFileChecker struct {
	driverType string
	moduleName string
}

func (checker executeSqlFileChecker) CheckIsSupport() bool {
	return driver.GetPluginManager().IsOptionalModuleEnabled(checker.driverType, driverV2.OptionalExecBatch)
}

type sqlOptimizationChecker struct{}

func (s sqlOptimizationChecker) CheckIsSupport() bool {
	return len(optimization.OptimizationRuleMap) > 0
}

type sqlBackupChecker struct {
	driverType string
}

func (s sqlBackupChecker) CheckIsSupport() bool {
	svc := BackupService{}
	return svc.CheckIsDbTypeSupportEnableBackup(s.driverType) == nil
}
