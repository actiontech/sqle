package model

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/jinzhu/gorm"
)

type Member struct {
	Model
	UserID    uint `gorm:"index;not null"`
	ProjectID uint `gorm:"index;not null"`
	IsManager bool `gorm:"not null"`
}

func (s *Storage) AddMember(userName, projectName string, isManager bool, bindRole []BindRole) error {
	user, exist, err := s.GetUserByName(userName)
	if err != nil {
		return errors.ConnectStorageErrWrapper(err)
	}
	if !exist {
		return errors.ConnectStorageErrWrapper(fmt.Errorf("user not exist"))
	}

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return errors.ConnectStorageErrWrapper(err)
	}
	if !exist {
		return errors.ConnectStorageErrWrapper(fmt.Errorf("project not exist"))
	}

	return errors.New(errors.ConnectStorageError, s.db.Transaction(func(tx *gorm.DB) error {

		if err = tx.Save(&Member{
			UserID:    user.ID,
			ProjectID: project.ID,
			IsManager: isManager,
		}).Error; err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}

		err = s.updateUserRoles(tx, user, projectName, bindRole)
		if err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}
		return nil
	}))
}

func (s *Storage) UpdateMemberByName(userName, projectName string, attrs ...interface{}) error {
	// JOIN表的方式更新会出现重复表字段, 这将导致入参 attrs 的 key 也要以 表名.字段名 的方式传入, 先查 userID 再更新可以降低后续开发心智负担
	member, _, err := s.GetMemberByName(userName, projectName)
	if err != nil {
		return err
	}

	err = s.db.Table("members").
		Where("user_id = ?", member.UserID).
		Where("project_id = ?", member.ProjectID).
		Update(attrs...).Error

	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetMemberByName(userName, projectName string) (*Member, bool, error) {
	m := &Member{}

	err := s.db.Joins("JOIN users ON users.id = members.user_id").
		Joins("JOIN projects ON projects.id = members.project_id").
		Where("users.login_name = ?", userName).
		Where("projects.name = ?", projectName).
		First(m).Error

	if gorm.ErrRecordNotFound == err {
		return &Member{}, false, nil
	}

	return m, true, err
}

func (s *Storage) RemoveMember(userName, projectName string) error {
	sql := `
id = (

SELECT id 
FROM members as m
JOIN users ON users.id = members.user_id
JOIN projects ON projects.id = members.project_id
WHERE 
	users.login_name = ?
AND 
	projects.name = ?
LIMIT 1

)
`

	return errors.ConnectStorageErrWrapper(s.db.Where(sql, userName, projectName).Delete(&Member{}).Error)
}
