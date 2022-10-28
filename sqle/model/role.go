package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/jinzhu/gorm"
)

// NOTE: related model:
// - RoleOperation, ProjectMemberRole
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
	instances []*Instance, opCodes []uint, us []*User, ugs []*UserGroup) (err error) {
	return s.Tx(func(txDB *gorm.DB) (err error) {

		// save role
		if err = txDB.Save(role).Error; err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}

		// save instances
		{
			if instances != nil {
				if err = txDB.Model(role).Association("Instances").Replace(instances).Error; err != nil {
					return errors.ConnectStorageErrWrapper(err)
				}
			}
		}

		// save users
		{
			if us != nil {
				if err = txDB.Model(role).Association("Users").Replace(us).Error; err != nil {
					return errors.ConnectStorageErrWrapper(err)
				}
			}
		}

		// save user groups
		{
			if ugs != nil {
				if err = txDB.Model(role).Association("UserGroups").Replace(ugs).Error; err != nil {
					return errors.ConnectStorageErrWrapper(err)
				}
			}
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
