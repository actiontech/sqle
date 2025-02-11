package model

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// 知识库标签
const (
	PredefineTagCustomizeKnowledgeBase TypeTag = "CustomizeKnowledgeBase"
	PredefineTagKnowledgeBase          TypeTag = "KnowledgeBase"
	WHERE                              TypeTag = "WHERE"
	SELECT                             TypeTag = "SELECT"
	EXPLAIN                            TypeTag = "EXPLAIN"
	TABLE                              TypeTag = "TABLE"
	INDEX                              TypeTag = "INDEX"
	UPDATE                             TypeTag = "UPDATE"
	DELETE                             TypeTag = "DELETE"
	TRANSACTION                        TypeTag = "TRANSACTION"
	ALTER                              TypeTag = "ALTER"
	AUTO_INCREMENT                     TypeTag = "AUTO_INCREMENT"
	COLUMN                             TypeTag = "COLUMN"
	PRIMARY_KEY                        TypeTag = "PRIMARY KEY"
	DATA_TYPE                          TypeTag = "DATA_TYPE"
	DECIMAL                            TypeTag = "DECIMAL"
	BIGINT                             TypeTag = "BIGINT"
	CHARSET                            TypeTag = "CHARSET"
	COLLATION                          TypeTag = "COLLATION"
	BLOB                               TypeTag = "BLOB"
	TEXT                               TypeTag = "TEXT"
	DEFAULT                            TypeTag = "DEFAULT"
	NOT_NULL                           TypeTag = "NOT NULL"
	JOIN                               TypeTag = "JOIN"
	SET                                TypeTag = "SET"
	ENGINE                             TypeTag = "ENGINE"
	// TODO 增加新的标签
)

// 语言标签
const (
	PredefineTagLanguage TypeTag = "language"
	PredefineTagChinese  TypeTag = "zh"
	PredefineTagEnglish  TypeTag = "en"
)

func GetTagMapPredefineLanguage() map[TypeTag]language.Tag {
	return map[TypeTag]language.Tag{
		PredefineTagChinese: language.Chinese,
		PredefineTagEnglish: language.English,
	}
}

// 数据库类型标签
const (
	PredefineTagDBType         TypeTag = "DBType"
	PredefineTagMySQL          TypeTag = driverV2.DriverTypeMySQL
	PredefineTagSQLServer      TypeTag = driverV2.DriverTypeSQLServer
	PredefineTagOracle         TypeTag = driverV2.DriverTypeOracle
	PredefineTagPostgreSQL     TypeTag = driverV2.DriverTypePostgreSQL
	PredefineTagTDSQLForInnoDB TypeTag = driverV2.DriverTypeTDSQLForInnoDB
	// TODO 增加新的数据库类型
)

// 获取数据库预定义标签映射
func GetTagMapPredefineDBType() map[TypeTag]struct{} {
	return map[TypeTag]struct{}{
		PredefineTagMySQL:          {},
		PredefineTagSQLServer:      {},
		PredefineTagOracle:         {},
		PredefineTagPostgreSQL:     {},
		PredefineTagTDSQLForInnoDB: {},
	}
}

func GetTagMapDefaultRuleKnowledge() map[string] /* rule_name */ []TypeTag {
	return map[string] /* rule_name */ []TypeTag{
		rulepkg.DMLCheckWhereIsInvalid:             {WHERE, SELECT, EXPLAIN, TABLE, INDEX},
		rulepkg.DMLCheckUpdateOrDeleteHasWhere:     {WHERE, UPDATE, DELETE, TABLE, INDEX, TRANSACTION},
		rulepkg.DDLCheckAlterTableNeedMerge:        {ALTER, TABLE, INDEX},
		rulepkg.DDLCheckAutoIncrement:              {AUTO_INCREMENT, COLUMN, PRIMARY_KEY, TABLE},
		rulepkg.DDLCheckBigintInsteadOfDecimal:     {DATA_TYPE, COLUMN, INDEX, DECIMAL, BIGINT},
		rulepkg.DDLCheckDatabaseCollation:          {CHARSET, COLLATION, TABLE, COLUMN, INDEX, JOIN},
		rulepkg.DDLCheckColumnBlobDefaultIsNotNull: {BLOB, TEXT, COLUMN, DEFAULT, NOT_NULL},
		rulepkg.DDLCheckColumnBlobNotice:           {BLOB, TEXT, TABLE, INDEX},
		rulepkg.DDLCheckColumnBlobWithNotNull:      {BLOB, TEXT, COLUMN, DEFAULT, NOT_NULL},
		rulepkg.DDLCheckColumnQuantity:             {COLUMN, TABLE},
		rulepkg.DDLCheckColumnQuantityInPK:         {PRIMARY_KEY, COLUMN, INDEX, ENGINE},
		rulepkg.DDLCheckColumnSetNotice:            {DATA_TYPE, COLUMN, SET},
		// TODO 增加新的规则和标签的映射
	}
}

