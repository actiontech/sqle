package model

import (
	"fmt"
	"strconv"
	"strings"

	"actiontech.cloud/sqle/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
)

const (
	RULE_LEVEL_NORMAL = "normal"
	RULE_LEVEL_NOTICE = "notice"
	RULE_LEVEL_WARN   = "warn"
	RULE_LEVEL_ERROR  = "error"
)

var RuleLevelMap = map[string]int{
	RULE_LEVEL_NORMAL: 0,
	RULE_LEVEL_NOTICE: 1,
	RULE_LEVEL_WARN:   2,
	RULE_LEVEL_ERROR:  3,
}

type RuleTemplate struct {
	Model
	Name      string             `json:"name"`
	Desc      string             `json:"desc"`
	Instances []Instance         `json:"instance_list" gorm:"many2many:instance_rule_template"`
	RuleList  []RuleTemplateRule `json:"rule_list" gorm:"foreignkey:rule_template_id;association_foreignkey:id"`
}

type Rule struct {
	Name  string `json:"name" gorm:"primary_key"`
	Desc  string `json:"desc"`
	Value string `json:"value"`
	Level string `json:"level" example:"error"` // notice, warn, error
	Typ   string `gorm:"column:type; not null" json:"type"`
}

func (r Rule) TableName() string {
	return "rules"
}

func (r *Rule) GetValue() string {
	if r == nil {
		return ""
	}
	return r.Value
}

func (r *Rule) GetValueInt(defaultRule *Rule) int64 {
	value := r.GetValue()
	i, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return i
	}
	i, err = strconv.ParseInt(defaultRule.GetValue(), 10, 64)
	if err == nil {
		return i
	}
	return 0
}

func (s *Storage) GetRuleTemplateByName(name string) (*RuleTemplate, bool, error) {
	t := &RuleTemplate{}
	err := s.db.Where("name = ?", name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetRuleTemplateDetailByName(name string) (*RuleTemplate, bool, error) {
	t := &RuleTemplate{Name: name}
	err := s.db.Preload("RuleList").Preload("RuleList.Rule").Preload("Instances").Where(t).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateRuleTemplateRules(tpl *RuleTemplate, rules ...RuleTemplateRule) error {
	if err := s.db.Where(&RuleTemplateRule{RuleTemplateId: tpl.ID}).Delete(&RuleTemplateRule{}).Error; err != nil {
		return errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	err := s.db.Model(tpl).Association("RuleList").Append(rules).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateRuleTemplateInstances(tpl *RuleTemplate, instances ...*Instance) error {
	err := s.db.Model(tpl).Association("Instances").Replace(instances).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) CloneRuleTemplateRules(source, destination *RuleTemplate) error {
	return s.UpdateRuleTemplateRules(destination, source.RuleList...)
}

func GetRuleMapFromAllArray(allRules ...[]Rule) map[string]Rule {
	ruleMap := map[string]Rule{}
	for _, rules := range allRules {
		for _, rule := range rules {
			ruleMap[rule.Name] = rule
		}
	}
	return ruleMap
}

func (s *Storage) GetRuleTemplateTips() ([]*RuleTemplate, error) {
	ruleTemplates := []*RuleTemplate{}
	err := s.db.Select("name").Find(&ruleTemplates).Error
	return ruleTemplates, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetAllRule() ([]Rule, error) {
	rules := []Rule{}
	err := s.db.Find(&rules).Error
	return rules, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetRulesByInstanceId(instanceId string) ([]Rule, error) {
	rules := []Rule{}
	instance, _, err := s.GetInstanceById(instanceId)
	if err != nil {
		return rules, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	templates := instance.RuleTemplates
	if len(templates) <= 0 {
		return rules, nil
	}
	templateIds := make([]string, len(templates))
	for n, template := range templates {
		templateIds[n] = fmt.Sprintf("%v", template.ID)
	}

	err = s.db.Table("rules").
		Joins("inner join rule_template_rule on rule_template_rule.rule_name = rules.name").
		Where("rule_template_rule.rule_template_id in (?)", templateIds).
		Group("rules.name").Scan(&rules).Error
	return rules, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetRuleTemplatesByNames(names []string) ([]*RuleTemplate, error) {
	templates := []*RuleTemplate{}
	err := s.db.Where("name in (?)", names).Find(&templates).Error
	return templates, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetAndCheckRuleTemplateExist(templateNames []string) (ruleTemplates []*RuleTemplate, err error) {
	ruleTemplates, err = s.GetRuleTemplatesByNames(templateNames)
	if err != nil {
		return ruleTemplates, err
	}
	existTemplateNames := map[string]struct{}{}
	for _, user := range ruleTemplates {
		existTemplateNames[user.Name] = struct{}{}
	}
	notExistTemplateNames := []string{}
	for _, userName := range templateNames {
		if _, ok := existTemplateNames[userName]; !ok {
			notExistTemplateNames = append(notExistTemplateNames, userName)
		}
	}
	if len(notExistTemplateNames) > 0 {
		return ruleTemplates, errors.New(errors.DataNotExist,
			fmt.Errorf("rule template %s not exist", strings.Join(notExistTemplateNames, ", ")))
	}
	return ruleTemplates, nil
}

func (s *Storage) GetRulesByNames(names []string) ([]Rule, error) {
	rules := []Rule{}
	err := s.db.Where("name in (?)", names).Find(&rules).Error
	return rules, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetAndCheckRuleExist(ruleNames []string) (rules []Rule, err error) {
	rules, err = s.GetRulesByNames(ruleNames)
	if err != nil {
		return rules, err
	}
	existRuleNames := map[string]struct{}{}
	for _, user := range rules {
		existRuleNames[user.Name] = struct{}{}
	}
	notExistRuleNames := []string{}
	for _, userName := range ruleNames {
		if _, ok := existRuleNames[userName]; !ok {
			notExistRuleNames = append(notExistRuleNames, userName)
		}
	}
	if len(notExistRuleNames) > 0 {
		return rules, errors.New(errors.DataNotExist,
			fmt.Errorf("rule %s not exist", strings.Join(notExistRuleNames, ", ")))
	}
	return rules, nil
}
