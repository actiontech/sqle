package v1

import (
	"context"
	"fmt"
	"strconv"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

func CheckCurrentUserCanOperateWorkflow(c echo.Context, projectUid string, workflow *model.Workflow, ops []dmsV1.OpPermissionType) error {
	userId := controller.GetUserID(c)
	up, err := dms.NewUserPermission(userId, projectUid)
	if err != nil {
		return err
	}
	if up.IsAdmin() {
		return nil
	}

	s := model.GetStorage()
	access, err := s.UserCanAccessWorkflow(userId, workflow)
	if err != nil {
		return err
	}
	if access {
		return nil
	}

	if len(ops) > 0 {
		for _, item := range workflow.Record.InstanceRecords {
			if !up.CanOpInstanceNoAdmin(item.Instance.GetIDStr(), ops...) {
				return ErrWorkflowNoAccess
			}
		}
		return nil
	}
	return ErrWorkflowNoAccess
}

func CheckCurrentUserCanOperateTasks(c echo.Context, projectUid string, workflow *model.Workflow, ops []dmsV1.OpPermissionType, taskIdList []uint) error {
	userId := controller.GetUserID(c)
	up, err := dms.NewUserPermission(userId, projectUid)
	if err != nil {
		return err
	}
	if up.IsAdmin() {
		return nil
	}

	s := model.GetStorage()

	access, err := s.UserCanAccessWorkflow(userId, workflow)
	if err != nil {
		return err
	}
	if access {
		return nil
	}

	if len(ops) > 0 {
		workflowInstances, err := s.GetWorkInstanceRecordByTaskIds(taskIdList)
		if err != nil {
			return err
		}

		instanceIds := make([]uint64, 0, len(workflowInstances))
		for _, item := range workflowInstances {
			instanceIds = append(instanceIds, item.InstanceId)
		}

		instances, err := dms.GetInstancesInProjectByIds(c.Request().Context(), projectUid, instanceIds)
		if err != nil {
			return err
		}
		for _, instance := range instances {
			if up.CanOpInstanceNoAdmin(instance.GetIDStr(), ops...) {
				return nil
			}
		}
	}

	return ErrWorkflowNoAccess
}

func checkCurrentUserCanAccessTask(c echo.Context, task *model.Task, ops []dmsV1.OpPermissionType) error {
	userId := controller.GetUserID(c)
	// todo issues-2005
	if task.Instance == nil || task.Instance.ProjectId == "" {
		return nil
	}
	up, err := dms.NewUserPermission(userId, task.Instance.ProjectId)
	if err != nil {
		return err
	}
	if up.IsAdmin() {
		return nil
	}
	if userId == fmt.Sprintf("%d", task.CreateUserId) {
		return nil
	}

	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowByTaskId(task.ID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.NewTaskNoExistOrNoAccessErr()
	}
	access, err := s.UserCanAccessWorkflow(userId, workflow)
	if err != nil {
		return err
	}
	if access {
		return nil
	}

	if up.CanOpInstanceNoAdmin(task.Instance.GetIDStr(), ops...) {
		return nil
	}

	return errors.NewTaskNoExistOrNoAccessErr()
}

func GetAuditPlanIfCurrentUserCanAccess(c echo.Context, projectId, auditPlanName string, opType v1.OpPermissionType) (*model.AuditPlan, bool, error) {
	storage := model.GetStorage()

	ap, exist, err := dms.GetAuditPlanWithInstanceFromProjectByName(projectId, auditPlanName, storage.GetAuditPlanFromProjectByName)
	if err != nil || !exist {
		return nil, exist, err
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return nil, false, err
	}

	if ap.CreateUserID == user.GetIDStr() {
		return ap, true, nil
	}

	_, isAdmin, err := dmsobject.GetUserOpPermission(c.Request().Context(), projectId, user.GetIDStr(), controller.GetDMSServerAddress())
	if err != nil {
		return nil, false, err
	}
	if isAdmin {
		return ap, true, nil
	}

	if opType != "" {
		instances, err := GetCanOperationInstances(c.Request().Context(), user, "", projectId, opType)
		if err != nil {
			return nil, false, errors.NewUserNotPermissionError(string(opType))
		}
		for _, instance := range instances {
			if ap.InstanceName == instance.Name {
				return ap, true, nil
			}
		}
	}
	return ap, false, errors.NewUserNotPermissionError(v1.GetOperationTypeDesc(opType))
}

func CheckCurrentUserCanAccessInstances(ctx context.Context, projectUID string, userId string, instances []*model.Instance) (bool, error) {
	up, err := dms.NewUserPermission(userId, projectUID)
	if err != nil {
		return false, fmt.Errorf("get user op permission from dms error: %v", err)
	}
	if up.IsAdmin() {
		return true, nil
	}
	for _, instance := range instances {
		if !up.CanOpInstanceNoAdmin(instance.GetIDStr(), dms.GetAllOpPermissions()...) {
			return false, nil
		}
	}
	return true, nil
}

func CheckCurrentUserCanCreateWorkflow(ctx context.Context, projectUID string, user *model.User, tasks []*model.Task) (bool, error) {
	up, err := dms.NewUserPermission(user.GetIDStr(), projectUID)
	if err != nil {
		return false, err
	}
	if up.IsAdmin() {
		return true, nil
	}

	instances := make([]*model.Instance, len(tasks))
	for i, task := range tasks {
		instances[i] = task.Instance
	}
	for _, instance := range instances {
		if !up.CanOpInstanceNoAdmin(instance.GetIDStr(), dmsV1.OpPermissionTypeCreateWorkflow) {
			return false, nil
		}
	}
	return true, nil
}

func CheckUserCanCreateAuditPlan(ctx context.Context, projectUID string, user *model.User, instances []*model.Instance) (bool, error) {
	up, err := dms.NewUserPermission(user.GetIDStr(), projectUID)
	if err != nil {
		return false, err
	}
	if up.IsAdmin() {
		return true, nil
	}
	for _, instance := range instances {
		if !up.CanOpInstanceNoAdmin(instance.GetIDStr(), dmsV1.OpPermissionTypeSaveAuditPlan) {
			return false, nil
		}
	}
	return true, nil
}

// 根据用户权限获取能访问/操作的实例列表
func GetCanOperationInstances(ctx context.Context, user *model.User, dbType, projectUid string, operationType v1.OpPermissionType) ([]*model.Instance, error) {
	// 获取当前项目下指定数据库类型的全部实例
	instances, err := dms.GetInstancesInProjectByType(ctx, projectUid, dbType)
	if err != nil {
		return nil, err
	}

	userOpPermissions, isAdmin, err := dmsobject.GetUserOpPermission(ctx, projectUid, user.GetIDStr(), controller.GetDMSServerAddress())
	if err != nil {
		return nil, err
	}

	if isAdmin || operationType == "" {
		return instances, nil
	}
	canOperationInstance := make([]*model.Instance, 0)
	for _, instance := range instances {
		if CanOperationInstance(userOpPermissions, []v1.OpPermissionType{operationType}, instance) {
			canOperationInstance = append(canOperationInstance, instance)
		}
	}
	return canOperationInstance, nil
}

func GetCanOpInstanceUsers(memberWithPermissions []*dmsV1.ListMembersForInternalItem, instance *model.Instance, opPermissioins []dmsV1.OpPermissionType) (opUsers []*model.User, err error) {
	opMapUsers := make(map[uint]struct{}, 0)
	for _, memberWithPermission := range memberWithPermissions {
		for _, memberOpPermission := range memberWithPermission.MemberOpPermissionList {
			if CanOperationInstance([]dmsV1.OpPermissionItem{memberOpPermission}, opPermissioins, instance) {
				opUser := new(model.User)
				userId, err := strconv.Atoi(memberWithPermission.User.Uid)
				if err != nil {
					return nil, err
				}
				opUser.ID = uint(userId)
				opUser.Name = memberWithPermission.User.Name
				if _, ok := opMapUsers[opUser.ID]; !ok {
					opMapUsers[opUser.ID] = struct{}{}
					opUsers = append(opUsers, opUser)
				}
			}
		}
	}
	return opUsers, nil
}

func CanOperationInstance(userOpPermissions []dmsV1.OpPermissionItem, needOpPermissionTypes []dmsV1.OpPermissionType, instance *model.Instance) bool {
	for _, userOpPermission := range userOpPermissions {
		// 对象权限(当前空间内所有对象)
		if userOpPermission.RangeType == dmsV1.OpRangeTypeProject {
			return true
		}

		// 动作权限(创建、审核、上线工单等)
		hasPrivilege := false
		for _, needOpPermissionType := range needOpPermissionTypes {
			if needOpPermissionType == userOpPermission.OpPermissionType {
				hasPrivilege = true
				break
			}
		}
		if !hasPrivilege {
			continue
		}
		// 对象权限(指定数据源)
		if userOpPermission.RangeType == dmsV1.OpRangeTypeDBService {
			for _, id := range userOpPermission.RangeUids {
				if id == instance.GetIDStr() {
					return true
				}
			}
		}
	}
	return false
}
