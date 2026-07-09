package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/errors"

	"gorm.io/gorm"
)

type SQLRuleException struct {
	Model
	ProjectId      ProjectUID `json:"project_id" gorm:"index;not null"`
	InstanceID     uint64     `json:"instance_id" gorm:"not null;index"`
	SQLFingerprint string     `json:"sql_fingerprint" gorm:"type:varchar(512);not null"`
	RuleName       string     `json:"rule_name" gorm:"type:varchar(255);not null"`
	RuleDesc       string     `json:"rule_desc" gorm:"type:varchar(1024)"`
	RuleLevel      string     `json:"rule_level" gorm:"type:varchar(64)"`
	Reason         string     `json:"reason" gorm:"type:varchar(255);not null"`
	CreatedBy      string     `json:"created_by" gorm:"type:varchar(128)"`
	UniqueKey      string     `json:"-" gorm:"type:char(64);not null;uniqueIndex:uniq_sql_rule_exception_effective"`
}

type SQLRuleExceptionListFilter struct {
	ProjectID        ProjectUID
	InstanceID       *uint64
	RuleName         *string
	CreatedBy        *string
	CreatedTimeFrom  *string
	CreatedTimeTo    *string
	SQLFingerprint   *string
	FuzzySearchValue *string
}

type SQLRuleExceptionListItem struct {
	SQLRuleException
	InstanceName string `json:"instance_name"`
}

const SQLRuleExceptionMissingFingerprintMessage = "SQL fingerprint is missing or cannot be generated; cannot add a SQL rule exception"

func (s SQLRuleException) TableName() string {
	return "sql_rule_exception"
}

