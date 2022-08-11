package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type TaskCountsV1 struct {
	Total      uint `json:"total"`
	TodayCount uint `json:"today_count"`
}

type GetTaskCountsResV1 struct {
	controller.BaseRes
	Data *TaskCountsV1 `json:"data"`
}

// GetTaskCountsV1
// @Summary 获取工单数量统计数据
// @Description get task counts
// @Tags statistic
// @Id getTaskCountV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskCountsResV1
// @router /v1/statistic/tasks/counts [get]
func GetTaskCountsV1(c echo.Context) error {
	return getTaskCounts(c)
}

type TaskStageDuration struct {
	Minutes uint `json:"minutes"`
}

type GetTaskDurationOfWaitingForAuditResV1 struct {
	controller.BaseRes
	Data *TaskStageDuration `json:"data"`
}

// GetTaskDurationOfWaitingForAuditV1
// @Summary 获取工单从创建到审核结束的平均时长
// @Description get duration from task being created to audited
// @Tags statistic
// @Id getTaskDurationOfWaitingForAuditV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskDurationOfWaitingForAuditResV1
// @router /v1/statistic/tasks/duration_of_waiting_for_audit [get]
func GetTaskDurationOfWaitingForAuditV1(c echo.Context) error {
	return getTaskDurationOfWaitingForAuditV1(c)
}

type GetTaskDurationOfWaitingForExecutionResV1 struct {
	controller.BaseRes
	Data *TaskStageDuration `json:"data"`
}

// GetTaskDurationOfWaitingForExecutionV1
// @Summary 获取工单各从审核完毕到执行上线的平均时长
// @Description get duration from task being created to executed
// @Tags statistic
// @Id getTaskDurationOfWaitingForExecutionV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskDurationOfWaitingForExecutionResV1
// @router /v1/statistic/tasks/duration_of_waiting_for_execution [get]
func GetTaskDurationOfWaitingForExecutionV1(c echo.Context) error {
	return getTaskDurationOfWaitingForExecutionV1(c)
}

type TaskPassPercentV1 struct {
	AuditPassPercent        float64 `json:"audit_pass_percent"`
	ExecutionSuccessPercent float64 `json:"execution_success_percent"`
}

type GetTaskPassPercentResV1 struct {
	controller.BaseRes
	Data *TaskPassPercentV1 `json:"data"`
}

// GetTaskPassPercentV1
// @Summary 获取工单通过率
// @Description get task pass percent
// @Tags statistic
// @Id getTaskPassPercentV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskPassPercentResV1
// @router /v1/statistic/tasks/pass_percent [get]
func GetTaskPassPercentV1(c echo.Context) error {
	return getTaskPassPercentV1(c)
}

type GetTaskCreatedCountsEachDayReqV1 struct {
	FilterDateFrom string `json:"filter_date_from" query:"filter_date_from" valid:"required"`
	FilterDateTo   string `json:"filter_date_to" query:"filter_date_to" valid:"required"`
}

type TaskCreatedCountsEachDayItem struct {
	Date  string `json:"date" example:"2022-08-24"`
	Value uint   `json:"value"`
}

type TaskCreatedCountsEachDayV1 struct {
	Samples []TaskCreatedCountsEachDayItem `json:"samples"`
}

type GetTaskCreatedCountsEachDayResV1 struct {
	controller.BaseRes
	Data *TaskCreatedCountsEachDayV1 `json:"data"`
}

// GetTaskCreatedCountsEachDayV1
// @Summary 获取每天工单创建数量
// @Description get counts of created task each day
// @Tags statistic
// @Id getTaskCreatedCountEachDayV1
// @Security ApiKeyAuth
// @Param filter_date_from query string true "filter date from.(format:yyyy-mm-dd)"
// @Param filter_date_to query string true "filter date to.(format:yyyy-mm-dd)"
// @Success 200 {object} v1.GetTaskCreatedCountsEachDayResV1
// @router /v1/statistic/tasks/each_day_counts [get]
func GetTaskCreatedCountsEachDayV1(c echo.Context) error {
	return getTaskCreatedCountsEachDayV1(c)
}

type TaskStatusCountV1 struct {
	ExecutionSuccessCount    int `json:"execution_success_count"`
	ExecutingCount           int `json:"executing_count"`
	WaitingForExecutionCount int `json:"waiting_for_execution_count"`
	RejectedCount            int `json:"rejected_count"`
	WaitingForAuditCount     int `json:"waiting_for_audit_count"`
	ClosedCount              int `json:"closed_count"`
}

