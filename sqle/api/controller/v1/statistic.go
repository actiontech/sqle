package v1

import (
	"context"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/auditplan"

	"time"

	"math"

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

type GetSqlAverageExecutionTimeReqV1 struct {
	Limit uint `json:"limit" query:"limit" valid:"required"`
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
// @Summary 获取sql上线平均耗时，按平均耗时降序排列
// @Description get average execution time of sql
// @Tags statistic
// @Id getSqlAverageExecutionTimeV1
// @Security ApiKeyAuth
// @Param limit query uint true "the limit of result item number"
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
// @Deprecated
// @Summary 获取各个数据源相关的工单驳回率，按驳回率降序排列
// @Description get workflow rejected percent group by instance. The result will be sorted by rejected percent in descending order
// @Tags statistic
// @Id getWorkflowRejectedPercentGroupByInstanceV1
// @Security ApiKeyAuth
// @Param limit query uint true "the limit of result item number"
// @Success 200 {object} v1.GetWorkflowRejectedPercentGroupByInstanceResV1
// @router /v1/statistic/workflows/rejected_percent_group_by_instance [get]
func GetWorkflowRejectedPercentGroupByInstanceV1(c echo.Context) error {
	return nil
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
	// dms-todo: 移除 license
	return c.JSON(http.StatusOK, &GetLicenseUsageResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &LicenseUsageV1{
			UsersUsage: LicenseUsageItem{
				ResourceType:     "user",
				ResourceTypeDesc: "用户",
				Used:             0,
				Limit:            0,
				IsLimited:        false,
			},
			InstancesUsage: []LicenseUsageItem{},
		},
	})
	// return getLicenseUsageV1(c)
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

type GetProjectStatisticsResV1 struct {
	controller.BaseRes
	Data GetProjectStatisticsResDataV1 `json:"data"`
}

type GetProjectStatisticsResDataV1 struct {
	WorkflowTotal     uint64 `json:"workflow_total,omitempty"`
	AuditPlanTotal    uint64 `json:"audit_plan_total,omitempty"`
	InstanceTotal     uint64 `json:"instance_total,omitempty"`
	MemberTotal       uint64 `json:"member_total,omitempty"`
	RuleTemplateTotal uint64 `json:"rule_template_total,omitempty"`
	WhitelistTotal    uint64 `json:"whitelist_total,omitempty"`
}

// GetProjectStatisticsV1
// @Summary 获取项目统计信息
// @Description get project statistics
// @Tags statistic
// @Id getProjectStatisticsV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetProjectStatisticsResV1
// @router /v1/projects/{project_name}/statistics [get]
func GetProjectStatisticsV1(c echo.Context) error {
	// TODO 暂时不处理,解决页面报错
	// projectName := c.Param("project_name")
	// err := CheckIsProjectMember(controller.GetUserName(c), projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	resp := GetProjectStatisticsResDataV1{}
	// s := model.GetStorage()

	// resp.MemberTotal, err = s.GetUserTotalInProjectByProjectName(projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// resp.WhitelistTotal, err = s.GetSqlWhitelistTotalByProjectName(projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// resp.RuleTemplateTotal, err = s.GetRuleTemplateTotalByProjectName(projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// resp.InstanceTotal, err = s.GetInstanceTotalByProjectName(projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// resp.AuditPlanTotal, err = s.GetAuditPlanTotalByProjectName(projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// resp.WorkflowTotal, err = s.GetWorkflowTotalByProjectName(projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	return c.JSON(http.StatusOK, GetProjectStatisticsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resp,
	})
}

type AuditedSQLCount struct {
	TotalSQL uint `json:"total_sql_count"`
	RiskSQL  uint `json:"risk_sql_count"`
}

type StatisticsAuditedSQLResV1 struct {
	controller.BaseRes
	Data     AuditedSQLCount `json:"data"`
	RiskRate int             `json:"risk_rate"`
}