func (s *SQLRuleException) BeforeSave(tx *gorm.DB) error {
	s.SQLFingerprint = strings.TrimSpace(s.SQLFingerprint)
	s.RuleName = strings.TrimSpace(s.RuleName)
	s.Reason = strings.TrimSpace(s.Reason)
	if s.SQLFingerprint == "" {
		return fmt.Errorf("sql fingerprint is required")
	}
	if s.RuleName == "" {
		return fmt.Errorf("rule name is required")
	}
	if s.Reason == "" {
		return fmt.Errorf("reason is required")
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	s.UniqueKey = BuildSQLRuleExceptionUniqueKey(s.ProjectId, s.InstanceID, s.SQLFingerprint, s.RuleName)
	return nil
}

func BuildSQLRuleExceptionUniqueKey(projectID ProjectUID, instanceID uint64, sqlFingerprint, ruleName string) string {
	rawValue := fmt.Sprintf("%s\x00%d\x00%s\x00%s", projectID, instanceID, strings.TrimSpace(sqlFingerprint), strings.TrimSpace(ruleName))
	sum := sha256.Sum256([]byte(rawValue))
	return hex.EncodeToString(sum[:])
}

func BuildDeletedSQLRuleExceptionUniqueKey(id uint, uniqueKey string) string {
	rawValue := fmt.Sprintf("deleted\x00%d\x00%s", id, uniqueKey)
	sum := sha256.Sum256([]byte(rawValue))
	return hex.EncodeToString(sum[:])
}

func (s *Storage) GetEffectiveSQLRuleException(projectID ProjectUID, instanceID uint64, sqlFingerprint, ruleName string) (*SQLRuleException, bool, error) {
	sqlRuleException := &SQLRuleException{}
	err := s.db.Where("project_id = ? AND instance_id = ? AND sql_fingerprint = ? AND rule_name = ?", projectID, instanceID, strings.TrimSpace(sqlFingerprint), strings.TrimSpace(ruleName)).
		First(sqlRuleException).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return sqlRuleException, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetEffectiveSQLRuleExceptions(projectID ProjectUID, instanceID uint64, sqlFingerprint string, ruleNames []string) (map[string]*SQLRuleException, error) {
	return s.GetEffectiveSQLRuleExceptionsByFingerprints(projectID, instanceID, []string{sqlFingerprint}, ruleNames)
}

func (s *Storage) GetEffectiveSQLRuleExceptionsByFingerprints(projectID ProjectUID, instanceID uint64, sqlFingerprints []string, ruleNames []string) (map[string]*SQLRuleException, error) {
	trimmedSQLFingerprints := make([]string, 0, len(sqlFingerprints))
	seenSQLFingerprints := map[string]struct{}{}
	for _, sqlFingerprint := range sqlFingerprints {
		trimmedSQLFingerprint := strings.TrimSpace(sqlFingerprint)
		if trimmedSQLFingerprint == "" {
			continue
		}
		if _, ok := seenSQLFingerprints[trimmedSQLFingerprint]; ok {
			continue
		}
		seenSQLFingerprints[trimmedSQLFingerprint] = struct{}{}
		trimmedSQLFingerprints = append(trimmedSQLFingerprints, trimmedSQLFingerprint)
	}
	if len(trimmedSQLFingerprints) == 0 {
		return map[string]*SQLRuleException{}, nil
	}

	trimmedRuleNames := make([]string, 0, len(ruleNames))
	seenRuleNames := map[string]struct{}{}
	for _, ruleName := range ruleNames {
		trimmedRuleName := strings.TrimSpace(ruleName)
		if trimmedRuleName == "" {
			continue
		}
		if _, ok := seenRuleNames[trimmedRuleName]; ok {
			continue
		}
		seenRuleNames[trimmedRuleName] = struct{}{}
		trimmedRuleNames = append(trimmedRuleNames, trimmedRuleName)
	}
	if len(trimmedRuleNames) == 0 {
		return map[string]*SQLRuleException{}, nil
	}

	sqlRuleExceptions := []*SQLRuleException{}
	err := s.db.Where("project_id = ? AND sql_fingerprint IN (?) AND rule_name IN (?)", projectID, trimmedSQLFingerprints, trimmedRuleNames).
		Find(&sqlRuleExceptions).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}

	ret := make(map[string]*SQLRuleException, len(sqlRuleExceptions))
	for _, sqlRuleException := range sqlRuleExceptions {
		if sqlRuleException.InstanceID == instanceID {
			ret[sqlRuleException.RuleName] = sqlRuleException
		}
	}
	if len(ret) > 0 {
		return ret, nil
	}
	for _, sqlRuleException := range sqlRuleExceptions {
		ret[sqlRuleException.RuleName] = sqlRuleException
	}
	return ret, nil
}

func (s *Storage) CreateSQLRuleExceptionIfNotExist(sqlRuleException *SQLRuleException) (*SQLRuleException, bool, error) {
	if strings.TrimSpace(sqlRuleException.SQLFingerprint) == "" {
		return nil, false, errors.NewDataInvalidErr(SQLRuleExceptionMissingFingerprintMessage)
	}

	existedSQLRuleException, exist, err := s.GetEffectiveSQLRuleException(sqlRuleException.ProjectId, sqlRuleException.InstanceID, sqlRuleException.SQLFingerprint, sqlRuleException.RuleName)
	if err != nil {
		return nil, false, err
	}
	if exist {
		return existedSQLRuleException, false, errors.New(errors.DataExist, fmt.Errorf("sql rule exception already exists"))
	}

	if err := s.db.Create(sqlRuleException).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			existedSQLRuleException, exist, getErr := s.GetEffectiveSQLRuleException(sqlRuleException.ProjectId, sqlRuleException.InstanceID, sqlRuleException.SQLFingerprint, sqlRuleException.RuleName)
			if getErr != nil {
				return nil, false, getErr
			}
			if exist {
				return existedSQLRuleException, false, errors.New(errors.DataExist, fmt.Errorf("sql rule exception already exists"))
			}
		}
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	return sqlRuleException, true, nil
}

