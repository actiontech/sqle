//go:build enterprise
// +build enterprise

package model

import (
	"context"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/pkg/i18nPkg"
	"gorm.io/gorm"
)

func (s *Storage) CreateOrUpdateRuleKnowledgeContent(ctx context.Context, ruleName, dbType, content string) error {
	lang := locale.GetLangTagFromCtx(ctx)
	rule := Rule{
		Name:   ruleName,
		DBType: dbType,
	}
	if err := s.db.Preload("Knowledge").Find(&rule).Error; err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	if rule.Knowledge == nil || rule.Knowledge.I18nContent == nil {
		rule.Knowledge = &RuleKnowledge{I18nContent: i18nPkg.I18nStr{
			lang: content,
		}}
	} else {
		rule.Knowledge.I18nContent.SetStrInLang(lang, content)
	}
	return errors.New(errors.ConnectStorageError, s.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&rule).Error)
}

func (s *Storage) CreateOrUpdateCustomRuleKnowledgeContent(ctx context.Context, ruleName, dbType, content string) error {
	lang := locale.GetLangTagFromCtx(ctx)
	rule := CustomRule{}
	if err := s.db.Preload("Knowledge").Where("rule_id = ?", ruleName).Where("db_type = ?", dbType).Find(&rule).Error; err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	if rule.Knowledge == nil || rule.Knowledge.I18nContent == nil {
		rule.Knowledge = &RuleKnowledge{I18nContent: i18nPkg.I18nStr{
			lang: content,
		}}
	} else {
		rule.Knowledge.I18nContent.SetStrInLang(lang, content)
	}
	return errors.New(errors.ConnectStorageError, s.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&rule).Error)
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