// StatisticsAuditedSQLV1
// @Summary 获取审核SQL总数，以及触发审核规则的SQL数量
// @Description statistics audited sql
// @Tags statistic
// @Id statisticsAuditedSQLV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.StatisticsAuditedSQLResV1
// @router /v1/projects/{project_name}/statistic/audited_sqls [get]
func StatisticsAuditedSQLV1(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	workflowSqlCount, err := s.GetSqlCountAndTriggerRuleCountFromWorkflowByProject(projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	auditPlanSqlCount, err := s.GetAuditPlanSQLCountAndTriggerRuleCountByProject(projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditedSQLCount := AuditedSQLCount{
		TotalSQL: workflowSqlCount.SqlCount + auditPlanSqlCount.SqlCount,
		RiskSQL:  workflowSqlCount.TriggerRuleCount + auditPlanSqlCount.TriggerRuleCount,
	}

	var riskRate float64
	if auditedSQLCount.RiskSQL > 0 {
		riskRate = float64(auditedSQLCount.RiskSQL) / float64(auditedSQLCount.TotalSQL)
	}

	return c.JSON(http.StatusOK, StatisticsAuditedSQLResV1{
		BaseRes:  controller.NewBaseReq(nil),
		Data:     auditedSQLCount,
		RiskRate: int(math.Round(riskRate * 100)),
	})
}

type dbErr struct {
	s   *model.Storage
	err error
}

func (d *dbErr) getWorkFlowStatusCountByProject(status string, projectUid string) (count int) {
	if d.err != nil {
		return 0
	}

	count, d.err = d.s.GetWorkflowCountByStatusAndProject(status, projectUid)

	return count
}

// StatisticWorkflowStatusV1
// @Summary 获取项目下工单各个状态的数量
// @Description statistic workflow status
// @Tags statistic
// @Id statisticWorkflowStatusV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetWorkflowStatusCountResV1
// @router /v1/projects/{project_name}/statistic/workflow_status [get]
func StatisticWorkflowStatusV1(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	d := &dbErr{s: model.GetStorage()}
	waitingForAuditCount := d.getWorkFlowStatusCountByProject(model.WorkflowStatusWaitForAudit, projectUid)
	waitingForExecutionCount := d.getWorkFlowStatusCountByProject(model.WorkflowStatusWaitForExecution, projectUid)
	executingCount := d.getWorkFlowStatusCountByProject(model.WorkflowStatusExecuting, projectUid)
	executionSuccessCount := d.getWorkFlowStatusCountByProject(model.WorkflowStatusFinish, projectUid)
	executingFailedCount := d.getWorkFlowStatusCountByProject(model.WorkflowStatusExecFailed, projectUid)
	rejectedCount := d.getWorkFlowStatusCountByProject(model.WorkflowStatusReject, projectUid)
	closedCount := d.getWorkFlowStatusCountByProject(model.WorkflowStatusCancel, projectUid)

	return c.JSON(http.StatusOK, &GetWorkflowStatusCountResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &WorkflowStatusCountV1{
			ExecutionSuccessCount:    executionSuccessCount,
			ExecutingCount:           executingCount,
			ExecutingFailedCount:     executingFailedCount,
			WaitingForExecutionCount: waitingForExecutionCount,
			RejectedCount:            rejectedCount,
			WaitingForAuditCount:     waitingForAuditCount,
			ClosedCount:              closedCount,
		},
	})
}

type RiskWorkflow struct {
	Name       string     `json:"workflow_name"`
	WorkflowID string     `json:"workflow_id"`
	Status     string     `json:"workflow_status"`
	CreateUser string     `json:"create_user_name"`
	UpdateTime *time.Time `json:"update_time"`
}

type StatisticRiskWorkflowResV1 struct {
	controller.BaseRes
	Data []*RiskWorkflow `json:"data"`
}

// StatisticRiskWorkflowV1
// @Summary 获取存在风险的工单
// @Description statistic risk workflow
// @Tags statistic
// @Id statisticRiskWorkflowV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.StatisticRiskWorkflowResV1
// @router /v1/projects/{project_name}/statistic/risk_workflow [get]
func StatisticRiskWorkflowV1(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	projectWorkflowStatusDetails, err := s.GetProjectWorkflowStatusDetail(projectUid, []string{model.WorkflowStatusReject, model.WorkflowStatusExecFailed})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	riskWorkflows := make([]*RiskWorkflow, len(projectWorkflowStatusDetails))
	for i, info := range projectWorkflowStatusDetails {
		user, err := func() (*model.User, error) {
			ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer cancel()

			return dms.GetUser(ctx, info.CreateUserId, controller.GetDMSServerAddress())
		}()

		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		riskWorkflows[i] = &RiskWorkflow{
			Name:       info.Subject,
			WorkflowID: info.WorkflowId,
			Status:     info.Status,
			CreateUser: user.Name,
			UpdateTime: info.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, StatisticRiskWorkflowResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    riskWorkflows,
	})
}

type AuditPlanCount struct {
	Type  string `json:"audit_plan_type"`
	Count uint   `json:"audit_plan_count"`
	Desc  string `json:"audit_plan_desc"`
}

type DBTypeAuditPlan struct {
	DBType string            `json:"db_type"`
	Data   []*AuditPlanCount `json:"data"`
}

type StatisticAuditPlanResV1 struct {
	controller.BaseRes
	Data []*DBTypeAuditPlan `json:"data"`
}

// StatisticAuditPlanV1
// @Summary 获取各类型数据源上的扫描任务数量
// @Description statistic audit plan
// @Tags statistic
// @Id statisticAuditPlanV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.StatisticAuditPlanResV1
// @router /v1/projects/{project_name}/statistic/audit_plans [get]
func StatisticAuditPlanV1(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	dBTypeAuditPlanCounts, err := s.GetDBTypeAuditPlanCountByProject(projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	dbTypeAuditPlanCountSliceMap := make(map[string][]*AuditPlanCount)
	for i := range dBTypeAuditPlanCounts {
		dbType := dBTypeAuditPlanCounts[i].DbType
		auditPlanCountSlice, exist := dbTypeAuditPlanCountSliceMap[dbType]
		if !exist {
			auditPlanCountSlice = []*AuditPlanCount{}
			dbTypeAuditPlanCountSliceMap[dbType] = auditPlanCountSlice
		}
		meta, err := auditplan.GetMeta(dBTypeAuditPlanCounts[i].Type)
		if err != nil {
			continue
		}
		newAuditPlanCount := &AuditPlanCount{
			Count: dBTypeAuditPlanCounts[i].AuditPlanCount,
			Type:  dBTypeAuditPlanCounts[i].Type,
			Desc:  meta.Desc,
		}
		dbTypeAuditPlanCountSliceMap[dbType] = append(auditPlanCountSlice, newAuditPlanCount)
	}

	dBTypeAuditPlanSlice := []*DBTypeAuditPlan{}
	for dbType := range dbTypeAuditPlanCountSliceMap {
		dBTypeAuditPlan := DBTypeAuditPlan{DBType: dbType, Data: dbTypeAuditPlanCountSliceMap[dbType]}
		dBTypeAuditPlanSlice = append(dBTypeAuditPlanSlice, &dBTypeAuditPlan)
	}

	return c.JSON(http.StatusOK, StatisticAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    dBTypeAuditPlanSlice,
	})
}

type RiskAuditPlan struct {
	ReportTimeStamp *time.Time `json:"audit_plan_report_timestamp"`
	ReportId        uint       `json:"audit_plan_report_id"`
	AuditPlanName   string     `json:"audit_plan_name"`
	TriggerTime     *time.Time `json:"trigger_audit_plan_time"`
	RiskSQLCount    uint       `json:"risk_sql_count"`
}

type GetRiskAuditPlanResV1 struct {
	controller.BaseRes
	Data []*RiskAuditPlan `json:"data"`
}

// GetRiskAuditPlanV1
// @Summary 获取扫描任务报告评分低于60的扫描任务
// @Description get risk audit plan
// @Tags statistic
// @Id getRiskAuditPlanV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetRiskAuditPlanResV1
// @router /v1/projects/{project_name}/statistic/risk_audit_plans [get]
func GetRiskAuditPlanV1(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	riskAuditPlanInfos, err := s.GetRiskAuditPlan(projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	riskAuditPlans := make([]*RiskAuditPlan, len(riskAuditPlanInfos))
	for i, info := range riskAuditPlanInfos {
		riskAuditPlans[i] = &RiskAuditPlan{
			ReportTimeStamp: info.ReportCreateAt,
			ReportId:        info.ReportId,
			AuditPlanName:   info.AuditPlanName,
			TriggerTime:     info.ReportCreateAt,
			RiskSQLCount:    info.RiskSqlCOUNT,
		}
	}

	return c.JSON(http.StatusOK, GetRiskAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    riskAuditPlans,
	})
}

type RoleUserCount struct {
	Role  string `json:"role"`
	Count uint   `json:"count"`
}

type GetRoleUserCountResV1 struct {
	controller.BaseRes
	Data []*RoleUserCount `json:"data"`
}

// GetRoleUserCountV1
// @Summary 获取各角色类型对应的成员数量
// @Description get role user count
// @Tags statistic
// @Id getRoleUserCountV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetRoleUserCountResV1
// @router /v1/projects/{project_name}/statistic/role_user [get]
func GetRoleUserCountV1(c echo.Context) error {
	// projectName := c.Param("project_name")

	// s := model.GetStorage()
	// userRoles, err := s.GetUserRoleByProjectName(projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }
	// userRoleFromUserGroup, err := s.GetUserRoleFromUserGroupByProjectName(projectName)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// // 若成员和成员组重复绑定相同角色，需要去重
	// userRoleMap := make(map[string]map[string]struct{}) // key 1: role_name, key 2: username
	// for i := range userRoles {
	// 	userRole := userRoles[i]
	// 	subMap, exist := userRoleMap[userRole.RoleName]
	// 	if !exist {
	// 		subMap = make(map[string]struct{})
	// 		userRoleMap[userRole.RoleName] = subMap
	// 	}
	// 	subMap[userRole.UserName] = struct{}{}
	// }
	// for i := range userRoleFromUserGroup {
	// 	userRole := userRoleFromUserGroup[i]
	// 	subMap, exist := userRoleMap[userRole.RoleName]
	// 	if !exist {
	// 		subMap = make(map[string]struct{})
	// 		userRoleMap[userRole.RoleName] = subMap
	// 	}
	// 	subMap[userRole.UserName] = struct{}{}
	// }

	// roleUserCountSlice := []*RoleUserCount{}
	// for k, v := range userRoleMap {
	// 	roleUserCountSlice = append(roleUserCountSlice, &RoleUserCount{Role: k, Count: uint(len(v))})
	// }

	return c.JSON(http.StatusOK, GetRoleUserCountResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    nil,
	})
}

