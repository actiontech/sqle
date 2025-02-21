package model

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	// "gorm.io/gorm"
)

// 知识库标签
const (
	PredefineTagCustomizeKnowledgeBase TypeTag = "CustomizeKnowledgeBase"
	PredefineTagKnowledgeBase          TypeTag = "KnowledgeBase"
)

// SQL对象关键字
const (
	DATABASE  TypeTag = "DATABASE"
	SCHEMA    TypeTag = "SCHEMA"
	TABLE     TypeTag = "TABLE"
	INDEX     TypeTag = "INDEX"
	COLUMN    TypeTag = "COLUMN"
	ROW       TypeTag = "ROW"
	VIEW      TypeTag = "VIEW"
	PROCEDURE TypeTag = "PROCEDURE"
	FUNCTION  TypeTag = "FUNCTION"
	TRIGGER   TypeTag = "TRIGGER"
	USER      TypeTag = "USER"
	ROLE      TypeTag = "ROLE"
	EVENT     TypeTag = "EVENT"
	PARTITION TypeTag = "PARTITION"
	SEQUENCE  TypeTag = "SEQUENCE"
	KEY       TypeTag = "KEY"
)

// SQL操作关键字
const (
	SELECT TypeTag = "SELECT"
	INSERT TypeTag = "INSERT"
	UPDATE TypeTag = "UPDATE"
	DELETE TypeTag = "DELETE"
	CREATE TypeTag = "CREATE"
	ALTER  TypeTag = "ALTER"
	DROP   TypeTag = "DROP"
	RENAME TypeTag = "RENAME"
	ADD    TypeTag = "ADD"
	SET    TypeTag = "SET"
	LIMIT  TypeTag = "LIMIT"
)

// SQL对象描述词
const (
	UNIQUE  TypeTag = "UNIQUE"
	PRIMARY TypeTag = "PRIMARY"
	FOREIGN TypeTag = "FOREIGN"
	ENUM    TypeTag = "ENUM"
	VARCHAR TypeTag = "VARCHAR"
)

// 其他SQL描述词
const (
	REFERENCES     TypeTag = "REFERENCES"
	MODIFY         TypeTag = "MODIFY"
	LENGTH         TypeTag = "LENGTH"
	WHERE          TypeTag = "WHERE"
	FROM           TypeTag = "FROM"
	CONSTRAINT     TypeTag = "CONSTRAINT"
	EXPLAIN        TypeTag = "EXPLAIN"
	STATEMENT      TypeTag = "STATEMENT"
	IF             TypeTag = "IF"
	NOT            TypeTag = "NOT"
	EXISTS         TypeTag = "EXISTS"
	DECIMAL        TypeTag = "DECIMAL"
	HAVING         TypeTag = "HAVING"
	AS             TypeTag = "AS"
	GROUP          TypeTag = "GROUP"
	BY             TypeTag = "BY"
	OR             TypeTag = "OR"
	SUBQUERY       TypeTag = "SUBQUERY"
	COUNT          TypeTag = "COUNT"
	CHARSET        TypeTag = "CHARSET"
	JOIN           TypeTag = "JOIN"
	AUTO_INCREMENT TypeTag = "AUTO_INCREMENT"
	COLLATION      TypeTag = "COLLATION"
	ASC            TypeTag = "ASC"
	DESC           TypeTag = "DESC"
	PRIVILEGE      TypeTag = "PRIVILEGE"
	TEMPORARY      TypeTag = "TEMPORARY"
	VALUES         TypeTag = "VALUES"
	TRANSACTION    TypeTag = "TRANSACTION"
	ISOLATION      TypeTag = "ISOLATION"
	LEVEL          TypeTag = "LEVEL"
	COMMENT        TypeTag = "COMMENT"
	IN             TypeTag = "IN"
	LIKE           TypeTag = "LIKE"
	RAND           TypeTag = "RAND"
	ENGINE         TypeTag = "ENGINE"
	BLOB           TypeTag = "BLOB"
	FOR            TypeTag = "FOR"
	GLOBAL         TypeTag = "GLOBAL"
	INT            TypeTag = "INT"
	TO             TypeTag = "TO"
	NOT_NULL       TypeTag = "NOT_NULL"
	BIGINT         TypeTag = "BIGINT"
	FLOAT          TypeTag = "FLOAT"
	CHAR           TypeTag = "CHAR"
	TIMESTAMP      TypeTag = "TIMESTAMP"
	HINT           TypeTag = "HINT"
	TRUNCATE       TypeTag = "TRUNCATE"
	GRANT          TypeTag = "GRANT"
	ORDER          TypeTag = "ORDER"
	ON             TypeTag = "ON"
	NULL           TypeTag = "NULL"
	IS             TypeTag = "IS"
	INTO           TypeTag = "INTO"
	UNION          TypeTag = "UNION"
	OFFSET         TypeTag = "OFFSET"
	EACH           TypeTag = "EACH"
)

