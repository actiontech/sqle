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
	Stat  uint    `json:"stat" gorm:"comment:'0:正常,1:禁用'"`
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

func (s *Storage) CheckIfUserGroupExistByName(userGroupName string) (
	ug *UserGroup, isExist bool, err error) {
	ug = &UserGroup{
		Name: userGroupName,
	}
	err = s.db.First(ug).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, false, nil
	}
	return ug, true, errors.NewConnectStorageErrWrapper(err)
}