type ProjectScore struct {
	Score int `json:"score"`
}

type GetProjectScoreResV1 struct {
	controller.BaseRes
	Data ProjectScore `json:"data"`
}

// GetProjectScoreV1
// @Summary 获取项目分数
// @Description get project score
// @Tags statistic
// @Id GetProjectScoreV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetProjectScoreResV1
// @router /v1/projects/{project_name}/statistic/project_score [get]
func GetProjectScoreV1(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	execFailedStatus := []string{model.WorkflowStatusExecFailed, model.WorkflowStatusFinish, model.WorkflowStatusExecuting}
	workflowsStatuses, err := s.GetProjectWorkflowStatusDetail(projectUid, execFailedStatus)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	auditReports, err := s.GetAuditPlanReportByProjectName(projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 项目分数=工单上线成功率*2/5+扫描任务通过率*3/5
	// 工单上线成功率=成功上线的工单数量/执行上线的工单总数
	// 扫描任务通过率=报告>60分的扫描任务数量/扫描任务报告总数
	var projectScore float64 = 1
	var workflowScore float64
	var auditReportScore float64

	workflowCount := len(workflowsStatuses)
	if workflowCount > 0 {
		workFlowFinishCount := 0
		for i := range workflowsStatuses {
			if workflowsStatuses[i].Status == model.WorkflowStatusFinish {
				workFlowFinishCount++
			}
		}
		workflowScore = float64(workFlowFinishCount) / float64(workflowCount)
	}
	auditReportCount := len(auditReports)
	if auditReportCount > 0 {
		auditReportPassCount := 0
		for i := range auditReports {
			if auditReports[i].Score > 60 {
				auditReportPassCount++
			}
		}
		auditReportScore = float64(auditReportPassCount) / float64(auditReportCount)
	}

	if workflowCount > 0 && auditReportCount > 0 {
		projectScore = workflowScore*2/5 + auditReportScore*3/5
	} else if workflowCount > 0 {
		projectScore = workflowScore
	} else if auditReportCount > 0 {
		projectScore = auditReportScore
	}

	return c.JSON(http.StatusOK, GetProjectScoreResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: ProjectScore{
			Score: int(math.Round(projectScore * 100)),
		},
	})
}

