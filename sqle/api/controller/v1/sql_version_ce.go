//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var ErrCommunityEditionNotSupportSqlVersion = errors.New(errors.EnterpriseEditionFeatures, e.New("sql version is enterprise version feature"))

func createSqlVersion(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func getSqlVersionList(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func getSqlVersionDetail(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func updateSqlVersion(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func lockSqlVersion(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func deleteSqlVersion(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func getDependenciesBetweenStageInstance(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func batchReleaseWorkflows(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func batchExecuteWorkflows(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func batchAssociateWorkflowsWithVersion(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}

func getWorkflowsThatCanBeAssociatedToVersion(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlVersion
}
