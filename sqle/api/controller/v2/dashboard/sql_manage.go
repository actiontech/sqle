package dashboard

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	dashboardsvc "github.com/actiontech/sqle/sqle/server/dashboard"
	"github.com/labstack/echo/v4"
)

// --- SQL Manage Structures ---

type GlobalSqlManageStatisticsResV2 struct {
	controller.BaseRes
	Data GlobalSqlManageStatisticsV2 `json:"data"`
}

type GlobalSqlManageStatisticsV2 struct {
	PendingSqlCount        uint64 `json:"pending_sql_count"`         // 待我优化卡片，展示需要我处理的SQL数量
	OptimizedThisWeekCount uint64 `json:"optimized_this_week_count"` // 优化完成卡片，展示历史优化完成的SQL数量
}

type GetGlobalSqlManageStatisticsReqV2 struct {
	FilterProjectUid string `json:"filter_project_uid" query:"filter_project_uid"` // 项目ID
	FilterInstanceId string `json:"filter_instance_id" query:"filter_instance_id"` // 实例ID
}

type GetGlobalSqlManageTaskListReqV2 struct {
	// 分页
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"` // 页码
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`   // 每页条数

	// 筛选与搜索
	Keyword          string `json:"keyword" query:"keyword"`                       // 关键词模糊搜索
	FilterProjectUid string `json:"filter_project_uid" query:"filter_project_uid"` // 项目ID
	FilterInstanceId string `json:"filter_instance_id" query:"filter_instance_id"` // 实例ID

	// 根据卡片类型过滤
	FilterCard dashboardsvc.SqlManageFilterCard `json:"filter_card" query:"filter_card" valid:"omitempty,oneof=pending optimized" enums:"pending,optimized"`
}

type GlobalSqlManageTaskListResV2 struct {
	controller.BaseRes
	Data      []*GlobalSqlManageTaskItemV2 `json:"data"`
	TotalNums int64                        `json:"total_nums"`
}

type GlobalSqlManageTaskItemV2 struct {
	SqlFingerprint string   `json:"sql_fingerprint"`                                             // SQL 指纹
	AvgTime        *float64 `json:"avg_time"`                                                    // 平均耗时(秒)，来自采集 info；无慢日志等指标时 JSON 为 null（非 0）
	Count          *uint64  `json:"count"`                                                       // 执行频次；info 中无 counter 等时为 null（非 0）
	Suggestion     string   `json:"suggestion"`                                                  // 管理员派单备注；无备注时为空字符串
	Source         string   `json:"source"`                                                      // 来源，如：慢日志、库表元数据等，根据数据源id和项目id跳转到智能扫描详情
	ProjectUid     string   `json:"project_uid"`                                                 // 项目ID
	ProjectName    string   `json:"project_name"`                                                // 项目名称
	InstanceId     string   `json:"instance_id"`                                                 // 实例ID
	InstanceName   string   `json:"instance_name"`                                               // 实例名称
	LastSeenAt     string   `json:"last_seen_at"`                                                // 最后一次捕捉时间
	Status         string   `json:"status" enums:"unhandled,solved,ignored,manual_audited,sent"` // SQL治理任务状态：待处理，已优化，已忽略，人工审核中，已分配
}

// --- Handler Stubs ---

// GetGlobalSqlManageStatisticsV2
// @Summary 获取全局 SQL 治理统计看板
// @Description get global sql manage statistics;
// @Tags GlobalDashboard
// @Id GetGlobalSqlManageStatisticsV2
// @Security ApiKeyAuth
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id"
// @Success 200 {object} GlobalSqlManageStatisticsResV2
// @Router /v2/dashboard/sql_manage/statistics [get]
func GetGlobalSqlManageStatisticsV2(c echo.Context) error {
	return getGlobalSqlManageStatisticsV2(c)
}

// GetGlobalSqlManageTaskListV2
// @Summary 获取全局 SQL 治理任务列表
// @Description get global sql manage task list with filtering and pagination.
// @Description Pagination: 1-based page_index with offset (page_index-1)*page_size; page_index and page_size are required query parameters (non-zero).
// @Tags GlobalDashboard
// @Id GetGlobalSqlManageTaskListV2
// @Security ApiKeyAuth
// @Param page_index query uint32 true "1-based page index"
// @Param page_size query uint32 true "Number of items per page"
// @Param filter_card query string false "filter by card; pending 待我优化, optimized 优化完成" Enums(pending,optimized)
// @Param keyword query string false "fuzzy search keyword"
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id"
// @Success 200 {object} GlobalSqlManageTaskListResV2
// @Router /v2/dashboard/sql_manage/tasks [get]
func GetGlobalSqlManageTaskListV2(c echo.Context) error {
	return getGlobalSqlManageTaskListV2(c)
}