type DBTypeHealth struct {
	DBType             string   `json:"db_type"`
	HealthInstances    []string `json:"health_instance_names"`
	UnhealthyInstances []string `json:"unhealth_instance_names"`
}

type GetInstanceHealthResV1 struct {
	controller.BaseRes
	Data []*DBTypeHealth `json:"data"`
}

func getStringMapFromMap(stringMap map[string]map[string]struct{}, name string) map[string]struct{} {
	queryMap, exist := stringMap[name]
	if !exist {
		queryMap = make(map[string]struct{})
		stringMap[name] = queryMap
	}
	return queryMap
}

func keysFromMap(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func generateMaps(instanceWorkFlowFailedStatus []*model.InstanceWorkFlowStatusCount,
	latestAuditPlanReportScores []*model.LatestAuditPlanReportScore) (map[string]map[string]struct{}, map[string]map[string]struct{}) {

	instanceHealthMap := make(map[string]map[string]struct{}) // key1: db_type key2: instance name
	instanceUnhealthyMap := make(map[string]map[string]struct{})

	for i := range instanceWorkFlowFailedStatus {
		instanceStatus := instanceWorkFlowFailedStatus[i]
		dbType := instanceStatus.DbType
		if instanceStatus.StatusCount > 0 {
			unHealthMap := getStringMapFromMap(instanceUnhealthyMap, dbType)
			unHealthMap[instanceStatus.InstanceName] = struct{}{}
		} else {
			healthMap := getStringMapFromMap(instanceHealthMap, dbType)
			healthMap[instanceStatus.InstanceName] = struct{}{}
		}
	}
	for i := range latestAuditPlanReportScores {
		latestReportScore := latestAuditPlanReportScores[i]
		dbType := latestReportScore.DbType
		if latestReportScore.Score >= 60 {
			healthMap := getStringMapFromMap(instanceHealthMap, dbType)
			healthMap[latestReportScore.InstanceName] = struct{}{}
		} else {
			unhealthMap := getStringMapFromMap(instanceUnhealthyMap, dbType)
			unhealthMap[latestReportScore.InstanceName] = struct{}{}
		}
	}
	return instanceHealthMap, instanceUnhealthyMap
}

// GetInstanceHealthV1
// @Summary 获取各类型数据源的健康情况
// @Description get instance health
// @Tags statistic
// @Id GetInstanceHealthV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetInstanceHealthResV1
// @router /v1/projects/{project_name}/statistic/instance_health [get]
func GetInstanceHealthV1(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	instances, err := dms.GetInstancesInProject(c.Request().Context(), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	instanceWorkFlowFailedStatus, err := s.GetInstanceWorkFlowStatusCountByProject(instances, []string{model.WorkflowStatusReject, model.WorkflowStatusExecFailed})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceNames, err := dms.GetInstanceNamesInProject(c.Request().Context(), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	latestAuditPlanReportScores, err := s.GetLatestAuditPlanReportScoreFromInstanceByProject(projectUid, instanceNames)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceHealthMap, instanceUnhealthyMap := generateMaps(instanceWorkFlowFailedStatus, latestAuditPlanReportScores)

	dBTypeHealthMap := make(map[string]DBTypeHealth)
	// 从instanceHealthMap中去除在instanceUnhealthyMap中重复的instance name
	for dbType := range instanceUnhealthyMap {
		dBTypeHealth := DBTypeHealth{DBType: dbType}
		unhealthyInstanceNames := instanceUnhealthyMap[dbType]
		dBTypeHealth.UnhealthyInstances = keysFromMap(unhealthyInstanceNames)

		healthInstanceNames, exist := instanceHealthMap[dbType]
		if !exist {
			dBTypeHealthMap[dbType] = dBTypeHealth
			continue
		}

		for instanceName := range unhealthyInstanceNames {
			_, exist := healthInstanceNames[instanceName]
			if exist {
				delete(healthInstanceNames, instanceName)
			}
		}
		dBTypeHealth.HealthInstances = keysFromMap(healthInstanceNames)
		dBTypeHealthMap[dbType] = dBTypeHealth
	}

	for dbType := range instanceHealthMap {
		_, exist := dBTypeHealthMap[dbType]
		if exist {
			continue
		}
		dBTypeHealth := DBTypeHealth{DBType: dbType}
		dBTypeHealth.HealthInstances = keysFromMap(instanceHealthMap[dbType])
		dBTypeHealthMap[dbType] = dBTypeHealth
	}

	dBTypeHealth := []*DBTypeHealth{}
	for i := range dBTypeHealthMap {
		tmp := dBTypeHealthMap[i]
		dBTypeHealth = append(dBTypeHealth, &tmp)
	}

	return c.JSON(http.StatusOK, GetInstanceHealthResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    dBTypeHealth,
	})
}
