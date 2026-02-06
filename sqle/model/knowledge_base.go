package model

type Knowledge struct {
	Model
	Title       string        `gorm:"type:varchar(255);not null" json:"title"`                       // 标题
	Description string        `gorm:"type:text" json:"description"`                                  // 描述
	Content     string        `gorm:"type:text" json:"content"`                // 内容
	Tags        []*Tag        `gorm:"many2many:knowledge_tag_relations" json:"tags"`                 // 标签
	Rules       []*Rule       `gorm:"many2many:rule_knowledge_relations" json:"rules"`               // 规则和知识的关系
	CustomRules []*CustomRule `gorm:"many2many:custom_rule_knowledge_relations" json:"custom_rules"` // 自定义规则和知识的关系
}

type MultiLanguageKnowledge []*Knowledge

func (Knowledge) TableName() string {
	return "knowledge"
}
