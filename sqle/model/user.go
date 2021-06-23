package model

import (
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/universe/ucommon/v4/util"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

const DefaultAdminUser = "admin"

type User struct {
	Model
	// has created composite index: [id, login_name] by gorm#AddIndex
	Name           string `gorm:"index;column:login_name"`
	Email          string
	Password       string  `json:"-" gorm:"-"`
	SecretPassword string  `json:"secret_password" gorm:"not null;column:password"`
	Roles          []*Role `gorm:"many2many:user_role;"`
}

type Role struct {
	Model
	Name      string `gorm:"index"`
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
		log.NewEntry().Errorf("decrypt password for user %s failed, error: %v", i.Name, err)
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
		return roles, errors.New(errors.DataNotExist,
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
		return users, errors.New(errors.DataNotExist,
			fmt.Errorf("user %s not exist", strings.Join(notExistUserNames, ", ")))
	}
	return users, nil
}

func (s *Storage) UserCanAccessInstance(user *User, instance *Instance) (bool, error) {
	query := `SELECT count(instances.id) FROM users
JOIN user_role AS ur ON users.id = ur.user_id
JOIN roles ON ur.role_id = roles.id
JOIN instance_role AS ir ON roles.id = ir.role_id
JOIN instances ON ir.instance_id = instances.id
WHERE users.id = ? AND instances.id = ?
LIMIT 1
`
	var count uint
	err := s.db.Raw(query, user.ID, instance.ID).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	return count > 0, nil
}

func (s *Storage) UserCanAccessWorkflow(user *User, workflow *Workflow) (bool, error) {
	query := `SELECT count(w.id) FROM workflows AS w
JOIN workflow_records AS wr ON w.workflow_record_id = wr.id AND w.id = ?
LEFT JOIN workflow_steps AS cur_ws ON wr.current_workflow_step_id = cur_ws.id
LEFT JOIN workflow_step_templates AS cur_wst ON cur_ws.workflow_step_template_id = cur_wst.id
LEFT JOIN workflow_step_template_user AS cur_wst_re_user ON cur_wst.id = cur_wst_re_user.workflow_step_template_id
LEFT JOIN users AS cur_ass_user ON cur_wst_re_user.user_id = cur_ass_user.id
LEFT JOIN workflow_steps AS op_ws ON w.id = op_ws.workflow_id AND op_ws.state != "initialized"
LEFT JOIN workflow_step_templates AS op_wst ON op_ws.workflow_step_template_id = op_wst.id
LEFT JOIN workflow_step_template_user AS op_wst_re_user ON op_wst.id = op_wst_re_user.workflow_step_template_id
LEFT JOIN users AS op_ass_user ON op_wst_re_user.user_id = op_ass_user.id
where w.deleted_at IS NULL
AND (w.create_user_id = ? OR cur_ass_user.id = ? OR op_ass_user.id = ?)
`
	var count uint
	err := s.db.Raw(query, workflow.ID, user.ID, user.ID, user.ID).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	return count > 0, nil
}

func (s *Storage) UpdatePassword(user *User, password string) error {
	user.Password = password
	// User{}.encryptPassword(): SecretPassword为空时才会对密码进行加密操作
	user.SecretPassword = ""
	return s.Save(user)
}
