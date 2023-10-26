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

func getWorkflowRejectedPercentGroupByCreatorV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getWorkflowCounts(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getWorkflowDurationOfWaitingForAuditV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getWorkflowDurationOfWaitingForExecutionV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getWorkflowAuditPassPercentV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getWorkflowCreatedCountsEachDayV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getWorkflowStatusCountV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getWorkflowPercentCountedByInstanceTypeV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getSqlAverageExecutionTimeV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}

func getSqlExecutionFailPercentV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportStatistic)
}
