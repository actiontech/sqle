package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type WorkflowCountsV1 struct {
	Total      uint `json:"total"`
	TodayCount uint `json:"today_count"`
}

type GetWorkflowCountsResV1 struct {
	controller.BaseRes
	Data *WorkflowCountsV1 `json:"data"`
}

// GetWorkflowCountsV1
// @Summary 获取工单数量统计数据
// @Description get workflow counts
// @Tags statistic
// @Id getWorkflowCountV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowCountsResV1
// @router /v1/statistic/workflows/counts [get]
func GetWorkflowCountsV1(c echo.Context) error {
	return getWorkflowCounts(c)
}

type WorkflowStageDuration struct {
	Minutes uint `json:"minutes"`
}

type GetWorkflowDurationOfWaitingForAuditResV1 struct {
	controller.BaseRes
	Data *WorkflowStageDuration `json:"data"`
}

// GetWorkflowDurationOfWaitingForAuditV1
// @Summary 获取工单从创建到审核结束的平均时长
// @Description get duration from workflow being created to audited
// @Tags statistic
// @Id getWorkflowDurationOfWaitingForAuditV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowDurationOfWaitingForAuditResV1
// @router /v1/statistic/workflows/duration_of_waiting_for_audit [get]
func GetWorkflowDurationOfWaitingForAuditV1(c echo.Context) error {
	return getWorkflowDurationOfWaitingForAuditV1(c)
}

type GetSqlAverageExecutionTimeResV1 struct {
	controller.BaseRes
	Data []SqlAverageExecutionTime `json:"data"`
}

type SqlAverageExecutionTime struct {
	InstanceName            string `json:"instance_name"`
	AverageExecutionSeconds uint   `json:"average_execution_seconds"`
	MaxExecutionSeconds     uint   `json:"max_execution_seconds"`
	MinExecutionSeconds     uint   `json:"min_execution_seconds"`
}

// GetSqlAverageExecutionTimeV1
// @Summary sql上线平均耗时top10
// @Description get average execution time of sql
// @Tags statistic
// @Id getSqlAverageExecutionTimeV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSqlAverageExecutionTimeResV1
// @router /v1/statistic/instances/sql_average_execution_time [get]
func GetSqlAverageExecutionTimeV1(c echo.Context) error {
	return getSqlAverageExecutionTimeV1(c)
}

type GetWorkflowDurationOfWaitingForExecutionResV1 struct {
	controller.BaseRes
	Data *WorkflowStageDuration `json:"data"`
}

// GetWorkflowDurationOfWaitingForExecutionV1
// @Deprecated
// @Summary 获取工单各从审核完毕到执行上线的平均时长
// @Description get duration from workflow being created to executed
// @Tags statistic
// @Id getWorkflowDurationOfWaitingForExecutionV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowDurationOfWaitingForExecutionResV1
// @router /v1/statistic/workflows/duration_of_waiting_for_execution [get]
func GetWorkflowDurationOfWaitingForExecutionV1(c echo.Context) error {
	return getWorkflowDurationOfWaitingForExecutionV1(c)
}

type WorkflowPassPercentV1 struct {
	AuditPassPercent        float64 `json:"audit_pass_percent"`
	ExecutionSuccessPercent float64 `json:"execution_success_percent"`
}

type GetWorkflowPassPercentResV1 struct {
	controller.BaseRes
	Data *WorkflowPassPercentV1 `json:"data"`
}

// GetWorkflowPassPercentV1
// @Deprecated
// @Summary 获取工单通过率
// @Description get workflow pass percent
// @Tags statistic
// @Id getWorkflowPassPercentV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowPassPercentResV1
// @router /v1/statistic/workflows/pass_percent [get]
func GetWorkflowPassPercentV1(c echo.Context) error {
	return nil
}

type WorkflowAuditPassPercentV1 struct {
	AuditPassPercent float64 `json:"audit_pass_percent"`
}

type GetWorkflowAuditPassPercentResV1 struct {
	controller.BaseRes
	Data *WorkflowAuditPassPercentV1 `json:"data"`
}

// GetWorkflowAuditPassPercentV1
// @Summary 获取工单审核通过率
// @Description get workflow audit pass percent
// @Tags statistic
// @Id getWorkflowAuditPassPercentV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowAuditPassPercentResV1
// @router /v1/statistic/workflows/audit_pass_percent [get]
func GetWorkflowAuditPassPercentV1(c echo.Context) error {
	return getWorkflowAuditPassPercentV1(c)
}

