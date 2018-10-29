package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type RuleTemplate struct {
	Model
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Rules []Rule `json:"-" gorm:"many2many:rule_template_rule"`
}

type Rule struct {
	Name  string `json:"name" gorm:"primary_key"`
	Desc  string `json:"desc"`
	Value string `json:"value"`
	// notice, warn, error
	Level string `json:"level" example:"error"`
}

func (r Rule) TableName() string {
	return "rules"
}

// RuleTemplateDetail use for http request and swagger docs;
// it is same as RuleTemplate, but display Rules in json format.
type RuleTemplateDetail struct {
	RuleTemplate
	Rules []Rule `json:"rule_list"`
}

func (r *RuleTemplate) Detail() RuleTemplateDetail {
	data := RuleTemplateDetail{
		RuleTemplate: *r,
		Rules:        r.Rules,
	}
	if r.Rules == nil {
		data.Rules = []Rule{}
	}
	return data
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
	SCHEMA_NOT_EXIST              = "schema_not_exist"
	SCHEMA_EXIST                  = "schema_exist"
	TABLE_NOT_EXIST               = "table_not_exist"
	TABLE_EXIST                   = "table_exist"
	DDL_CREATE_TABLE_NOT_EXIST    = "ddl_create_table_not_exist"
	DDL_CHECK_TABLE_NAME_LENGTH   = "ddl_check_table_name_length"
	DDL_CHECK_COLUMNS_NAME_LENGTH = "ddl_check_columns_name_length"
	DDL_CHECK_INDEX_NAME_LENGTH   = "ddl_check_schema_name_length"
	DDL_CHECK_PRIMARY_KEY_EXIST   = "ddl_check_primary_key_exist"
	DDL_CHECK_PRIMARY_KEY_TYPE    = "ddl_check_primary_key_type"
	DDL_DISABLE_VARCHAR_MAX       = "ddl_disable_varchar_max"
	DDL_CHECK_TYPE_CHAR_LENGTH    = "ddl_check_type_char_length"
)

var RuleMessageMap = map[string]string{
	SCHEMA_NOT_EXIST:              "schema %s 不存在",
	SCHEMA_EXIST:                  "schema %s 已存在",
	TABLE_NOT_EXIST:               "表 %s 不存在",
	TABLE_EXIST:                   "表 %s 已存在",
	DDL_CREATE_TABLE_NOT_EXIST:    "新建表必须加入if not exists create，保证重复执行不报错",
	DDL_CHECK_TABLE_NAME_LENGTH:   "%s 表名长度不能大于64字节",
	DDL_CHECK_COLUMNS_NAME_LENGTH: "%s 列名长度不能大于64字节",
	DDL_CHECK_INDEX_NAME_LENGTH:   "%s 索引名长度不能大于64字节",
	DDL_CHECK_PRIMARY_KEY_EXIST:   "表必须有主键",
	DDL_CHECK_PRIMARY_KEY_TYPE:    "主键建议使用自增，且为bigint无符号类型，即bigint unsigned",
	DDL_DISABLE_VARCHAR_MAX:       "禁止使用 varchar(max)",
	DDL_CHECK_TYPE_CHAR_LENGTH:    "char长度大于20时，必须使用varchar类型",
}

var DefaultRules = []Rule{
	Rule{
		Name:  SCHEMA_NOT_EXIST,
		Desc:  "操作数据库时，数据库必须存在",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  SCHEMA_EXIST,
		Desc:  "创建数据库时，数据库不能存在",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  TABLE_NOT_EXIST,
		Desc:  "操作表时，表必须不存在",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  TABLE_EXIST,
		Desc:  "创建表时，表不能存在",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_CREATE_TABLE_NOT_EXIST,
		Desc:  "新建表必须加入if not exists create，保证重复执行不报错",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_CHECK_TABLE_NAME_LENGTH,
		Desc:  "表名长度不能大于64字节",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_CHECK_COLUMNS_NAME_LENGTH,
		Desc:  "列名长度不能大于64字节",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_CHECK_INDEX_NAME_LENGTH,
		Desc:  "索引名长度不能大于64字节",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_CHECK_PRIMARY_KEY_EXIST,
		Desc:  "表必须有主键",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_CHECK_PRIMARY_KEY_TYPE,
		Desc:  "主键建议使用自增，且为bigint无符号类型，即bigint unsigned",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_DISABLE_VARCHAR_MAX,
		Desc:  "禁止使用 varchar(max)",
		Level: RULE_LEVEL_ERROR,
	},
	Rule{
		Name:  DDL_CHECK_TYPE_CHAR_LENGTH,
		Desc:  "char长度大于20时，必须使用varchar类型",
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

func (s *Storage) UpdateTemplateRules(tpl *RuleTemplate, rules ...Rule) error {
	return s.db.Model(tpl).Association("Rules").Replace(rules).Error
}

func (s *Storage) GetAllTemplate() ([]RuleTemplate, error) {
	ts := []RuleTemplate{}
	err := s.db.Find(&ts).Error
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

func (s *Storage) GetRulesByInstanceId(instanceId string) ([]Rule, error) {
	rules := []Rule{}
	instance, _, err := s.GetInstById(instanceId)
	if err != nil {
		return rules, err
	}
	templates := instance.RuleTemplates
	if len(templates) <= 0 {
		return rules, nil
	}
	templateIds := make([]string, len(templates))
	for n, template := range templates {
		templateIds[n] = fmt.Sprintf("%v", template.ID)
	}

	err = s.db.Table("rules").Select("rules.name, rules.value, rules.level").
		Joins("inner join rule_template_rule on rule_template_rule.rule_name = rules.name").
		Where("rule_template_rule.rule_template_id in (?)", templateIds).Scan(&rules).Error
	return rules, err
}
