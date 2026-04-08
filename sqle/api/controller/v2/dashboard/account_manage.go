package dashboard

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	dashboardsvc "github.com/actiontech/sqle/sqle/server/dashboard"
	"github.com/labstack/echo/v4"
)

// --- Account Management Structures ---

type GlobalAccountStatisticsData struct {
	ExpiringSoonCount  uint64 `json:"expiring_soon_count"`  // 即将过期卡片，展示即将过期的账号数量
	ActiveAccountCount uint64 `json:"active_account_count"` // 我的可用账号卡片，展示我的可用账号数量
}

type GlobalAccountStatisticsResV2 struct {
	controller.BaseRes
	Data GlobalAccountStatisticsData `json:"data"` // 统计数据
}

type GetGlobalAccountStatisticsReqV2 struct {
	FilterProjectUid string `json:"filter_project_uid" query:"filter_project_uid"` // 项目ID
	FilterInstanceId string `json:"filter_instance_id" query:"filter_instance_id"` // 实例ID
}

type GetGlobalAccountListReqV2 struct {
	// 分页
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"` // 页码
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`   // 每页条数

	// 筛选与搜索
	Keyword          string `json:"keyword" query:"keyword"`                       // 关键词模糊搜索
	FilterProjectUid string `json:"filter_project_uid" query:"filter_project_uid"` // 项目ID
	FilterInstanceId string `json:"filter_instance_id" query:"filter_instance_id"` // 实例ID

	// 根据卡片类型过滤
	FilterCard dashboardsvc.AccountFilterCard `json:"filter_card" query:"filter_card" valid:"omitempty,oneof=expiring_soon active" enums:"expiring_soon,active"`
}

type GlobalAccountListDataV2 struct {
	Accounts  []*GlobalAccountListItemV2 `json:"accounts"`
	CanManage bool                       `json:"can_manage"` // 是否有管理权限
}

type GlobalAccountListResV2 struct {
	controller.BaseRes
	Data GlobalAccountListDataV2 `json:"data"`
	TotalNums int64                      `json:"total_nums"`
}

type GlobalAccountListItemV2 struct {
	AccountUid   string `json:"account_uid"`                                   // 账号ID
	AccountName  string `json:"account_name"`                                  // 账号名称
	ProjectUid   string `json:"project_uid"`                                   // 项目ID
	ProjectName  string `json:"project_name"`                                  // 项目名称
	InstanceId   string `json:"instance_id"`                                   // 实例ID
	InstanceName string `json:"instance_name"`                                 // 实例名称
	ExpiredTime  string `json:"expired_time"`                                  // 过期时间，前端应显示剩余时间，若小于7天则需要用警告图标配合加粗红色文字警示
	Status       string `json:"status" enums:"active,expiring,expired,locked"` // 状态展示: 活跃，临期，已过期，已锁定
}

// --- Handler Stubs ---

// GetGlobalAccountStatisticsV2
// @Summary 获取全局账号管理统计数据
// @Description get global account statistics;
// @Tags GlobalDashboard
// @Id GetGlobalAccountStatisticsV2
// @Security ApiKeyAuth
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id"
// @Success 200 {object} GlobalAccountStatisticsResV2
// @Router /v2/dashboard/accounts/statistics [get]
func GetGlobalAccountStatisticsV2(c echo.Context) error {
	return getGlobalAccountStatisticsV2(c)
}

// GetGlobalAccountListV2
// @Summary 获取全局账号管理列表
// @Description get global account list with filtering and pagination.
// @Description Pagination: page_index is 1-based; offset is (page_index-1)*page_size. Values of page_index less than 1 are normalized to 1; page_size less than 1 may fall back to a default on the server.
// @Tags GlobalDashboard
// @Id GetGlobalAccountListV2
// @Security ApiKeyAuth
// @Param page_index query uint32 true "1-based page index"
// @Param page_size query uint32 true "Number of items per page"
// @Param filter_card query string false "filter by card; expiring_soon 即将过期, active 我的可用账号" Enums(expiring_soon,active)
// @Param keyword query string false "fuzzy search keyword"
// @Param filter_project_uid query string false "filter by project uid"
// @Param filter_instance_id query string false "filter by instance id"
// @Success 200 {object} GlobalAccountListResV2
// @Router /v2/dashboard/accounts [get]
func GetGlobalAccountListV2(c echo.Context) error {
	return getGlobalAccountListV2(c)
}
