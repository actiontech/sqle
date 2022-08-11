//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"
	"sort"

	"github.com/actiontech/sqle/sqle/model"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

func getTaskCounts(c echo.Context) error {
	s := model.GetStorage()
	total, err := s.GetTaskCounts()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	todayCount, err := s.GetTaskCountsToday()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetTaskCountsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &TaskCountsV1{
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

func getTaskDurationOfWaitingForAuditV1(c echo.Context) error {
	return nil
}

func getTaskDurationOfWaitingForExecutionV1(c echo.Context) error {
	return nil
}

func getTaskPassPercentV1(c echo.Context) error {
	auditPassPercent, err := getAuditPassPercent()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	executionSuccessPercent, err := getExecutionSuccessPercent()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetTaskPassPercentResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &TaskPassPercentV1{
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
	return float64(passCount) / float64(allCount), err
}

func getExecutionSuccessPercent() (float64, error) {
	s := model.GetStorage()
	successCount, err := s.GetWorkflowCountByTaskStatus([]string{model.TaskStatusExecuteSucceeded})
	if err != nil {
		return 0, err
	}
	allCount, err := s.GetAllWorkflowCount()
	return float64(successCount) / float64(allCount), err
}

type CreatorRejectedPercent struct {
	Creator         string
	TaskTotalNum    uint
	RejectedPercent float64
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

func getTaskRejectedPercentGroupByCreatorV1(c echo.Context) error {
	req := new(GetTaskRejectedPercentGroupByCreatorReqV1)
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
			Creator:         user.Name,
			TaskTotalNum:    uint(total),
			RejectedPercent: percent,
		})
	}

	if percents == nil {
		return c.JSON(http.StatusOK, &GetTaskRejectedPercentGroupByCreatorResV1{
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

	percentsRes := make([]*TaskRejectedPercentGroupByCreator, resItemCount)
	for i := 0; i < int(resItemCount); i++ {
		percentsRes[i] = &TaskRejectedPercentGroupByCreator{
			Creator:         percents[i].Creator,
			TaskTotalNum:    percents[i].TaskTotalNum,
			RejectedPercent: percents[i].RejectedPercent,
		}
	}

	return c.JSON(http.StatusOK, &GetTaskRejectedPercentGroupByCreatorResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    percentsRes,
	})
}

type InstanceRejectedPercent struct {
	InstanceName    string
	TaskTotalNum    uint
	RejectedPercent float64
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

func getTaskRejectedPercentGroupByInstanceV1(c echo.Context) error {
	req := new(GetTaskRejectedPercentGroupByInstanceReqV1)
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
			InstanceName:    inst.Name,
			TaskTotalNum:    uint(total),
			RejectedPercent: percent,
		})
	}

	if percents == nil {
		return c.JSON(http.StatusOK, &GetTaskRejectedPercentGroupByCreatorResV1{
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

	percentsRes := make([]*TaskRejectedPercentGroupByInstance, resItemCount)
	for i := 0; i < int(resItemCount); i++ {
		percentsRes[i] = &TaskRejectedPercentGroupByInstance{
			InstanceName:    percents[i].InstanceName,
			TaskTotalNum:    percents[i].TaskTotalNum,
			RejectedPercent: percents[i].RejectedPercent,
		}
	}

	return c.JSON(http.StatusOK, &GetTaskRejectedPercentGroupByInstanceResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    percentsRes,
	})
}
