package model

import ()

type Knowledge struct {
	Model
	RuleName    string `gorm:"type:varchar(255);default:''" json:"rule_name"`  // 规则名称
	Title       string `gorm:"type:varchar(255);not null" json:"title"`        // 标题
	Description string `gorm:"type:text" json:"description"`                   // 描述
	Content     string `gorm:"type:text;index:,class:FULLTEXT" json:"content"` // 内容
	Tags        []*Tag `gorm:"many2many:knowledge_tag_relations" json:"tags"`  // 标签
}

func (Knowledge) TableName() string {
	return "knowledge"
}
