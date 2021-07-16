package model

type RuleTemplateRule struct {
	RuleTemplateId uint   `json:"rule_template_id" gorm:"primary_key;auto_increment:false;"`
	RuleName       string `json:"name" gorm:"primary_key;"`
	RuleLevel      string `json:"level" gorm:"column:level;"`
	RuleValue      string `json:"value" gorm:"column:value;" `

	Rule Rule `json:"-" gorm:"foreignkey:Name;association_foreignkey:RuleName"`
}

func (rtr *RuleTemplateRule) TableName() string {
	return "rule_template_rule"
}
