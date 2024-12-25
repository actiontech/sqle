package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
	"time"
)

type AuditRuleCategory struct {
	ID        uint      `json:"id" gorm:"primary_key" example:"1"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp(3)" example:"2018-10-21T16:40:23+08:00"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp(3) on update current_timestamp(3)" example:"2018-10-21T16:40:23+08:00"`
	Category  string    `json:"category" gorm:"not null;type:varchar(255)"`
	Tag       string    `json:"tag" gorm:"not null;type:varchar(255)"`
}

func (s *Storage) GetAllCategories() ([]*AuditRuleCategory, error) {
	var auditRuleCategories []*AuditRuleCategory
	err := s.db.Find(&auditRuleCategories).Error
	if err != nil {
		return nil, err
	}
	return auditRuleCategories, nil
}

func (s *Storage) GetAuditRuleCategoryByCategory(category string) ([]*AuditRuleCategory, error) {
	var auditRuleCategory []*AuditRuleCategory
	err := s.db.Model(AuditRuleCategory{}).Where("category = ?", category).Find(&auditRuleCategory).Error
	if err == gorm.ErrRecordNotFound {
		return auditRuleCategory, err
	}
	return auditRuleCategory, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditRuleCategoryByTagIn(tag []string) ([]*AuditRuleCategory, error) {
	var auditRuleCategory []*AuditRuleCategory
	err := s.db.Model(AuditRuleCategory{}).Where("tag in (?)", tag).Find(&auditRuleCategory).Error
	if err == gorm.ErrRecordNotFound {
		return auditRuleCategory, err
	}
	return auditRuleCategory, errors.New(errors.ConnectStorageError, err)
}

type RuleCategoryStatistic struct {
	Category string `json:"category" enums:"audit_accuracy,audit_purpose,operand,sql"`
	Tag      string `json:"tag"`
	Count    int    `json:"count"`
}

func (s *Storage) GetAuditRuleCategoryStatistics() ([]*RuleCategoryStatistic, error) {
	var auditRuleCategoryStatistics []*RuleCategoryStatistic
	err := s.db.Model(AuditRuleCategory{}).
		Joins("left join audit_rule_category_rels on audit_rule_categories.id = audit_rule_category_rels.category_id").
		Joins("left join rules on audit_rule_category_rels.rule_name = rules.name and audit_rule_category_rels.rule_db_type = rules.db_type").
		Select("audit_rule_categories.category, audit_rule_categories.tag, count(audit_rule_categories.tag) as count").
		Group("category, tag").
		Scan(&auditRuleCategoryStatistics).Error
	return auditRuleCategoryStatistics, err
}

func (s *Storage) GetCustomRuleCategoryStatistics() ([]*RuleCategoryStatistic, error) {
	var auditRuleCategoryStatistics []*RuleCategoryStatistic
	err := s.db.Model(AuditRuleCategory{}).
		Joins("left join custom_rule_category_rels on audit_rule_categories.id = custom_rule_category_rels.category_id").
		Joins("left join custom_rules on custom_rule_category_rels.custom_rule_id = custom_rules.rule_id").
		Select("audit_rule_categories.category, audit_rule_categories.tag, count(audit_rule_categories.tag) as count").
		Group("category, tag").
		Scan(&auditRuleCategoryStatistics).Error
	return auditRuleCategoryStatistics, err
}

type AuditRuleCategoryRel struct {
	CategoryId uint   `json:"category_id" gorm:"not null;type:bigint unsigned;primary_key;autoIncrement:false"`
	RuleName   string `json:"rule_name" gorm:"not null;type:varchar(255);primary_key"`
	RuleDBType string `json:"db_type" gorm:"not null;column:rule_db_type;type:varchar(255);primary_key"`
	Rule       *Rule  `json:"-" gorm:"foreignkey:Name,DBType;references:RuleName,RuleDBType"`
}

func (s *Storage) FirstAuditRuleCategoryRelByRule(ruleName string, ruleDbType string) (*AuditRuleCategoryRel, bool, error) {
	categoryRel := &AuditRuleCategoryRel{}
	err := s.db.Where("rule_name = ? AND rule_db_type = ?", ruleName, ruleDbType).First(categoryRel).Error
	if err == gorm.ErrRecordNotFound {
		return categoryRel, false, nil
	}
	return categoryRel, true, errors.New(errors.ConnectStorageError, err)
}

type CustomRuleCategoryRel struct {
	CategoryId   uint        `json:"category_id" gorm:"not null;type:bigint unsigned;primary_key"`
	CustomRuleId string      `json:"custom_rule_id" gorm:"not null;type:varchar(255);primary_key"`
	CustomRule   *CustomRule `json:"-" gorm:"foreignkey:RuleId;references:CustomRuleId"`
}

func (s *Storage) FirstCustomRuleCategoryRelByCustomRuleId(customRuleId string) (*CustomRuleCategoryRel, bool, error) {
	categoryRel := &CustomRuleCategoryRel{}
	err := s.db.Where("custom_rule_id = ?", customRuleId).First(categoryRel).Error
	if err == gorm.ErrRecordNotFound {
		return categoryRel, false, nil
	}
	return categoryRel, true, errors.New(errors.ConnectStorageError, err)
}
