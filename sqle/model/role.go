package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/jinzhu/gorm"
)

type Role struct {
	Model
	Name           string `gorm:"index"`
	Desc           string
	Stat           uint         `json:"stat" gorm:"not null; default: 0; comment:'0:正常 1:被禁用'"`
	Users          []*User      `gorm:"many2many:user_role;"`
	Instances      []*Instance  `gorm:"many2many:instance_role; comment:'关联实例'"`
	UserGroups     []*UserGroup `gorm:"many2many:user_group_roles; comment:'关联用户组'"`
	OperationCodes []*Operation `json:"op_codes" gorm:"column:op_codes; many2many:role_operations; comment:'动作权限'"`
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
