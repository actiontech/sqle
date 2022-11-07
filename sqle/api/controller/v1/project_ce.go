//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/labstack/echo/v4"
)

var errCommunityEditionDoesNotSupportCreateProject = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support create project"))

var errCommunityEditionDoesNotSupportDeleteProject = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support delete project"))

func createProjectV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportCreateProject)
}

func deleteProjectV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportDeleteProject)
}
