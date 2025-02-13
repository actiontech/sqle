//go:build enterprise
// +build enterprise

package knowledge_base

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/model"
	"golang.org/x/text/language"
)

/*
	该文件的主要目的是为了统一管理和维护规则到知识的转换逻辑，避免重复代码和不一致的转换规则。
	该文件的结构如下：
		1. RuleInfoProvider 接口定义了获取规则信息的通用行为，包括获取规则名称、描述、内容等。
		2. BaseRuleWrapper 是 RuleInfoProvider 接口的基础实现，提供了获取规则名称、描述、内容等的通用逻辑。
		3. RuleWrapper 是 RuleInfoProvider 接口的实现，用于将 model.Rule 转换为 RuleInfoProvider。
		4. CustomRuleWrapper 是 RuleInfoProvider 接口的实现，用于将 model.CustomRule 转换为 RuleInfoProvider。
		5. getDefaultRuleKnowledgeWithTags 函数用于将默认规则转换为知识。
		6. getCustomRuleKnowledgeWithTags 函数用于将自定义规则转换为知识。
*/

// KnowledgeWithFilter 包含知识和过滤标签
type KnowledgeWithFilter struct {
	knowledge  *model.Knowledge             // 知识
	filterTags []*model.Tag                 // 用于过滤出唯一的知识
	tagMap     map[model.TypeTag]*model.Tag // 包含所有用于创建该知识的标签
}

// BaseRuleWrapper 作为模板方法的基类
type BaseRuleWrapper struct {
	predefineTags map[model.TypeTag]*model.Tag
}

// ToModelKnowledge 是模板方法，包含通用逻辑
func (brw *BaseRuleWrapper) ToModelKnowledge(rule RuleInfoProvider) ([]*KnowledgeWithFilter, error) {
	var result []*KnowledgeWithFilter

	for langTag, lang := range model.GetTagMapPredefineLanguage() {
		// 收集所有标签
		tagMap := make(map[model.TypeTag]*model.Tag)

		// 获取必需的基础标签
		knowledgeTag, i18nTag, dbTag, err := rule.GetRequiredTags(brw.predefineTags, langTag)
		if err != nil {
			return nil, err
		}

		// 添加必需标签
		tagMap[model.PredefineTagKnowledgeBase] = knowledgeTag
		tagMap[model.TypeTag(langTag)] = i18nTag
		tagMap[model.TypeTag(rule.GetDBType())] = dbTag

		// 额外的规则特定标签
		rule.AddExtraTags(tagMap, brw.predefineTags)

		// 创建知识对象
		knowledge := &model.Knowledge{
			Title:       rule.GetRuleName(),
			Description: rule.GetRuleDescription(lang),
			Content:     rule.GetRuleContent(lang),
		}

		// 必需的过滤标签
		filterTags := []*model.Tag{knowledgeTag, i18nTag, dbTag}

		result = append(result, &KnowledgeWithFilter{
			knowledge:  knowledge,
			filterTags: filterTags,
			tagMap:     tagMap,
		})
	}

	return result, nil
}

// RuleInfoProvider 作为策略接口，提供不同规则的获取方式
type RuleInfoProvider interface {
	GetRuleName() string
	GetRuleDescription(lang language.Tag) string
	GetDBType() model.TypeTag
	GetRuleContent(lang language.Tag) string
	GetRequiredTags(predefineTags map[model.TypeTag]*model.Tag, langTag model.TypeTag) (*model.Tag, *model.Tag, *model.Tag, error)
	AddExtraTags(tagMap map[model.TypeTag]*model.Tag, predefineTags map[model.TypeTag]*model.Tag)
}

// RuleWrapper 实现 RuleInfoProvider
type RuleWrapper struct {
	BaseRuleWrapper
	rule *model.Rule
}

// RuleWrapper 具体实现方法
func (rw *RuleWrapper) GetRuleName() string {
	return rw.rule.Name
}

func (rw *RuleWrapper) GetRuleDescription(lang language.Tag) string {
	return rw.rule.I18nRuleInfo.GetRuleInfoByLangTag(lang).Desc
}

