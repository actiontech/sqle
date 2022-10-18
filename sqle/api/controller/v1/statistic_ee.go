//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

func getWorkflowCounts(c echo.Context) error {
	s := model.GetStorage()
	total, err := s.GetAllWorkflowCount()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	t := time.Now()
	zeroClockToday := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	todayCount, err := s.GetWorkFlowCountBetweenStartTimeAndEndTime(zeroClockToday, time.Now())
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetWorkflowCountsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &WorkflowCountsV1{
			Total:      uint(total),
			TodayCount: uint(todayCount),
		},
	})
}

func getInstancesTypePercentV1(c echo.Context) error {
	s := model.GetStorage()

	typeCounts, err := s.GetAllInstanceCount()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	var total uint
	for _, count := range typeCounts {
		total += uint(count.Count)
	}

	percents := make([]InstanceTypePercent, len(typeCounts))
	for i, count := range typeCounts {
		percents[i] = InstanceTypePercent{
			Type:    count.DBType,
			Percent: float64(count.Count) / float64(total) * 100,
			Count:   uint(count.Count),
		}
	}

	return c.JSON(http.StatusOK, &GetInstancesTypePercentResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &InstancesTypePercentV1{
			InstanceTypePercents: percents,
			InstanceTotalNum:     total,
		},
	})
}

func getWorkflowDurationOfWaitingForAuditV1(c echo.Context) error {
	s := model.GetStorage()

	workFlowStepIdsHasAudit, err := s.GetWorkFlowStepIdsHasAudit()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	orderCount := len(workFlowStepIdsHasAudit)
	if orderCount == 0 {
		return c.JSON(http.StatusOK, &GetWorkflowDurationOfWaitingForAuditResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    &WorkflowStageDuration{Minutes: 0},
		})
	}

	durationMin, err := s.GetDurationMinHasAudit(workFlowStepIdsHasAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	averageMin := durationMin / orderCount

	return c.JSON(http.StatusOK, &GetWorkflowDurationOfWaitingForAuditResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    &WorkflowStageDuration{Minutes: uint(averageMin)},
	})
}

func getWorkflowDurationOfWaitingForExecutionV1(c echo.Context) error {
	//s := model.GetStorage()
	//
	//// 获取所有最后一位审核人审核通过的WorkStep
	//allStepsHasAudit, err := getAllFinalAuditedPassWorkStepBO(s)
	//if err != nil {
	//	return controller.JSONBaseErrorReq(c, err)
	//}
	//
	//// 获取所有上线成功的WorkStep
	//allStepsHasOnline, err := getAllExecutedSuccessWorkStepBO(s)
	//if err != nil {
	//	return controller.JSONBaseErrorReq(c, err)
	//}
	//
	//var durationMin float64
	//var count int
	//for _, stepsHasOnline := range allStepsHasOnline {
	//	for _, stepsHasAudit := range allStepsHasAudit {
	//		if stepsHasAudit.WorkflowId == stepsHasOnline.WorkflowId {
	//			count++
	//			durationMin += stepsHasOnline.OperateAt.Sub(*stepsHasAudit.OperateAt).Minutes()
	//		}
	//	}
	//}
	//
	//if count == 0 {
	//	return c.JSON(http.StatusOK, &GetWorkflowDurationOfWaitingForExecutionResV1{
	//		BaseRes: controller.NewBaseReq(nil),
	//		Data:    &WorkflowStageDuration{Minutes: 0},
	//	})
	//}
	//
	//averageOnlineMin := int(durationMin) / count

	return c.JSON(http.StatusOK, &GetWorkflowDurationOfWaitingForExecutionResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &WorkflowStageDuration{
			Minutes: 0,
		},
	})
}

func getAllExecutedSuccessWorkStepBO(s *model.Storage) ([]*model.WorkFlowStepsBO, error) {
	return s.GetWorkFlowReverseStepsByIndexAndState(0, model.WorkflowStepStateApprove)
}

func getAllFinalAuditedPassWorkStepBO(s *model.Storage) ([]*model.WorkFlowStepsBO, error) {
	return s.GetWorkFlowReverseStepsByIndexAndState(1, model.WorkflowStepStateApprove)
}

func getAuditPassPercent() (float64, error) {
	s := model.GetStorage()
	passCount, err := s.GetApprovedWorkflowCount()
	if err != nil {
		return 0, err
	}
	allCount, err := s.GetAllWorkflowCount()
	if allCount == 0 {
		return 0, nil
	}
	return float64(passCount) / float64(allCount), err
}

type CreatorRejectedPercent struct {
	Creator          string
	WorkflowTotalNum uint
	RejectedPercent  float64
}

type CreatorRejectedPercents []CreatorRejectedPercent

func (c CreatorRejectedPercents) Len() int {
	return len(c)
}
func (c CreatorRejectedPercents) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c CreatorRejectedPercents) Less(i, j int) bool {
	return c[j].RejectedPercent < c[i].RejectedPercent
}

