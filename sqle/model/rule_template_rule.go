package model

type RuleTemplateRule struct {
	RuleTemplateId uint   `json:"rule_template_id" gorm:"primary_key;auto_increment:false;not null;"`
	RuleName       string `json:"name" gorm:"primary_key;not null;"`
	RuleLevel      string `json:"level" gorm:"column:level;"`
	RuleValue      string `json:"value" gorm:"column:value;" `

	Rule Rule `json:"-" gorm:"foreignkey:name;association_foreignkey:rule_name"`
}

func (rtr *RuleTemplateRule) TableName() string {
	return "rule_template_rule"
}
