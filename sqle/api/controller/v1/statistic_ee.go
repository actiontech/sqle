//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

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
	return nil
}

func getLicenseUsageV1(c echo.Context) error {
	return nil
}

func getTaskRejectedPercentGroupByCreatorV1(c echo.Context) error {
	return nil
}

func getTaskRejectedPercentGroupByInstanceV1(c echo.Context) error {
	return nil
}
