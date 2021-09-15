package model

import "github.com/actiontech/sqle/sqle/errors"

type RuleTemplateRule struct {
	RuleTemplateId uint   `json:"rule_template_id" gorm:"primary_key;auto_increment:false;"`
	RuleName       string `json:"name" gorm:"primary_key;"`
	RuleLevel      string `json:"level" gorm:"column:level;"`
	RuleValue      string `json:"value" gorm:"column:value;" `
	RuleDBType     string `json:"rule_db_type" gorm:"column:db_type; not null; default:'mysql'"`

	Rule Rule `json:"-" gorm:"foreignkey:Name,DBType;association_foreignkey:RuleName,RuleDBType"`
}

func (rtr *RuleTemplateRule) TableName() string {
	return "rule_template_rule"
}

func NewRuleTemplateRule(t *RuleTemplate, r *Rule) RuleTemplateRule {
	return RuleTemplateRule{
		RuleTemplateId: t.ID,
		RuleName:       r.Name,
		RuleLevel:      r.Level,
		RuleValue:      r.Value,
		RuleDBType:     r.DBType,
	}
}

func (s *Storage) GetRuleTemplatesByInstance(inst *Instance) ([]RuleTemplate, error) {
	var associationRT []RuleTemplate
	err := s.db.Model(inst).Association("RuleTemplates").Find(&associationRT).Error
	return associationRT, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRulesFromRuleTemplateByName(name string) ([]Rule, error) {
	tpl, exist, err := s.GetRuleTemplateDetailByName(name)
	if !exist {
		return nil, errors.New(errors.DataNotExist, err)
	}
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}

	rules := make([]Rule, 0, len(tpl.RuleList))
	for _, r := range tpl.RuleList {
		if r.RuleLevel != "" {
			r.Rule.Level = r.RuleLevel
		}
		if r.RuleValue != "" {
			r.Rule.Value = r.RuleValue
		}
		rules = append(rules, r.Rule)
	}

	return rules, errors.New(errors.ConnectStorageError, err)
}