// 版本标签
const (
	PredefineTagVersion TypeTag = "version"
	PredefineTagV1      TypeTag = "v1"
	PredefineTagV2      TypeTag = "v2"
)

// GetTagMapPredefineVersion 获取预定义版本标签映射
func GetTagMapPredefineVersion() map[TypeTag]struct{} {
	return map[TypeTag]struct{}{
		PredefineTagV1: {},
		PredefineTagV2: {},
	}
}

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
	PredefineTagTiDB           TypeTag = driverV2.DriverTypeTiDB
	PredefineTagDB2            TypeTag = driverV2.DriverTypeDB2
	PredefineTagTBase          TypeTag = driverV2.DriverTypeTBase
	PredefineTagOceanBase      TypeTag = driverV2.DriverTypeOceanBase
	PredefineTagActionDB       TypeTag = "ActionDB"
)

// 获取数据库预定义标签映射
func GetTagMapPredefineDBType() map[TypeTag]struct{} {
	return map[TypeTag]struct{}{
		PredefineTagMySQL:          {},
		PredefineTagSQLServer:      {},
		PredefineTagOracle:         {},
		PredefineTagPostgreSQL:     {},
		PredefineTagTDSQLForInnoDB: {},
		PredefineTagTiDB:           {},
		PredefineTagDB2:            {},
		PredefineTagTBase:          {},
		PredefineTagOceanBase:      {},
	}
}

