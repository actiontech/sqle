package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

var storage *Storage

func InitStorage(s *Storage) {
	storage = s
}

func GetStorage() *Storage {
	return storage
}

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key" example:"1"`
	CreatedAt time.Time  `json:"create_at" example:"2018-10-21T16:40:23+08:00"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

func NewMysql(user, password, host, port, schema string) (*Storage, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, schema))
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.AutoMigrate(&Instance{}, &RuleTemplate{}, &Rule{}, &Task{}, &Sql{}, &Action{})

	//if err := createTable(db, &Instance{}); err != nil {
	//	return nil, err
	//}
	////if err := createTable(db, &Sql{}); err != nil {
	////	return nil, err
	////}
	////if err := createTable(db, &Task{}); err != nil {
	////	return nil, err
	////}
	//if err := createTable(db, &RuleTemplate{}); err != nil {
	//	return nil, err
	//}
	//if err := createTable(db, &Rule{}); err != nil {
	//	return nil, err
	//}
	return &Storage{
		db: db,
	}, nil
}

func createTable(db *gorm.DB, model interface{}) error {
	hasTable := db.HasTable(model)
	if db.Error != nil {
		return db.Error
	}
	if !hasTable {
		return db.CreateTable(model).Error
	}
	return nil
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) Exist(model interface{}) (bool, error) {
	var count int
	err := s.db.Model(model).Where(model).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Storage) Create(model interface{}) error {
	return s.db.Create(model).Error
}

func (s *Storage) Save(model interface{}) error {
	return s.db.Save(model).Error
}

func (s *Storage) Update(model interface{}, attrs ...interface{}) error {
	return s.db.Model(model).UpdateColumns(attrs).Error
}
