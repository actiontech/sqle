//go:build enterprise
// +build enterprise

package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
)

// 源项目id与sqle项目名称映射表
type IdProjectNamePair struct {
	ID              int64  `json:"id" gorm:"primary_key,autoIncrement"`     //自增主键
	OriginProjectID int64  `json:"origin_project_id" gorm:"not null;index"` //源系统的项目id
	SqleProjectName string `json:"sqle_project_name" gorm:"not null"`       //sqle项目名称
}

func (s *Storage) GetIdProjectNamePairsByIds(ids []int64) ([]IdProjectNamePair, bool, error) {
	var IdProjectNamePairs []IdProjectNamePair
	err := s.db.
		Model(&IdProjectNamePair{}).
		Where("origin_project_id IN (?)", ids).
		First(&IdProjectNamePairs).
		Error
	if err == gorm.ErrRecordNotFound {
		return IdProjectNamePairs, false, nil
	}
	return IdProjectNamePairs, true, errors.New(errors.ConnectStorageError, err)
}
