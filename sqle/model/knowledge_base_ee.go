package model

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
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
	LIMIT                              TypeTag = "LIMIT"
	CREATE_FUNCTION                    TypeTag = "CREATE_FUNCTION"
	ALTER_TABLE                        TypeTag = "ALTER_TABLE"
	DDL                                TypeTag = "DDL"
	UNIQUE_INDEX                       TypeTag = "UNIQUE_INDEX"
	ALTER                              TypeTag = "ALTER"
	CREATE                             TypeTag = "CREATE"
	RAND                               TypeTag = "RAND"
	ORDER_BY                           TypeTag = "ORDER_BY"
	CREATE_DATABASE                    TypeTag = "CREATE_DATABASE"
	ALTER_DATABASE                     TypeTag = "ALTER_DATABASE"
	CHARSET                            TypeTag = "CHARSET"
	TEMPORARY                          TypeTag = "TEMPORARY"
	DML                                TypeTag = "DML"
	COMMENT                            TypeTag = "COMMENT"
	EXPLAIN_FORMAT_TRADITIONAL         TypeTag = "EXPLAIN_FORMAT_TRADITIONAL"
	TIMESTAMP                          TypeTag = "TIMESTAMP"
	DATETIME                           TypeTag = "DATETIME"
	INSERT_VALUES                      TypeTag = "INSERT_VALUES"
	REPLACE_VALUES                     TypeTag = "REPLACE_VALUES"
	INSERT_SELECT                      TypeTag = "INSERT_SELECT"
	REPLACE_SELECT                     TypeTag = "REPLACE_SELECT"
	ADD                                TypeTag = "ADD"
	ALTER_TABLE_MODIFY                 TypeTag = "ALTER_TABLE_MODIFY"
	ALTER_TABLE_CHANGE                 TypeTag = "ALTER_TABLE_CHANGE"
	ALTER_TABLE_ADD                    TypeTag = "ALTER_TABLE_ADD"
	FIRST                              TypeTag = "FIRST"
	AFTER                              TypeTag = "AFTER"
	DROP                               TypeTag = "DROP"
	EXISTS                             TypeTag = "EXISTS"
	IN                                 TypeTag = "IN"
	COALESCE                           TypeTag = "COALESCE"
	OR                                 TypeTag = "OR"
	TRUNCATE                           TypeTag = "TRUNCATE"
	NOT_NULL                           TypeTag = "NOT_NULL"
	DEFAULT                            TypeTag = "DEFAULT"
	AUTO_INCREMENT                     TypeTag = "AUTO_INCREMENT"
	INSERT                             TypeTag = "INSERT"
	COUNT_STAR                         TypeTag = "COUNT_STAR"
	COUNT_1                            TypeTag = "COUNT_1"
	SUM                                TypeTag = "SUM"
	COUNT                              TypeTag = "COUNT"
	EXPLAIN_FORMAT_JSON                TypeTag = "EXPLAIN_FORMAT_JSON"
	JOIN                               TypeTag = "JOIN"
	FUNCTION                           TypeTag = "FUNCTION"
	VIEW                               TypeTag = "VIEW"
	FOREIGN_KEY                        TypeTag = "FOREIGN_KEY"
	PARTITION                          TypeTag = "PARTITION"
	CREATE_TRIGGER                     TypeTag = "CREATE_TRIGGER"
	HINT                               TypeTag = "HINT"
	FORCE_INDEX                        TypeTag = "FORCE_INDEX"
	USE_INDEX                          TypeTag = "USE_INDEX"
	IGNORE_INDEX                       TypeTag = "IGNORE_INDEX"
	STRAIGHT_JOIN                      TypeTag = "STRAIGHT_JOIN"
	DISTINCT                           TypeTag = "DISTINCT"
	GROUP_BY                           TypeTag = "GROUP_BY"
	UNION                              TypeTag = "UNION"
	IN_NULL                            TypeTag = "IN_NULL"
	NOT_IN_NULL                        TypeTag = "NOT_IN_NULL"
	RENAME                             TypeTag = "RENAME"
	DROP_CONSTRAINT                    TypeTag = "DROP_CONSTRAINT"
	SELECT_STAR                        TypeTag = "SELECT_STAR"
	CREATE_PROCEDURE                   TypeTag = "CREATE_PROCEDURE"
	ALTER_PROCEDURE                    TypeTag = "ALTER_PROCEDURE"
	CREATE_TABLE                       TypeTag = "CREATE_TABLE"
	CREATE_INDEX                       TypeTag = "CREATE_INDEX"
	ALTER_TABLE_ADD_INDEX              TypeTag = "ALTER_TABLE_ADD_INDEX"
	SUBQUERY                           TypeTag = "SUBQUERY"
	GRANT                              TypeTag = "GRANT"
	UNION_ALL                          TypeTag = "UNION_ALL"
	LIKE                               TypeTag = "LIKE"
	WITH                               TypeTag = "WITH"
	SET_TRANSACTION                    TypeTag = "SET_TRANSACTION"
	CREATE_VIEW                        TypeTag = "CREATE_VIEW"
	TRANSACTION_ISOLATION              TypeTag = "TRANSACTION_ISOLATION"
	ALTER_TABLE_ADD_KEY                TypeTag = "ALTER_TABLE_ADD_KEY"
	FLOAT                              TypeTag = "FLOAT"
	DOUBLE                             TypeTag = "DOUBLE"
	TIME                               TypeTag = "TIME"
	YEAR                               TypeTag = "YEAR"
	DATE                               TypeTag = "DATE"
	SET                                TypeTag = "SET"
	ALTER_TABLE_ADD_PRIMARY_KEY        TypeTag = "ALTER_TABLE_ADD_PRIMARY_KEY"
	// TODO 增加新的标签
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

