package model

import "github.com/actiontech/sqle/sqle/errors"

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

	// AuditPlan: 审核计划 reserved 30000-39999
	// NOTE: 用户默认可以查看自己创建的扫描任务，无需定义此项动作权限
	OP_AUDIT_PLAN_VIEW_OTHERS = 30100
	OP_AUDIT_PLAN_SAVE        = 30200 // including "CREATE" and "UPDATE"

	// SqlQuery: SQL查询 reserved 40000-49999
	// SqlQuery is implement in edition ee
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
		// Audit plan: 扫描任务
		OP_AUDIT_PLAN_VIEW_OTHERS,
		OP_AUDIT_PLAN_SAVE,
	}
}

func GetOperationCodeDesc(opCode uint) string {
	switch opCode {
	case OP_WORKFLOW_VIEW_OTHERS:
		return "查看他人创建的工单"
	case OP_WORKFLOW_SAVE:
		return "创建/编辑工单"
	case OP_WORKFLOW_AUDIT:
		return "审核/驳回工单"
	case OP_AUDIT_PLAN_VIEW_OTHERS:
		return "查看他人创建的扫描任务"
	case OP_AUDIT_PLAN_SAVE:
		return "创建扫描任务"
	default:
		return additionalOperationForEE(opCode)
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
