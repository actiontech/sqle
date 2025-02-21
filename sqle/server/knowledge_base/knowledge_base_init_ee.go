//go:build enterprise
// +build enterprise

package knowledge_base

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/license"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
)

func CheckKnowledgeBaseLicense() error {
	dmsLicense, err := getDMSLicense()
	if err != nil {
		log.Logger().Errorf("get dms license failed: %v", err)
		return err
	}
	support, err := dmsLicense.CheckSupportKnowledgeBase()
	if err != nil {
		log.Logger().Errorf("check support knowledge base failed: %v", err)
		return err
	}
	if !support {
		return fmt.Errorf("license not support knowledge base")
	}
	return nil
}

func getDMSLicense() (*license.License, error) {
	reply, err := dms.GetLicense(context.TODO(), config.GetOptions().SqleOptions.DMSServerAddress)
	if err != nil {
		log.Logger().Errorf("get license failed: %v", err)
		return nil, err
	}
	dmsLicense, err := license.GetDMSLicense(reply.Content)
	if err != nil {
		log.Logger().Errorf("get dms license failed: %v", err)
		return nil, err
	}
	return dmsLicense, nil
}

// 迁移规则知识库到知识库，并且关联标签
func LoadKnowledge(rulesMap map[string][]*model.Rule) error {
	// 创建系统保留标签
	storage := model.GetStorage()
	predefineTags, err := NewTagService(storage).GetOrCreatePredefinedTags()
	if err != nil {
		log.Logger().Errorf("get or create predefined tags failed: %v", err)
		return err
	}
	mysqlRuleKnowledgeMap, err := rule.GetDefaultRulesKnowledge()
	if err != nil {
		return fmt.Errorf("get default rules knowledge failed: %v", err)
	}
	// 加载购买License后的AI规则的知识
	if dmsLicense, err := getDMSLicense(); err == nil {
		// 仅初始化支持的数据库类型的知识库
		if support, err := dmsLicense.CheckSupportKnowledgeBase(); err == nil && support {
			for _, dbType := range dmsLicense.GetKnowledgeBaseDBTypes() {
				if dbType == driverV2.DriverTypeMySQL { // TODO 后续会增加其他知识库
					// 获取AI规则的知识
					aiRuleKnowledge, err := ai.GetAIRulesKnowledge()
					if err != nil {
						return fmt.Errorf("get ai rules knowledge failed: %v", err)
					}
					for ruleName, knowledgeContent := range aiRuleKnowledge {
						mysqlRuleKnowledgeMap[ruleName] = knowledgeContent
					}
				}
			}
		}
	} else {
		log.Logger().Errorf("init knowledge failed, get dms license failed: %v", err)
	}
	// load rule knowledge
	for _, rules := range rulesMap {
		for _, rule := range rules {

			warpedRule := &RuleWrapper{
				BaseRuleWrapper: BaseRuleWrapper{
					predefineTags:           predefineTags,
					ruleKnowledgeContentMap: mysqlRuleKnowledgeMap,
				},
				rule: rule,
			}

			knowledgeWithTags, err := warpedRule.ToModelKnowledge(warpedRule)
			if err != nil {
				return fmt.Errorf("failed to get knowledge for rule: %w", err)
			}
			for _, item := range knowledgeWithTags {
				modelKnowledge, err := storage.CreateKnowledgeWithTags(item.knowledge, rule.Name, item.tagMap, item.filterTags)
				if err != nil {
					return fmt.Errorf("failed to create knowledge: %w", err)
				}
				// 创建规则和知识的关系
				err = storage.CreateRuleKnowledgeRelation(uint64(modelKnowledge.ID), rule.Name, rule.DBType)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// 初始化一条自定义规则的知识
func InitCustomRuleKnowledge(customRule *model.CustomRule) error {
	predefineTags, err := NewTagService(model.GetStorage()).GetOrCreatePredefinedTags()
	if err != nil {
		log.Logger().Errorf("get or create predefined tags failed: %v", err)
		return err
	}
	warpedRule := &CustomRuleWrapper{
		BaseRuleWrapper: BaseRuleWrapper{
			predefineTags: predefineTags,
		},
		rule: customRule,
	}
	knowledgeWithTags, err := warpedRule.ToModelKnowledge(warpedRule)
	if err != nil {
		return fmt.Errorf("failed to get knowledge for rule: %w", err)
	}
	s := model.GetStorage()
	for _, item := range knowledgeWithTags {
		modelKnowledge, err := s.CreateKnowledgeWithTags(item.knowledge, customRule.DBType, item.tagMap, item.filterTags)
		if err != nil {
			return fmt.Errorf("failed to create knowledge: %w", err)
		}
		err = s.CreateCustomRuleKnowledgeRelation(uint64(modelKnowledge.ID), customRule)
		if err != nil {
			return err
		}
	}
	return nil
}
