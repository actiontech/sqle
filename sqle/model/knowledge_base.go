package model

import ()

func init() {
	autoMigrateList = append(autoMigrateList, &Knowledge{})
}

type Knowledge struct {
	Model
	Title       string `gorm:"type:varchar(255);not null" json:"title"`       // 标题
	Description string `gorm:"type:text" json:"description"`                  // 描述
	Content     string `gorm:"type:text" json:"content"`                      // 内容
	Tags        []*Tag  `gorm:"many2many:knowledge_tag_relations" json:"tags"` // 标签
}

func (Knowledge) TableName() string {
	return "knowledge"
}

// 为知识表添加全文索引
func (s *Storage) addFullTextIndexForKnowledge() error {
	if !s.db.Migrator().HasIndex(&Knowledge{}, "fulltext_content") {
		err := s.db.Exec(`
				ALTER TABLE knowledge
				ADD FULLTEXT INDEX fulltext_content (content);
			`).Error
		if err != nil {
			return err
		}
	}
	return nil
}
