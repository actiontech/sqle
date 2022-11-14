package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/jinzhu/gorm"
)

// NOTE: related model:
// - RoleOperation, ProjectMemberRole, ProjectMemberGroupRole
type Role struct {
	Model
	Name string `gorm:"index"`
	Desc string
	Stat uint `json:"stat" gorm:"not null; default: 0; comment:'0:正常 1:被禁用'"`

	// todo issue960 remove Users, Instances, UserGroups
	Users      []*User      `gorm:"many2many:user_role;"`
	Instances  []*Instance  `gorm:"many2many:instance_role; comment:'关联实例'"`
	UserGroups []*UserGroup `gorm:"many2many:user_group_roles; comment:'关联用户组'"`
}

// NOTE: related model:
// - Role, User, Instance
type ProjectMemberRole struct {
	Model
	UserID     uint `json:"user_id" gorm:"not null"`
	InstanceID uint `json:"instance_id" gorm:"not null"`
	RoleID     uint `json:"role_id" gorm:"not null"`
}

// NOTE: related model:
// - Role, UserGroup, Instance
type ProjectMemberGroupRole struct {
	Model
	UserGroupID uint `json:"user_group_id" gorm:"not null"`
	InstanceID  uint `json:"instance_id" gorm:"not null"`
	RoleID      uint `json:"role_id" gorm:"not null"`
}

type BindRole struct {
	InstanceName string   `json:"instance_name" valid:"required"`
	RoleNames    []string `json:"role_names" valid:"required"`
}

func (s *Storage) UpdateUserRoles(userName, projectName string, bindRoles []BindRole) error {
	user, exist, err := s.GetUserByName(userName)
	if err != nil {
		return errors.ConnectStorageErrWrapper(err)
	}
	if !exist {
		return errors.ConnectStorageErrWrapper(fmt.Errorf("user not exist"))
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		return errors.ConnectStorageErrWrapper(s.updateUserRoles(tx, user, projectName, bindRoles))
	})
}

