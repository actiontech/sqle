package model

import (
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/errors"

	"gorm.io/gorm"
)

const (
	SQLWhitelistExactMatch = "exact_match"
	SQLWhitelistFPMatch    = "fp_match"
)

type SqlWhitelist struct {
	Model
	ProjectId ProjectUID `gorm:"index; not null"`
	// Value store SQL text.
	Value            string `json:"value" gorm:"not null;type:text"`
	CapitalizedValue string `json:"-" gorm:"-"`
	Desc             string `json:"desc" gorm:"type:varchar(255)"`
	// MessageDigest deprecated after 1.1.0, keep it for compatibility.
	MessageDigest   string     `json:"message_digest" gorm:"type:char(32) not null comment 'md5 data';" `
	MatchType       string     `json:"match_type" gorm:"default:\"exact_match\""`
	MatchedCount    int        `json:"matched_count" gorm:"default:0"`
	LastMatchedTime *time.Time `json:"last_matched_time"`
}

// BeforeSave is a hook implement gorm model before exec create
func (s *SqlWhitelist) BeforeSave(tx *gorm.DB) error {
	tx.Statement.SetColumn("MessageDigest", "deprecated after 1.1.0")
	return nil
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db
func (s *SqlWhitelist) AfterFind(tx *gorm.DB) error {
	s.CapitalizedValue = strings.ToUpper(s.Value)
	return nil
}

func (s SqlWhitelist) TableName() string {
	return "sql_whitelist"
}

// func (s *Storage) GetSqlWhitelistByIdAndProjectName(sqlWhiteId, projectName string) (*SqlWhitelist, bool, error) {
// 	sqlWhitelist := &SqlWhitelist{}
// 	err := s.db.Table("sql_whitelist").
// 		Joins("LEFT JOIN projects ON projects.id = sql_whitelist.project_id").
// 		Where("sql_whitelist.id = ?", sqlWhiteId).
// 		Where("projects.name = ?", projectName).
// 		First(sqlWhitelist).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return sqlWhitelist, false, nil
// 	}
// 	return sqlWhitelist, true, errors.New(errors.ConnectStorageError, err)
// }

func (s *Storage) GetSqlWhitelistByIdAndProjectUID(sqlWhiteId string, projectUID ProjectUID) (*SqlWhitelist, bool, error) {
	sqlWhitelist := &SqlWhitelist{}
	err := s.db.Table("sql_whitelist").
		Where("sql_whitelist.id = ?", sqlWhiteId).
		Where("project_id = ?", projectUID).
		First(sqlWhitelist).Error
	if err == gorm.ErrRecordNotFound {
		return sqlWhitelist, false, nil
	}
	return sqlWhitelist, true, errors.New(errors.ConnectStorageError, err)
}

// func (s *Storage) GetSqlWhitelistByProjectName(pageIndex, pageSize uint32, projectName string) ([]SqlWhitelist, uint32, error) {
// 	var count uint32
// 	sqlWhitelist := []SqlWhitelist{}
// 	query := s.db.Table("sql_whitelist").
// 		Joins("LEFT JOIN projects ON projects.id = sql_whitelist.project_id").
// 		Where("projects.name = ?", projectName)
// 	if pageSize == 0 {
// 		err := query.Order("id desc").Find(&sqlWhitelist).Count(&count).Error
// 		return sqlWhitelist, count, errors.New(errors.ConnectStorageError, err)
// 	}
// 	err := query.Count(&count).Error
// 	if err != nil {
// 		return sqlWhitelist, 0, errors.New(errors.ConnectStorageError, err)
// 	}
// 	err = query.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&sqlWhitelist).Error
// 	return sqlWhitelist, count, errors.New(errors.ConnectStorageError, err)
// }

func (s *Storage) GetSqlWhitelistByProjectUID(pageIndex, pageSize uint32, projectUID ProjectUID, fuzzyValue, matchType *string) ([]SqlWhitelist, int64, error) {
	var count int64
	sqlWhitelist := []SqlWhitelist{}
	query := s.db.Table("sql_whitelist").
		Where("project_id = ?", projectUID).Where("deleted_at IS NULL")

	if fuzzyValue != nil {
		query = query.Where("value LIKE ?", "%"+*fuzzyValue+"%")
	}

	if matchType != nil {
		query = query.Where("match_type = ?", *matchType)
	}

	if pageSize == 0 {
		err := query.Order("id desc").Find(&sqlWhitelist).Count(&count).Error
		return sqlWhitelist, count, errors.New(errors.ConnectStorageError, err)
	}
	err := query.Count(&count).Error
	if err != nil {
		return sqlWhitelist, 0, errors.New(errors.ConnectStorageError, err)
	}
	err = query.Offset(int((pageIndex - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&sqlWhitelist).Error
	return sqlWhitelist, count, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetSqlWhitelistByProjectId(projectId string) ([]SqlWhitelist, error) {
	sqlWhitelist := []SqlWhitelist{}
	err := s.db.Table("sql_whitelist").
		Where("sql_whitelist.project_id = ?", projectId).
		Find(&sqlWhitelist).Error
	return sqlWhitelist, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateSqlWhitelistMatchedInfo(id uint, count int, lastMatchedTime time.Time) error {
	m := map[string]interface{}{
		"matched_count":     gorm.Expr("matched_count + ?", count),
		"last_matched_time": lastMatchedTime,
	}

	err := s.db.Model(&SqlWhitelist{}).Where("sql_whitelist.id = ?", id).UpdateColumns(m).Error
	return errors.New(errors.ConnectStorageError, err)
}

// func (s *Storage) GetSqlWhitelistTotalByProjectName(projectName string) (uint64, error) {
// 	var count uint64
// 	err := s.db.
// 		Table("sql_whitelist").
// 		Joins("LEFT JOIN projects ON sql_whitelist.project_id = projects.id").
// 		Where("projects.name = ?", projectName).
// 		Where("sql_whitelist.deleted_at IS NULL").
// 		Count(&count).
// 		Error
// 	return count, errors.ConnectStorageErrWrapper(err)
// }
