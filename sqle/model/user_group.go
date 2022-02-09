package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
)

type UserGroup struct {
	Model
	Name  string  `json:"name" gorm:"index"`
	Desc  string  `json:"desc" gorm:"column:description"`
	Users []*User `gorm:"many2many:user_group_users"`
	Stat  uint    `json:"stat" gorm:"comment:'0:active,1:disabled'"`
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
	err := s.db.Select("name").Find(&userGroups).Error
	return userGroups, errors.New(errors.ConnectStorageError, err)
}

type UserGroupDetail struct {
	Id        int
	Name      string `json:"user_group_name"`
	Desc      string
	Stat      uint    `json:"stat"`
	UserNames RowList `json:"user_names"`
	RoleNames RowList `json:"role_names"`
}

func (ugd *UserGroupDetail) IsDisabled() bool {
	return ugd.Stat == Disabled
}

// TODO: implementation
func (s *Storage) GetUserGroupsByReq(data map[string]interface{}) (
	userGroups []*UserGroupDetail, count uint64, err error) {
	return
}
