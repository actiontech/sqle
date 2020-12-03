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
func (s *Storage) GetSqlWhitelist() ([]SqlWhitelist, error) {
	sqlWhitelist := []SqlWhitelist{}
	err := s.db.Find(&sqlWhitelist).Error
	return sqlWhitelist, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
