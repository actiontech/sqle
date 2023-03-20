//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

func createProjectV1(c echo.Context) error {
	req := new(CreateProjectReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// check
	s := model.GetStorage()
	have, err := s.CheckUserHaveManagementPermission(user.ID, []uint{model.ManagementPermissionCreateProject})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !have {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("current user do not have permission to create project")))
	}

	_, exist, err := s.GetProjectByName(req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("project exist")))
	}

	// create
	err = s.CreateProject(req.Name, req.Desc, user.ID)

	return controller.JSONBaseErrorReq(c, err)
}

func deleteProjectV1(c echo.Context) error {
	userName := controller.GetUserName(c)

	projectName := c.Param("project_name")
	err := CheckIsProjectManager(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = checkProjectCanDelete(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	err = s.RemoveProject(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

func checkProjectCanDelete(projectName string) error {
	s := model.GetStorage()
	has, err := s.HasNotEndWorkflowByProjectName(projectName)
	if err != nil {
		return err
	}
	if has {
		return errors.New(errors.UserNotPermission, fmt.Errorf("there are unfinished work orders, and the current project cannot be deleted"))
	}
	return nil
}
