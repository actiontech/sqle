package model

import (
	"context"
	e "errors"
	"fmt"
	"strings"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type RuleTemplate struct {
	Model
	// global rule-template has no ProjectId
	ProjectId      ProjectUID               `gorm:"index"`
	Name           string                   `json:"name" gorm:"type:varchar(255)"`
	Desc           string                   `json:"desc" gorm:"type:varchar(255)"`
	DBType         string                   `json:"db_type" gorm:"type:varchar(255)"`
	RuleVersion    uint32                   `json:"rule_version" gorm:"not null;default:1"`
	RuleList       []RuleTemplateRule       `json:"rule_list" gorm:"foreignkey:rule_template_id;association_foreignkey:id"`
	CustomRuleList []RuleTemplateCustomRule `json:"custom_rule_list" gorm:"foreignkey:rule_template_id;association_foreignkey:id"`
}

func GenerateRuleByDriverRule(dr *driverV2.Rule, dbType string) *Rule {
	multiLanguageKnowledge := make(MultiLanguageKnowledge, 0, len(dr.I18nRuleInfo))
	r := &Rule{
		Name:         dr.Name,
		Level:        string(dr.Level),
		DBType:       dbType,
		Params:       dr.Params,
		Knowledge:    multiLanguageKnowledge,
		I18nRuleInfo: make(driverV2.I18nRuleInfo, len(dr.I18nRuleInfo)),
		AllowOffline: dr.AllowOffline,
		Version:      dr.Version,
	}
	for lang, v := range dr.I18nRuleInfo {
		r.Knowledge = append(r.Knowledge, &Knowledge{
			Title:       v.Desc,
			Content:     v.Knowledge.Content,
			Description: v.Annotation,
		})
		r.I18nRuleInfo[lang] = &driverV2.RuleInfo{
			Desc:       v.Desc,
			Annotation: v.Annotation,
			Category:   v.Category,
		}
	}
	for tagCate, tags := range dr.CategoryTags {
		for _, tag := range tags {
			r.Categories = append(r.Categories, &AuditRuleCategory{
				Category: tagCate,
				Tag:      tag,
			})
		}
	}
	return r
}

func ConvertRuleToDriverRule(r *Rule) *driverV2.Rule {
	dr := &driverV2.Rule{
		Name:         r.Name,
		Level:        driverV2.RuleLevel(r.Level),
		Params:       r.Params,
		I18nRuleInfo: make(map[language.Tag]*driverV2.RuleInfo, len(r.I18nRuleInfo)),
	}
	for lang, v := range r.I18nRuleInfo {
		dr.I18nRuleInfo[lang] = &driverV2.RuleInfo{
			Desc:       v.Desc,
			Annotation: v.Annotation,
			Category:   v.Category,
			Knowledge:  driverV2.RuleKnowledge{},
		}
	}
	return dr
}

type RuleKnowledge struct {
	Model
	I18nContent i18nPkg.I18nStr `gorm:"type:json"`
}

func (r *RuleKnowledge) TableName() string {
	return "rule_knowledge"
}

func (r *RuleKnowledge) GetContentByLangTag(lang language.Tag) string {
	if r == nil {
		return ""
	}
	return r.I18nContent.GetStrInLang(lang)
}

type Rule struct {
	Name            string                 `json:"name" gorm:"primary_key; not null;type:varchar(255)"`
	DBType          string                 `json:"db_type" gorm:"primary_key; not null; default:\"mysql\";type:varchar(255)"`
	Level           string                 `json:"level" example:"error" gorm:"type:varchar(255)"` // notice, warn, error
	Params          params.Params          `json:"params" gorm:"type:json"`
	HasAuditPower   bool                   `json:"has_audit_power" gorm:"type:bool" example:"true"`
	HasRewritePower bool                   `json:"has_rewrite_power" gorm:"type:bool" example:"true"`
	I18nRuleInfo    driverV2.I18nRuleInfo  `json:"i18n_rule_info" gorm:"type:json"`
	Categories      []*AuditRuleCategory   `gorm:"many2many:audit_rule_category_rels;joinForeignKey:RuleName,RuleDBType;joinReferences:CategoryId"`
	AllowOffline    bool                   `json:"allow_offline" gorm:"-"`
	Knowledge       MultiLanguageKnowledge `gorm:"many2many:rule_knowledge_relations" json:"rules"` // 规则和知识的关系
	Version         uint32                 `json:"version" gorm:"not null;default:1"`
}

func (r Rule) TableName() string {
	return "rules"
}

type RuleTemplateRule struct {
	RuleTemplateId uint          `json:"rule_template_id" gorm:"primary_key;auto_increment:false;"`
	RuleName       string        `json:"name" gorm:"primary_key;type:varchar(255)"`
	RuleLevel      string        `json:"level" gorm:"column:level;type:varchar(255)"`
	RuleParams     params.Params `json:"value" gorm:"column:rule_params;type:json"`
	RuleDBType     string        `json:"rule_db_type" gorm:"column:db_type; not null; type:varchar(255)"`

	Rule *Rule `json:"-" gorm:"foreignkey:Name,DBType;references:RuleName,RuleDBType"`
}

func (rtr *RuleTemplateRule) TableName() string {
	return "rule_template_rule"
}

func NewRuleTemplateRule(t *RuleTemplate, r *Rule) RuleTemplateRule {
	return RuleTemplateRule{
		RuleTemplateId: t.ID,
		RuleName:       r.Name,
		RuleLevel:      r.Level,
		RuleParams:     r.Params,
		RuleDBType:     r.DBType,
	}
}

func (rtr *RuleTemplateRule) GetRule() *Rule {
	rule := rtr.Rule
	if rtr.RuleLevel != "" {
		rule.Level = rtr.RuleLevel
	}
	if rtr.RuleParams != nil && len(rtr.RuleParams) > 0 {
		rule.Params = rtr.RuleParams
	}
	return rule
}

func (rtr *RuleTemplateRule) GetOptimizationRule() *Rule {
	rule := rtr.Rule
	if !rule.HasRewritePower {
		return nil
	}
	if rtr.RuleLevel != "" {
		rule.Level = rtr.RuleLevel
	}
	if rtr.RuleParams != nil && len(rtr.RuleParams) > 0 {
		rule.Params = rtr.RuleParams
	}
	return rule
}

type RuleTemplateCustomRule struct {
	RuleTemplateId uint   `json:"rule_template_id" gorm:"primary_key;auto_increment:false;"`
	RuleId         string `json:"rule_id" gorm:"primary_key;type:varchar(255)"`
	RuleLevel      string `json:"level" gorm:"column:level;type:varchar(255)"`
	RuleDBType     string `json:"rule_db_type" gorm:"column:db_type; not null;type:varchar(255)"`

	CustomRule *CustomRule `json:"-" gorm:"foreignkey:RuleId;references:RuleId"`
}

func NewRuleTemplateCustomRule(t *RuleTemplate, r *CustomRule) RuleTemplateCustomRule {
	return RuleTemplateCustomRule{
		RuleTemplateId: t.ID,
		RuleId:         r.RuleId,
		RuleLevel:      r.Level,
		RuleDBType:     r.DBType,
	}
}

func (rtr *RuleTemplateCustomRule) GetRule() *CustomRule {
	rule := rtr.CustomRule
	if rtr.RuleLevel != "" {
		rule.Level = rtr.RuleLevel
	}
	return rule
}

func (s *Storage) GetRuleTemplatesByInstance(inst *Instance) ([]RuleTemplate, error) {
	var associationRT []RuleTemplate
	err := s.db.Model(inst).Association("RuleTemplates").Find(&associationRT)
	return associationRT, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRuleTemplateNamesByProjectId(projectId string) ([]string, error) {
	records := []*RuleTemplate{}
	err := s.db.Model(&RuleTemplate{}).
		Select("name").
		Where("project_id = ?", projectId).
		Find(&records).
		Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}

	templateNames := make([]string, len(records))
	for i, record := range records {
		templateNames[i] = record.Name
	}
	return templateNames, nil
}

func (s *Storage) GetRuleTemplatesByInstanceNameAndProjectId(name string, projectId string) (*RuleTemplate, bool, error) {
	t := &RuleTemplate{}
	err := s.db.Joins("JOIN `instance_rule_template` ON `rule_templates`.`id` = `instance_rule_template`.`rule_template_id`").
		Joins("JOIN `instances` ON `instance_rule_template`.`instance_id` = `instances`.`id`").
		Where("`instances`.`name` = ?", name).
		Where("`instances`.`project_id` = ?", projectId).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRulesFromRuleTemplateByName(projectIds []string, name string) ([]*Rule, []*CustomRule, error) {
	tpl, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds(projectIds, name, "", "")
	if !exist {
		return nil, nil, errors.New(errors.DataNotExist, err)
	}
	if err != nil {
		return nil, nil, errors.New(errors.ConnectStorageError, err)
	}

	rules := make([]*Rule, 0, len(tpl.RuleList))
	for _, r := range tpl.RuleList {
		rules = append(rules, r.GetRule())
	}
	customRules := make([]*CustomRule, 0, len(tpl.CustomRuleList))
	for _, r := range tpl.CustomRuleList {
		customRules = append(customRules, r.GetRule())
	}
	return rules, customRules, nil
}

func (s *Storage) GetAllRulesByInstance(instance *Instance) ([]*Rule, []*CustomRule, error) {
	var template RuleTemplate

	if err := s.db.Where("id = ? ", instance.RuleTemplateId).First(&template).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, nil, errors.New(errors.ConnectStorageError, err)
		} else {
			return nil, nil, nil
		}
	}

	// 数据源可以绑定全局模板和项目模板
	return s.GetRulesFromRuleTemplateByName([]string{instance.ProjectId, ProjectIdForGlobalRuleTemplate}, instance.RuleTemplateName)
}

