package model

import "github.com/jinzhu/gorm"

type RuleTemplate struct {
	Model
	Name  string
	Desc  string
	Rules []Rule `gorm:"many2many:rule_template_rule"`
}

type Rule struct {
	Name  string `json:"name" gorm:"primary_key"`
	Desc  string `json:"desc"`
	Value string `json:"value"`
	// notice, warn, error
	Level string `json:"level" example:"error"`
}

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

// inspector rule code
const (
	DDL_ALL_CHECK_NAME_LENGTH  = "ddl_all_check_name_length"
	DDL_CREATE_TABLE_NOT_EXIST = "ddl_create_table_not_exist"
)

var DefaultRules = []Rule{
	Rule{
		Name:  DDL_ALL_CHECK_NAME_LENGTH,
		Desc:  "",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_CREATE_TABLE_NOT_EXIST,
		Desc:  "",
		Level: RULE_LEVEL_ERROR,
	},
}

func (s *Storage) GetTemplateById(templateId string) (RuleTemplate, bool, error) {
	t := RuleTemplate{}
	err := s.db.Preload("Rules").Where("id = ?", templateId).First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, err
}

func (s *Storage) GetTemplateByName(name string) (RuleTemplate, bool, error) {
	t := RuleTemplate{}
	err := s.db.Preload("Rules").Where("name = ?", name).First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, err
}

func (s *Storage) UpdateRules(tpl *RuleTemplate, rules ...Rule) error {
	return s.db.Model(tpl).Association("Rules").Append(rules).Error
}

func (s *Storage) GetAllTemplate() ([]RuleTemplate, error) {
	ts := []RuleTemplate{}
	err := s.db.Preload("Rules").Find(&ts).Error
	return ts, err
}

func (s *Storage) CreateDefaultRules() error {
	for _, rule := range DefaultRules {
		exist, err := s.Exist(&rule)
		if err != nil {
			return err
		}
		if exist {
			continue
		}
		err = s.Save(rule)
		if err != nil {
			return err
		}
	}
	return nil
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

func (s *Storage) GetAllRule() ([]Rule, error) {
	rules := []Rule{}
	err := s.db.Find(&rules).Error
	return rules, err
}
