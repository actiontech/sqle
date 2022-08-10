//go:build enterprise
// +build enterprise

package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
	"net/http"
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
