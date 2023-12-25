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

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	xerrors "github.com/pkg/errors"
)

var storage *Storage

var storageMutex sync.Mutex

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
	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		panic(err)
	}
	storage = &Storage{db: gormDB}

	// Custom NowFunc solve this problem:
	// 	When mock SQL which will update CreateAt/UpdateAt fields,
	// 	GORM will auto-update this field by NowFunc(when is is empty),
	// 	then it will never equal to our expectation(always later than our expectation).
	gorm.NowFunc = func() time.Time {
		return MockTime
	}
}

func GetStorage() *Storage {
	storageMutex.Lock()
	defer storageMutex.Unlock()
	return storage
}

func UpdateStorage(newStorage *Storage) {
	storageMutex.Lock()
	defer storageMutex.Unlock()
	storage.db.Close()
	storage = newStorage
}

func GetDb() *gorm.DB {
	return storage.db
}

func GetSqlxDb() *sqlx.DB {
	db := sqlx.NewDb(storage.db.DB(), dbDriver)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return db
}

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key" example:"1"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:current_timestamp" example:"2018-10-21T16:40:23+08:00"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:current_timestamp on update current_timestamp" example:"2018-10-21T16:40:23+08:00"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

func (m Model) GetIDStr() string {
	return fmt.Sprintf("%d", m.ID)
}

func NewStorage(user, password, host, port, schema string, debug bool) (*Storage, error) {
	log.Logger().Infof("connecting to storage, host: %s, port: %s, user: %s, schema: %s",
		host, port, user, schema)
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, schema))
	if err != nil {
		log.Logger().Errorf("connect to storage failed, error: %v", err)
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	if debug {
		db.SetLogger(log.Logger().WithField("type", "sql"))
		db.LogMode(true)
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
	&SqlWhitelist{},
	&SystemVariable{},
	&Task{},
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
	&SqlManageSqlAuditRecord{},
	&BlackListAuditPlanSQL{},
	&CompanyNotice{},
	&SqlManageEndpoint{},
}

func (s *Storage) AutoMigrate() error {
	err := s.db.AutoMigrate(autoMigrateList...).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	err = s.db.Model(AuditPlanSQLV2{}).AddUniqueIndex("uniq_audit_plan_sqls_v2_audit_plan_id_fingerprint_md5",
		"audit_plan_id", "fingerprint_md5").Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	err = s.db.Model(BlackListAuditPlanSQL{}).AddUniqueIndex("uniq_type_content", "filter_type", "filter_content").Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	err = s.db.Model(&SqlManage{}).AddIndex("idx_project_id_status_deleted_at", "project_id", "status", "deleted_at").Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}

	if s.db.Dialect().HasColumn(Rule{}.TableName(), "is_default") {
		if err = s.db.Model(&Rule{}).DropColumn("is_default").Error; err != nil {
			return errors.New(errors.ConnectStorageError, err)
		}
	}

	return nil
}

func (s *Storage) CreateRulesIfNotExist(rules map[string][]*driverV2.Rule) error {
	isRuleExistInDB := func(rulesInDB []*Rule, targetRuleName, dbType string) (*Rule, bool) {
		for i := range rulesInDB {
			rule := rulesInDB[i]
			if rule.DBType != dbType || rule.Name != targetRuleName {
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
	for dbType, rules := range rules {
		for _, rule := range rules {
			existedRule, exist := isRuleExistInDB(rulesInDB, rule.Name, dbType)
			// rule will be created or update if:
			// 1. rule not exist;
			if !exist {
				err := s.Save(GenerateRuleByDriverRule(rule, dbType))
				if err != nil {
					return err
				}
			} else {
				isRuleDescSame := existedRule.Desc == rule.Desc
				isRuleAnnotationSame := existedRule.Annotation == rule.Annotation
				isRuleLevelSame := existedRule.Level == string(rule.Level)
				isRuleTypSame := existedRule.Typ == rule.Category
				existRuleParam, err := existedRule.Params.Value()
				if err != nil {
					return err
				}
				pluginRuleParam, err := rule.Params.Value()
				if err != nil {
					return err
				}
				isParamSame := reflect.DeepEqual(existRuleParam, pluginRuleParam)

				if !isRuleDescSame || !isRuleAnnotationSame || !isRuleLevelSame || !isRuleTypSame || !isParamSame {
					if existedRule.Knowledge != nil && existedRule.Knowledge.Content != "" {
						// 知识库是可以在页面上编辑的，而插件里只是默认内容，以页面上编辑后的内容为准
						rule.Knowledge.Content = existedRule.Knowledge.Content
					}
					err := s.Save(GenerateRuleByDriverRule(rule, dbType))
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
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
	for dbType, r := range rules {
		templateName := s.GetDefaultRuleTemplateName(dbType)
		exist, err := s.IsRuleTemplateExistFromAnyProject(projectId, templateName)
		if err != nil {
			return xerrors.Wrap(err, "get rule template failed")
		}
		if exist {
			continue
		}

		t := &RuleTemplate{
			ProjectId: projectId,
			Name:      templateName,
			Desc:      "默认规则模板",
			DBType:    dbType,
		}
		if err := s.Save(t); err != nil {
			return err
		}

		ruleList := make([]RuleTemplateRule, 0, len(r))
		for _, rule := range r {
			if rule.Level != driverV2.RuleLevelError {
				continue
			}
			modelRule := GenerateRuleByDriverRule(rule, dbType)
			ruleList = append(ruleList, RuleTemplateRule{
				RuleTemplateId: t.ID,
				RuleName:       modelRule.Name,
				RuleLevel:      modelRule.Level,
				RuleParams:     modelRule.Params,
				RuleDBType:     dbType,
			})
		}
		if err := s.UpdateRuleTemplateRules(t, ruleList...); err != nil {
			return xerrors.Wrap(err, "update rule template rules failed")
		}
	}

	return nil
}

func (s *Storage) GetDefaultRuleTemplateName(dbType string) string {
	return fmt.Sprintf("default_%v", dbType)
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
	var count int
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
	db := s.db.DB()
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

	sqlxDb := GetSqlxDb()

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
