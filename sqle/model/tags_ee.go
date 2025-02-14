package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 创建标签的子标签
func (s *Storage) CreateSubTag(parentTagId uint, subTag *Tag) (*Tag, error) {
	return subTag, s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(subTag).Error
		if err != nil {
			return err
		}
		err = tx.Create(&TagTagRelation{
			TagID:    parentTagId,
			SubTagID: subTag.ID,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

// 获取或创建标签
func (s *Storage) GetOrCreateTag(tag *Tag) (*Tag, error) {
	var existTag Tag
	err := s.db.Where("name =?", tag.Name).First(&existTag).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	if existTag.ID == 0 {
		err = s.db.Create(tag).Error
		if err != nil {
			return nil, err
		}
		return tag, nil
	}
	return &existTag, nil
}

// 获取标签的子标签
func (s *Storage) GetSubTags(tagId uint) ([]*Tag, error) {
	var tags []*Tag
	err := s.db.Model(&Tag{}).
		Joins("LEFT JOIN tag_tag_relations ON tags.id = tag_tag_relations.sub_tag_id").
		Where("tag_id = ?", tagId).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}


// 根据Tag IDs 获取标签与标签的关联
func (s *Storage) GetTagRelationsByTagIds(tagIds []uint) ([]TagTagRelation, error) {
	var relations []TagTagRelation
	err := s.db.Model(&TagTagRelation{}).Where("sub_tag_id IN ?", tagIds).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// 为知识库添加标签
func (s *Storage) CreateKnowledgeTagRelation(knowledgeId uint, tagId uint) error {
	if knowledgeId == 0 || tagId == 0 {
		return nil
	}
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "knowledge_id"}, {Name: "tag_id"}},
		DoNothing: true, // 如果存在则不插入也不更新
	}).Create(&KnowledgeTagRelation{
		KnowledgeID: knowledgeId,
		TagID:       tagId,
	}).Error
}

// 创建标签，如果标签已存在则不插入
func (s *Storage) CreateTag(tag *Tag) (*Tag, error) {
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true, // 如果存在则不插入也不更新
	}).Create(tag).Error
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// 根据标签名称获取标签
func (s *Storage) GetTagByName(name TypeTag) (*Tag, error) {
	var tag Tag
	err := s.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// 根据标签ID获取标签
func (s *Storage) GetTagById(id uint) (*Tag, error) {
	var tag Tag
	err := s.db.Where("id =?", id).First(&tag).Error
	if err!= nil {
		return nil, err
	}
	return &tag, nil
}

// 获取标签关联关系
func (s *Storage) GetTagRelations() ([]TagTagRelation, error) {
	var relations []TagTagRelation
	err := s.db.Model(&TagTagRelation{}).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// 获取所有知识库标签的关联关系
func (s *Storage) GetKnowledgeBaseTagRelations() ([]*KnowledgeTagRelation, error) {
	var relations []*KnowledgeTagRelation
	err := s.db.Model(&KnowledgeTagRelation{}).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// 根据一个标签获取知识库
func (s *Storage) GetKnowledgeByTagAndTitle(tagId []uint, title string) (*Knowledge, error) {
	var knowledge Knowledge
	err := s.db.Model(&Knowledge{}).
		Joins("LEFT JOIN knowledge_tag_relations ON knowledge.id = knowledge_tag_relations.knowledge_id").
		Where("knowledge_tag_relations.tag_id =?", tagId).
		Where("knowledge.title =?", title).
		First(&knowledge).Error
	if err != nil {
		return nil, err
	}
	return &knowledge, nil
}
