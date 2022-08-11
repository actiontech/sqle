//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionDoesNotSupportStatistic = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support statistic"))

func getInstancesTypePercentV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getLicenseUsageV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getTaskRejectedPercentGroupByCreatorV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getTaskRejectedPercentGroupByInstanceV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getTaskCounts(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getTaskDurationOfWaitingForAuditV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getTaskDurationOfWaitingForExecutionV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getTaskPassPercentV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getTaskCreatedCountsEachDayV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getTasksPercentCountedByInstanceTypeV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}
