//go:build enterprise
// +build enterprise

package knowledge_base

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/license"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"golang.org/x/text/language"
)

// 迁移规则知识库到知识库，并且关联标签
func MigrateKnowledgeFromRules(rulesMap map[string] /* DBType */ []*model.Rule) error {
	if license.CheckKnowledgeBaseLicense(config.GetOptions().SqleOptions.KnowledgeBaseTempLicense) == nil {
		// 获取所有规则的知识
		defaultAIRuleKnowledge, err := ai.GetAIRulesKnowledge()
		if err != nil {
			panic(fmt.Errorf("get ai rules knowledge failed: %v", err))
		}
		// 初始化规则的知识
		for dbType, ruleList := range rulesMap {
			for idx, rule := range ruleList {
				if knowledge, exist := defaultAIRuleKnowledge[rule.Name]; exist {
					rulesMap[dbType][idx].Knowledge.I18nContent[language.Chinese] = knowledge
				}
			}
		}
	}
	predefineTags, err := NewTagService(model.GetStorage()).GetOrCreatePredefinedTags()
	if err != nil {
		log.Logger().Errorf("get or create predefined tags failed: %v", err)
		return err
	}
	err = createKnowledgeWithTag(rulesMap, predefineTags)
	if err != nil {
		log.Logger().Errorf("create knowledge failed: %v", err)
		return err
	}
	return nil
}

// 根据规则创建知识库并关联标签
func createKnowledgeWithTag(rulesMap map[string][]*model.Rule, predefineTags map[model.TypeTag]*model.Tag) error {
	storage := model.GetStorage()

	// 处理默认规则
	for _, rules := range rulesMap {
		if err := processRules(storage, rules, predefineTags, getDefaultRuleKnowledgeWithTags); err != nil {
			return err
		}
	}

	// 处理自定义规则
	customRules, err := storage.GetAllCustomRules()
	if err != nil {
		return fmt.Errorf("failed to get custom rules: %w", err)
	}
	return processRules(storage, customRules, predefineTags, getCustomRuleKnowledgeWithTags)
}

// 使用泛型来处理不同类型的规则
func processRules[T any](
	storage *model.Storage,
	rules []T,
	predefineTags map[model.TypeTag]*model.Tag,
	getKnowledgeWithTags func(T, map[model.TypeTag]*model.Tag) ([]*KnowledgeWithFilter, error),
) error {
	for _, rule := range rules {
		knowledgeWithTags, err := getKnowledgeWithTags(rule, predefineTags)
		if err != nil {
			return fmt.Errorf("failed to get knowledge for rule: %w", err)
		}
		for _, item := range knowledgeWithTags {
			if _, err := storage.CreateKnowledgeWithTags(item.knowledge, item.tagMap, item.filterTags); err != nil {
				return fmt.Errorf("failed to create knowledge: %w", err)
			}
		}
	}
	return nil
}

// 初始化一条自定义规则的知识
func InitCustomRuleKnowledge(rule *model.CustomRule) error {
	predefineTags, err := NewTagService(model.GetStorage()).GetOrCreatePredefinedTags()
	if err != nil {
		log.Logger().Errorf("get or create predefined tags failed: %v", err)
		return err
	}
	return processRules(model.GetStorage(), []*model.CustomRule{rule}, predefineTags, getCustomRuleKnowledgeWithTags)
}

// 获取所有知识库标签
func GetKnowledgeBaseTags() ([]*model.Tag, error) {
	s := model.GetStorage()
	// 获取标签：知识库预定义标签
	modelPredefineTag, err := s.GetTagByName(model.PredefineTagKnowledgeBase)
	if err != nil {
		return nil, err
	}
	// 获取所有知识库预定义标签
	modelKnowledgeTags, err := s.GetSubTags(modelPredefineTag.ID)
	if err != nil {
		return nil, err
	}
	return modelKnowledgeTags, nil
}

func UpdateRuleKnowledgeContent(ctx context.Context, ruleName, dbType, content string) error {
	return updateRuleKnowledgeContent(ctx, ruleName, dbType, content, model.PredefineTagKnowledgeBase)
}

func UpdateCustomRuleKnowledgeContent(ctx context.Context, ruleName, dbType, content string) error {
	return updateRuleKnowledgeContent(ctx, ruleName, dbType, content, model.PredefineTagCustomizeKnowledgeBase)
}

func updateRuleKnowledgeContent(ctx context.Context, ruleName, dbType, content string, ruleType model.TypeTag) error {
	s := model.GetStorage()
	tagFilter, err := getKnowledgeDefaultTag(ctx, dbType, ruleType)
	if err != nil {
		log.Logger().Errorf("get one predefine knowledge default tag failed, err: %v", err)
		return err
	}
	// 获取知识库
	knowledge, err := s.GetKnowledgeByTagsAndRuleName(tagFilter, ruleName)
	if err != nil {
		log.Logger().Errorf("get knowledge by langTag and title failed, err: %v", err)
		return err
	}
	err = s.UpdateKnowledgeContent(knowledge, content)
	if err != nil {
		log.Logger().Errorf("update knowledge content failed, err: %v", err)
		return err
	}

	return nil
}

