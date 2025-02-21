//go:build enterprise
// +build enterprise

package knowledge_base

import (
	"strings"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"golang.org/x/text/language"
)

/*
	该文件的主要目的是为了统一管理和维护规则到知识的转换逻辑，避免重复代码和不一致的转换规则。
	该文件的结构如下：
		1. KnowledgeInfoProvider 接口定义了获取规则信息的通用行为，包括获取规则名称、描述、内容等。
		2. BaseRuleWrapper 是 KnowledgeInfoProvider 接口的基础实现，提供了获取规则名称、描述、内容等的通用逻辑。
		3. RuleWrapper 是 KnowledgeInfoProvider 接口的实现，用于将 model.Rule 转换为 RuleInfoProvider。
		4. CustomRuleWrapper 是 KnowledgeInfoProvider 接口的实现，用于将 model.CustomRule 转换为 RuleInfoProvider。
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
	predefineTags           map[model.TypeTag]*model.Tag
	ruleKnowledgeContentMap map[string]string
}

// ToModelKnowledge 是模板方法，包含通用逻辑
func (brw *BaseRuleWrapper) ToModelKnowledge(rule KnowledgeInfoProvider) ([]*KnowledgeWithFilter, error) {
	var result []*KnowledgeWithFilter

	for langTag, lang := range model.GetTagMapPredefineLanguage() {
		// 收集所有标签
		tagMap := make(map[model.TypeTag]*model.Tag)

		// 获取必需的基础标签
		requiredTags, err := rule.GetRequiredTags(brw.predefineTags, langTag)
		if err != nil {
			return nil, err
		}

		// 添加必需标签
		for _, tag := range requiredTags {
			tagMap[tag.Name] = tag
		}

		// 额外的规则特定标签
		rule.AddExtraTags(tagMap, brw.predefineTags)

		// 创建知识对象
		knowledge := &model.Knowledge{
			Title:       rule.GetTitle(lang),
			Description: rule.GetDescription(lang),
			Content:     rule.GetContent(lang),
		}

		result = append(result, &KnowledgeWithFilter{
			knowledge:  knowledge,
			filterTags: requiredTags, // 必须的基础标签即为过滤标签，用于过滤出唯一的知识
			tagMap:     tagMap,
		})
	}

	return result, nil
}

// KnowledgeInfoProvider 作为策略接口，提供不同规则的获取方式
type KnowledgeInfoProvider interface {
	GetTitle(lang language.Tag) string
	GetDescription(lang language.Tag) string
	GetDBType() model.TypeTag
	GetContent(lang language.Tag) string
	GetRequiredTags(predefineTags map[model.TypeTag]*model.Tag, langTag model.TypeTag) ([]*model.Tag, error)
	AddExtraTags(tagMap map[model.TypeTag]*model.Tag, predefineTags map[model.TypeTag]*model.Tag)
}

// RuleWrapper 实现 KnowledgeInfoProvider
type RuleWrapper struct {
	BaseRuleWrapper
	rule *model.Rule
}

// RuleWrapper 具体实现方法
func (rw *RuleWrapper) GetTitle(lang language.Tag) string {
	return rw.rule.I18nRuleInfo.GetRuleInfoByLangTag(lang).Desc
}

func (rw *RuleWrapper) GetRuleName() string {
	return rw.rule.Name
}

func (rw *RuleWrapper) GetDescription(lang language.Tag) string {
	return rw.rule.I18nRuleInfo.GetRuleInfoByLangTag(lang).Annotation
}

func (rw *RuleWrapper) GetDBType() model.TypeTag {
	return model.TypeTag(rw.rule.DBType)
}

func (rw *RuleWrapper) GetContent(lang language.Tag) string {
	if rw.rule.DBType == driverV2.DriverTypeMySQL {
		// TODO 目前都是中文
		if lang != language.Chinese {
			return ""
		}
		return rw.ruleKnowledgeContentMap[rw.GetRuleName()]
	}
	if rw.rule.Knowledge == nil {
		return ""
	}
	knowledge := rw.rule.Knowledge.GetKnowledgeByLang(lang)
	if knowledge == nil {
		return ""
	}
	return knowledge.Content
}

func (rw *RuleWrapper) GetRequiredTags(predefineTags map[model.TypeTag]*model.Tag, langTag model.TypeTag) ([]*model.Tag, error) {
	var requiredTags []*model.Tag

	if knowledgeTag, ok := predefineTags[model.PredefineTagKnowledgeBase]; ok {
		requiredTags = append(requiredTags, knowledgeTag)
	}
	if i18nTag, ok := predefineTags[langTag]; ok {
		requiredTags = append(requiredTags, i18nTag)
	}
	if dbTag, ok := predefineTags[model.TypeTag(rw.rule.DBType)]; ok {
		requiredTags = append(requiredTags, dbTag)
	}
	// TODO 此处为了判断新的规则，临时写死SQLE前缀，后续需要优化
	versionTag := predefineTags[model.PredefineTagV1]
	if strings.HasPrefix(rw.rule.Name, "SQLE") {
		versionTag = predefineTags[model.PredefineTagV2]
	}
	if versionTag != nil {
		requiredTags = append(requiredTags, versionTag)
	}
	return requiredTags, nil
}

func (rw *RuleWrapper) AddExtraTags(tagMap map[model.TypeTag]*model.Tag, predefineTags map[model.TypeTag]*model.Tag) {
	if ruleTagsSlice, ok := model.GetTagMapDefaultRuleKnowledge()[rw.rule.Name]; ok {
		for _, ruleTags := range ruleTagsSlice {
			for _, tagName := range ruleTags {
				if tag := predefineTags[tagName]; tag != nil {
					tagMap[tagName] = tag
				}
			}
		}
	}
}

// CustomRuleWrapper 实现 KnowledgeInfoProvider
type CustomRuleWrapper struct {
	BaseRuleWrapper
	rule *model.CustomRule
}

// CustomRuleWrapper 具体实现方法
func (crw *CustomRuleWrapper) GetTitle(lang language.Tag) string {
	return crw.rule.Annotation
}

func (crw *CustomRuleWrapper) GetDescription(lang language.Tag) string {
	return crw.rule.Desc
}

func (crw *CustomRuleWrapper) GetDBType() model.TypeTag {
	return model.TypeTag(crw.rule.DBType)
}

func (crw *CustomRuleWrapper) GetRuleName() string {
	return crw.rule.RuleId
}

func (crw *CustomRuleWrapper) GetContent(lang language.Tag) string {
	if len(crw.rule.Knowledge) > 0 {
		return crw.rule.Knowledge.GetKnowledgeByLang(lang).Content
	}
	return ""
}

func (crw *CustomRuleWrapper) GetRequiredTags(predefineTags map[model.TypeTag]*model.Tag, langTag model.TypeTag) ([]*model.Tag, error) {
	var requiredTags []*model.Tag

	if knowledgeTag, ok := predefineTags[model.PredefineTagCustomizeKnowledgeBase]; ok {
		requiredTags = append(requiredTags, knowledgeTag)
	}
	if i18nTag, ok := predefineTags[langTag]; ok {
		requiredTags = append(requiredTags, i18nTag)
	}
	if dbTag, ok := predefineTags[model.TypeTag(crw.rule.DBType)]; ok {
		requiredTags = append(requiredTags, dbTag)
	}
	if versionTag, ok := predefineTags[model.PredefineTagV1]; ok {
		requiredTags = append(requiredTags, versionTag)
	}

	return requiredTags, nil
}

func (crw *CustomRuleWrapper) AddExtraTags(tagMap map[model.TypeTag]*model.Tag, predefineTags map[model.TypeTag]*model.Tag) {
	// 自定义规则没有额外的标签
}
