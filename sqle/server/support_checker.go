package server

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/config"
)

const (
	sqlOptimization string = "sql_optimization"
)

type StatusChecker interface {
	CheckIsSupport() bool
}

func NewModuleStatusChecker(driverType string, moduleName string) (StatusChecker, error) {
	switch moduleName {
	case sqlOptimization:
		return sqlOptimizationChecker{}, nil
	}
	return nil, fmt.Errorf("no checker matched")
}

type sqlOptimizationChecker struct{}

func (s sqlOptimizationChecker) CheckIsSupport() bool {
	return config.GetOptions().SqleOptions.OptimizationConfig.OptimizationKey != "" &&
		config.GetOptions().SqleOptions.OptimizationConfig.OptimizationURL != ""
}