type GetWorkflowCreatedCountsEachDayReqV1 struct {
	FilterDateFrom string `json:"filter_date_from" query:"filter_date_from" valid:"required"`
	FilterDateTo   string `json:"filter_date_to" query:"filter_date_to" valid:"required"`
}

type WorkflowCreatedCountsEachDayItem struct {
	Date  string `json:"date" example:"2022-08-24"`
	Value uint   `json:"value"`
}

type WorkflowCreatedCountsEachDayV1 struct {
	Samples []WorkflowCreatedCountsEachDayItem `json:"samples"`
}

type GetWorkflowCreatedCountsEachDayResV1 struct {
	controller.BaseRes
	Data *WorkflowCreatedCountsEachDayV1 `json:"data"`
}

// GetWorkflowCreatedCountsEachDayV1
// @Summary 获取每天工单创建数量
// @Description get counts of created workflow each day
// @Tags statistic
// @Id getWorkflowCreatedCountEachDayV1
// @Security ApiKeyAuth
// @Param filter_date_from query string true "filter date from.(format:yyyy-mm-dd)"
// @Param filter_date_to query string true "filter date to.(format:yyyy-mm-dd)"
// @Success 200 {object} v1.GetWorkflowCreatedCountsEachDayResV1
// @router /v1/statistic/workflows/each_day_counts [get]
func GetWorkflowCreatedCountsEachDayV1(c echo.Context) error {
	return getWorkflowCreatedCountsEachDayV1(c)
}

type WorkflowStatusCountV1 struct {
	ExecutionSuccessCount    int `json:"execution_success_count"`
	ExecutingCount           int `json:"executing_count"`
	ExecutingFailedCount     int `json:"executing_failed_count"`
	WaitingForExecutionCount int `json:"waiting_for_execution_count"`
	RejectedCount            int `json:"rejected_count"`
	WaitingForAuditCount     int `json:"waiting_for_audit_count"`
	ClosedCount              int `json:"closed_count"`
}

type GetWorkflowStatusCountResV1 struct {
	controller.BaseRes
	Data *WorkflowStatusCountV1 `json:"data"`
}

// GetWorkflowStatusCountV1
// @Summary 获取各种状态工单的数量
// @Description get count of workflow status
// @Tags statistic
// @Id getWorkflowStatusCountV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowStatusCountResV1
// @router /v1/statistic/workflows/status_count [get]
func GetWorkflowStatusCountV1(c echo.Context) error {
	return getWorkflowStatusCountV1(c)
}

type WorkflowPercentCountedByInstanceType struct {
	InstanceType string  `json:"instance_type"`
	Percent      float64 `json:"percent"`
	Count        uint    `json:"count"`
}

type WorkflowPercentCountedByInstanceTypeV1 struct {
	WorkflowPercents []WorkflowPercentCountedByInstanceType `json:"workflow_percents"`
	WorkflowTotalNum uint                                   `json:"workflow_total_num"`
}

type GetWorkflowPercentCountedByInstanceTypeResV1 struct {
	controller.BaseRes
	Data *WorkflowPercentCountedByInstanceTypeV1 `json:"data"`
}

// GetWorkflowPercentCountedByInstanceTypeV1
// @Summary 获取按数据源类型统计的工单百分比
// @Description get workflows percent counted by instance type
// @Tags statistic
// @Id getWorkflowPercentCountedByInstanceTypeV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowPercentCountedByInstanceTypeResV1
// @router /v1/statistic/workflows/instance_type_percent [get]
func GetWorkflowPercentCountedByInstanceTypeV1(c echo.Context) error {
	return getWorkflowPercentCountedByInstanceTypeV1(c)
}

type GetWorkflowRejectedPercentGroupByCreatorReqV1 struct {
	Limit uint `json:"limit" query:"limit" valid:"required"`
}

type WorkflowRejectedPercentGroupByCreator struct {
	Creator          string  `json:"creator"`
	WorkflowTotalNum uint    `json:"workflow_total_num"`
	RejectedPercent  float64 `json:"rejected_percent"`
}

type GetWorkflowRejectedPercentGroupByCreatorResV1 struct {
	controller.BaseRes
	Data []*WorkflowRejectedPercentGroupByCreator `json:"data"`
}

// GetWorkflowRejectedPercentGroupByCreatorV1
// @Summary 获取各个用户提交的工单驳回率，按驳回率降序排列
// @Description get workflows rejected percent group by creator. The result will be sorted by rejected percent in descending order
// @Tags statistic
// @Id getWorkflowRejectedPercentGroupByCreatorV1
// @Security ApiKeyAuth
// @Param limit query uint true "the limit of result item number"
// @Success 200 {object} v1.GetWorkflowRejectedPercentGroupByCreatorResV1
// @router /v1/statistic/workflows/rejected_percent_group_by_creator [get]
func GetWorkflowRejectedPercentGroupByCreatorV1(c echo.Context) error {
	return getWorkflowRejectedPercentGroupByCreatorV1(c)
}

