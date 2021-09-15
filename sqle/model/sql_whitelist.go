package model

import (
	"strings"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/jinzhu/gorm"
)

const (
	SQLWhitelistExactMatch = "exact_match"
	SQLWhitelistFPMatch    = "fp_match"
)

type SqlWhitelist struct {
	Model
	// Value store SQL text.
	Value            string `json:"value" gorm:"not null;type:text"`
	CapitalizedValue string `json:"-" gorm:"-"`
	Desc             string `json:"desc"`
	// MessageDigest deprecated after 1.1.0, keep it for compatibility.
	MessageDigest string `json:"message_digest" gorm:"type:char(32) not null comment 'md5 data';" `
	MatchType     string `json:"match_type" gorm:"default:\"exact_match\""`
}

// BeforeSave is a hook implement gorm model before exec create
func (s *SqlWhitelist) BeforeSave() error {
	s.MessageDigest = "deprecated after 1.1.0"
	return nil
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db
func (s *SqlWhitelist) AfterFind() error {
	s.CapitalizedValue = strings.ToUpper(s.Value)
	return nil
}

func (s SqlWhitelist) TableName() string {
	return "sql_whitelist"
}

func (s *Storage) GetSqlWhitelistById(sqlWhiteId string) (*SqlWhitelist, bool, error) {
	sqlWhitelist := &SqlWhitelist{}
	err := s.db.Table("sql_whitelist").Where("id = ?", sqlWhiteId).First(sqlWhitelist).Error
	if err == gorm.ErrRecordNotFound {
		return sqlWhitelist, false, nil
	}
	return sqlWhitelist, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetSqlWhitelist(pageIndex, pageSize uint32) ([]SqlWhitelist, uint32, error) {
	var count uint32
	sqlWhitelist := []SqlWhitelist{}
	if pageSize == 0 {
		err := s.db.Order("id desc").Find(&sqlWhitelist).Count(&count).Error
		return sqlWhitelist, count, errors.New(errors.ConnectStorageError, err)
	}
	err := s.db.Model(&SqlWhitelist{}).Count(&count).Error
	if err != nil {
		return sqlWhitelist, 0, errors.New(errors.ConnectStorageError, err)
	}
	err = s.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&sqlWhitelist).Error
	return sqlWhitelist, count, errors.New(errors.ConnectStorageError, err)
}
