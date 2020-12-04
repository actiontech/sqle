package model

import (
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"github.com/jinzhu/gorm"
)

type SqlWhitelist struct {
	Model
	Value string `json:"value" gorm:"not null"`
	Desc  string `json:"desc"`
}

func (s SqlWhitelist) TableName() string {
	return "sql_whitelist"
}
func (s *Storage) GetSqlWhitelistItemById(sqlWhiteId string) (*SqlWhitelist, bool, error) {
	sqlWhitelist := &SqlWhitelist{}
	err := s.db.Table("sql_whitelist").Where("id = ?", sqlWhiteId).First(sqlWhitelist).Error
	if err == gorm.ErrRecordNotFound {
		return sqlWhitelist, false, nil
	}
	return sqlWhitelist, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
func (s *Storage) GetSqlWhitelist(pageIndex, pageSize int) ([]SqlWhitelist, uint32, error) {
	var count uint32
	sqlWhitelist := []SqlWhitelist{}
	if pageSize == 0 {
		err := s.db.Order("id desc").Find(&sqlWhitelist).Count(&count).Error
		return sqlWhitelist, count, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	err := s.db.Model(&SqlWhitelist{}).Count(&count).Error
	if err != nil {
		return sqlWhitelist, 0, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	err = s.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&sqlWhitelist).Error
	return sqlWhitelist, count, errors.New(errors.CONNECT_STORAGE_ERROR, err)

}
