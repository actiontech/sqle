package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

// ========== AI 智能中心 ==========

// -------------------- 通用结构体 --------------------

// EfficiencyCard 效能卡片
type EfficiencyCard struct {
	MetricTitle string `json:"metric_title" enums:"security_defense,resource_cost,code_standard,rd_efficiency,query_performance"` // 效能指标标题

	MetricEvaluation string `json:"metric_evaluation" example:"S级,+450%,123h"` // 效能评价

	EvidenceValue string `json:"evidence_value" example:"高危拦截35次"` // 具体指标分值

	BusinessValue string `json:"business_value"` // 业务价值描述（仅战略价值接口使用）
}

// ProjectIOAnalysis 项目投入产出分析
type ProjectIOAnalysis struct {
	ProjectName     string  `json:"project_name"`                    // 项目名称
	ActiveMembers   int     `json:"active_members"`                  // 活跃人数
	InvokeCount     int64   `json:"invoke_count"`                    // 调用频次
	PerformanceGain string  `json:"performance_gain" example:"+35%"` // 性能提升
	TimeSaved       float64 `json:"time_saved" example:"45.5"`       // 节省工时
	HealthScore     float64 `json:"health_score" example:"92.3"`     // 健康分
}

// TopProblemDistribution Top问题分布
type TopProblemDistribution struct {
	ProblemType string  `json:"problem_type"` // 问题类型
	Percentage  float64 `json:"percentage"`   // 占比
}

// AIStrategicInsight AI战略洞察
type AIStrategicInsight struct {
	Title       string `json:"title"`       // 标题
	Description string `json:"description"` // 描述
}

// AIModuleAnalysis AI模块分析数据
type AIModuleAnalysis struct {
	AIModuleType string `json:"ai_module_type" enums:"smart_correction,performance_engine"` // AI模块类型

	ProjectIOAnalysis []ProjectIOAnalysis `json:"project_io_analysis"` // 项目组投入产出分析

	TopProblemDistribution []TopProblemDistribution `json:"top_problem_distribution"` // Top问题分布
}

// AIHubExecutionRecord 执行数据记录
type AIHubExecutionRecord struct {
	ID            uint64 `json:"id"`             // 记录ID
	SourceProject string `json:"source_project"` // 来源项目
	SQLSnippet    string `json:"sql_snippet"`    // SQL片段

	FunctionModule string `json:"function_module" enums:"smart_correction,performance_engine"` // 功能模块

	ValueDimension string `json:"value_dimension" enums:"security,performance,correction,maintenance,code_standard"` // 价值维度

	ProcessStatus string `json:"process_status" enums:"pending,running,completed,failed"` // 处理状态

	EstimatedUpgrade float64 `json:"estimated_upgrade"` // 预估提升
	OperationTime    string  `json:"operation_time"`    // 操作时间
}

// -------------------- 1. Banner 接口 --------------------

type GetAIHubBannerReq struct{}

type GetAIHubBannerResp struct {
	controller.BaseRes
	Data *AIHubBannerData `json:"data"`
}

type AIHubBannerData struct {
	Modules []AIModuleEfficiencyCards `json:"modules"` // 按模块分组的效能卡片
}

// AIModuleEfficiencyCards 模块效能卡片
type AIModuleEfficiencyCards struct {
	AIModuleType string `json:"ai_module_type" enums:"smart_correction,performance_engine"` // AI模块类型

	EfficiencyCards []EfficiencyCard `json:"efficiency_cards"` // 效能卡片列表
}

// GetAIHubBanner
// @Summary 获取AI智能中心Banner
// @Description 返回按模块分组的效能卡片列表（模块类型、效能指标、效能评价、具体指标）
// @Id GetAIHubBanner
// @Tags ai_hub
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} GetAIHubBannerResp
// @Router /v1/ai_hub/banner [get]
func GetAIHubBanner(c echo.Context) error {
	return getAIHubBanner(c)
}

// -------------------- 2. 战略价值接口 --------------------

type GetAIHubStrategicValueReq struct{}

type GetAIHubStrategicValueResp struct {
	controller.BaseRes
	Data *AIHubStrategicValueData `json:"data"`
}

type AIHubStrategicValueData struct {
	AIStrategicInsight *AIStrategicInsight `json:"ai_strategic_insight"` // AI战略洞察
	EfficiencyCards    []EfficiencyCard    `json:"efficiency_cards"`     // 效能卡片列表
}

// GetAIHubStrategicValue
// @Summary 获取AI智能中心战略价值
// @Description 返回AI战略洞察（标题、描述）和效能卡片列表（效能指标、效能评价、具体指标、价值描述）
// @Id GetAIHubStrategicValue
// @Tags ai_hub
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} GetAIHubStrategicValueResp
// @Router /v1/ai_hub/strategic_value [get]
func GetAIHubStrategicValue(c echo.Context) error {
	return getAIHubStrategicValue(c)
}

// -------------------- 3. 管理视图接口 --------------------

type GetAIHubManagementViewReq struct{}

type GetAIHubManagementViewResp struct {
	controller.BaseRes
	Data *AIHubManagementViewData `json:"data"`
}

type AIHubManagementViewData struct {
	Modules []AIModuleAnalysis `json:"modules"` // 按模块分组的分析数据
}

// GetAIHubManagementView
// @Summary 获取AI智能中心管理视图
// @Description 返回按模块分组的分析数据（模块类型、项目组投入产出分析、Top问题分布）
// @Id GetAIHubManagementView
// @Tags ai_hub
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} GetAIHubManagementViewResp
// @Router /v1/ai_hub/management_view [get]
func GetAIHubManagementView(c echo.Context) error {
	return getAIHubManagementView(c)
}

// -------------------- 4. 执行数据接口 --------------------

type GetAIHubExecutionDataResp struct {
	controller.BaseRes
	Data      []AIHubExecutionRecord `json:"data"`
	TotalNums uint64                 `json:"total_nums"` // 总数
}

// GetAIHubExecutionData
// @Summary 获取AI智能中心执行数据
// @Description 返回最新10条执行数据记录（来源项目、SQL片段、功能模块、价值维度、处理状态、预估提升、时间）
// @Id GetAIHubExecutionData
// @Tags ai_hub
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} GetAIHubExecutionDataResp
// @Router /v1/ai_hub/execution_data [get]
func GetAIHubExecutionData(c echo.Context) error {
	return getAIHubExecutionData(c)
}