type GetTaskStatusCountResV1 struct {
	controller.BaseRes
	Data *TaskStatusCountV1 `json:"data"`
}

// GetTaskStatusPercentV1
// @Summary 获取各种状态工单的数量
// @Description get count of task status
// @Tags statistic
// @Id getTaskStatusCountV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskStatusCountResV1
// @router /v1/statistic/tasks/status_count [get]
func GetTaskStatusCountV1(c echo.Context) error {
	return getTaskStatusCountV1(c)
}

type TasksPercentCountedByInstanceType struct {
	InstanceType string  `json:"instance_type"`
	Percent      float64 `json:"percent"`
}

type TasksPercentCountedByInstanceTypeV1 struct {
	TaskPercents []TasksPercentCountedByInstanceType `json:"task_percents"`
	TaskTotalNum uint                                `json:"task_total_num"`
}

type GetTasksPercentCountedByInstanceTypeResV1 struct {
	controller.BaseRes
	Data *TasksPercentCountedByInstanceTypeV1 `json:"data"`
}

// GetTasksPercentCountedByInstanceTypeV1
// @Summary 获取按数据源类型统计的工单百分比
// @Description get tasks percent counted by instance type
// @Tags statistic
// @Id getTasksPercentCountedByInstanceTypeV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTasksPercentCountedByInstanceTypeResV1
// @router /v1/statistic/tasks/instance_type_percent [get]
func GetTasksPercentCountedByInstanceTypeV1(c echo.Context) error {
	return nil
}

type GetTaskRejectedPercentGroupByCreatorReqV1 struct {
	Limit uint `json:"limit" query:"limit" valid:"required"`
}

type TaskRejectedPercentGroupByCreator struct {
	Creator         string  `json:"creator"`
	TaskTotalNum    uint    `json:"task_total_num"`
	RejectedPercent float64 `json:"rejected_percent"`
}

type GetTaskRejectedPercentGroupByCreatorResV1 struct {
	controller.BaseRes
	Data []*TaskRejectedPercentGroupByCreator `json:"data"`
}

// GetTaskRejectedPercentGroupByCreatorV1
// @Summary 获取各个用户提交的工单驳回率，按驳回率降序排列
// @Description get task rejected percent group by creator. The result will be sorted by rejected percent in descending order
// @Tags statistic
// @Id getTaskRejectedPercentGroupByCreatorV1
// @Security ApiKeyAuth
// @Param limit query uint true "the limit of result item number"
// @Success 200 {object} v1.GetTaskRejectedPercentGroupByCreatorResV1
// @router /v1/statistic/task/rejected_percent_group_by_creator [get]
func GetTaskRejectedPercentGroupByCreatorV1(c echo.Context) error {
	return getTaskRejectedPercentGroupByCreatorV1(c)
}

type GetTaskRejectedPercentGroupByInstanceReqV1 struct {
	Limit uint `json:"limit" query:"limit" valid:"required"`
}

type TaskRejectedPercentGroupByInstance struct {
	InstanceName    string  `json:"instance_name"`
	TaskTotalNum    uint    `json:"task_total_num"`
	RejectedPercent float64 `json:"rejected_percent"`
}

type GetTaskRejectedPercentGroupByInstanceResV1 struct {
	controller.BaseRes
	Data []*TaskRejectedPercentGroupByInstance `json:"data"`
}

// GetTaskRejectedPercentGroupByInstanceV1
// @Summary 获取各个数据源相关的工单驳回率，按驳回率降序排列
// @Description get task rejected percent group by instance. The result will be sorted by rejected percent in descending order
// @Tags statistic
// @Id getTaskRejectedPercentGroupByInstanceV1
// @Security ApiKeyAuth
// @Param limit query uint true "the limit of result item number"
// @Success 200 {object} v1.GetTaskRejectedPercentGroupByInstanceResV1
// @router /v1/statistic/task/rejected_percent_group_by_instance [get]
func GetTaskRejectedPercentGroupByInstanceV1(c echo.Context) error {
	return getTaskRejectedPercentGroupByInstanceV1(c)
}

type InstanceTypePercent struct {
	Type    string  `json:"type"`
	Percent float64 `json:"percent"`
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
