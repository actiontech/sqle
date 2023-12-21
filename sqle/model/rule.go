package model

import (
	"fmt"
	"strings"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/jinzhu/gorm"
)

type RuleTemplate struct {
	Model
	// global rule-template has no ProjectId
	ProjectId      ProjectUID               `gorm:"index"`
	Name           string                   `json:"name"`
	Desc           string                   `json:"desc"`
	DBType         string                   `json:"db_type"`
	RuleList       []RuleTemplateRule       `json:"rule_list" gorm:"foreignkey:rule_template_id;association_foreignkey:id"`
	CustomRuleList []RuleTemplateCustomRule `json:"custom_rule_list" gorm:"foreignkey:rule_template_id;association_foreignkey:id"`
}

func GenerateRuleByDriverRule(dr *driverV2.Rule, dbType string) *Rule {
	return &Rule{
		Name:       dr.Name,
		Desc:       dr.Desc,
		Annotation: dr.Annotation,
		Level:      string(dr.Level),
		Typ:        dr.Category,
		DBType:     dbType,
		Params:     dr.Params,
		Knowledge: &RuleKnowledge{
			Content: dr.Knowledge.Content,
		},
	}
}

func ConvertRuleToDriverRule(r *Rule) *driverV2.Rule {
	return &driverV2.Rule{
		Name:       r.Name,
		Desc:       r.Desc,
		Annotation: r.Annotation,
		Category:   r.Typ,
		Level:      driverV2.RuleLevel(r.Level),
		Params:     r.Params,
	}
}

type RuleKnowledge struct {
	Model
	Content string `gorm:"type:longtext"`
}

func (r *RuleKnowledge) TableName() string {
	return "rule_knowledge"
}

func (r *RuleKnowledge) GetContent() string {
	if r == nil {
		return ""
	}
	return r.Content
}

type Rule struct {
	Name        string         `json:"name" gorm:"primary_key; not null"`
	DBType      string         `json:"db_type" gorm:"primary_key; not null; default:\"mysql\""`
	Desc        string         `json:"desc"`
	Annotation  string         `json:"annotation" gorm:"column:annotation"`
	Level       string         `json:"level" example:"error"` // notice, warn, error
	Typ         string         `json:"type" gorm:"column:type; not null"`
	Params      params.Params  `json:"params" gorm:"type:varchar(1000)"`
	KnowledgeId uint           `json:"knowledge_id"`
	Knowledge   *RuleKnowledge `json:"knowledge" gorm:"foreignkey:KnowledgeId"`
}

func (r Rule) TableName() string {
	return "rules"
}