func GetTagMapDefaultRuleKnowledge() map[string] /* rule_name */ []TypeTag {
	return map[string] /* rule_name */ []TypeTag{
		ai.SQLE00014: {CREATE_FUNCTION},
		ai.SQLE00092: {DELETE, UPDATE, LIMIT},
		ai.SQLE00109: {UPDATE, DELETE, SELECT, LIMIT},
		ai.SQLE00086: {SELECT, INSERT, UPDATE, DELETE, UNION_ALL, LIKE},
		ai.SQLE00035: {CREATE_TABLE, ALTER_TABLE, DDL},
		ai.SQLE00063: {ALTER, CREATE, UNIQUE_INDEX},
		ai.SQLE00131: {SELECT, ORDER_BY, RAND},
		ai.SQLE00056: {CREATE_TABLE, ALTER_TABLE, CREATE_DATABASE, ALTER_DATABASE, CHARSET},
		ai.SQLE00042: {CREATE_TABLE, ALTER_TABLE, TEMPORARY},
		// ai.SQLE00104: {SELECT, ORDER_BY, INSERT, SELECT_FROM},
		ai.SQLE00094: {DML},
		ai.SQLE00027: {CREATE_TABLE, ALTER_TABLE, COMMENT},
		ai.SQLE00175: {SELECT, EXPLAIN_FORMAT_TRADITIONAL},
		ai.SQLE00033: {CREATE_TABLE, ALTER_TABLE, TIMESTAMP, DATETIME},
		ai.SQLE00080: {INSERT_VALUES, REPLACE_VALUES, INSERT_SELECT, REPLACE_SELECT},
		ai.SQLE00049: {CREATE_TABLE, ALTER_TABLE, ADD, RENAME},
		// ai.SQLE00006: {SELECT, JOIN},
		ai.SQLE00102: {UPDATE, DELETE, ORDER_BY},
		ai.SQLE00123: {TRUNCATE},
		ai.SQLE00065: {ALTER_TABLE_MODIFY, ALTER_TABLE_CHANGE, ALTER_TABLE_ADD, FIRST, AFTER},
		ai.SQLE00071: {ALTER_TABLE, DROP},
		ai.SQLE00001: {SELECT, INSERT, UPDATE, DELETE, WHERE, EXISTS, IN, COALESCE, OR},
		ai.SQLE00153: {CREATE_TABLE},
		ai.SQLE00087: {SELECT, IN, WITH, INSERT_SELECT, UPDATE, DELETE},
		ai.SQLE00108: {SELECT, INSERT, UPDATE, DELETE, SUBQUERY},
		ai.SQLE00034: {CREATE_TABLE, ALTER_TABLE, NOT_NULL, DEFAULT},
		ai.SQLE00020: {CREATE_TABLE, ALTER_TABLE},
		ai.SQLE00076: {UPDATE, DELETE},
		ai.SQLE00124: {DELETE},
		ai.SQLE00062: {SET_TRANSACTION, TRANSACTION_ISOLATION},
		ai.SQLE00043: {CREATE_INDEX, ALTER_TABLE_ADD_INDEX},
		ai.SQLE00039: {CREATE_INDEX, ALTER_TABLE_ADD_KEY, ALTER_TABLE_ADD_INDEX, SELECT},
		ai.SQLE00111: {SELECT, INSERT, UPDATE, DELETE, WHERE},
		ai.SQLE00032: {CREATE_DATABASE},
		ai.SQLE00174: {GRANT},
		ai.SQLE00048: {CREATE_TABLE, ALTER_TABLE, CREATE_DATABASE, ALTER_DATABASE},
		ai.SQLE00013: {CREATE_TABLE, ALTER_TABLE, FLOAT, DOUBLE},
		ai.SQLE00179: {SELECT, UPDATE, DELETE, INSERT},
		ai.SQLE00051: {CREATE_TABLE, ALTER_TABLE, AUTO_INCREMENT},
		ai.SQLE00098: {SELECT, UPDATE, DELETE, INSERT},
		ai.SQLE00064: {CREATE_TABLE, CREATE_INDEX, ALTER_TABLE_ADD_INDEX},
		ai.SQLE00122: {SELECT, SUM, COUNT},
		ai.SQLE00023: {CREATE_TABLE, ALTER_TABLE_ADD_PRIMARY_KEY},
		ai.SQLE00084: {CREATE, ALTER_TABLE, TEMPORARY},
		ai.SQLE00037: {CREATE_TABLE, ALTER_TABLE_ADD_INDEX, CREATE_INDEX},
		ai.SQLE00220: {SELECT, COUNT_STAR, COUNT_1},
		ai.SQLE00040: {CREATE_INDEX, ALTER_TABLE_ADD_INDEX, RENAME},
		ai.SQLE00061: {CREATE_TABLE},
		ai.SQLE00127: {SELECT, ORDER_BY, INSERT_SELECT},
		ai.SQLE00180: {SELECT, EXPLAIN_FORMAT_JSON},
		ai.SQLE00075: {CREATE_TABLE, ALTER_TABLE, CHARSET},
		ai.SQLE00004: {CREATE_TABLE, SET},
		ai.SQLE00096: {SELECT, JOIN},
		ai.SQLE00025: {CREATE_TABLE, ALTER_TABLE, DATE, TIMESTAMP, DATETIME, TIME, YEAR},
		ai.SQLE00177: {SELECT, ORDER_BY},
		ai.SQLE00031: {CREATE_VIEW},
		ai.SQLE00082: {SELECT, EXPLAIN_FORMAT_JSON},
		ai.SQLE00121: {SELECT, ORDER_BY, LIMIT},
		ai.SQLE00067: {CREATE_TABLE, ALTER_TABLE, FOREIGN_KEY},
		ai.SQLE00073: {ALTER_TABLE, CHARSET},
		ai.SQLE00009: {SELECT, WHERE, FUNCTION},
		ai.SQLE00052: {CREATE_TABLE, ALTER_TABLE, AUTO_INCREMENT},
		ai.SQLE00100: {SELECT, EXPLAIN_FORMAT_JSON},
		ai.SQLE00046: {CREATE_TABLE, ALTER_TABLE},
		ai.SQLE00170: {ALTER_TABLE_MODIFY, ALTER_TABLE_CHANGE},
		ai.SQLE00022: {CREATE_INDEX, ALTER_TABLE_ADD_KEY, ALTER_TABLE_ADD_INDEX, SELECT},
		ai.SQLE00091: {SELECT, JOIN},
		ai.SQLE00058: {CREATE_TABLE, ALTER_TABLE, PARTITION},
		ai.SQLE00003: {CREATE_TABLE, CREATE_INDEX, ALTER_TABLE_ADD_INDEX},
		ai.SQLE00041: {CREATE_INDEX, ALTER_TABLE_ADD_INDEX, RENAME},
		ai.SQLE00088: {INSERT},
		ai.SQLE00219: {CREATE_TABLE, ALTER_TABLE, TIMESTAMP},
		ai.SQLE00055: {CREATE_INDEX, ALTER_TABLE_ADD_INDEX},
		ai.SQLE00132: {SELECT, INSERT, UPDATE, DELETE, UNION_ALL},
		ai.SQLE00139: {SELECT, UPDATE, DELETE, INSERT, EXPLAIN_FORMAT_JSON},
		ai.SQLE00143: {SELECT, INSERT, DELETE, UPDATE, UNION, OR},
		ai.SQLE00030: {CREATE_TRIGGER},
		ai.SQLE00176: {SELECT, HINT, FORCE_INDEX, USE_INDEX, IGNORE_INDEX, STRAIGHT_JOIN},
		ai.SQLE00083: {SELECT, ORDER_BY, DISTINCT, GROUP_BY, UNION, EXPLAIN_FORMAT_JSON},
		ai.SQLE00097: {UPDATE, DELETE, INSERT_SELECT, SELECT, ORDER_BY, DISTINCT, GROUP_BY, UNION},
		ai.SQLE00134: {UPDATE},
		ai.SQLE00072: {ALTER_TABLE, DROP_CONSTRAINT},
		ai.SQLE00120: {SELECT, IN_NULL, NOT_IN_NULL},
		ai.SQLE00047: {CREATE_TABLE, ALTER_TABLE, RENAME},
		ai.SQLE00101: {SELECT, ORDER_BY, INSERT_SELECT},
		ai.SQLE00115: {SELECT, SUBQUERY},
		ai.SQLE00053: {SELECT_STAR},
		ai.SQLE00029: {CREATE_PROCEDURE, ALTER_PROCEDURE},
		// TODO 增加新的规则和标签的映射
	}
}

