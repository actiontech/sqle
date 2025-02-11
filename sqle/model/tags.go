package model

func init() {
	autoMigrateList = append(autoMigrateList, &Tag{})
}

// TypeTag 标签类型
type TypeTag string

// Tag 标签
type Tag struct {
	ID        uint        `json:"id" gorm:"primary_key;autoIncrement" example:"1"`
	Name      TypeTag     `json:"name" gorm:"type:varchar(50);not null;uniqueIndex:idx_tag_name;comment:标签名称"` // varchar50大约可以写16个汉字
	Knowledge []Knowledge `json:"knowledge" gorm:"many2many:knowledge_tag_relations"`
	SubTag    []*Tag      `json:"sub_tag" gorm:"many2many:tag_tag_relations"`
}

func (Tag) TableName() string {
	return "tags"
}

// KnowledgeTagRelation 知识库标签关联表，由Gorm自动创建，不需要AutoMigrate
type KnowledgeTagRelation struct {
	KnowledgeID uint `gorm:"column:knowledge_id;primaryKey;not null;type:bigint unsigned;comment:知识库ID"`
	TagID       uint `gorm:"column:tag_id;primaryKey;not null;type:bigint unsigned;comment:标签ID"`
}

func (KnowledgeTagRelation) TableName() string {
	return "knowledge_tag_relations"
}

// TagTagRelation 标签标签关联表，由Gorm自动创建，不需要AutoMigrate
type TagTagRelation struct {
	TagID    uint `gorm:"column:tag_id;primaryKey;not null;type:bigint unsigned;comment:标签ID"`
	SubTagID uint `gorm:"column:sub_tag_id;primaryKey;not null;type:bigint unsigned;comment:子标签ID"`
}