func getWorkflowRejectedPercentGroupByCreatorV1(c echo.Context) error {
	req := new(GetWorkflowRejectedPercentGroupByCreatorReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	users, err := s.GetAllUserTip()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var percents []CreatorRejectedPercent
	for _, user := range users {
		rejected, err := s.GetWorkflowCountByReq(map[string]interface{}{
			"filter_create_user_name": user.Name,
			"filter_status":           model.WorkflowStatusReject,
		})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		if rejected == 0 {
			continue
		}

		total, err := s.GetWorkflowCountByReq(map[string]interface{}{
			"filter_create_user_name": user.Name,
		})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		percent := float64(rejected) / float64(total) * 100
		percents = append(percents, CreatorRejectedPercent{
			Creator:          user.Name,
			WorkflowTotalNum: uint(total),
			RejectedPercent:  percent,
		})
	}

	if percents == nil {
		return c.JSON(http.StatusOK, &GetWorkflowRejectedPercentGroupByCreatorResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    nil,
		})
	}

	sort.Sort(CreatorRejectedPercents(percents))

	actualPercentsCount := uint(len(percents))
	resItemCount := req.Limit
	if req.Limit > actualPercentsCount {
		resItemCount = actualPercentsCount
	}

	percentsRes := make([]*WorkflowRejectedPercentGroupByCreator, resItemCount)
	for i := 0; i < int(resItemCount); i++ {
		percentsRes[i] = &WorkflowRejectedPercentGroupByCreator{
			Creator:          percents[i].Creator,
			WorkflowTotalNum: percents[i].WorkflowTotalNum,
			RejectedPercent:  percents[i].RejectedPercent,
		}
	}

	return c.JSON(http.StatusOK, &GetWorkflowRejectedPercentGroupByCreatorResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    percentsRes,
	})
}

type InstanceRejectedPercent struct {
	InstanceName     string
	WorkflowTotalNum uint
	RejectedPercent  float64
}

type InstanceRejectedPercents []InstanceRejectedPercent

func (c InstanceRejectedPercents) Len() int {
	return len(c)
}
func (c InstanceRejectedPercents) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c InstanceRejectedPercents) Less(i, j int) bool {
	return c[j].RejectedPercent < c[i].RejectedPercent
}

func getWorkflowCreatedCountsEachDayV1(c echo.Context) error {
	req := new(GetWorkflowCreatedCountsEachDayReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	// parse date string
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, err))
	}
	dateFrom, err := time.ParseInLocation("2006-01-02", req.FilterDateFrom, loc)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, fmt.Errorf("parse dateFrom failed: %v", err)))
	}
	dateTo, err := time.ParseInLocation("2006-01-02", req.FilterDateTo, loc)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, fmt.Errorf("parse dateTo failed: %v", err)))
	}
	if dateFrom.After(dateTo) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("dateFrom must before dateTo")))
	}
	dateTo = dateTo.Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second) // 假设接口要查询第1天(date from)到第3天(date to)的趋势，那么第3天的工单创建数量是第3天0点到第23:59:59之间的数量

	var datePoints []time.Time
	currentDate := dateFrom
	for !currentDate.After(dateTo) {
		datePoints = append(datePoints, currentDate)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	s := model.GetStorage()
	counts, err := s.GetWorkflowDailyCountBetweenStartTimeAndEndTime(dateFrom, dateTo)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	samples := make([]WorkflowCreatedCountsEachDayItem, len(datePoints))
	for i, datePoint := range datePoints {
		workflowCount := 0
		for _, count := range counts {
			if datePoint.Equal(count.Date) {
				workflowCount = count.Count
				break
			}
		}
		samples[i] = WorkflowCreatedCountsEachDayItem{
			Date:  datePoint.Format("2006-01-02"),
			Value: uint(workflowCount),
		}
	}

	return c.JSON(http.StatusOK, &GetWorkflowCreatedCountsEachDayResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &WorkflowCreatedCountsEachDayV1{
			Samples: samples,
		},
	})
}

