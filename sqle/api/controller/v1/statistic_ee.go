//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
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
	s := model.GetStorage()

	// 获取所有最后一位审核人审核通过的WorkStep
	allStepsHasAudit, err := getAllFinalAuditedPassWorkStepBO(s)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 获取所有上线成功的WorkStep
	allStepsHasOnline, err := getAllExecutedSuccessWorkStepBO(s)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var durationMin float64
	var count int
	for _, stepsHasOnline := range allStepsHasOnline {
		for _, stepsHasAudit := range allStepsHasAudit {
			if stepsHasAudit.WorkflowId == stepsHasOnline.WorkflowId {
				count++
				durationMin += stepsHasOnline.OperateAt.Sub(*stepsHasAudit.OperateAt).Minutes()
			}
		}
	}

	if count == 0 {
		return c.JSON(http.StatusOK, &GetWorkflowDurationOfWaitingForExecutionResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    &WorkflowStageDuration{Minutes: 0},
		})
	}

	averageOnlineMin := int(durationMin) / count

	return c.JSON(http.StatusOK, &GetWorkflowDurationOfWaitingForExecutionResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &WorkflowStageDuration{
			Minutes: uint(averageOnlineMin),
		},
	})
}

func getAllExecutedSuccessWorkStepBO(s *model.Storage) ([]*model.WorkFlowStepsBO, error) {
	return s.GetWorkFlowReverseStepsByIndexAndState(0, model.WorkflowStepStateApprove)
}

func getAllFinalAuditedPassWorkStepBO(s *model.Storage) ([]*model.WorkFlowStepsBO, error) {
	return s.GetWorkFlowReverseStepsByIndexAndState(1, model.WorkflowStepStateApprove)
}

func getWorkflowPassPercentV1(c echo.Context) error {
	auditPassPercent, err := getAuditPassPercent()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	executionSuccessPercent, err := getExecutionSuccessPercent()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetWorkflowPassPercentResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &WorkflowPassPercentV1{
			AuditPassPercent:        auditPassPercent * 100,
			ExecutionSuccessPercent: executionSuccessPercent * 100,
		},
	})
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

func getExecutionSuccessPercent() (float64, error) {
	s := model.GetStorage()
	successCount, err := s.GetWorkflowCountByTaskStatus([]string{model.TaskStatusExecuteSucceeded})
	if err != nil {
		return 0, err
	}
	allCount, err := s.GetAllWorkflowCount()
	if allCount == 0 {
		return 0, nil
	}
	return float64(successCount) / float64(allCount), err
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

func getWorkflowRejectedPercentGroupByInstanceV1(c echo.Context) error {
	req := new(GetWorkflowRejectedPercentGroupByInstanceReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	insts, err := s.GetAllInstanceTips("")
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var percents []InstanceRejectedPercent
	for _, inst := range insts {
		rejected, err := s.GetWorkflowCountByReq(map[string]interface{}{
			"filter_task_instance_name": inst.Name,
			"filter_status":             model.WorkflowStatusReject,
		})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		if rejected == 0 {
			continue
		}

		total, err := s.GetWorkflowCountByReq(map[string]interface{}{
			"filter_task_instance_name": inst.Name,
		})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		percent := float64(rejected) / float64(total) * 100
		percents = append(percents, InstanceRejectedPercent{
			InstanceName:     inst.Name,
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

	sort.Sort(InstanceRejectedPercents(percents))

	resItemCount := req.Limit
	actualPercentsCount := uint(len(percents))
	if req.Limit > actualPercentsCount {
		resItemCount = actualPercentsCount
	}

	percentsRes := make([]*WorkflowRejectedPercentGroupByInstance, resItemCount)
	for i := 0; i < int(resItemCount); i++ {
		percentsRes[i] = &WorkflowRejectedPercentGroupByInstance{
			InstanceName:     percents[i].InstanceName,
			WorkflowTotalNum: percents[i].WorkflowTotalNum,
			RejectedPercent:  percents[i].RejectedPercent,
		}
	}

	return c.JSON(http.StatusOK, &GetWorkflowRejectedPercentGroupByInstanceResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    percentsRes,
	})
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
	s := model.GetStorage()
	executionSuccessCount, err := s.GetWorkflowCountByTaskStatus([]string{model.TaskStatusExecuteSucceeded})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	executingCount, err := s.GetWorkflowCountByTaskStatus([]string{model.TaskStatusExecuting})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	executingFailedCount, err := s.GetWorkflowCountByTaskStatus([]string{model.TaskStatusExecuteFailed})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	waitingForExecutionCount, err := s.GetWorkflowCountByStepType([]string{model.WorkflowStepTypeSQLExecute})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	rejectedCount, err := s.GetWorkflowCountByStatus([]string{model.WorkflowStatusReject})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	waitingForAuditCount, err := s.GetWorkflowCountByStepType([]string{model.WorkflowStepTypeSQLReview})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	closedCount, err := s.GetWorkflowCountByStatus([]string{model.WorkflowStatusCancel})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
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
		instanceTypes := strings.Split(workflow.TaskInstanceType, ",")
		for _, instanceType := range instanceTypes {
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
