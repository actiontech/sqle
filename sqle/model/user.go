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
	err := s.db.Preload("UserGroups").
		Where("login_name = ?", name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
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
SELECT DISTINCT users.login_name 
FROM users
LEFT JOIN project_user on project_user.user_id = users.id
LEFT JOIN projects on project_user.project_id = projects.id
LEFT JOIN user_group_users on users.id = user_group_users.user_id 
LEFT JOIN project_user_group on user_group_users.user_group_id = project_user_group.user_group_id
LEFT JOIN projects as p on project_user_group.project_id = p.id
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

func (s *Storage) GetUserByID(id uint) (*User, bool, error) {
	u := &User{}
	err := s.db.Model(User{}).Where("id = ?", id).Find(u).Error
	if err == gorm.ErrRecordNotFound {
		return u, false, nil
	}
	return u, true, errors.New(errors.ConnectStorageError, err)
}
