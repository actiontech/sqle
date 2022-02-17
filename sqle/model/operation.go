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
	OP_WORKFLOW_DELETE      = 20300

	// AuditPlan: 审核计划 reserved 30000-39999
)

func GetConfigurableOperationCodeList() []uint {
	return []uint{
		// Workflow：工单
		OP_WORKFLOW_VIEW_OTHERS,
		OP_WORKFLOW_SAVE,
		OP_WORKFLOW_DELETE,
	}
}

func GetOperationCodeDesc(opCode uint) string {
	switch opCode {
	case OP_WORKFLOW_VIEW_OTHERS:
		return "查看他人创建的工单"
	case OP_WORKFLOW_SAVE:
		return "创建/编辑工单"
	case OP_WORKFLOW_DELETE:
		return "删除工单"
	}
	return "未知动作"
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

func (s *Storage) GetOperationCodesByRoleIDs(roleIDs []uint) (
	opCodes []uint, err error) {

	if len(roleIDs) == 0 {
		return opCodes, nil
	}

	err = s.db.Model(&RoleOperation{}).
		Where("role_id IN (?)", roleIDs).
		Group("op_code, role_id").
		Pluck("op_code", &opCodes).Error
	if err != nil {
		return opCodes, errors.ConnectStorageErrWrapper(err)
	}

	return opCodes, nil
}

func (s *Storage) CheckUserInstanceAccessByOpcodes(userID uint, instIDs, opCodes []uint) (
	missingInstIDs []uint, missingOpCodes []uint, ok bool, err error) {

	roles, err := s.GetRolesByUserID(int(userID))
	if err != nil {
		return missingInstIDs, missingOpCodes, false, err
	}

	if len(roles) == 0 {
		return instIDs, opCodes, false, nil
	}

	roleIDs := GetRoleIDsFromRoles(roles)

	return s.CheckRoleInstanceAccessByOpCodes(roleIDs, instIDs, opCodes)
}
