package model

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

func GetConfigurableOperationCodeList() []int {
	return []int{
		// Workflow：工单
		OP_WORKFLOW_VIEW_OTHERS,
		OP_WORKFLOW_SAVE,
		OP_WORKFLOW_DELETE,
	}
}

func GetOperationCodeDesc(opCode int) string {
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
