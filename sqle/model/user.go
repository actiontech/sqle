package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/jinzhu/gorm"
)

const (
	DefaultAdminUser = "admin"
)

func IsDefaultAdminUser(user string) bool {
	return user == DefaultAdminUser
}

type UserAuthenticationType string

const (
	UserAuthenticationTypeLDAP UserAuthenticationType = "ldap" // user verify through ldap
	UserAuthenticationTypeSQLE UserAuthenticationType = "sqle" //user verify through sqle
)

type User struct {
	Model
	// has created composite index: [id, login_name] by gorm#AddIndex
	Name                   string `gorm:"index;column:login_name"`
	Email                  string
	Password               string                 `json:"-" gorm:"-"`
	SecretPassword         string                 `json:"secret_password" gorm:"not null;column:password"`
	UserAuthenticationType UserAuthenticationType `json:"user_authentication_type" gorm:"not null"`
	Roles                  []*Role                `gorm:"many2many:user_role;"`
	Stat                   uint                   `json:"stat" gorm:"not null; default: 0; comment:'0:正常 1:被禁用'"`

	WorkflowStepTemplates []*WorkflowStepTemplate `gorm:"many2many:workflow_step_template_user"`
}

func (u *User) IsDisabled() bool {
	return u.Stat == Disabled
}

func (u *User) SetStat(stat uint) {
	u.Stat = stat
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
		data, err := utils.AesDecrypt(i.SecretPassword)
		if err != nil {
			return err
		} else {
			i.Password = data
		}
	}
	return nil
}

func (i *User) encryptPassword() error {
	if i == nil {
		return nil
	}
	if i.SecretPassword == "" {
		data, err := utils.AesEncrypt(i.Password)
		if err != nil {
			return err
		}
		i.SecretPassword = data
	}
	return nil
}

func (s *Storage) GetUserByName(name string) (*User, bool, error) {
	t := &User{}
	err := s.db.Where("login_name = ?", name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetUserDetailByName(name string) (*User, bool, error) {
	t := &User{}
	err := s.db.Preload("Roles").Where("login_name = ?", name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateUserRoles(user *User, rs ...*Role) error {
	err := s.db.Model(user).Association("Roles").Replace(rs).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetUsersByNames(names []string) ([]*User, error) {
	users := []*User{}
	err := s.db.Where("login_name in (?)", names).Find(&users).Error
	return users, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllUserTip() ([]*User, error) {
	users := []*User{}
	err := s.db.Select("login_name").Where("stat=0").Find(&users).Error
	return users, errors.New(errors.ConnectStorageError, err)
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

	// 1. find role ids
	roles, err := s.GetRolesByUserID(int(user.ID))
	if err != nil {
		return false, err
	}

	if len(roles) == 0 {
		return false, nil
	}
	roleIDs := GetRoleIDsFromRoles(roles)

	// 2. check user access instance
	query := `
SELECT count(1) FROM instances
LEFT JOIN instance_role ON instances.id = instance_role.instance_id
LEFT JOIN roles ON instance_role.role_id = roles.id AND roles.stat=0
WHERE roles.id IN (?) AND instances.id = ?
`

	var count uint
	err = s.db.Raw(query, roleIDs, instance.ID).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return count > 0, nil
}

func (s *Storage) UserCanAccessWorkflow(user *User, workflow *Workflow) (bool, error) {
	query := `SELECT count(w.id) FROM workflows AS w
JOIN workflow_records AS wr ON w.workflow_record_id = wr.id AND w.id = ?
LEFT JOIN workflow_steps AS cur_ws ON wr.current_workflow_step_id = cur_ws.id
LEFT JOIN workflow_step_templates AS cur_wst ON cur_ws.workflow_step_template_id = cur_wst.id
LEFT JOIN workflow_step_template_user AS cur_wst_re_user ON cur_wst.id = cur_wst_re_user.workflow_step_template_id
LEFT JOIN users AS cur_ass_user ON cur_wst_re_user.user_id = cur_ass_user.id AND cur_ass_user.stat=0
LEFT JOIN workflow_steps AS op_ws ON w.id = op_ws.workflow_id AND op_ws.state != "initialized"
LEFT JOIN workflow_step_templates AS op_wst ON op_ws.workflow_step_template_id = op_wst.id
LEFT JOIN workflow_step_template_user AS op_wst_re_user ON op_wst.id = op_wst_re_user.workflow_step_template_id
LEFT JOIN users AS op_ass_user ON op_wst_re_user.user_id = op_ass_user.id AND op_ass_user.stat=0
where w.deleted_at IS NULL
AND (w.create_user_id = ? OR cur_ass_user.id = ? OR op_ass_user.id = ?)
`
	var count uint
	err := s.db.Raw(query, workflow.ID, user.ID, user.ID, user.ID).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return count > 0, nil
}

func (s *Storage) UpdatePassword(user *User, newPassword string) error {
	user.Password = newPassword
	// User{}.encryptPassword(): SecretPassword为空时才会对密码进行加密操作
	user.SecretPassword = ""
	return s.Save(user)
}

func (s *Storage) UserHasRunningWorkflow(userId uint) (bool, error) {
	// count how many running workflows have been assigned to this user
	query := `SELECT COUNT(user_id) FROM users
LEFT JOIN workflow_step_template_user wstu ON users.id = wstu.user_id
LEFT JOIN workflow_steps ws ON wstu.workflow_step_template_id = ws.workflow_step_template_id
LEFT JOIN workflow_records wr ON ws.workflow_record_id = wr.id
WHERE users.id = ? AND wr.status = ? AND ws.state = ?;`
	var count uint
	err := s.db.Raw(query, userId, WorkflowStatusRunning, WorkflowStepStateInit).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	if count > 0 {
		return true, nil
	}

	// count how many running workflows have been created by this user
	var workflows []*Workflow
	err = s.db.Model(workflows).
		Preload("Record", "status = ?", WorkflowStatusRunning).
		Where("create_user_id = ?", userId).
		Find(&workflows).Error
	return len(workflows) > 0 && workflows[0].Record != nil, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UserHasBindWorkflowTemplate(user *User) (bool, error) {
	copyUser := *user
	// 1 WorkflowTemplate to many WorkflowStepTemplates (delete: set NULL=set WorkflowStepTemplate.WorkflowTemplateId = NULL)
	// Many Users to many WorkflowStepTemplates
	err := s.db.Model(&copyUser).Preload("WorkflowStepTemplates", "workflow_template_id IS NOT NULL").Find(&copyUser).Error
	return len(copyUser.WorkflowStepTemplates) > 0, errors.New(errors.ConnectStorageError, err)
}