func getWorkflowStatusCountV1(c echo.Context) error {
	d := &dbErr{s: model.GetStorage()}

	waitingForAuditCount := d.getWorkFlowStatusCount(model.WorkflowStatusWaitForAudit)
	waitingForExecutionCount := d.getWorkFlowStatusCount(model.WorkflowStatusWaitForExecution)
	executingCount := d.getWorkFlowStatusCount(model.WorkflowStatusExecuting)
	executionSuccessCount := d.getWorkFlowStatusCount(model.WorkflowStatusFinish)
	executingFailedCount := d.getWorkFlowStatusCount(model.WorkflowStatusExecFailed)
	rejectedCount := d.getWorkFlowStatusCount(model.WorkflowStatusReject)
	closedCount := d.getWorkFlowStatusCount(model.WorkflowStatusCancel)
	if d.err != nil {
		return controller.JSONBaseErrorReq(c, d.err)
	}

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

type dbErr struct {
	s   *model.Storage
	err error
}

func (d *dbErr) getWorkFlowStatusCount(status string) (count int) {
	if d.err != nil {
		return 0
	}

	count, d.err = d.s.GetWorkflowCountByStatus(status)

	return count
}

func getWorkflowPercentCountedByInstanceTypeV1(c echo.Context) error {
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	workflows, total, err := s.GetWorkflowsByReq(map[string]interface{}{}, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("get workflows failed: %v", err))
	}

	instanceTypesMap := make(map[string]int, 0)

	for _, workflow := range workflows {
		for _, instanceType := range workflow.TaskInstanceType {
			instanceTypesMap[instanceType]++
		}
	}

	percents := make([]WorkflowPercentCountedByInstanceType, len(instanceTypesMap))
	i := 0
	for instType, count := range instanceTypesMap {
		percents[i] = WorkflowPercentCountedByInstanceType{
			InstanceType: instType,
			Percent:      float64(count) / float64(total) * 100,
			Count:        uint(count),
		}
		i++
	}

	return c.JSON(http.StatusOK, &GetWorkflowPercentCountedByInstanceTypeResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &WorkflowPercentCountedByInstanceTypeV1{
			WorkflowPercents: percents,
			WorkflowTotalNum: uint(total),
		},
	})
}

func getSqlAverageExecutionTimeV1(c echo.Context) error {
	req := new(GetSqlAverageExecutionTimeReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	sqlExecuteStatistics, err := s.GetSqlAvgExecutionTimeStatistic(req.Limit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instIds := make([]uint, 0, len(sqlExecuteStatistics))
	for _, statistic := range sqlExecuteStatistics {
		instIds = append(instIds, statistic.InstanceID)
	}

	instances, err := s.GetInstancesByIds(instIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	InstIdNameMap := make(map[uint] /*instance id*/ string /*instance name*/, 0)
	for _, instance := range instances {
		InstIdNameMap[instance.ID] = instance.Name
	}

	sqlAverageExecutionTimes := make([]SqlAverageExecutionTime, len(sqlExecuteStatistics))
	for i, executeStatistic := range sqlExecuteStatistics {
		sqlAverageExecutionTimes[i] = SqlAverageExecutionTime{
			InstanceName:            InstIdNameMap[executeStatistic.InstanceID],
			AverageExecutionSeconds: sqlExecuteStatistics[i].AvgExecutionTime,
			MaxExecutionSeconds:     sqlExecuteStatistics[i].MaxExecutionTime,
			MinExecutionSeconds:     sqlExecuteStatistics[i].MinExecutionTime,
		}
	}

	return c.JSON(http.StatusOK, &GetSqlAverageExecutionTimeResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    sqlAverageExecutionTimes,
	})
}

func getWorkflowAuditPassPercentV1(c echo.Context) error {
	auditPassPercent, err := getAuditPassPercent()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetWorkflowAuditPassPercentResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &WorkflowAuditPassPercentV1{
			AuditPassPercent: auditPassPercent * 100,
		},
	})
}

func getSqlExecutionFailPercentV1(c echo.Context) error {
	req := new(GetSqlExecutionFailPercentReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	sqlExecFailCount, err := s.GetSqlExecutionFailCount()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	sqlExecTotalCount, err := s.GetSqlExecutionTotalCount()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instIdExecTotalCountMap := make(map[uint] /*instance id*/ uint /*execute fail total count*/, 0)
	for _, totalCount := range sqlExecTotalCount {
		instIdExecTotalCountMap[totalCount.InstanceID] = totalCount.Count
	}

	instIds := make([]uint, 0, len(sqlExecFailCount))
	for _, failCount := range sqlExecFailCount {
		instIds = append(instIds, failCount.InstanceID)
	}

	instances, err := s.GetInstancesByIds(instIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instIdNameMap := make(map[uint] /*instance id*/ string /*instance name*/, 0)
	for _, inst := range instances {
		instIdNameMap[inst.ID] = inst.Name
	}

	executionFailPercents := make([]SqlExecutionFailPercent, 0, len(sqlExecFailCount))
	for _, failCount := range sqlExecFailCount {
		execTotalCount, ok := instIdExecTotalCountMap[failCount.InstanceID]
		if !ok {
			continue
		}

		executionFailPercents = append(executionFailPercents, SqlExecutionFailPercent{
			InstanceName: instIdNameMap[failCount.InstanceID],
			Percent:      float64(failCount.Count) / float64(execTotalCount) * 100,
		})
	}

	sort.Slice(executionFailPercents, func(i, j int) bool {
		return executionFailPercents[i].Percent > executionFailPercents[j].Percent
	})

	if len(executionFailPercents) > int(req.Limit) {
		executionFailPercents = executionFailPercents[:req.Limit]
	}

	return c.JSON(http.StatusOK, &GetSqlExecutionFailPercentResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    executionFailPercents,
	})
}
