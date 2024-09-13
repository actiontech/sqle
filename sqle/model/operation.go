package model

import (
	"context"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
)

// NOTE: related model:
// - model.Role
type RoleOperation struct {
	Model
	RoleID uint `json:"role_id" gorm:"index"`
	Code   uint `json:"op_code" gorm:"column:op_code; comment:'动作权限'"`
}

const (
	// operation code list

	// System(系统)/Users(用户)/Roles(角色) reserved 0-19999

	// Workflow：工单 20000-29999
	// NOTE: 用户默认可以查看自己创建的工单，无需定义此项动作权限
	OP_WORKFLOW_VIEW_OTHERS = 20100
	OP_WORKFLOW_SAVE        = 20200 // including "CREATE" and "UPDATE"
	OP_WORKFLOW_AUDIT       = 20300 // including "PASSED" and "REJECT"
	OP_WORKFLOW_EXECUTE     = 20400 // 上线工单权限

	// AuditPlan: 审核计划 reserved 30000-39999
	// NOTE: 用户默认可以查看自己创建的扫描任务，无需定义此项动作权限
	OP_AUDIT_PLAN_VIEW_OTHERS = 30100
	OP_AUDIT_PLAN_SAVE        = 30200 // including "CREATE" and "UPDATE"

	// SqlQuery: SQL查询 reserved 40000-49999
	OP_SQL_QUERY_QUERY = 40100
)

func GetConfigurableOperationCodeList() []uint {
	opts := append(getConfigurableOperationCodeList(), getConfigurableOperationCodeListForEE()...)
	return opts
}

func getConfigurableOperationCodeList() []uint {
	return []uint{
		// Workflow：工单
		OP_WORKFLOW_VIEW_OTHERS,
		OP_WORKFLOW_SAVE,
		OP_WORKFLOW_AUDIT,
		OP_WORKFLOW_EXECUTE,
		// Audit plan: 扫描任务
		OP_AUDIT_PLAN_VIEW_OTHERS,
		OP_AUDIT_PLAN_SAVE,
		// Sql Query: SQL查询
		OP_SQL_QUERY_QUERY,
	}
}

func GetOperationCodeDesc(ctx context.Context, opCode uint) string {
	switch opCode {
	case OP_WORKFLOW_VIEW_OTHERS:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.OpWorkflowViewOthers)
	case OP_WORKFLOW_SAVE:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.OpWorkflowSave)
	case OP_WORKFLOW_AUDIT:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.OpWorkflowAudit)
	case OP_WORKFLOW_EXECUTE:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.OpWorkflowExecute)
	case OP_AUDIT_PLAN_VIEW_OTHERS:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.OpAuditPlanViewOthers)
	case OP_AUDIT_PLAN_SAVE:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.OpAuditPlanSave)
	case OP_SQL_QUERY_QUERY:
		return locale.Bundle.LocalizeMsgByCtx(ctx, locale.OpSqlQueryQuery)
	default:
		return additionalOperationForEE(ctx, opCode)
	}
}

func CheckIfOperationCodeValid(opCodes []uint) (err error) {

	invalidOpCodes := make([]uint, 0)

	for i := range opCodes {
		if !IsValidOperationCode(opCodes[i]) {
			invalidOpCodes = append(invalidOpCodes, opCodes[i])
		}
	}

	if len(invalidOpCodes) > 0 {
		return errors.NewDataInvalidErr("unknown operation code <%v>", invalidOpCodes)
	}

	return nil
}

func IsValidOperationCode(opCode uint) bool {
	validOpCodes := GetConfigurableOperationCodeList()
	for i := range validOpCodes {
		if opCode == validOpCodes[i] {
			return true
		}
	}
	return false
}

func (s *Storage) ReplaceRoleOperationsByOpCodes(roleID uint, opCodes []uint) (err error) {

	// Delete all current role operation records
	{
		err = s.db.Where("role_id = ?", roleID).
			Unscoped(). // Hard delete
			Delete(RoleOperation{}).Error
		if err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}
	}

	// Insert new role operation record
	if len(opCodes) > 0 {
		for i := range opCodes {
			roleOp := &RoleOperation{
				RoleID: roleID,
				Code:   opCodes[i],
			}
			if err = s.db.Create(roleOp).Error; err != nil {
				return errors.ConnectStorageErrWrapper(err)
			}
		}
	}

	return nil
}

func (s *Storage) GetRoleOperationsByRoleID(roleID uint) (roleOps []*RoleOperation, err error) {
	err = s.db.Where("role_id = ?", roleID).Find(&roleOps).Error
	return roleOps, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) DeleteRoleOperationByRoleID(roleID uint) (err error) {
	return errors.ConnectStorageErrWrapper(
		s.db.Where("role_id = ?", roleID).Delete(RoleOperation{}).Error)
}
