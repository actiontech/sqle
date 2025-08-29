package model

import (
	"bytes"
	"database/sql"
	sqlDriver "database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	xerrors "github.com/pkg/errors"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var storage *Storage

var storageMutex sync.Mutex

var pluginRules map[string][]*Rule

const dbDriver = "mysql"

func InitStorage(s *Storage) {
	storageMutex.Lock()
	defer storageMutex.Unlock()
	storage = s
}

var MockTime, _ = time.Parse("0000-00-00 00:00:00.0000000", "0000-00-00 00:00:00.0000000")

func InitMockStorage(db *sql.DB) {
	storageMutex.Lock()
	defer storageMutex.Unlock()
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	storage = &Storage{db: gormDB}

	// Custom NowFunc solve this problem:
	// 	When mock SQL which will update CreateAt/UpdateAt fields,
	// 	GORM will auto-update this field by NowFunc(when is is empty),
	// 	then it will never equal to our expectation(always later than our expectation).
	gormDB.NowFunc = func() time.Time {
		return MockTime
	}
}

func GetStorage() *Storage {
	storageMutex.Lock()
	defer storageMutex.Unlock()
	return storage
}

func GetDb() *gorm.DB {
	return storage.db
}

func GetSqlxDb() (*sqlx.DB, error) {
	sdb, err := storage.db.DB()
	if err != nil {
		return nil, err
	}
	db := sqlx.NewDb(sdb, dbDriver)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return db, nil
}

type Model struct {
	ID        uint           `json:"id" gorm:"primary_key" example:"1"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:current_timestamp(3)" example:"2018-10-21T16:40:23+08:00"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:current_timestamp(3) on update current_timestamp(3)" example:"2018-10-21T16:40:23+08:00"`
	DeletedAt gorm.DeletedAt `json:"-" sql:"index" gorm:"index"`
}

func (m Model) GetIDStr() string {
	return fmt.Sprintf("%d", m.ID)
}

func NewStorage(user, password, host, port, schema string, debug bool) (*Storage, error) {
	log.Logger().Infof("connecting to storage, host: %s, port: %s, user: %s, schema: %s",
		host, port, user, schema)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, schema)

	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	if debug {
		config.Logger = log.NewGormLogWrapper(logger.Info)
	} else {
		config.Logger = log.NewGormLogWrapper(logger.Silent)
	}
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		log.Logger().Errorf("connect to storage failed, error: %v", err)
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	log.Logger().Info("connected to storage")
	return &Storage{db: db}, errors.New(errors.ConnectStorageError, err)
}

type Storage struct {
	db *gorm.DB
}

var autoMigrateList = []interface{}{
	&AuditPlanReportSQLV2{},
	&AuditPlanReportV2{},
	&AuditPlanSQLV2{},
	&AuditPlan{},
	&ExecuteSQL{},
	&RoleOperation{},
	&RollbackSQL{},
	&RuleTemplateRule{},
	&RuleTemplate{},
	&Rule{},
	&AuditRuleCategory{},
	&AuditRuleCategoryRel{},
	&CustomRuleCategoryRel{},
	&SqlWhitelist{},
	&Task{},
	&AuditFile{},
	&WorkflowRecord{},
	&WorkflowStepTemplate{},
	&WorkflowStep{},
	&WorkflowTemplate{},
	&Workflow{},
	&SqlQueryExecutionSql{},
	&SqlQueryHistory{},
	&TaskGroup{},
	&WorkflowInstanceRecord{},
	&FeishuInstance{},
	&IM{},
	&DingTalkInstance{},
	&OperationRecord{},
	&CustomRule{},
	&RuleTemplateCustomRule{},
	&SQLAuditRecord{},
	&RuleKnowledge{},
	&SqlManage{},
	&BlackListAuditPlanSQL{},
	&CompanyNotice{},
	&SqlManageEndpoint{},
	&SQLDevRecord{},
	&WechatRecord{},
	&FeishuScheduledRecord{},
	&InstanceAuditPlan{},
	&AuditPlanV2{},
	&AuditPlanTaskInfo{},
	&SQLManageRecord{},
	&SQLManageRecordProcess{},
	&SQLManageQueue{},
	&SqlManageMetricRecord{},
	&SqlManageMetricValue{},
	&SqlManageMetricExecutePlanRecord{},
	&ReportPushConfig{},
	&ReportPushConfigRecord{},
	&SqlVersion{},
	&SqlVersionStage{},
	&SqlVersionStagesDependency{},
	&WorkflowVersionStage{},
	&Tag{},
	&Knowledge{},
	&DataLock{},
	&SQLManageRawSQL{},
}

