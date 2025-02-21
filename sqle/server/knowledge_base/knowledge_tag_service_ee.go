//go:build enterprise
// +build enterprise

package knowledge_base

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/model"
)

/*
	该文件用于定义知识库的预定义规则，并提供统一的接口用于创建和获取标签。
	该文件的主要目的是为了统一管理和维护标签的创建和获取逻辑，避免重复代码和不一致的标签定义。
	该文件的结构如下：
		1. TagManager 接口定义了标签处理的通用行为，包括获取父级标签、获取预定义标签和获取已存在的标签等。
		2. BaseTagManager 是 TagManager 接口的基础实现，提供了获取父级标签和获取已存在的标签的通用逻辑。
		3. KnowledgeBaseTagManager、LanguageTagManager 和 DBTypeTagManager 分别是知识库、语言和数据库类型的标签管理器，分别实现了 TagManager 接口。
	当需要增加新的标签类型时，只需要实现 TagManager 接口，并在 NewTagService 函数中添加对应的标签管理器即可。
*/

// TagManager 接口定义tag处理的通用行为
type TagManager interface {
	// GetParentTag 获取父级标签
	GetParentTag() (parentTag *model.Tag, err error)
	// GetPreDefinedTags 获取预定义的标签集合
	GetPreDefinedTags() map[model.TypeTag]struct{}
	// GetExistingTags 获取数据库中已存在的标签
	GetExistingTags(parentID uint) (map[model.TypeTag]*model.Tag, error)
}

// BaseTagManager 提供基础实现
type BaseTagManager struct {
	storage    *model.Storage
	parentName model.TypeTag
}

// NewBaseTagManager 创建基础tag管理器
func NewBaseTagManager(storage *model.Storage, parentName model.TypeTag) *BaseTagManager {
	return &BaseTagManager{
		storage:    storage,
		parentName: parentName,
	}
}

func (b *BaseTagManager) GetParentTag() (*model.Tag, error) {
	return b.storage.GetOrCreateTag(&model.Tag{Name: b.parentName})
}

func (b *BaseTagManager) GetExistingTags(parentID uint) (map[model.TypeTag]*model.Tag, error) {
	tags, err := b.storage.GetSubTags(parentID)
	if err != nil {
		return nil, fmt.Errorf("get sub tags failed: %w", err)
	}

	tagMap := make(map[model.TypeTag]*model.Tag)
	for _, tag := range tags {
		tagMap[tag.Name] = tag
	}
	return tagMap, nil
}

// KnowledgeBaseTagManager 知识库标签管理器
type KnowledgeBaseTagManager struct {
	*BaseTagManager
}

func NewKnowledgeBaseTagManager(storage *model.Storage) *KnowledgeBaseTagManager {
	return &KnowledgeBaseTagManager{
		BaseTagManager: NewBaseTagManager(storage, model.PredefineTagKnowledgeBase),
	}
}

func (k *KnowledgeBaseTagManager) GetPreDefinedTags() map[model.TypeTag]struct{} {
	defaultTagNameMap := make(map[model.TypeTag]struct{})
	for _, tagNamesSlice := range model.GetTagMapDefaultRuleKnowledge() {
		for _, tagNames := range tagNamesSlice {
			for _, tagName := range tagNames {
				defaultTagNameMap[tagName] = struct{}{}
			}
		}
	}
	return defaultTagNameMap
}

// 自定义标签管理器
type CustomTagManager struct {
	*BaseTagManager
}

func NewCustomTagManager(storage *model.Storage) *CustomTagManager {
	return &CustomTagManager{
		BaseTagManager: NewBaseTagManager(storage, model.PredefineTagCustomizeKnowledgeBase),
	}
}

// GetPreDefinedTags 获取预定义的标签集合
func (c *CustomTagManager) GetPreDefinedTags() map[model.TypeTag]struct{} {
	return nil
}

// LanguageTagManager 语言标签管理器
type LanguageTagManager struct {
	*BaseTagManager
}

func NewLanguageTagManager(storage *model.Storage) *LanguageTagManager {
	return &LanguageTagManager{
		BaseTagManager: NewBaseTagManager(storage, model.PredefineTagLanguage),
	}
}

func (l *LanguageTagManager) GetPreDefinedTags() map[model.TypeTag]struct{} {
	langTagMap := make(map[model.TypeTag]struct{})
	for tagName := range model.GetTagMapPredefineLanguage() {
		langTagMap[tagName] = struct{}{}
	}
	return langTagMap
}

// DBTypeTagManager 数据库类型标签管理器
type DBTypeTagManager struct {
	*BaseTagManager
}

func NewDBTypeTagManager(storage *model.Storage) *DBTypeTagManager {
	return &DBTypeTagManager{
		BaseTagManager: NewBaseTagManager(storage, model.PredefineTagDBType),
	}
}

func (d *DBTypeTagManager) GetPreDefinedTags() map[model.TypeTag]struct{} {
	return model.GetTagMapPredefineDBType()
}

// VersionTagManager 版本标签管理器
type VersionTagManager struct {
	*BaseTagManager
}

func NewVersionTagManager(storage *model.Storage) *VersionTagManager {
	return &VersionTagManager{
		BaseTagManager: NewBaseTagManager(storage, model.PredefineTagVersion),
	}
}

func (v *VersionTagManager) GetPreDefinedTags() map[model.TypeTag]struct{} {
	versionTagMap := make(map[model.TypeTag]struct{})
	for tagName := range model.GetTagMapPredefineVersion() {
		versionTagMap[tagName] = struct{}{}
	}
	return versionTagMap
}

// TagService 用于协调不同类型的tag管理
type TagService struct {
	storage     *model.Storage
	tagManagers []TagManager
}

func NewTagService(storage *model.Storage) *TagService {
	return &TagService{
		storage: storage,
		tagManagers: []TagManager{
			NewKnowledgeBaseTagManager(storage),
			NewLanguageTagManager(storage),
			NewDBTypeTagManager(storage),
			NewCustomTagManager(storage),
			NewVersionTagManager(storage),
		},
	}
}

// GetOrCreatePredefinedTags 重构后的主函数
func (service *TagService) GetOrCreatePredefinedTags() (map[model.TypeTag]*model.Tag, error) {
	result := make(map[model.TypeTag]*model.Tag)

	for _, manager := range service.tagManagers {
		// 获取父标签
		parentTag, err := manager.GetParentTag()
		if err != nil {
			return nil, fmt.Errorf("get parent tag failed: %w", err)
		}
		result[parentTag.Name] = parentTag

		// 获取现有标签
		existingTags, err := manager.GetExistingTags(parentTag.ID)
		if err != nil {
			return nil, fmt.Errorf("get existing tags failed: %w", err)
		}

		// 创建缺失的预定义标签
		for tagName := range manager.GetPreDefinedTags() {
			if _, exists := existingTags[tagName]; exists {
				result[tagName] = existingTags[tagName]
				continue
			}

			newTag, err := service.storage.CreateSubTag(parentTag.ID, &model.Tag{Name: tagName})
			if err != nil {
				return nil, fmt.Errorf("create tag failed: %w", err)
			}
			result[tagName] = newTag
		}
	}

	return result, nil
}
