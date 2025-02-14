package server

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/license"
	optimization "github.com/actiontech/sqle/sqle/server/optimization/rule"
)

const (
	executeSqlFileMode string = "execute_sql_file_mode"
	sqlOptimization    string = "sql_optimization"
	backup             string = "backup"
	knowledge_base     string = "knowledge_base"
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
	case knowledge_base:
		return knowledgeBaseChecker{driverType: driverType}, nil
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

// knowledgeBaseChecker 知识库检查器
type knowledgeBaseChecker struct {
	driverType string
}

func (s knowledgeBaseChecker) CheckIsSupport() bool {
	if license.CheckKnowledgeBaseLicense(config.GetOptions().SqleOptions.KnowledgeBaseTempLicense) != nil {
		return false
	}
	return true
}