// 在这里一个tag set 表示了这个规则的标签集合，每个子集是一个SQL的操作对象以及直接作用于其上的操作或修饰词
func GetTagMapDefaultRuleKnowledge() map[string] /* rule_name */ [][]TypeTag /* tag set */ {
	return map[string] /* rule_name */ [][]TypeTag{
		ai.SQLE00063: {{ALTER, TABLE}, {ADD, COLUMN}, {RENAME, COLUMN}},
		ai.SQLE00067: {{CREATE, TABLE}, {FOREIGN, KEY}, {REFERENCES, COLUMN}},
		ai.SQLE00019: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}},
		ai.SQLE00064: {{CREATE, TABLE}, {CREATE, INDEX}, {ALTER, TABLE}, {ADD, INDEX}},
		ai.SQLE00107: {{SELECT, COLUMN}, {LENGTH, STATEMENT}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00003: {{CREATE, TABLE}, {CREATE, INDEX}, {ALTER, TABLE}, {ADD, INDEX}},
		ai.SQLE00016: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}},
		ai.SQLE00072: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {CONSTRAINT}},
		ai.SQLE00082: {{EXPLAIN, STATEMENT}, {SELECT, STATEMENT}},
		ai.SQLE00061: {{CREATE, TABLE}, {IF, TABLE, NOT, EXISTS}},
		ai.SQLE00032: {{CREATE, DATABASE}},
		ai.SQLE00012: {{CREATE, TABLE}, {COLUMN, DECIMAL}},
		ai.SQLE00128: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {HAVING, STATEMENT}},
		ai.SQLE00079: {{SELECT, COLUMN}, {FROM, TABLE}, {TABLE, AS}},
		ai.SQLE00035: {{CREATE, TABLE}, {COLUMN, CONSTRAINT}},
		ai.SQLE00025: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}},
		ai.SQLE00097: {{SELECT, COLUMN}, {FROM, TABLE}, {ORDER, BY, COLUMN}},
		ai.SQLE00023: {{CREATE, TABLE}, {PRIMARY, KEY}},
		ai.SQLE00111: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {FUNCTION}},
		ai.SQLE00031: {{CREATE, VIEW}, {AS, SELECT, STATEMENT}},
		ai.SQLE00083: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {GROUP, BY, COLUMN}},
		ai.SQLE00094: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {FUNCTION}},
		ai.SQLE00151: {{CREATE, TABLE}, {ALTER, TABLE}},
		ai.SQLE00039: {{CREATE, INDEX}, {ALTER, TABLE}, {ADD, INDEX}},
		ai.SQLE00143: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN, OR}},
		ai.SQLE00055: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, INDEX}},
		ai.SQLE00058: {{CREATE, TABLE}, {ALTER, TABLE}, {PARTITION, BY}},
		ai.SQLE00049: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, COLUMN}},
		ai.SQLE00092: {{DELETE, ROW}, {UPDATE, ROW}, {SET, COLUMN}},
		ai.SQLE00071: {{ALTER, TABLE}, {DROP, COLUMN}},
		ai.SQLE00132: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {SUBQUERY}},
		ai.SQLE00115: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {SUBQUERY}},
		ai.SQLE00100: {{SELECT, COLUMN}, {FROM, TABLE}, {ROW, LIMIT}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00220: {{SELECT, COLUMN}, {COUNT, COLUMN}},
		ai.SQLE00075: {{CREATE, TABLE}, {ALTER, TABLE, CHARSET}},
		ai.SQLE00098: {{SELECT, COLUMN}, {FROM, TABLE}, {JOIN, TABLE}},
		ai.SQLE00004: {{CREATE, TABLE}, {COLUMN, AUTO_INCREMENT}},
		ai.SQLE00008: {{CREATE, TABLE}, {PRIMARY, KEY}, {ALTER, TABLE}, {DROP, PRIMARY, KEY}},
		ai.SQLE00015: {{CREATE, TABLE, COLLATION}, {ALTER, TABLE, COLLATION}},
		ai.SQLE00011: {{ALTER, TABLE}, {MODIFY, COLUMN}},
		ai.SQLE00029: {{CREATE, PROCEDURE}, {AS, SELECT, STATEMENT}},
		ai.SQLE00010: {{ALTER, TABLE}, {DROP, PRIMARY, KEY}},
		ai.SQLE00022: {{CREATE, INDEX}, {ALTER, TABLE}, {ADD, INDEX}},
		ai.SQLE00048: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, COLUMN}},
		ai.SQLE00174: {{GRANT, PRIVILEGE, ON, TABLE}, {TO, USER}},
		ai.SQLE00119: {{SELECT, COLUMN}, {FROM, TABLE}, {GROUP, BY, COLUMN}},
		ai.SQLE00037: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, INDEX}},
		ai.SQLE00084: {{CREATE, TEMPORARY, TABLE}, {SELECT, STATEMENT}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00088: {{INSERT, INTO, TABLE}, {INSERT, VALUES}},
		ai.SQLE00126: {{SELECT, COLUMN}, {FROM, TABLE}, {GROUP, BY, COLUMN}},
		ai.SQLE00062: {{SET, TRANSACTION, ISOLATION, LEVEL}},
		ai.SQLE00051: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, AUTO_INCREMENT}},
		ai.SQLE00021: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}},
		ai.SQLE00027: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN, COMMENT}},
		ai.SQLE00074: {{ALTER, TABLE}, {RENAME, TO}},
		ai.SQLE00020: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, COLUMN}},
		ai.SQLE00087: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN, IN}},
		ai.SQLE00180: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00041: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, UNIQUE, INDEX}},
		ai.SQLE00086: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN, LIKE}},
		ai.SQLE00140: {{SELECT, COLUMN}, {FROM, TABLE}},
		ai.SQLE00131: {{SELECT, COLUMN}, {FROM, TABLE}, {ORDER, BY, COLUMN, RAND}},
		ai.SQLE00057: {{CREATE, TABLE, ENGINE}},
		ai.SQLE00218: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00009: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {FUNCTION}},
		ai.SQLE00017: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, BLOB}},
		ai.SQLE00120: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {IN, COLUMN, VALUES}, {VALUES, NULL}},
		ai.SQLE00099: {{SELECT, COLUMN}, {FROM, TABLE}, {FOR, UPDATE}},
		ai.SQLE00127: {{SELECT, COLUMN}, {FROM, TABLE}, {ORDER, BY, COLUMN}, {FUNCTION}},
		ai.SQLE00043: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, INDEX}},
		ai.SQLE00090: {{SELECT, COLUMN}, {FROM, TABLE}, {UNION, TABLE}},
		ai.SQLE00112: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00177: {{SELECT, COLUMN}, {FROM, TABLE}, {ORDER, BY, COLUMN}},
		ai.SQLE00078: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {FUNCTION}},
		ai.SQLE00101: {{SELECT, COLUMN}, {FROM, TABLE}, {ORDER, BY, COLUMN}},
		ai.SQLE00122: {{SELECT, COLUMN}, {COUNT, ROW}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {VALUES, IS, NULL}},
		ai.SQLE00045: {{SELECT, COLUMN}, {FROM, TABLE}, {ROW, LIMIT}, {ROW, OFFSET}},
		ai.SQLE00113: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {COLUMN, NOT, IN, LIKE}, {VALUES, NOT, IN, LIKE}},
		ai.SQLE00068: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, TIMESTAMP}},
		ai.SQLE00007: {{CREATE, TABLE}, {COLUMN, AUTO_INCREMENT}},
		ai.SQLE00178: {{SELECT, COLUMN}, {FROM, TABLE}, {ORDER, BY, COLUMN}},
		ai.SQLE00085: {{SELECT, COLUMN}, {FROM, TABLE}},
		ai.SQLE00018: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, CHAR}},
		ai.SQLE00056: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {CHARSET}},
		ai.SQLE00091: {{SELECT, COLUMN}, {FROM, TABLE}, {JOIN, TABLE}, {ON, STATEMENT}},
		ai.SQLE00108: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {WHERE, SUBQUERY}},
		ai.SQLE00153: {{CREATE, TABLE}},
		ai.SQLE00170: {{ALTER, TABLE}, {MODIFY, COLUMN}},
		ai.SQLE00095: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00030: {{CREATE, TRIGGER, ON}, {ON, TABLE}, {FOR, EACH, ROW}, {SELECT, STATEMENT}},
		ai.SQLE00110: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}, {INDEX}},
		ai.SQLE00134: {{UPDATE, ROW}, {SET}},
		ai.SQLE00176: {{SELECT, COLUMN}, {FROM, TABLE}, {HINT, INDEX}, {HINT, JOIN, TABLE}},
		ai.SQLE00047: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, COLUMN}},
		ai.SQLE00219: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, TIMESTAMP}},
		ai.SQLE00118: {{DROP, TABLE}, {TRUNCATE, TABLE}},
		ai.SQLE00046: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, COLUMN}},
		ai.SQLE00066: {{ALTER, TABLE}, {DROP, COLUMN}, {MODIFY, COLUMN}},
		ai.SQLE00076: {{UPDATE, ROW}, {SET}, {DELETE, ROW}},
		ai.SQLE00161: {{SET, GLOBAL, VALUES}, {COLUMN, AUTO_INCREMENT}},
		ai.SQLE00026: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, INT}},
		ai.SQLE00042: {{CREATE, TEMPORARY, TABLE}, {ALTER, TABLE}, {RENAME, TO}},
		ai.SQLE00052: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, AUTO_INCREMENT}},
		ai.SQLE00089: {{INSERT, INTO, TABLE}, {SELECT, STATEMENT}},
		ai.SQLE00034: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, NOT_NULL}, {VALUES, NOT_NULL}},
		ai.SQLE00139: {{SELECT, COLUMN}, {FROM, TABLE}},
		ai.SQLE00123: {{TRUNCATE, TABLE}},
		ai.SQLE00014: {{CREATE, FUNCTION}, {TABLE, AS}, {SELECT, STATEMENT}},
		ai.SQLE00040: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, INDEX}},
		ai.SQLE00059: {{ALTER, TABLE}, {MODIFY, COLUMN}},
		ai.SQLE00096: {{SELECT, COLUMN}, {FROM, TABLE}, {JOIN, TABLE}},
		ai.SQLE00053: {{SELECT, COLUMN}, {FROM, TABLE}},
		ai.SQLE00080: {{INSERT, INTO, TABLE}, {VALUES}, {SELECT, STATEMENT}},
		ai.SQLE00054: {{CREATE, TABLE}, {COLUMN, BIGINT}},
		ai.SQLE00073: {{ALTER, TABLE}, {MODIFY, COLUMN}},
		ai.SQLE00109: {{ALTER, TABLE}, {CHARSET}},
		ai.SQLE00102: {{DELETE, ROW}, {ORDER, BY, COLUMN}},
		ai.SQLE00013: {{CREATE, TABLE}, {ALTER, TABLE}, {MODIFY, COLUMN}, {COLUMN, FLOAT}},
		ai.SQLE00060: {{CREATE, TABLE}, {COMMENT}},
		ai.SQLE00002: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00124: {{DELETE, ROW}},
		ai.SQLE00001: {{SELECT, COLUMN}, {FROM, TABLE}, {WHERE, COLUMN}, {WHERE, VALUES}},
		ai.SQLE00175: {{SELECT, COLUMN}, {FROM, TABLE}},
		ai.SQLE00005: {{CREATE, TABLE}, {ALTER, TABLE}, {ADD, INDEX}},
	}
}

