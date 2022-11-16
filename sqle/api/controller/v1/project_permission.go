package v1

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

/*

project permission.

*/

func CheckIsProjectMember(userName, projectName string) error {
	if userName == model.DefaultAdminUser {
		return nil
	}
	s := model.GetStorage()
	isMember, err := s.IsUserInProject(userName, projectName)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New(errors.UserNotPermission, fmt.Errorf("the project does not exist or user %v is not in project %v", userName, projectName))
	}
	return nil
}

func CheckIsProjectManager(userName, projectName string) error {
	if userName == model.DefaultAdminUser {
		return nil
	}
	s := model.GetStorage()
	isManager, err := s.IsProjectManager(userName, projectName)
	if err != nil {
		return err
	}
	if !isManager {
		return errors.New(errors.UserNotPermission, fmt.Errorf("the project does not exist or the user does not have permission to operate"))
	}
	return nil
}

/*

workflow permission.

*/

func CheckCurrentUserCanOperateWorkflow(c echo.Context, project *model.Project, workflow *model.Workflow, ops []uint) error {
	if controller.GetUserName(c) == model.DefaultAdminUser {
		return nil
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}

	s := model.GetStorage()

	isManager, err := s.IsProjectManager(user.Name, project.Name)
	if err != nil {
		return err
	}
	if isManager {
		return nil
	}

	access, err := s.UserCanAccessWorkflow(user, workflow)
	if err != nil {
		return err
	}
	if access {
		return nil
	}
	if len(ops) > 0 {
		instances, err := s.GetInstancesByWorkflowID(workflow.ID)
		if err != nil {
			return err
		}
		ok, err := s.CheckUserHasOpToInstances(user, instances, ops)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return ErrWorkflowNoAccess
}

func checkCurrentUserCanAccessTask(c echo.Context, task *model.Task, ops []uint) error {
	if controller.GetUserName(c) == model.DefaultAdminUser {
		return nil
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	if user.ID == task.CreateUserId {
		return nil
	}
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByTaskId(task.ID)
	if err != nil {
		return err
	}
	if !exist {
		return ErrTaskNoAccess
	}
	access, err := s.UserCanAccessWorkflow(user, workflow)
	if err != nil {
		return err
	}
	if access {
		return nil
	}
	if len(ops) > 0 {
		ok, err := s.CheckUserHasOpToInstances(user, []*model.Instance{task.Instance}, ops)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	return ErrTaskNoAccess
}

func CheckCurrentUserCanViewWorkflow(c echo.Context, workflow *model.Workflow) error {
	if controller.GetUserName(c) == model.DefaultAdminUser {
		return nil
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	s := model.GetStorage()
	access, err := s.UserCanAccessWorkflow(user, workflow)
	if err != nil {
		return err
	}
	if access {
		return nil
	}
	instances, err := s.GetInstancesByWorkflowID(workflow.ID)
	if err != nil {
		return err
	}
	ok, err := s.CheckUserHasOpToAnyInstance(user, instances, []uint{model.OP_WORKFLOW_VIEW_OTHERS})
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return ErrWorkflowNoAccess
}

func checkCurrentUserCanCreateWorkflow(user *model.User, tasks []*model.Task) error {
	if model.IsDefaultAdminUser(user.Name) {
		return nil
	}

	instances := make([]*model.Instance, len(tasks))
	for i, task := range tasks {
		instances[i] = task.Instance
	}

	s := model.GetStorage()
	ok, err := s.CheckUserHasOpToInstances(user, instances, []uint{model.OP_WORKFLOW_SAVE})
	if err != nil {
		return err
	}
	if !ok {
		return errors.NewAccessDeniedErr("user has no access to create workflow for instance")
	}

	return nil
}

/*

instance permission.

*/

// 1. admin user have all access to all instance
// 2. non-admin user have access to instance which is bound to one of his roles
func checkCurrentUserCanAccessInstance(c echo.Context, instance *model.Instance) (bool, error) {
	if controller.GetUserName(c) == model.DefaultAdminUser {
		return true, nil
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return false, err
	}
	s := model.GetStorage()
	access, err := s.UserCanAccessInstance(user, instance)
	if err != nil {
		return false, err
	}
	if !access {
		return false, nil
	}
	return true, nil
}

func checkCurrentUserCanAccessInstances(c echo.Context, instances []*model.Instance) (bool, error) {
	if len(instances) == 0 {
		return false, nil
	}

	for _, instance := range instances {
		can, err := checkCurrentUserCanAccessInstance(c, instance)
		if err != nil {
			return false, err
		}
		if !can {
			return false, nil
		}
	}

	return true, nil
}

/*

audit plan permission.

*/

func GetAuditPlanIfCurrentUserCanAccess(c echo.Context, projectName, auditPlanName string, opCode int) (*model.AuditPlan, bool, error) {
	storage := model.GetStorage()

	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return nil, false, err
	}

	ap, exist, err := storage.GetAuditPlanFromProjectByName(projectName, auditPlanName)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, true, nil
	}

	if controller.GetUserName(c) == model.DefaultAdminUser {
		return ap, true, nil
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return nil, false, err
	}

	if ap.CreateUserID == user.ID {
		return ap, true, nil
	}

	err = CheckIsProjectManager(userName, projectName)
	if err == nil {
		return ap, true, nil
	}

	if opCode > 0 {
		instances, err := storage.GetUserCanOpInstancesFromProject(user, projectName, []uint{uint(opCode)})
		if err != nil {
			return nil, false, errors.NewUserNotPermissionError(model.GetOperationCodeDesc(uint(opCode)))
		}
		for _, instance := range instances {
			if ap.InstanceName == instance.Name {
				return ap, true, nil
			}
		}
	}
	return nil, false, errors.NewUserNotPermissionError(model.GetOperationCodeDesc(uint(opCode)))
}