// 获取一条规则的知识库
func GetRuleWithKnowledge(ctx context.Context, ruleName, dbType string) (*model.Knowledge, error) {
	return getRuleWithKnowledge(ctx, ruleName, dbType, model.PredefineTagKnowledgeBase)
}

// 获取一条自定义规则的知识库
func GetCustomRuleWithKnowledge(ctx context.Context, ruleName, dbType string) (*model.Knowledge, error) {
	return getRuleWithKnowledge(ctx, ruleName, dbType, model.PredefineTagCustomizeKnowledgeBase)
}

func getRuleWithKnowledge(ctx context.Context, ruleName, dbType string, ruleType model.TypeTag) (*model.Knowledge, error) {
	s := model.GetStorage()
	tagFilter, err := getKnowledgeDefaultTag(ctx, dbType, ruleType)
	if err != nil {
		log.Logger().Errorf("get one predefine knowledge default tag failed, err: %v", err)
		return nil, err
	}
	// 获取知识库
	knowledge, err := s.GetKnowledgeByTagsAndRuleName(tagFilter, ruleName)
	if err != nil {
		log.Logger().Errorf("get knowledge by langTag and title failed, err: %v", err)
		return nil, err
	}
	return knowledge, nil
}

func getKnowledgeDefaultTag(ctx context.Context, dbType string, ruleType model.TypeTag) ([]*model.Tag, error) {
	var (
		langTag *model.Tag
		dbTag   *model.Tag
		err     error
		lang    = locale.Bundle.GetLangTagFromCtx(ctx)
		s       = model.GetStorage()
	)
	// 获取语言标签
	switch lang {
	case language.Chinese:
		langTag, err = s.GetTagByName(model.PredefineTagChinese)
		if err != nil {
			log.Logger().Errorf("get langTag by name failed, err: %v", err)
			return nil, err
		}
	case language.English:
		langTag, err = s.GetTagByName(model.PredefineTagEnglish)
		if err != nil {
			log.Logger().Errorf("get langTag by name failed, err: %v", err)
			return nil, err
		}
	}
	// 获取数据库类型
	dbTag, err = s.GetTagByName(model.TypeTag(dbType))
	if err != nil {
		log.Logger().Errorf("get dbTag by name failed, err: %v", err)
		return nil, err
	}
	// 获取知识库预定义标签
	knowledgePredefineTag, err := s.GetTagByName(ruleType)
	if err != nil {
		log.Logger().Errorf("get knowledge predefine tag by name failed, err: %v", err)
		return nil, err
	}
	return []*model.Tag{langTag, dbTag, knowledgePredefineTag}, nil
}

// 搜索知识列表
func SearchKnowledgeList(ctx context.Context, keyword string, tags []string, limit, offset int) ([]model.Knowledge, int64, error) {
	s := model.GetStorage()
	// 根据语言和数据库版本过滤
	tags = append(tags, locale.Bundle.GetLangTagFromCtx(ctx).String())
	tags = append(tags, string(model.PredefineTagV2))
	knowledge, count, err := s.SearchKnowledge(keyword, tags, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	// 获取所有tag id的Map最终转化为列表
	tagIDs := make([]uint, 0)
	tagIDMap := make(map[uint] /* tag id */ struct{} /* empty */)
	for _, item := range knowledge {
		for _, tag := range item.Tags {
			tagIDMap[tag.ID] = struct{}{}
		}
	}
	for tagID := range tagIDMap {
		tagIDs = append(tagIDs, tagID)
	}

	// 获取所有tag作为子标签的关联
	tagTagRelations, err := s.GetTagRelationsByTagIds(tagIDs)
	if err != nil {
		return nil, 0, err
	}
	// 获取所有父级标签的Map
	fatherTagCache := make(map[uint] /* father tag id */ *model.Tag /* father tag */)
	for _, tagTagRelation := range tagTagRelations {
		if _, exist := fatherTagCache[tagTagRelation.TagID]; !exist {
			fatherTagCache[tagTagRelation.TagID], err = s.GetTagById(tagTagRelation.TagID)
		}
	}
	// 获取所有子集标签到父级标签的Map
	subTagFatherTagCache := make(map[uint] /* sub tag id */ *model.Tag /* father tag */)
	for _, tagTagRelation := range tagTagRelations {
		subTagFatherTagCache[tagTagRelation.SubTagID] = fatherTagCache[tagTagRelation.TagID]
	}

	// 遍历知识列表，获取标签的父级标签，并根据父级标签/子级标签两层嵌套覆盖知识列表的标签
	for idx, item := range knowledge {
		knowledge[idx].Tags = nil
		// 用于覆盖知识列表标签的新标签Map
		fatherTagMap := make(map[*model.Tag] /* father tag */ []*model.Tag /* sub tag */)
		for _, tag := range item.Tags {
			if fatherTag, exist := subTagFatherTagCache[tag.ID]; exist {
				fatherTagMap[fatherTag] = append(fatherTagMap[fatherTag], tag)
			} else {
				knowledge[idx].Tags = append(knowledge[idx].Tags, tag)
				continue
			}
		}
		// 覆盖知识列表标签
		for fatherTag, subTags := range fatherTagMap {
			newTag := fatherTag
			newTag.SubTag = subTags
			knowledge[idx].Tags = append(knowledge[idx].Tags, newTag)
		}
	}
	return knowledge, count, nil
}
