package model

const (
	// operation code list

	// System(系统)/Users(用户)/Roles(角色) reserved 0-19999

	// Workflow：工单 20000-29999
	OP_WORKFLOW_VIEW        = 20100
	OP_WORKFLOW_VIEW_OTHERS = 20150
	OP_WORKFLOW_CREATE      = 20200
	OP_WORKFLOW_UPDATE      = 20300
	OP_WORKFLOW_DELETE      = 20400

	// AuditPlan: 审核计划 reserved 30000-39999
)

func GetOperationCodeList() []int {
	return []int{
		// Workflow：工单
		OP_WORKFLOW_VIEW,
		OP_WORKFLOW_VIEW_OTHERS,
		OP_WORKFLOW_CREATE,
		OP_WORKFLOW_UPDATE,
		OP_WORKFLOW_DELETE,
	}
}

func GetOperationCodeDesc(opCode int) string {
	switch opCode {
	case OP_WORKFLOW_VIEW:
		return "查看工单"
	case OP_WORKFLOW_VIEW_OTHERS:
		return "查看他人创建的工单"
	case OP_WORKFLOW_CREATE:
		return "创建工单"
	case OP_WORKFLOW_UPDATE:
		return "更新工单"
	case OP_WORKFLOW_DELETE:
		return "删除工单"
	}
	return "未知动作"
}
