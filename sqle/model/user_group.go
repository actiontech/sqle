package model

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/utils"
	"github.com/jackc/pgx/v4"
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

func (s *Storage) CheckIfUserGroupExistByName(userGroupName string) (isExist bool, err error) {
	query := `SELECT 1 FROM %v WHERE name = ? LIMIT 1`
	query = fmt.Sprintf(query, (&UserGroup{}).TableName())

	cnt := 0
	if err = s.db.Raw(query, userGroupName).Count(&cnt).Error; err != nil &&
		utils.IsErrorEqual(err, pgx.ErrNoRows) {
		return false, err
	}

	return cnt == 1, nil
}
