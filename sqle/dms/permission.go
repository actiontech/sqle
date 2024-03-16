package dms

import (
	"context"
	"fmt"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/ungerik/go-dry"
)

type UserPermission struct {
	userId           string
	projectId        string
	isAdmin          bool
	opPermissionItem []v1.OpPermissionItem
}

func NewUserPermission(userId, projectId string) (*UserPermission, error) {
	opPermissions, isAdmin, err := dmsobject.GetUserOpPermission(context.TODO(), projectId, userId, controller.GetDMSServerAddress())
	if err != nil {
		return nil, fmt.Errorf("get user op permission from dms error: %v", err)
	}
	return &UserPermission{
		userId:           userId,
		projectId:        projectId,
		isAdmin:          isAdmin,
		opPermissionItem: opPermissions,
	}, nil
}

func (p *UserPermission) IsAdmin() bool {
	return p.isAdmin
}

func (p *UserPermission) IsProjectAdmin() bool {
	for _, userOpPermission := range p.opPermissionItem {
		if userOpPermission.RangeType == v1.OpRangeTypeProject {
			return true
		}
	}
	return false
}

// dms-todo: 1. 判断用户是 project 成员的方式成本高，看是否可以优化. 2. 捕捉错误.
func (p *UserPermission) IsProjectMember() bool {
	members, _, err := dmsobject.ListMembersInProject(context.TODO(), controller.GetDMSServerAddress(), v1.ListMembersForInternalReq{PageSize: 999, PageIndex: 1, ProjectUid: p.projectId})
	if err != nil {
		log.NewEntry().WithField("project_id", p.projectId).Errorln("fail to list member in project from dms")
		return false
	}
	for _, member := range members {
		if member.User.Uid == p.userId {
			return true
		}
	}
	return false
}

func (p *UserPermission) isPermissionMatch(opType v1.OpPermissionType, opTypes ...v1.OpPermissionType) bool {
	for i := range opTypes {
		if opTypes[i] == opType {
			return true
		}
	}
	return false
}

func (p *UserPermission) CanOpInstanceNoAdmin(instanceId string, OpTypes ...v1.OpPermissionType) bool {
	for _, userOpPermission := range p.opPermissionItem {
		// 判断是否是数据源权限
		if userOpPermission.RangeType != v1.OpRangeTypeDBService {
			continue
		}
		// 判断权限类型是否一致
		if !p.isPermissionMatch(userOpPermission.OpPermissionType, OpTypes...) {
			continue
		}
		// 判断权限对应的资源内有指定的数据源
		if dry.StringInSlice(instanceId, userOpPermission.RangeUids) {
			return true
		}
	}
	return false
}

func (p *UserPermission) GetInstancesByOP(OpTypes ...v1.OpPermissionType) []string {
	instances := []string{}
	instanceMap := map[string]struct{}{}

	for _, userOpPermission := range p.opPermissionItem {
		// 判断是否是数据源权限
		if userOpPermission.RangeType != v1.OpRangeTypeDBService {
			continue
		}
		// 判断权限类型是否一致
		if !p.isPermissionMatch(userOpPermission.OpPermissionType, OpTypes...) {
			continue
		}
		for _, id := range userOpPermission.RangeUids {
			if _, ok := instanceMap[id]; ok {
				continue
			}
			instances = append(instances, id)
			instanceMap[id] = struct{}{}
		}
	}
	return instances
}

func GetAllOpPermissions() []v1.OpPermissionType {
	return []v1.OpPermissionType{
		v1.OpPermissionTypeAuditWorkflow,
		v1.OpPermissionTypeCreateWorkflow,
		v1.OpPermissionTypeExecuteWorkflow,
		v1.OpPermissionTypeViewOthersWorkflow,
		v1.OpPermissionTypeSaveAuditPlan,
		v1.OpPermissionTypeViewOtherAuditPlan,
		v1.OpPermissionTypeSQLQuery,
		v1.OpPermissionTypeExportCreate,
		v1.OpPermissionTypeAuditExportWorkflow,
	}
}
