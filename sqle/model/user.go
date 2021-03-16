package model

import (
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/ucommon/v4/util"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

type User struct {
	Model
	Name           string `gorm:"column:login_name"`
	Email          string
	Password       string  `json:"-" gorm:"-"`
	SecretPassword string  `json:"secret_password" gorm:"not null;column:password"`
	Roles          []*Role `gorm:"many2many:user_role;"`
}

type Role struct {
	Model
	Name      string
	Desc      string
	Users     []*User     `gorm:"many2many:user_role;"`
	Instances []*Instance `gorm:"many2many:instance_role;"`
}

// BeforeSave is a hook implement gorm model before exec create
func (i *User) BeforeSave() error {
	return i.encryptPassword()
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db
func (i *User) AfterFind() error {
	err := i.decryptPassword()
	if err != nil {
		log.NewEntry().Errorf("decrypt password for user %d failed, error: %v", i.Name, err)
	}
	return nil
}

func (i *User) decryptPassword() error {
	if i == nil {
		return nil
	}
	if i.Password == "" {
		data, err := util.AesDecrypt(i.SecretPassword)
		if err != nil {
			return err
		} else {
			i.Password = string(data)
		}
	}
	return nil
}

func (i *User) encryptPassword() error {
	if i == nil {
		return nil
	}
	if i.SecretPassword == "" {
		data, err := util.AesEncrypt(i.Password)
		if err != nil {
			return err
		}
		i.SecretPassword = string(data)
	}
	return nil
}

func (s *Storage) GetUserByName(name string) (*User, bool, error) {
	t := &User{}
	err := s.db.Where("login_name = ?", name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetUserDetailByName(name string) (*User, bool, error) {
	t := &User{}
	err := s.db.Preload("Roles").Where("login_name = ?", name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateUserRoles(user *User, rs ...*Role) error {
	err := s.db.Model(user).Association("Roles").Replace(rs).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetUsersByNames(names []string) ([]*User, error) {
	users := []*User{}
	err := s.db.Where("login_name in (?)", names).Find(&users).Error
	return users, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetAllUserTip() ([]*User, error) {
	users := []*User{}
	err := s.db.Select("login_name").Find(&users).Error
	return users, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetRoleByName(name string) (*Role, bool, error) {
	role := &Role{}
	err := s.db.Where("name = ?", name).Find(role).Error
	if err == gorm.ErrRecordNotFound {
		return role, false, nil
	}
	return role, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetRolesByNames(names []string) ([]*Role, error) {
	roles := []*Role{}
	err := s.db.Where("name in (?)", names).Find(&roles).Error
	return roles, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateRoleUsers(role *Role, users ...*User) error {
	err := s.db.Model(role).Association("Users").Replace(users).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateRoleInstances(role *Role, instances ...*Instance) error {
	err := s.db.Model(role).Association("Instances").Replace(instances).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetAllRoleTip() ([]*Role, error) {
	roles := []*Role{}
	err := s.db.Select("name").Find(&roles).Error
	return roles, errors.New(errors.CONNECT_STORAGE_ERROR, err)
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
		return roles, errors.New(errors.DATA_NOT_EXIST,
			fmt.Errorf("user role %s not exist", strings.Join(notExistRoleNames, ", ")))
	}
	return roles, nil
}

func (s *Storage) GetAndCheckUserExist(userNames []string) (users []*User, err error) {
	users, err = s.GetUsersByNames(userNames)
	if err != nil {
		return users, err
	}
	existUserNames := map[string]struct{}{}
	for _, user := range users {
		existUserNames[user.Name] = struct{}{}
	}
	notExistUserNames := []string{}
	for _, userName := range userNames {
		if _, ok := existUserNames[userName]; !ok {
			notExistUserNames = append(notExistUserNames, userName)
		}
	}
	if len(notExistUserNames) > 0 {
		return users, errors.New(errors.DATA_NOT_EXIST,
			fmt.Errorf("user %s not exist", strings.Join(notExistUserNames, ", ")))
	}
	return users, nil
}
