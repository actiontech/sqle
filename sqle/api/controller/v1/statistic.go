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
	return nil
}

type AverageDurationsV1 struct {
	AuditDurationMinutes     uint `json:"audit_duration_minutes"`
	ExecutionDurationMinutes uint `json:"execution_duration_minutes"`
}

type GetAverageOfTaskStageDurationsResV1 struct {
	controller.BaseRes
	Data *AverageDurationsV1 `json:"data"`
}

// GetAverageOfTaskStageDurationsV1
// @Summary 获取工单各阶段平均经历的时长
// @Description get average of durations of task's stages
// @Tags statistic
// @Id getAverageOfTaskStageDurationsV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetAverageOfTaskStageDurationsResV1
// @router /v1/statistic/tasks/stage_durations [get]
func GetAverageOfTaskStageDurationsV1(c echo.Context) error {
	return nil
}

type TaskPassPercentV1 struct {
	AuditPassPercent        uint `json:"audit_pass_percent"`
	ExecutionSuccessPercent uint `json:"execution_success_percent"`
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
	return nil
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
	return nil
}

type TaskStatusPercentV1 struct {
	ExecutionSuccessPercent    uint `json:"execution_success_percent"`
	ExecutingPercent           uint `json:"executing_percent"`
	WaitingForExecutionPercent uint `json:"waiting_for_execution_percent"`
	RejectedPercent            uint `json:"rejected_percent"`
	WaitingForAuditPercent     uint `json:"waiting_for_audit_percent"`
	ClosedPercent              uint `json:"closed_percent"`
}

type GetTaskStatusPercentResV1 struct {
	controller.BaseRes
	Data *TaskStatusPercentV1 `json:"data"`
}

// GetTaskStatusPercentV1
// @Summary 获取工单状态百分比
// @Description get percent of task status
// @Tags statistic
// @Id getTaskStatusPercentV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTaskStatusPercentResV1
// @router /v1/statistic/tasks/status_percent [get]
func GetTaskStatusPercentV1(c echo.Context) error {
	return nil
}

type TasksPercentCountedByInstanceType struct {
	InstanceType string `json:"instance_type"`
	Percent      uint   `json:"percent"`
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

type TaskRejectedPercentGroupByCreator struct {
	Creator         string `json:"creator"`
	TaskTotalNum    uint   `json:"task_total_num"`
	RejectedPercent uint   `json:"rejected_percent"`
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
	return nil
}

type TaskRejectedPercentGroupByInstance struct {
	InstanceName    string `json:"instance_name"`
	TaskTotalNum    uint   `json:"task_total_num"`
	RejectedPercent uint   `json:"rejected_percent"`
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
	return nil
}

type InstanceTypePercent struct {
	Type    string `json:"type"`
	Percent uint   `json:"percent"`
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
	return nil
}

type LicenseUsageV1 struct {
	UsedUsersPercent      uint                  `json:"used_users_percent"`
	UsedInstancesPercents []InstanceTypePercent `json:"used_instances_percents"`
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
	return nil
}
