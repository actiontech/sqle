package dashboard

// ViewType defines the global dashboard visibility scope (user personal vs admin global).
//
// swagger:enum ViewType
type ViewType string

const (
	// ViewTypeAdmin: global view — risks and governance across the estate
	ViewTypeAdmin ViewType = "admin"
	// ViewTypeUser: personal view — my tasks and my assets
	ViewTypeUser ViewType = "user"
)

// AccountFilterCard selects which account management dashboard card drives the list query (query param filter_card).
// Product cards: 即将过期, 我的可用账号.
//
// swagger:enum AccountFilterCard
type AccountFilterCard string

const (
	// AccountFilterCardExpiringSoon: 即将过期
	AccountFilterCardExpiringSoon AccountFilterCard = "expiring_soon"
	// AccountFilterCardMyAvailable: 我的可用账号
	AccountFilterCardMyAvailable AccountFilterCard = "active"
)

// GlobalWorkflowFilterCard selects which dashboard workflow card drives the list query (query param filter_card).
//
// swagger:enum GlobalWorkflowFilterCard
type GlobalWorkflowFilterCard string

const (
	// GlobalWorkflowFilterCardArchived: 已归档 / 历史
	GlobalWorkflowFilterCardArchived GlobalWorkflowFilterCard = "archived"
	// GlobalWorkflowFilterCardPendingForMe: 待我处理
	GlobalWorkflowFilterCardPendingForMe GlobalWorkflowFilterCard = "pending_for_me"
	// GlobalWorkflowFilterCardInitiatedByMe: 我发起的
	GlobalWorkflowFilterCardInitiatedByMe GlobalWorkflowFilterCard = "initiated_by_me"
)

// WorkflowType is the dashboard/API identifier for a workflow business line.
//
// swagger:enum WorkflowType
type WorkflowType string

const (
	// WorkflowTypeSQLRelease is the SQL release (上线) workflow backed by model.WorkflowRecord.
	WorkflowTypeSQLRelease WorkflowType = "sql_release"
	// WorkflowTypeDataExport is the data export workflow (dms data export APIs).
	WorkflowTypeDataExport WorkflowType = "data_export"
)

// UnifiedWorkflowStatus normalizes per-domain workflow states for the global dashboard list and filters.
//
// swagger:enum UnifiedWorkflowStatus
type UnifiedWorkflowStatus string

const (
	UnifiedWorkflowStatusPendingApproval UnifiedWorkflowStatus = "pending_approval" // 待我审核 / 待批准
	UnifiedWorkflowStatusPendingAction   UnifiedWorkflowStatus = "pending_action"   // 待我执行或后续操作（如待上线、待导出）
	UnifiedWorkflowStatusInProgress      UnifiedWorkflowStatus = "in_progress"      // SQL 上线执行中
	UnifiedWorkflowStatusExporting       UnifiedWorkflowStatus = "exporting"        // 数据导出中
	UnifiedWorkflowStatusRejected        UnifiedWorkflowStatus = "rejected"         // 已驳回
	UnifiedWorkflowStatusCancelled       UnifiedWorkflowStatus = "cancelled"        // 已关闭（统一 canceled / cancel）
	UnifiedWorkflowStatusFailed          UnifiedWorkflowStatus = "failed"           // 执行失败
	UnifiedWorkflowStatusCompleted       UnifiedWorkflowStatus = "completed"        // 已成功结束
	// UnifiedWorkflowStatusUnknown is returned when the native status is not mapped yet.
	UnifiedWorkflowStatusUnknown UnifiedWorkflowStatus = "unknown"
)

// SqlManageFilterCard selects which SQL manage dashboard card drives the list query (query param filter_card).
// Product cards: 待我优化, 优化完成.
//
// swagger:enum SqlManageFilterCard
type SqlManageFilterCard string

const (
	// SqlManageFilterCardPendingOptimize: 待我优化
	SqlManageFilterCardPendingOptimize SqlManageFilterCard = "pending"
	// SqlManageFilterCardOptimizeCompleted: 优化完成
	SqlManageFilterCardOptimizeCompleted SqlManageFilterCard = "optimized"
)

type GlobalWorkflowListItem struct {
	WorkflowId   string                `json:"workflow_id"`   // 工单ID
	WorkflowType WorkflowType          `json:"workflow_type" enums:"sql_release,data_export"` // 工单类型
	WorkflowName string                `json:"workflow_name"` // 工单名称
	WorkflowDesc string                `json:"workflow_desc"` // 工单描述
	ProjectUid   string                `json:"project_uid"`   // 项目ID
	ProjectName  string                `json:"project_name"`  // 项目名称
	InstanceId   string                `json:"instance_id"`   // 实例ID
	InstanceName string                `json:"instance_name"` // 实例名称
	Assignee     string                `json:"assignee"`      // 当前处理人姓名
	Priority     string                `json:"priority"`      // High, Medium, Low
	Status       UnifiedWorkflowStatus `json:"status" enums:"pending_approval,pending_action,in_progress,exporting,rejected,cancelled,failed,completed,unknown"` // 工单状态
	CreatedAt    string                `json:"created_at"`    // 创建时间
}