// 创建知识库
func (s *Storage) CreateKnowledgeWithTags(knowledge *Knowledge, tags map[TypeTag] /* tag name */ *Tag, filterTags []*Tag) (*Knowledge, error) {
	// 先查询是否存在同名的知识库
	modelKnowledge, err := s.GetKnowledgeByTagsAndRuleName(filterTags, knowledge.RuleName)
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
		Joins(`JOIN knowledge_tag_relations ktr ON knowledge.id = ktr.knowledge_id JOIN tags t ON ktr.tag_id = t.id`).
		Where(`t.id IN ? AND knowledge.rule_name = ?`, tagIds, ruleName).
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
		Select(`knowledge.id, rule_name, description, title`).
		Joins("LEFT JOIN knowledge_tag_relations ON knowledge_tag_relations.knowledge_id = knowledge.id LEFT JOIN tags ON tags.id = knowledge_tag_relations.tag_id")

	countClause := s.db.
		Table("knowledge").
		Joins("LEFT JOIN knowledge_tag_relations ON knowledge_tag_relations.knowledge_id = knowledge.id LEFT JOIN tags ON tags.id = knowledge_tag_relations.tag_id")
	// 如果有关键字，则根据关键字进行模糊查询+全文检索
	if len(keyword) > 0 {
		likeClause := "%%" + keyword + "%%"
		searchClause = searchClause.
			Select(`knowledge.id, rule_name, SUBSTRING(content, LOCATE(?, content), 50) AS content, description, title`, keyword).
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
