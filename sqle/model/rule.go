package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/jinzhu/gorm"
)

type RuleTemplate struct {
	Model
	Name      string             `json:"name"`
	Desc      string             `json:"desc"`
	DBType    string             `json:"db_type"`
	Instances []Instance         `json:"instance_list" gorm:"many2many:instance_rule_template"`
	RuleList  []RuleTemplateRule `json:"rule_list" gorm:"foreignkey:rule_template_id;association_foreignkey:id"`
}

func GenerateRuleByDriverRule(dr *driver.Rule, dbType string) *Rule {
	return &Rule{
		Name:   dr.Name,
		Desc:   dr.Desc,
		Level:  string(dr.Level),
		Typ:    dr.Category,
		DBType: dbType,
		Params: dr.Params,
	}
}

func ConvertRuleToDriverRule(r *Rule) *driver.Rule {
	return &driver.Rule{
		Name:     r.Name,
		Desc:     r.Desc,
		Category: r.Typ,
		Level:    driver.RuleLevel(r.Level),
		Params:   r.Params,
	}
}

type Rule struct {
	Name   string        `json:"name" gorm:"primary_key; not null"`
	DBType string        `json:"db_type" gorm:"primary_key; not null; default:\"mysql\""`
	Desc   string        `json:"desc"`
	Level  string        `json:"level" example:"error"` // notice, warn, error
	Typ    string        `json:"type" gorm:"column:type; not null"`
	Params params.Params `json:"params" gorm:"type:varchar(1000)"`
}

func (r Rule) TableName() string {
	return "rules"
}

type RuleTemplateRule struct {
	RuleTemplateId uint          `json:"rule_template_id" gorm:"primary_key;auto_increment:false;"`
	RuleName       string        `json:"name" gorm:"primary_key;"`
	RuleLevel      string        `json:"level" gorm:"column:level;"`
	RuleParams     params.Params `json:"value" gorm:"column:rule_params;type:varchar(1000)"`
	RuleDBType     string        `json:"rule_db_type" gorm:"column:db_type; not null;"`

	Rule *Rule `json:"-" gorm:"foreignkey:Name,DBType;association_foreignkey:RuleName,RuleDBType"`
}

func (rtr *RuleTemplateRule) TableName() string {
	return "rule_template_rule"
}

func NewRuleTemplateRule(t *RuleTemplate, r *Rule) RuleTemplateRule {
	return RuleTemplateRule{
		RuleTemplateId: t.ID,
		RuleName:       r.Name,
		RuleLevel:      r.Level,
		RuleParams:     r.Params,
		RuleDBType:     r.DBType,
	}
}

func (rtr *RuleTemplateRule) GetRule() *Rule {
	rule := rtr.Rule
	if rtr.RuleLevel != "" {
		rule.Level = rtr.RuleLevel
	}
	if rtr.RuleParams != nil && len(rtr.RuleParams) > 0 {
		rule.Params = rtr.RuleParams
	}
	return rule
}

