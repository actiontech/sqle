package model

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