func (s *Storage) AutoMigrate() error {
	err := s.db.AutoMigrate(autoMigrateList...)
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	err = s.db.SetupJoinTable(&Rule{}, "Categories", &AuditRuleCategoryRel{})
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	err = s.db.SetupJoinTable(&CustomRule{}, "Categories", &CustomRuleCategoryRel{})
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	if !s.db.Migrator().HasIndex(&SqlManage{}, "idx_project_id_status_deleted_at") {
		err = s.db.Exec("CREATE INDEX idx_project_id_status_deleted_at ON sql_manages (project_id, status, deleted_at)").Error
		if err != nil {
			return errors.New(errors.ConnectStorageError, err)
		}
	}
	return nil
}

func (s *Storage) CreateRuleCategoriesRelated() error {
	err := s.CreateRuleCategories()
	if err != nil {
		return err
	}
	// 创建自定义规则和分类的关系
	err = s.UpdateCustomRuleCategoryRels()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateCustomRuleCategoryRels() error {
	customRules := []*CustomRule{}
	s.db.Find(&customRules)
	for _, customRule := range customRules {
		if customRule.Typ == "" {
			// 新的规则分类Typ字段为""说明已经有了新的分类关系，直接忽略
			continue
		}
		_, existed, err := s.FirstCustomRuleCategoryRelByCustomRuleId(customRule.RuleId)
		if err != nil {
			return err
		}
		// 已存在规则关系直接忽略
		if existed {
			return nil
		}
		tags := mappingToNewCategory(customRule.Desc, customRule.Typ)
		// 获取分类表中的分类信息
		auditRuleCategories, err := s.GetAuditRuleCategoryByTagIn(tags)
		if err != nil {
			return err
		}
		for _, newCategory := range auditRuleCategories {
			customerCategoryRel := CustomRuleCategoryRel{CategoryId: newCategory.ID, CustomRuleId: customRule.RuleId}
			err = s.db.Create(&customerCategoryRel).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Storage) CreateRuleCategories() error {
	isCategoryExistInDB := func(categories []*AuditRuleCategory, targetCategory *AuditRuleCategory) (*AuditRuleCategory, bool) {
		for i := range categories {
			if categories[i].Category != targetCategory.Category || categories[i].Tag != targetCategory.Tag {
				continue
			}
			return categories[i], true
		}
		return nil, false
	}
	categories, err := s.GetAllCategories()
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	categoryTagMap := map[string][]string{
		plocale.RuleCategoryOperand.ID:              {plocale.RuleTagDatabase.ID, plocale.RuleTagTablespace.ID, plocale.RuleTagTable.ID, plocale.RuleTagColumn.ID, plocale.RuleTagIndex.ID, plocale.RuleTagView.ID, plocale.RuleTagProcedure.ID, plocale.RuleTagFunction.ID, plocale.RuleTagTrigger.ID, plocale.RuleTagEvent.ID, plocale.RuleTagUser.ID, plocale.RuleTagSequence.ID, plocale.RuleTagBusiness.ID},
		plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID, plocale.RuleTagDDL.ID, plocale.RuleTagDCL.ID, plocale.RuleTagIntegrity.ID, plocale.RuleTagComplete.ID, plocale.RuleTagQuery.ID, plocale.RuleTagJoin.ID, plocale.RuleTagTransaction.ID, plocale.RuleTagPrivilege.ID, plocale.RuleTagManagement.ID, plocale.RuleTagSQLTablespace.ID, plocale.RuleTagSQLFunction.ID, plocale.RuleTagSQLProcedure.ID, plocale.RuleTagSQLTrigger.ID, plocale.RuleTagSQLView.ID},
		plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID, plocale.RuleTagMaintenance.ID, plocale.RuleTagSecurity.ID, plocale.RuleTagCorrection.ID},
		plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID, plocale.RuleTagOffline.ID},
		plocale.RuleCategoryAuditPerformanceCost.ID: {plocale.RuleTagPerformanceCostHigh.ID, plocale.RuleTagPerformanceCostMedium.ID, plocale.RuleTagPerformanceCostLow.ID},
	}
	for category, tags := range categoryTagMap {
		for _, tag := range tags {
			auditRuleCategory := &AuditRuleCategory{Category: category, Tag: tag}
			_, existed := isCategoryExistInDB(categories, auditRuleCategory)
			if !existed {
				err := s.Save(auditRuleCategory)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Storage) CreateRulesIfNotExist(rulesMap map[string][]*Rule) error {
	isRuleExistInDB := func(rulesInDB []*Rule, targetRule *Rule, dbType string) (*Rule, bool) {
		for i := range rulesInDB {
			rule := rulesInDB[i]
			if rule.DBType != dbType || rule.Name != targetRule.Name {
				continue
			}
			return rule, true
		}
		return nil, false
	}

	rulesInDB, err := s.GetAllRules()
	if err != nil {
		return err
	}
	for dbType, rules := range rulesMap {
		for _, rule := range rules {
			existedRule, exist := isRuleExistInDB(rulesInDB, rule, dbType)
			// rule will be created or update if:
			// 1. rule not exist;
			if !exist {
				err := errors.New(errors.ConnectStorageError, s.db.Omit("Categories", "Knowledge").Save(rule).Error)
				if err != nil {
					return err
				}
			} else {
				isRuleLevelSame := existedRule.Level == rule.Level
				isI18nInfoSame := reflect.DeepEqual(existedRule.I18nRuleInfo, rule.I18nRuleInfo)
				isHasAuditPowerSame := existedRule.HasAuditPower == rule.HasAuditPower
				isHasRewritePowerSame := existedRule.HasRewritePower == rule.HasRewritePower
				isRuleVersionSame := existedRule.Version == rule.Version
				existRuleParam, err := existedRule.Params.Value()
				if err != nil {
					return err
				}
				pluginRuleParam, err := rule.Params.Value()
				if err != nil {
					return err
				}
				isParamSame := reflect.DeepEqual(existRuleParam, pluginRuleParam)

				if !isI18nInfoSame || !isRuleLevelSame || !isParamSame || !isHasAuditPowerSame || !isHasRewritePowerSame || !isRuleVersionSame {
					// 保存规则
					err := errors.New(errors.ConnectStorageError, s.db.Omit("Categories", "Knowledge").Save(rule).Error)
					if err != nil {
						return err
					}
					if !isParamSame {
						// 同步模板规则的参数
						err = s.UpdateRuleTemplateRulesParams(rule, dbType)
						if err != nil {
							return err
						}
					}
				}
			}

			// 持久化规则分类信息
			categoryError := s.UpdateRuleCategoryRels(rule, existedRule)
			if categoryError != nil {
				return fmt.Errorf("update rule category rels err: %w", categoryError)
			}
		}
	}
	return nil
}

func (s *Storage) UpdateRuleTemplateRulesParams(pluginRule *Rule, dbType string) error {
	ruleTemplateRules, err := s.GetRuleTemplateRuleByName(pluginRule.Name, dbType)
	if err != nil {
		return err
	}
	for _, ruleTemplateRule := range ruleTemplateRules {
		ruleTemplateRuleParamsMap := make(map[string]string)
		for _, p := range ruleTemplateRule.RuleParams {
			ruleTemplateRuleParamsMap[p.Key] = p.Value
		}
		for _, pluginParam := range pluginRule.Params {
			// 避免参数的值被还原成默认
			if value, ok := ruleTemplateRuleParamsMap[pluginParam.Key]; ok {
				pluginParam.Value = value
			}
		}
		ruleTemplateRule.RuleParams = pluginRule.Params
		err = s.Save(&ruleTemplateRule)
		if err != nil {
			return err
		}
	}
	return nil
}

// 为所有模板删除插件中已不存在的规则
func (s *Storage) DeleteRulesIfNotExist(rules map[string][]*Rule) error {
	pluginRules = rules
	// 避免清空规则
	if len(pluginRules) <= 0 {
		return nil
	}
	rulesInDB, err := s.GetAllRules()
	if err != nil {
		return err
	}
	for _, dbRule := range rulesInDB {
		// 判断Plugin是不是读取到了，防止模板里规则被清空
		if pluginExist := PluginIsExist(dbRule.DBType); !pluginExist {
			continue
		}
		// 判断规则是否被删除
		if ruleExist := DBRuleInPluginRule(dbRule); !ruleExist {
			err := s.DeleteCascadeRule(dbRule.Name, dbRule.DBType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func PluginIsExist(dbType string) bool {
	for pluginDBType := range pluginRules {
		if dbType == pluginDBType {
			return true
		}
	}
	return false
}

func DBRuleInPluginRule(dbRule *Rule) bool {
	for dbType, rules := range pluginRules {
		for _, rule := range rules {
			if dbRule.Name == rule.Name && dbRule.DBType == dbType {
				return true
			}
		}
	}
	return false
}

// 整合sql优化规则与插件规则，并赋予审核、重写能力
func MergeOptimizationRules(pluginRulesMap map[string][]*driverV2.Rule) map[string][]*Rule {
	resultAllRulesMap := map[string][]*Rule{}
	for dbType, pluginRules := range pluginRulesMap {
		resultAllRules := []*Rule{}
		for _, rule := range pluginRules {
			resultRule := GenerateRuleByDriverRule(rule, dbType)
			resultRule.HasAuditPower = true
			resultRule.HasRewritePower = false
			resultAllRules = append(resultAllRules, resultRule)
		}

		resultAllRulesMap[dbType] = resultAllRules
	}
	return resultAllRulesMap
}

// func (s *Storage) CreateDefaultRole() error {
// 	roles, err := s.GetAllRoleTip()
// 	if err != nil {
// 		return err
// 	}
// 	if len(roles) > 0 {
// 		return nil
// 	}

// 	// dev
// 	err = s.SaveRoleAndAssociations(&Role{
// 		Name: "dev",
// 		Desc: "dev",
// 	}, []uint{OP_WORKFLOW_SAVE, OP_AUDIT_PLAN_SAVE, OP_SQL_QUERY_QUERY})
// 	if err != nil {
// 		return err
// 	}

// 	// dba
// 	err = s.SaveRoleAndAssociations(&Role{
// 		Name: "dba",
// 		Desc: "dba",
// 	}, []uint{OP_WORKFLOW_AUDIT, OP_SQL_QUERY_QUERY})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

const DefaultProjectUid string = "700300"

func (s *Storage) CreateDefaultWorkflowTemplateIfNotExist() error {
	_, exist, err := s.GetWorkflowTemplateByProjectId(ProjectUID(DefaultProjectUid))
	if err != nil {
		return err
	}
	if !exist {
		td := DefaultWorkflowTemplate(DefaultProjectUid)
		err = s.SaveWorkflowTemplate(td)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) CreateDefaultTemplateIfNotExist(projectId ProjectUID, rules map[string][]*driverV2.Rule) error {
	for dbType, dbRules := range rules {
		versionRules := make(map[uint32][]*driverV2.Rule)
		for _, dbRule := range dbRules {
			versionRules[dbRule.Version] = append(versionRules[dbRule.Version], dbRule)
		}

		for version, perVersionRules := range versionRules {
			templateName := s.GetRuleTemplateName(dbType, version)
			exist, err := s.IsRuleTemplateExistFromAnyProject(projectId, templateName)
			if err != nil {
				return xerrors.Wrap(err, "get rule template failed")
			}
			if exist {
				continue
			}
			t := &RuleTemplate{
				ProjectId:   projectId,
				Name:        templateName,
				DBType:      dbType,
				RuleVersion: version,
			}
			if err := s.Save(t); err != nil {
				return err
			}

			ruleList := make([]RuleTemplateRule, 0, len(perVersionRules))
			for _, r := range perVersionRules {
				if r.Level != driverV2.RuleLevelError {
					continue
				}
				modelRule := GenerateRuleByDriverRule(r, dbType)
				ruleList = append(ruleList, RuleTemplateRule{
					RuleTemplateId: t.ID,
					RuleName:       modelRule.Name,
					RuleLevel:      modelRule.Level,
					RuleParams:     modelRule.Params,
					RuleDBType:     modelRule.DBType,
				})
			}
			if err := s.UpdateRuleTemplateRules(t, ruleList...); err != nil {
				return xerrors.Wrap(err, "update rule template rules failed")
			}
		}
	}

	return nil
}

func mappingToNewCategory(ruleDesc string, oldCategory string) []string {
	// 当旧规则是命名规范的映射关系
	if oldCategory == plocale.RuleTypeNamingConvention.Other {
		if strings.Contains(ruleDesc, plocale.RuleTagDatabase.Other) || strings.Contains(ruleDesc, "对象") {
			return []string{plocale.RuleTagDatabase.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagIndex.Other) || strings.Contains(ruleDesc, "主键") {
			return []string{plocale.RuleTagIndex.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagTable.Other) {
			return []string{plocale.RuleTagTable.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagView.Other) {
			return []string{plocale.RuleTagView.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagColumn.Other) {
			return []string{plocale.RuleTagColumn.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagProcedure.Other) {
			return []string{plocale.RuleTagProcedure.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagFunction.Other) {
			return []string{plocale.RuleTagFunction.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagTrigger.Other) {
			return []string{plocale.RuleTagTrigger.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagEvent.Other) {
			return []string{plocale.RuleTagEvent.ID}
		} else if strings.Contains(ruleDesc, plocale.RuleTagUser.Other) {
			return []string{plocale.RuleTagUser.ID}
		} else {
			return []string{
				plocale.RuleTagDatabase.ID, plocale.RuleTagTablespace.ID, plocale.RuleTagTable.ID, plocale.RuleTagColumn.ID, plocale.RuleTagIndex.ID, plocale.RuleTagView.ID, plocale.RuleTagProcedure.ID, plocale.RuleTagFunction.ID, plocale.RuleTagTrigger.ID, plocale.RuleTagEvent.ID, plocale.RuleTagUser.ID}
		}
	}
	newCategoryMap := categoryMapping[oldCategory]
	if newCategoryMap == nil {
		return []string{
			plocale.RuleTagDatabase.ID, plocale.RuleTagTablespace.ID, plocale.RuleTagTable.ID, plocale.RuleTagColumn.ID, plocale.RuleTagIndex.ID, plocale.RuleTagView.ID, plocale.RuleTagProcedure.ID, plocale.RuleTagFunction.ID, plocale.RuleTagTrigger.ID, plocale.RuleTagEvent.ID, plocale.RuleTagUser.ID}
	}
	tags := make([]string, 0)
	for _, newTags := range newCategoryMap {
		tags = append(tags, newTags...)
	}
	return tags
}

var categoryMapping = map[string]map[string][]string{
	plocale.RuleTypeGlobalConfig.Other: {
		plocale.RuleCategoryAuditPurpose.ID: {plocale.RuleTagPerformance.ID},
	},
	plocale.RuleTypeIndexingConvention.Other: {
		plocale.RuleCategoryOperand.ID: {plocale.RuleTagIndex.ID},
	},
	plocale.RuleTypeIndexOptimization.Other: {
		plocale.RuleCategoryOperand.ID: {plocale.RuleTagIndex.ID},
	},
	plocale.RuleTypeIndexInvalidation.Other: {
		plocale.RuleCategoryOperand.ID: {plocale.RuleTagIndex.ID},
	},
	plocale.RuleTypeDDLConvention.Other: {
		plocale.RuleCategorySQL.ID: {plocale.RuleTagDDL.ID},
	},
	plocale.RuleTypeDMLConvention.Other: {
		plocale.RuleCategorySQL.ID: {plocale.RuleTagDML.ID},
	},
	plocale.RuleTypeDQLConvention.Other: {
		plocale.RuleCategorySQL.ID: {plocale.RuleTagDML.ID},
	},
	plocale.RuleTypeUsageSuggestion.Other: {
		plocale.RuleCategoryAuditPurpose.ID: {plocale.RuleTagMaintenance.ID},
	},
	plocale.RuleTypeExecutePlan.Other: {
		plocale.RuleCategoryAuditPurpose.ID: {plocale.RuleTagPerformance.ID},
	},
	plocale.RuleTypeDistributedConvention.Other: {
		plocale.RuleCategoryAuditPurpose.ID: {plocale.RuleTagMaintenance.ID},
	},
}

func (s *Storage) UpdateRuleCategoryRels(rule, existedRule *Rule) error {
	var err error
	var auditRuleCategories []*AuditRuleCategory
	if len(rule.Categories) > 0 {
		// 规则定义了新的规则分类标签
		db := s.db.Model(AuditRuleCategory{})
		for _, v := range rule.Categories {
			db = db.Or("category = ? and tag = ?", v.Category, v.Tag)
		}
		if err = db.Find(&auditRuleCategories).Error; err != nil {
			return err
		}
	} else {
		// 旧的分类 映射到 新的规则分类标签
		var tags []string
		ruleInfo := rule.I18nRuleInfo.GetRuleInfoByLangTag(language.Chinese)
		oldCategory := ruleInfo.Category
		ruleDesc := ruleInfo.Desc
		tags = mappingToNewCategory(ruleDesc, oldCategory)
		tags = append(tags, plocale.RuleTagOnline.ID)
		if rule.AllowOffline {
			tags = append(tags, plocale.RuleTagOffline.ID)
		}

		// 根据特定规则名添加性能消耗分类
		if performanceCostId, ok := ruleNameToPerformanceCostId[rule.Name]; ok {
			tags = append(tags, performanceCostId)
		}

		// 获取分类表中的分类信息
		auditRuleCategories, err = s.GetAuditRuleCategoryByTagIn(tags)
		if err != nil {
			return err
		}
	}

	var existedRuleCategoryMap map[string]*AuditRuleCategory
	if existedRule != nil {
		// map缓存已存在的规则分类标签
		existedRuleCategoryMap = make(map[string]*AuditRuleCategory, len(existedRule.Categories))
		for k := range existedRule.Categories {
			existedRuleCategoryMap[existedRule.Categories[k].Category+existedRule.Categories[k].Tag] = existedRule.Categories[k]
		}
	}
	for _, newCategory := range auditRuleCategories {
		_, existed := existedRuleCategoryMap[newCategory.Category+newCategory.Tag]
		if !existed {
			// 创建规则与新增的标签分类的关联关系
			auditRuleCategoryRel := AuditRuleCategoryRel{CategoryId: newCategory.ID, RuleName: rule.Name, RuleDBType: rule.DBType}
			err = s.db.Save(&auditRuleCategoryRel).Error
			if err != nil {
				return err
			}
		} else {
			// 可复用的关联关系，在map缓存中删除
			delete(existedRuleCategoryMap, newCategory.Category+newCategory.Tag)
		}
	}
	for _, existedCategory := range existedRuleCategoryMap {
		// 删除在map缓存剩下（规则弃用的标签）的关联关系
		auditRuleCategoryRel := AuditRuleCategoryRel{CategoryId: existedCategory.ID, RuleName: rule.Name, RuleDBType: rule.DBType}
		err = s.db.Delete(&auditRuleCategoryRel).Error
		if err != nil {
			return err
		}
	}

	return err
}

// TODO : 避免直接用规则名称映射
var ruleNameToPerformanceCostId = map[string] /*ruleId*/ string /*auditPerformanceLevelCategorieID*/ {
	// old rule: MySQL
	rule.DMLCheckSelectRows:                plocale.RuleTagPerformanceCostHigh.ID,
	rule.DMLCheckAffectedRows:              plocale.RuleTagPerformanceCostHigh.ID,
	rule.ConfigOptimizeIndexEnabled:        plocale.RuleTagPerformanceCostHigh.ID,
	rule.DDLCheckIndexOption:               plocale.RuleTagPerformanceCostHigh.ID,
	rule.DDLCheckCompositeIndexDistinction: plocale.RuleTagPerformanceCostHigh.ID,

	// new rule: Oracle
	"Oracle_011": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_012": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_017": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_019": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_044": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_046": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_050": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_077": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_078": plocale.RuleTagPerformanceCostHigh.ID,
	"Oracle_080": plocale.RuleTagPerformanceCostHigh.ID,
}

func (s *Storage) GetDefaultRuleTemplateName(dbType string) string {
	metas := driver.GetPluginManager().GetDriverMetasOfPlugin(dbType)
	var latestVersion uint32
	for _, v := range metas.RuleVersionIncluded {
		if v > latestVersion {
			latestVersion = v
		}
	}
	return s.GetRuleTemplateName(dbType, latestVersion)
}

func (s *Storage) GetRuleTemplateName(dbType string, version uint32) string {
	return fmt.Sprintf("default_%v_V%dRules", dbType, version)
}

func (s *Storage) CreateDefaultReportPushConfigIfNotExist(projectUId string) error {
	_, exist, err := s.GetReportPushConfigByProjectId(ProjectUID(projectUId))
	if err != nil {
		return err
	}
	if !exist {
		err = s.InitReportPushConfigInProject(projectUId)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (s *Storage) CreateAdminUser() error {
// 	_, exist, err := s.GetUserByName(DefaultAdminUser)
// 	if err != nil {
// 		return err
// 	}
// 	if !exist {
// 		return s.Save(&User{
// 			Name:     DefaultAdminUser,
// 			Password: "admin",
// 		})
// 	}
// 	return nil
// }

const DefaultProject = "default"

// func (s *Storage) CreateDefaultProject() error {
// 	exist, err := s.IsProjectExist()
// 	if err != nil {
// 		return err
// 	}
// 	if exist {
// 		return nil
// 	}

// 	err = s.CreateProject(DefaultProject, "", 700200 /* TODO 从公共包传？*/)
// 	return err
// }

func (s *Storage) Exist(model interface{}) (bool, error) {
	var count int64
	err := s.db.Model(model).Where(model).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return count > 0, nil
}

func (s *Storage) Create(model interface{}) error {
	return errors.New(errors.ConnectStorageError, s.db.Create(model).Error)
}

func (s *Storage) Save(model interface{}) error {
	return errors.New(errors.ConnectStorageError, s.db.Save(model).Error)
}

func (s *Storage) BatchSave(value any, batchSize int) error {
	reflectValue := reflect.Indirect(reflect.ValueOf(value))
	if reflectValue.Kind() != reflect.Slice && reflectValue.Kind() != reflect.Array {
		return s.db.Save(value).Error
	}

	sliceLen := reflectValue.Len()
	if sliceLen == 0 {
		return nil
	}

	txDB := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			txDB.Rollback()
		}
	}()

	for i := 0; i < sliceLen; i += batchSize {
		end := i + batchSize
		if end > sliceLen {
			end = sliceLen
		}

		if err := txDB.Save(reflectValue.Slice(i, end).Interface()).Error; err != nil {
			txDB.Rollback()
			return errors.ConnectStorageErrWrapper(err)
		}
	}

	if err := txDB.Commit().Error; err != nil {
		txDB.Rollback()
		return errors.ConnectStorageErrWrapper(err)
	}

	return nil
}

func (s *Storage) Update(model interface{}, attrs ...interface{}) error {
	return errors.New(errors.ConnectStorageError, s.db.Model(model).UpdateColumns(attrs).Error)
}

func (s *Storage) Delete(model interface{}) error {
	return errors.New(errors.ConnectStorageError, s.db.Delete(model).Error)
}

func (s *Storage) HardDelete(model interface{}) error {
	return errors.New(errors.ConnectStorageError, s.db.Unscoped().Delete(model).Error)
}

func (s *Storage) TxExec(fn func(tx *sql.Tx) error) error {
	db, err := s.db.DB()
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	err = fn(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.NewEntry().Error("rollback sql transact failed, err:", err)
		}
		return errors.New(errors.ConnectStorageError, err)
	}
	err = tx.Commit()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.NewEntry().Error("rollback sql transact failed, err:", err)
		}
		return errors.New(errors.ConnectStorageError, err)
	}
	return nil
}

func (s *Storage) Tx(fn func(txDB *gorm.DB) error) (err error) {
	txDB := s.db.Begin()
	err = fn(txDB)
	if err != nil {
		txDB.Rollback()
		return errors.ConnectStorageErrWrapper(err)
	}

	err = txDB.Commit().Error
	if err != nil {
		txDB.Rollback()
		return errors.ConnectStorageErrWrapper(err)
	}
	return nil
}

type RowList []string

func (r *RowList) Scan(src interface{}) error {
	var data string
	switch src := src.(type) {
	case nil:
		data = ""
	case string:
		data = src
	case []byte:
		data = string(src)
	default:
		return fmt.Errorf("scan: unable to scan type %T into []string", src)
	}
	*r = []string{}
	if data != "" {
		l := strings.Split(data, ",")
		for _, v := range l {
			if v != "" {
				*r = append(*r, v)
			}
		}
	}
	return nil
}

func (r RowList) Value() (sqlDriver.Value, error) {
	return strings.Join(r, ","), nil
}

type JSON json.RawMessage

func (j JSON) OriginValue() (map[string]interface{}, error) {
	mp := map[string]interface{}{}
	return mp, json.Unmarshal(j, &mp)
}

// Value impl sql.driver.Valuer interface
func (j JSON) Value() (sqlDriver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	bytes, err := json.RawMessage(j).MarshalJSON()
	return string(bytes), err
}

// Scan impl sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("null")
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSON value: %s", value)
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (rl *RowList) ForceConvertIntSlice() []uint {
	res := make([]uint, len(*rl))
	for i := range *rl {
		n, _ := strconv.Atoi((*rl)[i])
		res[i] = uint(n)
	}
	return res
}

func (s *Storage) getTemplateQueryResult(data map[string]interface{}, result interface{}, queryTpl string, bodyTemplates ...string) error {
	var buff bytes.Buffer
	tpl := template.New("getQuery")
	var err error
	for _, bt := range bodyTemplates {
		if tpl, err = tpl.Parse(bt); err != nil {
			return err
		}
	}
	tpl, err = tpl.Parse(queryTpl)
	if err != nil {
		return err
	}
	err = tpl.Execute(&buff, data)
	if err != nil {
		return err
	}

	sqlxDb, err := GetSqlxDb()
	if err != nil {
		return err
	}

	query, args, err := sqlx.Named(buff.String(), data)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	query = sqlxDb.Rebind(query)
	err = sqlxDb.Select(result, query, args...)
	return err
}