func (s *Storage) GetRuleTemplatesByInstance(inst *Instance) ([]RuleTemplate, error) {
	var associationRT []RuleTemplate
	err := s.db.Model(inst).Association("RuleTemplates").Find(&associationRT).Error
	return associationRT, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRuleTemplatesByInstanceName(name string) (*RuleTemplate, error) {
	t := &RuleTemplate{}
	err := s.db.Joins("JOIN `instance_rule_template` ON `rule_templates`.`id` = `instance_rule_template`.`rule_template_id`").
		Joins("JOIN `instances` ON `instance_rule_template`.`instance_id` = `instances`.`id`").
		Where("`instances`.`name` = ?", name).Find(t).Error
	return t, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRulesFromRuleTemplateByName(name string) ([]*Rule, error) {
	tpl, exist, err := s.GetRuleTemplateDetailByName(name)
	if !exist {
		return nil, errors.New(errors.DataNotExist, err)
	}
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}

	rules := make([]*Rule, 0, len(tpl.RuleList))
	for _, r := range tpl.RuleList {
		rules = append(rules, r.GetRule())
	}
	return rules, nil
}

func (s *Storage) GetRulesByInstanceId(instanceId string) ([]*Rule, error) {
	instance, _, err := s.GetInstanceById(instanceId)
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	templates := instance.RuleTemplates
	if len(templates) <= 0 {
		return nil, nil
	}
	tplName := templates[0].Name
	return s.GetRulesFromRuleTemplateByName(tplName)
}

func (s *Storage) GetRuleTemplateByName(name string) (*RuleTemplate, bool, error) {
	t := &RuleTemplate{}
	err := s.db.Where("name = ?", name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRuleTemplateDetailByName(name string) (*RuleTemplate, bool, error) {
	dbOrder := func(db *gorm.DB) *gorm.DB {
		return db.Order("rule_template_rule.rule_name ASC")
	}
	t := &RuleTemplate{Name: name}
	err := s.db.Preload("RuleList", dbOrder).Preload("RuleList.Rule").Preload("Instances").
		Where(t).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateRuleTemplateRules(tpl *RuleTemplate, rules ...RuleTemplateRule) error {
	if err := s.db.Where(&RuleTemplateRule{RuleTemplateId: tpl.ID}).Delete(&RuleTemplateRule{}).Error; err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}

	//err := s.db.Debug().Model(tpl).Association("RuleList").Append(rules).Error //暂时没找到跳过更新关联表的方式
	//TODO gorm v1 没有预编译, 没有批量插入, 规则可能很多, 为了防止拼SQL批量插入导致SQL太长, 只能通过这种方式插入
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			rule.RuleTemplateId = tpl.ID
			err := tx.Omit("Rule").Create(&rule).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateRuleTemplateInstances(tpl *RuleTemplate, instances ...*Instance) error {
	err := s.db.Model(tpl).Association("Instances").Replace(instances).Error
	return errors.New(errors.ConnectStorageError, err)
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

func (s *Storage) GetRuleTemplateTips(dbType string) ([]*RuleTemplate, error) {
	ruleTemplates := []*RuleTemplate{}

	db := s.db.Select("name, db_type")
	if dbType != "" {
		db = db.Where("db_type = ?", dbType)
	}
	err := db.Find(&ruleTemplates).Error
	return ruleTemplates, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRule(name, dbType string) (*Rule, bool, error) {
	rule := Rule{Name: name, DBType: dbType}
	err := s.db.Find(&rule).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return &rule, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllRule() ([]*Rule, error) {
	rules := []*Rule{}
	err := s.db.Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllRuleByDBType(dbType string) ([]*Rule, error) {
	rules := []*Rule{}
	err := s.db.Where(&Rule{DBType: dbType}).Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRuleTemplatesByNames(names []string) ([]*RuleTemplate, error) {
	templates := []*RuleTemplate{}
	err := s.db.Where("name in (?)", names).Find(&templates).Error
	return templates, errors.New(errors.ConnectStorageError, err)
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

func (s *Storage) GetRulesByNames(names []string, dbType string) ([]Rule, error) {
	rules := []Rule{}
	err := s.db.Where("db_type = ?", dbType).Where("name in (?)", names).Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAndCheckRuleExist(ruleNames []string, dbType string) (map[string]Rule, error) {
	rules, err := s.GetRulesByNames(ruleNames, dbType)
	if err != nil {
		return nil, err
	}
	existRules := map[string]Rule{}
	for _, rule := range rules {
		existRules[rule.Name] = rule
	}
	notExistRuleNames := []string{}
	for _, userName := range ruleNames {
		if _, ok := existRules[userName]; !ok {
			notExistRuleNames = append(notExistRuleNames, userName)
		}
	}
	if len(notExistRuleNames) > 0 {
		return nil, errors.New(errors.DataNotExist,
			fmt.Errorf("rule %s not exist", strings.Join(notExistRuleNames, ", ")))
	}
	return existRules, nil
}

func (s *Storage) IsRuleTemplateExist(ruleTemplateName string) (bool, error) {
	var count int
	err := s.db.Table("rule_templates").Where("name = ?", ruleTemplateName).Count(&count).Error
	return count > 0, errors.New(errors.ConnectStorageError, err)
}
