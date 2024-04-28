package server

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

const (
	executeSqlFileMode string = "execute_sql_file_mode"
)

type StatusChecker interface {
	CheckIsSupport() bool
}

func NewModuleStatusChecker(driverType string, moduleName string) (StatusChecker, error) {
	switch moduleName {
	case executeSqlFileMode:
		return executeSqlFileChecker{driverType: driverType, moduleName: moduleName}, nil
	}
	return nil, fmt.Errorf("no checker mached")
}

type executeSqlFileChecker struct {
	driverType string
	moduleName string
}

func (checker executeSqlFileChecker) CheckIsSupport() bool {
	return driver.GetPluginManager().IsOptionalModuleEnabled(checker.driverType, driverV2.OptionalExecBatch)
}
