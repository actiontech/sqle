package dashboard

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	dashboard_svc "github.com/actiontech/sqle/sqle/server/dashboard"
	"github.com/labstack/echo/v4"
)

// --- Workflow Management Structures ---

type GetGlobalWorkflowStatisticsReqV2 struct {
	FilterProjectUid string `json:"filter_project_uid" query:"filter_project_uid"` // 项目ID
	FilterInstanceId string `json:"filter_instance_id" query:"filter_instance_id"` // 实例ID
}

type GlobalWorkflowStatisticsResV2 struct {
	controller.BaseRes
	Data GlobalWorkflowStatisticsV2 `json:"data"`
}

type GlobalWorkflowStatisticsV2 struct {
	ArchivedCount      uint64 `json:"archived_count"`        // 已完成的工单数量
	PendingForMeCount  uint64 `json:"pending_for_me_count"`  // 待我处理的工单数量
	InitiatedByMeCount uint64 `json:"initiated_by_me_count"` // 我发起的工单数量
}

type GetGlobalWorkflowListReqV2 struct {
	// 分页
	Cursor    string `json:"cursor" query:"cursor"`                        // 游标，仅在未指定 workflow_type 时有效，用于多源数据聚合分页
	PageIndex uint32 `json:"page_index" query:"page_index"`                // 页码，仅在指定了 workflow_type 时有效
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"` // 每页条数

	// 筛选与搜索
	Keyword          string `json:"keyword" query:"keyword"`                       // 关键词模糊搜索
	FilterProjectUid string `json:"filter_project_uid" query:"filter_project_uid"` // 项目ID
	FilterInstanceId string `json:"filter_instance_id" query:"filter_instance_id"` // 实例ID

	// 卡片过滤类型
	FilterCard dashboard_svc.GlobalWorkflowFilterCard `json:"filter_card" query:"filter_card" valid:"omitempty,oneof=archived pending_for_me initiated_by_me" enums:"archived,pending_for_me,initiated_by_me"`

	// 工单类型，sql_release: SQL上线工单，data_export: 数据导出工单
	WorkflowType dashboard_svc.WorkflowType `json:"workflow_type" query:"workflow_type" valid:"omitempty,oneof=sql_release data_export" enums:"sql_release,data_export"`
}

type GlobalWorkflowListResV2 struct {
	controller.BaseRes
	Data GlobalWorkflowListData `json:"data"`
}

type GlobalWorkflowListData struct {
	Workflows  []*dashboard_svc.GlobalWorkflowListItem `json:"workflows"`   // 工单列表
	NextCursor string                                  `json:"next_cursor"` // 游标分页：下一页游标
	TotalNums  int64                                   `json:"total_nums"`  // 总条数
	HasMore    bool                                    `json:"has_more"`    // 是否还有更多数据
}

// --- Handler Stubs ---

// GetGlobalWorkflowStatisticsV2
// @Summary 获取全局工单管理统计数据
// @Description get global workflow statistics, returns archived, pending_for_me, and initiated_by_me counts
// @Tags GlobalDashboard
// @Id GetGlobalWorkflowStatisticsV2
// @Security ApiKeyAuth
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id"
// @Success 200 {object} GlobalWorkflowStatisticsResV2
// @Router /v2/dashboard/workflows/statistics [get]
func GetGlobalWorkflowStatisticsV2(c echo.Context) error {
	return getGlobalWorkflowStatisticsV2(c)
}

// GetGlobalWorkflowListV2
// @Summary 获取全局工单管理列表
// @Description get global workflow list with filtering and pagination.
// @Description Pagination uses two mutually exclusive modes. Aggregate mode (omit workflow_type): merge sql_release and data_export ordered by time; use cursor (empty on first request, then data.next_cursor from the previous response) and page_size; page_index is not used for paging. Single-type mode (workflow_type set): list only that type with standard page_index (1-based) and page_size; cursor is ignored. Do not rely on mixing cursor tokens with workflow_type in one flow.
// @Tags GlobalDashboard
// @Id GetGlobalWorkflowListV2
// @Security ApiKeyAuth
// @Param cursor query string false "Aggregate mode only (omit workflow_type): opaque cursor string; leave empty on the first request, then pass the previous response's data.next_cursor. Ignored when workflow_type is set."
// @Param page_size query uint32 true "Page size; required in both aggregate and single-type modes."
// @Param page_index query uint32 true "Single-type mode only (workflow_type set): 1-based page index. When workflow_type is omitted, this value is ignored for pagination (use cursor instead)."
// @Param keyword query string false "fuzzy search keyword"
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id"
// @Param filter_card query string false "filter by card type; archived 已完成工单, pending_for_me 待我处理, initiated_by_me 我发起" Enums(archived,pending_for_me,initiated_by_me)
// @Param workflow_type query string false "filter by workflow type; sql_release SQL上线工单, data_export 数据导出工单" Enums(sql_release,data_export)
// @Success 200 {object} GlobalWorkflowListResV2
// @Router /v2/dashboard/workflows [get]
func GetGlobalWorkflowListV2(c echo.Context) error {
	return getGlobalWorkflowListV2(c)
}