type RuleTemplateRule struct {
	RuleTemplateId uint          `json:"rule_template_id" gorm:"primary_key;auto_increment:false;"`
	RuleName       string        `json:"name" gorm:"primary_key;"`
	RuleLevel      string        `json:"level" gorm:"column:level;"`
	RuleParams     params.Params `json:"value" gorm:"column:rule_params;type:varchar(1000)"`
	RuleDBType     string        `json:"rule_db_type" gorm:"column:db_type; not null;"`

	Rule *Rule `json:"-" gorm:"foreignkey:Name,DBType;association_foreignkey:RuleName,RuleDBType"`
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

type RuleTemplateCustomRule struct {
	RuleTemplateId uint   `json:"rule_template_id" gorm:"primary_key;auto_increment:false;"`
	RuleId         string `json:"rule_id" gorm:"primary_key;"`
	RuleLevel      string `json:"level" gorm:"column:level;"`
	RuleDBType     string `json:"rule_db_type" gorm:"column:db_type; not null;"`

	CustomRule *CustomRule `json:"-" gorm:"foreignkey:RuleId;association_foreignkey:RuleId"`
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
	err := s.db.Model(inst).Association("RuleTemplates").Find(&associationRT).Error
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
		Where("`instances`.`project_id` = ?", projectId).Find(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRulesFromRuleTemplateByName(projectIds []string, name string) ([]*Rule, []*CustomRule, error) {
	tpl, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds(projectIds, name, "")
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
	var count int
	err := s.db.Model(&RuleTemplate{}).Where("name = ? and project_id = ?", name, string(projectId)).Count(&count).Error
	return count > 0, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetRuleTemplateDetailByNameAndProjectIds(projectIds []string, name string, fuzzy_keyword_rule string) (*RuleTemplate, bool, error) {
	dbOrder := func(db *gorm.DB) *gorm.DB {
		return db.Order("rule_template_rule.rule_name ASC")
	}
	fuzzy_condition := func(db *gorm.DB) *gorm.DB {
		if fuzzy_keyword_rule == "" {
			return db
		}
		return db.Where("`desc` like ? OR annotation like ?", fmt.Sprintf("%%%s%%", fuzzy_keyword_rule), fmt.Sprintf("%%%s%%", fuzzy_keyword_rule))
	}
	t := &RuleTemplate{Name: name}
	err := s.db.Preload("RuleList", dbOrder).Preload("RuleList.Rule", fuzzy_condition).Preload("CustomRuleList.CustomRule", fuzzy_condition).
		Where(t).
		Where("project_id IN (?)", projectIds).
		First(t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.ConnectStorageError, err)
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

	db := s.db.Select("id,name, db_type").Where("project_id = ?", projectId)
	if dbType != "" {
		db = db.Where("db_type = ?", dbType)
	}
	err := db.Find(&ruleTemplates).Error
	return ruleTemplates, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetRule(name, dbType string) (*Rule, bool, error) {
	rule := Rule{Name: name, DBType: dbType}
	err := s.db.Find(&rule).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return &rule, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllRules() ([]*Rule, error) {
	rules := []*Rule{}
	err := s.db.Preload("Knowledge").Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
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

func (s *Storage) GetRulesByNamesAndDBType(names []string, dbType string) ([]Rule, error) {
	rules := []Rule{}
	err := s.db.Where("db_type = ?", dbType).Where("name in (?)", names).Find(&rules).Error
	return rules, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAndCheckRuleExist(ruleNames []string, dbType string) (map[string]Rule, error) {
	rules, err := s.GetRulesByNamesAndDBType(ruleNames, dbType)
	if err != nil {
		return nil, err
	}
	existRules := map[string]Rule{}
	for _, rule := range rules {
		existRules[rule.Name] = rule
	}
	notExistRuleNames := []string{}
	for _, userName := range ruleNames {
		if _, ok := existRules[userName]; !ok {
			notExistRuleNames = append(notExistRuleNames, userName)
		}
	}
	if len(notExistRuleNames) > 0 {
		return nil, errors.New(errors.DataNotExist,
			fmt.Errorf("rule %s not exist", strings.Join(notExistRuleNames, ", ")))
	}
	return existRules, nil
}

func (s *Storage) IsRuleTemplateExist(ruleTemplateName string, projectIds []string) (bool, error) {
	if len(projectIds) <= 0 {
		return false, nil
	}
	var count int
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

func (s *Storage) GetRuleTypeByDBType(DBType string) ([]string, error) {
	rules := []*Rule{}
	err := s.db.Select("type").Where("db_type = ?", DBType).Group("type").Find(&rules).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	ruleDBTypes := make([]string, len(rules))
	for i := range rules {
		ruleDBTypes[i] = rules[i].Typ
	}
	return ruleDBTypes, nil
}

type CustomRule struct {
	Model
	RuleId      string         `json:"rule_id" gorm:"unique; not null"`
	Desc        string         `json:"desc" gorm:"not null"`
	Annotation  string         `json:"annotation"`
	DBType      string         `json:"db_type" gorm:"not null; default:\"mysql\""`
	Level       string         `json:"level" example:"error"` // notice, warn, error
	Typ         string         `json:"type" gorm:"column:type; not null"`
	RuleScript  string         `json:"rule_script" gorm:"type:text"`
	ScriptType  string         `json:"script_type" gorm:"not null; default:\"regular\""`
	KnowledgeId uint           `json:"knowledge_id"`
	Knowledge   *RuleKnowledge `json:"knowledge" gorm:"foreignkey:KnowledgeId"`
}

func (s *Storage) GetCustomRuleByRuleId(ruleId string) (*CustomRule, bool, error) {
	rule := &CustomRule{}
	err := s.db.Where("rule_id = ?", ruleId).Find(rule).Error
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
	err := s.db.Select(queryFields).Find(&rules).Error
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

func (s *Storage) UpdateCustomRuleByRuleId(ruleId string, attrs ...interface{}) error {
	err := s.db.Table("custom_rules").Where("rule_id = ?", ruleId).Update(attrs...).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetCustomRuleByDBTypeAndScriptType(DBType, ScriptType string) ([]CustomRule, bool, error) {
	rules := []CustomRule{}
	err := s.db.Where("db_type = ? and script_type = ?", DBType, ScriptType).Find(&rules).Error
	if err == gorm.ErrRecordNotFound {
		return rules, false, nil
	}
	return rules, true, errors.New(errors.ConnectStorageError, err)
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
	if err == gorm.ErrRecordNotFound {
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