func (rw *RuleWrapper) GetDBType() model.TypeTag {
	return model.TypeTag(rw.rule.DBType)
}

func (rw *RuleWrapper) GetRuleContent(lang language.Tag) string {
	if rw.rule.Knowledge == nil {
		return ""
	}
	return rw.rule.Knowledge.I18nContent.GetStrInLang(lang)
}

func (rw *RuleWrapper) GetRequiredTags(predefineTags map[model.TypeTag]*model.Tag, langTag model.TypeTag) (*model.Tag, *model.Tag, *model.Tag, error) {
	knowledgeTag := predefineTags[model.PredefineTagKnowledgeBase]
	i18nTag := predefineTags[langTag]
	dbTag := predefineTags[model.TypeTag(rw.rule.DBType)]

	if knowledgeTag == nil || i18nTag == nil || dbTag == nil {
		return nil, nil, nil, fmt.Errorf("missing required tags for rule %s", rw.rule.Name)
	}

	return knowledgeTag, i18nTag, dbTag, nil
}

func (rw *RuleWrapper) AddExtraTags(tagMap map[model.TypeTag]*model.Tag, predefineTags map[model.TypeTag]*model.Tag) {
	if ruleTags, ok := model.GetTagMapDefaultRuleKnowledge()[rw.rule.Name]; ok {
		for _, tagName := range ruleTags {
			if tag := predefineTags[tagName]; tag != nil {
				tagMap[tagName] = tag
			}
		}
	}
}

// CustomRuleWrapper 实现 RuleInfoProvider
type CustomRuleWrapper struct {
	BaseRuleWrapper
	rule *model.CustomRule
}

// CustomRuleWrapper 具体实现方法
func (crw *CustomRuleWrapper) GetRuleName() string {
	return crw.rule.RuleId
}

func (crw *CustomRuleWrapper) GetRuleDescription(lang language.Tag) string {
	return crw.rule.Desc
}

func (crw *CustomRuleWrapper) GetDBType() model.TypeTag {
	return model.TypeTag(crw.rule.DBType)
}

func (crw *CustomRuleWrapper) GetRuleContent(lang language.Tag) string {
	if crw.rule.Knowledge == nil {
		return ""
	}
	return crw.rule.Knowledge.I18nContent.GetStrInLang(lang)
}

func (crw *CustomRuleWrapper) GetRequiredTags(predefineTags map[model.TypeTag]*model.Tag, langTag model.TypeTag) (*model.Tag, *model.Tag, *model.Tag, error) {
	knowledgeTag := predefineTags[model.PredefineTagCustomizeKnowledgeBase]
	i18nTag := predefineTags[langTag]
	dbTag := predefineTags[model.TypeTag(crw.rule.DBType)]

	if knowledgeTag == nil || i18nTag == nil || dbTag == nil {
		return nil, nil, nil, fmt.Errorf("missing required tags for rule %s", crw.rule.Annotation)
	}

	return knowledgeTag, i18nTag, dbTag, nil
}

func (crw *CustomRuleWrapper) AddExtraTags(tagMap map[model.TypeTag]*model.Tag, predefineTags map[model.TypeTag]*model.Tag) {
	// 自定义规则没有额外的标签
}

// 统一的知识获取方法
func getKnowledgeWithTags(rule RuleInfoProvider, predefineTags map[model.TypeTag]*model.Tag) ([]*KnowledgeWithFilter, error) {
	baseWrapper := &BaseRuleWrapper{predefineTags: predefineTags}
	return baseWrapper.ToModelKnowledge(rule)
}

// 将默认规则转换为知识
func getDefaultRuleKnowledgeWithTags(rule *model.Rule, predefineTags map[model.TypeTag]*model.Tag) ([]*KnowledgeWithFilter, error) {
	return getKnowledgeWithTags(&RuleWrapper{rule: rule}, predefineTags)
}

// 将自定义规则转换为知识
func getCustomRuleKnowledgeWithTags(rule *model.CustomRule, predefineTags map[model.TypeTag]*model.Tag) ([]*KnowledgeWithFilter, error) {
	return getKnowledgeWithTags(&CustomRuleWrapper{rule: rule}, predefineTags)
}
