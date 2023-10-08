//go:build enterprise
// +build enterprise

package model

import "github.com/actiontech/sqle/sqle/errors"

func (s *Storage) CreateOrUpdateRuleKnowledgeContent(ruleName, dbType, content string) error {
	rule := Rule{
		Name:   ruleName,
		DBType: dbType,
	}
	if err := s.db.Preload("Knowledge").Find(&rule).Error; err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	if rule.Knowledge == nil {
		rule.Knowledge = &RuleKnowledge{Content: content}
	} else {
		rule.Knowledge.Content = content
	}
	return errors.New(errors.ConnectStorageError, s.db.Save(&rule).Error)
}

func (s *Storage) CreateOrUpdateCustomRuleKnowledgeContent(ruleName, dbType, content string) error {
	rule := CustomRule{}
	if err := s.db.Preload("Knowledge").Where("rule_id = ?", ruleName).Where("db_type = ?", dbType).Find(&rule).Error; err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	if rule.Knowledge == nil {
		rule.Knowledge = &RuleKnowledge{Content: content}
	} else {
		rule.Knowledge.Content = content
	}
	return errors.New(errors.ConnectStorageError, s.db.Save(&rule).Error)
}

func (s *Storage) GetRuleWithKnowledge(ruleName, dbType string) (*Rule, error) {
	rule := Rule{
		Name:   ruleName,
		DBType: dbType,
	}
	if err := s.db.Preload("Knowledge").Find(&rule).Error; err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	return &rule, nil
}

func (s *Storage) GetCustomRuleWithKnowledge(ruleName, dbType string) (*CustomRule, error) {
	rule := CustomRule{}
	if err := s.db.Preload("Knowledge").Where("rule_id = ?", ruleName).Where("db_type = ?", dbType).Find(&rule).Error; err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	return &rule, nil
}
