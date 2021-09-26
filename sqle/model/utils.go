package model

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
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

func NewStorage(user, password, host, port, schema string, debug bool) (*Storage, error) {
	log.Logger().Infof("connecting to storage, host: %s, port: %s, user: %s, schema: %s",
		host, port, user, schema)
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
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

func (s *Storage) AutoMigrate() error {
	err := s.db.AutoMigrate(
		&RuleTemplateRule{},
		&Instance{},
		&RuleTemplate{},
		&Rule{},
		&Task{},
		&ExecuteSQL{},
		&RollbackSQL{},
		&SqlWhitelist{},
		&User{},
		&Role{},
		&WorkflowTemplate{},
		&WorkflowStepTemplate{},
		&Workflow{},
		&WorkflowRecord{},
		&WorkflowStep{},
		&SMTPConfiguration{},
		&SystemVariable{},
		&AuditPlan{},
		&AuditPlanReport{},
		&AuditPlanSQL{},
		&AuditPlanReportSQL{},
	).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	err = s.db.Model(&User{}).AddIndex("idx_users_id_name", "id", "login_name").Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	err = s.db.Model(&Instance{}).AddIndex("idx_instances_id_name", "id", "name").Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	return nil
}

func (s *Storage) CreateRulesIfNotExist(rules []*Rule) error {
	for _, rule := range rules {
		existedRule, exist, err := s.GetRule(rule.Name, rule.DBType)
		if err != nil {
			return err
		}
		if !exist || (existedRule.Value == "" && rule.Value != "") {
			err = s.Save(rule)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Storage) CreateDefaultTemplate(rules []*Rule) error {
	defaultTemplates := make(map[string][]*Rule)
	for _, rule := range rules {
		defaultTemplates[rule.DBType] = append(defaultTemplates[rule.DBType], rule)
	}

	for dbType, rs := range defaultTemplates {
		templateName := fmt.Sprintf("default_%v", dbType)
		_, exist, err := s.GetRuleTemplateByName(templateName)
		if err != nil {
			return err
		}
		if !exist {
			t := &RuleTemplate{
				Name:   templateName,
				Desc:   "默认规则模板",
				DBType: dbType,
			}
			if err := s.Save(t); err != nil {
				return err
			}

			ruleList := make([]RuleTemplateRule, 0, len(rs))
			for _, rule := range rs {
				if rule.IsDefault {
					ruleList = append(ruleList, RuleTemplateRule{
						RuleTemplateId: t.ID,
						RuleName:       rule.Name,
						RuleLevel:      rule.Level,
						RuleValue:      rule.Value,
						RuleDBType:     rule.DBType,
					})
				}
			}
			if err := s.UpdateRuleTemplateRules(t, ruleList...); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Storage) CreateAdminUser() error {
	_, exist, err := s.GetUserByName(DefaultAdminUser)
	if err != nil {
		return err
	}
	if !exist {
		return s.Save(&User{
			Name:     DefaultAdminUser,
			Password: "admin",
		})
	}
	return nil
}

var DefaultWorkflowTemplate = "default"

func (s *Storage) CreateDefaultWorkflowTemplate() error {
	user, exist, err := s.GetUserByName(DefaultAdminUser)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("admin user not exist")
	}
	_, exist, err = s.GetWorkflowTemplateByName(DefaultWorkflowTemplate)
	if err != nil {
		return err
	}
	if !exist {
		wt := &WorkflowTemplate{
			Name: DefaultWorkflowTemplate,
			Desc: "默认模板",
			Steps: []*WorkflowStepTemplate{
				{
					Number: 1,
					Typ:    WorkflowStepTypeSQLExecute,
					Users:  []*User{user},
				},
			},
		}
		err = s.SaveWorkflowTemplate(wt)
		if err != nil {
			return err
		}
	}
	return nil
}

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
		tx.Rollback()
		return errors.New(errors.ConnectStorageError, err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.New(errors.ConnectStorageError, err)
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
