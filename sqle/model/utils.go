package model

import (
	"fmt"
	"sync"
	"time"

	"actiontech.cloud/universe/sqle/v3/sqle/errors"
	"actiontech.cloud/universe/sqle/v3/sqle/log"
	"github.com/jinzhu/gorm"
)

var storage *Storage

var storageMutex sync.Mutex

func InitStorage(s *Storage) {
	storageMutex.Lock()
	defer storageMutex.Unlock()
	storage = s
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

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key" example:"1"`
	CreatedAt time.Time  `json:"created_at" example:"2018-10-21T16:40:23+08:00"`
	UpdatedAt time.Time  `json:"updated_at" example:"2018-10-21T16:40:23+08:00"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

func NewStorage(user, password, host, port, schema string, debug bool) (*Storage, error) {
	log.Logger().Infof("connecting to storage, host: %s, port: %s, user: %s, schema: %s",
		host, port, user, schema)
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, schema))
	if err != nil {
		log.Logger().Errorf("connect to storage failed, error: %v", err)
		return nil, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	if debug {
		db.SetLogger(log.Logger().WithField("type", "sql"))
		db.LogMode(true)
	}
	log.Logger().Info("connected to storage")
	return &Storage{db: db}, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) AutoMigrate() error {
	err := s.db.AutoMigrate(&Instance{}, &RuleTemplate{}, &Rule{}, &Task{}, &CommitSql{}, &RollbackSql{}).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) CreateRulesIfNotExist(rules []Rule) error {
	for _, rule := range rules {
		exist, err := s.Exist(&rule)
		if err != nil {
			return err
		}
		if exist {
			continue
		}
		err = s.Save(rule)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) CreateDefaultTemplate(rules []Rule) error {
	_, exist, err := s.GetTemplateByName("all")
	if err != nil {
		return err
	}
	if !exist {
		t := &RuleTemplate{
			Name: "all",
			Desc: "default template for all rule",
		}
		if err := s.Save(t); err != nil {
			return err
		}
		return s.UpdateTemplateRules(t, rules...)
	}
	return nil
}

func (s *Storage) Exist(model interface{}) (bool, error) {
	var count int
	err := s.db.Model(model).Where(model).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	return count > 0, nil
}

func (s *Storage) Create(model interface{}) error {
	return errors.New(errors.CONNECT_STORAGE_ERROR, s.db.Create(model).Error)
}

func (s *Storage) Save(model interface{}) error {
	return errors.New(errors.CONNECT_STORAGE_ERROR, s.db.Save(model).Error)
}

func (s *Storage) Update(model interface{}, attrs ...interface{}) error {
	return errors.New(errors.CONNECT_STORAGE_ERROR, s.db.Model(model).UpdateColumns(attrs).Error)
}

func (s *Storage) Delete(model interface{}) error {
	return errors.New(errors.CONNECT_STORAGE_ERROR, s.db.Unscoped().Delete(model).Error)
}