// 创建知识库
func (s *Storage) CreateKnowledgeWithTags(knowledge *Knowledge, ruleName string, tags map[TypeTag] /* tag name */ *Tag, filterTags []*Tag) (*Knowledge, error) {
	// 先查询是否存在同名的知识库
	modelKnowledge, err := s.GetKnowledgeByTagsAndRuleName(filterTags, ruleName)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	// 标签映射
	modelTagMap := make(map[TypeTag] /* tag name */ *Tag)

	// 若modelKnowledge的Content为空，则更新Content
	if modelKnowledge != nil && modelKnowledge.Content == "" && knowledge.Content != "" {
		err = s.UpdateKnowledgeContent(modelKnowledge, knowledge.Content)
		if err != nil {
			return nil, err
		}
	}

	// 如果不存在，则创建
	if modelKnowledge == nil {
		err = s.db.Omit("Tags").Create(knowledge).Error
		if err != nil {
			return nil, err
		}
		modelKnowledge = knowledge
	}

	for _, tag := range modelKnowledge.Tags {
		modelTagMap[tag.Name] = tag
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
func (s *Storage) GetKnowledgeByTagsAndRuleName(filterTags []*Tag, ruleName string) (*Knowledge, error) {
	var (
		modelKnowledge Knowledge
		tagIds         = make([]uint, 0, len(filterTags))
	)

	for _, tag := range filterTags {
		tagIds = append(tagIds, tag.ID)
	}

	err := s.db.Model(&Knowledge{}).Preload(`Tags`).
		Joins(`
			JOIN knowledge_tag_relations ktr ON knowledge.id = ktr.knowledge_id 
			JOIN tags t ON ktr.tag_id = t.id
			JOIN rule_knowledge_relations rkr ON knowledge.id = rkr.knowledge_id
		`).
		Where(`rkr.rule_name = ? AND t.id IN ?`, ruleName, tagIds).
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
		likeClause := "%%" + keyword + "%%"
		searchClause = searchClause.
			Select(`knowledge.id,SUBSTRING(content, LOCATE(?, content), 50) AS content, description, title`, keyword).
			Where("title LIKE ? OR description LIKE ? OR MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)", likeClause, likeClause, keyword)
		countClause = countClause.
			Where("title LIKE ? OR description LIKE ? OR MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)", likeClause, likeClause, keyword)
	}
	// 如果有标签，则根据标签查询带有标签的知识库
	if len(tags) > 0 {
		// 使用 HAVING 子句确保每个知识记录都包含所有指定的标签
		searchClause = searchClause.Group("knowledge.id").Having("SUM(CASE WHEN tags.name IN ? THEN 1 ELSE 0 END) = ?", tags, len(tags))
		countClause = countClause.Group("knowledge.id").Having("SUM(CASE WHEN tags.name IN ? THEN 1 ELSE 0 END) = ?", tags, len(tags))
	} else {
		searchClause = searchClause.Group("knowledge.id")
		countClause = countClause.Group("knowledge.id")
	}

	// 分页查询
	err := searchClause.Limit(limit).Offset(offset).Find(&results).Error
	if err != nil {
		return nil, 0, err
	}
	// 统计查询结果总数
	var count int64
	err = countClause.Distinct("knowledge.id").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

// 更新知识库内容
func (s *Storage) UpdateKnowledgeContent(knowledge *Knowledge, newContent string) error {
	return s.db.Model(knowledge).Updates(map[string]interface{}{"content": newContent}).Error
}

type RuleKnowledgeRelation struct {
	KnowledgeID uint64 `gorm:"primaryKey;autoIncrement:false"`
	RuleName    string `gorm:"primaryKey;size:255"`
	RuleDBType  string `gorm:"primaryKey;size:255"`
}

func (RuleKnowledgeRelation) TableName() string {
	return "rule_knowledge_relations"
}

type CustomRuleKnowledgeRelation struct {
	KnowledgeID  uint64 `gorm:"primaryKey;autoIncrement:false"`
	CustomRuleID uint64 `gorm:"primaryKey;autoIncrement:false"`
}

func (CustomRuleKnowledgeRelation) TableName() string {
	return "custom_rule_knowledge_relations"
}

// get rule knowledge relations
func (s *Storage) GetRuleKnowledgeRelationByKnowledgeID(knowledgeID uint) (*RuleKnowledgeRelation, error) {
	var relation *RuleKnowledgeRelation
	err := s.db.Where("knowledge_id = ?", knowledgeID).First(&relation).Error
	if err != nil {
		return nil, err
	}
	return relation, nil
}

// 创建规则和知识库的关联
func (s *Storage) CreateRuleKnowledgeRelation(knowledgeId uint64, ruleName string, ruleDBType string) error {
	if knowledgeId == 0 || ruleName == "" || ruleDBType == "" {
		return nil
	}
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "knowledge_id"}, {Name: "rule_name"}, {Name: "rule_db_type"}},
		DoNothing: true, // 如果存在则不插入也不更新
	}).Create(&RuleKnowledgeRelation{
		KnowledgeID: knowledgeId,
		RuleName:    ruleName,
		RuleDBType:  ruleDBType,
	}).Error
}

// 创建自定义规则和知识库的关联
func (s *Storage) CreateCustomRuleKnowledgeRelation(knowledgeId uint64, customRule *CustomRule) error {
	if knowledgeId == 0 || customRule == nil {
		return nil
	}
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "knowledge_id"}, {Name: "custom_rule_id"}},
		DoNothing: true, // 如果存在则不插入也不更新
	}).Create(&CustomRuleKnowledgeRelation{
		KnowledgeID:  knowledgeId,
		CustomRuleID: uint64(customRule.ID),
	}).Error
}

// 根据上下文获取对应语言的知识
func (m MultiLanguageKnowledge) GetKnowledgeByLang(lang language.Tag) *Knowledge {
	var defaultLangKnowledge *Knowledge
	for _, k := range m {
		for _, tag := range k.Tags {
			if tag.Name == TypeTag(lang.String()) {
				return k
			}
			if tag.Name == PredefineTagChinese {
				defaultLangKnowledge = k
			}
		}
	}
	if defaultLangKnowledge != nil {
		return defaultLangKnowledge
	}
	return &Knowledge{}
}