func (s *Storage) GetSQLRuleExceptionByIDAndProjectID(id string, projectID ProjectUID) (*SQLRuleException, bool, error) {
	sqlRuleException := &SQLRuleException{}
	err := s.db.Where("id = ? AND project_id = ?", id, projectID).First(sqlRuleException).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return sqlRuleException, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) DeleteSQLRuleException(sqlRuleException *SQLRuleException) error {
	if sqlRuleException == nil {
		return errors.New(errors.DataNotExist, fmt.Errorf("sql rule exception is not exist"))
	}
	return errors.New(errors.ConnectStorageError, s.db.Transaction(func(tx *gorm.DB) error {
		deletedUniqueKey := BuildDeletedSQLRuleExceptionUniqueKey(sqlRuleException.ID, sqlRuleException.UniqueKey)
		if err := tx.Table(sqlRuleException.TableName()).Where("id = ? AND deleted_at IS NULL", sqlRuleException.ID).Update("unique_key", deletedUniqueKey).Error; err != nil {
			return err
		}
		return tx.Delete(sqlRuleException).Error
	}))
}

func (s *Storage) GetSQLRuleExceptions(pageIndex, pageSize uint32, filter SQLRuleExceptionListFilter) ([]*SQLRuleExceptionListItem, int64, error) {
	var count int64
	sqlRuleExceptions := []*SQLRuleExceptionListItem{}
	query := s.db.Table("sql_rule_exception").
		Select("sql_rule_exception.*, instances.name AS instance_name").
		Joins("LEFT JOIN instances ON instances.id = sql_rule_exception.instance_id").
		Where("sql_rule_exception.project_id = ?", filter.ProjectID).
		Where("sql_rule_exception.deleted_at IS NULL")

	if filter.InstanceID != nil {
		query = query.Where("sql_rule_exception.instance_id = ?", *filter.InstanceID)
	}
	if filter.RuleName != nil && strings.TrimSpace(*filter.RuleName) != "" {
		query = query.Where("sql_rule_exception.rule_name LIKE ?", "%"+strings.TrimSpace(*filter.RuleName)+"%")
	}
	if filter.CreatedBy != nil && strings.TrimSpace(*filter.CreatedBy) != "" {
		query = query.Where("sql_rule_exception.created_by = ?", strings.TrimSpace(*filter.CreatedBy))
	}
	if filter.CreatedTimeFrom != nil && strings.TrimSpace(*filter.CreatedTimeFrom) != "" {
		query = query.Where("sql_rule_exception.created_at > ?", strings.TrimSpace(*filter.CreatedTimeFrom))
	}
	if filter.CreatedTimeTo != nil && strings.TrimSpace(*filter.CreatedTimeTo) != "" {
		query = query.Where("sql_rule_exception.created_at < ?", strings.TrimSpace(*filter.CreatedTimeTo))
	}
	if filter.SQLFingerprint != nil && strings.TrimSpace(*filter.SQLFingerprint) != "" {
		query = query.Where("sql_rule_exception.sql_fingerprint LIKE ?", "%"+strings.TrimSpace(*filter.SQLFingerprint)+"%")
	}
	if filter.FuzzySearchValue != nil && strings.TrimSpace(*filter.FuzzySearchValue) != "" {
		fuzzySearchValue := "%" + strings.TrimSpace(*filter.FuzzySearchValue) + "%"
		query = query.Where("sql_rule_exception.sql_fingerprint LIKE ? OR sql_rule_exception.rule_name LIKE ? OR sql_rule_exception.rule_desc LIKE ? OR sql_rule_exception.reason LIKE ? OR sql_rule_exception.created_by LIKE ?", fuzzySearchValue, fuzzySearchValue, fuzzySearchValue, fuzzySearchValue, fuzzySearchValue)
	}

	if pageSize == 0 {
		err := query.Order("sql_rule_exception.id desc").Find(&sqlRuleExceptions).Count(&count).Error
		return sqlRuleExceptions, count, errors.New(errors.ConnectStorageError, err)
	}
	err := query.Count(&count).Error
	if err != nil {
		return sqlRuleExceptions, 0, errors.New(errors.ConnectStorageError, err)
	}
	err = query.Offset(int((pageIndex - 1) * pageSize)).Limit(int(pageSize)).Order("sql_rule_exception.id desc").Find(&sqlRuleExceptions).Error
	return sqlRuleExceptions, count, errors.New(errors.ConnectStorageError, err)
}