func (s *Storage) GetOptimizationRulesFromRuleTemplateByName(projectIds []string, name string) ([]*Rule, error) {
	tpl, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds(projectIds, name, "", "")
	if !exist {
		return nil, errors.New(errors.DataNotExist, err)
	}
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}

	rules := make([]*Rule, 0, len(tpl.RuleList))
	for _, r := range tpl.RuleList {
		rule := r.GetOptimizationRule()
		if rule != nil {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func (s *Storage) GetAllOptimizationRulesByInstance(instance *Instance) ([]*Rule, error) {
	var template RuleTemplate
	if err := s.db.Where("id = ? ", instance.RuleTemplateId).First(&template).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ConnectStorageError, err)
		} else {
			return nil, nil
		}
	}
	return s.GetOptimizationRulesFromRuleTemplateByName([]string{instance.ProjectId, ProjectIdForGlobalRuleTemplate}, instance.RuleTemplateName)
}

func (s *Storage) GetRuleTemplateByProjectIdAndName(projectId, name string) (*RuleTemplate, bool, error) {
	t := &RuleTemplate{}
	err := s.db.Where("name = ?", name).Where("project_id = ?", projectId).First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetGlobalAndProjectRuleTemplateByNameAndProjectId(name string, projectId string) (*RuleTemplate, bool, error) {
	t := &RuleTemplate{}
	err := s.db.Table("rule_templates").
		Where("rule_templates.name = ?", name).
		Where("project_id = ? OR project_id = 0", projectId).
		First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) IsRuleTemplateExistFromAnyProject(projectId ProjectUID, name string) (bool, error) {
	var count int64
	err := s.db.Model(&RuleTemplate{}).Where("name = ? and project_id = ?", name, string(projectId)).Count(&count).Error
	return count > 0, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetRuleTemplateDetailByNameAndProjectIds(projectIds []string, name string, fuzzy_keyword_rule string, tags string) (*RuleTemplate, bool, error) {
	dbOrder := func(db *gorm.DB) *gorm.DB {
		return db.Order("rule_template_rule.rule_name ASC")
	}
	fuzzyCondition := func(field, keyword string) func(*gorm.DB) *gorm.DB {
		if field == "" || keyword == "" {
			return func(db *gorm.DB) *gorm.DB { return db }
		}
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s` like ?", field), fmt.Sprintf("%%%s%%", keyword))
		}
	}
	t := &RuleTemplate{Name: name}
	err := s.db.Preload("RuleList", dbOrder).Preload("RuleList.Rule", fuzzyCondition("i18n_rule_info", fuzzy_keyword_rule)).Preload("RuleList.Rule.Categories").
		Preload("CustomRuleList.CustomRule", fuzzyCondition("desc", fuzzy_keyword_rule)).Preload("CustomRuleList.CustomRule.Categories").
		Where(t).
		Where("project_id IN (?)", projectIds).
		First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	t.RuleList = filterRulesByTag(t.RuleList, tags)
	t.CustomRuleList = filterCustomRulesByTag(t.CustomRuleList, tags)
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func filterRulesByTag(ruleTemplateRules []RuleTemplateRule, tags string) []RuleTemplateRule {
	if tags == "" {
		return ruleTemplateRules
	}
	var filteredRuleResult []RuleTemplateRule
	for _, ruleTemplateRule := range ruleTemplateRules {
		if templateRuleByTagsCondition(ruleTemplateRule.Rule.Categories, tags) {
			filteredRuleResult = append(filteredRuleResult, ruleTemplateRule)
		}
	}
	return filteredRuleResult
}

func filterCustomRulesByTag(ruleTemplateRules []RuleTemplateCustomRule, tags string) []RuleTemplateCustomRule {
	if tags == "" {
		return ruleTemplateRules
	}
	var filteredRuleResult []RuleTemplateCustomRule
	for _, ruleTemplateRule := range ruleTemplateRules {
		if templateRuleByTagsCondition(ruleTemplateRule.CustomRule.Categories, tags) {
			filteredRuleResult = append(filteredRuleResult, ruleTemplateRule)
		}
	}
	return filteredRuleResult
}

func templateRuleByTagsCondition(ruleCategories []*AuditRuleCategory, tags string) bool {
	tagSlice := strings.Split(tags, ",")
	tagSliceLen := len(tagSlice)
	tagMatchCount := 0
	for _, tag := range tagSlice {
		for _, ruleCategory := range ruleCategories {
			if ruleCategory.Tag == tag {
				tagMatchCount++
				break
			}
		}
	}
	return tagSliceLen == tagMatchCount
}

func (s *Storage) UpdateRuleTemplateRules(tpl *RuleTemplate, rules ...RuleTemplateRule) error {
	if err := s.db.Where(&RuleTemplateRule{RuleTemplateId: tpl.ID}).Delete(&RuleTemplateRule{}).Error; err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}

	//err := s.db.Debug().Model(tpl).Association("RuleList").Append(rules).Error //暂时没找到跳过更新关联表的方式
	//TODO gorm v1 没有预编译, 没有批量插入, 规则可能很多, 为了防止拼SQL批量插入导致SQL太长, 只能通过这种方式插入
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			rule.RuleTemplateId = tpl.ID
			err := tx.Omit("Rule").Create(&rule).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) CloneRuleTemplateRules(source, destination *RuleTemplate) error {
	return s.UpdateRuleTemplateRules(destination, source.RuleList...)
}

func (s *Storage) GetRuleTemplateRuleByName(name string, dbType string) ([]RuleTemplateRule, error) {
	ruleTemplateRule := []RuleTemplateRule{}
	result := s.db.Where("rule_name = ?", name).Where("db_type = ?", dbType).Find(&ruleTemplateRule)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return ruleTemplateRule, errors.New(errors.ConnectStorageError, result.Error)
}

func (s *Storage) CloneRuleTemplateCustomRules(source, destination *RuleTemplate) error {
	return s.UpdateRuleTemplateCustomRules(destination, source.CustomRuleList...)
}

func GetRuleMapFromAllArray(allRules ...[]Rule) map[string]Rule {
	ruleMap := map[string]Rule{}
	for _, rules := range allRules {
		for _, rule := range rules {
			ruleMap[rule.Name] = rule
		}
	}
	return ruleMap
}

func (s *Storage) GetRuleTemplateById(id uint64) (*RuleTemplate, error) {
	ruleTemplate := new(RuleTemplate)

	err := s.db.Where("id = ?", id).First(ruleTemplate).Error
	return ruleTemplate, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRuleTemplateTips(projectId, dbType string) ([]*RuleTemplate, error) {
	ruleTemplates := []*RuleTemplate{}

	db := s.db.Where("project_id = ?", projectId)
	if dbType != "" {
		db = db.Where("db_type = ?", dbType)
	}
	err := db.Find(&ruleTemplates).Error
	return ruleTemplates, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRule(name, dbType string) (*Rule, bool, error) {
	rule := Rule{Name: name, DBType: dbType}
	err := s.db.First(&rule).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return &rule, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllRules() ([]*Rule, error) {
	rules := []*Rule{}
	err := s.db.Preload("Categories").Preload("Knowledge").Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) DeleteCascadeRule(name, dbType string) error {
	/*
		1. 删除rule_template_rule表中rule_name和db_type等于name和dbType的记录
		2. 删除rule表中rule_name和db_type等于name和dbType的记录
		3. 删除rule_knowledge_relations表中rule_name和db_type等于name和dbType的记录
		4. 删除audit_rule_category_rels表中rule_name和db_type等于name和dbType的记录
		5. 不删除knowledge表中的记录，因为knowledge表中的记录不与规则强绑定，后续允许存在脱离规则的知识，且删除规则时，需要保留客户积累的知识
	*/
	err := s.db.Exec(`delete u,t,k,c
					from rules u 
					left join rule_template_rule t on u.name = t.rule_name and u.db_type = t.db_type 
					left join rule_knowledge_relations k on u.name = k.rule_name AND u.db_type = k.rule_db_type
					left join audit_rule_category_rels c on u.name = c.rule_name and u.db_type = c.rule_db_type
					where u.name = ? AND u.db_type = ?`, name, dbType).Error
	return err
}

func (s *Storage) GetAllRuleByDBType(dbType string) ([]*Rule, error) {
	rules := []*Rule{}
	err := s.db.Where(&Rule{DBType: dbType}).Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRulesByNames(names []string) ([]*Rule, error) {
	rules := []*Rule{}
	err := s.db.Where("name in (?)", names).Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRulesByNamesAndDBType(names []string, dbType string) ([]*Rule, error) {
	rules := []*Rule{}
	err := s.db.Where("db_type = ?", dbType).Where("name in (?)", names).Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

var ErrRuleNotExist = e.New("rule not exist")

func (s *Storage) GetAndCheckRuleExist(ruleNames []string, dbType string) (map[string]Rule, map[string]struct{}, error) {
	rules, err := s.GetRulesByNamesAndDBType(ruleNames, dbType)
	if err != nil {
		return nil, nil, err
	}
	existRules := map[string]Rule{}
	for _, rule := range rules {
		existRules[rule.Name] = *rule
	}
	notExistRuleNames := make(map[string]struct{})
	for _, ruleName := range ruleNames {
		if _, ok := existRules[ruleName]; !ok {
			notExistRuleNames[ruleName] = struct{}{}
		}
	}

	if len(notExistRuleNames) > 0 {
		return nil, notExistRuleNames, fmt.Errorf("rule %v not exist:%w", notExistRuleNames, ErrRuleNotExist)
	}

	return existRules, notExistRuleNames, nil
}

func (s *Storage) IsRuleTemplateExist(ruleTemplateName string, projectIds []string) (bool, error) {
	if len(projectIds) <= 0 {
		return false, nil
	}
	var count int64
	err := s.db.Table("rule_templates").
		Where("name = ?", ruleTemplateName).
		Where("deleted_at IS NULL").
		Where("project_id IN (?)", projectIds).
		Count(&count).Error
	return count > 0, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanNamesByRuleTemplate(
	ruleTemplateName string) (auditPlanNames []string, err error) {

	var auditPlans []*AuditPlan

	err = s.db.Model(&AuditPlan{}).
		Select("name").
		Where("rule_template_name=?", ruleTemplateName).
		Find(&auditPlans).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	auditPlanNames = make([]string, len(auditPlans))
	for i := range auditPlans {
		auditPlanNames[i] = auditPlans[i].Name
	}

	return
}

func (s *Storage) GetAuditPlanNamesByRuleTemplateAndProject(
	ruleTemplateName string, projectID string) (auditPlanNames []string, err error) {

	var auditPlans []*AuditPlan

	err = s.db.Model(&AuditPlan{}).
		Select("name").
		Where("rule_template_name=?", ruleTemplateName).
		Where("project_id=?", projectID).
		Find(&auditPlans).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	auditPlanNames = make([]string, len(auditPlans))
	for i := range auditPlans {
		auditPlanNames[i] = auditPlans[i].Name
	}

	return auditPlanNames, nil
}

func (s *Storage) GetRuleTypeByDBType(ctx context.Context, DBType string) ([]string, error) {
	lang := locale.Bundle.GetLangTagFromCtx(ctx)
	rules := []*Rule{}
	err := s.db.Where("db_type = ?", DBType).Find(&rules).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	var ruleDBTypes []string
	categoryMap := make(map[string]struct{})
	for _, v := range rules {
		category := v.I18nRuleInfo.GetRuleInfoByLangTag(lang).Category
		if _, exist := categoryMap[category]; !exist {
			ruleDBTypes = append(ruleDBTypes, category)
			categoryMap[category] = struct{}{}
		}
	}
	return ruleDBTypes, nil
}

type CustomRule struct {
	Model
	RuleId string `json:"rule_id" gorm:"index:unique; not null; type:varchar(255)"`

	Desc       string                 `json:"desc" gorm:"not null; type:varchar(255)"`
	Annotation string                 `json:"annotation" gorm:"type:varchar(1024)"`
	DBType     string                 `json:"db_type" gorm:"not null; default:\"mysql\"; type:varchar(255)"`
	Level      string                 `json:"level" example:"error" gorm:"type:varchar(255)"` // notice, warn, error
	Typ        string                 `json:"type" gorm:"column:type; not null; type:varchar(255)"`
	RuleScript string                 `json:"rule_script" gorm:"type:text"`
	ScriptType string                 `json:"script_type" gorm:"not null; default:\"regular\"; type:varchar(255)"`
	Categories []*AuditRuleCategory   `json:"categories" gorm:"many2many:custom_rule_category_rels;foreignKey:RuleId;joinForeignKey:CustomRuleId;joinReferences:CategoryId"`
	Knowledge  MultiLanguageKnowledge `gorm:"many2many:custom_rule_knowledge_relations" json:"custom_rules"` // 自定义规则和知识的关系
}

func (s *Storage) GetCustomRuleByRuleId(ruleId string) (*CustomRule, bool, error) {
	rule := &CustomRule{}
	err := s.db.Preload("Categories").Where("rule_id = ?", ruleId).First(rule).Error
	if err == gorm.ErrRecordNotFound {
		return rule, false, nil
	}
	return rule, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetCustomRulesByDBTypeAndFuzzyDesc(queryFields, filterDbType, fuzzyDesc string) ([]*CustomRule, error) {
	rules := []*CustomRule{}
	db := s.db.Select(queryFields)
	if filterDbType != "" {
		db = db.Where("db_type=?", filterDbType)
	}
	if fuzzyDesc != "" {
		db = db.Where("`desc` like ?", "%"+fuzzyDesc+"%")
	}
	err := db.Find(&rules).Error

	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetCustomRules(queryFields string) ([]*CustomRule, error) {
	rules := []*CustomRule{}
	err := s.db.Preload("Categories").Select(queryFields).Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetCustomRulesByDescAndDBType(filterDesc, filterDbType string) (CustomRule, bool, error) {
	rule := CustomRule{}
	err := s.db.Where("db_type = ? and `desc` = ?", filterDbType, filterDesc).First(&rule).Error
	if err == gorm.ErrRecordNotFound {
		return rule, false, nil
	}
	return rule, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateCustomRuleByRuleId(ruleId string, attrs interface{}) error {
	err := s.db.Table("custom_rules").Where("rule_id = ?", ruleId).Updates(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateCustomRuleCategoriesByRuleId(ruleId string, tags []string) error {
	err := s.db.Where(&CustomRuleCategoryRel{CustomRuleId: ruleId}).Delete(&CustomRuleCategoryRel{}).Error
	if err != nil {
		return err
	}
	var auditRuleCategory []*AuditRuleCategory
	err = s.db.Model(AuditRuleCategory{}).Where("tag in (?)", tags).Find(&auditRuleCategory).Error
	if err != nil {
		return err
	}
	for _, category := range auditRuleCategory {
		err = s.db.Save(CustomRuleCategoryRel{CategoryId: category.ID, CustomRuleId: ruleId}).Error
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetCustomRuleByDBTypeAndScriptType(DBType, ScriptType string) ([]CustomRule, bool, error) {
	rules := []CustomRule{}
	result := s.db.Where("db_type = ? and script_type = ?", DBType, ScriptType).Find(&rules)
	if result.RowsAffected == 0 {
		return rules, false, nil
	}
	return rules, true, errors.New(errors.ConnectStorageError, result.Error)
}

type CustomTypeCount struct {
	Type      string `json:"type"`
	TypeCount uint   `json:"type_count"`
}

func (s *Storage) GetCustomRuleTypeCountByDBType(DBType string) ([]*CustomTypeCount, error) {
	typeCounts := []*CustomTypeCount{}
	err := s.db.Model(&CustomRule{}).
		Select("type, count(1) type_count, MIN(created_at) as min_created_at").
		Where("db_type = ?", DBType).
		Group("type").
		Order("min_created_at").Scan(&typeCounts).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	return typeCounts, nil
}

func (s *Storage) GetCustomRulesByDBType(filterDbType string) ([]*CustomRule, error) {
	rules := []*CustomRule{}
	err := s.db.Where("db_type = ?", filterDbType).Find(&rules).Error
	if len(rules) == 0 {
		return rules, nil
	}
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetCustomRulesByIds(ruleIds []string) ([]*CustomRule, error) {
	rules := []*CustomRule{}
	err := s.db.Where("rule_id in (?)", ruleIds).Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAndCheckCustomRuleExist(ruleIds []string) (map[string]CustomRule, error) {
	rules, err := s.GetCustomRulesByIds(ruleIds)
	if err != nil {
		return nil, err
	}
	existRules := map[string]CustomRule{}
	for _, rule := range rules {
		existRules[rule.RuleId] = *rule
	}
	notExistRuleNames := []string{}
	for _, ruleId := range ruleIds {
		if _, ok := existRules[ruleId]; !ok {
			notExistRuleNames = append(notExistRuleNames, ruleId)
		}
	}
	if len(notExistRuleNames) > 0 {
		return nil, errors.New(errors.DataNotExist,
			fmt.Errorf("rule %s not exist", strings.Join(notExistRuleNames, ", ")))
	}
	return existRules, nil
}

func (s *Storage) UpdateRuleTemplateCustomRules(tpl *RuleTemplate, rules ...RuleTemplateCustomRule) error {
	if err := s.db.Where(&RuleTemplateCustomRule{RuleTemplateId: tpl.ID}).Delete(&RuleTemplateCustomRule{}).Error; err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}

	//err := s.db.Debug().Model(tpl).Association("RuleList").Append(rules).Error //暂时没找到跳过更新关联表的方式
	//TODO gorm v1 没有预编译, 没有批量插入, 规则可能很多, 为了防止拼SQL批量插入导致SQL太长, 只能通过这种方式插入
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			rule.RuleTemplateId = tpl.ID
			err := tx.Omit("CustomRule").Create(&rule).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllRulesByTmpNameAndProjectIdInstanceDBType(ruleTemplateName string, projectId string,
	inst *Instance, dbType string) (rules []*Rule, customRules []*CustomRule, err error) {
	if ruleTemplateName != "" {
		if projectId == "" {
			return nil, nil, errors.New(errors.DataInvalid,
				fmt.Errorf("project id is needed when rule template name is given"))
		}
		rules, customRules, err = s.GetRulesFromRuleTemplateByName([]string{projectId, ProjectIdForGlobalRuleTemplate}, ruleTemplateName)
	} else {
		if inst != nil {
			rules, customRules, err = s.GetAllRulesByInstance(inst)
		} else {
			templateName := s.GetDefaultRuleTemplateName(dbType)
			// 默认规则模板从全局模板里拿
			rules, customRules, err = s.GetRulesFromRuleTemplateByName([]string{ProjectIdForGlobalRuleTemplate}, templateName)
		}
	}
	if err != nil {
		return nil, nil, err
	}
	return rules, customRules, nil
}

func (s *Storage) DeleteCustomRule(ruleId string) error {
	err := s.Tx(func(tx *gorm.DB) error {
		if err := tx.Where("rule_id = ?", ruleId).Delete(&CustomRule{}).Error; err != nil {
			return err
		}
		if err := tx.Where("rule_id = ?", ruleId).Delete(&RuleTemplateCustomRule{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.ConnectStorageErrWrapper(err)
	}

	return nil
}
