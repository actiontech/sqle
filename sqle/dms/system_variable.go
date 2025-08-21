package dms

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsobject "github.com/actiontech/dms/pkg/dms-common/dmsobject"
)

const (
	SystemVariableSqlManageRawExpiredHours    = "system_variable_sql_manage_raw_expired_hours"
	SystemVariableWorkflowExpiredHours        = "system_variable_workflow_expired_hours"
	SystemVariableSqleUrl                     = "system_variable_sqle_url"
	SystemVariableOperationRecordExpiredHours = "system_variable_operation_record_expired_hours"
	SystemVariableCbOperationLogsExpiredHours = "system_variable_cb_operation_logs_expired_hours"
	SystemVariableSSHPrimaryKey               = "system_variable_ssh_primary_key"
)

const (
	DefaultSqlManageRawExpiredHours    = 30 * 24
	DefaultOperationRecordExpiredHours = 90 * 24
	DefaultCbOperationLogsExpiredHours = 90 * 24
)

func GetWorkflowExpiredHoursOrDefault() (int64, error) {
	systemVariable, err := dmsobject.GetSystemVariables(context.TODO(), GetDMSServerAddress())
	if err != nil {
		return 0, err
	}
	if systemVariable.Code == 0 {
		return int64(systemVariable.Data.SystemVariableWorkflowExpiredHours), nil
	}

	return 30 * 24, nil
}
func GetSqlManageRawSqlExpiredHoursOrDefault() (int64, error) {
	systemVariable, err := dmsobject.GetSystemVariables(context.TODO(), GetDMSServerAddress())
	if err != nil {
		return 0, err
	}
	if systemVariable.Code == 0 {
		return int64(systemVariable.Data.SystemVariableWorkflowExpiredHours), nil
	}

	return DefaultSqlManageRawExpiredHours, nil
}

// GetSqleUrl 从DMS获取SQLE URL
func GetSqleUrl(ctx context.Context) (string, error) {
	reply, err := dmsobject.GetSystemVariables(ctx, GetDMSServerAddress())
	if err != nil {
		return "", fmt.Errorf("failed to get system variables from DMS: %v", err)
	}
	return fmt.Sprintf("%s/sqle", reply.Data.Url), nil
}

// UpdateSystemVariables 更新系统变量配置到DMS
func UpdateSystemVariables(ctx context.Context, key, value string) error {
	dmsAddr := GetDMSServerAddress()
	req := &v1.UpdateSystemVariablesReqV1{}
	switch key {
	case SystemVariableSqleUrl:
		req.Url = &value
	case SystemVariableOperationRecordExpiredHours:
		if val, err := strconv.Atoi(value); err == nil {
			req.OperationRecordExpiredHours = &val
		}
	case SystemVariableCbOperationLogsExpiredHours:
		if val, err := strconv.Atoi(value); err == nil {
			req.CbOperationLogsExpiredHours = &val
		}
	case SystemVariableSSHPrimaryKey:
		req.SystemVariableSSHPrimaryKey = &value
	case SystemVariableSqlManageRawExpiredHours:
		if val, err := strconv.Atoi(value); err == nil {
			req.SystemVariableSqlManageRawExpiredHours = &val
		}
	case SystemVariableWorkflowExpiredHours:
		if val, err := strconv.Atoi(value); err == nil {
			req.SystemVariableWorkflowExpiredHours = &val
		}
	}

	return dmsobject.UpdateSystemVariables(ctx, dmsAddr, req)
}