// 创建知识库
func (s *Storage) CreateKnowledgeWithTags(knowledge *Knowledge, tags map[TypeTag] /* tag name */ *Tag, filterTags []*Tag) (*Knowledge, error) {
	// 先查询是否存在同名的知识库
	modelKnowledge, err := s.GetKnowledgeByTagsAndTitle(filterTags, knowledge.Title)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	// 标签映射
	modelTagMap := make(map[TypeTag] /* tag name */ *Tag)

	// 如果不存在，则创建
	if modelKnowledge == nil {
		err = s.db.Omit("Tags").Create(knowledge).Error
		if err != nil {
			return nil, err
		}
		modelKnowledge = knowledge
	} else {
		for _, tag := range modelKnowledge.Tags {
			modelTagMap[tag.Name] = tag
		}
	}
	// 若标签不存在，则关联标签，若标签已存在则跳过
	for _, tag := range tags {
		if _, exist := modelTagMap[tag.Name]; exist {
			continue
		}
		err = s.CreateKnowledgeTagRelation(modelKnowledge.ID, tag.ID)
		if err != nil {
			return nil, err
		}
	}
	return knowledge, nil
}

// 根据标签和规则名称获取知识库
func (s *Storage) GetKnowledgeByTagsAndTitle(filterTags []*Tag, title string) (*Knowledge, error) {
	var (
		modelKnowledge Knowledge
		tagIds         []uint
	)

	for _, tag := range filterTags {
		tagIds = append(tagIds, tag.ID)
	}

	err := s.db.Model(&Knowledge{}).Preload(`Tags`).
		Joins(`JOIN knowledge_tag_relations ktr ON knowledge.id = ktr.knowledge_id JOIN tags t ON ktr.tag_id = t.id`).
		Where(`t.id IN ? AND knowledge.title = ?`, tagIds, title).
		Group(`knowledge.id`).
		Having("COUNT(DISTINCT t.id) = ?", len(tagIds)).
		First(&modelKnowledge).Error
	if err != nil {
		return nil, err
	}
	return &modelKnowledge, nil
}

// 查询知识库知识列表，支持使用关键字和标签进行搜索
func (s *Storage) SearchKnowledge(keyword string, tags []string, limit, offset int) ([]Knowledge, int64, error) {
	if keyword == "" && len(tags) == 0 {
		return nil, 0, nil
	}
	if limit <= 0 {
		limit = 20
	}
	var results []Knowledge
	// 该查询涉及knowledge_tag_relations表和tags表，根据关联表knowledge_tag_relations关联查询
	searchClause := s.db.
		Table("knowledge").
		Preload("Tags").
		Select(`knowledge.id, description, title`).
		Joins("LEFT JOIN knowledge_tag_relations ON knowledge_tag_relations.knowledge_id = knowledge.id LEFT JOIN tags ON tags.id = knowledge_tag_relations.tag_id")

	countClause := s.db.
		Table("knowledge").
		Joins("LEFT JOIN knowledge_tag_relations ON knowledge_tag_relations.knowledge_id = knowledge.id LEFT JOIN tags ON tags.id = knowledge_tag_relations.tag_id")
	// 如果有关键字，则根据关键字进行模糊查询+全文检索
	if len(keyword) > 0 {
		likeClause := "%" + keyword + "%"
		searchClause = searchClause.Select(`knowledge.id, SUBSTRING(content, LOCATE(?, content), 50) AS content, description, title`, keyword).Where("title LIKE ? OR description LIKE ? OR MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)", likeClause, likeClause, keyword)
		countClause = countClause.Where("title LIKE ? OR description LIKE ? OR MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)", likeClause, likeClause, keyword)
	}
	// 如果有标签，则根据标签查询带有标签的知识库
	if len(tags) > 0 {
		searchClause = searchClause.Where("tags.name IN ?", tags)
		countClause = countClause.Where("tags.name IN ?", tags)
	}

	// 分页查询
	err := searchClause.Group("knowledge.id").Limit(limit).Offset(offset).Find(&results).Error
	if err != nil {
		return nil, 0, err
	}
	// 统计查询结果总数
	var count int64
	err = countClause.Group("knowledge.id").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

// 更新知识库内容
func (s *Storage) UpdateKnowledgeContent(knowledge *Knowledge, newContent string) error {
	return s.db.Model(knowledge).Updates(map[string]interface{}{"content": newContent}).Error
}