type GetWorkflowRejectedPercentGroupByInstanceReqV1 struct {
	Limit uint `json:"limit" query:"limit" valid:"required"`
}

type WorkflowRejectedPercentGroupByInstance struct {
	InstanceName     string  `json:"instance_name"`
	WorkflowTotalNum uint    `json:"workflow_total_num"`
	RejectedPercent  float64 `json:"rejected_percent"`
}

type GetWorkflowRejectedPercentGroupByInstanceResV1 struct {
	controller.BaseRes
	Data []*WorkflowRejectedPercentGroupByInstance `json:"data"`
}

// GetWorkflowRejectedPercentGroupByInstanceV1
// @Summary 获取各个数据源相关的工单驳回率，按驳回率降序排列
// @Description get workflow rejected percent group by instance. The result will be sorted by rejected percent in descending order
// @Tags statistic
// @Id getWorkflowRejectedPercentGroupByInstanceV1
// @Security ApiKeyAuth
// @Param limit query uint true "the limit of result item number"
// @Success 200 {object} v1.GetWorkflowRejectedPercentGroupByInstanceResV1
// @router /v1/statistic/workflows/rejected_percent_group_by_instance [get]
func GetWorkflowRejectedPercentGroupByInstanceV1(c echo.Context) error {
	return getWorkflowRejectedPercentGroupByInstanceV1(c)
}

type InstanceTypePercent struct {
	Type    string  `json:"type"`
	Percent float64 `json:"percent"`
	Count   uint    `json:"count"`
}

type InstancesTypePercentV1 struct {
	InstanceTypePercents []InstanceTypePercent `json:"instance_type_percents"`
	InstanceTotalNum     uint                  `json:"instance_total_num"`
}

type GetInstancesTypePercentResV1 struct {
	controller.BaseRes
	Data *InstancesTypePercentV1 `json:"data"`
}

// GetInstancesTypePercentV1
// @Summary 获取数据源类型百分比
// @Description get database instances' types percent
// @Tags statistic
// @Id getInstancesTypePercentV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetInstancesTypePercentResV1
// @router /v1/statistic/instances/type_percent [get]
func GetInstancesTypePercentV1(c echo.Context) error {
	return getInstancesTypePercentV1(c)
}

type LicenseUsageItem struct {
	ResourceType     string `json:"resource_type"`
	ResourceTypeDesc string `json:"resource_type_desc"`
	Used             uint   `json:"used"`
	Limit            uint   `json:"limit"`
	IsLimited        bool   `json:"is_limited"`
}

type LicenseUsageV1 struct {
	UsersUsage     LicenseUsageItem   `json:"users_usage"`
	InstancesUsage []LicenseUsageItem `json:"instances_usage"`
}

type GetLicenseUsageResV1 struct {
	controller.BaseRes
	Data *LicenseUsageV1 `json:"data"`
}

// GetLicenseUsageV1
// @Summary 获取License使用情况
// @Description get usage of license
// @Tags statistic
// @Id getLicenseUsageV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetLicenseUsageResV1
// @router /v1/statistic/license/usage [get]
func GetLicenseUsageV1(c echo.Context) error {
	return getLicenseUsageV1(c)
}

type GetSqlExecutionFailPercentReqV1 struct {
	Limit uint `json:"limit" query:"limit" valid:"required"`
}

type SqlExecutionFailPercent struct {
	InstanceName string  `json:"instance_name"`
	Percent      float64 `json:"percent"`
}

type GetSqlExecutionFailPercentResV1 struct {
	controller.BaseRes
	Data []SqlExecutionFailPercent `json:"data"`
}

// GetSqlExecutionFailPercentV1
// @Summary 获取SQL上线失败率,按失败率降序排列
// @Description get sql execution fail percent
// @Tags statistic
// @Id getSqlExecutionFailPercentV1
// @Security ApiKeyAuth
// @Param limit query uint true "the limit of result item number"
// @Success 200 {object} v1.GetSqlExecutionFailPercentResV1
// @router /v1/statistic/instances/sql_execution_fail_percent [get]
func GetSqlExecutionFailPercentV1(c echo.Context) error {
	return getSqlExecutionFailPercentV1(c)
}
