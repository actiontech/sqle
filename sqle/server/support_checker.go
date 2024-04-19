package server

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

func NewFunctionSupportChecker(driverType string, functionType string) (SupportChecker, error) {
	instType, err := convertToDriverType(driverType)
	if err != nil {
		return nil, err
	}
	funcType, err := convertToFunctionType(functionType)
	if err != nil {
		return nil, err
	}
	switch funcType {
	case execute_sql_file_mode:
		return executeSqlFileChecker{driverType: instType, functionType: funcType}, nil
	}
	return nil, fmt.Errorf("no checker mached")
}

type SupportChecker interface {
	CheckIsSupport() bool
}

type executeSqlFileChecker struct {
	driverType   DriverType
	functionType FunctionType
}

var supportExecBatch map[DriverType]bool = make(map[DriverType]bool)

func (checker executeSqlFileChecker) CheckIsSupport() bool {
	if isSupport, exist := supportExecBatch[checker.driverType]; exist {
		return isSupport
	} else {
		// lazy load, when driver did not load in this map
		for _, meta := range driver.GetPluginManager().AllDriverMetas() {
			driverType, _ := convertToDriverType(meta.PluginName)
			for _, module := range meta.EnabledOptionalModule {
				if module.String() == "ExecBatch" {
					supportExecBatch[driverType] = true
				}
			}
		}
		return supportExecBatch[checker.driverType]
	}
}

type DriverType string

const (
	DriverTypeMySQL          DriverType = driverV2.DriverTypeMySQL
	DriverTypeOracle         DriverType = driverV2.DriverTypeOracle
	DriverTypeSqlServer      DriverType = driverV2.DriverTypeSQLServer
	DriverTypePostgreSQL     DriverType = driverV2.DriverTypePostgreSQL
	DriverTypeTiDB           DriverType = driverV2.DriverTypeTiDB
	DriverTypeDB2            DriverType = driverV2.DriverTypeDB2
	DriverTypeOceanBase      DriverType = driverV2.DriverTypeOceanBase
	DriverTypeTDSQLForInnoDB DriverType = driverV2.DriverTypeTDSQLForInnoDB
)

type FunctionType string

const (
	execute_sql_file_mode FunctionType = "execute_sql_file_mode"
)

func errUnknownDriverType() error {
	return fmt.Errorf("unknown driver type")
}

func errUnknownFunctionType() error {
	return fmt.Errorf("unknown function type")
}

func convertToDriverType(driverType string) (DriverType, error) {
	switch driverType {
	case "MySQL":
		return DriverTypeMySQL, nil
	case "Oracle":
		return DriverTypeOracle, nil
	case "SQL Server":
		return DriverTypeSqlServer, nil
	case "PostgreSQL":
		return DriverTypePostgreSQL, nil
	case "TiDB":
		return DriverTypeTiDB, nil
	case "DB2":
		return DriverTypeDB2, nil
	case "OceanBase For MySQL":
		return DriverTypeOceanBase, nil
	case "TDSQL For InnoDB":
		return DriverTypeTDSQLForInnoDB, nil
	default:
		return "", errUnknownDriverType()
	}
}

func convertToFunctionType(tp string) (FunctionType, error) {
	switch tp {
	case "execute_sql_file_mode":
		return execute_sql_file_mode, nil
	default:
		return "", errUnknownFunctionType()
	}
}
