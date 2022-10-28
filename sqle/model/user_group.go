package model

import (
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
)

type UserGroup struct {
	Model
	Name  string  `json:"name" gorm:"index"`
	Desc  string  `json:"desc" gorm:"column:description"`
	Users []*User `gorm:"many2many:user_group_users"`
	Stat  uint    `json:"stat" gorm:"comment:'0:active,1:disabled'"`
	// todo issue960 remove Roles
	Roles []*Role `gorm:"many2many:user_group_roles"`
}

func (ug *UserGroup) TableName() string {
	return "user_groups"
}

func (ug *UserGroup) SetStat(stat int) {
	ug.Stat = uint(stat)
}

func (ug *UserGroup) IsDisabled() bool {
	return ug.Stat == Disabled
}

func (s *Storage) GetUserGroupByName(name string) (
	userGroup *UserGroup, isExist bool, err error) {
	userGroup = &UserGroup{}

	err = s.db.Where("name = ?", name).First(userGroup).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, false, nil
	}

	return userGroup, true, err
}

// NOTE: parameter: us([]*Users) and rs([]*Role) need to be distinguished as nil or zero length slice.
func (s *Storage) SaveUserGroupAndAssociations(
	ug *UserGroup, us []*User, rs []*Role) (err error) {

	return s.Tx(func(txDB *gorm.DB) error {
		if err := txDB.Save(ug).Error; err != nil {
			return err
		}

		// save user group users
		if us != nil {
			if err := txDB.Model(ug).Association("Users").Replace(us).Error; err != nil {
				return err
			}
		}

		// save user group roles
		if rs != nil {
			if err := txDB.Model(ug).Association("Roles").Replace(rs).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Storage) GetAllUserGroupTip() ([]*UserGroup, error) {
	userGroups := []*UserGroup{}
	err := s.db.Select("name").Where("stat = 0").Find(&userGroups).Error
	return userGroups, errors.New(errors.ConnectStorageError, err)
}

var userGroupsQueryTpl = `SELECT 
user_groups.id, 
user_groups.name, 
user_groups.description,
user_groups.stat,
GROUP_CONCAT(DISTINCT COALESCE(users.login_name,'')) AS user_names,
GROUP_CONCAT(DISTINCT COALESCE(roles.name,'')) AS role_names
FROM user_groups
LEFT JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
LEFT JOIN users ON user_group_users.user_id = users.id AND users.deleted_at IS NULL
LEFT JOIN user_group_roles ON user_groups.id = user_group_roles.user_group_id
LEFT JOIN roles ON user_group_roles.role_id = roles.id AND roles.deleted_at IS NULL
WHERE 
user_groups.id IN (SELECT DISTINCT(user_groups.id)

{{- template "body" . -}}
)
GROUP BY user_groups.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var userGroupCountTpl = `SELECT COUNT(DISTINCT user_groups.id)

{{- template "body" . -}}
`
var userGroupsQueryBodyTpl = `
{{ define "body" }}
FROM user_groups
LEFT JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
LEFT JOIN users ON user_group_users.user_id = users.id AND users.deleted_at IS NULL
LEFT JOIN user_group_roles ON user_groups.id = user_group_roles.user_group_id
LEFT JOIN roles ON user_group_roles.role_id = roles.id AND roles.deleted_at IS NULL
WHERE
user_groups.deleted_at IS NULL

{{- if .filter_user_group_name }}
AND user_groups.name = :filter_user_group_name
{{- end -}}

{{- end }}
`

type UserGroupDetail struct {
	Id        int
	Name      string  `json:"name"`
	Desc      string  `json:"description"`
	Stat      uint    `json:"stat"`
	UserNames RowList `json:"user_names"`
	RoleNames RowList `json:"role_names"`
}

func (ugd *UserGroupDetail) IsDisabled() bool {
	return ugd.Stat == Disabled
}

func (s *Storage) GetUserGroupsByReq(data map[string]interface{}) (
	results []*UserGroupDetail, count uint64, err error) {

	err = s.getListResult(userGroupsQueryBodyTpl, userGroupsQueryTpl, data, &results)
	if err != nil {
		return nil, 0, err
	}

	count, err = s.getCountResult(userGroupsQueryBodyTpl, userGroupCountTpl, data)
	return results, count, err
}

func (s *Storage) GetUserGroupsByNames(names []string) (ugs []*UserGroup, err error) {
	ugs = []*UserGroup{}
	err = s.db.Where("name IN (?)", names).Find(&ugs).Error
	return ugs, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetAndCheckUserGroupExist(userGroupNames []string) (ugs []*UserGroup, err error) {
	ugs, err = s.GetUserGroupsByNames(userGroupNames)
	if err != nil {
		return nil, err
	}

	existUserGroupNames := map[string]struct{}{}
	{
		for i := range ugs {
			existUserGroupNames[ugs[i].Name] = struct{}{}
		}
	}

	notExistUserGroupNames := []string{}
	{
		for i := range userGroupNames {
			if _, ok := existUserGroupNames[userGroupNames[i]]; !ok {
				notExistUserGroupNames = append(notExistUserGroupNames, userGroupNames[i])
			}
		}
	}

	if len(notExistUserGroupNames) > 0 {
		return ugs, errors.NewDataNotExistErr("user group <%v> not exist",
			strings.Join(notExistUserGroupNames, ", "))
	}

	return ugs, nil
}
