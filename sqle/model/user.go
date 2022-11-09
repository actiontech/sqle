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
	UserAuthenticationTypeLDAP   UserAuthenticationType = "ldap"   // user verify through ldap
	UserAuthenticationTypeSQLE   UserAuthenticationType = "sqle"   //user verify through sqle
	UserAuthenticationTypeOAUTH2 UserAuthenticationType = "oauth2" //user verify through oauth2
)

// NOTE: related model:
// - ProjectMemberRole, ManagementPermission
type User struct {
	Model
	// has created composite index: [id, login_name] by gorm#AddIndex
	Name                   string `gorm:"index;column:login_name"`
	Email                  string
	WeChatID               string                 `json:"wechat_id" gorm:"column:wechat_id"`
	Password               string                 `json:"-" gorm:"-"`
	SecretPassword         string                 `json:"secret_password" gorm:"not null;column:password"`
	UserAuthenticationType UserAuthenticationType `json:"user_authentication_type" gorm:"not null"`
	// todo issue960 remove Roles
	Roles            []*Role      `gorm:"many2many:user_role;"`
	UserGroups       []*UserGroup `gorm:"many2many:user_group_users"`
	Stat             uint         `json:"stat" gorm:"not null; default: 0; comment:'0:正常 1:被禁用'"`
	ThirdPartyUserID string       `json:"third_party_user_id"`

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

func (i *User) FingerPrint() string {
	return fmt.Sprintf(`{"id":"%v", "secret_password":"%v" }`, i.ID, i.SecretPassword)
}

func (s *Storage) GetUserByThirdPartyUserID(thirdPartyUserID string) (*User, bool, error) {
	t := &User{}
	err := s.db.Where("third_party_user_id = ?", thirdPartyUserID).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
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
	err := s.db.Preload("Roles").Preload("UserGroups").
		Where("login_name = ?", name).First(t).Error
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

func (s *Storage) GetUserTipsByProject(projectName string) ([]*User, error) {
	if projectName == "" {
		return s.GetAllUserTip()
	}

	query := `
SELECT users.login_name 
FROM users
JOIN project_user on project_user.user_id = users.id
JOIN projects on project_user.project_id = projects.id
JOIN user_group_users on users.id = user_group_users.user_id 
JOIN project_user_group on user_group_users.user_group_id = project_user_group.user_group_id
JOIN projects as p on project_user_group.project_id = p.id
WHERE users.stat = 0
AND( 
	projects.name = ?
OR
	p.name = ?
)
`

	var users []*User
	err := s.db.Raw(query, projectName, projectName).Scan(&users).Error

	return users, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllUserTip() ([]*User, error) {
	users := []*User{}
	err := s.db.Select("login_name").Where("stat=0").Find(&users).Error
	return users, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllUserCount() (int64, error) {
	var count int64
	return count, s.db.Model(&User{}).Count(&count).Error
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

func (s *Storage) UserCanAccessInstance(user *User, instance *Instance) (
	ok bool, err error) {

	if IsDefaultAdminUser(user.Name) {
		return true, nil
	}

	type countStruct struct {
		Count int `json:"count"`
	}

	query := `
SELECT COUNT(1) AS count
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.stat = 0 AND roles.deleted_at IS NULL
LEFT JOIN user_role ON user_role.role_id = roles.id
LEFT JOIN users ON users.id = user_role.user_id AND users.stat = 0 AND users.deleted_at IS NULL
WHERE instances.deleted_at IS NULL
AND instances.id = ?
AND users.id = ?
GROUP BY instances.id
UNION
SELECT instances.id
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.stat = 0 AND roles.deleted_at IS NULL
JOIN user_group_roles ON roles.id = user_group_roles.role_id
JOIN user_groups ON user_groups.id = user_group_roles.user_group_id AND user_groups.stat = 0 AND user_groups.deleted_at IS NULL
JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
JOIN users ON users.id = user_group_users.user_id AND users.stat = 0 AND users.deleted_at IS NULL
WHERE instances.deleted_at IS NULL
AND instances.id = ?
AND users.id = ?
GROUP BY instances.id
`
	var cnt countStruct
	err = s.db.Unscoped().Raw(query, instance.ID, user.ID, instance.ID, user.ID).Scan(&cnt).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return cnt.Count > 0, nil
}

func (s *Storage) UserCanAccessWorkflow(user *User, workflow *Workflow) (bool, error) {
	query := `SELECT count(w.id) FROM workflows AS w
JOIN workflow_records AS wr ON w.workflow_record_id = wr.id AND w.id = ?
LEFT JOIN workflow_steps AS cur_ws ON wr.current_workflow_step_id = cur_ws.id
LEFT JOIN workflow_step_templates AS cur_wst ON cur_ws.workflow_step_template_id = cur_wst.id
LEFT JOIN workflow_step_user AS cur_wst_re_user ON cur_ws.id = cur_wst_re_user.workflow_step_id
LEFT JOIN users AS cur_ass_user ON cur_wst_re_user.user_id = cur_ass_user.id AND cur_ass_user.stat=0
LEFT JOIN workflow_steps AS op_ws ON w.id = op_ws.workflow_id AND op_ws.state != "initialized"
LEFT JOIN workflow_step_templates AS op_wst ON op_ws.workflow_step_template_id = op_wst.id
LEFT JOIN workflow_step_user AS op_wst_re_user ON op_ws.id = op_wst_re_user.workflow_step_id
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
LEFT JOIN workflow_step_user wstu ON users.id = wstu.user_id
LEFT JOIN workflow_steps ws ON wstu.workflow_step_id = ws.id
LEFT JOIN workflow_records wr ON ws.workflow_record_id = wr.id
WHERE users.id = ? AND wr.status IN (?) AND ws.state = ?;`
	var count uint
	err := s.db.Raw(query, userId, []string{WorkflowStatusWaitForAudit, WorkflowStatusWaitForExecution}, WorkflowStepStateInit).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	if count > 0 {
		return true, nil
	}

	// count how many running workflows have been created by this user
	var workflows []*Workflow
	err = s.db.Model(workflows).
		Preload("Record", "status IN (?)", []string{WorkflowStatusWaitForAudit, WorkflowStatusWaitForExecution}).
		Where("create_user_id = ?", userId).
		Find(&workflows).Error
	return len(workflows) > 0 && workflows[0].Record != nil, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UserHasBindWorkflowTemplate(user *User) (bool, error) {
	count := 0
	err := s.db.Table("workflow_templates").
		Joins("join workflow_step_templates on workflow_templates.id = workflow_step_templates.workflow_template_id").
		Joins("join workflow_step_user on workflow_step_templates.id = workflow_step_user.workflow_step_id").
		Where("workflow_templates.deleted_at is null").
		Where("workflow_step_user.user_id = ?", user.ID).
		Count(&count).Error
	return count > 0, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) SaveUserAndAssociations(
	user *User, userGroups []*UserGroup, managementPermissionCodes *[]uint) (err error) {
	return s.Tx(func(txDB *gorm.DB) error {

		// User
		if err := txDB.Save(user).Error; err != nil {
			txDB.Rollback()
			return errors.ConnectStorageErrWrapper(err)
		}

		// user groups
		if userGroups != nil {
			if err := txDB.Model(user).
				Association("UserGroups").
				Replace(userGroups).Error; err != nil {
				txDB.Rollback()
				return errors.ConnectStorageErrWrapper(err)
			}
		}

		// permission
		if managementPermissionCodes != nil {
			if err := updateManagementPermission(txDB, user.ID, *managementPermissionCodes); err != nil {
				txDB.Rollback()
				return errors.ConnectStorageErrWrapper(err)
			}
		}

		return nil
	})
}

// GetUsersByOperationCode will return admin user if no qualified user is found, preventing the process from being stuck because no user can operate
func (s *Storage) GetUsersByOperationCode(instance *Instance, opCode ...int) (users []*User, err error) {
	names := []string{}
	err = s.db.Model(&User{}).Select("DISTINCT users.login_name").
		Joins("LEFT JOIN user_role ON users.id = user_role.user_id "+
			"LEFT JOIN user_group_users ON users.id = user_group_users.user_id "+
			"LEFT JOIN user_group_roles ON user_group_users.user_group_id = user_group_roles.role_id "+
			"LEFT JOIN role_operations ON ( user_role.role_id = role_operations.role_id OR user_group_roles.role_id = role_operations.role_id ) AND role_operations.deleted_at IS NULL "+
			"LEFT JOIN instance_role ON instance_role.role_id = role_operations.role_id ").
		Where("instance_role.instance_id = ?", instance.ID).
		Where("role_operations.op_code in (?)", opCode).
		Group("users.id").
		Pluck("login_name", &names).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}
	if len(names) == 0 {
		names = append(names, DefaultAdminUser)
	}
	return s.GetUsersByNames(names)
}

func (s *Storage) GetUserByID(id uint) (*User, bool, error) {
	u := &User{}
	err := s.db.Model(User{}).Where("id = ?", id).Find(u).Error
	if err == gorm.ErrRecordNotFound {
		return u, false, nil
	}
	return u, true, errors.New(errors.ConnectStorageError, err)
}