// 每次更新都是全量更新 InstID+UserID 定位到的角色
func (s *Storage) updateUserRoles(tx *gorm.DB, user *User, projectName string, bindRoles []BindRole) error {
	if len(bindRoles) == 0 {
		return nil
	}

	// 获取实例ID和角色ID
	instNames := []string{}
	roleNames := []string{}
	for _, role := range bindRoles {
		instNames = append(instNames, role.InstanceName)
		for _, name := range role.RoleNames {
			roleNames = append(roleNames, name)
		}
	}

	roleNames = utils.RemoveDuplicate(roleNames)
	instNames = utils.RemoveDuplicate(instNames)

	roles, err := s.GetRolesByNames(roleNames)
	if err != nil {
		return err
	}
	insts, err := s.GetInstancesByNamesAndProjectName(instNames, projectName)
	if err != nil {
		return err
	}

	roleCache := map[string /*role name*/ ]uint /*role id*/ {}
	instCache := map[string /*inst name*/ ]uint /*inst id*/ {}
	instIDs := []uint{}
	for _, role := range roles {
		roleCache[role.Name] = role.ID
	}

	for _, inst := range insts {
		instCache[inst.Name] = inst.ID
		instIDs = append(instIDs, inst.ID)
	}

	// 删掉所有旧数据
	err = tx.Where("user_id = ?", user.ID).Where("instance_id in (?)", instIDs).Delete(&ProjectMemberRole{}).Error
	if err != nil {
		return err
	}

	// 写入新数据
	for _, role := range bindRoles {
		for _, name := range role.RoleNames {
			if err = tx.Save(&ProjectMemberRole{
				RoleID:     roleCache[name],
				InstanceID: instCache[role.InstanceName],
				UserID:     user.ID,
			}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Storage) GetBindRolesByMemberNames(names []string, projectName string) (map[string /*member name*/ ][]BindRole, error) {
	roles := []*struct {
		UserName     string `json:"user_name"`
		InstanceName string `json:"instance_name"`
		RoleName     string `json:"role_name"`
	}{}

	err := s.db.Table("project_member_role").
		Select("users.login_name AS user_name , instances.name AS instance_name , roles.name AS role_name").
		Joins("JOIN users ON users.id = project_member_role.user_id").
		Joins("JOIN instances ON instances.id = project_member_role.instance_id").
		Joins("JOIN projects ON projects.id = instances.project_id").
		Joins("JOIN roles ON roles.id = project_member_role.role_id").
		Where("projects.name = ?", projectName).
		Where("users.login_name in (?)", names).
		Scan(&roles).Error

	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	removeDuplicate := map[string]struct{}{}
	resp := map[string][]BindRole{}

A:
	for _, role := range roles {
		// 去重
		fg := role.RoleName + role.InstanceName + role.UserName
		if _, ok := removeDuplicate[fg]; ok {
			continue
		}
		removeDuplicate[fg] = struct{}{}

		// resp中已有此用户+实例的信息时走这里
		for i, bindRole := range resp[role.UserName] {
			if bindRole.InstanceName == role.InstanceName {
				resp[role.UserName][i].RoleNames = append(resp[role.UserName][i].RoleNames, role.RoleName)
				continue A
			}
		}
		// resp还没记录过此用户或此用户+实例的信息时走这里
		resp[role.UserName] = append(resp[role.UserName], BindRole{
			InstanceName: role.InstanceName,
			RoleNames:    []string{role.RoleName},
		})
	}

	return resp, nil
}

func (s *Storage) GetRoleByName(name string) (*Role, bool, error) {
	role := &Role{}
	err := s.db.Where("name = ?", name).Find(role).Error
	if err == gorm.ErrRecordNotFound {
		return role, false, nil
	}
	return role, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRolesByNames(names []string) ([]*Role, error) {
	roles := []*Role{}
	err := s.db.Where("name in (?)", names).Find(&roles).Error
	return roles, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateRoleUsers(role *Role, users ...*User) error {
	err := s.db.Model(role).Association("Users").Replace(users).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateRoleInstances(role *Role, instances ...*Instance) error {
	err := s.db.Model(role).Association("Instances").Replace(instances).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllRoleTip() ([]*Role, error) {
	roles := []*Role{}
	err := s.db.Select("name").Find(&roles).Error
	return roles, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAndCheckRoleExist(roleNames []string) (roles []*Role, err error) {
	roles, err = s.GetRolesByNames(roleNames)
	if err != nil {
		return roles, err
	}
	existRoleNames := map[string]struct{}{}
	for _, role := range roles {
		existRoleNames[role.Name] = struct{}{}
	}
	notExistRoleNames := []string{}
	for _, roleName := range roleNames {
		if _, ok := existRoleNames[roleName]; !ok {
			notExistRoleNames = append(notExistRoleNames, roleName)
		}
	}
	if len(notExistRoleNames) > 0 {
		return roles, errors.New(errors.DataNotExist,
			fmt.Errorf("user role %s not exist", strings.Join(notExistRoleNames, ", ")))
	}
	return roles, nil
}

func (s *Storage) SaveRoleAndAssociations(role *Role,
	opCodes []uint) (err error) {
	return s.Tx(func(txDB *gorm.DB) (err error) {

		// save role
		if err = txDB.Save(role).Error; err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}

		// sync operations
		{
			if opCodes != nil {
				if err := s.ReplaceRoleOperationsByOpCodes(role.ID, opCodes); err != nil {
					return err
				}
			}
		}

		return
	})
}

func (s *Storage) DeleteRoleAndAssociations(role *Role) error {
	return s.Tx(func(txDB *gorm.DB) (err error) {

		// delete role
		if err = txDB.Delete(role).Error; err != nil {
			txDB.Rollback()
			return errors.ConnectStorageErrWrapper(err)
		}

		// delete role operations
		if err = s.DeleteRoleOperationByRoleID(role.ID); err != nil {
			txDB.Rollback()
			return err
		}

		return nil
	})
}

func (s *Storage) CheckRolesExist(roleNames []string) (bool, error) {
	roleNames = utils.RemoveDuplicate(roleNames)

	var count int
	err := s.db.Model(&Role{}).Where("name in (?)", roleNames).Count(&count).Error
	return len(roleNames) == count, errors.ConnectStorageErrWrapper(err)
}
